// Package tailwind provides a Tailwind CSS theme implementation for mintyui
package tailwind

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
)

// TailwindTheme implements the Theme interface using Tailwind CSS classes
type TailwindTheme struct {
	name    string
	version string
}

// NewTailwindTheme creates a new Tailwind theme
func NewTailwindTheme() mui.Theme {
	return &TailwindTheme{
		name:    "Tailwind",
		version: "1.0.0",
	}
}

// GetName returns the theme name
func (t *TailwindTheme) GetName() string {
	return t.name
}

// GetVersion returns the theme version
func (t *TailwindTheme) GetVersion() string {
	return t.version
}

// =====================================================
// BASIC COMPONENTS
// =====================================================

// Button creates a Tailwind button
func (t *TailwindTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getButtonClass(variant)
		allAttrs := append([]mi.Attribute{mi.Class(class), mi.Type("button")}, attrs...)
		args := make([]interface{}, len(allAttrs)+1)
		for i, attr := range allAttrs {
			args[i] = attr
		}
		args[len(allAttrs)] = text
		return b.Button(args...)
	}
}

// Card creates a Tailwind card
func (t *TailwindTheme) Card(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("bg-white rounded-lg border border-gray-200 shadow-sm mb-4"),
			mi.NewFragment(
				func(b *mi.Builder) mi.Node {
					if title != "" {
						return b.Div(mi.Class("px-6 py-4 border-b border-gray-200"),
							b.H3(mi.Class("text-lg font-medium text-gray-900"), title),
						)
					}
					return mi.NewFragment()
				}(b),
				b.Div(mi.Class("p-6"),
					content(b),
				),
			),
		)
	}
}

// Badge creates a Tailwind badge
func (t *TailwindTheme) Badge(text, variant string) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getBadgeClass(variant)
		return b.Span(mi.Class(class), text)
	}
}

// =====================================================
// FORM COMPONENTS
// =====================================================

// FormInput creates a Tailwind form input with label
func (t *TailwindTheme) FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "input_" + name
		inputAttrs := append([]mi.Attribute{
			mi.Class("mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"),
			mi.ID(id),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		
		return b.Div(mi.Class("mb-4"),
			t.FormLabel(label, id)(b),
			b.Input(inputAttrs...),
		)
	}
}

// FormSelect creates a Tailwind select dropdown with label
func (t *TailwindTheme) FormSelect(label, name string, options []mui.SelectOption) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "select_" + name
		
		optionNodes := make([]mi.Node, len(options))
		for i, option := range options {
			var optAttrs []mi.Attribute
			optAttrs = append(optAttrs, mi.Value(option.Value))
			if option.Selected {
				optAttrs = append(optAttrs, mi.Selected())
			}
			if option.Disabled {
				optAttrs = append(optAttrs, mi.Disabled())
			}
			optionNodes[i] = b.Option(optAttrs, option.Text)
		}
		
		return b.Div(mi.Class("mb-4"),
			t.FormLabel(label, id)(b),
			b.Select(
				mi.Class("mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"),
				mi.ID(id), mi.Name(name),
				mi.NewFragment(optionNodes...),
			),
		)
	}
}

// FormTextarea creates a Tailwind textarea with label
func (t *TailwindTheme) FormTextarea(label, name string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "textarea_" + name
		textareaAttrs := append([]mi.Attribute{
			mi.Class("mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"),
			mi.ID(id),
			mi.Name(name),
			mi.Rows(3),
		}, attrs...)
		
		// Convert to []interface{} for Textarea
		args := make([]interface{}, len(textareaAttrs))
		for i, attr := range textareaAttrs {
			args[i] = attr
		}
		
		return b.Div(mi.Class("mb-4"),
			t.FormLabel(label, id)(b),
			b.Textarea(args...),
		)
	}
}

// FormLabel creates a Tailwind form label
func (t *TailwindTheme) FormLabel(text, forField string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Label(
			mi.Class("block text-sm font-medium text-gray-700"),
			mi.DataAttr("for", forField),
			text,
		)
	}
}

// Input creates a standalone Tailwind input
func (t *TailwindTheme) Input(name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		inputAttrs := append([]mi.Attribute{
			mi.Class("block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		return b.Input(inputAttrs...)
	}
}

// =====================================================
// LAYOUT COMPONENTS
// =====================================================

// Container creates a Tailwind container
func (t *TailwindTheme) Container(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
			content(b),
		)
	}
}

// Grid creates a Tailwind grid layout
func (t *TailwindTheme) Grid(columns int, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		gridClass := fmt.Sprintf("grid grid-cols-1 md:grid-cols-%d gap-6", columns)
		return b.Div(mi.Class(gridClass),
			content(b),
		)
	}
}

// Sidebar creates a Tailwind sidebar
func (t *TailwindTheme) Sidebar(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(
			mi.Class("bg-gray-50 border-r border-gray-200 px-4 py-6"),
			mi.Style("min-height: 100vh;"),
			content(b),
		)
	}
}

