# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2024

### Added
- **Try-Catch-Finally 流暢 API**: 全新設計，生成純粹的 try-catch-finally 語句
  - `js.Try(...).Catch(...).End()` - Try-Catch 模式
  - `js.Try(...).Catch(...).Finally(...)` - 完整錯誤處理
  - `js.Try(...).Finally(...)` - Try-Finally 模式
  - Try 不再自動包裝在 IIFE 中，更加靈活
  - 錯誤對象統一命名為 `error`（而非 `e`）
- **Do / AsyncDo**: 專門用於創建立即執行函數（IIFE）
  - `js.Do(...)` - 創建立即執行的普通函數
  - `js.AsyncDo(...)` - 創建立即執行的異步函數
  - 職責分離：Try 負責錯誤處理，Do/AsyncDo 負責 IIFE
- 新示例：`examples/07_trycatch_usage.go` - 展示所有 Try-Catch-Finally 和 Do/AsyncDo 用法
- 新文檔：`docs/TRY_CATCH_FINALLY.md` - 完整使用指南
- 新文檔：`docs/TRY_CATCH_QUICK_REF.md` - 快速參考手冊
- 新文檔：`docs/CHANGES_TRY_CATCH.md` - API 重新設計說明

### Changed
- Try-Catch-Finally 不再自動包裝在 async IIFE 中
- 需要 async/await 時，使用 AsyncFn 或 AsyncDo 包裝
- 更新 `examples/03_javascript_dsl.go` 使用新的 Try API
- 更新 `examples/05_foreach_usage.go` 使用新的 Try API
- 更新 README 加入 Try-Catch-Finally 和 Do/AsyncDo 說明

### Deprecated
- `TryCatch` 函數仍可用但建議使用新的流暢 API（Try + AsyncFn/AsyncDo）

## [1.1.0] - 2025-01-24

### Added
- **⚡ AsyncFn 異步函數支持** (`jsdsl/jsdsl.go`)
  - ✅ 新增 `AsyncFn()` 函數，用於創建支持 `await` 語法的異步函數
  - ✅ 解決 "await is only valid in async functions" 錯誤
  - ✅ 完全兼容現有的 `Fn()` API
  - ✅ 支持參數傳遞和函數體定義
  - 📝 任何包含 `await` 語句的函數都應使用 `AsyncFn` 而非 `Fn`
  - 📝 與重新設計的 `TryCatch` 完美配合處理異步錯誤

- **🔄 TryCatch 重新設計** (`jsdsl/jsdsl.go`)
  - ✅ 重新設計 `TryCatch(tryActions, catchActions, finallyActions)` - 接受動作列表而非函數包裝
  - ✅ 自動創建 async 函數包裝，完全支持 await 語法
  - ✅ 解決了之前 TryCatch 內部無法使用 AsyncFn 的問題
  - ✅ 更符合直覺的 API 設計
  - ✅ 立即執行，錯誤對象自動命名為 `e`
  - 📝 新用法：
    ```go
    js.TryCatch(
        []JSAction{
            js.Const("data", "await fetch('/api')"),
            js.Log("data"),
        },
        []JSAction{
            js.Log("'錯誤:', e.message"),
        },
        nil,
    )
    ```
  - 📝 舊的包裝式 API 已廢棄，請使用新的列表式 API

- **🔄 ForEach 列表渲染改進**
  - **後端渲染** (`vdom/tags.go`)
    - ✅ 新增 `ForEach[T](items []T, func(item T) VNode) []VNode` - 簡潔的列表渲染
    - ✅ 新增 `ForEachWithIndex[T](items []T, func(item T, i int) VNode) []VNode` - 帶索引的列表渲染
    - ✅ 提供更簡潔的語法，無需導入 control 包
    - ✅ 支持泛型，可用於任何類型的切片
    - 📝 用法：`Ul(ForEach(items, func(item string) VNode { return Li(item) }))`
  - **前端渲染** (`jsdsl/jsdsl.go`)
    - ✅ 重構 `ForEachJS(arrayExpr, itemVar, actions...)` - 更通用的前端遍歷
    - ✅ 新增 `ForEachWithIndexJS(arrayExpr, itemVar, indexVar, actions...)` - 帶索引的前端遍歷
    - ✅ 新增 `ForEachElement(arrayExpr, func(el Elem) JSAction)` - DOM 元素專用遍歷
    - ✅ 不再限於 DOM 元素，可遍歷任何 JavaScript 數組
    - 📝 用法：`js.ForEachJS("['A','B','C']", "item", js.Log("item"))`
  - **新示例程序** (`examples/05_foreach_usage.go`)
    - 展示後端和前端 ForEach 的完整用法
    - 包含 9 個實用示例
    - 提供最佳實踐和對比表格

