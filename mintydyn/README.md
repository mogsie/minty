# mintydyn

Dynamic UI components with automatic pattern detection for the Minty HTML generation library.

## Installation

```go
import mdy "github.com/ha1tch/minty/mintydyn"
```

Recommended alias: `mdy`

## Overview

mintydyn generates minimal, targeted JavaScript for client-side interactivity while maintaining server-side rendering as the source of truth. Components are built by declaring what data you have, and the system automatically selects the optimal client-side implementation.

## Quick Start

### Simple Tabs

```go
tabs := mdy.Tabs("profile", []mdy.ComponentState{
    mdy.ActiveState("info", "Info", infoContent),
    mdy.NewState("settings", "Settings", settingsContent),
})
```

### Filterable Data

```go
catalog := mdy.Filter("products", products, mdy.FilterSchema{
    Fields: []mdy.FilterableField{
        mdy.SelectField("category", "Category", categories),
        mdy.RangeField("price", "Price", 0, 1000, 10),
    },
})
```

### Form Dependencies

```go
form := mdy.Form("insurance", []mdy.DependencyRule{
    mdy.ShowWhen("marital-status", "equals", "married", "spouse-section"),
})
```

## External Library Integration

Integrate with Google Maps, D3.js, Jitsi, or any external JavaScript library:

```go
mapTabs := mdy.Dyn("location").
    States(locationStates).
    ExternalScript("https://maps.googleapis.com/maps/api/js?key=XXX", mdy.Required()).
    External("map").  // Reserve name in externals registry
    OnInit(`
        this.registerExternal('map', new google.maps.Map(
            document.getElementById('map-container'),
            { center: {lat: -34.397, lng: 150.644}, zoom: 8 }
        ));
    `).
    OnState("map-view", `
        google.maps.event.trigger(this.getExternal('map'), 'resize');
    `).
    OnDestroy(`
        // Cleanup map resources
    `).
    Build()
```

## Lifecycle Hooks

| Hook | Description |
|------|-------------|
| `BeforeInit(js)` | Runs before init, return `false` to cancel |
| `OnInit(js)` | Runs after initialization |
| `BeforeStateChange(js)` | Runs before state change, return `false` to cancel |
| `OnStateChange(js)` | Runs after state change |
| `OnState(id, js)` | Runs when specific state becomes active |
| `OnFilter(js)` | Runs after filter changes |
| `OnDestroy(js)` | Runs on cleanup |

## Themes

```go
// Bootstrap 5
tabs := mdy.TabsWithTheme("nav", states, mdy.NewBootstrapDynamicTheme())

// Tailwind CSS
tabs := mdy.TabsWithTheme("nav", states, mdy.NewTailwindDynamicTheme())

// Default with included CSS
tabs := mdy.WithDefaultCSS(mdy.Tabs("nav", states))
```

## Pattern Detection

The system automatically detects patterns based on data:

| Condition | Pattern | Description |
|-----------|---------|-------------|
| States ≤10 | `pre-rendered-states` | All content rendered upfront |
| States >10 | `dynamic-states` | Content loaded on demand |
| Data ≤50 | `client-filterable` | Browser-side filtering |
| Data >50 | `server-filterable` | Server-side filtering |
| Rules only | `dependency-only` | Form field dependencies |
| States + Data | `stateful-data` | Data context per state |
| States + Rules | `dependent-states` | Rule-controlled tabs |
| Data + Rules | `dependent-data` | Rule-controlled filters |
| All three | `complete` | Full coordination |

## JavaScript Access

```javascript
// Access component
const comp = window.DynComponent_myid;

// Events
comp.on('state:change', (e) => console.log(e.detail));
comp.on('data:filtered', (e) => console.log(e.detail.resultCount));

// External objects
const map = comp.getExternal('map');
comp.registerExternal('chart', myD3Chart);

// Actions
comp.switchToState('settings');
comp.destroy();
```

## CSS Builder

```go
css := mdy.NewCSSBuilder().
    Rule(".my-tabs",
        mdy.Display("flex"),
        mdy.Gap("1rem"),
        mdy.BorderBottom("2px solid #ccc"),
    ).
    Render()
```

## Bundle Size

| Component Type | Approx. JS Size |
|----------------|-----------------|
| Simple tabs | ~2KB |
| Filterable data | ~3KB |
| With rules | ~4KB |
| Complete (all patterns) | ~8KB |
| 20 components (shared managers) | ~10KB total |

Compare to React SPA: 200-400KB initial bundle.

## Files

| File | Lines | Purpose |
|------|-------|---------|
| `mintydyn.go` | 274 | Core types, constraints |
| `builder.go` | 428 | Generic builder, pattern detection |
| `generate.go` | 506 | HTML structure generation |
| `combined.go` | 375 | Combined pattern structures |
| `javascript.go` | 1021 | JS generation with hooks |
| `convenience.go` | 460 | Helper functions |
| `theme.go` | 190 | Theme interface and implementations |
| `css.go` | 290 | CSS builder utilities |
| `doc.go` | 173 | Package documentation |
| **Total** | **3717** | |

## License

Same as minty (to be determined).
