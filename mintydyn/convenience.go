package mintydyn

import (
	mi "github.com/ha1tch/minty"
)

// =============================================================================
// SINGLE-PATTERN CONVENIENCE FUNCTIONS
// =============================================================================

// Tabs creates a simple tabbed interface.
//
//	mdy.Tabs("profile", []mdy.ComponentState{
//	    {ID: "info", Label: "Info", Active: true, Content: infoContent},
//	    {ID: "settings", Label: "Settings", Content: settingsContent},
//	})
func Tabs(id string, states []ComponentState) mi.H {
	return New[[]ComponentState, []map[string]interface{}, []DependencyRule](id).
		WithStates(states).
		Build()
}

// TabsWithTheme creates tabs with a specific theme.
func TabsWithTheme(id string, states []ComponentState, theme DynamicTheme) mi.H {
	return New[[]ComponentState, []map[string]interface{}, []DependencyRule](id).
		WithStates(states).
		WithTheme(theme).
		Build()
}

// Filter creates a filterable data view.
//
//	mdy.Filter("products", productData, mdy.FilterSchema{
//	    Fields: []mdy.FilterableField{
//	        {Name: "category", Type: "select", Label: "Category", Options: categories},
//	    },
//	})
func Filter(id string, data []map[string]interface{}, schema FilterSchema) mi.H {
	dataset := FilterableDataset{
		Items:  data,
		Schema: schema,
	}
	return New[[]ComponentState, FilterableDataset, []DependencyRule](id).
		WithData(dataset).
		Build()
}

// FilterWithOptions creates a filterable view with custom options.
func FilterWithOptions(id string, data []map[string]interface{}, schema FilterSchema, opts FilterOptions) mi.H {
	dataset := FilterableDataset{
		Items:   data,
		Schema:  schema,
		Options: opts,
	}
	return New[[]ComponentState, FilterableDataset, []DependencyRule](id).
		WithData(dataset).
		Build()
}

// Form creates a form with dependency rules.
//
//	mdy.Form("insurance", []mdy.DependencyRule{
//	    {
//	        ID: "show-spouse",
//	        Trigger: mdy.TriggerCondition{ComponentID: "marital-status", Event: "change", Condition: "equals", Value: "married"},
//	        Actions: []mdy.DependencyAction{{TargetID: "spouse-section", Action: "show"}},
//	    },
//	})
func Form(id string, rules []DependencyRule) mi.H {
	return New[[]ComponentState, []map[string]interface{}, []DependencyRule](id).
		WithRules(rules).
		Build()
}

// =============================================================================
// TWO-PATTERN CONVENIENCE FUNCTIONS
// =============================================================================

// TabsWithData creates tabs where each tab shows filtered data.
func TabsWithData(id string, states []ComponentState, data []map[string]interface{}, schema FilterSchema) mi.H {
	dataset := FilterableDataset{
		Items:  data,
		Schema: schema,
	}
	return New[[]ComponentState, FilterableDataset, []DependencyRule](id).
		WithStates(states).
		WithData(dataset).
		Build()
}

// TabsWithRules creates tabs controlled by dependency rules.
func TabsWithRules(id string, states []ComponentState, rules []DependencyRule) mi.H {
	return New[[]ComponentState, []map[string]interface{}, []DependencyRule](id).
		WithStates(states).
		WithRules(rules).
		Build()
}

// FilterWithRules creates a filterable view with dependency rules.
func FilterWithRules(id string, data []map[string]interface{}, schema FilterSchema, rules []DependencyRule) mi.H {
	dataset := FilterableDataset{
		Items:  data,
		Schema: schema,
	}
	return New[[]ComponentState, FilterableDataset, []DependencyRule](id).
		WithData(dataset).
		WithRules(rules).
		Build()
}

// =============================================================================
// FULL COMBINATION
// =============================================================================

