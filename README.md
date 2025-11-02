# go-vdom

[![Go Version](https://img.shields.io/badge/Go-1.24.1+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`go-vdom` 是一個用 Go 語言實現的虛擬 DOM (Virtual DOM) 庫，專為服務器端渲染 HTML 和 JavaScript 而設計。它允許開發者在 Go 中以聲明式的方式構建動態網頁，無需手動編寫 HTML 和 JavaScript 代碼。

## 功能特性

- **虛擬 DOM 實現**: 提供高效的 DOM 操作和渲染機制。
- **組件系統**: 支持可重用的 UI 組件，類似於現代前端框架。
- **控制流**: 內建 `If/Then/Else`、`Repeat` 和 `For` 等控制結構，用於動態內容渲染。
- **JavaScript DSL**: 內建 JavaScript 代碼生成器，支持事件處理、API 調用等。
- **UI 組件庫**: 提供常見的 UI 組件，如按鈕、輸入框、下拉選單等。
- **服務器端渲染**: 直接在 Go HTTP 服務器中生成完整的 HTML 文檔。
- **Bootstrap 集成**: 內建支持 Bootstrap CSS 框架。
- **API 集成**: 支持 Fetch API 用於 GET 和 POST 請求。

## 安裝

確保您的系統已安裝 Go 1.24.1 或更高版本。

```bash
go get github.com/TimLai666/go-vdom
```

或者克隆倉庫：

```bash
git clone https://github.com/TimLai666/go-vdom.git
cd go-vdom
go mod tidy
```

## 快速開始

運行示例應用：

```bash
go run main.go
```

然後在瀏覽器中訪問 `http://localhost:8080` 查看演示頁面。

## 項目結構

```
go-vdom/
├── components/     # UI 組件庫
├── control/        # 控制流結構 (If, Repeat, For 等)
├── jsdsl/          # JavaScript DSL 生成器
├── vdom/           # 核心虛擬 DOM 實現
├── main.go         # 示例應用入口
├── go.mod          # Go 模塊定義
├── LICENSE         # MIT 許可證
└── README.md       # 本文件
```

## 使用指南

### 基本用法

```go
package main

import (
    . "github.com/TimLai666/go-vdom/vdom"
    . "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
    // 創建一個簡單的 HTML 文檔
    doc := Document(
        "我的頁面",
        nil, // 無外部鏈接
        nil, // 無腳本
        nil, // 無 meta
        Div(
            Props{"class": "container"},
            H1("歡迎使用 go-vdom"),
            P("這是一個簡單的示例"),
        ),
    )

    // 渲染為 HTML 字符串
    html := Render(doc)
    fmt.Println(html)
}
```

### 組件使用

```go
// 定義一個組件
Card := Component(
    Div(
        Props{"class": "card"},
        H2("{{title}}"),
        P("{{content}}"),
    ),
    PropsDef{"title": "", "content": ""},
)

// 使用組件
card := Card(Props{"title": "標題", "content": "內容"})
```

### 控制流

```go
import control "github.com/TimLai666/go-vdom/control"

// 條件渲染
show := true
content := control.If(show,
    control.Then(Div("顯示內容")),
    control.Else(Div("隱藏內容")),
)

// 循環渲染
items := []string{"項目1", "項目2", "項目3"}
list := Ul(control.For(items, func(item string, i int) VNode {
    return Li(item)
}))
```

### JavaScript 集成

```go
script := Script(Props{"type": "module"}, DomReady(
    El("#myButton").OnClick(
        Alert("按鈕被點擊了！"),
    ),
))
```

### UI 組件

```go
import comp "github.com/TimLai666/go-vdom/components"

// 文字輸入框
textField := comp.TextField(Props{
    "id": "username",
    "label": "用戶名",
    "placeholder": "請輸入用戶名",
})

// 下拉選單
dropdown := comp.Dropdown(Props{
    "id": "country",
    "label": "國家",
    "options": "台灣,中國,日本",
})
```

## API 參考

### 核心類型

- `VNode`: 虛擬 DOM 節點
- `Props`: 屬性映射
- `PropsDef`: 組件屬性定義

### 主要函數

- `Document(title, links, scripts, metas, body)`: 創建完整 HTML 文檔
- `Component(template, propsDef)`: 定義可重用組件
- `Render(vnode)`: 將虛擬 DOM 渲染為 HTML 字符串

### HTML 元素函數

所有標準 HTML 元素都有對應的函數，如 `Div`, `H1`, `P`, `A` 等。

### 控制流函數

- `control.If(condition, then, else)`: 條件渲染
- `control.Repeat(count, fn)`: 重複渲染
- `control.For(items, fn)`: 遍歷渲染

### JavaScript DSL

- `DomReady(actions...)`: 文檔就緒時執行的動作
- `El(selector)`: 選擇 DOM 元素
- `FetchRequest(url, options...)`: 發送 HTTP 請求

## 示例應用

倉庫中的 `main.go` 是一個完整的示例應用，展示了：

- 完整的 HTML 文檔生成
- 組件使用
- 控制流
- JavaScript 事件處理
- API 調用 (GET 和 POST)
- UI 組件展示
- Bootstrap 樣式集成

運行示例：

```bash
go run main.go
```

訪問 `http://localhost:8080` 查看效果。

## 貢獻

歡迎貢獻！請遵循以下步驟：

1. Fork 此倉庫
2. 創建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 許可證

本項目採用 MIT 許可證 - 詳見 [LICENSE](LICENSE) 文件。

## 聯繫

- 作者: TimLai666
- GitHub: [https://github.com/TimLai666/go-vdom](https://github.com/TimLai666/go-vdom)

---

**注意**: 此庫仍在開發中，API 可能會發生變化。請在使用前檢查兼容性。
