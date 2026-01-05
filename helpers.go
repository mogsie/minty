// helpers.go provides utility functions for working with dynamic data in minty.
// These helpers bridge the gap between JavaScript's dynamic object model and Go's
// static type system, making it easier to work with map[string]interface{} data
// commonly encountered when migrating from React/JavaScript.
//
// Import with: import mi "github.com/ha1tch/minty"
package minty

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// =============================================================================
// TIER 1: Type-safe map accessors
// =============================================================================

// Str safely extracts a string from a map. Returns empty string if key is missing
// or value is not a string (falls back to fmt.Sprint for other types).
//
// Usage:
//
//	post := map[string]interface{}{"title": "Hello", "views": 42}
//	mi.Str(post, "title")  // "Hello"
//	mi.Str(post, "views")  // "42" (converted)
//	mi.Str(post, "missing") // ""
func Str(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// Int safely extracts an int from a map. Returns 0 if key is missing or
// value cannot be converted to int. Handles int, float64, and string types.
//
// Usage:
//
//	post := map[string]interface{}{"views": 42, "rating": 4.5}
//	mi.Int(post, "views")   // 42
//	mi.Int(post, "rating")  // 4 (truncated)
//	mi.Int(post, "missing") // 0
func Int(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}

// Float safely extracts a float64 from a map. Returns 0 if key is missing or
// value cannot be converted to float64.
//
// Usage:
//
//	product := map[string]interface{}{"price": 19.99}
//	mi.Float(product, "price") // 19.99
func Float(m map[string]interface{}, key string) float64 {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f
		}
	}
	return 0
}

// Bool safely extracts a bool from a map. Returns false if key is missing or
// value is not a boolean.
//
// Usage:
//
//	post := map[string]interface{}{"published": true}
//	mi.Bool(post, "published") // true
//	mi.Bool(post, "missing")   // false
func Bool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// =============================================================================
// TIER 1: Truthy checks (JavaScript-like)
// =============================================================================

// Truthy performs a JavaScript-like truthy check on a value.
// Returns false for: nil, false, 0, 0.0, "", empty slices
// Returns true for everything else.
//
// Usage:
//
//	mi.Truthy(nil)           // false
//	mi.Truthy("")            // false
//	mi.Truthy(0)             // false
//	mi.Truthy("hello")       // true
//	mi.Truthy(42)            // true
//	mi.Truthy([]int{1,2,3})  // true
func Truthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case int:
		return t != 0
	case int64:
		return t != 0
	case float64:
		return t != 0
	case float32:
		return t != 0
	case string:
		return t != ""
	case []interface{}:
		return len(t) > 0
	case []string:
		return len(t) > 0
	case []int:
		return len(t) > 0
	case map[string]interface{}:
		return len(t) > 0
	}
	return true
}

// IfTruthy renders content only if the value is truthy (JavaScript-like).
// This is useful for conditionally showing content based on map values
// without explicit nil/empty checks.
//
// Usage:
//
//	mi.IfTruthy(post["category"], func(b *mi.Builder) mi.Node {
//	    return b.Span(mi.Class("category"), mi.Str(post, "category"))
//	})
func IfTruthy(v interface{}, then H) Node {
	if Truthy(v) {
		return then(B)
	}
	return nil
}

// IfTruthyElse renders 'then' if value is truthy, otherwise renders 'els'.
//
// Usage:
//
//	mi.IfTruthyElse(user["avatar"],
//	    func(b *mi.Builder) mi.Node { return b.Img(mi.Src(mi.Str(user, "avatar"))) },
//	    func(b *mi.Builder) mi.Node { return b.Span("No avatar") },
//	)
func IfTruthyElse(v interface{}, then, els H) Node {
	if Truthy(v) {
		return then(B)
	}
	return els(B)
}

// =============================================================================
// TIER 2: Map field comparisons
// =============================================================================

// Eq checks if a map field equals a string value.
//
// Usage:
//
//	if mi.Eq(post, "status", "published") { ... }
func Eq(m map[string]interface{}, key, val string) bool {
	return Str(m, key) == val
}

// Ne checks if a map field does not equal a string value.
//
// Usage:
//
//	if mi.Ne(post, "status", "draft") { ... }
func Ne(m map[string]interface{}, key, val string) bool {
	return Str(m, key) != val
}

