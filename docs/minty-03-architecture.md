# Minty Documentation - Part 3
## Core Architecture & Type System: Minty's Foundation

> **Part of the Minty System**: This document covers the foundational architecture and type system that underlies the entire Minty System. The Node interface, Builder pattern, and rendering architecture described here form the foundation for HTML generation (minty), business domain presentation (mintyfinui, mintymoveui, mintycartui), theme implementations (bootstrap, tailwind), and iterator-based rendering patterns. Understanding this foundation is essential for working with any component of the Minty System.

---

### Table of Contents
1. [The Node Interface and Its Implementations](#the-node-interface-and-its-implementations)
2. [Builder Pattern Design](#builder-pattern-design)
3. [Template Type System](#template-type-system)
4. [HTML Element Generation and Attributes](#html-element-generation-and-attributes)
5. [Text Handling and Automatic Escaping](#text-handling-and-automatic-escaping)
6. [Memory Management and Performance Characteristics](#memory-management-and-performance-characteristics)

---

## The Node Interface and Its Implementations

At the heart of Minty's architecture lies the `Node` interface, a simple but powerful abstraction that represents any piece of HTML content. This interface serves as the foundation for Minty's type safety, composability, and rendering efficiency.

### The Core Node Interface

```go
// Node represents any HTML content that can be rendered
type Node interface {
    Render(w io.Writer) error
}
```

This deliberately minimal interface encapsulates the fundamental operation of HTML generation: writing content to an output stream. By keeping the interface small and focused, Minty ensures that different implementations can be easily composed and that the rendering process remains efficient and predictable.

### Node Implementation Types

Minty provides several concrete implementations of the Node interface, each optimized for specific use cases while maintaining the same simple contract:

#### Element Nodes

Element nodes represent HTML tags with optional attributes and children:

```go
// Element represents an HTML tag with attributes and children
type Element struct {
    Tag        string
    Attributes map[string]string
    Children   []Node
    SelfClosing bool
}

func (e *Element) Render(w io.Writer) error {
    // Write opening tag
    if _, err := w.Write([]byte("<" + e.Tag)); err != nil {
        return err
    }
    
    // Write attributes
    for key, value := range e.Attributes {
        if _, err := fmt.Fprintf(w, ` %s="%s"`, key, html.EscapeString(value)); err != nil {
            return err
        }
    }
    
    if e.SelfClosing {
        _, err := w.Write([]byte(" />"))
        return err
    }
    
    if _, err := w.Write([]byte(">")); err != nil {
        return err
    }
    
    // Render children
    for _, child := range e.Children {
        if err := child.Render(w); err != nil {
            return err
        }
    }
    
    // Write closing tag
    _, err := fmt.Fprintf(w, "</%s>", e.Tag)
    return err
}
```

This implementation demonstrates several key architectural decisions. The rendering is streaming-based, writing directly to the output without building intermediate strings. Attributes are automatically escaped, providing security by default. The error handling is explicit and composable, allowing higher-level code to handle rendering failures appropriately.

#### Text Nodes

Text nodes represent plain text content with automatic HTML escaping:

```go
// TextNode represents escaped text content
type TextNode struct {
    Content string
}

func (t *TextNode) Render(w io.Writer) error {
    escaped := html.EscapeString(t.Content)
    _, err := w.Write([]byte(escaped))
    return err
}
```

Text nodes embody Minty's security-first philosophy by automatically escaping content. This makes XSS vulnerabilities much harder to introduce accidentally, since the safe behavior is the default behavior.

#### Raw Nodes

For cases where unescaped HTML is explicitly needed, Minty provides raw nodes:

```go
// RawNode represents unescaped HTML content (use with caution)
type RawNode struct {
    Content string
}

func (r *RawNode) Render(w io.Writer) error {
    _, err := w.Write([]byte(r.Content))
    return err
}
```

Raw nodes require explicit creation through a dedicated function, making it clear when potentially unsafe content is being included:

```go
// Explicit opt-in to unescaped content
func Raw(content string) Node {
    return &RawNode{Content: content}
}

// Usage requires deliberate choice
b.Div(
    b.P("This is safe text"),           // Automatically escaped
    minty.Raw("<em>This is raw HTML</em>"), // Explicitly unescaped
)
```

#### Fragment Nodes

Fragment nodes allow grouping multiple nodes without introducing wrapper elements:

```go
// Fragment represents a collection of nodes without a wrapper element
type Fragment struct {
    Children []Node
}

func (f *Fragment) Render(w io.Writer) error {
    for _, child := range f.Children {
        if err := child.Render(w); err != nil {
            return err
        }
    }
    return nil
}
```

Fragments are essential for conditional rendering and component composition where wrapper elements would interfere with styling or semantics:

```go
func ConditionalContent(showExtra bool) Node {
    nodes := []Node{
        b.H1("Always shown"),
    }
    
    if showExtra {
        nodes = append(nodes, 
            b.P("Extra content"),
            b.P("More extra content"),
        )
    }
    
    return minty.Fragment(nodes...)
}
```

### Node Composition Patterns

The Node interface enables powerful composition patterns that form the basis of Minty's component system:

#### Functional Composition

Nodes can be composed through pure functions, enabling reusable components:

```go
// Components are functions that return Nodes
func UserCard(user User) Node {
    return b.Div(minty.Class("user-card"),
        b.H3(user.Name),
        b.P(user.Email),
        b.Img(minty.Src(user.Avatar), minty.Alt("User avatar")),
    )
}

func UserList(users []User) Node {
    var userNodes []Node
    for _, user := range users {
        userNodes = append(userNodes, UserCard(user))
    }
    
    return b.Div(minty.Class("user-list"), userNodes...)
}
```

This functional approach leverages Go's strengths in composition and type safety while avoiding the complexity of object-oriented component hierarchies.

#### Higher-Order Components

The Node interface supports higher-order components that wrap or modify other nodes:

```go
// Higher-order component that adds a wrapper
func WithContainer(content Node) Node {
    return b.Div(minty.Class("container"),
        b.Div(minty.Class("row"),
            b.Div(minty.Class("col"),
                content,
            ),
        ),
    )
}

// Higher-order component that adds error boundaries
func WithErrorBoundary(content Node, fallback Node) Node {
    // In practice, this might check for rendering errors
    return content
}

// Usage
page := WithContainer(
    WithErrorBoundary(
        UserDashboard(users),
        b.P("Failed to load dashboard"),
    ),
)
```

#### Conditional Rendering

The Node interface makes conditional rendering natural and type-safe:

```go
func ConditionalButton(user User) Node {
    if user.IsAdmin {
        return b.Button(minty.Class("btn btn-danger"), "Admin Panel")
    }
    return minty.Fragment() // Empty fragment for no content
}

// Helper function for common conditional patterns
func If(condition bool, node Node) Node {
    if condition {
        return node
    }
    return minty.Fragment()
}

// Usage
b.Div(
    b.H1("Welcome"),
    If(user.IsLoggedIn, 
        b.P("Hello, " + user.Name),
    ),
    If(!user.IsLoggedIn,
        b.A(minty.Href("/login"), "Please log in"),
    ),
)
```

### Interface Benefits and Design Rationale

The Node interface's simplicity provides several architectural benefits that align with Minty's core principles:

**Composition Over Inheritance**: Rather than complex component hierarchies, Minty uses function composition. This aligns with Go's composition-over-inheritance philosophy and makes component behavior more predictable and testable.

**Streaming Rendering**: The `io.Writer` interface enables streaming rendering, which is more memory-efficient than building complete HTML strings in memory. This becomes particularly important for large pages or high-traffic applications.

**Error Handling**: By returning errors from the Render method, Node implementations can handle and propagate rendering errors appropriately. This fits Go's explicit error handling patterns and makes debugging easier.

**Type Safety**: All Node implementations are checked at compile time, eliminating the runtime template errors that plague string-based templating systems.

**Performance Predictability**: The interface contract makes performance characteristics predictable - rendering cost is proportional to content size, with no hidden allocations or surprising overhead.

The Node interface demonstrates how a simple, well-designed abstraction can provide the foundation for a powerful and flexible system while maintaining the simplicity and predictability that Go developers value.

---

## Builder Pattern Design

Minty's builder pattern serves as the primary interface between developers and the HTML generation system. The builder (`*minty.B`) provides a fluent, discoverable API that encapsulates HTML element creation while maintaining type safety and ultra-minimal syntax.

### The Builder Architecture

```go
// Builder provides methods for creating HTML elements
type Builder struct {
    // Internal state for rendering context, if needed
    context *RenderContext
}

// Standard constructor
func NewBuilder() *Builder {
    return &Builder{}
}

// Global builder instance for convenience
var B = NewBuilder()
```

The builder pattern in Minty serves multiple purposes beyond simple convenience. It provides a namespace for HTML element methods, enables future extensibility through context, and creates a consistent interface that IDEs can easily understand and autocomplete.

### Element Generation Methods

Each HTML element gets a corresponding method on the builder that follows consistent naming and signature patterns:

```go
// Container elements that accept children
func (b *Builder) Div(children ...Node) Node {
    return &Element{
        Tag:      "div",
        Children: children,
    }
}

func (b *Builder) P(children ...Node) Node {
    return &Element{
        Tag:      "p", 
        Children: children,
    }
}

func (b *Builder) H1(children ...Node) Node {
    return &Element{
        Tag:      "h1",
        Children: children,
    }
}

// Self-closing elements
func (b *Builder) Img(attrs ...Attribute) Node {
    element := &Element{
        Tag:         "img",
        SelfClosing: true,
        Attributes:  make(map[string]string),
    }
    
    for _, attr := range attrs {
        attr.Apply(element)
    }
    
    return element
}

func (b *Builder) Input(attrs ...Attribute) Node {
    element := &Element{
        Tag:         "input",
        SelfClosing: true,
        Attributes:  make(map[string]string),
    }
    
    for _, attr := range attrs {
        attr.Apply(element)
    }
    
    return element
}
```

This design provides several key benefits. The method names directly correspond to HTML element names, making the API immediately familiar to web developers. The varargs pattern for children enables natural nesting syntax. Self-closing elements are handled automatically, preventing malformed HTML.

### Attribute Handling Integration

The builder pattern integrates seamlessly with Minty's attribute system by accepting both child nodes and attributes in element methods:

```go
// Mixed attributes and children
func (b *Builder) A(attrs ...interface{}) Node {
    element := &Element{
        Tag:        "a",
        Attributes: make(map[string]string),
    }
    
    for _, attr := range attrs {
        switch v := attr.(type) {
        case Attribute:
            v.Apply(element)
        case Node:
            element.Children = append(element.Children, v)
        case string:
            // Automatic text node creation
            element.Children = append(element.Children, &TextNode{Content: v})
        }
    }
    
    return element
}

// Usage examples
b.A(minty.Href("/home"), minty.Class("nav-link"), "Home")
b.A(minty.Href("/about"), "About Us")
b.A(minty.Href("/contact"), 
    b.Span(minty.Class("icon"), "📧"),
    "Contact",
)
```

This flexible approach allows developers to specify attributes and content in whatever order feels most natural, while maintaining type safety through interface-based dispatch.

### Automatic Text Node Creation

One of the builder pattern's most important features is automatic text node creation for string parameters:

```go
// Automatic text node conversion
func processChild(child interface{}) Node {
    switch v := child.(type) {
    case Node:
        return v
    case string:
        return &TextNode{Content: v}
    case fmt.Stringer:
        return &TextNode{Content: v.String()}
    case int:
        return &TextNode{Content: strconv.Itoa(v)}
    case float64:
        return &TextNode{Content: strconv.FormatFloat(v, 'f', -1, 64)}
    default:
        return &TextNode{Content: fmt.Sprintf("%v", v)}
    }
}

// Applied in element methods
func (b *Builder) H1(children ...interface{}) Node {
    var nodes []Node
    for _, child := range children {
        nodes = append(nodes, processChild(child))
    }
    
    return &Element{
        Tag:      "h1",
        Children: nodes,
    }
}
```

This automatic conversion eliminates the explicit `Text()` wrappers required by other libraries while maintaining type safety and security through automatic HTML escaping.

### Builder Method Categories

Minty's builder methods fall into several categories, each optimized for their specific use patterns:

| Category | Examples | Characteristics | Usage Pattern |
|----------|----------|----------------|---------------|
| **Container Elements** | `Div()`, `P()`, `Section()` | Accept children, no restrictions | `b.Div(child1, child2, ...)` |
| **Text Elements** | `H1()`, `Span()`, `Strong()` | Accept text and inline elements | `b.H1("Title", b.Span("subtitle"))` |
| **Form Elements** | `Form()`, `Input()`, `Button()` | Mix of container and self-closing | `b.Form(b.Input(...), b.Button(...))` |
| **Void Elements** | `Img()`, `Br()`, `Hr()` | Self-closing, attribute-only | `b.Img(minty.Src(...), minty.Alt(...))` |
| **Document Structure** | `Html()`, `Head()`, `Body()` | Document-level containers | `b.Html(b.Head(...), b.Body(...))` |

Understanding these categories helps developers predict how unfamiliar elements will behave and choose the right elements for their use cases.

### Performance Optimizations in Builder Design

The builder pattern includes several performance optimizations that align with Minty's efficiency goals:

**Method Inlining**: Builder methods are designed to be simple enough for Go's compiler to inline them, eliminating function call overhead in performance-critical rendering paths.

**Minimal Allocations**: Element creation involves minimal allocations - typically just the element struct itself and its children slice. Attribute maps are only allocated when needed.

**Reusable Builders**: Builder instances are stateless and can be reused across multiple rendering operations, reducing allocation pressure in high-throughput scenarios.

```go
// Efficient reuse pattern
var b = minty.NewBuilder()

func handler1(w http.ResponseWriter, r *http.Request) {
    content := b.Div(b.H1("Page 1"))  // Reuses builder
    minty.Render(content, w)
}

func handler2(w http.ResponseWriter, r *http.Request) {
    content := b.Div(b.H1("Page 2"))  // Same builder instance
    minty.Render(content, w)
}
```

**Streaming-Friendly**: Builder methods create nodes that render directly to streams without intermediate string building, keeping memory usage constant regardless of output size.

### Extensibility Through Builder Context

The builder pattern provides a foundation for future extensibility through rendering context:

```go
// Future extensibility through context
type Builder struct {
    context *RenderContext
}

type RenderContext struct {
    Minify      bool
    Validate    bool
    BaseURL     string
    CSPNonce    string
    Theme       *Theme
}

// Context-aware element creation
func (b *Builder) Script(src string) Node {
    attrs := []Attribute{minty.Src(src)}
    
    if b.context != nil && b.context.CSPNonce != "" {
        attrs = append(attrs, minty.Nonce(b.context.CSPNonce))
    }
    
    return b.ScriptWithAttrs(attrs...)
}
```

This design allows Minty to evolve and add features like Content Security Policy support, theming, or validation without breaking existing code.

### Builder Pattern Benefits

The builder pattern provides several specific benefits that align with Minty's design principles:

**Discoverability**: IDE autocompletion can show all available HTML elements by typing `b.` and waiting for suggestions. This eliminates the need to memorize element names or consult documentation.

**Consistency**: All HTML elements follow the same creation pattern, reducing the cognitive load of learning the API. Once developers understand how `b.Div()` works, they understand how all container elements work.

**Type Safety**: Builder methods are strongly typed, preventing common errors like passing incorrect parameter types or using non-existent element names.

**Future-Proofing**: The builder pattern provides a stable interface that can evolve to support new HTML standards without breaking existing code.

**Performance**: The pattern enables optimizations like method inlining and allocation reduction while maintaining a convenient developer interface.

The builder pattern in Minty demonstrates how thoughtful API design can provide convenience and discoverability without sacrificing performance or type safety.

---

## Template Type System

Minty's template type system bridges the gap between Go's static typing and the dynamic nature of HTML generation. This system provides compile-time safety for templates while enabling powerful composition patterns that feel natural to Go developers.

### Core Template Types

The template type system revolves around two primary types that work together to provide flexibility and safety:

```go
// H represents a template function that generates HTML
type H func(*Builder) Node

// Node represents the result of template rendering
type Node interface {
    Render(io.Writer) error
}
```

This design separates template definition (`H`) from template result (`Node`), enabling different patterns of composition and reuse while maintaining type safety throughout.

### Template Function Patterns

Template functions in Minty follow several patterns that correspond to different use cases and complexity levels:

#### Simple Static Templates

Static templates generate the same content every time they're called:

```go
// Static template - no parameters needed
var homePage H = func(b *Builder) Node {
    return b.Html(
        b.Head(b.Title("Welcome")),
        b.Body(
            b.H1("Welcome to our site"),
            b.P("Thanks for visiting!"),
        ),
    )
}

// Usage
minty.Render(homePage(b), w)
```

Static templates are useful for unchanging content like error pages, static marketing pages, or legal documents.

#### Parameterized Templates

Parameterized templates accept data and generate dynamic content:

```go
// Parameterized template factory
func userProfile(user User) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("user-profile"),
            b.H1(user.Name),
            b.P("Email: " + user.Email),
            b.P("Joined: " + user.CreatedAt.Format("January 2006")),
        )
    }
}

// Usage
template := userProfile(currentUser)
minty.Render(template(b), w)
```

This pattern enables type-safe parameter passing while maintaining the template function interface.

#### Composite Templates

Composite templates combine multiple other templates into larger structures:

```go
// Layout template that accepts content
func layout(title string, content H) H {
    return func(b *Builder) Node {
        return b.Html(
            b.Head(
                b.Title(title),
                b.Meta("viewport", "width=device-width, initial-scale=1"),
            ),
            b.Body(
                header(b),
                b.Main(content(b)),  // Render content template
                footer(b),
            ),
        )
    }
}

// Usage
page := layout("User Profile", userProfile(currentUser))
minty.Render(page(b), w)
```

Composite templates enable sophisticated layout systems while preserving type safety and composition clarity.

### Template Composition Strategies

The template type system supports multiple composition strategies, each suited to different architectural needs:

#### Functional Composition

Pure functional composition treats templates as composable functions:

```go
// Compose templates functionally
func compose(outer H, inner H) H {
    return func(b *Builder) Node {
        return outer(b).WithChild(inner(b))  // Hypothetical API
    }
}

// Pipeline composition
func pipeline(templates ...H) H {
    return func(b *Builder) Node {
        var nodes []Node
        for _, template := range templates {
            nodes = append(nodes, template(b))
        }
        return minty.Fragment(nodes...)
    }
}
```

This approach leverages Go's functional programming capabilities while maintaining template modularity.

#### Context-Based Composition

Context-based composition passes shared data through template hierarchies:

```go
// Template context for shared data
type TemplateContext struct {
    User     User
    Request  *http.Request
    Session  Session
    Config   Config
}

// Context-aware template type
type ContextTemplate func(*Builder, TemplateContext) Node

// Layout with context
func layoutWithContext(title string, content ContextTemplate) ContextTemplate {
    return func(b *Builder, ctx TemplateContext) Node {
        return b.Html(
            b.Head(b.Title(title)),
            b.Body(
                navigation(b, ctx),    // Access to context
                b.Main(content(b, ctx)), // Pass context down
                footer(b, ctx),
            ),
        )
    }
}
```

Context-based composition provides access to shared data without explicit parameter threading.

#### Slot-Based Composition

Slot-based composition enables templates with multiple insertion points:

```go
// Template with multiple slots
type SlottedTemplate struct {
    Header H
    Main   H
    Sidebar H
    Footer H
}

func renderSlottedTemplate(template SlottedTemplate) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("layout"),
            b.Header(template.Header(b)),
            b.Div(minty.Class("content"),
                b.Main(template.Main(b)),
                b.Aside(template.Sidebar(b)),
            ),
            b.Footer(template.Footer(b)),
        )
    }
}

// Usage
page := SlottedTemplate{
    Header:  navigationTemplate,
    Main:    userProfileTemplate,
    Sidebar: sidebarAdsTemplate,
    Footer:  footerLinksTemplate,
}
```

Slot-based composition provides maximum flexibility for complex layouts while maintaining type safety.

### Template Inheritance and Extension

While Go doesn't have traditional inheritance, the template type system enables inheritance-like patterns through composition:

```go
// Base template with extension points
type BaseTemplate struct {
    Title       string
    Description string
    Scripts     []string
    Styles      []string
}

func (bt BaseTemplate) Extend(content H) H {
    return func(b *Builder) Node {
        return b.Html(
            b.Head(
                b.Title(bt.Title),
                b.Meta("description", bt.Description),
                // Render styles
                minty.Fragment(bt.renderStyles(b)...),
            ),
            b.Body(
                content(b),  // Extension point
                // Render scripts
                minty.Fragment(bt.renderScripts(b)...),
            ),
        )
    }
}

func (bt BaseTemplate) renderStyles(b *Builder) []Node {
    var nodes []Node
    for _, style := range bt.Styles {
        nodes = append(nodes, b.Link(minty.Rel("stylesheet"), minty.Href(style)))
    }
    return nodes
}

// Usage
baseTemplate := BaseTemplate{
    Title:       "User Profile",
    Description: "View and edit user profile information",
    Styles:      []string{"/css/base.css", "/css/profile.css"},
    Scripts:     []string{"/js/profile.js"},
}

userPage := baseTemplate.Extend(userProfileContent)
```

This pattern provides template inheritance benefits while leveraging Go's composition strengths.

### Type Safety Mechanisms

The template type system includes several mechanisms to ensure type safety throughout the template composition and rendering process:

#### Compile-Time Template Validation

Template functions are validated at compile time, catching errors before runtime:

```go
// Compile-time error detection
func brokenTemplate(user User) H {
    return func(b *Builder) Node {
        return b.Div(
            b.H1(user.Name),
            b.P(user.Emial),  // ✗ Compile error: field doesn't exist
        )
    }
}
```

This validation extends through the entire template composition chain, ensuring that complex template hierarchies remain type-safe.

#### Parameter Type Checking

Template parameters are type-checked at template creation time:

```go
// Type-safe parameter passing
func userCard(user User) H { /* ... */ }
func productCard(product Product) H { /* ... */ }

func itemList[T any](items []T, cardTemplate func(T) H) H {
    return func(b *Builder) Node {
        var cards []Node
        for _, item := range items {
            cards = append(cards, cardTemplate(item)(b))
        }
        return b.Div(minty.Class("item-list"), cards...)
    }
}

// Usage - type-checked at compile time
userList := itemList(users, userCard)        // ✓ Types match
productList := itemList(products, productCard) // ✓ Types match
broken := itemList(users, productCard)       // ✗ Compile error: type mismatch
```

Generic template functions provide type safety while enabling reusable patterns.

#### Return Type Validation

Template return types are validated to ensure they implement the Node interface correctly:

```go
// Return type validation
func validTemplate() H {
    return func(b *Builder) Node {
        return b.Div("Valid content")  // ✓ Returns Node
    }
}

func invalidTemplate() H {
    return func(b *Builder) Node {
        return "Invalid content"  // ✗ Compile error: string is not Node
    }
}
```

This validation prevents runtime errors that could occur from incorrect return types.

### Template Caching and Optimization

The template type system supports various caching and optimization strategies:

#### Template Result Caching

Template results can be cached for expensive-to-compute content:

```go
// Cached template wrapper
func cached(template H, key string) H {
    return func(b *Builder) Node {
        if cached := getFromCache(key); cached != nil {
            return cached
        }
        
        result := template(b)
        putInCache(key, result)
        return result
    }
}

// Usage
expensiveContent := cached(complexCalculationTemplate, "complex-calc-v1")
```

Caching integrates naturally with the template type system without requiring special syntax or annotations.

#### Template Compilation

The template type system enables ahead-of-time template compilation for maximum performance:

```go
// Compiled template representation
type CompiledTemplate struct {
    renderFunc func(io.Writer, map[string]interface{}) error
    paramTypes []reflect.Type
}

// Compilation process
func compileTemplate(template H) CompiledTemplate {
    // Analyze template structure
    // Generate optimized rendering function
    // Return compiled representation
}
```

Compilation can optimize away builder overhead and generate highly efficient rendering code.

### Integration with Go's Type System

Minty's template type system integrates seamlessly with Go's existing type system features:

#### Interface Satisfaction

Templates naturally work with Go interfaces:

```go
// Template interface
type Renderable interface {
    ToTemplate() H
}

// Types implementing the interface
func (u User) ToTemplate() H {
    return userTemplate(u)
}

func (p Product) ToTemplate() H {
    return productTemplate(p)
}

// Generic rendering function
func renderItems[T Renderable](items []T) H {
    return func(b *Builder) Node {
        var nodes []Node
        for _, item := range items {
            nodes = append(nodes, item.ToTemplate()(b))
        }
        return b.Div(nodes...)
    }
}
```

This integration enables polymorphic template rendering while maintaining type safety.

#### Error Handling Integration

Template errors integrate with Go's standard error handling patterns:

```go
// Template with error handling
type ErrorableTemplate func(*Builder) (Node, error)

func safeTemplate(template ErrorableTemplate) H {
    return func(b *Builder) Node {
        result, err := template(b)
        if err != nil {
            return b.Div(minty.Class("error"),
                b.P("Template error: " + err.Error()),
            )
        }
        return result
    }
}
```

This pattern enables graceful error handling in template hierarchies.

The template type system demonstrates how careful type design can provide both safety and flexibility, enabling complex template composition patterns while maintaining the simplicity and predictability that make Go attractive for web development.

---

## HTML Element Generation and Attributes

Minty's HTML generation system balances the competing demands of type safety, performance, and developer convenience. The system automatically handles HTML structure validation, attribute management, and content escaping while providing an intuitive API that mirrors HTML's natural structure.

### Element Generation Architecture

The element generation system is built around a consistent pattern that handles all HTML elements uniformly while accommodating their individual characteristics:

```go
// Core element structure
type Element struct {
    Tag         string
    Attributes  map[string]string
    Children    []Node
    SelfClosing bool
    Namespace   string  // For SVG, MathML, etc.
}

// Element creation pattern
func (b *Builder) createElementWithAttrs(tag string, selfClosing bool, args ...interface{}) Node {
    element := &Element{
        Tag:         tag,
        SelfClosing: selfClosing,
        Attributes:  make(map[string]string),
    }
    
    for _, arg := range args {
        switch v := arg.(type) {
        case Attribute:
            v.Apply(element)
        case Node:
            if !selfClosing {
                element.Children = append(element.Children, v)
            }
        case string:
            if !selfClosing {
                element.Children = append(element.Children, &TextNode{Content: v})
            }
        case int, float64, bool:
            if !selfClosing {
                element.Children = append(element.Children, &TextNode{
                    Content: fmt.Sprintf("%v", v),
                })
            }
        }
    }
    
    return element
}
```

This unified approach ensures consistent behavior across all HTML elements while handling the variations between container elements, void elements, and special cases.

### HTML5 Element Coverage

Minty provides complete coverage of HTML5 elements, organized by category for better developer understanding:

#### Document Structure Elements

```go
// Document structure
func (b *Builder) Html(children ...Node) Node {
    return b.createElementWithAttrs("html", false, children...)
}

func (b *Builder) Head(children ...Node) Node {
    return b.createElementWithAttrs("head", false, children...)
}

func (b *Builder) Body(children ...Node) Node {
    return b.createElementWithAttrs("body", false, children...)
}

func (b *Builder) Title(text string) Node {
    return b.createElementWithAttrs("title", false, text)
}
```

Document structure elements form the skeleton of HTML documents and require careful handling of content restrictions and nesting rules.

#### Content Sectioning Elements

```go
// Sectioning elements
func (b *Builder) Header(children ...Node) Node {
    return b.createElementWithAttrs("header", false, children...)
}

func (b *Builder) Nav(children ...Node) Node {
    return b.createElementWithAttrs("nav", false, children...)
}

func (b *Builder) Main(children ...Node) Node {
    return b.createElementWithAttrs("main", false, children...)
}

func (b *Builder) Section(children ...Node) Node {
    return b.createElementWithAttrs("section", false, children...)
}

func (b *Builder) Article(children ...Node) Node {
    return b.createElementWithAttrs("article", false, children...)
}

func (b *Builder) Aside(children ...Node) Node {
    return b.createElementWithAttrs("aside", false, children...)
}

func (b *Builder) Footer(children ...Node) Node {
    return b.createElementWithAttrs("footer", false, children...)
}
```

Sectioning elements provide semantic structure and are essential for accessibility and SEO.

#### Text Content Elements

```go
// Text and content elements
func (b *Builder) P(children ...Node) Node {
    return b.createElementWithAttrs("p", false, children...)
}

func (b *Builder) Div(children ...Node) Node {
    return b.createElementWithAttrs("div", false, children...)
}

func (b *Builder) Span(children ...Node) Node {
    return b.createElementWithAttrs("span", false, children...)
}

// Heading elements
func (b *Builder) H1(children ...Node) Node {
    return b.createElementWithAttrs("h1", false, children...)
}

func (b *Builder) H2(children ...Node) Node {
    return b.createElementWithAttrs("h2", false, children...)
}

// ... H3 through H6

// List elements
func (b *Builder) Ul(children ...Node) Node {
    return b.createElementWithAttrs("ul", false, children...)
}

func (b *Builder) Ol(children ...Node) Node {
    return b.createElementWithAttrs("ol", false, children...)
}

func (b *Builder) Li(children ...Node) Node {
    return b.createElementWithAttrs("li", false, children...)
}
```

Text content elements handle the majority of content structuring and benefit from Minty's automatic text node creation.

#### Form Elements

Form elements require special handling due to their complex attribute requirements and interaction patterns:

```go
// Form container
func (b *Builder) Form(args ...interface{}) Node {
    return b.createElementWithAttrs("form", false, args...)
}

// Input elements with type safety
func (b *Builder) Input(attrs ...Attribute) Node {
    return b.createElementWithAttrs("input", true, attrs...)
}

func (b *Builder) Button(args ...interface{}) Node {
    return b.createElementWithAttrs("button", false, args...)
}

func (b *Builder) Select(args ...interface{}) Node {
    return b.createElementWithAttrs("select", false, args...)
}

func (b *Builder) Option(value string, text string) Node {
    return b.createElementWithAttrs("option", false, 
        minty.Value(value), text)
}

func (b *Builder) Textarea(attrs ...Attribute) Node {
    return b.createElementWithAttrs("textarea", false, attrs...)
}

func (b *Builder) Label(args ...interface{}) Node {
    return b.createElementWithAttrs("label", false, args...)
}
```

Form elements demonstrate the attribute system's flexibility in handling complex requirements while maintaining type safety.

#### Void Elements

Void elements (also called empty or self-closing elements) require special handling since they cannot contain children:

```go
// Image element
func (b *Builder) Img(attrs ...Attribute) Node {
    return b.createElementWithAttrs("img", true, attrs...)
}

// Line break
func (b *Builder) Br() Node {
    return b.createElementWithAttrs("br", true)
}

// Horizontal rule
func (b *Builder) Hr(attrs ...Attribute) Node {
    return b.createElementWithAttrs("hr", true, attrs...)
}

// Meta elements
func (b *Builder) Meta(name string, content string) Node {
    return b.createElementWithAttrs("meta", true, 
        minty.Name(name), minty.Content(content))
}

// Link element
func (b *Builder) Link(attrs ...Attribute) Node {
    return b.createElementWithAttrs("link", true, attrs...)
}
```

Void elements showcase the type system's ability to prevent invalid HTML structures at compile time.

### Attribute System Design

Minty's attribute system provides type safety and convenience while handling the complexity of HTML attribute validation and formatting:

```go
// Attribute interface
type Attribute interface {
    Apply(*Element)
}

// String attribute implementation
type StringAttribute struct {
    Name  string
    Value string
}

func (sa StringAttribute) Apply(element *Element) {
    element.Attributes[sa.Name] = html.EscapeString(sa.Value)
}

// Boolean attribute implementation
type BooleanAttribute struct {
    Name string
}

func (ba BooleanAttribute) Apply(element *Element) {
    element.Attributes[ba.Name] = ba.Name
}

// Conditional attribute implementation
type ConditionalAttribute struct {
    Condition bool
    Attribute Attribute
}

func (ca ConditionalAttribute) Apply(element *Element) {
    if ca.Condition {
        ca.Attribute.Apply(element)
    }
}
```

This design enables different attribute types while maintaining a consistent interface and automatic escaping.

### Common Attribute Factories

Minty provides factory functions for all standard HTML attributes with appropriate type constraints:

```go
// Universal attributes
func Class(value string) Attribute {
    return StringAttribute{Name: "class", Value: value}
}

func ID(value string) Attribute {
    return StringAttribute{Name: "id", Value: value}
}

func Style(value string) Attribute {
    return StringAttribute{Name: "style", Value: value}
}

func Title(value string) Attribute {
    return StringAttribute{Name: "title", Value: value}
}

// Link and navigation attributes
func Href(url string) Attribute {
    return StringAttribute{Name: "href", Value: url}
}

func Target(value string) Attribute {
    return StringAttribute{Name: "target", Value: value}
}

// Form attributes with validation
func Name(value string) Attribute {
    return StringAttribute{Name: "name", Value: value}
}

func Value(value string) Attribute {
    return StringAttribute{Name: "value", Value: value}
}

func Type(inputType string) Attribute {
    // Could include validation for valid input types
    return StringAttribute{Name: "type", Value: inputType}
}

func Placeholder(value string) Attribute {
    return StringAttribute{Name: "placeholder", Value: value}
}

// Boolean attributes
func Required() Attribute {
    return BooleanAttribute{Name: "required"}
}

func Disabled() Attribute {
    return BooleanAttribute{Name: "disabled"}
}

func Checked() Attribute {
    return BooleanAttribute{Name: "checked"}
}

func Multiple() Attribute {
    return BooleanAttribute{Name: "multiple"}
}

// Numeric attributes
func TabIndex(index int) Attribute {
    return StringAttribute{Name: "tabindex", Value: strconv.Itoa(index)}
}

func ColSpan(span int) Attribute {
    return StringAttribute{Name: "colspan", Value: strconv.Itoa(span)}
}

func RowSpan(span int) Attribute {
    return StringAttribute{Name: "rowspan", Value: strconv.Itoa(span)}
}
```

These factory functions provide type-safe attribute creation while handling proper formatting and escaping automatically.

### Advanced Attribute Patterns

The attribute system supports advanced patterns for complex scenarios:

#### Conditional Attributes

```go
// Conditional attribute helper
func If(condition bool, attr Attribute) Attribute {
    return ConditionalAttribute{
        Condition: condition,
        Attribute: attr,
    }
}

// Usage
b.Button(
    Class("btn"),
    If(isPrimary, Class("btn-primary")),
    If(isDisabled, Disabled()),
    "Click me",
)
```

#### Attribute Merging

```go
// Mergeable attribute for complex cases like CSS classes
type MergeableAttribute struct {
    Name   string
    Values []string
}

func (ma MergeableAttribute) Apply(element *Element) {
    existing := element.Attributes[ma.Name]
    if existing != "" {
        ma.Values = append([]string{existing}, ma.Values...)
    }
    element.Attributes[ma.Name] = strings.Join(ma.Values, " ")
}

// CSS class helper that merges classes
func Classes(classes ...string) Attribute {
    return MergeableAttribute{
        Name:   "class",
        Values: classes,
    }
}

// Usage - classes are automatically merged
b.Div(
    Class("base"),
    If(isActive, Class("active")),
    Classes("responsive", "shadow"),
    "Content",
)
```

#### Data Attributes

```go
// Data attribute helper
func Data(name string, value string) Attribute {
    return StringAttribute{
        Name:  "data-" + name,
        Value: value,
    }
}

// ARIA attribute helper
func Aria(name string, value string) Attribute {
    return StringAttribute{
        Name:  "aria-" + name,
        Value: value,
    }
}

// Usage
b.Button(
    Data("action", "submit"),
    Data("confirm", "Are you sure?"),
    Aria("label", "Submit form"),
    "Submit",
)
```

### Element Validation and Error Prevention

The element generation system includes validation mechanisms to prevent common HTML errors:

#### Content Model Validation

```go
// Validation during element creation
func (b *Builder) validateContentModel(tag string, children []Node) error {
    rules := getContentModelRules(tag)
    
    for _, child := range children {
        if !rules.AllowsChild(child) {
            return fmt.Errorf("element <%s> cannot contain %T", tag, child)
        }
    }
    
    return nil
}

// Content model rules
type ContentModelRules struct {
    AllowedChildren   []string
    ForbiddenChildren []string
    RequiredChildren  []string
}

func (cmr ContentModelRules) AllowsChild(node Node) bool {
    // Implementation of content model checking
    return true
}
```

This validation can catch structural errors at compile time or development time, preventing invalid HTML generation.

#### Attribute Validation

```go
// Attribute validation during application
func validateAttribute(tag string, name string, value string) error {
    rules := getAttributeRules(tag, name)
    
    if !rules.IsValidValue(value) {
        return fmt.Errorf("invalid value '%s' for attribute '%s' on element <%s>", 
                         value, name, tag)
    }
    
    return nil
}

// Usage in attribute application
func (sa StringAttribute) Apply(element *Element) {
    if err := validateAttribute(element.Tag, sa.Name, sa.Value); err != nil {
        // Handle validation error (log, panic, or ignore based on configuration)
    }
    
    element.Attributes[sa.Name] = html.EscapeString(sa.Value)
}
```

Attribute validation helps catch common errors like invalid input types or malformed URLs.

The HTML element generation and attribute system demonstrates how comprehensive type safety can be achieved without sacrificing developer convenience, providing a foundation for reliable HTML generation that prevents common web development errors while maintaining the performance and simplicity that Go developers value.

---

## Text Handling and Automatic Escaping

Text content security and proper encoding form a critical foundation of any HTML generation system. Minty's text handling is designed around the principle of "secure by default," where the safest behavior requires the least effort from developers, while still providing escape hatches for legitimate use cases requiring unescaped content.

### Security-First Text Processing

Minty's text handling architecture prioritizes security through automatic HTML escaping as the default behavior:

```go
// TextNode with automatic escaping
type TextNode struct {
    Content string
    Escaped bool  // Track if content is already escaped
}

func (t *TextNode) Render(w io.Writer) error {
    var content string
    if t.Escaped {
        content = t.Content
    } else {
        content = html.EscapeString(t.Content)
    }
    
    _, err := w.Write([]byte(content))
    return err
}

// Automatic text node creation with escaping
func NewTextNode(content string) *TextNode {
    return &TextNode{
        Content: content,
        Escaped: false,  // Will be escaped during rendering
    }
}
```

This design ensures that all text content is automatically escaped unless explicitly marked as already safe, eliminating the most common source of XSS vulnerabilities in web applications.

### Automatic Content Type Detection

Minty automatically detects different content types and handles them appropriately:

```go
// Content type detection and conversion
func convertToNode(value interface{}) Node {
    switch v := value.(type) {
    case Node:
        return v
    case string:
        return NewTextNode(v)
    case []byte:
        return NewTextNode(string(v))
    case int:
        return NewTextNode(strconv.Itoa(v))
    case int64:
        return NewTextNode(strconv.FormatInt(v, 10))
    case float64:
        return NewTextNode(strconv.FormatFloat(v, 'f', -1, 64))
    case float32:
        return NewTextNode(strconv.FormatFloat(float64(v), 'f', -1, 32))
    case bool:
        return NewTextNode(strconv.FormatBool(v))
    case time.Time:
        return NewTextNode(v.Format(time.RFC3339))
    case fmt.Stringer:
        return NewTextNode(v.String())
    case error:
        return NewTextNode(v.Error())
    default:
        return NewTextNode(fmt.Sprintf("%v", v))
    }
}
```

This automatic conversion eliminates the need for manual string conversion while ensuring consistent escaping across all content types.

### HTML Escaping Implementation

Minty's HTML escaping goes beyond the standard library's basic escaping to handle edge cases and provide comprehensive protection:

```go
// Enhanced HTML escaping with additional safety measures
func escapeHTML(s string) string {
    // Start with standard library escaping
    escaped := html.EscapeString(s)
    
    // Additional escaping for enhanced security
    escaped = strings.ReplaceAll(escaped, "<!--", "&lt;!--")
    escaped = strings.ReplaceAll(escaped, "-->", "--&gt;")
    
    // Escape JavaScript protocol in potential URLs
    if strings.HasPrefix(strings.ToLower(s), "javascript:") {
        escaped = "javascript&colon;" + escaped[11:]
    }
    
    return escaped
}

// Context-aware escaping for different HTML contexts
type EscapeContext int

const (
    ContextHTML EscapeContext = iota
    ContextAttribute
    ContextURL
    ContextCSS
    ContextJS
)

func escapeForContext(s string, context EscapeContext) string {
    switch context {
    case ContextHTML:
        return escapeHTML(s)
    case ContextAttribute:
        return escapeHTMLAttribute(s)
    case ContextURL:
        return url.QueryEscape(s)
    case ContextCSS:
        return escapeCSSValue(s)
    case ContextJS:
        return escapeJSString(s)
    default:
        return escapeHTML(s)
    }
}
```

Context-aware escaping provides the appropriate level of protection for different parts of the HTML document.

### Raw Content Handling

While automatic escaping is the default, Minty provides controlled mechanisms for including unescaped content when necessary:

```go
// Raw content node for unescaped HTML
type RawNode struct {
    Content string
}

func (r *RawNode) Render(w io.Writer) error {
    _, err := w.Write([]byte(r.Content))
    return err
}

// Explicit raw content creation
func Raw(content string) Node {
    return &RawNode{Content: content}
}

// Trusted content creation with validation
func TrustedHTML(content string) Node {
    // Optional: validate that content is safe HTML
    if isValidHTML(content) {
        return &RawNode{Content: content}
    }
    
    // Fall back to escaped content if validation fails
    return NewTextNode(content)
}

// Pre-escaped content that doesn't need re-escaping
func PreEscaped(content string) Node {
    return &TextNode{
        Content: content,
        Escaped: true,
    }
}
```

These mechanisms require explicit opt-in, making it clear when potentially unsafe content is being used.

### Template String Processing

Minty handles template strings and interpolation securely:

```go
// Secure template string interpolation
func Interpolate(template string, values map[string]interface{}) Node {
    var result strings.Builder
    
    // Parse template and replace placeholders
    parts := parseTemplate(template)
    for _, part := range parts {
        if part.IsPlaceholder {
            value := values[part.Key]
            if value != nil {
                // Automatically escape interpolated values
                escapedValue := escapeHTML(fmt.Sprintf("%v", value))
                result.WriteString(escapedValue)
            }
        } else {
            result.WriteString(part.Text)
        }
    }
    
    return PreEscaped(result.String())
}

// Usage
greeting := Interpolate("Hello, {{name}}! Welcome to {{site}}.", map[string]interface{}{
    "name": userInput,  // Automatically escaped
    "site": "Our Site",
})
```

Template interpolation maintains security by escaping all interpolated values while allowing the template structure to remain unescaped.

### Content Sanitization

For scenarios requiring HTML input from users, Minty provides sanitization capabilities:

```go
// HTML sanitization policy
type SanitizationPolicy struct {
    AllowedTags       map[string]bool
    AllowedAttributes map[string][]string
    RemoveScripts     bool
    RemoveStyles      bool
}

// Default safe policy
var DefaultSafePolicy = SanitizationPolicy{
    AllowedTags: map[string]bool{
        "p": true, "br": true, "strong": true, "em": true,
        "ul": true, "ol": true, "li": true, "a": true,
    },
    AllowedAttributes: map[string][]string{
        "a": {"href", "title"},
    },
    RemoveScripts: true,
    RemoveStyles:  true,
}

// Sanitize HTML content
func SanitizeHTML(content string, policy SanitizationPolicy) Node {
    sanitized := sanitizeWithPolicy(content, policy)
    return Raw(sanitized)  // Safe to use raw after sanitization
}

// Usage
userContent := SanitizeHTML(userInput, DefaultSafePolicy)
```

Sanitization enables controlled acceptance of HTML input while maintaining security.

### Performance Optimizations

Minty's text handling includes several performance optimizations:

#### String Interning

```go
// String interning for common values
var stringInternPool = sync.Map{}

func internString(s string) string {
    if cached, exists := stringInternPool.Load(s); exists {
        return cached.(string)
    }
    
    stringInternPool.Store(s, s)
    return s
}

// Use interning for common CSS classes and other repeated strings
func Class(className string) Attribute {
    return StringAttribute{
        Name:  "class",
        Value: internString(className),
    }
}
```

String interning reduces memory usage for repeated values like CSS classes or common text content.

#### Escape Caching

```go
// Cache escaped versions of frequently used strings
var escapeCache = sync.Map{}
var escapeCacheSize int64

const maxEscapeCacheSize = 1000

func cachedEscapeHTML(s string) string {
    if cached, exists := escapeCache.Load(s); exists {
        return cached.(string)
    }
    
    escaped := escapeHTML(s)
    
    // Only cache if under size limit
    if atomic.LoadInt64(&escapeCacheSize) < maxEscapeCacheSize {
        escapeCache.Store(s, escaped)
        atomic.AddInt64(&escapeCacheSize, 1)
    }
    
    return escaped
}
```

Escape caching improves performance for repeated content while preventing unbounded memory growth.

#### Streaming Text Processing

```go
// Stream text directly to output without intermediate strings
func (t *TextNode) RenderStream(w io.Writer) error {
    if t.Escaped {
        _, err := io.WriteString(w, t.Content)
        return err
    }
    
    // Stream escaped content without building intermediate string
    return streamEscapedHTML(w, t.Content)
}

func streamEscapedHTML(w io.Writer, s string) error {
    start := 0
    for i, r := range s {
        var replacement string
        switch r {
        case '<':
            replacement = "&lt;"
        case '>':
            replacement = "&gt;"
        case '&':
            replacement = "&amp;"
        case '"':
            replacement = "&quot;"
        case '\'':
            replacement = "&#39;"
        default:
            continue
        }
        
        // Write unescaped portion
        if _, err := io.WriteString(w, s[start:i]); err != nil {
            return err
        }
        
        // Write escaped character
        if _, err := io.WriteString(w, replacement); err != nil {
            return err
        }
        
        start = i + 1
    }
    
    // Write remaining content
    _, err := io.WriteString(w, s[start:])
    return err
}
```

Streaming processing reduces memory allocation for large text content.

### Character Encoding Handling

Minty properly handles character encoding to prevent encoding-related security issues:

```go
// UTF-8 validation and normalization
func validateAndNormalizeUTF8(s string) string {
    if utf8.ValidString(s) {
        return s
    }
    
    // Replace invalid UTF-8 sequences with replacement character
    return strings.ToValidUTF8(s, "�")
}

// Content type handling with encoding specification
func (b *Builder) Meta(name, content string) Node {
    // Ensure charset is specified for content-type meta tags
    if name == "http-equiv" && strings.ToLower(content) == "content-type" {
        if !strings.Contains(content, "charset") {
            content += "; charset=utf-8"
        }
    }
    
    return b.createElementWithAttrs("meta", true,
        Name(name), Content(content))
}
```

Proper encoding handling prevents character encoding attacks and ensures consistent text rendering.

### Text Processing Error Handling

Minty's text processing includes comprehensive error handling:

```go
// Error-aware text processing
type TextProcessingError struct {
    Content   string
    Operation string
    Err       error
}

func (e TextProcessingError) Error() string {
    return fmt.Sprintf("text processing error in %s: %v", e.Operation, e.Err)
}

// Safe text processing with error recovery
func safeProcessText(content string, processor func(string) (string, error)) Node {
    processed, err := processor(content)
    if err != nil {
        // Log error and fall back to escaped original content
        log.Printf("Text processing failed: %v", err)
        return NewTextNode(content)
    }
    
    return PreEscaped(processed)
}
```

Error handling ensures that text processing failures don't break page rendering.

The text handling and automatic escaping system demonstrates how security can be built into the foundation of a templating system without sacrificing performance or developer convenience. By making the secure choice the default choice, Minty helps developers build secure web applications without requiring deep security expertise.

---

## Memory Management and Performance Characteristics

Minty's performance characteristics are fundamental to its design philosophy of providing a high-performance alternative to traditional templating systems. The architecture prioritizes predictable memory usage, minimal allocations, and efficient rendering to enable high-throughput web applications.

### Memory Allocation Patterns

Minty's memory allocation strategy focuses on minimizing allocations during the hot path of HTML generation:

```go
// Allocation-conscious element creation
type Element struct {
    Tag         string
    Attributes  map[string]string  // Allocated only when needed
    Children    []Node            // Pre-allocated with capacity
    SelfClosing bool
}

// Efficient element creation with pre-allocation
func (b *Builder) createElement(tag string, estimatedChildren int) *Element {
    element := &Element{
        Tag:         tag,
        SelfClosing: false,
    }
    
    // Pre-allocate children slice based on estimated size
    if estimatedChildren > 0 {
        element.Children = make([]Node, 0, estimatedChildren)
    }
    
    return element
}

// Lazy attribute map allocation
func (e *Element) ensureAttributes() {
    if e.Attributes == nil {
        e.Attributes = make(map[string]string, 4)  // Small initial capacity
    }
}

func (e *Element) AddAttribute(name, value string) {
    e.ensureAttributes()
    e.Attributes[name] = value
}
```

This approach minimizes allocations for simple elements while providing efficient growth for complex ones.

### Object Pooling for High-Throughput Scenarios

For high-traffic applications, Minty provides object pooling to reduce garbage collection pressure:

```go
// Element pool for reuse
var elementPool = sync.Pool{
    New: func() interface{} {
        return &Element{
            Children: make([]Node, 0, 4),
        }
    },
}

// Builder pool for concurrent requests
var builderPool = sync.Pool{
    New: func() interface{} {
        return &Builder{}
    },
}

// Pooled element creation
func (b *Builder) createPooledElement(tag string) *Element {
    element := elementPool.Get().(*Element)
    
    // Reset element state
    element.Tag = tag
    element.SelfClosing = false
    element.Children = element.Children[:0]  // Reset slice but keep capacity
    element.Attributes = nil  // Reset map
    
    return element
}

// Element cleanup and return to pool
func (e *Element) Release() {
    // Clear sensitive data if any
    for key := range e.Attributes {
        delete(e.Attributes, key)
    }
    
    elementPool.Put(e)
}

// High-throughput rendering with pooling
func PooledRender(template H, w io.Writer) error {
    builder := builderPool.Get().(*Builder)
    defer builderPool.Put(builder)
    
    content := template(builder)
    defer func() {
        if poolable, ok := content.(interface{ Release() }); ok {
            poolable.Release()
        }
    }()
    
    return content.Render(w)
}
```

Object pooling significantly reduces allocation overhead in high-throughput scenarios.

### Streaming Rendering Architecture

Minty's streaming rendering minimizes memory usage by writing output incrementally rather than building complete strings in memory:

```go
// Streaming renderer with minimal buffering
type StreamingRenderer struct {
    writer     io.Writer
    buffer     []byte
    bufferSize int
}

func NewStreamingRenderer(w io.Writer) *StreamingRenderer {
    return &StreamingRenderer{
        writer:     w,
        buffer:     make([]byte, 0, 4096),  // 4KB buffer
        bufferSize: 4096,
    }
}

func (sr *StreamingRenderer) Write(data []byte) error {
    if len(sr.buffer)+len(data) > sr.bufferSize {
        // Flush buffer when it would overflow
        if err := sr.Flush(); err != nil {
            return err
        }
    }
    
    if len(data) > sr.bufferSize {
        // Write large data directly without buffering
        _, err := sr.writer.Write(data)
        return err
    }
    
    sr.buffer = append(sr.buffer, data...)
    return nil
}

func (sr *StreamingRenderer) Flush() error {
    if len(sr.buffer) > 0 {
        _, err := sr.writer.Write(sr.buffer)
        sr.buffer = sr.buffer[:0]  // Reset buffer but keep capacity
        return err
    }
    return nil
}

// Element rendering with streaming
func (e *Element) RenderStreaming(sr *StreamingRenderer) error {
    // Write opening tag
    if err := sr.Write([]byte("<" + e.Tag)); err != nil {
        return err
    }
    
    // Write attributes
    for name, value := range e.Attributes {
        attr := fmt.Sprintf(` %s="%s"`, name, html.EscapeString(value))
        if err := sr.Write([]byte(attr)); err != nil {
            return err
        }
    }
    
    if e.SelfClosing {
        return sr.Write([]byte(" />"))
    }
    
    if err := sr.Write([]byte(">")); err != nil {
        return err
    }
    
    // Render children
    for _, child := range e.Children {
        if streamable, ok := child.(interface{ RenderStreaming(*StreamingRenderer) error }); ok {
            if err := streamable.RenderStreaming(sr); err != nil {
                return err
            }
        } else {
            // Fall back to regular rendering
            var buf bytes.Buffer
            if err := child.Render(&buf); err != nil {
                return err
            }
            if err := sr.Write(buf.Bytes()); err != nil {
                return err
            }
        }
    }
    
    // Write closing tag
    return sr.Write([]byte("</" + e.Tag + ">"))
}
```

Streaming rendering keeps memory usage constant regardless of output size.

### Performance Benchmarking and Optimization

Minty includes comprehensive benchmarking to ensure performance characteristics meet design goals:

```go
// Benchmark suite for performance monitoring
func BenchmarkSimpleElement(b *testing.B) {
    builder := NewBuilder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        element := builder.Div("Hello, World!")
        _ = element  // Prevent optimization
    }
}

func BenchmarkComplexElement(b *testing.B) {
    builder := NewBuilder()
    user := User{Name: "John", Email: "john@example.com"}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        element := builder.Div(minty.Class("user-card"),
            builder.H3(user.Name),
            builder.P(user.Email),
            builder.Button(minty.Class("btn btn-primary"), "Edit"),
        )
        _ = element
    }
}

func BenchmarkRendering(b *testing.B) {
    builder := NewBuilder()
    template := userListTemplate(generateTestUsers(100))
    content := template(builder)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        if err := content.Render(&buf); err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkStreamingRendering(b *testing.B) {
    builder := NewBuilder()
    template := userListTemplate(generateTestUsers(100))
    content := template(builder)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        renderer := NewStreamingRenderer(io.Discard)
        if err := content.RenderStreaming(renderer); err != nil {
            b.Fatal(err)
        }
        if err := renderer.Flush(); err != nil {
            b.Fatal(err)
        }
    }
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocation(b *testing.B) {
    builder := NewBuilder()
    
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        element := builder.Div(minty.Class("container"),
            builder.H1("Title"),
            builder.P("Content"),
        )
        
        var buf bytes.Buffer
        if err := element.Render(&buf); err != nil {
            b.Fatal(err)
        }
    }
}
```

Regular benchmarking ensures that performance improvements and optimizations have measurable impact.

### Garbage Collection Optimization

Minty's design minimizes garbage collection pressure through several strategies:

#### Allocation Reduction

```go
// Minimize allocations through careful design
type optimizedElement struct {
    tag         string
    attributes  [4]string     // Small fixed array for common case
    attrCount   int
    children    [8]Node       // Small fixed array for common case
    childCount  int
    overflow    *overflowData // Only allocated when needed
}

type overflowData struct {
    attributes map[string]string
    children   []Node
}

// Add attribute with minimal allocation
func (e *optimizedElement) addAttribute(name, value string) {
    if e.attrCount < len(e.attributes)-1 {
        e.attributes[e.attrCount] = name
        e.attributes[e.attrCount+1] = value
        e.attrCount += 2
    } else {
        e.ensureOverflow()
        e.overflow.attributes[name] = value
    }
}

func (e *optimizedElement) ensureOverflow() {
    if e.overflow == nil {
        e.overflow = &overflowData{
            attributes: make(map[string]string),
            children:   make([]Node, 0, 8),
        }
    }
}
```

Fixed arrays for common cases eliminate allocations for simple elements while providing overflow capacity for complex ones.

#### String Builder Optimization

```go
// Optimized string building for attribute rendering
type attributeBuilder struct {
    buffer []byte
    cap    int
}

func newAttributeBuilder() *attributeBuilder {
    return &attributeBuilder{
        buffer: make([]byte, 0, 256),  // Reasonable initial capacity
        cap:    256,
    }
}

func (ab *attributeBuilder) writeAttribute(name, value string) {
    needed := len(name) + len(value) + 4  // ` name="value"`
    if len(ab.buffer)+needed > ab.cap {
        ab.grow(needed)
    }
    
    ab.buffer = append(ab.buffer, ' ')
    ab.buffer = append(ab.buffer, name...)
    ab.buffer = append(ab.buffer, '=', '"')
    ab.buffer = appendEscaped(ab.buffer, value)
    ab.buffer = append(ab.buffer, '"')
}

func (ab *attributeBuilder) grow(needed int) {
    newCap := ab.cap * 2
    if newCap < ab.cap+needed {
        newCap = ab.cap + needed
    }
    
    newBuffer := make([]byte, len(ab.buffer), newCap)
    copy(newBuffer, ab.buffer)
    ab.buffer = newBuffer
    ab.cap = newCap
}

// Append escaped content without additional allocation
func appendEscaped(dst []byte, src string) []byte {
    for _, r := range src {
        switch r {
        case '<':
            dst = append(dst, "&lt;"...)
        case '>':
            dst = append(dst, "&gt;"...)
        case '&':
            dst = append(dst, "&amp;"...)
        case '"':
            dst = append(dst, "&quot;"...)
        default:
            dst = append(dst, string(r)...)
        }
    }
    return dst
}
```

Optimized string building reduces allocation overhead during attribute rendering.

### Performance Characteristics Summary

Minty's performance characteristics can be summarized in several key metrics:

| Operation | Time Complexity | Space Complexity | Allocations |
|-----------|----------------|------------------|-------------|
| **Element Creation** | O(1) | O(1) | 1 allocation per element |
| **Attribute Addition** | O(1) amortized | O(n) attributes | 0-1 allocations |
| **Child Addition** | O(1) amortized | O(n) children | 0-1 allocations |
| **Rendering** | O(n) nodes | O(1) with streaming | 0-2 allocations |
| **Template Execution** | O(n) complexity | O(d) depth | Function call overhead |

These characteristics enable predictable performance scaling and efficient resource utilization in production environments.

### Real-World Performance Testing

Minty includes real-world performance tests that simulate actual application usage patterns:

```go
// Real-world scenario: blog post rendering
func BenchmarkBlogPostRendering(b *testing.B) {
    post := generateTestBlogPost()
    template := blogPostTemplate(post)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        if err := RenderTemplate(template, &buf); err != nil {
            b.Fatal(err)
        }
    }
}

// Real-world scenario: data table rendering
func BenchmarkDataTableRendering(b *testing.B) {
    users := generateTestUsers(1000)
    template := userTableTemplate(users)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        if err := RenderTemplate(template, &buf); err != nil {
            b.Fatal(err)
        }
    }
}

// Memory pressure test
func TestMemoryPressure(t *testing.T) {
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // Render many templates to test memory usage
    for i := 0; i < 10000; i++ {
        template := complexPageTemplate(generateTestData())
        var buf bytes.Buffer
        if err := RenderTemplate(template, &buf); err != nil {
            t.Fatal(err)
        }
    }
    
    runtime.GC()
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    // Verify memory usage is reasonable
    allocatedMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024
    if allocatedMB > 50 {  // Reasonable threshold
        t.Errorf("Memory usage too high: %.2f MB", allocatedMB)
    }
}
```

Real-world testing ensures that performance characteristics hold up under actual usage conditions.

The memory management and performance characteristics demonstrate how careful architectural decisions can provide both developer convenience and production-ready performance. By focusing on allocation reduction, efficient rendering, and predictable scaling, Minty enables high-performance web applications without sacrificing the simplicity and type safety that make it attractive for development.

---

## Architecture Across the Minty System

The foundational architecture described in this document enables sophisticated patterns throughout the Minty System:

### Node Interface Supporting Complex Domains

The simple Node interface scales to support complex business domain presentations:

```go
// Business domain data flows through the same Node architecture
func FinancialDashboard(theme Theme, data mifi.DashboardData) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("dashboard"),
            // Complex business components using same Node foundation
            AccountSummarySection(theme, data.Accounts)(b),
            TransactionHistorySection(theme, data.Transactions)(b),
            InvoiceStatusSection(theme, data.Invoices)(b),
        )
    }
}
```

### Builder Pattern Enabling Iterator Integration

The builder pattern seamlessly integrates with iterator functions for data processing:

```go
// Iterator functions work naturally with the builder pattern
func UsersList(users []User) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Ul(mi.Class("users-list"),
            // Iterator-generated nodes integrate naturally
            miex.Map(users, func(u User) mi.H {
                return func(b *mi.Builder) mi.Node {
                    return b.Li(mi.Class("user-item"),
                        b.Strong(u.Name),
                        b.P(u.Email),
                    )
                }
            })...,
        )
    }
}
```

### Type System Supporting Theme Abstractions

The template type system enables flexible theme implementations while maintaining type safety:

```go
// Theme interface builds on Node foundation
type Theme interface {
    Button(text, variant string, attrs ...mi.Attribute) mi.H
    Card(title string, content mi.H) mi.H
    // All theme methods return mi.H (template functions)
}

// Theme implementations leverage the same architectural patterns
func (t *BootstrapTheme) Card(title string, content mi.H) mi.H {
    return func(b *mi.Builder) mi.Node {
        // Uses same Node interface and builder pattern
        return b.Div(mi.Class("card"),
            b.Div(mi.Class("card-header"), b.H5(title)),
            b.Div(mi.Class("card-body"), content(b)),
        )
    }
}
```

### Performance Characteristics Across System Layers

The memory management and performance optimizations benefit the entire system:

```go
// Complex multi-domain applications maintain performance
func CompleteApplication(services *ApplicationServices) mi.H {
    return func(b *mi.Builder) mi.Node {
        // Large applications with multiple domains
        // Still benefit from efficient Node rendering
        return b.Html(
            ApplicationHead()(b),
            b.Body(
                Navigation(services.User)(b),
                MainContent(
                    FinanceDashboard(services.Finance)(b),
                    LogisticsDashboard(services.Logistics)(b), 
                    EcommerceDashboard(services.Ecommerce)(b),
                )(b),
                ApplicationFooter()(b),
            ),
        )
    }
}
```

This foundational architecture enables the Minty System to scale from simple HTML generation to complex business applications while maintaining consistent patterns, type safety, and performance characteristics.

---

## Conclusion

Minty's core architecture and type system demonstrate how thoughtful design can bridge the gap between developer convenience and system efficiency. The Node interface provides a simple but powerful foundation for HTML generation, while the builder pattern offers an intuitive API that leverages Go's strengths in type safety and composition.

The template type system enables sophisticated composition patterns while maintaining compile-time safety, and the HTML generation system handles the complexity of modern web standards without exposing that complexity to developers. Automatic text escaping provides security by default, while the performance characteristics ensure that choosing the safe and convenient path doesn't compromise application efficiency.

Together, these architectural components create a foundation that supports not only Minty's core design philosophy of ultra-minimalism, type safety, and developer happiness, but also the **entire Minty System** including business domains, iterator functionality, theme systems, and presentation layers. This architectural foundation scales from simple HTML components to complex multi-domain applications while maintaining consistent patterns and performance characteristics.

The next part of this documentation series will explore how this foundation enables the ultra-concise syntax that makes Minty distinctive and productive across all components of the system.

---

*This document is part of the comprehensive Minty documentation series. Continue with [Part 4: Syntax Design & API](minty-04.md) to explore how these architectural foundations enable Minty's distinctive ultra-concise syntax.*