# Do/AsyncDo 參數說明

## 概述

從 v1.2.1 開始，`js.Do()` 和 `js.AsyncDo()` 的簽名與 `js.Fn()` 和 `js.AsyncFn()` 保持一致，第一個參數為參數列表 `[]string`。

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

## 使用方式

### 1. 無參數（最常見）

大多數事件處理器不需要顯式參數，傳入 `nil`：

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
- 參數名可以自定義（如 `e`、`evt`），但調用時傳入的始終是外部作用域的 `event`
- 第一個 `event`/`e`/`evt` 是參數名，調用時的 `event` 是傳入的值

### 3. 帶自定義參數（進階用法）

如果需要定義接受參數的 IIFE（較少見）：

```go
// 定義帶參數的 IIFE 並立即調用
js.Do([]string{"x", "y"},
    js.Const("sum", "x + y"),
    js.Log("sum"),
)
```

生成：
```javascript
((x,y)=>{const sum=x+y;console.log(sum)})()
```

**使用場景**：這種用法在事件處理器中不常見，因為 HTML 事件不會傳遞自定義參數。主要用於需要立即執行且接受參數的代碼塊。

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

**注意**：調用時傳入的始終是 `event`（HTML 事件處理器的標準名稱）。

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
A: 不一定！你可以使用任何名字（`e`、`evt`、`myEvent` 等），只要與代碼中使用的變量名一致即可。調用時傳入的始終是外部作用域的 `event`。

**Q: 什麼時候傳 nil，什麼時候傳 []string？**  
A: 
- 傳 `nil`：當你的代碼不使用 `event` 對象時（如簡單的 alert）
- 傳 `[]string{"event"}`（或 `[]string{"e"}`）：當你需要訪問 `event.target`、`event.preventDefault()` 等時

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
- `event` - 最清晰明確，與外部作用域名稱一致（推薦用於複雜邏輯）
- `e` - 最簡短，常見於 JavaScript 社區（推薦用於簡單邏輯）
- `evt` - 折衷方案，清晰且簡短

**Q: 如果忘記聲明 event 參數會怎樣？**  
A: 會出現運行時錯誤：`Cannot read properties of undefined (reading 'preventDefault')` 或類似的錯誤，因為 IIFE 內部的 `event` 是 `undefined`。

## 示例

完整示例請參考：
- [examples/09_event_handlers.go](../examples/09_event_handlers.go) - 事件處理器示例
- [examples/10_do_with_params.go](../examples/10_do_with_params.go) - 參數使用示例

## 相關文檔

- [EVENT_HANDLER_CHANGES.md](EVENT_HANDLER_CHANGES.md) - 事件處理器變更說明
- [EVENT_HANDLER_QUICK_REF.md](EVENT_HANDLER_QUICK_REF.md) - 快速參考
- [V1.2.1_SUMMARY.md](V1.2.1_SUMMARY.md) - 版本摘要