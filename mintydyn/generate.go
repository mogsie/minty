package mintydyn

import (
	mi "github.com/ha1tch/minty"
)

// =============================================================================
// MAIN COMPONENT GENERATION
// =============================================================================

// generateComponent builds the complete component based on detected pattern.
func (db *DynamicBuilder[S, D, R]) generateComponent(b *mi.Builder, pattern DetectedPattern) mi.Node {
	theme := db.getTheme()
	
	// Build container class
	containerClass := combineClasses(
		theme.ComponentClass(),
		theme.ComponentPatternClass(pattern.PrimaryPattern),
	)

	// Build container attributes
	containerAttrs := []interface{}{
		mi.ID(db.id),
		mi.Class(containerClass),
		mi.Data("client-managed", "dynamic"),
		mi.Data("pattern", pattern.PrimaryPattern),
	}

	// Add custom attributes
	for key, value := range db.options.CustomAttributes {
		containerAttrs = append(containerAttrs, mi.Data(key, value))
	}

	// Build children
	var children []interface{}

	// Inject theme CSS if provided
	if css := theme.InjectCSS(); css != "" {
		children = append(children, b.Style(mi.Raw(css)))
	}

	// Generate configuration script (JSON data for JS)
	children = append(children, db.generateConfigScript(b, pattern))

	// Generate pattern-specific HTML structure
	structureNodes := db.generatePatternStructure(b, pattern)
	for _, node := range structureNodes {
		children = append(children, node)
	}

	// Generate JavaScript
	children = append(children, mi.Raw(db.generateJavaScript(pattern)))

	// Combine all
	containerAttrs = append(containerAttrs, children...)
	return b.Div(containerAttrs...)
}

// generateConfigScript creates the JSON configuration for client-side JS.
func (db *DynamicBuilder[S, D, R]) generateConfigScript(b *mi.Builder, pattern DetectedPattern) mi.Node {
	config := map[string]interface{}{
		"id":      db.id,
		"pattern": pattern,
		"options": db.options,
	}

	// Add theme classes for JS to use when switching states
	theme := db.getTheme()
	config["themeClasses"] = map[string]string{
		"triggerBase":            theme.StateTriggerClass(),
		"triggerActive":          theme.StateTriggerActiveClass(),
		"triggerDisabled":        theme.StateTriggerDisabledClass(),
		"contentBase":            theme.StateContentClass(),
		"contentActive":          theme.StateContentActiveClass(),
		"contentHidden":          theme.StateContentHiddenClass(),
		"paginationButton":       theme.PaginationButtonClass(),
		"paginationButtonActive": theme.PaginationButtonActiveClass(),
	}

	// Add data based on what's provided
	if pattern.HasStates {
		// Only send state metadata, not content (content is pre-rendered)
		states := db.extractStates()
		stateMeta := make([]map[string]interface{}, len(states))
		for i, s := range states {
			stateMeta[i] = map[string]interface{}{
				"id":        s.ID,
				"label":     s.Label,
				"active":    s.Active,
				"disabled":  s.Disabled,
				"condition": s.Condition,
			}
		}
		config["states"] = stateMeta
	}

	if pattern.HasData {
		config["data"] = db.extractData()
		config["schema"] = db.extractFilterSchema()
		config["filterOptions"] = db.extractFilterOptions()
	}

	if pattern.HasRules {
		config["rules"] = db.extractRules()
	}

	// Add hooks if present
	if db.options.Hooks.BeforeInit != "" ||
		db.options.Hooks.AfterInit != "" ||
		db.options.Hooks.BeforeStateChange != "" ||
		db.options.Hooks.AfterStateChange != "" ||
		db.options.Hooks.OnDestroy != "" ||
		len(db.options.Hooks.StateHooks) > 0 {
		config["hooks"] = db.options.Hooks
	}

	// Add external scripts if present
	if len(db.options.ExternalScripts) > 0 {
		config["externalScripts"] = db.options.ExternalScripts
	}

	// Add external registry if present
	if len(db.options.ExternalRegistry) > 0 {
		config["externalRegistry"] = db.options.ExternalRegistry
	}

	return b.Script(
		mi.Type("application/json"),
		mi.ID(db.id+"-config"),
		mi.Raw(MustJSON(config)),
	)
}

