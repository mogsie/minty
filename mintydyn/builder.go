package mintydyn

import (
	"reflect"

	mi "github.com/ha1tch/minty"
)

// =============================================================================
// GENERIC DYNAMIC BUILDER
// =============================================================================

// DynamicBuilder constructs dynamic components with automatic pattern detection.
// Type parameters S, D, R constrain what can be provided for states, data, and rules.
type DynamicBuilder[S States, D Data, R Rules] struct {
	id       string
	states   S
	data     D
	rules    R
	renderer ComponentRenderer
	theme    DynamicTheme
	options  DynamicOptions
}

// =============================================================================
// CONSTRUCTORS
// =============================================================================

// New creates a new DynamicBuilder with the given ID.
// The type parameters determine what patterns are available.
func New[S States, D Data, R Rules](id string) *DynamicBuilder[S, D, R] {
	return &DynamicBuilder[S, D, R]{
		id:      id,
		options: DefaultOptions(),
	}
}

// DefaultOptions returns sensible default options.
func DefaultOptions() DynamicOptions {
	return DynamicOptions{
		MaxStates:       50,
		MaxDataSize:     1000,
		FallbackPattern: PatternEmpty,
		PerformanceMode: "balanced",
	}
}

// =============================================================================
// FLUENT BUILDER METHODS
// =============================================================================

// WithStates sets the states/tabs for the component.
func (db *DynamicBuilder[S, D, R]) WithStates(states S) *DynamicBuilder[S, D, R] {
	db.states = states
	return db
}

// WithData sets the filterable data for the component.
func (db *DynamicBuilder[S, D, R]) WithData(data D) *DynamicBuilder[S, D, R] {
	db.data = data
	return db
}

// WithRules sets the dependency rules for the component.
func (db *DynamicBuilder[S, D, R]) WithRules(rules R) *DynamicBuilder[S, D, R] {
	db.rules = rules
	return db
}

// WithRenderer sets the component renderer for data items.
func (db *DynamicBuilder[S, D, R]) WithRenderer(renderer ComponentRenderer) *DynamicBuilder[S, D, R] {
	db.renderer = renderer
	return db
}

// WithTheme sets the theme for styling.
func (db *DynamicBuilder[S, D, R]) WithTheme(theme DynamicTheme) *DynamicBuilder[S, D, R] {
	db.theme = theme
	return db
}

// WithOptions sets all options at once.
func (db *DynamicBuilder[S, D, R]) WithOptions(options DynamicOptions) *DynamicBuilder[S, D, R] {
	db.options = options
	return db
}

// Minified enables JavaScript minification for smaller output.
func (db *DynamicBuilder[S, D, R]) Minified() *DynamicBuilder[S, D, R] {
	db.options.MinifyJS = true
	return db
}

// =============================================================================
// EXTERNAL SCRIPT METHODS
// =============================================================================

// WithExternalScript adds an external script dependency.
func (db *DynamicBuilder[S, D, R]) WithExternalScript(src string, opts ...ScriptOption) *DynamicBuilder[S, D, R] {
	script := ExternalScript{Src: src}
	for _, opt := range opts {
		opt(&script)
	}
	db.options.ExternalScripts = append(db.options.ExternalScripts, script)
	return db
}

// WithExternal reserves a name in the externals registry.
func (db *DynamicBuilder[S, D, R]) WithExternal(name string) *DynamicBuilder[S, D, R] {
	db.options.ExternalRegistry = append(db.options.ExternalRegistry, name)
	return db
}

// ScriptOption configures an external script.
type ScriptOption func(*ExternalScript)

// Required marks the script as required (blocks init until loaded).
func Required() ScriptOption {
	return func(s *ExternalScript) { s.Required = true }
}

// Async sets the async attribute on the script tag.
func Async() ScriptOption {
	return func(s *ExternalScript) { s.Async = true }
}

// Defer sets the defer attribute on the script tag.
func ScriptDefer() ScriptOption {
	return func(s *ExternalScript) { s.Defer = true }
}

