package main

import (
	"fmt"
	"net/http"

	"github.com/TimLai666/go-vdom/dom"
)

// User 結構體示例
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

// Config 配置結構體
type Config struct {
	Theme      string          `json:"theme"`
	Language   string          `json:"language"`
	FontSize   int             `json:"fontSize"`
	Features   map[string]bool `json:"features"`
	Categories []string        `json:"categories"`
}

// DataCard 創建一個顯示 JSON 數據的卡片組件
func DataCard(title string, dataAttr string, data interface{}) dom.VNode {
	cardID := fmt.Sprintf("card-%s", dataAttr)

	// 創建顯示 JSON 的腳本
	displayScript := dom.JSAction{
		Code: fmt.Sprintf(`(function(){
var container = document.getElementById('%s');
var preElement = container.querySelector('pre');
var dataStr = container.dataset.%s;
if(dataStr){
var parsedData = JSON.parse(dataStr);
preElement.textContent = JSON.stringify(parsedData, null, 2);
}else{
preElement.textContent = '無數據';
}
})()`, cardID, dataAttr),
	}

	return dom.VNode{
		Tag: "div",
		Props: dom.Props{
			"id":               cardID,
			"class":            "data-card",
			"data-" + dataAttr: data,
			"onDOMReady":       displayScript,
		},
		Children: []dom.VNode{
			{Tag: "h3", Children: []dom.VNode{{Content: title}}},
			{Tag: "pre", Props: dom.Props{"class": "json-display"}},
		},
	}
}

// InteractiveList 創建可互動的列表組件
func InteractiveList(items []string) dom.VNode {
	listID := "interactive-list"

	clickHandler := dom.JSAction{
		Code: fmt.Sprintf(`(function(event){
var container = document.getElementById('%s');
var itemsStr = container.getAttribute('data-items');
var items = JSON.parse(itemsStr);
var output = document.getElementById('list-output');
if(event.target.tagName === 'LI'){
var index = Array.from(event.target.parentElement.children).indexOf(event.target);
output.textContent = '選擇了：' + items[index];
}
})`, listID),
	}

	return dom.VNode{
		Tag: "div",
		Props: dom.Props{
			"id":         listID,
			"data-items": items,
			"class":      "interactive-container",
		},
		Children: []dom.VNode{
			{Tag: "h3", Children: []dom.VNode{{Content: "可點擊的列表（數據來自 JSON）"}}},
			{
				Tag: "ul",
				Props: dom.Props{
					"class":   "item-list",
					"onClick": clickHandler,
				},
				Children: func() []dom.VNode {
					nodes := make([]dom.VNode, len(items))
					for i, item := range items {
						nodes[i] = dom.VNode{
							Tag:      "li",
							Props:    dom.Props{"class": "list-item"},
							Children: []dom.VNode{{Content: item}},
						}
					}
					return nodes
				}(),
			},
			{
				Tag:      "div",
				Props:    dom.Props{"id": "list-output", "class": "output"},
				Children: []dom.VNode{{Content: "點擊列表項查看"}},
			},
		},
	}
}

