package ui

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	mi "github.com/ha1tch/minty"
	mdy "github.com/ha1tch/minty/mintydyn"
	"github.com/ha1tch/minty/examples/insurance-quote/internal/models"
	"github.com/ha1tch/minty/examples/insurance-quote/internal/store"
)

// Handler handles HTTP requests.
type Handler struct {
	store  *store.Store
	logger *slog.Logger
	theme  mdy.DynamicTheme
}

// NewHandler creates a new handler.
func NewHandler(store *store.Store, logger *slog.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
		theme:  mdy.NewTailwindDynamicTheme(),
	}
}

// =============================================================================
// DARK MODE
// =============================================================================

var darkMode = mi.DarkModeTailwind(
	mi.DarkModeMinify(),
)

// =============================================================================
// LAYOUT
// =============================================================================

const globalCSS = `
::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: #f1f1f1; }
.dark ::-webkit-scrollbar-track { background: #1f2937; }
::-webkit-scrollbar-thumb { background: #c1c1c1; border-radius: 3px; }
.dark ::-webkit-scrollbar-thumb { background: #4b5563; }
`

func (h *Handler) pageLayout(activePage, title, subtitle string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			mi.Raw("<!DOCTYPE html>"),
			b.Html(mi.Lang("en"),
				b.Head(
					b.Title("InsureQuote - "+title),
					b.Meta(mi.Charset("UTF-8")),
					b.Meta(mi.Name("viewport"), mi.Content("width=device-width, initial-scale=1")),
					b.Script(mi.Src("https://cdn.tailwindcss.com")),
					b.Script(mi.Raw(`tailwind.config = { darkMode: 'class' }`)),
					b.Style(mi.Raw(globalCSS)),
					darkMode.Script(b),
				),
				b.Body(mi.Class("bg-gray-50 dark:bg-gray-900 min-h-screen transition-colors"),
					b.Div(mi.Class("flex"),
						h.sidebar(b, activePage),
						b.Div(mi.Class("flex-1 ml-64"),
							h.header(b, title, subtitle),
							b.Main(mi.Class("p-6"), content(b)),
						),
					),
				),
			),
		)
	}
}

func (h *Handler) sidebar(b *mi.Builder, activePage string) mi.Node {
	navItems := []struct{ IconName, Label, Href, ID string }{
		{"home", "Dashboard", "/", "dashboard"},
		{"shield-check", "Get Quote", "/quote", "quote"},
		{"clipboard-document-list", "My Quotes", "/quotes", "quotes"},
		{"document-text", "Claims", "/claims", "claims"},
		{"calculator", "Compare Plans", "/compare", "compare"},
		{"cog-6-tooth", "Settings", "/settings", "settings"},
	}

	var items []interface{}
	for _, item := range navItems {
		itemClass := "flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-colors "
		if item.ID == activePage {
			itemClass += "bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300"
		} else {
			itemClass += "text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
		}
		items = append(items,
			b.A(mi.Href(item.Href), mi.Class(itemClass),
				Icon(item.IconName, "w-5 h-5"),
				item.Label,
			),
		)
	}

	navArgs := []interface{}{mi.Class("p-4 space-y-1")}
	navArgs = append(navArgs, items...)

	return b.Aside(mi.Class("fixed left-0 top-0 w-64 h-screen bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700"),
		b.Div(mi.Class("p-4 border-b border-gray-200 dark:border-gray-700"),
			b.Div(mi.Class("flex items-center gap-3"),
				Icon("shield-check", "w-8 h-8 text-blue-600 dark:text-blue-400"),
				b.Span(mi.Class("text-xl font-bold text-gray-900 dark:text-white"), "InsureQuote"),
			),
		),
		b.Nav(navArgs...),
	)
}

func (h *Handler) header(b *mi.Builder, title, subtitle string) mi.Node {
	return b.Header(mi.Class("bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4"),
		b.Div(mi.Class("flex items-center justify-between"),
			b.Div(
				b.H1(mi.Class("text-2xl font-bold text-gray-900 dark:text-white"), title),
				b.P(mi.Class("text-sm text-gray-500 dark:text-gray-400"), subtitle),
			),
			b.Div(mi.Class("flex items-center gap-3"),
				darkMode.Toggle(b,
					mi.Class("p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"),
				),
				b.Button(mi.Class("p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 relative"),
					Icon("bell", "w-5 h-5"),
					b.Span(mi.Class("absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full")),
				),
			),
		),
	)
}

