# TryCatch 重新設計說明

## 📋 概述

TryCatch 已完全重新設計，從基於函數包裝的 API 改為基於動作列表的 API，使其更直觀且完全支持異步操作。

---

## ⚠️ 破壞性變更

**舊 API（v1.0.x）已移除：**
```go
// ❌ 舊的 API（不再支持）
js.TryCatch(
    js.AsyncFn(nil, ...),           // baseAction
    js.Ptr(js.Fn(nil, ...)),        // catchFn
    js.Ptr(js.Fn(nil, ...)),        // finallyFn
)
```

**新 API（v1.1.0+）：**
```go
// ✅ 新的 API
js.TryCatch(
    []JSAction{...},  // tryActions
    []JSAction{...},  // catchActions
    []JSAction{...},  // finallyActions (可選)
)
```

---

## 🎯 為什麼要重新設計？

### 舊設計的問題

1. **無法正確處理異步操作**
   ```go
   // 問題：內部的 AsyncFn 只會定義函數，不會執行
   js.TryCatch(
       js.AsyncFn(nil,
           js.Const("data", "await fetch('/api')"),  // 永遠不會執行！
       ),
       js.Ptr(js.Fn(nil, ...)),
       nil,
   )
   ```

2. **API 不直觀**
   - 需要使用 `js.Ptr()` 包裝
   - 需要嵌套 `AsyncFn` 或 `Fn`
   - 概念混亂（為什麼要包裝函數？）

3. **用戶經常犯錯**
   - 忘記使用 AsyncFn
   - 忘記使用 Ptr
   - 不理解為什麼代碼沒有執行

### 新設計的優勢

1. **✅ 完全支持異步操作**
   ```go
   js.TryCatch(
       []JSAction{
           js.Const("data", "await fetch('/api')"),  // ✅ 正確執行
       },
       []JSAction{
           js.Log("'錯誤:', e.message"),
       },
       nil,
   )
   ```

2. **✅ API 簡潔直觀**
   - 直接傳入動作列表
   - 不需要 Ptr 包裝
   - 不需要嵌套函數
   - 符合直覺的設計

3. **✅ 自動處理細節**
   - 自動創建 async 函數包裝
   - 自動添加分號
   - 自動格式化代碼
   - 錯誤對象自動命名為 `e`

---

## 📚 完整 API 說明

### 函數簽名

```go
func TryCatch(
    tryActions []JSAction,      // try 區塊中的動作列表
    catchActions []JSAction,    // catch 區塊中的動作列表（可選）
    finallyActions []JSAction,  // finally 區塊中的動作列表（可選）
) JSAction
```

### 參數說明

- **`tryActions`**: try 區塊中執行的動作列表，支持 `await` 語法
- **`catchActions`**: catch 區塊中執行的動作列表，可以訪問錯誤對象 `e`
- **`finallyActions`**: finally 區塊中執行的動作列表（可選）

**要求：** `catchActions` 和 `finallyActions` 至少需要提供一個。

### 生成的代碼

```javascript
(async () => {
  try {
    // tryActions 的語句
  } catch (e) {
    // catchActions 的語句
  } finally {
    // finallyActions 的語句
  }
})()
```

---

## 🔄 遷移指南

### 基本用法

**舊代碼：**
```go
js.TryCatch(
    js.AsyncFn(nil,
        js.Const("response", "await fetch('/api')"),
        js.Const("data", "await response.json()"),
        js.Log("data"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'錯誤:', e.message"),
        js.Alert("'發生錯誤'"),
    )),
    nil,
)
```

**新代碼：**
```go
js.TryCatch(
    []JSAction{
        js.Const("response", "await fetch('/api')"),
        js.Const("data", "await response.json()"),
        js.Log("data"),
    },
    []JSAction{
        js.Log("'錯誤:', e.message"),
        js.Alert("'發生錯誤'"),
    },
    nil,
)
```

### 帶 finally 的用法

**舊代碼：**
```go
js.TryCatch(
    js.AsyncFn(nil,
        js.Const("result", "await doSomething()"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'錯誤'"),
    )),
    js.Ptr(js.Fn(nil,
        js.Log("'清理資源'"),
    )),
)
```

**新代碼：**
```go
js.TryCatch(
    []JSAction{
        js.Const("result", "await doSomething()"),
    },
    []JSAction{
        js.Log("'錯誤'"),
    },
    []JSAction{
        js.Log("'清理資源'"),
    },
)
```

### 在事件處理器中使用

**舊代碼：**
```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.TryCatch(
            js.AsyncFn(nil,
                js.Const("data", "await fetch('/api')"),
            ),
            js.Ptr(js.Fn(nil,
                js.Alert("'錯誤'"),
            )),
            nil,
        ),
    ),
}, "點擊")
```

**新代碼：**
```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.TryCatch(
            []JSAction{
                js.Const("data", "await fetch('/api')"),
            },
            []JSAction{
                js.Alert("'錯誤'"),
            },
            nil,
        ),
    ),
}, "點擊")
```

---

## 💡 使用示例

### 示例 1：基本異步操作

```go
js.AsyncFn(nil,
    js.Log("'開始操作...'"),
    js.TryCatch(
        []JSAction{
            js.Const("response", "await fetch('/api/data')"),
            js.Log("'收到響應'"),
            JSAction{Code: "if (!response.ok) throw new Error('請求失敗')"},
            js.Const("data", "await response.json()"),
            js.Log("'數據:', data"),
        },
        []JSAction{
            js.Log("'錯誤:', e.message"),
            js.Alert("'操作失敗: ' + e.message"),
        },
        nil,
    ),
)
```

