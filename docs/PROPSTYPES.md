# PropsDef 類型支援說明

## 概述

從本版本開始，`PropsDef` 支援多種 Go 原生類型，不再僅限於字串。這使得組件的屬性定義更加靈活和類型安全。

## 類型定義

```go
type PropsDef map[string]interface{}
type Props map[string]interface{}
```

**重要：Props 和 PropsDef 使用相同的類型定義**

- ✅ 兩者都是 `map[string]interface{}`
- ✅ 兩者都支援所有 Go 原生類型
- ✅ 類型處理規則完全一致
- ✅ Props 可以覆寫 PropsDef 的預設值，並保持傳入時的類型

## 支援的類型

### 1. 字串 (string)

最常用的類型，支援模板插值。

```go
PropsDef{
    "title":       "預設標題",
    "placeholder": "請輸入內容",
    "color":       "#3b82f6",
}
```

### 2. 布林值 (bool)

用於控制組件的開關狀態。

```go
PropsDef{
    "disabled":  false,
    "checked":   true,
    "visible":   true,
    "required":  false,
}
```

**優點**：

- 類型安全，避免 `"true"` 和 `"false"` 字串混淆
- 可以直接用於 Go 層面的邏輯判斷
- 在模板插值時自動轉換為字串 `"true"` 或 `"false"`

### 3. 整數 (int, int8, int16, int32, int64)

用於數值計算和配置。

```go
PropsDef{
    "count":      10,
    "maxLength":  100,
    "age":        25,
    "level":      5,
}
```

### 4. 無符號整數 (uint, uint8, uint16, uint32, uint64)

```go
PropsDef{
    "id":     uint(12345),
    "port":   uint16(8080),
}
```

### 5. 浮點數 (float32, float64)

用於精確數值，如價格、百分比等。

```go
PropsDef{
    "price":      99.99,
    "discount":   0.15,
    "opacity":    0.8,
}
```

### 6. 切片 (slice)

用於列表數據。

```go
PropsDef{
    "tags":    []string{"Go", "VDOM", "Web"},
    "options": []int{1, 2, 3, 4, 5},
    "items":   []interface{}{"a", 1, true},
}
```

### 7. Map

用於複雜的配置對象。

```go
PropsDef{
    "config": map[string]string{
        "theme": "dark",
        "lang":  "zh-TW",
    },
    "metadata": map[string]interface{}{
        "version": "1.0",
        "count":   42,
        "active":  true,
    },
}
```

## 類型處理規則

### 核心原則：插值不改變類型

**重要：** 模板插值（`{{key}}`）**不會**改變 Props 中的值類型。類型轉換只發生在渲染成 HTML 字串時。

```go
// PropsDef 定義
PropsDef{
    "disabled": false,  // bool 類型
    "count": 10,        // int 類型
}

// 組件模板使用插值
Component(
    Input(Props{
        "disabled": "{{disabled}}",  // 模板引用
        "count": "{{count}}",        // 模板引用
    }),
    nil,
    PropsDef{...},
)

// 結果
result := MyComponent(Props{})
// result.Props["disabled"] 仍然是 bool 類型！
// result.Props["count"] 仍然是 int 類型！

// 只有在渲染成 HTML 時才轉換
html := Render(result)
// HTML: <input disabled count="10">
```

### Props 和 PropsDef 的一致性

Props 和 PropsDef 的類型處理規則完全相同：

```go
// 在 PropsDef 中定義預設值
PropsDef{
    "enabled": true,    // 布林值
    "count": 10,        // 整數
}

// 在 Props 中覆寫，類型保持一致
Props{
    "enabled": false,   // 布林值
    "count": 20,        // 整數
}

// 也可以使用不同類型覆寫
Props{
    "count": "20",      // 改為字串
}
```

### 規則 1：純模板引用保持原始類型

當屬性值在組件模板中是純粹的 `{{key}}`（沒有其他文字），則會保持原始類型。

```go
// 組件定義
MyComponent := Component(
    Div(Props{
        "disabled": "{{disabled}}", // 純模板引用
    }),
    nil,
    PropsDef{
        "disabled": false, // bool 類型
    },
)

// 結果
result := MyComponent(Props{})
// result.Props["disabled"] 仍是 bool: false ✓

// 渲染成 HTML 時才轉換
html := Render(result)
// HTML: <div></div>  （false 時不輸出）
```

### 規則 1.1：混合字串會轉換為字串

如果模板包含其他文字，則會轉換為字串：

```go
Component(
    Div(Props{
        "title": "Count: {{count}}", // 混合字串
    }),
    nil,
    PropsDef{
        "count": 10, // int 類型
    },
)

// 結果
result := MyComponent(Props{})
// result.Props["title"] 是字串: "Count: 10"
```

### 規則 2：未出現在模板中的屬性保持原始類型

如果屬性沒有在組件模板的 Props 中出現，它會保持 PropsDef 中定義的原始類型。

```go
// 組件定義
MyComponent := Component(
    Div(Props{"class": "container"}), // 沒有使用 disabled
    nil,
    PropsDef{
        "disabled": false, // 布林值
        "count":    10,    // 整數
    },
)

// 結果
result := MyComponent(Props{})
// result.Props["disabled"] 保持為 bool: false
// result.Props["count"] 保持為 int: 10
```

