package components

import (
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/dom"
)

// Btn 現代化按鈕組件
//
// 提供精緻的視覺效果，自適應內容，以及細膩的交互反饋設計。
//
// 參數:
//   - variant: 按鈕變種，可選 "filled"(填充)、"outlined"(邊框)、"text"(純文字)，預設 "filled"
//   - color: 主色調，預設現代藍 "#3b82f6"
//   - size: 尺寸，可選 "sm"(小)、"md"(中)、"lg"(大)，預設 "md"
//   - fullWidth: 是否填滿父容器寬度，預設 "false"
//   - rounded: 圓角程度，可選 "none"、"sm"、"md"、"lg"、"full"，預設 "md"
//   - disabled: 是否禁用，預設 "false"
//   - id: 按鈕ID，必填，用於JavaScript事件綁定
//   - name: 按鈕名稱，預設空
//   - type: 按鈕類型，預設 "button"
//   - weight: 字重，預設 "500"
//   - icon: 圖標HTML字符(如 "&#10003;")，預設為空
//   - iconPosition: 圖標位置，可選 "left" 或 "right"，預設 "left"
//
// 用法:
//
//	Btn(Props{"id": "submit-btn", "color": "#8b5cf6", "size": "lg"}, "點擊我")
//	Btn(Props{"id": "confirm-btn", "variant": "outlined", "icon": "&#10003;"}, "確認")
var Btn = Component(
	Div(
		Props{},
		Button(
			Props{"id": "btn-{{id}}",
				"name":     "{{name}}",
				"type":     "{{type}}",
				"disabled": "{{disabled}}",
				"style": `
					display: inline-flex;
					align-items: center;
					justify-content: center;
					gap: 0.5rem;
					font-family: inherit;
					font-size: {{fontSize}};
					font-weight: {{weight}};
					line-height: 1.5;
					text-decoration: none;
					vertical-align: middle;					cursor: {{cursor}};
					user-select: none;
					padding: {{paddingY}} {{paddingX}};
					border-radius: {{buttonRadius}};
					transition: all 180ms ease-out;
					position: relative;
					overflow: hidden;
					width: {{width}};
					height: {{height}};
					letter-spacing: 0.01em;
					box-shadow: {{boxShadow}};
					background: {{background}};
					color: {{textColor}};
					border: {{border}};
					text-align: center;
					opacity: {{opacity}};
					text-transform: {{textTransform}};
				`,
			},
			Span(
				Props{
					"style": `
						display: {{iconLeftDisplay}};
						margin-right: 0.35rem;
						margin-left: -0.15rem;
					`,
				},
				"{{iconLeft}}",
			),
			"{{children}}",
			Span(
				Props{
					"style": `
						display: {{iconRightDisplay}};
						margin-left: 0.35rem;
						margin-right: -0.15rem;
					`,
				},
				"{{iconRight}}",
			),
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `try {
					const btn = document.getElementById('btn-{{id}}');
					if (!btn) return;
					// 確保禁用狀態正確設置
					btn.disabled = btn.getAttribute('disabled') === 'true';

					// 添加點擊波紋效果
					btn.addEventListener('click', function(e) {
						if (!this.disabled) {
							const rect = this.getBoundingClientRect();
							const ripple = document.createElement('div');
							const size = Math.max(rect.width, rect.height);
							const x = e.clientX - rect.left - size/2;
							const y = e.clientY - rect.top - size/2;
							ripple.style.cssText =
								'position: absolute;' +
								'left: ' + x + 'px;' +
								'top: ' + y + 'px;' +
								'width: ' + size + 'px;' +
								'height: ' + size + 'px;' +
								'background: currentColor;' +
								'border-radius: 50%;' +
								'opacity: 0.3;' +
								'transform: scale(0);' +
								'animation: ripple 0.6s linear;' +
								'pointer-events: none;';

							this.appendChild(ripple);

							setTimeout(function(){ ripple.remove(); }, 600);
							// 發出自定義點擊事件
							this.dispatchEvent(new CustomEvent('btn:click', {
								detail: { id: '{{id}}' }
							}));
						}
					});

					// 添加懸停效果
					btn.addEventListener('mouseenter', function() {
						if (!this.disabled) {
							this.style.transform = 'translateY(-1px)';
							this.style.background = '{{hoverBackground}}';
							this.style.borderColor = '{{hoverBorderColor}}';
							this.style.color = '{{hoverTextColor}}';
							this.style.boxShadow = '{{hoverBoxShadow}}';
							this.dispatchEvent(new CustomEvent('btn:hover'));
						}
					});

					btn.addEventListener('mouseleave', function() {
						if (!this.disabled) {
							this.style.transform = 'translateY(0)';
							this.style.borderColor = '{{border}}'.split(' ')[2];
							this.style.color = '{{textColor}}';
							this.style.boxShadow = '{{boxShadow}}';
						}
					});

					// 添加 focus 效果
					btn.addEventListener('focus', function() {
						if (!this.disabled) {
							this.style.outline = 'none';
							this.style.boxShadow = '0 0 0 3px {{focusRingColor}}';
						}
					});

					btn.addEventListener('blur', function() {
						if (!this.disabled) {
							this.style.boxShadow = '{{boxShadow}}';
						}
					});
				} catch (err) {
					console.error('Btn init error for id={{id}}', err);
				}`})),
	PropsDef{ // 主要參數
		"id":            "1",       // 按鈕ID，必須提供
		"variant":       "filled",  // 按鈕樣式：filled, outlined, text
		"color":         "#3b82f6", // 主色調
		"size":          "md",      // 尺寸：sm, md, lg
		"fullWidth":     "false",   // 是否填滿父容器寬度
		"rounded":       "md",      // 圓角：none, sm, md, lg, full
		"disabled":      "false",   // 是否禁用
		"name":          "",        // 按鈕名稱
		"type":          "button",  // 按鈕類型
		"weight":        "500",     // 字重
		"icon":          "",        // 圖標HTML字符
		"iconPosition":  "left",    // 圖標位置：left, right
		"textTransform": "none",    // 文字轉換：none, uppercase

		// 計算屬性 (不應由用戶直接設置)
		"fontSize":         "0.95rem",
		"paddingX":         "1.25rem",
		"paddingY":         "0.5rem",
		"buttonRadius":     "0.5rem",
		"width":            "auto",
		"height":           "auto",
		"background":       "#3b82f6",
		"textColor":        "#ffffff",
		"border":           "1px solid transparent",
		"boxShadow":        "0 1px 3px rgba(0,0,0,0.1)",
		"opacity":          "1",
		"hoverBackground":  "#2563eb",
		"hoverTextColor":   "#ffffff",
		"hoverBorderColor": "transparent",
		"hoverBoxShadow":   "0 4px 6px rgba(0,0,0,0.12)",
		"focusRingColor":   "rgba(59,130,246,0.25)",
		"cursor":           "pointer",
		"iconLeft":         "",
		"iconRight":        "",
		"iconLeftDisplay":  "none",
		"iconRightDisplay": "none",
	},
)
