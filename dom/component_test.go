// component_test.go
package dom

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSerializeComplexType(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "int value",
			input:    42,
			expected: "42",
		},
		{
			name:     "float value",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "string slice",
			input:    []string{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
		{
			name:     "int slice",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "map[string]string",
			input:    map[string]string{"key": "value", "foo": "bar"},
			expected: `{"foo":"bar","key":"value"}`, // JSON marshal sorts keys
		},
		{
			name:     "map[string]interface{}",
			input:    map[string]interface{}{"name": "John", "age": 30, "active": true},
			expected: `{"active":true,"age":30,"name":"John"}`,
		},
		{
			name: "struct value",
			input: struct {
				Name string
				Age  int
			}{"Alice", 25},
			expected: `{"Name":"Alice","Age":25}`,
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: `[]`,
		},
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := serializeComplexType(tt.input)
			if result != tt.expected {
				t.Errorf("serializeComplexType() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestComponentWithComplexProps(t *testing.T) {
	tests := []struct {
		name      string
		template  VNode
		props     Props
		checkFunc func(*testing.T, VNode)
	}{
		{
			name: "array prop as data attribute",
			template: VNode{
				Tag:   "div",
				Props: Props{"data-items": "{{items}}"},
			},
			props: Props{"items": []string{"apple", "banana", "orange"}},
			checkFunc: func(t *testing.T, result VNode) {
				dataItems, ok := result.Props["data-items"].(string)
				if !ok {
					t.Fatalf("data-items is not a string")
				}
				expected := `["apple","banana","orange"]`
				if dataItems != expected {
					t.Errorf("data-items = %q, want %q", dataItems, expected)
				}
			},
		},
		{
			name: "map prop as data attribute",
			template: VNode{
				Tag:   "div",
				Props: Props{"data-config": "{{config}}"},
			},
			props: Props{"config": map[string]interface{}{"theme": "dark", "fontSize": 14}},
			checkFunc: func(t *testing.T, result VNode) {
				dataConfig, ok := result.Props["data-config"].(string)
				if !ok {
					t.Fatalf("data-config is not a string")
				}
				// Parse JSON to verify it's valid
				var config map[string]interface{}
				if err := json.Unmarshal([]byte(dataConfig), &config); err != nil {
					t.Errorf("data-config is not valid JSON: %v", err)
				}
				if config["theme"] != "dark" {
					t.Errorf("config.theme = %v, want 'dark'", config["theme"])
				}
			},
		},
		{
			name: "struct prop",
			template: VNode{
				Tag:   "div",
				Props: Props{"data-user": "{{user}}"},
			},
			props: Props{"user": struct {
				Name  string
				Email string
			}{"John Doe", "john@example.com"}},
			checkFunc: func(t *testing.T, result VNode) {
				dataUser, ok := result.Props["data-user"].(string)
				if !ok {
					t.Fatalf("data-user is not a string")
				}
				var user struct {
					Name  string
					Email string
				}
				if err := json.Unmarshal([]byte(dataUser), &user); err != nil {
					t.Errorf("data-user is not valid JSON: %v", err)
				}
				if user.Name != "John Doe" {
					t.Errorf("user.Name = %q, want 'John Doe'", user.Name)
				}
			},
		},
		{
			name: "mixed simple and complex props",
			template: VNode{
				Tag: "div",
				Props: Props{
					"class":       "{{className}}",
					"data-items":  "{{items}}",
					"data-count":  "{{count}}",
					"data-active": "{{active}}",
				},
			},
			props: Props{
				"className": "container",
				"items":     []int{1, 2, 3, 4, 5},
				"count":     5,
				"active":    true,
			},
			checkFunc: func(t *testing.T, result VNode) {
				if result.Props["class"] != "container" {
					t.Errorf("class = %q, want 'container'", result.Props["class"])
				}
				if result.Props["data-items"] != "[1,2,3,4,5]" {
					t.Errorf("data-items = %q, want '[1,2,3,4,5]'", result.Props["data-items"])
				}
				if result.Props["data-count"] != "5" {
					t.Errorf("data-count = %q, want '5'", result.Props["data-count"])
				}
				if result.Props["data-active"] != "true" {
					t.Errorf("data-active = %q, want 'true'", result.Props["data-active"])
				}
			},
		},
		{
			name: "nested array and map",
			template: VNode{
				Tag:   "div",
				Props: Props{"data-config": "{{config}}"},
			},
			props: Props{"config": map[string]interface{}{
				"users": []map[string]string{
					{"name": "Alice", "role": "admin"},
					{"name": "Bob", "role": "user"},
				},
				"settings": map[string]bool{
					"darkMode":      true,
					"notifications": false,
				},
			}},
			checkFunc: func(t *testing.T, result VNode) {
				dataConfig, ok := result.Props["data-config"].(string)
				if !ok {
					t.Fatalf("data-config is not a string")
				}
				var config map[string]interface{}
				if err := json.Unmarshal([]byte(dataConfig), &config); err != nil {
					t.Errorf("data-config is not valid JSON: %v", err)
				}
				users, ok := config["users"].([]interface{})
				if !ok || len(users) != 2 {
					t.Errorf("config.users is invalid")
				}
			},
		},
		{
			name: "array prop preserved as data attribute",
			template: VNode{
				Tag: "div",
				Props: Props{
					"data-options": "{{options}}",
				},
			},
			props: Props{
				"options": []string{},
			},
			checkFunc: func(t *testing.T, result VNode) {
				if result.Props["data-options"] != "[]" {
					t.Errorf("data-options = %q, want '[]'", result.Props["data-options"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			componentFn := Component(tt.template, nil)
			result := componentFn(tt.props)
			tt.checkFunc(t, result)
		})
	}
}

func TestComponentJSONInTemplate(t *testing.T) {
	// Test that complex props can be used in templates with JSON serialization
	template := VNode{
		Tag: "div",
		Props: Props{
			"data-config": "{{config}}",
		},
		Children: []VNode{
			{Tag: "span", Props: Props{"data-items": "{{items}}"}},
		},
	}

	componentFn := Component(template, nil)
	result := componentFn(Props{
		"config": map[string]string{"theme": "dark"},
		"items":  []int{1, 2, 3},
	})

	dataConfig := result.Props["data-config"].(string)
	if !strings.Contains(dataConfig, "theme") {
		t.Errorf("data-config should contain theme, got: %s", dataConfig)
	}

	dataItems := result.Children[0].Props["data-items"].(string)
	if dataItems != "[1,2,3]" {
		t.Errorf("data-items = %q, want '[1,2,3]'", dataItems)
	}
}

func TestComponentWithPointerProps(t *testing.T) {
	str := "test"
	num := 42

	template := VNode{
		Tag: "div",
		Props: Props{
			"data-str": "{{str}}",
			"data-num": "{{num}}",
		},
	}

	componentFn := Component(template, nil)
	result := componentFn(Props{
		"str": &str,
		"num": &num,
	})

	if result.Props["data-str"] != "test" {
		t.Errorf("data-str = %q, want 'test'", result.Props["data-str"])
	}
	if result.Props["data-num"] != "42" {
		t.Errorf("data-num = %q, want '42'", result.Props["data-num"])
	}
}

func TestInterpolateWithComplexTypes(t *testing.T) {
	tests := []struct {
		name     string
		template string
		props    Props
		expected string
	}{
		{
			name:     "array in string",
			template: "data: {{items}}",
			props:    Props{"items": []string{"a", "b"}},
			expected: `data: ["a","b"]`,
		},
		{
			name:     "map in string",
			template: "config: {{config}}",
			props:    Props{"config": map[string]int{"x": 1, "y": 2}},
			expected: `config: {"x":1,"y":2}`,
		},
		{
			name:     "mixed types",
			template: "name: {{name}}, items: {{items}}, count: {{count}}",
			props: Props{
				"name":  "test",
				"items": []int{1, 2, 3},
				"count": 3,
			},
			expected: "name: test, items: [1,2,3], count: 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateString(tt.template, tt.props)
			if result != tt.expected {
				t.Errorf("interpolateString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func BenchmarkSerializeComplexType(b *testing.B) {
	testData := []interface{}{
		[]string{"a", "b", "c", "d", "e"},
		map[string]interface{}{"name": "John", "age": 30, "active": true},
		struct {
			Name  string
			Email string
			Age   int
		}{"Alice", "alice@example.com", 25},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, data := range testData {
			_ = serializeComplexType(data)
		}
	}
}

func BenchmarkComponentWithComplexProps(b *testing.B) {
	template := VNode{
		Tag: "div",
		Props: Props{
			"data-items":  "{{items}}",
			"data-config": "{{config}}",
		},
	}

	componentFn := Component(template, nil)
	props := Props{
		"items":  []string{"a", "b", "c", "d", "e"},
		"config": map[string]interface{}{"theme": "dark", "fontSize": 14},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = componentFn(props)
	}
}
