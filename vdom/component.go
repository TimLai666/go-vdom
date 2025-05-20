// component.go
package vdom

import (
	"regexp"
	"strings"
)

// Component 創建一個新的組件函數，支援預設 props
func Component(template VNode, defaultProps ...Props) func(props Props, children ...VNode) VNode {
	return func(p Props, children ...VNode) VNode {
		mergedProps := make(Props)
		if len(defaultProps) > 0 {
			for k, v := range defaultProps[0] {
				mergedProps[k] = v
			}
		}
		for k, v := range p {
			mergedProps[k] = v
		}
		return interpolate(template, mergedProps, children)
	}
}

// interpolate 替換模板中的變量
func interpolate(template VNode, p Props, children []VNode) VNode {
	newProps := make(Props)
	for k, v := range template.Props {
		newProps[k] = interpolateString(v, p)
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
func interpolateString(s string, p Props) string {
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}"))
		if val, ok := p[key]; ok {
			return val
		}
		return match
	})
}
