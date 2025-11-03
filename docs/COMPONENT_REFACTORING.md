# 組件重構文檔

## 概述

本次重構移除了 `Component` 函數中的所有派生屬性計算邏輯，改為在組件模板中使用條件表達式。這使得：

1. **關注點分離**：每個組件自己管理自己的邏輯
2. **通用性**：`Component` 函數保持通用，不包含特定組件的邏輯
3. **可維護性**：派生邏輯直接在模板中可見，更容易理解和修改

## 模板表達式系統

### 支持的表達式

#### 1. 簡單條件判斷
```
${'{{prop}}' === 'value' ? 'result1' : 'result2'}
```

#### 2. 不等於判斷
```
${'{{prop}}' !== 'value' ? 'result1' : 'result2'}
```

#### 3. 字符串 trim 檢查
```
${'{{prop}}'.trim() ? 'hasValue' : 'empty'}
```

#### 4. 嵌套三元表達式
```
${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1.125rem' : '0.95rem'}
```

### 表達式處理流程

1. **替換 {{}} 變量**：先將 `{{prop}}` 替換為實際值
2. **評估 ${} 表達式**：遞歸處理嵌套的三元運算符
3. **返回結果**：將計算結果插入到模板中

## 已更新的組件

### 1. Alert 組件 ✅

**移除的派生屬性**：
- `closeDisplay`, `iconDisplay`, `titleDisplay`
- `padding`, `alertRadius`, `alignItems`, `boxShadow`
- `bgColor`, `textColor`, `border`, `titleColor`, `iconColor`
- `alertIcon`

**新增的模板表達式**：
```go
// closable 控制關閉按鈕顯示
display: ${'{{closable}}' === 'true' ? 'block' : 'none'}

// type 控制顏色主題
background-color: ${'{{type}}' === 'success' ? '#f0fdf4' : '{{type}}' === 'warning' ? '#fffbeb' : '{{type}}' === 'error' ? '#fef2f2' : '#eef2ff'}

// rounded 控制圓角
border-radius: ${'{{rounded}}' === 'none' ? '0' : '{{rounded}}' === 'sm' ? '0.25rem' : '{{rounded}}' === 'lg' ? '0.75rem' : '0.5rem'}

// compact 控制 padding
padding: ${'{{compact}}' === 'true' ? '0.75rem 1rem' : '1rem 1.25rem'}

// elevation 控制陰影
box-shadow: ${'{{elevation}}' === '0' ? 'none' : '{{elevation}}' === '1' ? '0 1px 2px rgba(0,0,0,0.05)' : '{{elevation}}' === '2' ? '0 1px 3px rgba(0,0,0,0.1)' : '{{elevation}}' === '3' ? '0 4px 6px rgba(0,0,0,0.1)' : '0 10px 15px rgba(0,0,0,0.1)'}
```

### 2. Button 組件 ✅

**移除的派生屬性**：
- `fontSize`, `paddingX`, `paddingY`, `buttonRadius`
- `width`, `height`, `cursor`, `opacity`
- `background`, `textColor`, `border`, `boxShadow`
- `hoverBackground`, `hoverTextColor`, `hoverBorderColor`, `hoverBoxShadow`
- `focusRingColor`
- `iconLeft`, `iconRight`, `iconLeftDisplay`, `iconRightDisplay`