- **➰ control.For 和 control.ForEach 重構** (`control/control.go`)
  - ✅ 將 `control.For` 改名為 `control.ForEach` - 更語義化
  - ✅ 新增傳統循環 `control.For(start, end, step, func(i) VNode)` - 類似 for i := start; i < end; i += step
  - ✅ 支持正向循環：`For(1, 11, 1, ...)` 生成 1-10
  - ✅ 支持倒序循環：`For(10, 0, -1, ...)` 生成 10-1
  - ✅ 支持步進循環：`For(0, 20, 2, ...)` 生成偶數 0-18
  - ✅ 將 `KeyedFor` 改名為 `KeyedForEach` 保持一致性
  - 📝 用法：
    - 遍歷集合：`control.ForEach(items, func(item, i) VNode {...})`
    - 數字循環：`control.For(1, 11, 1, func(i) VNode {...})`
  - **新示例程序** (`examples/06_control_loops.go`)
    - 展示 ForEach 和 For 的完整用法
    - 包含 8 個實用示例（分頁、表格、評分系統等）
    - 提供詳細的對比表格和選擇指南

### Changed
- **🔧 TryCatch API 重大改進**: 從函數包裝改為動作列表
  - ⚠️ **破壞性變更**：舊的 `TryCatch(baseAction, catchFn, finallyFn)` 已移除
  - ✅ 新的 `TryCatch(tryActions, catchActions, finallyActions)` 更直觀
  - ✅ 不再需要 `js.Ptr()` 包裝
  - ✅ 內部可以直接使用任何 JSAction，包括包含 await 的語句
  - ✅ 自動處理異步邏輯
  - 📝 遷移指南：將原本的 `js.AsyncFn(nil, ...actions)` 改為 `[]JSAction{...actions}`

- **🎯 Props 類型系統重大改進**: 從 `map[string]string` 改為 `map[string]interface{}`
  - ✅ 支持任意類型的值（string, bool, int, float64, JSAction 等）
  - ✅ 自動根據類型轉換，無需手動轉換
  - ✅ 布爾值語義更明確（`true` 渲染屬性，`false` 省略屬性）
  - ✅ 完全向後兼容，現有代碼無需修改
  - 📝 詳見 [IMPROVEMENTS.md](IMPROVEMENTS.md)

- **重構 main.go**: 將所有 JavaScript 代碼重構為使用 DSL
  - GET 請求完全使用 `js.TryCatch` 和其他 DSL 函數
  - POST 請求完全使用 DSL 而非原始 JavaScript 字符串
  - 提高了代碼的可讀性和類型安全性
  - 更容易維護和調試

### Added
- **重構 main.go**: 所有異步 JavaScript 代碼改用 `AsyncFn`
  - GET 請求使用 `AsyncFn` 包裝 `await fetch()`
  - POST 請求使用 `AsyncFn` 包裝表單提交邏輯
  - 解決了控制台 "await is only valid in async functions" 錯誤
  - 提高了異步代碼的正確性和可維護性

- **🔄 模板序列化功能** (`vdom/template.go`)
  - `ToGoTemplate()` - 將 VNode 轉換為 Go template 格式
  - `SaveTemplate()` - 保存為命名模板
  - `ToJSON()` / `FromJSON()` - JSON 序列化和反序列化
  - `ExecuteGoTemplate()` - 執行 Go template
  - `ExtractTemplateVars()` - 提取模板中的所有變數
  - `CloneVNode()` - 深度克隆 VNode
  - `MergeProps()` - 智能合併多個 Props
  - `WrapWithLayout()` - 將內容包裝到佈局模板中
  - 📝 支持與 Go `html/template` 無縫集成
  - 📝 可以將模板保存到文件並重用

- **文檔重組**
  - 將所有文檔移至 `docs/` 目錄
  - 新增 `docs/API_REFERENCE.md` - JavaScript DSL 完整 API 參考（包含 AsyncFn）
  - 更新 `docs/QUICK_START.md` - 快速入門指南
  - 保留 `docs/DOCUMENTATION.md` - 完整技術文檔
  - 保留 `docs/QUICK_REFERENCE.md` - 語法速查表
  - 刪除冗餘文檔（IMPROVEMENTS.md, UPDATE_SUMMARY.md）
  - 大幅精簡 README.md，保留核心內容並引用 docs/ 文檔

