package minty

import (
	"net/http"
)

// Layout creates a simple HTML page layout with title and content.
// This is the basic layout pattern for Week 2.
func Layout(title string, content H) H {
	return func(b *Builder) Node {
		return b.Html(
			b.Head(
				b.Meta(Charset("utf-8")),
				b.Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				b.Title(title),
			),
			b.Body(
				content(b),
			),
		)
	}
}

// LayoutWithCSS creates a layout with custom CSS styles.
func LayoutWithCSS(title string, cssURL string, content H) H {
	return func(b *Builder) Node {
		return b.Html(
			b.Head(
				b.Meta(Charset("utf-8")),
				b.Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				b.Title(title),
				b.Link(Rel("stylesheet"), Href(cssURL)),
			),
			b.Body(
				content(b),
			),
		)
	}
}

// LayoutWithMeta creates a layout with custom meta tags.
func LayoutWithMeta(title, description string, keywords []string, content H) H {
	return func(b *Builder) Node {
		headChildren := []interface{}{
			b.Meta(Charset("utf-8")),
			b.Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			b.Title(title),
		}
		
		if description != "" {
			headChildren = append(headChildren, 
				b.Meta(Name("description"), Content(description)))
		}
		
		if len(keywords) > 0 {
			keywordString := ""
			for i, keyword := range keywords {
				if i > 0 {
					keywordString += ", "
				}
				keywordString += keyword
			}
			headChildren = append(headChildren,
				b.Meta(Name("keywords"), Content(keywordString)))
		}
		
		return b.Html(
			b.Head(headChildren...),
			b.Body(
				content(b),
			),
		)
	}
}

// FullLayout creates a complete HTML5 layout with semantic structure.
func FullLayout(title string, nav, main, aside, footer H) H {
	return func(b *Builder) Node {
		return b.Html(Lang("en"),
			b.Head(
				b.Meta(Charset("utf-8")),
				b.Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				b.Title(title),
			),
			b.Body(
				If(nav != nil, func(b *Builder) Node {
					return b.Header(nav(b))
				})(b),
				
				b.Main(main(b)),
				
				If(aside != nil, func(b *Builder) Node {
					return aside(b)
				})(b),
				
				If(footer != nil, func(b *Builder) Node {
					return footer(b)
				})(b),
			),
		)
	}
}

// ArticleLayout creates a layout optimized for article content.
func ArticleLayout(title, author string, publishDate string, content H) H {
	return Layout(title, func(b *Builder) Node {
		return b.Article(
			b.Header(
				b.H1(title),
				If(author != "", func(b *Builder) Node {
					return b.P("By ", b.Strong(author))
				})(b),
				If(publishDate != "", func(b *Builder) Node {
					return b.Time(Datetime(publishDate), publishDate)
				})(b),
			),
			b.Section(content(b)),
		)
	})
}

// CardLayout creates a card-style layout component.
func CardLayout(title string, content H, actions H) H {
	return func(b *Builder) Node {
		return b.Article(Class("card"),
			If(title != "", func(b *Builder) Node {
				return b.Header(Class("card-header"),
					b.H3(Class("card-title"), title),
				)
			})(b),
			
			b.Section(Class("card-content"),
				content(b),
			),
			
			If(actions != nil, func(b *Builder) Node {
				return b.Footer(Class("card-actions"),
					actions(b),
				)
			})(b),
		)
	}
}

// ModalLayout creates a modal dialog layout.
func ModalLayout(id, title string, content, actions H) H {
	return func(b *Builder) Node {
		return b.Dialog(ID(id), Class("modal"),
			b.Header(Class("modal-header"),
				b.H2(Class("modal-title"), title),
				b.Button(Class("modal-close"), Type("button"), "×"),
			),
			
			b.Section(Class("modal-body"),
				content(b),
			),
			
			If(actions != nil, func(b *Builder) Node {
				return b.Footer(Class("modal-footer"),
					actions(b),
				)
			})(b),
		)
	}
}

// Component creates a reusable component wrapper.
func Component(className string, content H) H {
	return func(b *Builder) Node {
		return b.Div(Class(className), content(b))
	}
}

// Container creates a simple container wrapper.
func Container(content H) H {
	return Component("container", content)
}

// Section creates a semantic section with optional heading.
func SectionLayout(heading string, content H) H {
	return func(b *Builder) Node {
		return b.Section(
			If(heading != "", func(b *Builder) Node {
				return b.H2(heading)
			})(b),
			content(b),
		)
	}
}

// Grid creates a CSS Grid layout container.
func GridLayout(columns string, gap string, content H) H {
	return func(b *Builder) Node {
		gridStyle := "display: grid;"
		if columns != "" {
			gridStyle += " grid-template-columns: " + columns + ";"
		}
		if gap != "" {
			gridStyle += " gap: " + gap + ";"
		}
		
		return b.Div(Style(gridStyle), content(b))
	}
}

