// Package bulma provides a Bulma CSS theme implementation for mintyui
package bulma

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
)

// BulmaTheme implements the Theme interface using Bulma CSS classes
type BulmaTheme struct {
	name    string
	version string
}

// NewBulmaTheme creates a new Bulma theme
func NewBulmaTheme() mui.Theme {
	return &BulmaTheme{
		name:    "Bulma",
		version: "1.0.0",
	}
}

// GetName returns the theme name
func (t *BulmaTheme) GetName() string {
	return t.name
}

// GetVersion returns the theme version
func (t *BulmaTheme) GetVersion() string {
	return t.version
}

// =====================================================
// BASIC COMPONENTS
// =====================================================

// Button creates a Bulma button
func (t *BulmaTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
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

// Card creates a Bulma card
func (t *BulmaTheme) Card(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("card mb-5"),
			mi.NewFragment(
				func(b *mi.Builder) mi.Node {
					if title != "" {
						return b.Header(mi.Class("card-header"),
							b.P(mi.Class("card-header-title"), title),
						)
					}
					return mi.NewFragment()
				}(b),
				b.Div(mi.Class("card-content"),
					b.Div(mi.Class("content"),
						content(b),
					),
				),
			),
		)
	}
}

// Badge creates a Bulma tag (badge equivalent)
func (t *BulmaTheme) Badge(text, variant string) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getBadgeClass(variant)
		return b.Span(mi.Class(class), text)
	}
}

// =====================================================
// FORM COMPONENTS
// =====================================================

// FormInput creates a Bulma form input with label
func (t *BulmaTheme) FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "input_" + name
		inputAttrs := append([]mi.Attribute{
			mi.Class("input"),
			mi.ID(id),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		
		return b.Div(mi.Class("field"),
			t.FormLabel(label, id)(b),
			b.Div(mi.Class("control"),
				b.Input(inputAttrs...),
			),
		)
	}
}

// FormSelect creates a Bulma select dropdown with label
func (t *BulmaTheme) FormSelect(label, name string, options []mui.SelectOption) mi.H {
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
		
		return b.Div(mi.Class("field"),
			t.FormLabel(label, id)(b),
			b.Div(mi.Class("control"),
				b.Div(mi.Class("select"),
					b.Select(mi.ID(id), mi.Name(name),
						mi.NewFragment(optionNodes...),
					),
				),
			),
		)
	}
}

// FormTextarea creates a Bulma textarea with label
func (t *BulmaTheme) FormTextarea(label, name string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "textarea_" + name
		textareaAttrs := append([]mi.Attribute{
			mi.Class("textarea"),
			mi.ID(id),
			mi.Name(name),
		}, attrs...)
		
		// Convert to []interface{} for Textarea
		args := make([]interface{}, len(textareaAttrs))
		for i, attr := range textareaAttrs {
			args[i] = attr
		}
		
		return b.Div(mi.Class("field"),
			t.FormLabel(label, id)(b),
			b.Div(mi.Class("control"),
				b.Textarea(args...),
			),
		)
	}
}

// FormLabel creates a Bulma form label
func (t *BulmaTheme) FormLabel(text, forField string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Label(mi.Class("label"), mi.DataAttr("for", forField), text)
	}
}

// Input creates a standalone Bulma input
func (t *BulmaTheme) Input(name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		inputAttrs := append([]mi.Attribute{
			mi.Class("input"),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		return b.Input(inputAttrs...)
	}
}

// =====================================================
// LAYOUT COMPONENTS
// =====================================================

// Container creates a Bulma container
func (t *BulmaTheme) Container(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("container is-fluid"),
			content(b),
		)
	}
}

// Grid creates a Bulma columns layout
func (t *BulmaTheme) Grid(columns int, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		colClass := fmt.Sprintf("column is-%d", 12/columns)
		return b.Div(mi.Class("columns"),
			b.Div(mi.Class(colClass),
				content(b),
			),
		)
	}
}

// Sidebar creates a Bulma sidebar
func (t *BulmaTheme) Sidebar(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Aside(mi.Class("menu has-background-light"),
			mi.Style("min-height: 100vh; padding: 1rem;"),
			content(b),
		)
	}
}

// =====================================================
// NAVIGATION COMPONENTS
// =====================================================

// Nav creates a Bulma menu navigation
func (t *BulmaTheme) Nav(items []mui.NavItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		navItems := make([]mi.Node, len(items))
		for i, item := range items {
			class := "menu-item"
			if item.Active {
				class += " is-active"
			}
			
			navItems[i] = b.Li(
				b.A(mi.Class(class), mi.Href(item.URL),
					func(b *mi.Builder) mi.Node {
						if item.Icon != "" {
							return mi.NewFragment(
								b.Span(mi.Class("icon mr-2"), item.Icon),
								mi.Txt(item.Text),
							)
						}
						return mi.Txt(item.Text)
					}(b),
				),
			)
		}
		
		return b.Aside(mi.Class("menu"),
			b.Ul(mi.Class("menu-list"),
				mi.NewFragment(navItems...),
			),
		)
	}
}

// Breadcrumbs creates Bulma breadcrumbs
func (t *BulmaTheme) Breadcrumbs(items []mui.BreadcrumbItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		breadcrumbItems := make([]mi.Node, len(items))
		for i, item := range items {
			if item.Last {
				breadcrumbItems[i] = b.Li(mi.Class("is-active"),
					b.A(mi.AriaLabel("current"), item.Text))
			} else {
				breadcrumbItems[i] = b.Li(
					b.A(mi.Href(item.URL), item.Text))
			}
		}
		
		return b.Nav(mi.Class("breadcrumb"), mi.AriaLabel("breadcrumbs"),
			b.Ul(mi.NewFragment(breadcrumbItems...)),
		)
	}
}

