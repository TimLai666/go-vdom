# Catch 錯誤物件名稱自訂

## 概述

`Try().Catch()` 方法現在允許自訂錯誤物件的變數名稱，讓你的代碼更符合個人或團隊的編碼風格。

## 語法

```go
Try(actions...).Catch(errorName string, actions...).End()
```

**參數：**
- `errorName`: 錯誤物件的變數名稱（如 `"err"`, `"e"`, `"error"`, `"exception"` 等）
- 傳入空字串 `""` 時，使用預設名稱 `"error"`

## 使用示例

### 示例 1：使用簡短名稱 "err"

```go
Button(Props{
    "onClick": jsdsl.AsyncFn(nil,
        jsdsl.Try(
            jsdsl.Const("response", "await fetch('/api/data')"),
            jsdsl.Const("data", "await response.json()"),
            jsdsl.Log("data"),
        ).Catch("err",  // 使用 "err" 作為錯誤物件名稱
            jsdsl.Log("err.message"),
            jsdsl.Alert("'請求失敗: ' + err.message"),
        ).End(),
    ),
}, "獲取數據")
```

**生成的 JavaScript：**
```javascript
async () => {
  try {
    const response = await fetch('/api/data');
    const data = await response.json();
    console.log(data);
  } catch (err) {
    console.log(err.message);
    alert('請求失敗: ' + err.message);
  }
}
```

### 示例 2：使用單字母 "e"

```go
jsdsl.Try(
    jsdsl.Const("data", "JSON.parse(input)"),
).Catch("e",  // 使用 "e"
    jsdsl.Log("'解析錯誤:', e"),
).End()
```

**生成的 JavaScript：**
```javascript
try {
  const data = JSON.parse(input);
} catch (e) {
  console.log('解析錯誤:', e);
}
```

### 示例 3：使用預設名稱 "error"

```go
jsdsl.Try(
    jsdsl.Const("result", "riskyOperation()"),
).Catch("",  // 空字串使用預設名稱
    jsdsl.Log("error.message"),
).End()
```

**生成的 JavaScript：**
```javascript
try {
  const result = riskyOperation();
} catch (error) {
  console.log(error.message);
}
```

### 示例 4：使用描述性名稱 "exception"

```go
jsdsl.Try(
    jsdsl.Const("file", "fs.readFileSync('data.txt')"),
).Catch("exception",  // 使用 "exception"
    jsdsl.Log("'文件讀取錯誤:', exception"),
    jsdsl.Const("fallback", "''"),
).Finally(
    jsdsl.Log("'操作結束'"),
)
```

**生成的 JavaScript：**
```javascript
try {
  const file = fs.readFileSync('data.txt');
} catch (exception) {
  console.log('文件讀取錯誤:', exception);
  const fallback = '';
} finally {
  console.log('操作結束');
}
```

### 示例 5：表單提交錯誤處理

```go
Form(Props{
    "onSubmit": jsdsl.Fn([]string{"event"},
        jsdsl.CallMethod("event", "preventDefault"),
        jsdsl.AsyncDo(
            jsdsl.Try(
                jsdsl.Const("formData", "new FormData(event.target)"),
                jsdsl.Const("response", "await fetch('/api/submit', { method: 'POST', body: formData })"),
                jsdsl.If("!response.ok",
                    JSAction{Code: "throw new Error('HTTP ' + response.status)"},
                ),
                jsdsl.Const("result", "await response.json()"),
                jsdsl.Alert("'提交成功!'"),
            ).Catch("submitError",  // 描述性的錯誤名稱
                jsdsl.Log("'提交失敗:', submitError"),
                jsdsl.Alert("'提交失敗: ' + submitError.message"),
                jsdsl.CallMethod("event.target", "reset"),
            ).End(),
        ),
    ),
})
```

## 常用錯誤名稱

不同的錯誤物件名稱適用於不同的場景：

| 名稱 | 適用場景 | 說明 |
|------|----------|------|
| `"err"` | 一般錯誤處理 | 簡短、常用 |
| `"e"` | 簡單的 catch | 最簡短的形式 |
| `"error"` | 預設/標準 | 傳空字串時的預設值 |
| `"exception"` | Java 風格 | 更正式的命名 |
| `"httpError"` | HTTP 請求 | 描述性命名 |
| `"parseError"` | 解析錯誤 | 描述性命名 |
| `"dbError"` | 資料庫錯誤 | 描述性命名 |

## 最佳實踐

### ✅ 推薦做法

1. **團隊統一風格**
```go
// 團隊約定使用 "err"
.Catch("err", 
    jsdsl.Log("err"),
)
```

2. **使用描述性名稱（複雜場景）**
```go
// API 錯誤
.Catch("apiError",
    jsdsl.Log("apiError"),
)

// 解析錯誤
.Catch("parseError",
    jsdsl.Log("parseError"),
)
```

3. **簡單場景使用短名稱**
```go
// 簡單的錯誤處理
.Catch("e",
    jsdsl.Log("e"),
)
```

### ❌ 避免的做法

1. **不要使用保留字**
```go
.Catch("catch", ...)   // ✗ 不要使用 JavaScript 保留字
.Catch("function", ...)  // ✗
```

2. **避免誤導性名稱**
```go
.Catch("success", ...)  // ✗ 錯誤物件不應該叫 success
.Catch("data", ...)     // ✗ 可能與正常數據混淆
```

3. **保持一致性**
```go
// ✗ 在同一個文件中混用不同風格
jsdsl.Try(...).Catch("err", ...)
jsdsl.Try(...).Catch("e", ...)
jsdsl.Try(...).Catch("error", ...)

// ✓ 保持一致
jsdsl.Try(...).Catch("err", ...)
jsdsl.Try(...).Catch("err", ...)
jsdsl.Try(...).Catch("err", ...)
```

## 向後兼容性

**重要更新：** 舊版本的 `Catch()` 方法簽名已改變。

### 舊版本（已過時）

```go
.Catch(
    jsdsl.Log("error.message"),  // ✗ 不再支援
)
```

### 新版本（必須提供錯誤名稱）

```go
.Catch("err",  // ✓ 必須提供錯誤名稱
    jsdsl.Log("err.message"),
)

// 或使用預設名稱
.Catch("",  // ✓ 空字串使用預設 "error"
    jsdsl.Log("error.message"),
)
```

## 遷移指南

如果你的代碼使用了舊版本的 `Catch()`，需要進行以下修改：

### 步驟 1：添加錯誤名稱參數

```go
// 舊代碼
jsdsl.Try(
    jsdsl.Const("x", "1"),
).Catch(
    jsdsl.Log("error"),
)

// 新代碼
jsdsl.Try(
    jsdsl.Const("x", "1"),
).Catch("error",  // 添加錯誤名稱
    jsdsl.Log("error"),
)
```

### 步驟 2：更新錯誤物件引用（如需要）

```go
// 如果你想使用不同的名稱
jsdsl.Try(
    jsdsl.Const("x", "1"),
).Catch("err",  // 使用新名稱
    jsdsl.Log("err"),  // 更新引用
)
```

## 總結

- ✅ 第一個參數指定錯誤物件名稱
- ✅ 支援任意合法的 JavaScript 變數名
- ✅ 空字串使用預設名稱 `"error"`
- ✅ 讓代碼風格更靈活、更統一
- ✅ 支援描述性命名，提高可讀性

---

**文檔版本**: 1.0.0  
**最後更新**: 2025-01-24  
**作者**: TimLai666
