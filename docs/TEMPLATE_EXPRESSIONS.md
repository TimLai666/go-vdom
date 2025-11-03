# Template Expressions Guide

Complete guide to using expressions in go-vdom templates.

## Table of Contents

- [Basic Syntax](#basic-syntax)
- [Ternary Operators](#ternary-operators)
- [Logical Operators](#logical-operators)
- [Parentheses Support](#parentheses-support)
- [Comparisons](#comparisons)
- [Best Practices](#best-practices)
- [Examples](#examples)

---

## Basic Syntax

### Variable Interpolation

```go
// Simple variable
<div>{{name}}</div>

// In attributes
<div class="{{className}}">

// In JavaScript context (use ${...})
onClick: js.Do(
    js.Alert("'Hello ' + {{name}}"),
)
```

### Expression Interpolation

Use `${...}` for evaluated expressions:

```go
// Ternary operator
<div class="${{{isActive}} ? 'active' : 'inactive'}">

// With comparison
<div style="display: ${{{visible}} === true ? 'block' : 'none'}">
```

---

## Ternary Operators

### Basic Ternary

```go
${condition ? trueValue : falseValue}
```

**Example:**
```go
Template: `<div class="${{{isActive}} ? 'active' : 'inactive'}">`,
Props:    Props{"isActive": true},
Result:   <div class="active">
```

### Nested Ternary (Without Parentheses)

```go
${{{a}} ? {{b}} ? 'x' : 'y' : 'z'}
```

**Evaluation order:** `a ? (b ? 'x' : 'y') : 'z'`

**Example:**
```go
Template: `${{{hasIcon}} ? {{iconPosition}} === 'left' ? 'flex' : 'none' : 'none'}`,
Props:    Props{"hasIcon": true, "iconPosition": "left"},
Result:   flex
```

---

## Logical Operators

### AND Operator (`&&`)

```go
${{{a}} && {{b}} ? 'both true' : 'not both'}
```

**Example:**
```go
Template: `${{{enabled}} && {{visible}} ? 'show' : 'hide'}`,
Props:    Props{"enabled": true, "visible": true},
Result:   show
```

### OR Operator (`||`)

```go
${{{a}} || {{b}} ? 'at least one' : 'none'}
```

**Example:**
```go
Template: `${{{morning}} || {{evening}} ? 'greet' : 'ignore'}`,
Props:    Props{"morning": false, "evening": true},
Result:   greet
```

### Combined Operators

```go
${({{a}} && {{b}}) || {{c}} ? 'yes' : 'no'}
```

---

## Parentheses Support

**NEW in v1.3.0**: Full support for parentheses in template expressions!

### Why Use Parentheses?

1. **Clarity** - Makes complex logic easier to understand
2. **Control** - Explicit grouping of operations
3. **Nested Ternaries** - Cleaner nested conditional expressions

### Basic Grouping

```go
// Without parentheses (works but unclear)
${{{a}} ? {{b}} ? 'x' : 'y' : 'z'}

// With parentheses (clearer)
${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}
```

### Complex Nested Ternaries

```go
${{{hasIcon}} === true ? ({{iconPosition}} === 'left' ? 'flex' : 'none') : 'none'}
```

**Breakdown:**
1. Check if `hasIcon` is true
2. If yes → Check if `iconPosition` is 'left'
   - If yes → 'flex'
   - If no → 'none'
3. If no → 'none'

### Multiple Levels of Nesting

```go
${{{level1}} ? ({{level2}} ? ({{level3}} ? 'deep' : 'mid') : 'shallow') : 'none'}
```

### Both Branches with Parentheses

```go
${{{flag}} ? ({{x}} ? 'a' : 'b') : ({{y}} ? 'c' : 'd')}
```

**Evaluation:**
- If `flag` is true → Evaluate `(x ? 'a' : 'b')`
- If `flag` is false → Evaluate `(y ? 'c' : 'd')`

### Grouping Logical Operators

```go
// AND with parentheses
${({{a}} && {{b}}) ? 'both' : 'not both'}

// OR with parentheses
${({{a}} || {{b}}) ? 'one' : 'none'}

// Complex combination
${({{a}} && {{b}}) || ({{c}} && {{d}}) ? 'yes' : 'no'}
```

---

## Comparisons

### Equality (`===` or `==`)

```go
${{{status}} === 'active' ? 'on' : 'off'}
```

### Inequality (`!==` or `!=`)

```go
${{{role}} !== 'admin' ? 'limited' : 'full'}
```

### String Comparison

```go
${{{name}} === 'John' ? 'hello John' : 'hello stranger'}
```

### Boolean Comparison

```go
${{{isValid}} === true ? 'valid' : 'invalid'}

// Shorter form (recommended)
${{{isValid}} ? 'valid' : 'invalid'}
```

---

## Best Practices

### ✅ DO: Use Parentheses for Clarity

```go
// Good - clear intent
${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}

// Good - complex logic grouped
${({{enabled}} && {{visible}}) || {{forceShow}} ? 'show' : 'hide'}
```

### ✅ DO: Use Logical Operators

```go
// Good - concise and clear
${{{hasIcon}} && {{iconPosition}} === 'left' ? 'flex' : 'none'}

// Instead of
${{{hasIcon}} === true ? {{iconPosition}} === 'left' ? 'flex' : 'none' : 'none'}
```

### ✅ DO: Keep It Simple

```go
// Good - simple expression
${{{visible}} ? 'block' : 'none'}
```

### ⚠️ CONSIDER: Move Complex Logic to Go

```go
// If expression becomes too complex:
${({{a}} && {{b}}) ? ({{c}} || {{d}} ? ({{e}} ? 'x' : 'y') : 'z') : 'w'}

// Better: Handle in Go code
displayValue := "w"
if (a && b) {
    if (c || d) {
        if e {
            displayValue = "x"
        } else {
            displayValue = "y"
        }
    } else {
        displayValue = "z"
    }
}

Props{"displayValue": displayValue}

// Template
{{displayValue}}
```

### ❌ DON'T: Quote Template Variables

```go
// ❌ Wrong - quotes treat variable as string literal
${'{{flag}}' === 'true'}

// ✅ Correct - no quotes for variables
${{{flag}} === true}
```

### ❌ DON'T: Use String Concatenation in Expressions

```go
// ❌ Wrong - may not evaluate correctly
${'prefix-' + {{value}}}

// ✅ Correct - handle in Go
borderStyle := fmt.Sprintf("1px solid %s", color)
Props{"borderStyle": borderStyle}
```

---

## Examples

### Example 1: Icon Position

```go
Template: `
<div style="display: ${{{hasIcon}} && {{iconPosition}} === 'left' ? 'flex' : 'none'}">
    <i class="{{iconClass}}"></i>
</div>
`,

Props: Props{
    "hasIcon": true,
    "iconPosition": "left",
    "iconClass": "fa-home",
}

// Result: display: flex
```

### Example 2: Theme Selection

```go
Template: `
<div class="${{{theme}} === 'dark' ? 'bg-dark text-light' : 'bg-light text-dark'}">
`,

Props: Props{"theme": "dark"}

// Result: class="bg-dark text-light"
```

### Example 3: Multi-Level Conditional

```go
Template: `
<button class="${{{size}} === 'sm' ? 'btn-sm' : ({{size}} === 'lg' ? 'btn-lg' : 'btn-md')}">
`,

Props: Props{"size": "lg"}

// Result: class="btn-lg"
```

### Example 4: Form Validation State

```go
Template: `
<input
    class="${{{error}} ? 'border-red' : ({{success}} ? 'border-green' : 'border-gray')}"
    style="border-width: ${{{error}} || {{success}} ? '2px' : '1px'}"
>
`,

Props: Props{
    "error": false,
    "success": true,
}

// Result:
// class="border-green"
// style="border-width: 2px"
```

### Example 5: Complex Visibility Logic

```go
Template: `
<div style="display: ${
    ({{userRole}} === 'admin' || {{userRole}} === 'moderator') && {{isLoggedIn}}
        ? 'block'
        : 'none'
}">
    Admin Panel
</div>
`,

Props: Props{
    "userRole": "admin",
    "isLoggedIn": true,
}

// Result: display: block
```

### Example 6: Nested Conditions with Parentheses

```go
Template: `
<div class="${{{priority}} === 'high'
    ? ({{urgent}} ? 'alert-danger' : 'alert-warning')
    : ({{priority}} === 'low' ? 'alert-info' : 'alert-secondary')
}">
`,

Props: Props{
    "priority": "high",
    "urgent": true,
}

// Result: class="alert-danger"
```

---

## Supported Features Summary

| Feature | Supported | Example |
|---------|-----------|---------|
| Simple ternary | ✅ | `${{{a}} ? 'x' : 'y'}` |
| Nested ternary | ✅ | `${{{a}} ? {{b}} ? 'x' : 'y' : 'z'}` |
| Parentheses | ✅ | `${{{a}} ? ({{b}} ? 'x' : 'y') : 'z'}` |
| AND operator | ✅ | `${{{a}} && {{b}} ? 'yes' : 'no'}` |
| OR operator | ✅ | `${{{a}} \|\| {{b}} ? 'yes' : 'no'}` |
| Equality (`===`) | ✅ | `${{{x}} === 'test' ? 'a' : 'b'}` |
| Inequality (`!==`) | ✅ | `${{{x}} !== 'test' ? 'a' : 'b'}` |
| String literals | ✅ | `${'hello'}` |
| Nested parentheses | ✅ | `${(({{a}} && {{b}}) \|\| {{c}}) ? 'x' : 'y'}` |
| Multi-level nesting | ✅ | `${{{a}} ? ({{b}} ? ({{c}} ? 'd' : 'e') : 'f') : 'g'}` |

---

## Common Pitfalls

### 1. Forgetting to Use `${...}` for Expressions

```go
// ❌ Wrong - no evaluation
<div class="{{isActive}} ? 'active' : 'inactive'">

// ✅ Correct
<div class="${{{isActive}} ? 'active' : 'inactive'}">
```

### 2. Quoting Variables

```go
// ❌ Wrong - treats as string literal
${'{{name}}' === 'John'}

// ✅ Correct
${{{name}} === 'John'}
```

### 3. Missing Quotes for String Literals

```go
// ❌ Wrong - unquoted string
${{{status}} === active}

// ✅ Correct
${{{status}} === 'active'}
```

### 4. Parentheses in String Values

```go
// ✅ Correct - parentheses inside quotes are fine
${{{flag}} ? 'value (with parens)' : 'other'}
```

---

## Performance Considerations

- **Simple expressions** are evaluated at server-side render time
- **No runtime JavaScript overhead** - all evaluations happen in Go
- **Complex logic** should be moved to Go code for maintainability
- **Template compilation** is fast, but deeply nested expressions may impact readability

---

## Debugging Tips

### 1. Test Expressions Incrementally

```go
// Start simple
${{{a}}}

// Add condition
${{{a}} ? 'yes' : 'no'}

// Add complexity
${{{a}} ? ({{b}} ? 'yes' : 'maybe') : 'no'}
```

### 2. Use Template Linter

```bash
cd tools/template-linter
./template-linter -fix ../../
```

### 3. Check Prop Types

```go
// Boolean props should be bool, not string
Props{
    "enabled": true,      // ✅ Correct
    "enabled": "true",    // ❌ Wrong
}
```

### 4. Print Interpolated Results

```go
result := interpolateString(template, props)
fmt.Printf("Result: %q\n", result)
```

---

## Version History

- **v1.3.0** - Added full parentheses support
- **v1.2.0** - Added logical operators (`&&`, `||`)
- **v1.1.0** - Added nested ternary support
- **v1.0.0** - Initial template expression support

---

## Related Documentation

- [TEMPLATE_EXPRESSION_FIX_GUIDE.md](../TEMPLATE_EXPRESSION_FIX_GUIDE.md) - Migration guide
- [API_REFERENCE.md](API_REFERENCE.md) - Complete API documentation
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide

---

**Need help?** Check the [examples](../examples/) directory or create an issue on GitHub.
