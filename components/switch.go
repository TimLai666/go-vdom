package components

import (
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
					display: flex;
					align-items: center;
					cursor: pointer;
					user-select: none;
					order: {{labelOrder}};
					{{disabledStyle}}
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
				Props{"class": "switch-track",
					"style": `
						position: relative;
						display: inline-block;
						width: {{trackWidth}};
						height: {{trackHeight}};
						background-color: {{offColor}};
						border-radius: {{trackHeight}};
						transition: all 0.3s ease;
					`,
				},
			),
			Span(
				Props{
					"class": "switch-thumb",
					"style": `
						position: absolute;
						left: 3px;
						top: 3px;
						width: calc({{trackHeight}} - 6px);
						height: calc({{trackHeight}} - 6px);
						background-color: white;
						border-radius: 50%;
						transition: all 0.3s ease;
						box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
					`,
				},
			),
			Span(
				Props{
					"style": `
						margin-left: 0.5rem;
						font-size: {{fontSize}};
						color: #374151;
					`,
				},
				"{{labelText}}",
			),
		),
		Div(
			Props{
				"style": `
					display: {{helpDisplay}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: {{helpColor}};
				`,
			}, "{{helpMessage}}",
		),
		Script(`
			document.addEventListener('DOMContentLoaded', function() {
				const id = '{{id}}';
				const input = document.getElementById(id);
				const track = input.nextElementSibling;
				const thumb = track.nextElementSibling;

				function updateSwitchState() {
					if (input.checked) {
						track.style.backgroundColor = '{{onColor}}';
						thumb.style.transform = 'translateX(calc({{trackWidth}} - {{trackHeight}}))';
					} else {
						track.style.backgroundColor = '{{offColor}}';
						thumb.style.transform = 'translateX(0)';
					}

					if (input.disabled) {
						track.style.opacity = '0.6';
						track.style.cursor = 'not-allowed';
						thumb.style.opacity = '0.6';
					} else {
						track.style.opacity = '1';
						track.style.cursor = 'pointer';
						thumb.style.opacity = '1';
					}
				}

				// 初始化狀態
				updateSwitchState();

				// 切換狀態時更新
				input.addEventListener('change', updateSwitchState);

				// Focus 效果
				input.addEventListener('focus', function() {
					if (!this.disabled) {
						track.style.boxShadow = '0 0 0 3px rgba({{colorRgb}}, 0.15)';
					}
				});

				input.addEventListener('blur', function() {
					track.style.boxShadow = 'none';
				});

				// 觸發自定義事件
				input.addEventListener('change', function() {
					this.dispatchEvent(new CustomEvent('switch:change', {
						detail: {
							id: '{{id}}',
							checked: this.checked
						}
					}));
				});
			});
		`),
	),
	JSAction{},
	PropsDef{
		// 主要屬性
		"id":            "",        // 開關ID
		"name":          "",        // 開關名稱
		"label":         "",        // 標籤文字
		"checked":       "false",   // 是否開啟
		"required":      "false",   // 是否必填
		"disabled":      "false",   // 是否禁用
		"size":          "md",      // 尺寸: sm, md, lg
		"labelPosition": "right",   // 標籤位置: right, left
		"onColor":       "#3b82f6", // 開啟時的顏色
		"offColor":      "#d1d5db", // 關閉時的顏色
		"helpText":      "",        // 幫助文字

		// 計算屬性
		"trackWidth":    "2.75rem",
		"trackHeight":   "1.5rem",
		"fontSize":      "0.9375rem",
		"helpDisplay":   "none",
		"helpColor":     "#64748b",
		"labelText":     "",
		"helpMessage":   "",
		"flexDirection": "row",
		"labelOrder":    "0",
		"disabledStyle": "",
		"colorRgb":      "59, 130, 246",
	},
)
