package components

import (
	"strings"
	"testing"

	. "github.com/TimLai666/go-vdom/dom"
)

// TestSwitch tests the Switch component
func TestSwitch(t *testing.T) {
	tests := []struct {
		name     string
		props    Props
		expected []string
	}{
		{
			name: "basic switch",
			props: Props{
				"id":    "test-switch",
				"label": "Enable notifications",
			},
			expected: []string{
				`id="test-switch"`,
				`type="checkbox"`,
				"Enable notifications",
			},
		},
		{
			name: "checked switch",
			props: Props{
				"id":      "checked-switch",
				"checked": "true",
			},
			expected: []string{
				`id="checked-switch"`,
				`checked="true"`,
				"input.checked = 'true' === 'true'",
			},
		},
		{
			name: "disabled switch",
			props: Props{
				"id":       "disabled-switch",
				"disabled": "true",
			},
			expected: []string{
				`id="disabled-switch"`,
				`disabled="true"`,
				"input.disabled = 'true' === 'true'",
			},
		},
		{
			name: "custom colors",
			props: Props{
				"id":       "color-switch",
				"onColor":  "#10b981",
				"offColor": "#ef4444",
			},
			expected: []string{
				`id="color-switch"`,
				`data-on-color="#10b981"`,
				`data-off-color="#ef4444"`,
				"const onColor = track.getAttribute('data-on-color') || '#10b981'",
				"const offColor = track.getAttribute('data-off-color') || '#ef4444'",
			},
		},
		{
			name: "small size",
			props: Props{
				"id":   "small-switch",
				"size": "sm",
			},
			expected: []string{
				`id="small-switch"`,
				"width: 2.25rem",
				"height: 1.25rem",
				"const size = 'sm'",
				"const trackWidth = size === 'sm' ? '2.25rem'",
			},
		},
		{
			name: "large size",
			props: Props{
				"id":   "large-switch",
				"size": "lg",
			},
			expected: []string{
				`id="large-switch"`,
				"width: 3.25rem",
				"height: 1.75rem",
				"const size = 'lg'",
				"const trackHeight = size === 'sm' ? '1.25rem' : size === 'lg' ? '1.75rem' : '1.5rem'",
			},
		},
		{
			name: "with help text",
			props: Props{
				"id":       "help-switch",
				"helpText": "This will enable notifications",
			},
			expected: []string{
				`id="help-switch"`,
				"This will enable notifications",
				"display: block",
			},
		},
		{
			name: "left label position",
			props: Props{
				"id":            "left-label-switch",
				"label":         "Setting",
				"labelPosition": "left",
			},
			expected: []string{
				`id="left-label-switch"`,
				"Setting",
				"order: 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Switch(tt.props)
			html := Render(result)

			for _, exp := range tt.expected {
				if !strings.Contains(html, exp) {
					t.Errorf("Expected HTML to contain %q, but it didn't.\nGot: %s", exp, html)
				}
			}
		})
	}
}

// TestTextField tests the TextField (Input) component
func TestTextField(t *testing.T) {
	tests := []struct {
		name     string
		props    Props
		expected []string
	}{
		{
			name: "basic text input",
			props: Props{
				"id":          "test-input",
				"label":       "Email",
				"placeholder": "Enter your email",
			},
			expected: []string{
				`id="test-input"`,
				`type="text"`,
				`placeholder="Enter your email"`,
				"Email",
			},
		},
		{
			name: "email type",
			props: Props{
				"id":   "email-input",
				"type": "email",
			},
			expected: []string{
				`id="email-input"`,
				`type="email"`,
			},
		},
		{
			name: "required input",
			props: Props{
				"id":       "required-input",
				"required": "true",
			},
			expected: []string{
				`id="required-input"`,
				`required="true"`,
			},
		},
		{
			name: "disabled input",
			props: Props{
				"id":       "disabled-input",
				"disabled": "true",
			},
			expected: []string{
				`id="disabled-input"`,
				`disabled="true"`,
				"input.disabled = 'true' === 'true'",
			},
		},
		{
			name: "readonly input",
			props: Props{
				"id":       "readonly-input",
				"readonly": "true",
			},
			expected: []string{
				`id="readonly-input"`,
				`readonly="true"`,
				"input.readOnly = 'true' === 'true'",
			},
		},
		{
			name: "small size",
			props: Props{
				"id":   "small-input",
				"size": "sm",
			},
			expected: []string{
				`id="small-input"`,
				"padding: 0.5rem 0.75rem",
				"font-size: 0.875rem",
			},
		},
		{
			name: "large size",
			props: Props{
				"id":   "large-input",
				"size": "lg",
			},
			expected: []string{
				`id="large-input"`,
				"padding: 0.75rem 1rem",
				"font-size: 1rem",
			},
		},
		{
			name: "filled variant",
			props: Props{
				"id":      "filled-input",
				"variant": "filled",
			},
			expected: []string{
				`id="filled-input"`,
				"background: #f9fafb",
				"border: 1px solid transparent",
			},
		},
		{
			name: "underlined variant",
			props: Props{
				"id":      "underlined-input",
				"variant": "underlined",
			},
			expected: []string{
				`id="underlined-input"`,
				"border: none",
				"border-bottom: 1px solid #d1d5db",
				"border-radius: 0",
			},
		},
		{
			name: "with icon left",
			props: Props{
				"id":           "icon-input",
				"icon":         "📧",
				"iconPosition": "left",
			},
			expected: []string{
				`id="icon-input"`,
				"📧",
				"textfield-icon-left",
				"2.5rem",
			},
		},
		{
			name: "with icon right",
			props: Props{
				"id":           "icon-right-input",
				"icon":         "🔍",
				"iconPosition": "right",
			},
			expected: []string{
				`id="icon-right-input"`,
				"🔍",
				"textfield-icon-right",
				"2.5rem",
			},
		},
		{
			name: "with help text",
			props: Props{
				"id":       "help-input",
				"helpText": "Enter a valid email address",
			},
			expected: []string{
				`id="help-input"`,
				"Enter a valid email address",
				"color: #64748b",
			},
		},
		{
			name: "with error text",
			props: Props{
				"id":        "error-input",
				"errorText": "This field is required",
			},
			expected: []string{
				`id="error-input"`,
				"This field is required",
				"color: #ef4444",
			},
		},
		{
			name: "left label position",
			props: Props{
				"id":            "left-label-input",
				"label":         "Username",
				"labelPosition": "left",
			},
			expected: []string{
				`id="left-label-input"`,
				"Username",
				"display: flex",
				"width: 120px",
			},
		},
		{
			name: "custom color",
			props: Props{
				"id":    "custom-color-input",
				"color": "#8b5cf6",
			},
			expected: []string{
				`id="custom-color-input"`,
				`data-color="#8b5cf6"`,
			},
		},
		{
			name: "full width false",
			props: Props{
				"id":        "auto-width-input",
				"fullWidth": "false",
			},
			expected: []string{
				`id="auto-width-input"`,
				"width: auto",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TextField(tt.props)
			html := Render(result)

			for _, exp := range tt.expected {
				if !strings.Contains(html, exp) {
					t.Errorf("Expected HTML to contain %q, but it didn't.\nGot: %s", exp, html)
				}
			}
		})
	}
}

