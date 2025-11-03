package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
)

// Switch 開關組件
//
// 提供現代化的開關控件，適合表示開啟/關閉狀態。
//
// 參數:
//   - id: 開關ID，預設自動生成
//   - name: 開關名稱，預設為空
//   - label: 標籤文字，預設為空
//   - checked: 是否開啟，預設 "false"
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - labelPosition: 標籤位置，可選 "right"、"left"，預設 "right"
//   - onColor: 開啟時的顏色，預設現代藍 "#3b82f6"
//   - offColor: 關閉時的顏色，預設灰色 "#d1d5db"
//   - helpText: 幫助文字，預設為空
//
// 用法:
//
//	Switch(Props{
//	  "id": "notifications",
//	  "name": "notifications",
//	  "label": "啟用通知",
//	  "checked": "true",
//	  "onColor": "#10b981",
//	})
var Switch = Component(
	Div(
		Props{
			"style": `
				margin-bottom: 1rem;
			`,
		},
		Div(
			Props{
				"style": `
					display: flex;
					align-items: center;
					flex-direction: ${{{labelPosition}} === 'left' ? 'row' : 'row'};
					gap: 0.75rem;
				`,
			},
			Label(
				Props{
					"for": "{{id}}",
					"style": `
						display: ${{{label}}.trim() ? 'inline-flex' : 'none'};
						align-items: center;
						cursor: pointer;
						user-select: none;
						order: ${{{labelPosition}} === 'left' ? '0' : '1'};
						font-size: ${{{size}} === 'sm' ? '0.875rem' : {{size}} === 'lg' ? '1rem' : '0.9375rem'};
						color: #374151;
					`,
				},
				"{{label}}",
			),
			Div(
				Props{
					"style": `
						position: relative;
						display: inline-flex;
						align-items: center;
					`,
				},
				Input(
					Props{
						"id":       "{{id}}",
						"name":     "{{name}}",
						"type":     "checkbox",
						"checked":  "{{checked}}",
						"required": "{{required}}",
						"disabled": "{{disabled}}",
						"style": `
							position: absolute;
							opacity: 0;
							height: 1px;
							width: 1px;
							margin: -1px;
							padding: 0;
							border: 0;
							overflow: hidden;
							clip: rect(0 0 0 0);
							white-space: nowrap;
						`,
					},
				),
				Span(
					Props{
						"class":          "switch-track",
						"data-on-color":  "{{onColor}}",
						"data-off-color": "{{offColor}}",
						"style": `
							display: inline-block;
							width: ${{{size}} === 'sm' ? '2.25rem' : {{size}} === 'lg' ? '3.25rem' : '2.75rem'};
							height: ${{{size}} === 'sm' ? '1.25rem' : {{size}} === 'lg' ? '1.75rem' : '1.5rem'};
							background-color: {{offColor}};
							border-radius: 9999px;
							transition: all 0.2s ease;
							position: relative;
							cursor: ${{{disabled}} === true ? 'not-allowed' : 'pointer'};
							opacity: ${{{disabled}} === true ? '0.6' : '1'};
						`,
					},
					Span(
						Props{
							"class": "switch-thumb",
							"style": `
								display: block;
								width: calc(${{{size}} === 'sm' ? '1.25rem' : {{size}} === 'lg' ? '1.75rem' : '1.5rem'} - 4px);
								height: calc(${{{size}} === 'sm' ? '1.25rem' : {{size}} === 'lg' ? '1.75rem' : '1.5rem'} - 4px);
								background-color: white;
								border-radius: 50%;
								transition: all 0.2s ease;
								position: absolute;
								top: 2px;
								left: 2px;
								box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
							`,
						},
					),
				),
			),
		),
		Div(
			Props{
				"style": `
					display: ${{{helpText}}.trim() ? 'block' : 'none'};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: #64748b;
				`,
			},
			"{{helpText}}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
		const input = document.getElementById({{id}});
		if (!input) return;

		const track = input.nextElementSibling;
		const thumb = track.querySelector('.switch-thumb');
		if (!track || !thumb) return;

		const onColor = track.getAttribute('data-on-color') || {{onColor}};
		const offColor = track.getAttribute('data-off-color') || {{offColor}};
		const size = {{size}};
		const trackWidth = size === 'sm' ? '2.25rem' : size === 'lg' ? '3.25rem' : '2.75rem';
		const trackHeight = size === 'sm' ? '1.25rem' : size === 'lg' ? '1.75rem' : '1.5rem';

		function updateState() {
			const checked = input.checked;
			const disabled = input.disabled;

			if (checked) {
				track.style.backgroundColor = onColor;
				thumb.style.transform = 'translateX(calc(' + trackWidth + ' - ' + trackHeight + '))';
			} else {
				track.style.backgroundColor = offColor;
				thumb.style.transform = 'translateX(0)';
			}

			if (disabled) {
				track.style.opacity = '0.6';
				track.style.cursor = 'not-allowed';
			} else {
				track.style.opacity = '1';
				track.style.cursor = 'pointer';
			}
		}

		// 初始化
		input.checked = {{checked}};
		input.disabled = {{disabled}};
		updateState();

		// 點擊 track 切換狀態
		track.addEventListener('click', function(e) {
			e.preventDefault();
			if (!input.disabled) {
				input.checked = !input.checked;
				updateState();
				input.dispatchEvent(new Event('change', { bubbles: true }));
			}
		});

		// 監聽 change 事件
		input.addEventListener('change', function() {
			updateState();
			this.dispatchEvent(new CustomEvent('switch:change', {
				detail: { id: {{id}}, checked: this.checked },
				bubbles: true
			}));
		});

		// Focus 效果
		input.addEventListener('focus', function() {
			if (!this.disabled) {
				track.style.boxShadow = '0 0 0 3px rgba(59, 130, 246, 0.15)';
			}
		});

		input.addEventListener('blur', function() {
			track.style.boxShadow = 'none';
		});
	`})),
	PropsDef{
		"id":            "",
		"name":          "",
		"label":         "",
		"checked":       false,
		"disabled":      false,
		"required":      false,
		"size":          "md",
		"labelPosition": "right",
		"helpText":      "",
		"onColor":       "#3b82f6",
		"offColor":      "#d1d5db",
	},
)
