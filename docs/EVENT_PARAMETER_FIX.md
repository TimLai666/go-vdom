# Event Parameter Fix

## 問題描述

用戶報告在使用 `onSubmit` 事件處理器時出現錯誤：

```
Uncaught (in promise) TypeError: Cannot read properties of undefined (reading 'preventDefault')
```

## 根本原因

當使用 `js.Do(nil, ...)` 或 `js.AsyncDo(nil, ...)` 創建 IIFE 時，生成的代碼為：

```javascript
(()=>{
    event.preventDefault();  // ❌ event is undefined!
})()
```

**問題**：IIFE 創建了新的函數作用域，而 `event` 對象在外部作用域。在新作用域內訪問 `event` 會得到 `undefined`。

## 解決方案

### 必須顯式傳遞 event 參數

當代碼中需要使用 `event` 對象時，必須將其聲明為參數並傳遞：

```go
// ❌ 錯誤 - event 在 IIFE 內部是 undefined
"onSubmit": js.Do(nil,
    js.CallMethod("event", "preventDefault"),
)

// ✅ 正確 - 顯式聲明並傳遞 event
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
)
```

生成的代碼對比：

```javascript
// 錯誤的生成代碼
onsubmit="(()=>{event.preventDefault()})()"
// event 在 IIFE 內部找不到

// 正確的生成代碼
onsubmit="((event)=>{event.preventDefault()})(event)"
// event 作為參數傳入 IIFE
```

## 作用域說明

### HTML 事件處理器作用域

在 HTML 內聯事件處理器中：

```html
<button onclick="console.log(event)">Click</button>
```

`event` 對象在事件處理器的作用域中可用。

### IIFE 作用域隔離

但當創建 IIFE 時：

```html
<button onclick="(()=>{console.log(event)})()">Click</button>
```

IIFE 的箭頭函數 `()=>{}` 創建了新的作用域，`event` 必須顯式傳入：

```html
<button onclick="((event)=>{console.log(event)})(event)">Click</button>
```

## 何時需要聲明 event 參數

### 需要聲明 event 參數的情況

任何使用 `event` 對象的代碼都需要聲明參數：

```go
// ✅ 使用 event.preventDefault()
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
)

// ✅ 使用 event.target
"onClick": js.Do([]string{"event"},
    js.Const("target", "event.target"),
)

// ✅ 使用 event.target.value
"onInput": js.Do([]string{"event"},
    js.Const("value", "event.target.value"),
)

// ✅ 使用 event.target.checked
"onChange": js.Do([]string{"event"},
    js.Const("checked", "event.target.checked"),
)

// ✅ 表單數據
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    js.Const("formData", "new FormData(event.target)"),
)
```

### 不需要聲明 event 參數的情況

如果不使用 `event` 對象，傳入 `nil`：

```go
// ✅ 簡單的 alert
"onClick": js.Do(nil,
    js.Alert("'Clicked!'"),
)

// ✅ 訪問已知 ID 的元素
"onClick": js.Do(nil,
    js.Const("el", "document.getElementById('myDiv')"),
    JSAction{Code: "el.style.display = 'none'"},
)

// ✅ 執行預定義操作
"onClick": js.Do(nil,
    js.Log("'Button was clicked'"),
)
```

## 修復的文件

以下文件已修復，添加了正確的 event 參數聲明：

### examples/03_javascript_dsl.go
```go
// onSubmit 處理器
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    // ...
)
```

### examples/09_event_handlers.go
```go
// onInput 處理器
"onInput": js.Do([]string{"event"},
    js.Const("val", "event.target.value"),
    // ...
)

// onChange 處理器
"onChange": js.Do([]string{"event"},
    js.Const("val", "event.target.value"),
    // ...
)
```

### examples/10_do_with_params.go
```go
// 所有使用 event 的處理器都已更新
"onClick": js.Do([]string{"event"},
    js.Const("target", "event.target"),
    // ...
)

"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    // ...
)

"onInput": js.Do([]string{"event"},
    js.Const("value", "event.target.value"),
    // ...
)

"onChange": js.Do([]string{"event"},
    js.Const("selectedValue", "event.target.value"),
    // ...
)
```