// generatePatternStructure dispatches to pattern-specific generators.
func (db *DynamicBuilder[S, D, R]) generatePatternStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	switch pattern.PrimaryPattern {
	case PatternPreRenderedStates, PatternDynamicStates:
		return db.generateStatesStructure(b, pattern)

	case PatternClientFilterable, PatternServerFilterable:
		return db.generateFilterableStructure(b, pattern)

	case PatternDependencyOnly:
		return db.generateDependencyStructure(b, pattern)

	case PatternStatefulData:
		return db.generateStatefulDataStructure(b, pattern)

	case PatternFilterableStates:
		return db.generateFilterableStatesStructure(b, pattern)

	case PatternDependentStates:
		return db.generateDependentStatesStructure(b, pattern)

	case PatternDependentData:
		return db.generateDependentDataStructure(b, pattern)

	case PatternComplete:
		return db.generateCompleteStructure(b, pattern)

	default:
		return []mi.Node{b.Div(mi.Class("dyn-empty"), "Empty dynamic component")}
	}
}

// =============================================================================
// STATES STRUCTURE (Pattern 1)
// =============================================================================

// generateStatesStructure creates tabs/states HTML.
func (db *DynamicBuilder[S, D, R]) generateStatesStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	states := db.extractStates()
	if len(states) == 0 {
		return []mi.Node{b.Div(mi.Class("dyn-empty"), "No states provided")}
	}

	theme := db.getTheme()
	var children []mi.Node

	// Generate navigation (if more than one state)
	if len(states) > 1 {
		children = append(children, db.generateStateNavigation(b, states, theme))
	}

	// Generate state content panels
	children = append(children, db.generateStateContents(b, states, theme))

	return children
}

// generateStateNavigation creates the tab bar.
func (db *DynamicBuilder[S, D, R]) generateStateNavigation(b *mi.Builder, states []ComponentState, theme DynamicTheme) mi.Node {
	var buttons []interface{}

	for _, state := range states {
		btnClass := theme.StateTriggerClass()
		if state.Active {
			btnClass = combineClasses(btnClass, theme.StateTriggerActiveClass())
		}
		if state.Disabled {
			btnClass = combineClasses(btnClass, theme.StateTriggerDisabledClass())
		}

		btnAttrs := []interface{}{
			mi.Type("button"), // Prevent form submission
			mi.Class(btnClass),
			mi.Data("state-target", state.ID),
			mi.Data("client-action", "switch-state"),
			mi.Attr("role", "tab"),
			mi.Attr("aria-selected", boolStr(state.Active)),
			mi.Attr("aria-controls", "state-"+state.ID),
		}

		if state.Disabled {
			btnAttrs = append(btnAttrs, mi.Disabled())
		}

		// Add icon if present
		var content interface{}
		if state.Icon != "" {
			content = mi.NewFragment(
				mi.Raw(state.Icon+" "),
				mi.Txt(state.Label),
			)
		} else {
			content = state.Label
		}

		btnAttrs = append(btnAttrs, content)
		buttons = append(buttons, b.Button(btnAttrs...))
	}

	navAttrs := []interface{}{
		mi.Class(theme.StateNavigationClass()),
		mi.Attr("role", "tablist"),
	}
	navAttrs = append(navAttrs, buttons...)

	return b.Div(navAttrs...)
}

