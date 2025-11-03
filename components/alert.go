package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
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
			"id": "alert-{{id}}",
			"style": `
				position: relative;
				padding: ${'{{compact}}' === 'true' ? '0.75rem 1rem' : '1rem 1.25rem'};
				margin-bottom: 1rem;
				width: 100%;
				box-sizing: border-box;
				border-radius: ${'{{rounded}}' === 'none' ? '0' : '{{rounded}}' === 'sm' ? '0.25rem' : '{{rounded}}' === 'lg' ? '0.75rem' : '0.5rem'};
				background-color: ${'{{type}}' === 'success' ? '#f0fdf4' : '{{type}}' === 'warning' ? '#fffbeb' : '{{type}}' === 'error' ? '#fef2f2' : '#eef2ff'};
				color: ${'{{type}}' === 'success' ? '#16a34a' : '{{type}}' === 'warning' ? '#d97706' : '{{type}}' === 'error' ? '#dc2626' : '#4f46e5'};
				border: ${'{{bordered}}' === 'false' ? 'none' : '{{type}}' === 'success' ? '1px solid rgba(22, 163, 74, 0.2)' : '{{type}}' === 'warning' ? '1px solid rgba(217, 119, 6, 0.2)' : '{{type}}' === 'error' ? '1px solid rgba(220, 38, 38, 0.2)' : '1px solid rgba(79, 70, 229, 0.2)'};
				display: flex;
				align-items: ${'{{compact}}' === 'true' ? 'center' : 'flex-start'};
				box-shadow: ${'{{elevation}}' === '0' ? 'none' : '{{elevation}}' === '1' ? '0 1px 2px rgba(0,0,0,0.05)' : '{{elevation}}' === '2' ? '0 1px 3px rgba(0,0,0,0.1)' : '{{elevation}}' === '3' ? '0 4px 6px rgba(0,0,0,0.1)' : '0 10px 15px rgba(0,0,0,0.1)'};
			`,
		},
		Div(
			Props{
				"style": `
					display: ${'{{icon}}' === 'false' ? 'none' : 'block'};
					flex-shrink: 0;
					margin-right: 0.75rem;
					color: ${'{{type}}' === 'success' ? '#16a34a' : '{{type}}' === 'warning' ? '#d97706' : '{{type}}' === 'error' ? '#dc2626' : '#4f46e5'};
					font-size: 1.25rem;
					line-height: 1;
				`,
			},
			"${'{{customIcon}}'.trim() ? '{{customIcon}}' : '{{type}}' === 'success' ? '&#10003;' : '{{type}}' === 'warning' ? '&#9888;' : '{{type}}' === 'error' ? '&#10005;' : '&#8505;'}",
		),
		Div(
			Props{
				"style": "flex-grow: 1;",
			},
			Div(
				Props{
					"style": `
						display: ${'{{title}}'.trim() ? 'block' : 'none'};
						font-weight: 600;
						margin-bottom: 0.35rem;
						font-size: 1rem;
						color: ${'{{type}}' === 'success' ? '#15803d' : '{{type}}' === 'warning' ? '#b45309' : '{{type}}' === 'error' ? '#b91c1c' : '#4338ca'};
					`,
				},
				"{{title}}",
			),
			Div(Props{}, "{{children}}"),
		),
		Button(
			Props{
				"id": "close-{{id}}",
				"style": `
					display: ${'{{closable}}' === 'true' ? 'block' : 'none'};
					background: transparent;
					border: none;
					font-size: 1.25rem;
					line-height: 1;
					cursor: pointer;
					color: ${'{{type}}' === 'success' ? '#16a34a' : '{{type}}' === 'warning' ? '#d97706' : '{{type}}' === 'error' ? '#dc2626' : '#4f46e5'};
					opacity: 0.7;
					padding: 0;
					margin-left: 0.5rem;
					flex-shrink: 0;
					font-family: sans-serif;
					font-weight: 300;
					transition: all 0.2s;
				`,
				"aria-label": "關閉",
			},
			"×",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
const id = '{{id}}';
const alert = document.getElementById('alert-' + id);
const closeBtn = document.getElementById('close-' + id);

if (closeBtn && alert) {
	closeBtn.addEventListener('mouseenter', function() {
		this.style.opacity = '1';
	});

	closeBtn.addEventListener('mouseleave', function() {
		this.style.opacity = '0.7';
	});

	closeBtn.addEventListener('click', function() {
		alert.style.opacity = '0';
		alert.style.transform = 'scale(0.95)';
		setTimeout(() => {
			alert.style.display = 'none';
			alert.dispatchEvent(new CustomEvent('alert:close'));
		}, 200);
	});
}

if (alert) {
	alert.style.transition = 'all 0.2s ease-out';
	alert.style.opacity = '0';
	alert.style.transform = 'scale(0.95)';

	setTimeout(() => {
		alert.style.opacity = '1';
		alert.style.transform = 'scale(1)';
	}, 50);
}
	`})),
	PropsDef{
		"id":          "",     // 提示框ID，將自動生成
		"type":        "info", // info, success, warning, error
		"title":       "",     // 提示標題
		"bordered":    true,   // 是否顯示邊框
		"closable":    false,  // 是否可關閉
		"icon":        true,   // 是否顯示圖標
		"rounded":     "md",   // 圓角：none, sm, md, lg
		"elevation":   "0",    // 陰影高度 0-4
		"compact":     false,  // 是否緊湊模式
		"customIcon":  "",     // 自訂圖標HTML
		"customColor": "",     // 自訂主色調（暫未實現）
	},
)
