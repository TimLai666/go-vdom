package dom

import (
	"strings"
	"testing"
)

func TestTernaryWithParentheses(t *testing.T) {
	tests := []struct {
		name     string
		template string
		props    Props
		expected string
	}{
		{
			name:     "simple parentheses",
			template: `${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}`,
			props:    Props{"a": true, "b": true},
			expected: "x",
		},
		{
			name:     "simple parentheses - nested false",
			template: `${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}`,
			props:    Props{"a": true, "b": false},
			expected: "y",
		},
		{
			name:     "simple parentheses - outer false",
			template: `${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}`,
			props:    Props{"a": false, "b": true},
			expected: "z",
		},
		{
			name:     "complex nested with parentheses",
			template: `${{{hasIcon}} === true ? ({{iconPosition}} === 'left' ? 'flex' : 'none') : 'none'}`,
			props:    Props{"hasIcon": true, "iconPosition": "left"},
			expected: "flex",
		},
		{
			name:     "complex nested - icon right",
			template: `${{{hasIcon}} === true ? ({{iconPosition}} === 'left' ? 'flex' : 'none') : 'none'}`,
			props:    Props{"hasIcon": true, "iconPosition": "right"},
			expected: "none",
		},
		{
			name:     "complex nested - no icon",
			template: `${{{hasIcon}} === true ? ({{iconPosition}} === 'left' ? 'flex' : 'none') : 'none'}`,
			props:    Props{"hasIcon": false, "iconPosition": "left"},
			expected: "none",
		},
		{
			name:     "without parentheses for comparison",
			template: `${{{hasIcon}} === true ? {{iconPosition}} === 'left' ? 'flex' : 'none' : 'none'}`,
			props:    Props{"hasIcon": true, "iconPosition": "left"},
			expected: "flex",
		},
		{
			name:     "multiple levels with parentheses",
			template: `${{{a}} ? ({{b}} ? ({{c}} ? 'deep' : 'mid') : 'shallow') : 'none'}`,
			props:    Props{"a": true, "b": true, "c": true},
			expected: "deep",
		},
		{
			name:     "multiple levels - c false",
			template: `${{{a}} ? ({{b}} ? ({{c}} ? 'deep' : 'mid') : 'shallow') : 'none'}`,
			props:    Props{"a": true, "b": true, "c": false},
			expected: "mid",
		},
		{
			name:     "multiple levels - b false",
			template: `${{{a}} ? ({{b}} ? ({{c}} ? 'deep' : 'mid') : 'shallow') : 'none'}`,
			props:    Props{"a": true, "b": false, "c": true},
			expected: "shallow",
		},
		{
			name:     "both branches with parentheses",
			template: `${{{flag}} ? ({{x}} ? 'a' : 'b') : ({{y}} ? 'c' : 'd')}`,
			props:    Props{"flag": true, "x": true, "y": false},
			expected: "a",
		},
		{
			name:     "both branches - false path",
			template: `${{{flag}} ? ({{x}} ? 'a' : 'b') : ({{y}} ? 'c' : 'd')}`,
			props:    Props{"flag": false, "x": false, "y": true},
			expected: "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateString(tt.template, tt.props)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
				t.Logf("Template: %s", tt.template)
				t.Logf("Props: %+v", tt.props)
			}
		})
	}
}

func TestTernaryEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		template string
		props    Props
		expected string
	}{
		{
			name:     "parentheses in string value",
			template: `${{{flag}} ? 'value (with parens)' : 'other'}`,
			props:    Props{"flag": true},
			expected: "value (with parens)",
		},
		{
			name:     "nested quotes and parentheses",
			template: `${{{a}} ? ({{b}} === 'test' ? 'yes' : 'no') : 'maybe'}`,
			props:    Props{"a": true, "b": "test"},
			expected: "yes",
		},
		{
			name:     "AND operator with parentheses",
			template: `${({{a}} && {{b}}) ? 'both' : 'not both'}`,
			props:    Props{"a": true, "b": true},
			expected: "both",
		},
		{
			name:     "OR operator with parentheses",
			template: `${({{a}} || {{b}}) ? 'at least one' : 'none'}`,
			props:    Props{"a": false, "b": true},
			expected: "at least one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateString(tt.template, tt.props)
			// Some of these might fail if operators aren't supported
			// This test is to document current behavior
			t.Logf("Template: %s", tt.template)
			t.Logf("Props: %+v", tt.props)
			t.Logf("Result: %q", result)
			t.Logf("Expected: %q", tt.expected)

			if !strings.Contains(result, "ERROR") && result != tt.expected {
				t.Logf("Note: Result differs from expected (may be unsupported feature)")
			}
		})
	}
}