// generateStateContents creates the content panels.
func (db *DynamicBuilder[S, D, R]) generateStateContents(b *mi.Builder, states []ComponentState, theme DynamicTheme) mi.Node {
	var panels []interface{}

	for _, state := range states {
		panelClass := theme.StateContentClass()
		if state.Active {
			panelClass = combineClasses(panelClass, theme.StateContentActiveClass())
		} else {
			panelClass = combineClasses(panelClass, theme.StateContentHiddenClass())
		}

		panelAttrs := []interface{}{
			mi.ID("state-" + state.ID),
			mi.Class(panelClass),
			mi.Data("state-id", state.ID),
			mi.Attr("role", "tabpanel"),
			mi.Attr("aria-hidden", boolStr(!state.Active)),
		}

		// Add condition as data attribute if present
		if state.Condition != nil {
			panelAttrs = append(panelAttrs, mi.Data("condition", JSONOrEmpty(state.Condition)))
		}

		// Render content
		contentNode := db.renderStateContent(b, state.Content)
		panelAttrs = append(panelAttrs, contentNode)

		panels = append(panels, b.Div(panelAttrs...))
	}

	containerAttrs := []interface{}{mi.Class(theme.StateContainerClass())}
	containerAttrs = append(containerAttrs, panels...)

	return b.Div(containerAttrs...)
}

// renderStateContent converts various content types to a Node.
func (db *DynamicBuilder[S, D, R]) renderStateContent(b *mi.Builder, content interface{}) mi.Node {
	switch c := content.(type) {
	case mi.H:
		return c(b)
	case mi.Node:
		return c
	case string:
		return mi.Txt(c)
	case func(*mi.Builder) mi.Node:
		return c(b)
	default:
		if c == nil {
			return mi.Txt("")
		}
		return mi.Txt(JSONOrEmpty(c))
	}
}

// =============================================================================
// FILTERABLE STRUCTURE (Pattern 2)
// =============================================================================

// generateFilterableStructure creates filtering UI.
func (db *DynamicBuilder[S, D, R]) generateFilterableStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	theme := db.getTheme()
	var children []mi.Node

	// Filter controls
	children = append(children, db.generateFilterControls(b, theme))

	// Results summary
	children = append(children, b.Div(
		mi.ID(db.id+"-summary"),
		mi.Class(theme.ResultsSummaryClass()),
	))

	// Results container
	children = append(children, b.Div(
		mi.ID(db.id+"-results"),
		mi.Class(theme.ResultsClass()),
	))

	// Pagination
	opts := db.extractFilterOptions()
	if opts.EnablePagination {
		children = append(children, b.Div(
			mi.ID(db.id+"-pagination"),
			mi.Class(theme.PaginationClass()),
		))
	}

	return children
}

// generateFilterControls creates the filter input fields.
func (db *DynamicBuilder[S, D, R]) generateFilterControls(b *mi.Builder, theme DynamicTheme) mi.Node {
	schema := db.extractFilterSchema()
	if len(schema.Fields) == 0 {
		// Auto-generate schema from data if not provided
		schema = db.generateSchemaFromData()
	}

	var controls []interface{}
	for _, field := range schema.Fields {
		controls = append(controls, db.generateFilterControl(b, field, theme))
	}

	containerAttrs := []interface{}{
		mi.ID(db.id + "-filters"),
		mi.Class(theme.FilterControlsClass()),
	}
	containerAttrs = append(containerAttrs, controls...)

	return b.Div(containerAttrs...)
}

