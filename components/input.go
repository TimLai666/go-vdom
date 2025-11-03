package components

import (
	"strings"

	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
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
//	  "icon": "📧",
//	})
func TextField(props Props, children ...VNode) VNode {
	// Compute derived properties
	hasIcon := false
	if icon, ok := props["icon"]; ok {
		if iconStr, ok := icon.(string); ok && strings.TrimSpace(iconStr) != "" {
			hasIcon = true
		}
	}
	props["hasIcon"] = hasIcon

	hasError := false
	if errorText, ok := props["errorText"]; ok {
		if errorStr, ok := errorText.(string); ok && strings.TrimSpace(errorStr) != "" {
			hasError = true
		}
	}
	props["hasError"] = hasError

	hasHelp := false
	if helpText, ok := props["helpText"]; ok {
		if helpStr, ok := helpText.(string); ok && strings.TrimSpace(helpStr) != "" {
			hasHelp = true
		}
	}
	props["hasHelp"] = hasHelp

	return textFieldInternal(props, children...)
}

var textFieldInternal = Component(
	Div(
		Props{
			"class": "textfield-container",
			"style": `
				margin-bottom: 1.25rem;
				width: ${'{{fullWidth}}' === 'true' ? '100%' : 'auto'};
				display: ${'{{labelPosition}}' === 'left' ? 'flex' : 'block'};
				align-items: ${'{{labelPosition}}' === 'left' ? 'center' : 'flex-start'};
				gap: ${'{{labelPosition}}' === 'left' ? '1rem' : '0'};
			`,
		},
		Label(
			Props{
				"for": "{{id}}", "class": "textfield-label", "style": `
					display: ${'{{label}}' !== '' ? 'block' : 'none'};
					margin-bottom: ${'{{labelPosition}}' === 'top' ? '0.375rem' : '0'};
					font-weight: 500;
					font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1rem' : '0.9375rem'};
					color: #374151;
					width: ${'{{labelPosition}}' === 'left' ? '120px' : 'auto'};
					flex-shrink: 0;
				`,
			},
			"{{label}}",
		),
		Div(
			Props{
				"class": "textfield-wrapper",
				"style": `
					position: relative;
					width: ${'{{labelPosition}}' === 'left' ? 'calc(100% - 120px - 1rem)' : '100%'};
					flex: ${'{{labelPosition}}' === 'left' ? '1' : 'none'};
				`,
			},
			Div(
				Props{
					"class": "textfield-icon-left",
					"style": `
						display: ${'{{hasIcon}}' === 'true' ? '{{iconPosition}}' === 'left' ? 'flex' : 'none' : 'none'};
						position: absolute;
						top: 50%;
						left: 12px;
						transform: translateY(-50%);
						color: #64748b;
						z-index: 1;
						align-items: center;
						justify-content: center;
						pointer-events: none;
					`,
				},
				"{{icon}}",
			),
			Input(
				Props{
					"id":           "{{id}}",
					"name":         "{{name}}",
					"type":         "{{type}}",
					"placeholder":  "{{placeholder}}",
					"value":        "{{value}}",
					"required":     "{{required}}",
					"disabled":     "{{disabled}}",
					"readonly":     "{{readonly}}",
					"pattern":      "{{pattern}}",
					"min":          "{{min}}",
					"max":          "{{max}}",
					"maxlength":    "{{maxlength}}",
					"autofocus":    "{{autofocus}}",
					"autocomplete": "{{autocomplete}}",
					"class":        "textfield-input",
					"data-color":   "{{color}}",
					"style": `
						display: block;
						width: 100%;
						padding: ${'{{size}}' === 'sm' ? '0.5rem 0.75rem' : '{{size}}' === 'lg' ? '0.75rem 1rem' : '0.625rem 0.875rem'};
						padding-left: ${'{{hasIcon}}' === 'true' ? '{{iconPosition}}' === 'left' ? '2.5rem' : ${'{{size}}' === 'sm' ? '0.75rem' : '{{size}}' === 'lg' ? '1rem' : '0.875rem'} : ${'{{size}}' === 'sm' ? '0.75rem' : '{{size}}' === 'lg' ? '1rem' : '0.875rem'}};
						padding-right: ${'{{hasIcon}}' === 'true' ? '{{iconPosition}}' === 'right' ? '2.5rem' : ${'{{size}}' === 'sm' ? '0.75rem' : '{{size}}' === 'lg' ? '1rem' : '0.875rem'} : ${'{{size}}' === 'sm' ? '0.75rem' : '{{size}}' === 'lg' ? '1rem' : '0.875rem'}};
						font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1rem' : '0.9375rem'};
						line-height: 1.5;
						color: #333;
						background: ${'{{variant}}' === 'filled' ? '#f9fafb' : '#ffffff'};
						border: ${'{{variant}}' === 'outlined' ? '1px solid #d1d5db' : '{{variant}}' === 'filled' ? '1px solid transparent' : 'none'};
						border-bottom: ${'{{variant}}' === 'underlined' ? '1px solid #d1d5db' : ''};
						border-radius: ${'{{variant}}' === 'underlined' ? '0' : '0.375rem'};
						box-shadow: ${'{{variant}}' === 'outlined' ? '0 1px 2px rgba(0, 0, 0, 0.05)' : 'none'};
						transition: all 0.2s ease;
						outline: none;
						box-sizing: border-box;
						cursor: ${'{{disabled}}' === 'true' ? 'not-allowed' : '{{readonly}}' === 'true' ? 'default' : 'text'};
						opacity: ${'{{disabled}}' === 'true' ? '0.6' : '1'};
					`,
				},
			),
			Div(
				Props{
					"class": "textfield-icon-right",
					"style": `
						display: ${'{{hasIcon}}' === 'true' ? '{{iconPosition}}' === 'right' ? 'flex' : 'none' : 'none'};
						position: absolute;
						top: 50%;
						right: 12px;
						transform: translateY(-50%);
						color: #64748b;
						z-index: 1;
						align-items: center;
						justify-content: center;
						pointer-events: none;
					`,
				},
				"{{icon}}",
			),
		),
		Div(
			Props{"class": "textfield-help-text", "style": `
					display: ${'{{hasError}}' === 'true' ? 'block' : '{{hasHelp}}' === 'true' ? 'block' : 'none'};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: ${'{{hasError}}' === 'true' ? '#ef4444' : '#64748b'};
				`,
			},
			"${'{{hasError}}' === 'true' ? '{{errorText}}' : '{{helpText}}'}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `try {
    const input = document.getElementById('{{id}}');
    if (!input) return;

    const color = input.getAttribute('data-color') || '{{color}}';

    // 計算RGB值用於陰影
    function hexToRgb(hex) {
        const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ?
            parseInt(result[1], 16) + ', ' + parseInt(result[2], 16) + ', ' + parseInt(result[3], 16)
            : '59, 130, 246';
    }
    const colorRgb = hexToRgb(color);

    const handleFocus = function() {
      if (!input.disabled && !input.readOnly) {
        input.style.borderColor = color;
        input.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
      }
    };

    const handleBlur = function() {
      if (!input.disabled && !input.readOnly) {
        const variant = '{{variant}}';
        if (variant === 'outlined') {
          input.style.borderColor = '#d1d5db';
          input.style.boxShadow = '0 1px 2px rgba(0, 0, 0, 0.05)';
        } else if (variant === 'underlined') {
          input.style.borderBottomColor = '#d1d5db';
          input.style.boxShadow = 'none';
        } else {
          input.style.borderColor = 'transparent';
          input.style.boxShadow = 'none';
        }
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
      const variant = '{{variant}}';
      if (variant === 'filled') {
        input.style.backgroundColor = '#f9fafb';
      } else {
        input.style.backgroundColor = '#ffffff';
      }
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
		"required":      false,      // 是否必填
		"disabled":      false,      // 是否禁用
		"readonly":      false,      // 是否唯讀
		"pattern":       "",         // 驗證模式
		"min":           "",         // 最小值
		"max":           "",         // 最大值
		"maxlength":     "",         // 最大長度
		"autofocus":     false,      // 是否自動聚焦
		"autocomplete":  "",         // 自動完成
		"size":          "md",       // 尺寸: sm, md, lg
		"variant":       "outlined", // 變體: outlined, filled, underlined
		"fullWidth":     true,       // 是否填滿父容器寬度
		"icon":          "",         // 圖標HTML
		"iconPosition":  "left",     // 圖標位置: left, right
		"helpText":      "",         // 幫助文字
		"errorText":     "",         // 錯誤文字
		"labelPosition": "top",      // 標籤位置: top, left
		"color":         "#3b82f6",  // 主題色
		"hasIcon":       false,      // 計算屬性: 是否有圖標
		"hasError":      false,      // 計算屬性: 是否有錯誤
		"hasHelp":       false,      // 計算屬性: 是否有幫助文字
	},
)
