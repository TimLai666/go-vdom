package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
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
				margin-bottom: 1rem;
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
				},
				Span(
					Props{
						"class": "checkbox-checkmark",
						"style": `
							visibility: hidden;
							opacity: 0;
							transition: opacity 150ms ease-in-out, visibility 150ms step-end;
							width: calc({{checkboxSize}} / 2 - 1px);
							height: calc({{checkboxSize}} / 2 + 1px);
							border: solid white;
							border-width: 0 2px 2px 0;
							transform: rotate(45deg) translate(-2px, -2px);
						`,
					},
				),
			),
			Span(
				Props{
					"style": `
						font-size: {{fontSize}};
						color: #374151;
					`,
				},
				"{{label}}",
			),
		),
		Div(
			Props{
				"style": `
					display: {{helpDisplay}};
					font-size: 0.875rem;
					margin-top: 0.375rem;
					color: #64748b;
					margin-left: calc({{checkboxSize}} + 0.5rem);
				`,
			},
			"{{helpText}}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
		const input = document.getElementById('{{id}}');
		if (!input) return;

		const box = input.nextElementSibling;
		const checkmark = box.querySelector('.checkbox-checkmark');
		if (!box || !checkmark) return;

		function updateState() {
			const checked = input.checked;
			const disabled = input.disabled;

			if (checked) {
				box.style.borderColor = '{{color}}';
				box.style.background = '{{color}}';
				checkmark.style.visibility = 'visible';
				checkmark.style.opacity = '1';
			} else {
				box.style.borderColor = '#d1d5db';
				box.style.background = 'white';
				checkmark.style.visibility = 'hidden';
				checkmark.style.opacity = '0';
			}

			if (disabled) {
				box.style.borderColor = '#e5e7eb';
				box.style.background = '#f9fafb';
				box.style.cursor = 'not-allowed';
				if (checked) {
					box.style.background = '#d1d5db';
				}
			} else {
				box.style.cursor = 'pointer';
			}
		}

		// 初始化狀態
		updateState();

		// 點擊 box 時切換狀態
		box.addEventListener('click', function(e) {
			e.preventDefault();
			if (!input.disabled) {
				input.checked = !input.checked;
				updateState();
				input.dispatchEvent(new Event('change', { bubbles: true }));
			}
		});

		// 監聽 input 的 change 事件
		input.addEventListener('change', function() {
			updateState();
			this.dispatchEvent(new CustomEvent('checkbox:change', {
				detail: { id: '{{id}}', checked: this.checked, value: this.value },
				bubbles: true
			}));
		});

		// Focus 效果
		input.addEventListener('focus', function() {
			if (!this.disabled) {
				box.style.boxShadow = '0 0 0 3px rgba({{colorRgb}}, 0.15)';
			}
		});

		input.addEventListener('blur', function() {
			box.style.boxShadow = 'none';
		});
	`})),
	PropsDef{
		"id":           "",
		"name":         "",
		"value":        "",
		"label":        "",
		"checked":      false,
		"required":     false,
		"disabled":     false,
		"size":         "md",
		"helpText":     "",
		"color":        "#3b82f6",
		"checkboxSize": "1.25rem",
		"fontSize":     "0.9375rem",
		"borderRadius": "0.25rem",
		"helpDisplay":  "none",
		"colorRgb":     "59, 130, 246",
	},
)

// CheckboxGroup 勾選框組
//
// 提供一組相關的勾選框，適合多項選擇場景。
//
// 參數:
//   - id: 群組 ID，預設自動生成
//   - name: 勾選框組名稱，預設為空
//   - label: 標籤文字，預設為空
//   - options: 選項清單，以逗號分隔，如 "選項1,選項2,選項3"
//   - values: 已選中值，以逗號分隔，預設為空
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - direction: 排列方向，可選 "horizontal"、"vertical"，預設 "vertical"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - helpText: 幫助文字，預設為空
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
			Props{
				"style": `
					display: {{labelDisplay}};
					margin-bottom: 0.75rem;
					font-weight: 500;
					font-size: 0.9375rem;
					color: #374151;
				`,
			},
			"{{label}}",
		),
		Div(
			Props{
				"id":        "checkbox-group-{{id}}",
				"data-name": "{{name}}",
				"style": `
					display: flex;
					flex-direction: {{flexDirection}};
					gap: 0.75rem;
				`,
			},
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
		const container = document.getElementById('checkbox-group-{{id}}');
		if (!container) return;

		const options = '{{options}}'.split(',').map(s => s.trim()).filter(Boolean);
		const values = '{{values}}'.split(',').map(s => s.trim()).filter(Boolean);
		const name = '{{name}}' || '{{id}}';
		const disabled = '{{disabled}}' === 'true';
		const color = '{{color}}';
		const colorRgb = '{{colorRgb}}';

		options.forEach(function(option, index) {
			const id = 'checkbox-{{id}}-' + index;
			const isChecked = values.indexOf(option) !== -1;

			// 創建 label 容器
			const label = document.createElement('label');
			label.style.cssText = 'display:inline-flex; align-items:center; cursor:pointer; user-select:none;';

			// 創建隱藏的 input
			const input = document.createElement('input');
			input.type = 'checkbox';
			input.id = id;
			input.name = name;
			input.value = option;
			input.checked = isChecked;
			input.disabled = disabled;
			input.style.cssText = 'position:absolute; opacity:0; height:1px; width:1px; margin:-1px; padding:0; border:0; overflow:hidden; clip:rect(0 0 0 0); white-space:nowrap;';

			// 創建裝飾 box
			const box = document.createElement('span');
			box.className = 'checkbox-box';
			box.style.cssText = 'display:inline-flex; align-items:center; justify-content:center; width:1.25rem; height:1.25rem; border-radius:0.25rem; border:2px solid #d1d5db; margin-right:0.5rem; transition:all 0.2s ease; background:white;';

			// 創建 checkmark
			const checkmark = document.createElement('span');
			checkmark.className = 'checkbox-checkmark';
			checkmark.style.cssText = 'visibility:hidden; opacity:0; transition:opacity 150ms ease-in-out, visibility 150ms step-end; width:calc(1.25rem / 2 - 1px); height:calc(1.25rem / 2 + 1px); border:solid white; border-width:0 2px 2px 0; transform:rotate(45deg) translate(-2px, -2px);';

			// 創建文本
			const text = document.createElement('span');
			text.textContent = option;
			text.style.cssText = 'font-size:0.9375rem; color:#374151;';

			// 組裝
			box.appendChild(checkmark);
			label.appendChild(input);
			label.appendChild(box);
			label.appendChild(text);
			container.appendChild(label);

			// 更新狀態函數
			function updateState() {
				if (input.checked) {
					box.style.borderColor = color;
					box.style.background = color;
					checkmark.style.visibility = 'visible';
					checkmark.style.opacity = '1';
				} else {
					box.style.borderColor = '#d1d5db';
					box.style.background = 'white';
					checkmark.style.visibility = 'hidden';
					checkmark.style.opacity = '0';
				}

				if (input.disabled) {
					box.style.cursor = 'not-allowed';
					label.style.cursor = 'not-allowed';
				} else {
					box.style.cursor = 'pointer';
					label.style.cursor = 'pointer';
				}
			}

			// 初始化
			updateState();

			// 點擊 box 切換狀態
			box.addEventListener('click', function(e) {
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
				const selected = Array.from(container.querySelectorAll('input[type="checkbox"]:checked')).map(i => i.value);
				container.dispatchEvent(new CustomEvent('checkbox-group:change', {
					detail: { name: name, values: selected },
					bubbles: true
				}));
			});

			// Focus 效果
			input.addEventListener('focus', function() {
				if (!this.disabled) {
					box.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
				}
			});

			input.addEventListener('blur', function() {
				box.style.boxShadow = 'none';
			});
		});
	`})),
	PropsDef{
		"id":            "",
		"name":          "",
		"label":         "",
		"options":       "",
		"values":        "",
		"required":      false,
		"disabled":      false,
		"direction":     "vertical",
		"size":          "md",
		"helpText":      "",
		"color":         "#3b82f6",
		"flexDirection": "column",
		"labelDisplay":  "none",
		"helpDisplay":   "none",
		"colorRgb":      "59, 130, 246",
	},
)
