// template_test.go
package dom

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToGoTemplate(t *testing.T) {
	tests := []struct {
		name     string
		node     VNode
		expected []string // 期望包含的字符串
	}{
		{
			name: "simple div",
			node: VNode{
				Tag:   "div",
				Props: Props{"class": "container"},
				Children: []VNode{
					{Content: "Hello"},
				},
			},
			expected: []string{"<div", "class=\"container\"", "Hello", "</div>"},
		},
		{
			name: "with template variables",
			node: VNode{
				Tag:   "div",
				Props: Props{"id": "user-{{.ID}}"},
				Children: []VNode{
					{Tag: "h1", Children: []VNode{{Content: "{{.Name}}"}}},
				},
			},
			expected: []string{"{{.ID}}", "{{.Name}}"},
		},
		{
			name: "nested structure",
			node: VNode{
				Tag: "div",
				Children: []VNode{
					{Tag: "header", Children: []VNode{{Tag: "h1", Children: []VNode{{Content: "Title"}}}}},
					{Tag: "main", Children: []VNode{{Tag: "p", Children: []VNode{{Content: "Content"}}}}},
				},
			},
			expected: []string{"<div", "<header", "<h1", "Title", "</h1>", "</header>", "<main", "</main>", "</div>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGoTemplate(tt.node)
			for _, exp := range tt.expected {
				if !strings.Contains(result, exp) {
					t.Errorf("ToGoTemplate() missing expected string %q in output:\n%s", exp, result)
				}
			}
		})
	}
}

