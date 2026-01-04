// Package material provides a Material Design theme implementation for mintyui
package material

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
)

// MaterialTheme implements the Theme interface using Material Design Components
type MaterialTheme struct {
	name    string
	version string
}

// NewMaterialTheme creates a new Material Design theme
func NewMaterialTheme() mui.Theme {
	return &MaterialTheme{
		name:    "Material",
		version: "1.0.0",
	}
}

// GetName returns the theme name
func (t *MaterialTheme) GetName() string {
	return t.name
}

// GetVersion returns the theme version
func (t *MaterialTheme) GetVersion() string {
	return t.version
}

// =====================================================
// BASIC COMPONENTS
// =====================================================

// Button creates a Material Design button with ripple effect
func (t *MaterialTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getButtonClass(variant)
		allAttrs := append([]mi.Attribute{mi.Class(class), mi.Type("button")}, attrs...)
		args := make([]interface{}, len(allAttrs)+3)
		for i, attr := range allAttrs {
			args[i] = attr
		}
		// Material buttons have ripple effects and specific structure
		args[len(allAttrs)] = b.Span(mi.Class("mdc-button__ripple"))
		args[len(allAttrs)+1] = b.Span(mi.Class("mdc-button__focus-ring"))
		args[len(allAttrs)+2] = b.Span(mi.Class("mdc-button__label"), text)
		return b.Button(args...)
	}
}

// Card creates a Material Design card with elevation
func (t *MaterialTheme) Card(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mdc-card mdc-card--outlined"),
			mi.NewFragment(
				func(b *mi.Builder) mi.Node {
					if title != "" {
						return b.Div(mi.Class("mdc-card__primary-action"),
							b.Div(mi.Class("mdc-card__media mdc-card__media--16-9"),
								b.Div(mi.Class("mdc-card__media-content"),
									b.H2(mi.Class("mdc-typography--headline6"), title),
								),
							),
							b.Div(mi.Class("mdc-card__ripple")),
						)
					}
					return mi.NewFragment()
				}(b),
				b.Div(mi.Class("mdc-card__content"),
					content(b),
				),
			),
		)
	}
}

// Badge creates a Material Design chip (badge equivalent)
func (t *MaterialTheme) Badge(text, variant string) mi.H {
	return func(b *mi.Builder) mi.Node {
		class := t.getBadgeClass(variant)
		return b.Div(mi.Class(class), mi.Role("row"),
			b.Span(mi.Class("mdc-evolution-chip__cell mdc-evolution-chip__cell--primary"), mi.Role("gridcell"),
				b.Span(mi.Class("mdc-evolution-chip__action mdc-evolution-chip__action--primary"),
					b.Span(mi.Class("mdc-evolution-chip__text-label"), text),
				),
			),
		)
	}
}

// =====================================================
// FORM COMPONENTS
// =====================================================

// FormInput creates a Material Design text field with label
func (t *MaterialTheme) FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "input_" + name
		inputAttrs := append([]mi.Attribute{
			mi.Class("mdc-text-field__input"),
			mi.ID(id),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		
		return b.Div(mi.Class("mdc-text-field mdc-text-field--filled"),
			b.Span(mi.Class("mdc-text-field__ripple")),
			b.Span(mi.Class("mdc-floating-label"), mi.DataAttr("for", id), label),
			b.Input(inputAttrs...),
			b.Span(mi.Class("mdc-line-ripple")),
		)
	}
}

// FormSelect creates a Material Design select with label
func (t *MaterialTheme) FormSelect(label, name string, options []mui.SelectOption) mi.H {
	return func(b *mi.Builder) mi.Node {
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
		
		return b.Div(mi.Class("mdc-select mdc-select--filled"),
			b.Div(mi.Class("mdc-select__anchor"), mi.Role("button"), mi.AriaLabel("select"), mi.TabIndex(0),
				b.Span(mi.Class("mdc-select__ripple")),
				b.Span(mi.Class("mdc-floating-label"), label),
				b.Span(mi.Class("mdc-select__selected-text-container"),
					b.Span(mi.Class("mdc-select__selected-text")),
				),
				b.Span(mi.Class("mdc-select__dropdown-icon"),
					// Material design dropdown arrow SVG
					b.Svg(mi.Class("mdc-select__dropdown-icon-graphic"), 
						mi.ViewBox("7 10 10 5"), mi.DataAttr("focusable", "false"),
						b.Polygon(mi.Class("mdc-select__dropdown-icon-inactive"),
							mi.DataAttr("stroke", "none"), mi.DataAttr("fill-rule", "evenodd"),
							mi.DataAttr("points", "7 10 12 15 17 10")),
						b.Polygon(mi.Class("mdc-select__dropdown-icon-active"),
							mi.DataAttr("stroke", "none"), mi.DataAttr("fill-rule", "evenodd"),
							mi.DataAttr("points", "7 15 12 10 17 15")),
					),
				),
				b.Span(mi.Class("mdc-line-ripple")),
			),
			b.Div(mi.Class("mdc-select__menu mdc-menu mdc-menu-surface mdc-menu-surface--fullwidth"),
				b.Ul(mi.Class("mdc-deprecated-list"), mi.Role("listbox"), mi.AriaLabel(label),
					mi.NewFragment(func() []mi.Node {
						items := make([]mi.Node, len(options))
						for i, option := range options {
							itemClass := "mdc-deprecated-list-item"
							if option.Selected {
								itemClass += " mdc-deprecated-list-item--selected"
							}
							items[i] = b.Li(mi.Class(itemClass), mi.DataAttr("value", option.Value), mi.Role("option"),
								b.Span(mi.Class("mdc-deprecated-list-item__ripple")),
								b.Span(mi.Class("mdc-deprecated-list-item__text"), option.Text),
							)
						}
						return items
					}()...),
				),
			),
		)
	}
}