// ConfigPanel 配置面板組件
func ConfigPanel(config Config) dom.VNode {
	panelID := "config-panel"

	initScript := dom.JSAction{
		Code: fmt.Sprintf(`(function(){
var panel = document.getElementById('%s');
var configStr = panel.getAttribute('data-config');
var config = JSON.parse(configStr);

var themeEl = document.getElementById('theme-value');
themeEl.textContent = config.theme;

var langEl = document.getElementById('lang-value');
langEl.textContent = config.language;

var fontEl = document.getElementById('font-value');
fontEl.textContent = config.fontSize;

var featuresEl = document.getElementById('features-list');
var featuresHTML = Object.entries(config.features).map(function(entry){
  return '<li>' + entry[0] + ': ' + (entry[1] ? '✓' : '✗') + '</li>';
}).join('');
featuresEl.innerHTML = featuresHTML;

var categoriesEl = document.getElementById('categories-list');
var categoriesHTML = config.categories.map(function(cat){
  return '<li>' + cat + '</li>';
}).join('');
categoriesEl.innerHTML = categoriesHTML;
})()`, panelID),
	}

	return dom.VNode{
		Tag: "div",
		Props: dom.Props{
			"id":          panelID,
			"data-config": config,
			"class":       "config-panel",
			"onDOMReady":  initScript,
		},
		Children: []dom.VNode{
			{Tag: "h3", Children: []dom.VNode{{Content: "配置面板（從 JSON 讀取）"}}},
			{
				Tag:   "div",
				Props: dom.Props{"class": "config-row"},
				Children: []dom.VNode{
					{Tag: "strong", Children: []dom.VNode{{Content: "主題："}}},
					{Tag: "span", Props: dom.Props{"id": "theme-value"}},
				},
			},
			{
				Tag:   "div",
				Props: dom.Props{"class": "config-row"},
				Children: []dom.VNode{
					{Tag: "strong", Children: []dom.VNode{{Content: "語言："}}},
					{Tag: "span", Props: dom.Props{"id": "lang-value"}},
				},
			},
			{
				Tag:   "div",
				Props: dom.Props{"class": "config-row"},
				Children: []dom.VNode{
					{Tag: "strong", Children: []dom.VNode{{Content: "字體大小："}}},
					{Tag: "span", Props: dom.Props{"id": "font-value"}},
				},
			},
			{
				Tag:   "div",
				Props: dom.Props{"class": "config-section"},
				Children: []dom.VNode{
					{Tag: "h4", Children: []dom.VNode{{Content: "功能開關："}}},
					{Tag: "ul", Props: dom.Props{"id": "features-list", "class": "feature-list"}},
				},
			},
			{
				Tag:   "div",
				Props: dom.Props{"class": "config-section"},
				Children: []dom.VNode{
					{Tag: "h4", Children: []dom.VNode{{Content: "分類："}}},
					{Tag: "ul", Props: dom.Props{"id": "categories-list", "class": "category-list"}},
				},
			},
		},
	}
}

