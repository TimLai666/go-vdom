# 複雜類型 JSON 序列化功能更新

## 版本資訊
**功能版本**: v1.2.0
**更新日期**: 2024

## 📋 更新摘要

為 `Component` 函數添加了自動 JSON 序列化功能，現在可以輕鬆地將陣列（Array/Slice）、Map、結構體（Struct）等複雜類型作為 props 傳遞，系統會自動將它們序列化為 JSON 字符串。

## ✨ 新增功能

### 1. 自動 JSON 序列化

當你在 `Component` 中使用複雜類型的 props 時，系統會自動將它們序列化為 JSON 字符串：

```go
// 使用陣列作為 prop
componentFn(dom.Props{
    "data-items": []string{"apple", "banana", "orange"},
})
// → 渲染為: data-items='["apple","banana","orange"]'

// 使用 Map 作為 prop
componentFn(dom.Props{
    "data-config": map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
    },
})
// → 渲染為: data-config='{"fontSize":14,"theme":"dark"}'

// 使用結構體作為 prop
componentFn(dom.Props{
    "data-user": User{Name: "John", Email: "john@example.com"},
})
// → 渲染為: data-user='{"Name":"John","Email":"john@example.com"}'
```

### 2. 支持的複雜類型

- ✅ **陣列/切片** (`[]T`)
- ✅ **Map** (`map[K]V`)
- ✅ **結構體** (`struct`)
- ✅ **指針** (`*T` - 會自動解引用)
- ✅ **嵌套結構** (任意複雜度)

### 3. 簡單類型保持字符串轉換

簡單類型（如 `int`、`bool`、`float`、`string`）會被轉換為字符串：

```go
Props{
    "count": 42,        // → "42"
    "active": true,     // → "true"
    "price": 19.99,     // → "19.99"
    "name": "John",     // → "John"
}
```

## 🔧 技術實現

### 核心函數

新增 `serializeComplexType` 函數 (`dom/component.go`)：

```go
func serializeComplexType(v interface{}) string
```

此函數會：
1. 檢查值的類型（使用 `reflect` 包）
2. 對於簡單類型，返回字符串表示
3. 對於複雜類型（Slice、Array、Map、Struct），使用 `json.Marshal` 序列化為 JSON
4. 處理指針類型（遞歸處理指向的值）

### 修改的函數

**`interpolate` 函數**：
- 當模板引用為純 `{{key}}` 形式時，會調用 `serializeComplexType` 將值序列化
- 保持與 HTML 屬性的一致性（所有屬性值都是字符串）

**`interpolateString` 函數**：
- 在字符串插值時，對複雜類型調用 `serializeComplexType`

## 📝 使用示例

### 示例 1: 傳遞陣列數據

```go
template := dom.VNode{
    Tag: "div",
    Props: dom.Props{
        "data-items": "{{items}}",
    },
}

componentFn := dom.Component(template, nil)
result := componentFn(dom.Props{
    "items": []string{"Apple", "Banana", "Orange"},
})

// 在客戶端 JavaScript 中使用：
// const items = JSON.parse(element.dataset.items);
// console.log(items); // ["Apple", "Banana", "Orange"]
```

### 示例 2: 傳遞配置對象

```go
type Config struct {
    Theme    string `json:"theme"`
    FontSize int    `json:"fontSize"`
}

componentFn(dom.Props{
    "data-config": Config{
        Theme:    "dark",
        FontSize: 16,
    },
})

// 客戶端：
// const config = JSON.parse(element.dataset.config);
// console.log(config.theme); // "dark"
```

### 示例 3: 嵌套複雜結構

```go
componentFn(dom.Props{
    "data-app-state": map[string]interface{}{
        "users": []User{
            {Name: "Alice", Email: "alice@example.com"},
            {Name: "Bob", Email: "bob@example.com"},
        },
        "settings": map[string]bool{
            "darkMode": true,
            "notifications": false,
        },
    },
})
```

## 🧪 測試

新增完整的測試套件 (`dom/component_test.go`)：

- `TestSerializeComplexType`: 測試各種類型的序列化
  - 13 個測試案例，涵蓋所有類型
- `TestComponentWithComplexProps`: 測試組件中的複雜 props
  - 6 個測試案例，包括實際使用場景
- `TestComponentJSONInTemplate`: 測試模板中的 JSON 使用
- `TestComponentWithPointerProps`: 測試指針類型處理
- `TestInterpolateWithComplexTypes`: 測試字符串插值

**測試結果**: ✅ 所有測試通過

## 📚 文檔更新

### 更新的文檔

1. **`docs/DOCUMENTATION.md`**
   - 在 "Props 屬性系統" 章節新增複雜類型說明
   - 添加 JSON 序列化使用示例
   - 包含客戶端 JavaScript 使用範例

2. **`README.md`**
   - 更新 "Props 類型系統" 章節
   - 列出支持的複雜類型
   - 添加新示例的運行說明

3. **新增示例**: `examples/complex_props_demo.go`
   - 完整的可運行示例
   - 展示陣列、Map、結構體的使用
   - 包含交互式列表和配置面板
   - 運行於 `http://localhost:8084`

## 🎯 使用場景

此功能特別適合以下場景：

1. **數據驅動的 UI 組件**
   - 將數據列表傳遞給組件進行渲染

2. **配置對象**
   - 傳遞複雜的配置選項

3. **狀態管理**
   - 在 HTML 屬性中存儲應用狀態

4. **服務器端渲染**
   - 將服務器端數據嵌入到 HTML 中供客戶端使用

5. **漸進式增強**
   - 在 HTML 中嵌入數據，JavaScript 可選加載

## ⚠️ 注意事項

1. **JSON 大小**: 對於大型數據結構，JSON 字符串可能很長，建議考慮性能影響
2. **循環引用**: 結構體不應包含循環引用，否則 `json.Marshal` 會失敗
3. **導出字段**: 結構體字段必須是導出的（首字母大寫）才能被 JSON 序列化
4. **HTML 轉義**: JSON 字符串會被正確轉義以在 HTML 屬性中使用

## 🔄 向後兼容性

此更新完全向後兼容：

- ✅ 現有代碼無需修改
- ✅ 簡單類型的行為保持不變
- ✅ 所有現有測試通過
- ✅ 不影響已有的組件實現

## 📊 性能影響

- JSON 序列化使用 Go 標準庫的 `encoding/json`
- 僅在需要時進行序列化（按需處理）
- 對簡單類型沒有額外開銷
- 建議對大型數據結構進行基準測試

## 🚀 快速開始

```bash
# 運行示例
go run examples/complex_props_demo.go

# 訪問
open http://localhost:8084
```

## 📖 相關文檔

- [完整文檔](docs/DOCUMENTATION.md#props屬性系統)
- [快速入門](docs/QUICK_START.md)
- [API 參考](docs/API_REFERENCE.md)

## 🙏 貢獻

如有任何問題或建議，歡迎提交 Issue 或 Pull Request。
