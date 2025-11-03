# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Template Expression System**: Added support for conditional expressions in component templates
  - Supports ternary operators (`condition ? true : false`)
  - Supports nested ternary expressions
  - Supports comparison operators (`===`, `!==`, `==`, `!=`)
  - Supports `.trim()` for empty string checks
  - See `docs/COMPONENT_REFACTORING.md` for detailed documentation

### Changed
- **Component System Refactoring**: Major refactoring of component architecture
  - Removed all derived property calculations from `Component` function
  - Moved property derivation logic to component templates using `${}` expressions
  - `Component` function now only handles prop merging and template interpolation
  - Updated 8 components to use template expressions:
    - Alert: closable, type, title, icon, rounded, compact, elevation
    - Button: size, variant, disabled, fullWidth, rounded, icon
    - Radio: direction, label, helpText
    - Checkbox: direction, label, helpText
    - Switch: direction, helpText
    - Card: title, elevation, hoverable
    - Modal: open, size, centered, radius, hideHeader, closeButton, scrollable, animation
    - Table: footer display

### Improved
- **Maintainability**: Each component now manages its own derivation logic
- **Modularity**: Component function remains generic and reusable
- **Performance**: Expressions evaluated at render time, no runtime overhead
- **Type Safety**: All component APIs remain unchanged, backward compatible

### Fixed
- Alert `closable` parameter now correctly shows/hides close button
- Button `size`, `variant`, `disabled` parameters now work as expected
- Radio/Checkbox/Switch `direction` parameter properly controls layout
- Card `title` and `elevation` parameters apply correctly
- Modal `open`, `size`, and other parameters function properly

### Tests
- Added comprehensive test suite for components:
  - 6 Alert component tests
  - 5 Button component tests
  - 5 Radio/Checkbox/Switch tests
  - All 16 tests passing

### Documentation
- Added `docs/COMPONENT_REFACTORING.md` - Complete refactoring documentation
- Added template expression examples and best practices
- Updated component usage documentation
- Added migration guide for component developers

## [Previous Versions]

### Notable Features
- Virtual DOM implementation
- Component system
- Control flow structures (If/Then/Else, Repeat, For)
- JavaScript DSL with async/await support
- Code minification
- UI component library
- Server-side rendering
- Bootstrap integration
- Template serialization (Go templates, JSON)
- Type-flexible Props system
- Event handling with `js.Do()` and `js.AsyncDo()`

---

## Migration Guide

### For Users
**No changes required!** All component APIs remain backward compatible. Your existing code will continue to work without modifications.

### For Component Developers
If you're creating custom components or modifying existing ones:

**Before:**
```go
// ❌ Don't add component-specific logic to Component function
if myProp, ok := mergedProps["myProp"]; ok && myProp == "value" {
    mergedProps["myDerivedProp"] = "result"
}
```

**After:**
```go
// ✅ Use template expressions in component templates
Props{
    "style": `
        property: ${'{{myProp}}' === 'value' ? 'result' : 'default'};
    `,
}
```

For detailed information, see `docs/COMPONENT_REFACTORING.md`.
