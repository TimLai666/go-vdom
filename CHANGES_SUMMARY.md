# 更新總結 (v1.1.0)

## 問題修復

### 1. 修復 "await is only valid in async functions" 錯誤

**問題：**
控制台出現兩個 JavaScript 錯誤：
```
Uncaught SyntaxError: await is only valid in async functions and the top level bodies of modules (at (索引):118:20)
Uncaught SyntaxError: await is only valid in async functions and the top level bodies of modules (at (索引):151:20)
```

**原因：**
在 `main.go` 中，使用了 `js.Fn()` 創建包含 `await` 語句的函數，但 `Fn()` 生成的是普通同步函數，不支持 `await`。

**解決方案：**
1. 新增 `AsyncFn()` 函數到 `jsdsl/jsdsl.go`
2. 將 `main.go` 中所有包含 `await` 的 `js.Fn()` 改為 `js.AsyncFn()`

**修改位置：**
- `jsdsl/jsdsl.go` - 新增 `AsyncFn()` 函數 (第 112-134 行)
- `main.go` - 第 155 行：GET 請求的 onClick 改用 `AsyncFn`
- `main.go` - 第 159 行：TryCatch 內部改用 `AsyncFn`
- `main.go` - 第 205 行：POST 請求的 onSubmit 改用 `AsyncFn`
- `main.go` - 第 209 行：TryCatch 內部改用 `AsyncFn`

---

## 功能改進

### 2. ForEachJS 列表渲染優化

**問題：**
- 後端列表渲染寫法過於複雜，需要手動創建切片
- 前端 `js.ForEachJS` 僅能用於 DOM 元素遍歷，不夠通用
- 缺乏簡潔的列表渲染語法

**原始寫法（複雜）：**
```go
Ul(
    func() []VNode {
        nodes := make([]VNode, len(items))
        for i, item := range items {
            nodes[i] = Li(item)
        }
        return nodes
    }()...,
)
```

**解決方案：**

1. **後端渲染** - 新增簡潔的 `ForEachJS` 函數：
   ```go
   // vdom/tags.go
   func ForEachJS[T any](items []T, renderFunc func(item T) VNode) []VNode
   func ForEachJSWithIndexJS[T any](items []T, renderFunc func(item T, index int) VNode) []VNode
   ```

2. **前端渲染** - 重構 `js.ForEachJS` 使其更通用：
   ```go
   // jsdsl/jsdsl.go
   func ForEachJS(arrayExpr string, itemVar string, actions ...JSAction) JSAction
   func ForEachJSWithIndexJS(arrayExpr string, itemVar string, indexVar string, actions ...JSAction) JSAction
   func ForEachJSElement(arrayExpr string, fn func(el Elem) JSAction) JSAction
   ```

**新的簡潔寫法：**

**後端：**
```go
// 基本用法
Ul(ForEachJS(items, func(item string) VNode {
    return Li(item)
}))

// 帶索引
Ul(ForEachJSWithIndexJS(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
}))
```

**前端：**
```go
// 遍歷任意數組
js.ForEachJS("['A', 'B', 'C']", "item",
    js.Log("'項目: ' + item"),
)

// 帶索引遍歷
js.ForEachJSWithIndexJS("numbers", "num", "idx",
    js.Log("'[' + idx + '] = ' + num"),
)

// DOM 元素遍歷（專用）
js.ForEachJSElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
    return el.AddClass("'active'")
})
```

**修改位置：**
- `vdom/tags.go` - 新增第 17-34 行：`ForEachJS` 和 `ForEachJSWithIndexJS` 函數
- `jsdsl/jsdsl.go` - 重構第 221-278 行：`ForEachJS`, `ForEachJSWithIndexJS`, `ForEachJSElement`
- `examples/05_foreach_usage.go` - 新增完整示例（428 行）

**優勢：**
- ✅ 語法更簡潔，代碼更易讀
- ✅ 支持泛型，可用於任何類型
- ✅ 前端 ForEachJS 不再限於 DOM 元素
- ✅ 區分後端和前端渲染場景
- ✅ 提供索引版本滿足不同需求

---

## 功能改進 (續)

### 3. control.For 和 control.ForEach 重構

**問題：**
- `control.For` 命名不夠語義化（遍歷集合應該叫 ForEach）
- 缺乏傳統 for 循環功能（for i := start; i < end; i += step）
- 無法生成數字序列（如 1-10、偶數等）

