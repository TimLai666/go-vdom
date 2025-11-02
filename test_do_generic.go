package main

import (
	"fmt"
	"net/http"

	js "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/dom"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		doc := Document(
			"Do/AsyncDo 通用性測試",
			[]LinkInfo{
				{Rel: "stylesheet", Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"},
			},
			nil,
			nil,

			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4"}, "Do/AsyncDo 通用性測試"),
				P(Props{"class": "lead"}, "展示 Do/AsyncDo 不僅可用於事件處理器，還可用於其他場景"),

				Hr(),

				// 1. 事件處理器場景（自動傳入 event）
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "1. 事件處理器（自動傳入 event）"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "參數名為 event/e/evt/ev 時，自動傳入外部的 event 對象"),

							Button(Props{
								"class": "btn btn-primary me-2",
								"onClick": js.Do([]string{"event"},
									js.Const("id", "event.target.id"),
									js.Alert("'Event handler with event: ' + id"),
								),
								"id": "btn1",
							}, "使用 'event'"),

							Button(Props{
								"class": "btn btn-success me-2",
								"onClick": js.Do([]string{"e"},
									js.Const("id", "e.target.id"),
									js.Alert("'Event handler with e: ' + id"),
								),
								"id": "btn2",
							}, "使用 'e'"),

							Button(Props{
								"class": "btn btn-info",
								"onClick": js.Do([]string{"evt"},
									js.Const("id", "evt.target.id"),
									js.Alert("'Event handler with evt: ' + id"),
								),
								"id": "btn3",
							}, "使用 'evt'"),

							Pre(Props{"class": "bg-light p-3 mt-3"},
								Code(`// 生成：((event)=>{...})(event)
// 自動傳入外部的 event 對象`),
							),
						),
					),
				),

				// 2. 通用 IIFE - 創建作用域
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "2. 通用 IIFE - 創建獨立作用域"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "參數名不是事件相關時，生成純 IIFE，不傳入任何參數"),

							Script(nil, js.Do([]string{"x", "y"},
								js.Const("x", "10"),
								js.Const("y", "20"),
								js.Const("sum", "x + y"),
								js.Log("'IIFE with params x, y: sum = ' + sum"),
							)),

							Div(Props{"class": "alert alert-success"},
								"查看控制台，應該看到：IIFE with params x, y: sum = 30",
							),

							Pre(Props{"class": "bg-light p-3"},
								Code(`js.Do([]string{"x", "y"},
    js.Const("x", "10"),
    js.Const("y", "20"),
    js.Const("sum", "x + y"),
    js.Log("'sum = ' + sum"),
)

// 生成：((x,y)=>{const x=10;const y=20;const sum=x+y;console.log('sum = '+sum)})()
// 注意：不傳入任何參數，只是定義了參數佔位符`),
							),
						),
					),
				),

				// 3. 無參數 IIFE
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "3. 無參數 IIFE"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "不需要參數時，傳入 nil"),

							Script(nil, js.Do(nil,
								js.Const("timestamp", "new Date().toLocaleTimeString()"),
								js.Log("'Page loaded at: ' + timestamp"),
							)),

							Div(Props{"class": "alert alert-info"},
								"查看控制台，應該看到頁面載入時間",
							),

							Pre(Props{"class": "bg-light p-3"},
								Code(`js.Do(nil,
    js.Const("timestamp", "new Date().toLocaleTimeString()"),
    js.Log("'Page loaded at: ' + timestamp"),
)

// 生成：(()=>{const timestamp=new Date().toLocaleTimeString();console.log('Page loaded at: '+timestamp)})()
// 無參數，無調用實參`),
							),
						),
					),
				),

				// 4. 在 Script 標籤中使用（非事件）
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "4. 在 Script 標籤中使用"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "Do/AsyncDo 可以用在任何需要 IIFE 的地方"),

							Script(nil, JSAction{Code: `
// 使用 Do 創建模塊化代碼
` + js.Do([]string{"moduleA", "moduleB"},
								js.Const("moduleA", `{name: 'Module A', version: '1.0'}`),
								js.Const("moduleB", `{name: 'Module B', version: '2.0'}`),
								js.Log("'Modules initialized:'"),
								js.Log("moduleA"),
								js.Log("moduleB"),
							).Code}),

							Div(Props{"class": "alert alert-warning"},
								"查看控制台，應該看到模塊初始化信息",
							),

							Pre(Props{"class": "bg-light p-3"},
								Code(`// 在 <script> 標籤中使用
js.Do([]string{"moduleA", "moduleB"},
    js.Const("moduleA", "{...}"),
    js.Const("moduleB", "{...}"),
    js.Log("moduleA"),
    js.Log("moduleB"),
)

// 生成：((moduleA,moduleB)=>{...})()
// 參數名不是事件相關，所以不傳入 event`),
							),
						),
					),
				),

				// 5. 異步 IIFE
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "5. 異步 IIFE (AsyncDo)"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "AsyncDo 也支持通用場景"),

							Button(Props{
								"class": "btn btn-primary mb-3",
								"onClick": js.AsyncDo([]string{"e"},
									js.Alert("'開始異步操作（使用 e 參數）'"),
									JSAction{Code: "await new Promise(r => setTimeout(r, 1000))"},
									js.Alert("'異步操作完成！'"),
								),
							}, "異步事件處理器（使用 e）"),

							Script(nil, js.AsyncDo([]string{"data"},
								js.Log("'Starting async IIFE with data param...'"),
								JSAction{Code: "await new Promise(r => setTimeout(r, 500))"},
								js.Const("data", `{loaded: true, time: new Date().toISOString()}`),
								js.Log("'Async IIFE completed:'"),
								js.Log("data"),
							)),

							Div(Props{"class": "alert alert-info"},
								"查看控制台，應該看到異步 IIFE 的執行過程",
							),

							Pre(Props{"class": "bg-light p-3"},
								Code(`// 事件處理器（參數名為 e）
js.AsyncDo([]string{"e"}, ...)
// 生成：(async(e)=>{...})(event)  ← 傳入 event

// 通用 IIFE（參數名為 data）
js.AsyncDo([]string{"data"}, ...)
// 生成：(async(data)=>{...})()  ← 不傳入參數`),
							),
						),
					),
				),

				// 6. 智能判斷說明
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "6. 智能判斷機制"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							H5("自動傳入 event 的參數名（不區分大小寫）："),
							Ul(
								Li(Code("event"), " → 完整形式"),
								Li(Code("e"), " → 簡短形式"),
								Li(Code("evt"), " → 常用縮寫"),
								Li(Code("ev"), " → 超短形式"),
							),
							Hr(),
							H5("不傳入參數的其他名稱："),
							Ul(
								Li(Code("x"), ", ", Code("y"), ", ", Code("data"), ", ", Code("result"), " → 任何其他名稱"),
								Li("這些參數只是佔位符，調用時不傳入實參"),
								Li("在 IIFE 內部需要自己賦值"),
							),
						),
					),
				),

				// 7. 生成代碼對比
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "7. 生成代碼對比"),
					Div(Props{"class": "table-responsive"},
						Table(Props{"class": "table table-bordered"},
							Thead(Props{"class": "table-dark"},
								Tr(
									Th("Go 代碼"),
									Th("生成的 JavaScript"),
									Th("是否傳入 event"),
								),
							),
							Tbody(
								Tr(
									Td(Code("js.Do(nil, ...)")),
									Td(Code("(()=>{...})()")),
									Td(Span(Props{"class": "badge bg-secondary"}, "無參數")),
								),
								Tr(
									Td(Code("js.Do([]string{\"event\"}, ...)")),
									Td(Code("((event)=>{...})(event)")),
									Td(Span(Props{"class": "badge bg-success"}, "✅ 是")),
								),
								Tr(
									Td(Code("js.Do([]string{\"e\"}, ...)")),
									Td(Code("((e)=>{...})(event)")),
									Td(Span(Props{"class": "badge bg-success"}, "✅ 是")),
								),
								Tr(
									Td(Code("js.Do([]string{\"evt\"}, ...)")),
									Td(Code("((evt)=>{...})(event)")),
									Td(Span(Props{"class": "badge bg-success"}, "✅ 是")),
								),
								Tr(
									Td(Code("js.Do([]string{\"x\"}, ...)")),
									Td(Code("((x)=>{...})()")),
									Td(Span(Props{"class": "badge bg-danger"}, "❌ 否")),
								),
								Tr(
									Td(Code("js.Do([]string{\"data\"}, ...)")),
									Td(Code("((data)=>{...})()")),
									Td(Span(Props{"class": "badge bg-danger"}, "❌ 否")),
								),
								Tr(
									Td(Code("js.AsyncDo([]string{\"e\"}, ...)")),
									Td(Code("(async(e)=>{...})(event)")),
									Td(Span(Props{"class": "badge bg-success"}, "✅ 是")),
								),
								Tr(
									Td(Code("js.AsyncDo([]string{\"result\"}, ...)")),
									Td(Code("(async(result)=>{...})()")),
									Td(Span(Props{"class": "badge bg-danger"}, "❌ 否")),
								),
							),
						),
					),
				),

				// 總結
				Section(Props{"class": "mb-5"},
					Div(Props{"class": "alert alert-success"},
						H4(Props{"class": "alert-heading"}, "✅ 總結"),
						Hr(),
						Ul(Props{"class": "mb-0"},
							Li(Strong("事件處理器"), "：參數名為 event/e/evt/ev 時，自動傳入外部的 event 對象"),
							Li(Strong("通用 IIFE"), "：其他參數名時，生成純 IIFE，不傳入任何參數"),
							Li(Strong("靈活性"), "：同一個 API 支持多種使用場景"),
							Li(Strong("智能判斷"), "：根據參數名自動決定行為"),
						),
					),
				),
			),
		)

		fmt.Fprint(w, Render(doc))
	})

	port := ":8092"
	fmt.Printf("Do/AsyncDo 通用性測試服務器已啟動，請訪問 http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
