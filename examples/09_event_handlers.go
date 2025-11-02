package main

import (
	"fmt"
	"log"
	"net/http"

	js "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/vdom"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		doc := Document(
			"事件處理器測試",
			[]LinkInfo{
				{Rel: "stylesheet", Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"},
			},
			nil,
			nil,

			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4"}, "事件處理器測試"),
				P(Props{"class": "lead"}, "測試 Do() 和 AsyncDo() 的事件處理器"),

				Hr(),

				// 1. 同步事件處理器 - 使用 Do()
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "1. 同步事件處理器 (js.Do)"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 js.Do(nil,) 創建立即執行的同步代碼塊"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Button(Props{
    "onClick": js.Do(nil,
        js.Alert("'Hello from sync handler!'"),
    ),
}, "點擊我")`),
							),
							Button(Props{
								"class": "btn btn-primary",
								"onClick": js.Do(nil,
									js.Alert("'這是一個同步事件處理器！'"),
								),
							}, "測試同步處理器"),
						),
					),
				),

				// 2. 異步事件處理器 - 使用 AsyncDo()
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "2. 異步事件處理器 (js.AsyncDo)"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 js.AsyncDo(nil,) 創建異步 IIFE，可以使用 await"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Button(Props{
    "onClick": js.AsyncDo(nil,
        js.Const("response", "await fetch('/api/data')"),
        js.Alert("'Data loaded!'"),
    ),
}, "點擊我")`),
							),
							Button(Props{
								"class": "btn btn-success",
								"onClick": js.AsyncDo(nil,
									js.Alert("'開始異步操作...'"),
									JSAction{Code: "await new Promise(r => setTimeout(r, 1000))"},
									js.Alert("'異步操作完成！'"),
								),
							}, "測試異步處理器（1秒延遲）"),
						),
					),
				),

				// 3. 複雜的異步操作 - API 模擬
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "3. 複雜異步操作 - API 載入"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 AsyncDo + Try/Catch 處理複雜的異步邏輯"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Button(Props{
    "onClick": js.AsyncDo(nil,
        js.Const("container", "document.getElementById('result')"),
        JSAction{Code: "container.innerHTML = 'Loading...'"},
        js.Try(
            JSAction{Code: "await new Promise(r => setTimeout(r, 1000))"},
            js.Const("data", "{items: ['A', 'B', 'C']}"),
            JSAction{Code: "container.innerHTML = 'Loaded: ' + data.items.join(', ')"},
        ).Catch(
            JSAction{Code: "container.innerHTML = 'Error: ' + error.message"},
        ).End(),
    ),
}, "載入數據")`),
							),
							Button(Props{
								"class": "btn btn-info mb-3",
								"onClick": js.AsyncDo(nil,
									js.Const("container", "document.getElementById('apiResult')"),
									JSAction{Code: "container.innerHTML = '<div class=\"spinner-border spinner-border-sm\" role=\"status\"></div> 載入中...'"},
									js.Try(
										JSAction{Code: "await new Promise(r => setTimeout(r, 1500))"},
										js.Const("mockData", "{items: ['項目A', '項目B', '項目C'], count: 3}"),
										JSAction{Code: "container.innerHTML = '<div class=\"alert alert-success\">成功載入 ' + mockData.count + ' 個項目：' + mockData.items.join(', ') + '</div>'"},
									).Catch(
										JSAction{Code: "container.innerHTML = '<div class=\"alert alert-danger\">載入失敗: ' + error.message + '</div>'"},
									).End(),
								),
							}, "載入 API 數據"),
							Div(Props{
								"id":    "apiResult",
								"class": "border p-3 rounded bg-light",
								"style": "min-height: 60px;",
							}, "點擊上方按鈕載入數據..."),
						),
					),
				),

				// 4. DOM 操作
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "4. DOM 元素操作"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 Do() 進行同步 DOM 操作"),
							Div(Props{"class": "mb-3"},
								Span(Props{
									"id":    "counter",
									"class": "badge bg-secondary fs-5",
								}, "0"),
							),
							Div(Props{"class": "btn-group", "role": "group"},
								Button(Props{
									"class": "btn btn-success",
									"onClick": js.Do(nil,
										js.Const("el", "document.getElementById('counter')"),
										js.Const("count", "parseInt(el.textContent) + 1"),
										JSAction{Code: "el.textContent = count"},
									),
								}, "➕ 增加"),
								Button(Props{
									"class": "btn btn-danger",
									"onClick": js.Do(nil,
										js.Const("el", "document.getElementById('counter')"),
										js.Const("count", "parseInt(el.textContent) - 1"),
										JSAction{Code: "el.textContent = count"},
									),
								}, "➖ 減少"),
								Button(Props{
									"class": "btn btn-secondary",
									"onClick": js.Do(nil,
										JSAction{Code: "document.getElementById('counter').textContent = '0'"},
									),
								}, "🔄 重置"),
							),
						),
					),
				),

				// 5. 多個事件類型
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "5. 多種事件類型"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "測試不同的事件類型：click, mouseover, mouseout, focus, blur"),
							Div(Props{
								"id":    "eventBox",
								"class": "border p-4 rounded text-center bg-light",
								"style": "min-height: 100px; cursor: pointer; user-select: none;",
								"onClick": js.Do(nil,
									js.Const("el", "document.getElementById('eventLog')"),
									JSAction{Code: "el.innerHTML += '<div class=\"badge bg-primary me-2 mb-2\">Click</div>'"},
								),
								"onMouseOver": js.Do(nil,
									JSAction{Code: "document.getElementById('eventBox').style.backgroundColor = '#ffe'"},
								),
								"onMouseOut": js.Do(nil,
									JSAction{Code: "document.getElementById('eventBox').style.backgroundColor = '#f8f9fa'"},
								),
								"onDblClick": js.Do(nil,
									js.Const("el", "document.getElementById('eventLog')"),
									JSAction{Code: "el.innerHTML += '<div class=\"badge bg-danger me-2 mb-2\">Double Click</div>'"},
								),
							}, "與這個區域互動 (點擊、滑鼠移入/移出、雙擊)"),
							Div(Props{
								"id":    "eventLog",
								"class": "mt-3 p-3 border rounded bg-white",
								"style": "min-height: 80px;",
							}, "事件日誌："),
							Button(Props{
								"class": "btn btn-sm btn-secondary mt-2",
								"onClick": js.Do(nil,
									JSAction{Code: "document.getElementById('eventLog').innerHTML = '事件日誌：'"},
								),
							}, "清除日誌"),
						),
					),
				),

				// 6. 表單事件
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "6. 表單事件處理"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "測試表單相關事件：input, change, submit"),
							Div(Props{"class": "mb-3"},
								Label(Props{"class": "form-label"}, "輸入框（即時顯示）："),
								Input(Props{
									"type":  "text",
									"class": "form-control",
									"id":    "liveInput",
									"onInput": js.Do([]string{"event"},
										js.Const("val", "event.target.value"),
										JSAction{Code: "document.getElementById('liveOutput').textContent = val"},
									),
								}),
								Div(Props{
									"id":    "liveOutput",
									"class": "mt-2 text-muted",
								}, "你輸入的內容會即時顯示在這裡..."),
							),
							Div(Props{"class": "mb-3"},
								Label(Props{"class": "form-label"}, "選擇框："),
								Select(Props{
									"class": "form-select",
									"onChange": js.Do([]string{"event"},
										js.Const("val", "event.target.value"),
										js.Alert("'你選擇了: ' + val"),
									),
								},
									Option(Props{"value": ""}, "請選擇..."),
									Option(Props{"value": "A"}, "選項 A"),
									Option(Props{"value": "B"}, "選項 B"),
									Option(Props{"value": "C"}, "選項 C"),
								),
							),
							Div(Props{"class": "mb-3"},
								Label(Props{"class": "form-label"}, "核取框："),
								Div(Props{"class": "form-check"},
									Input(Props{
										"type":  "checkbox",
										"class": "form-check-input",
										"id":    "testCheckbox",
										"onChange": js.Do([]string{"event"},
											js.Const("checked", "event.target.checked"),
											js.Const("msg", "checked ? '已勾選' : '未勾選'"),
											JSAction{Code: "document.getElementById('checkboxStatus').textContent = msg"},
										),
									}),
									Label(Props{
										"class": "form-check-label",
										"for":   "testCheckbox",
									}, "勾選我"),
								),
								Div(Props{
									"id":    "checkboxStatus",
									"class": "mt-2 text-muted",
								}, "未勾選"),
							),
						),
					),
				),

				// 7. 錯誤處理示範
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "7. 錯誤處理"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 Try/Catch 處理錯誤"),
							Button(Props{
								"class": "btn btn-warning",
								"onClick": js.AsyncDo(nil,
									js.Try(
										js.Alert("'開始可能失敗的操作...'"),
										JSAction{Code: "throw new Error('這是一個測試錯誤')"},
									).Catch(
										js.Alert("'捕獲到錯誤: ' + error.message"),
										js.Call("console.error", "error"),
									).End(),
								),
							}, "觸發錯誤"),
							Button(Props{
								"class": "btn btn-danger ms-2",
								"onClick": js.AsyncDo(nil,
									js.Try(
										JSAction{Code: "await new Promise((resolve, reject) => setTimeout(() => reject(new Error('異步錯誤')), 1000))"},
									).Catch(
										js.Alert("'捕獲到異步錯誤: ' + error.message"),
									).End(),
								),
							}, "觸發異步錯誤"),
						),
					),
				),

				// 最佳實踐總結
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "✅ 最佳實踐"),
					Div(Props{"class": "alert alert-info"},
						H5("事件處理器指南："),
						Ul(
							Li(Strong("同步操作"), "：使用 ", Code("js.Do(nil,...)"), " - 適用於簡單的 DOM 操作、alert、console.log 等"),
							Li(Strong("異步操作"), "：使用 ", Code("js.AsyncDo(nil,...)"), " - 適用於 API 調用、setTimeout、fetch 等需要 await 的操作"),
							Li(Strong("錯誤處理"), "：在 Do/AsyncDo 內部使用 ", Code("js.Try(...).Catch(...).End()"), " 處理可能的錯誤"),
							Li(Strong("不要"), " 使用 ", Code("js.Fn()"), " 或 ", Code("js.AsyncFn()"), " 作為事件處理器 - 它們只創建函數定義但不執行"),
							Li(Strong("清晰明確"), "：Do() 和 AsyncDo() 的命名清楚表明了意圖，代碼更易讀"),
						),
					),
					Div(Props{"class": "alert alert-warning"},
						H5("⚠️ 注意事項："),
						Ul(
							Li("事件處理器中的代碼會直接注入到 HTML 屬性中"),
							Li("避免在事件處理器中使用過於複雜的邏輯"),
							Li("對於複雜邏輯，考慮抽取為獨立的 JavaScript 函數"),
							Li("記得處理可能的錯誤情況，避免未捕獲的異常"),
						),
					),
				),
			),
		)

		fmt.Fprint(w, Render(doc))
	})

	port := ":8089"
	fmt.Printf("事件處理器測試服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