// FormTextarea creates a Material Design textarea with label
func (t *MaterialTheme) FormTextarea(label, name string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		id := "textarea_" + name
		textareaAttrs := append([]mi.Attribute{
			mi.Class("mdc-text-field__input"),
			mi.ID(id),
			mi.Name(name),
			mi.Rows(4),
		}, attrs...)
		
		// Convert to []interface{} for Textarea
		args := make([]interface{}, len(textareaAttrs))
		for i, attr := range textareaAttrs {
			args[i] = attr
		}
		
		return b.Div(mi.Class("mdc-text-field mdc-text-field--textarea"),
			b.Span(mi.Class("mdc-text-field__ripple")),
			b.Span(mi.Class("mdc-text-field__resizer"),
				b.Textarea(args...),
			),
			b.Span(mi.Class("mdc-floating-label"), mi.DataAttr("for", id), label),
			b.Span(mi.Class("mdc-line-ripple")),
		)
	}
}

// FormLabel creates a Material Design form label
func (t *MaterialTheme) FormLabel(text, forField string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Label(mi.Class("mdc-floating-label"), mi.DataAttr("for", forField), text)
	}
}

// Input creates a standalone Material Design input
func (t *MaterialTheme) Input(name, inputType string, attrs ...mi.Attribute) mi.H {
	return func(b *mi.Builder) mi.Node {
		inputAttrs := append([]mi.Attribute{
			mi.Class("mdc-text-field__input"),
			mi.Name(name),
			mi.Type(inputType),
		}, attrs...)
		
		return b.Div(mi.Class("mdc-text-field mdc-text-field--filled"),
			b.Span(mi.Class("mdc-text-field__ripple")),
			b.Input(inputAttrs...),
			b.Span(mi.Class("mdc-line-ripple")),
		)
	}
}

// =====================================================
// LAYOUT COMPONENTS
// =====================================================

// Container creates a Material Design container
func (t *MaterialTheme) Container(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mdc-layout-grid"),
			b.Div(mi.Class("mdc-layout-grid__inner"),
				b.Div(mi.Class("mdc-layout-grid__cell mdc-layout-grid__cell--span-12"),
					content(b),
				),
			),
		)
	}
}

// Grid creates a Material Design grid layout
func (t *MaterialTheme) Grid(columns int, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		colSpan := 12 / columns
		colClass := fmt.Sprintf("mdc-layout-grid__cell mdc-layout-grid__cell--span-%d", colSpan)
		
		return b.Div(mi.Class("mdc-layout-grid"),
			b.Div(mi.Class("mdc-layout-grid__inner"),
				b.Div(mi.Class(colClass),
					content(b),
				),
			),
		)
	}
}

// Sidebar creates a Material Design drawer/sidebar
func (t *MaterialTheme) Sidebar(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Aside(mi.Class("mdc-drawer mdc-drawer--permanent"),
			b.Div(mi.Class("mdc-drawer__content"),
				content(b),
			),
		)
	}
}

// =====================================================
// NAVIGATION COMPONENTS
// =====================================================

// Nav creates a Material Design navigation drawer list
func (t *MaterialTheme) Nav(items []mui.NavItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		navItems := make([]mi.Node, len(items))
		for i, item := range items {
			itemClass := "mdc-deprecated-list-item"
			if item.Active {
				itemClass += " mdc-deprecated-list-item--activated"
			}
			
			navItems[i] = b.A(mi.Class(itemClass), mi.Href(item.URL),
				b.Span(mi.Class("mdc-deprecated-list-item__ripple")),
				func(b *mi.Builder) mi.Node {
					if item.Icon != "" {
						return b.I(mi.Class("material-icons mdc-deprecated-list-item__graphic"), mi.AriaHidden(true), 
							item.Icon)
					}
					return mi.NewFragment()
				}(b),
				b.Span(mi.Class("mdc-deprecated-list-item__text"), item.Text),
			)
		}
		
		return b.Nav(mi.Class("mdc-deprecated-list"),
			mi.NewFragment(navItems...),
		)
	}
}

