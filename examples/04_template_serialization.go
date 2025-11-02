// examples/04_template_serialization.go
// 模板序列化示例 - 展示如何導出和導入 VNode 模板

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TimLai666/go-vdom/control"
	. "github.com/TimLai666/go-vdom/vdom"
)

func main() {
	// ==========================================
	// 示例 1: 創建一個 VNode 並轉換為 Go Template
	// ==========================================

	fmt.Println("=== 示例 1: VNode 轉 Go Template ===")

	userCard := Div(
		Props{
			"class": "card",
			"id":    "user-{{.UserID}}",
		},
		Div(
			Props{"class": "card-header"},
			H3("{{.UserName}}"),
		),
		Div(
			Props{"class": "card-body"},
			P("Email: {{.Email}}"),
			P("Role: {{.Role}}"),
		),
	)

	// 轉換為 Go Template 格式
	goTemplate := ToGoTemplate(userCard)
	fmt.Println("Go Template 輸出:")
	fmt.Println(goTemplate)
	fmt.Println()

	// ==========================================
	// 示例 2: 保存為命名模板
	// ==========================================

	fmt.Println("=== 示例 2: 保存為命名模板 ===")

	savedTemplate := SaveTemplate("user-card", userCard)
	fmt.Println(savedTemplate)

	// 保存到文件
	err := os.WriteFile("user-card.tmpl", []byte(savedTemplate), 0644)
	if err != nil {
		log.Printf("警告: 無法保存模板文件: %v\n", err)
	} else {
		fmt.Println("✓ 模板已保存到 user-card.tmpl")
	}
	fmt.Println()

	// ==========================================
	// 示例 3: JSON 序列化和反序列化
	// ==========================================

	fmt.Println("=== 示例 3: JSON 序列化 ===")

	simpleNode := Div(
		Props{
			"class":   "container",
			"id":      "main",
			"style":   "padding: 20px;",
			"visible": true,
			"count":   42,
		},
		H1("標題"),
		P("這是一個段落"),
	)

	// 轉換為 JSON
	jsonStr, err := ToJSON(simpleNode)
	if err != nil {
		log.Printf("JSON 序列化錯誤: %v\n", err)
	} else {
		fmt.Println("JSON 格式:")
		fmt.Println(jsonStr)
		fmt.Println()

		// 保存 JSON 到文件
		err = os.WriteFile("simple-node.json", []byte(jsonStr), 0644)
		if err != nil {
			log.Printf("警告: 無法保存 JSON 文件: %v\n", err)
		} else {
			fmt.Println("✓ JSON 已保存到 simple-node.json")
		}
	}
	fmt.Println()

	// 從 JSON 反序列化
	fmt.Println("=== 示例 4: JSON 反序列化 ===")
	restoredNode, err := FromJSON(jsonStr)
	if err != nil {
		log.Printf("JSON 反序列化錯誤: %v\n", err)
	} else {
		fmt.Println("從 JSON 恢復的 VNode:")
		fmt.Printf("Tag: %s\n", restoredNode.Tag)
		fmt.Printf("Props: %v\n", restoredNode.Props)
		fmt.Printf("Children count: %d\n", len(restoredNode.Children))
		fmt.Println()

		// 渲染恢復的節點
		html := Render(restoredNode)
		fmt.Println("渲染的 HTML:")
		fmt.Println(html)
	}
	fmt.Println()

	// ==========================================
	// 示例 5: 提取模板變數
	// ==========================================

	fmt.Println("=== 示例 5: 提取模板變數 ===")

	templateNode := Div(
		Props{
			"class":   "user-profile",
			"data-id": "{{.ID}}",
		},
		H1("{{.Title}}"),
		P("Welcome, {{.UserName}}!"),
		Div("Your email is: {{.Email}}"),
		Span("Role: {{.Role}}"),
	)

	vars := ExtractTemplateVars(templateNode)
	fmt.Println("找到的模板變數:")
	for _, v := range vars {
		fmt.Printf("  - %s\n", v)
	}
	fmt.Println()

	// ==========================================
	// 示例 6: 克隆 VNode
	// ==========================================

	fmt.Println("=== 示例 6: 克隆 VNode ===")

	original := Div(
		Props{"class": "original", "id": "test"},
		H1("原始標題"),
		P("原始內容"),
	)

	cloned := CloneVNode(original)
	// 修改克隆的節點
	cloned.Props["class"] = "cloned"
	cloned.Props["id"] = "test-clone"

	fmt.Println("原始節點:")
	fmt.Println(Render(original))
	fmt.Println()
	fmt.Println("克隆節點:")
	fmt.Println(Render(cloned))
	fmt.Println()

	// ==========================================
	// 示例 7: 合併 Props
	// ==========================================

	fmt.Println("=== 示例 7: 合併 Props ===")

	baseProps := Props{
		"class": "btn",
		"type":  "button",
	}

	variantProps := Props{
		"class": "btn btn-primary",
	}

	extraProps := Props{
		"id":       "submit-btn",
		"disabled": false,
	}

	merged := MergeProps(baseProps, variantProps, extraProps)
	fmt.Println("合併後的 Props:")
	for k, v := range merged {
		fmt.Printf("  %s: %v (type: %T)\n", k, v, v)
	}
	fmt.Println()

	// ==========================================
	// 示例 8: HTTP 服務器示範
	// ==========================================

	fmt.Println("=== 啟動 HTTP 服務器 ===")
	fmt.Println("訪問 http://localhost:8083 查看模板序列化示範")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 創建一個包含模板變數的文檔
		doc := Document(
			"模板序列化示例",
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
				H1(Props{"class": "text-primary mb-4"}, "模板序列化示例"),

				// 示例 1: Go Template 格式
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("1. Go Template 格式"),
					),
					Div(
						Props{"class": "card-body"},
						P("將 VNode 轉換為 Go Template 格式，可以保存到文件後再使用。"),
						Pre(
							Props{"class": "bg-light p-3"},
							goTemplate,
						),
					),
				),

				// 示例 2: JSON 格式
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("2. JSON 格式"),
					),
					Div(
						Props{"class": "card-body"},
						P("VNode 可以序列化為 JSON，方便存儲和傳輸。"),
						Pre(
							Props{"class": "bg-light p-3", "style": "max-height: 300px; overflow-y: auto;"},
							jsonStr,
						),
					),
				),

				// 示例 3: Props 類型支持
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("3. Props 類型支持"),
					),
					Div(
						Props{"class": "card-body"},
						P("Props 現在支持任意類型，自動轉換："),
						Ul(
							Li("字符串 (string)"),
							Li("整數 (int, int64, uint, etc.)"),
							Li("浮點數 (float32, float64)"),
							Li("布爾值 (bool)"),
							Li("JSAction"),
							Li("任何其他類型 (interface{})"),
						),
						Pre(
							Props{"class": "bg-light p-3"},
							`Props{
    "class":   "container",        // string
    "visible": true,                // bool
    "count":   42,                  // int
    "price":   19.99,               // float64
    "onClick": js.Do(nil,...),         // JSAction
}`,
						),
					),
				),

				// 示例 4: 提取的模板變數
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("4. 提取模板變數"),
					),
					Div(
						Props{"class": "card-body"},
						P("從模板中自動提取所有變數："),
						Ul(
							control.ForEach(vars, func(item string, idx int) VNode {
								return Li(Props{"class": "badge bg-primary me-2"}, item)
							}),
						),
					),
				),

				// 示例 5: 使用說明
				Div(
					Props{"class": "card mb-4"},
					Div(
						Props{"class": "card-header"},
						H3("5. 使用說明"),
					),
					Div(
						Props{"class": "card-body"},
						H5("保存模板:"),
						Pre(
							Props{"class": "bg-light p-3"},
							`// 保存為 Go Template
template := ToGoTemplate(vnode)
os.WriteFile("template.tmpl", []byte(template), 0644)

// 保存為 JSON
jsonStr, _ := ToJSON(vnode)
os.WriteFile("template.json", []byte(jsonStr), 0644)`,
						),
						H5(Props{"class": "mt-3"}, "載入模板:"),
						Pre(
							Props{"class": "bg-light p-3"},
							`// 從 JSON 載入
data, _ := os.ReadFile("template.json")
vnode, _ := FromJSON(string(data))

// 渲染
html := Render(vnode)`,
						),
					),
				),

				// 頁腳
				Footer(
					Props{"class": "mt-5 pt-4 border-top text-center text-muted"},
					P("© 2025 Go VDOM 模板序列化示例"),
					P(Props{"class": "small"},
						"查看控制台輸出以了解更多詳情"),
				),
			),
		)

		html := Render(doc)
		fmt.Fprint(w, html)
	})

	port := ":8083"
	log.Fatal(http.ListenAndServe(port, nil))
}
