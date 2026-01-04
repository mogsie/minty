package minty

import (
	"net/http"
	"strings"
)

// HTMX Request Detection

// IsHTMX checks if the request was made by HTMX.
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// IsHTMXBoosted checks if the request was boosted by HTMX.
func IsHTMXBoosted(r *http.Request) bool {
	return r.Header.Get("HX-Boosted") == "true"
}

// GetHTMXTarget returns the target element ID from HTMX request headers.
func GetHTMXTarget(r *http.Request) string {
	return r.Header.Get("HX-Target")
}

// GetHTMXTrigger returns the ID of the element that triggered the request.
func GetHTMXTrigger(r *http.Request) string {
	return r.Header.Get("HX-Trigger")
}

// GetHTMXTriggerName returns the name of the event that triggered the request.
func GetHTMXTriggerName(r *http.Request) string {
	return r.Header.Get("HX-Trigger-Name")
}

// GetHTMXCurrentURL returns the current URL of the browser.
func GetHTMXCurrentURL(r *http.Request) string {
	return r.Header.Get("HX-Current-URL")
}

// HTMX Response Helpers

// SetHTMXTrigger sets the HX-Trigger response header to trigger client-side events.
func SetHTMXTrigger(w http.ResponseWriter, events string) {
	w.Header().Set("HX-Trigger", events)
}

// SetHTMXTriggerAfterSwap sets the HX-Trigger-After-Swap response header.
func SetHTMXTriggerAfterSwap(w http.ResponseWriter, events string) {
	w.Header().Set("HX-Trigger-After-Swap", events)
}

// SetHTMXTriggerAfterSettle sets the HX-Trigger-After-Settle response header.
func SetHTMXTriggerAfterSettle(w http.ResponseWriter, events string) {
	w.Header().Set("HX-Trigger-After-Settle", events)
}

// SetHTMXLocation sets the HX-Location response header for client-side redirects.
func SetHTMXLocation(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Location", url)
}

// SetHTMXPushURL sets the HX-Push-Url response header to update browser history.
func SetHTMXPushURL(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Push-Url", url)
}

// SetHTMXRedirect sets the HX-Redirect response header for full redirects.
func SetHTMXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

// SetHTMXRefresh sets the HX-Refresh response header to refresh the page.
func SetHTMXRefresh(w http.ResponseWriter) {
	w.Header().Set("HX-Refresh", "true")
}

// SetHTMXReplaceURL sets the HX-Replace-Url response header to replace current URL.
func SetHTMXReplaceURL(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Replace-Url", url)
}

// SetHTMXReswap sets the HX-Reswap response header to change swap behavior.
func SetHTMXReswap(w http.ResponseWriter, strategy string) {
	w.Header().Set("HX-Reswap", strategy)
}

// SetHTMXRetarget sets the HX-Retarget response header to change the target.
func SetHTMXRetarget(w http.ResponseWriter, selector string) {
	w.Header().Set("HX-Retarget", selector)
}

// Fragment Rendering

// RenderFragment renders an HTML fragment for HTMX responses.
func RenderFragment(template H, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return Render(template, w)
}

// RenderFragmentWithTrigger renders a fragment and sets trigger events.
func RenderFragmentWithTrigger(template H, w http.ResponseWriter, triggerEvents string) error {
	SetHTMXTrigger(w, triggerEvents)
	return RenderFragment(template, w)
}

// HTMXHandler creates an HTTP handler that automatically handles HTMX vs full page requests.
func HTMXHandler(fullPageTemplate H, fragmentTemplate H) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		if IsHTMX(r) {
			if err := Render(fragmentTemplate, w); err != nil {
				http.Error(w, "Fragment render error", http.StatusInternalServerError)
			}
		} else {
			if err := Render(fullPageTemplate, w); err != nil {
				http.Error(w, "Template render error", http.StatusInternalServerError)
			}
		}
	}
}

// HTMXHandlerFunc creates an HTTP handler from functions that return templates.
func HTMXHandlerFunc(fullPageFn func(*http.Request) H, fragmentFn func(*http.Request) H) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		if IsHTMX(r) {
			template := fragmentFn(r)
			if err := Render(template, w); err != nil {
				http.Error(w, "Fragment render error", http.StatusInternalServerError)
			}
		} else {
			template := fullPageFn(r)
			if err := Render(template, w); err != nil {
				http.Error(w, "Template render error", http.StatusInternalServerError)
			}
		}
	}
}

// HTMX Response Status Helpers

// HTMXStopPolling stops HTMX polling by returning 286 status.
func HTMXStopPolling(w http.ResponseWriter) {
	w.WriteHeader(286) // HTMX stop polling status
}

// Common HTMX Patterns

// HTMXSwapStrategies provides constants for common swap strategies.
var HTMXSwapStrategies = struct {
	InnerHTML   string
	OuterHTML   string
	BeforeBegin string
	AfterBegin  string
	BeforeEnd   string
	AfterEnd    string
	Delete      string
	None        string
}{
	InnerHTML:   "innerHTML",
	OuterHTML:   "outerHTML", 
	BeforeBegin: "beforebegin",
	AfterBegin:  "afterbegin",
	BeforeEnd:   "beforeend",
	AfterEnd:    "afterend",
	Delete:      "delete",
	None:        "none",
}