// Gt checks if a map field is greater than an int value.
//
// Usage:
//
//	if mi.Gt(post, "likes", 0) { ... }
func Gt(m map[string]interface{}, key string, val int) bool {
	return Int(m, key) > val
}

// Gte checks if a map field is greater than or equal to an int value.
func Gte(m map[string]interface{}, key string, val int) bool {
	return Int(m, key) >= val
}

// Lt checks if a map field is less than an int value.
func Lt(m map[string]interface{}, key string, val int) bool {
	return Int(m, key) < val
}

// Lte checks if a map field is less than or equal to an int value.
func Lte(m map[string]interface{}, key string, val int) bool {
	return Int(m, key) <= val
}

// Contains checks if a map field (as string) contains a substring.
// Case-sensitive by default.
//
// Usage:
//
//	if mi.Contains(post, "title", "Go") { ... }
func Contains(m map[string]interface{}, key, substr string) bool {
	return strings.Contains(Str(m, key), substr)
}

// ContainsI checks if a map field contains a substring (case-insensitive).
//
// Usage:
//
//	if mi.ContainsI(post, "title", "go") { ... }
func ContainsI(m map[string]interface{}, key, substr string) bool {
	return strings.Contains(strings.ToLower(Str(m, key)), strings.ToLower(substr))
}

// =============================================================================
// TIER 2: Predicate builders
// =============================================================================

// Predicate is a function that tests if an item matches a condition.
type Predicate func(interface{}) bool

// Where creates a predicate that checks if a map field equals a string value.
//
// Usage:
//
//	published := mi.Filter(posts, mi.Where("status", "published"))
func Where(key, val string) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Str(m, key) == val
		}
		return false
	}
}

// WhereNot creates a predicate that checks if a map field does NOT equal a value.
//
// Usage:
//
//	notDraft := mi.Filter(posts, mi.WhereNot("status", "draft"))
func WhereNot(key, val string) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Str(m, key) != val
		}
		return false
	}
}

// WhereGt creates a predicate that checks if a map field is > an int value.
//
// Usage:
//
//	popular := mi.Filter(posts, mi.WhereGt("views", 1000))
func WhereGt(key string, val int) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Int(m, key) > val
		}
		return false
	}
}

// WhereGte creates a predicate that checks if a map field is >= an int value.
func WhereGte(key string, val int) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Int(m, key) >= val
		}
		return false
	}
}

// WhereLt creates a predicate that checks if a map field is < an int value.
func WhereLt(key string, val int) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Int(m, key) < val
		}
		return false
	}
}

// WhereTruthy creates a predicate that checks if a map field is truthy.
//
// Usage:
//
//	withCategory := mi.Filter(posts, mi.WhereTruthy("category"))
func WhereTruthy(key string) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return Truthy(m[key])
		}
		return false
	}
}

// WhereContains creates a predicate that checks if a map field contains a substring.
//
// Usage:
//
//	goPosts := mi.Filter(posts, mi.WhereContains("title", "Go"))
func WhereContains(key, substr string) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return strings.Contains(Str(m, key), substr)
		}
		return false
	}
}

// WhereContainsI creates a case-insensitive contains predicate.
func WhereContainsI(key, substr string) Predicate {
	return func(item interface{}) bool {
		if m, ok := item.(map[string]interface{}); ok {
			return strings.Contains(strings.ToLower(Str(m, key)), strings.ToLower(substr))
		}
		return false
	}
}

