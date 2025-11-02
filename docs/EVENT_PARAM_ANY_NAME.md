# Event 參數名稱靈活性

## 概述

從 v1.2.1 開始，`js.Do()` 和 `js.AsyncDo()` 支援使用**任何參數名稱**來接收事件對象。你不必使用 `event`，可以使用 `e`、`evt` 或任何你喜歡的名稱。

## 基本原理

在 HTML 事件處理器中，外部作用域的事件對象總是叫 `event`。當你創建 IIFE 時，你可以選擇任何參數名來接收這個對象。

## 使用示例

### 使用 'event'（最明確）

```go
Button(Props{
    "onClick": js.Do([]string{"event"},
        js.Const("target", "event.target"),
        js.CallMethod("event", "preventDefault"),
    ),
}, "Click")
```

生成：
```javascript
onclick="((event)=>{const target=event.target;event.preventDefault()})(event)"
```

### 使用 'e'（最簡短）

```go
Button(Props{
    "onClick": js.Do([]string{"e"},
        js.Const("target", "e.target"),
        js.CallMethod("e", "preventDefault"),
    ),
}, "Click")
```

生成：
```javascript
onclick="((e)=>{const target=e.target;e.preventDefault()})(event)"
```

### 使用 'evt'（折衷方案）

```go
Button(Props{
    "onClick": js.Do([]string{"evt"},
        js.Const("target", "evt.target"),
        js.CallMethod("evt", "preventDefault"),
    ),
}, "Click")
```

生成：
```javascript
onclick="((evt)=>{const target=evt.target;evt.preventDefault()})(event)"
```

### 使用自定義名稱

```go
Button(Props{
    "onClick": js.Do([]string{"clickEvent"},
        js.Const("id", "clickEvent.target.id"),
    ),
}, "Click")
```

生成：
```javascript
onclick="((clickEvent)=>{const id=clickEvent.target.id})(event)"
```

## 關鍵要點

### 1. 參數名與代碼一致

參數名必須與代碼中使用的變量名一致：

```go
// ✅ 正確 - 參數名和使用的變量名都是 'e'
js.Do([]string{"e"},
    js.Const("id", "e.target.id"),
)

// ❌ 錯誤 - 參數名是 'e' 但代碼中使用 'event'
js.Do([]string{"e"},
    js.Const("id", "event.target.id"),  // event is undefined!
)
```

### 2. 調用時傳入的始終是 'event'

無論你選擇什麼參數名，生成的代碼在調用時都會傳入外部作用域的 `event`：

```javascript
// 參數名 'e'，但調用時傳入 'event'
((e) => { ... })(event)

// 參數名 'evt'，但調用時傳入 'event'
((evt) => { ... })(event)

// 參數名 'myEvent'，但調用時傳入 'event'
((myEvent) => { ... })(event)
```

這是因為 HTML 事件處理器的外部作用域中，事件對象的標準名稱是 `event`。

### 3. 作用域隔離

```javascript
// HTML 事件處理器的外部作用域
onclick="
  // 這裡 'event' 可用
  
  // IIFE 創建新作用域
  ((e) => {
    // 這裡使用參數 'e' 來訪問事件對象
    console.log(e.target);
  })(event)  // 將外部的 'event' 傳入，在內部作為 'e' 使用
"
```

## 實際應用

### 表單處理

```go
// 使用 'e' - 簡潔
Form(Props{
    "onSubmit": js.Do([]string{"e"},
        js.CallMethod("e", "preventDefault"),
        js.Const("data", "new FormData(e.target)"),
    ),
})

// 使用 'event' - 明確
Form(Props{
    "onSubmit": js.Do([]string{"event"},
        js.CallMethod("event", "preventDefault"),
        js.Const("data", "new FormData(event.target)"),
    ),
})
```

### 輸入處理

```go
// 使用 'e' - 常見於 JavaScript 社區
Input(Props{
    "onInput": js.Do([]string{"e"},
        js.Const("value", "e.target.value"),
    ),
})

// 使用 'evt' - 折衷方案
Input(Props{
    "onInput": js.Do([]string{"evt"},
        js.Const("value", "evt.target.value"),
    ),
})
```

### 滑鼠事件

```go
Div(Props{
    "onMouseOver": js.Do([]string{"e"},
        JSAction{Code: "e.target.style.backgroundColor = 'yellow'"},
    ),
    "onMouseOut": js.Do([]string{"e"},
        JSAction{Code: "e.target.style.backgroundColor = ''"},
    ),
})
```

## 推薦使用

根據不同場景選擇合適的參數名：

### 'event' - 最佳可讀性

**優點**：
- 最清晰明確
- 與外部作用域名稱一致
- 適合複雜邏輯或團隊協作

**使用場景**：
- 複雜的事件處理邏輯
- 需要多個 event 方法調用
- 團隊項目或文檔化代碼

```go
"onSubmit": js.Do([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    js.CallMethod("event", "stopPropagation"),
    js.Const("form", "event.target"),
    js.Const("data", "new FormData(event.target)"),
    // ... 更多邏輯
)
```

### 'e' - 最佳簡潔性

**優點**：
- 最簡短
- JavaScript 社區常見慣例
- 適合簡單邏輯

