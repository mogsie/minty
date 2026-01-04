// Package bootstrap provides a Bootstrap 5 theme implementation for mintyui
package bootstrap

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
)

// BootstrapTheme implements the Theme interface using Bootstrap 5 classes
type BootstrapTheme struct {
	name    string
	version string
}

// NewBootstrapTheme creates a new Bootstrap theme
func NewBootstrapTheme() mui.Theme {
	return &BootstrapTheme{
		name:    "Bootstrap",
		version: "1.0.0",
	}
}

// GetName returns the theme name
func (t *BootstrapTheme) GetName() string {
	return t.name
}

// GetVersion returns the theme version
func (t *BootstrapTheme) GetVersion() string {
	return t.version
}

// =====================================================
// BASIC COMPONENTS
// =====================================================

// Button creates a Bootstrap button
func (t *BootstrapTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
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

// Card creates a Bootstrap card
func (t *BootstrapTheme) Card(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("card mb-3"),
			mi.NewFragment(
				func(b *mi.Builder) mi.Node {
					if title != "" {
						return b.Div(mi.Class("card-header"),
							b.H5(mi.Class("card-title mb-0"), title),
						)
					}
					return mi.NewFragment()
				}(b),
				b.Div(mi.Class("card-body"),
					content(b),
				),
			),
		)
	}
}

// Badge creates a Bootstrap badge
func (t *BootstrapTheme) Badge(text, variant string) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getBadgeClass(variant)
		return b.Span(mi.Class(class), text)
	}
}

// =====================================================
// FORM COMPONENTS
// =====================================================

// FormInput creates a Bootstrap form input with label
func (t *BootstrapTheme) FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "input_" + name
		inputAttrs := append([]mi.Attribute{
			mi.Class("form-control"),
			mi.ID(id),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		
		return b.Div(mi.Class("mb-3"),
			t.FormLabel(label, id)(b),
			b.Input(inputAttrs...),
		)
	}
}

// FormSelect creates a Bootstrap select dropdown with label
func (t *BootstrapTheme) FormSelect(label, name string, options []mui.SelectOption) mi.H {
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
		
		return b.Div(mi.Class("mb-3"),
			t.FormLabel(label, id)(b),
			b.Select(mi.Class("form-select"), mi.ID(id), mi.Name(name),
				mi.NewFragment(optionNodes...),
			),
		)
	}
}

// FormTextarea creates a Bootstrap textarea with label
func (t *BootstrapTheme) FormTextarea(label, name string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "textarea_" + name
		textareaAttrs := append([]mi.Attribute{
			mi.Class("form-control"),
			mi.ID(id),
			mi.Name(name),
		}, attrs...)
		
		// Convert to []interface{} for Textarea
		args := make([]interface{}, len(textareaAttrs))
		for i, attr := range textareaAttrs {
			args[i] = attr
		}
		
		return b.Div(mi.Class("mb-3"),
			t.FormLabel(label, id)(b),
			b.Textarea(args...),
		)
	}
}

// FormLabel creates a Bootstrap form label
func (t *BootstrapTheme) FormLabel(text, forField string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Label(mi.Class("form-label"), mi.DataAttr("for", forField), text)
	}
}

// Input creates a standalone Bootstrap input
func (t *BootstrapTheme) Input(name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		inputAttrs := append([]mi.Attribute{
			mi.Class("form-control"),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		return b.Input(inputAttrs...)
	}
}

// =====================================================
// LAYOUT COMPONENTS
// =====================================================

// Container creates a Bootstrap container
func (t *BootstrapTheme) Container(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("container-fluid"),
			content(b),
		)
	}
}

// Grid creates a Bootstrap grid layout
func (t *BootstrapTheme) Grid(columns int, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		colClass := fmt.Sprintf("col-md-%d", 12/columns)
		return b.Div(mi.Class("row"),
			b.Div(mi.Class(colClass),
				content(b),
			),
		)
	}
}

// Sidebar creates a Bootstrap sidebar
func (t *BootstrapTheme) Sidebar(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("bg-light border-end"),
			mi.Style("min-height: 100vh;"),
			content(b),
		)
	}
}

