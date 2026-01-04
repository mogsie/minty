/*
Package mintydyn provides dynamic UI components with automatic pattern detection
for the Minty HTML generation library.

# Recommended Import

	import mdy "github.com/ha1tch/minty/mintydyn"

# Overview

mintydyn generates minimal, targeted JavaScript for client-side interactivity
while maintaining server-side rendering as the source of truth. Components are
built by declaring what data you have, and the system automatically selects
the optimal client-side implementation.

# Three Core Patterns

The package supports three primary patterns that can be combined:

  - States: Tab-like interfaces with pre-rendered or dynamic content
  - Data: Filterable datasets with client-side or server-side filtering
  - Rules: Dependency management for form fields and UI elements

Pattern detection is automatic based on what data you provide.

# Quick Start

Simple tabs:

	tabs := mdy.Tabs("profile", []mdy.ComponentState{
	    mdy.ActiveState("info", "Info", infoContent),
	    mdy.NewState("settings", "Settings", settingsContent),
	})

Filterable data:

	catalog := mdy.Filter("products", products, mdy.FilterSchema{
	    Fields: []mdy.FilterableField{
	        mdy.SelectField("category", "Category", categories),
	        mdy.RangeField("price", "Price", 0, 1000, 10),
	    },
	})

Form with dependencies:

	form := mdy.Form("insurance", []mdy.DependencyRule{
	    mdy.ShowWhen("marital-status", "equals", "married", "spouse-section"),
	})

# External Library Integration

mintydyn supports integration with external JavaScript libraries like
Google Maps, D3.js, Jitsi, etc.

	mapTabs := mdy.Dyn("location").
	    States(locationStates).
	    ExternalScript("https://maps.googleapis.com/maps/api/js?key=XXX", mdy.Required()).
	    External("map").
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

# Lifecycle Hooks

Components support lifecycle hooks for custom behavior:

  - BeforeInit: Runs before component initialization (can cancel)
  - OnInit (AfterInit): Runs after component initialization
  - BeforeStateChange: Runs before state changes (can cancel)
  - OnStateChange (AfterStateChange): Runs after state changes
  - OnState: Runs when a specific state becomes active
  - OnFilter (AfterFilter): Runs after filter changes
  - OnDestroy: Runs when component is destroyed

# Themes

Components can be styled using themes:

	// Use Bootstrap styling
	tabs := mdy.TabsWithTheme("nav", states, mdy.NewBootstrapDynamicTheme())

	// Use Tailwind styling
	tabs := mdy.TabsWithTheme("nav", states, mdy.NewTailwindDynamicTheme())

	// Use default semantic classes with included CSS
	tabs := mdy.WithDefaultCSS(mdy.Tabs("nav", states))

Custom themes can be created by implementing the DynamicTheme interface.

# CSS Builder

For custom styling, use the CSS builder:

	css := mdy.NewCSSBuilder().
	    Rule(".dyn-state-trigger",
	        mdy.Padding("1rem 2rem"),
	        mdy.BackgroundColor("#f0f0f0"),
	        mdy.BorderRadius("0.5rem"),
	    ).
	    Rule(".dyn-state-trigger.active",
	        mdy.BackgroundColor("#007bff"),
	        mdy.Color("white"),
	    ).
	    Render()

# Pattern Detection

The system automatically detects the optimal pattern based on:

  - Number of states (≤10 = pre-rendered, >10 = dynamic)
  - Data size (≤50 = client filtering, >50 = server filtering)
  - Combinations of states, data, and rules

Available patterns:

  - pre-rendered-states: All tab content rendered upfront
  - dynamic-states: Tab content loaded on demand
  - client-filterable: Data filtered in browser
  - server-filterable: Data filtered via server requests
  - dependency-only: Just form field dependencies
  - stateful-data: Tabs with data context per tab
  - filterable-states: Filtering with state context
  - dependent-states: Tabs controlled by rules
  - dependent-data: Filters controlled by rules
  - complete: All three patterns combined

# Accessing Components from JavaScript

Generated components are accessible globally:

	// Access component instance
	const comp = window.DynComponent_myid;

	// Listen to events
	comp.on('state:change', (e) => console.log(e.detail));

	// Access external objects
	const map = comp.getExternal('map');

	// Register external objects
	comp.registerExternal('chart', myD3Chart);

	// Trigger programmatic actions
	comp.switchToState('settings');

	// Cleanup
	comp.destroy();

# Generated JavaScript

The package generates minimal JavaScript specific to your component's needs.
A simple tabs component might generate ~2KB of JS, while a complex component
with all patterns might generate ~8KB. This is dramatically smaller than
framework-based approaches.

The generated code includes:

  - Base component class with event handling
  - Pattern-specific managers (States, Data, Rules)
  - Coordination logic between managers
  - External script loading and lifecycle hooks
*/
package mintydyn
