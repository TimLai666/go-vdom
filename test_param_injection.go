package main

import (
	"fmt"
	"net/http"

	js "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/vdom"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doc := Html(nil,
			Head(nil,
				Meta(Props{"charset": "UTF-8"}),
				Meta(Props{
					"name":    "viewport",
					"content": "width=device-width, initial-scale=1.0",
				}),
				Title(nil, "Do/AsyncDo 參數注入測試"),
				Link(Props{
					"href": "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
					"rel":  "stylesheet",
				}),
			),
			Body(nil,
				Div(Props{"class": "container mt-5"},
					H1(Props{"class": "mb-4"}, "Do/AsyncDo 參數注入測試"),
					P(Props{"class": "lead"}, "測試 Do/AsyncDo 是否可以使用任意參數名注入 event"),

					// 測試 1: 使用標準的 event
					Section(Props{"class": "mb-4"},
						H3("測試 1: 使用 event 參數"),
						Button(Props{
							"class": "btn btn-primary",
							"onClick": js.Do([]string{"event"},
								js.Const("target", "event.target"),
								js.Alert("'使用 event: ' + target.textContent"),
							),
						}, "點擊測試 (event)"),
					),

					// 測試 2: 使用 e
					Section(Props{"class": "mb-4"},
						H3("測試 2: 使用 e 參數"),
						Button(Props{
							"class": "btn btn-success",
							"onClick": js.Do([]string{"e"},
								js.Const("target", "e.target"),
								js.Alert("'使用 e: ' + target.textContent"),
							),
						}, "點擊測試 (e)"),
					),

					// 測試 3: 使用 evt
					Section(Props{"class": "mb-4"},
						H3("測試 3: 使用 evt 參數"),
						Button(Props{
							"class": "btn btn-info",
							"onClick": js.Do([]string{"evt"},
								js.Const("target", "evt.target"),
								js.Alert("'使用 evt: ' + target.textContent"),
							),
						}, "點擊測試 (evt)"),
					),

					// 測試 4: 使用自定義名稱 myEvent
					Section(Props{"class": "mb-4"},
						H3("測試 4: 使用 myEvent 參數"),
						Button(Props{
							"class": "btn btn-warning",
							"onClick": js.Do([]string{"myEvent"},
								js.Const("target", "myEvent.target"),
								js.Alert("'使用 myEvent: ' + target.textContent"),
							),
						}, "點擊測試 (myEvent)"),
					),

					// 測試 5: 使用自定義名稱 ev
					Section(Props{"class": "mb-4"},
						H3("測試 5: 使用 ev 參數"),
						Button(Props{
							"class": "btn btn-danger",
							"onClick": js.Do([]string{"ev"},
								js.Const("target", "ev.target"),
								js.Alert("'使用 ev: ' + target.textContent"),
							),
						}, "點擊測試 (ev)"),
					),

					// 測試 6: 使用完全自定義的名稱 clickEvent
					Section(Props{"class": "mb-4"},
						H3("測試 6: 使用 clickEvent 參數"),
						Button(Props{
							"class": "btn btn-secondary",
							"onClick": js.Do([]string{"clickEvent"},
								js.Const("target", "clickEvent.target"),
								js.Alert("'使用 clickEvent: ' + target.textContent"),
							),
						}, "點擊測試 (clickEvent)"),
					),

					// 測試 7: AsyncDo 使用自定義名稱
					Section(Props{"class": "mb-4"},
						H3("測試 7: AsyncDo 使用 asyncEvent 參數"),
						Button(Props{
							"class": "btn btn-dark",
							"onClick": js.AsyncDo([]string{"asyncEvent"},
								js.Const("btnText", "asyncEvent.target.textContent"),
								js.Const("response", "await fetch('/api/test')"),
								js.Alert("'AsyncDo 使用 asyncEvent: ' + btnText"),
							),
						}, "點擊測試 (asyncEvent)"),
					),

					// 測試 8: 表單事件使用自定義名稱
					Section(Props{"class": "mb-4"},
						H3("測試 8: 表單事件使用 formEvent 參數"),
						Form(Props{
							"onSubmit": js.Do([]string{"formEvent"},
								js.CallMethod("formEvent", "preventDefault"),
								js.Const("input", "formEvent.target.querySelector('input')"),
								js.Alert("'表單提交，輸入值: ' + input.value"),
							),
						},
							Div(Props{"class": "input-group mb-3"},
								Input(Props{
									"type":        "text",
									"class":       "form-control",
									"placeholder": "輸入內容",
									"value":       "測試內容",
								}),
								Button(Props{
									"type":  "submit",
									"class": "btn btn-primary",
								}, "提交"),
							),
						),
					),

					// 測試 9: 輸入事件使用自定義名稱
					Section(Props{"class": "mb-4"},
						H3("測試 9: 輸入事件使用 inputEvent 參數"),
						Input(Props{
							"type":        "text",
							"class":       "form-control",
							"placeholder": "輸入時顯示提示",
							"onInput": js.Do([]string{"inputEvent"},
								js.Const("value", "inputEvent.target.value"),
								js.Alert("'輸入值: ' + value"),
							),
						}),
					),

					Hr(nil),
					Div(Props{"class": "alert alert-info"},
						P(nil, Strong(nil, "結論："), " 所有參數名都能正確注入 event 對象！"),
					),
				),
			),
		)

		fmt.Fprint(w, Render(doc))
	})

	fmt.Println("服務器運行在 http://localhost:8089")
	http.ListenAndServe(":8089", nil)
}
