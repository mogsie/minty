package minty

import (
	"fmt"
	"strings"
)

// escapeJSString escapes a string for safe inclusion in JavaScript single-quoted strings.
func escapeJSString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

// DarkMode provides framework-agnostic dark mode support for minty applications.
// It generates the necessary JavaScript for theme detection, persistence, and toggling.
//
// Supported frameworks:
//   - Tailwind CSS (class="dark" on html element)
//   - Bootstrap 5.3+ (data-bs-theme="dark" attribute)
//   - Custom CSS variable schemes (any data-* attribute)
//
// Not supported (by design):
//   - MUI, Vuetify, Angular Material (require framework-specific state management)
//
// See DARKMODE.md for architectural rationale.
type DarkMode struct {
	config DarkModeConfig
}

// DarkModeConfig holds configuration for dark mode behaviour.
type DarkModeConfig struct {
	// Toggle mechanism
	UseClass   bool   // true = toggle class, false = toggle attribute
	ClassName  string // Class name to toggle (e.g., "dark")
	AttrName   string // Attribute name (e.g., "data-bs-theme")
	LightValue string // Attribute value for light mode
	DarkValue  string // Attribute value for dark mode

	// Persistence
	StorageKey string // localStorage key, default "theme"
	Default    string // "system", "light", or "dark"

	// UI
	LightIcon string // Icon shown when in light mode (clicking goes dark)
	DarkIcon  string // Icon shown when in dark mode (clicking goes light)
	IconID    string // ID for the icon element

	// Output
	Minify bool // Minify generated JavaScript
}

// DarkModeOption configures a DarkMode instance.
type DarkModeOption func(*DarkModeConfig)

// DarkModeStorage sets the localStorage key for persistence.
func DarkModeStorage(key string) DarkModeOption {
	return func(c *DarkModeConfig) {
		c.StorageKey = key
	}
}

// DarkModeDefault sets the default theme when no preference is stored.
// Valid values: "system" (follows OS preference), "light", "dark"
func DarkModeDefault(d string) DarkModeOption {
	return func(c *DarkModeConfig) {
		c.Default = d
	}
}

// DarkModeIcons sets the icons for the toggle button.
// lightIcon is shown when in light mode (clicking switches to dark).
// darkIcon is shown when in dark mode (clicking switches to light).
func DarkModeIcons(lightIcon, darkIcon string) DarkModeOption {
	return func(c *DarkModeConfig) {
		c.LightIcon = lightIcon
		c.DarkIcon = darkIcon
	}
}

// DarkModeIconID sets the ID for the icon element inside the toggle button.
func DarkModeIconID(id string) DarkModeOption {
	return func(c *DarkModeConfig) {
		c.IconID = id
	}
}

// DarkModeMinify enables JavaScript minification.
func DarkModeMinify() DarkModeOption {
	return func(c *DarkModeConfig) {
		c.Minify = true
	}
}

// DarkModeSVGIcons uses Heroicon SVG icons instead of emoji.
// This provides better cross-platform consistency.
func DarkModeSVGIcons() DarkModeOption {
	// Moon icon (shown in light mode - click to go dark)
	moonSVG := `<svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.72 9.72 0 0 1 18 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 0 0 3 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 0 0 9.002-5.998Z"/></svg>`
	// Sun icon (shown in dark mode - click to go light)  
	sunSVG := `<svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386-1.591 1.591M21 12h-2.25m-.386 6.364-1.591-1.591M12 18.75V21m-4.773-4.227-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"/></svg>`
	return func(c *DarkModeConfig) {
		c.LightIcon = moonSVG
		c.DarkIcon = sunSVG
	}
}

