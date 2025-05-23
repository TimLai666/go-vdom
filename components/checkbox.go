package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// Checkbox 勾選框組件
//
// 提供現代化的勾選框，支援獨立使用或群組使用。
//
// 參數:
//   - id: 勾選框ID，預設自動生成
//   - name: 勾選框名稱，預設為空
//   - value: 勾選框值，預設為空
//   - label: 標籤文字，預設為空
//   - checked: 是否選中，預設 "false"
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - indeterminate: 是否為不確定狀態，預設 "false"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - helpText: 幫助文字，預設為空
//   - color: 主題色，預設現代藍 "#3b82f6"
//
// 用法:
//
//	Checkbox(Props{
//	  "id": "agree",
//	  "name": "terms",
//	  "label": "我同意服務條款和隱私政策",
//	  "required": "true",
//	})
var Checkbox = Component(
	Div(
		Props{
			"style": `
				margin-bottom: {{marginBottom}};
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
					{{disabledStyle}}
				`,
			},
			Input(
				Props{
					"id":       "{{id}}",
					"name":     "{{name}}",
					"value":    "{{value}}",
					"type":     "checkbox",
					"checked":  "{{checked}}",
					"required": "{{required}}",
					"disabled": "{{disabled}}",
					"style": `
						position: absolute;
						opacity: 0;
						height: 0;
						width: 0;
								display: none;
					`,
				},
			),
			Span(Props{
				"class": "checkbox-box",
				"style": `
						display: inline-flex;
						align-items: center;
						justify-content: center;
						width: {{checkboxSize}};
						height: {{checkboxSize}};
						border-radius: {{borderRadius}};
						border: 2px solid #d1d5db;
						margin-right: 0.5rem;
						transition: all 0.2s ease;
						background: white;
					`,
			}),
			Span(
				Props{
					"class": "checkbox-checkmark",
					"style": `
						display: none;
						width: calc({{checkboxSize}} / 2 - 1px);
						height: calc({{checkboxSize}} / 2 + 1px);
						border: solid white;
						border-width: 0 2px 2px 0;
						transform: rotate(45deg) translate(-2px, -2px);
					`,
				},
			),
			Span(
				Props{"style": `
						font-size: {{fontSize}};
						color: #374151;						display: ${'{{label}}'.trim() ? 'inline' : 'none'};
					`,
				},
				"{{label}}",
			),
		),
		Div(
			Props{"style": `
					display: ${{'{{helpText}}'.trim() ? 'block' : 'none'}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: {{helpColor}};
					margin-left: {{checkboxSize}};
					margin-left: calc({{checkboxSize}} + 0.5rem);
				`,
			},
			"{{helpText}}",
		),
		Script(`
			document.addEventListener('DOMContentLoaded', function() {				const boxId = '{{id}}';
				const input = document.getElementById(boxId);
				if (!input) return;
				const box = input.nextElementSibling;
				const checkmark = box.nextElementSibling;
				
				function updateCheckboxState() {
					// 確保勾選和禁用狀態正確設置
					const disabled = '{{disabled}}' === 'true';
					const checked = '{{checked}}' === 'true';
					
					input.disabled = disabled;
					input.checked = checked;
					
					if (checked) {
						box.style.borderColor = '{{color}}';
						box.style.background = '{{color}}';
						checkmark.style.display = 'block';
					} else {
						box.style.borderColor = '#d1d5db';
						box.style.background = 'white';
						checkmark.style.display = 'none';
					}
					
					if (disabled) {
						box.style.borderColor = '#e5e7eb';
						box.style.background = '#f9fafb';
						if (checked) {
							box.style.background = '#d1d5db';
						}
						box.style.cursor = 'not-allowed';
					} else {
						box.style.cursor = 'pointer';
					}
				}
				
				// 初始化狀態
				updateCheckboxState();
				
				// 切換狀態時更新
				input.addEventListener('change', updateCheckboxState);
				
				// Focus 效果
				input.addEventListener('focus', function() {
					if (!this.disabled) {
						box.style.boxShadow = '0 0 0 3px rgba({{colorRgb}}, 0.15)';
					}
				});
				
				input.addEventListener('blur', function() {
					box.style.boxShadow = 'none';
				});
				
				// 觸發自定義事件
				input.addEventListener('change', function() {
					this.dispatchEvent(new CustomEvent('checkbox:change', {
						detail: { 
							id: '{{id}}',
							checked: this.checked,
							value: this.value
						}
					}));
				});
			});
		`),
	),
	PropsDef{
		// 主要屬性
		"id":            "",        // 勾選框ID
		"name":          "",        // 勾選框名稱
		"value":         "",        // 勾選框值
		"label":         "",        // 標籤文字
		"checked":       "false",   // 是否選中
		"required":      "false",   // 是否必填
		"disabled":      "false",   // 是否禁用
		"indeterminate": "false",   // 是否為不確定狀態 (暫未實作)
		"size":          "md",      // 尺寸: sm, md, lg
		"helpText":      "",        // 幫助文字
		"color":         "#3b82f6", // 主題色

		// 計算屬性
		"checkboxSize":  "1.25rem",
		"fontSize":      "0.9375rem",
		"borderRadius":  "0.25rem",
		"marginBottom":  "1rem",
		"disabledStyle": "",
		"helpDisplay":   "none",
		"helpColor":     "#64748b",
		"labelText":     "",
		"helpMessage":   "",
		"colorRgb":      "59, 130, 246",
	},
)

