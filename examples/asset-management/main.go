// Package main demonstrates a complex Asset Management System UI
package main

import (
	"fmt"
	"os"
	"time"

	mi "github.com/ha1tch/minty"
	"github.com/ha1tch/minty/themes/tailwind"
)

// =====================================================
// DOMAIN TYPES
// =====================================================

type Asset struct {
	ID           string
	Tag          string
	Name         string
	Category     string
	Status       string
	Location     string
	Department   string
	AssignedTo   string
	PurchaseDate string
	PurchaseCost float64
	CurrentValue float64
	Vendor       string
	SerialNumber string
	Model        string
	Warranty     string
	Notes        string
}

type MaintenanceRecord struct {
	ID          string
	AssetID     string
	Date        string
	Type        string
	Description string
	Cost        float64
	Technician  string
	Status      string
}

type AuditLog struct {
	ID        string
	Timestamp string
	User      string
	Action    string
	AssetID   string
	Details   string
}

// Sample data
var sampleAssets = []Asset{
	{"A001", "IT-LAP-001", "MacBook Pro 16\"", "Laptops", "active", "HQ Floor 3", "Engineering", "John Smith", "2024-01-15", 2499.00, 2100.00, "Apple Inc.", "C02XG123HKGY", "MacBook Pro 16 M3", "2027-01-15", "Primary development machine"},
	{"A002", "IT-LAP-002", "ThinkPad X1 Carbon", "Laptops", "active", "HQ Floor 2", "Sales", "Jane Doe", "2024-02-20", 1899.00, 1650.00, "Lenovo", "PF3ABCD1", "X1 Carbon Gen 11", "2027-02-20", ""},
	{"A003", "IT-MON-001", "Dell U2723QE", "Monitors", "active", "HQ Floor 3", "Engineering", "John Smith", "2024-01-15", 799.00, 700.00, "Dell", "CN0M2K831234", "UltraSharp 27 4K", "2027-01-15", "4K USB-C Hub Monitor"},
	{"A004", "IT-SRV-001", "Dell PowerEdge R750", "Servers", "active", "Data Center", "IT Operations", "Unassigned", "2023-06-01", 12500.00, 10000.00, "Dell", "SVCTAG001", "PowerEdge R750xs", "2026-06-01", "Primary application server"},
	{"A005", "IT-LAP-003", "MacBook Air M2", "Laptops", "maintenance", "IT Storage", "Marketing", "Bob Wilson", "2024-03-10", 1299.00, 1150.00, "Apple Inc.", "C02YH456JKLY", "MacBook Air 13 M2", "2027-03-10", "Battery replacement scheduled"},
	{"A006", "IT-NET-001", "Cisco Catalyst 9300", "Network", "active", "Server Room A", "IT Operations", "Unassigned", "2023-01-20", 4500.00, 3800.00, "Cisco", "FCW2345K001", "C9300-48P", "2026-01-20", "48-port PoE+ switch"},
	{"A007", "IT-PRN-001", "HP LaserJet Pro", "Printers", "active", "HQ Floor 2", "Shared", "Unassigned", "2024-04-05", 549.00, 500.00, "HP", "VNB3R12345", "M404dn", "2025-04-05", "Department printer"},
	{"A008", "IT-LAP-004", "Dell XPS 15", "Laptops", "retired", "IT Storage", "Finance", "Unassigned", "2021-05-15", 1799.00, 0.00, "Dell", "5CG123ABC", "XPS 15 9510", "2024-05-15", "End of life - data wiped"},
}

var theme = tailwind.NewTailwindTheme()

// =====================================================
// UI COMPONENTS
// =====================================================

// Icon renders a simple icon placeholder (in production, use actual SVG icons)
func Icon(name string) mi.H {
	icons := map[string]string{
		"dashboard":    "📊",
		"assets":       "💻",
		"maintenance":  "🔧",
		"reports":      "📈",
		"settings":     "⚙️",
		"users":        "👥",
		"search":       "🔍",
		"filter":       "⏳",
		"add":          "➕",
		"edit":         "✏️",
		"delete":       "🗑️",
		"view":         "👁️",
		"export":       "📤",
		"import":       "📥",
		"refresh":      "🔄",
		"notification": "🔔",
		"help":         "❓",
		"menu":         "☰",
		"close":        "✕",
		"check":        "✓",
		"warning":      "⚠️",
		"info":         "ℹ️",
		"chevron-down": "▼",
		"chevron-right": "▶",
	}
	icon := icons[name]
	if icon == "" {
		icon = "•"
	}
	return func(b *mi.Builder) mi.Node {
		return b.Span(mi.Class("icon"), icon)
	}
}

// StatusBadge renders a coloured status indicator
func StatusBadge(status string) mi.H {
	return func(b *mi.Builder) mi.Node {
		colors := map[string]string{
			"active":      "bg-green-100 text-green-800",
			"maintenance": "bg-yellow-100 text-yellow-800",
			"retired":     "bg-gray-100 text-gray-600",
			"pending":     "bg-blue-100 text-blue-800",
			"critical":    "bg-red-100 text-red-800",
		}
		colorClass := colors[status]
		if colorClass == "" {
			colorClass = "bg-gray-100 text-gray-600"
		}
		return b.Span(
			mi.Class("px-2 py-1 text-xs font-medium rounded-full "+colorClass),
			status,
		)
	}
}