### 示例 2：表單提交

```go
Form(Props{
    "onSubmit": js.AsyncFn([]string{"event"},
        js.CallMethod("event", "preventDefault"),
        js.TryCatch(
            []JSAction{
                js.Const("formData", "new FormData(event.target)"),
                js.Const("response", "await fetch('/api/submit', { method: 'POST', body: formData })"),
                JSAction{Code: "if (!response.ok) throw new Error('提交失敗')"},
                js.Const("result", "await response.json()"),
                js.Alert("'提交成功: ' + result.message"),
                JSAction{Code: "event.target.reset()"},
            },
            []JSAction{
                js.Alert("'提交失敗: ' + e.message"),
            },
            nil,
        ),
    ),
})
```

### 示例 3：帶 finally 的資源清理

```go
js.AsyncFn(nil,
    js.TryCatch(
        []JSAction{
            js.Const("file", "await openFile('data.txt')"),
            js.Const("content", "await file.read()"),
            js.Log("'內容:', content"),
        },
        []JSAction{
            js.Log("'讀取失敗:', e.message"),
        },
        []JSAction{
            js.Log("'關閉文件'"),
            JSAction{Code: "if (file) file.close()"},
        },
    ),
)
```

### 示例 4：API 數據載入並渲染

```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Const("container", "document.getElementById('result')"),
        JSAction{Code: "container.innerHTML = '載入中...'"},
        js.TryCatch(
            []JSAction{
                js.Const("response", "await fetch('/api/items')"),
                js.Const("items", "await response.json()"),
                JSAction{Code: "container.innerHTML = ''"},
                js.Const("ul", "document.createElement('ul')"),
                js.ForEachJS("items", "item",
                    js.Const("li", "document.createElement('li')"),
                    JSAction{Code: "li.textContent = item.name"},
                    JSAction{Code: "ul.appendChild(li)"},
                ),
                JSAction{Code: "container.appendChild(ul)"},
            },
            []JSAction{
                JSAction{Code: "container.innerHTML = '載入失敗: ' + e.message"},
            },
            nil,
        ),
    ),
}, "載入數據")
```

---

## 🎯 最佳實踐

### 1. 外層使用 AsyncFn，內部使用 TryCatch

```go
// ✅ 推薦
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.TryCatch(
            []JSAction{
                js.Const("data", "await fetch('/api')"),
            },
            []JSAction{
                js.Log("'錯誤:', e.message"),
            },
            nil,
        ),
    ),
})
```

### 2. 始終提供錯誤處理

```go
// ✅ 好的做法
js.TryCatch(
    []JSAction{...},
    []JSAction{
        js.Log("'錯誤:', e.message"),
        js.Alert("'操作失敗'"),
    },
    nil,
)

// ❌ 不好的做法（沒有錯誤處理）
js.AsyncFn(nil,
    js.Const("data", "await fetch('/api')"),
    // 沒有錯誤處理，失敗時用戶不知道發生什麼
)
```

### 3. 使用 finally 進行清理

```go
// ✅ 使用 finally 確保清理代碼執行
js.TryCatch(
    []JSAction{
        js.Const("loading", "true"),
        js.Const("data", "await fetch('/api')"),
    },
    []JSAction{
        js.Log("'錯誤:', e.message"),
    },
    []JSAction{
        JSAction{Code: "loading = false"},
        js.Log("'操作完成'"),
    },
)
```

### 4. 檢查響應狀態

```go
// ✅ 檢查 HTTP 狀態
js.TryCatch(
    []JSAction{
        js.Const("response", "await fetch('/api')"),
        JSAction{Code: "if (!response.ok) throw new Error('HTTP ' + response.status)"},
        js.Const("data", "await response.json()"),
    },
    []JSAction{
        js.Alert("'錯誤: ' + e.message"),
    },
    nil,
)
```

---

## 🔍 常見問題

### Q: 為什麼要破壞性變更？

A: 舊的 API 存在根本性設計問題，無法正確處理異步操作。新的 API 更直觀且功能完整。

### Q: 如何快速遷移代碼？

A: 
1. 找到所有 `js.TryCatch(` 使用
2. 將 `js.AsyncFn(nil, ...)` 改為 `[]JSAction{...}`
3. 將 `js.Ptr(js.Fn(nil, ...))` 改為 `[]JSAction{...}`
4. 移除 `js.Ptr()` 和內部的 `js.Fn()`

### Q: 可以不用 AsyncFn 包裝嗎？

A: TryCatch 會立即執行並返回，所以通常需要在外層用 AsyncFn 或 onClick 等事件處理器包裝。

### Q: 錯誤對象的名稱是什麼？

A: 錯誤對象自動命名為 `e`，可以在 catchActions 中直接使用。

### Q: finally 是必須的嗎？

A: 不是。catchActions 和 finallyActions 至少提供一個即可。

---

## 📖 相關文檔

- [API 參考 - TryCatch](docs/API_REFERENCE.md#trycatch)
- [快速參考 - 異步操作](docs/QUICK_REFERENCE.md)
- [完整文檔 - 錯誤處理](docs/DOCUMENTATION.md)

---

## 🎉 總結

新的 TryCatch API：

- ✅ 完全支持異步操作
- ✅ API 簡潔直觀
- ✅ 自動處理細節
- ✅ 更容易理解和使用
- ✅ 不容易出錯

**立即升級到 v1.1.0，享受更好的異步錯誤處理！**

---

**版本**: v1.1.0  
**更新日期**: 2025-01-24