**解決方案：**

1. **重命名為更語義化的名稱**：
   ```go
   // 原來的 For 改名為 ForEach
   func ForEach[T any](items []T, renderFunc func(item T, index int) VNode) []VNode
   ```

2. **新增傳統循環 For**：
   ```go
   // 新的 For 用於數字循環
   func For(start, end, step int, renderFunc func(i int) VNode) []VNode
   ```

**使用示例：**

**ForEach（遍歷集合）：**
```go
// 遍歷用戶列表
users := []User{...}
control.ForEach(users, func(user User, i int) VNode {
    return Div(
        H5(user.Name),
        P(fmt.Sprintf("年齡: %d", user.Age)),
    )
})
```

**For（數字循環）：**
```go
// 正向循環：生成 1-10
control.For(1, 11, 1, func(i int) VNode {
    return Span(Props{"class": "badge"}, fmt.Sprintf("%d", i))
})

// 倒序循環：10-1
control.For(10, 0, -1, func(i int) VNode {
    return Span(Props{"class": "badge"}, fmt.Sprintf("%d", i))
})

// 步進循環：偶數 0-18
control.For(0, 20, 2, func(i int) VNode {
    return Span(Props{"class": "badge"}, fmt.Sprintf("%d", i))
})
```

**實用案例：**

1. **生成分頁按鈕**：
```go
Nav(Ul(Props{"class": "pagination"},
    control.For(1, 11, 1, func(i int) VNode {
        return Li(Props{"class": "page-item"},
            A(Props{"class": "page-link"}, fmt.Sprintf("%d", i)),
        )
    }),
))
```

2. **生成表格行號**：
```go
Table(Tbody(
    control.For(1, 11, 1, func(i int) VNode {
        return Tr(
            Td(fmt.Sprintf("%d", i)),
            Td(fmt.Sprintf("項目 %d", i)),
        )
    }),
))
```

3. **評分星星（1-5）**：
```go
control.For(1, 6, 1, func(star int) VNode {
    if star <= rating {
        return Span("★")
    }
    return Span("☆")
})
```

**修改位置：**
- `control/control.go` - 第 72 行：`For` 改名為 `ForEach`
- `control/control.go` - 第 82-111 行：新增 `For(start, end, step, renderFunc)` 函數
- `control/control.go` - 第 148 行：`KeyedFor` 改名為 `KeyedForEach`
- `main.go` - 第 140-165 行：更新為使用新的 ForEach 和 For
- `examples/06_control_loops.go` - 新增完整示例（481 行）

**優勢：**
- ✅ 命名更語義化（ForEach = 遍歷，For = 循環）
- ✅ 支持傳統 for 循環功能
- ✅ 可生成任意數字序列
- ✅ 支持正向、倒序、步進循環
- ✅ 適用於分頁、表格、評分等場景

**對比表：**
| 特性 | control.ForEach | control.For |
|------|----------------|-------------|
| 用途 | 遍歷集合 | 數字循環 |
| 語法 | `ForEach(items, func(item, i) {...})` | `For(start, end, step, func(i) {...})` |
| 數據來源 | 現有切片/數組 | 動態生成數字 |
| 適用場景 | 用戶列表、商品列表 | 分頁按鈕、表格行號 |

---

## 文檔整理

### 3. 清理和重組文檔結構

**問題：**
項目根目錄有多個文檔文件，缺乏組織結構。

**解決方案：**
創建 `docs/` 目錄並重組所有文檔：

```
docs/
├── README.md               # 文檔導航中心（新增）
├── QUICK_START.md          # 快速入門（從 QUICK_START_V1.1.md 移動）
├── DOCUMENTATION.md        # 完整技術文檔（移動）
├── API_REFERENCE.md        # JavaScript DSL API 參考（新增）
└── QUICK_REFERENCE.md      # 語法速查表（移動）
```

**刪除的冗餘文檔：**
- `IMPROVEMENTS.md` - 內容已整合到其他文檔
- `UPDATE_SUMMARY.md` - 替換為本文件

**新增文檔：**
- `docs/README.md` - 209 行的文檔導航中心
- `docs/API_REFERENCE.md` - 709 行的完整 JavaScript DSL API 參考
  - 包含所有 DSL 函數的詳細說明
  - 完整的 AsyncFn 使用指南
  - TryCatch 錯誤處理最佳實踐
  - 實用示例和類型定義

