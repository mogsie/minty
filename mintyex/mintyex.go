// Package mintyex provides UI helper functions for the Minty System.
// It also re-exports business types from mintytypes for convenience.
//
// For pure business logic without UI dependencies, import mintytypes directly.
package mintyex

import (
	"fmt"
	"time"
	
	mi "github.com/ha1tch/minty"
	"github.com/ha1tch/minty/mintytypes"
)

// Type alias for HTML component functions (used by iterator helpers)
type H = mi.H

// =====================================================
// RE-EXPORTS FROM MINTYTYPES (for backward compatibility)
// =====================================================

// Type re-exports
type Money = mintytypes.Money
type Address = mintytypes.Address
type Customer = mintytypes.Customer
type Status = mintytypes.Status
type BaseStatus = mintytypes.BaseStatus
type ValidationError = mintytypes.ValidationError
type ValidationErrors = mintytypes.ValidationErrors

// Function re-exports
var (
	NewMoney            = mintytypes.NewMoney
	FormatDate          = mintytypes.FormatDate
	DaysAgo             = mintytypes.DaysAgo
	ValidateRequired    = mintytypes.ValidateRequired
	ValidateEmail       = mintytypes.ValidateEmail
	ValidateMoneyAmount = mintytypes.ValidateMoneyAmount
)

// Constant re-exports
const (
	StatusActive    = mintytypes.StatusActive
	StatusInactive  = mintytypes.StatusInactive
	StatusPending   = mintytypes.StatusPending
	StatusCompleted = mintytypes.StatusCompleted
	StatusCancelled = mintytypes.StatusCancelled
	StatusFailed    = mintytypes.StatusFailed
	StatusDraft     = mintytypes.StatusDraft
	StatusPublished = mintytypes.StatusPublished

	AddressBilling  = mintytypes.AddressBilling
	AddressShipping = mintytypes.AddressShipping
	AddressPickup   = mintytypes.AddressPickup
	AddressDelivery = mintytypes.AddressDelivery
	AddressOffice   = mintytypes.AddressOffice
	AddressHome     = mintytypes.AddressHome

	CurrencyUSD = mintytypes.CurrencyUSD
	CurrencyEUR = mintytypes.CurrencyEUR
	CurrencyGBP = mintytypes.CurrencyGBP
	CurrencyJPY = mintytypes.CurrencyJPY
	CurrencyCAD = mintytypes.CurrencyCAD
	CurrencyAUD = mintytypes.CurrencyAUD
)

// =====================================================
// UI HELPER TYPES (kept here, require minty)
// =====================================================

// =====================================================
// FRAMEWORK-AGNOSTIC UI UTILITIES
// =====================================================

// GridLayout creates a CSS grid layout utility
func GridLayout(columns int, gap string) func(content ...mi.H) mi.H {

	return func(content ...mi.H) mi.H {
		return func(b *mi.Builder) mi.Node {
			gridTemplate := fmt.Sprintf("repeat(%d, 1fr)", columns)
			style := fmt.Sprintf("display: grid; grid-template-columns: %s; gap: %s;", gridTemplate, gap)
			
			var children []mi.Node
			for _, item := range content {
				children = append(children, item(b))
			}
			
			return b.Div(mi.Style(style), 
				mi.NewFragment(children...))
		}
	}
}

// FlexLayout creates a flexible box layout utility
func FlexLayout(direction, justify, align string) func(content ...mi.H) mi.H {
	return func(content ...mi.H) mi.H {
		return func(b *mi.Builder) mi.Node {
			style := fmt.Sprintf("display: flex; flex-direction: %s; justify-content: %s; align-items: %s;", 
				direction, justify, align)
			
			var children []mi.Node
			for _, item := range content {
				children = append(children, item(b))
			}
			
			return b.Div(mi.Style(style),
				mi.NewFragment(children...))
		}
	}
}

// CardLayout creates a card-style container
func CardLayout(title string, content ...mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "border: 1px solid #ddd; border-radius: 8px; padding: 16px; margin: 8px 0; background: white; box-shadow: 0 2px 4px rgba(0,0,0,0.1);"
		
		var children []mi.Node
		if title != "" {
			titleStyle := "margin: 0 0 16px 0; font-size: 18px; font-weight: 600; border-bottom: 1px solid #eee; padding-bottom: 8px;"
			children = append(children, b.H3(mi.Style(titleStyle), title))
		}
		
		for _, item := range content {
			children = append(children, item(b))
		}
		
		return b.Div(mi.Style(style),
			mi.NewFragment(children...))
	}
}

// =====================================================
// CONTROL FLOW UTILITIES
// =====================================================

// Each applies a function to each item in a slice and returns HTML nodes
// The returned nodes can be spread directly into NewFragment
func Each[T any](items []T, fn func(T) mi.H) []mi.Node {
	if len(items) == 0 {
		return []mi.Node{}
	}
	nodes := make([]mi.Node, len(items))
	for i, item := range items {
		nodes[i] = fn(item)(mi.B)
	}
	return nodes
}

