// tags.go
package dom

import (
	"strings"

	"github.com/TimLai666/go-vdom/runtime"
)

// 基本標籤函數

// Text 創建一個文字節點
func Text(s string) VNode {
	return VNode{Content: s}
}

// ForEach 提供簡潔的列表渲染語法（後端渲染）
// 用法：Ul(ForEach(items, func(item string) VNode { return Li(item) }))
func ForEach[T any](items []T, renderFunc func(item T) VNode) []VNode {
	result := make([]VNode, len(items))
	for i, item := range items {
		result[i] = renderFunc(item)
	}
	return result
}

// ForEachWithIndex 提供帶索引的列表渲染（後端渲染）
// 用法：Ul(ForEachWithIndex(items, func(item string, i int) VNode { return Li(fmt.Sprintf("%d. %s", i+1, item)) }))
func ForEachWithIndex[T any](items []T, renderFunc func(item T, index int) VNode) []VNode {
	result := make([]VNode, len(items))
	for i, item := range items {
		result[i] = renderFunc(item, i)
	}
	return result
}

// 常用 HTML 標籤函數

// tag 是所有 HTML 標籤函數的基礎實現
func tag(name string, props Props, children ...any) VNode {
	// 如果沒有傳入 props，初始化為空 map 以便合併
	if props == nil {
		props = make(Props)
	}

	var chs []VNode
	var content string

	for _, child := range children {
		switch v := child.(type) {
		case Props:
			if v != nil {
				if props == nil {
					props = make(Props)
				}
				for k, val := range v {
					props[k] = val
				}
			}
		case VNode:
			chs = append(chs, v)
		case []VNode: // 自動展開
			chs = append(chs, v...)
		case string:
			chs = append(chs, Text(v))
		}
	}

	return VNode{
		Tag:      name,
		Props:    props,
		Children: chs,
		Content:  content,
	}
}

// Title 創建 <title> 標籤
func Title(s string) VNode { return tag("title", nil, Text(s)) }

// Meta 創建 <meta> 標籤（接受 Props）
func Meta(p Props) VNode { return tag("meta", p) }

// Link 創建 <link> 標籤（接受 Props）
func Link(p Props) VNode { return tag("link", p) }

// LinkInfo 定義外部資源鏈接的信息
type LinkInfo struct {
	Rel  string // 關聯類型，如 "stylesheet", "icon", "preload" 等
	Href string // 鏈接地址
	Type string // 媒體類型，如 "text/css"
}

// ScriptInfo 定義JavaScript腳本的信息
type ScriptInfo struct {
	Src   string // 腳本地址
	Async bool   // 是否異步加載
}

// Document 創建一個完整的HTML文檔結構
// 參數:
// - title: 頁面標題
// - links: 要加入的外部資源鏈接列表
// - scripts: 要加入的JavaScript腳本列表
// - metas: 要加入的meta標籤屬性列表
// - children: 頁面主體內容
func Document(title string, links []LinkInfo, scripts []ScriptInfo, metas []Props, children ...VNode) VNode {
	// 頭部元素列表
	var headElements []any

	// 添加標題
	headElements = append(headElements, Title(title))

	// 添加基本meta標籤
	headElements = append(headElements, Meta(Props{"charset": "UTF-8"}))
	headElements = append(headElements, Meta(Props{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}))

	// 添加自定義meta標籤
	for _, meta := range metas {
		headElements = append(headElements, Meta(meta))
	}

	// 添加外部資源鏈接
	for _, link := range links {
		props := Props{"rel": link.Rel, "href": link.Href}
		if link.Type != "" {
			props["type"] = link.Type
		}
		headElements = append(headElements, Link(props))
	}

	// 自動注入 go-vdom runtime 腳本（必須在其他腳本之前加載）
	headElements = append(headElements, Script(Props{}, runtime.ClientRuntime()))

	// 添加腳本
	for _, script := range scripts {
		props := Props{"src": script.Src}
		if script.Async {
			props["async"] = "true"
		}
		headElements = append(headElements, Script(props))
	}

	// 將children轉換為any類型
	var bodyElements []any
	for _, node := range children {
		bodyElements = append(bodyElements, node)
	}

	return Html(nil,
		Head(nil, headElements...),
		Body(nil, bodyElements...),
	)
}

// HTML 結構元素
func Html(p Props, children ...any) VNode { return tag("html", p, children...) }

func Head(p Props, children ...any) VNode { return tag("head", p, children...) }

func Body(p Props, children ...any) VNode { return tag("body", p, children...) }

func Main(p Props, children ...any) VNode { return tag("main", p, children...) }

func Header(p Props, children ...any) VNode { return tag("header", p, children...) }

func Footer(p Props, children ...any) VNode { return tag("footer", p, children...) }

func Nav(p Props, children ...any) VNode { return tag("nav", p, children...) }

func Aside(p Props, children ...any) VNode { return tag("aside", p, children...) }

func Section(p Props, children ...any) VNode { return tag("section", p, children...) }

func Article(p Props, children ...any) VNode { return tag("article", p, children...) }

func Address(p Props, children ...any) VNode { return tag("address", p, children...) }

func Hgroup(p Props, children ...any) VNode { return tag("hgroup", p, children...) }

// 文本區塊元素
func H1(p Props, children ...any) VNode { return tag("h1", p, children...) }

