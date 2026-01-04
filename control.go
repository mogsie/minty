package minty

// Control flow helpers for conditional rendering

// If returns the template if the condition is true, otherwise returns an empty fragment.
func If(condition bool, template H) H {
	if condition {
		return template
	}
	return func(b *Builder) Node {
		return NewFragment() // Return empty fragment
	}
}

// IfElse returns the trueTemplate if condition is true, otherwise returns falseTemplate.
func IfElse(condition bool, trueTemplate, falseTemplate H) H {
	if condition {
		return trueTemplate
	}
	return falseTemplate
}

// Unless returns the template if the condition is false (opposite of If).
func Unless(condition bool, template H) H {
	return If(!condition, template)
}

// Each renders a slice of items using the provided renderer function.
// Returns a slice of Nodes that can be spread into parent elements.
func Each[T any](items []T, renderer func(T) H) []Node {
	if len(items) == 0 {
		return []Node{}
	}
	
	nodes := make([]Node, 0, len(items))
	for _, item := range items {
		template := renderer(item)
		node := template(B) // Use global builder instance
		nodes = append(nodes, node)
	}
	return nodes
}

// EachWithIndex renders a slice of items with index using the provided renderer function.
func EachWithIndex[T any](items []T, renderer func(int, T) H) []Node {
	if len(items) == 0 {
		return []Node{}
	}
	
	nodes := make([]Node, 0, len(items))
	for i, item := range items {
		template := renderer(i, item)
		node := template(B) // Use global builder instance
		nodes = append(nodes, node)
	}
	return nodes
}

// Map transforms a slice of one type to a slice of Nodes using a renderer.
func Map[T any](items []T, renderer func(T) H) H {
	return func(b *Builder) Node {
		if len(items) == 0 {
			return NewFragment()
		}
		
		nodes := Each(items, renderer)
		return NewFragment(nodes...)
	}
}

// Filter renders only items that match the predicate condition.
func Filter[T any](items []T, predicate func(T) bool, renderer func(T) H) []Node {
	if len(items) == 0 {
		return []Node{}
	}
	
	var nodes []Node
	for _, item := range items {
		if predicate(item) {
			template := renderer(item)
			node := template(B)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// Range generates a sequence of numbers and renders each using the renderer.
func Range(start, end int, renderer func(int) H) []Node {
	if start >= end {
		return []Node{}
	}
	
	nodes := make([]Node, 0, end-start)
	for i := start; i < end; i++ {
		template := renderer(i)
		node := template(B)
		nodes = append(nodes, node)
	}
	return nodes
}

// When provides multiple condition/template pairs (like a switch statement).
type WhenCase[T comparable] struct {
	Value    T
	Template H
}

// When renders the first matching case based on the value.
func When[T comparable](value T, cases []WhenCase[T], defaultTemplate H) H {
	for _, c := range cases {
		if c.Value == value {
			return c.Template
		}
	}
	if defaultTemplate != nil {
		return defaultTemplate
	}
	return func(b *Builder) Node {
		return NewFragment()
	}
}

// Repeat renders the same template multiple times.
func Repeat(count int, template H) []Node {
	if count <= 0 {
		return []Node{}
	}
	
	nodes := make([]Node, 0, count)
	for i := 0; i < count; i++ {
		node := template(B)
		nodes = append(nodes, node)
	}
	return nodes
}

// Join renders items with a separator between them.
func Join[T any](items []T, renderer func(T) H, separator H) H {
	return func(b *Builder) Node {
		if len(items) == 0 {
			return NewFragment()
		}
		
		var nodes []Node
		for i, item := range items {
			if i > 0 && separator != nil {
				sepNode := separator(b)
				nodes = append(nodes, sepNode)
			}
			template := renderer(item)
			node := template(b)
			nodes = append(nodes, node)
		}
		
		return NewFragment(nodes...)
	}
}

// GroupBy groups items by a key function and renders each group.
func GroupBy[T any, K comparable](items []T, keyFn func(T) K, renderer func(K, []T) H) []Node {
	if len(items) == 0 {
		return []Node{}
	}
	
	groups := make(map[K][]T)
	var order []K
	
	for _, item := range items {
		key := keyFn(item)
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], item)
	}
	
	nodes := make([]Node, 0, len(groups))
	for _, key := range order {
		template := renderer(key, groups[key])
		node := template(B)
		nodes = append(nodes, node)
	}
	
	return nodes
}

// Partition splits items into two groups based on a predicate.
func Partition[T any](items []T, predicate func(T) bool, trueRenderer, falseRenderer func([]T) H) H {
	var trueItems, falseItems []T
	
	for _, item := range items {
		if predicate(item) {
			trueItems = append(trueItems, item)
		} else {
			falseItems = append(falseItems, item)
		}
	}
	
	return func(b *Builder) Node {
		var nodes []Node
		
		if len(trueItems) > 0 && trueRenderer != nil {
			nodes = append(nodes, trueRenderer(trueItems)(b))
		}
		
		if len(falseItems) > 0 && falseRenderer != nil {
			nodes = append(nodes, falseRenderer(falseItems)(b))
		}
		
		return NewFragment(nodes...)
	}
}

// Take renders only the first n items from a slice.
func Take[T any](items []T, n int, renderer func(T) H) []Node {
	if n <= 0 || len(items) == 0 {
		return []Node{}
	}
	
	end := n
	if end > len(items) {
		end = len(items)
	}
	
	return Each(items[:end], renderer)
}

// Skip renders items starting from the nth item.
func Skip[T any](items []T, n int, renderer func(T) H) []Node {
	if n >= len(items) || len(items) == 0 {
		return []Node{}
	}
	
	if n < 0 {
		n = 0
	}
	
	return Each(items[n:], renderer)
}

// Chunk renders items in chunks of specified size.
func Chunk[T any](items []T, chunkSize int, renderer func([]T) H) []Node {
	if chunkSize <= 0 || len(items) == 0 {
		return []Node{}
	}
	
	var nodes []Node
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		
		chunk := items[i:end]
		template := renderer(chunk)
		node := template(B)
		nodes = append(nodes, node)
	}
	
	return nodes
}

// WithDefault returns the template if items is not empty, otherwise returns the default template.
func WithDefault[T any](items []T, template H, defaultTemplate H) H {
	if len(items) > 0 {
		return template
	}
	return defaultTemplate
}

// Paginate renders items with pagination support.
type PaginationInfo struct {
	CurrentPage int
	PerPage     int
	TotalItems  int
	TotalPages  int
}

func Paginate[T any](items []T, page, perPage int, renderer func([]T, PaginationInfo) H) H {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	
	totalItems := len(items)
	totalPages := (totalItems + perPage - 1) / perPage
	
	if page > totalPages {
		page = totalPages
	}
	
	start := (page - 1) * perPage
	if start < 0 {
		start = 0
	}
	
	end := start + perPage
	if end > totalItems {
		end = totalItems
	}
	
	paginatedItems := items[start:end]
	paginationInfo := PaginationInfo{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
	
	return renderer(paginatedItems, paginationInfo)
}
