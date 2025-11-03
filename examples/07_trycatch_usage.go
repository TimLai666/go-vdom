package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/dom"
	js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
	http.HandleFunc("/", tryCatchExampleHandler)
	fmt.Println("Server running on http://localhost:8086")
	log.Fatal(http.ListenAndServe(":8086", nil))
}

func tryCatchExampleHandler(w http.ResponseWriter, r *http.Request) {
	page := Html(nil,
		Head(nil,
			Meta(Props{"charset": "UTF-8"}),
			Meta(Props{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
			Title(nil, "Try-Catch-Finally 示例"),
			Link(Props{
				"href": "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css",
				"rel":  "stylesheet",
			}),
		),
		Body(nil,
			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4 text-center"}, "Try-Catch-Finally 示例"),
				P(Props{"class": "lead text-center text-muted mb-5"},
					"展示新的流暢 API：Try 生成純粹的 try-catch-finally，Do/AsyncDo 創建立即執行函數",
				),

				// 示例 1：同步 Try-Catch
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-primary text-white"},
						H5(Props{"class": "mb-0"}, "1. 同步 Try-Catch"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "純粹的 try-catch 語句，不包裝在函數中"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.Try(
    js.Const("x", "parseInt('abc')"),
    js.Log("x"),
).Catch(
    js.Log("'錯誤: ' + error.message"),
).End()`),
						),
						Button(Props{
							"class": "btn btn-primary",
							"onClick": js.Do(nil,
								js.Try(
									js.Const("x", "Math.random()"),
									JSAction{Code: "if (x < 0.5) throw new Error('數字太小')"},
									js.Alert("'成功: ' + x"),
								).Catch(
									js.Alert("'錯誤: ' + error.message"),
								).End(),
							),
						}, "測試同步 Try-Catch"),
					),
				),

				// 示例 2：AsyncFn 中的 Try-Catch-Finally
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-success text-white"},
						H5(Props{"class": "mb-0"}, "2. AsyncFn 中的 Try-Catch-Finally"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "在 AsyncFn 中使用 Try，支持 await"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).Finally(
        js.Log("'清理完成'"),
    ),
)`),
						),
						Button(Props{
							"class": "btn btn-success mb-3",
							"onClick": js.AsyncDo(nil,
								js.Const("status", "document.getElementById('status2')"),
								js.Try(
									JSAction{Code: "status.innerHTML = '⏳ 處理中...'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 1000))"},
									js.Const("random", "Math.random()"),
									JSAction{Code: "if (random < 0.5) throw new Error('隨機失敗')"},
									JSAction{Code: "status.innerHTML = '✅ 成功！'"},
								).Catch(
									JSAction{Code: "status.innerHTML = '❌ 失敗: ' + error.message"},
								).Finally(
									js.Log("'清理完成'"),
									JSAction{Code: "setTimeout(() => status.innerHTML = '等待操作...', 2000)"},
								),
							),
						}, "測試 AsyncFn + Try"),
						Div(Props{
							"id":    "status2",
							"class": "alert alert-info",
						}, "等待操作..."),
					),
				),

				// 示例 3：AsyncDo - 立即執行的異步函數
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-warning text-dark"},
						H5(Props{"class": "mb-0"}, "3. AsyncDo - 立即執行異步函數"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "使用 AsyncDo 創建立即執行的 async IIFE"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.AsyncDo(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).End(),
)`),
						),
						Button(Props{
							"class": "btn btn-warning",
							"onClick": js.AsyncDo(nil,
								js.Const("status", "document.getElementById('status3')"),
								js.Try(
									JSAction{Code: "status.innerHTML = '📂 載入中...'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 800))"},
									js.Const("data", "{message: 'Hello from AsyncDo!'}"),
									JSAction{Code: "status.innerHTML = '✅ ' + data.message"},
								).Catch(
									JSAction{Code: "status.innerHTML = '❌ 錯誤: ' + error.message"},
								).Finally(
									JSAction{Code: "setTimeout(() => status.innerHTML = '等待操作...', 2000)"},
								),
							),
						}, "測試 AsyncDo"),
						Div(Props{
							"id":    "status3",
							"class": "alert alert-warning",
						}, "等待操作..."),
					),
				),

				// 示例 4：Do - 立即執行的普通函數
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-info text-white"},
						H5(Props{"class": "mb-0"}, "4. Do - 立即執行普通函數"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "使用 Do 創建獨立作用域"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.Do(nil,
    js.Const("x", "1"),
    js.Const("y", "2"),
    js.Log("'x + y = ' + (x + y)"),
)`),
						),
						Button(Props{
							"class": "btn btn-info mb-3",
							"onClick": js.Do(nil,
								js.Const("timestamp", "Date.now()"),
								js.Const("message", "'點擊時間: ' + new Date(timestamp).toLocaleTimeString()"),
								js.Alert("message"),
							),
						}, "測試 Do"),
						P(Props{"class": "text-muted mt-3"}, "Do 創建立即執行函數，適合需要獨立作用域的場景"),
					),
				),

				// 示例 5：API 請求與錯誤處理
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-danger text-white"},
						H5(Props{"class": "mb-0"}, "5. API 請求完整示例"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "在 AsyncFn 中使用 Try 處理 API 請求"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.AsyncFn(nil,
    js.Const("container", "document.getElementById('result')"),
    js.Try(
        JSAction{Code: "container.innerHTML = '載入中...'"},
        js.Const("response", "await fetch('/api/users')"),
        js.Const("users", "await response.json()"),
        // 渲染列表
    ).Catch(
        JSAction{Code: "container.innerHTML = '錯誤'"},
    ).Finally(
        js.Log("'完成'"),
    ),
)`),
						),
						Button(Props{
							"class": "btn btn-info mb-3",
							"onClick": js.AsyncDo(nil,
								js.Const("container", "document.getElementById('apiResult')"),
								js.Try(
									JSAction{Code: "container.innerHTML = '<div class=\"spinner-border spinner-border-sm\"></div> 載入中...'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 1000))"},
									// 模擬 API 響應
									js.Const("users", "[{id: 1, name: '張三'}, {id: 2, name: '李四'}, {id: 3, name: '王五'}]"),
									JSAction{Code: "container.innerHTML = ''"},
									js.Const("ul", "document.createElement('ul')"),
									JSAction{Code: "ul.className = 'list-group'"},
									js.ForEachJS("users", "user",
										js.Const("li", "document.createElement('li')"),
										JSAction{Code: "li.className = 'list-group-item'"},
										JSAction{Code: "li.textContent = user.name"},
										JSAction{Code: "ul.appendChild(li)"},
									),
									JSAction{Code: "container.appendChild(ul)"},
								).Catch(
									JSAction{Code: "container.innerHTML = '<div class=\"alert alert-danger\">載入失敗: ' + error.message + '</div>'"},
								).Finally(
									js.Log("'API 請求完成'"),
								),
							),
						}, "載入用戶列表"),
						Div(Props{
							"id":    "apiResult",
							"class": "border p-3 rounded bg-light",
							"style": "min-height: 100px;",
						}, "點擊按鈕載入數據..."),
					),
				),

				// 示例 6：多個異步操作
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-secondary text-white"},
						H5(Props{"class": "mb-0"}, "6. 多個異步操作"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "在 AsyncFn 中處理多個異步操作"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`js.AsyncFn(nil,
    js.Try(
        js.Const("data1", "await fetch('/api/1')"),
        js.Const("data2", "await fetch('/api/2')"),
        js.Const("data3", "await fetch('/api/3')"),
    ).Catch(
        js.Alert("'請求失敗'"),
    ).End(),
)`),
						),
						Button(Props{
							"class": "btn btn-danger mb-3",
							"onClick": js.AsyncDo(nil,
								js.Const("log", "document.getElementById('log5')"),
								js.Try(
									JSAction{Code: "log.innerHTML += '<div>⏳ 開始請求 1...</div>'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 300))"},
									js.Const("data1", "{value: 100}"),
									JSAction{Code: "log.innerHTML += '<div>✅ 請求 1 完成</div>'"},

									JSAction{Code: "log.innerHTML += '<div>⏳ 開始請求 2...</div>'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 300))"},
									js.Const("data2", "{value: 200}"),
									JSAction{Code: "log.innerHTML += '<div>✅ 請求 2 完成</div>'"},

									JSAction{Code: "log.innerHTML += '<div>⏳ 開始請求 3...</div>'"},
									JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 300))"},
									// 隨機失敗
									JSAction{Code: "if (Math.random() < 0.5) throw new Error('請求 3 失敗')"},
									js.Const("data3", "{value: 300}"),
									JSAction{Code: "log.innerHTML += '<div>✅ 請求 3 完成</div>'"},

									js.Const("total", "data1.value + data2.value + data3.value"),
									JSAction{Code: "log.innerHTML += '<div class=\"text-success fw-bold\">✅ 所有請求完成，總計: ' + total + '</div>'"},
								).Catch(
									JSAction{Code: "log.innerHTML += '<div class=\"text-danger fw-bold\">❌ ' + error.message + '</div>'"},
								).Finally(
									JSAction{Code: "log.innerHTML += '<div class=\"text-muted\">--- 操作結束 ---</div>'"},
								),
							),
						}, "執行多個請求"),
						Button(Props{
							"class": "btn btn-secondary mb-3 ms-2",
							"onClick": js.Do(nil,
								JSAction{Code: "document.getElementById('log5').innerHTML = ''"},
							),
						}, "清空日誌"),
						Div(Props{
							"id":    "log5",
							"class": "border p-3 rounded bg-light",
							"style": "min-height: 150px; max-height: 300px; overflow-y: auto; font-family: monospace; font-size: 0.9rem;",
						}),
					),
				),

				// API 說明
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-dark text-white"},
						H5(Props{"class": "mb-0"}, "📚 API 說明"),
					),
					Div(Props{"class": "card-body"},
						H6(nil, "Try-Catch-Finally："),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`// 純粹的 try-catch-finally 語句（不包裝）
js.Try(
    js.Const("x", "1"),
).Catch(
    js.Log("error.message"),
).End()

// 在 AsyncFn 中使用（支持 await）
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(...).End(),
)`),
						),

						H6(Props{"class": "mt-3"}, "Do / AsyncDo："),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`// Do - 立即執行普通函數
js.Do(nil,
    js.Const("x", "1"),
    js.Log("x"),
)

// AsyncDo - 立即執行異步函數
js.AsyncDo(nil,
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)`),
						),

						Div(Props{"class": "alert alert-success mt-3"},
							Strong(nil, "設計理念："),
							Ul(nil,
								Li(nil, "✅ Try 生成純粹的 try-catch-finally，不包裝在函數中"),
								Li(nil, "✅ 需要 await 時，用 AsyncFn 或 AsyncDo 包裝"),
								Li(nil, "✅ Do/AsyncDo 專門用於創建立即執行函數（IIFE）"),
								Li(nil, "✅ 更靈活、更清晰的職責分離"),
								Li(nil, "✅ 錯誤對象統一命名為 error"),
							),
						),
					),
				),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, Render(page))
}
