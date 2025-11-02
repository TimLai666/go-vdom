package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
)

// Modal 現代化對話框組件
//
// 提供精美的浮動對話框，適合展示重要信息或需要用戶交互的內容。
//
// 參數:
//   - title: 對話框標題，預設為空
//   - open: 是否顯示，可選 "true"、"false"，預設 "false"
//   - size: 尺寸，可選 "xs"、"sm"、"md"、"lg"、"xl"、"full"，預設 "md"
//   - closeButton: 是否顯示關閉按鈕，預設 "true"
//   - closeOnEsc: 是否按ESC鍵關閉，預設 "true"
//   - closeOnOverlayClick: 是否點擊遮罩層關閉，預設 "true"
//   - centered: 是否垂直居中，預設 "true"
//   - scrollable: 內容是否可滾動，預設 "true"
//   - animation: 動畫效果，可選 "fade"、"slide"、"zoom"，預設 "fade"
//   - footer: 底部內容，預設為空
//   - overlayColor: 遮罩層顏色，預設 "rgba(0,0,0,0.5)"
//   - radius: 圓角大小，預設 "md"
//   - elevation: 陰影高度，預設 "3"
//   - hideHeader: 是否隱藏頭部，預設 "false"
//   - hideFooter: 是否隱藏底部，預設 "false"
//   - zIndex: 層級，預設 "1050"
//
// 用法:
//
//	Modal(Props{
//	  "title": "確認操作",
//	  "open": "true",
//	  "size": "sm"
//	},
//	  P("您確定要執行此操作嗎？"),
//	  Div(Props{"slot": "footer"},
//	    Btn(Props{"variant": "outlined"}, "取消"),
//	    Btn(Props{}, "確認")
//	  ),
//	)
var Modal = Component(
	Div(
		Props{
			"style": `
				position: fixed;
				top: 0;
				left: 0;
				width: 100%;
				height: 100%;
				display: {{display}};
				z-index: {{zIndex}};
				overflow-x: hidden;
				overflow-y: {{overlayOverflow}};
				animation: {{overlayAnimation}} 0.3s ease;
				pointer-events: {{pointerEvents}};

				@keyframes modalOverlayFadeIn {
					0% { opacity: 0; }
					100% { opacity: 1; }
				}

				@keyframes modalOverlayFadeOut {
					0% { opacity: 1; }
					100% { opacity: 0; }
				}
			`,
		},
		Div(
			Props{
				"style": `
					position: fixed;
					top: 0;
					left: 0;
					width: 100%;
					height: 100%;
					background: {{overlayColor}};
					backdrop-filter: blur(2px);
				`,
			},
		),
		Div(
			Props{
				"style": `
					position: relative;
					width: {{modalWidth}};
					max-width: {{maxWidth}};
					margin: {{margin}};
					background: #ffffff;
					border-radius: {{borderRadius}};
					box-shadow: {{boxShadow}};
					display: flex;
					flex-direction: column;
					max-height: {{maxHeight}};
					animation: {{contentAnimation}} 0.3s ease;

					@keyframes modalFadeIn {
						0% { opacity: 0; transform: translate(0, -20px); }
						100% { opacity: 1; transform: translate(0, 0); }
					}

					@keyframes modalSlideIn {
						0% { opacity: 0; transform: translate(0, 40px); }
						100% { opacity: 1; transform: translate(0, 0); }
					}

					@keyframes modalZoomIn {
						0% { opacity: 0; transform: scale(0.95); }
						100% { opacity: 1; transform: scale(1); }
					}
				`,
			},
			Div(
				Props{
					"style": `
						display: {{headerDisplay}};
						padding: 1rem 1.5rem;
						border-bottom: 1px solid #f1f5f9;
						display: flex;
						align-items: center;
						justify-content: space-between;
					`,
				},
				H3(
					Props{
						"style": `
							font-size: 1.125rem;
							font-weight: 600;
							color: #0f172a;
							margin: 0;
						`,
					},
					"{{title}}",
				),
				Button(
					Props{
						"style": `
							display: {{closeButtonDisplay}};
							background: transparent;
							border: none;
							font-size: 1.5rem;
							line-height: 1;
							padding: 0.25rem;
							cursor: pointer;
							color: #64748b;
							transition: color 0.15s;
							margin: -0.5rem -0.5rem -0.5rem 0.5rem;
							border-radius: 0.25rem;
						`,
						"aria-label": "關閉",
					},
					"×",
				),
			),
			Div(
				Props{
					"style": `
						padding: 1.5rem;
						overflow-y: {{contentOverflow}};
					`,
				},
				"{{children}}",
			),
			Div(
				Props{
					"style": `
						display: {{footerDisplay}};
						padding: 1rem 1.5rem;
						border-top: 1px solid #f1f5f9;
						display: flex;
						align-items: center;
						justify-content: flex-end;
						gap: 0.75rem;
					`,
				},
				"{{footer}}",
			),
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
const modal = document.querySelector('[style*="z-index: {{zIndex}}"]');
const overlay = modal.children[0];
const content = modal.children[1];
const closeBtn = content.querySelector('button[aria-label="關閉"]');

function closeModal() {
	modal.style.display = 'none';
	modal.dispatchEvent(new CustomEvent('modal:close'));
}

function handleKeydown(event) {
	if (event.key === 'Escape' && '{{closeOnEsc}}' === 'true') {
		closeModal();
	}
}

if (closeBtn) {
	closeBtn.addEventListener('mouseenter', function() {
		this.style.color = '#334155';
		this.style.background = '#f1f5f9';
	});

	closeBtn.addEventListener('mouseleave', function() {
		this.style.color = '#64748b';
		this.style.background = 'transparent';
	});

	closeBtn.addEventListener('click', closeModal);
}

if ('{{closeOnOverlayClick}}' === 'true') {
	overlay.addEventListener('click', function(e) {
		if (e.target === overlay) {
			closeModal();
		}
	});
}

document.addEventListener('keydown', handleKeydown);

modal.addEventListener('DOMNodeRemoved', function() {
	document.removeEventListener('keydown', handleKeydown);
});

if ('{{open}}' === 'true') {
	modal.style.display = 'block';

	switch('{{animation}}') {
		case 'fade':
			content.style.opacity = '0';
			content.style.transform = 'translate(0, -20px)';
			break;
		case 'slide':
			content.style.opacity = '0';
			content.style.transform = 'translate(0, 40px)';
			break;
		case 'zoom':
			content.style.opacity = '0';
			content.style.transform = 'scale(0.95)';
			break;
	}

	requestAnimationFrame(() => {
		content.style.opacity = '1';
		content.style.transform = 'none';
	});
}
	`})),
	nil,
	PropsDef{
		// 主要參數
		"title":               "",                // 對話框標題
		"open":                false,             // 是否顯示
		"size":                "md",              // 尺寸: xs, sm, md, lg, xl, full
		"closeButton":         true,              // 是否顯示關閉按鈕
		"closeOnEsc":          true,              // 是否按ESC鍵關閉
		"closeOnOverlayClick": true,              // 是否點擊遮罩層關閉
		"centered":            true,              // 是否垂直居中
		"scrollable":          true,              // 內容是否可滾動
		"animation":           "fade",            // 動畫效果: fade, slide, zoom
		"footer":              "",                // 底部內容
		"overlayColor":        "rgba(0,0,0,0.5)", // 遮罩層顏色
		"radius":              "md",              // 圓角大小: none, sm, md, lg
		"elevation":           "3",               // 陰影高度: 0-5
		"hideHeader":          false,             // 是否隱藏頭部
		"hideFooter":          false,             // 是否隱藏底部
		"zIndex":              "1050",            // 層級

		// 計算屬性
		"display":            "none",
		"pointerEvents":      "none",
		"modalWidth":         "500px",
		"maxWidth":           "calc(100% - 2rem)",
		"margin":             "3.75rem auto",
		"maxHeight":          "calc(100vh - 7.5rem)",
		"borderRadius":       "0.5rem",
		"boxShadow":          "0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)",
		"headerDisplay":      "flex",
		"footerDisplay":      "none",
		"closeButtonDisplay": "block",
		"contentOverflow":    "auto",
		"overlayOverflow":    "auto",
		"contentAnimation":   "modalFadeIn",
		"overlayAnimation":   "modalOverlayFadeIn",
	},
)
