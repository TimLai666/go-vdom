# Development Guide

This guide covers the development workflow, tools, and best practices for contributing to go-vdom.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Tools](#development-tools)
- [Testing](#testing)
- [Template Linting](#template-linting)
- [Code Quality](#code-quality)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setup

1. **Clone the repository**

```bash
git clone https://github.com/TimLai666/go-vdom.git
cd go-vdom
```

2. **Install dependencies**

```bash
go mod download
```

3. **Install development tools**

```bash
make install-tools
```

4. **Set up git hooks**

```bash
chmod +x .githooks/install.sh
./.githooks/install.sh
```

5. **Verify setup**

```bash
make check
```

## Development Tools

### Template Linter

The template linter is a critical tool that detects dangerous patterns in component templates.

#### Building the Linter

```bash
cd tools/template-linter
go build
```

#### Running the Linter

```bash
# Check entire project
./template-linter ../../

# Check specific directory
./template-linter ../../components

# Show fix suggestions
./template-linter -fix ../../

# Verbose output
./template-linter -v ../../
```

#### What It Detects

1. **Quoted Variables in Expressions**
   - ❌ `${'{{variable}}'}`
   - ✅ `${{{variable}}}`

2. **String Boolean Comparisons**
   - ❌ `'{{flag}}' === 'true'`
   - ✅ `{{flag}} === true`

3. **String Concatenation in Expressions**
   - ❌ `${'text' + {{variable}}}`
   - ✅ Handle in Go code

4. **Quoted Variable Comparisons**
   - ❌ `'{{value}}' === 'something'`
   - ✅ `{{value}} === 'something'`

5. **Double-Quoted Variables**
   - Warning: Verify context (HTML vs JS)

#### Adding New Rules

To add a new linting rule:

1. Define a regex pattern in `tools/template-linter/main.go`:

```go
var myNewPattern = regexp.MustCompile(`pattern_here`)
```

2. Add detection in `checkLine()`:

```go
if match := myNewPattern.FindString(line); match != "" {
    col := strings.Index(line, match) + 1
    issues = append(issues, Issue{
        File:    filePath,
        Line:    lineNum,
        Column:  col,
        Type:    "my-new-issue",
        Message: "Description of the problem",
        Context: match,
    })
}
```

3. Add fix suggestion in `suggestFix()`:

```go
case "my-new-issue":
    return "How to fix this issue"
```

### Makefile Commands

The project includes a Makefile with common development tasks:

```bash
# Testing
make test              # Run all tests
make test-coverage     # Generate coverage report

# Linting
make lint              # Run template linter
make fmt               # Format code
make vet               # Run go vet

# Building
make build-linter      # Build template linter
make build-examples    # Build all examples

# Running
make run               # Run main demo
make run-example EXAMPLE=forms_demo  # Run specific example

# Quality checks
make check             # Run all checks (fmt, vet, lint, test)

# Cleanup
make clean             # Remove build artifacts

# Help
make help              # Show all available commands
```

### Git Hooks

The pre-commit hook ensures code quality before commits:

**What it checks:**
- Code formatting (`go fmt`)
- Static analysis (`go vet`)
- Template linting

**Bypass (not recommended):**
```bash
git commit --no-verify
```

**Manual installation:**
```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**Configure git to use .githooks directory:**
```bash
git config core.hooksPath .githooks
```

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./dom
go test ./components

# With verbose output
go test -v ./dom

# With coverage
go test -cover ./...

# With race detection
go test -race ./...

# Using Makefile
make test
make test-coverage
```

### Writing Tests

#### Component Tests

```go
func TestMyComponent(t *testing.T) {
    comp := MyComponent(Props{
        "title": "Test Title",
        "visible": true,
    })

    html := Render(comp)

    // Check expected output
    if !strings.Contains(html, "Test Title") {
        t.Errorf("Expected title not found in output")
    }
}
```

#### Template Interpolation Tests

```go
func TestTemplateInterpolation(t *testing.T) {
    template := `<div>{{message}}</div>`
    props := Props{"message": "Hello"}

    result := interpolateString(template, props)

    expected := `<div>Hello</div>`
    if result != expected {
        t.Errorf("Expected %q, got %q", expected, result)
    }
}
```

#### JavaScript Generation Tests

```go
func TestJavaScriptGeneration(t *testing.T) {
    action := js.Do(
        js.Const("x", "1"),
        js.Log("x"),
    )

    code := action.Generate()

    if !strings.Contains(code, "const x=1") {
        t.Errorf("Expected constant declaration in generated code")
    }
}
```

### Test Organization

```
package_test.go          # External black-box tests
package_internal_test.go # Internal white-box tests
testdata/                # Test fixtures
```

## Template Linting

### Common Issues and Fixes

#### Issue 1: Quoted Variable in Expression

**Problem:**
```go
Template: `
    <div data-value="${'{{myValue}}'}">
`,
```

**Why it's wrong:**
The variable is treated as a string literal, not interpolated.

**Fix:**
```go
Template: `
    <div data-value="${{{myValue}}}">
`,
```

#### Issue 2: String Boolean Comparison

**Problem:**
```go
Template: `
    <div data-active="${'{{isActive}}' === 'true' ? 'yes' : 'no'}">
`,
```

**Why it's wrong:**
Compares strings, not booleans. Won't work correctly.

**Fix:**
```go
Template: `
    <div data-active="${{{isActive}} === true ? 'yes' : 'no'}">
`,
```

And pass boolean in Go:
```go
Props{"isActive": true}  // not "true"
```

#### Issue 3: String Concatenation in Expression

**Problem:**
```go
Template: `
    <div style="${'1px solid ' + {{color}}}">
`,
```

**Why it's wrong:**
The expression evaluator may not handle this correctly.

**Fix:**
Handle in Go code:
```go
// In component function
borderStyle := fmt.Sprintf("1px solid %s", color)

// In template
Template: `
    <div style="{{borderStyle}}">
`,
```

### Best Practices for Templates

1. **Use correct prop types**
   ```go
   // ✅ Good
   Props{
       "enabled": true,
       "count": 42,
       "name": "John",
   }

   // ❌ Bad
   Props{
       "enabled": "true",
       "count": "42",
       "name": "John",
   }
   ```

2. **Boolean comparisons**
   ```go
   // ✅ Good
   ${{{flag}} === true}

   // ❌ Bad
   ${'{{flag}}' === 'true'}
   ```

3. **String operations**
   ```go
   // ✅ Good - in Go
   value := strings.TrimSpace(input)
   Props{"value": value}

   // ❌ Bad - in template
   ${{'{{value}}'.trim()}
   ```

4. **Complex calculations**
   ```go
   // ✅ Good - in Go
   result := calculate(a, b, c)
   Props{"result": result}

   // ❌ Bad - in template
   ${{{a}} + {{b}} * {{c}}}
   ```

## Code Quality

### Code Formatting

Always format code before committing:

```bash
go fmt ./...
# or
make fmt
```

### Static Analysis

Run `go vet` to catch common mistakes:

```bash
go vet ./...
# or
make vet
```

### Linting

Use `golangci-lint` for comprehensive linting:

```bash
golangci-lint run
```

Configuration in `.golangci.yml` (create if needed):

```yaml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode

linters-settings:
  gofmt:
    simplify: true
```

### Complete Check

Run all quality checks:

```bash
make check
```

This runs:
1. `go fmt`
2. `go vet`
3. Template linter
4. All tests

## Common Patterns

### Creating a New Component

```go
// 1. Define the component template
func MyComponent(props Props, children ...VNode) VNode {
    return Component(
        // Template
        Div(Props{"class": "my-component {{className}}"},
            H2("{{title}}"),
            Div(Props{"class": "content"},
                Text("{{content}}"),
            ),
            Div(Props{"class": "children"},
                children...,
            ),
        ),
        children,
        // Default props
        PropsDef{
            "title": "",
            "content": "",
            "className": "",
            "visible": true,
        },
    )(props)
}

// 2. Add tests
func TestMyComponent(t *testing.T) {
    comp := MyComponent(Props{
        "title": "Test",
        "content": "Hello",
    })

    html := Render(comp)

    if !strings.Contains(html, "Test") {
        t.Error("Title not found")
    }
}

// 3. Document usage in component file
/*
Usage:

    MyComponent(Props{
        "title": "My Title",
        "content": "My content",
        "className": "custom-class",
    },
        P("Child content"),
    )
*/
```

### Adding JavaScript Functionality

```go
// Simple click handler
Button(Props{
    "onClick": js.Do(
        js.Alert("'Button clicked!'"),
    ),
}, "Click Me")

// Async API call
Button(Props{
    "onClick": js.AsyncDo(
        js.Try(
            js.Const("response", "await fetch('/api/data')"),
            js.Const("data", "await response.json()"),
            js.Log("data"),
        ).Catch(
            js.Log("'Error: ' + error.message"),
        ).End(),
    ),
}, "Load Data")

// With parameters
Button(Props{
    "onClick": js.Do([]string{"event"},
        js.Call("event.preventDefault"),
        js.Alert("'Form submitted!'"),
    ),
}, "Submit")
```

### Complex Props

```go
// Array
Props{
    "items": []string{"apple", "banana", "orange"},
}

// Map
Props{
    "config": map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
    },
}

// Struct
type User struct {
    Name  string
    Email string
}

Props{
    "user": User{Name: "John", Email: "john@example.com"},
}
```

## Troubleshooting

### Template Variables Not Interpolating

**Symptom:** `{{variable}}` appears literally in output

**Causes:**
1. Variable not in props
2. Wrong prop name
3. Template not being processed

**Solution:**
```go
// Ensure prop is passed
comp := MyComponent(Props{
    "variable": "value",  // Must match template
})

// Check template uses correct name
Template: `<div>{{variable}}</div>`
```

### JavaScript Errors

**Symptom:** Browser console shows JavaScript errors

**Common causes:**

1. **Missing quotes in strings**
   ```go
   // ❌ Wrong
   js.Log("Hello")

   // ✅ Correct
   js.Log("'Hello'")
   ```

2. **Using await without AsyncFn/AsyncDo**
   ```go
   // ❌ Wrong
   js.Do(js.Const("data", "await fetch('/api')"))

   // ✅ Correct
   js.AsyncDo(js.Const("data", "await fetch('/api')"))
   ```

3. **Missing event parameter**
   ```go
   // ❌ Wrong - using event without declaring
   js.Do(js.Call("event.preventDefault"))

   // ✅ Correct
   js.Do([]string{"event"}, js.Call("event.preventDefault"))
   ```

### Build Errors in Examples

**Symptom:** `main redeclared` error

**Cause:** Each example file has its own `main()` function

**Solution:**
```bash
# Don't build all examples at once
go build ./examples  # ❌ Wrong

# Build individually
go build examples/forms_demo.go  # ✅ Correct

# Or use make
make build-examples  # ✅ Correct
```

### Template Linter Fails

**Symptom:** Linter reports issues

**Solution:**
1. Read the error message carefully
2. Use `-fix` flag for suggestions: `./template-linter -fix ../../`
3. Refer to [TEMPLATE_EXPRESSION_FIX_GUIDE.md](../TEMPLATE_EXPRESSION_FIX_GUIDE.md)
4. Ask for help if unclear

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- `MAJOR`: Incompatible API changes
- `MINOR`: New functionality (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Release Checklist

1. **Update version**
   - Update `CHANGELOG.md`
   - Update version in documentation

2. **Run all checks**
   ```bash
   make check
   make test-coverage
   ```

3. **Update documentation**
   - Update API reference if needed
   - Update examples
   - Update README if needed

4. **Create git tag**
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```

5. **Create GitHub release**
   - Go to GitHub releases page
   - Create new release from tag
   - Add release notes from CHANGELOG

6. **Verify**
   - Check that CI passes
   - Verify package can be installed: `go get github.com/TimLai666/go-vdom@v1.2.0`

## Contributing

### Workflow

1. **Create feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes**
   - Write code
   - Write tests
   - Update documentation

3. **Run checks**
   ```bash
   make check
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "Add my feature"
   ```
   (Pre-commit hook runs automatically)

5. **Push and create PR**
   ```bash
   git push origin feature/my-feature
   ```

### Commit Messages

Follow conventional commits:

```
feat: Add new component
fix: Fix template interpolation bug
docs: Update API reference
test: Add tests for Button component
refactor: Simplify render logic
chore: Update dependencies
```

### Pull Request Guidelines

- Clear description of changes
- Link to related issues
- All checks passing
- Up-to-date with main branch
- Reviewed and approved

## Resources

- [Quick Start Guide](QUICK_START.md)
- [API Reference](API_REFERENCE.md)
- [Template Expression Fix Guide](../TEMPLATE_EXPRESSION_FIX_GUIDE.md)
- [JSON Serialization Changelog](../CHANGELOG_JSON_SERIALIZATION.md)

## Getting Help

- Check documentation
- Search existing issues
- Create new issue with:
  - Clear description
  - Minimal reproduction
  - Expected vs actual behavior
  - Environment details

---

**Happy coding! 🚀**