// OnLoad sets JS code to run when the script loads.
func OnLoad(jsCode string) ScriptOption {
	return func(s *ExternalScript) { s.OnLoad = jsCode }
}

// =============================================================================
// LIFECYCLE HOOK METHODS
// =============================================================================

// OnInit sets the afterInit hook (runs after component initializes).
func (db *DynamicBuilder[S, D, R]) OnInit(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.AfterInit = jsCode
	return db
}

// BeforeInit sets the beforeInit hook (runs before component initializes).
func (db *DynamicBuilder[S, D, R]) BeforeInit(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.BeforeInit = jsCode
	return db
}

// OnStateChange sets the afterStateChange hook.
func (db *DynamicBuilder[S, D, R]) OnStateChange(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.AfterStateChange = jsCode
	return db
}

// BeforeStateChange sets the beforeStateChange hook (can cancel by returning false).
func (db *DynamicBuilder[S, D, R]) BeforeStateChange(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.BeforeStateChange = jsCode
	return db
}

// OnState sets a hook for a specific state becoming active.
func (db *DynamicBuilder[S, D, R]) OnState(stateID, jsCode string) *DynamicBuilder[S, D, R] {
	if db.options.Hooks.StateHooks == nil {
		db.options.Hooks.StateHooks = make(map[string]string)
	}
	db.options.Hooks.StateHooks[stateID] = jsCode
	return db
}

// OnFilter sets the afterFilter hook.
func (db *DynamicBuilder[S, D, R]) OnFilter(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.AfterFilter = jsCode
	return db
}

// OnDestroy sets the onDestroy hook (cleanup).
func (db *DynamicBuilder[S, D, R]) OnDestroy(jsCode string) *DynamicBuilder[S, D, R] {
	db.options.Hooks.OnDestroy = jsCode
	return db
}

// =============================================================================
// BUILD
// =============================================================================

// Build generates the final component as a minty H function.
func (db *DynamicBuilder[S, D, R]) Build() mi.H {
	return func(b *mi.Builder) mi.Node {
		pattern := db.detectPattern()
		return db.generateComponent(b, pattern)
	}
}

// ID returns the component's ID.
func (db *DynamicBuilder[S, D, R]) ID() string {
	return db.id
}

// =============================================================================
// PATTERN DETECTION
// =============================================================================

// detectPattern analyzes the provided data to determine the optimal pattern.
func (db *DynamicBuilder[S, D, R]) detectPattern() DetectedPattern {
	pattern := DetectedPattern{}

	// Use reflection to check if values are actually provided
	statesVal := reflect.ValueOf(db.states)
	dataVal := reflect.ValueOf(db.data)
	rulesVal := reflect.ValueOf(db.rules)

	// Check states
	if !statesVal.IsZero() {
		pattern.HasStates = true
		pattern.StateCount = db.getStateCount()
	}

	// Check data
	if !dataVal.IsZero() {
		pattern.HasData = true
		pattern.DataSize = db.getDataSize()
	}

	// Check rules
	if !rulesVal.IsZero() {
		pattern.HasRules = true
		pattern.RuleCount = db.getRuleCount()
	}

	// Check renderer
	if db.renderer != nil {
		pattern.HasRenderer = true
	}

	// Determine primary pattern
	pattern.PrimaryPattern = db.determinePrimaryPattern(pattern)

	// Add optimizations based on size and complexity
	pattern.Optimizations = db.determineOptimizations(pattern)

	return pattern
}

// determinePrimaryPattern selects the best pattern based on what's provided.
func (db *DynamicBuilder[S, D, R]) determinePrimaryPattern(p DetectedPattern) string {
	switch {
	case p.HasStates && p.HasData && p.HasRules:
		return PatternComplete

	case p.HasStates && p.HasData:
		if p.DataSize > 100 {
			return PatternFilterableStates // Data-heavy, use filtering with state context
		}
		return PatternStatefulData // States-heavy, show data within states

	case p.HasStates && p.HasRules:
		return PatternDependentStates

	case p.HasData && p.HasRules:
		return PatternDependentData

	case p.HasStates:
		if p.StateCount <= 10 {
			return PatternPreRenderedStates // Pre-render all states
		}
		return PatternDynamicStates // Too many states, use dynamic management

	case p.HasData:
		if p.DataSize <= 50 {
			return PatternClientFilterable // Client-side filtering
		}
		return PatternServerFilterable // Too much data for client

	case p.HasRules:
		return PatternDependencyOnly

	default:
		return PatternEmpty
	}
}

