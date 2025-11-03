package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/dom"
	js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
	http.HandleFunc("/", minifiedJSHandler)
	fmt.Println("Server running on http://localhost:8087")
	log.Fatal(http.ListenAndServe(":8087", nil))
}

func minifiedJSHandler(w http.ResponseWriter, r *http.Request) {
	page := Html(nil,
		Head(nil,
			Meta(Props{"charset": "UTF-8"}),
			Meta(Props{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
			Title(nil, "最小化 JavaScript 輸出"),
			Link(Props{
				"href": "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css",
				"rel":  "stylesheet",
			}),
			Style(nil, `
				.code-block {
					background: #f8f9fa;
					border: 1px solid #dee2e6;
					border-radius: 0.375rem;
					padding: 1rem;
					margin: 1rem 0;
					font-family: 'Courier New', monospace;
					font-size: 0.9rem;
					overflow-x: auto;
					white-space: pre-wrap;
					word-break: break-all;
				}
				.size-badge {
					display: inline-block;
					padding: 0.25rem 0.5rem;
					border-radius: 0.25rem;
					font-size: 0.875rem;
					font-weight: 600;
				}
				.size-before {
					background-color: #ffc107;
					color: #000;
				}
				.size-after {
					background-color: #28a745;
					color: #fff;
				}
			`),
		),
		Body(nil,
			Div(Props{"class": "container mt-5"},
				H1(Props{"class": "mb-4 text-center"}, "JavaScript 代碼最小化"),
				P(Props{"class": "lead text-center text-muted mb-5"},
					"展示最小化的 JavaScript 輸出，減少傳輸大小",
				),

				// 示例 1：基本函數最小化
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-primary text-white"},
						H5(Props{"class": "mb-0"}, "1. 基本函數最小化"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "函數去除換行和多餘空格"),
						Div(Props{"class": "row"},
							Div(Props{"class": "col-md-6"},
								H6(nil, "舊格式（可讀）："),
								Pre(Props{"class": "code-block"}, `(x) => {
  const result = x * 2;
  console.log(result);
  return result;
}`),
								Span(Props{"class": "size-badge size-before"}, "84 字節"),
							),
							Div(Props{"class": "col-md-6"},
								H6(nil, "新格式（最小化）："),
								Div(Props{"class": "code-block", "id": "minified1"}),
								Span(Props{"class": "size-badge size-after", "id": "size1"}),
							),
						),
						Button(Props{
							"class": "btn btn-primary mt-3",
							"onClick": js.Do(nil,
								js.Alert("'這是最小化後的代碼示例'"),
							),
						}, "測試按鈕"),
						Script(nil, JSAction{Code: fmt.Sprintf(
							"document.getElementById('minified1').textContent=%s;document.getElementById('size1').textContent=%s.length+' 字節'",
							quote(js.Fn([]string{"x"},
								js.Const("result", "x*2"),
								js.Log("result"),
								JSAction{Code: "return result"},
							).Code),
							quote(js.Fn([]string{"x"},
								js.Const("result", "x*2"),
								js.Log("result"),
								JSAction{Code: "return result"},
							).Code),
						)}),
					),
				),

				// 示例 2：Const/Let 接受 JSAction
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-success text-white"},
						H5(Props{"class": "mb-0"}, "2. Const/Let 接受 JSAction"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "Const 和 Let 現在可以接受 JSAction 作為值"),
						Pre(Props{"class": "bg-light p-3 rounded"},
							Code(`// 傳入字符串
js.Const("x", "1")

// 傳入 JSAction
js.Const("result", js.Call("Math.random"))
js.Const("doubled", JSAction{Code: "x * 2"})

// 組合使用
js.Const("data", js.Ident("response.data"))
js.Let("count", js.Call("items.length"))`),
						),
						Button(Props{
							"class": "btn btn-success mb-3",
							"onClick": js.AsyncDo(nil,
								// 使用 JSAction 作為值
								js.Const("random", JSAction{Code: "Math.random()"}),
								js.Const("doubled", JSAction{Code: "random*2"}),
								js.Const("message", JSAction{Code: "'隨機數: '+random+', 加倍: '+doubled"}),
								js.Alert("message"),
							),
						}, "測試 JSAction 值"),
						Div(Props{"class": "code-block", "id": "minified2"}),
						Span(Props{"class": "size-badge size-after", "id": "size2"}),
						Script(nil, JSAction{Code: fmt.Sprintf(
							"document.getElementById('minified2').textContent=%s;document.getElementById('size2').textContent=%s.length+' 字節'",
							quote(js.AsyncFn(nil,
								js.Const("random", JSAction{Code: "Math.random()"}),
								js.Const("doubled", JSAction{Code: "random*2"}),
								js.Const("message", JSAction{Code: "'隨機數: '+random+', 加倍: '+doubled"}),
								js.Alert("message"),
							).Code),
							quote(js.AsyncFn(nil,
								js.Const("random", JSAction{Code: "Math.random()"}),
								js.Const("doubled", JSAction{Code: "random*2"}),
								js.Const("message", JSAction{Code: "'隨機數: '+random+', 加倍: '+doubled"}),
								js.Alert("message"),
							).Code),
						)}),
					),
				),

				// 示例 3：Try-Catch-Finally 最小化
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-info text-white"},
						H5(Props{"class": "mb-0"}, "3. Try-Catch-Finally 最小化"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "錯誤處理代碼也被最小化"),
						Div(Props{"class": "row"},
							Div(Props{"class": "col-md-6"},
								H6(nil, "舊格式："),
								Pre(Props{"class": "code-block"}, `try {
  const data = await fetch('/api');
  console.log(data);
} catch (error) {
  console.error(error);
} finally {
  console.log('done');
}`),
								Span(Props{"class": "size-badge size-before"}, "~150 字節"),
							),
							Div(Props{"class": "col-md-6"},
								H6(nil, "新格式："),
								Div(Props{"class": "code-block", "id": "minified3"}),
								Span(Props{"class": "size-badge size-after", "id": "size3"}),
							),
						),
						Button(Props{
							"class": "btn btn-info mt-3",
							"onClick": js.AsyncDo(nil,
								js.Try(
									js.Const("data", "await fetch('/api')"),
									js.Log("data"),
								).Catch(
									JSAction{Code: "console.error(error)"},
								).Finally(
									js.Log("'done'"),
								),
							),
						}, "測試 Try-Catch-Finally"),
						Script(nil, JSAction{Code: fmt.Sprintf(
							"document.getElementById('minified3').textContent=%s;document.getElementById('size3').textContent=%s.length+' 字節'",
							quote(js.Try(
								js.Const("data", "await fetch('/api')"),
								js.Log("data"),
							).Catch(
								JSAction{Code: "console.error(error)"},
							).Finally(
								js.Log("'done'"),
							).Code),
							quote(js.Try(
								js.Const("data", "await fetch('/api')"),
								js.Log("data"),
							).Catch(
								JSAction{Code: "console.error(error)"},
							).Finally(
								js.Log("'done'"),
							).Code),
						)}),
					),
				),

				// 示例 4：AsyncDo 最小化
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-warning text-dark"},
						H5(Props{"class": "mb-0"}, "4. AsyncDo 最小化"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "立即執行函數也被最小化"),
						Div(Props{"class": "row"},
							Div(Props{"class": "col-md-6"},
								H6(nil, "舊格式："),
								Pre(Props{"class": "code-block"}, `(async () => {
  const data = await fetch('/api');
  console.log(data);
})()`),
								Span(Props{"class": "size-badge size-before"}, "~80 字節"),
							),
							Div(Props{"class": "col-md-6"},
								H6(nil, "新格式："),
								Div(Props{"class": "code-block", "id": "minified4"}),
								Span(Props{"class": "size-badge size-after", "id": "size4"}),
							),
						),
						Script(nil, JSAction{Code: fmt.Sprintf(
							"document.getElementById('minified4').textContent=%s;document.getElementById('size4').textContent=%s.length+' 字節'",
							quote(js.AsyncDo(nil,
								js.Const("data", "await fetch('/api')"),
								js.Log("data"),
							).Code),
							quote(js.AsyncDo(nil,
								js.Const("data", "await fetch('/api')"),
								js.Log("data"),
							).Code),
						)}),
					),
				),

				// 示例 5：完整示例對比
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-danger text-white"},
						H5(Props{"class": "mb-0"}, "5. 完整示例對比"),
					),
					Div(Props{"class": "card-body"},
						P(Props{"class": "text-muted"}, "實際應用中的大小對比"),
						Button(Props{
							"class": "btn btn-danger mb-3",
							"onClick": js.AsyncDo(nil,
								js.Const("container", "document.getElementById('result5')"),
								js.Try(
									JSAction{Code: "container.innerHTML='載入中...'"},
									js.Const("response", "await fetch('https://jsonplaceholder.typicode.com/users/1')"),
									js.Const("user", "await response.json()"),
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-success\">載入成功: '+user.name+'</div>'"},
								).Catch(
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-danger\">錯誤: '+error.message+'</div>'"},
								).Finally(
									js.Log("'請求完成'"),
								),
							),
						}, "執行完整示例"),
						Div(Props{"id": "result5", "class": "border p-3 rounded bg-light mb-3"}),
						Div(Props{"class": "row"},
							Div(Props{"class": "col-12"},
								H6(nil, "生成的代碼："),
								Div(Props{"class": "code-block", "id": "minified5"}),
								Span(Props{"class": "size-badge size-after", "id": "size5"}),
							),
						),
						Script(nil, JSAction{Code: fmt.Sprintf(
							"document.getElementById('minified5').textContent=%s;document.getElementById('size5').textContent=%s.length+' 字節'",
							quote(js.AsyncFn(nil,
								js.Const("container", "document.getElementById('result5')"),
								js.Try(
									JSAction{Code: "container.innerHTML='載入中...'"},
									js.Const("response", "await fetch('https://jsonplaceholder.typicode.com/users/1')"),
									js.Const("user", "await response.json()"),
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-success\">載入成功: '+user.name+'</div>'"},
								).Catch(
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-danger\">錯誤: '+error.message+'</div>'"},
								).Finally(
									js.Log("'請求完成'"),
								),
							).Code),
							quote(js.AsyncFn(nil,
								js.Const("container", "document.getElementById('result5')"),
								js.Try(
									JSAction{Code: "container.innerHTML='載入中...'"},
									js.Const("response", "await fetch('https://jsonplaceholder.typicode.com/users/1')"),
									js.Const("user", "await response.json()"),
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-success\">載入成功: '+user.name+'</div>'"},
								).Catch(
									JSAction{Code: "container.innerHTML='<div class=\"alert alert-danger\">錯誤: '+error.message+'</div>'"},
								).Finally(
									js.Log("'請求完成'"),
								),
							).Code),
						)}),
					),
				),

				// 優勢總結
				Div(Props{"class": "card mb-4"},
					Div(Props{"class": "card-header bg-dark text-white"},
						H5(Props{"class": "mb-0"}, "📊 優化總結"),
					),
					Div(Props{"class": "card-body"},
						H6(nil, "最小化優勢："),
						Ul(nil,
							Li(nil, "✅ 減少傳輸大小（約 30-50%）"),
							Li(nil, "✅ 加快頁面載入速度"),
							Li(nil, "✅ 降低帶寬消耗"),
							Li(nil, "✅ 不影響功能，只去除空白"),
						),
						H6(Props{"class": "mt-3"}, "JSAction 支持優勢："),
						Ul(nil,
							Li(nil, "✅ Const/Let 可以接受 JSAction 參數"),
							Li(nil, "✅ 更靈活的代碼組合"),
							Li(nil, "✅ 可以直接傳入函數調用結果"),
							Li(nil, "✅ 減少字符串拼接"),
						),
						Div(Props{"class": "alert alert-info mt-3"},
							Strong(nil, "注意："),
							" 最小化的代碼犧牲了可讀性，但對於生產環境是最佳選擇。開發時可以使用工具格式化查看。",
						),
					),
				),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, Render(page))
}

// quote 幫助函數：將字符串轉為 JavaScript 字符串字面量
func quote(s string) string {
	return fmt.Sprintf("'%s'", escapeJS(s))
}

// escapeJS 轉義 JavaScript 字符串
func escapeJS(s string) string {
	s = s
	// 替換反斜杠
	s = replaceAll(s, `\`, `\\`)
	// 替換單引號
	s = replaceAll(s, `'`, `\'`)
	// 替換換行
	s = replaceAll(s, "\n", `\n`)
	// 替換回車
	s = replaceAll(s, "\r", `\r`)
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		found := true
		if i+len(old) <= len(s) {
			for j := 0; j < len(old); j++ {
				if s[i+j] != old[j] {
					found = false
					break
				}
			}
		} else {
			found = false
		}

		if found {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
