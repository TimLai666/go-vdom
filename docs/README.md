# Go VDOM 文檔中心

歡迎來到 go-vdom 的文檔中心！這裡包含了所有你需要了解的關於 go-vdom 的信息。

## 📚 文檔導航

### 🚀 快速開始

**[快速入門指南 (QUICK_START.md)](QUICK_START.md)**

- 5 分鐘快速上手
- 基本概念介紹
- 簡單示例
- v1.1.0 新特性介紹
- 適合初學者

### 📖 核心文檔

**[完整技術文檔 (DOCUMENTATION.md)](DOCUMENTATION.md)**

- 深入的技術細節
- 核心功能介紹
- 組件系統詳解
- JavaScript DSL 完整指南
- 模板表達式系統
- 控制流和模板序列化
- 最佳實踐和常見問題

### 🔧 API 參考

**[JavaScript DSL API 參考 (API_REFERENCE.md)](API_REFERENCE.md)**

- 完整的 JavaScript DSL API 列表
- `Fn()` 和 `AsyncFn()` 使用指南
- `Try()` 錯誤處理（流暢 API）
- DOM 操作函數
- Fetch API 集成
- 實用工具函數
- 類型定義
- 完整示例

### ⚡ 快速參考

**[語法速查表 (QUICK_REFERENCE.md)](QUICK_REFERENCE.md)**

- 常用語法快速查詢
- HTML 元素列表
- Props 用法
- 組件定義
- 控制流語法
- JavaScript DSL 常用函數
- 適合作為速查工具

---

## 📋 按主題查找

### HTML 和 DOM

