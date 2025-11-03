package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/TimLai666/go-vdom/components"
	. "github.com/TimLai666/go-vdom/dom"
)

func main() {
	http.HandleFunc("/", formsHandler)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func formsHandler(w http.ResponseWriter, r *http.Request) {
	page := Html(
		Props{},
		Head(
			Props{},
			Title(Props{}, Text("Form Components Demo")),
			Meta(Props{"charset": "UTF-8"}),
			Meta(Props{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
			Style(Props{}, Text(`
				body {
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
					background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
					margin: 0;
					padding: 20px;
					min-height: 100vh;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background: white;
					padding: 40px;
					border-radius: 16px;
					box-shadow: 0 20px 60px rgba(0,0,0,0.3);
				}
				h1 {
					color: #1a202c;
					margin-top: 0;
					margin-bottom: 10px;
					font-size: 2rem;
				}
				.subtitle {
					color: #718096;
					margin-bottom: 40px;
					font-size: 1rem;
				}
				.section {
					margin-bottom: 40px;
					padding-bottom: 40px;
					border-bottom: 1px solid #e2e8f0;
				}
				.section:last-child {
					border-bottom: none;
					padding-bottom: 0;
					margin-bottom: 0;
				}
				.section-title {
					font-size: 1.25rem;
					font-weight: 600;
					color: #2d3748;
					margin-bottom: 20px;
				}
				.demo-row {
					margin-bottom: 20px;
				}
			`)),
		),
		Body(
			Props{},
			Div(
				Props{"class": "container"},
				H1(Props{}, Text("🎨 Form Components Demo")),
				Div(Props{"class": "subtitle"}, Text("Showcasing all form component features")),

				// TextField Section
				Div(
					Props{"class": "section"},
					Div(Props{"class": "section-title"}, Text("📝 Text Fields")),
					Div(
						Props{"class": "demo-row"},
						TextField(Props{
							"id":          "email-input",
							"label":       "Email Address",
							"type":        "email",
							"placeholder": "your.email@example.com",
							"icon":        "📧",
							"helpText":    "We'll never share your email",
							"required":    "true",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						TextField(Props{
							"id":           "password-input",
							"label":        "Password",
							"type":         "password",
							"placeholder":  "Enter password",
							"icon":         "🔒",
							"iconPosition": "left",
							"variant":      "filled",
							"size":         "md",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						TextField(Props{
							"id":           "search-input",
							"label":        "Search",
							"placeholder":  "Search something...",
							"icon":         "🔍",
							"iconPosition": "right",
							"variant":      "outlined",
							"size":         "sm",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						TextField(Props{
							"id":        "error-input",
							"label":     "Username",
							"value":     "ab",
							"errorText": "Username must be at least 3 characters",
							"variant":   "underlined",
						}),
					),
				),

				// Dropdown Section
				Div(
					Props{"class": "section"},
					Div(Props{"class": "section-title"}, Text("🎯 Dropdowns")),
					Div(
						Props{"class": "demo-row"},
						Dropdown(Props{
							"id":          "country-select",
							"label":       "Country",
							"options":     "Taiwan,Japan,South Korea,USA,Canada,UK,Germany,France",
							"placeholder": "Select your country",
							"required":    "true",
							"helpText":    "Choose your country of residence",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						Dropdown(Props{
							"id":           "size-select",
							"label":        "T-Shirt Size",
							"options":      "XS,S,M,L,XL,XXL",
							"defaultValue": "M",
							"size":         "sm",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						Dropdown(Props{
							"id":      "color-select",
							"label":   "Favorite Color",
							"options": "Red,Green,Blue,Purple,Orange,Yellow",
							"color":   "#8b5cf6",
							"size":    "lg",
						}),
					),
				),

				// Switch Section
				Div(
					Props{"class": "section"},
					Div(Props{"class": "section-title"}, Text("🔘 Switches")),
					Div(
						Props{"class": "demo-row"},
						Switch(Props{
							"id":      "notifications",
							"label":   "Enable Notifications",
							"checked": "true",
							"onColor": "#10b981",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						Switch(Props{
							"id":       "dark-mode",
							"label":    "Dark Mode",
							"onColor":  "#6366f1",
							"offColor": "#94a3b8",
							"size":     "lg",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						Switch(Props{
							"id":       "marketing",
							"label":    "Marketing Emails",
							"helpText": "Receive updates about new features",
							"size":     "sm",
						}),
					),
					Div(
						Props{"class": "demo-row"},
						Switch(Props{
							"id":       "disabled-switch",
							"label":    "Disabled Switch",
							"disabled": true,
							"helpText": "This option is not available",
						}),
					),
				),

				// Button Section
				Div(
					Props{"class": "section"},
					Div(Props{"class": "section-title"}, Text("🔵 Buttons")),
					Div(
						Props{"style": "display: flex; gap: 1rem; flex-wrap: wrap;"},
						Btn(Props{
							"id":      "submit-btn",
							"variant": "filled",
							"color":   "#3b82f6",
							"size":    "md",
						}, Text("Submit")),
						Btn(Props{
							"id":      "cancel-btn",
							"variant": "outlined",
							"color":   "#ef4444",
							"size":    "md",
						}, Text("Cancel")),
						Btn(Props{
							"id":      "text-btn",
							"variant": "text",
							"color":   "#8b5cf6",
							"size":    "sm",
						}, Text("Learn More")),
						Btn(Props{
							"id":       "disabled-btn",
							"variant":  "filled",
							"color":    "#6b7280",
							"disabled": "true",
						}, Text("Disabled")),
					),
				),

				// Success message
				Alert(Props{
					"id":       "success-alert",
					"severity": "success",
					"title":    "All Components Working!",
					"closable": "true",
				}, Text("All form components have been successfully migrated and are functioning correctly. ✨")),
			),

			Script(Props{}, Text(`
				// Demo event listeners
				document.addEventListener('textfield:change', (e) => {
					console.log('TextField changed:', e.detail);
				});

				document.addEventListener('dropdown:change', (e) => {
					console.log('Dropdown changed:', e.detail);
				});

				document.addEventListener('switch:change', (e) => {
					console.log('Switch toggled:', e.detail);
				});

				console.log('✅ Form Components Demo loaded successfully!');
			`)),
		),
	)

	html := Render(page)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