// Dynamic creates a component with all three patterns.
func Dynamic(id string, states []ComponentState, data []map[string]interface{}, schema FilterSchema, rules []DependencyRule) mi.H {
	dataset := FilterableDataset{
		Items:  data,
		Schema: schema,
	}
	return New[[]ComponentState, FilterableDataset, []DependencyRule](id).
		WithStates(states).
		WithData(dataset).
		WithRules(rules).
		Build()
}

// =============================================================================
// FLEXIBLE BUILDER (for advanced use)
// =============================================================================

// Dyn starts a flexible builder chain.
// Use when the convenience functions don't fit your needs.
//
//	mdy.Dyn("complex").
//	    States(myStates).
//	    Data(myData).
//	    Rules(myRules).
//	    OnInit("this.registerExternal('map', new google.maps.Map(...))").
//	    Build()
func Dyn(id string) *FlexBuilder {
	return &FlexBuilder{
		id:            id,
		options:       DefaultOptions(),
		filterOptions: FilterOptions{},
	}
}

// FlexBuilder provides a flexible, runtime-typed builder.
type FlexBuilder struct {
	id            string
	states        interface{}
	data          interface{}
	rules         interface{}
	renderer      ComponentRenderer
	theme         DynamicTheme
	options       DynamicOptions
	filterOptions FilterOptions
	filterSchema  FilterSchema
}

// States sets the states (validates at build time).
func (fb *FlexBuilder) States(s interface{}) *FlexBuilder {
	fb.states = s
	return fb
}

// Data sets the data (validates at build time).
func (fb *FlexBuilder) Data(d interface{}) *FlexBuilder {
	fb.data = d
	return fb
}

// Rules sets the rules (validates at build time).
func (fb *FlexBuilder) Rules(r interface{}) *FlexBuilder {
	fb.rules = r
	return fb
}

// ServerRenderedData configures filtering for pre-rendered HTML rows.
// rowSelector: CSS selector for data rows (e.g., ".asset-row")
// counterSelector: CSS selector for count display (e.g., "#asset-count")
func (fb *FlexBuilder) ServerRenderedData(rowSelector, counterSelector string) *FlexBuilder {
	fb.filterOptions.ServerRendered = true
	fb.filterOptions.RowSelector = rowSelector
	fb.filterOptions.CounterSelector = counterSelector
	return fb
}

// FilterField adds a filter field to the schema.
func (fb *FlexBuilder) FilterField(field FilterableField) *FlexBuilder {
	fb.filterSchema.Fields = append(fb.filterSchema.Fields, field)
	return fb
}

// TextFilter adds a text search filter field.
func (fb *FlexBuilder) TextFilter(name, label string) *FlexBuilder {
	fb.filterSchema.Fields = append(fb.filterSchema.Fields, FilterableField{
		Name:       name,
		Type:       "text",
		Label:      label,
		Searchable: true,
	})
	return fb
}

// SelectFilter adds a select dropdown filter field.
func (fb *FlexBuilder) SelectFilter(name, label string, options []string) *FlexBuilder {
	fb.filterSchema.Fields = append(fb.filterSchema.Fields, FilterableField{
		Name:    name,
		Type:    "select",
		Label:   label,
		Options: options,
	})
	return fb
}

// Renderer sets the component renderer.
func (fb *FlexBuilder) Renderer(r ComponentRenderer) *FlexBuilder {
	fb.renderer = r
	return fb
}

// Theme sets the theme.
func (fb *FlexBuilder) Theme(t DynamicTheme) *FlexBuilder {
	fb.theme = t
	return fb
}

// Options sets all options.
func (fb *FlexBuilder) Options(o DynamicOptions) *FlexBuilder {
	fb.options = o
	return fb
}

// ExternalScript adds an external script.
func (fb *FlexBuilder) ExternalScript(src string, opts ...ScriptOption) *FlexBuilder {
	script := ExternalScript{Src: src}
	for _, opt := range opts {
		opt(&script)
	}
	fb.options.ExternalScripts = append(fb.options.ExternalScripts, script)
	return fb
}