// EachH applies a function to each item and returns H functions (lazy evaluation)
// Use this with GridLayout and other functions that expect ...H
func EachH[T any](items []T, fn func(T) mi.H) []mi.H {
	if len(items) == 0 {
		return []mi.H{}
	}
	hs := make([]mi.H, len(items))
	for i, item := range items {
		hs[i] = fn(item)
	}
	return hs
}

// NodesToH converts a slice of Nodes to a slice of H functions
func NodesToH(nodes []mi.Node) []mi.H {
	hs := make([]mi.H, len(nodes))
	for i, node := range nodes {
		n := node // capture for closure
		hs[i] = func(b *mi.Builder) mi.Node { return n }
	}
	return hs
}

// If conditionally returns content or empty fragment
func If(condition bool, content mi.H) mi.H {
	if condition {
		return content
	}
	return func(b *mi.Builder) mi.Node { 
		return mi.NewFragment() 
	}
}

// Unless conditionally returns content when condition is false
func Unless(condition bool, content mi.H) mi.H {
	return If(!condition, content)
}

// IfElse conditionally returns one of two content options
func IfElse(condition bool, ifTrue mi.H, ifFalse mi.H) mi.H {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// WrapNode wraps a Node in an H function for use with conditional helpers
func WrapNode(node mi.Node) mi.H {
	return func(b *mi.Builder) mi.Node {
		return node
	}
}

// =====================================================
// ENHANCED SLICE OPERATIONS (Iterator Extensions)
// =====================================================

// Filter returns a new slice containing only elements that match the predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	if slice == nil {
		return nil
	}
	
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms each element using the provided function
func Map[T, U any](slice []T, transform func(T) U) []U {
	if slice == nil {
		return nil
	}
	
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = transform(item)
	}
	return result
}

// Take returns the first n elements
func Take[T any](slice []T, n int) []T {
	if slice == nil || n <= 0 {
		return nil
	}
	
	if n > len(slice) {
		n = len(slice)
	}
	return slice[:n]
}

// Skip returns all elements after the first n
func Skip[T any](slice []T, n int) []T {
	if slice == nil || n <= 0 {
		return slice
	}
	
	if n >= len(slice) {
		return []T{}
	}
	return slice[n:]
}