// =============================================================================
// DASHBOARD
// =============================================================================

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	page := h.pageLayout("dashboard", "Dashboard", "Overview of your insurance portfolio", func(b *mi.Builder) mi.Node {
		return b.Div(
			// Stats cards
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6"),
				h.statCard(b, "Active Policies", "4", "shield-check", "text-green-600 dark:text-green-400", "bg-green-50 dark:bg-green-900/20"),
				h.statCard(b, "Pending Quotes", "2", "clock", "text-yellow-600 dark:text-yellow-400", "bg-yellow-50 dark:bg-yellow-900/20"),
				h.statCard(b, "Open Claims", "3", "exclamation-circle", "text-red-600 dark:text-red-400", "bg-red-50 dark:bg-red-900/20"),
				h.statCard(b, "Monthly Premium", "$485", "currency-dollar", "text-blue-600 dark:text-blue-400", "bg-blue-50 dark:bg-blue-900/20"),
			),
			// Quick actions
			b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 mb-6"),
				b.H2(mi.Class("text-lg font-semibold text-gray-900 dark:text-white mb-4"), "Quick Actions"),
				b.Div(mi.Class("grid grid-cols-2 md:grid-cols-4 gap-4"),
					h.actionCard(b, "Get New Quote", "Start a new insurance quote", "plus", "/quote"),
					h.actionCard(b, "File a Claim", "Report an incident", "document-text", "/claims/new"),
					h.actionCard(b, "Compare Plans", "Find the best coverage", "calculator", "/compare"),
					h.actionCard(b, "Contact Support", "Get help from our team", "phone", "/support"),
				),
			),
			// Coverage types
			b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6"),
				b.H2(mi.Class("text-lg font-semibold text-gray-900 dark:text-white mb-4"), "Available Coverage"),
				func() mi.Node {
					args := []interface{}{mi.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4")}
					args = append(args, h.coverageCards(b)...)
					return b.Div(args...)
				}(),
			),
		)
	})
	h.render(w, page)
}

func (h *Handler) statCard(b *mi.Builder, label, value, iconName, iconColor, bgColor string) mi.Node {
	return b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4"),
		b.Div(mi.Class("flex items-center justify-between"),
			b.Div(
				b.P(mi.Class("text-sm text-gray-500 dark:text-gray-400"), label),
				b.P(mi.Class("text-2xl font-bold text-gray-900 dark:text-white mt-1"), value),
			),
			b.Div(mi.Class("p-3 rounded-lg "+bgColor),
				Icon(iconName, "w-6 h-6 "+iconColor),
			),
		),
	)
}

func (h *Handler) actionCard(b *mi.Builder, title, desc, iconName, href string) mi.Node {
	return b.A(mi.Href(href), mi.Class("block p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors group"),
		b.Div(mi.Class("flex items-center gap-3"),
			b.Div(mi.Class("p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg group-hover:bg-blue-200 dark:group-hover:bg-blue-800/50 transition-colors"),
				Icon(iconName, "w-5 h-5 text-blue-600 dark:text-blue-400"),
			),
			b.Div(
				b.P(mi.Class("font-medium text-gray-900 dark:text-white"), title),
				b.P(mi.Class("text-xs text-gray-500 dark:text-gray-400"), desc),
			),
		),
	)
}

func (h *Handler) coverageCards(b *mi.Builder) []interface{} {
	var cards []interface{}
	for _, cov := range h.store.Coverages {
		cards = append(cards, b.A(mi.Href("/quote?type="+cov.ID),
			mi.Class("block p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 hover:shadow-md transition-all"),
			b.Div(mi.Class("flex items-start gap-3"),
				b.Div(mi.Class("p-2 bg-gray-100 dark:bg-gray-700 rounded-lg"),
					Icon(cov.Icon, "w-6 h-6 text-gray-600 dark:text-gray-300"),
				),
				b.Div(
					b.P(mi.Class("font-medium text-gray-900 dark:text-white"), cov.Name),
					b.P(mi.Class("text-xs text-gray-500 dark:text-gray-400 mt-1 line-clamp-2"), cov.Description),
					b.P(mi.Class("text-sm font-medium text-blue-600 dark:text-blue-400 mt-2"),
						fmt.Sprintf("From $%.2f/mo", cov.BasePrice),
					),
				),
			),
		))
	}
	return cards
}

// =============================================================================
// QUOTE WIZARD - Demonstrates STATES + RULES patterns
// =============================================================================