// Dropdown menu component
func Dropdown(label string, items []string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("relative inline-block"),
			b.Button(
				mi.Class("inline-flex items-center gap-1 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"),
				mi.Type("button"),
				label,
				Icon("chevron-down")(b),
			),
			// Dropdown menu (hidden by default, shown via JS/HTMX)
			b.Div(
				mi.Class("hidden absolute right-0 z-10 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5"),
				mi.ID("dropdown-"+label),
				func() mi.Node {
					menuItems := make([]mi.Node, len(items))
					for i, item := range items {
						menuItems[i] = b.A(
							mi.Href("#"),
							mi.Class("block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"),
							item,
						)
					}
					return mi.NewFragment(menuItems...)
				}(),
			),
		)
	}
}

// Tab component
func Tab(id, label string, active bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		baseClass := "tab-button px-4 py-2 text-sm font-medium border-b-2 transition-colors cursor-pointer"
		if active {
			baseClass += " text-blue-600 border-blue-600 active"
		} else {
			baseClass += " text-gray-500 border-transparent hover:text-gray-700 hover:border-gray-300"
		}
		
		return b.Button(
			mi.Class(baseClass),
			mi.Type("button"),
			mi.DataAttr("tab", id),
			mi.Attr("onclick", "switchTab('"+id+"')"),
			label,
		)
	}
}

// TabBar component
func TabBar(tabs []struct{ ID, Label string }, activeTab string) mi.H {
	return func(b *mi.Builder) mi.Node {
		tabNodes := make([]mi.Node, len(tabs))
		for i, tab := range tabs {
			tabNodes[i] = Tab(tab.ID, tab.Label, tab.ID == activeTab)(b)
		}
		return b.Div(mi.Class("flex border-b border-gray-200 mb-4"),
			mi.NewFragment(tabNodes...),
		)
	}
}

// FormField renders a labelled form field
func FormField(label, name, fieldType, placeholder, value string, required bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "field-" + name
		attrs := []mi.Attribute{
			mi.Class("w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"),
			mi.ID(id),
			mi.Name(name),
			mi.Type(fieldType),
			mi.Placeholder(placeholder),
			mi.Value(value),
		}
		if required {
			attrs = append(attrs, mi.Required())
		}
		
		labelClass := "block text-sm font-medium text-gray-700 mb-1"
		if required {
			labelClass += " after:content-['*'] after:ml-0.5 after:text-red-500"
		}
		
		return b.Div(mi.Class("mb-4"),
			b.Label(mi.Class(labelClass), mi.For(id), label),
			b.Input(attrs...),
		)
	}
}

// SelectField renders a select dropdown
func SelectField(label, name string, options []struct{ Value, Text string }, selected string, required bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "field-" + name
		
		optionNodes := make([]mi.Node, len(options)+1)
		optionNodes[0] = b.Option(mi.Value(""), "Select...")
		for i, opt := range options {
			attrs := []interface{}{mi.Value(opt.Value)}
			if opt.Value == selected {
				attrs = append(attrs, mi.Selected())
			}
			attrs = append(attrs, opt.Text)
			optionNodes[i+1] = b.Option(attrs...)
		}
		
		labelClass := "block text-sm font-medium text-gray-700 mb-1"
		if required {
			labelClass += " after:content-['*'] after:ml-0.5 after:text-red-500"
		}
		
		selectAttrs := []interface{}{
			mi.Class("w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"),
			mi.ID(id),
			mi.Name(name),
		}
		if required {
			selectAttrs = append(selectAttrs, mi.Required())
		}
		selectAttrs = append(selectAttrs, mi.NewFragment(optionNodes...))
		
		return b.Div(mi.Class("mb-4"),
			b.Label(mi.Class(labelClass), mi.For(id), label),
			b.Select(selectAttrs...),
		)
	}
}

// TextareaField renders a textarea
func TextareaField(label, name, placeholder, value string, rows int) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "field-" + name
		return b.Div(mi.Class("mb-4"),
			b.Label(mi.Class("block text-sm font-medium text-gray-700 mb-1"), mi.For(id), label),
			b.Textarea(
				mi.Class("w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"),
				mi.ID(id),
				mi.Name(name),
				mi.Placeholder(placeholder),
				mi.Rows(rows),
				value,
			),
		)
	}
}

// Card component
func Card(title string, content mi.H, actions mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("bg-white rounded-lg shadow-sm border border-gray-200"),
			mi.If(title != "", func(b *mi.Builder) mi.Node {
				return b.Div(mi.Class("px-4 py-3 border-b border-gray-200 flex items-center justify-between"),
					b.H3(mi.Class("text-lg font-medium text-gray-900"), title),
					mi.If(actions != nil, actions)(b),
				)
			})(b),
			b.Div(mi.Class("p-4"),
				content(b),
			),
		)
	}
}

// StatCard for dashboard metrics
func StatCard(title, value, change string, positive bool, icon string) mi.H {
	return func(b *mi.Builder) mi.Node {
		changeColor := "text-green-600"
		if !positive {
			changeColor = "text-red-600"
		}
		return b.Div(mi.Class("bg-white rounded-lg shadow-sm border border-gray-200 p-4"),
			b.Div(mi.Class("flex items-center justify-between"),
				b.Div(
					b.P(mi.Class("text-sm font-medium text-gray-500"), title),
					b.P(mi.Class("text-2xl font-semibold text-gray-900 mt-1"), value),
					b.P(mi.Class("text-sm mt-1 "+changeColor), change),
				),
				b.Div(mi.Class("text-3xl opacity-20"), Icon(icon)(b)),
			),
		)
	}
}

// =====================================================
// MAIN LAYOUT SECTIONS
// =====================================================