// And combines multiple predicates with AND logic.
//
// Usage:
//
//	result := mi.Filter(posts, mi.And(
//	    mi.Where("status", "published"),
//	    mi.WhereGt("views", 100),
//	))
func And(predicates ...Predicate) Predicate {
	return func(item interface{}) bool {
		for _, p := range predicates {
			if !p(item) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple predicates with OR logic.
//
// Usage:
//
//	result := mi.Filter(posts, mi.Or(
//	    mi.Where("status", "published"),
//	    mi.Where("status", "featured"),
//	))
func Or(predicates ...Predicate) Predicate {
	return func(item interface{}) bool {
		for _, p := range predicates {
			if p(item) {
				return true
			}
		}
		return false
	}
}

// =============================================================================
// TIER 2: Collection helpers
// =============================================================================

// FilterItems returns items from a slice that match the predicate.
//
// Usage:
//
//	published := mi.FilterItems(posts, mi.Where("status", "published"))
//	popular := mi.FilterItems(posts, func(item interface{}) bool {
//	    p := item.(map[string]interface{})
//	    return mi.Int(p, "views") > 1000
//	})
func FilterItems(items []interface{}, pred Predicate) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range items {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}

// Count returns the number of items that match the predicate.
//
// Usage:
//
//	publishedCount := mi.Count(posts, mi.Where("status", "published"))
func Count(items []interface{}, pred Predicate) int {
	n := 0
	for _, item := range items {
		if pred(item) {
			n++
		}
	}
	return n
}

// Sum returns the sum of an integer field across all items.
//
// Usage:
//
//	totalViews := mi.Sum(posts, "views")
func Sum(items []interface{}, key string) int {
	total := 0
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			total += Int(m, key)
		}
	}
	return total
}

// SumFloat returns the sum of a float field across all items.
//
// Usage:
//
//	totalRevenue := mi.SumFloat(orders, "amount")
func SumFloat(items []interface{}, key string) float64 {
	total := 0.0
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			total += Float(m, key)
		}
	}
	return total
}

// Avg returns the average of an integer field across all items.
//
// Usage:
//
//	avgViews := mi.Avg(posts, "views")
func Avg(items []interface{}, key string) float64 {
	if len(items) == 0 {
		return 0
	}
	return float64(Sum(items, key)) / float64(len(items))
}

// MaxInt returns the maximum value of an integer field.
//
// Usage:
//
//	maxViews := mi.MaxInt(posts, "views")
func MaxInt(items []interface{}, key string) int {
	if len(items) == 0 {
		return 0
	}
	max := Int(items[0].(map[string]interface{}), key)
	for _, item := range items[1:] {
		if m, ok := item.(map[string]interface{}); ok {
			if v := Int(m, key); v > max {
				max = v
			}
		}
	}
	return max
}

// MinInt returns the minimum value of an integer field.
//
// Usage:
//
//	minViews := mi.MinInt(posts, "views")
func MinInt(items []interface{}, key string) int {
	if len(items) == 0 {
		return 0
	}
	min := Int(items[0].(map[string]interface{}), key)
	for _, item := range items[1:] {
		if m, ok := item.(map[string]interface{}); ok {
			if v := Int(m, key); v < min {
				min = v
			}
		}
	}
	return min
}

// Find returns the first item that matches the predicate, or nil if not found.
//
// Usage:
//
//	featured := mi.Find(posts, mi.Where("featured", "true"))
func Find(items []interface{}, pred Predicate) interface{} {
	for _, item := range items {
		if pred(item) {
			return item
		}
	}
	return nil
}

// Any returns true if any item matches the predicate.
//
// Usage:
//
//	hasPublished := mi.Any(posts, mi.Where("status", "published"))
func Any(items []interface{}, pred Predicate) bool {
	for _, item := range items {
		if pred(item) {
			return true
		}
	}
	return false
}

// All returns true if all items match the predicate.
//
// Usage:
//
//	allPublished := mi.All(posts, mi.Where("status", "published"))
func All(items []interface{}, pred Predicate) bool {
	for _, item := range items {
		if !pred(item) {
			return false
		}
	}
	return true
}

// =============================================================================
// TIER 3: Sort helpers
// =============================================================================

// SortDir represents sort direction.
type SortDir int

const (
	// Asc sorts in ascending order (A-Z, 0-9, oldest first)
	Asc SortDir = iota
	// Desc sorts in descending order (Z-A, 9-0, newest first)
	Desc
)