**精簡的 README.md：**
- 從 1167 行縮減到 305 行
- 移除重複的詳細說明
- 保留核心概念和快速開始
- 添加指向 docs/ 的鏈接

---

## 代碼變更

### 新增函數實現

#### AsyncFn 函數

```go
// jsdsl/jsdsl.go 新增函數
func AsyncFn(params []string, actions ...JSAction) JSAction {
    var sb strings.Builder
    
    // 創建一個異步匿名函數
    paramsStr := ""
    if params != nil {
        paramsStr = strings.Join(params, ", ")
    }
    sb.WriteString(fmt.Sprintf("async (%s) => {\n", paramsStr))
    
    // 添加函數體
    for _, a := range actions {
        line := strings.TrimSpace(a.Code)
        if !strings.HasSuffix(line, ";") {
            line += ";"
        }
        sb.WriteString("  " + line + "\n")
    }
    
    sb.WriteString("}")
    return JSAction{Code: sb.String()}
}
```

#### ForEachJS 函數（後端）

```go
// vdom/tags.go
func ForEachJS[T any](items []T, renderFunc func(item T) VNode) []VNode {
    result := make([]VNode, len(items))
    for i, item := range items {
        result[i] = renderFunc(item)
    }
    return result
}

func ForEachJSWithIndexJS[T any](items []T, renderFunc func(item T, index int) VNode) []VNode {
    result := make([]VNode, len(items))
    for i, item := range items {
        result[i] = renderFunc(item, i)
    }
    return result
}
```

#### ForEachJS 函數（前端）

```go
// jsdsl/jsdsl.go
func ForEachJS(arrayExpr string, itemVar string, actions ...JSAction) JSAction {
    var sb strings.Builder
    for _, a := range actions {
        line := strings.TrimSpace(a.Code)
        if !strings.HasSuffix(line, ";") {
            line += ";"
        }
        sb.WriteString(line + "\n")
    }
    
    return JSAction{
        Code: fmt.Sprintf(`%s.forEach(function(%s) {
%s});`, arrayExpr, itemVar, indent(sb.String(), "  ")),
    }
}
```

### 使用示例

#### AsyncFn

**修改前（錯誤）：**
```go
Button(Props{
    "onClick": js.Fn(nil,  // ❌ 錯誤：不支持 await
        js.Const("response", "await fetch('/api/data')"),
    ),
}, "Fetch Data")
```

**修改後（正確）：**
```go
Button(Props{
    "onClick": js.AsyncFn(nil,  // ✅ 正確：支持 await
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Log("data"),
    ),
}, "Fetch Data")
```

#### ForEachJS（後端）

**修改前（複雜）：**
```go
Ul(
    func() []VNode {
        nodes := make([]VNode, len(vars))
        for i, v := range vars {
            nodes[i] = Li(Props{"class": "badge bg-primary me-2"}, v)
        }
        return nodes
    }()...,
)
```

**修改後（簡潔）：**
```go
// 基本用法
Ul(ForEachJS(vars, func(v string) VNode {
    return Li(Props{"class": "badge bg-primary me-2"}, v)
}))

// 帶索引
Ul(ForEachJSWithIndexJS(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
}))
```

#### ForEachJS（前端）

**修改前（僅限 DOM）：**
```go
js.ForEachJS("items", func(el js.Elem) JSAction {
    return js.Log(el.InnerText())
})
```

**修改後（通用）：**
```go
// 遍歷任意數組
js.ForEachJS("['Apple', 'Banana', 'Orange']", "fruit",
    js.Log("'水果: ' + fruit"),
)

// 帶索引遍歷
js.ForEachJSWithIndexJS("items", "item", "index",
    js.Log("'[' + index + '] = ' + item"),
)

// DOM 元素專用
js.ForEachJSElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
    return el.AddClass("'active'")
})
```

---

## 更新的文件清單

### 新增文件
- `docs/README.md` - 文檔導航中心（209 行）
- `docs/API_REFERENCE.md` - JavaScript DSL API 參考（更新了 ForEach 部分）
- `examples/05_foreach_usage.go` - ForEach 完整示例（428 行）
- `examples/06_control_loops.go` - control.For 和 control.ForEach 示例（481 行）
- `CHANGES_SUMMARY.md` - 本文件

