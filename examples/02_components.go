// examples/02_components.go
// 組件示例 - 展示如何創建和使用可重用的組件

package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/dom"
)

func main() {
	// 定義一個簡單的 Alert 組件
	Alert := Component(
		Div(
			Props{
				"class": "alert alert-{{type}} {{className}}",
				"role":  "alert",
			},
			Strong("{{title}}"),
			Span(" {{message}}"),
		),
		nil,
		PropsDef{
			"type":      "info", // info, success, warning, danger
			"title":     "",
			"message":   "",
			"className": "",
		},
	)

	// 定義一個 Card 組件
	Card := Component(
		Div(
			Props{"class": "card mb-3 {{className}}", "style": "{{style}}"},
			Div(
				Props{"class": "card-header bg-{{variant}}"},
				H5(Props{"class": "card-title mb-0"}, "{{title}}"),
			),
			Div(
				Props{"class": "card-body"},
				P(Props{"class": "card-text"}, "{{description}}"),
				Div("{{children}}"),
			),
			Div(
				Props{"class": "card-footer text-muted"},
				Small("{{footer}}"),
			),
		),
		nil,
		PropsDef{
			"title":       "卡片標題",
			"description": "",
			"footer":      "",
			"variant":     "primary", // primary, secondary, success, danger, warning, info
			"className":   "",
			"style":       "",
		},
	)

	// 定義一個 Badge 組件
	Badge := Component(
		Span(
			Props{"class": "badge bg-{{color}} {{className}}"},
			"{{text}}",
		),
		nil,
		PropsDef{
			"text":      "Badge",
			"color":     "primary",
			"className": "",
		},
	)

	// 定義一個 Button 組件
	CustomButton := Component(
		Button(
			Props{
				"type":  "{{type}}",
				"class": "btn btn-{{variant}} {{size}} {{className}}",
			},
			"{{text}}",
		),
		nil,
		PropsDef{
			"text":      "Button",
			"type":      "button",
			"variant":   "primary",
			"size":      "",
			"className": "",
		},
	)

	// 定義一個 UserCard 組件
	UserCard := Component(
		Div(
			Props{"class": "card text-center {{className}}"},
			Div(
				Props{"class": "card-body"},
				Img(
					Props{
						"src":   "{{avatar}}",
						"alt":   "{{name}}",
						"class": "rounded-circle mb-3",
						"style": "width: 100px; height: 100px; object-fit: cover;",
					},
				),
				H5(Props{"class": "card-title"}, "{{name}}"),
				P(Props{"class": "card-text text-muted"}, "{{role}}"),
				P(Props{"class": "card-text"}, "{{bio}}"),
				Div("{{children}}"),
			),
		),
		nil,
		PropsDef{
			"name":      "使用者",
			"avatar":    "https://via.placeholder.com/100",
			"role":      "角色",
			"bio":       "",
			"className": "",
		},
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		doc := Document(
			"組件示例",
			[]LinkInfo{
				{
					Rel:  "stylesheet",
					Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
					Type: "text/css",
				},
			},
			nil,
			[]Props{
				{"name": "viewport", "content": "width=device-width, initial-scale=1"},
				{"charset": "UTF-8"},
			},
			Div(
				Props{"class": "container mt-5"},

				H1(Props{"class": "text-primary mb-4"}, "Go VDOM 組件示例"),

				// Alert 組件展示
				H2(Props{"class": "mt-4 mb-3"}, "Alert 組件"),
				Alert(Props{
					"type":    "success",
					"title":   "成功！",
					"message": "您的操作已成功完成。",
				}),
				Alert(Props{
					"type":    "warning",
					"title":   "警告！",
					"message": "請注意這個重要訊息。",
				}),
				Alert(Props{
					"type":    "danger",
					"title":   "錯誤！",
					"message": "操作失敗，請重試。",
				}),
				Alert(Props{
					"type":    "info",
					"title":   "提示：",
					"message": "這是一條資訊提示。",
				}),

				// Card 組件展示
				H2(Props{"class": "mt-5 mb-3"}, "Card 組件"),
				Div(
					Props{"class": "row"},
					Div(
						Props{"class": "col-md-4"},
						Card(
							Props{
								"title":       "基礎卡片",
								"description": "這是一個基礎的卡片組件",
								"footer":      "2025-01-24",
								"variant":     "primary",
							},
							CustomButton(Props{
								"text":    "查看更多",
								"variant": "primary",
								"size":    "btn-sm",
							}),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						Card(
							Props{
								"title":       "進階卡片",
								"description": "這個卡片包含更多內容",
								"footer":      "最後更新: 今天",
								"variant":     "success",
							},
							Badge(Props{"text": "新功能", "color": "success"}),
							Span(" "),
							Badge(Props{"text": "熱門", "color": "danger"}),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						Card(
							Props{
								"title":       "自定義樣式",
								"description": "可以自定義各種樣式",
								"footer":      "v1.0.0",
								"variant":     "warning",
								"className":   "shadow",
							},
						),
					),
				),

				// Badge 組件展示
				H2(Props{"class": "mt-5 mb-3"}, "Badge 組件"),
				Div(
					Badge(Props{"text": "Primary", "color": "primary"}),
					Span(" "),
					Badge(Props{"text": "Secondary", "color": "secondary"}),
					Span(" "),
					Badge(Props{"text": "Success", "color": "success"}),
					Span(" "),
					Badge(Props{"text": "Danger", "color": "danger"}),
					Span(" "),
					Badge(Props{"text": "Warning", "color": "warning"}),
					Span(" "),
					Badge(Props{"text": "Info", "color": "info"}),
					Span(" "),
					Badge(Props{"text": "Light", "color": "light"}),
					Span(" "),
					Badge(Props{"text": "Dark", "color": "dark"}),
				),

				// Button 組件展示
				H2(Props{"class": "mt-5 mb-3"}, "Button 組件"),
				Div(
					Props{"class": "mb-3"},
					CustomButton(Props{"text": "Primary", "variant": "primary"}),
					Span(" "),
					CustomButton(Props{"text": "Secondary", "variant": "secondary"}),
					Span(" "),
					CustomButton(Props{"text": "Success", "variant": "success"}),
					Span(" "),
					CustomButton(Props{"text": "Danger", "variant": "danger"}),
				),
				Div(
					Props{"class": "mb-3"},
					CustomButton(Props{"text": "Large", "variant": "primary", "size": "btn-lg"}),
					Span(" "),
					CustomButton(Props{"text": "Normal", "variant": "primary"}),
					Span(" "),
					CustomButton(Props{"text": "Small", "variant": "primary", "size": "btn-sm"}),
				),

				// UserCard 組件展示
				H2(Props{"class": "mt-5 mb-3"}, "UserCard 組件"),
				Div(
					Props{"class": "row"},
					Div(
						Props{"class": "col-md-4"},
						UserCard(
							Props{
								"name":   "張小明",
								"avatar": "https://i.pravatar.cc/100?img=1",
								"role":   "前端工程師",
								"bio":    "熱愛編程，專注於前端開發",
							},
							CustomButton(Props{
								"text":    "查看資料",
								"variant": "primary",
								"size":    "btn-sm",
							}),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						UserCard(
							Props{
								"name":   "李小華",
								"avatar": "https://i.pravatar.cc/100?img=5",
								"role":   "後端工程師",
								"bio":    "擅長 Go 語言和微服務架構",
							},
							CustomButton(Props{
								"text":    "查看資料",
								"variant": "success",
								"size":    "btn-sm",
							}),
						),
					),
					Div(
						Props{"class": "col-md-4"},
						UserCard(
							Props{
								"name":   "王小美",
								"avatar": "https://i.pravatar.cc/100?img=9",
								"role":   "UI/UX 設計師",
								"bio":    "致力於創造優秀的用戶體驗",
							},
							CustomButton(Props{
								"text":    "查看資料",
								"variant": "info",
								"size":    "btn-sm",
							}),
						),
					),
				),

				// 組件組合示例
				H2(Props{"class": "mt-5 mb-3"}, "組件組合"),
				Card(
					Props{
						"title":       "團隊成員",
						"description": "這是一個組合多個組件的示例",
						"footer":      "共 3 位成員",
						"variant":     "info",
					},
					Div(
						Props{"class": "d-flex gap-2"},
						Badge(Props{"text": "3 人", "color": "primary"}),
						Badge(Props{"text": "全職", "color": "success"}),
						Badge(Props{"text": "遠程", "color": "info"}),
					),
				),

				// 代碼示例
				H2(Props{"class": "mt-5 mb-3"}, "組件定義示例"),
				Pre(
					Props{"class": "bg-light p-3"},
					Code(`// 定義組件
Card := Component(
    Div(
        Props{"class": "card"},
        H5("{{title}}"),
        P("{{content}}"),
    ),
    nil,
    PropsDef{
        "title": "默認標題",
        "content": "",
    },
)

// 使用組件
myCard := Card(Props{
    "title": "我的卡片",
    "content": "這是內容",
})`),
				),

				// 頁腳
				Footer(
					Props{"class": "mt-5 pt-4 border-top text-center text-muted"},
					P("© 2025 Go VDOM 組件示例"),
				),
			),
		)

		html := Render(doc)
		fmt.Fprint(w, html)
	})

	port := ":8081"
	fmt.Printf("組件示例服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
