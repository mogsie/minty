# Minty System Documentation - Part 9
## Theme System: Pluggable Styling and Component Architecture

---

### Table of Contents
1. [Theme System Overview](#theme-system-overview)
2. [Theme Interface Architecture](#theme-interface-architecture)
3. [Bootstrap Theme Implementation](#bootstrap-theme-implementation)
4. [Tailwind Theme Implementation](#tailwind-theme-implementation)
5. [Creating Custom Themes](#creating-custom-themes)
6. [Domain-Specific Theming](#domain-specific-theming)
7. [Advanced Theme Patterns](#advanced-theme-patterns)
8. [Theme Integration Best Practices](#theme-integration-best-practices)

---

## Theme System Overview

The Minty System's theme architecture provides complete separation between business logic, component structure, and visual styling. This enables applications to support multiple visual themes, white-labeling, and brand customization without changing any business code or component logic.

### Core Theme Architecture

```
┌─────────────────────────────────────┐
│           Business Logic            │ ← Domain services (mintyfin, etc.)
│        (Theme Independent)          │
├─────────────────────────────────────┤
│         Component Logic             │ ← Presentation adapters  
│        (Theme Aware)                │
├─────────────────────────────────────┤
│         Theme Interface             │ ← Pluggable styling layer
│    (Bootstrap, Tailwind, Custom)    │
├─────────────────────────────────────┤
│          HTML Output                │ ← Rendered components
│      (Theme Specific CSS)           │
└─────────────────────────────────────┘
```

### Key Benefits

**Complete Separation**: Business logic has zero knowledge of visual styling or CSS frameworks.

**Runtime Theme Switching**: Applications can switch themes dynamically without restarting or recompiling.

**White-labeling Support**: Different themes for different customers/brands within the same application.

**CSS Framework Agnostic**: Theme interface abstracts away specific CSS framework details.

**Consistent Components**: Same business components work identically across all themes.

### Theme Philosophy

The Minty theme system follows the principle that **visual presentation should be pluggable without affecting functionality**. A financial dashboard should work identically whether styled with Bootstrap, Tailwind, or a custom corporate theme - only the visual appearance changes.

---

## Theme Interface Architecture

### Core Theme Interface

The theme interface defines a comprehensive set of UI components that any theme must implement:

```go
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
```

### Component Configuration Types

```go
// Select dropdown configuration
type SelectOption struct {
    Value    string
    Text     string
    Selected bool
    Disabled bool
}

// Navigation configuration
type NavItem struct {
    Text   string
    URL    string
    Active bool
    Icon   string
}

// Breadcrumb configuration
type BreadcrumbItem struct {
    Text string
    URL  string
    Last bool
}
```

### Theme Usage Pattern

```go
func UserDashboard(theme mui.Theme, user User, accounts []Account) mi.H {
    return theme.Container(
        func(b *mi.Builder) mi.Node {
            return b.Div(
                theme.Nav([]mui.NavItem{
                    {Text: "Dashboard", URL: "/dashboard", Active: true},
                    {Text: "Accounts", URL: "/accounts"},
                    {Text: "Reports", URL: "/reports"},
                }),
                
                theme.Card("Account Summary", func(b *mi.Builder) mi.Node {
                    return b.Div(
                        mintyex.Map(accounts, func(a Account) mi.H {
                            return theme.Card(a.Name, func(b *mi.Builder) mi.Node {
                                return b.P("Balance: ", a.Balance.Format())
                            })
                        })...,
                    )
                }),
            )
        },
    )
}
```

The same component code works with any theme implementation - switching themes changes only visual appearance, not functionality.

---

## Bootstrap Theme Implementation

The Bootstrap theme provides comprehensive Bootstrap 5 styling with proper component implementations and CSS class management.

### Bootstrap Theme Structure

```go
type BootstrapTheme struct {
    name    string
    version string
}

func NewBootstrapTheme() mui.Theme {
    return &BootstrapTheme{
        name:    "Bootstrap",
        version: "1.0.0",
    }
}
```

### Button Implementation

```go
func (t *BootstrapTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        class := t.getButtonClass(variant)
        allAttrs := append([]mi.Attribute{mi.Class(class), mi.Type("button")}, attrs...)
        return b.Button(allAttrs, text)
    }
}

func (t *BootstrapTheme) getButtonClass(variant string) string {
    baseClass := "btn"
    
    switch variant {
    case "primary":
        return baseClass + " btn-primary"
    case "secondary":
        return baseClass + " btn-secondary"
    case "success":
        return baseClass + " btn-success"
    case "danger":
        return baseClass + " btn-danger"
    case "warning":
        return baseClass + " btn-warning"
    case "info":
        return baseClass + " btn-info"
    case "light":
        return baseClass + " btn-light"
    case "dark":
        return baseClass + " btn-dark"
    case "link":
        return baseClass + " btn-link"
    default:
        return baseClass + " btn-primary"
    }
}
```

### Card Component

```go
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
```

### Form Components

```go
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
```

### Table Implementation

```go
func (t *BootstrapTheme) Table(headers []string, rows [][]string) mi.H {
    return func(b *mi.Builder) mi.Node {
        // Create header nodes
        headerNodes := mintyex.Map(headers, func(header string) mi.Node {
            return b.Th(mi.Scope("col"), header)
        })
        
        // Create row nodes
        rowNodes := mintyex.Map(rows, func(row []string) mi.Node {
            cellNodes := mintyex.Map(row, func(cell string) mi.Node {
                return b.Td(mi.RawHTML(cell)) // Allow HTML in cells
            })
            return b.Tr(cellNodes...)
        })
        
        return b.Div(mi.Class("table-responsive"),
            b.Table(mi.Class("table table-striped table-hover"),
                b.Thead(mi.Class("table-dark"),
                    b.Tr(headerNodes...),
                ),
                b.Tbody(rowNodes...),
            ),
        )
    }
}
```

### Navigation Components

```go
func (t *BootstrapTheme) Nav(items []mui.NavItem) mi.H {
    return func(b *mi.Builder) mi.Node {
        navItems := mintyex.Map(items, func(item mui.NavItem) mi.Node {
            linkClass := "nav-link"
            if item.Active {
                linkClass += " active"
            }
            
            return b.Li(mi.Class("nav-item"),
                b.A(mi.Class(linkClass), mi.Href(item.URL), 
                    func() string {
                        if item.Icon != "" {
                            return fmt.Sprintf(`<i class="%s"></i> %s`, item.Icon, item.Text)
                        }
                        return item.Text
                    }(),
                ),
            )
        })
        
        return b.Nav(mi.Class("navbar navbar-expand-lg navbar-dark bg-dark"),
            b.Div(mi.Class("container-fluid"),
                b.Ul(mi.Class("navbar-nav me-auto"),
                    navItems...,
                ),
            ),
        )
    }
}
```

---

## Tailwind Theme Implementation

The Tailwind theme provides utility-first CSS styling with Tailwind's comprehensive class system.

### Tailwind Theme Structure

```go
type TailwindTheme struct {
    name    string
    version string
}

func NewTailwindTheme() mui.Theme {
    return &TailwindTheme{
        name:    "Tailwind",
        version: "1.0.0",
    }
}
```

### Button Implementation

```go
func (t *TailwindTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        class := t.getButtonClass(variant)
        allAttrs := append([]mi.Attribute{mi.Class(class), mi.Type("button")}, attrs...)
        return b.Button(allAttrs, text)
    }
}

func (t *TailwindTheme) getButtonClass(variant string) string {
    baseClass := "px-4 py-2 rounded font-medium focus:outline-none focus:ring-2 focus:ring-opacity-75 transition-colors duration-200"
    
    switch variant {
    case "primary":
        return baseClass + " bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500"
    case "secondary":
        return baseClass + " bg-gray-600 hover:bg-gray-700 text-white focus:ring-gray-500"
    case "success":
        return baseClass + " bg-green-600 hover:bg-green-700 text-white focus:ring-green-500"
    case "danger":
        return baseClass + " bg-red-600 hover:bg-red-700 text-white focus:ring-red-500"
    case "warning":
        return baseClass + " bg-yellow-500 hover:bg-yellow-600 text-black focus:ring-yellow-400"
    case "info":
        return baseClass + " bg-cyan-600 hover:bg-cyan-700 text-white focus:ring-cyan-500"
    case "light":
        return baseClass + " bg-gray-100 hover:bg-gray-200 text-gray-900 focus:ring-gray-300"
    case "dark":
        return baseClass + " bg-gray-900 hover:bg-gray-800 text-white focus:ring-gray-700"
    default:
        return baseClass + " bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500"
    }
}
```

### Card Component

```go
func (t *TailwindTheme) Card(title string, content mi.H) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden mb-4"),
            func() mi.Node {
                if title != "" {
                    return b.Div(mi.Class("px-6 py-4 bg-gray-50 border-b border-gray-200"),
                        b.H3(mi.Class("text-lg font-semibold text-gray-900"), title),
                    )
                }
                return mi.NewFragment()
            }(),
            b.Div(mi.Class("px-6 py-4"),
                content(b),
            ),
        )
    }
}
```

### Form Components

```go
func (t *TailwindTheme) FormInput(label, name, inputType string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        id := "input_" + name
        inputAttrs := append([]mi.Attribute{
            mi.Class("mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"),
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

func (t *TailwindTheme) FormLabel(text, forField string) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Label(
            mi.Class("block text-sm font-medium text-gray-700"),
            mi.DataAttr("for", forField),
            text,
        )
    }
}
```

### Table Implementation

```go
func (t *TailwindTheme) Table(headers []string, rows [][]string) mi.H {
    return func(b *mi.Builder) mi.Node {
        // Create header nodes
        headerNodes := mintyex.Map(headers, func(header string) mi.Node {
            return b.Th(
                mi.Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"),
                header,
            )
        })
        
        // Create row nodes
        rowNodes := mintyex.Map(rows, func(row []string) mi.Node {
            cellNodes := mintyex.Map(row, func(cell string) mi.Node {
                return b.Td(
                    mi.Class("px-6 py-4 whitespace-nowrap text-sm text-gray-900"),
                    mi.RawHTML(cell),
                )
            })
            return b.Tr(mi.Class("even:bg-gray-50"), cellNodes...)
        })
        
        return b.Div(mi.Class("flex flex-col"),
            b.Div(mi.Class("-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8"),
                b.Div(mi.Class("py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"),
                    b.Div(mi.Class("shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"),
                        b.Table(mi.Class("min-w-full divide-y divide-gray-200"),
                            b.Thead(mi.Class("bg-gray-50"),
                                b.Tr(headerNodes...),
                            ),
                            b.Tbody(mi.Class("bg-white divide-y divide-gray-200"),
                                rowNodes...,
                            ),
                        ),
                    ),
                ),
            ),
        )
    }
}
```

---

## Creating Custom Themes

### Basic Custom Theme Structure

```go
// Custom corporate theme
type CorporateTheme struct {
    name         string
    version      string
    primaryColor string
    fontFamily   string
}

func NewCorporateTheme(primaryColor, fontFamily string) mui.Theme {
    return &CorporateTheme{
        name:         "Corporate",
        version:      "1.0.0",
        primaryColor: primaryColor,
        fontFamily:   fontFamily,
    }
}

func (t *CorporateTheme) GetName() string    { return t.name }
func (t *CorporateTheme) GetVersion() string { return t.version }
```

### Custom Button Implementation

```go
func (t *CorporateTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        style := t.getButtonStyle(variant)
        allAttrs := append([]mi.Attribute{
            mi.Style(style),
            mi.Class("corp-button"),
            mi.Type("button"),
        }, attrs...)
        
        return b.Button(allAttrs, text)
    }
}

func (t *CorporateTheme) getButtonStyle(variant string) string {
    baseStyle := fmt.Sprintf(`
        font-family: %s;
        padding: 12px 24px;
        border: none;
        border-radius: 4px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
    `, t.fontFamily)
    
    switch variant {
    case "primary":
        return baseStyle + fmt.Sprintf(`
            background-color: %s;
            color: white;
        `, t.primaryColor) + `
            &:hover { opacity: 0.9; transform: translateY(-1px); }
        `
    case "secondary":
        return baseStyle + fmt.Sprintf(`
            background-color: transparent;
            color: %s;
            border: 2px solid %s;
        `, t.primaryColor, t.primaryColor) + `
            &:hover { background-color: ` + t.primaryColor + `; color: white; }
        `
    default:
        return baseStyle + fmt.Sprintf(`
            background-color: %s;
            color: white;
        `, t.primaryColor)
    }
}
```

### CSS-in-Go Pattern for Custom Themes

```go
func (t *CorporateTheme) Card(title string, content mi.H) mi.H {
    return func(b *mi.Builder) mi.Node {
        cardStyle := fmt.Sprintf(`
            font-family: %s;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            overflow: hidden;
        `, t.fontFamily)
        
        headerStyle := fmt.Sprintf(`
            background: linear-gradient(135deg, %s, %s);
            color: white;
            padding: 20px;
            margin: 0;
            font-size: 1.25rem;
            font-weight: 600;
        `, t.primaryColor, t.adjustColor(t.primaryColor, -20))
        
        return b.Div(mi.Style(cardStyle),
            func() mi.Node {
                if title != "" {
                    return b.Div(mi.Style(headerStyle), title)
                }
                return mi.NewFragment()
            }(),
            b.Div(mi.Style("padding: 20px;"),
                content(b),
            ),
        )
    }
}

// Helper function to adjust color brightness
func (t *CorporateTheme) adjustColor(color string, adjustment int) string {
    // Implementation would parse hex color and adjust brightness
    // This is a simplified example
    return color
}
```

### Advanced Custom Theme with Configuration

```go
type ConfigurableTheme struct {
    config ThemeConfig
}

type ThemeConfig struct {
    Name            string            `json:"name"`
    PrimaryColor    string            `json:"primary_color"`
    SecondaryColor  string            `json:"secondary_color"`
    FontFamily      string            `json:"font_family"`
    BorderRadius    string            `json:"border_radius"`
    ButtonStyle     string            `json:"button_style"` // flat, raised, outlined
    CardStyle       string            `json:"card_style"`   // minimal, shadowed, bordered
    CustomCSS       map[string]string `json:"custom_css"`
}

func NewConfigurableTheme(config ThemeConfig) mui.Theme {
    return &ConfigurableTheme{config: config}
}

func (t *ConfigurableTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        style := t.buildButtonStyle(variant)
        class := fmt.Sprintf("theme-%s-button theme-button-%s", t.config.Name, variant)
        
        allAttrs := append([]mi.Attribute{
            mi.Style(style),
            mi.Class(class),
            mi.Type("button"),
        }, attrs...)
        
        return b.Button(allAttrs, text)
    }
}

func (t *ConfigurableTheme) buildButtonStyle(variant string) string {
    baseStyle := fmt.Sprintf(`
        font-family: %s;
        border-radius: %s;
        padding: 12px 24px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
        border: none;
    `, t.config.FontFamily, t.config.BorderRadius)
    
    switch t.config.ButtonStyle {
    case "flat":
        return baseStyle + t.getFlatButtonVariant(variant)
    case "raised":
        return baseStyle + t.getRaisedButtonVariant(variant) + `
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        `
    case "outlined":
        return baseStyle + t.getOutlinedButtonVariant(variant)
    default:
        return baseStyle + t.getFlatButtonVariant(variant)
    }
}
```

---

## Domain-Specific Theming

The Minty System includes domain-specific theming utilities for consistent styling across business domains.

### Domain CSS Classes

```go
// Semantic classes for domain-specific styling
func NewSemanticClasses(domain string) SemanticClasses {
    return SemanticClasses{domain: domain}
}

type SemanticClasses struct {
    domain string
}

func (sc SemanticClasses) Container(element string) string {
    return fmt.Sprintf("%s_%s_container", sc.domain, element)
}

func (sc SemanticClasses) Header(element string) string {
    return fmt.Sprintf("%s_%s_header", sc.domain, element)
}

func (sc SemanticClasses) Content(element string) string {
    return fmt.Sprintf("%s_%s_content", sc.domain, element)
}

func (sc SemanticClasses) Primary(element string) string {
    return fmt.Sprintf("%s_%s_primary", sc.domain, element)
}

func (sc SemanticClasses) Secondary(element string) string {
    return fmt.Sprintf("%s_%s_secondary", sc.domain, element)
}
```

### Domain-Specific Component Helpers

```go
// Domain card with semantic styling
func DomainCard(theme mui.Theme, domain, title string, content mi.H) mi.H {
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

// Domain button with semantic styling
func DomainButton(theme mui.Theme, domain, text, variant string, attrs ...mi.Attribute) mi.H {
    classes := mintyex.NewSemanticClasses(domain)
    
    return func(b *mi.Builder) mi.Node {
        var buttonClass string
        switch variant {
        case "primary":
            buttonClass = classes.Primary("button")
        case "secondary":
            buttonClass = classes.Secondary("button")
        default:
            buttonClass = classes.Secondary("button")
        }
        
        // Add domain class to existing attributes
        allAttrs := append([]mi.Attribute{mi.Class(buttonClass)}, attrs...)
        return b.Button(allAttrs, text)
    }
}
```

### Usage in Domain UIs

```go
// Finance domain using semantic classes
func AccountSummaryCard(theme mui.Theme, account mifi.Account) mi.H {
    displayData := mifi.PrepareAccountForDisplay(account)
    
    return mui.DomainCard(theme, "mifi", account.Name, func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("mifi_account_summary"),
            b.Div(mi.Class("mifi_account_info"),
                b.P(mi.Class("mifi_account_type"), displayData.TypeDisplay),
                b.Div(mi.Class("mifi_account_balance"),
                    b.Strong(displayData.FormattedBalance),
                ),
            ),
            b.Div(mi.Class("mifi_account_actions"),
                mui.DomainButton(theme, "mifi", "View Details", "secondary",
                    mi.Href("/accounts/"+account.ID)),
            ),
        )
    })
}
```

This generates CSS classes like:
- `mifi_card_container`
- `mifi_title_header`
- `mifi_content_body`
- `mifi_account_summary`
- `mifi_account_info`
- `mifi_secondary_button`

---

## Advanced Theme Patterns

### Theme Composition and Inheritance

```go
// Base theme providing common functionality
type BaseTheme struct {
    name    string
    version string
}

func (t *BaseTheme) GetName() string    { return t.name }
func (t *BaseTheme) GetVersion() string { return t.version }

// Common utility methods
func (t *BaseTheme) buildCommonStyle(baseStyle string, variant string, colorMap map[string]string) string {
    color, exists := colorMap[variant]
    if !exists {
        color = colorMap["default"]
    }
    return baseStyle + color
}

// Enhanced Bootstrap theme inheriting from base
type EnhancedBootstrapTheme struct {
    BaseTheme
    animations bool
    darkMode   bool
}

func NewEnhancedBootstrapTheme(animations, darkMode bool) mui.Theme {
    return &EnhancedBootstrapTheme{
        BaseTheme:  BaseTheme{name: "Enhanced Bootstrap", version: "2.0.0"},
        animations: animations,
        darkMode:   darkMode,
    }
}

func (t *EnhancedBootstrapTheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        class := t.getEnhancedButtonClass(variant)
        
        // Add animation classes if enabled
        if t.animations {
            class += " animate-on-hover"
        }
        
        // Add dark mode classes if enabled
        if t.darkMode {
            class += " dark-mode-button"
        }
        
        allAttrs := append([]mi.Attribute{mi.Class(class), mi.Type("button")}, attrs...)
        return b.Button(allAttrs, text)
    }
}
```

### Theme Provider Pattern

```go
// Theme provider for dependency injection
type ThemeProvider struct {
    themes map[string]mui.Theme
    default string
}

func NewThemeProvider() *ThemeProvider {
    provider := &ThemeProvider{
        themes: make(map[string]mui.Theme),
    }
    
    // Register built-in themes
    provider.Register("bootstrap", bootstrap.NewBootstrapTheme())
    provider.Register("tailwind", tailwind.NewTailwindTheme())
    provider.SetDefault("bootstrap")
    
    return provider
}

func (tp *ThemeProvider) Register(name string, theme mui.Theme) {
    tp.themes[name] = theme
}

func (tp *ThemeProvider) SetDefault(name string) {
    tp.default = name
}

func (tp *ThemeProvider) Get(name string) mui.Theme {
    if theme, exists := tp.themes[name]; exists {
        return theme
    }
    return tp.themes[tp.default]
}

func (tp *ThemeProvider) GetDefault() mui.Theme {
    return tp.themes[tp.default]
}

// Usage in application
type Application struct {
    themeProvider *ThemeProvider
}

func (app *Application) RenderPage(themeName string, content mi.H) mi.H {
    theme := app.themeProvider.Get(themeName)
    return theme.Container(content)
}
```

### Runtime Theme Switching

```go
// Theme switcher component
func ThemeSwitcher(currentTheme string, availableThemes []string) mi.H {
    return func(b *mi.Builder) mi.Node {
        options := mintyex.Map(availableThemes, func(themeName string) mui.SelectOption {
            return mui.SelectOption{
                Value:    themeName,
                Text:     strings.Title(themeName),
                Selected: themeName == currentTheme,
            }
        })
        
        return b.Form(mi.Action("/set-theme"), mi.Method("POST"),
            b.Label("Theme:"),
            b.Select(mi.Name("theme"), mi.OnChange("this.form.submit()"),
                mintyex.Map(options, func(opt mui.SelectOption) mi.Node {
                    var attrs []mi.Attribute
                    attrs = append(attrs, mi.Value(opt.Value))
                    if opt.Selected {
                        attrs = append(attrs, mi.Selected())
                    }
                    return b.Option(attrs, opt.Text)
                })...,
            ),
        )
    }
}

// Handler for theme switching
func (app *Application) SetThemeHandler(w http.ResponseWriter, r *http.Request) {
    themeName := r.FormValue("theme")
    
    // Validate theme exists
    theme := app.themeProvider.Get(themeName)
    if theme == nil {
        http.Error(w, "Invalid theme", http.StatusBadRequest)
        return
    }
    
    // Set theme cookie
    http.SetCookie(w, &http.Cookie{
        Name:   "theme",
        Value:  themeName,
        Path:   "/",
        MaxAge: 86400 * 30, // 30 days
    })
    
    // Redirect back to referrer
    referer := r.Header.Get("Referer")
    if referer == "" {
        referer = "/"
    }
    http.Redirect(w, r, referer, http.StatusSeeOther)
}
```

---

## Theme Integration Best Practices

### 1. Theme-Agnostic Component Design

Always design components to work with any theme:

```go
// Good: Theme-agnostic design
func UserProfile(theme mui.Theme, user User) mi.H {
    return theme.Card("User Profile", func(b *mi.Builder) mi.Node {
        return b.Div(
            b.P("Name: ", user.Name),
            b.P("Email: ", user.Email),
            theme.PrimaryButton("Edit Profile", mi.Href("/profile/edit")),
        )
    })
}

// Bad: Theme-specific assumptions
func UserProfile(user User) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("card"), // Assumes Bootstrap classes
            b.H3(mi.Class("card-title"), "User Profile"),
            b.Div(mi.Class("card-body"),
                b.P("Name: ", user.Name),
                b.Button(mi.Class("btn btn-primary"), "Edit"), // Bootstrap-specific
            ),
        )
    }
}
```

### 2. Consistent Theme Interface Usage

Use theme methods consistently rather than mixing with direct HTML:

```go
// Good: Consistent theme usage
func ProductCard(theme mui.Theme, product Product) mi.H {
    return theme.Card(product.Name, func(b *mi.Builder) mi.Node {
        return b.Div(
            b.P(product.Description),
            b.P(product.Price.Format()),
            theme.PrimaryButton("Buy Now", mi.Href("/buy/"+product.ID)),
        )
    })
}

// Bad: Mixing theme and direct HTML
func ProductCard(theme mui.Theme, product Product) mi.H {
    return theme.Card(product.Name, func(b *mi.Builder) mi.Node {
        return b.Div(
            b.P(product.Description),
            b.Div(mi.Class("price-display"), // Direct styling
                b.Strong(product.Price.Format()),
            ),
            b.Button(mi.Class("custom-button"), "Buy Now"), // Bypasses theme
        )
    })
}
```

### 3. Theme Fallback Handling

Always handle missing or invalid themes gracefully:

```go
func GetThemeOrDefault(themeProvider *ThemeProvider, requestedTheme string) mui.Theme {
    if requestedTheme == "" {
        return themeProvider.GetDefault()
    }
    
    theme := themeProvider.Get(requestedTheme)
    if theme == nil {
        log.Printf("Theme '%s' not found, using default", requestedTheme)
        return themeProvider.GetDefault()
    }
    
    return theme
}
```

### 4. Performance Considerations

Cache theme instances and avoid recreating them unnecessarily:

```go
// Theme instance caching
var (
    bootstrapThemeInstance mui.Theme
    tailwindThemeInstance  mui.Theme
    themeOnce              sync.Once
)

func GetBootstrapTheme() mui.Theme {
    themeOnce.Do(func() {
        bootstrapThemeInstance = bootstrap.NewBootstrapTheme()
        tailwindThemeInstance = tailwind.NewTailwindTheme()
    })
    return bootstrapThemeInstance
}
```

### 5. Testing Across Multiple Themes

Test components with different themes to ensure compatibility:

```go
func TestUserProfileWithMultipleThemes(t *testing.T) {
    user := User{Name: "Alice", Email: "alice@example.com"}
    
    themes := []mui.Theme{
        bootstrap.NewBootstrapTheme(),
        tailwind.NewTailwindTheme(),
        NewCustomTheme(),
    }
    
    for _, theme := range themes {
        t.Run(fmt.Sprintf("Theme_%s", theme.GetName()), func(t *testing.T) {
            component := UserProfile(theme, user)
            html := mi.Render(component)
            
            // Verify essential content regardless of theme
            assert.Contains(t, html, "Alice")
            assert.Contains(t, html, "alice@example.com")
            assert.Contains(t, html, "Edit Profile")
        })
    }
}
```

---

## Summary

The Minty System's theme architecture provides:

**Complete Separation**: Business logic remains independent of visual styling and CSS framework choices.

**Pluggable Design**: Themes can be swapped at runtime without code changes or application restarts.

**Comprehensive Interface**: Rich set of UI components covering forms, navigation, data display, and layout needs.

**Implementation Flexibility**: Themes can use CSS frameworks (Bootstrap, Tailwind) or custom CSS-in-Go approaches.

**Domain Integration**: Semantic class generation enables domain-specific styling while maintaining theme abstraction.

**Enterprise Features**: Theme providers, inheritance patterns, and runtime switching support complex application requirements.

The theme system enables applications to maintain consistent functionality while supporting diverse visual requirements, from simple Bootstrap styling to complex corporate branding and white-labeling scenarios.
