// component.go
package dom

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
)

var componentIDCounter uint64

func genComponentID() string {
	v := atomic.AddUint64(&componentIDCounter, 1)
	return fmt.Sprintf("vdom-%d", v)
}

type PropsDefault map[string]any

// Component 創建一個新的組件函數，支援預設 props
//   - template: 組件模板 VNode
//   - onDOMReadyCallback: (可選) 指向 JSAction 的指標（建議先用 jsdsl.Fn 建立一個 JSAction 變數，然後傳該變數的地址）
//     如果提供且非 nil，會在建立 node 後注入到 node.Props["onDOMReady"]
//   - defaultProps: 可選的預設 PropsDefault
//
// 使用範例：
// act := jsdsl.Fn(nil, JSAction{Code: "/* ... */"})
// Component(template, &act, PropsDefault{"id":"", ...})
// Component(template, nil, PropsDefault{"id":"", ...}) // 傳 nil 表示不注入 onDOMReadyCallback
func Component(template VNode, onDOMReadyCallback *JSAction, defaultProps ...PropsDefault) func(props Props, children ...VNode) VNode {
	return func(p Props, children ...VNode) VNode {
		mergedProps := make(Props)

		// 若有提供 defaultProps，合併進 mergedProps
		if len(defaultProps) > 0 && defaultProps[0] != nil {
			for k, v := range defaultProps[0] {
				mergedProps[k] = v
			}
		}

		// 合併使用者傳入的 props（Props 已為 map[string]interface{}）
		for k, v := range p {
			mergedProps[k] = v
		}

		// 如果沒有提供 id，生成一個穩定的唯一 id（便於組件內腳本綁定 container）
		if idv, ok := mergedProps["id"]; !ok || strings.TrimSpace(fmt.Sprint(idv)) == "" {
			mergedProps["id"] = genComponentID()
		}

		// 使用模板與合併後的 props 產生 VNode (先進行模板插值)
		node := interpolate(template, mergedProps, children)

		// 若提供了 onDOMReadyCallback（指標）且其內容非空，且使用者未透過 props 顯式覆寫 onDOMReadyCallback，則將其注入為 node.Props["onDOMReadyCallback"]
		if onDOMReadyCallback != nil && strings.TrimSpace(onDOMReadyCallback.Code) != "" {
			if node.Props == nil {
				node.Props = make(Props)
			}
			// 只在使用者沒有明確提供 onDOMReady 的情況下注入（避免覆寫）
			if _, exists := node.Props["onDOMReady"]; !exists {
				// 將 JSAction 值放入 Props，並對其中的 {{...}} 模板進行插值
				// 使用 interpolateStringForJS 保持 JSON 格式
				interpolatedCode := interpolateStringForJS(onDOMReadyCallback.Code, mergedProps)
				node.Props["onDOMReady"] = JSAction{Code: interpolatedCode}
			}
		}

		return node
	}
}

