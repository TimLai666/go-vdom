# Go DOM JavaScript DSL API Reference

完整的 JavaScript DSL API 參考文檔，包含所有函數、類型和使用示例。

## 目錄

- [核心函數](#核心函數)
- [DOM 操作](#dom-操作)
- [事件處理](#事件處理)
- [異步操作](#異步操作)
- [錯誤處理](#錯誤處理)
- [實用工具](#實用工具)

---

## 核心函數

### `Fn(params []string, actions ...JSAction) JSAction`

創建一個同步 JavaScript 函數。

**參數：**

- `params`: 函數參數列表，如不需要參數可傳入 `nil`
- `actions`: 函數體內的動作序列

**返回：** `JSAction`

**示例：**

```go
// 無參數函數
js.Fn(nil,
    js.Log("'Hello World'"),
)

// 帶參數函數
js.Fn([]string{"event"},
    js.CallMethod("event", "preventDefault"),
    js.Log("'Button clicked'"),
)
```

**生成的 JavaScript：**

```javascript
() => {
  console.log("Hello World");
};

(event) => {
  event.preventDefault();
  console.log("Button clicked");
};
```

---

### `AsyncFn(params []string, actions ...JSAction) JSAction`

創建一個異步 JavaScript 函數，支持 `await` 語法。

**參數：**

- `params`: 函數參數列表，如不需要參數可傳入 `nil`
- `actions`: 函數體內的動作序列（可包含 `await` 語句）

**返回：** `JSAction`

**示例：**

```go
// 異步 API 調用
js.AsyncFn(nil,
    js.Const("response", "await fetch('/api/data')"),
    js.Const("data", "await response.json()"),
    js.Log("data"),
)

// 帶參數的異步函數
js.AsyncFn([]string{"url"},
    js.Const("response", "await fetch(url)"),
    JSAction{Code: "return await response.json()"},
)
```

**生成的 JavaScript：**

```javascript
async () => {
  const response = await fetch("/api/data");
  const data = await response.json();
  console.log(data);
};

async (url) => {
  const response = await fetch(url);
  return await response.json();
};
```

---

### `Call(fnExpr any, args ...any) JSAction`

調用一個函數。

**參數：**

- `fnExpr`: 函數表達式（可以是 `string` 或 `JSAction`）
- `args`: 函數參數

**返回：** `JSAction`

**示例：**

```go
js.Call("myFunction", "'arg1'", "123")
js.Call("window.alert", "'Hello'")
```

---

### `CallMethod(objExpr string, methodName string, args ...any) JSAction`

調用對象的方法。

**參數：**

- `objExpr`: 對象表達式
- `methodName`: 方法名稱
- `args`: 方法參數

**返回：** `JSAction`

**示例：**

```go
js.CallMethod("document", "getElementById", "'myId'")
js.CallMethod("console", "log", "'message'", "data")
```

---

## DOM 操作

### `El(selector string) Elem`

通過 CSS 選擇器獲取單個元素。

**參數：**

- `selector`: CSS 選擇器字符串

**返回：** `Elem` 對象

**示例：**

```go
element := js.El("#myButton")
element.SetText("'Click me'")
element.AddClass("'active'")
```

---

### `Els(selector string) ElemList`

通過 CSS 選擇器獲取多個元素。

**參數：**

- `selector`: CSS 選擇器字符串

**返回：** `ElemList` 對象

**示例：**

```go
elements := js.Els(".item")
js.QueryEach(elements, func(el js.Elem) JSAction {
    return el.AddClass("'highlighted'")
})
```

---

### `CreateEl(tagName string, varName ...string) (Elem, JSAction)`

動態創建 DOM 元素。

**參數：**

- `tagName`: HTML 標籤名稱
- `varName`: 可選的變數名稱

**返回：**

- `Elem`: 元素對象
- `JSAction`: 創建元素的動作

**示例：**

```go
div, createAction := js.CreateEl("div", "myDiv")
js.NewJSActionBuilder().
    Add(createAction).
    SetElementText(div, "'Hello'").
    AddElementClass(div, "'card'").
    Build()
```

---

### Elem 方法

#### `SetText(text string) JSAction`

設置元素的 innerText。

```go
js.El("#title").SetText("'新標題'")
```

#### `SetHTML(html string) JSAction`

設置元素的 innerHTML。

```go
js.El("#content").SetHTML("'<p>HTML 內容</p>'")
```

#### `AddClass(class string) JSAction`

添加 CSS 類別。

```go
js.El("#element").AddClass("'active'")
```

#### `RemoveClass(class string) JSAction`

移除 CSS 類別。

```go
js.El("#element").RemoveClass("'hidden'")
```

#### `AppendChild(child Elem) JSAction`

添加子元素。

```go
parent := js.El("#parent")
child := js.El("#child")
parent.AppendChild(child)
```

#### `InnerText() string`

獲取元素的 innerText 表達式。

```go
text := js.El("#input").InnerText()
js.Log(text)
```

#### `InnerHTML() string`

獲取元素的 innerHTML 表達式。

```go
html := js.El("#container").InnerHTML()
js.Const("content", html)
```

---

### `DomReady(actions ...JSAction) JSAction`

在 DOM 加載完成後執行代碼。

**參數：**

- `actions`: 要執行的動作序列

**返回：** `JSAction`

**示例：**

```go
js.DomReady(
    js.Log("'DOM is ready'"),
    js.El("#app").SetText("'Application loaded'"),
)
```

**生成的 JavaScript：**

```javascript
document.addEventListener("DOMContentLoaded", () => {
  console.log("DOM is ready");
  document.querySelector("#app").innerText = "Application loaded";
});
```

---

## 事件處理

### `OnClick(action JSAction) JSAction`

為元素添加點擊事件監聽器。

**示例：**

```go
js.El("#button").OnClick(
    js.Alert("'Button clicked!'"),
)
```

---

## 異步操作

### Try-Catch-Finally（流暢 API）

請參閱 [Try-Catch-Finally 指南](TRY_CATCH_FINALLY.md) 獲取完整說明。

**基本用法：**

```go
// Try-Catch
js.Try(
    js.Const("x", "parseInt('abc')"),
    js.Log("x"),
).Catch("error",
    js.Log("'錯誤: ' + error.message"),
).End()

// Try-Catch-Finally
js.AsyncFn(nil,
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch("error",
        js.Log("'錯誤: ' + error.message"),
    ).Finally(
        js.Log("'清理完成'"),
    ),
)

// 在 AsyncDo 中使用
js.AsyncDo(
    js.Try(
        js.Const("data", "await fetch('/api')"),
    ).Catch("error",
        js.Log("'錯誤: ' + error.message"),
    ).End(),
)
```

**特點：**

- ✅ 生成純粹的 try-catch-finally 語句（不包裝在函數中）
- ✅ 錯誤對象統一命名為 `error`
- ✅ 支持 Try-Catch、Try-Finally、Try-Catch-Finally
- ✅ 流暢的鏈式調用 API
- ✅ 需要 async 時，用 AsyncFn 或 AsyncDo 包裝

---

### Fetch API 輔助函數

#### `FetchRequest(url string, options ...FetchOption) JSAction`

創建一個 fetch 請求。

**示例：**

```go
js.FetchRequest("/api/data",
    js.WithMethod("POST"),
    js.WithContentType("application/json"),
    js.WithBody("JSON.stringify({name: 'test'})"),
)
```

#### Fetch 選項

- `WithMethod(method string) FetchOption` - 設置 HTTP 方法
- `WithHeader(name, value string) FetchOption` - 設置請求頭
- `WithBody(body string) FetchOption` - 設置請求體
- `WithContentType(contentType string) FetchOption` - 設置 Content-Type
- `WithJSON(jsonObject string) []FetchOption` - 設置 JSON 請求
- `WithFormData(formData map[string]string) []FetchOption` - 設置表單數據

---

## 錯誤處理

### `Alert(jsExpr string) JSAction`

顯示警告對話框。

```go
js.Alert("'Hello World'")
js.Alert("data.message")
```

### `Log(msg string) JSAction`

輸出到控制台。

```go
js.Log("'Debug message'")
js.Log("JSON.stringify(data)")
```

---

## 實用工具

### `Let(varName string, value string) JSAction`

聲明一個可變變數。

```go
js.Let("counter", "0")
```

### `Const(varName string, value string) JSAction`

聲明一個常量。

```go
js.Const("apiUrl", "'/api/data'")
js.Const("response", "await fetch(apiUrl)")
```

### `V(varName string) variable`

創建一個變數引用對象。

```go
v := js.V("myVar")
v.SetHTML("'<p>Content</p>'")
v.AddClass("'active'")
```

### `VRef(varName string) JSAction`

直接引用變數名。

```go
js.VRef("myVariable")
```

### `Ptr(a JSAction) *JSAction`

將 JSAction 轉換為指針（用於需要可選參數的函數）。

```go
js.Ptr(js.Fn(nil, js.Log("'error'")))
```

### `Redirect(url string) JSAction`

重定向到指定 URL。

```go
js.Redirect("/home")
```

### `SetTimeout(action JSAction, delayMs int) JSAction`

延遲執行代碼。

```go
js.SetTimeout(
    js.Log("'Delayed message'"),
    1000,
)
```

### `SetInterval(action JSAction, intervalMs int) JSAction`

定期執行代碼。

```go
js.SetInterval(
    js.Log("'Periodic message'"),
    5000,
)
```

---

## 迭代和循環

### `ForEach(arrayExpr string, itemVar string, actions ...JSAction) JSAction`

遍歷任意 JavaScript 數組或可迭代對象（前端渲染）。

**參數：**

- `arrayExpr`: 數組表達式（如 `"myArray"`, `"[1,2,3]"`, `"data.items"`）
- `itemVar`: 項目變數名稱（如 `"item"`, `"user"`）
- `actions`: 對每個項目執行的動作

**示例：**

```go
// 遍歷數組並輸出每個項目
js.ForEachJS("['Apple', 'Banana', 'Orange']", "fruit",
    js.Log("'水果: ' + fruit"),
)

// 遍歷 API 返回的數據
js.ForEachJS("data.items", "item",
    js.Log("'ID: ' + item.id"),
    js.Log("'Name: ' + item.name"),
)

// 動態創建 DOM 元素
js.Const("container", "document.getElementById('list')"),
js.ForEachJS("colors", "color",
    js.Const("div", "document.createElement('div')"),
    JSAction{Code: "div.textContent = color"},
    JSAction{Code: "container.appendChild(div)"},
)
```

**生成的 JavaScript：**

```javascript
["Apple", "Banana", "Orange"].forEach(function (fruit) {
  console.log("水果: " + fruit);
});
```

---

### `ForEachWithIndex(arrayExpr string, itemVar string, indexVar string, actions ...JSAction) JSAction`

遍歷數組並提供索引。

**參數：**

- `arrayExpr`: 數組表達式
- `itemVar`: 項目變數名稱
- `indexVar`: 索引變數名稱
- `actions`: 對每個項目執行的動作

**示例：**

```go
js.ForEachWithIndexJS("items", "item", "index",
    js.Log("'[' + index + '] = ' + item"),
)

// 創建編號列表
js.ForEachWithIndexJS("names", "name", "i",
    js.Const("li", "document.createElement('li')"),
    JSAction{Code: "li.textContent = (i + 1) + '. ' + name"},
    JSAction{Code: "list.appendChild(li)"},
)
```

**生成的 JavaScript：**

```javascript
items.forEach(function (item, index) {
  console.log("[" + index + "] = " + item);
});
```

---

### `ForEachElement(arrayExpr string, fn func(el Elem) JSAction) JSAction`

遍歷 DOM 元素列表（專門用於 DOM 操作）。

**示例：**

```go
js.ForEachElement("document.querySelectorAll('.item')", func(el js.Elem) JSAction {
    return el.AddClass("'highlighted'")
})
```

---

### `QueryEach(list ElemList, fn func(el Elem) JSAction) JSAction`

遍歷元素列表。

```go
elements := js.Els(".item")
js.QueryEach(elements, func(el js.Elem) JSAction {
    return el.AddClass("'processed'")
})
```

---

## 進階用法

### JSActionBuilder

用於構建複雜的 JavaScript 動作序列。

```go
builder := js.NewJSActionBuilder()

// 創建元素
div := builder.CreateElement("div", "myDiv")

// 鏈式操作
builder.
    SetElementText(div, "'Hello'").
    AddElementClass(div, "'card'").
    SetElementHTML(div, "'<strong>Bold</strong>'")

// 構建最終動作
action := builder.Build()
```

#### 方法

- `Add(action JSAction)` - 添加單個動作
- `AddMany(actions ...JSAction)` - 添加多個動作
- `CreateElement(tagName, varName)` - 創建元素
- `SetElementText(elem, text)` - 設置元素文本
- `SetElementHTML(elem, html)` - 設置元素 HTML
- `AddElementClass(elem, class)` - 添加類別
- `RemoveElementClass(elem, class)` - 移除類別
- `AppendChild(parent, child)` - 添加子元素
- `GetActions()` - 獲取所有動作
- `Build()` - 構建 DomReady 包裝的最終動作

---

## 完整示例

### 異步表單提交

```go
Props{
    "onSubmit": js.AsyncFn([]string{"event"},
        js.CallMethod("event", "preventDefault"),
        js.Try(
            js.Const("formData", "new FormData(event.target)"),
            js.Const("response", "await fetch('/api/submit', { method: 'POST', body: formData })"),
            JSAction{Code: "if (!response.ok) throw new Error('提交失敗')"},
            js.Const("result", "await response.json()"),
            js.Alert("'提交成功: ' + result.message"),
        ).Catch("error",
            js.Alert("'提交失敗: ' + error.message"),
        ).End(),
    ),
}
```

### 動態內容加載

```go
Props{
    "onClick": js.AsyncFn(nil,
        js.Const("container", "document.getElementById('content')"),
        JSAction{Code: "container.innerHTML = '載入中...'"},
        js.Try(
            js.Const("response", "await fetch('/api/content')"),
            js.Const("html", "await response.text()"),
            JSAction{Code: "container.innerHTML = html"},
        ).Catch("error",
            JSAction{Code: "container.innerHTML = '載入失敗'"},
        ).End(),
    ),
}
```

### 複雜 DOM 操作

```go
js.DomReady(
    js.Const("items", "['項目1', '項目2', '項目3']"),
    js.Const("list", "document.createElement('ul')"),
    js.ForEachJS("items", "item",
        js.Const("li", "document.createElement('li')"),
        JSAction{Code: "li.textContent = item"},
        JSAction{Code: "list.appendChild(li)"},
    ),
    js.CallMethod("document.body", "appendChild", "list"),
)
```

### 從 API 獲取並渲染列表

```go
Props{
    "onClick": js.AsyncFn(nil,
        js.Const("container", "document.getElementById('list')"),
        js.Try(
            js.Const("response", "await fetch('/api/users')"),
            js.Const("users", "await response.json()"),
            JSAction{Code: "container.innerHTML = ''"},
            js.ForEachJS("users", "user",
                js.Const("div", "document.createElement('div')"),
                JSAction{Code: "div.className = 'user-card'"},
                JSAction{Code: "div.innerHTML = '<h3>' + user.name + '</h3><p>' + user.email + '</p>'"},
                js.CallMethod("container", "appendChild", js.Ident("div")),
            ),
        ).Catch("error",
            JSAction{Code: "container.innerHTML = '<p>載入失敗</p>'"},
        ).End(),
    ),
}
```

---

## 類型定義

### `JSAction`

```go
type JSAction struct {
    Code string
}
```

JavaScript 動作的基本類型，包含要執行的 JavaScript 代碼。

### `Elem`

```go
type Elem struct {
    Selector string
    VarName  string
}
```

代表一個 DOM 元素的引用。

### `ElemList`

```go
type ElemList struct {
    Selector string
}
```

代表多個 DOM 元素的引用。

### `FetchOption`

```go
type FetchOption struct {
    Key   string
    Value string
}
```

Fetch 請求的配置選項。

### `ResponseType`

```go
type ResponseType string

const (
    JSONResponse ResponseType = "json"
    TextResponse ResponseType = "text"
    BlobResponse ResponseType = "blob"
)
```

---

## 最佳實踐

1. **使用 AsyncFn 處理異步操作**
   - 任何包含 `await` 的代碼必須使用 `AsyncFn`
   - 外層點擊事件等也需要使用 `AsyncFn`

2. **錯誤處理**
   - 使用 `Try(...).Catch(...)` 包裝異步操作
   - Try 生成純粹的 try-catch-finally 語句
   - 需要 async 時，用 AsyncFn 或 AsyncDo 包裝
   - 錯誤對象統一命名為 `error`
   - 始終提供 catch 或 finally 處理器

3. **字符串引用**
   - JavaScript 字符串需要使用單引號：`"'text'"`
   - 變數引用不需要引號：`"variableName"`

4. **列表渲染選擇**
   - 前端動態數據使用 `js.ForEach` 或 `js.ForEachWithIndex`
   - DOM 元素操作使用 `js.ForEachElement`
   - 後端靜態數據使用 Go 的 `ForEach` 或 `ForEachWithIndex`（見 dom 包）

5. **代碼組織**
   - 複雜邏輯使用 `JSActionBuilder`
   - 將可重用的代碼提取為函數

6. **性能優化**
   - 使用 `DomReady` 確保 DOM 已加載
   - 批量 DOM 操作使用 DocumentFragment
   - 大列表優先考慮後端渲染

---

## 版本歷史

- **v1.1.0** - 新增 `AsyncFn` 支持異步函數
- **v1.0.0** - 初始發布

---

## 相關文檔

- [快速入門](QUICK_START.md)
- [完整文檔](DOCUMENTATION.md)
- [快速參考](QUICK_REFERENCE.md)
