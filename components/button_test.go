package components

import (
	"strings"
	"testing"

	. "github.com/TimLai666/go-vdom/dom"
)

func TestButtonSize(t *testing.T) {
	tests := []struct {
		size     string
		fontSize string
		padding  string
	}{
		{"sm", "0.875rem", "0.375rem 1rem"},
		{"md", "0.95rem", "0.5rem 1.25rem"},
		{"lg", "1.125rem", "0.625rem 1.5rem"},
	}

	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			result := Btn(Props{"id": "test", "size": tt.size}, Text("按鈕"))
			html := Render(result)

			if !strings.Contains(html, "font-size: "+tt.fontSize) {
				t.Errorf("Button size=%s should have font-size: %s", tt.size, tt.fontSize)
			}

			if !strings.Contains(html, "padding: "+tt.padding) {
				t.Errorf("Button size=%s should have padding: %s", tt.size, tt.padding)
			}
		})
	}
}

func TestButtonVariant(t *testing.T) {
	// Test filled variant
	result1 := Btn(Props{"id": "test1", "variant": "filled", "color": "#3b82f6"}, Text("填充"))
	html1 := Render(result1)

	if !strings.Contains(html1, "background: #3b82f6") {
		t.Error("Button variant=filled should have colored background")
	}

	if !strings.Contains(html1, "color: #ffffff") {
		t.Error("Button variant=filled should have white text")
	}

	// Test outlined variant
	result2 := Btn(Props{"id": "test2", "variant": "outlined", "color": "#8b5cf6"}, Text("外框"))
	html2 := Render(result2)

	if !strings.Contains(html2, "background: transparent") {
		t.Error("Button variant=outlined should have transparent background")
	}

	if !strings.Contains(html2, "border: 1px solid #8b5cf6") {
		t.Error("Button variant=outlined should have colored border")
	}

	if !strings.Contains(html2, "color: #8b5cf6") {
		t.Error("Button variant=outlined should have colored text")
	}

	// Test text variant
	result3 := Btn(Props{"id": "test3", "variant": "text", "color": "#10b981"}, Text("文字"))
	html3 := Render(result3)

	if !strings.Contains(html3, "background: transparent") {
		t.Error("Button variant=text should have transparent background")
	}

	if !strings.Contains(html3, "color: #10b981") {
		t.Error("Button variant=text should have colored text")
	}
}

func TestButtonDisabled(t *testing.T) {
	// Test disabled=true
	result := Btn(Props{"id": "test", "disabled": true}, Text("禁用"))
	html := Render(result)

	if !strings.Contains(html, "cursor: not-allowed") {
		t.Error("Button with disabled=true should have cursor: not-allowed")
	}

	if !strings.Contains(html, "opacity: 0.6") {
		t.Error("Button with disabled=true should have opacity: 0.6")
	}

	// Test disabled=false (default)
	result2 := Btn(Props{"id": "test2", "disabled": false}, Text("啟用"))
	html2 := Render(result2)

	if !strings.Contains(html2, "cursor: pointer") {
		t.Error("Button with disabled=false should have cursor: pointer")
	}

	if !strings.Contains(html2, "opacity: 1") {
		t.Error("Button with disabled=false should have opacity: 1")
	}
}

func TestButtonFullWidth(t *testing.T) {
	// Test fullWidth=true
	result := Btn(Props{"id": "test", "fullWidth": true}, Text("全寬"))
	html := Render(result)

	if !strings.Contains(html, "width: 100%") {
		t.Error("Button with fullWidth=true should have width: 100%")
	}

	// Test fullWidth=false (default)
	result2 := Btn(Props{"id": "test2", "fullWidth": false}, Text("自動"))
	html2 := Render(result2)

	if !strings.Contains(html2, "width: auto") {
		t.Error("Button with fullWidth=false should have width: auto")
	}
}

func TestButtonRounded(t *testing.T) {
	tests := []struct {
		rounded string
		radius  string
	}{
		{"none", "0"},
		{"sm", "0.25rem"},
		{"md", "0.5rem"},
		{"lg", "0.75rem"},
		{"full", "9999px"},
	}

	for _, tt := range tests {
		t.Run(tt.rounded, func(t *testing.T) {
			result := Btn(Props{"id": "test", "rounded": tt.rounded}, Text("按鈕"))
			html := Render(result)

			if !strings.Contains(html, "border-radius: "+tt.radius) {
				t.Errorf("Button rounded=%s should have border-radius: %s", tt.rounded, tt.radius)
			}
		})
	}
}
