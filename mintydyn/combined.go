package mintydyn

import (
	mi "github.com/ha1tch/minty"
)

// =============================================================================
// STATEFUL DATA STRUCTURE (States + Data)
// =============================================================================

// generateStatefulDataStructure creates tabs with data content per tab.
func (db *DynamicBuilder[S, D, R]) generateStatefulDataStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	states := db.extractStates()
	theme := db.getTheme()
	var children []mi.Node

	// State navigation
	if len(states) > 1 {
		children = append(children, db.generateStateNavigation(b, states, theme))
	}

	// Each state contains its own data context
	var stateContents []interface{}
	for _, state := range states {
		panelClass := combineClasses(
			theme.StateContentClass(),
			"dyn-filterable-state",
		)
		if state.Active {
			panelClass = combineClasses(panelClass, theme.StateContentActiveClass())
		} else {
			panelClass = combineClasses(panelClass, theme.StateContentHiddenClass())
		}

		panel := b.Div(
			mi.ID("state-"+state.ID),
			mi.Class(panelClass),
			mi.Data("state-id", state.ID),
			// Filters scoped to this state
			b.Div(
				mi.Class("dyn-state-filters"),
				mi.Data("state-context", state.ID),
			),
			// Results scoped to this state
			b.Div(
				mi.Class("dyn-state-results"),
				mi.Data("state-context", state.ID),
			),
		)

		stateContents = append(stateContents, panel)
	}

	containerAttrs := []interface{}{mi.Class("dyn-stateful-data-container")}
	containerAttrs = append(containerAttrs, stateContents...)
	children = append(children, b.Div(containerAttrs...))

	return children
}

// =============================================================================
// FILTERABLE STATES STRUCTURE (States + Data, data-heavy)
// =============================================================================

// generateFilterableStatesStructure creates filtering with state context.
func (db *DynamicBuilder[S, D, R]) generateFilterableStatesStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	states := db.extractStates()
	theme := db.getTheme()
	var children []mi.Node

	// State navigation as filter context switcher
	if len(states) > 1 {
		children = append(children, db.generateStateNavigation(b, states, theme))
	}

	// Shared filter controls
	children = append(children, db.generateFilterControls(b, theme))

	// Results summary
	children = append(children, b.Div(
		mi.ID(db.id+"-summary"),
		mi.Class(theme.ResultsSummaryClass()),
	))

	// State-contextualized results
	var stateResults []interface{}
	for _, state := range states {
		resultClass := "dyn-state-results"
		if state.Active {
			resultClass = combineClasses(resultClass, theme.StateContentActiveClass())
		} else {
			resultClass = combineClasses(resultClass, theme.StateContentHiddenClass())
		}

		stateResults = append(stateResults, b.Div(
			mi.ID("results-"+state.ID),
			mi.Class(resultClass),
			mi.Data("state-context", state.ID),
		))
	}

	resultsContainer := []interface{}{
		mi.ID(db.id + "-results"),
		mi.Class(theme.ResultsClass()),
	}
	resultsContainer = append(resultsContainer, stateResults...)
	children = append(children, b.Div(resultsContainer...))

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

// =============================================================================
// DEPENDENT STATES STRUCTURE (States + Rules)
// =============================================================================

// generateDependentStatesStructure creates tabs controlled by rules.
func (db *DynamicBuilder[S, D, R]) generateDependentStatesStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	states := db.extractStates()
	rules := db.extractRules()
	theme := db.getTheme()
	var children []mi.Node

	// Enhanced navigation with dependency indicators
	children = append(children, db.generateDependentStateNavigation(b, states, rules, theme))

	// State contents with dependency attributes
	children = append(children, db.generateDependentStateContents(b, states, rules, theme))

	return children
}

// generateDependentStateNavigation creates navigation aware of dependency rules.
func (db *DynamicBuilder[S, D, R]) generateDependentStateNavigation(b *mi.Builder, states []ComponentState, rules []DependencyRule, theme DynamicTheme) mi.Node {
	var buttons []interface{}

	for _, state := range states {
		btnClass := combineClasses(theme.StateTriggerClass(), "dyn-dependent-trigger")
		if state.Active {
			btnClass = combineClasses(btnClass, theme.StateTriggerActiveClass())
		}
		if state.Disabled {
			btnClass = combineClasses(btnClass, theme.StateTriggerDisabledClass())
		}

		btnAttrs := []interface{}{
			mi.Class(btnClass),
			mi.Data("state-target", state.ID),
			mi.Data("client-action", "switch-state"),
			mi.Attr("role", "tab"),
			mi.Attr("aria-selected", boolStr(state.Active)),
			mi.Attr("aria-controls", "state-"+state.ID),
		}

		// Check if this state is controlled by rules
		if isStateControlledByRules(state.ID, rules) {
			btnAttrs = append(btnAttrs, mi.Data("dependent", "true"))
		}

		if state.Disabled {
			btnAttrs = append(btnAttrs, mi.Disabled())
		}

		btnAttrs = append(btnAttrs, state.Label)
		buttons = append(buttons, b.Button(btnAttrs...))
	}

	navAttrs := []interface{}{
		mi.Class(combineClasses(theme.StateNavigationClass(), "dyn-dependent-navigation")),
		mi.Attr("role", "tablist"),
	}
	navAttrs = append(navAttrs, buttons...)

	return b.Div(navAttrs...)
}

