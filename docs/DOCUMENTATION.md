# go-vdom 完整文檔

## 目錄

1. [簡介](#簡介)
2. [架構設計](#架構設計)
3. [核心概念](#核心概念)
4. [VDOM 模塊](#vdom-模塊)
5. [Control 模塊](#control-模塊)
6. [JavaScript DSL 模塊](#javascript-dsl-模塊)
7. [Components 模塊](#components-模塊)
8. [進階用法](#進階用法)
9. [性能優化](#性能優化)
10. [故障排除](#故障排除)

---

## 簡介

`go-vdom` 是一個純 Go 語言實現的虛擬 DOM 庫，專注於服務器端 HTML 和 JavaScript 的生成。它不同於傳統的前端虛擬 DOM 框架（如 React、Vue），而是提供了一套完整的 DSL，讓開發者能夠在 Go 中以聲明式、類型安全的方式構建網頁。

### 設計理念

- **類型安全**: 利用 Go 的類型系統在編譯時捕獲錯誤
- **零運行時**: 生成純靜態 HTML/JS，無需客戶端框架
- **DSL 優先**: 提供直觀的 DSL 而非字符串模板
- **組件化**: 支持可重用的組件系統
- **服務器優先**: 專為服務器端渲染設計

### 使用場景

- ✅ 服務器端渲染（SSR）應用
- ✅ 傳統 Web 應用（MPA）
- ✅ 動態生成 HTML 郵件
- ✅ 管理後台頁面
- ✅ 靜態網站生成器
- ❌ 單頁應用（SPA）的客戶端渲染
- ❌ 實時響應式更新（建議使用 htmx 等搭配）

---

## 架構設計

### 整體架構

```
┌─────────────────────────────────────────┐
│           Application Layer             │
│      (Your Go HTTP Handlers)            │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│         go-vdom Public API              │
├─────────────────────────────────────────┤
│  Components  │  Control  │  JSDSL       │
├──────────────┼───────────┼──────────────┤
│              VDOM Core                   │
│  (VNode, Props, Rendering)               │
└─────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│         Generated Output                │
│     HTML + JavaScript Strings            │
└─────────────────────────────────────────┘
```

### 模塊說明

#### vdom (核心)
- 虛擬 DOM 節點定義
- HTML 元素構建函數
- 渲染引擎
- 組件系統基礎

#### control (控制流)
- 條件渲染 (If/Then/Else)
- 循環渲染 (For/Repeat)
- 邏輯控制

#### jsdsl (JavaScript DSL)
- JavaScript 代碼生成
- DOM 操作
- 事件處理
- Fetch API 封裝
- 異步處理 (Try/Catch)

#### components (UI 組件)
- 表單組件
- 輸入組件
- 交互組件

---

## 核心概念

### VNode (虛擬節點)

VNode 是 go-vdom 的核心數據結構，代表一個 HTML 元素。

```go
type VNode struct {
    Tag      string    // HTML 標籤名
    Props    Props     // 元素屬性
    Children []VNode   // 子元素
    Text     string    // 文本內容
}
```

#### 創建 VNode

```go
// 方式 1: 純文本節點
text := VNode{Text: "Hello"}

// 方式 2: 使用 HTML 函數
div := Div("Hello")

// 方式 3: 帶屬性
div := Div(Props{"class": "container"}, "Hello")

// 方式 4: 嵌套子元素
div := Div(
    Props{"class": "card"},
    H1("標題"),
    P("內容"),
)
```

### Props (屬性)

Props 是一個字符串映射，用於設置 HTML 屬性。

```go
type Props map[string]string
```

#### 常用屬性

```go
Props{
    // 基本屬性
    "id":    "myElement",
    "class": "btn btn-primary",
    "style": "color: red;",
    
    // 數據屬性
    "data-id":    "123",
    "data-value": "test",
    
    // 布爾屬性
    "disabled": "true",
    "required": "true",
    "checked":  "true",
    
    // 事件屬性（接受 JSAction）
    "onClick":  jsAction,
    "onChange": jsAction,
    "onSubmit": jsAction,
}
```

#### 特殊處理

某些屬性會被特殊處理：

```go
// 事件屬性（on* 開頭）
"onClick": js.Fn(nil, js.Alert("'Hi'")),
// 渲染為: onclick="(function() { alert('Hi'); })()"

// JSAction 類型
"onClick": JSAction{Code: "alert('Hi')"},
// 渲染為: onclick="alert('Hi')"
```

### PropsDef (組件屬性定義)

PropsDef 定義組件的默認 props。

```go
type PropsDef map[string]string

// 使用示例
PropsDef{
    "title":   "默認標題",
    "color":   "blue",
    "visible": "true",
}
```

---

## VDOM 模塊

### Document 函數

創建完整的 HTML 文檔結構。

```go
func Document(
    title string,
    links []LinkInfo,
    scripts []ScriptInfo,
    metas []Props,
    body ...VNode,
) VNode
```

#### 參數說明

- `title`: 頁面標題（<title> 標籤）
- `links`: 外部資源連結（CSS、圖標等）
- `scripts`: JavaScript 腳本
- `metas`: Meta 標籤
- `body`: 頁面主體內容

#### 使用示例

```go
doc := Document(
    "我的網站",
    []LinkInfo{
        {
            Rel:  "stylesheet",
            Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
            Type: "text/css",
        },
        {
            Rel:  "icon",
            Href: "/favicon.ico",
            Type: "image/x-icon",
        },
    },
    []ScriptInfo{
        {
            Src:   "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js",
            Async: true,
        },
    },
    []Props{
        {"name": "description", "content": "網站描述"},
        {"name": "viewport", "content": "width=device-width, initial-scale=1"},
        {"charset": "UTF-8"},
    },
    // Body 內容
    Div(
        Props{"class": "container"},
        H1("歡迎"),
        P("這是主要內容"),
    ),
)
```

### HTML 元素函數

所有標準 HTML 元素都有對應的函數。

#### 結構元素

```go
// 容器
Div(children ...any) VNode
Span(children ...any) VNode
Section(children ...any) VNode
Article(children ...any) VNode
Aside(children ...any) VNode

// 語義化標籤
Header(children ...any) VNode
Footer(children ...any) VNode
Main(children ...any) VNode
Nav(children ...any) VNode
```

#### 文本元素

```go
// 標題
H1(children ...any) VNode
H2(children ...any) VNode
H3(children ...any) VNode
H4(children ...any) VNode
H5(children ...any) VNode
H6(children ...any) VNode

// 文本
P(children ...any) VNode
Span(children ...any) VNode
Strong(children ...any) VNode
Em(children ...any) VNode
Code(children ...any) VNode
Pre(children ...any) VNode
```

#### 列表元素

```go
Ul(children ...any) VNode  // 無序列表
Ol(children ...any) VNode  // 有序列表
Li(children ...any) VNode  // 列表項
Dl(children ...any) VNode  // 定義列表
Dt(children ...any) VNode  // 定義術語
Dd(children ...any) VNode  // 定義描述
```

#### 表單元素

```go
Form(children ...any) VNode
Input(children ...any) VNode
Textarea(children ...any) VNode
Select(children ...any) VNode
Option(children ...any) VNode
Button(children ...any) VNode
Label(children ...any) VNode
Fieldset(children ...any) VNode
Legend(children ...any) VNode
```

#### 多媒體元素

```go
Img(children ...any) VNode
Video(children ...any) VNode
Audio(children ...any) VNode
Source(children ...any) VNode
```

#### 表格元素

```go
Table(children ...any) VNode
Thead(children ...any) VNode
Tbody(children ...any) VNode
Tfoot(children ...any) VNode
Tr(children ...any) VNode
Th(children ...any) VNode
Td(children ...any) VNode
```

#### 其他元素

```go
A(children ...any) VNode       // 連結
Br(children ...any) VNode      // 換行
Hr(children ...any) VNode      // 分隔線
Script(children ...any) VNode  // 腳本
Style(children ...any) VNode   // 樣式
```

### 元素函數參數規則

元素函數接受可變參數 `children ...any`，支持以下類型：

```go
// 1. Props（必須是第一個參數）
Div(Props{"class": "container"}, ...)

// 2. 字符串（作為文本節點）
P("這是文本")

// 3. VNode（子元素）
Div(H1("標題"), P("段落"))

// 4. []VNode（多個子元素）
Div([]VNode{H1("標題"), P("段落")})

// 5. 混合使用
Div(
    Props{"class": "card"},
    H1("標題"),
    P("段落"),
    Button(Props{"class": "btn"}, "按鈕"),
)
```

### Component 函數

定義可重用的組件。

```go
func Component(
    template VNode,
    jsAction *JSAction,
    propsDef PropsDef,
) func(Props, ...VNode) VNode
```

#### 參數說明

- `template`: 組件模板（使用 `{{propName}}` 作為占位符）
- `jsAction`: 可選的 JavaScript 代碼（通常用於組件初始化）
- `propsDef`: Props 默認值定義

#### 模板占位符

```go
// {{propName}} - 替換為 prop 值
H1("{{title}}")

// {{children}} - 特殊占位符，替換為子元素
Div("{{children}}")
```

#### 完整示例

```go
// 1. 定義組件
Card := Component(
    Div(
        Props{"class": "card {{className}}", "style": "{{style}}"},
        Div(
            Props{"class": "card-header"},
            H3("{{title}}"),
        ),
        Div(
            Props{"class": "card-body"},
            P("{{description}}"),
            Div("{{children}}"),
        ),
        Div(
            Props{"class": "card-footer"},
            Small("{{footer}}"),
        ),
    ),
    nil, // 無 JSAction
    PropsDef{
        "title":       "無標題",
        "description": "",
        "footer":      "",
        "className":   "",
        "style":       "",
    },
)

// 2. 使用組件
cardInstance := Card(
    Props{
        "title":       "我的卡片",
        "description": "這是一個卡片組件",
        "footer":      "2025",
        "className":   "shadow",
    },
    // children
    Button(Props{"class": "btn btn-primary"}, "操作"),
    Button(Props{"class": "btn btn-secondary"}, "取消"),
)
```

### Render 函數

將 VNode 渲染為 HTML 字符串。

```go
func Render(vnode VNode) string
```

```go
html := Render(Div(
    Props{"class": "container"},
    H1("Hello"),
))
// 輸出: <div class="container"><h1>Hello</h1></div>
```

---

## Control 模塊

Control 模塊提供控制流結構，用於條件渲染和循環渲染。

### If/Then/Else

條件渲染。

```go
func If(condition bool, thenBlock VNode, elseBlock ...VNode) VNode
func Then(content ...VNode) VNode
func Else(content ...VNode) VNode
```

#### 基本用法

```go
isLoggedIn := true

content := control.If(isLoggedIn,
    control.Then(
        Div("歡迎回來！"),
    ),
    control.Else(
        Div("請登入"),
    ),
)
```

#### 多條件嵌套

```go
userRole := "admin"

content := control.If(userRole == "admin",
    control.Then(
        Div("管理員面板"),
    ),
    control.Else(
        control.If(userRole == "user",
            control.Then(
                Div("用戶面板"),
            ),
            control.Else(
                Div("訪客面板"),
            ),
        ),
    ),
)
```

#### 只有 Then

```go
showAlert := true

content := control.If(showAlert,
    control.Then(
        Div(Props{"class": "alert alert-warning"}, "警告訊息"),
    ),
)
```

### Repeat

重複渲染相同的元素。

```go
func Repeat(count int, fn func(int) VNode) VNode
```

#### 使用示例

```go
// 生成 5 個項目
items := control.Repeat(5, func(i int) VNode {
    return Div(
        Props{"class": "item"},
        fmt.Sprintf("項目 #%d", i+1),
    )
})

// 生成表格行
rows := Table(
    Tbody(
        control.Repeat(10, func(i int) VNode {
            return Tr(
                Td(fmt.Sprintf("行 %d", i+1)),
                Td(fmt.Sprintf("數據 %d", i+1)),
            )
        }),
    ),
)
```

### For

遍歷切片並渲染。

```go
func For[T any](items []T, fn func(T, int) VNode) VNode
```

#### 基本用法

```go
fruits := []string{"蘋果", "香蕉", "橘子"}

list := Ul(
    control.For(fruits, func(fruit string, i int) VNode {
        return Li(fmt.Sprintf("%d. %s", i+1, fruit))
    }),
)
```

#### 結構體切片

```go
type User struct {
    Name  string
    Email string
    Age   int
}

users := []User{
    {Name: "Alice", Email: "alice@example.com", Age: 25},
    {Name: "Bob", Email: "bob@example.com", Age: 30},
}

userList := Div(
    control.For(users, func(user User, i int) VNode {
        return Div(
            Props{"class": "user-card"},
            H3(user.Name),
            P(user.Email),
            Span(fmt.Sprintf("年齡: %d", user.Age)),
        )
    }),
)
```

#### 複雜數據結構

```go
type Product struct {
    ID    int
    Name  string
    Price float64
    Tags  []string
}

products := []Product{
    {ID: 1, Name: "筆記本", Price: 29.99, Tags: []string{"文具", "辦公"}},
    {ID: 2, Name: "鋼筆", Price: 15.50, Tags: []string{"文具", "書寫"}},
}

productGrid := Div(
    Props{"class": "row"},
    control.For(products, func(product Product, i int) VNode {
        return Div(
            Props{"class": "col-md-4"},
            Div(
                Props{"class": "card"},
                Div(
                    Props{"class": "card-body"},
                    H5(product.Name),
                    P(fmt.Sprintf("$%.2f", product.Price)),
                    Div(
                        control.For(product.Tags, func(tag string, j int) VNode {
                            return Span(
                                Props{"class": "badge bg-secondary me-1"},
                                tag,
                            )
                        }),
                    ),
                ),
            ),
        )
    }),
)
```

---

## JavaScript DSL 模塊

JavaScript DSL 模塊提供了類型安全的 JavaScript 代碼生成功能。

### JSAction

JSAction 是 JavaScript 代碼的載體。

```go
type JSAction struct {
    Code string
}
```

### 基本操作

#### Log (控制台日誌)

```go
func Log(msg string) JSAction
```

```go
js.Log("'Hello, World!'")
// 生成: console.log('Hello, World!')

js.Log("myVariable")
// 生成: console.log(myVariable)

js.Log("'User:', user")
// 生成: console.log('User:', user)
```

#### Alert (警告框)

```go
func Alert(jsExpr string) JSAction
```

```go
js.Alert("'歡迎！'")
// 生成: alert('歡迎！')

js.Alert("user.name")
// 生成: alert(user.name)
```

#### Redirect (頁面重定向)

```go
func Redirect(url string) JSAction
```

```go
js.Redirect("/home")
// 生成: location.href = '/home'
```

### 變數定義

#### Let (可變變數)

```go
func Let(varName string, value string) JSAction
```

```go
js.Let("counter", "0")
// 生成: let counter = 0

js.Let("name", "'Alice'")
// 生成: let name = 'Alice'

js.Let("data", "{ id: 1, name: 'Test' }")
// 生成: let data = { id: 1, name: 'Test' }
```

#### Const (常量)

```go
func Const(varName string, value string) JSAction
```

```go
js.Const("API_URL", "'https://api.example.com'")
// 生成: const API_URL = 'https://api.example.com'

js.Const("user", "await fetchUser()")
// 生成: const user = await fetchUser()
```

### DOM 操作

#### El (選擇單個元素)

```go
func El(selector string) Elem
```

```go
button := js.El("#myButton")
form := js.El("form.login")
title := js.El("h1:first-child")
```

#### Els (選擇多個元素)

```go
func Els(selector string) ElemList
```

```go
buttons := js.Els(".btn")
items := js.Els("li.item")
```

#### Elem 方法

```go
type Elem struct {
    Selector string
    VarName  string
}
```

##### SetText (設置文本)

```go
js.El("#title").SetText("'新標題'")
// 生成: document.querySelector('#title').innerText = '新標題'
```

##### SetHTML (設置 HTML)

```go
js.El("#content").SetHTML("'<strong>粗體</strong>'")
// 生成: document.querySelector('#content').innerHTML = '<strong>粗體</strong>'
```

##### AddClass (添加 class)

```go
js.El("#box").AddClass("active")
// 生成: document.querySelector('#box').classList.add('active')
```

##### RemoveClass (移除 class)

```go
js.El("#box").RemoveClass("hidden")
// 生成: document.querySelector('#box').classList.remove('hidden')
```

##### OnClick (點擊事件)

```go
js.El("#button").OnClick(
    js.Alert("'點擊了！'"),
)
// 生成: document.querySelector('#button').addEventListener('click', function() {
//   alert('點擊了！');
// });
```

##### InnerText / InnerHTML (訪問屬性)

```go
text := js.El("#input").InnerText()
// 返回字符串: "document.querySelector('#input').innerText"

html := js.El("#content").InnerHTML()
// 返回字符串: "document.querySelector('#content').innerHTML"
```

### 函數定義

#### Fn (定義函數)

```go
func Fn(params []string, actions ...JSAction) JSAction
```

```go
// 無參數函數
myFunc := js.Fn(nil,
    js.Log("'執行中'"),
    js.Alert("'完成'"),
)
// 生成: () => {
//   console.log('執行中');
//   alert('完成');
// }

// 有參數函數
greet := js.Fn([]string{"name"},
    js.Log("'Hello, ' + name"),
)
// 生成: (name) => {
//   console.log('Hello, ' + name);
// }

// 多參數函數
add := js.Fn([]string{"a", "b"},
    JSAction{Code: "return a + b"},
)
// 生成: (a, b) => {
//   return a + b;
// }
```

#### Call (調用函數)

```go
func Call(fnExpr any, args ...any) JSAction
```

```go
js.Call("myFunction", "'arg1'", "arg2")
// 生成: myFunction('arg1', arg2)

js.Call("fetch", "'/api/data'")
// 生成: fetch('/api/data')
```

#### CallMethod (調用對象方法)

```go
func CallMethod(objExpr string, methodName string, args ...any) JSAction
```

```go
js.CallMethod("evt", "preventDefault")
// 生成: evt.preventDefault()

js.CallMethod("arr", "push", "1", "2", "3")
// 生成: arr.push(1, 2, 3)
```

### DomReady (DOM 就緒)

```go
func DomReady(actions ...JSAction) JSAction
```

```go
script := Script(Props{"type": "text/javascript"},
    js.DomReady(
        js.El("#button").OnClick(
            js.Alert("'點擊'"),
        ),
        js.Log("'DOM 已就緒'"),
    ),
)
// 生成: document.addEventListener("DOMContentLoaded", () => {
//   document.querySelector('#button').addEventListener('click', function() {
//     alert('點擊');
//   });
//   console.log('DOM 已就緒');
// });
```

### TryCatch (異步錯誤處理)

```go
func TryCatch(baseAction JSAction, catchFn *JSAction, finallyFn *JSAction) JSAction
```

#### 基本用法

```go
js.TryCatch(
    js.Fn(nil,
        js.Const("data", "await fetch('/api/data')"),
        js.Log("'成功:', data"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'錯誤:', e"),
        js.Alert("'操作失敗'"),
    )),
    nil,
)
// 生成: (async () => { try {
//   const data = await fetch('/api/data');
//   console.log('成功:', data);
// } catch (e) {
//   console.log('錯誤:', e);
//   alert('操作失敗');
// } })()
```

#### 帶 finally

```go
js.TryCatch(
    js.Fn(nil,
        js.Log("'開始處理'"),
        js.Const("result", "await process()"),
    ),
    js.Ptr(js.Fn(nil,
        js.Log("'錯誤:', e"),
    )),
    js.Ptr(js.Fn(nil,
        js.Log("'清理資源'"),
    )),
)
```

### 高級功能

#### CreateEl (創建元素)

```go
func CreateEl(tagName string, varName ...string) (Elem, JSAction)
```

```go
div, createAction := js.CreateEl("div", "myDiv")
// createAction: const myDiv = document.createElement('div');

// 使用創建的元素
actions := []JSAction{
    createAction,
    div.SetText("'內容'"),
    div.AddClass("box"),
}
```

#### AppendChild (添加子元素)

```go
parent.AppendChild(child)
```

```go
container, createContainer := js.CreateEl("div", "container")
item, createItem := js.CreateEl("span", "item")

actions := []JSAction{
    createContainer,
    createItem,
    item.SetText("'文本'"),
    container.AppendChild(item),
    js.El("#root").AppendChild(container),
}
```

#### SetTimeout / SetInterval

```go
func SetTimeout(action JSAction, delayMs int) JSAction
func SetInterval(action JSAction, intervalMs int) JSAction
```

```go
// 延遲執行
js.SetTimeout(
    js.Alert("'3 秒後顯示'"),
    3000,
)

// 定時執行
js.SetInterval(
    js.Log("'每秒執行'"),
    1000,
)
```

### Fetch API 示例

雖然 jsdsl 模塊有 FetchRequest 等函數，但推薦使用 DSL 方式構建 Fetch 請求：

#### GET 請求

```go
Button(Props{
    "onClick": js.Fn(nil,
        js.Log("'開始獲取'"),
        js.TryCatch(
            js.Fn(nil,
                js.Const("response", "await fetch('/api/data')"),
                JSAction{Code: "if (!response.ok) throw new Error('HTTP ' + response.status)"},
                js.Const("data", "await response.json()"),
                js.Log("'數據:', data"),
            ),
            js.Ptr(js.Fn(nil,
                js.Log("'錯誤:', e"),
                js.Alert("'獲取失敗: ' + e.message"),
            )),
            nil,
        ),
    ),
}, "獲取數據")
```

#### POST 請求

```go
Form(Props{
    "onSubmit": js.Fn([]string{"evt"},
        js.CallMethod("evt", "preventDefault"),
        js.TryCatch(
            js.Fn(nil,
                js.Const("formData", "{ name: document.getElementById('name').value }"),
                js.Const("response", "await fetch('/api/submit', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(formData) })"),
                JSAction{Code: "if (!response.ok) throw new Error('提交失敗')"},
                js.Const("result", "await response.json()"),
                js.Alert("'提交成功'"),
            ),
            js.Ptr(js.Fn(nil,
                js.Alert("'提交失敗: ' + e.message"),
            )),
            nil,
        ),
    ),
},
    Input(Props{"id": "name", "type": "text"}),
    Button(Props{"type": "submit"}, "提交"),
)
```

---

## Components 模塊

Components 模塊提供了常用的表單組件。

### TextField (文字輸入框)

```go
func TextField(props Props) VNode
```

#### 支持的 Props

- `id`: 輸入框 ID（必填）
- `name`: 表單名稱（默認同 id）
- `label`: 標籤文本
- `placeholder`: 占位符
- `type`: 輸入類型（默認 "text"）
- `value`: 默認值
- `required`: 是否必填（"true" / "false"）
- `disabled`: 是否禁用（"true" / "false"）
- `helpText`: 幫助文本
- `class`: 額外的 CSS 類

#### 示例

```go
comp.TextField(Props{
    "id":          "email",
    "label":       "電子郵件",
    "type":        "email",
    "placeholder": "請輸入電子郵件",
    "required":    "true",
    "helpText":    "我們不會分享您的郵件地址",
})
```

### Dropdown (下拉選單)

```go
func Dropdown(props Props) VNode
```

#### 支持的 Props

- `id`: 選單 ID（必填）
- `name`: 表單名稱（默認同 id）
- `label`: 標籤文本
- `options`: 選項，逗號分隔（例如："選項1,選項2,選項3"）
- `value`: 默認選中值
- `required`: 是否必填
- `disabled`: 是否禁用
- `helpText`: 幫助文本

#### 示例

```go
comp.Dropdown(Props{
    "id":       "country",
    "label":    "選擇國家",
    "options":  "台灣,中國,日本,美國",
    "value":    "台灣",
    "required": "true",
    "helpText": "請選擇您的所在國家",
})
```

### RadioGroup (單選按鈕組)

```go
func RadioGroup(props Props) VNode
```

#### 支持的 Props

- `id`: 組 ID（必填）
- `name`: 表單名稱（必填）
- `label`: 組標籤
- `options`: 選項，逗號分隔
- `defaultValue`: 默認選中值
- `direction`: 佈局方向（"horizontal" / "vertical"，默認 "vertical"）
- `required`: 是否必填
- `disabled`: 是否禁用
- `helpText`: 幫助文本

#### 示例

```go
comp.RadioGroup(Props{
    "id":           "gender",
    "name":         "gender",
    "label":        "性別",
    "options":      "男性,女性,其他",
    "defaultValue": "男性",
    "direction":    "horizontal",
    "required":     "true",
})
```

### Checkbox (勾選框)

```go
func Checkbox(props Props) VNode
```

#### 支持的 Props

- `id`: 勾選框 ID（必填）
- `name`: 表單名稱
- `label`: 標籤文本
- `value`: 勾選框值
- `checked`: 是否預選（"true" / "false"）
- `required`: 是否必填
- `disabled`: 是否禁用
- `helpText`: 幫助文本

#### 示例

```go
comp.Checkbox(Props{
    "id":       "terms",
    "name":     "terms",
    "label":    "我同意服務條款和隱私政策",
    "required": "true",
    "checked":  "false",
    "helpText": "您必須同意條款才能繼續",
})
```

### CheckboxGroup (勾選框組)

```go
func CheckboxGroup(props Props) VNode
```

#### 支持的 Props

- `id`: 組 ID（必填）
- `name`: 表單名稱
- `label`: 組標籤
- `options`: 選項，逗號分隔
- `values`: 預選值，逗號分隔
- `required`: 是否必填
- `disabled`: 是否禁用
- `helpText`: 幫助文本

#### 示例

```go
comp.CheckboxGroup(Props{
    "id":      "hobbies",
    "name":    "hobbies",
    "label":   "選擇愛好",
    "options": "閱讀,運動,音樂,繪畫,旅行",
    "values":  "閱讀,音樂",
    "helpText": "可多選",
})
```

### Switch (開關)

```go
func Switch(props Props) VNode
```

#### 支持的 Props

- `id`: 開關 ID（必填）
- `name`: 表單名稱
- `label`: 標籤文本
- `checked`: 是否開啟（"true" / "false"）
- `disabled`: 是否禁用
- `labelPosition`: 標籤位置（"left" / "right"，默認 "left"）
- `helpText`: 幫助文本

#### 示例

```go
comp.Switch(Props{
    "id":            "notifications",
    "name":          "notifications",
    "label":         "啟用電子郵件通知",
    "checked":       "true",
    "labelPosition": "right",
    "helpText":      "開啟以接收重要通知",
})
```

---

## 進階用法

### 自定義組件庫

創建您自己的組件庫：

```go
package mycomponents

import (
    . "github.com/TimLai666/go-vdom/vdom"
    js "github.com/TimLai666/go-vdom/jsdsl"
)

// Alert 組件
var Alert = Component(
    Div(
        Props{
            "class": "alert alert-{{type}} {{className}}",
            "role":  "alert",
        },
        Strong("{{title}}"),
        Span(" {{message}}"),
        Button(
            Props{
                "type":  "button",
                "class": "btn-close",
                "data-bs-dismiss": "alert",
            },
        ),
    ),
    nil,
    PropsDef{
        "type":      "info",  // info, success, warning, danger
        "title":     "",
        "message":   "",
        "className": "",
    },
)

// Modal 組件
var Modal = Component(
    Div(
        Props{
            "class": "modal fade",
            "id":    "{{id}}",
            "tabindex": "-1",
        },
        Div(
            Props{"class": "modal-dialog"},
            Div(
                Props{"class": "modal-content"},
                Div(
                    Props{"class": "modal-header"},
                    H5(Props{"class": "modal-title"}, "{{title}}"),
                    Button(
                        Props{
                            "type":  "button",
                            "class": "btn-close",
                            "data-bs-dismiss": "modal",
                        },
                    ),
                ),
                Div(
                    Props{"class": "modal-body"},
                    "{{children}}",
                ),
                Div(
                    Props{"class": "modal-footer"},
                    Button(
                        Props{
                            "type":  "button",
                            "class": "btn btn-secondary",
                            "data-bs-dismiss": "modal",
                        },
                        "關閉",
                    ),
                    Button(
                        Props{
                            "type":  "button",
                            "class": "btn btn-primary",
                        },
                        "確定",
                    ),
                ),
            ),
        ),
    ),
    nil,
    PropsDef{
        "id":    "myModal",
        "title": "Modal",
    },
)
```

### 高階組件模式

```go
// WithLoading - 添加載入狀態的高階組件
func WithLoading(component func(Props, ...VNode) VNode) func(Props, ...VNode) VNode {
    return func(props Props, children ...VNode) VNode {
        loading := props["loading"] == "true"
        
        if loading {
            return Div(
                Props{"class": "loading-wrapper"},
                Div(
                    Props{"class": "spinner-border", "role": "status"},
                    Span(Props{"class": "visually-hidden"}, "載入中..."),
                ),
            )
        }
        
        return component(props, children...)
    }
}

// 使用
EnhancedCard := WithLoading(Card)

instance := EnhancedCard(
    Props{
        "title":   "我的卡片",
        "loading": "false",
    },
    P("內容"),
)
```

### 組合組件

```go
// UserProfile 組合多個組件
func UserProfile(user User) VNode {
    return Div(
        Props{"class": "user-profile"},
        
        // 頭像卡片
        Card(Props{
            "title": "個人信息",
        },
            Img(Props{
                "src":   user.Avatar,
                "class": "avatar",
                "alt":   user.Name,
            }),
            H4(user.Name),
            P(user.Email),
        ),
        
        // 編輯表單
        Card(Props{
            "title": "編輯資料",
        },
            Form(
                comp.TextField(Props{
                    "id":    "name",
                    "label": "姓名",
                    "value": user.Name,
                }),
                comp.TextField(Props{
                    "id":    "email",
                    "label": "電子郵件",
                    "type":  "email",
                    "value": user.Email,
                }),
                comp.Dropdown(Props{
                    "id":      "country",
                    "label":   "國家",
                    "options": "台灣,中國,日本",
                    "value":   user.Country,
                }),
                Button(
                    Props{"type": "submit", "class": "btn btn-primary"},
                    "保存",
                ),
            ),
        ),
    )
}
```

### 條件樣式

```go
func StatusBadge(status string) VNode {
    var badgeClass string
    var text string
    
    switch status {
    case "active":
        badgeClass = "bg-success"
        text = "活躍"
    case "pending":
        badgeClass = "bg-warning"
        text = "待處理"
    case "inactive":
        badgeClass = "bg-secondary"
        text = "未活躍"
    default:
        badgeClass = "bg-info"
        text = "未知"
    }
    
    return Span(
        Props{"class": fmt.Sprintf("badge %s", badgeClass)},
        text,
    )
}
```

### 動態 Props

```go
func buildProps(base Props, conditional map[string]bool) Props {
    result := Props{}
    
    // 複製基礎 props
    for k, v := range base {
        result[k] = v
    }
    
    // 根據條件添加
    for k, v := range conditional {
        if v {
            result[k] = "true"
        }
    }
    
    return result
}

// 使用
props := buildProps(
    Props{"class": "btn", "type": "button"},
    map[string]bool{
        "disabled": isDisabled,
        "required": isRequired,
    },
)

button := Button(props, "提交")
```

---

## 性能優化

### 1. 組件複用

```go
// ✅ 好的做法 - 定義一次，多次使用
var UserCard = Component(...)

for _, user := range users {
    cards = append(cards, UserCard(Props{...}))
}

// ❌ 不好的做法 - 每次都定義
for _, user := range users {
    card := Component(...)  // 重複定義
    cards = append(cards, card(Props{...}))
}
```

### 2. Props 預分配

```go
// ✅ 好的做法
props := make(Props, 10)  // 預分配容量
props["id"] = "myId"
props["class"] = "container"
// ...

// ❌ 不好的做法
props := Props{}  // 可能需要多次重新分配
props["id"] = "myId"
props["class"] = "container"
// ...
```

### 3. 字符串構建

```go
// ✅ 好的做法 - 使用 strings.Builder
var sb strings.Builder
for _, item := range items {
    sb.WriteString(Render(Li(item)))
}
html := sb.String()

// ❌ 不好的做法 - 字符串拼接
var html string
for _, item := range items {
    html += Render(Li(item))  // 每次都分配新字符串
}
```

### 4. 避免深度嵌套

```go
// ✅ 好的做法 - 分解組件
Header := Component(...)
Content := Component(...)
Footer := Component(...)

page := Div(
    Header(Props{}),
    Content(Props{}),
    Footer(Props{}),
)

// ❌ 不好的做法 - 深度嵌套
page := Div(
    Div(
        Div(
            Div(
                Div(
                    Div("內容"),
                ),
            ),
        ),
    ),
)
```

### 5. 批量渲染

```go
// ✅ 好的做法 - 一次渲染
items := control.For(data, func(item Data, i int) VNode {
    return Li(item.Name)
})
html := Render(Ul(items))

// ❌ 不好的做法 - 多次渲染
var htmlParts []string
for _, item := range data {
    htmlParts = append(htmlParts, Render(Li(item.Name)))
}
html := strings.Join(htmlParts, "")
```

---

## 故障排除

### 常見問題

#### 1. Props 未生效

```go
// ❌ 錯誤：Props 不是第一個參數
Div("文本", Props{"class": "container"})

// ✅ 正確：Props 必須是第一個參數
Div(Props{"class": "container"}, "文本")
```

#### 2. 組件 Props 未替換

```go
// ❌ 錯誤：忘記定義 PropsDef
MyComponent := Component(
    Div("{{title}}"),
    nil,
    PropsDef{},  // 空的 PropsDef
)

// ✅ 正確：定義所有使用的 props
MyComponent := Component(
    Div("{{title}}"),
    nil,
    PropsDef{"title": "默認標題"},
)
```

#### 3. JavaScript 事件不觸發

```go
// ❌ 錯誤：忘記調用 Fn
Button(Props{
    "onClick": js.Alert("'Hi'"),  // 直接使用 JSAction
}, "按鈕")

// ✅ 正確：包裝在 Fn 中
Button(Props{
    "onClick": js.Fn(nil, js.Alert("'Hi'")),
}, "按鈕")
```

#### 4. TryCatch 錯誤

```go
// ❌ 錯誤：沒有使用 Ptr
js.TryCatch(
    js.Fn(...),
    js.Fn(...),  // 應該是 *JSAction
    nil,
)

// ✅ 正確：使用 Ptr
js.TryCatch(
    js.Fn(...),
    js.Ptr(js.Fn(...)),
    nil,
)
```

#### 5. 字符串轉義問題

```go
// ❌ 錯誤：JavaScript 字符串沒有引號
js.Log("Hello")  // 生成: console.log(Hello) - 錯誤！

// ✅ 正確：添加引號
js.Log("'Hello'")  // 生成: console.log('Hello')
```

### 調試技巧

#### 1. 檢查生成的 HTML

```go
vnode := Div(Props{"class": "test"}, "內容")
html := Render(vnode)
fmt.Println(html)  // 輸出生成的 HTML
```

#### 2. 檢查 JavaScript 代碼

```go
action := js.Fn(nil,
    js.Log("'Test'"),
    js.Alert("'Hi'"),
)
fmt.Println(action.Code)  // 輸出生成的 JavaScript
```

#### 3. 分步構建

```go
// 分步構建，便於調試
header := H1("標題")
content := P("內容")
footer := Footer("頁腳")

page := Div(
    Props{"class": "page"},
    header,
    content,
    footer,
)

// 檢查每個部分
fmt.Println(Render(header))
fmt.Println(Render(content))
fmt.Println(Render(footer))
fmt.Println(Render(page))
```

#### 4. 使用瀏覽器開發工具

生成 HTML 後，在瀏覽器中：
- 查看元素（Elements/Inspector）
- 查看控制台（Console）
- 查看網絡請求（Network）
- 使用 Source Maps（如果有）

---

## 總結

go-vdom 提供了一套完整的工具鏈，讓你能夠在 Go 中以類型安全的方式構建動態網頁。通過合理使用組件化、控制流和 JavaScript DSL，你可以創建出維護性高、性能優的 Web 應用。

### 關鍵要點

1. **使用 DSL** - 盡量使用 DSL 而非原始字符串
2. **組件化** - 將重複的 UI 邏輯封裝成組件
3. **類型安全** - 利用 Go 的類型系統避免運行時錯誤
4. **錯誤處理** - 在 JavaScript 代碼中使用 TryCatch
5. **性能優化** - 注意組件複用和批量渲染

### 下一步

- 查看 `main.go` 中的完整示例
- 探索 `components` 包中的組件實現
- 嘗試創建自己的組件庫
- 與現有的 Go web 框架集成

---

**文檔版本**: 1.0.0  
**最後更新**: 2025-01-24  
**作者**: TimLai666