// =====================================================
// NAVIGATION COMPONENTS
// =====================================================

// Nav creates a Bootstrap navigation menu
func (t *BootstrapTheme) Nav(items []mui.NavItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		navItems := make([]mi.Node, len(items))
		for i, item := range items {
			class := "nav-link"
			if item.Active {
				class += " active"
			}
			
			navItems[i] = b.Li(mi.Class("nav-item"),
				b.A(mi.Class(class), mi.Href(item.URL),
					func(b *mi.Builder) mi.Node {
						if item.Icon != "" {
							return mi.NewFragment(
								b.Span(mi.Class("me-2"), item.Icon),
								mi.Txt(item.Text),
							)
						}
						return mi.Txt(item.Text)
					}(b),
				),
			)
		}
		
		return b.Ul(mi.Class("nav nav-pills flex-column"),
			mi.NewFragment(navItems...),
		)
	}
}

// Breadcrumbs creates Bootstrap breadcrumbs
func (t *BootstrapTheme) Breadcrumbs(items []mui.BreadcrumbItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		breadcrumbItems := make([]mi.Node, len(items))
		for i, item := range items {
			if item.Last {
				breadcrumbItems[i] = b.Li(mi.Class("breadcrumb-item active"),
					mi.AriaLabel("current"), item.Text)
			} else {
				breadcrumbItems[i] = b.Li(mi.Class("breadcrumb-item"),
					b.A(mi.Href(item.URL), item.Text))
			}
		}
		
		return b.Nav(mi.AriaLabel("breadcrumb"),
			b.Ol(mi.Class("breadcrumb"),
				mi.NewFragment(breadcrumbItems...),
			),
		)
	}
}

// Pagination creates Bootstrap pagination
func (t *BootstrapTheme) Pagination(currentPage, totalPages int, baseURL string) mi.H {
	return func(b *mi.Builder) mi.Node {
		pageItems := make([]mi.Node, 0)
		
		// Previous button
		prevClass := "page-item"
		if currentPage <= 1 {
			prevClass += " disabled"
		}
		pageItems = append(pageItems, b.Li(mi.Class(prevClass),
			b.A(mi.Class("page-link"), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage-1)),
				"Previous"),
		))
		
		// Page numbers
		start := maxInt(1, currentPage-2)
		end := minInt(totalPages, currentPage+2)
		
		for i := start; i <= end; i++ {
			pageClass := "page-item"
			if i == currentPage {
				pageClass += " active"
			}
			
			pageItems = append(pageItems, b.Li(mi.Class(pageClass),
				b.A(mi.Class("page-link"), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, i)),
					fmt.Sprintf("%d", i)),
			))
		}
		
		// Next button
		nextClass := "page-item"
		if currentPage >= totalPages {
			nextClass += " disabled"
		}
		pageItems = append(pageItems, b.Li(mi.Class(nextClass),
			b.A(mi.Class("page-link"), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage+1)),
				"Next"),
		))
		
		return b.Nav(mi.AriaLabel("Page navigation"),
			b.Ul(mi.Class("pagination justify-content-center"),
				mi.NewFragment(pageItems...),
			),
		)
	}
}

// =====================================================
// DATA COMPONENTS
// =====================================================

// Table creates a Bootstrap table
func (t *BootstrapTheme) Table(headers []string, rows [][]string) mi.H {
	return func(b *mi.Builder) mi.Node {
		// Create header row
		headerCells := make([]mi.Node, len(headers))
		for i, header := range headers {
			headerCells[i] = b.Th(mi.Scope("col"), header)
		}
		
		// Create data rows
		dataRows := make([]mi.Node, len(rows))
		for i, row := range rows {
			cells := make([]mi.Node, len(row))
			for j, cell := range row {
				cells[j] = b.Td(mi.RawHTML(cell))
			}
			dataRows[i] = b.Tr(mi.NewFragment(cells...))
		}
		
		return b.Div(mi.Class("table-responsive"),
			b.Table(mi.Class("table table-striped table-hover"),
				b.Thead(mi.Class("table-dark"),
					b.Tr(mi.NewFragment(headerCells...)),
				),
				b.Tbody(mi.NewFragment(dataRows...)),
			),
		)
	}
}

