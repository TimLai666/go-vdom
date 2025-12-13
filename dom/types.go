// types.go
package dom

// Props 是一個用於存儲元素屬性的映射
// 調整：Props 現在允許儲存多種型別（例如 string、JSAction、ServerHandlerRef 等）
// 這使得你可以在 Props 中直接傳入 JSAction（代表要在 client 端執行的 JS 片段）
// 或其他未來擴展用的引用型別。
type Props map[string]any

// VNode 表示虛擬DOM中的一個節點
type VNode struct {
	Tag      string
	Props    Props
	Children []VNode
	Content  string
}

// JSAction 代表一段要在客戶端執行的 JavaScript 代碼片段。
// 在 renderer 中遇到 Props 傳入 JSAction 時，會將其轉為 client-side handler。
type JSAction struct {
	Code string
}

// ServerHandlerRef 是一個簡單的引用型別，用於在 Props 中指向
// 由伺服器端註冊的 handler（例如 RegisterServerHandler 回傳的 id）。
// 這樣可以在 renderer 中識別並產生相對應的 data 屬性（例如 data-gvd-server-handler）。
type ServerHandlerRef struct {
	ID string
}