// FlexLayout creates a flexbox layout container.
func FlexLayout(direction, justify, align string, content H) H {
	return func(b *Builder) Node {
		flexStyle := "display: flex;"
		if direction != "" {
			flexStyle += " flex-direction: " + direction + ";"
		}
		if justify != "" {
			flexStyle += " justify-content: " + justify + ";"
		}
		if align != "" {
			flexStyle += " align-items: " + align + ";"
		}
		
		return b.Div(Style(flexStyle), content(b))
	}
}

// HTTP Integration Helpers

// RenderHandler creates an HTTP handler that renders a Minty template.
func RenderHandler(template H) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := Render(template, w); err != nil {
			http.Error(w, "Template render error", http.StatusInternalServerError)
		}
	}
}

// RenderHandlerFunc creates an HTTP handler from a function that returns a template.
func RenderHandlerFunc(fn func(*http.Request) H) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		template := fn(r)
		if err := Render(template, w); err != nil {
			http.Error(w, "Template render error", http.StatusInternalServerError)
		}
	}
}

// Common layout components

// Navigation creates a semantic navigation bar.
func Navigation(links []NavLink) H {
	return func(b *Builder) Node {
		navItems := make([]interface{}, 0, len(links)*2+1)
		navItems = append(navItems, Class("navigation"))
		for i, link := range links {
			if i > 0 {
				navItems = append(navItems, " | ")
			}
			navItems = append(navItems, b.A(Href(link.URL), link.Text))
		}
		
		return b.Nav(navItems...)
	}
}

// BreadcrumbNavigation creates a breadcrumb navigation.
func BreadcrumbNavigation(links []NavLink) H {
	return func(b *Builder) Node {
		return b.Nav(Class("breadcrumb"), Role("navigation"), AriaLabel("Breadcrumb"),
			b.Ol(
				NewFragment(Each(links, func(link NavLink) H {
					return func(b *Builder) Node {
						return b.Li(b.A(Href(link.URL), link.Text))
					}
				})...),
			),
		)
	}
}

// NavLink represents a navigation link.
type NavLink struct {
	URL  string
	Text string
}

// Message components

// ErrorMessage creates a styled error message.
func ErrorMessage(message string) H {
	return func(b *Builder) Node {
		if message == "" {
			return NewFragment() // Empty fragment for no error
		}
		return b.Div(
			Class("message error"),
			Role("alert"),
			Style("color: #d32f2f; background: #ffebee; border: 1px solid #ffcdd2; padding: 0.75rem; border-radius: 4px; margin: 0.5rem 0;"),
			b.P(message),
		)
	}
}

// SuccessMessage creates a styled success message.
func SuccessMessage(message string) H {
	return func(b *Builder) Node {
		if message == "" {
			return NewFragment()
		}
		return b.Div(
			Class("message success"),
			Role("status"),
			Style("color: #2e7d32; background: #e8f5e8; border: 1px solid #c8e6c9; padding: 0.75rem; border-radius: 4px; margin: 0.5rem 0;"),
			b.P(message),
		)
	}
}

// WarningMessage creates a styled warning message.
func WarningMessage(message string) H {
	return func(b *Builder) Node {
		if message == "" {
			return NewFragment()
		}
		return b.Div(
			Class("message warning"),
			Role("alert"),
			Style("color: #f57c00; background: #fff8e1; border: 1px solid #ffecb3; padding: 0.75rem; border-radius: 4px; margin: 0.5rem 0;"),
			b.P(message),
		)
	}
}

// InfoMessage creates a styled info message.
func InfoMessage(message string) H {
	return func(b *Builder) Node {
		if message == "" {
			return NewFragment()
		}
		return b.Div(
			Class("message info"),
			Role("status"),
			Style("color: #1976d2; background: #e3f2fd; border: 1px solid #bbdefb; padding: 0.75rem; border-radius: 4px; margin: 0.5rem 0;"),
			b.P(message),
		)
	}
}

// Accessibility helpers

// SkipLink creates an accessibility skip link.
func SkipLink(href, text string) H {
	return func(b *Builder) Node {
		return b.A(
			Href(href),
			Class("skip-link"),
			Style("position: absolute; top: -40px; left: 6px; z-index: 999; color: white; background: #000; padding: 8px; text-decoration: none; border-radius: 0 0 4px 4px;"),
			text,
		)
	}
}

// VisuallyHidden creates content that's hidden visually but available to screen readers.
func VisuallyHidden(content H) H {
	return func(b *Builder) Node {
		return b.Span(
			Class("visually-hidden"),
			Style("position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0,0,0,0); white-space: nowrap; border: 0;"),
			content(b),
		)
	}
}