// Sidebar navigation
func Sidebar() mi.H {
	return func(b *mi.Builder) mi.Node {
		navItems := []struct {
			Icon, Label, Href string
			Active            bool
		}{
			{"dashboard", "Dashboard", "/", false},
			{"assets", "Assets", "/assets", true},
			{"maintenance", "Maintenance", "/maintenance", false},
			{"reports", "Reports", "/reports", false},
			{"users", "Users", "/users", false},
			{"settings", "Settings", "/settings", false},
		}
		
		navNodes := make([]mi.Node, len(navItems))
		for i, item := range navItems {
			class := "flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-lg transition-colors"
			if item.Active {
				class += " bg-blue-50 text-blue-700"
			} else {
				class += " text-gray-600 hover:bg-gray-100 hover:text-gray-900"
			}
			navNodes[i] = b.A(
				mi.Href(item.Href),
				mi.Class(class),
				Icon(item.Icon)(b),
				item.Label,
			)
		}
		
		return b.Aside(mi.Class("w-64 bg-white border-r border-gray-200 min-h-screen"),
			// Logo
			b.Div(mi.Class("p-4 border-b border-gray-200"),
				b.H1(mi.Class("text-xl font-bold text-gray-900"), "AssetTrack"),
				b.P(mi.Class("text-xs text-gray-500"), "Enterprise Asset Management"),
			),
			// Navigation
			b.Nav(mi.Class("p-4 space-y-1"),
				mi.NewFragment(navNodes...),
			),
			// Bottom section
			b.Div(mi.Class("absolute bottom-0 left-0 w-64 p-4 border-t border-gray-200 bg-white"),
				b.Div(mi.Class("flex items-center gap-3"),
					b.Div(mi.Class("w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-medium"), "JD"),
					b.Div(
						b.P(mi.Class("text-sm font-medium text-gray-900"), "John Doe"),
						b.P(mi.Class("text-xs text-gray-500"), "Administrator"),
					),
				),
			),
		)
	}
}

// Header bar
func Header() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Header(mi.Class("bg-white border-b border-gray-200 px-6 py-3"),
			b.Div(mi.Class("flex items-center justify-between"),
				// Search
				b.Div(mi.Class("flex-1 max-w-lg"),
					b.Div(mi.Class("relative"),
						b.Span(mi.Class("absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"),
							Icon("search")(b),
						),
						b.Input(
							mi.Type("search"),
							mi.Placeholder("Search assets, users, locations..."),
							mi.Class("w-full pl-10 pr-4 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"),
							mi.HtmxGet("/search"),
							mi.HtmxTrigger("keyup changed delay:300ms"),
							mi.HtmxTarget("#search-results"),
						),
					),
				),
				// Actions
				b.Div(mi.Class("flex items-center gap-4"),
					b.Button(
						mi.Class("p-2 text-gray-400 hover:text-gray-600 relative"),
						mi.Type("button"),
						Icon("notification")(b),
						b.Span(mi.Class("absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full")),
					),
					b.Button(
						mi.Class("p-2 text-gray-400 hover:text-gray-600"),
						mi.Type("button"),
						Icon("help")(b),
					),
					Dropdown("Quick Actions", []string{"Add Asset", "Run Report", "Export Data", "Import Data"})(b),
				),
			),
		)
	}
}

