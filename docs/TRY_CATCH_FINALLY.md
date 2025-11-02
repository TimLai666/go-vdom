# Try-Catch-Finally 流暢 API

Go VDOM 提供了優雅的流暢 API 來處理 JavaScript 的 try-catch-finally 語句。

## 設計理念

### 核心原則

1. **Try 生成純粹的 try-catch-finally 語句**：不包裝在自執行函數中
2. **需要 async 時由用戶決定**：使用 `AsyncFn` 或 `AsyncDo` 包裝
3. **Do/AsyncDo 專門創建 IIFE**：立即執行函數表達式

### 為什麼這樣設計？

- **更靈活**：Try 可以在任何上下文中使用（同步或異步函數內）
- **更清晰**：職責分離，Try 只負責 try-catch-finally，IIFE 由 Do/AsyncDo 負責
- **更可控**：用戶明確決定何時需要 async/await

## API 概覽

### Try-Catch-Finally

```go
// Try-Catch
js.Try(
    js.Const("x", "1"),
    js.Log("x"),
).Catch(
    js.Log("'錯誤: ' + error.message"),
).End()  // ⚠️ 必須調用 .End()

// Try-Catch-Finally
js.Try(
    js.Const("x", "1"),
).Catch(
    js.Log("'錯誤: ' + error.message"),
).Finally(
    js.Log("'清理完成'"),
)  // ✅ Finally 自動返回 JSAction

// Try-Finally
js.Try(
    js.Const("x", "1"),
).Finally(
    js.Log("'清理完成'"),
)
```

**生成的 JavaScript：**

```javascript
try {
  const x = 1;
  console.log(x);
} catch (error) {
  console.log('錯誤: ' + error.message);
}
```

### Do / AsyncDo

```go
// Do - 立即執行的普通函數
js.Do(
    js.Const("x", "1"),
    js.Log("x"),
)

// AsyncDo - 立即執行的異步函數
js.AsyncDo(
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)
```

**生成的 JavaScript：**

```javascript
// Do
(() => {
  const x = 1;
  console.log(x);
})()

// AsyncDo
(async () => {
  const data = await fetch('/api');
  console.log(data);
})()
```

## 使用場景

### 1. 同步 Try-Catch

當不需要 async/await 時，直接使用 Try：

```go
Button(Props{
    "onClick": js.Fn(nil,
        js.Try(
            js.Const("x", "Math.random()"),
            JSAction{Code: "if (x < 0.5) throw new Error('太小')"},
            js.Alert("'成功: ' + x"),
        ).Catch(
            js.Alert("'錯誤: ' + error.message"),
        ).End(),
    ),
}, "測試")
```

### 2. 在 AsyncFn 中使用 Try

當需要在事件處理器中使用 async/await：

```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Try(
            js.Const("response", "await fetch('/api/data')"),
            js.Const("data", "await response.json()"),
            js.Alert("'成功載入'"),
        ).Catch(
            js.Alert("'錯誤: ' + error.message"),
        ).End(),
    ),
}, "載入數據")
```

**生成的 JavaScript：**

```javascript
async () => {
  try {
    const response = await fetch('/api/data');
    const data = await response.json();
    alert('成功載入');
  } catch (error) {
    alert('錯誤: ' + error.message);
  }
}
```

### 3. 使用 AsyncDo 立即執行

當需要在頂層立即執行異步代碼：

```go
Script(nil,
    js.AsyncDo(
        js.Try(
            js.Const("response", "await fetch('/api/init')"),
            js.Const("config", "await response.json()"),
            JSAction{Code: "window.appConfig = config"},
        ).Catch(
            js.Log("'初始化失敗: ' + error.message"),
        ).End(),
    ),
)
```

**生成的 JavaScript：**

```javascript
<script>
(async () => {
  try {
    const response = await fetch('/api/init');
    const config = await response.json();
    window.appConfig = config;
  } catch (error) {
    console.log('初始化失敗: ' + error.message);
  }
})()
</script>
```

### 4. 使用 Do 創建獨立作用域

當需要避免變數污染全局作用域：

```go
Button(Props{
    "onClick": js.Do(
        js.Const("timestamp", "Date.now()"),
        js.Const("message", "'點擊時間: ' + new Date(timestamp).toLocaleTimeString()"),
        js.Alert("message"),
    ),
}, "顯示時間")
```

