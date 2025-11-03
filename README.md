# go-vdom

[![Go Version](https://img.shields.io/badge/Go-1.24.1+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`go-vdom` 是一個用 Go 語言實現的虛擬 DOM (Virtual DOM) 庫，專為服務器端渲染 HTML 和 JavaScript 而設計。它提供了一套完整的 DSL（Domain Specific Language），讓開發者能夠在 Go 中以聲明式、類型安全的方式構建動態網頁。

## 功能特性

- ✨ **虛擬 DOM 實現**: 高效的 DOM 操作和渲染機制
- 🧩 **組件系統**: 支持可重用的 UI 組件，類似於現代前端框架
- 🔀 **控制流**: 內建 `If/Then/Else`、`Repeat` 和 `For` 等控制結構
- 📝 **JavaScript DSL**: 完整的 JavaScript 代碼生成器，支持同步/異步函數、事件處理、API 調用
- ⚡ **代碼最小化**: 自動最小化生成的 JavaScript，減少 30-50% 體積
- 🎨 **UI 組件庫**: 提供常見的 UI 組件（按鈕、輸入框、下拉選單等）
- 🖥️ **服務器端渲染**: 直接在 Go HTTP 服務器中生成完整的 HTML 文檔
- 🎯 **類型安全**: 利用 Go 的類型系統確保代碼正確性
- 🚀 **Bootstrap 集成**: 內建支持 Bootstrap CSS 框架
- 🌐 **API 集成**: 支持 Fetch API 用於 GET 和 POST 請求
- ⚡ **高性能**: 零運行時依賴，純靜態 HTML/JS 生成
- 🔄 **模板序列化**: 支持導出/導入 VNode 為 Go template、JSON 格式
- 📦 **類型靈活**: Props 支持任意類型值，自動類型轉換
- ⚡ **異步支持**: JavaScript DSL 完整支持 async/await 語法

## 安裝

```bash
go get github.com/TimLai666/go-vdom
```

## ⚠️ 重要說明

**事件處理器**：事件處理器必須使用 `js.Do()` 或 `js.AsyncDo()`，不要使用 `js.Fn()` 或 `js.AsyncFn()`。

```go
// ❌ 舊方式 (不再有效)
"onClick": js.Fn(nil, js.Alert("'Hello'")),

// ✅ 新方式 (正確)
"onClick": js.Do(nil, js.Alert("'Hello'")),                    // 同步事件（不使用 event）
"onClick": js.Do([]string{"event"}, ...),                       // 使用 event 對象時必須聲明
"onClick": js.AsyncDo(nil, js.Alert("'Hello'")),               // 異步事件（不使用 event）
"onClick": js.AsyncDo([]string{"event"}, ...),                 // 異步且使用 event
```

⚠️ **重要**：

- 當需要使用 `event` 對象時（如 `event.target`、`event.preventDefault()`），必須聲明為參數
- 參數名可以任意命名：`[]string{"event"}`、`[]string{"e"}`、`[]string{"evt"}`、`[]string{"myEvent"}` 等都可以
- 只要宣告了參數，就會自動注入 event 對象

詳見：

- [API 參考文檔](docs/API_REFERENCE.md) - 完整 API 說明
- [快速入門指南](docs/QUICK_START.md) - 基礎教學
- [examples/09_event_handlers.go](examples/09_event_handlers.go) - 實例演示
- [examples/10_do_with_params.go](examples/10_do_with_params.go) - 參數使用示例

## 快速開始

```go
package main

import (
    "fmt"
    "net/http"

    js "github.com/TimLai666/go-vdom/jsdsl"
    . "github.com/TimLai666/go-vdom/dom"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        doc := Document(
            "我的網頁",
            nil, nil, nil,
            Div(Props{"class": "container"},
                H1("Hello, go-vdom!"),
                Button(Props{
                    "onClick": js.AsyncDo(nil,
                        js.Const("response", "await fetch('/api/data')"),
                        js.Const("data", "await response.json()"),
                        js.Alert("'Data loaded: ' + JSON.stringify(data)"),
                    ),
                }, "Load Data"),
            ),
        )

        fmt.Fprint(w, Render(doc))
    })

    http.ListenAndServe(":8080", nil)
}
```

## 項目結構

```
go-vdom/
├── components/          # UI 組件庫
│   ├── checkbox.go
│   ├── dropdown.go
│   ├── radio.go
│   ├── switch.go
│   └── textfield.go
├── control/             # 控制流結構
│   └── control.go       # If/Then/Else/Repeat/For
├── jsdsl/               # JavaScript DSL 生成器
│   ├── jsdsl.go         # 核心 DSL 函數（Fn, AsyncFn, TryCatch 等）
│   └── builder.go       # JSAction 建構器
├── dom/                 # 核心 DOM 實現
│   ├── types.go         # VNode、Props 類型定義
│   ├── tags.go          # HTML 標籤函數
│   ├── render.go        # 渲染邏輯
│   ├── component.go     # 組件系統
│   ├── template.go      # 模板序列化（JSON、Go template）
│   └── template_test.go # 單元測試
├── runtime/             # 運行時支持
├── examples/            # 示例代碼
│   ├── 01_basic_usage.go
│   ├── 02_components.go
│   ├── 03_javascript_dsl.go
│   └── 04_template_serialization.go
├── docs/                # 完整文檔
│   ├── QUICK_START.md      # 快速入門
│   ├── DOCUMENTATION.md    # 完整技術文檔
│   ├── API_REFERENCE.md    # JavaScript DSL API 參考
│   └── QUICK_REFERENCE.md  # 語法速查表
├── main.go              # 完整示例應用
├── go.mod
├── CHANGELOG.md
├── LICENSE
└── README.md
```

## 核心概念

### 虛擬 DOM (VNode)

```go
// 創建元素
Div(Props{"class": "container"},
    H1("標題"),
    P("段落內容"),
)
```

### 組件系統

```go
Card := Component(
    Div(Props{"class": "card"},
        H2("{{title}}"),
        P("{{content}}"),
    ),
    nil,
    PropsDef{
        "title": "",
        "content": "",
        "visible": true,  // 布林值
        "priority": 1,    // 整數
    },
)

// 使用組件
Card(Props{"title": "我的卡片", "content": "卡片內容"})
```

### 控制流

```go
import "github.com/TimLai666/go-vdom/control"

// 條件渲染
control.If(isLoggedIn,
    control.Then(Div("歡迎回來")),
    control.Else(Div("請登入")),
)

// 列表渲染
control.For(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
})
```

### JavaScript DSL

```go
import js "github.com/TimLai666/go-vdom/jsdsl"

// 同步函數
js.Fn(nil,
    js.Log("'Hello'"),
    js.Alert("'World'"),
)

// 異步函數（支持 await）
js.AsyncFn(nil,
    js.Const("response", "await fetch('/api')"),
    js.Const("data", "await response.json()"),
    js.Log("data"),
)

// 錯誤處理 - 流暢 API（Try 生成純粹的 try-catch-finally）
js.AsyncFn(nil,
    js.Try(
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).Finally(
        js.Log("'清理完成'"),
    ),
)

// 或使用 AsyncDo 立即執行
js.AsyncDo(
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).End(),
)
```

## 文檔

完整文檔位於 `docs/` 目錄：

- **[快速入門](docs/QUICK_START.md)** - 5 分鐘上手指南
- **[完整文檔](docs/DOCUMENTATION.md)** - 深入技術文檔
- **[API 參考](docs/API_REFERENCE.md)** - JavaScript DSL 完整 API
- **[快速參考](docs/QUICK_REFERENCE.md)** - 語法速查表
- **[Try-Catch-Finally 指南](docs/TRY_CATCH_FINALLY.md)** - 錯誤處理完整說明
- **[代碼優化指南](docs/OPTIMIZATION.md)** - 最小化和性能優化

## 重要更新 (v1.2.0)

### 1. Try-Catch-Finally 與 Do/AsyncDo

全新設計的錯誤處理和 IIFE 創建 API：

**Try-Catch-Finally** - 生成純粹的 try-catch-finally 語句（不包裝在函數中）：

```go
// 在 AsyncFn 中使用（支持 await）
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).End(),
)

// Try-Catch-Finally
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).Finally(
        js.Log("'清理完成'"),
    ),
)
```

**Do / AsyncDo** - 創建立即執行函數（IIFE）：

```go
// Do - 普通 IIFE
js.Do(
    js.Const("x", "1"),
    js.Log("x"),
)

// AsyncDo - 異步 IIFE
js.AsyncDo(
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)
```

**設計理念：**

- ✅ Try 生成純粹的 try-catch-finally，不包裝
- ✅ 需要 async 時，用 AsyncFn 或 AsyncDo 包裝
- ✅ Do/AsyncDo 專門用於創建 IIFE
- ✅ 職責分離清晰，更靈活
- ✅ 錯誤對象統一命名為 `error`

詳細說明請參考 [Try-Catch-Finally 指南](docs/TRY_CATCH_FINALLY.md)

### 2. JavaScript 代碼最小化

所有生成的 JavaScript 代碼自動最小化：

```go
js.AsyncFn(nil,
    js.Const("x", "1"),
    js.Const("y", "2"),
    js.Log("x+y"),
)

// 生成：async()=>{const x=1;const y=2;console.log(x+y)}
// 而非：async () => {\n  const x = 1;\n  const y = 2;\n  console.log(x+y);\n}
```

**優勢：**

- ✅ 減少 30-50% 的代碼體積
- ✅ 加快頁面載入速度
- ✅ 降低帶寬消耗
- ✅ 無需配置，自動應用

### 3. Const/Let 支持 JSAction

`Const` 和 `Let` 現在可以接受 `JSAction` 類型：

```go
// 傳入字符串（舊方式，仍然支持）
js.Const("x", "1")

// 傳入 JSAction（新方式）
js.Const("random", js.Call("Math.random"))
js.Const("doubled", JSAction{Code: "x * 2"})
js.Const("data", js.Ident("response.data"))
```

**優勢：**

- ✅ 更靈活的值賦值
- ✅ 更好的代碼組合
- ✅ 減少字符串拼接
- ✅ 向後兼容

詳細說明請參考 [代碼優化指南](docs/OPTIMIZATION.md)

### AsyncFn - 異步函數支持

可以使用 `AsyncFn` 創建支持 `await` 的異步函數：

```go
Button(Props{
    "onClick": js.AsyncFn(nil,  // 使用 AsyncFn 而非 Fn
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Alert("'Success!'"),
    ),
}, "Fetch Data")
```

**重要：** 任何包含 `await` 的代碼都必須使用 `AsyncFn`！

### Props 類型系統

Props 現在支持任意類型，會自動轉換為 HTML 屬性：

```go
Props{
    "class": "btn",           // string
    "disabled": true,         // bool → "true"
    "count": 42,              // int → "42"
    "price": 19.99,           // float64 → "19.99"
    "onClick": js.Fn(...),    // JSAction

    // 複雜類型自動序列化為 JSON
    "data-items": []string{"apple", "banana", "orange"},
    // → '["apple","banana","orange"]'

    "data-config": map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
    },
    // → '{"fontSize":14,"theme":"dark"}'

    "data-user": User{Name: "John", Email: "john@example.com"},
    // → '{"Name":"John","Email":"john@example.com"}'
}
```

**複雜類型支持：**

- ✅ 陣列（Array/Slice）
- ✅ Map
- ✅ 結構體（Struct）
- ✅ 嵌套的複雜結構

這些複雜類型會自動序列化為 JSON 字符串，可在客戶端 JavaScript 中使用 `JSON.parse()` 解析。

### 模板序列化

支持導出/導入 VNode 為 JSON 或 Go template：

```go
// 導出為 JSON
jsonStr := ToJSON(vnode)

// 從 JSON 導入
vnode := FromJSON(jsonStr)

// 導出為 Go template
tmpl := ToGoTemplate(vnode)
```

## 運行示例

```bash
# 運行主示例（包含所有功能）
go run main.go
# 訪問 http://localhost:8080

# 運行單獨示例
go run examples/01_basic_usage.go          # http://localhost:8080
go run examples/02_components.go           # http://localhost:8081
go run examples/03_javascript_dsl.go       # http://localhost:8082
go run examples/04_template_serialization.go  # http://localhost:8083
go run examples/complex_props_demo.go      # http://localhost:8084 - 複雜 Props 示例
```

## 最佳實踐

1. **使用 AsyncFn 處理異步操作**

   ```go
   // ✅ 正確
   js.AsyncFn(nil, js.Const("data", "await fetch('/api')"))

   // ❌ 錯誤（會導致 "await is only valid in async functions" 錯誤）
   js.Fn(nil, js.Const("data", "await fetch('/api')"))
   ```

2. **始終使用 TryCatch 包裝異步操作**

   ```go
   js.TryCatch(
       js.AsyncFn(nil, /* 異步操作 */),
       js.Ptr(js.Fn(nil, /* 錯誤處理 */)),
       nil,
   )
   ```

3. **JavaScript 字符串需要單引號**

   ```go
   js.Log("'This is a string'")    // ✅ 正確
   js.Log("This is a string")      // ❌ 錯誤
   ```

4. **組件化重用代碼**
   ```go
   // 定義一次，多處使用
   MyCard := Component(template, nil, propsDef)
   ```

## 完整示例

查看 `main.go` 了解包含以下功能的完整應用：

- ✅ 基本 HTML 渲染
- ✅ 組件系統（卡片、表單等）
- ✅ 控制流（If/Repeat/For）
- ✅ 異步 API 調用（GET/POST）
- ✅ 錯誤處理（TryCatch）
- ✅ UI 組件庫（TextField, Dropdown, Checkbox 等）
- ✅ Bootstrap 集成

## 貢獻

歡迎提交 Issue 和 Pull Request！

## 許可證

MIT License - 詳見 [LICENSE](LICENSE) 文件

---

**作者**: TimLai666
**倉庫**: https://github.com/TimLai666/go-vdom