// Asset table
func AssetTable(assets []Asset) mi.H {
	return func(b *mi.Builder) mi.Node {
		rows := make([]mi.Node, len(assets))
		for i, asset := range assets {
			rows[i] = b.Tr(mi.Class("hover:bg-gray-50"),
				b.Td(mi.Class("px-4 py-3"),
					b.Input(mi.Type("checkbox"), mi.Class("rounded border-gray-300")),
				),
				b.Td(mi.Class("px-4 py-3"),
					b.Div(
						b.P(mi.Class("font-medium text-gray-900"), asset.Name),
						b.P(mi.Class("text-xs text-gray-500"), asset.Tag),
					),
				),
				b.Td(mi.Class("px-4 py-3 text-sm text-gray-600"), asset.Category),
				b.Td(mi.Class("px-4 py-3"), StatusBadge(asset.Status)(b)),
				b.Td(mi.Class("px-4 py-3 text-sm text-gray-600"), asset.Location),
				b.Td(mi.Class("px-4 py-3 text-sm text-gray-600"), asset.AssignedTo),
				b.Td(mi.Class("px-4 py-3 text-sm text-gray-600"), fmt.Sprintf("$%.2f", asset.CurrentValue)),
				b.Td(mi.Class("px-4 py-3"),
					b.Div(mi.Class("flex items-center gap-2"),
						b.Button(
							mi.Class("p-1 text-gray-400 hover:text-blue-600"),
							mi.Type("button"),
							mi.HtmxGet("/assets/"+asset.ID),
							mi.HtmxTarget("#modal-content"),
							mi.Attr("title", "View"),
							Icon("view")(b),
						),
						b.Button(
							mi.Class("p-1 text-gray-400 hover:text-blue-600"),
							mi.Type("button"),
							mi.HtmxGet("/assets/"+asset.ID+"/edit"),
							mi.HtmxTarget("#modal-content"),
							mi.Attr("title", "Edit"),
							Icon("edit")(b),
						),
						b.Button(
							mi.Class("p-1 text-gray-400 hover:text-red-600"),
							mi.Type("button"),
							mi.Attr("title", "Delete"),
							Icon("delete")(b),
						),
					),
				),
			)
		}
		
		return b.Div(mi.Class("overflow-x-auto"),
			b.Table(mi.Class("w-full"),
				b.Thead(mi.Class("bg-gray-50 border-y border-gray-200"),
					b.Tr(
						b.Th(mi.Class("px-4 py-3 text-left"),
							b.Input(mi.Type("checkbox"), mi.Class("rounded border-gray-300")),
						),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Asset"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Category"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Status"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Location"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Assigned To"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Value"),
						b.Th(mi.Class("px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), "Actions"),
					),
				),
				b.Tbody(mi.Class("divide-y divide-gray-200"),
					mi.NewFragment(rows...),
				),
			),
		)
	}
}

// Asset detail form
func AssetDetailForm(asset Asset) mi.H {
	return func(b *mi.Builder) mi.Node {
		categories := []struct{ Value, Text string }{
			{"Laptops", "Laptops"},
			{"Monitors", "Monitors"},
			{"Servers", "Servers"},
			{"Network", "Network Equipment"},
			{"Printers", "Printers"},
			{"Mobile", "Mobile Devices"},
			{"Software", "Software Licenses"},
			{"Other", "Other"},
		}
		
		statuses := []struct{ Value, Text string }{
			{"active", "Active"},
			{"maintenance", "Under Maintenance"},
			{"retired", "Retired"},
			{"pending", "Pending"},
		}
		
		departments := []struct{ Value, Text string }{
			{"Engineering", "Engineering"},
			{"Sales", "Sales"},
			{"Marketing", "Marketing"},
			{"Finance", "Finance"},
			{"HR", "Human Resources"},
			{"IT Operations", "IT Operations"},
			{"Shared", "Shared/Common"},
		}
		
		tabs := []struct{ ID, Label string }{
			{"details", "Details"},
			{"financial", "Financial"},
			{"maintenance", "Maintenance"},
			{"history", "History"},
			{"documents", "Documents"},
		}
		
		return b.Form(
			mi.Class("space-y-6"),
			mi.HtmxPost("/assets/"+asset.ID),
			mi.HtmxTarget("#form-result"),
			
			// Tabs
			TabBar(tabs, "details")(b),
			
			// Tab panels container
			b.Div(mi.ID("tab-content"),
				// Details Tab Panel
				b.Div(mi.Class("tab-panel"), mi.ID("panel-details"),
					// Basic Information Section
					b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-6"),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Basic Information"),
							FormField("Asset Tag", "tag", "text", "e.g., IT-LAP-001", asset.Tag, true)(b),
							FormField("Asset Name", "name", "text", "e.g., MacBook Pro 16\"", asset.Name, true)(b),
							SelectField("Category", "category", categories, asset.Category, true)(b),
							SelectField("Status", "status", statuses, asset.Status, true)(b),
						),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Assignment"),
							SelectField("Department", "department", departments, asset.Department, true)(b),
							FormField("Assigned To", "assigned_to", "text", "Employee name", asset.AssignedTo, false)(b),
							FormField("Location", "location", "text", "e.g., HQ Floor 3", asset.Location, true)(b),
						),
					),
					
					// Hardware Details Section
					b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-6 mt-6 pt-6 border-t border-gray-200"),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Hardware Details"),
							FormField("Manufacturer/Vendor", "vendor", "text", "e.g., Apple Inc.", asset.Vendor, false)(b),
							FormField("Model", "model", "text", "e.g., MacBook Pro 16 M3", asset.Model, false)(b),
							FormField("Serial Number", "serial_number", "text", "Device serial number", asset.SerialNumber, false)(b),
						),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Warranty"),
							FormField("Purchase Date", "purchase_date", "date", "", asset.PurchaseDate, false)(b),
							FormField("Warranty Expiry", "warranty", "date", "", asset.Warranty, false)(b),
						),
					),
					
					// Notes Section
					b.Div(mi.Class("mt-6 pt-6 border-t border-gray-200"),
						TextareaField("Notes", "notes", "Additional notes about this asset...", asset.Notes, 4)(b),
					),
				),
				
				// Financial Tab Panel
				b.Div(mi.Class("tab-panel hidden"), mi.ID("panel-financial"),
					b.Div(mi.Class("grid grid-cols-1 md:grid-cols-3 gap-6"),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Purchase Information"),
							FormField("Purchase Date", "fin_purchase_date", "date", "", asset.PurchaseDate, false)(b),
							FormField("Purchase Cost", "purchase_cost", "number", "0.00", fmt.Sprintf("%.2f", asset.PurchaseCost), false)(b),
							FormField("Purchase Order #", "po_number", "text", "PO-0000", "", false)(b),
							FormField("Invoice #", "invoice_number", "text", "INV-0000", "", false)(b),
						),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Current Valuation"),
							FormField("Current Value", "current_value", "number", "0.00", fmt.Sprintf("%.2f", asset.CurrentValue), false)(b),
							SelectField("Depreciation Method", "depreciation_method", []struct{ Value, Text string }{
								{"straight-line", "Straight Line"},
								{"declining-balance", "Declining Balance"},
								{"sum-of-years", "Sum of Years' Digits"},
							}, "straight-line", false)(b),
							FormField("Useful Life (years)", "useful_life", "number", "5", "5", false)(b),
							FormField("Salvage Value", "salvage_value", "number", "0.00", "0.00", false)(b),
						),
						b.Div(
							b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Budget & Cost Center"),
							FormField("Cost Center", "cost_center", "text", "e.g., CC-1001", "", false)(b),
							FormField("GL Account", "gl_account", "text", "e.g., 1520-00", "", false)(b),
							FormField("Budget Code", "budget_code", "text", "", "", false)(b),
						),
					),
					// Depreciation schedule
					b.Div(mi.Class("mt-6 pt-6 border-t border-gray-200"),
						b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Depreciation Schedule"),
						b.Table(mi.Class("w-full text-sm"),
							b.Thead(mi.Class("bg-gray-50"),
								b.Tr(
									b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Year"),
									b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Beginning Value"),
									b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Depreciation"),
									b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Ending Value"),
								),
							),
							b.Tbody(mi.Class("divide-y divide-gray-200"),
								depreciationRow(b, "2024", 2499.00, 499.80, 1999.20),
								depreciationRow(b, "2025", 1999.20, 499.80, 1499.40),
								depreciationRow(b, "2026", 1499.40, 499.80, 999.60),
								depreciationRow(b, "2027", 999.60, 499.80, 499.80),
								depreciationRow(b, "2028", 499.80, 499.80, 0.00),
							),
						),
					),
				),
				
				// Maintenance Tab Panel
				b.Div(mi.Class("tab-panel hidden"), mi.ID("panel-maintenance"),
					// Schedule maintenance button
					b.Div(mi.Class("flex justify-between items-center mb-4"),
						b.H4(mi.Class("text-sm font-medium text-gray-900"), "Maintenance History"),
						b.Button(
							mi.Class("inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"),
							mi.Type("button"),
							Icon("add")(b),
							"Schedule Maintenance",
						),
					),
					// Maintenance records table
					b.Table(mi.Class("w-full text-sm"),
						b.Thead(mi.Class("bg-gray-50"),
							b.Tr(
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Date"),
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Type"),
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Description"),
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Technician"),
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Cost"),
								b.Th(mi.Class("px-4 py-2 text-left text-xs font-medium text-gray-500"), "Status"),
							),
						),
						b.Tbody(mi.Class("divide-y divide-gray-200"),
							maintenanceRow(b, "2025-01-02", "Scheduled", "Annual checkup and cleaning", "Mike Tech", 75.00, "completed"),
							maintenanceRow(b, "2024-09-15", "Repair", "Battery replacement", "Sarah Fix", 199.00, "completed"),
							maintenanceRow(b, "2024-06-20", "Upgrade", "RAM upgrade to 32GB", "Mike Tech", 350.00, "completed"),
							maintenanceRow(b, "2024-03-10", "Scheduled", "Initial setup and configuration", "IT Team", 0.00, "completed"),
						),
					),
					// Maintenance summary
					b.Div(mi.Class("mt-6 pt-6 border-t border-gray-200 grid grid-cols-3 gap-4"),
						b.Div(mi.Class("text-center p-4 bg-gray-50 rounded-lg"),
							b.P(mi.Class("text-2xl font-semibold text-gray-900"), "4"),
							b.P(mi.Class("text-sm text-gray-500"), "Total Records"),
						),
						b.Div(mi.Class("text-center p-4 bg-gray-50 rounded-lg"),
							b.P(mi.Class("text-2xl font-semibold text-gray-900"), "$624.00"),
							b.P(mi.Class("text-sm text-gray-500"), "Total Cost"),
						),
						b.Div(mi.Class("text-center p-4 bg-gray-50 rounded-lg"),
							b.P(mi.Class("text-2xl font-semibold text-gray-900"), "Jan 2025"),
							b.P(mi.Class("text-sm text-gray-500"), "Last Service"),
						),
					),
				),
				
				// History Tab Panel
				b.Div(mi.Class("tab-panel hidden"), mi.ID("panel-history"),
					b.H4(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Asset History & Audit Trail"),
					b.Div(mi.Class("space-y-4"),
						historyItem(b, "2025-01-03 14:32", "John Doe", "Updated", "Changed status from 'maintenance' to 'active'"),
						historyItem(b, "2025-01-02 09:15", "System", "Maintenance", "Scheduled maintenance completed"),
						historyItem(b, "2024-12-15 11:20", "Jane Smith", "Reassigned", "Transferred from Bob Wilson to John Smith"),
						historyItem(b, "2024-09-15 16:45", "Mike Tech", "Repair", "Battery replacement completed"),
						historyItem(b, "2024-06-20 10:30", "Mike Tech", "Upgrade", "RAM upgraded from 16GB to 32GB"),
						historyItem(b, "2024-03-10 08:00", "IT Team", "Deployed", "Initial deployment to Engineering department"),
						historyItem(b, "2024-01-15 09:00", "System", "Created", "Asset record created"),
					),
				),
				
				// Documents Tab Panel
				b.Div(mi.Class("tab-panel hidden"), mi.ID("panel-documents"),
					// Upload button
					b.Div(mi.Class("flex justify-between items-center mb-4"),
						b.H4(mi.Class("text-sm font-medium text-gray-900"), "Attached Documents"),
						b.Button(
							mi.Class("inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"),
							mi.Type("button"),
							Icon("add")(b),
							"Upload Document",
						),
					),
					// Documents grid
					b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 gap-4"),
						documentCard(b, "Purchase Invoice", "invoice-2024-001.pdf", "PDF", "245 KB", "2024-01-15"),
						documentCard(b, "Warranty Certificate", "warranty-macbook.pdf", "PDF", "128 KB", "2024-01-15"),
						documentCard(b, "User Manual", "macbook-pro-manual.pdf", "PDF", "4.2 MB", "2024-01-15"),
						documentCard(b, "Service Report - Sep 2024", "service-report-092024.pdf", "PDF", "89 KB", "2024-09-15"),
					),
					// Drop zone
					b.Div(mi.Class("mt-6 border-2 border-dashed border-gray-300 rounded-lg p-8 text-center"),
						b.Div(mi.Class("text-gray-400 text-4xl mb-2"), "📄"),
						b.P(mi.Class("text-sm text-gray-600"), "Drag and drop files here, or"),
						b.Button(
							mi.Class("mt-2 text-sm text-blue-600 hover:text-blue-700 font-medium"),
							mi.Type("button"),
							"browse to upload",
						),
						b.P(mi.Class("text-xs text-gray-400 mt-2"), "PDF, DOC, XLS, PNG, JPG up to 10MB"),
					),
				),
			),
			
			// Form Actions
			b.Div(mi.Class("flex items-center justify-end gap-3 pt-6 border-t border-gray-200"),
				b.Button(
					mi.Class("px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"),
					mi.Type("button"),
					"Cancel",
				),
				b.Button(
					mi.Class("px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"),
					mi.Type("submit"),
					"Save Changes",
				),
			),
			b.Div(mi.ID("form-result")),
		)
	}
}

