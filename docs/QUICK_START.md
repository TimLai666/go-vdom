# Go VDOM v1.1.0 快速入門指南

本指南介紹 go-vdom v1.1.0 的新特性和使用方法。

---

## 🎯 新特性概覽

### v1.1.0 主要更新

1. **Props 類型系統** - 支持任意類型值
2. **模板序列化** - Go Template 和 JSON 支持
3. **完整 DSL** - 所有示例都使用 DSL

---

## 📦 安裝

```bash
go get github.com/TimLai666/go-vdom@v1.1.0
```

---

## 🚀 快速開始

### 1. Props 類型系統 - 新特性！

**不再需要手動轉換類型！**

#### 之前的方式
```go
import (
    "strconv"
    . "github.com/TimLai666/go-vdom/vdom"
)

// 需要手動轉換
Props{
    "width":    strconv.Itoa(800),
    "disabled": "true",
    "count":    fmt.Sprintf("%d", 42),
}
```

#### 現在的方式 ✨
```go
import . "github.com/TimLai666/go-vdom/vdom"

// 直接使用原始類型
Props{
    "width":    800,        // int - 自動轉換
    "disabled": true,       // bool - true 渲染，false 省略
    "count":    42,         // int - 自動轉換
    "opacity":  0.8,        // float64 - 自動轉換
}
```

#### 支持的類型

```go
Props{
    // 字符串
    "class": "container",
    
    // 布爾值（true 渲染屬性，false 省略）
    "disabled": true,      // 渲染為 disabled="true"
    "hidden":   false,     // 不渲染此屬性
    
    // 整數
    "width":    800,       // int
    "height":   600,       // int
    "tabindex": 0,         // int
    
    // 浮點數
    "opacity": 0.8,        // float64
    "price":   19.99,      // float64
    
    // JSAction（事件處理）
    "onClick": js.Fn(nil, js.Alert("'Hi'")),
}
```

---

### 2. 模板序列化 - 新特性！

**現在可以保存和重用模板了！**

#### 保存為 Go Template

```go
package main

import (
    "os"
    . "github.com/TimLai666/go-vdom/vdom"
)

func main() {
    // 1. 創建帶模板變數的 VNode
    userCard := Div(
        Props{
            "class": "user-card",
            "id":    "user-{{.ID}}",
        },
        H3("{{.Name}}"),
        P("Email: {{.Email}}"),
        P("Role: {{.Role}}"),
    )
    
    // 2. 保存為 Go Template
    template := SaveTemplate("user-card", userCard)
    os.WriteFile("user-card.tmpl", []byte(template), 0644)
    
    // 3. 保存為 JSON
    jsonStr, _ := ToJSON(userCard)
    os.WriteFile("user-card.json", []byte(jsonStr), 0644)
}
```

#### 生成的 Go Template

```html
{{/* Template: user-card */}}
{{define "user-card"}}
<div class="user-card" id="user-{{.ID}}">
  <h3>
    {{.Name}}
  </h3>
  <p>
    Email: {{.Email}}
  </p>
  <p>
    Role: {{.Role}}
  </p>
</div>
{{end}}
```

#### 從 JSON 載入

```go
// 從文件載入
data, _ := os.ReadFile("user-card.json")
restored, _ := FromJSON(string(data))

// 渲染
html := Render(restored)
```

---

### 3. 完整示例

