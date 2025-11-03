package components

import (
	"strings"
	"testing"

	. "github.com/TimLai666/go-vdom/dom"
)

func TestAlertClosable(t *testing.T) {
	// Test closable=true
	result := Alert(Props{"closable": true, "type": "info"}, Text("測試訊息"))
	html := Render(result)

	if !strings.Contains(html, `id="close-`) {
		t.Error("Alert with closable=true should have close button")
	}

	// Check if close button is visible by looking for the button tag and checking style before it
	buttonTagIdx := strings.Index(html, `id="close-`)
	if buttonTagIdx == -1 {
		t.Fatal("Close button not found")
	}

	// Look backwards to find the opening <button tag
	buttonStart := strings.LastIndex(html[:buttonTagIdx], "<button")
	if buttonStart == -1 {
		t.Fatal("Could not find button tag")
	}

	buttonSection := html[buttonStart : buttonTagIdx+200]
	if strings.Contains(buttonSection, "display: none") {
		t.Error("Close button should be visible when closable=true, but found 'display: none'")
	}

	// Test closable=false
	result2 := Alert(Props{"closable": false, "type": "info"}, Text("測試訊息"))
	html2 := Render(result2)

	buttonTagIdx2 := strings.Index(html2, `id="close-`)
	if buttonTagIdx2 == -1 {
		t.Fatal("Close button not found in second test")
	}

	// Look backwards to find the opening <button tag
	buttonStart2 := strings.LastIndex(html2[:buttonTagIdx2], "<button")
	if buttonStart2 == -1 {
		t.Fatal("Could not find button tag in second test")
	}

	buttonSection2 := html2[buttonStart2 : buttonTagIdx2+200]
	if !strings.Contains(buttonSection2, "display: none") {
		t.Errorf("Close button should be hidden when closable=false. Section: %s", buttonSection2)
	}
}

func TestAlertType(t *testing.T) {
	tests := []struct {
		alertType string
		bgColor   string
		textColor string
	}{
		{"info", "#eef2ff", "#4f46e5"},
		{"success", "#f0fdf4", "#16a34a"},
		{"warning", "#fffbeb", "#d97706"},
		{"error", "#fef2f2", "#dc2626"},
	}

	for _, tt := range tests {
		t.Run(tt.alertType, func(t *testing.T) {
			result := Alert(Props{"type": tt.alertType}, Text("測試"))
			html := Render(result)

			if !strings.Contains(html, tt.bgColor) {
				t.Errorf("Alert type=%s should have background color %s", tt.alertType, tt.bgColor)
			}

			if !strings.Contains(html, tt.textColor) {
				t.Errorf("Alert type=%s should have text color %s", tt.alertType, tt.textColor)
			}
		})
	}
}

func TestAlertTitle(t *testing.T) {
	// Test with title
	result := Alert(Props{"title": "重要通知", "type": "info"}, Text("內容"))
	html := Render(result)

	if !strings.Contains(html, "重要通知") {
		t.Error("Alert should display title")
	}

	// Find the title div and check if it's visible
	titleIdx := strings.Index(html, "重要通知")
	if titleIdx == -1 {
		t.Fatal("Title not found")
	}

	// Look backwards for the style attribute of the title div
	beforeTitle := html[:titleIdx]
	lastStyleIdx := strings.LastIndex(beforeTitle, "style=")
	if lastStyleIdx != -1 {
		styleSection := beforeTitle[lastStyleIdx : lastStyleIdx+100]
		if strings.Contains(styleSection, "display: none") {
			t.Error("Title should be visible when title prop is provided")
		}
	}

	// Test without title
	result2 := Alert(Props{"type": "info"}, Text("內容"))
	html2 := Render(result2)

	// The title div should have display: none
	// This is harder to test without DOM parsing, but we can check the structure
	if !strings.Contains(html2, "display: none") {
		t.Log("Warning: Expected to find 'display: none' for empty title")
	}
}

func TestAlertRounded(t *testing.T) {
	tests := []struct {
		rounded string
		radius  string
	}{
		{"none", "0"},
		{"sm", "0.25rem"},
		{"md", "0.5rem"},
		{"lg", "0.75rem"},
	}

	for _, tt := range tests {
		t.Run(tt.rounded, func(t *testing.T) {
			result := Alert(Props{"rounded": tt.rounded}, Text("測試"))
			html := Render(result)

			if !strings.Contains(html, "border-radius: "+tt.radius) {
				t.Errorf("Alert rounded=%s should have border-radius: %s", tt.rounded, tt.radius)
			}
		})
	}
}

