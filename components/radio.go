package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// RadioGroup 單選按鈕組
//
// 提供現代化的單選按鈕組，允許用戶在多個選項中選擇其一。
//
// 參數:
//   - name: 單選按鈕組名稱，預設為空
//   - label: 標籤文字，預設為空
//   - options: 選項清單，以逗號分隔，如 "選項1,選項2,選項3"
//   - defaultValue: 預選值，預設為空
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
//	RadioGroup(Props{
//	  "label": "選擇性別",
//	  "name": "gender",
//	  "options": "男性,女性,其他",
//	  "defaultValue": "男性",
//	  "direction": "horizontal",
//	})
var RadioGroup = Component(
	Div(
		Props{
			"style": `
				margin-bottom: 1.25rem;
				width: 100%;
			`,
		},
		Div(
			Props{
				"style": `
					display: block;
					margin-bottom: 0.75rem;
					font-weight: 500;
					font-size: {{labelSize}};
					color: #374151;
				`,
			},
			"{{labelText}}",
		),
		Div(
			Props{
				"style": `
					display: flex;
					flex-direction: {{flexDirection}};
					gap: {{gap}};
				`,
			},
			"{{radioItems}}",
		),
		Div(
			Props{
				"style": `
					display: {{helpDisplay}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: {{helpColor}};
				`,
			},
			"{{helpMessage}}",
		),
	),
	PropsDef{
		// 主要屬性
		"name":         "",         // 單選按鈕組名稱
		"label":        "",         // 標籤文字
		"options":      "",         // 選項清單，逗號分隔
		"defaultValue": "",         // 預選值
		"required":     "false",    // 是否必填
		"disabled":     "false",    // 是否禁用
		"direction":    "vertical", // 排列方向: horizontal, vertical
		"size":         "md",       // 尺寸: sm, md, lg
		"helpText":     "",         // 幫助文字
		"errorText":    "",         // 錯誤文字
		"color":        "#3b82f6",  // 主題色

		// 計算屬性
		"labelSize":     "0.9375rem",
		"helpDisplay":   "none",
		"helpColor":     "#64748b",
		"labelText":     "",
		"helpMessage":   "",
		"radioItems":    "",
		"flexDirection": "column",
		"gap":           "0.75rem",
		"colorRgb":      "59, 130, 246",
	},
)

// Radio 單選按鈕組件
//
// 提供現代化的單選按鈕，可單獨使用或作為RadioGroup的部分。
//
// 參數:
//   - id: 按鈕ID，預設自動生成
//   - name: 按鈕名稱，預設為空
//   - value: 按鈕值，預設為空
//   - label: 標籤文字，預設為空
//   - checked: 是否選中，預設 "false"
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - color: 主題色，預設現代藍 "#3b82f6"
//
// 用法:
//
//	Radio(Props{
//	  "id": "male",
//	  "name": "gender",
//	  "value": "male",
//	  "label": "男性",
//	  "checked": "true",
//	})
var Radio = Component(
	Label(
		Props{
			"for": "{{id}}",
			"style": `
				display: flex;
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
				"type":     "radio",
				"checked":  "{{checked}}",
				"required": "{{required}}",
				"disabled": "{{disabled}}",
				"style": `
					position: absolute;
					opacity: 0;
					height: 0;
					width: 0;
					
					& + span {
						display: inline-flex;
						align-items: center;
						justify-content: center;
						width: {{radioSize}};
						height: {{radioSize}};
						border-radius: 50%;
						border: 2px solid #d1d5db;
						margin-right: 0.5rem;
						transition: all 0.2s ease;
						background: white;
					}
					
					&:checked + span {
						border-color: {{color}};
						background: white;
					}
					
					&:checked + span:after {
						content: "";
						display: block;
						width: calc({{radioSize}} - 10px);
						height: calc({{radioSize}} - 10px);
						border-radius: 50%;
						background: {{color}};
					}
					
					&:focus + span {
						box-shadow: 0 0 0 3px rgba({{colorRgb}}, 0.15);
					}
					
					&:disabled + span {
						border-color: #e5e7eb;
						background-color: #f9fafb;
					}
					
					&:disabled:checked + span:after {
						background-color: #d1d5db;
					}
				`,
			},
		),
		Span(Props{}),
		Span(
			Props{
				"style": `
					font-size: {{fontSize}};
					color: #374151;
				`,
			},
			"{{labelText}}",
		),
	),
	PropsDef{
		// 主要屬性
		"id":       "",        // 按鈕ID
		"name":     "",        // 按鈕名稱
		"value":    "",        // 按鈕值
		"label":    "",        // 標籤文字
		"checked":  "false",   // 是否選中
		"required": "false",   // 是否必填
		"disabled": "false",   // 是否禁用
		"size":     "md",      // 尺寸: sm, md, lg
		"color":    "#3b82f6", // 主題色

		// 計算屬性
		"radioSize":     "1.25rem",
		"fontSize":      "0.9375rem",
		"disabledStyle": "",
		"labelText":     "",
		"colorRgb":      "59, 130, 246",
	},
)
