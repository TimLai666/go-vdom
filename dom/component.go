// component.go
package dom

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync/atomic"
)

var componentIDCounter uint64

func genComponentID() string {
	v := atomic.AddUint64(&componentIDCounter, 1)
	return fmt.Sprintf("vdom-%d", v)
}

type PropsDef map[string]interface{}

// Component 創建一個新的組件函數，支援預設 props
//   - template: 組件模板 VNode
//   - onDOMReadyCallback: (可選) 指向 JSAction 的指標（建議先用 jsdsl.Fn 建立一個 JSAction 變數，然後傳該變數的地址）
//     如果提供且非 nil，會在建立 node 後注入到 node.Props["onDOMReady"]
//   - defaultProps: 可選的預設 PropsDef
//
// 使用範例：
// act := jsdsl.Fn(nil, JSAction{Code: "/* ... */"})
// Component(template, &act, PropsDef{"id":"", ...})
// Component(template, nil, PropsDef{"id":"", ...}) // 傳 nil 表示不注入 onDOMReadyCallback
func Component(template VNode, onDOMReadyCallback *JSAction, defaultProps ...PropsDef) func(props Props, children ...VNode) VNode {
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
				interpolatedCode := interpolateString(onDOMReadyCallback.Code, mergedProps)
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
				// 純模板引用：取值並序列化為字符串（因為 HTML 屬性都是字符串）
				key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{{"), "}}"))
				if val, ok := p[key]; ok {
					// 所有類型都序列化為字符串，以保持與 HTML 屬性的一致性
					newProps[k] = serializeComplexType(val)
				} else {
					newProps[k] = "" // 找不到則為空字串
				}
			} else {
				// 混合字串或複雜模板：進行字串插值
				newProps[k] = interpolateString(t, p)
			}
		case JSAction:
			// 對 JSAction 的 Code 字串進行插值處理，保留為 JSAction
			newProps[k] = JSAction{Code: interpolateString(t.Code, p)}
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
// ${'{{label}}'.trim() ? 'inline' : 'none'}
// ${'{{direction}}' === 'horizontal' ? 'row' : 'column'}
// 流程：先替換 {{...}}，再解析並評估簡單的 ${...} 三元式
func interpolateString(s string, p Props) string {
	// 先替換 {{key}}
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)
	res := re.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
		if val, ok := p[key]; ok {
			// p[key] 可能是各種型別
			// 對於布林值，轉換為 "true" 或 "false"
			// 對於其他類型，統一以 fmt.Sprint 轉字串
			switch v := val.(type) {
			case bool:
				if v {
					return "true"
				}
				return "false"
			default:
				// 對於複雜類型，嘗試序列化為 JSON 字符串
				return serializeComplexType(val)
			}
		}
		return ""
	})

	// 處理 ${...} 表達式，支持嵌套三元運算符
	// 使用遞歸方式處理嵌套
	dollarBraceRe := regexp.MustCompile(`\$\{([^}]+)\}`)
	for dollarBraceRe.MatchString(res) {
		res = dollarBraceRe.ReplaceAllStringFunc(res, func(m string) string {
			// 提取 ${} 內的表達式
			expr := m[2 : len(m)-1] // 去掉 ${ 和 }
			return evaluateExpression(strings.TrimSpace(expr))
		})
	}

	return res
}

// evaluateExpression 評估簡單的條件表達式
// 支持格式：
// - 'value'.trim() ? 'A' : 'B'
// - 'X' === 'Y' ? 'A' : 'B'
// - 'X' !== 'Y' ? 'A' : 'B'
// - 嵌套三元: 'X' === 'Y' ? 'A' : 'Z' === 'W' ? 'B' : 'C'
func evaluateExpression(expr string) string {
	expr = strings.TrimSpace(expr)

	// 查找最外層的 ? 和 :
	questionIdx := -1
	colonIdx := -1
	depth := 0
	inQuotes := false

	for i := 0; i < len(expr); i++ {
		if expr[i] == '\'' {
			inQuotes = !inQuotes
		}
		if inQuotes {
			continue
		}

		if expr[i] == '?' {
			if questionIdx == -1 {
				questionIdx = i
			}
			depth++
		} else if expr[i] == ':' && depth > 0 {
			depth--
			if depth == 0 && colonIdx == -1 {
				colonIdx = i
			}
		}
	}

	// 如果不是三元表達式，直接返回去除引號的值
	if questionIdx == -1 || colonIdx == -1 {
		if strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'") {
			return expr[1 : len(expr)-1]
		}
		return expr
	}

	// 分解三元表達式: condition ? trueValue : falseValue
	condition := strings.TrimSpace(expr[:questionIdx])
	trueValue := strings.TrimSpace(expr[questionIdx+1 : colonIdx])
	falseValue := strings.TrimSpace(expr[colonIdx+1:])

	// 評估條件
	if evaluateCondition(condition) {
		return evaluateExpression(trueValue)
	}
	return evaluateExpression(falseValue)
}

// evaluateCondition 評估條件表達式
func evaluateCondition(condition string) bool {
	condition = strings.TrimSpace(condition)

	// 檢查 .trim() 語法: 'value'.trim()
	trimRe := regexp.MustCompile(`^'([^']*)'\.trim\(\)$`)
	if match := trimRe.FindStringSubmatch(condition); match != nil {
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

	// 默認：非空字符串為 true
	unquoted := unquote(condition)
	return unquoted != "" && unquoted != "false"
}

// indexOfOperator 找到運算符的位置（不在引號內）
func indexOfOperator(s, op string) int {
	inQuotes := false
	for i := 0; i <= len(s)-len(op); i++ {
		if s[i] == '\'' {
			inQuotes = !inQuotes
		}
		if !inQuotes && s[i:i+len(op)] == op {
			return i
		}
	}
	return -1
}

// unquote 移除字符串兩端的引號
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") && len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

// serializeComplexType 將複雜類型（slice、map、struct 等）序列化為 JSON 字符串
// 如果無法序列化，則返回 fmt.Sprint 的結果
func serializeComplexType(v interface{}) string {
	if v == nil {
		return ""
	}

	// 檢查是否為複雜類型（需要 JSON 序列化）
	rv := reflect.ValueOf(v)
	kind := rv.Kind()

	// 簡單類型直接返回字符串表示
	switch kind {
	case reflect.String:
		return rv.String()
	case reflect.Bool:
		if rv.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", rv.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", rv.Float())
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		// 複雜類型：序列化為 JSON
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			// 序列化失敗，返回默認字符串表示
			return fmt.Sprint(v)
		}
		return string(jsonBytes)
	case reflect.Ptr:
		// 指針類型：遞歸處理指向的值
		if rv.IsNil() {
			return ""
		}
		return serializeComplexType(rv.Elem().Interface())
	default:
		// 其他類型使用默認字符串表示
		return fmt.Sprint(v)
	}
}
