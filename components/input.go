package components

import (
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
	. "github.com/TimLai666/go-vdom/dom"
)

// TextField 現代化輸入框組件
//
// 提供極簡主義設計的輸入元素，注重可用性和美感。
//
// 參數:
//   - type: 輸入類型，如 "text"、"password"、"email" 等，預設 "text"
//   - label: 標籤文字，預設為空
//   - placeholder: 提示文字，預設為空
//   - value: 預設值，預設為空
//   - id: 輸入框ID，預設自動生成
//   - name: 輸入框名稱，預設為空
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - readonly: 是否唯讀，預設 "false"
//   - pattern: 驗證模式，預設為空
//   - min: 最小值，預設為空
//   - max: 最大值，預設為空
//   - maxlength: 最大長度，預設為空
//   - autofocus: 是否自動聚焦，預設 "false"
//   - autocomplete: 自動完成，預設為空
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - variant: 變體，可選 "outlined"、"filled"、"underlined"，預設 "outlined"
//   - fullWidth: 是否填滿父容器寬度，預設 "true"
//   - icon: 圖標HTML，預設為空
//   - iconPosition: 圖標位置，可選 "left"、"right"，預設 "left"
//   - helpText: 幫助文字，預設為空
//   - errorText: 錯誤文字，預設為空
//   - labelPosition: 標籤位置，可選 "top"、"left"，預設 "top"
//   - color: 主題色，預設現代藍 "#3b82f6"
//
// 用法:
//
//	TextField(Props{
//	  "label": "電子郵件",
//	  "type": "email",
//	  "placeholder": "請輸入您的電子郵件",
//	  "required": "true",
//	})
var TextField = Component(
	Div(
		Props{
			"class": "textfield-container",
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
				"for": "{{id}}", "class": "textfield-label", "style": `
					display: {{label}};
					margin-bottom: {{labelMargin}};
					font-weight: 500;
					font-size: {{labelSize}};
					color: #374151;
					width: {{labelWidth}};
				`,
			},
			"{{label}}",
		),
		Div(
			Props{
				"class": "textfield-wrapper",
				"style": `
					position: relative;
					width: {{inputWrapWidth}};
				`,
			},
			Div(
				Props{
					"class": "textfield-icon-left",
					"style": `
						display: {{iconLeftDisplay}};
						position: absolute;
						top: 50%;
						left: 12px;
						transform: translateY(-50%);
						color: #64748b;
						z-index: 1;
					`,
				},
				"{{iconLeft}}",
			),
			Input(
				Props{
					"id":             "{{id}}",
					"name":           "{{name}}",
					"type":           "{{type}}",
					"placeholder":    "{{placeholder}}",
					"value":          "{{value}}",
					"required":       "{{required}}",
					"disabled":       "{{disabled}}",
					"readonly":       "{{readonly}}",
					"pattern":        "{{pattern}}",
					"min":            "{{min}}",
					"max":            "{{max}}",
					"maxlength":      "{{maxlength}}",
					"autofocus":      "{{autofocus}}",
					"autocomplete":   "{{autocomplete}}",
					"class":          "textfield-input",
					"data-color":     "{{color}}",
					"data-color-rgb": "{{colorRgb}}",
					"style": `
						display: block;
						width: 100%;
						padding: {{inputPadding}};
						font-size: {{fontSize}};
						line-height: 1.5;
						color: #333;
						background: {{inputBg}};
						border: {{inputBorder}};
						border-radius: {{inputRadius}};
						box-shadow: {{inputShadow}};
						transition: all 0.2s ease;
						outline: none;
						box-sizing: border-box;
					`,
				},
			),
			Div(
				Props{
					"class": "textfield-icon-right",
					"style": `
						display: {{iconRightDisplay}};
						position: absolute;
						top: 50%;
						right: 12px;
						transform: translateY(-50%);
						color: #64748b;
						z-index: 1;
					`,
				},
				"{{iconRight}}",
			),
		),
		Div(
			Props{"class": "textfield-help-text", "style": `
					display: {{helpText}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: {{helpColor}};
				`,
			},
			"{{helpText}}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `try {
    const input = document.getElementById('{{id}}');
    if (!input) return;

    const handleFocus = function() {
      if (!input.disabled && !input.readOnly) {
        input.style.borderColor = '{{color}}';
        input.style.boxShadow = '0 0 0 3px rgba(' + '{{colorRgb}}' + ', 0.15)';
      }
    };

    const handleBlur = function() {
      if (!input.disabled && !input.readOnly) {
        input.style.borderColor = '#d1d5db';
        input.style.boxShadow = '0 1px 2px rgba(0, 0, 0, 0.05)';
      }
    };

    const handleInput = function(e) {
      input.dispatchEvent(new CustomEvent('textfield:input', {
        detail: {
          id: input.id,
          value: input.value,
          type: input.type
        },
        bubbles: true
      }));
    };

    const handleChange = function(e) {
      input.dispatchEvent(new CustomEvent('textfield:change', {
        detail: {
          id: input.id,
          value: input.value,
          type: input.type
        },
        bubbles: true
      }));
    };

    input.addEventListener('focus', handleFocus);
    input.addEventListener('blur', handleBlur);
    input.addEventListener('input', handleInput);
    input.addEventListener('change', handleChange);

    // 設置禁用和唯讀狀態（使用字串比較，因為組件屬性是字串）
    input.disabled = '{{disabled}}' === 'true';
    input.readOnly = '{{readonly}}' === 'true';

    // 更新樣式
    if (input.disabled) {
      input.style.backgroundColor = '#f9fafb';
      input.style.color = '#9ca3af';
      input.style.cursor = 'not-allowed';
      input.style.pointerEvents = 'none';
    } else if (input.readOnly) {
      input.style.backgroundColor = '#f9fafb';
      input.style.cursor = 'default';
      input.style.color = '#374151';
      input.style.pointerEvents = 'auto';
    } else {
      input.style.backgroundColor = '#ffffff';
      input.style.color = '#333333';
      input.style.cursor = 'text';
      input.style.pointerEvents = 'auto';
    }
  			} catch (err) {
  				console.error('Input init error for id={{id}}', err);
  			}`})),
	PropsDef{
		// 主要屬性
		"type":          "text",     // 輸入類型
		"label":         "",         // 標籤文字
		"placeholder":   "",         // 提示文字
		"value":         "",         // 預設值
		"id":            "",         // 輸入框ID
		"name":          "",         // 輸入框名稱
		"required":      "false",    // 是否必填
		"disabled":      "false",    // 是否禁用
		"readonly":      "false",    // 是否唯讀
		"pattern":       "",         // 驗證模式
		"min":           "",         // 最小值
		"max":           "",         // 最大值
		"maxlength":     "",         // 最大長度
		"autofocus":     "false",    // 是否自動聚焦
		"autocomplete":  "",         // 自動完成
		"size":          "md",       // 尺寸: sm, md, lg
		"variant":       "outlined", // 變體: outlined, filled, underlined
		"fullWidth":     "true",     // 是否填滿父容器寬度
		"icon":          "",         // 圖標HTML
		"iconPosition":  "left",     // 圖標位置
		"helpText":      "",         // 幫助文字
		"errorText":     "",         // 錯誤文字
		"labelPosition": "top",      // 標籤位置: top, left
		"color":         "#3b82f6",  // 主題色

		// 計算屬性
		"width":            "100%",
		"flexDisplay":      "block",
		"flexAlign":        "flex-start",
		"flexGap":          "0",
		"labelWidth":       "auto",
		"labelMargin":      "0.375rem",
		"labelSize":        "0.9375rem",
		"inputWrapWidth":   "100%",
		"inputPadding":     "0.625rem 0.875rem",
		"fontSize":         "0.9375rem",
		"inputRadius":      "0.375rem",
		"inputBg":          "#ffffff",
		"inputBorder":      "1px solid #d1d5db",
		"inputShadow":      "0 1px 2px rgba(0, 0, 0, 0.05)",
		"iconLeftDisplay":  "none",
		"iconRightDisplay": "none",
		"helpDisplay":      "none",
		"helpColor":        "#64748b",
		"labelText":        "",
		"helpMessage":      "",
		"iconLeft":         "",
		"iconRight":        "",
		"colorRgb":         "59, 130, 246",
	},
)
