package components

import (
	. "github.com/TimLai666/go-vdom/vdom"
)

// TableComponent 現代化表格組件
//
// 提供高度美觀和功能性的數據表格，適合展示結構化數據。
//
// 參數:
//   - stripped: 是否顯示條紋，預設 "false"
//   - bordered: 是否顯示邊框，預設 "false"
//   - hoverable: 是否顯示懸停效果，預設 "true"
//   - small: 是否使用緊湊布局，預設 "false"
//   - responsive: 是否響應式，預設 "true"
//   - fullWidth: 是否填滿容器寬度，預設 "true"
//   - header: 表頭內容，格式為 HTML 字符串如 "<tr><th>名稱</th><th>價格</th></tr>"
//   - footer: 表尾內容，預設為空
//   - highlightColor: 高亮色，預設現代藍 "#3b82f6"
//   - borderColor: 邊框顏色，預設 "#e5e7eb"
//   - headerBg: 表頭背景色，預設 "#f9fafb"
//   - stripped: 條紋背景色，預設 "#f9fafb"
//   - radius: 圓角大小，可選 "none"、"sm"、"md"、"lg"，預設 "md"
//   - shadow: 陰影強度，可選 "none"、"sm"、"md"、"lg"，預設 "sm"
//
// 用法:
//
//	TableComponent(Props{
//	  "hoverable": "true",
//	  "stripped": "true",
//	  "header": "<tr><th>名稱</th><th>價格</th><th>庫存</th></tr>"
//	},
//	  Tr(
//	    Td("iPhone"),
//	    Td("$999"),
//	    Td("有貨"),
//	  ),
//	  Tr(
//	    Td("iPad"),
//	    Td("$799"),
//	    Td("無貨"),
//	  ),
//	)
var TableComponent = Component(
	Div(
		Props{
			"style": `
				width: 100%;
				overflow-x: {{overflowX}};
				border-radius: {{tableBorderRadius}};
				box-shadow: {{tableBoxShadow}};
			`,
		},
		Table(
			Props{
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
					"style": `
						background-color: {{headerBgColor}};
						color: #1e293b;
						vertical-align: bottom;
						border-bottom: 2px solid {{borderColor}};
						
						tr {
							border-style: solid;
							border-width: 0;
							border-color: {{borderColor}};
						}
						
						th {
							padding: {{thPadding}};
							font-weight: 600;
							text-align: left;
							vertical-align: bottom;
							border-right-style: solid;
							border-right-width: {{rightBorderWidth}};
							border-right-color: {{borderColor}};
							border-bottom-style: solid;
							border-bottom-width: {{bottomBorderWidth}};
							border-bottom-color: {{borderColor}};
						}
						
						th:last-child {
							border-right: {{lastColumnRightBorder}};
						}
					`,
				},
				"{{header}}",
			),
			Tbody(
				Props{
					"style": `
						vertical-align: inherit;
						
						tr {
							border-style: solid;
							border-width: 0;
							border-color: {{borderColor}};
						}
						
						tr:last-child {
							border-bottom: {{lastRowBorder}};
						}
						
						tr:nth-child(even) {
							background-color: {{evenRowBgColor}};
						}
						
						tr:hover {
							background-color: {{hoverBgColor}};
						}
						
						td {
							padding: {{tdPadding}};
							vertical-align: middle;
							border-right-style: solid;
							border-right-width: {{rightBorderWidth}};
							border-right-color: {{borderColor}};
							border-bottom-style: solid;
							border-bottom-width: {{bottomBorderWidth}};
							border-bottom-color: {{borderColor}};
						}
						
						td:last-child {
							border-right: {{lastColumnRightBorder}};
						}
					`,
				},
				"{{children}}",
			),
			Tfoot(
				Props{
					"style": `
						display: {{tfootDisplay}};
						background-color: {{headerBgColor}};
						color: #1e293b;
						vertical-align: bottom;
						border-top: 2px solid {{borderColor}};
						
						tr {
							border-style: solid;
							border-width: 0;
							border-color: {{borderColor}};
						}
						
						td {
							padding: {{thPadding}};
							font-weight: 600;
							text-align: left;
							vertical-align: bottom;
							border-right-style: solid;
							border-right-width: {{rightBorderWidth}};
							border-right-color: {{borderColor}};
						}
						
						td:last-child {
							border-right: {{lastColumnRightBorder}};
						}
					`,
				},
				"{{footer}}",
			),
		),
	),
	PropsDef{
		// 主要參數
		"stripped":       "false",   // 是否顯示條紋
		"bordered":       "false",   // 是否顯示邊框
		"hoverable":      "true",    // 是否顯示懸停效果
		"small":          "false",   // 是否使用緊湊布局
		"responsive":     "true",    // 是否響應式
		"fullWidth":      "true",    // 是否填滿容器寬度
		"header":         "",        // 表頭內容
		"footer":         "",        // 表尾內容
		"highlightColor": "#3b82f6", // 高亮色
		"borderColor":    "#e5e7eb", // 邊框顏色
		"headerBg":       "#f9fafb", // 表頭背景色
		"stripeColor":    "#f9fafb", // 條紋背景色
		"radius":         "md",      // 圓角大小
		"shadow":         "sm",      // 陰影強度

		// 計算屬性
		"overflowX":             "auto",
		"tableWidth":            "100%",
		"tableBorder":           "none",
		"thPadding":             "0.75rem 1rem",
		"tdPadding":             "0.75rem 1rem",
		"headerBgColor":         "#f9fafb",
		"evenRowBgColor":        "transparent",
		"hoverBgColor":          "rgba(0, 0, 0, 0.02)",
		"rightBorderWidth":      "0",
		"bottomBorderWidth":     "1px",
		"tableBorderRadius":     "0.5rem",
		"tableBoxShadow":        "0 1px 3px 0 rgba(0, 0, 0, 0.1)",
		"lastColumnRightBorder": "none",
		"lastRowBorder":         "none",
		"tfootDisplay":          "none",
	},
)