func (h *Handler) QuoteWizard(w http.ResponseWriter, r *http.Request) {
	coverageType := r.URL.Query().Get("type")
	if coverageType == "" {
		coverageType = "auto"
	}

	page := h.pageLayout("quote", "Get a Quote", "Complete the form to receive your personalized quote", func(b *mi.Builder) mi.Node {
		// PATTERN: States (wizard steps)
		wizardStates := []mdy.ComponentState{
			{ID: "coverage", Label: "Coverage Type", Active: true},
			{ID: "details", Label: "Your Details"},
			{ID: "customize", Label: "Customize Plan"},
			{ID: "review", Label: "Review & Submit"},
		}

		wizard := mdy.Dyn("quote-wizard").
			States(wizardStates).
			Theme(h.theme).
			Minified().
			Build()

		return b.Div(mi.Class("max-w-4xl mx-auto"),
			// Progress indicator
			b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 mb-6"),
				wizard(b),
			),
			// Form content - demonstrates RULES pattern
			b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6"),
				h.quoteFormWithRules(b, coverageType),
			),
		)
	})
	h.render(w, page)
}

// quoteFormWithRules demonstrates the RULES (dependency) pattern.
// Fields show/hide based on coverage type selection.
func (h *Handler) quoteFormWithRules(b *mi.Builder, initialType string) mi.Node {
	// PATTERN: Rules (form field dependencies)
	// When coverage type changes, show/hide relevant field sections
	formRules := mdy.Dyn("quote-form-rules").
		Rules([]mdy.DependencyRule{
			// Auto insurance fields
			mdy.ShowWhen("coverage-type", "equals", "auto", "auto-fields"),
			// Home insurance fields
			mdy.ShowWhen("coverage-type", "equals", "home", "home-fields"),
			// Life insurance fields
			mdy.ShowWhen("coverage-type", "equals", "life", "life-fields"),
			// Business insurance fields
			mdy.ShowWhen("coverage-type", "equals", "business", "business-fields"),
			// Accident details shown when "has accidents" is checked
			mdy.ShowWhen("has-accidents", "equals", true, "accident-details"),
			// Pool coverage shown when "has pool" is checked
			mdy.ShowWhen("has-pool", "equals", true, "pool-coverage"),
			// Smoker surcharge notice
			mdy.ShowWhen("is-smoker", "equals", true, "smoker-notice"),
			// Business premises fields
			mdy.ShowWhen("has-premises", "equals", true, "premises-fields"),
		}).
		Theme(h.theme).
		Minified().
		Build()

	// Coverage type icons for visual selection
	coverageOptions := []struct {
		Value, Label, Icon, Desc string
	}{
		{"auto", "Auto", "truck", "Vehicle coverage"},
		{"home", "Home", "home-modern", "Property protection"},
		{"life", "Life", "heart", "Family security"},
		{"business", "Business", "building-office", "Business protection"},
	}

	var coverageButtons []interface{}
	for _, opt := range coverageOptions {
		selected := opt.Value == initialType
		btnClass := "flex flex-col items-center p-4 rounded-lg border-2 transition-all cursor-pointer "
		if selected {
			btnClass += "border-blue-500 bg-blue-50 dark:bg-blue-900/30"
		} else {
			btnClass += "border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600"
		}

		inputAttrs := []mi.Attribute{
			mi.Type("radio"), mi.Name("coverage-type"), mi.ID("coverage-type-" + opt.Value),
			mi.Value(opt.Value),
			mi.Class("sr-only"),
			mi.Data("dependency-trigger", "coverage-type"),
		}
		if selected {
			inputAttrs = append(inputAttrs, mi.Attr("checked", "checked"))
		}

		coverageButtons = append(coverageButtons,
			b.Label(mi.Class(btnClass),
				b.Input(inputAttrs...),
				Icon(opt.Icon, "w-8 h-8 text-gray-600 dark:text-gray-300 mb-2"),
				b.Span(mi.Class("font-medium text-gray-900 dark:text-white"), opt.Label),
				b.Span(mi.Class("text-xs text-gray-500 dark:text-gray-400"), opt.Desc),
			),
		)
	}

	return b.Form(mi.Method("POST"), mi.Action("/quote/submit"),
		formRules(b),
		// Coverage type selection
		b.Div(mi.Class("mb-6"),
			b.Label(mi.Class("block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3"), "Select Coverage Type"),
			func() mi.Node {
				args := []interface{}{mi.Class("grid grid-cols-2 md:grid-cols-4 gap-4")}
				args = append(args, coverageButtons...)
				return b.Div(args...)
			}(),
		),

		// Basic contact info (always shown)
		b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-6"),
			h.formField(b, "First Name", "firstName", "text", "John", true),
			h.formField(b, "Last Name", "lastName", "text", "Smith", true),
			h.formField(b, "Email", "email", "email", "john@example.com", true),
			h.formField(b, "Phone", "phone", "tel", "(555) 123-4567", true),
		),

		// === AUTO INSURANCE FIELDS ===
		b.Div(mi.ID("auto-fields"), mi.Class("border-t border-gray-200 dark:border-gray-700 pt-6 mb-6"),
			mi.Data("dependency-target", "auto-fields"),
			b.H3(mi.Class("text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center gap-2"),
				Icon("truck", "w-5 h-5"), "Vehicle Information",
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-3 gap-4 mb-4"),
				h.formField(b, "Vehicle Make", "vehicleMake", "text", "Toyota", false),
				h.formField(b, "Vehicle Model", "vehicleModel", "text", "Camry", false),
				h.formField(b, "Vehicle Year", "vehicleYear", "number", "2022", false),
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "VIN", "vin", "text", "1HGBH41JXMN109186", false),
				h.formField(b, "Years Driving", "drivingYears", "number", "10", false),
			),
			// Conditional: accidents
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("flex items-center gap-2 cursor-pointer"),
					b.Input(mi.Type("checkbox"), mi.ID("has-accidents"), mi.Name("hasAccidents"),
						mi.Class("rounded border-gray-300 text-blue-600 focus:ring-blue-500"),
						mi.Data("dependency-trigger", "has-accidents"),
					),
					b.Span(mi.Class("text-sm text-gray-700 dark:text-gray-300"), "I have had accidents in the past 5 years"),
				),
			),
			b.Div(mi.ID("accident-details"), mi.Class("ml-6 p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg hidden"),
				mi.Data("dependency-target", "accident-details"),
				b.Div(mi.Class("flex items-start gap-2 mb-3"),
					Icon("exclamation-triangle", "w-5 h-5 text-yellow-600 dark:text-yellow-400 flex-shrink-0 mt-0.5"),
					b.P(mi.Class("text-sm text-yellow-800 dark:text-yellow-200"), "Accident history may affect your premium. Please provide details."),
				),
				h.formField(b, "Number of Accidents", "accidentCount", "number", "1", false),
			),
		),

		// === HOME INSURANCE FIELDS ===
		b.Div(mi.ID("home-fields"), mi.Class("border-t border-gray-200 dark:border-gray-700 pt-6 mb-6 hidden"),
			mi.Data("dependency-target", "home-fields"),
			b.H3(mi.Class("text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center gap-2"),
				Icon("home-modern", "w-5 h-5"), "Property Information",
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-3 gap-4 mb-4"),
				h.formSelect(b, "Property Type", "propertyType", []string{"House", "Condo", "Townhouse", "Apartment"}),
				h.formField(b, "Year Built", "yearBuilt", "number", "1995", false),
				h.formField(b, "Square Feet", "squareFeet", "number", "2000", false),
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "Property Value ($)", "propertyValue", "number", "350000", false),
				h.formField(b, "Zip Code", "zipCode", "text", "90210", false),
			),
			// Conditional: pool
			b.Div(mi.Class("space-y-3"),
				b.Label(mi.Class("flex items-center gap-2 cursor-pointer"),
					b.Input(mi.Type("checkbox"), mi.ID("has-pool"), mi.Name("hasPool"),
						mi.Class("rounded border-gray-300 text-blue-600 focus:ring-blue-500"),
						mi.Data("dependency-trigger", "has-pool"),
					),
					b.Span(mi.Class("text-sm text-gray-700 dark:text-gray-300"), "Property has a swimming pool"),
				),
				b.Label(mi.Class("flex items-center gap-2 cursor-pointer"),
					b.Input(mi.Type("checkbox"), mi.Name("hasAlarm"),
						mi.Class("rounded border-gray-300 text-blue-600 focus:ring-blue-500"),
					),
					b.Span(mi.Class("text-sm text-gray-700 dark:text-gray-300"), "Property has a security alarm (discount available)"),
				),
			),
			b.Div(mi.ID("pool-coverage"), mi.Class("mt-4 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg hidden"),
				mi.Data("dependency-target", "pool-coverage"),
				b.Div(mi.Class("flex items-start gap-2"),
					Icon("information-circle", "w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5"),
					b.P(mi.Class("text-sm text-blue-800 dark:text-blue-200"), "Pool coverage includes liability protection and equipment coverage. Additional premium of $15/month applies."),
				),
			),
		),

		// === LIFE INSURANCE FIELDS ===
		b.Div(mi.ID("life-fields"), mi.Class("border-t border-gray-200 dark:border-gray-700 pt-6 mb-6 hidden"),
			mi.Data("dependency-target", "life-fields"),
			b.H3(mi.Class("text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center gap-2"),
				Icon("heart", "w-5 h-5"), "Health Information",
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "Date of Birth", "dateOfBirth", "date", "", false),
				h.formSelect(b, "Health Status", "healthStatus", []string{"Excellent", "Good", "Fair", "Poor"}),
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "Coverage Amount ($)", "coverageAmount", "number", "250000", false),
				h.formField(b, "Number of Beneficiaries", "beneficiaries", "number", "2", false),
			),
			// Conditional: smoker
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("flex items-center gap-2 cursor-pointer"),
					b.Input(mi.Type("checkbox"), mi.ID("is-smoker"), mi.Name("isSmoker"),
						mi.Class("rounded border-gray-300 text-blue-600 focus:ring-blue-500"),
						mi.Data("dependency-trigger", "is-smoker"),
					),
					b.Span(mi.Class("text-sm text-gray-700 dark:text-gray-300"), "I am a smoker or have used tobacco in the past 12 months"),
				),
			),
			b.Div(mi.ID("smoker-notice"), mi.Class("p-4 bg-orange-50 dark:bg-orange-900/20 rounded-lg hidden"),
				mi.Data("dependency-target", "smoker-notice"),
				b.Div(mi.Class("flex items-start gap-2"),
					Icon("exclamation-triangle", "w-5 h-5 text-orange-600 dark:text-orange-400 flex-shrink-0 mt-0.5"),
					b.P(mi.Class("text-sm text-orange-800 dark:text-orange-200"), "Tobacco use may result in higher premiums. Consider our smoking cessation program for potential discounts."),
				),
			),
		),

		// === BUSINESS INSURANCE FIELDS ===
		b.Div(mi.ID("business-fields"), mi.Class("border-t border-gray-200 dark:border-gray-700 pt-6 mb-6 hidden"),
			mi.Data("dependency-target", "business-fields"),
			b.H3(mi.Class("text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center gap-2"),
				Icon("building-office", "w-5 h-5"), "Business Information",
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "Business Name", "businessName", "text", "Acme Corp", false),
				h.formSelect(b, "Business Type", "businessType", []string{"Retail", "Restaurant", "Office", "Manufacturing", "Service", "Other"}),
			),
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4 mb-4"),
				h.formField(b, "Number of Employees", "employees", "number", "25", false),
				h.formField(b, "Annual Revenue ($)", "annualRevenue", "number", "500000", false),
			),
			// Conditional: premises
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("flex items-center gap-2 cursor-pointer"),
					b.Input(mi.Type("checkbox"), mi.ID("has-premises"), mi.Name("hasPremises"),
						mi.Class("rounded border-gray-300 text-blue-600 focus:ring-blue-500"),
						mi.Data("dependency-trigger", "has-premises"),
					),
					b.Span(mi.Class("text-sm text-gray-700 dark:text-gray-300"), "Business has physical premises open to customers"),
				),
			),
			b.Div(mi.ID("premises-fields"), mi.Class("ml-6 space-y-4 hidden"),
				mi.Data("dependency-target", "premises-fields"),
				h.formField(b, "Premises Address", "premisesAddress", "text", "123 Main St", false),
				h.formField(b, "Premises Square Feet", "premisesSqft", "number", "5000", false),
			),
		),

		// Submit button
		b.Div(mi.Class("flex justify-end gap-3 pt-6 border-t border-gray-200 dark:border-gray-700"),
			b.Button(mi.Type("button"), mi.Class("px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600"),
				"Save Draft",
			),
			b.Button(mi.Type("submit"), mi.Class("px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 flex items-center gap-2"),
				"Continue",
				Icon("arrow-right", "w-4 h-4"),
			),
		),
	)
}

