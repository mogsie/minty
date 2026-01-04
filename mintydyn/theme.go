package mintydyn

// =============================================================================
// DYNAMIC THEME INTERFACE
// =============================================================================

// DynamicTheme provides theme-specific styling for dynamic components.
// Return empty string from any method to use the default class.
type DynamicTheme interface {
	// Container classes
	ComponentClass() string        // default: "dyn-component"
	ComponentPatternClass(pattern string) string // default: "dyn-{pattern}"

	// State/tab navigation
	StateNavigationClass() string      // default: "dyn-state-navigation"
	StateTriggerClass() string         // default: "dyn-state-trigger"
	StateTriggerActiveClass() string   // default: "active"
	StateTriggerDisabledClass() string // default: "disabled"
	StateContentClass() string         // default: "dyn-state-content"
	StateContentActiveClass() string   // default: "active"
	StateContentHiddenClass() string   // default: "hidden"
	StateContainerClass() string       // default: "dyn-state-container"

	// Filter controls
	FilterControlsClass() string   // default: "dyn-filter-controls"
	FilterGroupClass() string      // default: "dyn-filter-group"
	FilterLabelClass() string      // default: "dyn-filter-label"
	FilterInputClass() string      // default: "dyn-filter-input"
	FilterSelectClass() string     // default: "dyn-filter-select"
	FilterCheckboxClass() string   // default: "dyn-filter-checkbox"
	FilterRangeClass() string      // default: "dyn-filter-range"

	// Results
	ResultsClass() string          // default: "dyn-results"
	ResultsEmptyClass() string     // default: "dyn-no-results"
	ResultsSummaryClass() string   // default: "dyn-results-summary"
	PaginationClass() string       // default: "dyn-pagination"
	PaginationButtonClass() string       // default: "dyn-page-btn"
	PaginationButtonActiveClass() string // default: "active"

	// Utility
	HiddenClass() string           // default: "hidden"
	DisabledClass() string         // default: "disabled"

	// Optional: inject theme-specific CSS
	InjectCSS() string             // default: ""
}

// =============================================================================
// DEFAULT THEME
// =============================================================================

// DefaultTheme provides semantic class names that work standalone.
type DefaultTheme struct{}

// NewDefaultTheme creates a new default theme.
func NewDefaultTheme() DynamicTheme {
	return &DefaultTheme{}
}

func (t *DefaultTheme) ComponentClass() string              { return "dyn-component" }
func (t *DefaultTheme) ComponentPatternClass(p string) string { return "dyn-" + p }
func (t *DefaultTheme) StateNavigationClass() string        { return "dyn-state-navigation" }
func (t *DefaultTheme) StateTriggerClass() string           { return "dyn-state-trigger" }
func (t *DefaultTheme) StateTriggerActiveClass() string     { return "active" }
func (t *DefaultTheme) StateTriggerDisabledClass() string   { return "disabled" }
func (t *DefaultTheme) StateContentClass() string           { return "dyn-state-content" }
func (t *DefaultTheme) StateContentActiveClass() string     { return "active" }
func (t *DefaultTheme) StateContentHiddenClass() string     { return "hidden" }
func (t *DefaultTheme) StateContainerClass() string         { return "dyn-state-container" }
func (t *DefaultTheme) FilterControlsClass() string         { return "dyn-filter-controls" }
func (t *DefaultTheme) FilterGroupClass() string            { return "dyn-filter-group" }
func (t *DefaultTheme) FilterLabelClass() string            { return "dyn-filter-label" }
func (t *DefaultTheme) FilterInputClass() string            { return "dyn-filter-input" }
func (t *DefaultTheme) FilterSelectClass() string           { return "dyn-filter-select" }
func (t *DefaultTheme) FilterCheckboxClass() string         { return "dyn-filter-checkbox" }
func (t *DefaultTheme) FilterRangeClass() string            { return "dyn-filter-range" }
func (t *DefaultTheme) ResultsClass() string                { return "dyn-results" }
func (t *DefaultTheme) ResultsEmptyClass() string           { return "dyn-no-results" }
func (t *DefaultTheme) ResultsSummaryClass() string         { return "dyn-results-summary" }
func (t *DefaultTheme) PaginationClass() string             { return "dyn-pagination" }
func (t *DefaultTheme) PaginationButtonClass() string       { return "dyn-page-btn" }
func (t *DefaultTheme) PaginationButtonActiveClass() string { return "active" }
func (t *DefaultTheme) HiddenClass() string                 { return "hidden" }
func (t *DefaultTheme) DisabledClass() string               { return "disabled" }
func (t *DefaultTheme) InjectCSS() string                   { return "" }

// =============================================================================
// BOOTSTRAP DYNAMIC THEME
// =============================================================================

// BootstrapDynamicTheme provides Bootstrap 5 compatible classes.
type BootstrapDynamicTheme struct{}

// NewBootstrapDynamicTheme creates a Bootstrap-compatible theme.
func NewBootstrapDynamicTheme() DynamicTheme {
	return &BootstrapDynamicTheme{}
}

