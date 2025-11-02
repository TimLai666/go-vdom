# Go VDOM 示例集合

這個目錄包含了各種 go-vdom 的使用示例，從基礎到進階，幫助你快速上手。

## 示例列表

### 01_basic_usage.go - 基本用法
展示如何創建簡單的 HTML 頁面，包括：
- Document 函數的使用
- 基本 HTML 元素
- Bootstrap 樣式集成
- 頁面佈局和結構

**運行方式:**
```bash
go run examples/01_basic_usage.go
```
訪問: http://localhost:8080

---

### 02_components.go - 組件系統
展示如何創建和使用可重用的組件，包括：
- Alert 組件
- Card 組件
- Badge 組件
- Button 組件
- UserCard 組件
- 組件組合和嵌套

**運行方式:**
```bash
go run examples/02_components.go
```
訪問: http://localhost:8081

---

### 03_javascript_dsl.go - JavaScript DSL
展示如何使用 JavaScript DSL 創建交互式頁面，包括：
- DOM 操作（SetText, SetHTML, AddClass, RemoveClass）
- 變數定義（Let, Const）
- 事件處理（OnClick）
- 表單處理
- 動態創建元素
- Try/Catch 錯誤處理
- DomReady 初始化

**運行方式:**
```bash
go run examples/03_javascript_dsl.go
```
訪問: http://localhost:8082

---

## 快速開始

1. 確保已安裝 Go 1.24.1 或更高版本
2. 克隆倉庫並進入項目目錄
3. 運行任意示例文件
4. 在瀏覽器中訪問對應的端口

## 學習路徑

建議按以下順序學習示例：

1. **01_basic_usage.go** - 了解基本概念和結構
2. **02_components.go** - 學習組件化開發
3. **03_javascript_dsl.go** - 掌握 JavaScript DSL 和交互功能

## 自定義示例

你可以基於這些示例創建自己的應用：

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
            "我的應用",
            nil, nil, nil,
            Div(
                Props{"class": "container"},
                H1("歡迎使用 go-vdom"),
            ),
        )
        
        fmt.Fprint(w, Render(doc))
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 更多資源

- [完整文檔](../DOCUMENTATION.md)
- [主 README](../README.md)
- [GitHub 倉庫](https://github.com/TimLai666/go-vdom)

## 提示

- 所有示例使用不同的端口，可以同時運行多個示例
- 打開瀏覽器控制台查看 JavaScript 日誌輸出
- 每個示例都是獨立的，可以單獨運行
- 示例代碼包含詳細註釋，便於理解

## 需要幫助？

如果遇到問題：
1. 查看 [DOCUMENTATION.md](../DOCUMENTATION.md) 中的詳細文檔
2. 查看 [故障排除](../DOCUMENTATION.md#故障排除) 章節
3. 在 GitHub 上提交 Issue

---

**Happy Coding! 🚀**