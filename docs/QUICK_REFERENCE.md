# Go VDOM 快速參考

快速查找常用操作和語法。

## 目錄

- [基本結構](#基本結構)
- [HTML 元素](#html-元素)
- [Props 屬性](#props-屬性)
- [組件](#組件)
- [控制流](#控制流)
- [JavaScript DSL](#javascript-dsl)
- [UI 組件](#ui-組件)

---

## 基本結構

### 創建完整 HTML 文檔

```go
doc := Document(
    "頁面標題",
    []LinkInfo{{Rel: "stylesheet", Href: "style.css"}},
    []ScriptInfo{{Src: "script.js"}},
    []Props{{"charset": "UTF-8"}},
    Div("內容"),
)
```

### 渲染為 HTML

```go
html := Render(doc)
```

---

## HTML 元素

### 常用元素

```go
Div()                                    // <div></div>
Div("文本")                               // <div>文本</div>
Div(Props{"class": "box"})               // <div class="box"></div>
Div(Props{"class": "box"}, "文本")       // <div class="box">文本</div>

H1("標題")                                // <h1>標題</h1>
P("段落")                                 // <p>段落</p>
Span("文本")                              // <span>文本</span>
A(Props{"href": "/url"}, "連結")         // <a href="/url">連結</a>
```

### 嵌套元素

```go
Div(
    Props{"class": "container"},
    H1("標題"),
    P("段落"),
    Ul(
        Li("項目1"),
        Li("項目2"),
    ),
)
```

---

## Props 屬性

### 基本屬性

```go
Props{
    "id":    "myId",
    "class": "btn btn-primary",
    "style": "color: red;",
}
```

### 事件屬性

```go
Props{
    "onClick": js.Fn(nil, js.Alert("'Hi'")),
    "onChange": js.Fn([]string{"e"}, js.Log("e.target.value")),
}
```

---

## 組件

### 定義組件

```go
MyComponent := Component(
    Div(
        Props{"class": "{{className}}"},
        H2("{{title}}"),
        P("{{content}}"),
        Div("{{children}}"),
    ),
    nil,
    PropsDef{
        "title":     "默認標題",
        "content":   "",
        "className": "",
    },
)
```

### 使用組件

```go
instance := MyComponent(
    Props{
        "title":   "我的標題",
        "content": "內容",
    },
    P("子元素1"),
    P("子元素2"),
)
```

---

## 控制流

### If/Then/Else

```go
import control "github.com/TimLai666/go-vdom/control"

control.If(condition,
    control.Then(Div("顯示")),
    control.Else(Div("隱藏")),
)
```

### Repeat

```go
control.Repeat(5, func(i int) VNode {
    return Div(fmt.Sprintf("項目 %d", i))
})
```

### ForEach（遍歷集合）

```go
items := []string{"A", "B", "C"}
control.ForEach(items, func(item string, i int) VNode {
    return Li(item)
})
```

### For（傳統循環）

```go
// 正向：1 到 10
control.For(1, 11, 1, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})

// 倒序：10 到 1
control.For(10, 0, -1, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})

// 步進：偶數 0-18
control.For(0, 20, 2, func(i int) VNode {
    return Span(fmt.Sprintf("%d", i))
})
```

### ForEach（後端渲染）

```go
// 簡潔的列表渲染
Ul(ForEach(fruits, func(fruit string) VNode {
    return Li(fruit)
}))

// 帶索引
Ul(ForEachWithIndex(items, func(item string, i int) VNode {
    return Li(fmt.Sprintf("%d. %s", i+1, item))
}))
```

**對比：**
| 用法 | control.For | ForEach |
|------|-------------|---------|
| 語法 | `control.For(items, func(item, i) {...})` | `ForEach(items, func(item) {...})` |
| 索引 | 總是提供 | `ForEachWithIndex` 才提供 |
| 簡潔性 | 需要導入 control 包 | 直接可用 |
| 推薦 | 需要索引時 | 不需要索引時 |

---

## JavaScript DSL

### 導入

```go
import js "github.com/TimLai666/go-vdom/jsdsl"
```

### 基本操作

```go
js.Log("'訊息'")                          // console.log('訊息')
js.Alert("'警告'")                        // alert('警告')
js.Redirect("/home")                     // location.href = '/home'
```

### 變數定義

```go
js.Let("counter", "0")                   // let counter = 0
js.Const("name", "'Alice'")              // const name = 'Alice'
```

### DOM 操作

```go
element := js.El("#myId")
element.SetText("'文本'")                 // element.innerText = '文本'
element.SetHTML("'<b>HTML</b>'")         // element.innerHTML = '<b>HTML</b>'
element.AddClass("active")               // element.classList.add('active')
element.RemoveClass("hidden")            // element.classList.remove('hidden')
```

### 事件處理

```go
js.El("#btn").OnClick(
    js.Alert("'點擊'"),
)
```

### 列表遍歷（前端渲染）

```go
// 基本遍歷
js.ForEachJS("['A', 'B', 'C']", "item",
    js.Log("'項目: ' + item"),
)

// 帶索引遍歷
js.ForEachWithIndexJS("items", "item", "index",
    js.Log("'[' + index + '] = ' + item"),
)

// DOM 元素遍歷
js.ForEachElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
    return el.AddClass("'active'")
})
```

**前端 vs 後端：**
- **後端 ForEach**: `ForEach(items, func(item) VNode {...})`  
  → 在伺服器生成 HTML，SEO 友好
- **前端 js.ForEach**: `js.ForEachJS("array", "item", ...actions)`  
  → 在瀏覽器執行，適合動態數據

### 函數定義

```go
// 無參數
js.Fn(nil,
    js.Log("'執行'"),
)

// 有參數
js.Fn([]string{"name"},
    js.Log("'Hello, ' + name"),
)
```

### Try/Catch

```go
js.TryCatch(
    js.Fn(nil,
        js.Const("data", "await fetch('/api')"),
        js.Log("data"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'錯誤:', e"),
    )),
    nil,
)
```

### DomReady

```go
Script(Props{"type": "text/javascript"},
    js.DomReady(
        js.Log("'DOM 已就緒'"),
    ),
)
```

### Fetch GET

```go
Button(Props{
    "onClick": js.Fn(nil,
        js.TryCatch(
            js.Fn(nil,
                js.Const("response", "await fetch('/api/data')"),
                js.Const("data", "await response.json()"),
                js.Log("data"),
            ),
            js.Ptr(js.Fn(nil,
                js.Alert("'錯誤: ' + e.message"),
            )),
            nil,
        ),
    ),
}, "獲取數據")
```

### Fetch POST

```go
Form(Props{
    "onSubmit": js.Fn([]string{"evt"},
        js.CallMethod("evt", "preventDefault"),
        js.TryCatch(
            js.Fn(nil,
                js.Const("data", "{ name: 'test' }"),
                js.Const("response", "await fetch('/api', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) })"),
                js.Const("result", "await response.json()"),
                js.Log("result"),
            ),
            js.Ptr(js.Fn(nil,
                js.Alert("'提交失敗'"),
            )),
            nil,
        ),
    ),
}, /* form fields */)
```

---

## UI 組件

### 導入

```go
import comp "github.com/TimLai666/go-vdom/components"
```

### TextField

```go
comp.TextField(Props{
    "id":          "email",
    "label":       "電子郵件",
    "type":        "email",
    "placeholder": "請輸入郵箱",
    "required":    "true",
    "helpText":    "幫助文本",
})
```

### Dropdown

```go
comp.Dropdown(Props{
    "id":      "country",
    "label":   "國家",
    "options": "台灣,中國,日本",
    "value":   "台灣",
})
```

### RadioGroup

```go
comp.RadioGroup(Props{
    "id":           "gender",
    "name":         "gender",
    "label":        "性別",
    "options":      "男,女,其他",
    "defaultValue": "男",
    "direction":    "horizontal",
})
```

### Checkbox

```go
comp.Checkbox(Props{
    "id":       "agree",
    "label":    "我同意",
    "required": "true",
    "checked":  "true",
})
```

### CheckboxGroup

```go
comp.CheckboxGroup(Props{
    "id":      "hobbies",
    "label":   "愛好",
    "options": "閱讀,運動,音樂",
    "values":  "閱讀,音樂",
})
```

### Switch

```go
comp.Switch(Props{
    "id":            "notifications",
    "label":         "通知",
    "checked":       "true",
    "labelPosition": "right",
})
```

---

## 完整示例

### 簡單頁面

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    . "github.com/TimLai666/go-vdom/vdom"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        
        doc := Document(
            "我的網站",
            nil, nil, nil,
            Div(
                Props{"class": "container"},
                H1("歡迎"),
                P("這是內容"),
            ),
        )
        
        fmt.Fprint(w, Render(doc))
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 帶組件和 JavaScript

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    comp "github.com/TimLai666/go-vdom/components"
    control "github.com/TimLai666/go-vdom/control"
    js "github.com/TimLai666/go-vdom/jsdsl"
    . "github.com/TimLai666/go-vdom/vdom"
)

func main() {
    Card := Component(
        Div(
            Props{"class": "card"},
            H3("{{title}}"),
            P("{{content}}"),
        ),
        nil,
        PropsDef{"title": "", "content": ""},
    )
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        
        items := []string{"A", "B", "C"}
        
        doc := Document(
            "完整示例",
            []LinkInfo{{
                Rel: "stylesheet",
                Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
            }},
            nil, nil,
            Div(
                Props{"class": "container"},
                H1("完整示例"),
                
                Card(Props{
                    "title": "我的卡片",
                    "content": "內容",
                }),
                
                control.If(true,
                    control.Then(P("顯示此段落")),
                ),
                
                Ul(
                    control.ForEach(items, func(item string, i int) VNode {
                        return Li(item)
                    }),
                ),
                
                comp.TextField(Props{
                    "id": "name",
                    "label": "姓名",
                }),
                
                Button(Props{
                    "onClick": js.Fn(nil,
                        js.Alert("'Hello!'"),
                    ),
                }, "點擊我"),
            ),
        )
        
        fmt.Fprint(w, Render(doc))
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 提示

- Props 必須是第一個參數
- JavaScript 字符串需要加引號：`js.Log("'text'")`
- TryCatch 的 catch 和 finally 參數需要使用 `js.Ptr()`
- 組件模板使用 `{{propName}}` 占位符
- `{{children}}` 是特殊占位符，用於子元素
- **列表渲染**：
  - 後端遍歷集合 → `ForEach()` 或 `control.ForEach()`
  - 後端數字循環 → `control.For(start, end, step, ...)`
  - 前端動態數據 → `js.ForEachJS()`
  - DOM 元素操作 → `js.ForEachElement()`

---

**更多詳細信息請參閱 [DOCUMENTATION.md](DOCUMENTATION.md)**