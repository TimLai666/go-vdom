package jsdsl

import (
	"fmt"
	"math/rand"
	"strings"

	. "github.com/TimLai666/go-vdom/vdom"
)

// variable 代表一個 JS 變數，可以呼叫常用屬性與方法
type variable struct {
	Name string
}

// variable 建立一個 variable 物件
func V(varName string) variable {
	return variable{Name: varName}
}

// SetHTML 設定 innerHTML
func (v variable) SetHTML(html string) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.innerHTML = %s", v.Name, html)}
}

// SetText 設定 innerText
func (v variable) SetText(text string) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.innerText = %s", v.Name, text)}
}

// AddClass 新增 class
func (v variable) AddClass(class string) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.classList.add(%s)", v.Name, class)}
}

// RemoveClass 移除 class
func (v variable) RemoveClass(class string) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.classList.remove(%s)", v.Name, class)}
}

// CallMethod 呼叫任意方法
func (v variable) CallMethod(method string, args ...string) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.%s(%s)", v.Name, method, strings.Join(args, ", "))}
}

// VRef 代表一個 JS 變數（Variable Reference）for 只取變數名
func VRef(varName string) JSAction {
	return JSAction{Code: varName}
}

// Ptr returns a pointer to the provided JSAction.
// This is a small helper so callers can easily pass *JSAction where needed,
// e.g. Component(template, jsdsl.Ptr(jsdsl.Fn(nil, ...)), ...)
// The parameter is safe to take address of because it will escape to the heap.
func Ptr(a JSAction) *JSAction {
	return &a
}

type Elem struct {
	Selector string
	VarName  string
}

type ElemList struct {
	Selector string
}

func El(selector string) Elem {
	return Elem{Selector: selector}
}

func Els(selector string) ElemList {
	return ElemList{Selector: selector}
}

func (e Elem) Ref() string {
	if e.VarName != "" {
		return e.VarName
	}
	return fmt.Sprintf("document.querySelector('%s')", e.Selector)
}

func (e ElemList) Ref() string {
	return fmt.Sprintf("document.querySelectorAll('%s')", e.Selector)
}

// Fn 創建一個函數，支援傳入參數
// 如果不需要參數，可以傳入 nil 作為 params 參數
func Fn(params []string, actions ...JSAction) JSAction {
	var sb strings.Builder

	// 創建一個匿名函數
	paramsStr := ""
	if params != nil {
		paramsStr = strings.Join(params, ", ")
	}
	sb.WriteString(fmt.Sprintf("(%s) => {\n", paramsStr))

	// 添加函數體
	for _, a := range actions {
		line := strings.TrimSpace(a.Code)
		if !strings.HasSuffix(line, ";") {
			line += ";"
		}
		sb.WriteString("  " + line + "\n")
	}

	sb.WriteString("}")
	return JSAction{Code: sb.String()}
}

// AsyncFn 創建一個異步函數，支援 await 語法
// 如果不需要參數，可以傳入 nil 作為 params 參數
func AsyncFn(params []string, actions ...JSAction) JSAction {
	var sb strings.Builder

	// 創建一個異步匿名函數
	paramsStr := ""
	if params != nil {
		paramsStr = strings.Join(params, ", ")
	}
	sb.WriteString(fmt.Sprintf("async (%s) => {\n", paramsStr))

	// 添加函數體
	for _, a := range actions {
		line := strings.TrimSpace(a.Code)
		if !strings.HasSuffix(line, ";") {
			line += ";"
		}
		sb.WriteString("  " + line + "\n")
	}

	sb.WriteString("}")
	return JSAction{Code: sb.String()}
}

// Call 調用一個函數，傳入參數
func Call(fnExpr any, args ...any) JSAction {
	var processedArgs []string
	var fnExprStr string

	// 處理函數表達式
	switch v := fnExpr.(type) {
	case string:
		fnExprStr = v
	case JSAction:
		fnExprStr = v.Code
	default:
		fnExprStr = fmt.Sprintf("%v", fnExpr)
	}

	// 處理參數
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			processedArgs = append(processedArgs, v)
		case JSAction:
			processedArgs = append(processedArgs, v.Code)
		default:
			processedArgs = append(processedArgs, fmt.Sprintf("%v", v))
		}
	}

	argsStr := strings.Join(processedArgs, ", ")
	return JSAction{Code: fmt.Sprintf("%s(%s)", fnExprStr, argsStr)}
}

