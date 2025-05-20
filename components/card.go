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
//   - hoverable: 是否啟用懸停效果，預設 "true"
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
			"class":          "modern-card",
			"data-elevation": "{{elevation}}",
			"data-hoverable": "{{hoverable}}",
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
				transition: all 0.2s ease;
				overflow: hidden;
				display: flex;
				flex-direction: column;
				gap: {{contentGap}};
				border-top: 3px solid {{accentColor}};
				cursor: {{cursor}};
			`,
			"onmount": `
				(() => {
					const card = document.querySelector('.modern-card');
					if (!card) return;
					
					const isHoverable = card.dataset.hoverable !== 'false';
					const elevation = parseInt(card.dataset.elevation) || 2;
					
					if (isHoverable) {
						card.addEventListener('mouseenter', () => {
							const newElevation = Math.min(elevation + 2, 5);
							card.style.transform = 'translateY(-2px)';
							card.style.boxShadow = '0 ' + (newElevation * 2 + 4) + 'px ' + 
								(newElevation * 4 + 16) + 'px rgba(0,0,0,' + 
								(0.08 + newElevation * 0.01) + ')';
						});
						
						card.addEventListener('mouseleave', () => {
							card.style.transform = 'translateY(0)';
							card.style.boxShadow = '0 {{shadowY}}px {{shadowBlur}}px rgba(0,0,0,{{shadowOpacity}})';
						});
						
						card.addEventListener('click', () => {
							card.dispatchEvent(new CustomEvent('card:click', {
								detail: {
									title: card.querySelector('h3')?.textContent,
									content: card.querySelector('.card-content')?.textContent
								},
								bubbles: true
							}));
						});
					}
					
					// 添加淡入動畫
					card.style.opacity = '0';
					card.style.transform = 'translateY(10px)';
					
					requestAnimationFrame(() => {
						card.style.opacity = '1';
						card.style.transform = 'translateY(0)';
					});
				})();
			`,
		},
		Div(
			Props{
				"class": "card-header",
				"style": `
					display: {{titleDisplay}};
					margin: 0 0 0.5rem 0;
					padding-bottom: 0.75rem;
					border-bottom: 1px solid rgba(0,0,0,0.06);
				`,
			},
			H3(
				Props{
					"class": "card-title",
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
				"class": "card-content",
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
		// 主要屬性
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
		"hoverable":    "true",    // 是否啟用懸停效果

		// 計算屬性
		"titleDisplay":  "block",   // 標題顯示方式
		"shadowY":       "4",       // 陰影Y偏移
		"shadowBlur":    "16",      // 陰影模糊半徑
		"shadowOpacity": "0.08",    // 陰影不透明度
		"cursor":        "pointer", // 滑鼠游標樣式
	},
)