### 規則 3：渲染時的類型轉換

Props 中的值在渲染成 HTML 時才會轉換為字串：

```go
result := MyComponent(Props{
    "disabled": true,   // bool
    "count": 42,        // int
    "price": 99.99,     // float64
})

// Props 中保持原始類型
fmt.Printf("%T\n", result.Props["disabled"]) // bool
fmt.Printf("%T\n", result.Props["count"])    // int
fmt.Printf("%T\n", result.Props["price"])    // float64

// 渲染成 HTML 時轉換
html := Render(result)
// <div disabled count="42" price="99.99">
```

**HTML 布林屬性特殊處理：**

- `true` → 只輸出屬性名（如 `disabled`）
- `false` → 不輸出屬性

### 規則 4：用戶傳入的值會覆蓋預設值

Props 中傳入的值會直接覆寫 PropsDef 的預設值，並保持傳入時的類型：

```go
MyComponent := Component(
    Div(nil),
    nil,
    PropsDef{
        "count": 10, // int
    },
)

// 傳入 int（類型相同）
result1 := MyComponent(Props{"count": 20})
// result1.Props["count"] = 20 (int)

// 傳入 float64（類型不同）
result2 := MyComponent(Props{"count": 20.5})
// result2.Props["count"] = 20.5 (float64)

// 傳入字串（類型不同）
result3 := MyComponent(Props{"count": "30"})
// result3.Props["count"] = "30" (string)
```

**提示**：Props 和 PropsDef 的類型處理完全一致，沒有任何區別。

## 最佳實踐

### ✅ 推薦做法

1. **使用正確的類型**

```go
PropsDef{
    "disabled": false,    // ✓ 使用布林值
    "count":    10,       // ✓ 使用整數
    "price":    99.99,    // ✓ 使用浮點數
}
```

2. **在 Go 層面進行類型判斷和計算**

```go
result := MyComponent(Props{...})

// 類型斷言
if disabled, ok := result.Props["disabled"].(bool); ok && disabled {
    // 處理禁用狀態
}

// 數值計算
price := result.Props["price"].(float64)
quantity := result.Props["quantity"].(int)
total := price * float64(quantity)
```

3. **為布林值提供明確的預設值**

```go
PropsDef{
    "visible":  true,   // 預設顯示
    "disabled": false,  // 預設啟用
    "loading":  false,  // 預設非載入狀態
}
```

### ❌ 避免的做法

1. **不要使用字串表示布林值**

```go
PropsDef{
    "disabled": "false", // ✗ 避免使用字串
}

// 應該使用
PropsDef{
    "disabled": false,   // ✓ 使用布林值
}
```

2. **不要使用字串表示數字**

```go
PropsDef{
    "count": "10",       // ✗ 避免使用字串
    "price": "99.99",    // ✗ 避免使用字串
}

// 應該使用
PropsDef{
    "count": 10,         // ✓ 使用整數
    "price": 99.99,      // ✓ 使用浮點數
}
```

## 實際示例

### 示例 1：配置面板組件

```go
ConfigPanel := Component(
    Div(Props{"class": "config-panel"},
        H3("設定"),
        P("最大連線數: {{maxConnections}}"),
        P("超時時間: {{timeout}} 秒"),
    ),
    nil,
    PropsDef{
        "debugMode":      false,  // 布林值
        "maxConnections": 100,    // 整數
        "timeout":        30,     // 整數
        "retryCount":     3,      // 整數
    },
)

// 使用預設值
panel1 := ConfigPanel(Props{})

// 覆寫部分值
panel2 := ConfigPanel(Props{
    "debugMode":      true,
    "maxConnections": 200,
})

// 在 Go 層面進行邏輯判斷
if debugMode := panel2.Props["debugMode"].(bool); debugMode {
    fmt.Println("調試模式已啟用")
}
```

### 示例 2：價格計算組件

```go
PriceCard := Component(
    Div(Props{"class": "price-card"},
        H4("{{productName}}"),
        P("單價: ${{price}}"),
        P("數量: {{quantity}}"),
    ),
    nil,
    PropsDef{
        "productName": "商品",
        "price":       99.99,  // float64
        "quantity":    1,      // int
        "taxRate":     0.05,   // float64 (5%)
    },
)

result := PriceCard(Props{
    "productName": "筆記本電腦",
    "price":       1299.99,
    "quantity":    2,
})

// 在 Go 層面計算總價
price := result.Props["price"].(float64)
quantity := result.Props["quantity"].(int)
taxRate := result.Props["taxRate"].(float64)

subtotal := price * float64(quantity)
tax := subtotal * taxRate
total := subtotal + tax

fmt.Printf("小計: $%.2f\n", subtotal)
fmt.Printf("稅金: $%.2f\n", tax)
fmt.Printf("總計: $%.2f\n", total)
```

### 示例 3：表單組件

