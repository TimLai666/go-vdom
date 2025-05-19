// control.go
package control

import (
	"github.com/TimLai666/go-vdom/vdom"
)

// ThenBlock 表示If條件為真時要渲染的內容
type ThenBlock struct {
	Content []vdom.VNode
}

// ElseBlock 表示If條件為假時要渲染的內容
type ElseBlock struct {
	Content []vdom.VNode
}

// ElseIfBlock 表示If條件為假且符合ElseIf條件時要渲染的內容
type ElseIfBlock struct {
	Condition bool
	Content   []vdom.VNode
}

// Then 創建一個ThenBlock
func Then(nodes ...vdom.VNode) ThenBlock {
	return ThenBlock{Content: nodes}
}

// Else 創建一個ElseBlock
func Else(nodes ...vdom.VNode) ElseBlock {
	return ElseBlock{Content: nodes}
}

// ElseIf 創建一個 ElseIfBlock
func ElseIf(condition bool, nodes ...vdom.VNode) ElseIfBlock {
	return ElseIfBlock{Condition: condition, Content: nodes}
}

// If 條件渲染
// 支援 If(condition, Then(...), ElseIf(cond, ...), ..., Else(...))
func If(condition bool, thenBlock ThenBlock, elseIfOrElse ...any) []vdom.VNode {
	if condition {
		return thenBlock.Content
	}
	for _, block := range elseIfOrElse {
		switch b := block.(type) {
		case ElseIfBlock:
			if b.Condition {
				return b.Content
			}
		case ElseBlock:
			return b.Content
		}
	}
	return []vdom.VNode{}
}

// Repeat 重複渲染
// 重複指定次數的渲染
func Repeat(count int, renderFunc func(index int) vdom.VNode) []vdom.VNode {
	if count <= 0 {
		return []vdom.VNode{}
	}

	result := make([]vdom.VNode, count)
	for i := 0; i < count; i++ {
		result[i] = renderFunc(i)
	}
	return result
}

// For 循環渲染
// 對數據集合中的每一項應用渲染函數
func For[T any](items []T, renderFunc func(item T, index int) vdom.VNode) []vdom.VNode {
	result := make([]vdom.VNode, len(items))
	for i, item := range items {
		result[i] = renderFunc(item, i)
	}
	return result
}

// Map 映射渲染
// 與 For 類似，但更強調數據轉換
func Map[T any, U any](items []T, mapFunc func(item T, index int) U, renderFunc func(mappedItem U, index int) vdom.VNode) []vdom.VNode {
	result := make([]vdom.VNode, len(items))
	for i, item := range items {
		mappedItem := mapFunc(item, i)
		result[i] = renderFunc(mappedItem, i)
	}
	return result
}

// Switch-Case 條件渲染
// 類似於 switch-case 語句的條件渲染
type Case struct {
	Condition bool
	Content   []vdom.VNode
}

func Switch(cases []Case, defaultContent []vdom.VNode) []vdom.VNode {
	for _, c := range cases {
		if c.Condition {
			return c.Content
		}
	}
	return defaultContent
}

// 帶有鍵值的節點，用於提高列表渲染性能
type KeyedNode struct {
	Key  string
	Node vdom.VNode
}

// KeyedFor 帶鍵值的循環渲染
// 為列表中的每個項目創建帶有唯一鍵的節點
func KeyedFor[T any](items []T, keyFunc func(item T, index int) string, renderFunc func(item T, index int) vdom.VNode) []KeyedNode {
	result := make([]KeyedNode, len(items))
	for i, item := range items {
		key := keyFunc(item, i)
		result[i] = KeyedNode{
			Key:  key,
			Node: renderFunc(item, i),
		}
	}
	return result
}

// ToNodes 將 KeyedNode 轉換為普通的 VNode 數組
func ToNodes(keyedNodes []KeyedNode) []vdom.VNode {
	result := make([]vdom.VNode, len(keyedNodes))
	for i, kn := range keyedNodes {
		result[i] = kn.Node
	}
	return result
}