// TestDropdown tests the Dropdown component
func TestDropdown(t *testing.T) {
	tests := []struct {
		name     string
		props    Props
		expected []string
	}{
		{
			name: "basic dropdown",
			props: Props{
				"id":      "test-dropdown",
				"label":   "Select Country",
				"options": "Taiwan,Japan,USA",
			},
			expected: []string{
				`id="test-dropdown"`,
				"Select Country",
				"const options = 'Taiwan,Japan,USA'.split(',').filter(opt => opt.trim())",
			},
		},
		{
			name: "with default value",
			props: Props{
				"id":           "default-dropdown",
				"options":      "Red,Green,Blue",
				"defaultValue": "Green",
			},
			expected: []string{
				`id="default-dropdown"`,
				"const defaultValue = 'Green'",
				"opt.selected = true",
			},
		},
		{
			name: "custom placeholder",
			props: Props{
				"id":          "placeholder-dropdown",
				"options":     "A,B,C",
				"placeholder": "Choose an option",
			},
			expected: []string{
				`id="placeholder-dropdown"`,
				"Choose an option",
			},
		},
		{
			name: "required dropdown",
			props: Props{
				"id":       "required-dropdown",
				"options":  "Yes,No",
				"required": "true",
			},
			expected: []string{
				`id="required-dropdown"`,
				`required="true"`,
			},
		},
		{
			name: "disabled dropdown",
			props: Props{
				"id":       "disabled-dropdown",
				"options":  "One,Two",
				"disabled": "true",
			},
			expected: []string{
				`id="disabled-dropdown"`,
				`disabled="true"`,
				"select.disabled = true",
			},
		},
		{
			name: "small size",
			props: Props{
				"id":      "small-dropdown",
				"options": "A,B,C",
				"size":    "sm",
			},
			expected: []string{
				`id="small-dropdown"`,
				"padding: 0.5rem 2.5rem 0.5rem 0.75rem",
				"font-size: 0.875rem",
			},
		},
		{
			name: "large size",
			props: Props{
				"id":      "large-dropdown",
				"options": "X,Y,Z",
				"size":    "lg",
			},
			expected: []string{
				`id="large-dropdown"`,
				"padding: 0.75rem 2.75rem 0.75rem 1rem",
				"font-size: 1rem",
			},
		},
		{
			name: "with help text",
			props: Props{
				"id":       "help-dropdown",
				"options":  "Apple,Banana",
				"helpText": "Select your favorite fruit",
			},
			expected: []string{
				`id="help-dropdown"`,
				"Select your favorite fruit",
				"color: #64748b",
			},
		},
		{
			name: "with error text",
			props: Props{
				"id":        "error-dropdown",
				"options":   "Cat,Dog",
				"errorText": "Please select an option",
			},
			expected: []string{
				`id="error-dropdown"`,
				"Please select an option",
				"color: #ef4444",
			},
		},
		{
			name: "left label position",
			props: Props{
				"id":            "left-label-dropdown",
				"label":         "Category",
				"options":       "Tech,Sports,Music",
				"labelPosition": "left",
			},
			expected: []string{
				`id="left-label-dropdown"`,
				"Category",
				"display: flex",
				"width: 120px",
			},
		},
		{
			name: "custom color",
			props: Props{
				"id":      "custom-color-dropdown",
				"options": "One,Two,Three",
				"color":   "#10b981",
			},
			expected: []string{
				`id="custom-color-dropdown"`,
				`data-color="#10b981"`,
			},
		},
		{
			name: "full width false",
			props: Props{
				"id":        "auto-width-dropdown",
				"options":   "Yes,No",
				"fullWidth": "false",
			},
			expected: []string{
				`id="auto-width-dropdown"`,
				"width: auto",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Dropdown(tt.props)
			html := Render(result)

			for _, exp := range tt.expected {
				if !strings.Contains(html, exp) {
					t.Errorf("Expected HTML to contain %q, but it didn't.\nGot: %s", exp, html)
				}
			}
		})
	}
}
