# Template Linter

A static analysis tool for detecting dangerous patterns in go-vdom component templates.

## Purpose

This linter helps prevent common mistakes when working with template variables and expressions in go-vdom components, especially after the JSON serialization migration.

## What It Detects

### 1. Quoted Variables in Expressions
**Problem:** `${'{{variable}}'}` treats the variable as a string literal
**Fix:** Use `${{{variable}}}` to properly interpolate the value

### 2. String Boolean Comparisons
**Problem:** `'{{flag}}' === 'true'` compares strings instead of booleans
**Fix:** Use `{{flag}} === true` for proper boolean comparison

### 3. String Concatenation in Expressions
**Problem:** `${'text' + {{variable}}}` may not evaluate correctly
**Fix:** Handle string concatenation in Go code and pass the result as a prop

### 4. Quoted Variable Comparisons
**Problem:** `'{{value}}' === 'something'` treats everything as strings
**Fix:** Use `{{value}} === 'something'` for proper type comparison

### 5. Double-Quoted Variables
**Warning:** `"{{variable}}"` may be incorrect depending on context
**Note:** This is valid in HTML attributes but suspicious in JS expressions

## Installation

```bash
cd tools/template-linter
go build
```

## Usage

### Basic Usage
```bash
# Check current directory
./template-linter

# Check specific directory
./template-linter ../../components

# Check with verbose output
./template-linter -v

# Show fix suggestions
./template-linter -fix
```

### In CI/CD

Add to your CI pipeline (e.g., GitHub Actions):

```yaml
- name: Run Template Linter
  run: |
    cd tools/template-linter
    go build
    ./template-linter ../../
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/sh
cd tools/template-linter
go run main.go ../../
if [ $? -ne 0 ]; then
    echo "Template linter found issues. Please fix before committing."
    exit 1
fi
```

## Examples

### ❌ Bad
```go
Template: `
    <div data-active="${'{{isActive}}' === 'true' ? 'yes' : 'no'}">
        ${{'{{label}}'.trim() ? '{{label}}' : 'Default'}}
    </div>
`,
```

### ✅ Good
```go
Template: `
    <div data-active="${{{isActive}} === true ? 'yes' : 'no'}">
        ${{{label}}.trim() ? {{label}} : 'Default'}
    </div>
`,
```

## Exit Codes

- `0`: No issues found
- `1`: Issues detected or error occurred

## Flags

- `-v`: Verbose output (shows more details)
- `-fix`: Show suggested fixes for each issue

## Limitations

- May produce false positives in some edge cases
- Does not validate Go syntax, only template patterns
- Cannot detect all logical errors, only syntactic patterns

## Contributing

If you find a pattern that should be detected but isn't, please:
1. Document the pattern and why it's problematic
2. Add a test case
3. Update the detection logic

## Related Documentation

- [TEMPLATE_EXPRESSION_FIX_GUIDE.md](../../TEMPLATE_EXPRESSION_FIX_GUIDE.md)
- [CHANGELOG_JSON_SERIALIZATION.md](../../CHANGELOG_JSON_SERIALIZATION.md)
