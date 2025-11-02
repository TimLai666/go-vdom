package components

import (
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
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
						/* visually hidden when not active, but remains in layout for accessibility/interaction */
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
			Span(
				Props{"style": `
						font-size: {{fontSize}};
						color: #374151;
						display: ${'{{label}}'.trim() ? 'inline' : 'none'};
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
	),
	// onDOMReady: 初始化單一 Checkbox 的互動行為（保持原先邏輯）
	jsdsl.Fn(nil, JSAction{Code: `try {
				const boxId = '{{id}}';
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
						/* show checkmark by making it visible and opaque */
						checkmark.style.visibility = 'visible';
						checkmark.style.opacity = '1';
					} else {
						box.style.borderColor = '#d1d5db';
						box.style.background = 'white';
						/* hide checkmark visually but keep it in layout for accessibility */
						checkmark.style.opacity = '0';
						checkmark.style.visibility = 'hidden';
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
			} catch (err) {
				console.error('Checkbox init error for id={{id}}', err);
			}`}),
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
				"id":        "checkbox-group-{{id}}",
				"data-name": "{{name}}",
				"style": `
					display: flex;
					flex-direction: {{flexDirection}};
					gap: {{gap}};
				`,
			},
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
	// onDOMReady: 初始化 CheckboxGroup（由原先的 onmount 移入，使用 jsdsl.Fn 以確保是函數表達式）
	jsdsl.Fn(nil, JSAction{Code: `try {
					const rawId = '{{id}}';
					const rawName = '{{name}}';
					let container = null;
					// 優先使用明確提供的 id（checkbox-group-<id>）
					if (rawId && rawId.trim()) {
						container = document.getElementById('checkbox-group-' + rawId);
					}

					// 如果沒有 id，使用 name（如果提供），並允許透過 data-name 或 class 選取
					if (!container && rawName && rawName.trim()) {
						container = document.getElementById('checkbox-group-' + rawName) ||
									document.querySelector('[data-name=\"' + rawName + '\"]') ||
									document.querySelector('.checkbox-group-options-' + rawName);
					}

					if (!container) return;

					// 解析 options / values
					const rawOptions = '{{options}}';
					const rawValues = '{{values}}';
					const options = rawOptions ? rawOptions.split(',').map(s=>s.trim()).filter(Boolean) : [];
					const values = rawValues ? rawValues.split(',').map(s=>s.trim()).filter(Boolean) : [];
					const required = '{{required}}' === 'true';
					const disabled = '{{disabled}}' === 'true';
					const color = '{{color}}';
					const colorRgb = '{{colorRgb}}';

					// 清空舊內容（保證 id 可重複使用或重新渲染）
					container.innerHTML = '';

					options.forEach(function(opt, idx){
						const val = opt;
						const safeName = (rawName && rawName.trim()) ? rawName : (rawId && rawId.trim() ? rawId : 'checkboxgroup');
						const id = 'checkbox-' + safeName + '-' + idx;

						// label 容器
						const label = document.createElement('label');
						label.className = 'checkbox-label';
						label.style.cssText = 'display:flex; align-items:center; cursor:pointer; user-select:none; gap:0.5rem;';

						// 原生 input（視覺隱藏但仍可聚焦與互動）
						const input = document.createElement('input');
						input.type = 'checkbox';
						input.id = id;
						input.name = safeName;
						input.value = val;
						input.required = required;
						input.disabled = disabled;
						if (values.indexOf(val) !== -1) {
							input.checked = true;
						}
						input.style.cssText = 'position:absolute; opacity:0; height:1px; width:1px; margin:-1px; padding:0; border:0; overflow:hidden; clip:rect(0 0 0 0); white-space:nowrap;';

						// 裝飾性方塊
						const box = document.createElement('span');
						box.className = 'checkbox-box';
						box.style.cssText = 'display:inline-flex; align-items:center; justify-content:center; width: {{checkboxSize}}; height: {{checkboxSize}}; border-radius: {{borderRadius}}; border: 2px solid #d1d5db; margin-right: 0.5rem; transition: all 0.2s ease; background: white;';

						// 勾選符號
						const checkmark = document.createElement('span');
						checkmark.className = 'checkbox-checkmark';
						checkmark.style.cssText = 'visibility:hidden; opacity:0; transition: opacity 150ms ease-in-out, visibility 150ms step-end; width: calc({{checkboxSize}} / 2 - 1px); height: calc({{checkboxSize}} / 2 + 1px); border: solid white; border-width: 0 2px 2px 0; transform: rotate(45deg) translate(-2px, -2px);';

						// 文本
						const text = document.createElement('span');
						text.className = 'checkbox-label-text';
						text.textContent = opt;
						text.style.cssText = 'font-size: {{fontSize}}; color: #374151;';

						// 組裝節點順序： input（隱藏） -> 裝飾 box (內含 checkmark) -> 文本
						label.appendChild(input);
						label.appendChild(box);
						box.appendChild(checkmark);
						label.appendChild(text);
						container.appendChild(label);

						// 更新 UI 狀態的函數
						function updateState() {
							if (input.checked) {
								box.style.borderColor = color;
								box.style.background = color;
								checkmark.style.visibility = 'visible';
								checkmark.style.opacity = '1';
							} else {
								box.style.borderColor = '#d1d5db';
								box.style.background = 'white';
								checkmark.style.opacity = '0';
								checkmark.style.visibility = 'hidden';
							}
							if (input.disabled) {
								label.style.cursor = 'not-allowed';
							} else {
								label.style.cursor = 'pointer';
							}
						}

						// 綁定事件
						input.addEventListener('change', function() {
							updateState();
							// 發出群組事件（包含目前所有被選中的值）
							const selected = Array.from(container.querySelectorAll('input[type=\"checkbox\"]:checked')).map(i=>i.value);
							container.dispatchEvent(new CustomEvent('checkbox-group:change', {
								detail: { name: safeName, values: selected },
								bubbles: true
							}));
						});

						input.addEventListener('focus', function() {
							if (!input.disabled) {
								box.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
							}
						});
						input.addEventListener('blur', function() {
							box.style.boxShadow = 'none';
						});

						// 初始化狀態
						updateState();
					});
				} catch (err) {
					console.error('CheckboxGroup init error for ' + (rawName || rawId), err);
				}`}),
	PropsDef{
		// 主要屬性
		"id":        "",         // 群組 ID（可選，用於明確指定容器）
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
