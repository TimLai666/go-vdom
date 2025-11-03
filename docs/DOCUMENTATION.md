# go-vdom 完整文檔

> 純 Go 語言實現的虛擬 DOM 庫，專注於服務器端 HTML 和 JavaScript 的生成

**版本**: v1.1.0
**更新日期**: 2025-01-24

---

## 📚 目錄

1. [簡介](#簡介)
2. [快速開始](#快速開始)
3. [核心功能](#核心功能)
4. [組件系統](#組件系統)
5. [JavaScript DSL](#javascript-dsl)
6. [模板表達式](#模板表達式)
7. [控制流](#控制流)
8. [模板序列化](#模板序列化)
9. [API 參考](#api-參考)
10. [最佳實踐](#最佳實踐)

---

## 簡介

### 什麼是 go-vdom？

`go-vdom` 是一個純 Go 語言實現的虛擬 DOM 庫，讓您能夠在 Go 中以聲明式、類型安全的方式構建網頁。不同於傳統的前端虛擬 DOM 框架（如 React、Vue），go-vdom 專注於服務器端渲染。

### 設計理念

- **類型安全**: 利用 Go 的類型系統在編譯時捕獲錯誤
- **零運行時**: 生成純靜態 HTML/JS，無需客戶端框架
- **DSL 優先**: 提供直觀的 DSL 而非字符串模板
- **組件化**: 支持可重用的組件系統
- **服務器優先**: 專為服務器端渲染設計

### 適用場景

✅ **推薦使用**

- 服務器端渲染（SSR）應用
- 傳統 Web 應用（MPA）
- 動態生成 HTML 郵件
- 管理後台頁面
- 靜態網站生成器

❌ **不推薦**

- 單頁應用（SPA）的客戶端渲染
- 實時響應式更新（建議使用 htmx 等搭配）

### 安裝

```bash
go get github.com/TimLai666/go-vdom@v1.1.0
```

---

## 快速開始

### Hello World

```go
package main

import (
    "fmt"
    "net/http"
    . "github.com/TimLai666/go-vdom/dom"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        page := Html(Props{},
            Head(Props{},
                Title(Props{}, Text("Hello World")),
            ),
            Body(Props{},
                H1(Props{}, Text("Hello, go-vdom!")),
                P(Props{}, Text("這是我的第一個頁面")),
            ),
        )

        html := Render(page)
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprint(w, html)
    })

    http.ListenAndServe(":8080", nil)
}
```

### 帶交互的示例

```go
package main

import (
    "fmt"
    "net/http"
    . "github.com/TimLai666/go-vdom/dom"
    js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        page := Html(Props{},
            Head(Props{},
                Title(Props{}, Text("互動示例")),
            ),
            Body(Props{},
                H1(Props{}, Text("計數器")),
                Div(Props{},
                    Button(Props{
                        "id": "counter-btn",
                        "onClick": js.Fn(nil,
                            js.Const("span", "document.getElementById('count')"),
                            js.Const("current", "parseInt(span.innerText)"),
                            js.SetText("span", "(current + 1).toString()"),
                        ),
                    }, Text("點擊 +1")),
                    Text(" 計數: "),
                    Span(Props{"id": "count"}, Text("0")),
                ),
            ),
        )

        html := Render(page)
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprint(w, html)
    })

    http.ListenAndServe(":8080", nil)
}
```

---

## 核心功能

### VNode（虛擬節點）

VNode 是 go-vdom 的核心數據結構，代表一個 HTML 元素或文本節點。

```go
type VNode struct {
    Tag      string         // HTML 標籤名（空字符串表示文本節點）
    Props    Props          // 屬性映射
    Children []VNode        // 子節點列表
    Content  string         // 文本內容
}
```

### Props（屬性系統）

Props 支持多種類型的值，會自動轉換為 HTML 屬性。

#### 支持的類型

```go
Props{
    // 字符串 - 直接使用
    "class": "container",
    "id":    "main",

    // 布爾值 - true 渲染為 "true"，false 渲染為 "false"
    "disabled": true,      // 渲染為 disabled="true"
    "hidden":   false,     // 渲染為 hidden="false"
    "required": true,      // 渲染為 required="true"

    // 整數 - 自動轉換為字符串
    "width":    800,       // 渲染為 width="800"
    "height":   600,       // 渲染為 height="600"
    "tabindex": 0,         // 渲染為 tabindex="0"

    // 浮點數 - 自動轉換為字符串
    "opacity": 0.8,        // 渲染為 opacity="0.8"
    "price":   19.99,      // 渲染為 price="19.99"

    // 陣列 - 自動序列化為 JSON 字符串
    "data-items": []string{"apple", "banana", "orange"},
    // 渲染為 data-items='["apple","banana","orange"]'

    "data-numbers": []int{1, 2, 3, 4, 5},
    // 渲染為 data-numbers='[1,2,3,4,5]'

    // Map - 自動序列化為 JSON 字符串
    "data-config": map[string]interface{}{
        "theme":    "dark",
        "fontSize": 14,
        "enabled":  true,
    },
    // 渲染為 data-config='{"enabled":true,"fontSize":14,"theme":"dark"}'

    // 結構體 - 自動序列化為 JSON 字符串
    "data-user": struct {
        Name  string
        Email string
    }{"John Doe", "john@example.com"},
    // 渲染為 data-user='{"Name":"John Doe","Email":"john@example.com"}'

    // JSAction - 事件處理（特殊處理）
    "onClick": js.Fn(nil, js.Alert("'Hello'")),
}
```

**複雜類型的 JSON 序列化**

當你傳遞陣列、map、或結構體等複雜類型作為 props 時，`Component` 函數會自動將它們序列化為 JSON 字符串。這使得你可以輕鬆地將 Go 的數據結構傳遞到 HTML 屬性中，並在客戶端 JavaScript 中使用。

```go
// 示例：傳遞複雜數據到組件
template := dom.VNode{
    Tag: "div",
    Props: dom.Props{
        "data-items":  "{{items}}",
        "data-config": "{{config}}",
    },
}

componentFn := dom.Component(template, nil)

// 使用複雜類型的 props
result := componentFn(dom.Props{
    "items": []string{"Apple", "Banana", "Orange"},
    "config": map[string]interface{}{
        "theme": "dark",
        "language": "zh-TW",
    },
})

// result 的 Props 會包含：
// "data-items": '["Apple","Banana","Orange"]'
// "data-config": '{"language":"zh-TW","theme":"dark"}'
```

在客戶端 JavaScript 中使用這些數據：

```javascript
// 從 data 屬性讀取並解析 JSON
const element = document.querySelector("[data-items]");
const items = JSON.parse(element.dataset.items);
console.log(items); // ["Apple", "Banana", "Orange"]

const config = JSON.parse(element.dataset.config);
console.log(config.theme); // "dark"
```

#### Props 工具函數

```go
// 合併多個 Props（後面的覆蓋前面的）
merged := MergeProps(props1, props2, props3)

// 克隆 Props（深拷貝）
cloned := CloneProps(original)
```

### HTML 元素

所有標準 HTML 元素都有對應的構造函數：

```go
// 基本結構
Html(props Props, children ...VNode) VNode
Head(props Props, children ...VNode) VNode
Body(props Props, children ...VNode) VNode

// 標題
H1, H2, H3, H4, H5, H6(props Props, children ...VNode) VNode

// 文本
P(props Props, children ...VNode) VNode
Span(props Props, children ...VNode) VNode
Text(content string) VNode

// 容器
Div(props Props, children ...VNode) VNode
Section(props Props, children ...VNode) VNode
Article(props Props, children ...VNode) VNode

// 列表
Ul(props Props, children ...VNode) VNode
Ol(props Props, children ...VNode) VNode
Li(props Props, children ...VNode) VNode

// 表單
Form(props Props, children ...VNode) VNode
Input(props Props) VNode
Button(props Props, children ...VNode) VNode
Select(props Props, children ...VNode) VNode
Option(props Props, children ...VNode) VNode
Textarea(props Props, children ...VNode) VNode
Label(props Props, children ...VNode) VNode

// 表格
Table(props Props, children ...VNode) VNode
Thead(props Props, children ...VNode) VNode
Tbody(props Props, children ...VNode) VNode
Tr(props Props, children ...VNode) VNode
Th(props Props, children ...VNode) VNode
Td(props Props, children ...VNode) VNode

// 媒體
Img(props Props) VNode
A(props Props, children ...VNode) VNode
Script(props Props, children ...VNode) VNode
Style(props Props, children ...VNode) VNode
Link(props Props) VNode

// 其他
Code(props Props, children ...VNode) VNode
Pre(props Props, children ...VNode) VNode
Strong(props Props, children ...VNode) VNode
Em(props Props, children ...VNode) VNode
```

### 渲染

```go
// 渲染 VNode 為 HTML 字符串
html := Render(vnode)

// 創建完整的 HTML 文檔（包含 doctype）
doc := Document(
    "頁面標題",
    []LinkInfo{
        {Rel: "stylesheet", Href: "/style.css"},
    },
    []string{"/script.js"}, // 外部腳本
    []VNode{Script(Props{}, Text("console.log('內聯腳本')"))}, // 內聯腳本
    Body(Props{},
        H1(Props{}, Text("內容")),
    ),
)
html := Render(doc)
```

---

## 組件系統

### 創建組件

組件是一個返回 VNode 的函數。

#### 簡單組件

```go
// 無狀態組件
func Card(title, content string) VNode {
    return Div(Props{"class": "card"},
        Div(Props{"class": "card-header"},
            H3(Props{}, Text(title)),
        ),
        Div(Props{"class": "card-body"},
            P(Props{}, Text(content)),
        ),
    )
}

// 使用
card := Card("標題", "內容")
```

#### 可配置組件

```go
// 接受 Props 和 children
func Alert(props Props, children ...VNode) VNode {
    severity := "info"
    if s, ok := props["severity"].(string); ok {
        severity = s
    }

    return Div(Props{
        "class": "alert alert-" + severity,
        "role":  "alert",
    }, children...)
}

// 使用
alert := Alert(Props{"severity": "success"},
    Text("操作成功！"),
)
```

#### 使用 Component 函數

go-vdom 提供了 `Component` 函數來創建可重用的組件，支持預設屬性和模板插值。

```go
// 定義組件模板和預設屬性
var MyButton = Component(
    Button(Props{
        "class": "btn btn-{{variant}}",
        "type":  "{{type}}",
        "disabled": "{{disabled}}",
    }, Text("{{label}}")),
    nil, // 可選的 JavaScript 回調
    PropsDef{ // 預設屬性
        "variant":  "primary",
        "type":     "button",
        "disabled": false,
        "label":    "按鈕",
    },
)

// 使用組件
btn1 := MyButton(Props{"label": "提交", "variant": "success"})
btn2 := MyButton(Props{"label": "取消", "variant": "danger"})
```

### 模板插值

組件模板支持 `{{key}}` 語法進行屬性插值：

```go
// 模板中的 {{name}} 會被替換為 props["name"] 的值
Div(Props{"id": "user-{{id}}"},
    H1(Props{}, Text("{{name}}")),
    P(Props{}, Text("Email: {{email}}")),
)

// 使用時
component(Props{
    "id":    "123",
    "name":  "張三",
    "email": "zhang@example.com",
})
```

### 內建 UI 組件

go-vdom 提供了一套完整的 UI 組件庫：

#### 按鈕組件 (Btn)

```go
import . "github.com/TimLai666/go-vdom/components"

Btn(Props{
    "id":       "submit-btn",
    "variant":  "filled",    // filled, outlined, text
    "color":    "#3b82f6",   // 自定義顏色
    "size":     "md",        // sm, md, lg
    "rounded":  "md",        // none, sm, md, lg, full
    "disabled": false,
    "fullWidth": false,
}, Text("提交"))
```

#### 輸入框組件 (TextField)

```go
TextField(Props{
    "id":          "email",
    "label":       "電子郵件",
    "type":        "email",
    "placeholder": "your@email.com",
    "icon":        "📧",
    "iconPosition": "left",  // left, right
    "variant":     "outlined", // outlined, filled, underlined
    "size":        "md",      // sm, md, lg
    "helpText":    "我們不會分享您的郵件",
    "errorText":   "",
    "required":    true,
    "disabled":    false,
})
```

#### 下拉選單 (Dropdown)

```go
Dropdown(Props{
    "id":           "country",
    "label":        "國家",
    "options":      "台灣,日本,美國,英國", // 逗號分隔
    "defaultValue": "台灣",
    "placeholder":  "請選擇",
    "required":     true,
})
```

#### 開關組件 (Switch)

```go
Switch(Props{
    "id":      "notifications",
    "label":   "啟用通知",
    "checked": true,
    "onColor": "#10b981",  // 開啟時的顏色
    "offColor": "#d1d5db", // 關閉時的顏色
    "size":    "md",       // sm, md, lg
})
```

#### 單選框 (Radio)

```go
Radio(Props{
    "id":      "option1",
    "name":    "choice",
    "label":   "選項 1",
    "checked": true,
    "color":   "#3b82f6",
})
```

#### 複選框 (Checkbox)

```go
Checkbox(Props{
    "id":      "agree",
    "label":   "我同意條款",
    "checked": false,
    "color":   "#3b82f6",
})
```

#### 警告框 (Alert)

```go
Alert(Props{
    "id":       "success-msg",
    "severity": "success", // success, info, warning, error
    "title":    "成功",
    "closable": true,
}, Text("操作已成功完成！"))
```

#### 卡片 (Card)

```go
Card(Props{
    "title":    "卡片標題",
    "subtitle": "副標題",
    "elevated": true,
},
    P(Props{}, Text("卡片內容")),
)
```

#### 模態框 (Modal)

```go
Modal(Props{
    "id":         "confirm-modal",
    "title":      "確認刪除",
    "size":       "md", // sm, md, lg
    "closeButton": true,
},
    P(Props{}, Text("確定要刪除嗎？")),
)
```

---

## JavaScript DSL

### 基本函數

#### Fn - 普通函數

```go
js.Fn(params []string, actions ...JSAction) JSAction
```

創建普通 JavaScript 函數。

```go
// 無參數
js.Fn(nil,
    js.Log("'Hello'"),
    js.Alert("'World'"),
)

// 有參數
js.Fn([]string{"event", "data"},
    js.Log("event"),
    js.Const("value", "data.value"),
)
```

#### AsyncFn - 異步函數

```go
js.AsyncFn(params []string, actions ...JSAction) JSAction
```

創建異步函數，支持 `await` 語法。

```go
// ✅ 正確 - 使用 AsyncFn
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Log("data"),
    ),
}, Text("載入數據"))

// ❌ 錯誤 - 使用 Fn 會報錯
Button(Props{
    "onClick": js.Fn(nil,
        js.Const("response", "await fetch('/api/data')"), // 錯誤！
    ),
}, Text("載入數據"))
```

### DOM 操作

#### 選擇器

```go
// 通過選擇器獲取元素
js.El("#id")           // document.querySelector('#id')
js.ElAll(".class")     // document.querySelectorAll('.class')

// 通過 ID 獲取元素
js.GetById("myId")     // document.getElementById('myId')
```

#### 元素操作

```go
// 設置文本
js.SetText("element", "'新文本'")

// 設置 HTML
js.SetHTML("element", "'<b>HTML</b>'")

// 設置屬性
js.SetAttr("element", "disabled", "true")

// 添加/移除類
js.AddClass("element", "active")
js.RemoveClass("element", "hidden")
js.ToggleClass("element", "selected")

// 設置樣式
js.SetStyle("element", "color", "'red'")

// 鏈式調用
js.El("#btn").SetText("'點擊'").AddClass("active")
```

### 變量聲明

```go
// const 聲明
js.Const("name", "'value'")
js.Const("num", "42")

// let 聲明
js.Let("counter", "0")

// var 聲明
js.Var("global", "true")
```

### 控制流

```go
// if 語句
js.If("x > 0",
    js.Log("'正數'"),
)

// if-else 語句
js.IfElse("x > 0",
    js.Log("'正數'"),
    js.Log("'非正數'"),
)

// switch 語句
js.Switch("value",
    []js.Case{
        {Value: "'a'", Actions: []JSAction{js.Log("'A'")}},
        {Value: "'b'", Actions: []JSAction{js.Log("'B'")}},
    },
    []JSAction{js.Log("'默認'")}, // default case
)

// for 循環
js.For("let i = 0", "i < 10", "i++",
    js.Log("i"),
)

// while 循環
js.While("condition",
    js.Log("'循環中'"),
)
```

### 錯誤處理

```go
// try-catch
js.TryCatch(
    js.AsyncFn(nil,
        js.Const("response", "await fetch('/api')"),
        js.Const("data", "await response.json()"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'Error:', e.message"),
        js.Alert("'請求失敗'"),
    )),
    nil, // finally (可選)
)

// try-catch-finally
js.TryCatch(
    js.Fn(nil, js.Log("'嘗試'")),
    js.Ptr(js.Fn(nil, js.Log("'錯誤'"))),
    js.Ptr(js.Fn(nil, js.Log("'總是執行'"))),
)
```

### Fetch API

```go
// GET 請求
js.AsyncFn(nil,
    js.Const("response", "await fetch('/api/users')"),
    js.Const("data", "await response.json()"),
    js.Log("data"),
)

// POST 請求
js.AsyncFn(nil,
    js.Const("response", `await fetch('/api/users', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({name: 'John'})
    })`),
    js.IfElse("response.ok",
        js.Log("'成功'"),
        js.Log("'失敗'"),
    ),
)
```

### 事件處理

```go
// 添加事件監聽器
js.AddEventListener("document", "DOMContentLoaded", js.Fn(nil,
    js.Log("'頁面已載入'"),
))

// 移除事件監聽器
js.RemoveEventListener("element", "click", "handler")

// 觸發事件
js.DispatchEvent("element", "new CustomEvent('myEvent', {detail: {}})")
```

### 實用函數

```go
// 日誌
js.Log("'消息'", "變量")
js.Warn("'警告'")
js.Error("'錯誤'")

// 定時器
js.SetTimeout(js.Fn(nil, js.Log("'延遲'")), "1000")
js.SetInterval(js.Fn(nil, js.Log("'重複'")), "1000")

// 其他
js.Alert("'提示'")
js.Confirm("'確認?'")
js.Prompt("'輸入:'")
js.ConsoleDir("object")
```

---

## 模板表達式

組件模板支持強大的表達式系統，在服務器端評估條件邏輯。

### 基本語法

```go
// 條件表達式（三元運算符）
${'{{prop}}' === 'value' ? 'result1' : 'result2'}

// 嵌套三元運算符
${'{{size}}' === 'sm' ? '0.875rem' :
  '{{size}}' === 'md' ? '1rem' :
  '{{size}}' === 'lg' ? '1.125rem' : '1rem'}

// 比較運算符
${'{{count}}' > '0' ? 'block' : 'none'}
${'{{name}}' !== '' ? 'visible' : 'hidden'}
```

### 字符串檢查

```go
// 檢查是否為空（需要手動 trim）
${'{{text}}'.trim() ? 'block' : 'none'}

// 不等於空字符串
${'{{value}}' !== '' ? 'show' : 'hide'}
```

### 實際應用示例

```go
// 按鈕樣式根據變體變化
Button(Props{
    "style": `
        background: ${'{{variant}}' === 'filled' ? '{{color}}' : 'transparent'};
        color: ${'{{variant}}' === 'filled' ? 'white' : '{{color}}'};
        border: ${'{{variant}}' === 'outlined' ? '1px solid {{color}}' : 'none'};
    `,
})

// 根據尺寸調整字體
Div(Props{
    "style": `
        font-size: ${'{{size}}' === 'sm' ? '0.875rem' : '{{size}}' === 'lg' ? '1.125rem' : '1rem'};
        padding: ${'{{size}}' === 'sm' ? '0.5rem' : '{{size}}' === 'lg' ? '0.75rem' : '0.625rem'};
    `,
})

// 條件顯示
Div(Props{
    "style": `
        display: ${'{{visible}}' === 'true' ? 'block' : 'none'};
    `,
})
```

### 注意事項

1. **引號很重要**: 表達式中的字符串必須用引號包圍

   ```go
   // ✅ 正確
   ${'{{value}}' === 'test' ? 'a' : 'b'}

   // ❌ 錯誤
   ${{{value}} === test ? a : b}
   ```

2. **不支持邏輯運算符**: 使用嵌套三元運算符代替

   ```go
   // ❌ 不支持
   ${'{{a}}' && '{{b}}' ? 'yes' : 'no'}

   // ✅ 使用嵌套三元
   ${'{{a}}' ? '{{b}}' ? 'yes' : 'no' : 'no'}
   ```

3. **服務器端評估**: 表達式在渲染時（服務器端）評估，不是在客戶端

---

## 控制流

### 條件渲染

使用 `control` 包進行條件渲染。

```go
import ctrl "github.com/TimLai666/go-vdom/control"

// If-Then
ctrl.If(user != nil,
    ctrl.Then(
        P(Props{}, Text("歡迎, " + user.Name)),
    ),
)

// If-Then-Else
ctrl.If(user != nil,
    ctrl.Then(
        P(Props{}, Text("歡迎, " + user.Name)),
    ),
    ctrl.Else(
        P(Props{}, Text("請先登入")),
    ),
)

// 多條件
ctrl.If(score >= 90,
    ctrl.Then(H3(Props{}, Text("優秀"))),
    ctrl.ElseIf(score >= 80,
        ctrl.Then(H3(Props{}, Text("良好"))),
        ctrl.Else(H3(Props{}, Text("需要努力"))),
    ),
)
```

### 列表渲染

```go
import ctrl "github.com/TimLai666/go-vdom/control"

// For - 遍歷切片
users := []User{{Name: "Alice"}, {Name: "Bob"}}

Ul(Props{},
    ctrl.For(users, func(user User, index int) VNode {
        return Li(Props{}, Text(user.Name))
    }),
)

// Repeat - 重複 n 次
Div(Props{},
    ctrl.Repeat(5, func(i int) VNode {
        return P(Props{}, Text(fmt.Sprintf("第 %d 項", i+1)))
    }),
)

// Map - 映射函數
items := []string{"a", "b", "c"}
mapped := ctrl.Map(items, func(item string, i int) VNode {
    return Span(Props{}, Text(item))
})
```

---

## 模板序列化

### 保存為 Go Template

```go
import . "github.com/TimLai666/go-vdom/dom"

// 創建帶模板變數的 VNode
vnode := Div(Props{"id": "user-{{.ID}}"},
    H3(Props{}, Text("{{.Name}}")),
    P(Props{}, Text("Email: {{.Email}}")),
)

// 保存為 Go Template 格式
template := SaveTemplate("user-card", vnode)
// 寫入文件
os.WriteFile("user-card.tmpl", []byte(template), 0644)
```

生成的模板：

```html
{{/* Template: user-card */}} {{define "user-card"}}
<div id="user-{{.ID}}">
  <h3>{{.Name}}</h3>
  <p>Email: {{.Email}}</p>
</div>
{{end}}
```

### JSON 序列化

```go
// 轉為 JSON
jsonStr, err := ToJSON(vnode)

// 從 JSON 載入
restored, err := FromJSON(jsonStr)

// 渲染
html := Render(restored)
```

### 提取模板變數

```go
vnode := Div(Props{"id": "user-{{.ID}}"},
    H1(Props{}, Text("{{.Name}}")),
    P(Props{}, Text("{{.Email}}")),
)

// 提取所有模板變數
vars := ExtractTemplateVars(vnode)
// 返回: [".ID", ".Name", ".Email"]
```

### VNode 克隆

```go
// 創建原始 VNode
original := Button(Props{"class": "btn"}, Text("按鈕"))

// 克隆並修改
cloned := CloneVNode(original)
cloned.Props["class"] = "btn btn-primary"

// 原始 VNode 不受影響
```

---

## API 參考

### VNode 構造函數

所有 HTML 元素的完整列表請參考[核心功能 - HTML 元素](#html-元素)。

### Props 工具函數

```go
// 合併 Props
MergeProps(props ...Props) Props

// 克隆 Props
CloneProps(p Props) Props

// 轉換 Props 值類型
ConvertPropsToAny(p map[string]interface{}) Props
```

### 渲染函數

```go
// 渲染 VNode 為 HTML
Render(node VNode) string

// 創建完整 HTML 文檔
Document(title string, links []LinkInfo, scripts []string,
         inlineScripts []VNode, body VNode) VNode
```

### Component 函數

```go
Component(template VNode, onDOMReadyCallback *JSAction,
          defaultProps ...PropsDef) func(props Props, children ...VNode) VNode
```

### 控制流函數

```go
// 條件渲染
ctrl.If(condition bool, branches ...VNode) []VNode
ctrl.Then(nodes ...VNode) VNode
ctrl.Else(nodes ...VNode) VNode
ctrl.ElseIf(condition bool, branches ...VNode) VNode

// 列表渲染
ctrl.For[T any](items []T, fn func(T, int) VNode) []VNode
ctrl.Repeat(count int, fn func(int) VNode) []VNode
ctrl.Map[T any](items []T, fn func(T, int) VNode) []VNode
```

### JavaScript DSL 完整 API

請參考 [JavaScript DSL](#javascript-dsl) 章節。

---

## 最佳實踐

### 組件設計

#### 1. 保持組件簡單

```go
// ✅ 好：單一職責
func UserAvatar(url string, size int) VNode {
    return Img(Props{
        "src":    url,
        "width":  size,
        "height": size,
        "class":  "avatar",
    })
}

// ❌ 壞：做太多事情
func UserProfile(user User) VNode {
    // 包含頭像、個人信息、帖子列表等...
}
```

#### 2. 使用 Props 使組件可配置

```go
// ✅ 好：通過 Props 配置
func Card(props Props, children ...VNode) VNode {
    elevated := false
    if e, ok := props["elevated"].(bool); ok {
        elevated = e
    }

    shadow := "none"
    if elevated {
        shadow = "0 4px 6px rgba(0,0,0,0.1)"
    }

    return Div(Props{
        "class": "card",
        "style": "box-shadow: " + shadow,
    }, children...)
}
```

#### 3. 提取可重用的樣式

```go
// 定義樣式常量
var (
    PrimaryColor   = "#3b82f6"
    SuccessColor   = "#10b981"
    ErrorColor     = "#ef4444"

    ButtonBase = Props{
        "class": "btn",
        "style": "padding: 0.5rem 1rem; border-radius: 0.375rem;",
    }
)

// 使用
btn := Button(MergeProps(ButtonBase, Props{
    "style": "background: " + PrimaryColor,
}), Text("按鈕"))
```

### 性能優化

#### 1. 避免不必要的重新渲染

```go
// ✅ 好：緩存不變的部分
var cachedHeader = Header(Props{},
    H1(Props{}, Text("網站標題")),
    Nav(Props{}, /* ... */),
)

func Page(content VNode) VNode {
    return Html(Props{},
        Head(Props{}, /* ... */),
        Body(Props{},
            cachedHeader,  // 重用緩存的 header
            content,
        ),
    )
}
```

#### 2. 使用條件渲染避免生成不必要的 HTML

```go
// ✅ 好：使用控制流
ctrl.If(user != nil,
    ctrl.Then(UserDashboard(user)),
)

// ❌ 壞：總是生成 HTML 再用 CSS 隱藏
Div(Props{
    "style": func() string {
        if user == nil {
            return "display: none"
        }
        return ""
    }(),
}, UserDashboard(user))
```

#### 3. 大列表使用虛擬滾動或分頁

```go
// ✅ 好：分頁
func ItemList(items []Item, page, pageSize int) VNode {
    start := page * pageSize
    end := start + pageSize
    if end > len(items) {
        end = len(items)
    }

    return Ul(Props{},
        ctrl.For(items[start:end], func(item Item, i int) VNode {
            return Li(Props{}, Text(item.Name))
        }),
    )
}

// ❌ 壞：一次渲染所有項目
func ItemList(items []Item) VNode {
    return Ul(Props{},
        ctrl.For(items, func(item Item, i int) VNode {
            return Li(Props{}, Text(item.Name))
        }),
    )
}
```

### 錯誤處理

#### 1. 使用 TryCatch 處理異步錯誤

```go
Button(Props{
    "onClick": js.TryCatch(
        js.AsyncFn(nil,
            js.Const("response", "await fetch('/api/data')"),
            js.Const("data", "await response.json()"),
            js.Log("data"),
        ),
        js.Ptr(js.Fn(nil,
            js.Log("'Error:', e.message"),
            js.Alert("'請求失敗，請稍後再試'"),
        )),
        nil,
    ),
}, Text("載入數據"))
```

#### 2. 驗證用戶輸入

```go
Form(Props{
    "onSubmit": js.Fn([]string{"e"},
        js.Call("e.preventDefault", nil),
        js.Const("email", "document.getElementById('email').value"),
        js.If("!email.includes('@')",
            js.Alert("'請輸入有效的郵件地址'"),
            js.Call("return", nil),
        ),
        // 提交表單...
    ),
}, /* ... */)
```

### 代碼組織

#### 1. 按功能組織文件

```
/components
  /auth
    login.go
    register.go
  /layout
    header.go
    footer.go
  /user
    profile.go
    settings.go
```

#### 2. 使用包級別變數存儲組件

```go
package components

// 導出組件供其他包使用
var (
    Header  = headerComponent
    Footer  = footerComponent
    Sidebar = sidebarComponent
)

func headerComponent(props Props) VNode {
    // 實現...
}
```

#### 3. 使用工廠函數創建相似組件

```go
func makeButton(variant string) func(Props, ...VNode) VNode {
    return func(props Props, children ...VNode) VNode {
        mergedProps := MergeProps(Props{
            "class": "btn btn-" + variant,
        }, props)
        return Button(mergedProps, children...)
    }
}

var (
    PrimaryButton   = makeButton("primary")
    SecondaryButton = makeButton("secondary")
    DangerButton    = makeButton("danger")
)
```

---

## 常見問題

### Q: await 語法錯誤怎麼辦？

**A:** 使用 `AsyncFn` 而不是 `Fn`。

```go
// ✅ 正確
js.AsyncFn(nil, js.Const("data", "await fetch('/api')"))

// ❌ 錯誤
js.Fn(nil, js.Const("data", "await fetch('/api')"))
```

### Q: 如何處理表單提交？

**A:** 使用 `onSubmit` 事件和 `e.preventDefault()`。

```go
Form(Props{
    "onSubmit": js.AsyncFn([]string{"e"},
        js.Call("e.preventDefault", nil),
        js.Const("formData", "new FormData(e.target)"),
        js.Const("response", "await fetch('/api/submit', {method: 'POST', body: formData})"),
        js.IfElse("response.ok",
            js.Alert("'提交成功'"),
            js.Alert("'提交失敗'"),
        ),
    ),
}, /* 表單內容 */)
```

### Q: Props 支持哪些類型？

**A:** 支持字符串、布爾值、整數、浮點數和 JSAction。詳見 [Props 屬性系統](#props屬性系統)。

### Q: 如何優化性能？

**A:**

1. 緩存不變的組件
2. 使用條件渲染
3. 大列表使用分頁
4. 避免在循環中創建函數

### Q: 可以用於單頁應用（SPA）嗎？

**A:** 不推薦。go-vdom 是為服務器端渲染設計的。對於 SPA，建議使用 React、Vue 等客戶端框架。

---

## 相關資源

- **[GitHub 倉庫](https://github.com/TimLai666/go-vdom)** - 源代碼和 Issues
- **[示例程序](../examples/)** - 可運行的完整示例
- **[CHANGELOG](../CHANGELOG.md)** - 版本更新歷史
- **[快速參考](QUICK_REFERENCE.md)** - 語法速查表

---

**版本**: v1.1.0
**作者**: TimLai666
**許可**: MIT License
