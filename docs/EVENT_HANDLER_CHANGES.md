# Event Handler Changes - v1.2.1

## Summary

As of version 1.2.1, the automatic function wrapping and detection mechanism for event handlers has been **removed**. This change simplifies the rendering logic and gives developers explicit control over event handler execution.

## What Changed

### Before (v1.2.0 and earlier)

The renderer attempted to intelligently detect whether a `JSAction` was a function definition or statement, and would automatically wrap it:

```go
// This would be auto-detected and wrapped
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Alert("'Hello'"),
    ),
}, "Click me")

// Renderer would detect it's a function and wrap it in function(evt,el){...}
```

**Problems:**
- Complex detection logic that could fail
- `AsyncDo()` generates `(async()=>{...})()` which was misdetected as a function definition
- This caused errors like "...is not a function"
- Developers had less control over execution context

### After (v1.2.1)

Event handlers are rendered **exactly as provided** - no automatic wrapping or detection:

```go
// Use js.Do() for synchronous IIFE
Button(Props{
    "onClick": js.Do(
        js.Alert("'Hello'"),
    ),
}, "Click me")

// Use js.AsyncDo() for async IIFE
Button(Props{
    "onClick": js.AsyncDo(
        js.Const("response", "await fetch('/api/data')"),
        js.Alert("'Done!'"),
    ),
}, "Click me")
```

## Migration Guide

### ❌ Don't Use (Will Not Work)

```go
// DON'T: js.Fn() creates a function definition but doesn't execute it
Button(Props{
    "onClick": js.Fn(nil,
        js.Alert("'Hello'"),
    ),
}, "Click me")
// Result: onclick="()=>{alert('Hello')}" - defines function but doesn't call it

// DON'T: js.AsyncFn() also creates a function definition
Button(Props{
    "onClick": js.AsyncFn(nil,
        js.Alert("'Hello'"),
    ),
}, "Click me")
// Result: onclick="async ()=>{alert('Hello')}" - defines function but doesn't call it
```

### ✅ Do Use (Correct Approach)

```go
// ✅ DO: Use js.Do() for synchronous operations
Button(Props{
    "onClick": js.Do(nil,
        js.Alert("'Hello'"),
    ),
}, "Click me")
// Result: onclick="(()=>{alert('Hello')})()" - IIFE that executes immediately

// ✅ DO: Use js.AsyncDo() for async operations
Button(Props{
    "onClick": js.AsyncDo(nil,
        js.Const("response", "await fetch('/api/data')"),
        js.Alert("'Done!'"),
    ),
}, "Click me")
// Result: onclick="(async()=>{...})()" - async IIFE that executes immediately
```

## Common Patterns

### 1. Simple Alert/Action

```go
Button(Props{
    "onClick": js.Do(nil,
        js.Alert("'Button clicked!'"),
    ),
}, "Click me")
```

### 2. DOM Manipulation

```go
Button(Props{
    "onClick": js.Do(nil,
        js.Const("el", "document.getElementById('myDiv')"),
        JSAction{Code: "el.style.display = 'none'"},
    ),
}, "Hide element")
```

### 3. Async API Call with Error Handling

```go
Button(Props{
    "onClick": js.AsyncDo(nil,
        js.Const("container", "document.getElementById('result')"),
        JSAction{Code: "container.innerHTML = 'Loading...'"},
        js.Try(
            js.Const("response", "await fetch('/api/data')"),
            js.Const("data", "await response.json()"),
            JSAction{Code: "container.innerHTML = 'Success: ' + data.message"},
        ).Catch(
            JSAction{Code: "container.innerHTML = 'Error: ' + error.message"},
        ).End(),
    ),
}, "Load data")
```

### 4. Multiple Events on Same Element

```go
Div(Props{
    "onClick": js.Do(nil,
        js.Log("'Clicked'"),
    ),
    "onMouseOver": js.Do(nil,
        JSAction{Code: "event.target.style.backgroundColor = 'yellow'"},
    ),
    "onMouseOut": js.Do(nil,
        JSAction{Code: "event.target.style.backgroundColor = ''"},
    ),
}, "Interactive div")
```

### 5. Form Event Handlers

```go
Input(Props{
    "type": "text",
    "onInput": js.Do(nil,
        js.Const("value", "event.target.value"),
        JSAction{Code: "document.getElementById('output').textContent = value"},
    ),
})
```

## Benefits of This Change

1. **Explicit Control**: Developers explicitly choose `Do()` or `AsyncDo()`, making intent clear
2. **No Magic**: What you write is what gets rendered - no hidden transformations
3. **Predictable Behavior**: No complex detection logic that might fail
4. **Better Error Messages**: Easier to debug because generated code matches source code
5. **Consistency**: All event handlers follow the same pattern

## Technical Details

### Old Behavior (Removed)

```go
// Old render.go logic (REMOVED):
if isArrowFunction || isFunctionKeyword {
    // Detected as function, wrap in another function
    handlerFn = fmt.Sprintf("function(evt,el){(%s)(evt,el);}", trimmedCode)
} else {
    // Not a function, wrap in function
    handlerFn = fmt.Sprintf("function(evt,el){%s}", safeCode)
}
```

### New Behavior (Current)

```go
// New render.go logic (SIMPLIFIED):
case JSAction:
    // Direct rendering - no wrapping or detection
    safeCode := strings.ReplaceAll(t.Code, "\"", "&quot;")
    sb.WriteString(fmt.Sprintf(" %s=\"%s\"", k, safeCode))
```

## When to Use What

| Use Case | Helper | Example |
|----------|--------|---------|
| Simple synchronous action | `js.Do()` | Alert, console.log, simple DOM manipulation |
| Async operation with await | `js.AsyncDo()` | fetch, setTimeout, async API calls |
| Define reusable function | `js.Fn()` | Store in variable, pass to other functions |
| Define async function | `js.AsyncFn()` | Store in variable for later use |

## Breaking Changes

If you have existing code that uses `js.Fn()` or `js.AsyncFn()` for event handlers:

1. **Find all event handlers** using `Fn` or `AsyncFn`
2. **Replace with `Do()` or `AsyncDo()`**:
   - `js.Fn(nil, ...)` → `js.Do(nil, ...)`
   - `js.AsyncFn(nil, ...)` → `js.AsyncDo(nil, ...)`
3. **Test all event handlers** to ensure they execute correctly

## See Also

- [TRY_CATCH_FINALLY.md](TRY_CATCH_FINALLY.md) - Try/Catch/Finally patterns
- [OPTIMIZATION.md](OPTIMIZATION.md) - Code generation optimizations
- [examples/09_event_handlers.go](../examples/09_event_handlers.go) - Comprehensive event handler examples

## Version History

- **v1.2.1**: Removed automatic function wrapping, introduced explicit `Do()`/`AsyncDo()` pattern
- **v1.2.0**: Introduced `Try()` builder API, `Do()` and `AsyncDo()` helpers
- **v1.1.x**: Had automatic function detection (now removed)