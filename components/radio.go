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
			"class": "radio-group",
			"style": `
				margin-bottom: 1.25rem;
				width: 100%;
			`,
		},
		Div(
			Props{"class": "radio-group-label",
				"style": `					display: ${'{{label}}'.trim() ? 'block' : 'none'};
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
				"class": "radio-group-options", "style": `
					display: flex;
					flex-direction: ${'{{direction}}' === 'horizontal' ? 'row' : 'column'};
					gap: {{gap}};
				`,
				"onmount": `
					(() => {
						const container = document.querySelector('.radio-group-options');
						if (!container) return;						const name = '{{name}}';
						const options = '{{options}}'.split(',').filter(opt => opt.trim());
						const defaultValue = '{{defaultValue}}';
						const required = '{{required}}' === 'true';
						const disabled = '{{disabled}}' === 'true';
						const color = '{{color}}';
						const colorRgb = '{{colorRgb}}';
						
						options.forEach((option, index) => {
							const id = 'radio-' + name + '-' + index;
							const label = document.createElement('label');
							label.className = 'radio-label';
							label.htmlFor = id;
							label.style.cssText = 'display: flex; align-items: center; cursor: pointer; user-select: none;';
							
							const input = document.createElement('input');
							input.type = 'radio';
							input.id = id;
							input.name = name;
							input.value = option.trim();
							input.required = required;
							input.disabled = disabled;
							input.checked = option.trim() === defaultValue;
							input.style.cssText = 'position: absolute; opacity: 0; height: 0; width: 0; display: none;';
							
							const circle = document.createElement('span');
							circle.className = 'radio-circle';
							circle.style.cssText = 'display: inline-flex; align-items: center; justify-content: center; width: 1.25rem; height: 1.25rem; border-radius: 50%; border: 2px solid #d1d5db; margin-right: 0.5rem; transition: all 0.2s ease; background: white;';
							
							const dot = document.createElement('span');
							dot.className = 'radio-dot';
							dot.style.cssText = 'width: 0.75rem; height: 0.75rem; border-radius: 50%; background: ' + color + '; display: none;';
							
							const text = document.createElement('span');
							text.textContent = option.trim();
							text.style.cssText = 'font-size: 0.9375rem; color: #374151;';
							
							label.appendChild(input);
							label.appendChild(circle);
							label.appendChild(dot);
							label.appendChild(text);
							container.appendChild(label);
							
							// 事件處理
							input.addEventListener('change', () => {
								const allCircles = container.querySelectorAll('.radio-circle');
								const allDots = container.querySelectorAll('.radio-dot');
								
								allCircles.forEach(c => {
									c.style.borderColor = '#d1d5db';
									c.style.background = 'white';
								});
								allDots.forEach(d => d.style.display = 'none');
								
								if (input.checked) {
									circle.style.borderColor = color;
									circle.style.background = 'white';
									dot.style.display = 'block';
									
									// 觸發自定義事件
									container.dispatchEvent(new CustomEvent('radio-group:change', {
										detail: {
											name: name,
											value: input.value
										},
										bubbles: true
									}));
								}
							});
							
							input.addEventListener('focus', () => {
								if (!input.disabled) {
									circle.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
								}
							});
							
							input.addEventListener('blur', () => {
								circle.style.boxShadow = 'none';
							});
							
							// 初始化選中狀態
							if (input.checked) {
								circle.style.borderColor = color;
								dot.style.display = 'block';
							}
							
							// 禁用狀態
							if (disabled) {
								label.style.cursor = 'not-allowed';
								circle.style.borderColor = '#e5e7eb';
								circle.style.background = '#f9fafb';
								if (input.checked) {
									dot.style.background = '#d1d5db';
								}
							}
						});
					})();
				`,
			},
		),
		Div(
			Props{"class": "radio-group-help-text",
				"style": `					display: ${'{{helpText}}'.trim() ? 'block' : 'none'};
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
		"labelSize":   "0.9375rem",
		"helpDisplay": "none",
		"helpColor":   "#64748b",
		"gap":         "0.75rem",
		"colorRgb":    "59, 130, 246",
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
					display: none;
				`,
			},
		),
		Span(Props{
			"class": "radio-circle",
			"style": `
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
					`,
		}),
		Span(
			Props{
				"class": "radio-dot",
				"style": `
						display: none;
						width: calc({{radioSize}} - 10px);
						height: calc({{radioSize}} - 10px);
						border-radius: 50%;
						background: {{color}};
					`,
			},
		),
		Span(
			Props{"style": `
					font-size: {{fontSize}};
					color: #374151;
				`},
			"{{labelText}}",
		),
		Script(`
			document.addEventListener('DOMContentLoaded', function() {				const id = {{id}};
				const input = document.getElementById(id);
				if (!input) return; // 確保元素存在
				const circle = input.nextElementSibling;
				const dot = circle.nextElementSibling;
				function updateRadioState() {
					if (input.checked) {
						circle.style.borderColor = '{{color}}';
						circle.style.background = 'white';
						dot.style.display = 'block';
					} else {
						circle.style.borderColor = '#d1d5db';
						circle.style.background = 'white';
						dot.style.display = 'none';
					}
					
					if (input.disabled) {
						circle.style.borderColor = '#e5e7eb';
						circle.style.background = '#f9fafb';
						if (input.checked) {
							dot.style.background = '#d1d5db';
						}
						circle.style.cursor = 'not-allowed';
					} else {
						circle.style.cursor = 'pointer';
						if (input.checked) {
							dot.style.background = '{{color}}';
						}
					}
				}
				
				// 初始化狀態
				updateRadioState();
				
				// 切換狀態時更新
				input.addEventListener('change', updateRadioState);
				
				// Focus 效果
				input.addEventListener('focus', function() {
					if (!this.disabled) {
						circle.style.boxShadow = '0 0 0 3px rgba({{colorRgb}}, 0.15)';
					}
				});
				
				input.addEventListener('blur', function() {
					circle.style.boxShadow = 'none';
				});
				
				// 觸發自定義事件
				input.addEventListener('change', function() {
					this.dispatchEvent(new CustomEvent('radio:change', {
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
