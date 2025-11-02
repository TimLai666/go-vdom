# Try-Catch-Finally 設計理念

## 核心設計決策

### Try 不包裝在自執行函數中

**決策：** `js.Try()` 生成純粹的 `try-catch-finally` 語句，不自動包裝在 IIFE 或 async 函數中。

**理由：**

1. **最大靈活性**
   - Try 可以在任何上下文中使用：同步函數、異步函數、事件處理器、全局作用域
   - 用戶可以自由決定是否需要 async/await
   - 避免產生不必要的函數包裝

2. **職責分離**
   - Try 只負責錯誤處理邏輯
   - IIFE 創建由專門的 `Do` 和 `AsyncDo` 負責
   - 每個工具只做一件事，做好一件事

3. **可組合性**
   ```go
   // 可以在已有的 async 函數中使用 Try
   js.AsyncFn(nil,
       js.Const("x", "1"),
       js.Try(
           js.Const("data", "await fetch('/api')"),
       ).Catch(
           js.Log("error"),
       ).End(),
       js.Log("'完成'"),
   )
   ```

4. **明確性**
   - 代碼意圖更清晰：看到 `AsyncFn` 就知道是 async 函數
   - 不會產生隱藏的函數包裝
   - 降低認知負擔

### Do 和 AsyncDo 的引入

**決策：** 新增 `js.Do()` 和 `js.AsyncDo()` 專門用於創建立即執行函數（IIFE）。

**理由：**

1. **明確的語義**
   - `Do` 表示「立即執行這些動作」
   - `AsyncDo` 表示「立即執行這些異步動作」
   - 比自動包裝更容易理解

2. **通用性**
   - 不僅限於 Try-Catch，任何需要 IIFE 的場景都可使用
   - 例如：創建獨立作用域、避免變數污染

3. **對稱性**
   ```go
   js.Fn(...)      // 定義函數
   js.AsyncFn(...) // 定義異步函數
   js.Do(...)      // 立即執行函數
   js.AsyncDo(...) // 立即執行異步函數
   ```

## 使用場景分析

### 場景 1：同步錯誤處理

**需求：** 在同步代碼中處理可能的錯誤

**解決方案：**
```go
js.Fn(nil,
    js.Try(
        js.Const("x", "Math.random()"),
        JSAction{Code: "if (x < 0.5) throw new Error('太小')"},
    ).Catch(
        js.Alert("error.message"),
    ).End(),
)
```

**為什麼這樣設計：**
- Try 不包裝，可以直接在 Fn 中使用
- 不會產生額外的函數嵌套

### 場景 2：事件處理器中的異步操作

**需求：** 在按鈕點擊時執行異步請求並處理錯誤

**解決方案：**
```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Try(
            js.Const("data", "await fetch('/api')"),
        ).Catch(
            js.Alert("error.message"),
        ).End(),
    ),
})
```

**為什麼這樣設計：**
- AsyncFn 明確表示這是一個異步事件處理器
- Try 可以直接在其中使用 await
- 職責清晰：AsyncFn 管理 async 上下文，Try 管理錯誤

### 場景 3：頁面載入時立即執行

**需求：** 頁面載入時立即執行異步初始化代碼

**解決方案：**
```go
Script(nil,
    js.AsyncDo(
        js.Try(
            js.Const("config", "await fetch('/api/config').then(r => r.json())"),
            JSAction{Code: "window.appConfig = config"},
        ).Catch(
            js.Log("'初始化失敗'"),
        ).End(),
    ),
)
```

**為什麼這樣設計：**
- AsyncDo 明確表示「立即執行這段異步代碼」
- 語義清晰：Do = 立即做
- 與 AsyncFn 語義對稱

### 場景 4：創建獨立作用域

**需求：** 避免變數污染全局作用域

**解決方案：**
```go
js.Do(
    js.Const("temp", "calculateSomething()"),
    js.Log("temp"),
    // temp 不會污染外部作用域
)
```

**為什麼這樣設計：**
- Do 可以用於任何需要 IIFE 的場景，不僅限於錯誤處理
- 通用性強

## 與其他設計的對比

### 方案 A：Try 自動包裝（舊設計）