func (t *BootstrapDynamicTheme) ComponentClass() string              { return "dyn-component" }
func (t *BootstrapDynamicTheme) ComponentPatternClass(p string) string { return "dyn-" + p }
func (t *BootstrapDynamicTheme) StateNavigationClass() string        { return "nav nav-tabs" }
func (t *BootstrapDynamicTheme) StateTriggerClass() string           { return "nav-link" }
func (t *BootstrapDynamicTheme) StateTriggerActiveClass() string     { return "active" }
func (t *BootstrapDynamicTheme) StateTriggerDisabledClass() string   { return "disabled" }
func (t *BootstrapDynamicTheme) StateContentClass() string           { return "tab-pane fade" }
func (t *BootstrapDynamicTheme) StateContentActiveClass() string     { return "show active" }
func (t *BootstrapDynamicTheme) StateContentHiddenClass() string     { return "" } // Bootstrap uses fade
func (t *BootstrapDynamicTheme) StateContainerClass() string         { return "tab-content" }
func (t *BootstrapDynamicTheme) FilterControlsClass() string         { return "row g-3 mb-3" }
func (t *BootstrapDynamicTheme) FilterGroupClass() string            { return "col-md-4" }
func (t *BootstrapDynamicTheme) FilterLabelClass() string            { return "form-label" }
func (t *BootstrapDynamicTheme) FilterInputClass() string            { return "form-control" }
func (t *BootstrapDynamicTheme) FilterSelectClass() string           { return "form-select" }
func (t *BootstrapDynamicTheme) FilterCheckboxClass() string         { return "form-check-input" }
func (t *BootstrapDynamicTheme) FilterRangeClass() string            { return "form-range" }
func (t *BootstrapDynamicTheme) ResultsClass() string                { return "dyn-results" }
func (t *BootstrapDynamicTheme) ResultsEmptyClass() string           { return "alert alert-info" }
func (t *BootstrapDynamicTheme) ResultsSummaryClass() string         { return "text-muted mb-2" }
func (t *BootstrapDynamicTheme) PaginationClass() string             { return "pagination" }
func (t *BootstrapDynamicTheme) PaginationButtonClass() string       { return "page-link" }
func (t *BootstrapDynamicTheme) PaginationButtonActiveClass() string { return "active" }
func (t *BootstrapDynamicTheme) HiddenClass() string                 { return "d-none" }
func (t *BootstrapDynamicTheme) DisabledClass() string               { return "disabled" }
func (t *BootstrapDynamicTheme) InjectCSS() string                   { return "" }

// =============================================================================
// TAILWIND DYNAMIC THEME
// =============================================================================

// TailwindDynamicTheme provides Tailwind CSS compatible classes.
type TailwindDynamicTheme struct{}

// NewTailwindDynamicTheme creates a Tailwind-compatible theme.
func NewTailwindDynamicTheme() DynamicTheme {
	return &TailwindDynamicTheme{}
}

func (t *TailwindDynamicTheme) ComponentClass() string              { return "dyn-component" }
func (t *TailwindDynamicTheme) ComponentPatternClass(p string) string { return "dyn-" + p }
func (t *TailwindDynamicTheme) StateNavigationClass() string        { return "flex border-b border-gray-200 mb-4" }
func (t *TailwindDynamicTheme) StateTriggerClass() string           { return "px-4 py-2.5 text-sm font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 border-b-2 border-transparent -mb-px cursor-pointer transition-colors" }
func (t *TailwindDynamicTheme) StateTriggerActiveClass() string     { return "text-blue-600 !border-b-4 !border-blue-600 !bg-blue-50 font-semibold" }
func (t *TailwindDynamicTheme) StateTriggerDisabledClass() string   { return "opacity-50 cursor-not-allowed" }
func (t *TailwindDynamicTheme) StateContentClass() string           { return "dyn-state-content" }
func (t *TailwindDynamicTheme) StateContentActiveClass() string     { return "block" }
func (t *TailwindDynamicTheme) StateContentHiddenClass() string     { return "hidden" }
func (t *TailwindDynamicTheme) StateContainerClass() string         { return "mt-4" }
func (t *TailwindDynamicTheme) FilterControlsClass() string         { return "grid grid-cols-1 md:grid-cols-3 gap-4 mb-4" }
func (t *TailwindDynamicTheme) FilterGroupClass() string            { return "flex flex-col" }
func (t *TailwindDynamicTheme) FilterLabelClass() string            { return "text-sm font-medium text-gray-700 mb-1" }
func (t *TailwindDynamicTheme) FilterInputClass() string            { return "mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" }
func (t *TailwindDynamicTheme) FilterSelectClass() string           { return "mt-1 block w-full px-3 py-2 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" }
func (t *TailwindDynamicTheme) FilterCheckboxClass() string         { return "h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded" }
func (t *TailwindDynamicTheme) FilterRangeClass() string            { return "w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer" }
func (t *TailwindDynamicTheme) ResultsClass() string                { return "dyn-results" }
func (t *TailwindDynamicTheme) ResultsEmptyClass() string           { return "text-center py-8 text-gray-500" }
func (t *TailwindDynamicTheme) ResultsSummaryClass() string         { return "text-sm text-gray-500 mb-2" }
func (t *TailwindDynamicTheme) PaginationClass() string             { return "flex justify-center space-x-2 mt-4" }
func (t *TailwindDynamicTheme) PaginationButtonClass() string       { return "px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100" }
func (t *TailwindDynamicTheme) PaginationButtonActiveClass() string { return "bg-blue-600 text-white border-blue-600 hover:bg-blue-700" }
func (t *TailwindDynamicTheme) HiddenClass() string                 { return "hidden" }
func (t *TailwindDynamicTheme) DisabledClass() string               { return "opacity-50 cursor-not-allowed" }
func (t *TailwindDynamicTheme) InjectCSS() string                   { return "" }

