# Try-Catch-Finally API 重新設計

## 變更日期
2024

## 摘要

將 Try-Catch-Finally API 重新設計為更靈活的架構：**Try 生成純粹的 try-catch-finally 語句，不包裝在自執行函數中**。新增 **Do** 和 **AsyncDo** 專門用於創建立即執行函數（IIFE）。

## 設計理念

### 核心原則

1. **職責分離**：Try 只負責錯誤處理，IIFE 由 Do/AsyncDo 負責
2. **靈活性優先**：用戶明確決定何時需要 async/await
3. **可組合性**：Try 可以在任何上下文中使用（同步或異步函數內）

### 為什麼要這樣設計？

舊設計的問題：
- Try 自動包裝在 async IIFE 中，不夠靈活
- 無法在普通函數或已有的 async 函數中使用 Try
- 職責不清：Try 同時負責錯誤處理和創建 IIFE

新設計的優勢：
- ✅ Try 生成純粹的 try-catch-finally，可以在任何地方使用
- ✅ 需要 async 時，由用戶用 AsyncFn 或 AsyncDo 包裝
- ✅ 職責清晰分離：Try 管錯誤，Do/AsyncDo 管 IIFE
- ✅ 更加靈活和可控

## 新 API

### Try-Catch-Finally

生成純粹的 try-catch-finally 語句（不包裝）：

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
    js.Const("lock", "acquireLock()"),
).Finally(
    JSAction{Code: "lock.release()"},
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

### Do - 立即執行的普通函數

創建立即執行函數表達式（IIFE）：

```go
js.Do(
    js.Const("x", "1"),
    js.Log("x"),
)
```

**生成的 JavaScript：**

```javascript
(() => {
  const x = 1;
  console.log(x);
})()
```

### AsyncDo - 立即執行的異步函數

創建立即執行的異步函數表達式（async IIFE）：

```go
js.AsyncDo(
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)
```

**生成的 JavaScript：**

```javascript
(async () => {
  const data = await fetch('/api');
  console.log(data);
})()
```

## 使用場景

### 1. 同步 Try-Catch

當不需要 async/await 時：

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

### 2. AsyncFn 中的 Try

在事件處理器中使用 async/await：

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

### 3. AsyncDo 立即執行

在頁面載入時立即執行異步代碼：

```go
Script(nil,
    js.AsyncDo(
        js.Try(
            js.Const("response", "await fetch('/api/config')"),
            js.Const("config", "await response.json()"),
            JSAction{Code: "window.appConfig = config"},
        ).Catch(
            js.Log("'配置載入失敗'"),
        ).End(),
    ),
)
```

**生成的 JavaScript：**

```javascript
<script>
(async () => {
  try {
    const response = await fetch('/api/config');
    const config = await response.json();
    window.appConfig = config;
  } catch (error) {
    console.log('配置載入失敗');
  }
})()
</script>
```

### 4. Do 創建獨立作用域

避免變數污染全局作用域：

```go
Button(Props{
    "onClick": js.Do(
        js.Const("timestamp", "Date.now()"),
        js.Const("message", "'點擊時間: ' + new Date(timestamp).toLocaleTimeString()"),
        js.Alert("message"),
    ),
}, "顯示時間")
```

## 遷移指南

### 從舊 API 遷移

#### 舊設計（自動包裝在 async IIFE 中）

```go
// 舊：Try 自動包裝在 async IIFE 中
Button(Props{
    "onClick": js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("error.message"),
    ).End(),
})
```

這在舊設計中會生成：
```javascript
(async () => {
  try {
    const data = await fetch('/api');
  } catch (error) {
    console.log(error.message);
  }
})()
```

#### 新設計（不包裝，由用戶決定）

```go
// 新：明確使用 AsyncFn 包裝
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(
            js.Log("error.message"),
        ).End(),
    ),
})
```

這會生成：
```javascript
async () => {
  try {
    const data = await fetch('/api');
  } catch (error) {
    console.log(error.message);
  }
}
```

### 遷移步驟

1. **識別 Try 的使用場景**：
   - 是否在事件處理器中？→ 用 AsyncFn 包裝
   - 是否需要立即執行？→ 用 AsyncDo 包裝
   - 是否在已有的 async 函數中？→ 直接使用 Try

2. **更新代碼**：

```go
// 舊：直接在 onClick 中使用 Try
"onClick": js.Try(...).Catch(...).End()

// 新：用 AsyncFn 包裝
"onClick": js.AsyncFn(nil, js.Try(...).Catch(...).End())
```

3. **更新錯誤對象名稱**（如果從更舊的 TryCatch 遷移）：
   - 舊：`e.message`
   - 新：`error.message`

## 使用場景對照表

| 場景 | 使用方式 | 說明 |
|------|----------|------|
| 同步錯誤處理 | `js.Try(...).Catch(...).End()` | 不需要 await |
| 事件處理器 + async | `js.AsyncFn(nil, js.Try(...).Catch(...).End())` | 用 AsyncFn 包裝 |
| 頁面載入時執行 | `js.AsyncDo(js.Try(...).Catch(...).End())` | 立即執行 |
| 獨立作用域 | `js.Do(...)` | 創建 IIFE |
| 在 async 函數中 | `js.Try(...).Catch(...).End()` | 直接使用 |
| 資源清理 | `js.Try(...).Finally(...)` | 不需要 catch |

## API 完整對比

### 舊 TryCatch API

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

- 自動包裝在 async IIFE 中
- 需要陣列包裝
- 必須傳入所有三個參數（即使為 nil）
- 錯誤對象使用 `e`

### 新 Try + AsyncFn/AsyncDo API

```go
// 在事件處理器中
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),  // 使用 error
    ).Finally(
        js.Log("'清理'"),
    ),
)

// 或立即執行
js.AsyncDo(
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(
        js.Log("'錯誤: ' + error.message"),
    ).Finally(
        js.Log("'清理'"),
    ),
)
```

- 不自動包裝，更靈活
- 流暢的鏈式調用
- 可選的 catch/finally
- 錯誤對象統一使用 `error`
- 職責分離清晰

## 實現細節

### tryBuilder 結構

```go
type tryBuilder struct {
    tryActions     []JSAction
    catchActions   []JSAction
    finallyActions []JSAction
}

func Try(actions ...JSAction) *tryBuilder
func (tb *tryBuilder) Catch(actions ...JSAction) *tryBuilder
func (tb *tryBuilder) Finally(actions ...JSAction) JSAction
func (tb *tryBuilder) End() JSAction
```

### 關鍵設計決策

1. **鏈式返回**：`Catch()` 返回 `*tryBuilder`，允許繼續調用 `Finally()`
2. **自動終結**：`Finally()` 直接返回 `JSAction`，不需要 `.End()`
3. **顯式終結**：只有 `Catch` 時需要調用 `.End()` 來獲得 `JSAction`
4. **錯誤對象**：統一使用 `error` 而非 `e`
5. **不包裝**：生成純粹的 try-catch-finally，不包裝在函數中

### Do 和 AsyncDo 實現

```go
func Do(actions ...JSAction) JSAction {
    // 生成：(() => { ...actions })()
}

func AsyncDo(actions ...JSAction) JSAction {
    // 生成：(async () => { ...actions })()
}
```

## 優勢總結

| 特性 | 舊設計 | 新設計 |
|------|--------|--------|
| 自動包裝 | ✅ 自動 async IIFE | ❌ 不包裝，更靈活 |
| 職責分離 | ❌ Try 管太多 | ✅ Try 只管錯誤 |
| 可組合性 | ❌ 只能頂層使用 | ✅ 任何地方使用 |
| 明確性 | ❌ 隱式 async | ✅ 顯式 async |
| 靈活性 | ⚠️ 較低 | ✅ 非常高 |
| IIFE 控制 | ❌ 自動創建 | ✅ 用戶控制（Do/AsyncDo） |

## 向後兼容性

### 舊 TryCatch 函數

舊的 `TryCatch` 函數**仍然保留**並可以正常使用：

```go
func TryCatch(tryActions []JSAction, catchActions []JSAction, finallyActions []JSAction) JSAction
```

但建議在新代碼中使用新的流暢 API。

## 文檔

完整文檔：
- `docs/TRY_CATCH_FINALLY.md` - 詳細使用指南
- `docs/TRY_CATCH_QUICK_REF.md` - 快速參考手冊
- `examples/07_trycatch_usage.go` - 完整交互式示例

運行示例：
```bash
go run examples/07_trycatch_usage.go
```
訪問 http://localhost:8086

## 總結

新的 Try-Catch-Finally API 設計遵循以下原則：

1. **職責分離**：Try 只負責錯誤處理，IIFE 由 Do/AsyncDo 負責
2. **靈活性優先**：不自動包裝，用戶明確決定何時需要 async
3. **可組合性**：Try 可以在任何上下文中使用
4. **清晰明確**：代碼意圖更加明確，易於理解和維護

這種設計提供了最大的靈活性，同時保持了代碼的清晰和易讀性。