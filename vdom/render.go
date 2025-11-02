// render.go
package vdom

import (
	"fmt"
	"html"
	"strings"
	"time"
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
	// 收集要注入到頁面的 handler registry 腳本片段（對應於 props 中的 JSAction 或命名 handler）
	var injectedHandlers []string

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
				// 為此內聯 JSAction 產生一個唯一 id，並把 handler 注入到 injectedHandlers
				handlerID := fmt.Sprintf("h-%d", time.Now().UnixNano())
				// 在元素上寫入 data-gvd-handler 屬性來引用 handlerID 與事件類型
				sb.WriteString(fmt.Sprintf(" data-gvd-handler=\"%s|%s\"", handlerID, eventName))
				// 安全處理：避免內嵌 </script> 破壞文檔結構
				safeCode := strings.ReplaceAll(t.Code, "</script>", "</scr\" + \"ipt>")
				// 把 handler 註冊片段加入 injectedHandlers，等會會被注入為一段 script
				reg := fmt.Sprintf("window.__gvd=window.__gvd||{};window.__gvd.handlers=window.__gvd.handlers||{};window.__gvd.handlers['%s']={fn:function(evt,el){%s},eventType:'%s'};", handlerID, safeCode, eventName)
				injectedHandlers = append(injectedHandlers, reg)
			case string:
				// 當作命名的全域函式參考（named handler），在 client 端 runtime 會解析並呼叫
				namedID := fmt.Sprintf("named:%s", t)
				sb.WriteString(fmt.Sprintf(" data-gvd-handler=\"%s|%s\"", namedID, eventName))
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

	// 如果有注入的 handler registry 腳本，或 onDOMReady，注入對應的 <script>
	if len(injectedHandlers) > 0 || onDOMReady != "" {
		var scripts []string

		// 如果有要註冊的 handler，先產生一段 registry 初始化腳本
		if len(injectedHandlers) > 0 {
			var regBuilder strings.Builder
			regBuilder.WriteString("(function(){window.__gvd=window.__gvd||{};window.__gvd.handlers=window.__gvd.handlers||{};")
			for _, s := range injectedHandlers {
				regBuilder.WriteString(s)
			}
			regBuilder.WriteString("})();")
			scripts = append(scripts, regBuilder.String())
		}

		// 處理 onDOMReady（如果有）
		if onDOMReady != "" {
			// 以簡單方式避免原始 onDOMReady 中出現 "</script>" 導致 HTML 結構中斷
			safeScript := strings.ReplaceAll(onDOMReady, "</script>", "</scr\" + \"ipt>")
			// 建立包裹：嘗試將 safeScript 當作函數表達式賦值並呼叫；若不是函數則以原始腳本方式執行。
			// 這樣可以支援由 jsdsl.Fn 產生的函數表達式 (e.g. "(...) => { ... }")，同時對於裸腳本也能執行。
			onReadyWrapper := "(function(){function __gvd_call_onready(){try{var __gvd_onready = " + safeScript + "; if(typeof __gvd_onready === 'function'){ __gvd_onready(); return; }}catch(err){} try{ " + safeScript + " }catch(err2){ console.error('onDOMReady execution error', err2); }}; if(document.readyState==='loading'){document.addEventListener('DOMContentLoaded',__gvd_call_onready);}else{__gvd_call_onready();}})();"
			scripts = append(scripts, onReadyWrapper)
		}

		// 注入所有 script 片段
		for _, sc := range scripts {
			sb.WriteString("<script>")
			sb.WriteString(sc)
			sb.WriteString("</script>")
		}
	}

	return sb.String()
}
