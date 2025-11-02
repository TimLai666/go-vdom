# Event Handler Quick Reference

## Basic Patterns

### ✅ Correct Usage

```go
// Synchronous event handler (no event needed)
Button(Props{
    "onClick": js.Do(nil,
        js.Alert("'Hello!'"),
    ),
}, "Click me")

// Event handler with event parameter
Button(Props{
    "onClick": js.Do([]string{"event"},
        js.Const("target", "event.target"),
        js.Alert("'Clicked!'"),
    ),
}, "Click me")

// Async event handler
Button(Props{
    "onClick": js.AsyncDo(nil,
        js.Const("data", "await fetch('/api')"),
        js.Alert("'Done!'"),
    ),
}, "Click me")
```

### ❌ Incorrect Usage (Will NOT Work)

```go
// ❌ DON'T: Creates function but doesn't execute it
Button(Props{
    "onClick": js.Fn(nil, js.Alert("'Hello!'")),
}, "Click me")

// ❌ DON'T: Forgot event parameter when using event object
Button(Props{
    "onClick": js.Do(nil,
        js.Const("val", "event.target.value"),  // Error: event is undefined!
    ),
}, "Click me")
```

## Common Use Cases

### Simple Action
```go
"onClick": js.Do(nil,
    js.Alert("'Clicked!'"),
)
```

### DOM Manipulation (no event)
```go
"onClick": js.Do(nil,
    js.Const("el", "document.getElementById('myDiv')"),
    JSAction{Code: "el.style.display = 'none'"},
)
```

### DOM Manipulation (with event)
```go
"onClick": js.Do([]string{"event"},
    js.Const("target", "event.target"),
    JSAction{Code: "target.style.display = 'none'"},
)
```

### API Call with Error Handling
```go
"onClick": js.AsyncDo(nil,
    js.Try(
        js.Const("response", "await fetch('/api/data')"),
        js.Const("data", "await response.json()"),
        js.Alert("'Success!'"),
    ).Catch(
        js.Alert("'Error: ' + error.message"),
    ).End(),
)
```

### Update Counter
```go
"onClick": js.Do(nil,
    js.Const("el", "document.getElementById('counter')"),
    js.Const("val", "parseInt(el.textContent) + 1"),
    JSAction{Code: "el.textContent = val"},
)
```

### Form Input Handler
```go
"onInput": js.Do([]string{"event"},
    js.Const("value", "event.target.value"),
    JSAction{Code: "document.getElementById('output').textContent = value"},
)
```

### Mouse Events
```go
"onMouseOver": js.Do([]string{"event"},
    JSAction{Code: "event.target.style.backgroundColor = 'yellow'"},
)
"onMouseOut": js.Do([]string{"event"},
    JSAction{Code: "event.target.style.backgroundColor = ''"},
)
```

## Event Types

| Event | Description | Example Use Case |
|-------|-------------|------------------|
| `onClick` | Mouse click | Buttons, links |
| `onDblClick` | Double click | Special actions |
| `onInput` | Input value changes | Live search, validation |
| `onChange` | Select/checkbox changes | Form handling |
| `onMouseOver` | Mouse enters element | Hover effects |
| `onMouseOut` | Mouse leaves element | Unhover effects |
| `onFocus` | Element receives focus | Input highlighting |
| `onBlur` | Element loses focus | Form validation |
| `onSubmit` | Form submission | Form handling |
| `onKeyDown` | Key pressed | Keyboard shortcuts |

## Decision Tree

```
Need to handle an event?
├─ Is it async (uses await/fetch/setTimeout)?
│  └─ YES → Use js.AsyncDo(...)
│
└─ Is it synchronous?
   └─ YES → Use js.Do(...)

Need to define a reusable function?
├─ Is it async?
│  └─ YES → Use js.AsyncFn([params], ...)
│
└─ Is it synchronous?
   └─ YES → Use js.Fn([params], ...)
```

## Quick Rules

1. **Event handlers = Use Do/AsyncDo**
2. **Async operations = Use AsyncDo**
3. **Simple actions = Use Do**
4. **Never use Fn/AsyncFn for events** (they don't execute)
5. **Using event object = Must declare []string{"event"} parameter**
6. **Always handle errors** with Try/Catch in AsyncDo

## Generated Code Examples

### Without event parameter
**Input:**
```go
"onClick": js.Do(nil, js.Alert("'Hi'"))
```

**Output:**
```html
onclick="(()=>{alert('Hi')})()"
```

### With event parameter
**Input:**
```go
"onClick": js.Do([]string{"event"},
    js.Const("val", "event.target.value"),
    js.Alert("val"),
)
```

**Output:**
```html
onclick="((event)=>{const val=event.target.value;alert(val)})(event)"
```

### Async without event
**Input:**
```go
"onClick": js.AsyncDo(nil,
    js.Const("x", "await fetch('/api')"),
    js.Alert("'Done'"),
)
```

**Output:**
```html
onclick="(async()=>{const x=await fetch('/api');alert('Done')})()"
```

## Important Notes

⚠️ **Event Parameter Requirement**

When using the `event` object in your code, you **must** declare it as a parameter:

```go
// ❌ Wrong - event is undefined in IIFE scope
js.Do(nil,
    js.Const("val", "event.target.value"),
)

// ✅ Correct - event is passed as parameter
js.Do([]string{"event"},
    js.Const("val", "event.target.value"),
)
```

**Why?** IIFE creates a new scope. The external `event` object must be explicitly passed in as a parameter.

## See Also

- [EVENT_HANDLER_CHANGES.md](EVENT_HANDLER_CHANGES.md) - Detailed migration guide
- [DO_ASYNCDO_PARAMS.md](DO_ASYNCDO_PARAMS.md) - Parameter usage guide
- [TRY_CATCH_QUICK_REF.md](TRY_CATCH_QUICK_REF.md) - Try/Catch patterns
- [examples/09_event_handlers.go](../examples/09_event_handlers.go) - Live examples
- [examples/10_do_with_params.go](../examples/10_do_with_params.go) - Parameter examples