## 實際應用示例

### API 請求完整示例

```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Const("container", "document.getElementById('result')"),
        js.Try(
            JSAction{Code: "container.innerHTML = '載入中...'"},
            js.Const("response", "await fetch('/api/users')"),
            js.Const("users", "await response.json()"),
            JSAction{Code: "if (!response.ok) throw new Error('請求失敗')"},
            
            // 渲染數據
            JSAction{Code: "container.innerHTML = ''"},
            js.Const("ul", "document.createElement('ul')"),
            JSAction{Code: "ul.className = 'list-group'"},
            js.ForEachJS("users", "user",
                js.Const("li", "document.createElement('li')"),
                JSAction{Code: "li.className = 'list-group-item'"},
                JSAction{Code: "li.textContent = user.name"},
                JSAction{Code: "ul.appendChild(li)"},
            ),
            JSAction{Code: "container.appendChild(ul)"},
        ).Catch(
            JSAction{Code: "container.innerHTML = '<div class=\"alert alert-danger\">' + error.message + '</div>'"},
            js.Log("'API 錯誤: ' + error.message"),
        ).Finally(
            js.Log("'請求完成'"),
            JSAction{Code: "hideLoadingSpinner()"},
        ),
    ),
}, "載入用戶")
```

### 表單提交

```go
Form(Props{
    "onSubmit": js.AsyncFn([]string{"event"},
        JSAction{Code: "event.preventDefault()"},
        js.Const("formData", "new FormData(event.target)"),
        js.Try(
            js.Const("response", "await fetch('/api/submit', { method: 'POST', body: formData })"),
            js.Const("result", "await response.json()"),
            JSAction{Code: "if (!response.ok) throw new Error(result.message)"},
            js.Alert("'提交成功'"),
            JSAction{Code: "event.target.reset()"},
        ).Catch(
            js.Alert("'提交失敗: ' + error.message"),
        ).Finally(
            js.Log("'提交處理完成'"),
        ),
    ),
}, /* form fields... */)
```

### 資源管理

```go
js.AsyncFn(nil,
    js.Const("connection", "null"),
    js.Try(
        js.Assign("connection", "await database.connect()"),
        js.Const("data", "await connection.query('SELECT * FROM users')"),
        JSAction{Code: "processData(data)"},
    ).Finally(
        JSAction{Code: "if (connection) await connection.close()"},
        js.Log("'連接已關閉'"),
    ),
)
```

### 初始化腳本（頁面載入時執行）

```go
Script(nil,
    js.AsyncDo(
        js.Try(
            js.Const("response", "await fetch('/api/config')"),
            js.Const("config", "await response.json()"),
            JSAction{Code: "window.appConfig = config"},
            js.Log("'配置載入完成'"),
        ).Catch(
            js.Log("'配置載入失敗，使用預設值'"),
            JSAction{Code: "window.appConfig = { theme: 'light' }"},
        ).End(),
    ),
)
```

## 重要規則

### 1. 錯誤對象名稱

在 `Catch` 區塊中，錯誤對象的變數名為 **`error`**：

```go
js.Try(
    js.Const("result", "await riskyOperation()"),
).Catch(
    // ✅ 正確：使用 error
    js.Log("'錯誤: ' + error.message"),
    JSAction{Code: "console.error(error.stack)"},
).End()
```

### 2. Try-Catch 必須調用 .End()

當只使用 Try-Catch（沒有 Finally）時，必須調用 `.End()` 來獲得 `JSAction`：

```go
// ❌ 錯誤：缺少 .End()
js.Try(...).Catch(...)

// ✅ 正確
js.Try(...).Catch(...).End()

// ✅ 正確：Finally 會自動返回 JSAction，不需要 .End()
js.Try(...).Catch(...).Finally(...)
```

### 3. 必須有 Catch 或 Finally

你不能只有 Try 而沒有 Catch 或 Finally：

```go
// ❌ 錯誤：會引發 panic
js.Try(
    js.Const("data", "await fetch('/api')"),
)

// ✅ 正確：至少要有 Catch 或 Finally
js.Try(
    js.Const("data", "await fetch('/api')"),
).Catch(
    js.Log("'錯誤'"),
).End()
```

### 4. Try 不包裝在函數中

Try 生成純粹的 try-catch-finally 語句，不會自動包裝在 IIFE 中：