func (h *Handler) formField(b *mi.Builder, label, name, inputType, placeholder string, required bool) mi.Node {
	labelContent := []interface{}{
		mi.For(name), mi.Class("block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"),
		label,
	}
	if required {
		labelContent = append(labelContent, b.Span(mi.Class("text-red-500 ml-1"), "*"))
	}

	inputAttrs := []mi.Attribute{
		mi.Type(inputType), mi.ID(name), mi.Name(name), mi.Placeholder(placeholder),
		mi.Class("w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
	}
	if required {
		inputAttrs = append(inputAttrs, mi.Attr("required", "required"))
	}

	return b.Div(
		b.Label(labelContent...),
		b.Input(inputAttrs...),
	)
}

func (h *Handler) formSelect(b *mi.Builder, label, name string, options []string) mi.Node {
	var opts []interface{}
	opts = append(opts, b.Option(mi.Value(""), "Select..."))
	for _, opt := range options {
		opts = append(opts, b.Option(mi.Value(strings.ToLower(opt)), opt))
	}

	return b.Div(
		b.Label(mi.For(name), mi.Class("block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"), label),
		b.Select(append([]interface{}{
			mi.ID(name), mi.Name(name),
			mi.Class("w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
		}, opts...)...),
	)
}

// =============================================================================
// CLAIMS - Demonstrates CLIENT-SIDE FILTERABLE pattern
// =============================================================================

func (h *Handler) Claims(w http.ResponseWriter, r *http.Request) {
	page := h.pageLayout("claims", "Claims", "View and manage your insurance claims", func(b *mi.Builder) mi.Node {
		// PATTERN: ClientFilterable - JSON data with client-side filtering
		claimsFilter := mdy.Dyn("claims-filter").
			Data(mdy.FilterableDataset{
				Items: h.store.ClaimsAsMapSlice(),
				Schema: mdy.FilterSchema{
					Fields: []mdy.FilterableField{
						mdy.TextField("customerName", "Customer"),
						mdy.SelectField("status", "Status", []string{"open", "in-progress", "approved", "denied", "closed"}),
						mdy.SelectField("type", "Type", []string{"collision", "theft", "fire", "water", "weather", "liability", "medical", "glass"}),
					},
				},
				Options: mdy.FilterOptions{
					EnableSearch:     true,
					EnablePagination: true,
					ItemsPerPage:     5,
				},
			}).
			Theme(h.theme).
			Minified().
			Build()

		return b.Div(
			// Toolbar
			b.Div(mi.Class("flex items-center justify-between mb-4"),
				b.A(mi.Href("/claims/new"), mi.Class("inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"),
					Icon("plus", "w-4 h-4"), "File New Claim",
				),
			),
			// Filter component (generates controls and filters JSON data client-side)
			b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6"),
				claimsFilter(b),
			),
		)
	})
	h.render(w, page)
}

// =============================================================================
// COMPARE PLANS - Demonstrates TABS WITH DATA pattern
// =============================================================================

func (h *Handler) ComparePlans(w http.ResponseWriter, r *http.Request) {
	page := h.pageLayout("compare", "Compare Plans", "Find the perfect coverage for your needs", func(b *mi.Builder) mi.Node {
		// PATTERN: TabsWithData - Each tab shows filtered subset of plans
		// Build states for each coverage type
		states := []mdy.ComponentState{
			{ID: "all", Label: "All Plans", Active: true, Content: h.planGrid(b, "")},
			{ID: "auto", Label: "Auto", Icon: IconHTML("truck", "w-4 h-4 inline"), Content: h.planGrid(b, "auto")},
			{ID: "home", Label: "Home", Icon: IconHTML("home-modern", "w-4 h-4 inline"), Content: h.planGrid(b, "home")},
			{ID: "life", Label: "Life", Icon: IconHTML("heart", "w-4 h-4 inline"), Content: h.planGrid(b, "life")},
			{ID: "business", Label: "Business", Icon: IconHTML("building-office", "w-4 h-4 inline"), Content: h.planGrid(b, "business")},
		}

		planTabs := mdy.Dyn("plan-comparison").
			States(states).
			Theme(h.theme).
			Minified().
			Build()

		return b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6"),
			planTabs(b),
		)
	})
	h.render(w, page)
}