// =====================================================
// NAVIGATION COMPONENTS
// =====================================================

// Nav creates a Tailwind navigation menu
func (t *TailwindTheme) Nav(items []mui.NavItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		navItems := make([]mi.Node, len(items))
		for i, item := range items {
			class := "block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100"
			if item.Active {
				class = "block px-3 py-2 rounded-md text-base font-medium bg-blue-100 text-blue-700"
			}
			
			navItems[i] = b.A(mi.Class(class), mi.Href(item.URL),
				func(b *mi.Builder) mi.Node {
					if item.Icon != "" {
						return mi.NewFragment(
							b.Span(mi.Class("mr-3"), item.Icon),
							mi.Txt(item.Text),
						)
					}
					return mi.Txt(item.Text)
				}(b),
			)
		}
		
		return b.Nav(mi.Class("space-y-1"),
			mi.NewFragment(navItems...),
		)
	}
}

// Breadcrumbs creates Tailwind breadcrumbs
func (t *TailwindTheme) Breadcrumbs(items []mui.BreadcrumbItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		breadcrumbItems := make([]mi.Node, len(items))
		for i, item := range items {
			if item.Last {
				breadcrumbItems[i] = b.Span(
					mi.Class("text-gray-500"),
					item.Text,
				)
			} else {
				breadcrumbItems[i] = mi.NewFragment(
					b.A(
						mi.Class("text-blue-600 hover:text-blue-800"),
						mi.Href(item.URL),
						item.Text,
					),
					b.Span(mi.Class("mx-2 text-gray-400"), "/"),
				)
			}
		}
		
		return b.Nav(mi.Class("flex"), mi.AriaLabel("breadcrumb"),
			b.Ol(mi.Class("inline-flex items-center space-x-1 md:space-x-3"),
				mi.NewFragment(breadcrumbItems...),
			),
		)
	}
}

// Pagination creates Tailwind pagination
func (t *TailwindTheme) Pagination(currentPage, totalPages int, baseURL string) mi.H {
	return func(b *mi.Builder) mi.Node {
		pageItems := make([]mi.Node, 0)
		
		// Previous button
		prevClass := "relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
		if currentPage <= 1 {
			prevClass += " cursor-not-allowed opacity-50"
		}
		pageItems = append(pageItems, 
			b.A(mi.Class(prevClass), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage-1)),
				"Previous"),
		)
		
		// Page numbers
		start := maxInt(1, currentPage-2)
		end := minInt(totalPages, currentPage+2)
		
		for i := start; i <= end; i++ {
			pageClass := "relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
			if i == currentPage {
				pageClass = "relative inline-flex items-center px-4 py-2 border border-blue-500 bg-blue-50 text-sm font-medium text-blue-600"
			}
			
			pageItems = append(pageItems, 
				b.A(mi.Class(pageClass), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, i)),
					fmt.Sprintf("%d", i)),
			)
		}
		
		// Next button
		nextClass := "relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
		if currentPage >= totalPages {
			nextClass += " cursor-not-allowed opacity-50"
		}
		pageItems = append(pageItems, 
			b.A(mi.Class(nextClass), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage+1)),
				"Next"),
		)
		
		return b.Nav(mi.Class("flex justify-center"), mi.AriaLabel("pagination"),
			b.Div(mi.Class("relative z-0 inline-flex rounded-md shadow-sm -space-x-px"),
				mi.NewFragment(pageItems...),
			),
		)
	}
}

// =====================================================
// DATA COMPONENTS
// =====================================================

// Table creates a Tailwind table
func (t *TailwindTheme) Table(headers []string, rows [][]string) mi.H {
	return func(b *mi.Builder) mi.Node {
		// Create header row
		headerCells := make([]mi.Node, len(headers))
		for i, header := range headers {
			headerCells[i] = b.Th(
				mi.Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"),
				header,
			)
		}
		
		// Create data rows
		dataRows := make([]mi.Node, len(rows))
		for i, row := range rows {
			cells := make([]mi.Node, len(row))
			for j, cell := range row {
				cells[j] = b.Td(
					mi.Class("px-6 py-4 whitespace-nowrap text-sm text-gray-900"),
					mi.RawHTML(cell),
				)
			}
			rowClass := "bg-white"
			if i%2 == 1 {
				rowClass = "bg-gray-50"
			}
			dataRows[i] = b.Tr(mi.Class(rowClass), mi.NewFragment(cells...))
		}
		
		return b.Div(mi.Class("overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg"),
			b.Table(mi.Class("min-w-full divide-y divide-gray-300"),
				b.Thead(mi.Class("bg-gray-50"),
					b.Tr(mi.NewFragment(headerCells...)),
				),
				b.Tbody(mi.Class("divide-y divide-gray-200"),
					mi.NewFragment(dataRows...),
				),
			),
		)
	}
}