// HTMXTriggers provides constants for common trigger patterns.
var HTMXTriggers = struct {
	Click          string
	Change         string
	KeyUp          string
	KeyDown        string
	Input          string
	Load           string
	Revealed       string
	Intersect      string
	SubmitForm     string
	ChangeDelayed  string
	KeyUpDelayed   string
	InputDelayed   string
}{
	Click:          "click",
	Change:         "change", 
	KeyUp:          "keyup",
	KeyDown:        "keydown",
	Input:          "input",
	Load:           "load",
	Revealed:       "revealed",
	Intersect:      "intersect",
	SubmitForm:     "submit",
	ChangeDelayed:  "change delay:500ms",
	KeyUpDelayed:   "keyup changed delay:300ms", 
	InputDelayed:   "input delay:300ms",
}

// HTMX Helper Components

// HTMXLoadingIndicator creates a loading indicator for HTMX requests.
func HTMXLoadingIndicator(text string) H {
	return func(b *Builder) Node {
		return b.Div(
			Class("htmx-indicator"),
			Style("display: none;"), // Hidden by default, shown by HTMX
			text,
		)
	}
}

// HTMXLoadingSpinner creates a CSS-based loading spinner.
func HTMXLoadingSpinner() H {
	return func(b *Builder) Node {
		return b.Div(
			Class("htmx-indicator spinner"),
			Style(`
				display: none;
				width: 20px;
				height: 20px;
				border: 2px solid #f3f3f3;
				border-top: 2px solid #3498db;
				border-radius: 50%;
				animation: spin 1s linear infinite;
				margin: 0 auto;
			`),
		)
	}
}

// LiveSearch creates a live search input with HTMX.
func LiveSearch(searchURL, targetSelector, placeholder string) H {
	return func(b *Builder) Node {
		return b.Input(
			Type("text"),
			Name("q"),
			Placeholder(placeholder),
			HtmxGet(searchURL),
			HtmxTarget(targetSelector),
			HtmxTrigger(HTMXTriggers.KeyUpDelayed),
			HtmxIndicator("#search-spinner"),
		)
	}
}

// InfiniteScroll creates an element that triggers infinite scroll loading.
func InfiniteScroll(loadURL, targetSelector string) H {
	return func(b *Builder) Node {
		return b.Div(
			HtmxGet(loadURL),
			HtmxTrigger("revealed"),
			HtmxTarget(targetSelector),
			HtmxSwap("beforeend"),
		)
	}
}

// HTMXForm creates a form that submits via HTMX.
func HTMXForm(method, action, targetSelector string, children ...interface{}) H {
	return func(b *Builder) Node {
		args := []interface{}{Method(method), Action(action)}
		
		// Add HTMX attributes based on method
		switch strings.ToUpper(method) {
		case "GET":
			args = append(args, HtmxGet(action))
		case "POST":
			args = append(args, HtmxPost(action))
		case "PUT":
			args = append(args, HtmxPut(action))
		case "DELETE":
			args = append(args, HtmxDelete(action))
		case "PATCH":
			args = append(args, HtmxPatch(action))
		}
		
		if targetSelector != "" {
			args = append(args, HtmxTarget(targetSelector))
		}
		
		// Add children
		args = append(args, children...)
		
		return b.Form(args...)
	}
}

// AutoRefresh creates an element that automatically refreshes via HTMX.
func AutoRefresh(url string, interval string) H {
	return func(b *Builder) Node {
		return b.Div(
			HtmxGet(url),
			HtmxTrigger("load, every "+interval),
			HtmxTarget("this"),
			HtmxSwap("outerHTML"),
		)
	}
}

// ProgressBar creates a progress bar that updates via HTMX.
func ProgressBar(progressURL string, targetSelector string) H {
	return func(b *Builder) Node {
		return b.Progress(
			HtmxGet(progressURL),
			HtmxTrigger("load, every 1s"),
			HtmxTarget(targetSelector),
			Max(100),
			Value("0"),
		)
	}
}

// Notification creates a notification that can be dismissed.
func Notification(message, notificationType string) H {
	return func(b *Builder) Node {
		return b.Div(
			Class("notification "+notificationType),
			ID("notification"),
			message,
			b.Button(
				Class("close"),
				HtmxDelete("/dismiss-notification"),
				HtmxTarget("#notification"),
				HtmxSwap("delete"),
				"×",
			),
		)
	}
}

// TabContainer creates a tabbed interface with HTMX.
type Tab struct {
	ID    string
	Title string
	URL   string
}

func TabContainer(tabs []Tab, activeTab string) H {
	return func(b *Builder) Node {
		return b.Div(Class("tab-container"),
			// Tab navigation
			b.Nav(Class("tabs"),
				NewFragment(Each(tabs, func(tab Tab) H {
					return func(b *Builder) Node {
						activeClass := ""
						if tab.ID == activeTab {
							activeClass = " active"
						}
						
						return b.Button(
							Class("tab"+activeClass),
							HtmxGet(tab.URL),
							HtmxTarget("#tab-content"),
							tab.Title,
						)
					}
				})...),
			),
			
			// Tab content area
			b.Div(ID("tab-content")),
		)
	}
}

// SearchWithFilters creates a search interface with filter options.
func SearchWithFilters(searchURL string, filters []FilterOption) H {
	return func(b *Builder) Node {
		return b.Form(
			HtmxGet(searchURL),
			HtmxTarget("#search-results"),
			HtmxTrigger("submit, change"),
			
			// Search input
			b.Input(
				Type("text"),
				Name("q"),
				Placeholder("Search..."),
			),
			
			// Filters
			NewFragment(Each(filters, func(filter FilterOption) H {
				return func(b *Builder) Node {
					return b.Label(
						b.Input(
							Type("checkbox"),
							Name("filter"),
							Value(filter.Value),
						),
						filter.Label,
					)
				}
			})...),
			
			// Results area
			b.Div(ID("search-results")),
		)
	}
}

type FilterOption struct {
	Value string
	Label string
}