// Find returns the first element matching predicate, or zero value and false
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element matching predicate, or -1
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// Any returns true if any element matches predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All returns true if all elements match predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Reduce applies an accumulator function over a slice
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// GroupBy groups elements by the result of the key function
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Unique returns a slice with duplicate elements removed
func Unique[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	
	seen := make(map[T]bool, len(slice))
	var result []T
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// UniqueBy returns a slice with duplicate elements removed based on key function
func UniqueBy[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	if slice == nil {
		return nil
	}
	
	seen := make(map[K]bool)
	var result []T
	
	for _, item := range slice {
		key := keyFunc(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	
	return result
}

// Reverse returns a new slice with elements in reverse order
func Reverse[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	
	length := len(slice)
	result := make([]T, length)
	
	for i, item := range slice {
		result[length-1-i] = item
	}
	
	return result
}

// Chain wraps a slice for fluent operations
type Chain[T any] struct {
	data []T
}

// ChainSlice starts a fluent operation chain
func ChainSlice[T any](slice []T) *Chain[T] {
	return &Chain[T]{data: slice}
}

// Filter applies a filter predicate to the chain
func (c *Chain[T]) Filter(predicate func(T) bool) *Chain[T] {
	return &Chain[T]{data: Filter(c.data, predicate)}
}

// Take returns the first n elements in the chain
func (c *Chain[T]) Take(n int) *Chain[T] {
	return &Chain[T]{data: Take(c.data, n)}
}

// Skip skips the first n elements in the chain
func (c *Chain[T]) Skip(n int) *Chain[T] {
	return &Chain[T]{data: Skip(c.data, n)}
}

// Unique removes duplicates from the chain (uses string representation for comparison)
func (c *Chain[T]) Unique() *Chain[T] {
	seen := make(map[string]bool)
	result := make([]T, 0, len(c.data))
	for _, item := range c.data {
		key := fmt.Sprintf("%v", item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return &Chain[T]{data: result}
}

// Reverse reverses the order of elements in the chain
func (c *Chain[T]) Reverse() *Chain[T] {
	return &Chain[T]{data: Reverse(c.data)}
}

// ToSlice returns the final slice result
func (c *Chain[T]) ToSlice() []T {
	return c.data
}

// Map transforms elements in the chain to a different type
func (c *Chain[T]) Map(transform func(T) any) *Chain[any] {
	result := make([]any, len(c.data))
	for i, item := range c.data {
		result[i] = transform(item)
	}
	return &Chain[any]{data: result}
}

// Count returns the number of elements in the chain
func (c *Chain[T]) Count() int {
	return len(c.data)
}

// First returns the first element, or zero value and false if empty
func (c *Chain[T]) First() (T, bool) {
	if len(c.data) == 0 {
		var zero T
		return zero, false
	}
	return c.data[0], true
}

// Last returns the last element, or zero value and false if empty
func (c *Chain[T]) Last() (T, bool) {
	if len(c.data) == 0 {
		var zero T
		return zero, false
	}
	return c.data[len(c.data)-1], true
}

// =====================================================
// HTML-SPECIFIC ITERATOR HELPERS
// =====================================================

// FilterAndRender filters elements and renders them to HTML components
func FilterAndRender[T any](slice []T, predicate func(T) bool, render func(T) H) []H {
	return Map(Filter(slice, predicate), render)
}

// RenderIf conditionally renders elements based on a boolean condition
func RenderIf[T any](slice []T, condition bool, render func(T) H) []H {
	if !condition || len(slice) == 0 {
		return []H{}
	}
	return Map(slice, render)
}

// RenderFirst renders only the first n elements
func RenderFirst[T any](slice []T, n int, render func(T) H) []H {
	return Map(Take(slice, n), render)
}

// RenderWhen renders elements that match a condition
func RenderWhen[T any](slice []T, when func(T) bool, render func(T) H) []H {
	return FilterAndRender(slice, when, render)
}

// EachWithIndex renders elements with their index
func EachWithIndex[T any](slice []T, render func(T, int) H) []H {
	result := make([]H, len(slice))
	for i, item := range slice {
		result[i] = render(item, i)
	}
	return result
}

// Partition splits a slice into two based on a predicate
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	var truthy, falsy []T
	for _, item := range slice {
		if predicate(item) {
			truthy = append(truthy, item)
		} else {
			falsy = append(falsy, item)
		}
	}
	return truthy, falsy
}

// Chunk splits a slice into chunks of specified size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 || len(slice) == 0 {
		return [][]T{}
	}
	
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// ChunkAndRender splits elements into chunks and renders each chunk
func ChunkAndRender[T any](slice []T, size int, renderChunk func([]T) H) []H {
	chunks := Chunk(slice, size)
	return Map(chunks, renderChunk)
}

// =====================================================
// SEMANTIC CSS CLASS HELPERS
// =====================================================


// SemanticClasses provides semantic CSS class naming
type SemanticClasses struct {
	Prefix string
}

func NewSemanticClasses(prefix string) SemanticClasses {
	return SemanticClasses{Prefix: prefix}
}

// Primary action or element
func (sc SemanticClasses) Primary(element string) string {
	return fmt.Sprintf("%s_%s_primary", sc.Prefix, element)
}

// Secondary action or element  
func (sc SemanticClasses) Secondary(element string) string {
	return fmt.Sprintf("%s_%s_secondary", sc.Prefix, element)
}

// Success state
func (sc SemanticClasses) Success(element string) string {
	return fmt.Sprintf("%s_%s_success", sc.Prefix, element)
}

// Warning state
func (sc SemanticClasses) Warning(element string) string {
	return fmt.Sprintf("%s_%s_warning", sc.Prefix, element)
}

// Error/danger state
func (sc SemanticClasses) Error(element string) string {
	return fmt.Sprintf("%s_%s_error", sc.Prefix, element)
}

// Info state
func (sc SemanticClasses) Info(element string) string {
	return fmt.Sprintf("%s_%s_info", sc.Prefix, element)
}

// Container class for layout containers
func (sc SemanticClasses) Container(element string) string {
	return fmt.Sprintf("%s_%s_container", sc.Prefix, element)
}

// Header class for header elements
func (sc SemanticClasses) Header(element string) string {
	return fmt.Sprintf("%s_%s_header", sc.Prefix, element)
}

// Content class for content areas
func (sc SemanticClasses) Content(element string) string {
	return fmt.Sprintf("%s_%s_content", sc.Prefix, element)
}

// Footer class for footer elements
func (sc SemanticClasses) Footer(element string) string {
	return fmt.Sprintf("%s_%s_footer", sc.Prefix, element)
}

// Item class for list items or repeated elements
func (sc SemanticClasses) Item(element string) string {
	return fmt.Sprintf("%s_%s_item", sc.Prefix, element)
}

// =====================================================
// DATE/TIME UTILITIES
// =====================================================

// FormatDateTime formats a datetime string for display
func FormatDateTime(datetime string) string {
	if t, err := time.Parse(time.RFC3339, datetime); err == nil {
		return t.Format("January 2, 2006 at 3:04 PM")
	}
	return datetime // Return original if parsing fails
}

// IsToday checks if a date string represents today
func IsToday(date string) bool {
	if t, err := time.Parse("2006-01-02", date); err == nil {
		today := time.Now().Format("2006-01-02")
		return t.Format("2006-01-02") == today
	}
	return false
}