// interpolate 替換模板中的變量
func interpolate(template VNode, p Props, children []VNode) VNode {
	newProps := make(Props)

	// 先處理模板中定義的 Props
	for k, v := range template.Props {
		// template.Props 的值可能為非 string 型別
		// 若值為 JSAction，則對其 Code 內容做插值（替換 {{...}}），並保留為 JSAction 類型
		switch t := v.(type) {
		case string:
			// 檢查是否為純模板引用（如 "{{key}}"）
			trimmed := strings.TrimSpace(t)
			if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") && strings.Count(trimmed, "{{") == 1 {
				// 純模板引用：取值並轉換為字符串
				key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{{"), "}}"))
				if val, ok := p[key]; ok {
					// 用於 HTML 屬性，需要去除 JSON 字符串的引號
					jsonStr := serializeComplexType(val)
					// 如果是 JSON 字符串格式（以 " 開頭和結尾），去除引號
					if len(jsonStr) >= 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
						// 去除引號並反轉義
						unquoted := jsonStr[1 : len(jsonStr)-1]
						unquoted = strings.ReplaceAll(unquoted, `\"`, `"`)
						unquoted = strings.ReplaceAll(unquoted, `\\`, `\`)
						newProps[k] = unquoted
					} else {
						newProps[k] = jsonStr
					}
				} else {
					newProps[k] = "" // 找不到則為空字串
				}
			} else {
				// 混合字串或複雜模板：進行字串插值
				newProps[k] = interpolateString(t, p)
			}
		case JSAction:
			// 對 JSAction 的 Code 字串進行插值處理，保留為 JSAction
			// 在 JavaScript 代碼中，保持 JSON 格式（字符串帶引號）
			newProps[k] = JSAction{Code: interpolateStringForJS(t.Code, p)}
		case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			// 保留原始類型的數字和布林值
			newProps[k] = v
		default:
			// 其他複雜類型（如 slice, map, struct 等）需要序列化為 JSON
			// 以便在客戶端 JavaScript 中使用
			newProps[k] = serializeComplexType(v)
		}
	}

	newChildren := []VNode{}
	for _, c := range template.Children {
		if c.Tag == "" {
			// 純文字節點要替換 {{children}}
			content := strings.TrimSpace(c.Content)
			if content == "{{children}}" {
				newChildren = append(newChildren, children...)
			} else {
				newChildren = append(newChildren, VNode{
					Content: interpolateString(c.Content, p),
				})
			}
		} else {
			newChildren = append(newChildren, interpolate(c, p, children))
		}
	}

	return VNode{
		Tag:      template.Tag,
		Props:    newProps,
		Children: newChildren,
		Content:  interpolateString(template.Content, p),
	}
}

// interpolateString 替換字符串中的變量
// 支援額外處理少量常見的 JS-style ternary 表達式，例如：
// ${{{label}}.trim() ? 'inline' : 'none'}
// ${{{direction}} === 'horizontal' ? 'row' : 'column'}
// 流程：先處理 ${...} 表達式（其中可能包含 {{...}}），再替換剩餘的 {{...}}
func interpolateString(s string, p Props) string {
	// 先處理 ${...} 表達式，支持嵌套三元運算符
	// 在表達式內部的 {{...}} 會在 evaluateExpression 中處理
	result := s
	// 手動查找並處理 ${...} 表達式
	for {
		startIdx := strings.Index(result, "${")
		if startIdx == -1 {
			break
		}

		// 從 ${ 之後開始，找到匹配的 }
		depth := 1
		endIdx := -1
		for i := startIdx + 2; i < len(result); i++ {
			if result[i] == '{' {
				depth++
			} else if result[i] == '}' {
				depth--
				if depth == 0 {
					endIdx = i
					break
				}
			}
		}

		if endIdx == -1 {
			// 沒有找到匹配的 }，跳過
			break
		}

		// 提取 ${} 內的表達式
		expr := result[startIdx+2 : endIdx]
		// 在表達式中替換 {{...}}
		exprWithValues := replaceTemplateVarsInExpression(expr, p)
		evaluated := evaluateExpression(strings.TrimSpace(exprWithValues))

		// 替換整個 ${...} 為評估結果
		result = result[:startIdx] + evaluated + result[endIdx+1:]
	}

	// 再替換剩餘的 {{key}}（不在 ${...} 內的）
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
		if val, ok := p[key]; ok {
			// 在字符串插值中（用於文本內容），需要去除 JSON 字符串的引號
			jsonStr := serializeComplexType(val)
			// 如果是 JSON 字符串格式（以 " 開頭和結尾），去除引號
			if len(jsonStr) >= 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
				// 去除引號並反轉義
				unquoted := jsonStr[1 : len(jsonStr)-1]
				unquoted = strings.ReplaceAll(unquoted, `\"`, `"`)
				unquoted = strings.ReplaceAll(unquoted, `\\`, `\`)
				return unquoted
			}
			return jsonStr
		}
		return ""
	})

	return result
}

// replaceTemplateVarsInExpression 在表達式中替換 {{...}} 為 JSON 值
func replaceTemplateVarsInExpression(expr string, p Props) string {
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)
	return re.ReplaceAllStringFunc(expr, func(match string) string {
		key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
		if val, ok := p[key]; ok {
			// 返回 JSON 格式（字符串帶引號）
			return serializeComplexType(val)
		}
		return "null"
	})
}

// evaluateExpression 評估簡單的條件表達式
// 支持格式：
// - 'value'.trim() ? 'A' : 'B'
// - 'X' === 'Y' ? 'A' : 'B'
// - 'X' !== 'Y' ? 'A' : 'B'
// - 嵌套三元: 'X' === 'Y' ? 'A' : 'Z' === 'W' ? 'B' : 'C'
// - 括號分組: 'X' ? ('Y' ? 'A' : 'B') : 'C'
func evaluateExpression(expr string) string {
	expr = strings.TrimSpace(expr)

	// 查找最外層的 ? 和 :
	questionIdx := -1
	colonIdx := -1
	depth := 0
	parenDepth := 0
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		// 處理引號
		if (ch == '\'' || ch == '"') && (i == 0 || expr[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
		}

		if inQuotes {
			continue
		}

		// 追蹤括號深度
		if ch == '(' {
			parenDepth++
		} else if ch == ')' {
			parenDepth--
		}

		// 只有在括號深度為 0 時才處理 ? 和 :
		if parenDepth == 0 {
			if ch == '?' {
				if questionIdx == -1 {
					questionIdx = i
				}
				depth++
			} else if ch == ':' && depth > 0 {
				depth--
				if depth == 0 && colonIdx == -1 {
					colonIdx = i
				}
			}
		}
	}

	// 如果不是三元表達式，直接返回去除引號的值
	if questionIdx == -1 || colonIdx == -1 {
		if strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'") && len(expr) >= 2 {
			return expr[1 : len(expr)-1]
		}
		if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") && len(expr) >= 2 {
			return expr[1 : len(expr)-1]
		}
		return expr
	}

	// 分解三元表達式: condition ? trueValue : falseValue
	condition := strings.TrimSpace(expr[:questionIdx])
	trueValue := strings.TrimSpace(expr[questionIdx+1 : colonIdx])
	falseValue := strings.TrimSpace(expr[colonIdx+1:])

	// 移除 trueValue 和 falseValue 的外層括號
	trueValue = stripOuterParentheses(trueValue)
	falseValue = stripOuterParentheses(falseValue)

	// 評估條件
	if evaluateCondition(condition) {
		return evaluateExpression(trueValue)
	}
	return evaluateExpression(falseValue)
}

// evaluateCondition 評估條件表達式
func evaluateCondition(condition string) bool {
	condition = strings.TrimSpace(condition)

	// 移除最外層的括號（如果存在）
	condition = stripOuterParentheses(condition)

	// 檢查 && 運算符
	if idx := indexOfOperator(condition, "&&"); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+2:])
		return evaluateCondition(left) && evaluateCondition(right)
	}

	// 檢查 || 運算符
	if idx := indexOfOperator(condition, "||"); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+2:])
		return evaluateCondition(left) || evaluateCondition(right)
	}

	// 檢查 .trim() 語法: 'value'.trim() 或 "value".trim()
	trimReSingle := regexp.MustCompile(`^'([^']*)'\.trim\(\)$`)
	if match := trimReSingle.FindStringSubmatch(condition); match != nil {
		return strings.TrimSpace(match[1]) != ""
	}
	trimReDouble := regexp.MustCompile(`^"([^"]*)"\.trim\(\)$`)
	if match := trimReDouble.FindStringSubmatch(condition); match != nil {
		return strings.TrimSpace(match[1]) != ""
	}

	// 檢查 === 或 == 比較
	if idx := indexOfOperator(condition, "==="); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+3:])
		return unquote(left) == unquote(right)
	}

	if idx := indexOfOperator(condition, "=="); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+2:])
		return unquote(left) == unquote(right)
	}

	// 檢查 !== 或 != 比較
	if idx := indexOfOperator(condition, "!=="); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+3:])
		return unquote(left) != unquote(right)
	}

	if idx := indexOfOperator(condition, "!="); idx != -1 {
		left := strings.TrimSpace(condition[:idx])
		right := strings.TrimSpace(condition[idx+2:])
		return unquote(left) != unquote(right)
	}

	// 布林值字面量
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}

	// 默認：非空字符串為 true
	unquoted := unquote(condition)
	return unquoted != "" && unquoted != "false" && unquoted != "null"
}