func (h *Handler) planGrid(b *mi.Builder, coverageType string) mi.Node {
	var plans []models.Plan
	if coverageType == "" {
		plans = h.store.Plans
	} else {
		plans = h.store.GetPlansByType(coverageType)
	}

	var cards []interface{}
	for _, plan := range plans {
		cards = append(cards, h.planCard(b, plan))
	}

	args := []interface{}{mi.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6")}
	args = append(args, cards...)
	return b.Div(args...)
}

func (h *Handler) planCard(b *mi.Builder, plan models.Plan) mi.Node {
	tierColors := map[string]string{
		"basic":    "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300",
		"standard": "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300",
		"premium":  "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300",
	}

	featureArgs := []interface{}{mi.Class("space-y-2 mb-6")}
	for _, f := range plan.Features {
		featureArgs = append(featureArgs, b.Li(mi.Class("flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400"),
			Icon("check", "w-4 h-4 text-green-500"),
			f,
		))
	}

	cardClass := "relative p-6 rounded-xl border transition-all hover:shadow-lg "
	if plan.Popular {
		cardClass += "border-blue-500 dark:border-blue-400"
	} else {
		cardClass += "border-gray-200 dark:border-gray-700"
	}

	cardContent := []interface{}{mi.Class(cardClass)}

	if plan.Popular {
		cardContent = append(cardContent,
			b.Div(mi.Class("absolute -top-3 left-1/2 transform -translate-x-1/2"),
				b.Span(mi.Class("px-3 py-1 text-xs font-medium text-white bg-blue-600 rounded-full"), "Most Popular"),
			),
		)
	}

	cardContent = append(cardContent,
		b.Div(mi.Class("flex items-center justify-between mb-4"),
			b.Span(mi.Class("text-lg font-semibold text-gray-900 dark:text-white"), plan.Name),
			b.Span(mi.Class("px-2 py-1 text-xs font-medium rounded "+tierColors[plan.Tier]), strings.Title(plan.Tier)),
		),
		b.Div(mi.Class("mb-4"),
			b.Span(mi.Class("text-3xl font-bold text-gray-900 dark:text-white"), fmt.Sprintf("$%.0f", plan.Price)),
			b.Span(mi.Class("text-gray-500 dark:text-gray-400"), "/month"),
		),
		b.Div(mi.Class("mb-4 text-sm text-gray-600 dark:text-gray-400"),
			b.P("Coverage: ", b.Span(mi.Class("font-medium"), fmt.Sprintf("$%,.0f", plan.Coverage))),
			b.P("Deductible: ", b.Span(mi.Class("font-medium"), fmt.Sprintf("$%,.0f", plan.Deductible))),
		),
		b.Ul(featureArgs...),
		b.A(mi.Href("/quote?plan="+plan.ID), mi.Class("block w-full text-center px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"),
			"Select Plan",
		),
	)

	return b.Div(cardContent...)
}

// =============================================================================
// SETTINGS - Demonstrates STATES pattern
// =============================================================================

func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	page := h.pageLayout("settings", "Settings", "Manage your account and preferences", func(b *mi.Builder) mi.Node {
		states := []mdy.ComponentState{
			{ID: "profile", Label: "Profile", Active: true, Content: h.settingsProfile(b)},
			{ID: "notifications", Label: "Notifications", Content: h.settingsNotifications(b)},
			{ID: "security", Label: "Security", Content: h.settingsSecurity(b)},
			{ID: "billing", Label: "Billing", Content: h.settingsBilling(b)},
		}

		settingsTabs := mdy.Dyn("settings-tabs").
			States(states).
			Theme(h.theme).
			Minified().
			Build()

		return b.Div(mi.Class("bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700"),
			b.Div(mi.Class("p-6"), settingsTabs(b)),
		)
	})
	h.render(w, page)
}