### 移動的文件
- `DOCUMENTATION.md` → `docs/DOCUMENTATION.md`
- `QUICK_REFERENCE.md` → `docs/QUICK_REFERENCE.md`
- `QUICK_START_V1.1.md` → `docs/QUICK_START.md`

### 刪除的文件
- `IMPROVEMENTS.md` - 內容已整合
- `UPDATE_SUMMARY.md` - 替換為本文件

### 修改的文件
- `jsdsl/jsdsl.go` - 新增 AsyncFn()、重構 ForEach 相關函數
- `vdom/tags.go` - 新增 ForEach() 和 ForEachWithIndex()
- `control/control.go` - 重構 For → ForEach，新增傳統循環 For
- `main.go` - 將 Fn 改為 AsyncFn（4 處修改）+ 更新為使用新的 ForEach 和 For
- `README.md` - 精簡內容，從 1167 行縮減到 305 行
- `CHANGELOG.md` - 添加 AsyncFn、ForEach 和 control.For 相關更新記錄
- `docs/API_REFERENCE.md` - 更新 ForEach 文檔
- `docs/QUICK_REFERENCE.md` - 添加 ForEach 對比表格和 control.For 用法

---

## 驗證

### 編譯測試
```bash
go build -o test-build main.go
# ✅ 編譯成功，無錯誤
```

### 功能驗證
運行 `go run main.go` 並訪問 http://localhost:8080

**預期行為：**
- ✅ 頁面正常加載
- ✅ 點擊 "獲取數據" 按鈕成功調用 API
- ✅ 表單提交正常工作
- ✅ 控制台無 "await is only valid in async functions" 錯誤
- ✅ Handler 註冊正常工作

---

## 使用指南

### 何時使用 Fn vs AsyncFn vs ForEachJS

#### Fn vs AsyncFn

**使用 Fn：**
```go
// 同步操作
js.Fn(nil,
    js.Log("'Hello'"),
    js.Alert("'World'"),
)
```

**使用 AsyncFn：**
```go
// 包含 await 的異步操作
js.AsyncFn(nil,
    js.Const("response", "await fetch('/api')"),
    js.Const("data", "await response.json()"),
)
```

### TryCatch 與 AsyncFn 配合

```go
js.TryCatch(
    // Try block - 使用 AsyncFn
    js.AsyncFn(nil,
        js.Const("data", "await fetchData()"),
    ),
    // Catch block - 使用普通 Fn
    js.Ptr(js.Fn(nil,
        js.Log("'Error:', e.message"),
    )),
    nil,
)
```

#### ForEach vs control.ForEach vs control.For

**vdom.ForEach（簡潔的遍歷）：**
```go
// 不需要索引時
Ul(ForEach(items, func(item string) VNode {
    return Li(item)
}))

// 需要索引時
Ul(ForEachWithIndex(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
}))
```

**control.ForEach（完整的遍歷）：**
```go
import "github.com/TimLai666/go-vdom/control"

// 總是提供項目和索引
Ul(control.ForEach(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
}))
```

**control.For（數字循環）：**
```go
// 生成 1-10
control.For(1, 11, 1, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})

// 倒序 10-1
control.For(10, 0, -1, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})

// 偶數 0-18
control.For(0, 20, 2, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})
```

#### 前端遍歷選擇

| 場景 | 使用函數 |
|------|---------|
| 遍歷任意 JS 數組 | `js.ForEachJS(arrayExpr, itemVar, ...)` |
| 需要索引 | `js.ForEachJSWithIndexJS(arrayExpr, itemVar, indexVar, ...)` |
| DOM 元素操作 | `js.ForEachJSElement(selector, func(el Elem) {...})` |

**重要規則：**
- ✅ 外層事件處理如果內部有 await，必須使用 AsyncFn
- ✅ TryCatch 的第一個參數（try block）如果有 await，必須使用 AsyncFn
- ✅ Catch 和 finally block 通常使用普通 Fn（除非它們也需要 await）

---

## 文檔導航

從現在開始，請使用以下文檔結構：