// Helper functions for tab content

func depreciationRow(b *mi.Builder, year string, beginning, depreciation, ending float64) mi.Node {
	return b.Tr(
		b.Td(mi.Class("px-4 py-2 text-gray-900"), year),
		b.Td(mi.Class("px-4 py-2 text-gray-600"), fmt.Sprintf("$%.2f", beginning)),
		b.Td(mi.Class("px-4 py-2 text-gray-600"), fmt.Sprintf("$%.2f", depreciation)),
		b.Td(mi.Class("px-4 py-2 text-gray-900 font-medium"), fmt.Sprintf("$%.2f", ending)),
	)
}

func maintenanceRow(b *mi.Builder, date, mtype, desc, tech string, cost float64, status string) mi.Node {
	return b.Tr(mi.Class("hover:bg-gray-50"),
		b.Td(mi.Class("px-4 py-2 text-gray-900"), date),
		b.Td(mi.Class("px-4 py-2"),
			b.Span(mi.Class("px-2 py-0.5 text-xs rounded bg-blue-100 text-blue-800"), mtype),
		),
		b.Td(mi.Class("px-4 py-2 text-gray-600"), desc),
		b.Td(mi.Class("px-4 py-2 text-gray-600"), tech),
		b.Td(mi.Class("px-4 py-2 text-gray-600"), fmt.Sprintf("$%.2f", cost)),
		b.Td(mi.Class("px-4 py-2"), StatusBadge(status)(b)),
	)
}