// stripOuterParentheses 移除表達式最外層的括號
func stripOuterParentheses(expr string) string {
	expr = strings.TrimSpace(expr)
	if len(expr) < 2 || expr[0] != '(' {
		return expr
	}

	// 檢查是否整個表達式被一對括號包圍
	depth := 0
	for i := 0; i < len(expr); i++ {
		if expr[i] == '(' {
			depth++
		} else if expr[i] == ')' {
			depth--
			// 如果在結尾之前深度就變成 0，說明不是被一對括號完全包圍
			if depth == 0 && i < len(expr)-1 {
				return expr
			}
		}
	}

	// 如果整個表達式被一對括號包圍，移除它們並遞歸檢查
	if depth == 0 {
		return stripOuterParentheses(expr[1 : len(expr)-1])
	}

	return expr
}

// indexOfOperator 找到運算符的位置（不在引號內且不在括號內）
func indexOfOperator(s, op string) int {
	inQuotes := false
	quoteChar := byte(0)
	parenDepth := 0

	for i := 0; i <= len(s)-len(op); i++ {
		ch := s[i]

		// 處理引號
		if (ch == '\'' || ch == '"') && (i == 0 || s[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
		}

		if inQuotes {
			continue
		}

		// 追蹤括號深度
		if ch == '(' {
			parenDepth++
		} else if ch == ')' {
			parenDepth--
		}

		// 只有在引號外且括號深度為 0 時才匹配運算符
		if !inQuotes && parenDepth == 0 && s[i:i+len(op)] == op {
			return i
		}
	}
	return -1
}

// unquote 移除字符串兩端的引號（支持單引號和雙引號）
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") && len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") && len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