// Method 調用對象的方法，更符合直觀的呼叫方式
// 用法：Method("object", "methodName", arg1, arg2, ...)
func CallMethod(objExpr string, methodName string, args ...any) JSAction {
	return Call(fmt.Sprintf("%s.%s", objExpr, methodName), args...)
}

func DomReady(actions ...JSAction) JSAction {
	return JSAction{Code: fmt.Sprintf(`document.addEventListener("DOMContentLoaded",
%s);`, indent(Fn(nil, actions...).Code, "  "))}
}

func (el Elem) SetText(text string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.innerText = %s`, el.Ref(), text)}
}

func (el Elem) SetHTML(html string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.innerHTML = %s`, el.Ref(), html)}
}

func (el Elem) AddClass(class string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.classList.add('%s')`, el.Ref(), class)}
}

func (el Elem) RemoveClass(class string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.classList.remove('%s')`, el.Ref(), class)}
}

func Log(msg string) JSAction {
	return JSAction{Code: fmt.Sprintf(`console.log(%s)`, msg)}
}

func Redirect(url string) JSAction {
	return JSAction{Code: fmt.Sprintf(`location.href = '%s'`, url)}
}

func (el Elem) OnClick(action JSAction) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.addEventListener('click', function() {
%s
});`, el.Ref(), indent(action.Code, "  "))}
}

func Alert(jsExpr string) JSAction {
	return JSAction{Code: fmt.Sprintf(`alert(%s)`, jsExpr)}
}

func (el Elem) InnerText() string {
	return fmt.Sprintf("%s.innerText", el.Ref())
}

func (el Elem) InnerHTML() string {
	return fmt.Sprintf("%s.innerHTML", el.Ref())
}

// ForEachJS 遍歷任意 JavaScript 數組或可迭代對象（前端渲染）
// 參數：
// - arrayExpr: 數組表達式（如 "myArray", "[1,2,3]", "data.items"）
// - itemVar: 項目變數名稱（如 "item", "user"）
// - actions: 對每個項目執行的動作
// 用法：js.ForEachJS("items", "item", js.Log("item"))
func ForEachJS(arrayExpr string, itemVar string, actions ...JSAction) JSAction {
	var sb strings.Builder
	for _, a := range actions {
		line := strings.TrimSpace(a.Code)
		if !strings.HasSuffix(line, ";") {
			line += ";"
		}
		sb.WriteString(line + "\n")
	}

	return JSAction{
		Code: fmt.Sprintf(`%s.forEach(function(%s) {
%s});`, arrayExpr, itemVar, indent(sb.String(), "  ")),
	}
}

// ForEachWithIndexJS 遍歷數組並提供索引（前端渲染）
// 參數：
// - arrayExpr: 數組表達式
// - itemVar: 項目變數名稱
// - indexVar: 索引變數名稱
// - actions: 對每個項目執行的動作
func ForEachWithIndexJS(arrayExpr string, itemVar string, indexVar string, actions ...JSAction) JSAction {
	var sb strings.Builder
	for _, a := range actions {
		line := strings.TrimSpace(a.Code)
		if !strings.HasSuffix(line, ";") {
			line += ";"
		}
		sb.WriteString(line + "\n")
	}

	return JSAction{
		Code: fmt.Sprintf(`%s.forEach(function(%s, %s) {
%s});`, arrayExpr, itemVar, indexVar, indent(sb.String(), "  ")),
	}
}

// ForEachElement 遍歷 DOM 元素列表（保留向後兼容）
// 這是專門用於 DOM 元素操作的版本
func ForEachElement(arrayExpr string, fn func(el Elem) JSAction) JSAction {
	el := "el"
	return JSAction{
		Code: fmt.Sprintf(`%s.forEach(function(%s) {
%s
});`, arrayExpr, el, indent(fn(Elem{VarName: el}).Code, "  ")),
	}
}

func QueryEach(list ElemList, fn func(el Elem) JSAction) JSAction {
	return ForEachElement(list.Ref(), fn)
}

func indent(code string, prefix string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func escapeJS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func Let(varName string, value string) JSAction {
	return JSAction{Code: fmt.Sprintf("let %s = %s", varName, value)}
}

func Const(varName string, value string) JSAction {
	return JSAction{Code: fmt.Sprintf("const %s = %s", varName, value)}
}

// FetchOption 代表一個 fetch 請求的選項
type FetchOption struct {
	Key   string
	Value string
}

// WithMethod 設定 HTTP 方法
func WithMethod(method string) FetchOption {
	return FetchOption{Key: "method", Value: method}
}

// WithHeader 設定 HTTP 頭
func WithHeader(name, value string) FetchOption {
	return FetchOption{Key: "headers." + name, Value: value}
}

// WithBody 設定請求的主體
func WithBody(body string) FetchOption {
	return FetchOption{Key: "body", Value: body}
}

// WithContentType 設定 Content-Type 頭
func WithContentType(contentType string) FetchOption {
	return WithHeader("Content-Type", contentType)
}

// WithJSON 設定 Content-Type 為 application/json 並且將主體設定為 JSON 字符串
func WithJSON(jsonObject string) []FetchOption {
	return []FetchOption{
		WithContentType("application/json"),
		WithBody(jsonObject),
	}
}

// WithFormData 設定 Content-Type 為 application/x-www-form-urlencoded
func WithFormData(formData map[string]string) []FetchOption {
	var values []string
	for key, value := range formData {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}
	formBody := strings.Join(values, "&")

	return []FetchOption{
		WithContentType("application/x-www-form-urlencoded"),
		WithBody(formBody),
	}
}

// ResponseType 定義了 fetch 響應的解析方式
type ResponseType string

const (
	JSONResponse ResponseType = "json"
	TextResponse ResponseType = "text"
	BlobResponse ResponseType = "blob"
)

// FetchRequest 創建一個通用的 fetch 請求
func FetchRequest(url string, options ...FetchOption) JSAction {
	return buildFetch(url, JSONResponse, "", "", "", options...)
}

// WithThen 添加 then 處理器
func WithThen(thenCodes ...interface{}) JSAction {
	var sb strings.Builder

	for _, thenCode := range thenCodes {
		var codeStr string
		switch v := thenCode.(type) {
		case string:
			codeStr = v
		case JSAction:
			codeStr = v.Code
		default:
			codeStr = fmt.Sprintf("%v", thenCode)
		}

		// 添加代碼，確保每行都有適當的縮進
		lines := strings.Split(codeStr, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				sb.WriteString("  " + line + "\n")
			}
		}
	}

	return JSAction{Code: fmt.Sprintf("then(data => {\n%s})", sb.String())}
}

// Deprecated: WithCatch was removed. Use TryCatch(baseAction JSAction, catchFn *JSAction, finallyFn *JSAction) instead.
// This placeholder comment remains to indicate the old API has been intentionally removed.

// Deprecated: WithFinally was removed. Use TryCatch(baseAction JSAction, catchFn *JSAction, finallyFn *JSAction) instead.
// This placeholder comment remains to indicate the old API has been intentionally removed.

// WithResponseType 設定響應類型
func WithResponseType(responseType ResponseType) JSAction {
	return JSAction{Code: fmt.Sprintf("response_type:%s", string(responseType))}
}

// 處理多種可能的代碼輸入類型
func processCodeInput(code interface{}) string {
	switch v := code.(type) {
	case string:
		return v
	case JSAction:
		return v.Code
	default:
		return fmt.Sprintf("%v", code)
	}
}

// 構建完整的 fetch 請求
func buildFetch(url string, responseType ResponseType, then interface{}, catch interface{}, finally interface{}, options ...FetchOption) JSAction {
	// 處理輸入代碼
	thenStr := ""
	if then != nil {
		thenStr = processCodeInput(then)
	}

	catchStr := ""
	if catch != nil {
		catchStr = processCodeInput(catch)
	}

	finallyStr := ""
	if finally != nil {
		finallyStr = processCodeInput(finally)
	}

	var optsBuilder strings.Builder
	optsBuilder.WriteString("{\n")

	// 處理所有的選項
	headers := false
	for _, opt := range options {
		if strings.HasPrefix(opt.Key, "headers.") {
			if !headers {
				optsBuilder.WriteString("  headers: {\n")
				headers = true
			}
			headerName := strings.TrimPrefix(opt.Key, "headers.")
			optsBuilder.WriteString(fmt.Sprintf("    '%s': '%s',\n", headerName, escapeJS(opt.Value)))
		} else {
			optsBuilder.WriteString(fmt.Sprintf("  %s: '%s',\n", opt.Key, escapeJS(opt.Value)))
		}
	}

	if headers {
		optsBuilder.WriteString("  },\n")
	}

	optsBuilder.WriteString("}")

	// 構建 fetch 鏈
	var codeBuilder strings.Builder

	codeBuilder.WriteString(fmt.Sprintf("fetch('%s', %s)\n", url, optsBuilder.String()))
	codeBuilder.WriteString("  .then(response => {\n")
	codeBuilder.WriteString("    if (!response.ok) {\n")
	codeBuilder.WriteString("      throw new Error('Network response was not ok: ' + response.status + ' ' + response.statusText);\n")
	codeBuilder.WriteString("    }\n")

	// 根據響應類型返回不同格式的數據
	switch responseType {
	case TextResponse:
		codeBuilder.WriteString("    return response.text();\n")
	case BlobResponse:
		codeBuilder.WriteString("    return response.blob();\n")
	default:
		codeBuilder.WriteString("    return response.json();\n")
	}

	codeBuilder.WriteString("  })")

	// 添加 then 處理器
	if then != "" {
		codeBuilder.WriteString("\n  .then(data => {\n")
		codeBuilder.WriteString(indent(thenStr, "    "))
		codeBuilder.WriteString("\n  })")
	}

	// 添加 catch 處理器
	if catch != "" {
		codeBuilder.WriteString("\n  .catch(error => {\n")
		codeBuilder.WriteString(indent(catchStr, "    "))
		codeBuilder.WriteString("\n  })")
	}

	// 添加 finally 處理器
	if finally != "" {
		codeBuilder.WriteString("\n  .finally(() => {\n")
		codeBuilder.WriteString(indent(finallyStr, "    "))
		codeBuilder.WriteString("\n  })")
	}

	return JSAction{Code: codeBuilder.String()}
}

// TryCatch 生成一個自動執行的 async 函數，含 try/catch/finally 邏輯
//   - baseAction: 由 jsdsl.Fn 產生的函數表達式，放在 try { ... } 中
//   - catchFn: 可選，由 jsdsl.Fn 產生的函數表達式，放在 catch (e) { ... } 中；
//     若提供，會在異常發生時執行（可訪問 error 對象為 `e`）
//   - finallyFn: 可選，由 jsdsl.Fn 產生的函數表達式，放在 finally { ... } 中
//
// 要求：catchFn 與 finallyFn 不能同時為 nil（至少提供一個處理路徑）
//
// 生成的代碼形如：(async () => { try { <baseAction> } catch (e) { <catchFn> } finally { <finallyFn> } })()
// 會立即執行該函數表達式。
func TryCatch(baseAction JSAction, catchFn *JSAction, finallyFn *JSAction) JSAction {
	// 驗證輸入
	if strings.TrimSpace(baseAction.Code) == "" {
		// 沒有有意義的 try 主體，返回空 JSAction
		return JSAction{Code: ""}
	}
	if (catchFn == nil || strings.TrimSpace(catchFn.Code) == "") && (finallyFn == nil || strings.TrimSpace(finallyFn.Code) == "") {
		panic("TryCatch requires at least one of catchFn or finallyFn to be non-nil")
	}

	// 準備 base 代碼（由 jsdsl.Fn 產生，期望是函數表達式或代碼片段）
	baseCode := strings.TrimSpace(baseAction.Code)

	var sb strings.Builder
	// 生成自調用的 async 函數包裝
	sb.WriteString("(async () => { try { ")
	sb.WriteString(baseCode)
	sb.WriteString(" } ")

	// catch 部分（如果提供）
	if catchFn != nil && strings.TrimSpace(catchFn.Code) != "" {
		sb.WriteString("catch (e) { ")
		catchCode := strings.TrimSpace(catchFn.Code)
		sb.WriteString(catchCode)
		sb.WriteString(" } ")
	}

	// finally 部分（如果提供）
	if finallyFn != nil && strings.TrimSpace(finallyFn.Code) != "" {
		sb.WriteString("finally { ")
		finallyCode := strings.TrimSpace(finallyFn.Code)
		sb.WriteString(finallyCode)
		sb.WriteString(" } ")
	}

	// 立即執行該函數
	sb.WriteString("})()")

	return JSAction{Code: sb.String()}
}

// StoreResult 將 fetch 的結果存儲到指定的變數中，並可以執行額外的動作
// 用法：WithThen(StoreResult("resultVar", ...其他動作))
func StoreResult(varName string, additionalActions ...interface{}) JSAction {
	var actionCodes []string

	// 首先將結果存儲到指定的變數中
	actionCodes = append(actionCodes, fmt.Sprintf("const %s = data;", varName))

	// 處理額外的動作
	for _, action := range additionalActions {
		var code string
		switch v := action.(type) {
		case string:
			code = v
		case JSAction:
			code = v.Code
		default:
			code = fmt.Sprintf("%v", action)
		}
		actionCodes = append(actionCodes, code)
	}

	return JSAction{Code: strings.Join(actionCodes, "\n")}
}

// CreateEl 創建一個 DOM 元素，並返回一個 Elem 物件以及創建元素的 JSAction
// tagName：要創建的 HTML 元素標籤名
// varName：可選參數，為創建的元素指定一個變數名稱
func CreateEl(tagName string, varName ...string) (Elem, JSAction) {
	var vName string
	if len(varName) > 0 {
		vName = varName[0]
	} else {
		vName = fmt.Sprintf("el_%s_%d", tagName, generateRandomID())
	}

	jsAction := JSAction{Code: fmt.Sprintf("const %s = document.createElement('%s');", vName, tagName)}
	return Elem{VarName: vName}, jsAction
}

// AppendChild 將子元素添加到父元素中
func (el Elem) AppendChild(child Elem) JSAction {
	return JSAction{Code: fmt.Sprintf("%s.appendChild(%s)", el.Ref(), child.Ref())}
}

// Pipe 將元素和創建元素的動作傳遞給函數，並返回函數執行的結果
// 這允許在一個流暢的鏈式操作中同時處理元素和創建元素的 JSAction
func (el Elem) Pipe(fn func(Elem, JSAction) []JSAction) []JSAction {
	// 創建一個虛擬的 JSAction，因為 Pipe 通常與 CreateEl 一起使用
	// 在這種情況下，實際的創建動作會由調用者傳入
	dummyAction := JSAction{Code: fmt.Sprintf("// Reference to %s", el.Ref())}
	return fn(el, dummyAction)
}

// generateRandomID 生成一個隨機的 ID 用於元素命名
func generateRandomID() int {
	return rand.Intn(10000)
}

// SetTimeout 產生 setTimeout 語法
func SetTimeout(action JSAction, delayMs int) JSAction {
	return JSAction{Code: fmt.Sprintf("setTimeout(\n%s\n, %d)", indent(action.Code, "  "), delayMs)}
}

// SetInterval 產生 setInterval 語法
func SetInterval(action JSAction, intervalMs int) JSAction {
	return JSAction{Code: fmt.Sprintf("setInterval(\n%s\n, %d)", indent(action.Code, "  "), intervalMs)}
}
