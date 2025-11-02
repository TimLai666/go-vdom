# 插值改進：參數本身不轉換類型

## 用戶觀點

> "我覺得不應該有出現在插值的就轉字串，而是只有插值時轉字串，參數本身不轉"

## 問題說明

### 舊行為（已修正）

在之前的實作中，如果屬性在模板中使用了插值語法 `{{key}}`，那麼在 `interpolate` 函數中會立即將值轉換為字串：

```go
// 舊行為
template.Props["disabled"] = "{{disabled}}"  // 模板
mergedProps["disabled"] = false              // 布林值

// 經過 interpolate() 後
newProps["disabled"] = "false"  // ❌ 被轉換為字串！
```

**問題：**
- Props 中的值失去了類型資訊
- 無法在 Go 層面進行類型判斷
- 布林值、數字等都變成了字串

### 新行為（已改進）✅

現在的實作中，**插值不會改變 Props 中的值類型**。只有在渲染成 HTML 字串時才進行類型轉換：

```go
// 新行為
template.Props["disabled"] = "{{disabled}}"  // 模板
mergedProps["disabled"] = false              // 布林值

// 經過 interpolate() 後
newProps["disabled"] = false  // ✓ 保持布林值！

// 只有在渲染時才轉換
html := Render(result)
// HTML: <div></div>  (false 時不輸出)
```

**優點：**
- ✅ Props 中保持原始類型
- ✅ 可以進行類型判斷和邏輯運算
- ✅ 延遲轉換，只在必要時（渲染 HTML）才轉換
- ✅ 更符合預期的行為

## 實作細節

### 修改 1：`interpolate` 函數

```go
// dom/component.go

case string:
    // 檢查是否為純模板引用（如 "{{key}}"）
    trimmed := strings.TrimSpace(t)
    if strings.HasPrefix(trimmed, "{{") &&
       strings.HasSuffix(trimmed, "}}") &&
       strings.Count(trimmed, "{{") == 1 {
        // 純模板引用：直接取值，保持原始類型
        key := strings.TrimSpace(strings.TrimSuffix(
            strings.TrimPrefix(trimmed, "{{"), "}}"))
        if val, ok := p[key]; ok {
            newProps[k] = val // ✓ 保持原始類型
        } else {
            newProps[k] = ""
        }
    } else {
        // 混合字串或複雜模板：進行字串插值
        newProps[k] = interpolateString(t, p)
    }
```

**邏輯：**
1. 如果是純模板引用 `"{{key}}"`，直接從 Props 中取值，保持類型
2. 如果是混合字串 `"Count: {{count}}"`，則進行字串插值
3. 如果不在模板中，直接複製，保持類型

### 修改 2：`Render` 函數

```go
// dom/render.go

case bool:
    // 布林值：false 時不輸出屬性，true 時輸出屬性名
    isBool = true
    boolVal = t
    if !t {
        continue // false 時跳過
    }
    valStr = k // true 時只輸出屬性名
```

**HTML 布林屬性處理：**
- `disabled: true` → `<input disabled>`
- `disabled: false` → `<input>` （不輸出）
- `maxlength: 100` → `<input maxlength="100">`

## 行為對比

### 案例 1：純模板引用

```go
Component(
    Input(Props{
        "disabled": "{{disabled}}",  // 純模板引用
    }),
    nil,
    PropsDef{
        "disabled": false,  // bool
    },
)

result := MyComponent(Props{})
```

| 層面 | 舊行為 ❌ | 新行為 ✅ |
|------|-----------|-----------|
| Props 類型 | `string: "false"` | `bool: false` |
| 可類型判斷 | ❌ 需要字串比較 | ✅ 直接布林判斷 |
| HTML 輸出 | `disabled="false"` | （不輸出） |

### 案例 2：混合字串

```go
Component(
    Div(Props{
        "title": "Count: {{count}}",  // 混合字串
    }),
    nil,
    PropsDef{
        "count": 10,  // int
    },
)
```

| 層面 | 行為 |
|------|------|
| Props 類型 | `string: "Count: 10"` |
| 說明 | 混合字串需要插值，所以轉為字串 ✓ |

### 案例 3：不在模板中

```go
Component(
    Div(Props{"class": "container"}),
    nil,
    PropsDef{
        "enabled": true,  // bool
        "count": 10,      // int
    },
)
```

| 層面 | 行為 |
|------|------|
| Props 類型 | `bool: true`, `int: 10` |
| 說明 | 不在模板中，直接保持原始類型 ✓ |

## 使用示例

### 示例 1：類型判斷

