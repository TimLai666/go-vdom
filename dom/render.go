// render.go
package dom

import (
	"fmt"
	"html"
	"strings"
)

// Render 將虛擬DOM節點轉換為HTML字符串
// 改進：
//   - 對屬性值做最小轉義，並在屬性值為字串 "false" 時省略該屬性（便於處理布林屬性表示法）
//   - 特別處理 `onDOMReady` 屬性：只接受透過 Component 第二參數注入的 JS 函數（建議由 jsdsl.Fn 建立）；該函數會在 DOMContentLoaded 時被呼叫。
//     注意：不再支援舊的 `onMount` / `onmount` 屬性；所有初始化邏輯必須透過 Component 的第二個參數注入。
func Render(v VNode) string {
	if v.Tag == "" {
		return v.Content
	}

	var sb strings.Builder
	sb.WriteString("<" + v.Tag)

	// 收集 onDOMReady（如果有），但不要直接作為屬性輸出
	var onDOMReady string

	for k, rawVal := range v.Props {
		// 當屬性名是 onDOMReady 時，保留其 JS 函數內容以便在 DOMContentLoaded 時呼叫，並跳過將其作為 HTML 屬性輸出
		// 注意：renderer 僅支援 `onDOMReady`，且該屬性應由 Component 的第二個參數注入（通常由 jsdsl.Fn 產生）。
		if k == "onDOMReady" {
			switch t := rawVal.(type) {
			case string:
				onDOMReady = t
			case JSAction:
				onDOMReady = t.Code
			default:
				// 不處理其他類型
			}
			continue
		}

		// 如果屬性是事件處理器（以 on 開頭，例如 onClick/onChange）
		if len(k) > 2 && strings.HasPrefix(k, "on") {
			eventName := strings.ToLower(k[2:])
			switch t := rawVal.(type) {
			case JSAction:
				// JSAction 直接作為內聯事件處理器
				// 用戶應使用 js.Do() 或 js.AsyncDo() 來創建 IIFE
				// 安全處理：避免內嵌引號破壞屬性結構
				safeCode := strings.ReplaceAll(t.Code, "\"", "&quot;")
				safeCode = strings.ReplaceAll(safeCode, "\n", " ")
				safeCode = strings.ReplaceAll(safeCode, "\r", " ")
				sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, safeCode))
			case string:
				// 字符串直接作為內聯事件處理器
				escaped := html.EscapeString(t)
				escaped = strings.ReplaceAll(escaped, "\n", " ")
				escaped = strings.ReplaceAll(escaped, "\r", " ")
				sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, escaped))
			case ServerHandlerRef:
				// 伺服器端 handler 引用，產生 data-gvd-server-handler 屬性
				sb.WriteString(fmt.Sprintf(" data-gvd-server-handler=\"%s|%s\"", t.ID, eventName))
			default:
				// fallback：將值轉為字串並當普通屬性輸出
				valStr := fmt.Sprint(rawVal)
				if valStr == "false" {
					continue
				}
				escaped := html.EscapeString(valStr)
				escaped = strings.ReplaceAll(escaped, "\n", " ")
				escaped = strings.ReplaceAll(escaped, "\r", " ")
				sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, escaped))
			}
			// 事件處理器已處理，跳過一般屬性處理
			continue
		}

		// 處理一般屬性（如果不是事件，也不是 onmount）
		var valStr string
		switch t := rawVal.(type) {
		case string:
			valStr = t
		case JSAction:
			// 若使用者傳入 JSAction 作為屬性（較少見），我們用其 Code 的字串形式
			valStr = t.Code
		default:
			valStr = fmt.Sprint(t)
		}

		// 當屬性值為字串 "false" 時，視為不設置該布林屬性，跳過
		if valStr == "false" {
			continue
		}

		// 對屬性值做最小的 HTML 轉義，並將換行符替換成空格，避免破壞屬性語法
		escaped := html.EscapeString(valStr)
		escaped = strings.ReplaceAll(escaped, "\n", " ")
		escaped = strings.ReplaceAll(escaped, "\r", " ")

		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, escaped))
	}
	sb.WriteString(">")

	if v.Content != "" {
		sb.WriteString(v.Content)
	}
	for _, c := range v.Children {
		sb.WriteString(Render(c))
	}

	sb.WriteString(fmt.Sprintf("</%s>", v.Tag))

	// 如果有 onDOMReady，注入對應的 <script>
	if onDOMReady != "" {
		// 以簡單方式避免原始 onDOMReady 中出現 "</script>" 導致 HTML 結構中斷
		safeScript := strings.ReplaceAll(onDOMReady, "</script>", "</scr\" + \"ipt>")
		// onDOMReady 應由 jsdsl.Fn 產生函數表達式，直接在 DOMContentLoaded 時呼叫
		// 使用立即執行函數來確保只執行一次
		onReadyWrapper := "(function(){var fn=" + safeScript + ";if(document.readyState==='loading'){document.addEventListener('DOMContentLoaded',fn);}else{fn();}})();"
		sb.WriteString("<script>")
		sb.WriteString(onReadyWrapper)
		sb.WriteString("</script>")
	}

	return sb.String()
}
