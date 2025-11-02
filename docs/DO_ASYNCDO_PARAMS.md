# Do/AsyncDo 參數說明

## 概述

從 v1.2.1 開始，`js.Do()` 和 `js.AsyncDo()` 的簽名與 `js.Fn()` 和 `js.AsyncFn()` 保持一致，第一個參數為參數列表 `[]string`。

**重要更新**：Do/AsyncDo 現在具有智能參數檢測功能，可以自動判斷是否為事件處理器場景，讓同一個 API 支持多種用途。

## API 簽名

```go
// 同步 IIFE
func Do(params []string, actions ...JSAction) JSAction

// 異步 IIFE  
func AsyncDo(params []string, actions ...JSAction) JSAction

// 對比：函數定義
func Fn(params []string, actions ...JSAction) JSAction
func AsyncFn(params []string, actions ...JSAction) JSAction
```

## 智能參數檢測

Do/AsyncDo 會根據參數名自動判斷使用場景：

### 事件處理器模式
當參數名為常見的事件參數名時（`event`、`e`、`evt`、`ev`，不區分大小寫），自動傳入外部作用域的 `event` 對象：
- `js.Do([]string{"event"}, ...)` → `((event)=>{...})(event)`
- `js.Do([]string{"e"}, ...)` → `((e)=>{...})(event)`

### 通用 IIFE 模式
當參數名為其他名稱時，生成純 IIFE，不傳入任何參數：
- `js.Do([]string{"x"}, ...)` → `((x)=>{...})()`
- `js.Do([]string{"data"}, ...)` → `((data)=>{...})()`

這使得 Do/AsyncDo 既可以用於事件處理器，也可以用於創建獨立作用域、模塊化代碼等其他場景。

## 使用方式

### 1. 無參數

不需要任何參數時，傳入 `nil`：

```go
Button(Props{
    "onClick": js.Do(nil,
        js.Alert("'Hello!'"),
    ),
}, "Click me")
```

生成：
```javascript
onclick="(()=>{alert('Hi')})()"
```

### 2. 使用 event 對象

**重要**：在 IIFE 中使用 `event` 對象時，**必須**將其聲明為參數並傳遞：

```go
Button(Props{
    "onClick": js.Do([]string{"event"},
        js.Const("target", "event.target"),
        js.Const("text", "target.textContent"),
        js.Alert("'點擊了: ' + text"),
    ),
}, "Click me")
```

生成：
```javascript
onclick="((event)=>{const target=event.target;const text=target.textContent;alert('點擊了: '+text)})(event)"
```

**可以使用任何參數名**：

```go
// 使用 'e' 作為參數名
Button(Props{
    "onClick": js.Do([]string{"e"},
        js.Const("target", "e.target"),
        js.Alert("'點擊了'"),
    ),
}, "Click")
```

生成：
```javascript
onclick="((e)=>{const target=e.target;alert('點擊了')})(event)"
```

**關鍵點**：
- IIFE `(()=>{...})()` 創建了新的作用域
- 在新作用域內，外部的 `event` 對象不可訪問
- 必須通過參數傳遞：`((event)=>{...})(event)`
- 參數名為 `event`/`e`/`evt`/`ev` 時，自動傳入外部的 `event`
- 第一個 `event`/`e`/`evt` 是參數名，調用時的 `event` 是傳入的值
- 智能檢測基於參數名，讓 API 同時支持事件處理器和通用 IIFE 場景

### 3. 通用 IIFE（非事件場景）

當參數名不是事件相關時，生成純 IIFE，適用於創建獨立作用域、模塊化代碼等：

```go
// 創建獨立作用域，參數作為佔位符
js.Do([]string{"x", "y"},
    js.Const("x", "10"),
    js.Const("y", "20"),
    js.Const("sum", "x + y"),
    js.Log("sum"),
)
```

生成：
```javascript
((x,y)=>{const x=10;const y=20;const sum=x+y;console.log(sum)})()
```

**關鍵點**：
- 參數名為 `x`、`y`、`data` 等非事件名稱時，不會自動傳入 `event`
- 參數只是佔位符，在 IIFE 內部需要自己賦值
- 主要用於創建獨立作用域，避免變量污染全局空間
- 可用於模塊化代碼、初始化腳本等場景

## 與 Fn/AsyncFn 的對比