func (h *Handler) settingsProfile(b *mi.Builder) mi.Node {
	return b.Div(mi.Class("space-y-6"),
		b.Div(mi.Class("flex items-center gap-4"),
			b.Div(mi.Class("w-20 h-20 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center"),
				Icon("user", "w-10 h-10 text-gray-500 dark:text-gray-400"),
			),
			b.Div(
				b.Button(mi.Class("text-sm text-blue-600 dark:text-blue-400 hover:underline"), "Change photo"),
			),
		),
		b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4"),
			h.formField(b, "First Name", "firstName", "text", "John", false),
			h.formField(b, "Last Name", "lastName", "text", "Smith", false),
			h.formField(b, "Email", "email", "email", "john@example.com", false),
			h.formField(b, "Phone", "phone", "tel", "(555) 123-4567", false),
		),
		b.Div(mi.Class("flex justify-end"),
			b.Button(mi.Class("px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"), "Save Changes"),
		),
	)
}

func (h *Handler) settingsNotifications(b *mi.Builder) mi.Node {
	notifications := []struct{ Label, Desc string }{
		{"Email notifications", "Receive updates about your policies via email"},
		{"SMS alerts", "Get text messages for important updates"},
		{"Payment reminders", "Reminder before payment is due"},
		{"Claim updates", "Notifications when claim status changes"},
		{"Marketing emails", "Special offers and new products"},
	}

	var items []interface{}
	for i, n := range notifications {
		checkboxAttrs := []mi.Attribute{mi.Type("checkbox"), mi.Class("sr-only peer")}
		if i < 4 {
			checkboxAttrs = append(checkboxAttrs, mi.Attr("checked", "checked"))
		}

		items = append(items, b.Div(mi.Class("flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700 last:border-0"),
			b.Div(
				b.P(mi.Class("font-medium text-gray-900 dark:text-white"), n.Label),
				b.P(mi.Class("text-sm text-gray-500 dark:text-gray-400"), n.Desc),
			),
			b.Label(mi.Class("relative inline-flex items-center cursor-pointer"),
				b.Input(checkboxAttrs...),
				b.Div(mi.Class("w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600")),
			),
		))
	}

	args := []interface{}{mi.Class("space-y-2")}
	args = append(args, items...)
	return b.Div(args...)
}

