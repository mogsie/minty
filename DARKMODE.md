# Dark Mode Support in Minty

## Overview

Minty provides built-in dark mode support through the `DarkMode` type. This document explains the design decisions, supported frameworks, and architectural rationale.

## Quick Start

```go
import mi "github.com/ha1tch/minty"

// Create dark mode handler (Tailwind by default)
darkMode := mi.NewDarkMode()

// In your page layout:
func layout(b *mi.Builder) mi.Node {
    return b.Html(mi.Lang("en"),
        b.Head(
            b.Script(mi.Raw(`tailwind.config = { darkMode: 'class' }`)),
            darkMode.Script(b),  // Generates init + toggle functions
        ),
        b.Body(mi.Class("bg-white dark:bg-gray-900"),
            // Header with toggle button
            b.Header(
                darkMode.Toggle(b, 
                    mi.Class("p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700"),
                    mi.Attr("title", "Toggle dark mode"),
                ),
            ),
            // ... rest of page
        ),
    )
}
```

## Supported Frameworks

### Tailwind CSS (Default)

Tailwind uses a class on the `<html>` element:

```go
darkMode := mi.DarkModeTailwind()
// or simply:
darkMode := mi.NewDarkMode()
```

Generated behaviour:
- Adds/removes `class="dark"` on `<html>`
- CSS responds via `.dark` ancestor selectors

Required Tailwind config:
```javascript
tailwind.config = { darkMode: 'class' }
```

### Bootstrap 5.3+

Bootstrap 5.3 introduced native dark mode via data attributes:

```go
darkMode := mi.DarkModeBootstrap()
```

Generated behaviour:
- Sets `data-bs-theme="light"` or `data-bs-theme="dark"` on `<html>`
- CSS responds via `[data-bs-theme="dark"]` selectors

### Custom CSS Variables

For custom CSS or other frameworks using data attributes:

```go
darkMode := mi.DarkModeAttr("data-theme", "light", "dark")
```

Your CSS:
```css
:root[data-theme="light"] {
    --bg-color: #ffffff;
    --text-color: #1a1a1a;
}

:root[data-theme="dark"] {
    --bg-color: #1a1a1a;
    --text-color: #ffffff;
}

body {
    background: var(--bg-color);
    color: var(--text-color);
}
```

## Not Supported (By Design)

The following frameworks are **not supported**:

| Framework | Reason |
|-----------|--------|
| MUI (Material-UI) | Requires React ThemeProvider context |
| Vuetify | Requires Vue reactivity system (`$vuetify.theme.dark`) |
| Angular Material | Requires Angular service injection |
| Material Components Web | Complex JS-driven CSS variable swapping |

### Architectural Rationale

These frameworks couple theme state to their JavaScript framework's reactivity system. The theme doesn't live in the DOM—it lives in React/Vue/Angular state, which then updates the DOM.

**Minty's model is fundamentally different:**

1. Minty generates HTML on the server
2. Dark mode is a CSS concern, triggered by DOM attributes
3. A small script handles toggle + persistence
4. No framework runtime required

If you're using MUI, Vuetify, or Angular Material, you're writing React/Vue/Angular code, not Go templates. You'd use those frameworks' native theme systems, not minty.

**Minty's natural companions are CSS-first frameworks:**
- Tailwind CSS
- Bootstrap (via CDN)
- Vanilla CSS with variables
- Lightweight CSS frameworks (Pico, Water.css, etc.)

## Configuration Options

```go
darkMode := mi.NewDarkMode(
    mi.DarkModeStorage("theme"),        // localStorage key (default: "theme")
    mi.DarkModeDefault("system"),       // "system", "light", or "dark"
    mi.DarkModeIcons("🌙", "☀️"),       // light-mode icon, dark-mode icon
    mi.DarkModeIconID("theme-icon"),    // ID for icon element
    mi.DarkModeMinify(),                // Minify generated JS
)
```

### Default Behaviour

| Option | Default | Description |
|--------|---------|-------------|
| Storage key | `"theme"` | localStorage key for persistence |
| Default | `"system"` | Follow OS preference when no saved value |
| Light icon | `"🌙"` | Shown when in light mode |
| Dark icon | `"☀️"` | Shown when in dark mode |
| Icon ID | `"dark-mode-icon"` | ID of the icon span element |

### Default Modes

- `"system"` — Follows OS preference (`prefers-color-scheme`), respects saved choice
- `"light"` — Defaults to light, only goes dark if explicitly saved
- `"dark"` — Defaults to dark, only goes light if explicitly saved