// CheckboxGroup 勾選框組
//
// 提供一組相關的勾選框，適合多項選擇場景。
//
// 參數:
//   - name: 勾選框組名稱，預設為空
//   - label: 標籤文字，預設為空
//   - options: 選項清單，以逗號分隔，如 "選項1,選項2,選項3"
//   - values: 已選中值，以逗號分隔，預設為空
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - direction: 排列方向，可選 "horizontal"、"vertical"，預設 "vertical"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - helpText: 幫助文字，預設為空
//   - errorText: 錯誤文字，預設為空
//   - color: 主題色，預設現代藍 "#3b82f6"
//
// 用法:
//
//	CheckboxGroup(Props{
//	  "label": "選擇愛好",
//	  "name": "hobbies",
//	  "options": "閱讀,運動,音樂,繪畫,旅行",
//	  "values": "閱讀,音樂",
//	})
var CheckboxGroup = Component(
	Div(
		Props{
			"style": `
				margin-bottom: 1.25rem;
				width: 100%;
			`,
		},
		Div(
			Props{"style": `
					display: ${'{{label}}'.trim() ? 'block' : 'none'};
					margin-bottom: 0.75rem;
					font-weight: 500;
					font-size: {{labelSize}};
					color: #374151;
				`,
			},
			"{{label}}",
		),
		Div(
			Props{
				"style": `
					display: flex;
					flex-direction: {{flexDirection}};
					gap: {{gap}};
				`,
			},
			"{{checkboxItems}}",
		),
		Div(
			Props{"style": `
					display: ${'{{helpText}}'.trim() ? 'block' : 'none'};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: {{helpColor}};
				`,
			},
			"{{helpText}}",
		),
	),
	PropsDef{
		// 主要屬性
		"name":      "",         // 勾選框組名稱
		"label":     "",         // 標籤文字
		"options":   "",         // 選項清單，逗號分隔
		"values":    "",         // 已選中值，逗號分隔
		"required":  "false",    // 是否必填
		"disabled":  "false",    // 是否禁用
		"direction": "vertical", // 排列方向: horizontal, vertical
		"size":      "md",       // 尺寸: sm, md, lg
		"helpText":  "",         // 幫助文字
		"errorText": "",         // 錯誤文字
		"color":     "#3b82f6",  // 主題色

		// 計算屬性
		"labelSize":     "0.9375rem",
		"helpDisplay":   "none",
		"helpColor":     "#64748b",
		"labelText":     "",
		"helpMessage":   "",
		"checkboxItems": "",
		"flexDirection": "column",
		"gap":           "0.75rem",
		"colorRgb":      "59, 130, 246",
	},
)