// External reserves a name in the externals registry.
func (fb *FlexBuilder) External(name string) *FlexBuilder {
	fb.options.ExternalRegistry = append(fb.options.ExternalRegistry, name)
	return fb
}

// OnInit sets the afterInit hook.
func (fb *FlexBuilder) OnInit(jsCode string) *FlexBuilder {
	fb.options.Hooks.AfterInit = jsCode
	return fb
}

// BeforeInit sets the beforeInit hook.
func (fb *FlexBuilder) BeforeInit(jsCode string) *FlexBuilder {
	fb.options.Hooks.BeforeInit = jsCode
	return fb
}

// OnStateChange sets the afterStateChange hook.
func (fb *FlexBuilder) OnStateChange(jsCode string) *FlexBuilder {
	fb.options.Hooks.AfterStateChange = jsCode
	return fb
}

// BeforeStateChange sets the beforeStateChange hook.
func (fb *FlexBuilder) BeforeStateChange(jsCode string) *FlexBuilder {
	fb.options.Hooks.BeforeStateChange = jsCode
	return fb
}

// OnState sets a per-state hook.
func (fb *FlexBuilder) OnState(stateID, jsCode string) *FlexBuilder {
	if fb.options.Hooks.StateHooks == nil {
		fb.options.Hooks.StateHooks = make(map[string]string)
	}
	fb.options.Hooks.StateHooks[stateID] = jsCode
	return fb
}

// OnFilter sets the afterFilter hook.
func (fb *FlexBuilder) OnFilter(jsCode string) *FlexBuilder {
	fb.options.Hooks.AfterFilter = jsCode
	return fb
}

// OnDestroy sets the onDestroy hook.
func (fb *FlexBuilder) OnDestroy(jsCode string) *FlexBuilder {
	fb.options.Hooks.OnDestroy = jsCode
	return fb
}

// Minified enables JavaScript minification for smaller output.
func (fb *FlexBuilder) Minified() *FlexBuilder {
	fb.options.MinifyJS = true
	return fb
}

// Build creates the component.
func (fb *FlexBuilder) Build() mi.H {
	// Convert to the appropriate generic builder based on what's provided
	// This uses type assertions and falls back to sensible defaults

	var states []ComponentState
	var data FilterableDataset
	var rules []DependencyRule

	// Extract states
	switch s := fb.states.(type) {
	case []ComponentState:
		states = s
	case ComponentStateCollection:
		states = s.States
	case map[string]ComponentState:
		for _, state := range s {
			states = append(states, state)
		}
	}

	// Extract data
	switch d := fb.data.(type) {
	case []map[string]interface{}:
		data = FilterableDataset{Items: d}
	case FilterableDataset:
		data = d
	case DataCollection:
		data = FilterableDataset{Items: d.Items, Schema: d.Schema}
	}

	// Merge filterOptions from FlexBuilder
	if fb.filterOptions.ServerRendered || fb.filterOptions.RowSelector != "" {
		data.Options = fb.filterOptions
	}

	// Merge filterSchema from FlexBuilder
	if len(fb.filterSchema.Fields) > 0 {
		data.Schema = fb.filterSchema
	}

	// Extract rules
	switch r := fb.rules.(type) {
	case []DependencyRule:
		rules = r
	case RuleCollection:
		rules = r.Rules
	}

	// Build with the generic builder
	builder := New[[]ComponentState, FilterableDataset, []DependencyRule](fb.id).
		WithOptions(fb.options)

	if len(states) > 0 {
		builder = builder.WithStates(states)
	}
	// Include data if we have items, schema fields, or server-rendered mode configured
	if len(data.Items) > 0 || len(data.Schema.Fields) > 0 || data.Options.ServerRendered {
		builder = builder.WithData(data)
	}
	if len(rules) > 0 {
		builder = builder.WithRules(rules)
	}
	if fb.renderer != nil {
		builder = builder.WithRenderer(fb.renderer)
	}
	if fb.theme != nil {
		builder = builder.WithTheme(fb.theme)
	}

	return builder.Build()
}

