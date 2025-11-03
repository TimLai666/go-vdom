package components

import (
	. "github.com/TimLai666/go-vdom/dom"
	jsdsl "github.com/TimLai666/go-vdom/jsdsl"
)

// TableComponent 現代化表格組件
//
// 提供高度美觀和功能性的數據表格，適合展示結構化數據。
var TableComponent = Component(
	Div(
		Props{
			"class": "table-container",
			"style": `
				width: 100%;
				overflow-x: {{overflowX}};
				border-radius: {{tableBorderRadius}};
				box-shadow: {{tableBoxShadow}};
			`,
		},
		Table(
			Props{
				"class":                "modern-table",
				"data-highlight-color": "{{highlightColor}}",
				"data-hoverable":       "{{hoverable}}",
				"style": `
					width: {{tableWidth}};
					border-collapse: separate;
					border-spacing: 0;
					font-size: 0.9375rem;
					border: {{tableBorder}};
					overflow: hidden;
				`,
			},
			Thead(
				Props{
					"class": "table-header",
					"style": `
						background-color: {{headerBgColor}};
						color: #1e293b;
						vertical-align: bottom;
						border-bottom: 2px solid {{borderColor}};
					`,
				},
				"{{header}}",
			),
			Tbody(
				Props{
					"class": "table-body",
					"style": `
						vertical-align: inherit;
					`,
				},
				"{{children}}",
			),
			Tfoot(
				Props{
					"class": "table-footer",
					"style": `
						display: ${{{footer}}.trim() ? 'table-footer-group' : 'none'};
						background-color: {{headerBgColor}};
						color: #1e293b;
						vertical-align: bottom;
						border-top: 2px solid {{borderColor}};
					`,
				},
				"{{footer}}",
			),
			// helper sorting function kept as inline Script node (it becomes part of DOM output)
			Script(`
				function sortTable(table, columnIndex, ascending) {
					const tbody = table.querySelector('tbody');
					const rows = Array.from(tbody.querySelectorAll('tr'));
					const sortedRows = rows.sort((a, b) => {
						const aValue = a.cells[columnIndex].textContent.trim();
						const bValue = b.cells[columnIndex].textContent.trim();

						// 判斷是否為數字
						const aNum = Number(aValue);
						const bNum = Number(bValue);

						if (!isNaN(aNum) && !isNaN(bNum)) {
							return ascending ? aNum - bNum : bNum - aNum;
						}

						return ascending
							? aValue.localeCompare(bValue)
							: bValue.localeCompare(aValue);
					});

					tbody.innerHTML = '';
					sortedRows.forEach(row => tbody.appendChild(row));
				}
			`),
		),
	),
	// onDOMReady: 初始化表格互動（hover, row click, striped rows, sorting）
	jsdsl.Ptr(jsdsl.Fn(nil, JSAction{Code: `try {
		const table = document.querySelector('.modern-table');
		if (!table) return;

		const rows = table.querySelectorAll('tbody tr');
		const isHoverable = {{hoverable}} === true;
		const highlightColor = {{highlightColor}};

		rows.forEach(row => {
			if (isHoverable) {
				row.addEventListener('mouseenter', () => {
					row.style.backgroundColor = 'rgba(0, 0, 0, 0.02)';
					row.style.transition = 'background-color 0.2s ease';
				});
				row.addEventListener('mouseleave', () => {
					if (row.classList.contains('even-row')) {
						row.style.backgroundColor = '{{evenRowBgColor}}';
					} else {
						row.style.backgroundColor = 'transparent';
					}
				});
			}

			// row click event: emit detailed event on table
			row.addEventListener('click', () => {
				table.dispatchEvent(new CustomEvent('table:row-click', {
					detail: {
						rowIndex: Array.from(rows).indexOf(row),
						rowData: Array.from(row.cells).map(cell => cell.textContent.trim())
					},
					bubbles: true
				}));
			});

			// apply striped style for even rows if enabled
			if (row.rowIndex % 2 === 0) {
				row.classList.add('even-row');
				if ({{stripped}} === true) {
					row.style.backgroundColor = {{evenRowBgColor}};
				}
			}
		});

		// add sorting handlers to headers with .sortable
		const headers = table.querySelectorAll('th');
		headers.forEach((header, index) => {
			if (header.classList.contains('sortable')) {
				header.style.cursor = 'pointer';
				header.addEventListener('click', () => {
					const isAscending = header.classList.toggle('asc');
					try {
						sortTable(table, index, isAscending);
					} catch (err) {
						console.error('sortTable error', err);
					}
					table.dispatchEvent(new CustomEvent('table:sort', {
						detail: {
							columnIndex: index,
							isAscending: isAscending
						},
						bubbles: true
					}));
				});
			}
		});
	} catch (err) {
		console.error('Table init error', err);
	}`})),
	PropsDef{
		// 主要屬性
		"stripped":       false,     // 是否顯示條紋
		"bordered":       false,     // 是否顯示邊框
		"hoverable":      true,      // 是否顯示懸停效果
		"small":          false,     // 是否使用緊湊布局
		"responsive":     true,      // 是否響應式
		"fullWidth":      true,      // 是否填滿容器寬度
		"header":         "",        // 表頭內容
		"footer":         "",        // 表尾內容
		"highlightColor": "#3b82f6", // 高亮色
		"borderColor":    "#e5e7eb", // 邊框顏色
		"headerBg":       "#f9fafb", // 表頭背景色
		"stripeColor":    "#f9fafb", // 條紋背景色
		"radius":         "md",      // 圓角大小
		"shadow":         "sm",      // 陰影強度

		// 計算屬性 (這些是固定值，保留用於一致性)
		"overflowX":         "auto",
		"tableWidth":        "100%",
		"tableBorder":       "none",
		"thPadding":         "0.75rem 1rem",
		"tdPadding":         "0.75rem 1rem",
		"headerBgColor":     "#f9fafb",
		"evenRowBgColor":    "#fbfdff",
		"tableBorderRadius": "0.5rem",
		"tableBoxShadow":    "0 1px 2px rgba(0,0,0,0.04)",
	},
)