### Added (from v1.0.0)
- **完整文檔** (`DOCUMENTATION.md`)
  - 詳細的架構設計說明
  - 完整的 API 參考
  - 進階用法指南
  - 性能優化建議
  - 故障排除章節
  - 超過 1700 行的綜合文檔

- **增強的 README**
  - 添加詳細的目錄結構
  - 完整的核心概念解釋
  - 詳細的 HTML 元素、Props、組件定義說明
  - 條件渲染和列表渲染的完整示例
  - JavaScript 事件處理指南
  - Fetch API 集成示例
  - UI 組件庫詳細說明
  - 完整的 HTTP 服務器示例
  - 完整的 API 參考
  - 最佳實踐指南
  - 常見問題解答

- **示例程序集合** (`examples/` 目錄)
  - `01_basic_usage.go` - 基本用法示例
    - Document 函數使用
    - 基本 HTML 元素
    - Bootstrap 集成
    - 頁面佈局
  
  - `02_components.go` - 組件系統示例
    - Alert 組件
    - Card 組件
    - Badge 組件
    - Button 組件
    - UserCard 組件
    - 組件組合和嵌套
  
  - `03_javascript_dsl.go` - JavaScript DSL 示例
    - DOM 操作（SetText, SetHTML, AddClass, RemoveClass）
    - 變數定義（Let, Const）
    - 事件處理（OnClick）
    - 表單處理
    - 動態創建元素
    - Try/Catch 錯誤處理
    - DomReady 初始化
  
  - `README.md` - 示例文檔
    - 每個示例的詳細說明
    - 運行指南
    - 學習路徑建議

### Improved
- **代碼質量**
  - 更好的 DSL 使用示範
  - 更清晰的代碼結構
  - 更詳細的註釋

- **文檔質量**
  - 從簡略文檔擴展到超過 2000 行的完整文檔
  - 添加了大量實用示例
  - 提供了詳細的 API 參考
  - 包含了最佳實踐和故障排除指南

- **測試文件** (`vdom/template_test.go`)
  - 完整的單元測試覆蓋
  - 性能基準測試
  - 測試所有序列化功能

- **新示例程序** (`examples/04_template_serialization.go`)
  - 展示 Go Template 導出和導入
  - 展示 JSON 序列化和反序列化
  - 展示模板變數提取
  - 展示 VNode 克隆和 Props 合併
  - HTTP 服務器示範

### Documentation
- **新增 IMPROVEMENTS.md** - Props 類型系統和模板序列化的詳細說明
  - Props 類型系統改進
  - 模板序列化功能
  - 遷移指南
  - 使用示例
  - 性能影響分析
  - 常見問題解答

- 所有文檔均使用繁體中文
- 添加了完整的代碼示例
- 提供了從基礎到進階的學習路徑
- 包含了實際可運行的示例程序
- 更新 README.md 添加 Props 類型系統和模板序列化章節</parameter>
- 包含了實際可運行的示例程序

## [1.0.0] - 2025-01-24

### Added
- 初始版本發布
- 虛擬 DOM 核心實現
- 組件系統
- 控制流（If/Then/Else, Repeat, For）
- JavaScript DSL
- UI 組件庫
- Bootstrap 集成
- Fetch API 支持

### Features
- 類型安全的 HTML 生成
- 聲明式 API
- 服務器端渲染
- 可重用組件
- JavaScript 代碼生成
- 表單組件

---

## 貢獻指南

在提交更改時，請更新此 CHANGELOG：

1. 在 `[Unreleased]` 部分添加你的更改
2. 使用以下類別之一：
   - `Added` - 新功能
   - `Changed` - 現有功能的更改
   - `Deprecated` - 即將移除的功能
   - `Removed` - 已移除的功能
   - `Fixed` - Bug 修復
   - `Security` - 安全性相關更改

3. 簡要描述更改內容
4. 如果適用，添加相關的 Issue 或 PR 編號

---

**注意**: 此項目遵循 [語義化版本控制](https://semver.org/lang/zh-TW/)。

- **主版本號（MAJOR）**: 不兼容的 API 更改
- **次版本號（MINOR）**: 向後兼容的新功能
- **修訂號（PATCH）**: 向後兼容的 Bug 修復