| 函數 | 第一參數 | 立即執行？ | 主要用途 |
|------|---------|-----------|---------|
| `Fn(params, ...)` | `[]string` | ❌ 否 | 定義可重用函數 |
| `AsyncFn(params, ...)` | `[]string` | ❌ 否 | 定義異步函數 |
| `Do(params, ...)` | `[]string` | ✅ 是 | 事件處理器（同步） |
| `AsyncDo(params, ...)` | `[]string` | ✅ 是 | 事件處理器（異步） |

## 常見模式

### 參數名示例

你可以使用任何喜歡的參數名：

```go
// 使用 'event'（最清晰）
js.Do([]string{"event"}, js.Const("id", "event.target.id"))

// 使用 'e'（最簡短）
js.Do([]string{"e"}, js.Const("id", "e.target.id"))

// 使用 'evt'（折衷方案）
js.Do([]string{"evt"}, js.Const("id", "evt.target.id"))

// 使用自定義名稱
js.Do([]string{"myEvent"}, js.Const("id", "myEvent.target.id"))
```

所有都會生成類似的代碼，只是參數名不同：
```javascript
((event)=>{...})(event)
((e)=>{...})(event)
((evt)=>{...})(event)
((myEvent)=>{...})(event)
```

**注意**：
- 事件參數名（event/e/evt/ev）：調用時傳入 `event`
- 其他參數名：不傳入任何參數，生成純 IIFE

### 按鈕點擊（不使用 event）
```go
"onClick": js.Do(nil,
    js.Alert("'Clicked!'"),
)
```

### 按鈕點擊（使用 event）
```go
"onClick": js.Do([]string{"event"},
    js.Const("btnId", "event.target.id"),
    js.Alert("'Clicked button: ' + btnId"),
)
```

### 表單提交
```go
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    js.Const("formData", "new FormData(event.target)"),
    js.Const("name", "formData.get('name')"),
    js.Alert("'提交: ' + name"),
)
```

### 輸入框變化
```go
"onInput": js.Do([]string{"event"},
    js.Const("value", "event.target.value"),
    JSAction{Code: "document.getElementById('output').textContent = value"},
)
```

### 異步 API 調用（不使用 event）
```go
"onClick": js.AsyncDo(nil,
    js.Const("response", "await fetch('/api/data')"),
    js.Const("data", "await response.json()"),
    js.Alert("'載入完成'"),
)
```

### 異步 API 調用（使用 event）
```go
"onClick": js.AsyncDo([]string{"event"},
    js.Const("btnId", "event.target.id"),
    js.Const("response", "await fetch('/api/data?btn=' + btnId)"),
    js.Const("data", "await response.json()"),
    js.Alert("'載入完成'"),
)
```

### 選擇框
```go
"onChange": js.Do([]string{"event"},
    js.Const("selected", "event.target.value"),
    js.Alert("'選擇了: ' + selected"),
)
```

## 為什麼需要 nil？

### 設計理念

1. **API 一致性**：`Fn()`、`AsyncFn()`、`Do()`、`AsyncDo()` 都使用相同的簽名
2. **明確性**：顯式表明是否有參數，而非隱式判斷
3. **可擴展性**：未來如果需要帶參數的 IIFE，API 已經支持

### 對比其他設計

❌ **自動判斷（已廢棄）**：
```go
// 舊設計：使用 interface{} 自動判斷
js.Do(action1, action2)              // 無參數
js.Do([]string{"x"}, action1)        // 有參數
```
問題：類型不明確，IDE 提示差，容易混淆

✅ **明確參數（當前設計）**：
```go
// 新設計：明確的參數列表
js.Do(nil, action1, action2)         // 無參數（明確）
js.Do([]string{"x"}, action1)        // 有參數（明確）
```
優點：類型明確，IDE 提示好，與其他函數一致

## 遷移指南

### 從 v1.2.0 遷移到 v1.2.1

**查找替換**：

1. 查找：`js.Do(`
   替換：`js.Do(nil,`

2. 查找：`js.AsyncDo(`
   替換：`js.AsyncDo(nil,`

**批量替換命令**（Linux/Mac）：
```bash
find . -name "*.go" -exec sed -i 's/js\.Do(/js.Do(nil,/g' {} \;
find . -name "*.go" -exec sed -i 's/js\.AsyncDo(/js.AsyncDo(nil,/g' {} \;
```

**批量替換命令**（Windows PowerShell）：
```powershell
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'js\.Do\(', 'js.Do(nil,' | Set-Content $_.FullName
    (Get-Content $_.FullName) -replace 'js\.AsyncDo\(', 'js.AsyncDo(nil,' | Set-Content $_.FullName
}
```