// generateFilterControl creates a single filter control based on field type.
func (db *DynamicBuilder[S, D, R]) generateFilterControl(b *mi.Builder, field FilterableField, theme DynamicTheme) mi.Node {
	var control mi.Node

	switch field.Type {
	case "text":
		control = b.Input(
			mi.Type("text"),
			mi.ID(db.id+"-filter-"+field.Name),
			mi.Class(theme.FilterInputClass()),
			mi.Data("filter-field", field.Name),
			mi.Data("filter-type", "text"),
			mi.Placeholder("Search "+field.Label+"..."),
		)

	case "select":
		var options []interface{}
		options = append(options, b.Option(mi.Value(""), "All"))
		for _, opt := range field.Options {
			options = append(options, b.Option(mi.Value(opt), opt))
		}
		control = b.Select(append([]interface{}{
			mi.ID(db.id + "-filter-" + field.Name),
			mi.Class(theme.FilterSelectClass()),
			mi.Data("filter-field", field.Name),
			mi.Data("filter-type", "select"),
		}, options...)...)

	case "multiselect":
		var checkboxes []interface{}
		for _, opt := range field.Options {
			checkboxes = append(checkboxes, b.Label(
				mi.Class("dyn-checkbox-label"),
				b.Input(
					mi.Type("checkbox"),
					mi.Value(opt),
					mi.Class(theme.FilterCheckboxClass()),
					mi.Data("filter-field", field.Name),
					mi.Data("filter-type", "multiselect"),
				),
				" "+opt,
			))
		}
		control = b.Div(append([]interface{}{mi.Class("dyn-multiselect-control")}, checkboxes...)...)

	case "boolean":
		control = b.Input(
			mi.Type("checkbox"),
			mi.ID(db.id+"-filter-"+field.Name),
			mi.Class(theme.FilterCheckboxClass()),
			mi.Data("filter-field", field.Name),
			mi.Data("filter-type", "boolean"),
		)

	case "range":
		if field.Range != nil {
			control = b.Div(
				mi.Class("dyn-range-control"),
				b.Input(
					mi.Type("range"),
					mi.ID(db.id+"-filter-"+field.Name+"-min"),
					mi.Class(theme.FilterRangeClass()),
					mi.Attr("min", floatStr(field.Range.Min)),
					mi.Attr("max", floatStr(field.Range.Max)),
					mi.Attr("step", floatStr(field.Range.Step)),
					mi.Data("filter-field", field.Name),
					mi.Data("filter-type", "range-min"),
				),
				b.Input(
					mi.Type("range"),
					mi.ID(db.id+"-filter-"+field.Name+"-max"),
					mi.Class(theme.FilterRangeClass()),
					mi.Attr("min", floatStr(field.Range.Min)),
					mi.Attr("max", floatStr(field.Range.Max)),
					mi.Attr("step", floatStr(field.Range.Step)),
					mi.Data("filter-field", field.Name),
					mi.Data("filter-type", "range-max"),
				),
			)
		} else {
			control = mi.Txt("Range configuration missing")
		}

	default:
		control = mi.Txt("Unknown filter type: " + field.Type)
	}

	return b.Div(
		mi.Class(theme.FilterGroupClass()),
		b.Label(
			mi.Class(theme.FilterLabelClass()),
			mi.For(db.id+"-filter-"+field.Name),
			field.Label,
		),
		control,
	)
}

// generateSchemaFromData infers a filter schema from the data structure.
func (db *DynamicBuilder[S, D, R]) generateSchemaFromData() FilterSchema {
	data := db.extractData()
	if len(data) == 0 {
		return FilterSchema{}
	}

	// Analyze first item to determine field types
	first := data[0]
	var fields []FilterableField

	for key, value := range first {
		field := FilterableField{
			Name:  key,
			Label: key, // Could be improved with title case conversion
		}

		switch value.(type) {
		case bool:
			field.Type = "boolean"
		case float64, float32, int, int64, int32:
			field.Type = "range"
			// Could analyze all data to find min/max
		case string:
			field.Type = "text"
			field.Searchable = true
		default:
			field.Type = "text"
		}

		fields = append(fields, field)
	}

	return FilterSchema{Fields: fields}
}

// =============================================================================
// DEPENDENCY STRUCTURE (Pattern 3)
// =============================================================================

// generateDependencyStructure creates dependency-only UI.
func (db *DynamicBuilder[S, D, R]) generateDependencyStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	// Dependency rules don't generate visible structure
	// They just inject rules into the config and JS handles it
	return []mi.Node{
		b.Div(
			mi.ID(db.id+"-dependency-status"),
			mi.Class("dyn-dependency-status "+db.getTheme().HiddenClass()),
		),
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func floatStr(v float64) string {
	return JSONOrEmpty(v)
}