func historyItem(b *mi.Builder, timestamp, user, action, details string) mi.Node {
	return b.Div(mi.Class("flex gap-4 p-3 bg-gray-50 rounded-lg"),
		b.Div(mi.Class("flex-shrink-0 w-2 h-2 mt-2 rounded-full bg-blue-500")),
		b.Div(mi.Class("flex-1"),
			b.Div(mi.Class("flex items-center gap-2 mb-1"),
				b.Span(mi.Class("text-sm font-medium text-gray-900"), user),
				b.Span(mi.Class("px-2 py-0.5 text-xs rounded bg-gray-200 text-gray-700"), action),
			),
			b.P(mi.Class("text-sm text-gray-600"), details),
			b.P(mi.Class("text-xs text-gray-400 mt-1"), timestamp),
		),
	)
}

func documentCard(b *mi.Builder, title, filename, doctype, size, date string) mi.Node {
	return b.Div(mi.Class("flex items-center gap-3 p-3 border border-gray-200 rounded-lg hover:bg-gray-50"),
		b.Div(mi.Class("flex-shrink-0 w-10 h-10 bg-red-100 rounded flex items-center justify-center text-red-600 text-xs font-medium"),
			doctype,
		),
		b.Div(mi.Class("flex-1 min-w-0"),
			b.P(mi.Class("text-sm font-medium text-gray-900 truncate"), title),
			b.P(mi.Class("text-xs text-gray-500"), filename),
			b.P(mi.Class("text-xs text-gray-400"), size+" • "+date),
		),
		b.Div(mi.Class("flex items-center gap-1"),
			b.Button(mi.Class("p-1 text-gray-400 hover:text-blue-600"), mi.Type("button"), mi.Attr("title", "Download"), "↓"),
			b.Button(mi.Class("p-1 text-gray-400 hover:text-red-600"), mi.Type("button"), mi.Attr("title", "Delete"), "×"),
		),
	)
}

// Filter sidebar
func FilterSidebar() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("w-64 bg-white border-r border-gray-200 p-4"),
			b.H3(mi.Class("text-sm font-medium text-gray-900 mb-4"), "Filters"),
			
			// Category filter
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("block text-xs font-medium text-gray-500 mb-2"), "Category"),
				b.Div(mi.Class("space-y-2"),
					filterCheckbox(b, "cat-laptops", "Laptops", true),
					filterCheckbox(b, "cat-monitors", "Monitors", true),
					filterCheckbox(b, "cat-servers", "Servers", true),
					filterCheckbox(b, "cat-network", "Network", true),
					filterCheckbox(b, "cat-printers", "Printers", false),
					filterCheckbox(b, "cat-other", "Other", false),
				),
			),
			
			// Status filter
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("block text-xs font-medium text-gray-500 mb-2"), "Status"),
				b.Div(mi.Class("space-y-2"),
					filterCheckbox(b, "status-active", "Active", true),
					filterCheckbox(b, "status-maintenance", "Maintenance", true),
					filterCheckbox(b, "status-retired", "Retired", false),
				),
			),
			
			// Department filter
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("block text-xs font-medium text-gray-500 mb-2"), "Department"),
				b.Select(
					mi.Class("w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"),
					mi.Name("department"),
					b.Option(mi.Value(""), "All Departments"),
					b.Option(mi.Value("engineering"), "Engineering"),
					b.Option(mi.Value("sales"), "Sales"),
					b.Option(mi.Value("marketing"), "Marketing"),
					b.Option(mi.Value("finance"), "Finance"),
					b.Option(mi.Value("it"), "IT Operations"),
				),
			),
			
			// Value range
			b.Div(mi.Class("mb-4"),
				b.Label(mi.Class("block text-xs font-medium text-gray-500 mb-2"), "Value Range"),
				b.Div(mi.Class("flex items-center gap-2"),
					b.Input(
						mi.Type("number"),
						mi.Placeholder("Min"),
						mi.Class("w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"),
					),
					b.Span(mi.Class("text-gray-400"), "-"),
					b.Input(
						mi.Type("number"),
						mi.Placeholder("Max"),
						mi.Class("w-full px-2 py-1.5 text-sm border border-gray-300 rounded-md"),
					),
				),
			),
			
			// Actions
			b.Div(mi.Class("flex gap-2 pt-4 border-t border-gray-200"),
				b.Button(
					mi.Class("flex-1 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"),
					mi.Type("button"),
					"Clear",
				),
				b.Button(
					mi.Class("flex-1 px-3 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"),
					mi.Type("button"),
					mi.HtmxPost("/assets/filter"),
					mi.HtmxTarget("#asset-table"),
					"Apply",
				),
			),
		)
	}
}