// SortBy returns a new slice sorted by the given field.
// The original slice is not modified.
//
// Usage:
//
//	byDate := mi.SortBy(posts, "date", mi.Desc)
//	byTitle := mi.SortBy(posts, "title", mi.Asc)
func SortBy(items []interface{}, key string, dir SortDir) []interface{} {
	result := make([]interface{}, len(items))
	copy(result, items)

	sort.Slice(result, func(i, j int) bool {
		mi, oki := result[i].(map[string]interface{})
		mj, okj := result[j].(map[string]interface{})
		if !oki || !okj {
			return false
		}

		vi, vj := mi[key], mj[key]

		// Compare based on type
		switch a := vi.(type) {
		case string:
			b, ok := vj.(string)
			if !ok {
				b = fmt.Sprint(vj)
			}
			if dir == Desc {
				return a > b
			}
			return a < b

		case int:
			var b int
			switch bv := vj.(type) {
			case int:
				b = bv
			case float64:
				b = int(bv)
			default:
				return false
			}
			if dir == Desc {
				return a > b
			}
			return a < b

		case int64:
			var b int64
			switch bv := vj.(type) {
			case int64:
				b = bv
			case int:
				b = int64(bv)
			case float64:
				b = int64(bv)
			default:
				return false
			}
			if dir == Desc {
				return a > b
			}
			return a < b

		case float64:
			var b float64
			switch bv := vj.(type) {
			case float64:
				b = bv
			case int:
				b = float64(bv)
			default:
				return false
			}
			if dir == Desc {
				return a > b
			}
			return a < b
		}

		// Fallback: convert to string
		as, bs := fmt.Sprint(vi), fmt.Sprint(vj)
		if dir == Desc {
			return as > bs
		}
		return as < bs
	})

	return result
}

// SortByMulti sorts by multiple fields in order of priority.
// Each field can have its own sort direction.
//
// Usage:
//
//	sorted := mi.SortByMulti(posts,
//	    mi.SortField{"status", mi.Asc},
//	    mi.SortField{"date", mi.Desc},
//	)
type SortField struct {
	Key string
	Dir SortDir
}

func SortByMulti(items []interface{}, fields ...SortField) []interface{} {
	result := make([]interface{}, len(items))
	copy(result, items)

	sort.Slice(result, func(i, j int) bool {
		mi, oki := result[i].(map[string]interface{})
		mj, okj := result[j].(map[string]interface{})
		if !oki || !okj {
			return false
		}

		for _, field := range fields {
			vi, vj := mi[field.Key], mj[field.Key]
			cmp := compareValues(vi, vj)
			if cmp != 0 {
				if field.Dir == Desc {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})

	return result
}

// compareValues returns -1, 0, or 1 for comparison.
func compareValues(a, b interface{}) int {
	// Handle nil
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	switch av := a.(type) {
	case int:
		bv, ok := toInt(b)
		if ok {
			if av < bv {
				return -1
			}
			if av > bv {
				return 1
			}
			return 0
		}
	case float64:
		bv, ok := toFloat(b)
		if ok {
			if av < bv {
				return -1
			}
			if av > bv {
				return 1
			}
			return 0
		}
	}

	// Fallback to string comparison
	as, bs := fmt.Sprint(a), fmt.Sprint(b)
	if as < bs {
		return -1
	}
	if as > bs {
		return 1
	}
	return 0
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

// =============================================================================
// Utility: Pluck / Map field extraction
// =============================================================================

// Pluck extracts a single field from each item as a slice of strings.
//
// Usage:
//
//	titles := mi.Pluck(posts, "title") // []string{"Hello", "World", ...}
func Pluck(items []interface{}, key string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, Str(m, key))
		}
	}
	return result
}

// PluckInt extracts a single integer field from each item.
//
// Usage:
//
//	viewCounts := mi.PluckInt(posts, "views") // []int{100, 200, ...}
func PluckInt(items []interface{}, key string) []int {
	result := make([]int, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, Int(m, key))
		}
	}
	return result
}

// GroupItems groups items by a string field value.
//
// Usage:
//
//	byStatus := mi.GroupItems(posts, "status")
//	// map[string][]interface{}{"published": [...], "draft": [...]}
func GroupItems(items []interface{}, key string) map[string][]interface{} {
	result := make(map[string][]interface{})
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			k := Str(m, key)
			result[k] = append(result[k], item)
		}
	}
	return result
}

// =============================================================================
// Utility: Safe type conversion for builders
// =============================================================================

// ToMap safely converts an interface{} to map[string]interface{}.
// Returns nil if conversion fails.
//
// Usage:
//
//	if m := mi.ToMap(item); m != nil {
//	    title := mi.Str(m, "title")
//	}
func ToMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// ToSlice safely converts an interface{} to []interface{}.
// Returns nil if conversion fails.
func ToSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return nil
}
