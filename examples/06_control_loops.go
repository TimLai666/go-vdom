package main

import (
	"fmt"
	"log"
	"net/http"

	control "github.com/TimLai666/go-vdom/control"
	. "github.com/TimLai666/go-vdom/vdom"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 示例數據
		fruits := []string{"蘋果", "香蕉", "橘子", "葡萄", "西瓜"}
		users := []struct {
			Name string
			Age  int
			Role string
		}{
			{"Alice", 25, "開發者"},
			{"Bob", 30, "設計師"},
			{"Charlie", 35, "產品經理"},
			{"David", 28, "測試工程師"},
		}

		doc := Document(
			"Control 循環控制示例",
			[]LinkInfo{
				{Rel: "stylesheet", Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"},
			},
			nil,
			nil,

			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4"}, "Control 循環控制示例"),
				P(Props{"class": "lead"}, "展示 control.For 和 control.ForEach 的強大功能"),

				Hr(),

				// ========== control.ForEach 示例 ==========
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "🔄 control.ForEach - 遍歷集合"),
					P(Props{"class": "text-muted"}, "用於遍歷切片、數組等集合數據"),

					// 示例 1：基本遍歷
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-primary text-white"},
							H5(Props{"class": "mb-0"}, "1. 基本遍歷 - 字符串切片"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`fruits := []string{"蘋果", "香蕉", "橘子", "葡萄", "西瓜"}

Ul(control.ForEach(fruits, func(fruit string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, fruit))
}))`),
							),
							H6("渲染結果："),
							Ul(Props{"class": "list-group"},
								control.ForEach(fruits, func(fruit string, i int) VNode {
									return Li(Props{"class": "list-group-item"}, fmt.Sprintf("%d. %s", i+1, fruit))
								}),
							),
						),
					),

					// 示例 2：結構體切片遍歷
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-success text-white"},
							H5(Props{"class": "mb-0"}, "2. 結構體切片遍歷"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`control.ForEach(users, func(user User, i int) VNode {
    return Div(
        H5(user.Name),
        P(fmt.Sprintf("年齡: %d | 職位: %s", user.Age, user.Role)),
    )
})`),
							),
							H6("渲染結果："),
							Div(Props{"class": "row"},
								control.ForEach(users, func(user struct {
									Name string
									Age  int
									Role string
								}, i int) VNode {
									return Div(Props{"class": "col-md-6 mb-3"},
										Div(Props{"class": "card h-100"},
											Div(Props{"class": "card-body"},
												Div(Props{"class": "d-flex justify-content-between align-items-start mb-2"},
													H5(Props{"class": "card-title mb-0"}, user.Name),
													Span(Props{"class": "badge bg-secondary"}, fmt.Sprintf("#%d", i+1)),
												),
												P(Props{"class": "card-text mb-2"},
													Strong("年齡："), fmt.Sprintf("%d 歲", user.Age),
												),
												P(Props{"class": "card-text mb-0"},
													Strong("職位："), user.Role,
												),
											),
										),
									)
								}),
							),
						),
					),

					// 示例 3：只需要項目，不需要索引
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-info text-white"},
							H5(Props{"class": "mb-0"}, "3. 忽略索引（使用底線）"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`control.ForEach(fruits, func(fruit string, _ int) VNode {
    return Span(Props{"class": "badge"}, fruit)
})`),
							),
							H6("渲染結果："),
							Div(Props{"class": "d-flex gap-2 flex-wrap"},
								control.ForEach(fruits, func(fruit string, _ int) VNode {
									return Span(Props{"class": "badge bg-warning text-dark fs-6"}, fruit)
								}),
							),
						),
					),
				),

				Hr(),

				// ========== control.For 示例 ==========
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "➰ control.For - 傳統循環"),
					P(Props{"class": "text-muted"}, "類似傳統的 for 循環：for i := start; i < end; i += step"),

					// 示例 4：正向循環
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-primary text-white"},
							H5(Props{"class": "mb-0"}, "4. 正向循環 - 生成數字序列"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`// 語法：control.For(start, end, step, renderFunc)
// 從 1 到 10（不包含 10）
control.For(1, 11, 1, func(i int) VNode {
    return Span(Props{"class": "badge"}, fmt.Sprintf("%d", i))
})`),
							),
							H6("1 到 10："),
							Div(Props{"class": "d-flex gap-2 mb-3 flex-wrap"},
								control.For(1, 11, 1, func(i int) VNode {
									return Span(Props{"class": "badge bg-primary fs-6"}, fmt.Sprintf("%d", i))
								}),
							),

							H6("1 到 20："),
							Div(Props{"class": "d-flex gap-2 flex-wrap"},
								control.For(1, 21, 1, func(i int) VNode {
									return Span(Props{"class": "badge bg-success"}, fmt.Sprintf("%d", i))
								}),
							),
						),
					),

					// 示例 5：倒序循環
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-danger text-white"},
							H5(Props{"class": "mb-0"}, "5. 倒序循環 - 倒數計時效果"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`// 使用負數步進實現倒序
// 從 10 到 1
control.For(10, 0, -1, func(i int) VNode {
    return Span(Props{"class": "badge"}, fmt.Sprintf("%d", i))
})`),
							),
							H6("倒數 10 到 1："),
							Div(Props{"class": "d-flex gap-2 mb-3 flex-wrap"},
								control.For(10, 0, -1, func(i int) VNode {
									return Span(Props{"class": "badge bg-danger fs-5"}, fmt.Sprintf("%d", i))
								}),
							),

							H6("倒數 20 到 1："),
							Div(Props{"class": "d-flex gap-2 flex-wrap"},
								control.For(20, 0, -1, func(i int) VNode {
									return Span(Props{"class": "badge bg-warning text-dark"}, fmt.Sprintf("%d", i))
								}),
							),
						),
					),

					// 示例 6：步進循環
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-success text-white"},
							H5(Props{"class": "mb-0"}, "6. 步進循環 - 跳躍式渲染"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`// 偶數：步進 2
control.For(0, 20, 2, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})

// 5 的倍數：步進 5
control.For(0, 50, 5, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})`),
							),
							H6("偶數 0-18（步進 2）："),
							Div(Props{"class": "d-flex gap-2 mb-3 flex-wrap"},
								control.For(0, 20, 2, func(i int) VNode {
									return Span(Props{"class": "badge bg-success fs-6"}, fmt.Sprintf("%d", i))
								}),
							),

							H6("5 的倍數 0-45（步進 5）："),
							Div(Props{"class": "d-flex gap-2 mb-3 flex-wrap"},
								control.For(0, 50, 5, func(i int) VNode {
									return Span(Props{"class": "badge bg-info fs-6"}, fmt.Sprintf("%d", i))
								}),
							),

							H6("10 的倍數 0-100（步進 10）："),
							Div(Props{"class": "d-flex gap-2 flex-wrap"},
								control.For(0, 101, 10, func(i int) VNode {
									return Span(Props{"class": "badge bg-secondary"}, fmt.Sprintf("%d", i))
								}),
							),
						),
					),

					// 示例 7：實用案例 - 分頁
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-dark text-white"},
							H5(Props{"class": "mb-0"}, "7. 實用案例 - 生成分頁按鈕"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`// 生成分頁按鈕 1-10
Nav(
    Ul(Props{"class": "pagination"},
        control.For(1, 11, 1, func(i int) VNode {
            return Li(Props{"class": "page-item"},
                A(Props{"class": "page-link", "href": "#"},
                    fmt.Sprintf("%d", i)),
            )
        }),
    ),
)`),
							),
							H6("渲染結果："),
							Nav(
								Ul(Props{"class": "pagination"},
									Li(Props{"class": "page-item"},
										A(Props{"class": "page-link", "href": "#"}, "«"),
									),
									control.For(1, 11, 1, func(i int) VNode {
										class := "page-item"
										if i == 1 {
											class += " active"
										}
										return Li(Props{"class": class},
											A(Props{"class": "page-link", "href": "#"}, fmt.Sprintf("%d", i)),
										)
									}),
									Li(Props{"class": "page-item"},
										A(Props{"class": "page-link", "href": "#"}, "»"),
									),
								),
							),
						),
					),

					// 示例 8：實用案例 - 表格行
					Div(Props{"class": "card mb-3"},
						Div(Props{"class": "card-header bg-warning text-dark"},
							H5(Props{"class": "mb-0"}, "8. 實用案例 - 生成表格行"),
						),
						Div(Props{"class": "card-body"},
							Pre(Props{"class": "bg-light p-3 rounded"},
								Code(`Table(
    Tbody(
        control.For(1, 6, 1, func(i int) VNode {
            return Tr(
                Td(fmt.Sprintf("項目 %d", i)),
                Td(fmt.Sprintf("值 %d", i*10)),
            )
        }),
    ),
)`),
							),
							H6("渲染結果："),
							Table(Props{"class": "table table-striped"},
								Thead(Props{"class": "table-dark"},
									Tr(
										Th("#"),
										Th("項目"),
										Th("數值"),
										Th("平方"),
									),
								),
								Tbody(
									control.For(1, 11, 1, func(i int) VNode {
										return Tr(
											Td(fmt.Sprintf("%d", i)),
											Td(fmt.Sprintf("項目 %d", i)),
											Td(fmt.Sprintf("%d", i*10)),
											Td(fmt.Sprintf("%d", i*i)),
										)
									}),
								),
							),
						),
					),
				),

				Hr(),

				// 對比表格
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "📊 ForEach vs For 對比"),
					Div(Props{"class": "table-responsive"},
						Table(Props{"class": "table table-bordered"},
							Thead(Props{"class": "table-dark"},
								Tr(
									Th("特性"),
									Th("control.ForEach"),
									Th("control.For"),
								),
							),
							Tbody(
								Tr(
									Td(Strong("用途")),
									Td("遍歷集合（切片、數組）"),
									Td("傳統數字循環"),
								),
								Tr(
									Td(Strong("語法")),
									Td(Code("ForEach(items, func(item, i) VNode {...})")),
									Td(Code("For(start, end, step, func(i) VNode {...})")),
								),
								Tr(
									Td(Strong("參數")),
									Td("集合 + 渲染函數（接收項目和索引）"),
									Td("起始值、結束值、步進 + 渲染函數（接收索引）"),
								),
								Tr(
									Td(Strong("數據來源")),
									Td("現有的切片/數組"),
									Td("動態生成的數字序列"),
								),
								Tr(
									Td(Strong("適用場景")),
									Td("用戶列表、商品列表等數據集合"),
									Td("分頁按鈕、表格行號、倒數計時等"),
								),
								Tr(
									Td(Strong("靈活性")),
									Td("取決於數據"),
									Td("完全控制起止和步進"),
								),
								Tr(
									Td(Strong("典型示例")),
									Td("遍歷用戶列表、展示商品"),
									Td("生成 1-10 的數字、偶數序列"),
								),
							),
						),
					),
				),

				// 最佳實踐
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "✅ 最佳實踐與選擇指南"),
					Div(Props{"class": "row"},
						Div(Props{"class": "col-md-6 mb-3"},
							Div(Props{"class": "card h-100 border-success"},
								Div(Props{"class": "card-header bg-success text-white"},
									H5(Props{"class": "mb-0"}, "使用 control.ForEach"),
								),
								Div(Props{"class": "card-body"},
									Ul(Props{"class": "mb-0"},
										Li("已有數據集合需要渲染"),
										Li("需要訪問具體的項目內容"),
										Li("數據來自 API、數據庫等"),
										Li("列表項目是複雜對象（結構體）"),
										Li("需要同時使用項目和索引"),
									),
									Hr(),
									H6("示例場景："),
									Ul(Props{"class": "mb-0"},
										Li("用戶列表"),
										Li("商品目錄"),
										Li("文章列表"),
										Li("評論展示"),
									),
								),
							),
						),
						Div(Props{"class": "col-md-6 mb-3"},
							Div(Props{"class": "card h-100 border-primary"},
								Div(Props{"class": "card-header bg-primary text-white"},
									H5(Props{"class": "mb-0"}, "使用 control.For"),
								),
								Div(Props{"class": "card-body"},
									Ul(Props{"class": "mb-0"},
										Li("需要生成數字序列"),
										Li("只需要索引值，不需要項目"),
										Li("需要精確控制循環範圍"),
										Li("需要倒序或特殊步進"),
										Li("動態生成重複的 UI 元素"),
									),
									Hr(),
									H6("示例場景："),
									Ul(Props{"class": "mb-0"},
										Li("分頁按鈕（1-10）"),
										Li("倒數計時器"),
										Li("表格行號"),
										Li("評分星星（1-5）"),
									),
								),
							),
						),
					),
				),

				// 組合使用示例
				Section(Props{"class": "mb-5"},
					H2(Props{"class": "mb-3"}, "🎯 組合使用示例"),
					Div(Props{"class": "card"},
						Div(Props{"class": "card-header bg-info text-white"},
							H5(Props{"class": "mb-0"}, "組合 For 和 ForEach 創建評分系統"),
						),
						Div(Props{"class": "card-body"},
							P("為每個用戶顯示 5 星評分系統："),
							control.ForEach(users, func(user struct {
								Name string
								Age  int
								Role string
							}, i int) VNode {
								rating := (i%5 + 1) // 模擬評分 1-5
								return Div(Props{"class": "border p-3 mb-3 rounded"},
									Div(Props{"class": "d-flex justify-content-between align-items-center"},
										Div(
											H6(Props{"class": "mb-1"}, user.Name),
											P(Props{"class": "text-muted mb-2 small"}, user.Role),
										),
										Div(Props{"class": "text-end"},
											// 使用 For 生成星星
											Div(Props{"class": "mb-1"},
												control.For(1, 6, 1, func(star int) VNode {
													if star <= rating {
														return Span(Props{"class": "text-warning"}, "★")
													}
													return Span(Props{"class": "text-muted"}, "☆")
												}),
											),
											Small(Props{"class": "text-muted"}, fmt.Sprintf("%d/5", rating)),
										),
									),
								)
							}),
						),
					),
				),
			),
		)

		fmt.Fprint(w, Render(doc))
	})

	port := ":8085"
	fmt.Printf("Control 循環控制示例服務器已啟動，請訪問 http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