## API Reference

### Constructors

```go
// Tailwind (default)
mi.NewDarkMode(opts ...DarkModeOption) *DarkMode

// Framework-specific
mi.DarkModeTailwind(opts ...DarkModeOption) *DarkMode
mi.DarkModeBootstrap(opts ...DarkModeOption) *DarkMode
mi.DarkModeAttr(attrName, lightValue, darkValue string, opts ...DarkModeOption) *DarkMode
```

### Methods

```go
// Generate <script> tag with init + toggle functions
// Place in <head> to prevent flash of wrong theme
darkMode.Script(b *Builder) Node

// Get raw JS string (for combining with other scripts)
darkMode.ScriptRaw() string

// Generate toggle button with icon
darkMode.Toggle(b *Builder, attrs ...interface{}) Node

// Get icon element ID (for custom implementations)
darkMode.ToggleID() string
```

## Generated JavaScript

The `Script()` method generates:

1. **`toggleDarkMode()`** — Called by the toggle button
2. **`updateDarkModeIcon(isDark)`** — Updates the icon display
3. **DOMContentLoaded handler** — Initialises theme on page load

Example output (Tailwind, unminified):

```javascript
function toggleDarkMode() {
    const html = document.documentElement;
    const isDark = html.classList.toggle('dark');
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
    updateDarkModeIcon(isDark);
}

function updateDarkModeIcon(isDark) {
    const icon = document.getElementById('dark-mode-icon');
    if (icon) icon.textContent = isDark ? '☀️' : '🌙';
}

document.addEventListener('DOMContentLoaded', function() {
    const saved = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const isDark = saved === 'dark' || (saved === null && prefersDark);
    if (isDark) html.classList.add('dark');
    else html.classList.remove('dark');
    updateDarkModeIcon(isDark);
});
```

## Flash Prevention

The script runs on `DOMContentLoaded`, which is early but not instant. For zero-flash, you can inline a minimal script in `<head>` before stylesheets:

```go
b.Head(
    // Inline flash prevention (runs immediately)
    b.Script(mi.Raw(`
        (function() {
            const saved = localStorage.getItem('theme');
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            if (saved === 'dark' || (saved === null && prefersDark)) {
                document.documentElement.classList.add('dark');
            }
        })();
    `)),
    // Tailwind config
    b.Script(mi.Raw(`tailwind.config = { darkMode: 'class' }`)),
    // Full dark mode script (handles toggle + icon)
    darkMode.Script(b),
)
```

## Why Not mintydyn?

Dark mode is intentionally in minty core, not mintydyn, because:

1. **mintydyn is for component-level interactivity** — States, Data filtering, Rules/dependencies within a component container

2. **Dark mode is document-level** — It affects `<html>`, not a component

3. **Dark mode is HTML generation** — A script tag and a button, both just HTML elements

4. **No component state** — Dark mode state lives in localStorage and the DOM, not in a mintydyn component

mintydyn patterns (States/Data/Rules) don't fit document-level concerns. Keeping dark mode in minty core maintains clean architectural separation.

## Complete Example

```go
package main

import (
    "net/http"
    mi "github.com/ha1tch/minty"
)

func main() {
    darkMode := mi.DarkModeTailwind(
        mi.DarkModeMinify(),
    )

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        b := mi.NewBuilder()
        page := b.Html(mi.Lang("en"),
            b.Head(
                b.Title("My App"),
                b.Script(mi.Raw(`tailwind.config = { darkMode: 'class' }`)),
                b.Script(mi.Src("https://cdn.tailwindcss.com")),
                darkMode.Script(b),
            ),
            b.Body(mi.Class("bg-gray-100 dark:bg-gray-900 min-h-screen transition-colors"),
                b.Header(mi.Class("p-4 flex justify-between items-center"),
                    b.H1(mi.Class("text-xl font-bold text-gray-900 dark:text-white"), "My App"),
                    darkMode.Toggle(b,
                        mi.Class("p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"),
                        mi.Attr("title", "Toggle dark mode"),
                    ),
                ),
                b.Main(mi.Class("p-4"),
                    b.P(mi.Class("text-gray-700 dark:text-gray-300"),
                        "This page supports dark mode!",
                    ),
                ),
            ),
        )
        mi.Render(w, page)
    })

    http.ListenAndServe(":8080", nil)
}
```