func filterCheckbox(b *mi.Builder, id, label string, checked bool) mi.Node {
	attrs := []mi.Attribute{
		mi.Type("checkbox"),
		mi.ID(id),
		mi.Class("rounded border-gray-300 text-blue-600"),
	}
	if checked {
		attrs = append(attrs, mi.Checked())
	}
	return b.Label(mi.Class("flex items-center gap-2 text-sm text-gray-600"),
		b.Input(attrs...),
		label,
	)
}

// Pagination
func Pagination(current, total int) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("flex items-center justify-between px-4 py-3 border-t border-gray-200"),
			b.Div(mi.Class("text-sm text-gray-500"),
				fmt.Sprintf("Showing 1 to %d of %d results", min(10, total), total),
			),
			b.Div(mi.Class("flex items-center gap-1"),
				b.Button(
					mi.Class("px-3 py-1.5 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"),
					mi.Type("button"),
					mi.Disabled(),
					"Previous",
				),
				pageButton(b, 1, true),
				pageButton(b, 2, false),
				pageButton(b, 3, false),
				b.Span(mi.Class("px-2 text-gray-400"), "..."),
				pageButton(b, 10, false),
				b.Button(
					mi.Class("px-3 py-1.5 text-sm border border-gray-300 rounded-md hover:bg-gray-50"),
					mi.Type("button"),
					"Next",
				),
			),
		)
	}
}

func pageButton(b *mi.Builder, page int, active bool) mi.Node {
	class := "px-3 py-1.5 text-sm rounded-md"
	if active {
		class += " bg-blue-600 text-white"
	} else {
		class += " border border-gray-300 hover:bg-gray-50"
	}
	return b.Button(
		mi.Class(class),
		mi.Type("button"),
		fmt.Sprintf("%d", page),
	)
}

// Modal for asset details
func Modal(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(
			mi.Class("fixed inset-0 z-50 overflow-y-auto"),
			mi.ID("modal"),
			// Backdrop
			b.Div(mi.Class("fixed inset-0 bg-black bg-opacity-50 transition-opacity")),
			// Modal panel
			b.Div(mi.Class("flex min-h-full items-center justify-center p-4"),
				b.Div(mi.Class("relative bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto"),
					// Header
					b.Div(mi.Class("flex items-center justify-between px-6 py-4 border-b border-gray-200"),
						b.H2(mi.Class("text-lg font-semibold text-gray-900"), title),
						b.Button(
							mi.Class("p-1 text-gray-400 hover:text-gray-600"),
							mi.Type("button"),
							Icon("close")(b),
						),
					),
					// Content
					b.Div(mi.Class("px-6 py-4"), mi.ID("modal-content"),
						content(b),
					),
				),
			),
		)
	}
}

// Dashboard content
func DashboardContent() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(
			// Stats row
			b.Div(mi.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6"),
				StatCard("Total Assets", "1,234", "+12% from last month", true, "assets")(b),
				StatCard("Active Assets", "1,156", "+8% from last month", true, "check")(b),
				StatCard("Under Maintenance", "45", "-5% from last month", true, "maintenance")(b),
				StatCard("Total Value", "$2.4M", "+15% from last month", true, "dashboard")(b),
			),
			
			// Main content grid
			b.Div(mi.Class("grid grid-cols-1 lg:grid-cols-3 gap-6"),
				// Assets by category chart placeholder
				b.Div(mi.Class("lg:col-span-2"),
					Card("Assets by Category", func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("h-64 flex items-center justify-center text-gray-400"),
							"[Chart: Asset distribution by category]",
						)
					}, nil)(b),
				),
				// Recent activity
				Card("Recent Activity", func(b *mi.Builder) mi.Node {
					activities := []struct{ Icon, Text, Time string }{
						{"add", "MacBook Pro added", "2 min ago"},
						{"edit", "Server SRV-001 updated", "15 min ago"},
						{"maintenance", "Laptop LAP-003 sent for repair", "1 hour ago"},
						{"check", "Audit completed for Floor 2", "3 hours ago"},
						{"delete", "Old printer retired", "Yesterday"},
					}
					nodes := make([]mi.Node, len(activities))
					for i, act := range activities {
						nodes[i] = b.Div(mi.Class("flex items-start gap-3 py-2"),
							b.Span(mi.Class("text-gray-400"), Icon(act.Icon)(b)),
							b.Div(
								b.P(mi.Class("text-sm text-gray-900"), act.Text),
								b.P(mi.Class("text-xs text-gray-500"), act.Time),
							),
						)
					}
					return b.Div(mi.Class("divide-y divide-gray-100"), mi.NewFragment(nodes...))
				}, nil)(b),
			),
		)
	}
}

