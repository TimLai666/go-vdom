package jsdsl

import . "github.com/TimLai666/go-vdom/vdom"

// JSActionBuilder 用於收集和處理 JSAction
type JSActionBuilder struct {
	actions []JSAction
}

// NewJSActionBuilder 創建一個新的 JSActionBuilder
func NewJSActionBuilder() *JSActionBuilder {
	return &JSActionBuilder{
		actions: []JSAction{},
	}
}

// Add 添加一個 JSAction 到 builder
func (b *JSActionBuilder) Add(action JSAction) *JSActionBuilder {
	b.actions = append(b.actions, action)
	return b
}

// AddMany 添加多個 JSAction 到 builder
func (b *JSActionBuilder) AddMany(actions ...JSAction) *JSActionBuilder {
	b.actions = append(b.actions, actions...)
	return b
}

// CreateElement 創建一個 DOM 元素並添加到 builder
// 返回創建的元素，可以繼續操作
func (b *JSActionBuilder) CreateElement(tagName string, varName ...string) Elem {
	elem, action := CreateEl(tagName, varName...)
	b.Add(action)
	return elem
}

// SetElementText 設置元素文本並添加到 builder
func (b *JSActionBuilder) SetElementText(elem Elem, text string) *JSActionBuilder {
	b.Add(elem.SetText(text))
	return b
}

// SetElementHTML 設置元素 HTML 並添加到 builder
func (b *JSActionBuilder) SetElementHTML(elem Elem, html string) *JSActionBuilder {
	b.Add(elem.SetHTML(html))
	return b
}

// AddElementClass 為元素添加類別並添加到 builder
func (b *JSActionBuilder) AddElementClass(elem Elem, class string) *JSActionBuilder {
	b.Add(elem.AddClass(class))
	return b
}

// RemoveElementClass 移除元素類別並添加到 builder
func (b *JSActionBuilder) RemoveElementClass(elem Elem, class string) *JSActionBuilder {
	b.Add(elem.RemoveClass(class))
	return b
}

// AppendChild 將子元素添加到父元素並添加到 builder
func (b *JSActionBuilder) AppendChild(parent Elem, child Elem) *JSActionBuilder {
	b.Add(parent.AppendChild(child))
	return b
}

// GetActions 獲取所有收集的 JSAction
func (b *JSActionBuilder) GetActions() []JSAction {
	return b.actions
}

// Build 生成最終的 DomReady JSAction
func (b *JSActionBuilder) Build() JSAction {
	return DomReady(b.actions...)
}