// Pagination creates Bulma pagination
func (t *BulmaTheme) Pagination(currentPage, totalPages int, baseURL string) mi.H {
	return func(b *mi.Builder) mi.Node {
		pageItems := make([]mi.Node, 0)
		
		// Previous button
		prevClass := "pagination-previous"
		var prevAttrs []mi.Attribute
		if currentPage <= 1 {
			prevAttrs = append(prevAttrs, mi.Disabled())
		} else {
			prevAttrs = append(prevAttrs, mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage-1)))
		}
		prevAttrs = append(prevAttrs, mi.Class(prevClass))
		
		pageItems = append(pageItems, b.A(prevAttrs, "Previous"))
		
		// Next button
		nextClass := "pagination-next"
		var nextAttrs []mi.Attribute
		if currentPage >= totalPages {
			nextAttrs = append(nextAttrs, mi.Disabled())
		} else {
			nextAttrs = append(nextAttrs, mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage+1)))
		}
		nextAttrs = append(nextAttrs, mi.Class(nextClass))
		
		pageItems = append(pageItems, b.A(nextAttrs, "Next"))
		
		// Page numbers
		paginationList := make([]mi.Node, 0)
		start := maxInt(1, currentPage-2)
		end := minInt(totalPages, currentPage+2)
		
		for i := start; i <= end; i++ {
			pageClass := "pagination-link"
			var pageAttrs []mi.Attribute
			if i == currentPage {
				pageClass += " is-current"
				pageAttrs = append(pageAttrs, mi.AriaLabel("current"))
			}
			pageAttrs = append(pageAttrs, 
				mi.Class(pageClass),
				mi.Href(fmt.Sprintf("%s?page=%d", baseURL, i)))
			
			paginationList = append(paginationList, b.Li(
				b.A(pageAttrs, fmt.Sprintf("%d", i)),
			))
		}
		
		return b.Nav(mi.Class("pagination"), mi.Role("navigation"), mi.AriaLabel("pagination"),
			mi.NewFragment(pageItems...),
			b.Ul(mi.Class("pagination-list"),
				mi.NewFragment(paginationList...),
			),
		)
	}
}

// =====================================================
// DATA COMPONENTS
// =====================================================

// Table creates a Bulma table
func (t *BulmaTheme) Table(headers []string, rows [][]string) mi.H {
	return func(b *mi.Builder) mi.Node {
		// Create header row
		headerCells := make([]mi.Node, len(headers))
		for i, header := range headers {
			headerCells[i] = b.Th(header)
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
		
		return b.Div(mi.Class("table-container"),
			b.Table(mi.Class("table is-striped is-hoverable is-fullwidth"),
				b.Thead(
					b.Tr(mi.NewFragment(headerCells...)),
				),
				b.Tbody(mi.NewFragment(dataRows...)),
			),
		)
	}
}

// List creates a Bulma list
func (t *BulmaTheme) List(items []string, ordered bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		listItems := make([]mi.Node, len(items))
		for i, item := range items {
			listItems[i] = b.Li(mi.Class("mb-2"), item)
		}
		
		if ordered {
			return b.Div(mi.Class("content"),
				b.Ol(mi.NewFragment(listItems...)))
		}
		
		return b.Div(mi.Class("content"),
			b.Ul(mi.NewFragment(listItems...)))
	}
}

// =====================================================
// UTILITY METHODS
// =====================================================

// PrimaryButton creates a primary Bulma button
func (t *BulmaTheme) PrimaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "primary", attrs...)
}

// SecondaryButton creates a secondary Bulma button
func (t *BulmaTheme) SecondaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "secondary", attrs...)
}

// DangerButton creates a danger Bulma button
func (t *BulmaTheme) DangerButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "danger", attrs...)
}

// =====================================================
// HELPER METHODS
// =====================================================

// getButtonClass returns Bulma button class for variant
func (t *BulmaTheme) getButtonClass(variant string) string {
	switch variant {
	case "primary":
		return "button is-primary"
	case "secondary":
		return "button is-light"
	case "success":
		return "button is-success"
	case "warning":
		return "button is-warning"
	case "danger", "error":
		return "button is-danger"
	case "info":
		return "button is-info"
	case "light":
		return "button is-white"
	case "dark":
		return "button is-dark"
	case "link":
		return "button is-text"
	case "view":
		return "button is-primary is-outlined is-small"
	case "payment":
		return "button is-success"
	default:
		return "button"
	}
}

// getBadgeClass returns Bulma tag class for variant (Bulma uses tags instead of badges)
func (t *BulmaTheme) getBadgeClass(variant string) string {
	switch variant {
	case "primary":
		return "tag is-primary"
	case "secondary":
		return "tag is-light"
	case "success":
		return "tag is-success"
	case "warning":
		return "tag is-warning"
	case "danger", "error":
		return "tag is-danger"
	case "info":
		return "tag is-info"
	case "light":
		return "tag is-white"
	case "dark":
		return "tag is-dark"
	default:
		return "tag"
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

// CDNLinks returns Bulma CSS CDN links for inclusion in HTML head
func CDNLinks() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			// Bulma CSS
			b.Link(
				mi.Href("https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css"),
				mi.Rel("stylesheet"),
			),
			// Font Awesome for icons
			b.Link(
				mi.Href("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"),
				mi.Rel("stylesheet"),
			),
		)
	}
}

// BulmaDocument creates a complete HTML document with Bulma styling
func BulmaDocument(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.Document(title,
			[]mi.Node{
				CDNLinks()(b),
			},
			b.Body(
				b.Section(mi.Class("section"),
					b.Div(mi.Class("container"),
						content(b),
					),
				),
			),
		)(b)
	}
}
