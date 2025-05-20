// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	comp "github.com/TimLai666/go-vdom/components" // 使用 dot import
	. "github.com/TimLai666/go-vdom/jsdsl"         // 使用 dot import
	. "github.com/TimLai666/go-vdom/vdom"          // 使用 dot import

	// 注意：不要用 dot import 匯入 control，請用 control.xxx
	control "github.com/TimLai666/go-vdom/control"
)

// 定義一個簡單的數據結構用於 API 響應
type ApiData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	// 定義一個卡片組件
	Card := Component(
		Div(
			Props{"class": "card"},
			H1("{{title}}"),
			Div("{{children}}"),
		),
		PropsDef{"title": ""}, // 預設 props
	)

	// 測試 control: If/Then/Else/Repeat
	show := false
	items := []string{"蘋果", "香蕉", "橘子"}

	// 處理API請求的函數
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// 設置內容類型為JSON
		w.Header().Set("Content-Type", "application/json")

		// 創建一些測試數據
		data := []ApiData{
			{Id: 1, Name: "項目一", Message: "這是從API獲取的第一條消息"},
			{Id: 2, Name: "項目二", Message: "這是從API獲取的第二條消息"},
			{Id: 3, Name: "項目三", Message: "這是從API獲取的第三條消息"},
		}

		// 編碼為JSON並發送
		json.NewEncoder(w).Encode(data)
	})

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
				// 添加 Fetch API 示例區塊 - GET 請求
				Div(
					Props{"class": "mt-4"},
					H4("Fetch GET API 測試"),
					P("點擊下方按鈕從 API 獲取數據："),
					Button(Props{
						"id":    "fetchButton",
						"class": "btn btn-primary mb-3",
					}, "獲取數據"),
					Div(Props{
						"id":    "dataContainer",
						"class": "border p-3 bg-light",
					}, "數據將顯示在這裡..."),
				),
				// 添加 Fetch POST API 示例區塊
				Div(
					Props{"class": "mt-4"},
					H4("Fetch POST API 測試"),
					P("填寫並提交表單以發送 POST 請求："),
					Form(Props{
						"id":     "postForm",
						"class":  "mb-3",
						"action": "#",
					},
						Div(Props{"class": "mb-3"},
							Label(Props{"for": "nameInput"}, "姓名"),
							Input(Props{
								"id":          "nameInput",
								"type":        "text",
								"class":       "form-control",
								"placeholder": "請輸入姓名",
							}),
						),
						Div(Props{"class": "mb-3"},
							Label(Props{"for": "messageInput"}, "訊息"),
							Textarea(Props{
								"id":          "messageInput",
								"class":       "form-control",
								"rows":        "3",
								"placeholder": "請輸入訊息",
							}),
						),
						Button(Props{
							"id":    "submitButton",
							"type":  "submit",
							"class": "btn btn-success",
						}, "提交表單"),
					),
					Div(Props{
						"id":    "postResponseContainer",
						"class": "border p-3 bg-light mt-3",
					}, "提交結果將顯示在這裡..."),
				),
				// 添加組件示例區塊
				Div(
					Props{"class": "mt-5 mb-5"},
					H3("UI 組件庫展示", Props{"class": "mb-4"}),
					Div(Props{"class": "row g-4"},
						// 左側欄
						Div(Props{"class": "col-md-6"},
							// 文字輸入框
							comp.TextField(Props{
								"id":          "nameInput",
								"label":       "用戶名稱",
								"placeholder": "請輸入您的用戶名",
								"required":    "true",
								"helpText":    "用戶名應為 3-16 個字符",
							}),
							// 下拉選單
							comp.Dropdown(Props{
								"id":       "country",
								"label":    "選擇國家",
								"options":  "台灣,中國,日本,美國,韓國",
								"helpText": "請選擇您的所在國家",
								"required": "true",
							}),
							// 單選按鈕組
							comp.RadioGroup(Props{
								"id":           "sex",
								"label":        "選擇性別",
								"name":         "gender",
								"options":      "男性,女性,其他",
								"defaultValue": "男性",
								"direction":    "horizontal",
								"helpText":     "請選擇您的性別",
							}),
						),
						// 右側欄
						Div(Props{"class": "col-md-6"},
							// 勾選框
							comp.Checkbox(Props{
								"id":       "terms",
								"name":     "terms",
								"label":    "我同意服務條款和隱私政策",
								"required": "true",
								"helpText": "您必須同意條款才能繼續",
							}),
							// 勾選框組
							comp.CheckboxGroup(Props{
								"id":       "hobbies",
								"label":    "選擇愛好",
								"name":     "hobbies",
								"options":  "閱讀,運動,音樂,繪畫,旅行",
								"values":   "閱讀,音樂",
								"helpText": "可多選",
							}),
							// 開關
							comp.Switch(Props{
								"id":            "notifications",
								"name":          "notifications",
								"label":         "啟用電子郵件通知",
								"checked":       "true",
								"helpText":      "開啟以接收重要通知",
								"labelPosition": "right",
							}),
						),
					),
					// 模擬表單按鈕
					Div(Props{"class": "d-flex justify-content-center mt-4"},
						Button(Props{
							"type":  "button",
							"class": "btn btn-primary me-2",
						}, "提交表單"),
						Button(Props{
							"type":  "button",
							"class": "btn btn-outline-secondary",
						}, "取消"),
					),
				),
			),
			Footer(Props{"class": "container bg-light p-4 mt-4"},
				P(Props{"class": "text-center"},
					"© 2025 ", Span(Props{"style": "color:red;"}, "Go VDOM"), " 示範網站 | 使用 Go 和 VDOM 製作",
				),
			),
			Script(Props{"type": "module"}, DomReady(
				// 導航欄點擊事件
				QueryEach(Els(".navbar-nav .nav-link"), func(el Elem) JSAction {
					return el.OnClick(JSAction{Code: `alert('你點擊了 ' + this.innerText)`})
				}), // Fetch GET 範例
				El("#fetchButton").OnClick(
					Try(
						FetchRequest("/api/data",
							WithMethod("GET"),
							WithContentType("application/json"),
						),
						WithResponseType(JSONResponse),
						WithThen(
							// 使用 StoreResult 將結果存儲到 apiData 變數中
							StoreResult("apiData",
								// 打印存儲的數據到控制台，以便驗證
								Log("'存儲的 API 數據:' + JSON.stringify(apiData)"), // 以下是原來的數據處理邏輯，但使用存儲的 apiData 變數
								El("#dataContainer").SetHTML("''"),
								Const("ul", "document.createElement('ul')"),
								CallMethod("ul", "classList.add", "'list-group'"), // 使用儲存的 apiData 變數而不是直接使用 data
								CallMethod("apiData", "forEach", Fn(
									[]string{"item"},
									Const("li", "document.createElement('li')"),
									CallMethod("li", "classList.add", "'list-group-item'"),
									V("li").SetHTML("'<strong>' + item.name + '</strong>: ' + item.message"),
									CallMethod("ul", "appendChild", "li"),
								)),
								CallMethod("dataContainer", "appendChild", "ul"), // 添加一個示例，展示如何再次使用存儲的數據
								Log("'API 數據項目數量: ' + apiData.length"),
							)), WithCatch(
							Log("'獲取數據時出錯:' + error.message"),
							El("#dataContainer").SetHTML(
								"'<div class=\"alert alert-danger\">獲取數據時出錯: ' + error.message + '</div>'",
							),
						),
					),
				), // 表單提交監聽器
				Let("postForm", "document.querySelector('#postForm')"),
				CallMethod("postForm", "addEventListener", "'submit'", Fn(
					[]string{"event"},
					CallMethod("event", "preventDefault"),
					Const("nameValue", "document.getElementById('nameInput').value"),
					Const("messageValue", "document.getElementById('messageInput').value"), Try(
						FetchRequest("/api/data",
							WithMethod("POST"),
							WithContentType("application/json"),
							WithBody("JSON.stringify({ name: nameValue, message: messageValue })"),
						),
						WithResponseType(JSONResponse),
						WithThen(
							// 使用 StoreResult 將 POST 結果存儲到 postResponse 變數中
							StoreResult("postResponse",
								// 打印存儲的 POST 響應到控制台
								Log("'POST 請求結果:' + JSON.stringify(postResponse)"),

								// 清空表單
								Call("document.getElementById('postForm').reset"),

								// 顯示成功消息，並添加響應的細節
								El("#postResponseContainer").SetHTML(
									"'<div class=\"alert alert-success\">表單提交成功！回應包含 ' + postResponse.length + ' 個項目</div>'",
								),
							),
						), WithCatch(
							Log("'提交表單時出錯:' + error.message"),
							El("#postResponseContainer").SetHTML(
								"'<div class=\"alert alert-danger\">提交表單時出錯: ' + error.message + '</div>'",
							),
						),
					),
				)),
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
