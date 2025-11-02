# Props 和 PropsDef 類型處理一致性

## 問題回答

**問：props 跟 propsDef 的類型處理一致嗎？**

**答：是的，完全一致！** ✅

## 類型定義

兩者使用完全相同的類型定義：

```go
type Props map[string]interface{}
type PropsDef map[string]interface{}
```

## 一致性證明

### 1. 底層類型相同

兩者都是 `map[string]interface{}`，這意味著：

- ✅ 都可以存儲任何 Go 類型
- ✅ 都支援字串、布林、數字、切片、map 等所有類型
- ✅ 沒有任何類型限制上的差異

### 2. 類型處理規則相同

無論是 Props 還是 PropsDef，類型處理遵循相同的規則：

#### 規則 A：模板插值轉換為字串

```go
// PropsDef 定義
PropsDef{
    "disabled": false,  // bool
}

// Props 覆寫
Props{
    "disabled": true,   // bool
}

// 模板使用
Div(Props{"disabled": "{{disabled}}"})

// 結果：無論來自 PropsDef 還是 Props，
// 純模板引用會保持原始類型 bool
// 只有在渲染成 HTML 時才轉換為字串
```

#### 規則 B：非模板屬性保持原始類型

```go
// PropsDef 定義
PropsDef{
    "count": 10,        // int
    "enabled": true,    // bool
}

// Props 覆寫
Props{
    "count": 20,        // int
    "enabled": false,   // bool
}

// 模板不使用這些屬性
Div(Props{"class": "container"})

// 結果：無論來自 PropsDef 還是 Props，
// 純模板引用會保持原始類型 (int, bool)
```

#### 規則 C：值覆寫規則相同

```go
Component(
    Div(nil),
    nil,
    PropsDef{
        "value": 100,  // int - 預設值
    },
)

// Props 可以覆寫為任何類型
Props{"value": 200}         // int - 保持相同類型
Props{"value": 200.5}       // float64 - 改變類型
Props{"value": "200"}       // string - 改變類型
Props{"value": true}        // bool - 改變類型

// 覆寫後的值會保持傳入時的類型
```

## 測試驗證

### 測試 1：布林值

```go
// PropsDef
PropsDef{"enabled": true}
// Props
Props{"enabled": false}
// 結果：兩者都保持 bool 類型
```

### 測試 2：數字

```go
// PropsDef
PropsDef{"count": 10, "price": 99.99}
// Props
Props{"count": 20, "price": 199.99}
// 結果：兩者都保持原始數字類型 (int, float64)
```

### 測試 3：複雜類型

```go
// PropsDef
PropsDef{
    "tags": []string{"a", "b"},
    "config": map[string]int{"max": 100},
}
// Props
Props{
    "tags": []string{"x", "y"},
    "config": map[string]int{"min": 0},
}
// 結果：兩者都保持原始類型 ([]string, map[string]int)
```

## 使用示例

### 示例 1：兩者可互換使用

```go
// 定義組件時使用 PropsDef
MyComponent := Component(
    Div(Props{"class": "container"}),
    nil,
    PropsDef{
        "disabled": false,
        "count": 10,
    },
)

// 使用時用 Props 覆寫
result := MyComponent(Props{
    "disabled": true,
    "count": 20,
})

// 類型保持一致
fmt.Printf("%T\n", result.Props["disabled"]) // bool
fmt.Printf("%T\n", result.Props["count"])    // int
```

### 示例 2：直接使用 Props

```go
// 不使用組件，直接用 Props 創建元素
btn := Button(Props{
    "disabled": false,    // bool
    "tabindex": 0,        // int
    "type": "submit",     // string
})

// 類型處理方式與 PropsDef 完全相同
```

### 示例 3：類型混用

```go
Component(
    Div(nil),
    nil,
    PropsDef{
        "value": 100,  // int
    },
)

// Props 可以傳入不同類型
result1 := MyComponent(Props{"value": 200})      // int
result2 := MyComponent(Props{"value": 200.5})    // float64
result3 := MyComponent(Props{"value": "200"})    // string

// 每個結果都保持傳入時的類型
```

## 核心原則

### ✅ 完全一致的特性

1. **類型定義一致**
   - 都是 `map[string]interface{}`
   - 都支援所有 Go 類型

2. **處理規則一致**
   - 純模板引用 `{{key}}` 保持原始類型
   - 混合字串模板會轉為字串
   - 非模板屬性保持原始類型
   - 值覆寫規則完全相同
   - 只有渲染成 HTML 時才轉換為字串

3. **使用方式一致**
   - 可以用相同的語法定義
   - 可以存儲相同的類型
   - 可以執行相同的操作

### 📝 唯一的區別

唯一的區別是**用途**，而非類型處理：

- **PropsDef**：定義組件的預設屬性值
- **Props**：傳入實際的屬性值

但從類型系統的角度來看，它們是完全相同的。

## 最佳實踐

### ✅ 推薦做法

```go
// 在 PropsDef 中使用正確的類型
PropsDef{
    "disabled": false,    // ✓ bool
    "count": 10,          // ✓ int
    "price": 99.99,       // ✓ float64
}

// 在 Props 中也使用正確的類型
Props{
    "disabled": true,     // ✓ bool
    "count": 20,          // ✓ int
    "price": 199.99,      // ✓ float64
}
```

### ❌ 避免的做法

```go
// 不要在 PropsDef 或 Props 中使用字串表示其他類型
PropsDef{
    "disabled": "false",  // ✗ 不要用字串
    "count": "10",        // ✗ 不要用字串
}

Props{
    "disabled": "true",   // ✗ 不要用字串
    "count": "20",        // ✗ 不要用字串
}
```

## 總結

### 問題答案

**props 跟 propsDef 的類型處理一致嗎？**

✅ **是的，完全一致！**

### 關鍵要點

1. ✅ 兩者都是 `map[string]interface{}`
2. ✅ 支援相同的所有類型
3. ✅ 遵循相同的類型處理規則
4. ✅ 純模板引用 `{{key}}` 保持原始類型
5. ✅ 只有渲染成 HTML 時才轉換為字串
6. ✅ 類型保留規則相同
7. ✅ 沒有任何處理上的差異

### 設計理念

Props 和 PropsDef 使用統一的類型系統，這樣設計的好處：

- 🎯 **簡單**：只需要理解一套規則
- 🎯 **一致**：PropsDef 預設值和 Props 傳入值行為相同
- 🎯 **靈活**：可以自由選擇任何類型
- 🎯 **類型安全**：利用 Go 的類型系統
- 🎯 **延遲轉換**：只在渲染時轉換，Props 中保持類型

---

**文檔版本**: 1.0.0
**最後更新**: 2025-01-24
**作者**: TimLai666
