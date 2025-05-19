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
		case []VNode: // 自動展開
			children = append(children, v...)
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
// - bodyContent: 頁面主體內容
func Document(title string, links []LinkInfo, scripts []ScriptInfo, metas []Props, bodyContent ...VNode) VNode {
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

	// 添加腳本
	for _, script := range scripts {
		props := Props{"src": script.Src}
		if script.Async {
			props["async"] = "true"
		}
		headElements = append(headElements, Script(props))
	}

	// 將body內容轉換為any類型
	var bodyElements []any
	for _, node := range bodyContent {
		bodyElements = append(bodyElements, node)
	}

	return Html(
		Head(headElements...),
		Body(bodyElements...),
	)
}

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
func Address(args ...any) VNode { return tag("address", args...) }
func Hgroup(args ...any) VNode  { return tag("hgroup", args...) }

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
func Datalist(args ...any) VNode { return tag("datalist", args...) }
func Optgroup(args ...any) VNode { return tag("optgroup", args...) }
func Option(args ...any) VNode   { return tag("option", args...) }
func Textarea(args ...any) VNode { return tag("textarea", args...) }
func Output(args ...any) VNode   { return tag("output", args...) }
func Progress(args ...any) VNode { return tag("progress", args...) }
func Meter(args ...any) VNode    { return tag("meter", args...) }
func Fieldset(args ...any) VNode { return tag("fieldset", args...) }
func Legend(args ...any) VNode   { return tag("legend", args...) }

// 表格元素
func Table(args ...any) VNode    { return tag("table", args...) }
func Thead(args ...any) VNode    { return tag("thead", args...) }
func Tbody(args ...any) VNode    { return tag("tbody", args...) }
func Tfoot(args ...any) VNode    { return tag("tfoot", args...) }
func Tr(args ...any) VNode       { return tag("tr", args...) }
func Th(args ...any) VNode       { return tag("th", args...) }
func Td(args ...any) VNode       { return tag("td", args...) }
func Caption(args ...any) VNode  { return tag("caption", args...) }
func Colgroup(args ...any) VNode { return tag("colgroup", args...) }
func Col(args ...any) VNode      { return tag("col", args...) }

// 列表元素
func Ul(args ...any) VNode { return tag("ul", args...) }
func Ol(args ...any) VNode { return tag("ol", args...) }
func Li(args ...any) VNode { return tag("li", args...) }
func Dl(args ...any) VNode { return tag("dl", args...) }
func Dt(args ...any) VNode { return tag("dt", args...) }
func Dd(args ...any) VNode { return tag("dd", args...) }

// 媒體元素
func Img(args ...any) VNode        { return tag("img", args...) }
func Audio(args ...any) VNode      { return tag("audio", args...) }
func Video(args ...any) VNode      { return tag("video", args...) }
func Source(args ...any) VNode     { return tag("source", args...) }
func Track(args ...any) VNode      { return tag("track", args...) }
func Map(args ...any) VNode        { return tag("map", args...) }
func Area(args ...any) VNode       { return tag("area", args...) }
func Canvas(args ...any) VNode     { return tag("canvas", args...) }
func Figure(args ...any) VNode     { return tag("figure", args...) }
func Figcaption(args ...any) VNode { return tag("figcaption", args...) }
func Picture(args ...any) VNode    { return tag("picture", args...) }
func Svg(args ...any) VNode        { return tag("svg", args...) }

// 互動元素
func A(args ...any) VNode       { return tag("a", args...) }
func Details(args ...any) VNode { return tag("details", args...) }
func Summary(args ...any) VNode { return tag("summary", args...) }
func Dialog(args ...any) VNode  { return tag("dialog", args...) }
func Menu(args ...any) VNode    { return tag("menu", args...) }

// 腳本元素
func Script(args ...any) VNode   { return tag("script", args...) }
func Noscript(args ...any) VNode { return tag("noscript", args...) }
func Template(args ...any) VNode { return tag("template", args...) }
func Slot(args ...any) VNode     { return tag("slot", args...) }
func Style(args ...any) VNode    { return tag("style", args...) }
func Link(args ...any) VNode     { return tag("link", args...) }
func Meta(args ...any) VNode     { return tag("meta", args...) }
func Title(args ...any) VNode    { return tag("title", args...) }
func Base(args ...any) VNode     { return tag("base", args...) }

// 內聯文本元素
func Abbr(args ...any) VNode   { return tag("abbr", args...) }
func B(args ...any) VNode      { return tag("b", args...) }
func Bdi(args ...any) VNode    { return tag("bdi", args...) }
func Bdo(args ...any) VNode    { return tag("bdo", args...) }
func Br(args ...any) VNode     { return tag("br", args...) }
func Cite(args ...any) VNode   { return tag("cite", args...) }
func Data(args ...any) VNode   { return tag("data", args...) }
func Dfn(args ...any) VNode    { return tag("dfn", args...) }
func Em(args ...any) VNode     { return tag("em", args...) }
func I(args ...any) VNode      { return tag("i", args...) }
func Kbd(args ...any) VNode    { return tag("kbd", args...) }
func Mark(args ...any) VNode   { return tag("mark", args...) }
func Q(args ...any) VNode      { return tag("q", args...) }
func Rb(args ...any) VNode     { return tag("rb", args...) }
func Rp(args ...any) VNode     { return tag("rp", args...) }
func Rt(args ...any) VNode     { return tag("rt", args...) }
func Rtc(args ...any) VNode    { return tag("rtc", args...) }
func Ruby(args ...any) VNode   { return tag("ruby", args...) }
func S(args ...any) VNode      { return tag("s", args...) }
func Samp(args ...any) VNode   { return tag("samp", args...) }
func Small(args ...any) VNode  { return tag("small", args...) }
func Strong(args ...any) VNode { return tag("strong", args...) }
func Sub(args ...any) VNode    { return tag("sub", args...) }
func Sup(args ...any) VNode    { return tag("sup", args...) }
func Time(args ...any) VNode   { return tag("time", args...) }
func U(args ...any) VNode      { return tag("u", args...) }
func Var(args ...any) VNode    { return tag("var", args...) }
func Wbr(args ...any) VNode    { return tag("wbr", args...) }

// 其他元素
func Hr(args ...any) VNode     { return tag("hr", args...) }
func Iframe(args ...any) VNode { return tag("iframe", args...) }
func Object(args ...any) VNode { return tag("object", args...) }
func Param(args ...any) VNode  { return tag("param", args...) }
func Embed(args ...any) VNode  { return tag("embed", args...) }
func Math(args ...any) VNode   { return tag("math", args...) }
