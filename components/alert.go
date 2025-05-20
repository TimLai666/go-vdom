package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// Alert 現代化提示框組件
//
// 簡潔優雅的提示框，用於重要信息傳達，支持多種風格。
//
// 參數:
//   - type: 提示類型，可選 "info"、"success"、"warning"、"error"，預設 "info"
//   - title: 提示標題，預設為空
//   - bordered: 是否顯示邊框，預設 "true"
//   - closable: 是否可關閉，預設 "false"
//   - icon: 是否顯示圖標，預設 "true"
//   - rounded: 圓角程度，可選 "none"、"sm"、"md"、"lg"，預設 "md"
//   - elevation: 陰影高度，數值越高陰影越深，預設 "0"
//   - compact: 是否緊湊模式，預設 "false"
//   - customIcon: 自訂圖標HTML，預設為空
//   - customColor: 自訂主色調，預設為空 (根據type選擇)
//
// 用法:
//
//	Alert(Props{"type": "success", "title": "操作成功"}, "檔案已上傳")
//	Alert(Props{"type": "error", "closable": "true"}, "發生錯誤，請重試")
var Alert = Component(
	Div(
		Props{
			"style": `
				position: relative;
				padding: {{padding}};
				margin-bottom: 1rem;
				width: 100%;
				box-sizing: border-box;
				border-radius: {{alertRadius}};
				background-color: {{bgColor}};
				color: {{textColor}};
				border: {{border}};
				display: flex;
				align-items: {{alignItems}};
				box-shadow: {{boxShadow}};
			`,
		},
		Div(
			Props{
				"style": `
					display: {{iconDisplay}};
					flex-shrink: 0;
					margin-right: 0.75rem;
					color: {{iconColor}};
					font-size: 1.25rem;
					line-height: 1;
				`,
			},
			"{{alertIcon}}",
		),
		Div(
			Props{
				"style": "flex-grow: 1;",
			},
			Div(
				Props{
					"style": `
						display: {{titleDisplay}};
						font-weight: 600;
						margin-bottom: 0.35rem;
						font-size: 1rem;
						color: {{titleColor}};
					`,
				},
				"{{title}}",
			),
			Div(Props{}, "{{children}}"),
		),
		Button(
			Props{
				"style": `
					display: {{closeDisplay}};
					background: transparent;
					border: none;
					font-size: 1.25rem;
					line-height: 1;
					cursor: pointer;
					color: {{iconColor}};
					opacity: 0.7;
					padding: 0;
					margin-left: 0.5rem;
					flex-shrink: 0;
					font-family: sans-serif;
					font-weight: 300;
					transition: all 0.2s;
					
					&:hover {
						opacity: 1;
					}
				`,
				"aria-label": "關閉",
			},
			"×",
		),
	),
	PropsDef{
		// 主要參數
		"type":        "info",  // info, success, warning, error
		"title":       "",      // 提示標題
		"bordered":    "true",  // 是否顯示邊框
		"closable":    "false", // 是否可關閉
		"icon":        "true",  // 是否顯示圖標
		"rounded":     "md",    // 圓角：none, sm, md, lg
		"elevation":   "0",     // 陰影高度
		"compact":     "false", // 是否緊湊模式
		"customIcon":  "",      // 自訂圖標HTML
		"customColor": "",      // 自訂主色調

		// 計算屬性
		"padding":      "1rem 1.25rem",
		"alertRadius":  "0.5rem",
		"bgColor":      "#eef2ff",
		"textColor":    "#4f46e5",
		"border":       "1px solid rgba(79, 70, 229, 0.2)",
		"titleColor":   "#4338ca",
		"iconColor":    "#4f46e5",
		"boxShadow":    "none",
		"iconDisplay":  "block",
		"closeDisplay": "none",
		"titleDisplay": "block",
		"alignItems":   "flex-start",
		"alertIcon":    "&#8505;", // info icon
	},
)