// List creates a Tailwind list
func (t *TailwindTheme) List(items []string, ordered bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		listItems := make([]mi.Node, len(items))
		for i, item := range items {
			listItems[i] = b.Li(
				mi.Class("py-2 px-4 border-b border-gray-200 last:border-b-0"),
				item,
			)
		}
		
		listClass := "bg-white border border-gray-200 rounded-lg divide-y divide-gray-200"
		
		if ordered {
			return b.Ol(mi.Class(listClass+" list-decimal list-inside"),
				mi.NewFragment(listItems...))
		}
		
		return b.Ul(mi.Class(listClass),
			mi.NewFragment(listItems...))
	}
}

// =====================================================
// UTILITY METHODS
// =====================================================

// PrimaryButton creates a primary Tailwind button
func (t *TailwindTheme) PrimaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "primary", attrs...)
}

// SecondaryButton creates a secondary Tailwind button
func (t *TailwindTheme) SecondaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "secondary", attrs...)
}

// DangerButton creates a danger Tailwind button
func (t *TailwindTheme) DangerButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "danger", attrs...)
}

// =====================================================
// HELPER METHODS
// =====================================================

// getButtonClass returns Tailwind button class for variant
func (t *TailwindTheme) getButtonClass(variant string) string {
	baseClass := "inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2"
	
	switch variant {
	case "primary":
		return baseClass + " border-transparent text-white bg-blue-600 hover:bg-blue-700 focus:ring-blue-500"
	case "secondary":
		return baseClass + " border-gray-300 text-gray-700 bg-white hover:bg-gray-50 focus:ring-blue-500"
	case "success":
		return baseClass + " border-transparent text-white bg-green-600 hover:bg-green-700 focus:ring-green-500"
	case "warning":
		return baseClass + " border-transparent text-black bg-yellow-400 hover:bg-yellow-500 focus:ring-yellow-500"
	case "danger", "error":
		return baseClass + " border-transparent text-white bg-red-600 hover:bg-red-700 focus:ring-red-500"
	case "info":
		return baseClass + " border-transparent text-white bg-cyan-600 hover:bg-cyan-700 focus:ring-cyan-500"
	case "light":
		return baseClass + " border-gray-300 text-gray-700 bg-gray-100 hover:bg-gray-200 focus:ring-gray-500"
	case "dark":
		return baseClass + " border-transparent text-white bg-gray-800 hover:bg-gray-900 focus:ring-gray-700"
	case "link":
		return "inline-flex items-center px-2 py-1 text-sm font-medium text-blue-600 hover:text-blue-800"
	case "view":
		return baseClass + " border-blue-300 text-blue-700 bg-blue-50 hover:bg-blue-100 focus:ring-blue-500 text-xs px-3 py-1"
	case "payment":
		return baseClass + " border-transparent text-white bg-green-600 hover:bg-green-700 focus:ring-green-500"
	default:
		return baseClass + " border-gray-300 text-gray-700 bg-white hover:bg-gray-50 focus:ring-blue-500"
	}
}

// getBadgeClass returns Tailwind badge class for variant
func (t *TailwindTheme) getBadgeClass(variant string) string {
	baseClass := "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
	
	switch variant {
	case "primary":
		return baseClass + " bg-blue-100 text-blue-800"
	case "secondary":
		return baseClass + " bg-gray-100 text-gray-800"
	case "success":
		return baseClass + " bg-green-100 text-green-800"
	case "warning":
		return baseClass + " bg-yellow-100 text-yellow-800"
	case "danger", "error":
		return baseClass + " bg-red-100 text-red-800"
	case "info":
		return baseClass + " bg-cyan-100 text-cyan-800"
	case "light":
		return baseClass + " bg-gray-50 text-gray-600"
	case "dark":
		return baseClass + " bg-gray-800 text-gray-100"
	default:
		return baseClass + " bg-gray-100 text-gray-800"
	}
}

// Helper functions
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// =====================================================
// CSS AND SCRIPTS
// =====================================================

// CDNLinks returns Tailwind CSS CDN links for inclusion in HTML head
func CDNLinks() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			// Tailwind CSS
			b.Link(
				mi.Href("https://cdn.tailwindcss.com"),
				mi.Rel("stylesheet"),
			),
			// Optional: Add custom Tailwind config
			b.Script(`
				tailwind.config = {
					theme: {
						extend: {
							colors: {
								primary: {
									50: '#eff6ff',
									500: '#3b82f6',
									600: '#2563eb',
									700: '#1d4ed8',
								}
							}
						}
					}
				}
			`),
		)
	}
}

// TailwindDocument creates a complete HTML document with Tailwind styling
func TailwindDocument(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.Document(title,
			[]mi.Node{
				CDNLinks()(b),
			},
			b.Body(mi.Class("bg-gray-50 min-h-screen"),
				b.Div(mi.Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8"),
					content(b),
				),
			),
		)(b)
	}
}
