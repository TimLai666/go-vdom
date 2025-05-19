// render.go
package vdom

import (
	"fmt"
	"strings"
)

// Render 將虛擬DOM節點轉換為HTML字符串
func Render(v VNode) string {
	if v.Tag == "" {
		return v.Content
	}

	var sb strings.Builder
	sb.WriteString("<" + v.Tag)
	for k, val := range v.Props {
		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, val))
	}
	sb.WriteString(">")

	if v.Content != "" {
		sb.WriteString(v.Content)
	}
	for _, c := range v.Children {
		sb.WriteString(Render(c))
	}

	sb.WriteString(fmt.Sprintf("</%s>", v.Tag))
	return sb.String()
}