func (h *Handler) settingsSecurity(b *mi.Builder) mi.Node {
	return b.Div(mi.Class("space-y-6"),
		b.Div(
			b.H3(mi.Class("font-medium text-gray-900 dark:text-white mb-4"), "Change Password"),
			b.Div(mi.Class("space-y-4 max-w-md"),
				h.formField(b, "Current Password", "currentPassword", "password", "", false),
				h.formField(b, "New Password", "newPassword", "password", "", false),
				h.formField(b, "Confirm Password", "confirmPassword", "password", "", false),
			),
			b.Button(mi.Class("mt-4 px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"), "Update Password"),
		),
		b.Div(mi.Class("pt-6 border-t border-gray-200 dark:border-gray-700"),
			b.H3(mi.Class("font-medium text-gray-900 dark:text-white mb-4"), "Two-Factor Authentication"),
			b.Div(mi.Class("flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg"),
				b.Div(mi.Class("flex items-center gap-3"),
					Icon("shield-check", "w-8 h-8 text-green-600 dark:text-green-400"),
					b.Div(
						b.P(mi.Class("font-medium text-gray-900 dark:text-white"), "2FA is enabled"),
						b.P(mi.Class("text-sm text-gray-500 dark:text-gray-400"), "Your account is protected"),
					),
				),
				b.Button(mi.Class("text-sm text-blue-600 dark:text-blue-400 hover:underline"), "Manage"),
			),
		),
	)
}

