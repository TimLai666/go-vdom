package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
)

// Btn 現代化按鈕組件
//
// 提供精緻的視覺效果，自適應內容，以及細膩的交互反饋設計。
//
// 參數:
//   - variant: 按鈕變種，可選 "filled"(填充)、"outlined"(邊框)、"text"(純文字)，預設 "filled"
//   - color: 主色調，預設現代藍 "#3b82f6"
//   - size: 尺寸，可選 "sm"(小)、"md"(中)、"lg"(大)，預設 "md"
//   - fullWidth: 是否填滿父容器寬度，預設 false
//   - rounded: 圓角程度，可選 "none"、"sm"、"md"、"lg"、"full"，預設 "md"
//   - disabled: 是否禁用，預設 false
//   - id: 按鈕ID，用於JavaScript事件綁定
//   - name: 按鈕名稱，預設空
//   - type: 按鈕類型，預設 "button"
//   - weight: 字重，預設 "500"
//   - icon: 圖標HTML字符(如 "&#10003;")，預設為空
//   - iconPosition: 圖標位置，可選 "left" 或 "right"，預設 "left"
//
// 用法:
//
//	Btn(Props{"id": "submit-btn", "color": "#8b5cf6", "size": "lg"}, Text("點擊我"))
//	Btn(Props{"id": "confirm-btn", "variant": "outlined", "icon": "&#10003;"}, Text("確認"))
func Btn(props Props, children ...VNode) VNode {
	// 預計算 borderValue 以避免在模板表達式中使用字串拼接
	color := "#3b82f6" // 預設顏色
	if c, ok := props["color"]; ok {
		if colorStr, ok := c.(string); ok {
			color = colorStr
		}
	}
	props["borderValue"] = "1px solid " + color

	return btnInternal(props, children...)
}

var btnInternal = Component(
	Div(
		Props{},
		Button(
			Props{
				"id":       "btn-{{id}}",
				"name":     "{{name}}",
				"type":     "{{type}}",
				"disabled": "{{disabled}}",
				"style": `
					display: inline-flex;
					align-items: center;
					justify-content: center;
					gap: 0.5rem;
					font-family: inherit;
					font-size: ${{{size}} === "sm" ? '0.875rem' : {{size}} === "lg" ? '1.125rem' : '0.95rem'};
					font-weight: {{weight}};
					line-height: 1.5;
					text-decoration: none;
					vertical-align: middle;
					cursor: ${{{disabled}} === true ? 'not-allowed' : 'pointer'};
					user-select: none;
					padding: ${{{size}} === "sm" ? '0.375rem 1rem' : {{size}} === "lg" ? '0.625rem 1.5rem' : '0.5rem 1.25rem'};
					border-radius: ${{{rounded}} === "none" ? '0' : {{rounded}} === "sm" ? '0.25rem' : {{rounded}} === "lg" ? '0.75rem' : {{rounded}} === "full" ? '9999px' : '0.5rem'};
					transition: all 180ms ease-out;
					position: relative;
					overflow: hidden;
					width: ${{{fullWidth}} === true ? '100%' : 'auto'};
					height: auto;
					letter-spacing: 0.01em;
					box-shadow: ${{{variant}} === "outlined" ? 'none' : {{variant}} === "text" ? 'none' : '0 1px 3px rgba(0,0,0,0.1)'};
					background: ${{{variant}} === "outlined" ? 'transparent' : {{variant}} === "text" ? 'transparent' : {{color}}};
					color: ${{{variant}} === "outlined" ? {{color}} : {{variant}} === "text" ? {{color}} : '#ffffff'};
					border: ${{{variant}} === "outlined" ? {{borderValue}} : '1px solid transparent'};
					text-align: center;
					opacity: ${{{disabled}} === true ? '0.6' : '1'};
					text-transform: {{textTransform}};
				`,
			},
			Span(
				Props{
					"style": `
						display: ${{{icon}}.trim() !== '' && {{iconPosition}} !== "right" ? 'inline' : 'none'};
						margin-right: 0.35rem;
						margin-left: -0.15rem;
					`,
				},
				"{{icon}}",
			),
			"{{children}}",
			Span(
				Props{
					"style": `
						display: ${{{icon}}.trim() !== '' && {{iconPosition}} === "right" ? 'inline' : 'none'};
						margin-left: 0.35rem;
						margin-right: -0.15rem;
					`,
				},
				"{{icon}}",
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

		// 添加懸停效果 - 根據 variant 調整
		const variant = {{variant}};
		const color = {{color}};
		const isDisabled = {{disabled}} === true;

		if (!isDisabled) {
			btn.addEventListener('mouseenter', function() {
				this.style.transform = 'translateY(-1px)';

				if (variant === 'filled') {
					// filled 變體：背景變深
					this.style.background = darkenColor(color);
					this.style.boxShadow = '0 4px 6px rgba(0,0,0,0.12)';
				} else if (variant === 'outlined') {
					// outlined 變體：添加淺色背景
					this.style.background = color + '10';
					this.style.boxShadow = '0 2px 4px rgba(0,0,0,0.08)';
				} else {
					// text 變體：只添加淺色背景
					this.style.background = color + '10';
				}

				this.dispatchEvent(new CustomEvent('btn:hover'));
			});

			btn.addEventListener('mouseleave', function() {
				this.style.transform = 'translateY(0)';

				if (variant === 'filled') {
					this.style.background = color;
					this.style.boxShadow = '0 1px 3px rgba(0,0,0,0.1)';
				} else if (variant === 'outlined') {
					this.style.background = 'transparent';
					this.style.boxShadow = 'none';
				} else {
					this.style.background = 'transparent';
				}
			});

			// 添加 focus 效果
			btn.addEventListener('focus', function() {
				this.style.outline = 'none';
				this.style.boxShadow = '0 0 0 3px ' + color + '40';
			});

			btn.addEventListener('blur', function() {
				if (variant === 'filled') {
					this.style.boxShadow = '0 1px 3px rgba(0,0,0,0.1)';
				} else {
					this.style.boxShadow = 'none';
				}
			});
		}

		// 輔助函數：使顏色變深
		function darkenColor(hex) {
			// 簡單的顏色變深邏輯
			const rgb = hexToRgb(hex);
			if (!rgb) return hex;

			const factor = 0.85;
			const r = Math.floor(rgb.r * factor);
			const g = Math.floor(rgb.g * factor);
			const b = Math.floor(rgb.b * factor);

			return rgbToHex(r, g, b);
		}

		function hexToRgb(hex) {
			const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
			return result ? {
				r: parseInt(result[1], 16),
				g: parseInt(result[2], 16),
				b: parseInt(result[3], 16)
			} : null;
		}

		function rgbToHex(r, g, b) {
			return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
		}
	} catch (err) {
		console.error('Btn init error for id={{id}}', err);
	}`})),
	PropsDefault{
		"id":            "",
		"variant":       "filled",
		"color":         "#3b82f6",
		"size":          "md",
		"fullWidth":     false,
		"rounded":       "md",
		"disabled":      false,
		"name":          "",
		"type":          "button",
		"weight":        "500",
		"icon":          "",
		"iconPosition":  "left",
		"textTransform": "none",
	},
)
