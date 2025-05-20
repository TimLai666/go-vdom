// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/vdom" // 使用 dot import

	// 注意：不要用 dot import 匯入 control，請用 control.xxx
	control "github.com/TimLai666/go-vdom/control"
)

func main() {
	// 定義一個卡片組件
	Card := Component(
		Div(
			Props{"class": "card"},
			H1("{{title}}"),
			Div("{{children}}"),
		),
	)

	// 測試 control: If/Then/Else/Repeat
	show := false
	items := []string{"蘋果", "香蕉", "橘子"}

	// 處理HTTP請求的函數
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 設置內容類型為HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 使用 Document 函數創建一個完整的 HTML 文檔
		doc := Document(
			"我的網頁", // 頁面標題
			[]LinkInfo{ // 外部資源鏈接
				{Rel: "stylesheet", Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css", Type: "text/css"},
				{Rel: "icon", Href: "https://golang.org/favicon.ico", Type: "image/x-icon"},
			},
			[]ScriptInfo{ // JavaScript 腳本
				{Src: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js", Async: true},
			},
			[]Props{ // Meta 標籤
				{"name": "description", "content": "這是一個使用go-vdom創建的示例網頁"},
				{"name": "keywords", "content": "Go, vdom, HTML"},
			},
			// 頁面主體內容
			Header(
				Props{"class": "container bg-light p-4 mb-4"},
				H1("Go VDOM 示範網站", Props{"class": "text-primary"}),
				Nav(
					Props{"class": "navbar navbar-expand"},
					Ul(
						Props{"class": "navbar-nav"},
						Li(Props{"class": "nav-item"}, A(Props{"href": "#", "class": "nav-link"}, "首頁")),
						Li(Props{"class": "nav-item"}, A(Props{"href": "#about", "class": "nav-link"}, "關於我們")),
						Li(Props{"class": "nav-item"}, A(Props{"href": "#contact", "class": "nav-link"}, "聯繫我們")),
					),
				),
			),
			Main(
				Props{"class": "container"},
				H2("歡迎訪問", Props{"class": "mt-4"}),
				P("這是一個使用 go-vdom 創建的網頁示例。該頁面直接從Go HTTP伺服器產生。"),
				Card(Props{"title": "使用組件"},
					Div(Props{"class": "mb-3"}, "這是卡片的內容，展示了組件的使用方法"),
					Div(Props{"class": "text-muted"}, "可以在組件中傳入多個子元素"),
				),
				Div(
					Props{"class": "row mt-4"},
					Div(Props{"class": "col-md-6"},
						H3("左側內容"),
						P("這是左側的一些內容，展示了多欄排版的效果。"),
						Ul(
							Li("項目1"),
							Li("項目2"),
							Li("項目3"),
						),
					),
					Div(Props{"class": "col-md-6"},
						H3("右側內容"),
						P("這是右側的一些內容。"),
						A(Props{"href": "https://github.com/TimLai666/go-vdom", "class": "btn btn-primary"},
							"查看源碼",
						),
					),
				),
				// 測試 control 區塊
				Div(
					Props{"class": "mt-4"},
					H4("If/Then/Else 測試"),
					control.If(show,
						control.Then(
							Div(Props{"class": "alert alert-success"},
								"這是 If 條件為真時顯示的內容",
							),
						),
						control.Else(
							Div(Props{"class": "alert alert-warning"},
								"這是 If 條件為假時顯示的內容",
							),
						),
					),
					Div(Props{"class": "mt-4"},
						H4("Repeat 測試"),
						control.Repeat(3, func(i int) VNode {
							return Div(Props{"class": "border p-2 mb-2"},
								fmt.Sprintf("Repeat 渲染第 %d 次", i+1),
							)
						}),
					),
					Div(Props{"class": "mt-4"},
						H4("For 測試"),
						Ul(control.For(items, func(item string, i int) VNode {
							return Li(fmt.Sprintf("第%d個: %s", i+1, item))
						})),
					),
				),
			),
			Footer(Props{"class": "container bg-light p-4 mt-4"},
				P(Props{"class": "text-center"},
					"© 2025 ", Span(Props{"style": "color:red;"}, "Go VDOM"), " 示範網站 | 使用 Go 和 VDOM 製作",
				),
			),
			Script(Props{"type": "module"},
				jsdsl.DomReady(
					jsdsl.QueryEach(jsdsl.Els(".navbar-nav .nav-link"), func(el jsdsl.Elem) JSAction {
						return jsdsl.OnClick(el, jsdsl.Alert(`'你點擊了 ' + `+jsdsl.InnerText(el)))
					}),
				),
			),
		)

		// 渲染 HTML 並寫入 HTTP 回應
		html := Render(doc)
		fmt.Fprint(w, html)
	})

	// 啟動 HTTP 伺服器
	port := ":8080"
	fmt.Printf("伺服器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
