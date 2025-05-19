// tags.go
package vdom

// 基本標籤函數

// Text 創建一個文字節點
func Text(s string) VNode {
	return VNode{Content: s}
}

// 常用 HTML 標籤函數

// tag 是所有 HTML 標籤函數的基礎實現
func tag(name string, args ...any) VNode {
	var props Props
	var children []VNode
	var content string

	for _, arg := range args {
		switch v := arg.(type) {
		case Props:
			props = v
		case VNode:
			children = append(children, v)
		case string:
			children = append(children, Text(v))
		}
	}

	return VNode{
		Tag:      name,
		Props:    props,
		Children: children,
		Content:  content,
	}
}

// HTML 標籤函數
// HTML 結構元素
func Html(args ...any) VNode    { return tag("html", args...) }
func Head(args ...any) VNode    { return tag("head", args...) }
func Body(args ...any) VNode    { return tag("body", args...) }
func Main(args ...any) VNode    { return tag("main", args...) }
func Header(args ...any) VNode  { return tag("header", args...) }
func Footer(args ...any) VNode  { return tag("footer", args...) }
func Nav(args ...any) VNode     { return tag("nav", args...) }
func Aside(args ...any) VNode   { return tag("aside", args...) }
func Section(args ...any) VNode { return tag("section", args...) }
func Article(args ...any) VNode { return tag("article", args...) }

// 文本區塊元素
func H1(args ...any) VNode         { return tag("h1", args...) }
func H2(args ...any) VNode         { return tag("h2", args...) }
func H3(args ...any) VNode         { return tag("h3", args...) }
func H4(args ...any) VNode         { return tag("h4", args...) }
func H5(args ...any) VNode         { return tag("h5", args...) }
func H6(args ...any) VNode         { return tag("h6", args...) }
func P(args ...any) VNode          { return tag("p", args...) }
func Div(args ...any) VNode        { return tag("div", args...) }
func Span(args ...any) VNode       { return tag("span", args...) }
func Pre(args ...any) VNode        { return tag("pre", args...) }
func Code(args ...any) VNode       { return tag("code", args...) }
func Blockquote(args ...any) VNode { return tag("blockquote", args...) }

// 表單元素
func Form(args ...any) VNode     { return tag("form", args...) }
func Input(args ...any) VNode    { return tag("input", args...) }
func Label(args ...any) VNode    { return tag("label", args...) }
func Button(args ...any) VNode   { return tag("button", args...) }
func Select(args ...any) VNode   { return tag("select", args...) }
func Option(args ...any) VNode   { return tag("option", args...) }
func Textarea(args ...any) VNode { return tag("textarea", args...) }
func Fieldset(args ...any) VNode { return tag("fieldset", args...) }
func Legend(args ...any) VNode   { return tag("legend", args...) }

// 表格元素
func Table(args ...any) VNode   { return tag("table", args...) }
func Thead(args ...any) VNode   { return tag("thead", args...) }
func Tbody(args ...any) VNode   { return tag("tbody", args...) }
func Tfoot(args ...any) VNode   { return tag("tfoot", args...) }
func Tr(args ...any) VNode      { return tag("tr", args...) }
func Th(args ...any) VNode      { return tag("th", args...) }
func Td(args ...any) VNode      { return tag("td", args...) }
func Caption(args ...any) VNode { return tag("caption", args...) }

// 列表元素
func Ul(args ...any) VNode { return tag("ul", args...) }
func Ol(args ...any) VNode { return tag("ol", args...) }
func Li(args ...any) VNode { return tag("li", args...) }
func Dl(args ...any) VNode { return tag("dl", args...) }
func Dt(args ...any) VNode { return tag("dt", args...) }
func Dd(args ...any) VNode { return tag("dd", args...) }

// 媒體元素
func Img(args ...any) VNode    { return tag("img", args...) }
func Audio(args ...any) VNode  { return tag("audio", args...) }
func Video(args ...any) VNode  { return tag("video", args...) }
func Source(args ...any) VNode { return tag("source", args...) }
func Canvas(args ...any) VNode { return tag("canvas", args...) }
func Svg(args ...any) VNode    { return tag("svg", args...) }

// 互動元素
func A(args ...any) VNode       { return tag("a", args...) }
func Details(args ...any) VNode { return tag("details", args...) }
func Summary(args ...any) VNode { return tag("summary", args...) }
func Dialog(args ...any) VNode  { return tag("dialog", args...) }
func Menu(args ...any) VNode    { return tag("menu", args...) }

// 腳本元素
func Script(args ...any) VNode { return tag("script", args...) }
func Style(args ...any) VNode  { return tag("style", args...) }
func Link(args ...any) VNode   { return tag("link", args...) }
func Meta(args ...any) VNode   { return tag("meta", args...) }
func Title(args ...any) VNode  { return tag("title", args...) }

// 其他常用元素
func Br(args ...any) VNode       { return tag("br", args...) }
func Hr(args ...any) VNode       { return tag("hr", args...) }
func Iframe(args ...any) VNode   { return tag("iframe", args...) }
func Noscript(args ...any) VNode { return tag("noscript", args...) }
func Time(args ...any) VNode     { return tag("time", args...) }
func Abbr(args ...any) VNode     { return tag("abbr", args...) }
func Strong(args ...any) VNode   { return tag("strong", args...) }
func Em(args ...any) VNode       { return tag("em", args...) }
