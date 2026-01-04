// Package mintyui provides theme-based UI components and utilities
// for building consistent user interfaces across domains.
package mintyui

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	"github.com/ha1tch/minty/mintyex"
)

// =====================================================
// THEME INTERFACE
// =====================================================

// Theme represents a UI theme that provides consistent styling
type Theme interface {
	// Basic components
	Button(text, variant string, attrs ...mi.Attribute) mi.H
	Card(title string, content mi.H) mi.H
	Badge(text, variant string) mi.H
	
	// Form components
	FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H
	FormSelect(label, name string, options []SelectOption) mi.H
	FormTextarea(label, name string, attrs ...mi.Attribute) mi.H
	FormLabel(text, forField string) mi.H
	Input(name, inputType string, attrs ...mi.Attribute) mi.H
	
	// Layout components
	Container(content mi.H) mi.H
	Grid(columns int, content mi.H) mi.H
	Sidebar(content mi.H) mi.H
	
	// Navigation components
	Nav(items []NavItem) mi.H
	Breadcrumbs(items []BreadcrumbItem) mi.H
	Pagination(currentPage, totalPages int, baseURL string) mi.H
	
	// Data components
	Table(headers []string, rows [][]string) mi.H
	List(items []string, ordered bool) mi.H
	
	// Utility methods
	PrimaryButton(text string, attrs ...mi.Attribute) mi.H
	SecondaryButton(text string, attrs ...mi.Attribute) mi.H
	DangerButton(text string, attrs ...mi.Attribute) mi.H
	
	// Theme metadata
	GetName() string
	GetVersion() string
}

// SelectOption represents an option in a select dropdown
type SelectOption struct {
	Value    string
	Text     string
	Selected bool
	Disabled bool
}

// NavItem represents a navigation menu item
type NavItem struct {
	Text   string
	URL    string
	Active bool
	Icon   string
}

// BreadcrumbItem represents a breadcrumb navigation item
type BreadcrumbItem struct {
	Text string
	URL  string
	Last bool
}

// =====================================================
// DOMAIN-SPECIFIC COMPONENTS
// =====================================================

// DomainCard creates a card component with domain-specific styling
func DomainCard(theme Theme, domain, title string, content mi.H) mi.H {
	classes := mintyex.NewSemanticClasses(domain)
	
	return func(b *mi.Builder) mi.Node {
		cardClass := classes.Container("card")
		titleClass := classes.Header("title")
		
		return b.Div(mi.Class(cardClass),
			b.Header(mi.Class(titleClass),
				b.H3(title),
			),
			b.Div(mi.Class(classes.Content("body")),
				content(b),
			),
		)
	}
}

// DomainButton creates a button with domain-specific styling
func DomainButton(theme Theme, domain, text, variant string, attrs ...mi.Attribute) mi.H {
	classes := mintyex.NewSemanticClasses(domain)
	
	return func(b *mi.Builder) mi.Node {
		var buttonClass string
		switch variant {
		case "primary":
			buttonClass = classes.Primary("button")
		case "secondary":
			buttonClass = classes.Secondary("button")
		case "success":
			buttonClass = classes.Success("button")
		case "warning":
			buttonClass = classes.Warning("button")
		case "error", "danger":
			buttonClass = classes.Error("button")
		case "info":
			buttonClass = classes.Info("button")
		default:
			buttonClass = classes.Secondary("button")
		}
		
		// Add button class to existing attributes
		allAttrs := append([]mi.Attribute{mi.Class(buttonClass)}, attrs...)
		return b.Button(allAttrs, text)
	}
}

// StatusIndicator creates a status indicator with appropriate styling
func StatusIndicator(theme Theme, status mintyex.Status) mi.H {
	return theme.Badge(status.GetDisplay(), status.GetSeverity())
}

// StatsCard creates a statistics display card
func StatsCard(theme Theme, title, value, description string) mi.H {
	return func(b *mi.Builder) mi.Node {
		cardStyle := "border: 1px solid #e2e8f0; border-radius: 8px; padding: 20px; text-align: center; background: white;"
		titleStyle := "font-size: 14px; color: #64748b; margin: 0 0 8px 0; text-transform: uppercase; letter-spacing: 0.05em;"
		valueStyle := "font-size: 32px; font-weight: 700; color: #1e293b; margin: 0 0 8px 0;"
		descStyle := "font-size: 12px; color: #64748b; margin: 0;"
		
		return b.Div(mi.Style(cardStyle),
			b.P(mi.Style(titleStyle), title),
			b.Div(mi.Style(valueStyle), value),
			b.P(mi.Style(descStyle), description),
		)
	}
}

