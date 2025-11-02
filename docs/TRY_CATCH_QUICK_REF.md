# Try-Catch-Finally 快速參考

## 核心概念

### Try-Catch-Finally
生成純粹的 try-catch-finally 語句，**不包裝在函數中**

### Do / AsyncDo
創建立即執行函數表達式（IIFE）

## 基本語法

### 1. Try-Catch
```go
js.Try(
    js.Const("x", "1"),
    js.Log("x"),
).Catch(
    js.Log("'錯誤: ' + error.message"),
).End()  // ⚠️ 必須調用 .End()
```

**生成：**
```javascript
try {
  const x = 1;
  console.log(x);
} catch (error) {
  console.log('錯誤: ' + error.message);
}
```

### 2. Try-Catch-Finally
```go
js.Try(
    js.Const("x", "1"),
).Catch(
    js.Log("'錯誤: ' + error.message"),
).Finally(
    js.Log("'清理完成'"),
)  // ✅ Finally 自動返回 JSAction，不需要 .End()
```

**生成：**
```javascript
try {
  const x = 1;
} catch (error) {
  console.log('錯誤: ' + error.message);
} finally {
  console.log('清理完成');
}
```

### 3. Try-Finally
```go
js.Try(
    js.Const("lock", "acquireLock()"),
    JSAction{Code: "doWork()"},
).Finally(
    JSAction{Code: "lock.release()"},
)
```

### 4. Do（立即執行函數）
```go
js.Do(
    js.Const("x", "1"),
    js.Log("x"),
)
```

**生成：**
```javascript
(() => {
  const x = 1;
  console.log(x);
})()
```

### 5. AsyncDo（立即執行異步函數）
```go
js.AsyncDo(
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)
```

**生成：**
```javascript
(async () => {
  const data = await fetch('/api');
  console.log(data);
})()
```

## 重要規則

| 規則 | 說明 |
|------|------|
| **錯誤對象名稱** | 必須使用 `error`（不是 `e`） |
| **Try-Catch 結尾** | 必須調用 `.End()` |
| **Try-Catch-Finally 結尾** | `.Finally()` 自動返回，不需要 `.End()` |
| **至少一個區塊** | 必須有 `.Catch()` 或 `.Finally()`，不能只有 Try |
| **Try 不包裝** | Try 生成純粹的 try-catch-finally，不包裝在函數中 |
| **需要 async** | 使用 AsyncFn 或 AsyncDo 包裝 |

## 常見使用模式

### 1. 同步 Try-Catch
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

### 2. AsyncFn 中的 Try（事件處理器）
```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Try(
            js.Const("response", "await fetch('/api/users')"),
            js.Const("users", "await response.json()"),
            // 處理數據...
        ).Catch(
            js.Alert("'載入失敗: ' + error.message"),
        ).End(),
    ),
}, "載入")
```

### 3. AsyncDo 立即執行
```go
Script(nil,
    js.AsyncDo(
        js.Try(
            js.Const("response", "await fetch('/api/config')"),
            js.Const("config", "await response.json()"),
            JSAction{Code: "window.appConfig = config"},
        ).Catch(
            js.Log("'初始化失敗: ' + error.message"),
        ).End(),
    ),
)
```

### 4. Do 創建獨立作用域
```go
Button(Props{
    "onClick": js.Do(
        js.Const("timestamp", "Date.now()"),
        js.Const("message", "'點擊時間: ' + new Date(timestamp).toLocaleTimeString()"),
        js.Alert("message"),
    ),
}, "顯示時間")
```

### 5. 資源管理（Try-Finally）
```go
js.AsyncFn(nil,
    js.Const("connection", "null"),
    js.Try(
        js.Assign("connection", "await db.connect()"),
        js.Const("result", "await connection.query('SELECT * FROM users')"),
    ).Finally(
        JSAction{Code: "if (connection) connection.close()"},
    ),
)
```

## 使用場景對照表

| 場景 | 使用方式 |
|------|----------|
| 同步錯誤處理 | `js.Try(...).Catch(...).End()` |
| 事件處理器 + async | `js.AsyncFn(nil, js.Try(...).Catch(...).End())` |
| 頁面載入時執行 async | `js.AsyncDo(js.Try(...).Catch(...).End())` |
| 創建獨立作用域 | `js.Do(...)` |
| 資源清理 | `js.Try(...).Finally(...)` |
| 完整錯誤處理 | `js.Try(...).Catch(...).Finally(...)` |

## 常見錯誤

### ❌ 錯誤：只有 Try，沒有 Catch 或 Finally
```go
js.Try(
    js.Const("data", "await fetch('/api')"),
)  // 會引發 panic
```

### ✅ 正確：至少要有 Catch 或 Finally
```go
js.Try(
    js.Const("data", "await fetch('/api')"),
).Catch(
    js.Log("'錯誤'"),
).End()
```

---

### ❌ 錯誤：Try-Catch 沒有調用 .End()
```go
Button(Props{
    "onClick": js.Try(...).Catch(...),  // 類型錯誤
})
```

### ✅ 正確：調用 .End()
```go
Button(Props{
    "onClick": js.Try(...).Catch(...).End(),
})
```

---

### ❌ 錯誤：使用 e 而非 error
```go
js.Try(...).Catch(
    js.Log("e.message"),  // e 未定義
).End()
```

### ✅ 正確：使用 error
```go
js.Try(...).Catch(
    js.Log("error.message"),
).End()
```

---

### ❌ 錯誤：在同步函數中使用 await
```go
Button(Props{
    "onClick": js.Fn(nil,  // ❌ Fn 不支持 await
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(...).End(),
    ),
})
```

### ✅ 正確：使用 AsyncFn
```go
Button(Props{
    "onClick": js.AsyncFn(nil,  // ✅
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(...).End(),
    ),
})
```

## API 設計理念

### 職責分離

- **Try**：只負責生成 try-catch-finally 語句
- **AsyncFn**：創建異步事件處理函數
- **Do**：創建立即執行的普通函數（IIFE）
- **AsyncDo**：創建立即執行的異步函數（async IIFE）

### 為什麼 Try 不包裝在 IIFE 中？

1. **更靈活**：Try 可以在任何上下文中使用（同步或異步函數內）
2. **更清晰**：職責分離，Try 只管錯誤處理
3. **更可控**：用戶明確決定何時需要 async/await

## 完整示例

### API 請求完整流程
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
            js.ForEachJS("users", "user",
                js.Const("div", "document.createElement('div')"),
                JSAction{Code: "div.textContent = user.name"},
                JSAction{Code: "container.appendChild(div)"},
            ),
        ).Catch(
            JSAction{Code: "container.innerHTML = '<div class=\"error\">' + error.message + '</div>'"},
            js.Log("'錯誤: ' + error.message"),
        ).Finally(
            js.Log("'請求完成'"),
            JSAction{Code: "hideLoadingSpinner()"},
        ),
    ),
}, "載入用戶")
```

## 參考資料

- `examples/07_trycatch_usage.go` - 完整交互式示例
- `docs/TRY_CATCH_FINALLY.md` - 詳細說明文檔

運行示例：
```bash
go run examples/07_trycatch_usage.go
```
訪問 http://localhost:8086