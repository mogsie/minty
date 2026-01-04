# Mintyex Iterator Extensions - Installation Guide

This package contains functional programming helpers for the mintyex package.

## Files Included

```
mintyex-iterators/
├── iterators.go              # Core iterator functions
├── iterators_test.go         # Comprehensive test suite  
├── integration_example.go    # Examples showing usage with Minty components
└── README_ITERATORS.md       # Detailed documentation
```

## Installation Instructions

### 1. Copy Files to Mintyex Package

Copy these files to your existing `mintyex` package directory:

```bash
# Assuming your mintyex package is at github.com/ha1tch/mintyex
cp iterators.go /path/to/mintyex/
cp iterators_test.go /path/to/mintyex/
```

### 2. Update Go Module Requirements

Ensure your `go.mod` requires Go 1.18+ for generics support:

```go
module github.com/ha1tch/mintyex

go 1.21 // or whatever version >= 1.18
```

### 3. Update Existing Each Function (Optional)

If you have an existing `Each` function in your codebase, you can optionally update it to use the new `Map` function internally:

```go
// Before
func Each[T any](slice []T, render func(T) mi.H) []mi.H {
    result := make([]mi.H, len(slice))
    for i, item := range slice {
        result[i] = render(item)
    }
    return result
}

// After - using new Map function  
func Each[T any](slice []T, render func(T) mi.H) []mi.H {
    return Map(slice, render)
}
```

### 4. Run Tests

Verify everything works:

```bash
cd /path/to/mintyex
go test -v
go test -bench=. -benchmem  # Run benchmarks
```

### 5. Import and Use

In your applications:

```go
import (
    "github.com/ha1tch/mintyex"
)

// Use the new functions
activeUsers := mintyex.Filter(users, func(u User) bool { return u.Active })
userCards := mintyex.Map(activeUsers, func(u User) mi.H { return UserCard(theme, u) })
```

## What's Added

### Core Functions
- `Filter[T]` - Filter elements matching predicate
- `Map[T, U]` - Transform elements to new type
- `Take[T]` - Take first n elements  
- `Skip[T]` - Skip first n elements
- `Find[T]` - Find first matching element
- `Any[T]` - Check if any element matches
- `All[T]` - Check if all elements match
- `Reduce[T, U]` - Fold/reduce over elements
- `GroupBy[T, K]` - Group elements by key function
- `Unique[T]` - Remove duplicates
- `UniqueBy[T, K]` - Remove duplicates by key
- `Reverse[T]` - Reverse slice
- `Partition[T]` - Split into two slices
- `Chunk[T]` - Split into chunks

### Chainable API
- `ChainSlice[T]` - Start fluent operations
- Chain methods: `Filter`, `Take`, `Skip`, `Unique`, `Reverse`, `Map`
- Terminal methods: `ToSlice`, `Count`, `First`, `Last`

### HTML-Specific Helpers
- `FilterAndRender[T]` - Filter and render in one step
- `RenderIf[T]` - Conditional rendering
- `RenderFirst[T]` - Render first n elements
- `RenderWhen[T]` - Render elements matching condition
- `EachWithIndex[T]` - Render with index
- `ChunkAndRender[T]` - Chunk and render each chunk

## Backward Compatibility

These additions are **completely additive**. Existing code will continue to work unchanged. The new functions provide additional capabilities for when you want more ergonomic slice operations.

## Performance Notes

- All functions create new slices (no mutation)
- `Map` pre-allocates result slices for better performance  
- Chain operations create intermediate slices
- No reflection used - fully compile-time type-safe
- Benchmarks included to validate performance

## Migration Strategy

1. **Start gradual**: Use in new code or when refactoring
2. **Replace manual loops**: Where iterator functions improve readability
3. **Use chains**: For complex multi-step transformations
4. **Leverage HTML helpers**: For common UI rendering patterns

## Examples

See `integration_example.go` for comprehensive examples showing how to use these functions with Minty components for real-world UI development scenarios.

## Questions or Issues

If you encounter any issues with the integration, check:
1. Go version is 1.18+ 
2. Files are in the correct package directory
3. Tests pass: `go test -v`
4. No naming conflicts with existing functions

The functions are designed to feel familiar to JavaScript developers while maintaining Go's explicit, performant characteristics.