// Main assets page content
func AssetsPageContent() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("flex"),
			// Filter sidebar
			FilterSidebar()(b),
			// Main content
			b.Div(mi.Class("flex-1"),
				// Toolbar
				b.Div(mi.Class("flex items-center justify-between px-4 py-3 bg-white border-b border-gray-200"),
					b.Div(mi.Class("flex items-center gap-2"),
						b.Button(
							mi.Class("inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"),
							mi.Type("button"),
							mi.HtmxGet("/assets/new"),
							mi.HtmxTarget("#modal-content"),
							Icon("add")(b),
							"Add Asset",
						),
						b.Button(
							mi.Class("inline-flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"),
							mi.Type("button"),
							Icon("import")(b),
							"Import",
						),
						b.Button(
							mi.Class("inline-flex items-center gap-2 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"),
							mi.Type("button"),
							Icon("export")(b),
							"Export",
						),
					),
					b.Div(mi.Class("flex items-center gap-2"),
						b.Span(mi.Class("text-sm text-gray-500"), "8 assets selected"),
						b.Button(
							mi.Class("px-3 py-1.5 text-sm text-red-600 hover:text-red-700"),
							mi.Type("button"),
							"Delete Selected",
						),
					),
				),
				// Table
				b.Div(mi.ID("asset-table"),
					AssetTable(sampleAssets)(b),
				),
				// Pagination
				Pagination(1, 87)(b),
			),
		)
	}
}

// =====================================================
// MAIN PAGE
// =====================================================

func AssetManagementUI() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			mi.Raw("<!DOCTYPE html>"),
			b.Html(mi.Lang("en"),
				b.Head(
					b.Title("AssetTrack - Asset Management System"),
					b.Meta(mi.Charset("UTF-8")),
					b.Meta(mi.Name("viewport"), mi.Content("width=device-width, initial-scale=1")),
					b.Script(mi.Src("https://unpkg.com/htmx.org@1.9.10")),
					b.Script(mi.Src("https://cdn.tailwindcss.com")),
					b.Style(mi.Raw(customCSS)),
				),
				b.Body(mi.Class("bg-gray-100"),
					b.Div(mi.Class("flex"),
						// Sidebar
						Sidebar()(b),
						// Main area
						b.Div(mi.Class("flex-1 flex flex-col min-h-screen"),
							// Header
							Header()(b),
							// Page content
							b.Main(mi.Class("flex-1 p-6"),
								// Page title
								b.Div(mi.Class("mb-6"),
									b.H2(mi.Class("text-2xl font-bold text-gray-900"), "Asset Inventory"),
									b.P(mi.Class("text-sm text-gray-500"), "Manage and track all company assets"),
								),
								// Content area
								b.Div(mi.Class("bg-white rounded-lg shadow-sm border border-gray-200"),
									AssetsPageContent()(b),
								),
							),
						),
					),
					// Hidden search results container
					b.Div(mi.ID("search-results"), mi.Class("hidden")),
					// Modal container (hidden by default)
					b.Div(mi.ID("modal-container"), mi.Class("hidden")),
				),
			),
		)
	}
}

const customCSS = `
.icon { font-style: normal; }
.sidebar-item { transition: all 0.15s ease; }

/* Tab panels */
.tab-panel { display: block; }
.tab-panel.hidden { display: none; }

/* Custom scrollbar */
::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: #f1f1f1; }
::-webkit-scrollbar-thumb { background: #c1c1c1; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #a1a1a1; }

/* HTMX loading states */
.htmx-request { opacity: 0.7; pointer-events: none; }
.htmx-request::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    width: 20px;
    height: 20px;
    margin: -10px 0 0 -10px;
    border: 2px solid #e5e7eb;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
`

const tabScript = `
function switchTab(tabId) {
    // Hide all panels
    document.querySelectorAll('.tab-panel').forEach(panel => {
        panel.classList.add('hidden');
    });
    // Show selected panel
    const panel = document.getElementById('panel-' + tabId);
    if (panel) {
        panel.classList.remove('hidden');
    }
    // Update tab button styles
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('text-blue-600', 'border-blue-600', 'active');
        btn.classList.add('text-gray-500', 'border-transparent');
    });
    const activeBtn = document.querySelector('.tab-button[data-tab="' + tabId + '"]');
    if (activeBtn) {
        activeBtn.classList.remove('text-gray-500', 'border-transparent');
        activeBtn.classList.add('text-blue-600', 'border-blue-600', 'active');
    }
}
`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Generate static HTML
	f, err := os.Create("asset-management.html")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	if err := mi.Render(AssetManagementUI(), f); err != nil {
		fmt.Println("Error rendering:", err)
		return
	}
	fmt.Println("Generated asset-management.html")
	
	// Also generate the detail form modal for review
	f2, err := os.Create("asset-detail-modal.html")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f2.Close()
	
	modalContent := Modal("Edit Asset: MacBook Pro 16\"", func(b *mi.Builder) mi.Node {
		return AssetDetailForm(sampleAssets[0])(b)
	})
	
	fullPage := func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			mi.Raw("<!DOCTYPE html>"),
			b.Html(mi.Lang("en"),
				b.Head(
					b.Title("Asset Detail Modal"),
					b.Meta(mi.Charset("UTF-8")),
					b.Meta(mi.Name("viewport"), mi.Content("width=device-width, initial-scale=1")),
					b.Script(mi.Src("https://cdn.tailwindcss.com")),
					b.Style(mi.Raw(customCSS)),
				),
				b.Body(mi.Class("bg-gray-500"),
					modalContent(b),
					b.Script(mi.Raw(tabScript)),
				),
			),
		)
	}
	
	if err := mi.Render(fullPage, f2); err != nil {
		fmt.Println("Error rendering modal:", err)
		return
	}
	fmt.Println("Generated asset-detail-modal.html")
	
	fmt.Printf("\nGenerated at %s\n", time.Now().Format("15:04:05"))
}