func main() {
	// 示例數據
	users := []User{
		{ID: 1, Name: "張三", Email: "zhang@example.com", Role: "管理員", Active: true},
		{ID: 2, Name: "李四", Email: "li@example.com", Role: "用戶", Active: true},
		{ID: 3, Name: "王五", Email: "wang@example.com", Role: "用戶", Active: false},
	}

	config := Config{
		Theme:    "dark",
		Language: "zh-TW",
		FontSize: 16,
		Features: map[string]bool{
			"notifications": true,
			"darkMode":      true,
			"autoSave":      false,
			"analytics":     true,
		},
		Categories: []string{"電子產品", "書籍", "服飾", "食品"},
	}

	products := []string{
		"筆記型電腦",
		"智慧手機",
		"平板電腦",
		"無線耳機",
		"智能手錶",
	}

	// 創建頁面
	page := dom.VNode{
		Tag: "html",
		Children: []dom.VNode{
			{
				Tag: "head",
				Children: []dom.VNode{
					{Tag: "meta", Props: dom.Props{"charset": "UTF-8"}},
					{Tag: "meta", Props: dom.Props{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}},
					{Tag: "title", Children: []dom.VNode{{Content: "複雜 Props JSON 序列化示例"}}},
					{
						Tag: "style",
						Children: []dom.VNode{{Content: `
							* { margin: 0; padding: 0; box-sizing: border-box; }
							body {
								font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
								background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
								padding: 20px;
								min-height: 100vh;
							}
							.container {
								max-width: 1200px;
								margin: 0 auto;
							}
							h1 {
								color: white;
								text-align: center;
								margin-bottom: 30px;
								text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
							}
							.description {
								background: rgba(255,255,255,0.95);
								padding: 20px;
								border-radius: 10px;
								margin-bottom: 30px;
								box-shadow: 0 4px 6px rgba(0,0,0,0.1);
							}
							.description h2 {
								color: #667eea;
								margin-bottom: 10px;
							}
							.description ul {
								margin-left: 20px;
								line-height: 1.8;
							}
							.grid {
								display: grid;
								grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
								gap: 20px;
								margin-bottom: 20px;
							}
							.data-card, .interactive-container, .config-panel {
								background: white;
								padding: 20px;
								border-radius: 10px;
								box-shadow: 0 4px 6px rgba(0,0,0,0.1);
							}
							h3 {
								color: #667eea;
								margin-bottom: 15px;
								font-size: 1.2em;
							}
							h4 {
								color: #764ba2;
								margin: 10px 0;
								font-size: 1em;
							}
							pre {
								background: #f5f5f5;
								padding: 15px;
								border-radius: 5px;
								overflow-x: auto;
								font-size: 0.9em;
								border: 1px solid #ddd;
							}
							.item-list {
								list-style: none;
								margin: 15px 0;
							}
							.list-item {
								padding: 10px 15px;
								background: #f8f9fa;
								margin: 5px 0;
								border-radius: 5px;
								cursor: pointer;
								transition: all 0.2s;
								border-left: 3px solid #667eea;
							}
							.list-item:hover {
								background: #667eea;
								color: white;
								transform: translateX(5px);
							}
							.output {
								padding: 15px;
								background: #e8f4f8;
								border-radius: 5px;
								margin-top: 15px;
								border: 2px solid #667eea;
								font-weight: bold;
								color: #667eea;
							}
							.config-row {
								padding: 8px 0;
								border-bottom: 1px solid #eee;
							}
							.config-row strong {
								color: #667eea;
								margin-right: 10px;
							}
							.config-section {
								margin-top: 15px;
								padding-top: 15px;
								border-top: 2px solid #f0f0f0;
							}
							.feature-list, .category-list {
								list-style: none;
								margin: 10px 0;
							}
							.feature-list li, .category-list li {
								padding: 5px 10px;
								background: #f8f9fa;
								margin: 5px 0;
								border-radius: 3px;
							}
							.json-display {
								max-height: 300px;
								overflow-y: auto;
							}
						`}},
					},
				},
			},
			{
				Tag: "body",
				Children: []dom.VNode{
					{
						Tag:   "div",
						Props: dom.Props{"class": "container"},
						Children: []dom.VNode{
							{Tag: "h1", Children: []dom.VNode{{Content: "🚀 複雜 Props JSON 序列化示例"}}},
							{
								Tag:   "div",
								Props: dom.Props{"class": "description"},
								Children: []dom.VNode{
									{Tag: "h2", Children: []dom.VNode{{Content: "功能說明"}}},
									{
										Tag: "ul",
										Children: []dom.VNode{
											{Tag: "li", Children: []dom.VNode{{Content: "自動將陣列、Map、結構體等複雜類型序列化為 JSON"}}},
											{Tag: "li", Children: []dom.VNode{{Content: "在 HTML data 屬性中存儲 JSON 數據"}}},
											{Tag: "li", Children: []dom.VNode{{Content: "客戶端 JavaScript 可以輕鬆讀取和解析這些數據"}}},
											{Tag: "li", Children: []dom.VNode{{Content: "支持嵌套的複雜數據結構"}}},
										},
									},
								},
							},
							{
								Tag:   "div",
								Props: dom.Props{"class": "grid"},
								Children: []dom.VNode{
									DataCard("用戶列表（結構體陣列）", "users", users),
									DataCard("配置對象（Map + 陣列）", "config", config),
									DataCard("產品陣列（字符串陣列）", "products", products),
								},
							},
							InteractiveList(products),
							ConfigPanel(config),
						},
					},
				},
			},
		},
	}

	// 設置路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := dom.Render(page)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	// 啟動服務器
	fmt.Println("🌟 複雜 Props 示例服務器啟動於 http://localhost:8084")
	fmt.Println("📝 此示例展示了如何使用 JSON 序列化傳遞複雜數據類型")
	if err := http.ListenAndServe(":8084", nil); err != nil {
		panic(err)
	}
}