**使用場景**：
- 簡單的事件處理
- 只需要 event.target
- 個人項目或快速原型

```go
"onClick": js.Do([]string{"e"},
    js.Const("id", "e.target.id"),
    js.Alert("'Clicked: ' + id"),
)
```

### 'evt' - 平衡方案

**優點**：
- 清晰且簡短
- 避免與其他變量名衝突
- 折衷方案

**使用場景**：
- 中等複雜度的邏輯
- 需要平衡可讀性和簡潔性
- 有命名規範的項目

```go
"onChange": js.Do([]string{"evt"},
    js.Const("value", "evt.target.value"),
    js.Const("name", "evt.target.name"),
)
```

## 一致性建議

### 項目級別一致性

在同一個項目中，建議保持一致的命名風格：

```go
// ✅ 好 - 整個項目都使用 'e'
"onClick": js.Do([]string{"e"}, ...)
"onSubmit": js.Do([]string{"e"}, ...)
"onChange": js.Do([]string{"e"}, ...)

// ⚠️ 可以，但不一致
"onClick": js.Do([]string{"event"}, ...)
"onSubmit": js.Do([]string{"e"}, ...)
"onChange": js.Do([]string{"evt"}, ...)
```

### 團隊協作

如果是團隊項目，建議：
1. 在項目文檔中明確規定使用哪個參數名
2. 在代碼審查中檢查一致性
3. 考慮使用 linter 或代碼風格檢查工具

## 對比表

| 參數名 | 長度 | 可讀性 | 慣例 | 推薦場景 |
|--------|------|--------|------|----------|
| `event` | 5 字符 | ⭐⭐⭐⭐⭐ | Go 風格 | 複雜邏輯、團隊協作 |
| `e` | 1 字符 | ⭐⭐⭐ | JS 風格 | 簡單邏輯、個人項目 |
| `evt` | 3 字符 | ⭐⭐⭐⭐ | 折衷 | 中等邏輯、平衡需求 |
| 自定義 | 可變 | 取決於名稱 | 不常見 | 特殊需求 |

## 常見問題

### Q1: 我應該使用哪個參數名？

**A**: 取決於你的偏好和項目需求：
- 個人項目/簡單邏輯 → `e`
- 團隊項目/複雜邏輯 → `event`
- 折衷方案 → `evt`

### Q2: 為什麼調用時傳入的是 'event' 而不是我選擇的參數名？

**A**: 因為 HTML 事件處理器的外部作用域中，事件對象的標準名稱是 `event`。我們只是在 IIFE 內部給它一個別名（你選擇的參數名）。

### Q3: 可以使用中文參數名嗎？

**A**: 技術上可以，但強烈不推薦：
```go
// ❌ 不推薦
js.Do([]string{"事件"}, js.Const("目標", "事件.target"))
```
應該使用英文參數名以保持代碼的可移植性和國際化。

### Q4: 多參數時怎麼辦？

**A**: 在事件處理器中通常只有一個 event 參數。如果你需要多個參數（非事件處理器場景），每個參數位置都會傳入 `event`：
```go
// 不常見，但支持
js.Do([]string{"e", "extra"},  // 不推薦在事件處理器中使用
    js.Log("e"),
    js.Log("extra"),  // 也會接收到 event
)
// 生成：((e, extra) => {...})(event, event)
```

### Q5: 為什麼不自動使用 'event'？

**A**: 為了給開發者選擇的自由。不同的項目和開發者有不同的編碼風格偏好。

## 技術實現

### 生成邏輯

Do/AsyncDo 的實現會：
1. 使用你指定的參數名創建函數簽名
2. 在調用時傳入外部作用域的 `event`

```go
// Go 代碼
js.Do([]string{"e"}, actions...)

// 生成 JavaScript
((e) => {
    // 你的 actions
})(event)  // 調用時傳入外部的 event
```

### 源碼參考

`jsdsl/jsdsl.go`:
```go
func Do(params []string, actions ...JSAction) JSAction {
    sb.WriteString("((")
    if params != nil {
        sb.WriteString(strings.Join(params, ","))  // 你的參數名
    }
    sb.WriteString(")=>{")
    // ... 函數體
    sb.WriteString("})(")
    if params != nil && len(params) > 0 {
        for i := range params {
            if i > 0 {
                sb.WriteString(",")
            }
            sb.WriteString("event")  // 調用時傳入 event
        }
    }
    sb.WriteString(")")
}
```

## 相關文檔

- [DO_ASYNCDO_PARAMS.md](DO_ASYNCDO_PARAMS.md) - 參數使用完整指南
- [EVENT_PARAMETER_FIX.md](EVENT_PARAMETER_FIX.md) - Event 參數修復說明
- [EVENT_HANDLER_QUICK_REF.md](EVENT_HANDLER_QUICK_REF.md) - 事件處理器快速參考
- [test_event_param.go](../test_event_param.go) - 參數名稱測試示例

## 總結

✅ **你可以使用任何參數名稱**  
✅ **參數名必須與代碼中使用的變量名一致**  
✅ **調用時傳入的始終是外部作用域的 'event'**  
✅ **選擇適合你的項目和風格的參數名**  
✅ **保持項目內的一致性**  

選擇 `event`、`e` 或 `evt`，都是有效且正確的做法！