### 驗證遷移

編譯所有文件確保沒有錯誤：
```bash
go build ./...
```

## FAQ

**Q: 為什麼需要顯式聲明 event 參數？**  
A: 因為 IIFE 創建了新的作用域。外部的 `event` 對象在新作用域內不可訪問，必須通過參數傳遞進去。

**Q: 參數名一定要叫 'event' 嗎？**  
A: 不一定！你可以使用 `event`、`e`、`evt`、`ev`（不區分大小寫），它們都會被識別為事件參數並自動傳入 `event`。使用其他名稱（如 `x`、`data`）則生成通用 IIFE。

**Q: 什麼時候傳 nil，什麼時候傳 []string？**  
A: 
- 傳 `nil`：不需要任何參數時（如簡單的 alert、console.log）
- 傳 `[]string{"event"}`/`[]string{"e"}`：事件處理器中需要使用 event 對象時
- 傳 `[]string{"x", "y"}` 等：需要創建獨立作用域的通用 IIFE 時

**Q: 如何判斷是否會傳入 event？**  
A: 檢查第一個參數名（不區分大小寫）：
- `event`、`e`、`evt`、`ev` → 自動傳入 event
- 其他任何名稱 → 不傳入參數，生成純 IIFE

**Q: 什麼時候需要傳 []string？**  
A: 當你需要定義接受特定參數的 IIFE 時。這種情況較少見，主要用於需要參數化的立即執行代碼。

**Q: 能否讓 nil 作為默認值，省略這個參數？**  
A: 雖然可以通過函數重載實現，但會使 API 不一致（與 Fn/AsyncFn 不同），且降低代碼的明確性。

**Q: 為什麼生成的代碼是 ((event)=>{...})(event)？**  
A: 
- 第一個 `(event)` 是函數參數聲明
- 第二個 `(event)` 是調用時傳入的實參
- HTML 事件機制提供外部作用域的 `event`，我們將它傳入 IIFE

**Q: 推薦使用哪個參數名？**  
A: 
- **事件處理器**：
  - `event` - 最清晰明確（推薦用於複雜邏輯）
  - `e` - 最簡短，JavaScript 慣例（推薦用於簡單邏輯）
  - `evt` - 折衷方案
- **通用 IIFE**：
  - 使用有意義的名稱，如 `x`, `y`, `data`, `config` 等
  - 避免使用事件相關的名稱（除非確實是事件處理器）

**Q: 如果忘記聲明 event 參數會怎樣？**  
A: 會出現運行時錯誤：`Cannot read properties of undefined (reading 'preventDefault')` 或類似的錯誤，因為 IIFE 內部的 `event` 是 `undefined`。

**Q: 為什麼要設計智能檢測而不是總是傳入 event？**  
A: 因為 Do/AsyncDo 不僅用於事件處理器，還可用於創建獨立作用域、模塊化代碼、初始化腳本等多種場景。智能檢測讓同一個 API 支持多種用途，更加通用和靈活。

## 示例

完整示例請參考：
- [examples/09_event_handlers.go](../examples/09_event_handlers.go) - 事件處理器示例
- [examples/10_do_with_params.go](../examples/10_do_with_params.go) - 參數使用示例

## 使用場景總結

| 場景 | 參數 | 生成代碼 | 用途 |
|------|------|---------|------|
| 簡單操作 | `nil` | `(()=>{...})()` | alert、console.log |
| 事件處理器 | `[]string{"event"}` | `((event)=>{...})(event)` | onClick、onSubmit 等 |
| 事件處理器（簡短） | `[]string{"e"}` | `((e)=>{...})(event)` | 簡單的事件邏輯 |
| 獨立作用域 | `[]string{"x", "y"}` | `((x,y)=>{...})()` | 避免全局污染 |
| 模塊初始化 | `[]string{"module"}` | `((module)=>{...})()` | 模塊化代碼 |

## 相關文檔

- [EVENT_HANDLER_CHANGES.md](EVENT_HANDLER_CHANGES.md) - 事件處理器變更說明
- [EVENT_HANDLER_QUICK_REF.md](EVENT_HANDLER_QUICK_REF.md) - 快速參考
- [examples/10_do_with_params.go](../examples/10_do_with_params.go) - 參數使用示例
- [test_do_generic.go](../test_do_generic.go) - 通用性測試示例
- [V1.2.1_SUMMARY.md](V1.2.1_SUMMARY.md) - 版本摘要