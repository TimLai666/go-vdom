# Go VDOM 文檔中心

歡迎來到 go-vdom 的文檔中心！這裡包含了所有你需要了解的關於 go-vdom 的信息。

## 📚 文檔導航

### 🚀 快速開始

**[快速入門指南 (QUICK_START.md)](QUICK_START.md)**
- 5 分鐘快速上手
- 基本概念介紹
- 簡單示例
- 適合初學者

### 📖 核心文檔

**[完整技術文檔 (DOCUMENTATION.md)](DOCUMENTATION.md)**
- 深入的技術細節
- 架構設計說明
- 完整的功能介紹
- 進階用法指南
- 性能優化建議
- 故障排除

### 🔧 API 參考

**[JavaScript DSL API 參考 (API_REFERENCE.md)](API_REFERENCE.md)**
- 完整的 JavaScript DSL API 列表
- `Fn()` 和 `AsyncFn()` 使用指南
- `TryCatch()` 錯誤處理
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
- [Props 屬性系統](DOCUMENTATION.md#props-屬性) - 屬性類型和用法
- [DOM 操作](API_REFERENCE.md#dom-操作) - JavaScript DSL DOM 操作

### 組件系統

- [組件定義](DOCUMENTATION.md#組件系統) - 如何創建可重用組件
- [組件最佳實踐](DOCUMENTATION.md#最佳實踐) - 組件設計指南
- [UI 組件庫](DOCUMENTATION.md#ui-組件庫) - 內建組件使用

### JavaScript 和交互

- [JavaScript DSL 基礎](API_REFERENCE.md#核心函數) - Fn, AsyncFn, Call 等
- [事件處理](API_REFERENCE.md#事件處理) - onClick, onSubmit 等
- [異步操作](API_REFERENCE.md#異步操作) - async/await, TryCatch
- [Fetch API](API_REFERENCE.md#fetch-api-輔助函數) - HTTP 請求

### 控制流

- [條件渲染](DOCUMENTATION.md#條件渲染) - If/Then/Else
- [列表渲染](DOCUMENTATION.md#列表渲染) - For, Repeat
- [控制流最佳實踐](QUICK_REFERENCE.md#控制流) - 常見模式

### 進階功能

- [模板序列化](DOCUMENTATION.md#模板序列化) - JSON, Go template
- [錯誤處理](API_REFERENCE.md#錯誤處理) - TryCatch 詳解
- [性能優化](DOCUMENTATION.md#性能優化) - 優化技巧

---

## 🎯 學習路徑

### 初學者路徑

1. 閱讀 **[快速入門](QUICK_START.md)** 了解基本概念
2. 運行 `examples/01_basic_usage.go` 查看基本用法
3. 運行 `examples/02_components.go` 學習組件系統
4. 閱讀 **[完整文檔](DOCUMENTATION.md)** 的核心概念部分

### 中級開發者路徑

1. 閱讀 **[API 參考](API_REFERENCE.md)** 了解 JavaScript DSL
2. 運行 `examples/03_javascript_dsl.go` 學習 DOM 操作和事件處理
3. 學習 **AsyncFn** 和異步操作
4. 探索 **[UI 組件庫](DOCUMENTATION.md#ui-組件庫)**

### 高級開發者路徑

1. 運行 `examples/04_template_serialization.go` 學習模板序列化
2. 閱讀 **[進階用法](DOCUMENTATION.md#進階用法)**
3. 學習 **[性能優化](DOCUMENTATION.md#性能優化)**
4. 閱讀 **[最佳實踐](DOCUMENTATION.md#最佳實踐)**

---

## 🔥 重點功能

### AsyncFn - 異步函數 (v1.1.0 新增)

解決 "await is only valid in async functions" 錯誤：

```go
// ✅ 正確 - 使用 AsyncFn
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Alert("'Success!'"),
    ),
}, "Load Data")

// ❌ 錯誤 - 使用 Fn 會導致錯誤
Button(Props{
    "onClick": js.Fn(nil,
        js.Const("response", "await fetch('/api/data')"), // 錯誤！
    ),
}, "Load Data")
```

詳見：[API 參考 - AsyncFn](API_REFERENCE.md#asyncfnparams-string-actions-jsaction-jsaction)

### Props 類型系統 (v1.1.0 更新)

支持任意類型的值：

```go
Props{
    "class": "btn",           // string
    "disabled": true,         // bool
    "count": 42,              // int
    "price": 19.99,           // float64
    "onClick": js.Fn(...),    // JSAction
}
```

詳見：[完整文檔 - Props 類型系統](DOCUMENTATION.md#props-類型系統)

### TryCatch - 錯誤處理

完整的異步錯誤處理：

```go
js.TryCatch(
    js.AsyncFn(nil,
        // 異步操作
        js.Const("data", "await fetchData()"),
    ),
    js.Ptr(js.Fn(nil,
        // 錯誤處理
        js.Log("'Error:', e.message"),
    )),
    nil,
)
```

詳見：[API 參考 - TryCatch](API_REFERENCE.md#trycatchbaseaction-jsaction-catchfn-jsaction-finallyfn-jsaction-jsaction)

---

## 💡 常見問題速查

| 問題 | 參考文檔 |
|------|---------|
| 如何開始使用 go-vdom？ | [快速入門](QUICK_START.md) |
| await 語法錯誤 | [API 參考 - AsyncFn](API_REFERENCE.md#asyncfnparams-string-actions-jsaction-jsaction) |
| 如何創建組件？ | [完整文檔 - 組件系統](DOCUMENTATION.md#組件系統) |
| 如何處理表單？ | [API 參考 - 事件處理](API_REFERENCE.md#事件處理) |
| 如何發送 API 請求？ | [API 參考 - Fetch API](API_REFERENCE.md#fetch-api-輔助函數) |
| 如何條件渲染？ | [快速參考 - 控制流](QUICK_REFERENCE.md#控制流) |
| Props 支持哪些類型？ | [完整文檔 - Props 類型系統](DOCUMENTATION.md#props-類型系統) |
| 如何優化性能？ | [完整文檔 - 性能優化](DOCUMENTATION.md#性能優化) |

---

## 🔗 相關資源

- **[主 README](../README.md)** - 項目概述和基本信息
- **[CHANGELOG](../CHANGELOG.md)** - 版本更新歷史
- **[示例代碼](../examples/)** - 可運行的示例程序
- **[GitHub 倉庫](https://github.com/TimLai666/go-vdom)** - 源代碼和 Issues

---

## 📝 文檔貢獻

發現文檔問題或想要改進文檔？

1. 在 [GitHub Issues](https://github.com/TimLai666/go-vdom/issues) 報告問題
2. 提交 Pull Request 改進文檔
3. 建議新增內容或示例

---

**版本**: v1.1.0  
**最後更新**: 2025-01-24