```go
FormField := Component(
    Div(Props{"class": "form-field"},
        Label(Props{"for": "{{id}}"}, "{{label}}"),
        Input(Props{
            "id":       "{{id}}",
            "type":     "{{type}}",
            "required": "{{required}}",
            "disabled": "{{disabled}}",
        }),
    ),
    nil,
    PropsDef{
        "id":       "",
        "label":    "",
        "type":     "text",
        "required": false,    // 布林值
        "disabled": false,    // 布林值
        "minLength": 0,       // 整數
        "maxLength": 255,     // 整數
    },
)

// 使用
emailField := FormField(Props{
    "id":       "email",
    "label":    "電子郵件",
    "type":     "email",
    "required": true,
})
```

## Props 使用示例

由於 Props 和 PropsDef 類型一致，你可以在任何地方使用相同的方式定義屬性：

```go
// 直接使用 Props 創建元素
btn := Button(Props{
    "disabled": false,    // 布林值
    "tabindex": 0,        // 整數
    "type": "submit",     // 字串
})

// 使用 PropsDef 定義組件預設值
MyButton := Component(
    Button(Props{"type": "{{type}}"}),
    nil,
    PropsDef{
        "type": "button",     // 字串
        "disabled": false,    // 布林值
        "tabindex": 0,        // 整數
    },
)

// 使用 Props 覆寫
customBtn := MyButton(Props{
    "type": "submit",
    "disabled": true,
})
```

## 遷移指南

如果你的專案之前使用字串表示布林值或數字，可以按照以下步驟遷移：

### 步驟 1：識別需要修改的 PropsDef

尋找使用字串表示布林值或數字的地方：

```go
// 舊寫法
PropsDef{
    "disabled": "false",
    "visible":  "true",
    "count":    "10",
}
```

### 步驟 2：修改為正確的類型

```go
// 新寫法
PropsDef{
    "disabled": false,
    "visible":  true,
    "count":    10,
}
```

### 步驟 3：更新邏輯判斷

如果程式碼中有字串比較邏輯，需要修改：

```go
// 舊寫法
if props["disabled"] == "true" {
    // ...
}

// 新寫法
if disabled, ok := props["disabled"].(bool); ok && disabled {
    // ...
}
```

## 常見問題

### Q1: Props 和 PropsDef 的類型處理有什麼區別？

**A**: 完全沒有區別！兩者使用相同的類型定義 `map[string]interface{}`，類型處理規則也完全一致。PropsDef 用於定義組件的預設值，Props 用於傳入實際值，但它們的類型系統是統一的。

### Q2: 模板中的 `{{disabled}}` 還是布林值嗎？

**A**: 是的！純模板引用（如 `"{{disabled}}"`）會保持原始類型。Props 中的值不會因為插值而改變類型。只有在渲染成 HTML 字串時，才會轉換為字串。

```go
// Props 中保持布林
result.Props["disabled"] // bool: true

// 渲染時轉換
Render(result) // HTML: <input disabled>
```

**例外：** 如果是混合字串（如 `"Status: {{disabled}}"`），則會轉換為字串。

### Q3: 如何在組件內部訪問原始類型？

**A**: 對於沒有出現在模板 Props 中的屬性，直接從 Props 中讀取即可：

```go
result := MyComponent(Props{...})
count := result.Props["count"].(int)
```

### Q4: 可以使用自訂結構體嗎？

**A**: 可以！PropsDef 支援 `interface{}`，所以你可以存儲任何類型：

```go
type Config struct {
    Host string
    Port int
}

PropsDef{
    "config": Config{Host: "localhost", Port: 8080},
}
```

### Q5: 切片和 Map 在模板中如何使用？

**A**: 切片和 Map 不能直接在模板插值中展開，但可以在 Go 層面使用控制流組件（如 `For`、`Map`）來處理它們。

### Q6: Props 可以傳入與 PropsDef 不同的類型嗎？

**A**: 可以！Props 會直接覆寫 PropsDef 的值，並保持傳入時的類型。例如，PropsDef 定義 `"count": 10` (int)，但你可以在 Props 中傳入 `"count": "10"` (string) 或 `"count": 10.5` (float64)。

### Q7: 什麼時候會轉換為字串？

**A**: 只有在以下情況才會轉換：

1. **渲染成 HTML 時**：所有屬性值都會轉為字串（HTML 的要求）
2. **混合字串模板**：如 `"Count: {{count}}"` 會在插值時轉換

**不會轉換的情況：**

- 純模板引用：`"{{key}}"` 保持原始類型
- 不在模板中的屬性：保持原始類型

## 總結

- ✅ **Props 和 PropsDef 使用相同的類型系統**
- ✅ 兩者都支援所有 Go 原生類型
- ✅ 類型處理規則完全一致
- ✅ 布林值應使用 `true`/`false` 而非字串
- ✅ 數字應使用數值類型而非字串
- ✅ **純模板引用 `{{key}}` 保持原始類型**
- ✅ **只有渲染成 HTML 時才轉換為字串**
- ✅ 混合字串模板會在插值時轉換
- ✅ 類型安全，減少錯誤
- ✅ 更好的程式碼可讀性和可維護性

---

**文檔版本**: 1.0.0
**最後更新**: 2025-01-24
**作者**: TimLai666
