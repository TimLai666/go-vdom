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
				display: ${{{open}} === true ? 'block' : 'none'};
				z-index: {{zIndex}};
				overflow-x: hidden;
				overflow-y: auto;
				animation: modalOverlayFadeIn 0.3s ease;
				pointer-events: ${{{open}} === true ? 'auto' : 'none'};

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
					width: ${{{size}} === 'xs' ? '300px' : {{size}} === 'sm' ? '400px' : {{size}} === 'md' ? '500px' : {{size}} === 'lg' ? '700px' : {{size}} === 'xl' ? '900px' : '100%'};
					max-width: calc(100% - 2rem);
					margin: ${{{centered}} === true ? '3.75rem auto' : '1rem auto'};
					background: #ffffff;
					border-radius: ${{{radius}} === 'none' ? '0' : {{radius}} === 'sm' ? '0.25rem' : {{radius}} === 'lg' ? '0.75rem' : '0.5rem'};
					box-shadow: ${{{elevation}} === '0' ? 'none' : {{elevation}} === '1' ? '0 1px 3px rgba(0,0,0,0.1)' : {{elevation}} === '2' ? '0 4px 6px rgba(0,0,0,0.1)' : {{elevation}} === '3' ? '0 10px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)' : {{elevation}} === '4' ? '0 20px 30px -10px rgba(0,0,0,0.15)' : '0 25px 40px -15px rgba(0,0,0,0.2)'};
					display: flex;
					flex-direction: column;
					max-height: calc(100vh - 7.5rem);
					animation: ${{{animation}} === 'slide' ? 'modalSlideIn' : {{animation}} === 'zoom' ? 'modalZoomIn' : 'modalFadeIn'} 0.3s ease;

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
						display: ${{{hideHeader}} === true || {{title}}.trim() === '' ? 'none' : 'flex'};
						padding: 1rem 1.5rem;
						border-bottom: 1px solid #f1f5f9;
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
							display: ${{{closeButton}} === false ? 'none' : 'block'};
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
						overflow-y: ${{{scrollable}} === true ? 'auto' : 'visible'};
					`,
				},
				"{{children}}",
			),
			Div(
				Props{
					"style": `
						display: ${{{hideFooter}} === true || {{footer}}.trim() === '' ? 'none' : 'flex'};
						padding: 1rem 1.5rem;
						border-top: 1px solid #f1f5f9;
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
	if (event.key === 'Escape' && {{closeOnEsc}} === true) {
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

if ({{closeOnOverlayClick}} === true) {
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

if ({{open}} === true) {
	modal.style.display = 'block';

	switch({{animation}}) {
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
	},
)