## 檢查清單

修復 event 參數問題時，檢查以下項目：

- [ ] 代碼中是否使用了 `event.target`？
- [ ] 代碼中是否調用了 `event.preventDefault()`？
- [ ] 代碼中是否調用了 `event.stopPropagation()`？
- [ ] 代碼中是否訪問了 `event` 的任何屬性或方法？
- [ ] 如果以上任一項為是，確保使用 `[]string{"event"}` 而非 `nil`

## 常見錯誤模式

### 錯誤 1：忘記聲明 event 參數

```go
// ❌ 錯誤
Form(Props{
    "onSubmit": js.Do(nil,
        js.CallMethod("event", "preventDefault"),
    ),
})

// ✅ 正確
Form(Props{
    "onSubmit": js.Do([]string{"event"},
        js.CallMethod("event", "preventDefault"),
    ),
})
```

**錯誤信息**：`Cannot read properties of undefined (reading 'preventDefault')`

### 錯誤 2：在不需要時聲明 event 參數

雖然這不會導致錯誤，但不必要：

```go
// ⚠️ 不必要但無害
"onClick": js.Do([]string{"event"},
    js.Alert("'Hello'"),  // 不使用 event
)

// ✅ 更簡潔
"onClick": js.Do(nil,
    js.Alert("'Hello'"),
)
```

### 錯誤 3：參數名拼寫錯誤

```go
// ❌ 錯誤 - 參數名與使用的變量名不匹配
"onClick": js.Do([]string{"evt"},
    js.Const("target", "event.target"),  // event is undefined!
)

// ✅ 正確 - 保持一致
"onClick": js.Do([]string{"event"},
    js.Const("target", "event.target"),
)
```

## 快速參考

| 使用場景 | 參數 | 範例 |
|---------|------|------|
| 不使用 event | `nil` | `js.Do(nil, js.Alert("'Hi'"))` |
| 使用 event.target | `[]string{"event"}` | `js.Do([]string{"event"}, js.Const("el", "event.target"))` |
| 使用 preventDefault | `[]string{"event"}` | `js.Do([]string{"event"}, js.CallMethod("event", "preventDefault"))` |
| 表單處理 | `[]string{"event"}` | `js.Do([]string{"event"}, js.Const("data", "new FormData(event.target)"))` |
| 輸入值變化 | `[]string{"event"}` | `js.Do([]string{"event"}, js.Const("val", "event.target.value"))` |

## 技術細節

### 為什麼 IIFE 需要參數傳遞？

JavaScript 的作用域規則：

```javascript
// 外部作用域有 event 對象
onclick="
  // 這裡可以訪問 event
  console.log(event); // ✅ 正常工作
  
  // 但 IIFE 創建新作用域
  (()=>{
    console.log(event); // ❌ undefined!
  })();
  
  // 需要顯式傳遞
  ((event)=>{
    console.log(event); // ✅ 正常工作
  })(event);
"
```

### 生成代碼的結構

```javascript
((event) => {        // 函數參數聲明
  // 函數體
  event.preventDefault();
})(event)            // 函數調用，傳入外部作用域的 event
```

## 相關文檔

- [DO_ASYNCDO_PARAMS.md](DO_ASYNCDO_PARAMS.md) - 完整參數使用指南
- [EVENT_HANDLER_QUICK_REF.md](EVENT_HANDLER_QUICK_REF.md) - 事件處理器快速參考
- [examples/10_do_with_params.go](../examples/10_do_with_params.go) - 參數使用示例

## 總結

**關鍵要點**：

1. ✅ IIFE 創建新作用域，`event` 必須顯式傳遞
2. ✅ 使用 `event` 時必須聲明：`[]string{"event"}`
3. ✅ 不使用 `event` 時傳入：`nil`
4. ✅ 參數名必須與代碼中使用的變量名一致
5. ✅ 所有事件類型（onClick、onSubmit、onChange 等）規則相同

**記住**：當你看到 `Cannot read properties of undefined (reading 'xxx')` 錯誤時，檢查是否忘記聲明 `event` 參數！