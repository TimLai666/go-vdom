// examples/03_javascript_dsl.go
// JavaScript DSL 示例 - 展示如何使用 DSL 生成交互式 JavaScript 代碼

package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/dom"
	js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		doc := Document(
			"JavaScript DSL 示例",
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

				H1(Props{"class": "text-primary mb-4"}, "JavaScript DSL 示例"),

				// 1. 基本 DOM 操作
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("1. 基本 DOM 操作"),
					),
					Div(
						Props{"class": "card-body"},
						P("點擊按鈕改變文本和樣式："),
						Button(Props{
							"id":    "btn1",
							"class": "btn btn-primary me-2",
							"onClick": js.Do(nil,
								js.El("#text1").SetText("'文本已改變！'"),
								js.El("#text1").AddClass("text-success"),
								js.El("#text1").AddClass("fw-bold"),
							),
						}, "改變文本"),
						Button(Props{
							"id":    "btn2",
							"class": "btn btn-secondary me-2",
							"onClick": js.Do(nil,
								js.El("#text1").SetHTML("'<em>HTML 已更新</em>'"),
								js.El("#text1").RemoveClass("text-success"),
								js.El("#text1").RemoveClass("fw-bold"),
							),
						}, "設置 HTML"),
						Button(Props{
							"id":    "btn3",
							"class": "btn btn-warning",
							"onClick": js.Do(nil,
								js.El("#text1").SetText("'原始文本'"),
								js.El("#text1").RemoveClass("text-success"),
								js.El("#text1").RemoveClass("fw-bold"),
							),
						}, "重置"),
						Div(Props{
							"id":    "text1",
							"class": "mt-3 p-3 border rounded",
						}, "原始文本"),
					),
				),

				// 2. 變數操作
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("2. 變數操作與計數器"),
					),
					Div(
						Props{"class": "card-body"},
						P("使用 let 和 const 定義變數："),
						Div(
							Props{"class": "d-flex gap-2 mb-3"},
							Button(Props{
								"class": "btn btn-success",
								"onClick": js.Do(nil,
									js.Const("currentValue", "parseInt(document.getElementById('counter').innerText)"),
									js.Const("newValue", "currentValue + 1"),
									js.El("#counter").SetText("newValue.toString()"),
									js.Log("'計數器增加: ' + newValue"),
								),
							}, "增加 +1"),
							Button(Props{
								"class": "btn btn-danger",
								"onClick": js.Do(nil,
									js.Const("currentValue", "parseInt(document.getElementById('counter').innerText)"),
									js.Const("newValue", "currentValue - 1"),
									js.El("#counter").SetText("newValue.toString()"),
									js.Log("'計數器減少: ' + newValue"),
								),
							}, "減少 -1"),
							Button(Props{
								"class": "btn btn-secondary",
								"onClick": js.Do(nil,
									js.El("#counter").SetText("'0'"),
									js.Log("'計數器已重置'"),
								),
							}, "重置"),
						),
						Div(
							Props{"class": "alert alert-info"},
							H4("計數器: ", Span(Props{"id": "counter", "class": "badge bg-primary"}, "0")),
						),
					),
				),

				// 3. 控制台日誌和警告框
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("3. 控制台日誌與警告框"),
					),
					Div(
						Props{"class": "card-body"},
						P("打開瀏覽器控制台查看日誌輸出："),
						Button(Props{
							"class": "btn btn-info me-2",
							"onClick": js.Do(nil,
								js.Log("'這是一條控制台日誌'"),
								js.Log("'當前時間:', new Date().toLocaleString()"),
							),
						}, "輸出日誌"),
						Button(Props{
							"class": "btn btn-warning me-2",
							"onClick": js.Do(nil,
								js.Alert("'這是一個警告框！'"),
							),
						}, "顯示警告框"),
						Button(Props{
							"class": "btn btn-secondary",
							"onClick": js.Do(nil,
								js.Const("userName", "prompt('請輸入您的名字:')"),
								JSAction{Code: "if (userName) { alert('您好, ' + userName + '!'); }"},
							),
						}, "輸入對話框"),
					),
				),

				// 4. 表單處理
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("4. 表單處理"),
					),
					Div(
						Props{"class": "card-body"},
						Form(Props{
							"id": "demoForm",
							"onSubmit": js.Do([]string{"event"},
								js.CallMethod("event", "preventDefault"),
								js.Log("'表單提交事件觸發'"),
								js.Const("nameValue", "document.getElementById('nameInput').value"),
								js.Const("emailValue", "document.getElementById('emailInput').value"),
								js.Const("messageValue", "document.getElementById('messageInput').value"),
								js.Log("'姓名:', nameValue"),
								js.Log("'郵箱:', emailValue"),
								js.Log("'訊息:', messageValue"),
								js.Const("resultHTML", "`<strong>提交成功！</strong><br>姓名: ${nameValue}<br>郵箱: ${emailValue}<br>訊息: ${messageValue}`"),
								js.El("#formResult").SetHTML("resultHTML"),
								js.El("#formResult").AddClass("alert-success"),
								js.El("#formResult").RemoveClass("d-none"),
							),
						},
							Div(
								Props{"class": "mb-3"},
								Label(Props{"for": "nameInput", "class": "form-label"}, "姓名"),
								Input(Props{
									"type":        "text",
									"class":       "form-control",
									"id":          "nameInput",
									"placeholder": "請輸入姓名",
									"required":    "true",
								}),
							),
							Div(
								Props{"class": "mb-3"},
								Label(Props{"for": "emailInput", "class": "form-label"}, "電子郵件"),
								Input(Props{
									"type":        "email",
									"class":       "form-control",
									"id":          "emailInput",
									"placeholder": "請輸入郵箱",
									"required":    "true",
								}),
							),
							Div(
								Props{"class": "mb-3"},
								Label(Props{"for": "messageInput", "class": "form-label"}, "訊息"),
								Textarea(Props{
									"class":       "form-control",
									"id":          "messageInput",
									"rows":        "3",
									"placeholder": "請輸入訊息",
								}),
							),
							Button(Props{
								"type":  "submit",
								"class": "btn btn-primary",
							}, "提交表單"),
						),
						Div(Props{
							"id":    "formResult",
							"class": "alert mt-3 d-none",
						}),
					),
				),

				// 5. 動態創建元素
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("5. 動態創建元素"),
					),
					Div(
						Props{"class": "card-body"},
						P("動態添加元素到頁面："),
						Input(Props{
							"type":        "text",
							"class":       "form-control mb-2",
							"id":          "itemInput",
							"placeholder": "輸入項目名稱",
						}),
						Button(Props{
							"class": "btn btn-success me-2",
							"onClick": js.Do(nil,
								js.Const("itemText", "document.getElementById('itemInput').value"),
								JSAction{Code: "if (!itemText.trim()) { alert('請輸入項目名稱'); return; }"},
								js.Const("newItem", "document.createElement('div')"),
								JSAction{Code: "newItem.className = 'alert alert-info alert-dismissible fade show'"},
								JSAction{Code: "newItem.innerHTML = itemText + '<button type=\"button\" class=\"btn-close\" data-bs-dismiss=\"alert\"></button>'"},
								js.Const("container", "document.getElementById('itemsContainer')"),
								JSAction{Code: "container.appendChild(newItem)"},
								js.El("#itemInput").SetText("''"),
								js.Log("'添加項目: ' + itemText"),
							),
						}, "添加項目"),
						Button(Props{
							"class": "btn btn-danger",
							"onClick": js.Do(nil,
								js.El("#itemsContainer").SetHTML("''"),
								js.Log("'清空所有項目'"),
							),
						}, "清空列表"),
						Div(Props{
							"id":    "itemsContainer",
							"class": "mt-3",
						}),
					),
				),

				// 6. Try-Catch 錯誤處理
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("6. Try/Catch 錯誤處理"),
					),
					Div(
						Props{"class": "card-body"},
						P("模擬異步操作和錯誤處理："),
						Button(Props{
							"class": "btn btn-primary me-2",
							"onClick": js.AsyncDo(nil,
								js.Try(
									js.Log("'開始異步操作...'"),
									js.El("#asyncResult").SetHTML("'<div class=\"spinner-border spinner-border-sm\"></div> 處理中...'"),
									js.Const("delay", "new Promise(resolve => setTimeout(resolve, 2000))"),
									JSAction{Code: "await delay"},
									js.Const("randomNum", "Math.random()"),
									JSAction{Code: "if (randomNum < 0.5) throw new Error('隨機錯誤發生 (運氣不好)')"},
									js.El("#asyncResult").SetHTML("'<div class=\"alert alert-success\">操作成功完成！</div>'"),
									js.Log("'異步操作完成'"),
								).Catch("error",
									js.Log("'捕獲錯誤:', error"),
									js.El("#asyncResult").SetHTML("'<div class=\"alert alert-danger\">錯誤: ' + error.message + '</div>'"),
								).End(),
							),
						}, "執行異步操作"),
						Div(Props{"id": "asyncResult", "class": "mt-3"}),
					),
				),

				// 7. DomReady 事件
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("7. DomReady 初始化"),
					),
					Div(
						Props{"class": "card-body"},
						P("頁面載入時執行的初始化代碼（查看控制台）："),
						Div(Props{
							"id":    "initStatus",
							"class": "alert alert-info",
						}, "等待初始化..."),
					),
				),

				// 頁腳
				Footer(
					Props{"class": "mt-5 pt-4 border-top text-center text-muted"},
					P("© 2025 Go VDOM JavaScript DSL 示例"),
					P(Props{"class": "small"}, "打開瀏覽器控制台查看更多日誌輸出"),
				),
			),

			// 在頁面底部添加初始化腳本
			Script(Props{"type": "text/javascript"},
				js.DomReady(
					js.Log("'頁面已載入完成'"),
					js.Log("'Go VDOM JavaScript DSL 正常運行'"),
					js.Const("now", "new Date().toLocaleString()"),
					js.Log("'當前時間:', now"),
					js.El("#initStatus").SetHTML("'<strong>✓ 初始化完成</strong><br>時間: ' + now"),
					js.El("#initStatus").AddClass("alert-success"),
					js.El("#initStatus").RemoveClass("alert-info"),
				),
			),
		)

		html := Render(doc)
		fmt.Fprint(w, html)
	})

	port := ":8082"
	fmt.Printf("JavaScript DSL 示例服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
