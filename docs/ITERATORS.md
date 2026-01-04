# Mintyex Iterator Extensions

This document describes the functional programming helpers added to the mintyex package to provide JavaScript-like array operations for Go slices.

## Overview

These iterator functions provide a functional programming API for common slice operations, making it easier for developers (especially those coming from JavaScript) to work with collections in Go while maintaining type safety and performance.

## Installation

Add these files to your existing mintyex package:
- `iterators.go` - Core functionality
- `iterators_test.go` - Test suite

Ensure your `go.mod` requires Go 1.18+ for generics support.

## Core Functions

### Basic Operations

```go
// Filter elements that match a predicate
activeUsers := mintyex.Filter(users, func(u User) bool { return u.Active })

// Transform elements to a new type
names := mintyex.Map(users, func(u User) string { return u.Name })

// Take the first n elements
firstFive := mintyex.Take(users, 5)

// Skip the first n elements
remaining := mintyex.Skip(users, 10)

// Find the first matching element
user, found := mintyex.Find(users, func(u User) bool { return u.Name == "Alice" })

// Check if any elements match
hasActive := mintyex.Any(users, func(u User) bool { return u.Active })

// Check if all elements match
allActive := mintyex.All(users, func(u User) bool { return u.Active })
```

### Advanced Operations

```go
// Group elements by a key
grouped := mintyex.GroupBy(users, func(u User) bool { return u.Active })
// Returns: map[bool][]User

// Remove duplicates
unique := mintyex.Unique([]string{"a", "b", "a", "c"})
// Returns: []string{"a", "b", "c"}

// Remove duplicates by key
uniqueUsers := mintyex.UniqueBy(users, func(u User) int { return u.Age })

// Reverse a slice
reversed := mintyex.Reverse(numbers)

// Reduce/fold over elements
sum := mintyex.Reduce(numbers, 0, func(acc, num int) int { return acc + num })

// Split into two groups
evens, odds := mintyex.Partition(numbers, func(n int) bool { return n%2 == 0 })

// Split into chunks
chunks := mintyex.Chunk(numbers, 3)
// [1,2,3,4,5,6,7] becomes [[1,2,3], [4,5,6], [7]]
```

## Chainable Operations

For complex transformations, use the fluent Chain API:

```go
result := mintyex.ChainSlice(users).
    Filter(func(u User) bool { return u.Active }).
    Filter(func(u User) bool { return u.Age > 25 }).
    Take(10).
    Unique().
    Reverse().
    ToSlice()

// Or get metadata about the chain
count := mintyex.ChainSlice(users).
    Filter(func(u User) bool { return u.Active }).
    Count()

first, hasFirst := mintyex.ChainSlice(users).
    Filter(func(u User) bool { return u.Active }).
    First()
```

## HTML Integration

Special helpers for working with Minty HTML generation:

```go
// Filter and render in one step
userCards := mintyex.FilterAndRender(users,
    func(u User) bool { return u.Active },
    func(u User) mi.H { return UserCard(theme, u) },
)

// Conditional rendering
adminPanel := mintyex.RenderIf(adminUsers, currentUser.IsAdmin, 
    func(u User) mi.H { return AdminUserCard(theme, u) },
)

// Render first n items
topUsers := mintyex.RenderFirst(users, 5,
    func(u User) mi.H { return UserSummary(theme, u) },
)

// Render with index
numberedList := mintyex.EachWithIndex(items,
    func(item Item, index int) mi.H {
        return b.Li(fmt.Sprintf("%d. %s", index+1, item.Name))
    },
)

// Render in chunks (for grid layouts)
userGrid := mintyex.ChunkAndRender(users, 3,
    func(chunk []User) mi.H {
        return b.Div(mi.Class("user-row"),
            mintyex.Map(chunk, func(u User) mi.H { return UserCard(theme, u) })...,
        )
    },
)
```

## Integration with Existing Code

### Enhanced Each Function

The existing `Each` function can be updated to use the new `Map` internally:

```go
// Before (if this exists)
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

### Migration Examples

**Traditional Go:**
```go
var activeUsers []User
for _, user := range users {
    if user.Active {
        activeUsers = append(activeUsers, user)
    }
}