// =============================================================================
// TAILWIND DARK THEME (with dark: variants)
// =============================================================================

// TailwindDarkTheme provides Tailwind CSS classes with dark mode support.
// Uses Tailwind's dark: prefix so both light and dark modes work automatically.
// Requires either:
// - class="dark" on <html> element, OR
// - User's system preference set to dark mode (with Tailwind CDN)
type TailwindDarkTheme struct{}

// NewTailwindDarkTheme creates a dark-mode-aware Tailwind theme.
func NewTailwindDarkTheme() DynamicTheme {
	return &TailwindDarkTheme{}
}

func (t *TailwindDarkTheme) ComponentClass() string              { return "dyn-component" }
func (t *TailwindDarkTheme) ComponentPatternClass(p string) string { return "dyn-" + p }
func (t *TailwindDarkTheme) StateNavigationClass() string        { return "flex border-b border-gray-200 dark:border-gray-700 mb-4" }
func (t *TailwindDarkTheme) StateTriggerClass() string           { return "px-4 py-2.5 text-sm font-medium text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700/50 border-b-2 border-transparent -mb-px cursor-pointer transition-colors" }
func (t *TailwindDarkTheme) StateTriggerActiveClass() string     { return "!text-blue-600 dark:!text-blue-400 !border-blue-600 dark:!border-blue-400 !bg-blue-50 dark:!bg-blue-900/30 font-semibold" }
func (t *TailwindDarkTheme) StateTriggerDisabledClass() string   { return "opacity-50 cursor-not-allowed" }
func (t *TailwindDarkTheme) StateContentClass() string           { return "dyn-state-content" }
func (t *TailwindDarkTheme) StateContentActiveClass() string     { return "block" }
func (t *TailwindDarkTheme) StateContentHiddenClass() string     { return "hidden" }
func (t *TailwindDarkTheme) StateContainerClass() string         { return "mt-4" }
func (t *TailwindDarkTheme) FilterControlsClass() string         { return "grid grid-cols-1 md:grid-cols-3 gap-4 mb-4" }
func (t *TailwindDarkTheme) FilterGroupClass() string            { return "flex flex-col" }
func (t *TailwindDarkTheme) FilterLabelClass() string            { return "text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" }
func (t *TailwindDarkTheme) FilterInputClass() string            { return "mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" }
func (t *TailwindDarkTheme) FilterSelectClass() string           { return "mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" }
func (t *TailwindDarkTheme) FilterCheckboxClass() string         { return "h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded" }
func (t *TailwindDarkTheme) FilterRangeClass() string            { return "w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer" }
func (t *TailwindDarkTheme) ResultsClass() string                { return "dyn-results" }
func (t *TailwindDarkTheme) ResultsEmptyClass() string           { return "text-center py-8 text-gray-500 dark:text-gray-400" }
func (t *TailwindDarkTheme) ResultsSummaryClass() string         { return "text-sm text-gray-500 dark:text-gray-400 mb-2" }
func (t *TailwindDarkTheme) PaginationClass() string             { return "flex justify-center space-x-2 mt-4" }
func (t *TailwindDarkTheme) PaginationButtonClass() string       { return "px-3 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800" }
func (t *TailwindDarkTheme) PaginationButtonActiveClass() string { return "!bg-blue-600 !text-white !border-blue-600 hover:!bg-blue-700" }
func (t *TailwindDarkTheme) HiddenClass() string                 { return "hidden" }
func (t *TailwindDarkTheme) DisabledClass() string               { return "opacity-50 cursor-not-allowed" }
func (t *TailwindDarkTheme) InjectCSS() string                   { return "" }

// =============================================================================
// THEME HELPER FUNCTIONS
// =============================================================================

// getTheme returns the theme or default if nil.
func (db *DynamicBuilder[S, D, R]) getTheme() DynamicTheme {
	if db.theme != nil {
		return db.theme
	}
	return NewDefaultTheme()
}

// getClass returns the theme class or fallback if empty.
func getClass(themeClass, fallback string) string {
	if themeClass != "" {
		return themeClass
	}
	return fallback
}

// combineClasses joins class names, filtering empty strings.
func combineClasses(classes ...string) string {
	var result string
	for _, c := range classes {
		if c != "" {
			if result != "" {
				result += " "
			}
			result += c
		}
	}
	return result
}