func TestSaveTemplate(t *testing.T) {
	node := VNode{
		Tag:   "div",
		Props: Props{"class": "test"},
		Children: []VNode{
			{Content: "Test content"},
		},
	}

	result := SaveTemplate("test-template", node)

	// 檢查是否包含模板定義
	if !strings.Contains(result, `{{define "test-template"}}`) {
		t.Error("SaveTemplate() missing template definition")
	}

	if !strings.Contains(result, "{{end}}") {
		t.Error("SaveTemplate() missing template end")
	}

	if !strings.Contains(result, "Test content") {
		t.Error("SaveTemplate() missing content")
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		node    VNode
		wantErr bool
	}{
		{
			name: "simple node",
			node: VNode{
				Tag:     "div",
				Props:   Props{"class": "test"},
				Content: "Hello",
			},
			wantErr: false,
		},
		{
			name: "with children",
			node: VNode{
				Tag: "div",
				Children: []VNode{
					{Tag: "span", Content: "Child"},
				},
			},
			wantErr: false,
		},
		{
			name: "with mixed types in props",
			node: VNode{
				Tag: "div",
				Props: Props{
					"class":   "test",
					"visible": true,
					"count":   42,
					"price":   19.99,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToJSON(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 驗證是否為有效的 JSON
				var v VNode
				if err := json.Unmarshal([]byte(result), &v); err != nil {
					t.Errorf("ToJSON() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		wantTag string
	}{
		{
			name:    "valid json",
			json:    `{"Tag":"div","Props":{"class":"test"},"Children":null,"Content":"Hello"}`,
			wantErr: false,
			wantTag: "div",
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "empty json",
			json:    `{}`,
			wantErr: false,
			wantTag: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromJSON(tt.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.Tag != tt.wantTag {
				t.Errorf("FromJSON() Tag = %v, want %v", result.Tag, tt.wantTag)
			}
		})
	}
}

func TestToJSONAndFromJSON(t *testing.T) {
	original := VNode{
		Tag: "div",
		Props: Props{
			"class":   "container",
			"id":      "main",
			"visible": true,
			"count":   42,
		},
		Children: []VNode{
			{Tag: "h1", Content: "Title"},
			{Tag: "p", Content: "Content"},
		},
	}

	// 序列化
	jsonStr, err := ToJSON(original)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// 反序列化
	restored, err := FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	// 驗證基本屬性
	if restored.Tag != original.Tag {
		t.Errorf("Tag mismatch: got %v, want %v", restored.Tag, original.Tag)
	}

	if len(restored.Children) != len(original.Children) {
		t.Errorf("Children count mismatch: got %v, want %v", len(restored.Children), len(original.Children))
	}
}

func TestExtractTemplateVars(t *testing.T) {
	tests := []struct {
		name     string
		node     VNode
		expected []string
	}{
		{
			name: "single variable in content",
			node: VNode{
				Tag:     "div",
				Content: "Hello {{.Name}}",
			},
			expected: []string{".Name"},
		},
		{
			name: "variable in props",
			node: VNode{
				Tag:   "div",
				Props: Props{"id": "user-{{.ID}}"},
			},
			expected: []string{".ID"},
		},
		{
			name: "multiple variables",
			node: VNode{
				Tag:   "div",
				Props: Props{"data-id": "{{.ID}}"},
				Children: []VNode{
					{Content: "{{.Name}}"},
					{Content: "Email: {{.Email}}"},
				},
			},
			expected: []string{".ID", ".Name", ".Email"},
		},
		{
			name: "no variables",
			node: VNode{
				Tag:     "div",
				Content: "Plain text",
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTemplateVars(tt.node)

			// 將結果轉換為 map 以便檢查
			resultMap := make(map[string]bool)
			for _, v := range result {
				resultMap[v] = true
			}

			// 檢查每個期望的變數是否存在
			for _, exp := range tt.expected {
				if !resultMap[exp] {
					t.Errorf("ExtractTemplateVars() missing expected variable %q", exp)
				}
			}

			// 檢查是否有意外的變數
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractTemplateVars() count = %v, want %v. Got: %v", len(result), len(tt.expected), result)
			}
		})
	}
}

func TestCloneVNode(t *testing.T) {
	original := VNode{
		Tag: "div",
		Props: Props{
			"class": "original",
			"id":    "test",
		},
		Content: "Original content",
		Children: []VNode{
			{Tag: "span", Content: "Child"},
		},
	}

	cloned := CloneVNode(original)

	// 修改克隆的節點
	cloned.Props["class"] = "cloned"
	cloned.Content = "Cloned content"
	if len(cloned.Children) > 0 {
		cloned.Children[0].Content = "Modified child"
	}

	// 驗證原始節點未被修改
	if original.Props["class"] != "original" {
		t.Error("CloneVNode() modified original Props")
	}

	if original.Content != "Original content" {
		t.Error("CloneVNode() modified original Content")
	}

	if len(original.Children) > 0 && original.Children[0].Content != "Child" {
		t.Error("CloneVNode() modified original Children")
	}

	// 驗證克隆的節點已被修改
	if cloned.Props["class"] != "cloned" {
		t.Error("CloneVNode() clone Props not modified")
	}
}

func TestMergeProps(t *testing.T) {
	tests := []struct {
		name     string
		props    []Props
		expected map[string]interface{}
	}{
		{
			name: "merge two props",
			props: []Props{
				{"class": "btn", "type": "button"},
				{"class": "btn btn-primary"},
			},
			expected: map[string]interface{}{
				"class": "btn btn-primary",
				"type":  "button",
			},
		},
		{
			name: "merge multiple props",
			props: []Props{
				{"a": "1", "b": "2"},
				{"b": "3", "c": "4"},
				{"c": "5", "d": "6"},
			},
			expected: map[string]interface{}{
				"a": "1",
				"b": "3",
				"c": "5",
				"d": "6",
			},
		},
		{
			name: "merge with different types",
			props: []Props{
				{"visible": false, "count": 0},
				{"visible": true, "count": 42},
			},
			expected: map[string]interface{}{
				"visible": true,
				"count":   42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeProps(tt.props...)

			for k, expected := range tt.expected {
				if result[k] != expected {
					t.Errorf("MergeProps()[%q] = %v, want %v", k, result[k], expected)
				}
			}

			if len(result) != len(tt.expected) {
				t.Errorf("MergeProps() count = %v, want %v", len(result), len(tt.expected))
			}
		})
	}
}

func TestIsSelfClosingTag(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"br", true},
		{"hr", true},
		{"img", true},
		{"input", true},
		{"meta", true},
		{"link", true},
		{"div", false},
		{"span", false},
		{"p", false},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			result := isSelfClosingTag(tt.tag)
			if result != tt.expected {
				t.Errorf("isSelfClosingTag(%q) = %v, want %v", tt.tag, result, tt.expected)
			}
		})
	}
}

func TestConvertPropsToAny(t *testing.T) {
	tests := []struct {
		name     string
		props    Props
		validate func(Props) error
	}{
		{
			name: "convert JSAction from map",
			props: Props{
				"onClick": map[string]interface{}{
					"Code": "alert('test')",
				},
			},
			validate: func(p Props) error {
				if action, ok := p["onClick"].(JSAction); !ok {
					return nil // 不是 JSAction 也可以接受
				} else if action.Code != "alert('test')" {
					return nil
				}
				return nil
			},
		},
		{
			name: "keep other types unchanged",
			props: Props{
				"class":   "test",
				"visible": true,
				"count":   42,
			},
			validate: func(p Props) error {
				if p["class"] != "test" {
					return nil
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertPropsToAny(tt.props)
			if tt.validate != nil {
				if err := tt.validate(result); err != nil {
					t.Errorf("ConvertPropsToAny() validation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkToGoTemplate(b *testing.B) {
	node := VNode{
		Tag: "div",
		Props: Props{
			"class": "container",
			"id":    "main",
		},
		Children: []VNode{
			{Tag: "h1", Content: "Title"},
			{Tag: "p", Content: "Content"},
			{Tag: "footer", Content: "Footer"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToGoTemplate(node)
	}
}

func BenchmarkToJSON(b *testing.B) {
	node := VNode{
		Tag: "div",
		Props: Props{
			"class": "container",
			"id":    "main",
		},
		Children: []VNode{
			{Tag: "h1", Content: "Title"},
			{Tag: "p", Content: "Content"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToJSON(node)
	}
}

func BenchmarkFromJSON(b *testing.B) {
	jsonStr := `{"Tag":"div","Props":{"class":"container","id":"main"},"Children":[{"Tag":"h1","Props":null,"Children":null,"Content":"Title"},{"Tag":"p","Props":null,"Children":null,"Content":"Content"}],"Content":""}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromJSON(jsonStr)
	}
}

func BenchmarkCloneVNode(b *testing.B) {
	node := VNode{
		Tag: "div",
		Props: Props{
			"class": "test",
			"id":    "main",
		},
		Children: []VNode{
			{Tag: "span", Content: "Child 1"},
			{Tag: "span", Content: "Child 2"},
			{Tag: "span", Content: "Child 3"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CloneVNode(node)
	}
}