var names []string
for _, user := range activeUsers {
    names = append(names, user.Name)
}

if len(names) > 10 {
    names = names[:10]
}
```

**With Iterator helpers:**
```go
names := mintyex.ChainSlice(users).
    Filter(func(u User) bool { return u.Active }).
    Map(func(u User) string { return u.Name }).
    Take(10).
    ToSlice().([]string)
```

**Even simpler:**
```go
names := mintyex.Take(
    mintyex.Map(
        mintyex.Filter(users, func(u User) bool { return u.Active }),
        func(u User) string { return u.Name },
    ), 
    10,
)
```

## Performance Considerations

- All functions create new slices rather than mutating input
- `Map` pre-allocates the result slice for better performance
- Chain operations create intermediate slices (consider direct function composition for hot paths)
- Functions handle nil and empty slices safely
- No reflection is used - everything is compile-time type-safe

## Comparison with JavaScript

| JavaScript | Mintyex Go | Notes |
|------------|------------|-------|
| `arr.filter(fn)` | `Filter(arr, fn)` | Same semantics |
| `arr.map(fn)` | `Map(arr, fn)` | Same semantics |
| `arr.slice(0, n)` | `Take(arr, n)` | Take first n |
| `arr.slice(n)` | `Skip(arr, n)` | Skip first n |
| `arr.find(fn)` | `Find(arr, fn)` | Returns value + boolean |
| `arr.some(fn)` | `Any(arr, fn)` | Same semantics |
| `arr.every(fn)` | `All(arr, fn)` | Same semantics |
| `arr.reduce(fn, init)` | `Reduce(arr, init, fn)` | Parameter order different |
| `[...new Set(arr)]` | `Unique(arr)` | Remove duplicates |
| `arr.reverse()` | `Reverse(arr)` | Returns new slice |

## Examples in Context

### Simple Todo App with Filtering

```go
func TodoList(todos []Todo, filter string) mi.H {
    var filteredTodos []Todo
    
    switch filter {
    case "active":
        filteredTodos = mintyex.Filter(todos, func(t Todo) bool { return !t.Completed })
    case "completed":
        filteredTodos = mintyex.Filter(todos, func(t Todo) bool { return t.Completed })
    default:
        filteredTodos = todos
    }
    
    return func(b *mi.Builder) mi.Node {
        return b.Ul(mi.Class("todo-list"),
            mintyex.Map(filteredTodos, func(todo Todo) mi.H {
                return TodoItem(todo)
            })...,
        )
    }
}
```

### User Dashboard with Complex Logic

```go
func UserDashboard(users []User, currentUser User) mi.H {
    return func(b *mi.Builder) mi.Node {
        // Get active users excluding current user
        others := mintyex.Filter(users, func(u User) bool { 
            return u.Active && u.ID != currentUser.ID 
        })
        
        // Group by department
        byDept := mintyex.GroupBy(others, func(u User) string { return u.Department })
        
        // Render each department
        var deptSections []mi.H
        for dept, deptUsers := range byDept {
            // Top 5 users by last activity
            topUsers := mintyex.ChainSlice(deptUsers).
                Take(5).
                ToSlice()
                
            deptSection := b.Div(mi.Class("department"),
                b.H3(dept),
                b.Div(
                    mintyex.Map(topUsers, func(u User) mi.H {
                        return UserCard(u)
                    })...,
                ),
            )
            deptSections = append(deptSections, deptSection)
        }
        
        return b.Div(mi.Class("dashboard"),
            b.H1("Team Dashboard"),
            mi.NewFragment(deptSections...),
        )
    }
}
```

## Testing

Run the tests with:

```bash
go test ./... -v
```

Benchmarks are included to ensure performance is acceptable:

```bash
go test -bench=. -benchmem
```

## Migration Path

1. **Add the files** to your mintyex package
2. **Update go.mod** to require Go 1.18+
3. **Gradually adopt** in new code or when refactoring
4. **Replace manual loops** with iterator functions where it improves readability
5. **Use chains** for complex multi-step transformations

The beauty of this approach is that it's **completely additive** - existing code continues to work unchanged, but you get access to more ergonomic operations when you need them.
