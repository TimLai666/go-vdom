package components

import (
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/vdom"
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
					flex-direction: {{flexDirection}};
					gap: 0.75rem;
				`,
			},
			Label(
				Props{
					"for": "{{id}}",
					"style": `
						display: inline-flex;
						align-items: center;
						cursor: pointer;
						user-select: none;
						order: {{labelOrder}};
						font-size: {{fontSize}};
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
						"class": "switch-track",
						"style": `
							display: inline-block;
							width: {{trackWidth}};
							height: {{trackHeight}};
							background-color: {{offColor}};
							border-radius: calc({{trackHeight}} / 2);
							transition: all 0.2s ease;
							position: relative;
							cursor: pointer;
						`,
					},
					Span(
						Props{
							"class": "switch-thumb",
							"style": `
								display: block;
								width: calc({{trackHeight}} - 4px);
								height: calc({{trackHeight}} - 4px);
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
					display: {{helpDisplay}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: #64748b;
				`,
			},
			"{{helpText}}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
		const input = document.getElementById('{{id}}');
		if (!input) return;

		const track = input.nextElementSibling;
		const thumb = track.querySelector('.switch-thumb');
		if (!track || !thumb) return;

		const onColor = '{{onColor}}';
		const offColor = '{{offColor}}';
		const trackWidth = '{{trackWidth}}';
		const trackHeight = '{{trackHeight}}';
		const colorRgb = '{{colorRgb}}';

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
				thumb.style.opacity = '0.6';
			} else {
				track.style.opacity = '1';
				track.style.cursor = 'pointer';
				thumb.style.opacity = '1';
			}
		}

		// 初始化
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
				detail: { id: '{{id}}', checked: this.checked },
				bubbles: true
			}));
		});

		// Focus 效果
		input.addEventListener('focus', function() {
			if (!this.disabled) {
				track.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
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
		"checked":       "false",
		"required":      "false",
		"disabled":      "false",
		"size":          "md",
		"labelPosition": "right",
		"onColor":       "#3b82f6",
		"offColor":      "#d1d5db",
		"helpText":      "",
		"trackWidth":    "2.75rem",
		"trackHeight":   "1.5rem",
		"fontSize":      "0.9375rem",
		"flexDirection": "row",
		"labelOrder":    "1",
		"helpDisplay":   "none",
		"colorRgb":      "59, 130, 246",
	},
)