```go
FormField := Component(
    Input(Props{
        "required": "{{required}}",
        "disabled": "{{disabled}}",
    }),
    nil,
    PropsDef{
        "required": true,
        "disabled": false,
    },
)

result := FormField(Props{})

// ✓ 可以直接進行布林判斷
if required, ok := result.Props["required"].(bool); ok && required {
    fmt.Println("此欄位為必填")
}

// ✓ 可以直接比較
if disabled := result.Props["disabled"].(bool); !disabled {
    fmt.Println("欄位已啟用")
}
```

### 示例 2：數值運算

```go
PriceCard := Component(
    Div(Props{
        "price": "{{price}}",
        "quantity": "{{quantity}}",
    }),
    nil,
    PropsDef{
        "price": 99.99,
        "quantity": 1,
    },
)

result := PriceCard(Props{"quantity": 3})

// ✓ 可以直接進行數值運算
price := result.Props["price"].(float64)
quantity := result.Props["quantity"].(int)
total := price * float64(quantity)
fmt.Printf("總價: $%.2f\n", total)
```

### 示例 3：HTML 渲染

```go
input := Input(Props{
    "type": "email",
    "required": true,    // bool
    "disabled": false,   // bool
    "maxlength": 100,    // int
})

html := Render(input)
// HTML: <input type="email" required maxlength="100">
//
// 注意：
// - required=true 渲染為 "required" (HTML 布林屬性)
// - disabled=false 不輸出
// - maxlength=100 渲染為 "maxlength=\"100\""
```

## 測試驗證

### 測試 1：純模板引用保持類型

```go
✓ disabled (布林): bool = false
✓ count (整數):    int = 42
✓ price (浮點):    float64 = 99.99
```

### 測試 2：混合字串轉換

```go
✓ title:       string = Title: 測試
✓ data-status: string = Status is true
✓ class:       string = btn-primary
```

### 測試 3：HTML 渲染正確

```go
✓ disabled=true 渲染為 HTML 布林屬性 (只有屬性名)
✓ disabled=false 不輸出到 HTML
✓ 數字類型渲染為帶值的屬性
```

## 優勢總結

### 1. 類型安全

```go
// ✓ 可以使用類型斷言
if disabled, ok := props["disabled"].(bool); ok {
    // 類型安全的處理
}

// ❌ 舊方式需要字串比較
if props["disabled"] == "true" {
    // 容易出錯
}
```

### 2. 邏輯清晰

```go
// ✓ 直觀的布林邏輯
if !disabled && required {
    // 清晰的邏輯運算
}

// ❌ 舊方式複雜
if props["disabled"] != "true" && props["required"] == "true" {
    // 不直觀
}
```

### 3. 數值運算

```go
// ✓ 直接運算
total := price * float64(quantity)

// ❌ 舊方式需要轉換
price, _ := strconv.ParseFloat(props["price"], 64)
quantity, _ := strconv.Atoi(props["quantity"])
total := price * float64(quantity)
```

### 4. HTML 渲染優化

```go
// ✓ HTML 布林屬性正確處理
<input disabled>           // true
<input>                    // false (不輸出)

// ❌ 舊方式錯誤
<input disabled="true">    // 錯誤的 HTML
<input disabled="false">   // 錯誤的 HTML
```

## 設計原則

### 延遲轉換原則

**只在必要時轉換類型**

1. **Props 層面**：保持原始類型，方便邏輯處理
2. **HTML 層面**：轉換為字串，符合 HTML 規範

### 最小驚訝原則

**行為符合預期**

```go
PropsDef{
    "disabled": false,  // 開發者定義為 bool
}

// 開發者期望：
result.Props["disabled"]  // 應該還是 bool
```

### 類型一致性原則

**PropsDef 和 Props 處理一致**

- 兩者都使用 `map[string]interface{}`
- 兩者都遵循相同的插值規則
- 兩者都在渲染時才轉換

## 總結

✅ **改進完成**

用戶的觀點完全正確，現在的實作已經做到：

1. ✅ **插值不改變參數類型**
   - 純模板引用 `{{key}}` 保持原始類型
   - Props 中的值保持類型資訊

2. ✅ **只在渲染時轉換**
   - Props 層面：保持原始類型
   - HTML 層面：轉換為字串

3. ✅ **類型安全**
   - 可以進行類型斷言
   - 可以進行邏輯運算
   - 可以進行數值計算

4. ✅ **HTML 正確性**
   - 布林屬性正確渲染
   - 數字屬性正確轉換
   - 符合 HTML 規範

**這是一個重要的設計改進，讓框架更加直觀、類型安全，並符合開發者的預期！** 🎉

---

**文檔版本**: 1.0.0
**最後更新**: 2025-01-24
**作者**: TimLai666