**新增的模板表達式**：
```go
// size 控制字體大小和 padding
font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1.125rem' : '0.95rem'}
padding: ${'{{size}}' === 'sm' ? '0.375rem 1rem' : '{{size}}' === 'lg' ? '0.625rem 1.5rem' : '0.5rem 1.25rem'}

// variant 控制樣式
background: ${'{{variant}}' === 'outlined' ? 'transparent' : '{{variant}}' === 'text' ? 'transparent' : '{{color}}'}
color: ${'{{variant}}' === 'outlined' ? '{{color}}' : '{{variant}}' === 'text' ? '{{color}}' : '#ffffff'}
border: ${'{{variant}}' === 'outlined' ? '1px solid {{color}}' : '1px solid transparent'}

// disabled 控制狀態
cursor: ${'{{disabled}}' === 'true' ? 'not-allowed' : 'pointer'}
opacity: ${'{{disabled}}' === 'true' ? '0.6' : '1'}

// fullWidth 控制寬度
width: ${'{{fullWidth}}' === 'true' ? '100%' : 'auto'}

// rounded 控制圓角
border-radius: ${'{{rounded}}' === 'none' ? '0' : '{{rounded}}' === 'sm' ? '0.25rem' : '{{rounded}}' === 'lg' ? '0.75rem' : '{{rounded}}' === 'full' ? '9999px' : '0.5rem'}
```

### 3. Radio 組件 ✅

**移除的派生屬性**：
- `flexDirection`, `labelDisplay`, `helpDisplay`

**新增的模板表達式**：
```go
// direction 控制排列方向
flex-direction: ${'{{direction}}' === 'horizontal' ? 'row' : 'column'}

// label 控制標籤顯示
display: ${'{{label}}'.trim() ? 'block' : 'none'}

// helpText 控制幫助文字顯示
display: ${'{{helpText}}'.trim() ? 'block' : 'none'}
```

### 4. Checkbox 組件 ✅

**移除的派生屬性**：
- `flexDirection`, `labelDisplay`, `helpDisplay`

**新增的模板表達式**：（同 Radio 組件）

### 5. Switch 組件 ✅

**移除的派生屬性**：
- `flexDirection`, `helpDisplay`

**新增的模板表達式**：
```go
// direction 控制排列方向（Switch 默認為 row）
flex-direction: ${'{{direction}}' === 'column' ? 'column' : 'row'}

// helpText 控制幫助文字顯示
display: ${'{{helpText}}'.trim() ? 'block' : 'none'}
```

### 6. Card 組件 ✅

**移除的派生屬性**：
- `titleDisplay`, `shadowY`, `shadowBlur`, `shadowOpacity`, `cursor`

**新增的模板表達式**：
```go
// title 控制標題顯示
display: ${'{{title}}'.trim() ? 'block' : 'none'}

// elevation 控制陰影
box-shadow: ${'{{elevation}}' === '0' ? 'none' : '{{elevation}}' === '1' ? '0 1px 3px rgba(0,0,0,0.05)' : '{{elevation}}' === '2' ? '0 4px 16px rgba(0,0,0,0.08)' : '{{elevation}}' === '3' ? '0 8px 24px rgba(0,0,0,0.10)' : '{{elevation}}' === '4' ? '0 12px 32px rgba(0,0,0,0.12)' : '0 16px 40px rgba(0,0,0,0.14)'}

// hoverable 控制游標
cursor: ${'{{hoverable}}' === 'false' ? 'default' : 'pointer'}
```

### 7. Modal 組件 ✅

**移除的派生屬性**：
- `display`, `pointerEvents`, `modalWidth`, `maxWidth`, `margin`
- `maxHeight`, `borderRadius`, `boxShadow`
- `headerDisplay`, `footerDisplay`, `closeButtonDisplay`
- `contentOverflow`, `overlayOverflow`
- `contentAnimation`, `overlayAnimation`

