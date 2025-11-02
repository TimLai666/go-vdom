# 最近更新

## Do/AsyncDo 參數注入改進

### 變更內容

**之前**：Do/AsyncDo 只會在參數名為 `event`、`e`、`evt`、`ev`（不區分大小寫）時才注入 event 對象。

**現在**：Do/AsyncDo 只要宣告了參數（無論參數名是什麼），就會自動注入 event 對象。

### 影響

這個改變讓開發者可以使用任意喜歡的參數名，不再受限於特定的幾個名稱：

```go
// 以下全部都能正常工作，都會注入 event 對象
js.Do([]string{"event"}, ...)      // ✅
js.Do([]string{"e"}, ...)          // ✅
js.Do([]string{"evt"}, ...)        // ✅
js.Do([]string{"myEvent"}, ...)    // ✅
js.Do([]string{"clickEvent"}, ...) // ✅
js.Do([]string{"formEvent"}, ...)  // ✅
js.Do([]string{"任意名稱"}, ...)     // ✅
```

### 代碼示例

```go
// 使用自定義的參數名
Button(Props{
    "onClick": js.Do([]string{"clickEvent"},
        js.Const("btnText", "clickEvent.target.textContent"),
        js.Alert("'點擊了: ' + btnText"),
    ),
}, "點我")

// 表單事件使用語義化的名稱
Form(Props{
    "onSubmit": js.Do([]string{"submitEvent"},
        js.CallMethod("submitEvent", "preventDefault"),
        js.Const("formData", "new FormData(submitEvent.target)"),
        // 處理表單...
    ),
})

// 輸入事件
Input(Props{
    "onInput": js.Do([]string{"inputEvent"},
        js.Const("value", "inputEvent.target.value"),
        // 處理輸入...
    ),
})
```

### 技術細節

修改的文件：`jsdsl/jsdsl.go`

- 移除了參數名檢查邏輯
- 簡化為：只要 `params` 不為 `nil` 且長度大於 0，就注入 event
- 適用於 `Do()` 和 `AsyncDo()` 兩個函數

---

## 文檔整理

### 刪除的文檔

為了保持文檔簡潔，刪除了以下版本相關和重複的文檔：

- `docs/CHANGES_TRY_CATCH.md`
- `docs/EVENT_HANDLER_CHANGES.md`
- `docs/EVENT_PARAM_ANY_NAME.md`
- `docs/EVENT_PARAMETER_FIX.md`
- `docs/FIXES_2024_EVENT_HANDLERS.md`
- `docs/V1.2.0_SUMMARY.md`
- `docs/V1.2.1_API_CONSISTENCY.md`
- `docs/V1.2.1_SUMMARY.md`
- `docs/DO_ASYNCDO_PARAMS.md`
- `docs/EVENT_HANDLER_QUICK_REF.md`
- `docs/TRY_CATCH_QUICK_REF.md`
- `docs/TRY_CATCH_FINALLY.md`
- `docs/DESIGN_RATIONALE.md`
- `docs/OPTIMIZATION.md`

### 保留的核心文檔

現在只保留以下核心文檔：

- `docs/README.md` - 文檔導航
- `docs/QUICK_START.md` - 快速入門
- `docs/DOCUMENTATION.md` - 完整技術文檔
- `docs/API_REFERENCE.md` - API 參考
- `docs/QUICK_REFERENCE.md` - 快速參考

### 移除版本號

- 從所有文檔中移除版本號標記
- 保持文檔內容為最新狀態，不再強調具體版本

---

## 測試

新增測試文件 `test_param_injection.go`，用於驗證 Do/AsyncDo 可以使用任意參數名注入 event 對象。

運行測試：
```bash
go run test_param_injection.go
```

訪問 http://localhost:8089 查看測試結果。

---

## 兼容性

這些更改完全向後兼容，不會破壞現有代碼：

- 使用 `event`、`e`、`evt`、`ev` 作為參數名的代碼仍然正常工作
- 只是擴展了支持，允許使用任意參數名
- 所有現有示例和測試都能正常運行

---

## 建議

建議開發者使用語義化的參數名，讓代碼更易讀：

```go
// 點擊事件
"onClick": js.Do([]string{"clickEvent"}, ...)

// 表單提交
"onSubmit": js.Do([]string{"submitEvent"}, ...)

// 輸入變化
"onInput": js.Do([]string{"inputEvent"}, ...)

// 滑鼠事件
"onMouseOver": js.Do([]string{"mouseEvent"}, ...)
```

或者繼續使用簡短的標準名稱：

```go
// 簡短但標準
"onClick": js.Do([]string{"event"}, ...)
"onClick": js.Do([]string{"e"}, ...)
```