#### HTTP 服務器 + Props 類型 + 模板序列化

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    . "github.com/TimLai666/go-vdom/vdom"
    js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        
        // 使用新的 Props 類型系統
        doc := Document(
            "我的網站",
            []LinkInfo{
                {
                    Rel:  "stylesheet",
                    Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css",
                },
            },
            nil, nil,
            Div(
                Props{
                    "class":   "container",
                    "style":   "padding: 20px;",
                    "data-id": 123,           // int - 自動轉換
                    "visible": true,          // bool - 渲染
                },
                H1("歡迎使用 Go VDOM v1.1.0"),
                
                // 顯示新特性
                Div(
                    Props{"class": "alert alert-info"},
                    H4("✨ 新特性：Props 類型系統"),
                    P("現在支持任意類型的 Props 值！"),
                    Ul(
                        Li("布爾值：", Code("true"), " / ", Code("false")),
                        Li("整數：", Code("42")),
                        Li("浮點數：", Code("19.99")),
                        Li("字符串：", Code("'text'")),
                    ),
                ),
                
                // 互動示例
                Div(
                    Props{"class": "card mt-3"},
                    Div(
                        Props{"class": "card-body"},
                        H5("計數器示例"),
                        Button(Props{
                            "class": "btn btn-primary",
                            "onClick": js.Fn(nil,
                                js.Const("counter", "document.getElementById('counter')"),
                                js.Const("current", "parseInt(counter.innerText)"),
                                js.Const("newValue", "current + 1"),
                                js.El("#counter").SetText("newValue.toString()"),
                            ),
                        }, "增加 +1"),
                        Span(" 計數: "),
                        Span(Props{"id": "counter", "class": "badge bg-primary"}, "0"),
                    ),
                ),
            ),
        )
        
        fmt.Fprint(w, Render(doc))
    })
    
    log.Println("服務器啟動於 http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 🔧 新功能使用指南

### Props 類型系統

#### 布爾值處理
```go
// ✅ 推薦：使用布爾值
Props{
    "disabled": true,   // 渲染為 disabled="true"
    "required": true,   // 渲染為 required="true"
    "hidden":   false,  // 不渲染
}

// ❌ 舊方式：使用字符串
Props{
    "disabled": "true",
    "required": "true",
}
```

#### 數字處理
```go
// ✅ 推薦：直接使用數字
Props{
    "width":    800,
    "height":   600,
    "opacity":  0.8,
    "z-index":  10,
}

// ❌ 舊方式：手動轉換
Props{
    "width":  strconv.Itoa(800),
    "height": fmt.Sprintf("%d", 600),
}
```

### 模板序列化

#### 提取模板變數
```go
vnode := Div(
    Props{"id": "user-{{.ID}}"},
    H1("{{.Name}}"),
    P("{{.Email}}"),
)

vars := ExtractTemplateVars(vnode)
// 返回: [".ID", ".Name", ".Email"]

fmt.Println("模板需要的數據:")
for _, v := range vars {
    fmt.Printf("  - %s\n", v)
}
```

#### VNode 克隆
```go
// 創建基礎按鈕
baseButton := Button(
    Props{"class": "btn", "type": "button"},
    "按鈕",
)

// 克隆並創建變體
primaryButton := CloneVNode(baseButton)
primaryButton.Props["class"] = "btn btn-primary"

secondaryButton := CloneVNode(baseButton)
secondaryButton.Props["class"] = "btn btn-secondary"

// 原始按鈕不受影響
```

#### Props 合併
```go
baseStyle := Props{
    "class": "btn",
    "type":  "button",
}

primaryStyle := Props{
    "class": "btn btn-primary",
}

extraProps := Props{
    "id":       "submit",
    "disabled": false,
}

// 合併（後面的覆蓋前面的）
merged := MergeProps(baseStyle, primaryStyle, extraProps)
// 結果: {
//   "class": "btn btn-primary",  // 被覆蓋
//   "type": "button",
//   "id": "submit",
//   "disabled": false,
// }
```

---

## 📚 完整示例

### 創建組件模板庫

```go
package main

import (
    "fmt"
    "os"
    . "github.com/TimLai666/go-vdom/vdom"
)

func main() {
    // 定義組件模板
    components := map[string]VNode{
        "card": Div(
            Props{"class": "card"},
            Div(
                Props{"class": "card-header"},
                H3("{{.Title}}"),
            ),
            Div(
                Props{"class": "card-body"},
                P("{{.Content}}"),
            ),
        ),
        
        "alert": Div(
            Props{"class": "alert alert-{{.Type}}"},
            Strong("{{.Title}}"),
            Span(" {{.Message}}"),
        ),
        
        "button": Button(
            Props{
                "class":    "btn btn-{{.Variant}}",
                "type":     "{{.Type}}",
                "disabled": "{{.Disabled}}",
            },
            "{{.Text}}",
        ),
    }
    
    // 保存所有組件
    os.MkdirAll("templates", 0755)
    
    for name, vnode := range components {
        // 保存為 Go Template
        template := SaveTemplate(name, vnode)
        filename := fmt.Sprintf("templates/%s.tmpl", name)
        os.WriteFile(filename, []byte(template), 0644)
        fmt.Printf("✓ 已保存: %s\n", filename)
        
        // 保存為 JSON
        jsonStr, _ := ToJSON(vnode)
        jsonFile := fmt.Sprintf("templates/%s.json", name)
        os.WriteFile(jsonFile, []byte(jsonStr), 0644)
        
        // 提取變數
        vars := ExtractTemplateVars(vnode)
        fmt.Printf("  變數: %v\n", vars)
    }
}
```

---

## 🎓 學習路徑

### 新手入門
1. 閱讀本文檔（QUICK_START_V1.1.md）
2. 運行 `go run examples/01_basic_usage.go`
3. 嘗試使用新的 Props 類型系統

### 進階使用
1. 運行 `go run examples/04_template_serialization.go`
2. 閱讀 [IMPROVEMENTS.md](IMPROVEMENTS.md)
3. 閱讀 [DOCUMENTATION.md](DOCUMENTATION.md)

### 完整參考
- [README.md](README.md) - 完整說明
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - 快速參考
- [IMPROVEMENTS.md](IMPROVEMENTS.md) - 改進說明

---

## ❓ 常見問題

### Q: 現有代碼需要修改嗎？

**A:** 不需要！完全向後兼容。

```go
// 舊代碼仍然有效
Props{"class": "container", "id": "main"}

// 也可以使用新特性
Props{"class": "container", "visible": true, "count": 42}
```

### Q: Props 支持哪些類型？

**A:** 支持所有類型：
- `string` - 直接使用
- `bool` - true 渲染，false 省略
- `int`, `int64`, `uint` 等 - 自動轉換
- `float32`, `float64` - 自動轉換
- `JSAction` - 特殊處理
- 任何其他類型 - 使用 `fmt.Sprint()` 轉換

### Q: 為什麼選擇 Go Template？

**A:** 因為：
- ✅ 與 Go `html/template` 無縫集成
- ✅ 標準庫支持，無需額外依賴
- ✅ 支持條件、循環等控制流
- ✅ 良好的性能和安全性

### Q: JSON 序列化有什麼限制？

**A:** 
- ✅ VNode 結構完全保留
- ✅ Props 值正確序列化
- ⚠️ JSAction 保留為 `{"Code": "..."}`
- ⚠️ 函數無法序列化（需特殊處理）

---

## 🔗 相關資源

- [GitHub 倉庫](https://github.com/TimLai666/go-vdom)
- [完整文檔](DOCUMENTATION.md)
- [改進說明](IMPROVEMENTS.md)
- [變更日誌](CHANGELOG.md)
- [示例程序](examples/)

---

## 🎉 開始使用

```bash
# 安裝
go get github.com/TimLai666/go-vdom@v1.1.0

# 運行示例
go run examples/01_basic_usage.go      # 端口 8080
go run examples/02_components.go       # 端口 8081
go run examples/03_javascript_dsl.go   # 端口 8082
go run examples/04_template_serialization.go  # 端口 8083

# 創建你的第一個應用
# 複製上面的完整示例，開始編碼！
```

**享受 go-vdom v1.1.0 帶來的全新體驗！** 🚀

---

**版本**: v1.1.0  
**更新日期**: 2025-01-24  
**作者**: TimLai666