// Breadcrumbs creates Material Design breadcrumbs
func (t *MaterialTheme) Breadcrumbs(items []mui.BreadcrumbItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		breadcrumbItems := make([]mi.Node, len(items))
		for i, item := range items {
			if item.Last {
				breadcrumbItems[i] = b.Span(mi.Class("mdc-typography--body2"), item.Text)
			} else {
				if i > 0 {
					breadcrumbItems = append(breadcrumbItems[:i], 
						append([]mi.Node{b.I(mi.Class("material-icons"), "chevron_right")}, 
							breadcrumbItems[i:]...)...)
				}
				breadcrumbItems[i] = b.A(
					mi.Class("mdc-typography--body2"), 
					mi.Style("color: var(--mdc-theme-primary); text-decoration: none;"),
					mi.Href(item.URL), 
					item.Text,
				)
			}
		}
		
		return b.Nav(mi.Class("mdc-typography--body2"),
			mi.Style("display: flex; align-items: center; gap: 8px;"),
			mi.AriaLabel("breadcrumb"),
			mi.NewFragment(breadcrumbItems...),
		)
	}
}

// Pagination creates Material Design pagination
func (t *MaterialTheme) Pagination(currentPage, totalPages int, baseURL string) mi.H {
	return func(b *mi.Builder) mi.Node {
		pageItems := make([]mi.Node, 0)
		
		// Previous button
		prevClass := "mdc-icon-button"
		var prevAttrs []mi.Attribute
		if currentPage <= 1 {
			prevAttrs = append(prevAttrs, mi.Disabled(), mi.Class(prevClass+" mdc-icon-button--disabled"))
		} else {
			prevAttrs = append(prevAttrs, mi.Class(prevClass), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage-1)))
		}
		
		pageItems = append(pageItems, 
			b.A(prevAttrs, 
				b.Div(mi.Class("mdc-icon-button__ripple")),
				b.I(mi.Class("material-icons"), "chevron_left"),
			),
		)
		
		// Page numbers
		start := maxInt(1, currentPage-2)
		end := minInt(totalPages, currentPage+2)
		
		for i := start; i <= end; i++ {
			if i == currentPage {
				pageItems = append(pageItems,
					b.Span(mi.Class("mdc-typography--body1"),
						mi.Style("padding: 8px 16px; background-color: var(--mdc-theme-primary); color: white; border-radius: 4px;"),
						fmt.Sprintf("%d", i)),
				)
			} else {
				pageItems = append(pageItems,
					b.A(mi.Class("mdc-typography--body1"),
						mi.Style("padding: 8px 16px; text-decoration: none; color: var(--mdc-theme-primary); border-radius: 4px;"),
						mi.Href(fmt.Sprintf("%s?page=%d", baseURL, i)),
						fmt.Sprintf("%d", i)),
				)
			}
		}
		
		// Next button
		nextClass := "mdc-icon-button"
		var nextAttrs []mi.Attribute
		if currentPage >= totalPages {
			nextAttrs = append(nextAttrs, mi.Disabled(), mi.Class(nextClass+" mdc-icon-button--disabled"))
		} else {
			nextAttrs = append(nextAttrs, mi.Class(nextClass), mi.Href(fmt.Sprintf("%s?page=%d", baseURL, currentPage+1)))
		}
		
		pageItems = append(pageItems, 
			b.A(nextAttrs,
				b.Div(mi.Class("mdc-icon-button__ripple")),
				b.I(mi.Class("material-icons"), "chevron_right"),
			),
		)
		
		return b.Nav(mi.Class("mdc-typography--body1"),
			mi.Style("display: flex; align-items: center; gap: 8px; justify-content: center;"),
			mi.AriaLabel("pagination"),
			mi.NewFragment(pageItems...),
		)
	}
}

// =====================================================
// DATA COMPONENTS
// =====================================================

// Table creates a Material Design data table
func (t *MaterialTheme) Table(headers []string, rows [][]string) mi.H {
	return func(b *mi.Builder) mi.Node {
		// Create header row
		headerCells := make([]mi.Node, len(headers))
		for i, header := range headers {
			headerCells[i] = b.Th(mi.Class("mdc-data-table__header-cell"), mi.Role("columnheader"), mi.Scope("col"), 
				header)
		}
		
		// Create data rows
		dataRows := make([]mi.Node, len(rows))
		for i, row := range rows {
			cells := make([]mi.Node, len(row))
			for j, cell := range row {
				cells[j] = b.Td(mi.Class("mdc-data-table__cell"), mi.RawHTML(cell))
			}
			dataRows[i] = b.Tr(mi.Class("mdc-data-table__row"), mi.NewFragment(cells...))
		}
		
		return b.Div(mi.Class("mdc-data-table"),
			b.Div(mi.Class("mdc-data-table__table-container"),
				b.Table(mi.Class("mdc-data-table__table"), mi.AriaLabel("Data table"),
					b.Thead(
						b.Tr(mi.Class("mdc-data-table__header-row"), mi.NewFragment(headerCells...)),
					),
					b.Tbody(mi.Class("mdc-data-table__content"), mi.NewFragment(dataRows...)),
				),
			),
		)
	}
}

