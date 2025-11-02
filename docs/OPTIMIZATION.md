# JavaScript 代碼優化

Go VDOM 生成的 JavaScript 代碼經過優化，最小化輸出以減少傳輸大小。

## 代碼最小化

### 特性

所有生成的 JavaScript 代碼都經過最小化處理：

1. **去除換行**：函數體中的語句使用分號連接
2. **去除縮排**：移除所有不必要的空格
3. **去除多餘空格**：只保留語法必需的空格
4. **保持功能**：最小化不影響代碼功能

### 優勢

- ✅ **減少傳輸大小**：約減少 30-50% 的代碼體積
- ✅ **加快載入速度**：更少的數據傳輸
- ✅ **降低帶寬消耗**：對移動設備特別重要
- ✅ **無功能損失**：只移除空白，不改變邏輯

### 示例

#### Fn / AsyncFn

```go
// Go 代碼
js.AsyncFn(nil,
    js.Const("x", "1"),
    js.Log("x"),
    js.Const("y", "2"),
)

// 生成的 JavaScript（最小化）
async()=>{const x=1;console.log(x);const y=2}
```

**對比未最小化的版本：**
```javascript
// 未最小化（約 70 字節）
async () => {
  const x = 1;
  console.log(x);
  const y = 2;
}

// 最小化後（約 45 字節）
async()=>{const x=1;console.log(x);const y=2}
```

#### Try-Catch-Finally

```go
// Go 代碼
js.Try(
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
).Catch(
    js.Log("error.message"),
).Finally(
    js.Log("'done'"),
)

// 生成的 JavaScript（最小化）
try{const data=await fetch('/api');console.log(data)}catch(error){console.log(error.message)}finally{console.log('done')}
```

#### Do / AsyncDo

```go
// Go 代碼
js.AsyncDo(
    js.Const("config", "await fetch('/api/config')"),
    JSAction{Code: "window.appConfig=config"},
)

// 生成的 JavaScript（最小化）
(async()=>{const config=await fetch('/api/config');window.appConfig=config})()
```

## Const 和 Let 支持 JSAction

### 新特性

`Const` 和 `Let` 現在可以接受多種類型的值：

```go
func Const(varName string, value any) JSAction
func Let(varName string, value any) JSAction
```

支持的類型：
- `string` - 字符串（直接使用）
- `JSAction` - JavaScript 動作（提取 Code）
- 其他類型 - 轉換為字符串

### 使用示例

#### 傳入字符串

```go
js.Const("x", "1")
// 生成：const x=1

js.Const("name", "'Alice'")
// 生成：const name='Alice'

js.Let("count", "0")
// 生成：let count=0
```

#### 傳入 JSAction

```go
// 使用函數調用
js.Const("random", js.Call("Math.random"))
// 生成：const random=Math.random()

js.Const("timestamp", js.Call("Date.now"))
// 生成：const timestamp=Date.now()

// 使用標識符
js.Const("data", js.Ident("response.data"))
// 生成：const data=response.data

// 使用自定義 JSAction
js.Const("doubled", JSAction{Code: "x * 2"})
// 生成：const doubled=x * 2
```

#### 組合使用

```go
js.AsyncFn(nil,
    js.Const("response", "await fetch('/api')"),
    js.Const("data", js.Ident("response.json()")),
    js.Const("firstItem", JSAction{Code: "data[0]"}),
    js.Log("firstItem.name"),
)

// 生成：
// async()=>{const response=await fetch('/api');const data=response.json();const firstItem=data[0];console.log(firstItem.name)}
```

#### 實際應用示例

```go
Button(Props{
    "onClick": js.AsyncFn(nil,
        // 獲取響應
        js.Const("response", "await fetch('/api/users')"),
        
        // 解析 JSON（使用 JSAction）
        js.Const("users", JSAction{Code: "await response.json()"}),
        
        // 獲取第一個用戶
        js.Const("firstUser", JSAction{Code: "users[0]"}),
        
        // 計算總數
        js.Const("total", js.Call("users.length")),
        
        // 顯示信息
        js.Alert("'總共 '+total+' 個用戶，第一個是 '+firstUser.name"),
    ),
}, "載入用戶")
```

### 優勢

#### 1. 更靈活的值賦值

```go
// 舊方式：需要手動拼接字符串
js.Const("result", "Math.random() * 10")

// 新方式：可以使用 JSAction
js.Const("result", JSAction{Code: "Math.random() * 10"})
```

#### 2. 更好的代碼組合