// List creates a Bootstrap list
func (t *BootstrapTheme) List(items []string, ordered bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		listItems := make([]mi.Node, len(items))
		for i, item := range items {
			listItems[i] = b.Li(mi.Class("list-group-item"), item)
		}
		
		if ordered {
			return b.Ol(mi.Class("list-group list-group-numbered"),
				mi.NewFragment(listItems...))
		}
		
		return b.Ul(mi.Class("list-group"),
			mi.NewFragment(listItems...))
	}
}

// =====================================================
// UTILITY METHODS
// =====================================================

// PrimaryButton creates a primary Bootstrap button
func (t *BootstrapTheme) PrimaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "primary", attrs...)
}

// SecondaryButton creates a secondary Bootstrap button
func (t *BootstrapTheme) SecondaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "secondary", attrs...)
}

// DangerButton creates a danger Bootstrap button
func (t *BootstrapTheme) DangerButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "danger", attrs...)
}

// =====================================================
// HELPER METHODS
// =====================================================

// getButtonClass returns Bootstrap button class for variant
func (t *BootstrapTheme) getButtonClass(variant string) string {
	switch variant {
	case "primary":
		return "btn btn-primary"
	case "secondary":
		return "btn btn-secondary"
	case "success":
		return "btn btn-success"
	case "warning":
		return "btn btn-warning"
	case "danger", "error":
		return "btn btn-danger"
	case "info":
		return "btn btn-info"
	case "light":
		return "btn btn-light"
	case "dark":
		return "btn btn-dark"
	case "link":
		return "btn btn-link"
	case "outline-primary":
		return "btn btn-outline-primary"
	case "outline-secondary":
		return "btn btn-outline-secondary"
	case "view":
		return "btn btn-outline-primary btn-sm"
	case "payment":
		return "btn btn-success"
	default:
		return "btn btn-secondary"
	}
}

// getBadgeClass returns Bootstrap badge class for variant
func (t *BootstrapTheme) getBadgeClass(variant string) string {
	switch variant {
	case "primary":
		return "badge bg-primary"
	case "secondary":
		return "badge bg-secondary"
	case "success":
		return "badge bg-success"
	case "warning":
		return "badge bg-warning text-dark"
	case "danger", "error":
		return "badge bg-danger"
	case "info":
		return "badge bg-info text-dark"
	case "light":
		return "badge bg-light text-dark"
	case "dark":
		return "badge bg-dark"
	default:
		return "badge bg-secondary"
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

// CDNLinks returns Bootstrap 5 CDN links for inclusion in HTML head
func CDNLinks() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			// Bootstrap CSS
			b.Link(
				mi.Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"),
				mi.Rel("stylesheet"),
				mi.DataAttr("integrity", "sha384-9ndCyUa+IgAaWhp066j+EugYkAuULlhkTAP0O7D4C/ZIyIbOrN4ySHXKdZJh3jP6"),
				mi.DataAttr("crossorigin", "anonymous"),
			),
			// Bootstrap Icons
			b.Link(
				mi.Href("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css"),
				mi.Rel("stylesheet"),
			),
		)
	}
}

// CDNScripts returns Bootstrap 5 CDN scripts for inclusion before closing body tag
func CDNScripts() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Script(
			mi.Src("https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"),
			mi.DataAttr("integrity", "sha384-geWF76RCwLtnZ8qwWowPQNguL3RmwHVBC9FhGdlKrxdiJJigb/j/68SIy3Te4Bkz"),
			mi.DataAttr("crossorigin", "anonymous"),
		)
	}
}

// BootstrapDocument creates a complete HTML document with Bootstrap styling
func BootstrapDocument(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.Document(title,
			[]mi.Node{
				CDNLinks()(b),
			},
			b.Body(
				b.Div(mi.Class("container-fluid"),
					content(b),
				),
				CDNScripts()(b),
			),
		)(b)
	}
}
