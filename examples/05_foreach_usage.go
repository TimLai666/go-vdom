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

		// 示例數據
		fruits := []string{"蘋果", "香蕉", "橘子", "葡萄", "西瓜"}
		numbers := []int{1, 2, 3, 4, 5}
		users := []struct {
			Name string
			Age  int
		}{
			{"Alice", 25},
			{"Bob", 30},
			{"Charlie", 35},
		}

		doc := Document(
			"ForEach 使用示例",
			[]LinkInfo{
				{Rel: "stylesheet", Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"},
			},
			nil,
			nil,

			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4"}, "ForEach 使用示例"),
				P(Props{"class": "lead"}, "展示後端和前端的列表渲染方法"),

				Hr(),

				// ========== 後端渲染示例 ==========
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "🔧 後端渲染（Go）"),

					// 示例 1：基本 ForEach
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-primary text-white"},
							H5(Props{"class": "mb-0"}, "1. ForEach - 基本用法"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "最簡潔的列表渲染方式"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Ul(ForEach(fruits, func(fruit string) VNode {
    return Li(fruit)
}))`),
							),
							H6("渲染結果："),
							Ul(Props{"class": "list-group"},
								// ✅ 使用 ForEach - 簡潔！
								ForEach(fruits, func(fruit string) VNode {
									return Li(Props{"class": "list-group-item"}, fruit)
								}),
							),
						),
					),

					// 示例 2：ForEachWithIndex
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-success text-white"},
							H5(Props{"class": "mb-0"}, "2. ForEachWithIndex - 帶索引"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "需要索引時使用"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Ul(ForEachWithIndex(fruits, func(fruit string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, fruit))
}))`),
							),
							H6("渲染結果："),
							Ol(Props{"class": "list-group list-group-numbered"},
								ForEachWithIndex(fruits, func(fruit string, i int) VNode {
									return Li(Props{"class": "list-group-item"}, fmt.Sprintf("%s (索引: %d)", fruit, i))
								}),
							),
						),
					),

					// 示例 3：複雜對象渲染
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-info text-white"},
							H5(Props{"class": "mb-0"}, "3. 複雜對象渲染"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "渲染結構體切片"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Div(ForEach(users, func(user User) VNode {
    return Div(
        H5(user.Name),
        P(fmt.Sprintf("年齡: %d", user.Age)),
    )
}))`),
							),
							H6("渲染結果："),
							Div(Props{"class": "row"},
								ForEach(users, func(user struct {
									Name string
									Age  int
								}) VNode {
									return Div(Props{"class": "col-md-4"},
										Div(Props{"class": "card"},
											Div(Props{"class": "card-body"},
												H5(Props{"class": "card-title"}, user.Name),
												P(Props{"class": "card-text"}, fmt.Sprintf("年齡: %d 歲", user.Age)),
												Span(Props{"class": "badge bg-primary"}, "用戶"),
											),
										),
									)
								}),
							),
						),
					),

					// 示例 4：數字序列
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-warning text-dark"},
							H5(Props{"class": "mb-0"}, "4. 數字序列渲染"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "渲染數字切片"),
							H6("渲染結果："),
							Div(Props{"class": "d-flex gap-2"},
								ForEach(numbers, func(num int) VNode {
									return Span(Props{"class": "badge bg-secondary fs-5"}, fmt.Sprintf("%d", num))
								}),
							),
						),
					),
				),

				Hr(),

				// ========== 前端渲染示例 ==========
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "⚡ 前端渲染（JavaScript）"),

					// 示例 5：前端 ForEach
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-dark text-white"},
							H5(Props{"class": "mb-0"}, "5. js.ForEach - 遍歷數組"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "點擊按鈕在前端遍歷數組並輸出到控制台"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`js.ForEachJS("['A', 'B', 'C']", "item",
    js.Log("'項目: ' + item"),
)`),
							),
							Button(Props{
								"class": "btn btn-primary mb-2",
								"onClick": js.Do(nil,
									js.Const("arr", "['A', 'B', 'C']"),
									js.ForEachJS("arr", "item",
										js.Log("'項目: ' + item"),
									),
									js.Alert("'查看控制台輸出！'"),
								),
							}, "執行 ForEach"),
							Div(Props{"class": "alert alert-info"},
								Strong("提示："), " 打開瀏覽器控制台 (F12) 查看輸出",
							),
						),
					),

					// 示例 6：前端 ForEachWithIndex
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-secondary text-white"},
							H5(Props{"class": "mb-0"}, "6. js.ForEachWithIndex - 帶索引遍歷"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "遍歷時同時獲取索引"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`js.Const("numbers", "[10, 20, 30, 40, 50]")
js.ForEachWithIndexJS("numbers", "num", "idx",
    js.Log("'索引 ' + idx + ': ' + num"),
)`),
							),
							Button(Props{
								"class": "btn btn-success mb-2",
								"onClick": js.Do(nil,
									js.Const("numbers", "[10, 20, 30, 40, 50]"),
									js.ForEachWithIndexJS("numbers", "num", "idx",
										js.Log("'[' + idx + '] = ' + num"),
									),
									js.Alert("'查看控制台輸出！'"),
								),
							}, "執行 ForEachWithIndex"),
							Div(Props{"class": "alert alert-info"},
								Strong("提示："), " 打開瀏覽器控制台 (F12) 查看輸出",
							),
						),
					),

					// 示例 7：動態創建 DOM 元素
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-danger text-white"},
							H5(Props{"class": "mb-0"}, "7. 動態創建 DOM 元素"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "使用 ForEach 動態創建並添加 DOM 元素"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`js.ForEachJS("colors", "color",
    js.Const("div", "document.createElement('div')"),
    JSAction{Code: "div.textContent = color"},
    JSAction{Code: "div.className = 'badge bg-' + color + ' me-2'"},
    JSAction{Code: "container.appendChild(div)"},
)`),
							),
							Button(Props{
								"class": "btn btn-danger mb-3",
								"onClick": js.Do(nil,
									js.Const("container", "document.getElementById('dynamicContainer')"),
									JSAction{Code: "container.innerHTML = ''"},
									js.Const("colors", "['primary', 'secondary', 'success', 'danger', 'warning', 'info']"),
									js.ForEachJS("colors", "color",
										js.Const("div", "document.createElement('div')"),
										JSAction{Code: "div.textContent = color"},
										JSAction{Code: "div.className = 'badge bg-' + color + ' me-2 mb-2'"},
										JSAction{Code: "container.appendChild(div)"},
									),
								),
							}, "動態生成徽章"),
							Div(Props{
								"id":    "dynamicContainer",
								"class": "border p-3 rounded bg-light",
								"style": "min-height: 60px;",
							}, "點擊上方按鈕生成內容..."),
						),
					),

					// 示例 8：前端 API 數據遍歷
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-primary text-white"},
							H5(Props{"class": "mb-0"}, "8. API 數據遍歷 (異步)"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "從 API 獲取數據並使用 ForEach 渲染"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`js.Try(
    js.Const("response", "await fetch('/api/items')"),
    js.Const("items", "await response.json()"),
    js.ForEachJS("items", "item",
        // 處理每個項目
    ),
).Catch(
    js.Call("console.error", js.Ident("error")),
).End()`),
							),
							Button(Props{
								"class": "btn btn-primary mb-3",
								"onClick": js.AsyncDo(nil,
									js.Const("container", "document.getElementById('apiContainer')"),
									JSAction{Code: "container.innerHTML = '<div class=\"spinner-border\" role=\"status\"></div> 載入中...'"},
									js.Try(
										// Try 區塊 - 模擬 API 調用
										js.Const("mockData", "[{name: '項目A', value: 100}, {name: '項目B', value: 200}, {name: '項目C', value: 300}]"),
										JSAction{Code: "await new Promise(resolve => setTimeout(resolve, 1000))"},
										JSAction{Code: "container.innerHTML = ''"},
										js.Const("ul", "document.createElement('ul')"),
										JSAction{Code: "ul.className = 'list-group'"},
										js.ForEachJS("mockData", "item",
											js.Const("li", "document.createElement('li')"),
											JSAction{Code: "li.className = 'list-group-item d-flex justify-content-between align-items-center'"},
											JSAction{Code: "li.innerHTML = item.name"},
											js.Const("badge", "document.createElement('span')"),
											JSAction{Code: "badge.className = 'badge bg-primary rounded-pill'"},
											JSAction{Code: "badge.textContent = item.value"},
											JSAction{Code: "li.appendChild(badge)"},
											JSAction{Code: "ul.appendChild(li)"},
										),
										JSAction{Code: "container.appendChild(ul)"},
									).Catch(
										// Catch 區塊（錯誤對象為 error）
										JSAction{Code: "container.innerHTML = '<div class=\"alert alert-danger\">載入失敗: ' + error.message + '</div>'"},
									).End(),
								),
							}, "從 API 載入數據"),
							Div(Props{
								"id":    "apiContainer",
								"class": "border p-3 rounded bg-light",
								"style": "min-height: 100px;",
							}, "點擊上方按鈕載入數據..."),
						),
					),

					// 示例 9：ForEachElement (DOM 元素專用)
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-success text-white"},
							H5(Props{"class": "mb-0"}, "9. js.ForEachElement - DOM 元素操作"),
						),
						Div(Props{"class": "card-body"},
							P(Props{"class": "text-muted"}, "專門用於操作 DOM 元素列表"),
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`js.ForEachElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
    return el.AddClass("'highlighted'")
})`),
							),
							Div(Props{"class": "mb-3"},
								Span(Props{"class": "item badge bg-secondary me-2"}, "項目 1"),
								Span(Props{"class": "item badge bg-secondary me-2"}, "項目 2"),
								Span(Props{"class": "item badge bg-secondary me-2"}, "項目 3"),
								Span(Props{"class": "item badge bg-secondary me-2"}, "項目 4"),
							),
							Button(Props{
								"class": "btn btn-success",
								"onClick": js.Do(nil,
									js.ForEachElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
										return JSAction{Code: el.Ref() + ".classList.toggle('bg-warning');" + el.Ref() + ".classList.toggle('bg-secondary')"}
									}),
								),
							}, "切換顏色"),
						),
					),
				),

				Hr(),

				// 對比總結
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "📊 使用對比"),
					Div(Props{"class": "table-responsive"},
						Table(Props{"class": "table table-bordered"},
							Thead(Props{"class": "table-dark"},
								Tr(
									Th("場景"),
									Th("後端渲染（Go）"),
									Th("前端渲染（JavaScript）"),
								),
							),
							Tbody(
								Tr(
									Td(Strong("基本列表")),
									Td(Code("ForEach(items, func(item) VNode {...})")),
									Td(Code("js.ForEachJS('array', 'item', ...actions)")),
								),
								Tr(
									Td(Strong("帶索引")),
									Td(Code("ForEachWithIndex(items, func(item, i) VNode {...})")),
									Td(Code("js.ForEachWithIndexJS('array', 'item', 'i', ...actions)")),
								),
								Tr(
									Td(Strong("DOM 操作")),
									Td("N/A（不適用）"),
									Td(Code("js.ForEachElement(selector, func(el Elem) JSAction {...})")),
								),
								Tr(
									Td(Strong("異步數據")),
									Td("N/A（後端處理）"),
									Td(Code("js.AsyncFn + js.ForEach")),
								),
								Tr(
									Td(Strong("性能")),
									Td(Span(Props{"class": "badge bg-success"}, "快（伺服器端生成）")),
									Td(Span(Props{"class": "badge bg-warning"}, "依賴客戶端")),
								),
								Tr(
									Td(Strong("SEO")),
									Td(Span(Props{"class": "badge bg-success"}, "友好")),
									Td(Span(Props{"class": "badge bg-danger"}, "不友好")),
								),
								Tr(
									Td(Strong("動態性")),
									Td(Span(Props{"class": "badge bg-warning"}, "靜態")),
									Td(Span(Props{"class": "badge bg-success"}, "動態")),
								),
							),
						),
					),
				),

				// 最佳實踐
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "✅ 最佳實踐"),
					Div(Props{"class": "row"},
						Div(Props{"class": "col-md-6"},
							Div(Props{"class": "card h-100"},
								Div(Props{"class": "card-header bg-success text-white"},
									H5(Props{"class": "mb-0"}, "✅ 推薦做法"),
								),
								Div(Props{"class": "card-body"},
									Ul(
										Li("靜態內容使用後端 ForEach（SEO 友好）"),
										Li("動態內容使用前端 js.ForEach"),
										Li("大列表優先考慮後端渲染"),
										Li("實時更新的數據使用前端渲染"),
										Li("混合使用以達到最佳效果"),
									),
								),
							),
						),
						Div(Props{"class": "col-md-6"},
							Div(Props{"class": "card h-100"},
								Div(Props{"class": "card-header bg-danger text-white"},
									H5(Props{"class": "mb-0"}, "❌ 避免做法"),
								),
								Div(Props{"class": "card-body"},
									Ul(
										Li("不要在前端渲染大量靜態列表"),
										Li("不要混淆後端和前端的 ForEach"),
										Li("不要在 js.ForEach 中使用 Go 變數"),
										Li("不要忘記處理空數組情況"),
										Li("不要過度使用前端渲染影響 SEO"),
									),
								),
							),
						),
					),
				),
			),
		)

		fmt.Fprint(w, Render(doc))
	})

	port := ":8084"
	fmt.Printf("ForEach 示例服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