```go
// 可以組合多個 JSAction
js.AsyncFn(nil,
    js.Const("x", js.Call("getX")),
    js.Const("y", js.Call("getY")),
    js.Const("sum", JSAction{Code: "x + y"}),
    js.Log("sum"),
)
```

#### 3. 減少字符串操作

```go
// 不需要手動處理引號和轉義
js.Const("element", js.Call("document.getElementById", "'myId'"))
// vs
js.Const("element", "document.getElementById('myId')")  // 手動加引號
```

#### 4. 類型安全

```go
// JSAction 提供了更好的類型檢查
fetchResult := js.Call("fetch", "'/api'")
js.Const("response", fetchResult)  // 類型安全
```

## 最小化策略

### 語句連接

多個語句使用分號連接，不使用換行：

```go
js.Fn(nil,
    js.Const("a", "1"),
    js.Const("b", "2"),
    js.Const("c", "a+b"),
)

// 生成：(x)=>{const a=1;const b=2;const c=a+b}
```

### 空格移除

只保留語法必需的空格：

```go
// 移除：箭頭函數前後的空格
() => {  →  ()=>{

// 移除：大括號內的空格
{ code }  →  {code}

// 保留：關鍵字後的空格（語法要求）
try {  →  try{
catch (e) {  →  catch(e){
```

### 分號處理

自動添加和移除分號：

```go
// 自動移除尾部分號，避免重複
js.Const("x", "1;")  // 輸入有分號
// 生成：const x=1（移除了分號）

// 自動在語句間添加分號
js.Fn(nil,
    js.Const("x", "1"),
    js.Log("x"),
)
// 生成：()=>{const x=1;console.log(x)}（添加分號）
```

## 開發建議

### 開發環境

在開發時，如果需要查看可讀的 JavaScript 代碼：

1. **使用瀏覽器開發者工具**：
   - 大多數瀏覽器的開發者工具會自動格式化顯示的 JavaScript
   - 可以在 Sources/Debugger 面板中查看格式化後的代碼

2. **使用在線工具**：
   - [Prettier Playground](https://prettier.io/playground/)
   - [JS Beautifier](https://beautifier.io/)

3. **調試技巧**：
   ```go
   // 在生成代碼時添加 console.log
   js.AsyncFn(nil,
       js.Log("'=== 開始執行 ==='"),
       js.Const("data", "await fetch('/api')"),
       js.Log("'data:', data"),
       // ... 其他代碼
   )
   ```

### 生產環境

生產環境中使用最小化代碼的最佳實踐：

1. **直接使用**：Go VDOM 生成的代碼已經最小化，無需額外處理
2. **組合壓縮**：如果使用 gzip/brotli 壓縮，效果更佳
3. **CDN 分發**：最小化代碼適合通過 CDN 分發

### 性能對比

假設一個中等規模的頁面：

| 項目 | 未最小化 | 最小化 | 節省 |
|------|---------|--------|------|
| Fn 函數 (10個) | ~800 字節 | ~500 字節 | 37.5% |
| AsyncFn (5個) | ~600 字節 | ~400 字節 | 33.3% |
| Try-Catch (3個) | ~450 字節 | ~280 字節 | 37.8% |
| **總計** | ~1850 字節 | ~1180 字節 | **36.2%** |

加上 gzip 壓縮後：
- 未最小化 + gzip: ~600 字節
- 最小化 + gzip: ~400 字節
- **總節省**: ~33%

## 向後兼容

所有現有代碼無需修改即可享受最小化優勢：

```go
// 舊代碼自動獲得最小化
js.AsyncFn(nil,
    js.Const("data", "await fetch('/api')"),
    js.Log("data"),
)

// Const/Let 仍然支持字符串
js.Const("x", "1")  // 仍然有效

// 現在還支持 JSAction
js.Const("y", JSAction{Code: "x * 2"})  // 新功能
```

## 示例

查看完整的最小化示例：

```bash
go run examples/08_minified_js.go
```

訪問 http://localhost:8087 查看：
- 最小化前後的代碼對比
- 實際大小節省
- JSAction 在 Const/Let 中的使用
- 實時測試和演示

## 總結

Go VDOM 的代碼優化策略：

1. **自動最小化**：所有生成的 JavaScript 代碼都經過最小化
2. **顯著減少體積**：平均減少 30-40% 的代碼大小
3. **零配置**：無需額外設置，自動應用
4. **不影響功能**：只移除空白，保持完整功能
5. **更靈活的 API**：Const/Let 支持 JSAction，更易組合

這些優化讓 Go VDOM 生成的代碼既高效又緊湊，非常適合生產環境使用。