func H2(p Props, children ...any) VNode { return tag("h2", p, children...) }

func H3(p Props, children ...any) VNode { return tag("h3", p, children...) }

func H4(p Props, children ...any) VNode { return tag("h4", p, children...) }

func H5(p Props, children ...any) VNode { return tag("h5", p, children...) }

func H6(p Props, children ...any) VNode { return tag("h6", p, children...) }

func P(p Props, children ...any) VNode { return tag("p", p, children...) }

func Div(p Props, children ...any) VNode { return tag("div", p, children...) }

func Span(p Props, children ...any) VNode { return tag("span", p, children...) }

func Pre(p Props, children ...any) VNode { return tag("pre", p, children...) }

func Code(p Props, children ...any) VNode { return tag("code", p, children...) }

func Blockquote(p Props, children ...any) VNode { return tag("blockquote", p, children...) }

// 表單元素
func Form(p Props, children ...any) VNode { return tag("form", p, children...) }

func Input(p Props, children ...any) VNode { return tag("input", p, children...) }

func Label(p Props, children ...any) VNode { return tag("label", p, children...) }

func Button(p Props, children ...any) VNode { return tag("button", p, children...) }

func Select(p Props, children ...any) VNode { return tag("select", p, children...) }

func Datalist(p Props, children ...any) VNode { return tag("datalist", p, children...) }

func Optgroup(p Props, children ...any) VNode { return tag("optgroup", p, children...) }

func Option(p Props, children ...any) VNode { return tag("option", p, children...) }

func Textarea(p Props, children ...any) VNode { return tag("textarea", p, children...) }

func Output(p Props, children ...any) VNode { return tag("output", p, children...) }

func Progress(p Props, children ...any) VNode { return tag("progress", p, children...) }

func Meter(p Props, children ...any) VNode { return tag("meter", p, children...) }

func Fieldset(p Props, children ...any) VNode { return tag("fieldset", p, children...) }

func Legend(p Props, children ...any) VNode { return tag("legend", p, children...) }

// 表格元素
func Table(p Props, children ...any) VNode { return tag("table", p, children...) }

func Thead(p Props, children ...any) VNode { return tag("thead", p, children...) }

func Tbody(p Props, children ...any) VNode { return tag("tbody", p, children...) }

func Tfoot(p Props, children ...any) VNode { return tag("tfoot", p, children...) }

func Tr(p Props, children ...any) VNode { return tag("tr", p, children...) }

func Th(p Props, children ...any) VNode { return tag("th", p, children...) }

func Td(p Props, children ...any) VNode { return tag("td", p, children...) }

func Caption(p Props, children ...any) VNode { return tag("caption", p, children...) }

func Colgroup(p Props, children ...any) VNode { return tag("colgroup", p, children...) }

func Col(p Props, children ...any) VNode { return tag("col", p, children...) }

// 列表元素
func Ul(p Props, children ...any) VNode { return tag("ul", p, children...) }

func Ol(p Props, children ...any) VNode { return tag("ol", p, children...) }

func Li(p Props, children ...any) VNode { return tag("li", p, children...) }

func Dl(p Props, children ...any) VNode { return tag("dl", p, children...) }

func Dt(p Props, children ...any) VNode { return tag("dt", p, children...) }

func Dd(p Props, children ...any) VNode { return tag("dd", p, children...) }

// 媒體元素
func Img(p Props, children ...any) VNode { return tag("img", p, children...) }

func Audio(p Props, children ...any) VNode { return tag("audio", p, children...) }

func Video(p Props, children ...any) VNode { return tag("video", p, children...) }

func Source(p Props, children ...any) VNode { return tag("source", p, children...) }

func Track(p Props, children ...any) VNode { return tag("track", p, children...) }

func Map(p Props, children ...any) VNode { return tag("map", p, children...) }

func Area(p Props, children ...any) VNode { return tag("area", p, children...) }

func Canvas(p Props, children ...any) VNode { return tag("canvas", p, children...) }

func Figure(p Props, children ...any) VNode { return tag("figure", p, children...) }

func Figcaption(p Props, children ...any) VNode { return tag("figcaption", p, children...) }

func Picture(p Props, children ...any) VNode { return tag("picture", p, children...) }

func Svg(p Props, children ...any) VNode { return tag("svg", p, children...) }

// 互動元素
func A(p Props, children ...any) VNode { return tag("a", p, children...) }

func Details(p Props, children ...any) VNode { return tag("details", p, children...) }

func Summary(p Props, children ...any) VNode { return tag("summary", p, children...) }

func Dialog(p Props, children ...any) VNode { return tag("dialog", p, children...) }

func Menu(p Props, children ...any) VNode { return tag("menu", p, children...) }

// 腳本元素
func Script(p Props, children ...any) VNode {
	var content strings.Builder
	props := p
	if props == nil {
		props = make(Props)
	}

	for _, child := range children {
		switch v := child.(type) {
		case Props:
			for k, val := range v {
				props[k] = val
			}
		case string:
			content.WriteString(v)
		case JSAction:
			content.WriteString(v.Code)
		case VNode:
			// If someone passes VNode(s) as script children, render them and append.
			content.WriteString(Render(v))
		}
	}

	return VNode{
		Tag:     "script",
		Props:   props,
		Content: content.String(),
	}
}