```go
// 生成的是純粹的 try-catch，不是 (async () => { try... })()
js.Try(
    js.Const("x", "1"),
).Catch(
    js.Log("error"),
).End()
```

如果需要 IIFE，使用 `Do` 或 `AsyncDo`：

```go
// 包裝在 IIFE 中
js.AsyncDo(
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("error"),
    ).End(),
)
```

### 5. 在 AsyncFn 中使用 await

當需要在事件處理器中使用 await 時，用 `AsyncFn` 包裝：

```go
Button(Props{
    "onClick": js.AsyncFn(nil,  // ✅ 使用 AsyncFn
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(...).End(),
    ),
})
```

## 常見模式對照表

| 場景 | 使用方式 |
|------|----------|
| 同步 try-catch | `js.Try(...).Catch(...).End()` |
| 事件處理器 + async | `js.AsyncFn(nil, js.Try(...).Catch(...).End())` |
| 立即執行 async | `js.AsyncDo(js.Try(...).Catch(...).End())` |
| 獨立作用域 | `js.Do(...)` |
| 資源清理 | `js.Try(...).Finally(...)` |
| 完整錯誤處理 | `js.Try(...).Catch(...).Finally(...)` |

## 常見錯誤

### ❌ 錯誤 1：只有 Try，沒有 Catch 或 Finally

```go
js.Try(
    js.Const("data", "await fetch('/api')"),
)  // panic: Try requires at least Catch() or Finally()
```

**解決方案：**
```go
js.Try(
    js.Const("data", "await fetch('/api')"),
).Catch(
    js.Log("'錯誤'"),
).End()
```

### ❌ 錯誤 2：Try-Catch 沒有調用 .End()

```go
Button(Props{
    "onClick": js.Try(...).Catch(...),  // 類型錯誤
})
```

**解決方案：**
```go
Button(Props{
    "onClick": js.Try(...).Catch(...).End(),  // ✅
})
```

### ❌ 錯誤 3：使用 e 而非 error

```go
js.Try(...).Catch(
    js.Log("e.message"),  // e 未定義
).End()
```

**解決方案：**
```go
js.Try(...).Catch(
    js.Log("error.message"),  // ✅ 使用 error
).End()
```

### ❌ 錯誤 4：在同步函數中使用 await

```go
Button(Props{
    "onClick": js.Fn(nil,  // ❌ 普通函數不支持 await
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(...).End(),
    ),
})
```

**解決方案：**
```go
Button(Props{
    "onClick": js.AsyncFn(nil,  // ✅ 使用 AsyncFn
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(...).End(),
    ),
})
```

## 與舊 API 的對比

### 舊 API（TryCatch）

```go
js.TryCatch(
    []JSAction{
        js.Const("data", "await fetch('/api')"),
    },
    []JSAction{
        js.Log("'錯誤: ' + e.message"),  // 使用 e
    },
    []JSAction{
        js.Log("'清理'"),
    },
)
```

**問題：**
- 自動包裝在 IIFE 中，不夠靈活
- 需要陣列包裝
- 必須傳入所有三個參數（即使為 nil）
- 錯誤對象使用 `e` 而非 `error`

### 新 API（推薦）

```go
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),  // 使用 error
    ).Finally(
        js.Log("'清理'"),
    ),
)
```

**優勢：**
- ✅ 不自動包裝，更靈活
- ✅ 流暢的鏈式調用
- ✅ 可選的 catch/finally
- ✅ 錯誤對象統一命名為 `error`
- ✅ 職責分離清晰（Try 負責錯誤處理，AsyncFn/AsyncDo 負責 async）

## 完整示例

參見 `examples/07_trycatch_usage.go` 查看所有用法的交互式示例。

運行示例：

```bash
go run examples/07_trycatch_usage.go
```

然後訪問 http://localhost:8086

## 總結

新的 Try-Catch-Finally API 設計遵循以下原則：

1. **Try 生成純粹的 try-catch-finally**：不包裝在函數中
2. **需要 async 時由用戶決定**：使用 AsyncFn 或 AsyncDo
3. **Do/AsyncDo 專門創建 IIFE**：職責清晰分離
4. **流暢的鏈式調用**：更接近 JavaScript 原生語法
5. **錯誤對象統一命名**：使用 `error`

這種設計提供了最大的靈活性，同時保持了代碼的清晰和易讀性。