// =============================================================================
// STATE HELPERS
// =============================================================================

// NewState creates a ComponentState with the given ID and label.
func NewState(id, label string, content interface{}) ComponentState {
	return ComponentState{
		ID:      id,
		Label:   label,
		Content: content,
	}
}

// ActiveState creates an active ComponentState.
func ActiveState(id, label string, content interface{}) ComponentState {
	return ComponentState{
		ID:      id,
		Label:   label,
		Content: content,
		Active:  true,
	}
}

// =============================================================================
// RULE HELPERS
// =============================================================================

// ShowWhen creates a rule that shows a target when a condition is met.
func ShowWhen(triggerID, condition string, value interface{}, targetID string) DependencyRule {
	return DependencyRule{
		ID: "show-" + targetID + "-when-" + triggerID,
		Trigger: TriggerCondition{
			ComponentID: triggerID,
			Event:       "change",
			Condition:   condition,
			Value:       value,
		},
		Actions: []DependencyAction{
			{TargetID: targetID, Action: "show"},
		},
	}
}

// HideWhen creates a rule that hides a target when a condition is met.
func HideWhen(triggerID, condition string, value interface{}, targetID string) DependencyRule {
	return DependencyRule{
		ID: "hide-" + targetID + "-when-" + triggerID,
		Trigger: TriggerCondition{
			ComponentID: triggerID,
			Event:       "change",
			Condition:   condition,
			Value:       value,
		},
		Actions: []DependencyAction{
			{TargetID: targetID, Action: "hide"},
		},
	}
}

// EnableWhen creates a rule that enables a target when a condition is met.
func EnableWhen(triggerID, condition string, value interface{}, targetID string) DependencyRule {
	return DependencyRule{
		ID: "enable-" + targetID + "-when-" + triggerID,
		Trigger: TriggerCondition{
			ComponentID: triggerID,
			Event:       "change",
			Condition:   condition,
			Value:       value,
		},
		Actions: []DependencyAction{
			{TargetID: targetID, Action: "enable"},
		},
	}
}

// DisableWhen creates a rule that disables a target when a condition is met.
func DisableWhen(triggerID, condition string, value interface{}, targetID string) DependencyRule {
	return DependencyRule{
		ID: "disable-" + targetID + "-when-" + triggerID,
		Trigger: TriggerCondition{
			ComponentID: triggerID,
			Event:       "change",
			Condition:   condition,
			Value:       value,
		},
		Actions: []DependencyAction{
			{TargetID: targetID, Action: "disable"},
		},
	}
}

// =============================================================================
// FILTER HELPERS
// =============================================================================

// TextField creates a text filter field.
func TextField(name, label string) FilterableField {
	return FilterableField{
		Name:       name,
		Type:       "text",
		Label:      label,
		Searchable: true,
	}
}

// SelectField creates a select filter field.
func SelectField(name, label string, options []string) FilterableField {
	return FilterableField{
		Name:    name,
		Type:    "select",
		Label:   label,
		Options: options,
	}
}

// MultiSelectField creates a multi-select filter field.
func MultiSelectField(name, label string, options []string) FilterableField {
	return FilterableField{
		Name:    name,
		Type:    "multiselect",
		Label:   label,
		Options: options,
	}
}

// BoolField creates a boolean filter field.
func BoolField(name, label string) FilterableField {
	return FilterableField{
		Name:  name,
		Type:  "boolean",
		Label: label,
	}
}

// RangeField creates a range filter field.
func RangeField(name, label string, min, max, step float64) FilterableField {
	return FilterableField{
		Name:  name,
		Type:  "range",
		Label: label,
		Range: &RangeInfo{Min: min, Max: max, Step: step},
	}
}