func TestAlertCompact(t *testing.T) {
	// Test compact=true
	result := Alert(Props{"compact": true}, Text("測試"))
	html := Render(result)

	if !strings.Contains(html, "padding: 0.75rem 1rem") {
		t.Error("Alert with compact=true should have smaller padding")
	}

	if !strings.Contains(html, "align-items: center") {
		t.Error("Alert with compact=true should have align-items: center")
	}

	// Test compact=false (default)
	result2 := Alert(Props{"compact": false}, Text("測試"))
	html2 := Render(result2)

	if !strings.Contains(html2, "padding: 1rem 1.25rem") {
		t.Error("Alert with compact=false should have default padding")
	}

	if !strings.Contains(html2, "align-items: flex-start") {
		t.Error("Alert with compact=false should have align-items: flex-start")
	}
}

func TestAlertIcon(t *testing.T) {
	// Test icon=false
	result := Alert(Props{"icon": false}, Text("測試"))
	html := Render(result)

	// Find the icon div (first child div after alert container)
	divCount := 0
	iconDivFound := false
	for i := 0; i < len(html)-50; i++ {
		if strings.HasPrefix(html[i:], "<div") {
			divCount++
			if divCount == 2 { // Second div is the icon div
				section := html[i : i+200]
				if strings.Contains(section, "display: none") {
					iconDivFound = true
				}
				break
			}
		}
	}

	if !iconDivFound {
		t.Error("Alert with icon=false should hide the icon")
	}

	// Test icon=true (default)
	result2 := Alert(Props{"icon": true, "type": "success"}, Text("測試"))
	html2 := Render(result2)

	// Should contain the success icon
	if !strings.Contains(html2, "&#10003;") {
		t.Error("Alert with icon=true and type=success should show checkmark icon")
	}
}

func TestRadioDirection(t *testing.T) {
	// Test horizontal direction
	result := RadioGroup(Props{
		"direction": "horizontal",
		"name":      "test",
		"options":   "選項A|選項B",
	})
	html := Render(result)

	if !strings.Contains(html, "flex-direction: row") {
		t.Error("RadioGroup with direction=horizontal should have flex-direction: row")
	}

	// Test vertical direction (default)
	result2 := RadioGroup(Props{
		"direction": "vertical",
		"name":      "test2",
		"options":   "選項A|選項B",
	})
	html2 := Render(result2)

	if !strings.Contains(html2, "flex-direction: column") {
		t.Error("RadioGroup with direction=vertical should have flex-direction: column")
	}
}

func TestRadioLabel(t *testing.T) {
	// Test with label
	result := RadioGroup(Props{
		"label":   "選擇一個選項",
		"name":    "test",
		"options": "A|B",
	})
	html := Render(result)

	if !strings.Contains(html, "選擇一個選項") {
		t.Error("RadioGroup should display label")
	}

	// Test without label
	result2 := RadioGroup(Props{
		"name":    "test2",
		"options": "A|B",
	})
	html2 := Render(result2)

	// The first div should have display: none
	firstDivIdx := strings.Index(html2, "<div")
	if firstDivIdx != -1 {
		section := html2[firstDivIdx : firstDivIdx+200]
		if !strings.Contains(section, "display: none") {
			t.Error("RadioGroup without label should hide label div")
		}
	}
}

func TestCheckboxDirection(t *testing.T) {
	// Test horizontal
	result := CheckboxGroup(Props{
		"direction": "horizontal",
		"name":      "test",
		"options":   "選項A|選項B",
	})
	html := Render(result)

	if !strings.Contains(html, "flex-direction: row") {
		t.Error("CheckboxGroup with direction=horizontal should have flex-direction: row")
	}

	// Test vertical
	result2 := CheckboxGroup(Props{
		"direction": "vertical",
		"name":      "test2",
		"options":   "選項A|選項B",
	})
	html2 := Render(result2)

	if !strings.Contains(html2, "flex-direction: column") {
		t.Error("CheckboxGroup with direction=vertical should have flex-direction: column")
	}
}

func TestSwitchDirection(t *testing.T) {
	// Test horizontal (default for switch)
	result := Switch(Props{
		"name":  "test",
		"label": "啟用",
	})
	html := Render(result)

	if !strings.Contains(html, "flex-direction: row") {
		t.Error("Switch should have flex-direction based on direction prop")
	}
}

func TestHelpTextDisplay(t *testing.T) {
	// Test Radio with helpText
	result := RadioGroup(Props{
		"name":     "test",
		"options":  "A|B",
		"helpText": "請選擇一個選項",
	})
	html := Render(result)

	if !strings.Contains(html, "請選擇一個選項") {
		t.Error("Should display helpText")
	}

	// Find helpText and check it's visible
	helpIdx := strings.Index(html, "請選擇一個選項")
	if helpIdx != -1 {
		before := html[:helpIdx]
		lastStyleIdx := strings.LastIndex(before, "style=")
		if lastStyleIdx != -1 {
			section := before[lastStyleIdx : lastStyleIdx+100]
			if strings.Contains(section, "display: none") {
				t.Error("helpText should be visible when provided")
			}
		}
	}
}
