package mintydyn

import (
	"fmt"
	"strings"

	mi "github.com/ha1tch/minty"
)

// =============================================================================
// CSS BUILDER
// =============================================================================

// CSSBuilder provides a fluent API for building CSS.
type CSSBuilder struct {
	rules []cssRule
}

type cssRule struct {
	selector   string
	properties []cssProperty
}

type cssProperty struct {
	name  string
	value string
}

// NewCSSBuilder creates a new CSS builder.
func NewCSSBuilder() *CSSBuilder {
	return &CSSBuilder{
		rules: make([]cssRule, 0),
	}
}

// Rule adds a CSS rule.
func (c *CSSBuilder) Rule(selector string, properties ...CSSProperty) *CSSBuilder {
	props := make([]cssProperty, len(properties))
	for i, p := range properties {
		props[i] = cssProperty{name: p.Name, value: p.Value}
	}
	c.rules = append(c.rules, cssRule{
		selector:   selector,
		properties: props,
	})
	return c
}

// Render produces the CSS string.
func (c *CSSBuilder) Render() string {
	var css strings.Builder

	for _, rule := range c.rules {
		css.WriteString(rule.selector)
		css.WriteString(" {\n")

		for _, prop := range rule.properties {
			css.WriteString(fmt.Sprintf("    %s: %s;\n", prop.name, prop.value))
		}

		css.WriteString("}\n\n")
	}

	return css.String()
}

// ToStyleNode creates a minty Style node with the CSS.
func (c *CSSBuilder) ToStyleNode(b *mi.Builder) mi.Node {
	return b.Style(mi.Raw(c.Render()))
}

// =============================================================================
// CSS PROPERTY TYPE AND CONSTRUCTORS
// =============================================================================

// CSSProperty represents a single CSS property.
type CSSProperty struct {
	Name  string
	Value string
}

// Prop creates a CSS property.
func Prop(name, value string) CSSProperty {
	return CSSProperty{Name: name, Value: value}
}

// Common CSS properties

func Display(v string) CSSProperty         { return Prop("display", v) }
func Position(v string) CSSProperty        { return Prop("position", v) }
func Width(v string) CSSProperty           { return Prop("width", v) }
func Height(v string) CSSProperty          { return Prop("height", v) }
func MinWidth(v string) CSSProperty        { return Prop("min-width", v) }
func MinHeight(v string) CSSProperty       { return Prop("min-height", v) }
func MaxWidth(v string) CSSProperty        { return Prop("max-width", v) }
func MaxHeight(v string) CSSProperty       { return Prop("max-height", v) }

func Margin(v string) CSSProperty          { return Prop("margin", v) }
func MarginTop(v string) CSSProperty       { return Prop("margin-top", v) }
func MarginRight(v string) CSSProperty     { return Prop("margin-right", v) }
func MarginBottom(v string) CSSProperty    { return Prop("margin-bottom", v) }
func MarginLeft(v string) CSSProperty      { return Prop("margin-left", v) }

func Padding(v string) CSSProperty         { return Prop("padding", v) }
func PaddingTop(v string) CSSProperty      { return Prop("padding-top", v) }
func PaddingRight(v string) CSSProperty    { return Prop("padding-right", v) }
func PaddingBottom(v string) CSSProperty   { return Prop("padding-bottom", v) }
func PaddingLeft(v string) CSSProperty     { return Prop("padding-left", v) }

func Background(v string) CSSProperty      { return Prop("background", v) }
func BackgroundColor(v string) CSSProperty { return Prop("background-color", v) }
func Color(v string) CSSProperty           { return Prop("color", v) }

func Border(v string) CSSProperty          { return Prop("border", v) }
func BorderTop(v string) CSSProperty       { return Prop("border-top", v) }
func BorderRight(v string) CSSProperty     { return Prop("border-right", v) }
func BorderBottom(v string) CSSProperty    { return Prop("border-bottom", v) }
func BorderLeft(v string) CSSProperty      { return Prop("border-left", v) }
func BorderRadius(v string) CSSProperty    { return Prop("border-radius", v) }
func BorderColor(v string) CSSProperty     { return Prop("border-color", v) }

func FontFamily(v string) CSSProperty      { return Prop("font-family", v) }
func FontSize(v string) CSSProperty        { return Prop("font-size", v) }
func FontWeight(v string) CSSProperty      { return Prop("font-weight", v) }
func LineHeight(v string) CSSProperty      { return Prop("line-height", v) }
func TextAlign(v string) CSSProperty       { return Prop("text-align", v) }
func TextDecoration(v string) CSSProperty  { return Prop("text-decoration", v) }

func FlexDirection(v string) CSSProperty   { return Prop("flex-direction", v) }
func JustifyContent(v string) CSSProperty  { return Prop("justify-content", v) }
func AlignItems(v string) CSSProperty      { return Prop("align-items", v) }
func FlexWrap(v string) CSSProperty        { return Prop("flex-wrap", v) }
func Gap(v string) CSSProperty             { return Prop("gap", v) }