// List creates a Material Design list
func (t *MaterialTheme) List(items []string, ordered bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		listItems := make([]mi.Node, len(items))
		for i, item := range items {
			listItems[i] = b.Li(mi.Class("mdc-deprecated-list-item"),
				b.Span(mi.Class("mdc-deprecated-list-item__ripple")),
				b.Span(mi.Class("mdc-deprecated-list-item__text"), item),
			)
		}
		
		return b.Ul(mi.Class("mdc-deprecated-list"), mi.Role("list"),
			mi.NewFragment(listItems...))
	}
}

// =====================================================
// UTILITY METHODS
// =====================================================

// PrimaryButton creates a primary Material Design button
func (t *MaterialTheme) PrimaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "primary", attrs...)
}

// SecondaryButton creates a secondary Material Design button
func (t *MaterialTheme) SecondaryButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "secondary", attrs...)
}

// DangerButton creates a danger Material Design button
func (t *MaterialTheme) DangerButton(text string, attrs ...mi.Attribute) mi.H {
	return t.Button(text, "danger", attrs...)
}

// =====================================================
// HELPER METHODS
// =====================================================

// getButtonClass returns Material Design button class for variant
func (t *MaterialTheme) getButtonClass(variant string) string {
	switch variant {
	case "primary":
		return "mdc-button mdc-button--raised"
	case "secondary":
		return "mdc-button mdc-button--outlined"
	case "success":
		return "mdc-button mdc-button--raised"
	case "warning":
		return "mdc-button mdc-button--raised"
	case "danger", "error":
		return "mdc-button mdc-button--raised"
	case "info":
		return "mdc-button mdc-button--raised"
	case "light":
		return "mdc-button"
	case "dark":
		return "mdc-button mdc-button--raised"
	case "link":
		return "mdc-button"
	case "view":
		return "mdc-button mdc-button--outlined"
	case "payment":
		return "mdc-button mdc-button--raised"
	default:
		return "mdc-button"
	}
}

// getBadgeClass returns Material Design chip class for variant
func (t *MaterialTheme) getBadgeClass(variant string) string {
	return "mdc-evolution-chip mdc-evolution-chip--selectable"
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

// CDNLinks returns Material Design Components CDN links for inclusion in HTML head
func CDNLinks() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			// Material Design Components CSS
			b.Link(
				mi.Href("https://unpkg.com/material-components-web@latest/dist/material-components-web.min.css"),
				mi.Rel("stylesheet"),
			),
			// Material Icons
			b.Link(
				mi.Href("https://fonts.googleapis.com/icon?family=Material+Icons"),
				mi.Rel("stylesheet"),
			),
			// Roboto Font (Material Design typography)
			b.Link(
				mi.Href("https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap"),
				mi.Rel("stylesheet"),
			),
		)
	}
}

// CDNScripts returns Material Design Components CDN scripts for inclusion before closing body tag
func CDNScripts() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			b.Script(
				mi.Src("https://unpkg.com/material-components-web@latest/dist/material-components-web.min.js"),
			),
			// Initialize Material Components
			b.Script(`
				// Auto-initialize all MDC components
				window.mdc.autoInit();
				
				// Initialize ripples for buttons
				const buttons = document.querySelectorAll('.mdc-button');
				buttons.forEach((button) => {
					if (button.querySelector('.mdc-button__ripple')) {
						window.mdc.ripple.MDCRipple.attachTo(button);
					}
				});
			`),
		)
	}
}

// MaterialDocument creates a complete HTML document with Material Design styling
func MaterialDocument(title string, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.Document(title,
			[]mi.Node{
				CDNLinks()(b),
				b.Style(`
					body {
						font-family: 'Roboto', sans-serif;
						margin: 0;
						background-color: #fafafa;
					}
					:root {
						--mdc-theme-primary: #6200ee;
						--mdc-theme-secondary: #03dac6;
						--mdc-theme-error: #b00020;
						--mdc-theme-surface: #ffffff;
						--mdc-theme-on-surface: #000000;
					}
				`),
			},
			b.Body(
				b.Div(mi.Class("mdc-typography"),
					content(b),
				),
				CDNScripts()(b),
			),
		)(b)
	}
}