// Dashboard creates a dashboard layout with sidebar and main content
func Dashboard(theme Theme, title string, sidebar mi.H, content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		containerStyle := "display: grid; grid-template-columns: 250px 1fr; min-height: 100vh; background: #f8fafc;"
		sidebarStyle := "background: white; border-right: 1px solid #e2e8f0; padding: 20px; overflow-y: auto;"
		mainStyle := "padding: 20px; overflow-y: auto;"
		headerStyle := "margin: 0 0 20px 0; font-size: 28px; color: #1e293b;"
		
		return b.Div(mi.Style(containerStyle),
			b.Aside(mi.Style(sidebarStyle),
				sidebar(b),
			),
			b.Main(mi.Style(mainStyle),
				b.H1(mi.Style(headerStyle), title),
				content(b),
			),
		)
	}
}

// =====================================================
// UTILITY COMPONENTS
// =====================================================

// Empty returns an empty HTML fragment
func Empty() mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment()
	}
}

// Loading creates a loading spinner component
func Loading(text string) mi.H {
	return func(b *mi.Builder) mi.Node {
		spinnerStyle := "display: inline-block; width: 20px; height: 20px; border: 3px solid #f3f3f3; border-top: 3px solid #3498db; border-radius: 50%; animation: spin 1s linear infinite;"
		containerStyle := "display: flex; align-items: center; gap: 10px; padding: 20px; justify-content: center;"
		
		return b.Div(mi.Style(containerStyle),
			b.Div(mi.Style(spinnerStyle)),
			mintyex.If(text != "", func(b *mi.Builder) mi.Node {
				return b.Span(text)
			})(b),
		)
	}
}

// ErrorMessage creates an error message component
func ErrorMessage(message string) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "background: #fee2e2; border: 1px solid #fecaca; color: #dc2626; padding: 12px; border-radius: 6px; margin: 10px 0;"
		return b.Div(mi.Style(style), "⚠️ ", message)
	}
}

// SuccessMessage creates a success message component
func SuccessMessage(message string) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "background: #d1fae5; border: 1px solid #a7f3d0; color: #047857; padding: 12px; border-radius: 6px; margin: 10px 0;"
		return b.Div(mi.Style(style), "✅ ", message)
	}
}

// InfoMessage creates an info message component
func InfoMessage(message string) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "background: #dbeafe; border: 1px solid #93c5fd; color: #1d4ed8; padding: 12px; border-radius: 6px; margin: 10px 0;"
		return b.Div(mi.Style(style), "ℹ️ ", message)
	}
}

// Modal creates a modal dialog component
func Modal(id, title string, content mi.H, showCloseButton bool) mi.H {
	return func(b *mi.Builder) mi.Node {
		overlayStyle := "position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000;"
		modalStyle := "background: white; border-radius: 8px; padding: 0; max-width: 500px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);"
		headerStyle := "padding: 20px; border-bottom: 1px solid #e2e8f0; display: flex; justify-content: space-between; align-items: center;"
		titleStyle := "margin: 0; font-size: 18px; color: #1e293b;"
		contentStyle := "padding: 20px;"
		closeStyle := "background: none; border: none; font-size: 24px; cursor: pointer; color: #64748b;"
		
		return b.Div(mi.ID(id), mi.Style(overlayStyle),
			b.Div(mi.Style(modalStyle),
				b.Header(mi.Style(headerStyle),
					b.H2(mi.Style(titleStyle), title),
					mintyex.If(showCloseButton, func(b *mi.Builder) mi.Node {
						return b.Button(mi.Style(closeStyle), mi.Type("button"), "×")
					})(b),
				),
				b.Div(mi.Style(contentStyle),
					content(b),
				),
			),
		)
	}
}

// Tabs creates a tabbed interface component
func Tabs(activeTab string, tabs []TabItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		tabsStyle := "border-bottom: 1px solid #e2e8f0; margin-bottom: 20px;"
		tabListStyle := "display: flex; gap: 0; margin: 0; padding: 0; list-style: none;"
		tabStyle := "padding: 10px 20px; border: none; background: none; cursor: pointer; color: #64748b; border-bottom: 2px solid transparent;"
		activeTabStyle := "padding: 10px 20px; border: none; background: none; cursor: pointer; color: #3b82f6; border-bottom: 2px solid #3b82f6;"
		
		var tabButtons []mi.Node
		for _, tab := range tabs {
			style := tabStyle
			if tab.ID == activeTab {
				style = activeTabStyle
			}
			
			tabButtons = append(tabButtons, b.Button(
				mi.Style(style),
				mi.Type("button"),
				mi.DataAttr("tab", tab.ID),
				tab.Label,
			))
		}
		
		return b.Div(mi.Style(tabsStyle),
			b.Div(mi.Style(tabListStyle),
				mi.NewFragment(tabButtons...),
			),
		)
	}
}

// TabItem represents a tab in a tabbed interface
type TabItem struct {
	ID     string
	Label  string
	Active bool
}