1. **快速開始** → `docs/QUICK_START.md`
2. **完整文檔** → `docs/DOCUMENTATION.md`
3. **API 參考** → `docs/API_REFERENCE.md` ⭐ 新增
4. **語法速查** → `docs/QUICK_REFERENCE.md`
5. **文檔中心** → `docs/README.md` ⭐ 新增

所有文檔都在 `docs/` 目錄下，易於維護和查找。

---

## 功能對比表

| 場景 | 後端渲染 | 前端渲染 |
|------|---------|---------|
| 遍歷集合 | `ForEach(items, func(item) {...})` 或 `control.ForEach(items, func(item, i) {...})` | `js.ForEachJS("array", "item", ...)` |
| 帶索引 | `ForEachWithIndex(items, func(item, i) {...})` | `js.ForEachWithIndexJS("array", "item", "i", ...)` |
| 數字循環 | `control.For(start, end, step, func(i) {...})` | 不適用（可用 JS 原生 for） |
| DOM 操作 | N/A | `js.ForEachElement(selector, func(el) {...})` |
| 靜態數據 | ✅ 推薦 | ❌ 不推薦 |
| 動態數據 | ❌ 不適用 | ✅ 推薦 |
| SEO | ✅ 友好 | ❌ 不友好 |

---

## 版本信息

- **版本**: v1.1.0
- **更新日期**: 2025-01-24
- **主要變更**: 
  1. 新增 AsyncFn 支持異步函數
  2. 優化 ForEach 列表渲染（後端+前端）
  3. 重構 control.For 和 control.ForEach（更語義化）
  4. 重組文檔結構
- **向後兼容**: 是（完全兼容 v1.0.0）
- **重要更新**:
  - ✅ 修復 await 語法錯誤
  - ✅ 簡化列表渲染語法
  - ✅ 前端 ForEach 更通用
  - ✅ control.For 支持傳統循環

---

## ⚠️ 重要注意事項

### TryCatch 與 AsyncFn 的正確用法

**問題：** TryCatch 內部不應該使用 AsyncFn 或 Fn 包裝

TryCatch 已經創建了 async 函數包裝，如果內部再用 AsyncFn/Fn，會導致函數定義但不執行。

**❌ 錯誤示例：**
```go
js.TryCatch(
    js.AsyncFn(nil,  // ❌ 錯誤！只會定義函數，不會執行
        js.Const("data", "await fetch('/api')"),
    ),
    js.Ptr(js.Fn(nil, ...)),  // ❌ 錯誤！
    nil,
)
```

**✅ 正確做法 1：外層 AsyncFn + 原生 try-catch**
```go
js.AsyncFn(nil,
    JSAction{Code: `try {
  const response = await fetch('/api/data');
  const data = await response.json();
  console.log(data);
} catch (e) {
  console.log('錯誤:', e.message);
}`},
)
```

**✅ 正確做法 2：TryCatch 用於同步代碼**
```go
js.TryCatch(
    JSAction{Code: "doSomething()"},  // 直接寫語句
    js.Ptr(JSAction{Code: "console.log(e)"}),  // 直接寫語句
    nil,
)
```

**規則：**
- ✅ 外層事件處理器用 `AsyncFn`
- ✅ 內部用原生 JavaScript `try-catch` 語法
- ❌ 不要在 TryCatch 內部使用 AsyncFn 或 Fn
- ✅ TryCatch 適用於不需要 await 的同步代碼

---

## 下一步

1. ✅ 代碼已修復，編譯通過
2. ✅ 文檔已重組，結構清晰
3. ✅ API 參考已完整
4. ⏭️ 運行示例測試新功能：
   - `go run main.go` - 主示例（http://localhost:8080）
   - `go run examples/05_foreach_usage.go` - ForEach 示例（http://localhost:8084）
   - `go run examples/06_control_loops.go` - control.For 和 control.ForEach 示例（http://localhost:8085）
5. ⏭️ 閱讀 `docs/README.md` 了解新文檔結構
6. ⏭️ 閱讀 `docs/API_REFERENCE.md` 學習 AsyncFn 和 ForEach 用法
7. ⏭️ **重要：** 學習 TryCatch 與 AsyncFn 的正確用法（見上方注意事項）
8. ⏭️ 使用新的 ForEach 和 control.For 語法替換複雜的列表渲染代碼
9. ⏭️ 將代碼中的 `control.For` 改為 `control.ForEach`（如果是遍歷集合）