```go
js.Try(
    js.Const("data", "await fetch('/api')"),
).Catch(...).End()

// 自動生成：
// (async () => {
//   try { const data = await fetch('/api'); }
//   catch (error) { ... }
// })()
```

**問題：**
- ❌ 不夠靈活：無法在已有的 async 函數中使用
- ❌ 隱式行為：看不出會創建 IIFE
- ❌ 職責不清：Try 同時管錯誤和 IIFE
- ❌ 限制了使用場景

### 方案 B：Try 不包裝 + Do/AsyncDo（新設計）

```go
// 事件處理器
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(...).End(),
)

// 立即執行
js.AsyncDo(
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch(...).End(),
)
```

**優勢：**
- ✅ 靈活：可以在任何地方使用 Try
- ✅ 明確：看代碼就知道是否有函數包裝
- ✅ 職責清晰：Try 管錯誤，AsyncFn/AsyncDo 管 async
- ✅ 可組合：各個部分可以自由組合

### 方案 C：區分 Try 和 AsyncTry

```go
js.Try(...)       // 同步 try-catch
js.AsyncTry(...)  // 異步 try-catch（自動包裝）
```

**為什麼不採用：**
- ❌ 語義重複：Try 本身不應該關心同步/異步
- ❌ 不夠靈活：AsyncTry 仍然有自動包裝的問題
- ❌ API 複雜度增加：需要記住兩個不同的函數

## 設計原則

### 1. 單一職責原則

每個函數只負責一件事：
- `Try`：錯誤處理
- `AsyncFn`：創建異步函數
- `Do`：創建 IIFE
- `AsyncDo`：創建異步 IIFE

### 2. 明確優於隱式

```go
// ✅ 明確：看到 AsyncFn 就知道是 async
js.AsyncFn(nil, js.Try(...).Catch(...).End())

// ❌ 隱式：不知道 Try 會創建 async IIFE
js.Try(...).Catch(...).End()
```

### 3. 可組合性

所有元件都可以自由組合：
```go
// Try 在 AsyncFn 中
js.AsyncFn(nil, js.Try(...).Catch(...).End())

// Try 在 AsyncDo 中
js.AsyncDo(js.Try(...).Catch(...).End())

// Try 在普通 Fn 中（不使用 await）
js.Fn(nil, js.Try(...).Catch(...).End())

// AsyncFn 可以包含多個 Try
js.AsyncFn(nil,
    js.Try(...).Catch(...).End(),
    js.Log("'中間步驟'"),
    js.Try(...).Catch(...).End(),
)
```

### 4. 最小驚訝原則

用戶看到代碼應該能立即理解其行為：
- `Try` 生成 `try-catch-finally`
- `Do` 生成 `(() => {...})()`
- `AsyncDo` 生成 `(async () => {...})()`
- `AsyncFn` 生成 `async () => {...}`

沒有隱藏的行為，沒有意外的包裝。

## 實踐建議

### 何時使用 Try

- 在任何需要錯誤處理的地方
- 同步或異步都適用
- 可以嵌套使用

### 何時使用 AsyncFn

- 事件處理器需要 async/await
- 創建可重用的異步函數
- 作為回調函數傳遞

### 何時使用 AsyncDo

- 頁面載入時立即執行異步代碼
- 在 `<script>` 標籤中執行一次性的異步初始化
- 不需要事件觸發，直接執行

### 何時使用 Do

- 創建獨立作用域
- 避免變數污染
- 立即執行一組同步操作

## 總結

新的設計通過**職責分離**和**明確性**，提供了更靈活、更清晰的 API：

1. **Try 專注於錯誤處理**：生成純粹的 try-catch-finally
2. **Do/AsyncDo 專注於 IIFE**：明確的立即執行語義
3. **AsyncFn 專注於異步函數**：清晰的 async 邊界
4. **可自由組合**：各個部分可以靈活組合使用

這種設計雖然需要用戶多寫一點代碼（例如顯式使用 AsyncFn），但換來的是：
- 更高的靈活性
- 更清晰的代碼意圖
- 更少的意外行為
- 更好的可維護性

**設計哲學：給用戶完全的控制權，而不是做隱式的「聰明」決定。**