// ProgressBar creates a progress bar component
func ProgressBar(value, max int, label string) mi.H {
	return func(b *mi.Builder) mi.Node {
		percentage := float64(value) / float64(max) * 100
		
		containerStyle := "margin: 10px 0;"
		labelStyle := "display: flex; justify-content: space-between; margin-bottom: 5px; font-size: 14px; color: #64748b;"
		barStyle := "width: 100%; height: 8px; background: #e2e8f0; border-radius: 4px; overflow: hidden;"
		fillStyle := fmt.Sprintf("height: 100%%; background: #3b82f6; width: %.1f%%; transition: width 0.3s ease;", percentage)
		
		return b.Div(mi.Style(containerStyle),
			b.Div(mi.Style(labelStyle),
				b.Span(label),
				b.Span(fmt.Sprintf("%d/%d (%.1f%%)", value, max, percentage)),
			),
			b.Div(mi.Style(barStyle),
				b.Div(mi.Style(fillStyle)),
			),
		)
	}
}

// Tooltip creates a tooltip component
func Tooltip(content mi.H, tooltipText string) mi.H {
	return func(b *mi.Builder) mi.Node {
		containerStyle := "position: relative; display: inline-block;"
		tooltipStyle := "position: absolute; bottom: 100%; left: 50%; transform: translateX(-50%); background: #1e293b; color: white; padding: 5px 10px; border-radius: 4px; font-size: 12px; white-space: nowrap; opacity: 0; pointer-events: none; transition: opacity 0.3s; z-index: 100;"
		
		return b.Div(mi.Style(containerStyle),
			mi.DataAttr("tooltip", tooltipText),
			content(b),
			b.Div(mi.Style(tooltipStyle), tooltipText),
		)
	}
}

// =====================================================
// FORM UTILITIES
// =====================================================

// FormGroup creates a form group with label and input
func FormGroup(label string, required bool, input mi.H, help string, errors []string) mi.H {
	return func(b *mi.Builder) mi.Node {
		groupStyle := "margin-bottom: 20px;"
		labelStyle := "display: block; margin-bottom: 5px; font-weight: 500; color: #374151;"
		helpStyle := "margin-top: 5px; font-size: 14px; color: #6b7280;"
		errorStyle := "margin-top: 5px; font-size: 14px; color: #dc2626;"
		
		labelText := label
		if required {
			labelText += " *"
		}
		
		var children []mi.Node
		children = append(children, b.Label(mi.Style(labelStyle), labelText))
		children = append(children, input(b))
		
		if help != "" {
			children = append(children, b.P(mi.Style(helpStyle), help))
		}
		
		for _, err := range errors {
			children = append(children, b.P(mi.Style(errorStyle), err))
		}
		
		return b.Div(mi.Style(groupStyle),
			mi.NewFragment(children...),
		)
	}
}

// FormActions creates a form actions section with buttons
func FormActions(buttons ...mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "display: flex; gap: 10px; justify-content: flex-end; margin-top: 20px; padding-top: 20px; border-top: 1px solid #e2e8f0;"
		
		var buttonNodes []mi.Node
		for _, button := range buttons {
			buttonNodes = append(buttonNodes, button(b))
		}
		
		return b.Div(mi.Style(style),
			mi.NewFragment(buttonNodes...),
		)
	}
}

// =====================================================
// RESPONSIVE UTILITIES
// =====================================================

// ResponsiveGrid creates a responsive grid layout
func ResponsiveGrid(content ...mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px;"
		
		var children []mi.Node
		for _, item := range content {
			children = append(children, item(b))
		}
		
		return b.Div(mi.Style(style),
			mi.NewFragment(children...),
		)
	}
}

// Stack creates a vertical stack layout
func Stack(gap string, content ...mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := fmt.Sprintf("display: flex; flex-direction: column; gap: %s;", gap)
		
		var children []mi.Node
		for _, item := range content {
			children = append(children, item(b))
		}
		
		return b.Div(mi.Style(style),
			mi.NewFragment(children...),
		)
	}
}

// HStack creates a horizontal stack layout
func HStack(gap, align string, content ...mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := fmt.Sprintf("display: flex; gap: %s; align-items: %s;", gap, align)
		
		var children []mi.Node
		for _, item := range content {
			children = append(children, item(b))
		}
		
		return b.Div(mi.Style(style),
			mi.NewFragment(children...),
		)
	}
}

// =====================================================
// ACCESSIBILITY UTILITIES
// =====================================================

// ScreenReaderOnly creates content only visible to screen readers
func ScreenReaderOnly(content string) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0;"
		return b.Span(mi.Style(style), content)
	}
}

// FocusRing adds a focus ring for keyboard navigation
func FocusRing(content mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		style := "outline: 2px solid transparent; outline-offset: 2px;"
		focusStyle := "outline: 2px solid #3b82f6; outline-offset: 2px;"
		
		return b.Div(
			mi.Style(style),
			mi.DataAttr("focus-style", focusStyle),
			content(b),
		)
	}
}