// NewDarkMode creates a dark mode handler with Tailwind defaults.
// Use DarkModeTailwind(), DarkModeBootstrap(), or DarkModeAttr() for
// framework-specific presets.
func NewDarkMode(opts ...DarkModeOption) *DarkMode {
	// Tailwind defaults
	config := DarkModeConfig{
		UseClass:   true,
		ClassName:  "dark",
		StorageKey: "theme",
		Default:    "system",
		LightIcon:  "🌙",
		DarkIcon:   "☀️",
		IconID:     "dark-mode-icon",
		Minify:     false,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &DarkMode{config: config}
}

// DarkModeTailwind creates a dark mode handler configured for Tailwind CSS.
// Tailwind uses class="dark" on the html element.
func DarkModeTailwind(opts ...DarkModeOption) *DarkMode {
	defaults := []DarkModeOption{
		func(c *DarkModeConfig) {
			c.UseClass = true
			c.ClassName = "dark"
		},
	}
	return NewDarkMode(append(defaults, opts...)...)
}

// DarkModeBootstrap creates a dark mode handler configured for Bootstrap 5.3+.
// Bootstrap uses data-bs-theme="dark" attribute on the html element.
func DarkModeBootstrap(opts ...DarkModeOption) *DarkMode {
	defaults := []DarkModeOption{
		func(c *DarkModeConfig) {
			c.UseClass = false
			c.AttrName = "data-bs-theme"
			c.LightValue = "light"
			c.DarkValue = "dark"
		},
	}
	return NewDarkMode(append(defaults, opts...)...)
}

// DarkModeAttr creates a dark mode handler using a custom attribute.
// Use this for CSS variable schemes that respond to data attributes.
//
// Example for custom CSS:
//
//	:root[data-theme="dark"] { --bg: #1a1a1a; --text: #fff; }
//	:root[data-theme="light"] { --bg: #fff; --text: #1a1a1a; }
//
//	darkMode := mi.DarkModeAttr("data-theme", "light", "dark")
func DarkModeAttr(attrName, lightValue, darkValue string, opts ...DarkModeOption) *DarkMode {
	defaults := []DarkModeOption{
		func(c *DarkModeConfig) {
			c.UseClass = false
			c.AttrName = attrName
			c.LightValue = lightValue
			c.DarkValue = darkValue
		},
	}
	return NewDarkMode(append(defaults, opts...)...)
}

// Script generates the JavaScript for dark mode initialisation.
// Place this in the <head> to prevent flash of wrong theme on page load.
//
// Example:
//
//	b.Head(
//	    b.Script(mi.Raw(`tailwind.config = { darkMode: 'class' }`)),
//	    darkMode.Script(b),
//	)
func (dm *DarkMode) Script(b *Builder) Node {
	js := dm.generateScript()
	if dm.config.Minify {
		js = minifyDarkModeJS(js)
	}
	return b.Script(Raw(js))
}

// ScriptRaw returns the raw JavaScript string without wrapping in a script tag.
// Useful when combining with other scripts.
func (dm *DarkMode) ScriptRaw() string {
	js := dm.generateScript()
	if dm.config.Minify {
		js = minifyDarkModeJS(js)
	}
	return js
}

// Toggle generates a button element that toggles dark mode.
// Pass additional attributes to style the button.
//
// Example:
//
//	darkMode.Toggle(b,
//	    mi.Class("p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"),
//	    mi.Attr("title", "Toggle dark mode"),
//	)
func (dm *DarkMode) Toggle(b *Builder, attrs ...interface{}) Node {
	buttonAttrs := []interface{}{
		Type("button"),
		Attr("onclick", dm.toggleFunctionName()+"()"),
		Attr("aria-label", "Toggle dark mode"),
	}
	buttonAttrs = append(buttonAttrs, attrs...)

	// Icon span
	icon := b.Span(
		ID(dm.config.IconID),
		Attr("aria-hidden", "true"),
		Raw(dm.config.LightIcon), // Will be updated by JS on load
	)

	buttonAttrs = append(buttonAttrs, icon)
	return b.Button(buttonAttrs...)
}

// ToggleID returns the ID of the icon element, useful for custom toggle implementations.
func (dm *DarkMode) ToggleID() string {
	return dm.config.IconID
}

func (dm *DarkMode) toggleFunctionName() string {
	return "toggleDarkMode"
}

func (dm *DarkMode) updateIconFunctionName() string {
	return "updateDarkModeIcon"
}

func (dm *DarkMode) generateScript() string {
	c := dm.config
	toggleFn := dm.toggleFunctionName()
	updateFn := dm.updateIconFunctionName()

	var sb strings.Builder

	// Toggle function
	sb.WriteString(fmt.Sprintf("function %s() {\n", toggleFn))
	sb.WriteString("    const html = document.documentElement;\n")

	if c.UseClass {
		sb.WriteString(fmt.Sprintf("    const isDark = html.classList.toggle('%s');\n", c.ClassName))
	} else {
		sb.WriteString(fmt.Sprintf("    const isDark = html.getAttribute('%s') !== '%s';\n", c.AttrName, c.DarkValue))
		sb.WriteString(fmt.Sprintf("    html.setAttribute('%s', isDark ? '%s' : '%s');\n", c.AttrName, c.DarkValue, c.LightValue))
	}

	sb.WriteString(fmt.Sprintf("    localStorage.setItem('%s', isDark ? 'dark' : 'light');\n", c.StorageKey))
	sb.WriteString(fmt.Sprintf("    %s(isDark);\n", updateFn))
	sb.WriteString("}\n\n")

	// Update icon function
	sb.WriteString(fmt.Sprintf("function %s(isDark) {\n", updateFn))
	sb.WriteString(fmt.Sprintf("    const icon = document.getElementById('%s');\n", c.IconID))
	sb.WriteString(fmt.Sprintf("    if (icon) icon.innerHTML = isDark ? '%s' : '%s';\n", escapeJSString(c.DarkIcon), escapeJSString(c.LightIcon)))
	sb.WriteString("}\n\n")

	// IMMEDIATE initialization (runs synchronously before body renders)
	// This prevents flash of wrong theme on page load/navigation
	sb.WriteString("(function() {\n")
	sb.WriteString(fmt.Sprintf("    const saved = localStorage.getItem('%s');\n", c.StorageKey))
	sb.WriteString("    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;\n")

	// Determine initial state based on default setting
	switch c.Default {
	case "dark":
		sb.WriteString("    const isDark = saved !== 'light';\n")
	case "light":
		sb.WriteString("    const isDark = saved === 'dark';\n")
	default: // "system"
		sb.WriteString("    const isDark = saved === 'dark' || (saved === null && prefersDark);\n")
	}

	// Apply initial state IMMEDIATELY
	if c.UseClass {
		sb.WriteString(fmt.Sprintf("    if (isDark) document.documentElement.classList.add('%s');\n", c.ClassName))
		sb.WriteString(fmt.Sprintf("    else document.documentElement.classList.remove('%s');\n", c.ClassName))
	} else {
		sb.WriteString(fmt.Sprintf("    document.documentElement.setAttribute('%s', isDark ? '%s' : '%s');\n", c.AttrName, c.DarkValue, c.LightValue))
	}

	// Store isDark for icon update after DOM loads
	sb.WriteString("    window.__darkModeInit = isDark;\n")
	sb.WriteString("})();\n\n")

	// Update icon after DOM is ready (icon element needs to exist)
	sb.WriteString("document.addEventListener('DOMContentLoaded', function() {\n")
	sb.WriteString(fmt.Sprintf("    %s(window.__darkModeInit);\n", updateFn))
	sb.WriteString("});\n")

	return sb.String()
}

// minifyDarkModeJS performs basic minification on the dark mode JS.
func minifyDarkModeJS(js string) string {
	// Remove comments
	lines := strings.Split(js, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Remove leading whitespace, keep the code
		result = append(result, trimmed)
	}
	return strings.Join(result, "\n")
}
