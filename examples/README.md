# Go VDOM 示例集合

這個目錄包含了各種 go-vdom 的使用示例，從基礎到進階，幫助你快速上手。

## 示例列表

### 01_basic_usage.go - 基本用法
展示如何創建簡單的 HTML 頁面，包括：
- Document 函數的使用
- 基本 HTML 元素
- Bootstrap 樣式集成
- 頁面佈局和結構

**運行方式:**
```bash
go run examples/01_basic_usage.go
```
訪問: http://localhost:8080

---

### 02_components.go - 組件系統
展示如何創建和使用可重用的組件，包括：
- Alert 組件
- Card 組件
- Badge 組件
- Button 組件
- UserCard 組件
- 組件組合和嵌套

**運行方式:**
```bash
go run examples/02_components.go
```
訪問: http://localhost:8081

---

### 03_javascript_dsl.go - JavaScript DSL
展示如何使用 JavaScript DSL 創建交互式頁面，包括：
- DOM 操作（SetText, SetHTML, AddClass, RemoveClass）
- 變數定義（Let, Const）
- 事件處理（OnClick）
- 表單處理
- 動態創建元素
- Try/Catch 錯誤處理
- DomReady 初始化

**運行方式:**
```bash
go run examples/03_javascript_dsl.go
```
訪問: http://localhost:8082

---

### 04_template_serialization.go - 模板序列化
展示如何將 VNode 序列化為不同格式，包括：
- JSON 序列化和反序列化
- Go template 格式導出
- 模板變數提取
- VNode 克隆
- Props 合併

**運行方式:**
```bash
go run examples/04_template_serialization.go
```
訪問: http://localhost:8083

---

### 05_foreach_usage.go - ForEach 使用
展示後端和前端的列表渲染方法，包括：
- Go 後端 ForEach 和 ForEachWithIndex
- JavaScript 前端 ForEachJS 和 ForEachWithIndexJS
- 複雜對象渲染
- 動態 DOM 元素創建
- API 數據遍歷（異步）
- ForEachElement DOM 操作

**運行方式:**
```bash
go run examples/05_foreach_usage.go
```
訪問: http://localhost:8084

---

### 06_control_loops.go - 控制流和循環
展示各種控制流結構的使用，包括：
- If/Then/Else 條件渲染
- Repeat 重複渲染
- For 範圍循環
- 複雜嵌套結構

**運行方式:**
```bash
go run examples/06_control_loops.go
```
訪問: http://localhost:8085

---

### 07_trycatch_usage.go - Try/Catch/Finally
展示 v1.2.0 新的 Try/Catch/Finally API，包括：
- Try-Catch 基本用法
- Try-Catch-Finally 完整錯誤處理
- Try-Finally 資源清理
- Do() 和 AsyncDo() IIFE 使用
- 異步錯誤處理
- 嵌套 Try 語句

**運行方式:**
```bash
go run examples/07_trycatch_usage.go
```
訪問: http://localhost:8086

---

### 08_minified_js.go - JavaScript 最小化
展示 v1.2.0 的代碼最小化功能，包括：
- 自動去除空白和換行
- Const/Let 支持 JSAction
- 代碼體積優化（30-50% 減少）
- 可讀性與體積的對比

**運行方式:**
```bash
go run examples/08_minified_js.go
```
訪問: http://localhost:8087

---

### 09_event_handlers.go - 事件處理器 (v1.2.1)
**⭐ 重要：v1.2.1 新增**

展示事件處理器的正確使用方式，包括：
- js.Do() 同步事件處理器
- js.AsyncDo() 異步事件處理器
- 複雜異步操作（API 調用）
- DOM 操作和計數器
- 多種事件類型（click, mouseover, input 等）
- 表單事件處理
- 錯誤處理最佳實踐

**運行方式:**
```bash
go run examples/09_event_handlers.go
```
訪問: http://localhost:8089

**注意**: 此示例展示 v1.2.1 的重要變更 - 事件處理器必須使用 `js.Do()` 或 `js.AsyncDo()`，不再支援 `js.Fn()` 或 `js.AsyncFn()`。

---

## 快速開始

1. 確保已安裝 Go 1.24.1 或更高版本
2. 克隆倉庫並進入項目目錄
3. 運行任意示例文件
4. 在瀏覽器中訪問對應的端口

## 學習路徑

建議按以下順序學習示例：

1. **01_basic_usage.go** - 了解基本概念和結構
2. **02_components.go** - 學習組件化開發
3. **03_javascript_dsl.go** - 掌握 JavaScript DSL 和交互功能
4. **04_template_serialization.go** - 學習模板序列化
5. **05_foreach_usage.go** - 掌握列表渲染
6. **06_control_loops.go** - 學習控制流結構
7. **07_trycatch_usage.go** - 掌握錯誤處理
8. **08_minified_js.go** - 了解代碼優化
9. **09_event_handlers.go** - ⭐ 掌握 v1.2.1 事件處理器（重要）

## 自定義示例

你可以基於這些示例創建自己的應用：

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    . "github.com/TimLai666/go-vdom/vdom"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        
        doc := Document(
            "我的應用",
            nil, nil, nil,
            Div(
                Props{"class": "container"},
                H1("歡迎使用 go-vdom"),
            ),
        )
        
        fmt.Fprint(w, Render(doc))
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 更多資源

- [完整文檔](../DOCUMENTATION.md)
- [主 README](../README.md)
- [GitHub 倉庫](https://github.com/TimLai666/go-vdom)

## 版本特定說明

### v1.2.1 重要變更
從 v1.2.1 開始，事件處理器的使用方式已變更：

```go
// ❌ 舊方式 (不再有效)
"onClick": js.Fn(nil, js.Alert("'Hello'")),

// ✅ 新方式 (正確)
"onClick": js.Do(js.Alert("'Hello'")),           // 同步
"onClick": js.AsyncDo(js.Alert("'Hello'")),      // 異步
```

請參閱 **09_event_handlers.go** 了解詳細用法。

## 提示

- 所有示例使用不同的端口（8080-8089），可以同時運行多個示例
- 打開瀏覽器控制台（F12）查看 JavaScript 日誌輸出
- 每個示例都是獨立的，可以單獨運行
- 示例代碼包含詳細註釋，便於理解
- 從 example 09 開始學習 v1.2.1 的新特性

## 需要幫助？

如果遇到問題：
1. 查看 [DOCUMENTATION.md](../DOCUMENTATION.md) 中的詳細文檔
2. 查看 [故障排除](../DOCUMENTATION.md#故障排除) 章節
3. 在 GitHub 上提交 Issue

---

**Happy Coding! 🚀**