package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// Dropdown 下拉式選單組件
//
// 提供現代化的下拉選單，支援單選及搜尋功能。
//
// 參數:
//   - label: 標籤文字，預設為空
//   - id: 選單ID，預設自動生成
//   - name: 選單名稱，預設為空
//   - options: 選項清單，以逗號分隔，如 "選項1,選項2,選項3"
//   - defaultValue: 預選值，預設為空
//   - placeholder: 提示文字，預設為 "請選擇"
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - searchable: 是否可搜尋，預設 "false"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - fullWidth: 是否填滿父容器寬度，預設 "true"
//   - icon: 圖標HTML，預設為下拉箭頭
//   - helpText: 幫助文字，預設為空
//   - errorText: 錯誤文字，預設為空
//   - color: 主題色，預設現代藍 "#3b82f6"
//
// 用法:
//
//	Dropdown(Props{
//	  "label": "選擇國家",
//	  "options": "台灣,中國,日本,美國,韓國",
//	  "placeholder": "請選擇國家",
//	  "required": "true",
//	})
var Dropdown = Component(
	Div(
		Props{
			"style": `
				margin-bottom: 1.25rem;
				width: {{width}};
				display: {{flexDisplay}};
				align-items: {{flexAlign}};
				gap: {{flexGap}};
			`,
		},
		Label(
			Props{
				"for": "{{id}}",
				"style": `
					display: block;
					margin-bottom: {{labelMargin}};
					font-weight: 500;
					font-size: {{labelSize}};
					color: #374151;
					width: {{labelWidth}};
				`,
			},
			"{{labelText}}",
		),
		Div(
			Props{
				"style": `
					position: relative;
					width: {{inputWrapWidth}};
				`,
			},
			Select(
				Props{
					"id":       "{{id}}",
					"name":     "{{name}}",
					"required": "{{required}}",
					"disabled": "{{disabled}}",
					"style": `
						display: block;
						width: 100%;
						padding: {{selectPadding}};
						font-size: {{fontSize}};
						line-height: 1.5;
						background: {{selectBg}};
						color: #333;
						border: {{selectBorder}};
						border-radius: {{selectRadius}};
						box-shadow: {{selectShadow}};
						transition: all 0.2s ease;
						appearance: none;
						background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23637381' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
						background-repeat: no-repeat;
						background-position: right 0.75rem center;
						background-size: 1rem;
						padding-right: 2.5rem;
						outline: none;
						box-sizing: border-box;
						
						&:focus {
							border-color: {{color}};
							box-shadow: 0 0 0 3px rgba({{colorRgb}}, 0.15);
						}
						
						&:disabled {
							background-color: #f9fafb;
							color: #9ca3af;
							cursor: not-allowed;
						}
					`,
				},
				Option(
					Props{
						"value": "",
						"style": `
							color: #94a3b8;
						`,
					},
					"{{placeholder}}",
				),
				"{{optionItems}}",
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
	),
	PropsDef{
		// 主要屬性
		"label":         "",        // 標籤文字
		"id":            "",        // 選單ID
		"name":          "",        // 選單名稱
		"options":       "",        // 選項清單，逗號分隔
		"defaultValue":  "",        // 預選值
		"placeholder":   "請選擇",     // 提示文字
		"required":      "false",   // 是否必填
		"disabled":      "false",   // 是否禁用
		"searchable":    "false",   // 是否可搜尋（目前無實作）
		"size":          "md",      // 尺寸: sm, md, lg
		"fullWidth":     "true",    // 是否填滿父容器寬度
		"helpText":      "",        // 幫助文字
		"errorText":     "",        // 錯誤文字
		"labelPosition": "top",     // 標籤位置: top, left
		"color":         "#3b82f6", // 主題色

		// 計算屬性
		"width":          "100%",
		"flexDisplay":    "block",
		"flexAlign":      "flex-start",
		"flexGap":        "0",
		"labelWidth":     "auto",
		"labelMargin":    "0.375rem",
		"labelSize":      "0.9375rem",
		"inputWrapWidth": "100%",
		"selectPadding":  "0.625rem 0.875rem",
		"fontSize":       "0.9375rem",
		"selectRadius":   "0.375rem",
		"selectBg":       "#ffffff",
		"selectBorder":   "1px solid #d1d5db",
		"selectShadow":   "0 1px 2px rgba(0, 0, 0, 0.05)",
		"helpDisplay":    "none",
		"helpColor":      "#64748b",
		"labelText":      "",
		"helpMessage":    "",
		"optionItems":    "",
		"colorRgb":       "59, 130, 246",
	},
)
