package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// Card 現代化卡片組件
//
// 提供優雅的陰影效果、精心設計的間距和精緻的圓角。
//
// 參數:
//   - title: 卡片標題，若不需要則留空
//   - titleWeight: 標題字重，默認 "500"
//   - titleColor: 標題顏色，默認優雅灰黑 "#1a2b4a"
//   - elevation: 陰影高度，數值越高陰影越深，默認 "2"
//   - accentColor: 強調色，默認現代藍 "#3b82f6"
//   - maxWidth: 卡片最大寬度，默認 "480px"
//   - padding: 內邊距，默認 "1.75rem"
//   - borderRadius: 圓角大小，默認 "12px"
//   - background: 背景色，默認純白 "#ffffff"
//   - contentGap: 內容間距，默認 "1.25rem"
//
// 用法:
//
//	Card(Props{"title": "我的卡片", "accentColor": "#6366f1"},
//	    P("這是卡片內容"),
//	    P("更多內容..."),
//	)
var Card = Component(
	Div(
		Props{
			"style": `
				position: relative;
				border: none;
				border-radius: {{borderRadius}};
				background: {{background}};
				padding: {{padding}};
				margin-bottom: 2rem;
				width: 100%;
				max-width: {{maxWidth}};
				box-sizing: border-box;
				box-shadow: 0 {{shadowY}}px {{shadowBlur}}px rgba(0,0,0,{{shadowOpacity}});
				transition: transform 0.2s, box-shadow 0.2s;
				overflow: hidden;
				display: flex;
				flex-direction: column;
				gap: {{contentGap}};
				border-top: 3px solid {{accentColor}};
			`,
		},
		Div(
			Props{
				"style": `
					display: {{titleDisplay}};
					margin: 0 0 0.5rem 0;
					padding-bottom: 0.75rem;
					border-bottom: 1px solid rgba(0,0,0,0.06);
				`,
			},
			H3(
				Props{
					"style": `
						margin: 0;
						padding: 0;
						font-size: 1.35rem;
						font-weight: {{titleWeight}};
						color: {{titleColor}};
						letter-spacing: -0.02em;
						line-height: 1.3;
					`,
				},
				"{{title}}",
			),
		),
		Div(
			Props{
				"style": `
					display: flex;
					flex-direction: column;
					gap: 1rem;
				`,
			},
			"{{children}}",
		),
	),
	PropsDef{
		"title":        "",        // 卡片標題
		"titleWeight":  "500",     // 標題字重
		"titleColor":   "#1a2b4a", // 標題顏色
		"elevation":    "2",       // 陰影高度
		"accentColor":  "#3b82f6", // 強調色
		"maxWidth":     "480px",   // 最大寬度
		"padding":      "1.75rem", // 內邊距
		"borderRadius": "12px",    // 圓角大小
		"background":   "#ffffff", // 背景色
		"contentGap":   "1.25rem", // 內容間距
		// 以下是計算屬性
		"titleDisplay":  "block", // 標題顯示方式
		"shadowY":       "4",     // 陰影Y偏移
		"shadowBlur":    "16",    // 陰影模糊半徑
		"shadowOpacity": "0.08",  // 陰影不透明度
	},
)