func (h *Handler) settingsBilling(b *mi.Builder) mi.Node {
	return b.Div(mi.Class("space-y-6"),
		b.Div(
			b.H3(mi.Class("font-medium text-gray-900 dark:text-white mb-4"), "Payment Method"),
			b.Div(mi.Class("flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg"),
				b.Div(mi.Class("flex items-center gap-3"),
					b.Div(mi.Class("w-12 h-8 bg-blue-900 rounded flex items-center justify-center text-white text-xs font-bold"), "VISA"),
					b.Div(
						b.P(mi.Class("font-medium text-gray-900 dark:text-white"), "•••• •••• •••• 4242"),
						b.P(mi.Class("text-sm text-gray-500 dark:text-gray-400"), "Expires 12/25"),
					),
				),
				b.Button(mi.Class("text-sm text-blue-600 dark:text-blue-400 hover:underline"), "Update"),
			),
		),
		b.Div(mi.Class("pt-6 border-t border-gray-200 dark:border-gray-700"),
			b.H3(mi.Class("font-medium text-gray-900 dark:text-white mb-4"), "Billing History"),
			b.Table(mi.Class("w-full"),
				b.Thead(
					b.Tr(mi.Class("text-left text-sm text-gray-500 dark:text-gray-400"),
						b.Th(mi.Class("pb-3"), "Date"),
						b.Th(mi.Class("pb-3"), "Description"),
						b.Th(mi.Class("pb-3 text-right"), "Amount"),
						b.Th(mi.Class("pb-3"), ""),
					),
				),
				b.Tbody(mi.Class("text-sm"),
					h.billingRow(b, "Dec 1, 2024", "Monthly Premium", "$485.00"),
					h.billingRow(b, "Nov 1, 2024", "Monthly Premium", "$485.00"),
					h.billingRow(b, "Oct 1, 2024", "Monthly Premium", "$485.00"),
				),
			),
		),
	)
}

func (h *Handler) billingRow(b *mi.Builder, date, desc, amount string) mi.Node {
	return b.Tr(mi.Class("border-t border-gray-200 dark:border-gray-700"),
		b.Td(mi.Class("py-3 text-gray-600 dark:text-gray-400"), date),
		b.Td(mi.Class("py-3 text-gray-900 dark:text-white"), desc),
		b.Td(mi.Class("py-3 text-right text-gray-900 dark:text-white"), amount),
		b.Td(mi.Class("py-3"),
			b.A(mi.Href("#"), mi.Class("text-blue-600 dark:text-blue-400 hover:underline flex items-center gap-1"),
				Icon("arrow-down-tray", "w-4 h-4"), "PDF",
			),
		),
	)
}

// =============================================================================
// RENDER
// =============================================================================

func (h *Handler) render(w http.ResponseWriter, page mi.H) {
	var buf bytes.Buffer
	if err := mi.Render(page, &buf); err != nil {
		h.logger.Error("render failed", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
