// examples/01_basic_usage.go
// 基本用法示例 - 展示如何創建簡單的 HTML 頁面

package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/dom"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 創建一個簡單的 HTML 文檔
		doc := Document(
			"基本用法示例", // 頁面標題
			[]LinkInfo{
				{
					Rel:  "stylesheet",
					Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
					Type: "text/css",
				},
			},
			nil, // 無額外腳本
			[]Props{
				{"name": "viewport", "content": "width=device-width, initial-scale=1"},
				{"charset": "UTF-8"},
			},
			// 頁面內容
			Div(
				Props{"class": "container mt-5"},

				// 標題
				H1(Props{"class": "text-primary"}, "Go VDOM 基本用法"),

				// 副標題
				H2(Props{"class": "mt-4"}, "歡迎使用 go-vdom"),

				// 段落
				P(Props{"class": "lead"},
					"這是一個使用 go-vdom 創建的簡單網頁示例。",
					"所有的 HTML 都是在 Go 中生成的。",
				),

				// 分隔線
				Hr(),

				// 內容區塊
				Div(
					Props{"class": "row mt-4"},

					// 左側欄
					Div(
						Props{"class": "col-md-6"},
						H3("特性列表"),
						Ul(
							Li("類型安全的 HTML 生成"),
							Li("聲明式的 API"),
							Li("組件化設計"),
							Li("服務器端渲染"),
						),
					),

					// 右側欄
					Div(
						Props{"class": "col-md-6"},
						H3("快速開始"),
						P("只需幾行代碼即可創建網頁："),
						Pre(
							Props{"class": "bg-light p-3"},
							Code(`doc := Document(
    "標題",
    nil, nil, nil,
    Div("內容"),
)`),
						),
					),
				),

				// 卡片展示
				Div(
					Props{"class": "row mt-4"},
					Div(
						Props{"class": "col-md-4"},
						Div(
							Props{"class": "card"},
							Div(
								Props{"class": "card-body"},
								H5(Props{"class": "card-title"}, "簡單"),
								P(Props{"class": "card-text"}, "簡潔的 API，易於學習和使用"),
							),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						Div(
							Props{"class": "card"},
							Div(
								Props{"class": "card-body"},
								H5(Props{"class": "card-title"}, "強大"),
								P(Props{"class": "card-text"}, "支持組件、控制流、JavaScript DSL"),
							),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						Div(
							Props{"class": "card"},
							Div(
								Props{"class": "card-body"},
								H5(Props{"class": "card-title"}, "高效"),
								P(Props{"class": "card-text"}, "零運行時開銷，純靜態 HTML 生成"),
							),
						),
					),
				),

				// 連結
				Div(
					Props{"class": "mt-5 text-center"},
					A(
						Props{
							"href":  "https://github.com/TimLai666/go-vdom",
							"class": "btn btn-primary btn-lg",
						},
						"查看 GitHub 倉庫",
					),
				),

				// 頁腳
				Footer(
					Props{"class": "mt-5 pt-4 border-top text-center text-muted"},
					P("© 2025 Go VDOM - MIT License"),
				),
			),
		)

		// 渲染並輸出 HTML
		html := Render(doc)
		fmt.Fprint(w, html)
	})

	// 啟動服務器
	port := ":8080"
	fmt.Printf("服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