**新增的模板表達式**：
```go
// open 控制顯示
display: ${'{{open}}' === 'true' ? 'block' : 'none'}
pointer-events: ${'{{open}}' === 'true' ? 'auto' : 'none'}

// size 控制寬度
width: ${'{{size}}' === 'xs' ? '300px' : '{{size}}' === 'sm' ? '400px' : '{{size}}' === 'md' ? '500px' : '{{size}}' === 'lg' ? '700px' : '{{size}}' === 'xl' ? '900px' : '100%'}

// centered 控制位置
margin: ${'{{centered}}' === 'true' ? '3.75rem auto' : '1rem auto'}

// radius 控制圓角
border-radius: ${'{{radius}}' === 'none' ? '0' : '{{radius}}' === 'sm' ? '0.25rem' : '{{radius}}' === 'lg' ? '0.75rem' : '0.5rem'}

// hideHeader 和 title 控制頭部顯示
display: ${'{{hideHeader}}' === 'true' || '{{title}}'.trim() === '' ? 'none' : 'flex'}

// closeButton 控制關閉按鈕
display: ${'{{closeButton}}' === 'false' ? 'none' : 'block'}

// scrollable 控制內容滾動
overflow-y: ${'{{scrollable}}' === 'true' ? 'auto' : 'visible'}

// animation 控制動畫效果
animation: ${'{{animation}}' === 'slide' ? 'modalSlideIn' : '{{animation}}' === 'zoom' ? 'modalZoomIn' : 'modalFadeIn'} 0.3s ease
```

### 8. Table 組件 ✅

**移除的派生屬性**：
- `tfootDisplay`

**新增的模板表達式**：
```go
// footer 控制表尾顯示
display: ${'{{footer}}'.trim() ? 'table-footer-group' : 'none'}
```

### 9. Dropdown 組件 ⚠️

**狀態**：需要更新（包含複雜的響應式佈局計算）

### 10. Input 組件 ⚠️

**狀態**：需要更新（包含複雜的圖標和標籤定位）

## dom/component.go 更新

### 移除的派生屬性計算

**之前**：
```go
// 計算派生屬性（避免在模板中使用 JS-style 的 ${...} 表達式）
// labelDisplay: 預設用於群組 label（block / none）
if lbl, ok := mergedProps["label"]; ok && strings.TrimSpace(fmt.Sprint(lbl)) != "" {
    mergedProps["labelDisplay"] = "block"
} else {
    mergedProps["labelDisplay"] = "none"
}

// helpDisplay: 當 helpText 有內容時為 block
if ht, ok := mergedProps["helpText"]; ok && strings.TrimSpace(fmt.Sprint(ht)) != "" {
    mergedProps["helpDisplay"] = "block"
} else {
    mergedProps["helpDisplay"] = "none"
}

// flexDirection: direction => row / column
if dir, ok := mergedProps["direction"]; ok && strings.TrimSpace(fmt.Sprint(dir)) == "horizontal" {
    mergedProps["flexDirection"] = "row"
} else {
    mergedProps["flexDirection"] = "column"
}
```

**現在**：
```go
// Component 函數只負責合併 props 和模板插值
// 所有派生邏輯都在組件模板中使用 ${} 表達式處理
```

### 新增的表達式評估函數

```go
// evaluateExpression 評估簡單的條件表達式
// 支持格式：
// - 'value'.trim() ? 'A' : 'B'
// - 'X' === 'Y' ? 'A' : 'B'
// - 'X' !== 'Y' ? 'A' : 'B'
// - 嵌套三元: 'X' === 'Y' ? 'A' : 'Z' === 'W' ? 'B' : 'C'
func evaluateExpression(expr string) string

// evaluateCondition 評估條件表達式
func evaluateCondition(condition string) bool

// indexOfOperator 找到運算符的位置（不在引號內）
func indexOfOperator(s, op string) int

// unquote 移除字符串兩端的引號
func unquote(s string) string
```

## 測試覆蓋

### Alert 組件測試
- ✅ `TestAlertClosable` - closable 參數
- ✅ `TestAlertType` - 四種類型的顏色主題
- ✅ `TestAlertTitle` - title 顯示/隱藏
- ✅ `TestAlertRounded` - 四種圓角大小
- ✅ `TestAlertCompact` - compact 模式
- ✅ `TestAlertIcon` - icon 顯示/隱藏

