package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
)

// RadioGroup 單選按鈕組
//
// 提供現代化的單選按鈕組，允許用戶在多個選項中選擇其一。
//
// 參數:
//   - id: 群組 ID，預設自動生成
//   - name: 單選按鈕組名稱，預設為空
//   - label: 標籤文字，預設為空
//   - options: 選項清單，以逗號分隔，如 "選項1,選項2,選項3"
//   - defaultValue: 預選值，預設為空
//   - required: 是否必填，預設 "false"
//   - disabled: 是否禁用，預設 "false"
//   - direction: 排列方向，可選 "horizontal"、"vertical"，預設 "vertical"
//   - size: 尺寸，可選 "sm"、"md"、"lg"，預設 "md"
//   - helpText: 幫助文字，預設為空
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
				display: ${'{{label}}'.trim() ? 'block' : 'none'};
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
				"style": `
				display: flex;
				flex-direction: ${'{{direction}}' === 'horizontal' ? 'row' : 'column'};
				gap: 0.75rem;
			`,
			},
			"{{children}}",
		),
		Div(
			Props{
				"style": `
				display: ${'{{helpText}}'.trim() ? 'block' : 'none'};
				font-size: 0.875rem;
				margin-top: 0.375rem;
				color: #64748b;
			`,
			},
			"{{helpText}}",
		),
	),
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
		const container = document.getElementById('radio-group-{{id}}');
		if (!container) return;

		const options = '{{options}}'.split(',').map(s => s.trim()).filter(Boolean);
		const defaultValue = '{{defaultValue}}'.trim();
		const name = '{{name}}' || '{{id}}';
		const disabled = '{{disabled}}' === 'true';
		const color = '{{color}}';
		const colorRgb = '{{colorRgb}}';

		options.forEach(function(option, index) {
			const id = 'radio-{{id}}-' + index;
			const isChecked = option === defaultValue;

			// 創建 label 容器
			const label = document.createElement('label');
			label.style.cssText = 'display:inline-flex; align-items:center; cursor:pointer; user-select:none;';

			// 創建隱藏的 input
			const input = document.createElement('input');
			input.type = 'radio';
			input.id = id;
			input.name = name;
			input.value = option;
			input.checked = isChecked;
			input.disabled = disabled;
			input.style.cssText = 'position:absolute; opacity:0; height:1px; width:1px; margin:-1px; padding:0; border:0; overflow:hidden; clip:rect(0 0 0 0); white-space:nowrap;';

			// 創建裝飾圓形
			const circle = document.createElement('span');
			circle.className = 'radio-circle';
			circle.style.cssText = 'display:inline-flex; align-items:center; justify-content:center; width:1.25rem; height:1.25rem; border-radius:50%; border:2px solid #d1d5db; margin-right:0.5rem; transition:all 0.2s ease; background:white;';

			// 創建中心點
			const dot = document.createElement('span');
			dot.className = 'radio-dot';
			dot.style.cssText = 'visibility:hidden; opacity:0; transition:opacity 150ms ease-in-out, visibility 150ms step-end; width:calc(1.25rem - 10px); height:calc(1.25rem - 10px); border-radius:50%; background:' + color + ';';

			// 創建文本
			const text = document.createElement('span');
			text.textContent = option;
			text.style.cssText = 'font-size:0.9375rem; color:#374151;';

			// 組裝
			circle.appendChild(dot);
			label.appendChild(input);
			label.appendChild(circle);
			label.appendChild(text);
			container.appendChild(label);

			// 更新狀態函數
			function updateState() {
				if (input.checked) {
					circle.style.borderColor = color;
					dot.style.visibility = 'visible';
					dot.style.opacity = '1';
				} else {
					circle.style.borderColor = '#d1d5db';
					dot.style.visibility = 'hidden';
					dot.style.opacity = '0';
				}

				if (input.disabled) {
					circle.style.cursor = 'not-allowed';
					label.style.cursor = 'not-allowed';
				} else {
					circle.style.cursor = 'pointer';
					label.style.cursor = 'pointer';
				}
			}

			// 初始化
			updateState();

			// 點擊 circle 選中此選項
			circle.addEventListener('click', function(e) {
				e.preventDefault();
				if (!input.disabled) {
					// 設置為 checked（瀏覽器會自動 uncheck 同名的其他 radio 並觸發 change）
					input.checked = true;
					// 手動觸發 change 事件（因為程式設置不會自動觸發）
					input.dispatchEvent(new Event('change', { bubbles: true }));
				}
			});

			// 監聽 change 事件 - 更新所有同組 radio 的 UI
			input.addEventListener('change', function() {
				// 更新所有同組 radio 的 UI（確保 UI 與實際狀態同步）
				container.querySelectorAll('input[type="radio"]').forEach(function(r) {
					if (r.name === name) {
						const rLabel = r.parentElement;
						const rCircle = rLabel ? rLabel.querySelector('.radio-circle') : null;
						const rDot = rLabel ? rLabel.querySelector('.radio-dot') : null;
						if (rCircle && rDot) {
							if (r.checked) {
								rCircle.style.borderColor = color;
								rDot.style.visibility = 'visible';
								rDot.style.opacity = '1';
							} else {
								rCircle.style.borderColor = '#d1d5db';
								rDot.style.visibility = 'hidden';
								rDot.style.opacity = '0';
							}
						}
					}
				});

				// 發送自定義事件
				const selected = container.querySelector('input[type="radio"]:checked');
				if (selected) {
					container.dispatchEvent(new CustomEvent('radio-group:change', {
						detail: { name: name, value: selected.value },
						bubbles: true
					}));
				}
			});

			// Focus 效果
			input.addEventListener('focus', function() {
				if (!this.disabled) {
					circle.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
				}
			});

			input.addEventListener('blur', function() {
				circle.style.boxShadow = 'none';
			});
		});
	`})),
	PropsDef{
		"id":           "",
		"name":         "",
		"label":        "",
		"options":      "",
		"defaultValue": "",
		"required":     false,
		"disabled":     false,
		"direction":    "vertical",
		"size":         "md",
		"helpText":     "",
		"color":        "#3b82f6",
		"colorRgb":     "59, 130, 246",
	},
)

// Radio 單選按鈕
//
// 提供單個 radio 按鈕，可以獨立使用或與其他 Radio 組合使用。
// 注意：通常建議使用 RadioGroup 而非手動組合多個 Radio。
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
				"type":     "radio",
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
			},
			Span(
				Props{
					"class": "radio-dot",
					"style": `
						visibility: hidden;
						opacity: 0;
						transition: opacity 150ms ease-in-out, visibility 150ms step-end;
						width: calc({{radioSize}} - 10px);
						height: calc({{radioSize}} - 10px);
						border-radius: 50%;
						background: {{color}};
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
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `
		const input = document.getElementById('{{id}}');
		if (!input) return;

		const circle = input.nextElementSibling;
		const dot = circle.querySelector('.radio-dot');
		if (!circle || !dot) return;

		const color = '{{color}}';
		const colorRgb = '{{colorRgb}}';

		function updateState() {
			if (input.checked) {
				circle.style.borderColor = color;
				dot.style.visibility = 'visible';
				dot.style.opacity = '1';
			} else {
				circle.style.borderColor = '#d1d5db';
				dot.style.visibility = 'hidden';
				dot.style.opacity = '0';
			}

			if (input.disabled) {
				circle.style.cursor = 'not-allowed';
			} else {
				circle.style.cursor = 'pointer';
			}
		}

		// 初始化
		updateState();

		// 點擊 circle 選中
		circle.addEventListener('click', function(e) {
			e.preventDefault();
			if (!input.disabled) {
				// 設置為 checked（瀏覽器會自動 uncheck 同名的其他 radio）
				input.checked = true;
				// 手動觸發 change 事件（因為程式設置不會自動觸發）
				input.dispatchEvent(new Event('change', { bubbles: true }));
			}
		});

		// 監聽 change 事件 - 更新所有同組 radio 的 UI
		input.addEventListener('change', function() {
			// 更新所有同名 radio 的 UI（確保 UI 與實際狀態同步）
			const radios = document.querySelectorAll('input[type="radio"][name="' + this.name + '"]');
			const color = '{{color}}';
			radios.forEach(function(r) {
				const rLabel = r.parentElement;
				if (rLabel) {
					const rCircle = rLabel.querySelector('.radio-circle');
					const rDot = rLabel.querySelector('.radio-dot');
					if (rCircle && rDot) {
						if (r.checked) {
							rCircle.style.borderColor = color;
							rDot.style.visibility = 'visible';
							rDot.style.opacity = '1';
						} else {
							rCircle.style.borderColor = '#d1d5db';
							rDot.style.visibility = 'hidden';
							rDot.style.opacity = '0';
						}
					}
				}
			});

			// 發送自定義事件
			this.dispatchEvent(new CustomEvent('radio:change', {
				detail: { id: '{{id}}', checked: this.checked, value: this.value },
				bubbles: true
			}));
		});

		// Focus 效果
		input.addEventListener('focus', function() {
			if (!this.disabled) {
				circle.style.boxShadow = '0 0 0 3px rgba(' + colorRgb + ', 0.15)';
			}
		});

		input.addEventListener('blur', function() {
			circle.style.boxShadow = 'none';
		});
	`})),
	PropsDef{
		"id":        "",
		"name":      "",
		"value":     "",
		"label":     "",
		"checked":   false,
		"required":  false,
		"disabled":  false,
		"size":      "md",
		"helpText":  "",
		"color":     "#3b82f6",
		"radioSize": "1.25rem",
		"fontSize":  "0.9375rem",
		"colorRgb":  "59, 130, 246",
	},
)
