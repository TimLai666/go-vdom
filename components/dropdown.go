package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
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
//   - helpText: 幫助文字，預設為空
//   - errorText: 錯誤文字，預設為空
//   - labelPosition: 標籤位置，可選 "top"、"left"，預設 "top"
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
			"class": "dropdown-container",
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
				"for": "{{id}}", "class": "dropdown-label", "style": `
					display: ${'{{label}}'.trim() ? 'block' : 'none'};
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
				"class": "dropdown-wrapper",
				"style": `
					position: relative;
					width: ${'{{labelPosition}}' === 'left' ? 'calc(100% - 120px - 1rem)' : '100%'};
					flex: ${'{{labelPosition}}' === 'left' ? '1' : 'none'};
				`,
			},
			Select(
				Props{
					"id":         "{{id}}",
					"name":       "{{name}}",
					"required":   "{{required}}",
					"disabled":   "{{disabled}}",
					"class":      "dropdown-select",
					"data-color": "{{color}}",
					"style": `
						display: block;
						width: 100%;
						padding: ${'{{size}}' === 'sm' ? '0.5rem 2.5rem 0.5rem 0.75rem' : '{{size}}' === 'lg' ? '0.75rem 2.75rem 0.75rem 1rem' : '0.625rem 2.5rem 0.625rem 0.875rem'};
						font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1rem' : '0.9375rem'};
						line-height: 1.5;
						background: #ffffff;
						color: #333;
						border: 1px solid #d1d5db;
						border-radius: 0.375rem;
						box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
						transition: all 0.2s ease;
						cursor: ${'{{disabled}}' === 'true' ? 'not-allowed' : 'pointer'};
						outline: none;
						box-sizing: border-box;
						-webkit-appearance: none;
						-moz-appearance: none;
						background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23637381' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
						background-repeat: no-repeat;
						background-position: right 0.75rem center;
						background-size: 1rem;
						opacity: ${'{{disabled}}' === 'true' ? '0.6' : '1'};
					`,
				},
				Option(
					Props{
						"value": "",
						"class": "dropdown-placeholder",
						"style": "color: #94a3b8;",
					},
					"{{placeholder}}",
				),
			),
		),
		Div(
			Props{"class": "dropdown-help-text", "style": `
					display: ${'{{errorText}}'.trim() ? 'block' : '{{helpText}}'.trim() ? 'block' : 'none'};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: ${'{{errorText}}'.trim() ? '#ef4444' : '#64748b'};
				`,
			},
			"${'{{errorText}}'.trim() ? '{{errorText}}' : '{{helpText}}'}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `try {
    const selectId = '{{id}}';
    const select = document.getElementById(selectId);
    if (!select) return;

    const color = select.getAttribute('data-color') || '{{color}}';

    // 計算RGB值用於陰影
    function hexToRgb(hex) {
        const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ?
            parseInt(result[1], 16) + ', ' + parseInt(result[2], 16) + ', ' + parseInt(result[3], 16)
            : '59, 130, 246';
    }
    const colorRgb = hexToRgb(color);

    // 清除現有選項
    while (select.firstChild) {
        select.removeChild(select.firstChild);
    }

    // 添加預設選項
    const placeholder = document.createElement('option');
    placeholder.value = '';
    placeholder.textContent = '{{placeholder}}';
    placeholder.disabled = true;
    placeholder.selected = true;
    select.appendChild(placeholder);

    // 解析選項
    const options = '{{options}}'.split(',').filter(opt => opt.trim());
    const defaultValue = '{{defaultValue}}';

    // 創建選項
    options.forEach(option => {
        const opt = document.createElement('option');
        opt.value = option.trim();
        opt.textContent = option.trim();
        if (option.trim() === defaultValue) {
            opt.selected = true;
        }
        select.appendChild(opt);
    });

    select.addEventListener('focus', function() {
        if (!select.disabled) {
            select.style.borderColor = color;
            select.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
        }
    });

    select.addEventListener('blur', function() {
        if (!select.disabled) {
            select.style.borderColor = '#d1d5db';
            select.style.boxShadow = '0 1px 2px rgba(0, 0, 0, 0.05)';
        }
    });

    select.addEventListener('change', function() {
        select.dispatchEvent(new CustomEvent('dropdown:change', {
            detail: {
                value: select.value,
                id: select.id
            },
            bubbles: true
        }));
    });

    // 設置初始狀態
    const disabled = '{{disabled}}' === 'true';
    if (disabled) {
        select.disabled = true;
        select.style.backgroundColor = '#f9fafb';
        select.style.color = '#9ca3af';
        select.style.cursor = 'not-allowed';
    }
} catch (err) {
    console.error('Dropdown init error for id={{id}}', err);
}`})),
	PropsDef{
		// 主要屬性
		"label":         "",
		"id":            "",
		"name":          "",
		"options":       "",
		"defaultValue":  "",
		"placeholder":   "請選擇",
		"required":      false,
		"disabled":      false,
		"searchable":    false,
		"size":          "md",
		"fullWidth":     true,
		"helpText":      "",
		"errorText":     "",
		"labelPosition": "top",
		"color":         "#3b82f6",
	},
)
