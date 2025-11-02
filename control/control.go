// control.go
package control

import (
	"github.com/TimLai666/go-vdom/dom"
)

// ThenBlock 表示If條件為真時要渲染的內容
type ThenBlock struct {
	Content []dom.VNode
}

// ElseBlock 表示If條件為假時要渲染的內容
type ElseBlock struct {
	Content []dom.VNode
}

// ElseIfBlock 表示If條件為假且符合ElseIf條件時要渲染的內容
type ElseIfBlock struct {
	Condition bool
	Content   []dom.VNode
}

// Then 創建一個ThenBlock
func Then(nodes ...dom.VNode) ThenBlock {
	return ThenBlock{Content: nodes}
}

// Else 創建一個ElseBlock
func Else(nodes ...dom.VNode) ElseBlock {
	return ElseBlock{Content: nodes}
}

// ElseIf 創建一個 ElseIfBlock
func ElseIf(condition bool, nodes ...dom.VNode) ElseIfBlock {
	return ElseIfBlock{Condition: condition, Content: nodes}
}

// If 條件渲染
// 支援 If(condition, Then(...), ElseIf(cond, ...), ..., Else(...))
func If(condition bool, thenBlock ThenBlock, elseIfOrElse ...any) []dom.VNode {
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
	return []dom.VNode{}
}

// Repeat 重複渲染
// 重複指定次數的渲染
func Repeat(count int, renderFunc func(index int) dom.VNode) []dom.VNode {
	if count <= 0 {
		return []dom.VNode{}
	}

	result := make([]dom.VNode, count)
	for i := 0; i < count; i++ {
		result[i] = renderFunc(i)
	}
	return result
}

// ForEach 遍歷渲染
// 對數據集合中的每一項應用渲染函數
func ForEach[T any](items []T, renderFunc func(item T, index int) dom.VNode) []dom.VNode {
	result := make([]dom.VNode, len(items))
	for i, item := range items {
		result[i] = renderFunc(item, i)
	}
	return result
}

// For 傳統循環渲染
// 類似於 for i := start; 條件; i += step 的循環
// 參數：
// - start: 起始值
// - end: 結束值（不包含）
// - step: 步進值（可以是負數用於倒序）
// 用法：
//
//	For(0, 10, 1, func(i int) VNode { return Div(fmt.Sprintf("項目 %d", i)) })  // 0-9
//	For(10, 0, -1, func(i int) VNode { return Div(fmt.Sprintf("項目 %d", i)) }) // 10-1 倒序
//	For(0, 100, 10, func(i int) VNode { return Div(fmt.Sprintf("項目 %d", i)) }) // 0, 10, 20, ..., 90
func For(start, end, step int, renderFunc func(i int) dom.VNode) []dom.VNode {
	if step == 0 {
		return []dom.VNode{}
	}

	var result []dom.VNode

	if step > 0 {
		// 正向循環
		for i := start; i < end; i += step {
			result = append(result, renderFunc(i))
		}
	} else {
		// 反向循環
		for i := start; i > end; i += step {
			result = append(result, renderFunc(i))
		}
	}

	return result
}

// Map 映射渲染
// 與 For 類似，但更強調數據轉換
func Map[T any, U any](items []T, mapFunc func(item T, index int) U, renderFunc func(mappedItem U, index int) dom.VNode) []dom.VNode {
	result := make([]dom.VNode, len(items))
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
	Content   []dom.VNode
}

func Switch(cases []Case, defaultContent []dom.VNode) []dom.VNode {
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
	Node dom.VNode
}

// KeyedForEach 帶鍵值的循環渲染
// 為列表中的每個項目創建帶有唯一鍵的節點
func KeyedForEach[T any](items []T, keyFunc func(item T, index int) string, renderFunc func(item T, index int) dom.VNode) []KeyedNode {
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
func ToNodes(keyedNodes []KeyedNode) []dom.VNode {
	result := make([]dom.VNode, len(keyedNodes))
	for i, kn := range keyedNodes {
		result[i] = kn.Node
	}
	return result
}
