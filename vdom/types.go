// types.go
package vdom

// Props 是一個用於存儲元素屬性的映射
type Props map[string]string

// VNode 表示虛擬DOM中的一個節點
type VNode struct {
	Tag      string
	Props    Props
	Children []VNode
	Content  string
}

type JSAction struct {
	Code string
}