// interpolateStringForJS 替換 JavaScript 代碼中的變量，保持 JSON 格式
func interpolateStringForJS(s string, p Props) string {
	// 先處理 ${...} 表達式
	result := s
	// 手動查找並處理 ${...} 表達式
	for {
		startIdx := strings.Index(result, "${")
		if startIdx == -1 {
			break
		}

		// 從 ${ 之後開始，找到匹配的 }
		depth := 1
		endIdx := -1
		for i := startIdx + 2; i < len(result); i++ {
			if result[i] == '{' {
				depth++
			} else if result[i] == '}' {
				depth--
				if depth == 0 {
					endIdx = i
					break
				}
			}
		}

		if endIdx == -1 {
			// 沒有找到匹配的 }，跳過
			break
		}

		// 提取 ${} 內的表達式
		expr := result[startIdx+2 : endIdx]
		// 在表達式中替換 {{...}}
		exprWithValues := replaceTemplateVarsInExpression(expr, p)
		evaluated := evaluateExpression(strings.TrimSpace(exprWithValues))

		// 替換整個 ${...} 為評估結果
		result = result[:startIdx] + evaluated + result[endIdx+1:]
	}

	// 再替換剩餘的 {{key}}（不在 ${...} 內的）
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
		if val, ok := p[key]; ok {
			// 在 JavaScript 代碼中，保持 JSON 格式（字符串帶引號）
			return serializeComplexType(val)
		}
		return "null"
	})

	return result
}

// serializeComplexType 將所有值統一序列化為 JSON 格式
// 這樣在 JavaScript 中可以直接使用：const value = {{prop}};
func serializeComplexType(v interface{}) string {
	if v == nil {
		return "null"
	}

	// 統一使用 JSON 序列化
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		// 序列化失敗，返回默認字符串表示
		return fmt.Sprint(v)
	}
	return string(jsonBytes)
}