// determineOptimizations suggests optimizations based on data characteristics.
func (db *DynamicBuilder[S, D, R]) determineOptimizations(p DetectedPattern) []string {
	var optimizations []string

	// Performance optimizations
	if p.DataSize > 500 {
		optimizations = append(optimizations, "pagination")
	}
	if p.StateCount > 20 {
		optimizations = append(optimizations, "lazy-state-loading")
	}
	if p.RuleCount > 50 {
		optimizations = append(optimizations, "rule-grouping")
	}

	// Coordination optimizations
	if p.HasStates && p.HasData {
		optimizations = append(optimizations, "state-data-coordination")
	}
	if p.HasRules && (p.HasStates || p.HasData) {
		optimizations = append(optimizations, "rule-coordination")
	}

	return optimizations
}

// =============================================================================
// COUNT HELPERS
// =============================================================================

// getStateCount returns the number of states.
func (db *DynamicBuilder[S, D, R]) getStateCount() int {
	switch states := any(db.states).(type) {
	case []ComponentState:
		return len(states)
	case map[string]ComponentState:
		return len(states)
	case ComponentStateCollection:
		return len(states.States)
	default:
		return 0
	}
}

// getDataSize returns the number of data items.
func (db *DynamicBuilder[S, D, R]) getDataSize() int {
	switch data := any(db.data).(type) {
	case []map[string]interface{}:
		return len(data)
	case []interface{}:
		return len(data)
	case FilterableDataset:
		return len(data.Items)
	case DataCollection:
		return len(data.Items)
	default:
		return 0
	}
}

// getRuleCount returns the number of rules.
func (db *DynamicBuilder[S, D, R]) getRuleCount() int {
	switch rules := any(db.rules).(type) {
	case []DependencyRule:
		return len(rules)
	case RuleCollection:
		return len(rules.Rules)
	default:
		return 0
	}
}

// =============================================================================
// DATA EXTRACTION HELPERS
// =============================================================================

// extractStates converts the generic states to a slice.
func (db *DynamicBuilder[S, D, R]) extractStates() []ComponentState {
	switch states := any(db.states).(type) {
	case []ComponentState:
		return states
	case map[string]ComponentState:
		result := make([]ComponentState, 0, len(states))
		for _, state := range states {
			result = append(result, state)
		}
		return result
	case ComponentStateCollection:
		return states.States
	default:
		return nil
	}
}

// extractData converts the generic data to a slice.
func (db *DynamicBuilder[S, D, R]) extractData() []map[string]interface{} {
	switch data := any(db.data).(type) {
	case []map[string]interface{}:
		return data
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(data))
		for _, item := range data {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	case FilterableDataset:
		return data.Items
	case DataCollection:
		return data.Items
	default:
		return nil
	}
}

// extractRules converts the generic rules to a slice.
func (db *DynamicBuilder[S, D, R]) extractRules() []DependencyRule {
	switch rules := any(db.rules).(type) {
	case []DependencyRule:
		return rules
	case RuleCollection:
		return rules.Rules
	default:
		return nil
	}
}

// extractFilterSchema gets the filter schema from data.
func (db *DynamicBuilder[S, D, R]) extractFilterSchema() FilterSchema {
	switch data := any(db.data).(type) {
	case FilterableDataset:
		return data.Schema
	case DataCollection:
		return data.Schema
	default:
		return FilterSchema{}
	}
}

// extractFilterOptions gets the filter options from data.
func (db *DynamicBuilder[S, D, R]) extractFilterOptions() FilterOptions {
	switch data := any(db.data).(type) {
	case FilterableDataset:
		return data.Options
	default:
		return FilterOptions{}
	}
}
