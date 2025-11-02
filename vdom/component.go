// component.go
package vdom

import (
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

type PropsDef map[string]string

// Component 創建一個新的組件函數，支援預設 props
// - template: 組件模板 VNode
// - onDOMReady: (可選) JSAction（建議由 jsdsl.Fn 建立），如果提供，會在建立 node 後注入到 node.Props["onDOMReady"]
// - defaultProps: 可選的預設 PropsDef
//
// 使用範例：
// Component(template, jsdsl.Fn(nil, JSAction{Code: "/* ... */"}), PropsDef{"id":"", ...})
// Component(template, JSAction{}, PropsDef{"id":"", ...}) // 傳空 JSAction 表示不注入 onDOMReady
func Component(template VNode, onDOMReady JSAction, defaultProps ...PropsDef) func(props Props, children ...VNode) VNode {
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

		// 計算派生屬性（避免在模板中使用 JS-style 的 ${...} 表達式）
		// labelDisplay: 預設用於群組 label（block / none）
		if lbl, ok := mergedProps["label"]; ok && strings.TrimSpace(fmt.Sprint(lbl)) != "" {
			mergedProps["labelDisplay"] = "block"
		} else {
			mergedProps["labelDisplay"] = "none"
		}

		// helpDisplay: 當 helpText 有內容時為 block
		if ht, ok := mergedProps["helpText"]; ok && strings.TrimSpace(fmt.Sprint(ht)) != "" {
			mergedProps["helpDisplay"] = "block"
		} else {
			mergedProps["helpDisplay"] = "none"
		}

		// flexDirection: direction => row / column
		if dir, ok := mergedProps["direction"]; ok && strings.TrimSpace(fmt.Sprint(dir)) == "horizontal" {
			mergedProps["flexDirection"] = "row"
		} else {
			mergedProps["flexDirection"] = "column"
		}

		// 使用模板與合併後的 props 產生 VNode (先進行模板插值)
		node := interpolate(template, mergedProps, children)

		// 若提供了 onDOMReady JSAction 且其內容非空，且使用者未透過 props 顯式覆寫 onDOMReady，則將其注入為 node.Props["onDOMReady"]
		if strings.TrimSpace(onDOMReady.Code) != "" {
			if node.Props == nil {
				node.Props = make(Props)
			}
			// 只在使用者沒有明確提供 onDOMReady 的情況下注入（避免覆寫）
			if _, exists := node.Props["onDOMReady"]; !exists {
				node.Props["onDOMReady"] = onDOMReady
			}
		}

		return node
	}
}

// interpolate 替換模板中的變量
func interpolate(template VNode, p Props, children []VNode) VNode {
	newProps := make(Props)
	for k, v := range template.Props {
		// template.Props 的值可能為非 string 型別，先格式化為字串再進行插值
		// 若值為 JSAction，則對其 Code 內容做插值（替換 {{...}}），並保留為 JSAction 類型
		switch t := v.(type) {
		case string:
			newProps[k] = interpolateString(t, p)
		case JSAction:
			// 對 JSAction 的 Code 字串進行插值處理，保留為 JSAction
			newProps[k] = JSAction{Code: interpolateString(t.Code, p)}
		default:
			// 其他類型以字串形式處理
			valStr := fmt.Sprint(t)
			newProps[k] = interpolateString(valStr, p)
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
			// p[key] 可能是各種型別，統一以 fmt.Sprint 轉字串
			return fmt.Sprint(val)
		}
		return ""
	})

	// 處理形如 ${'somevalue'.trim() ? 'A' : 'B'} 的情況
	trimTernaryRe := regexp.MustCompile(`\$\{\s*'([^']*)'\.trim\(\)\s*\?\s*'([^']*)'\s*:\s*'([^']*)'\s*\}`)
	res = trimTernaryRe.ReplaceAllStringFunc(res, func(m string) string {
		sub := trimTernaryRe.FindStringSubmatch(m)
		if len(sub) >= 4 {
			val := strings.TrimSpace(sub[1])
			if val != "" {
				return sub[2]
			}
			return sub[3]
		}
		return m
	})

	// 處理形如 ${'X' === 'Y' ? 'A' : 'B'} 或 ${'X' == 'Y' ? 'A' : 'B'}
	eqTernaryRe := regexp.MustCompile(`\$\{\s*'([^']*)'\s*(?:===|==)\s*'([^']*)'\s*\?\s*'([^']*)'\s*:\s*'([^']*)'\s*\}`)
	res = eqTernaryRe.ReplaceAllStringFunc(res, func(m string) string {
		sub := eqTernaryRe.FindStringSubmatch(m)
		if len(sub) >= 5 {
			left := sub[1]
			right := sub[2]
			if left == right {
				return sub[3]
			}
			return sub[4]
		}
		return m
	})

	return res
}