// generateDependentStateContents creates panels with rule metadata.
func (db *DynamicBuilder[S, D, R]) generateDependentStateContents(b *mi.Builder, states []ComponentState, rules []DependencyRule, theme DynamicTheme) mi.Node {
	var panels []interface{}

	for _, state := range states {
		panelClass := combineClasses(theme.StateContentClass(), "dyn-dependent-state")
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

		// Add rules that affect this state as metadata
		affectingRules := findRulesAffectingTarget(rules, state.ID)
		if len(affectingRules) > 0 {
			panelAttrs = append(panelAttrs, mi.Data("dependent-rules", JSONOrEmpty(affectingRules)))
		}

		// Render content
		contentNode := db.renderStateContent(b, state.Content)
		panelAttrs = append(panelAttrs, contentNode)

		panels = append(panels, b.Div(panelAttrs...))
	}

	containerAttrs := []interface{}{mi.Class("dyn-dependent-states-container")}
	containerAttrs = append(containerAttrs, panels...)

	return b.Div(containerAttrs...)
}

// =============================================================================
// DEPENDENT DATA STRUCTURE (Data + Rules)
// =============================================================================

// generateDependentDataStructure creates filtering controlled by rules.
func (db *DynamicBuilder[S, D, R]) generateDependentDataStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	theme := db.getTheme()
	var children []mi.Node

	// Dependent filter controls
	children = append(children, db.generateDependentFilterControls(b, theme))

	// Results summary
	children = append(children, b.Div(
		mi.ID(db.id+"-summary"),
		mi.Class(theme.ResultsSummaryClass()),
	))

	// Results with dependency support
	children = append(children, b.Div(
		mi.ID(db.id+"-results"),
		mi.Class(combineClasses(theme.ResultsClass(), "dyn-dependent-results")),
	))

	return children
}

// generateDependentFilterControls creates filters that may be controlled by rules.
func (db *DynamicBuilder[S, D, R]) generateDependentFilterControls(b *mi.Builder, theme DynamicTheme) mi.Node {
	schema := db.extractFilterSchema()
	rules := db.extractRules()

	if len(schema.Fields) == 0 {
		schema = db.generateSchemaFromData()
	}

	var controls []interface{}
	for _, field := range schema.Fields {
		// Check if this filter is controlled by rules
		affectingRules := findRulesAffectingTarget(rules, field.Name)

		groupClass := theme.FilterGroupClass()
		if len(affectingRules) > 0 {
			groupClass = combineClasses(groupClass, "dyn-dependent-filter")
		}

		control := db.generateFilterControl(b, field, theme)

		groupAttrs := []interface{}{mi.Class(groupClass)}
		if len(affectingRules) > 0 {
			groupAttrs = append(groupAttrs, mi.Data("dependent-rules", JSONOrEmpty(affectingRules)))
		}
		groupAttrs = append(groupAttrs,
			b.Label(mi.Class(theme.FilterLabelClass()), mi.For(db.id+"-filter-"+field.Name), field.Label),
			control,
		)

		controls = append(controls, b.Div(groupAttrs...))
	}

	containerAttrs := []interface{}{
		mi.ID(db.id + "-filters"),
		mi.Class(combineClasses(theme.FilterControlsClass(), "dyn-dependent-filters")),
	}
	containerAttrs = append(containerAttrs, controls...)

	return b.Div(containerAttrs...)
}

// =============================================================================
// COMPLETE STRUCTURE (States + Data + Rules)
// =============================================================================

// generateCompleteStructure creates the full-featured component.
func (db *DynamicBuilder[S, D, R]) generateCompleteStructure(b *mi.Builder, pattern DetectedPattern) []mi.Node {
	states := db.extractStates()
	rules := db.extractRules()
	theme := db.getTheme()
	var children []mi.Node

	// State navigation with dependency awareness
	if len(states) > 1 {
		children = append(children, db.generateDependentStateNavigation(b, states, rules, theme))
	}

	// Each state has its own filters and results
	var stateContents []interface{}
	for _, state := range states {
		panelClass := combineClasses(theme.StateContentClass(), "dyn-complete-state")
		if state.Active {
			panelClass = combineClasses(panelClass, theme.StateContentActiveClass())
		} else {
			panelClass = combineClasses(panelClass, theme.StateContentHiddenClass())
		}

		panel := b.Div(
			mi.ID("state-"+state.ID),
			mi.Class(panelClass),
			mi.Data("state-id", state.ID),
			// Dependent filters within state context
			b.Div(
				mi.Class("dyn-state-filters dyn-dependent-filters"),
				mi.Data("state-context", state.ID),
			),
			// Results with full feature support
			b.Div(
				mi.Class("dyn-state-results dyn-complete-results"),
				mi.Data("state-context", state.ID),
			),
		)

		stateContents = append(stateContents, panel)
	}

	containerAttrs := []interface{}{mi.Class("dyn-complete-container")}
	containerAttrs = append(containerAttrs, stateContents...)
	children = append(children, b.Div(containerAttrs...))

	return children
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// isStateControlledByRules checks if a state is a target of any rule.
func isStateControlledByRules(stateID string, rules []DependencyRule) bool {
	for _, rule := range rules {
		for _, action := range rule.Actions {
			if action.TargetID == stateID || action.TargetID == "state-"+stateID {
				return true
			}
		}
	}
	return false
}

// findRulesAffectingTarget finds all rules that affect a specific target.
func findRulesAffectingTarget(rules []DependencyRule, targetID string) []DependencyRule {
	var affecting []DependencyRule

	for _, rule := range rules {
		for _, action := range rule.Actions {
			if action.TargetID == targetID || action.TargetID == "state-"+targetID {
				affecting = append(affecting, rule)
				break
			}
		}
	}

	return affecting
}
