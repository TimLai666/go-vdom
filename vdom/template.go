// template.go
package vdom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// TemplateData 用於序列化 VNode 到模板格式
type TemplateData struct {
	Nodes []VNode `json:"nodes"`
}

// ToGoTemplate 將 VNode 轉換為 Go template 格式
// 這樣可以將模板保存到文件，之後再加載使用
func ToGoTemplate(v VNode) string {
	var sb strings.Builder
	renderToGoTemplate(&sb, v, 0)
	return sb.String()
}

// renderToGoTemplate 遞歸渲染 VNode 為 Go template 語法
func renderToGoTemplate(sb *strings.Builder, v VNode, depth int) {
	indent := strings.Repeat("  ", depth)

	if v.Tag == "" {
		// 文本節點
		if v.Content != "" {
			sb.WriteString(indent)
			sb.WriteString(v.Content)
			sb.WriteString("\n")
		}
		return
	}

	// 開始標籤
	sb.WriteString(indent)
	sb.WriteString("<")
	sb.WriteString(v.Tag)

	// 渲染屬性
	for k, val := range v.Props {
		sb.WriteString(" ")
		sb.WriteString(k)
		sb.WriteString("=\"")

		// 根據值的類型進行處理
		switch t := val.(type) {
		case string:
			// 檢查是否是模板變數（{{...}}）
			if strings.Contains(t, "{{") && strings.Contains(t, "}}") {
				sb.WriteString(t)
			} else {
				sb.WriteString(template.HTMLEscapeString(t))
			}
		case JSAction:
			sb.WriteString("{{/* JSAction */}}")
			sb.WriteString(template.HTMLEscapeString(t.Code))
		case bool:
			if t {
				sb.WriteString("true")
			} else {
				sb.WriteString("false")
			}
		case int, int8, int16, int32, int64:
			sb.WriteString(fmt.Sprintf("%d", t))
		case uint, uint8, uint16, uint32, uint64:
			sb.WriteString(fmt.Sprintf("%d", t))
		case float32, float64:
			sb.WriteString(fmt.Sprintf("%f", t))
		default:
			sb.WriteString(fmt.Sprintf("%v", val))
		}

		sb.WriteString("\"")
	}

	sb.WriteString(">")

	// 如果是自閉合標籤，直接結束
	if isSelfClosingTag(v.Tag) {
		sb.WriteString("\n")
		return
	}

	sb.WriteString("\n")

	// 渲染內容
	if v.Content != "" {
		sb.WriteString(indent)
		sb.WriteString("  ")
		sb.WriteString(v.Content)
		sb.WriteString("\n")
	}

	// 渲染子元素
	for _, child := range v.Children {
		renderToGoTemplate(sb, child, depth+1)
	}

	// 結束標籤
	sb.WriteString(indent)
	sb.WriteString("</")
	sb.WriteString(v.Tag)
	sb.WriteString(">\n")
}

// isSelfClosingTag 判斷是否為自閉合標籤
func isSelfClosingTag(tag string) bool {
	selfClosing := map[string]bool{
		"area": true, "base": true, "br": true, "col": true,
		"embed": true, "hr": true, "img": true, "input": true,
		"link": true, "meta": true, "param": true, "source": true,
		"track": true, "wbr": true,
	}
	return selfClosing[strings.ToLower(tag)]
}

// ToJSON 將 VNode 序列化為 JSON
func ToJSON(v VNode) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal VNode to JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON 從 JSON 反序列化 VNode
func FromJSON(jsonStr string) (VNode, error) {
	var v VNode
	err := json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		return VNode{}, fmt.Errorf("failed to unmarshal JSON to VNode: %w", err)
	}
	return v, nil
}

// ExecuteGoTemplate 執行 Go template 並返回渲染結果
// tmplStr: Go template 字符串
// data: 模板數據
func ExecuteGoTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("vdom").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ToCompactJSON 將 VNode 序列化為緊湊的 JSON（無縮進）
func ToCompactJSON(v VNode) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal VNode to compact JSON: %w", err)
	}
	return string(data), nil
}

// VNodeToTemplateFunc 創建一個可以在 Go template 中使用的函數
// 用於將 VNode 渲染為 HTML
func VNodeToTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"renderVNode": func(v VNode) template.HTML {
			return template.HTML(Render(v))
		},
		"toGoTemplate": func(v VNode) string {
			return ToGoTemplate(v)
		},
	}
}