### Button 組件測試
- ✅ `TestButtonSize` - 三種尺寸
- ✅ `TestButtonVariant` - 三種變體樣式
- ✅ `TestButtonDisabled` - disabled 狀態
- ✅ `TestButtonFullWidth` - fullWidth 參數
- ✅ `TestButtonRounded` - 五種圓角大小

### Radio/Checkbox/Switch 組件測試
- ✅ `TestRadioDirection` - 方向控制
- ✅ `TestRadioLabel` - 標籤顯示
- ✅ `TestCheckboxDirection` - 方向控制
- ✅ `TestSwitchDirection` - 方向控制
- ✅ `TestHelpTextDisplay` - helpText 顯示

## 遷移指南

### 對於使用者

**無需任何更改！** 組件的 API 保持不變，所有參數的使用方式都相同。

```go
// 之前和現在的用法完全相同
Alert(Props{"closable": true, "type": "success"}, Text("操作成功"))
Btn(Props{"size": "lg", "variant": "outlined"}, Text("按鈕"))
```

### 對於組件開發者

如果要創建新組件或修改現有組件：

1. **不要在 `Component` 函數中添加派生屬性計算**
2. **在組件模板中使用 `${}` 表達式**

**示例**：

```go
// ❌ 錯誤：不要這樣做
// 在 component.go 中添加特定組件的邏輯
if myProp, ok := mergedProps["myProp"]; ok && myProp == "value" {
    mergedProps["myDerivedProp"] = "result"
}

// ✅ 正確：在組件模板中使用表達式
Props{
    "style": `
        property: ${'{{myProp}}' === 'value' ? 'result' : 'default'};
    `,
}
```

## 性能考慮

### 優點
- **更少的 Go 代碼執行**：不需要在 Go 中進行大量的字符串比較和條件判斷
- **延遲計算**：只有在渲染時才計算需要的屬性
- **可緩存**：模板表達式的結果可以在 HTML 生成時緩存

### 注意事項
- **表達式複雜度**：避免過於複雜的嵌套表達式
- **重複計算**：如果同一個計算在多處使用，考慮在 PropsDef 中定義

## 最佳實踐

### 1. 表達式可讀性

**推薦**：使用換行提高可讀性
```go
Props{
    "style": `
        font-size: ${'{{size}}' === 'sm' ? '0.875rem' :
                    '{{size}}' === 'lg' ? '1.125rem' :
                    '0.95rem'};
    `,
}
```

### 2. 默認值處理

**推薦**：最後一個分支作為默認值
```go
${'{{prop}}' === 'value1' ? 'result1' :
  '{{prop}}' === 'value2' ? 'result2' :
  'defaultResult'}
```

### 3. 布尔值判斷

**推薦**：明確比較字符串
```go
// ✅ 好
${'{{disabled}}' === 'true' ? 'not-allowed' : 'pointer'}

// ❌ 避免（因為所有 prop 值都是字符串）
${'{{disabled}}' ? 'not-allowed' : 'pointer'}
```

### 4. 空值檢查

**推薦**：使用 `.trim()` 檢查非空字符串
```go
${'{{title}}'.trim() ? 'block' : 'none'}
```

## 未來改進

### 短期
- [ ] 完成 Dropdown 組件更新
- [ ] 完成 Input 組件更新
- [ ] 添加更多組件測試

### 中期
- [ ] 支持更複雜的表達式（如算術運算）
- [ ] 添加表達式語法檢查工具
- [ ] 優化表達式解析性能

### 長期
- [ ] 考慮引入模板編譯器，在構建時預處理表達式
- [ ] 支持自定義表達式函數
- [ ] 提供可視化組件編輯器

## 總結

本次重構成功地將組件的派生邏輯從通用的 `Component` 函數中分離出來，移到了各個組件的模板定義中。這提高了代碼的模塊化程度、可維護性，並保持了向後兼容性。

**主要成果**：
- ✅ 8 個組件完全更新
- ✅ 16 個測試全部通過
- ✅ API 向後兼容
- ✅ 代碼更清晰、更易維護