func GridTemplateColumns(v string) CSSProperty { return Prop("grid-template-columns", v) }
func GridTemplateRows(v string) CSSProperty    { return Prop("grid-template-rows", v) }
func GridGap(v string) CSSProperty             { return Prop("grid-gap", v) }

func BoxShadow(v string) CSSProperty       { return Prop("box-shadow", v) }
func Opacity(v string) CSSProperty         { return Prop("opacity", v) }
func Overflow(v string) CSSProperty        { return Prop("overflow", v) }
func Cursor(v string) CSSProperty          { return Prop("cursor", v) }
func ZIndex(v string) CSSProperty          { return Prop("z-index", v) }

func Transition(v string) CSSProperty      { return Prop("transition", v) }
func Transform(v string) CSSProperty       { return Prop("transform", v) }
func Animation(v string) CSSProperty       { return Prop("animation", v) }

// =============================================================================
// DEFAULT CSS FOR DYNAMIC COMPONENTS
// =============================================================================

// DefaultCSS returns minimal CSS for dynamic components to work without a theme.
func DefaultCSS() string {
	return NewCSSBuilder().
		// Hidden utility
		Rule(".hidden",
			Display("none !important"),
		).
		// Component container
		Rule(".dyn-component",
			Position("relative"),
		).
		// State navigation (tabs)
		Rule(".dyn-state-navigation",
			Display("flex"),
			Gap("0.25rem"),
			BorderBottom("2px solid #e5e7eb"),
			MarginBottom("1rem"),
		).
		Rule(".dyn-state-trigger",
			Padding("0.75rem 1.25rem"),
			Border("none"),
			Background("transparent"),
			Cursor("pointer"),
			FontSize("0.875rem"),
			FontWeight("500"),
			Color("#6b7280"),
			BorderBottom("2px solid transparent"),
			MarginBottom("-2px"),
			Transition("all 0.15s ease"),
		).
		Rule(".dyn-state-trigger:hover",
			Color("#374151"),
			BackgroundColor("#f9fafb"),
		).
		Rule(".dyn-state-trigger.active",
			Color("#2563eb"),
			BorderBottom("2px solid #2563eb"),
		).
		Rule(".dyn-state-trigger.disabled",
			Opacity("0.5"),
			Cursor("not-allowed"),
		).
		// State content
		Rule(".dyn-state-content",
			Display("none"),
		).
		Rule(".dyn-state-content.active",
			Display("block"),
		).
		// Filter controls
		Rule(".dyn-filter-controls",
			Display("grid"),
			GridTemplateColumns("repeat(auto-fit, minmax(200px, 1fr))"),
			Gap("1rem"),
			MarginBottom("1rem"),
		).
		Rule(".dyn-filter-group",
			Display("flex"),
			FlexDirection("column"),
		).
		Rule(".dyn-filter-label",
			FontSize("0.875rem"),
			FontWeight("500"),
			Color("#374151"),
			MarginBottom("0.25rem"),
		).
		Rule(".dyn-filter-input, .dyn-filter-select",
			Padding("0.5rem 0.75rem"),
			Border("1px solid #d1d5db"),
			BorderRadius("0.375rem"),
			FontSize("0.875rem"),
			Transition("border-color 0.15s ease"),
		).
		Rule(".dyn-filter-input:focus, .dyn-filter-select:focus",
			Prop("outline", "none"),
			BorderColor("#2563eb"),
			BoxShadow("0 0 0 3px rgba(37, 99, 235, 0.1)"),
		).
		// Results
		Rule(".dyn-results",
			MinHeight("100px"),
		).
		Rule(".dyn-no-results",
			TextAlign("center"),
			Padding("2rem"),
			Color("#6b7280"),
		).
		Rule(".dyn-results-summary",
			FontSize("0.875rem"),
			Color("#6b7280"),
			MarginBottom("0.5rem"),
		).
		// Pagination
		Rule(".dyn-pagination",
			Display("flex"),
			JustifyContent("center"),
			Gap("0.25rem"),
			MarginTop("1rem"),
		).
		Rule(".dyn-page-btn",
			Padding("0.5rem 0.75rem"),
			Border("1px solid #d1d5db"),
			Background("white"),
			Cursor("pointer"),
			BorderRadius("0.25rem"),
			FontSize("0.875rem"),
			Transition("all 0.15s ease"),
		).
		Rule(".dyn-page-btn:hover",
			BackgroundColor("#f3f4f6"),
		).
		Rule(".dyn-page-btn.active",
			BackgroundColor("#2563eb"),
			BorderColor("#2563eb"),
			Color("white"),
		).
		Render()
}

// DefaultCSSNode returns the default CSS as a style node.
func DefaultCSSNode(b *mi.Builder) mi.Node {
	return b.Style(mi.Raw(DefaultCSS()))
}

// =============================================================================
// CSS INJECTION HELPER
// =============================================================================

// WithDefaultCSS wraps a component with the default CSS.
// Use this if you're not using a CSS framework.
func WithDefaultCSS(component mi.H) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.NewFragment(
			DefaultCSSNode(b),
			component(b),
		)
	}
}