// SaveTemplate 保存 VNode 為 Go template 格式到字符串
// 包含模板定義和註釋
func SaveTemplate(name string, v VNode) string {
	var sb strings.Builder

	// 添加模板定義
	sb.WriteString(fmt.Sprintf("{{/* Template: %s */}}\n", name))
	sb.WriteString("{{define \"" + name + "\"}}\n")

	// 渲染模板內容
	sb.WriteString(ToGoTemplate(v))

	// 結束模板定義
	sb.WriteString("{{end}}\n")

	return sb.String()
}

// WrapWithLayout 將 VNode 包裝到佈局模板中
// layoutName: 佈局模板名稱
// contentName: 內容模板名稱
func WrapWithLayout(layoutName, contentName string, content VNode) string {
	var sb strings.Builder

	// 定義內容模板
	sb.WriteString(SaveTemplate(contentName, content))
	sb.WriteString("\n")

	// 使用佈局模板
	sb.WriteString(fmt.Sprintf("{{template \"%s\" .}}\n", layoutName))

	return sb.String()
}

// MergeProps 合併多個 Props，後面的會覆蓋前面的
func MergeProps(propsList ...Props) Props {
	result := make(Props)
	for _, props := range propsList {
		for k, v := range props {
			result[k] = v
		}
	}
	return result
}

// CloneVNode 深度克隆一個 VNode
func CloneVNode(v VNode) VNode {
	clone := VNode{
		Tag:     v.Tag,
		Content: v.Content,
		Props:   make(Props),
	}

	// 克隆 Props
	for k, v := range v.Props {
		clone.Props[k] = v
	}

	// 克隆子節點
	if len(v.Children) > 0 {
		clone.Children = make([]VNode, len(v.Children))
		for i, child := range v.Children {
			clone.Children[i] = CloneVNode(child)
		}
	}

	return clone
}

// ExtractTemplateVars 從 VNode 中提取所有模板變數（{{...}}）
func ExtractTemplateVars(v VNode) []string {
	vars := make(map[string]bool)
	extractVarsRecursive(v, vars)

	result := make([]string, 0, len(vars))
	for k := range vars {
		result = append(result, k)
	}
	return result
}

// extractVarsRecursive 遞歸提取模板變數
func extractVarsRecursive(v VNode, vars map[string]bool) {
	// 檢查內容
	if v.Content != "" {
		extractVarsFromString(v.Content, vars)
	}

	// 檢查屬性
	for _, val := range v.Props {
		if str, ok := val.(string); ok {
			extractVarsFromString(str, vars)
		}
	}

	// 遞歸檢查子節點
	for _, child := range v.Children {
		extractVarsRecursive(child, vars)
	}
}

// extractVarsFromString 從字符串中提取 {{...}} 變數
func extractVarsFromString(s string, vars map[string]bool) {
	start := 0
	for {
		idx := strings.Index(s[start:], "{{")
		if idx == -1 {
			break
		}
		idx += start

		endIdx := strings.Index(s[idx:], "}}")
		if endIdx == -1 {
			break
		}
		endIdx += idx

		varName := strings.TrimSpace(s[idx+2 : endIdx])
		// 排除註釋和控制結構
		if !strings.HasPrefix(varName, "/*") &&
			!strings.HasPrefix(varName, "if") &&
			!strings.HasPrefix(varName, "range") &&
			!strings.HasPrefix(varName, "with") &&
			!strings.HasPrefix(varName, "end") &&
			!strings.HasPrefix(varName, "define") &&
			!strings.HasPrefix(varName, "template") {
			vars[varName] = true
		}

		start = endIdx + 2
	}
}

// ConvertPropsToAny 確保 Props 中的值都是正確的類型
// 這個函數可以用來處理從 JSON 反序列化後的 Props
func ConvertPropsToAny(props Props) Props {
	result := make(Props)
	for k, v := range props {
		switch t := v.(type) {
		case map[string]interface{}:
			// 如果是 map，可能是 JSAction 的 JSON 表示
			if code, ok := t["Code"].(string); ok {
				result[k] = JSAction{Code: code}
			} else {
				result[k] = v
			}
		default:
			result[k] = v
		}
	}
	return result
}