- [HTML 元素創建](DOCUMENTATION.md#html-元素) - 基本元素用法
- [Props 屬性系統](DOCUMENTATION.md#props屬性系統) - 屬性類型和用法
- [DOM 操作](API_REFERENCE.md#dom-操作) - JavaScript DSL DOM 操作

### 組件系統

- [組件創建](DOCUMENTATION.md#創建組件) - 如何創建可重用組件
- [模板表達式](DOCUMENTATION.md#模板表達式) - 模板插值和條件表達式
- [內建組件](DOCUMENTATION.md#內建-ui-組件) - 完整的 UI 組件庫
- [組件最佳實踐](DOCUMENTATION.md#組件設計) - 組件設計指南

### JavaScript 和交互

- [JavaScript DSL 基礎](API_REFERENCE.md#核心函數) - Fn, AsyncFn, Call 等
- [事件處理](API_REFERENCE.md#事件處理) - onClick, onSubmit 等
- [異步操作](API_REFERENCE.md#異步操作) - async/await, Try-Catch-Finally
- [Fetch API](API_REFERENCE.md#fetch-api-輔助函數) - HTTP 請求

### 控制流

- [條件渲染](DOCUMENTATION.md#條件渲染) - If/Then/Else
- [列表渲染](DOCUMENTATION.md#列表渲染) - For, Repeat
- [控制流最佳實踐](QUICK_REFERENCE.md#控制流) - 常見模式

### 進階功能

- [模板表達式系統](DOCUMENTATION.md#模板表達式) - 服務器端條件表達式
- [模板序列化](DOCUMENTATION.md#模板序列化) - JSON, Go template
- [錯誤處理](API_REFERENCE.md#錯誤處理) - Try-Catch-Finally 詳解
- [性能優化](DOCUMENTATION.md#性能優化) - 優化技巧

---

## 🎯 學習路徑

### 初學者路徑

1. 閱讀 **[快速入門](QUICK_START.md)** 了解基本概念
2. 運行 `examples/01_basic_usage.go` 查看基本用法
3. 運行 `examples/02_components.go` 學習組件系統
4. 閱讀 **[完整文檔](DOCUMENTATION.md)** 的核心功能部分

### 中級開發者路徑

1. 閱讀 **[API 參考](API_REFERENCE.md)** 了解 JavaScript DSL
2. 運行 `examples/03_javascript_dsl.go` 學習 DOM 操作和事件處理
3. 學習 **AsyncFn** 和異步操作
4. 探索 **[內建組件庫](DOCUMENTATION.md#內建-ui-組件)**

### 高級開發者路徑

1. 運行 `examples/04_template_serialization.go` 學習模板序列化
2. 學習 **[模板表達式系統](DOCUMENTATION.md#模板表達式)**
3. 閱讀 **[性能優化](DOCUMENTATION.md#性能優化)**
4. 閱讀 **[最佳實踐](DOCUMENTATION.md#最佳實踐)**

---

## 🔥 重點功能

### Props 類型系統 (v1.1.0 新增)

支持任意類型的值，無需手動轉換：

```go
Props{
    "class": "container",    // 字符串
    "disabled": true,        // 布爾值 - 自動處理
    "width": 800,            // 整數 - 自動轉換
    "opacity": 0.8,          // 浮點數 - 自動轉換
    "onClick": js.Fn(...),   // JSAction
}
```

詳見：[完整文檔 - Props 屬性系統](DOCUMENTATION.md#props屬性系統)

### AsyncFn - 異步函數

解決 "await is only valid in async functions" 錯誤：

```go
// ✅ 正確 - 使用 AsyncFn
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Alert("'Success!'"),
    ),
}, Text("Load Data"))

// ❌ 錯誤 - 使用 Fn 會導致錯誤
Button(Props{
    "onClick": js.Fn(nil,
        js.Const("response", "await fetch('/api/data')"), // 錯誤！
    ),
}, Text("Load Data"))
```

詳見：[API 參考 - AsyncFn](API_REFERENCE.md#asyncfn)

### 模板表達式系統

在組件模板中使用條件表達式：

```go
Button(Props{
    "style": `
        background: ${'{{variant}}' === 'filled' ? '{{color}}' : 'transparent'};
        font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1.125rem' : '1rem'};
    `,
})
```

詳見：[完整文檔 - 模板表達式](DOCUMENTATION.md#模板表達式)

### Try-Catch-Finally - 錯誤處理

完整的異步錯誤處理：

```go
js.AsyncFn(nil,
    js.Try(
        // 異步操作
        js.Const("data", "await fetchData()"),
    ).Catch("error",
        // 錯誤處理
        js.Log("'Error:', error.message"),
    ).End(),
)
```

詳見：[API 參考 - Try-Catch-Finally](API_REFERENCE.md#try-catch-finally流暢-api)

### 內建 UI 組件庫

完整的表單和 UI 組件：

```go
import . "github.com/TimLai666/go-vdom/components"

// 按鈕
Btn(Props{"variant": "filled", "color": "#3b82f6"}, Text("按鈕"))

// 輸入框
TextField(Props{
    "label": "電子郵件",
    "type": "email",
    "icon": "📧",
})

// 下拉選單
Dropdown(Props{
    "label": "國家",
    "options": "台灣,日本,美國",
})

// 開關
Switch(Props{"label": "啟用通知", "checked": true})

// 更多組件...
```

詳見：[完整文檔 - 內建 UI 組件](DOCUMENTATION.md#內建-ui-組件)

---

## 💡 常見問題速查

| 問題                   | 參考文檔                                                                   |
| ---------------------- | -------------------------------------------------------------------------- |
| 如何開始使用 go-vdom？ | [快速入門](QUICK_START.md)                                                 |
| await 語法錯誤         | [API 參考 - AsyncFn](API_REFERENCE.md#asyncfn)                             |
| 如何創建組件？         | [完整文檔 - 組件系統](DOCUMENTATION.md#組件系統)                           |
| 如何處理表單？         | [API 參考 - 事件處理](API_REFERENCE.md#事件處理)                           |
| 如何發送 API 請求？    | [API 參考 - Fetch API](API_REFERENCE.md#fetch-api)                         |
| 如何條件渲染？         | [完整文檔 - 控制流](DOCUMENTATION.md#控制流)                               |
| Props 支持哪些類型？   | [完整文檔 - Props 屬性系統](DOCUMENTATION.md#props屬性系統)                |
| 如何優化性能？         | [完整文檔 - 性能優化](DOCUMENTATION.md#性能優化)                           |
| 如何處理錯誤？         | [API 參考 - Try-Catch-Finally](API_REFERENCE.md#try-catch-finally流暢-api) |
| 如何使用模板表達式？   | [完整文檔 - 模板表達式](DOCUMENTATION.md#模板表達式)                       |

---

## 🔗 相關資源

- **[主 README](../README.md)** - 項目概述和基本信息
- **[CHANGELOG](../CHANGELOG.md)** - 版本更新歷史
- **[示例代碼](../examples/)** - 可運行的示例程序
- **[GitHub 倉庫](https://github.com/TimLai666/go-vdom)** - 源代碼和 Issues

---

## 🎓 文檔結構

```
docs/
├── README.md              # 本文件 - 文檔導航
├── QUICK_START.md         # 快速入門指南
├── DOCUMENTATION.md       # 完整技術文檔
├── API_REFERENCE.md       # JavaScript DSL API 參考
└── QUICK_REFERENCE.md     # 語法速查表
```

---

## 📝 文檔貢獻

發現文檔問題或想要改進文檔？

1. 在 [GitHub Issues](https://github.com/TimLai666/go-vdom/issues) 報告問題
2. 提交 Pull Request 改進文檔
3. 建議新增內容或示例

---

**版本**: v1.1.0
**最後更新**: 2025-01-24
**維護者**: TimLai666
