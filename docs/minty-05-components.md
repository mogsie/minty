# **Minty Documentation - Part 5**
## **Template Composition & Patterns**
### *"Building Reusable Components"*

> **Part of the Minty System**: This document covers template composition patterns that are fundamental to the entire Minty System. These patterns apply to HTML generation (minty), business domain presentation (mintyfinui, mintymoveui, mintycartui), theme implementations (bootstrap, tailwind), and iterator-based component generation. The composition techniques described here scale from simple HTML components to complex multi-domain applications.

---

## **Executive Summary**

Template composition represents the architectural heart of sophisticated web applications. While individual templates handle specific rendering tasks, composition patterns enable the construction of complex, maintainable, and reusable component systems. Minty's composition architecture transforms the traditional challenge of HTML template organization into an elegant exercise in functional programming and type-safe component design.

This document explores how Minty's template system enables powerful composition patterns that scale from simple parameter injection to sophisticated component libraries. Unlike traditional templating approaches that rely on string interpolation or complex inheritance hierarchies, Minty's functional composition model provides both simplicity and power through Go's native type system and functional programming capabilities.

---

## **Foundation: Template Functions as First-Class Citizens**

### **The Template Function Paradigm**

Minty treats templates as first-class functions, enabling powerful composition patterns through Go's native functional programming support. This architectural decision creates a foundation where templates can be passed as parameters, returned from functions, and composed into larger structures with complete type safety.

```go
// Basic template function signature
type H func(*Builder) Node

// Template as a simple function
func simpleGreeting() H {
    return func(b *Builder) Node {
        return b.H1("Hello, World!")
    }
}

// Template with parameter capture
func personalGreeting(name string) H {
    return func(b *Builder) Node {
        return b.H1("Hello, " + name + "!")
    }
}

// Template with complex parameter structure
func userCard(user User) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("user-card"),
            b.Img(minty.Src(user.Avatar), minty.Alt("Profile")),
            b.H3(user.Name),
            b.P(minty.Class("email"), user.Email),
            b.P(minty.Class("joined"), 
                "Member since " + user.JoinDate.Format("January 2006")),
        )
    }
}
```

The template function paradigm enables clean parameter passing while maintaining the simple `H` interface that can be composed with other templates regardless of their internal complexity.

### **Template Factory Patterns**

Template factories provide a systematic approach to creating families of related templates with shared behavior and consistent interfaces. This pattern proves particularly valuable for component libraries and design systems where consistency across similar components is crucial.

```go
// Generic template factory with typed parameters
func TemplateFactory[T any](renderer func(*Builder, T) Node) func(T) H {
    return func(param T) H {
        return func(b *Builder) Node {
            return renderer(b, param)
        }
    }
}

// Specialized factory for button components
func ButtonFactory(variant ButtonVariant, size ButtonSize) func(string) H {
    return func(text string) H {
        return func(b *Builder) Node {
            classes := []string{"btn"}
            
            switch variant {
            case ButtonPrimary:
                classes = append(classes, "btn-primary")
            case ButtonSecondary:
                classes = append(classes, "btn-secondary")
            case ButtonDanger:
                classes = append(classes, "btn-danger")
            }
            
            switch size {
            case ButtonSmall:
                classes = append(classes, "btn-sm")
            case ButtonLarge:
                classes = append(classes, "btn-lg")
            }
            
            return b.Button(
                minty.Class(strings.Join(classes, " ")),
                text,
            )
        }
    }
}

// Usage of button factory
var (
    PrimaryButton = ButtonFactory(ButtonPrimary, ButtonMedium)
    SmallButton   = ButtonFactory(ButtonSecondary, ButtonSmall)
    DangerButton  = ButtonFactory(ButtonDanger, ButtonMedium)
)

// Creating specific button instances
submitButton := PrimaryButton("Submit Form")
cancelButton := SmallButton("Cancel")
deleteButton := DangerButton("Delete Account")
```

Template factories enable the creation of component families with consistent behavior while allowing customization through parameter injection and variant selection.

### **Parameterized Template Systems**

Sophisticated applications require templates that can accept complex data structures and render them consistently across different contexts. Minty's parameterized template system provides type-safe parameter passing while maintaining composition flexibility.

```go
// Complex parameter structure
type ProductDisplayConfig struct {
    ShowPrice       bool
    ShowDescription bool
    ShowReviews     bool
    ShowAddToCart   bool
    ImageSize       string
    PriceFormat     string
}

// Parameterized template with configuration
func productCard(product Product, config ProductDisplayConfig) H {
    return func(b *Builder) Node {
        cardClasses := []string{"product-card"}
        if config.ImageSize != "" {
            cardClasses = append(cardClasses, "image-"+config.ImageSize)
        }
        
        children := []Node{
            b.Img(
                minty.Src(product.ImageURL),
                minty.Alt(product.Name),
                minty.Class("product-image"),
            ),
            b.H3(minty.Class("product-name"), product.Name),
        }
        
        if config.ShowDescription && product.Description != "" {
            children = append(children,
                b.P(minty.Class("product-description"), product.Description),
            )
        }
        
        if config.ShowPrice {
            priceText := formatPrice(product.Price, config.PriceFormat)
            children = append(children,
                b.P(minty.Class("product-price"), priceText),
            )
        }
        
        if config.ShowReviews && product.ReviewCount > 0 {
            children = append(children,
                starRating(product.AverageRating)(b),
                b.P(minty.Class("review-count"), 
                    fmt.Sprintf("(%d reviews)", product.ReviewCount)),
            )
        }
        
        if config.ShowAddToCart {
            children = append(children,
                b.Button(
                    minty.Class("btn btn-primary add-to-cart"),
                    minty.Data("product-id", strconv.Itoa(product.ID)),
                    "Add to Cart",
                ),
            )
        }
        
        return b.Div(minty.Class(strings.Join(cardClasses, " ")), children...)
    }
}

// Configuration presets for different contexts
var (
    ProductGridConfig = ProductDisplayConfig{
        ShowPrice:     true,
        ShowReviews:   true,
        ShowAddToCart: true,
        ImageSize:     "medium",
        PriceFormat:   "standard",
    }
    
    ProductSearchConfig = ProductDisplayConfig{
        ShowPrice:       true,
        ShowDescription: true,
        ShowReviews:     false,
        ShowAddToCart:   false,
        ImageSize:       "small",
        PriceFormat:     "compact",
    }
    
    ProductDetailConfig = ProductDisplayConfig{
        ShowPrice:       true,
        ShowDescription: true,
        ShowReviews:     true,
        ShowAddToCart:   true,
        ImageSize:       "large",
        PriceFormat:     "detailed",
    }
)
```

Parameterized template systems enable the same component to render appropriately in different contexts while maintaining consistency and type safety throughout the application.

---

## **Layout Systems and Content Injection**

### **Content Injection Architecture**

Layout systems require a mechanism for injecting dynamic content into predefined structural templates. Minty's approach to content injection combines functional composition with explicit parameter passing to create flexible and maintainable layout systems.

```go
// Layout template with content injection points
func basicLayout(title string, content H) H {
    return func(b *Builder) Node {
        return b.Html(
            b.Head(
                b.Title(title),
                b.Meta(minty.Name("viewport"), minty.Content("width=device-width, initial-scale=1")),
                b.Link(minty.Rel("stylesheet"), minty.Href("/static/css/app.css")),
            ),
            b.Body(
                header()(b),
                b.Main(minty.Class("main-content"), content(b)),
                footer()(b),
            ),
        )
    }
}

// Multi-zone layout with multiple injection points
func dashboardLayout(title string, sidebar H, main H, notifications H) H {
    return func(b *Builder) Node {
        return b.Html(
            b.Head(
                b.Title(title + " - Dashboard"),
                b.Link(minty.Rel("stylesheet"), minty.Href("/static/css/dashboard.css")),
            ),
            b.Body(minty.Class("dashboard-layout"),
                b.Header(minty.Class("dashboard-header"),
                    navigationBar()(b),
                    b.Div(minty.Class("notifications-area"), notifications(b)),
                ),
                b.Div(minty.Class("dashboard-content"),
                    b.Aside(minty.Class("sidebar"), sidebar(b)),
                    b.Main(minty.Class("main-panel"), main(b)),
                ),
            ),
        )
    }
}

// Usage examples
homePage := basicLayout("Welcome", func(b *Builder) Node {
    return b.Div(
        b.H1("Welcome to Our Site"),
        b.P("This is the home page content."),
    )
})

userDashboard := dashboardLayout(
    "User Dashboard",
    userSidebar(currentUser),
    userMainContent(currentUser),
    userNotifications(currentUser),
)
```

Content injection architecture enables layouts to accept dynamic content while maintaining structural consistency across the application.

### **Nested Layout Hierarchies**

Complex applications often require nested layout systems where specialized layouts build upon more general base layouts. Minty's composition system enables elegant layout hierarchies through function composition.

```go
// Base application layout
func appLayout(title string, bodyClass string, content H) H {
    return func(b *Builder) Node {
        return b.Html(
            b.Head(
                b.Title(title + " - My Application"),
                commonHeadElements()(b),
            ),
            b.Body(minty.Class("app-body "+bodyClass),
                appHeader()(b),
                content(b),
                appFooter()(b),
            ),
        )
    }
}

// Admin-specific layout building on app layout
func adminLayout(title string, content H) H {
    adminContent := func(b *Builder) Node {
        return b.Div(minty.Class("admin-container"),
            adminSidebar()(b),
            b.Main(minty.Class("admin-main"), content(b)),
        )
    }
    
    return appLayout(title, "admin-layout", adminContent)
}

// User-specific layout building on app layout
func userLayout(title string, user User, content H) H {
    userContent := func(b *Builder) Node {
        return b.Div(minty.Class("user-container"),
            userNavigation(user)(b),
            b.Main(minty.Class("user-main"), content(b)),
        )
    }
    
    return appLayout(title, "user-layout", userContent)
}

// Specialized admin page layout
func adminReportsLayout(title string, reportContent H) H {
    reportsContent := func(b *Builder) Node {
        return b.Div(minty.Class("reports-container"),
            reportsToolbar()(b),
            b.Div(minty.Class("reports-content"), reportContent(b)),
        )
    }
    
    return adminLayout(title, reportsContent)
}

// Usage in specific pages
salesReportPage := adminReportsLayout("Sales Report", func(b *Builder) Node {
    return b.Div(
        salesChart(salesData)(b),
        salesTable(salesData)(b),
    )
})
```

Nested layout hierarchies enable code reuse while allowing each layer to add its specific structural and styling requirements.

### **Layout Composition Patterns**

Advanced layout systems require flexible composition patterns that can adapt to different content requirements while maintaining consistency. Minty provides several approaches to layout composition that balance flexibility with maintainability.

```go
// Slot-based layout composition
type LayoutSlots struct {
    Header   H
    Sidebar  H
    Main     H
    Footer   H
    Modal    H
}

func slottedLayout(title string, slots LayoutSlots) H {
    return func(b *Builder) Node {
        children := []Node{
            b.Head(
                b.Title(title),
                commonHeadElements()(b),
            ),
        }
        
        bodyChildren := []Node{}
        
        if slots.Header != nil {
            bodyChildren = append(bodyChildren, 
                b.Header(minty.Class("page-header"), slots.Header(b)))
        }
        
        contentArea := []Node{}
        
        if slots.Sidebar != nil {
            contentArea = append(contentArea,
                b.Aside(minty.Class("page-sidebar"), slots.Sidebar(b)))
        }
        
        if slots.Main != nil {
            contentArea = append(contentArea,
                b.Main(minty.Class("page-main"), slots.Main(b)))
        }
        
        if len(contentArea) > 0 {
            bodyChildren = append(bodyChildren,
                b.Div(minty.Class("content-area"), contentArea...))
        }
        
        if slots.Footer != nil {
            bodyChildren = append(bodyChildren,
                b.Footer(minty.Class("page-footer"), slots.Footer(b)))
        }
        
        if slots.Modal != nil {
            bodyChildren = append(bodyChildren, slots.Modal(b))
        }
        
        children = append(children, b.Body(bodyChildren...))
        
        return b.Html(children...)
    }
}

// Builder pattern for layout construction
type LayoutBuilder struct {
    title   string
    slots   LayoutSlots
    classes []string
}

func NewLayout(title string) *LayoutBuilder {
    return &LayoutBuilder{
        title: title,
        slots: LayoutSlots{},
    }
}

func (lb *LayoutBuilder) WithHeader(header H) *LayoutBuilder {
    lb.slots.Header = header
    return lb
}

func (lb *LayoutBuilder) WithSidebar(sidebar H) *LayoutBuilder {
    lb.slots.Sidebar = sidebar
    return lb
}

func (lb *LayoutBuilder) WithMain(main H) *LayoutBuilder {
    lb.slots.Main = main
    return lb
}

func (lb *LayoutBuilder) WithFooter(footer H) *LayoutBuilder {
    lb.slots.Footer = footer
    return lb
}

func (lb *LayoutBuilder) WithModal(modal H) *LayoutBuilder {
    lb.slots.Modal = modal
    return lb
}

func (lb *LayoutBuilder) WithClass(class string) *LayoutBuilder {
    lb.classes = append(lb.classes, class)
    return lb
}

func (lb *LayoutBuilder) Build() H {
    return slottedLayout(lb.title, lb.slots)
}

// Usage with builder pattern
profilePage := NewLayout("User Profile").
    WithHeader(navigationHeader()).
    WithSidebar(profileSidebar(user)).
    WithMain(profileContent(user)).
    WithFooter(standardFooter()).
    WithClass("profile-page").
    Build()
```

Layout composition patterns provide multiple approaches to constructing complex page structures while maintaining clear separation of concerns and reusability.

---

## **Component Composition Patterns**

### **Functional Composition Strategy**

Functional composition treats templates as pure functions that can be combined using higher-order functions and composition utilities. This approach leverages Go's functional programming capabilities to create powerful component systems.

```go
// Template composition utilities
func Compose(templates ...H) H {
    return func(b *Builder) Node {
        var nodes []Node
        for _, template := range templates {
            if template != nil {
                nodes = append(nodes, template(b))
            }
        }
        return minty.Fragment(nodes...)
    }
}

func Wrap(wrapper func(H) H, inner H) H {
    return wrapper(inner)
}

func Chain(base H, transforms ...func(H) H) H {
    result := base
    for _, transform := range transforms {
        result = transform(result)
    }
    return result
}

// Component transformation functions
func WithClass(className string) func(H) H {
    return func(template H) H {
        return func(b *Builder) Node {
            return b.Div(minty.Class(className), template(b))
        }
    }
}

func WithErrorBoundary(errorMessage string) func(H) H {
    return func(template H) H {
        return func(b *Builder) Node {
            // In a real implementation, this might include error handling
            return b.Div(minty.Class("error-boundary"),
                template(b),
            )
        }
    }
}

func WithLoading(isLoading bool, loadingTemplate H) func(H) H {
    return func(template H) H {
        return func(b *Builder) Node {
            if isLoading {
                return loadingTemplate(b)
            }
            return template(b)
        }
    }
}

// Usage of functional composition
userProfile := Chain(
    basicUserInfo(user),
    WithClass("user-profile"),
    WithErrorBoundary("Failed to load user profile"),
    WithLoading(user.IsLoading, loadingSpinner()),
)

navigationMenu := Compose(
    logoSection(),
    mainNavigation(navigationItems),
    userActions(currentUser),
    searchBox(),
)
```

Functional composition strategy enables the construction of complex components through the combination of simpler, reusable building blocks.

### **Context-Aware Composition**

Many applications require components that can access shared data without explicit parameter threading. Context-aware composition provides a mechanism for passing shared state through component hierarchies.

```go
// Template context for shared application state
type TemplateContext struct {
    User      User
    Request   *http.Request
    Session   Session
    Config    AppConfig
    Theme     Theme
    Language  string
    TimeZone  *time.Location
}

// Context-aware template type
type ContextTemplate func(*Builder, TemplateContext) Node

// Context provider wrapper
func WithContext(ctx TemplateContext, template ContextTemplate) H {
    return func(b *Builder) Node {
        return template(b, ctx)
    }
}

// Context-aware components
func contextualNavigation(b *Builder, ctx TemplateContext) Node {
    return b.Nav(minty.Class("main-navigation"),
        b.Div(minty.Class("nav-brand"),
            b.A(minty.Href("/"), ctx.Config.AppName),
        ),
        b.Ul(minty.Class("nav-links"),
            conditionalNavLink(ctx.User.HasPermission("admin"), "/admin", "Admin", b),
            conditionalNavLink(ctx.User.IsAuthenticated(), "/dashboard", "Dashboard", b),
            b.Li(
                b.A(minty.Href("/profile"), 
                    fmt.Sprintf("Welcome, %s", ctx.User.DisplayName())),
            ),
        ),
    )
}

func contextualFooter(b *Builder, ctx TemplateContext) Node {
    currentYear := time.Now().In(ctx.TimeZone).Year()
    
    return b.Footer(minty.Class("site-footer"),
        b.P(fmt.Sprintf("© %d %s", currentYear, ctx.Config.CompanyName)),
        b.P(fmt.Sprintf("Current time: %s", 
            time.Now().In(ctx.TimeZone).Format("3:04 PM MST"))),
    )
}

// Context-aware layout
func contextualLayout(title string, content ContextTemplate) ContextTemplate {
    return func(b *Builder, ctx TemplateContext) Node {
        return b.Html(
            b.Head(
                b.Title(fmt.Sprintf("%s - %s", title, ctx.Config.AppName)),
                b.Meta(minty.Name("language"), minty.Content(ctx.Language)),
            ),
            b.Body(
                contextualNavigation(b, ctx),
                b.Main(content(b, ctx)),
                contextualFooter(b, ctx),
            ),
        )
    }
}

// Usage with context injection
ctx := TemplateContext{
    User:     currentUser,
    Request:  req,
    Config:   appConfig,
    Theme:    userTheme,
    Language: "en",
    TimeZone: userTimeZone,
}

profilePage := WithContext(ctx, contextualLayout("Profile", 
    func(b *Builder, ctx TemplateContext) Node {
        return userProfileContent(ctx.User)(b)
    }))
```

Context-aware composition enables components to access shared application state while maintaining clean interfaces and avoiding excessive parameter threading.

### **Hierarchical Component Systems**

Large applications benefit from hierarchical component organization that promotes reuse while maintaining clear dependency relationships. Minty's composition system supports sophisticated component hierarchies.

```go
// Base component interfaces
type Component interface {
    Render(*Builder) Node
}

type ParameterizedComponent[T any] interface {
    Render(*Builder, T) Node
}

// Base UI components
type Button struct {
    Text    string
    Variant ButtonVariant
    Size    ButtonSize
    OnClick string
}

func (btn Button) Render(b *Builder) Node {
    classes := []string{"btn"}
    
    switch btn.Variant {
    case ButtonPrimary:
        classes = append(classes, "btn-primary")
    case ButtonSecondary:
        classes = append(classes, "btn-secondary")
    }
    
    switch btn.Size {
    case ButtonSmall:
        classes = append(classes, "btn-sm")
    case ButtonLarge:
        classes = append(classes, "btn-lg")
    }
    
    attributes := []minty.Attr{minty.Class(strings.Join(classes, " "))}
    if btn.OnClick != "" {
        attributes = append(attributes, minty.Attr{Name: "onclick", Value: btn.OnClick})
    }
    
    return b.Button(attributes[0], btn.Text) // Simplified for example
}

// Composite components built from base components
type SearchForm struct {
    Placeholder   string
    SubmitText    string
    OnSubmit      string
}

func (sf SearchForm) Render(b *Builder) Node {
    return b.Form(minty.Class("search-form"),
        b.Div(minty.Class("search-input-group"),
            b.Input(
                minty.Type("text"),
                minty.Name("query"),
                minty.Placeholder(sf.Placeholder),
                minty.Class("search-input"),
            ),
            Button{
                Text:    sf.SubmitText,
                Variant: ButtonPrimary,
                Size:    ButtonMedium,
            }.Render(b),
        ),
    )
}

// Page-level components
type ProductListPage struct {
    Products    []Product
    SearchForm  SearchForm
    Pagination  Pagination
}

func (plp ProductListPage) Render(b *Builder) Node {
    return b.Div(minty.Class("product-list-page"),
        b.Header(minty.Class("page-header"),
            b.H1("Our Products"),
            plp.SearchForm.Render(b),
        ),
        b.Main(minty.Class("products-main"),
            plp.renderProductGrid(b),
            plp.Pagination.Render(b),
        ),
    )
}

func (plp ProductListPage) renderProductGrid(b *Builder) Node {
    if len(plp.Products) == 0 {
        return b.Div(minty.Class("empty-state"),
            b.P("No products found."),
        )
    }
    
    var productCards []Node
    for _, product := range plp.Products {
        productCards = append(productCards, ProductCard{
            Product:     product,
            ShowPrice:   true,
            ShowReviews: true,
        }.Render(b))
    }
    
    return b.Div(minty.Class("product-grid"), productCards...)
}
```

Hierarchical component systems enable clear organization and reuse while maintaining type safety and composability throughout the application architecture.

---

## **Fragment Rendering for HTMX Integration**

### **Fragment Architecture Design**

HTMX applications require the ability to render both complete pages and individual page fragments for dynamic updates. Minty's fragment system provides first-class support for this dual rendering requirement.

```go
// Template with fragment support
type FragmentedTemplate struct {
    Main      H
    Fragments map[string]H
}

func NewFragmentedTemplate(main H) *FragmentedTemplate {
    return &FragmentedTemplate{
        Main:      main,
        Fragments: make(map[string]H),
    }
}

func (ft *FragmentedTemplate) WithFragment(name string, fragment H) *FragmentedTemplate {
    ft.Fragments[name] = fragment
    return ft
}

func (ft *FragmentedTemplate) Render() string {
    var buf strings.Builder
    minty.Render(ft.Main(minty.NewBuilder()), &buf)
    return buf.String()
}

func (ft *FragmentedTemplate) RenderFragment(name string) string {
    fragment, exists := ft.Fragments[name]
    if !exists {
        return ""
    }
    
    var buf strings.Builder
    minty.Render(fragment(minty.NewBuilder()), &buf)
    return buf.String()
}

func (ft *FragmentedTemplate) HasFragment(name string) bool {
    _, exists := ft.Fragments[name]
    return exists
}

// Fragment-aware component
func userDashboard(user User, recentPosts []Post, notifications []Notification) *FragmentedTemplate {
    main := func(b *Builder) Node {
        return b.Div(minty.Class("dashboard"),
            b.Header(minty.Class("dashboard-header"),
                b.H1("Welcome back, " + user.Name),
                b.Div(minty.Id("notification-area")),
            ),
            b.Main(minty.Class("dashboard-main"),
                b.Section(minty.Id("recent-posts"),
                    b.H2("Recent Posts"),
                    b.Div(minty.Id("posts-list")),
                ),
                b.Section(
                    b.H2("Quick Actions"),
                    quickActionsPanel(user)(b),
                ),
            ),
        )
    }
    
    return NewFragmentedTemplate(main).
        WithFragment("notifications", notificationsList(notifications)).
        WithFragment("posts", recentPostsList(recentPosts)).
        WithFragment("quick-actions", quickActionsPanel(user))
}

// HTTP handler integration
func handleDashboard(w http.ResponseWriter, r *http.Request) {
    user := getCurrentUser(r)
    posts := getRecentPosts(user.ID)
    notifications := getNotifications(user.ID)
    
    dashboard := userDashboard(user, posts, notifications)
    
    // Check if this is a fragment request
    if fragment := r.URL.Query().Get("fragment"); fragment != "" {
        if dashboard.HasFragment(fragment) {
            w.Header().Set("Content-Type", "text/html")
            w.Write([]byte(dashboard.RenderFragment(fragment)))
            return
        }
        http.Error(w, "Fragment not found", http.StatusNotFound)
        return
    }
    
    // Render full page
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(dashboard.Render()))
}
```

Fragment architecture design enables seamless integration with HTMX while maintaining clean separation between full page rendering and partial updates.

### **Dynamic Fragment Updates**

HTMX applications often require fragments that can update based on changing data or user interactions. Minty's fragment system supports dynamic fragment generation and caching strategies.

```go
// Fragment cache for performance
type FragmentCache struct {
    cache map[string]CachedFragment
    mutex sync.RWMutex
}

type CachedFragment struct {
    Content   string
    GeneratedAt time.Time
    TTL       time.Duration
}

func NewFragmentCache() *FragmentCache {
    return &FragmentCache{
        cache: make(map[string]CachedFragment),
    }
}

func (fc *FragmentCache) Get(key string) (string, bool) {
    fc.mutex.RLock()
    defer fc.mutex.RUnlock()
    
    fragment, exists := fc.cache[key]
    if !exists {
        return "", false
    }
    
    if time.Since(fragment.GeneratedAt) > fragment.TTL {
        return "", false
    }
    
    return fragment.Content, true
}

func (fc *FragmentCache) Set(key, content string, ttl time.Duration) {
    fc.mutex.Lock()
    defer fc.mutex.Unlock()
    
    fc.cache[key] = CachedFragment{
        Content:     content,
        GeneratedAt: time.Now(),
        TTL:         ttl,
    }
}

// Dynamic fragment generator
type DynamicFragments struct {
    cache    *FragmentCache
    user     User
    database Database
}

func NewDynamicFragments(user User, db Database) *DynamicFragments {
    return &DynamicFragments{
        cache:    NewFragmentCache(),
        user:     user,
        database: db,
    }
}

func (df *DynamicFragments) NotificationCount() H {
    return func(b *Builder) Node {
        cacheKey := fmt.Sprintf("notification-count-%d", df.user.ID)
        
        if cached, found := df.cache.Get(cacheKey); found {
            return minty.RawHTML(cached)
        }
        
        count := df.database.GetUnreadNotificationCount(df.user.ID)
        
        content := b.Span(
            minty.Class("notification-badge"),
            minty.Id("notification-count"),
            strconv.Itoa(count),
        )
        
        var buf strings.Builder
        minty.Render(content, &buf)
        rendered := buf.String()
        
        df.cache.Set(cacheKey, rendered, 5*time.Minute)
        
        return minty.RawHTML(rendered)
    }
}

func (df *DynamicFragments) LiveFeed() H {
    return func(b *Builder) Node {
        posts := df.database.GetRecentPosts(df.user.ID, 10)
        
        if len(posts) == 0 {
            return b.Div(minty.Class("empty-feed"),
                b.P("No recent activity"),
            )
        }
        
        var feedItems []Node
        for _, post := range posts {
            feedItems = append(feedItems, 
                feedItem(post, df.user)(b),
            )
        }
        
        return b.Div(
            minty.Class("live-feed"),
            minty.Id("live-feed"),
            feedItems...,
        )
    }
}

// HTMX-specific fragment helpers
func HTMXFragment(fragmentName string, template H) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Id(fragmentName),
            minty.Attr{Name: "hx-swap-oob", Value: "true"},
            template(b),
        )
    }
}

func HTMXUpdate(target string, content H) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Id(target),
            content(b),
        )
    }
}

// Usage with HTMX attributes
liveNotifications := HTMXFragment("notifications", func(b *Builder) Node {
    return b.Div(
        minty.Attr{Name: "hx-get", Value: "/api/notifications"},
        minty.Attr{Name: "hx-trigger", Value: "every 30s"},
        minty.Attr{Name: "hx-target", Value: "#notifications"},
        df.NotificationCount()(b),
    )
})
```

Dynamic fragment updates enable real-time user interfaces while maintaining performance through intelligent caching and selective updates.

---

## **Conditional Rendering and Control Flow**

### **Conditional Rendering Patterns**

Web applications require sophisticated conditional rendering based on user permissions, application state, and data availability. Minty provides several patterns for implementing conditional logic within templates.

```go
// Basic conditional rendering helper
func If(condition bool, template H) H {
    if condition {
        return template
    }
    return func(b *Builder) Node { return minty.Fragment() }
}

func IfElse(condition bool, trueTemplate, falseTemplate H) H {
    if condition {
        return trueTemplate
    }
    return falseTemplate
}

func Unless(condition bool, template H) H {
    return If(!condition, template)
}

// Advanced conditional rendering with multiple conditions
func Switch[T comparable](value T, cases map[T]H, defaultCase H) H {
    if template, exists := cases[value]; exists {
        return template
    }
    if defaultCase != nil {
        return defaultCase
    }
    return func(b *Builder) Node { return minty.Fragment() }
}

// Permission-based conditional rendering
func IfPermission(user User, permission string, template H) H {
    return If(user.HasPermission(permission), template)
}

func IfRole(user User, role UserRole, template H) H {
    return If(user.Role == role, template)
}

func IfAuthenticated(user User, template H) H {
    return If(user.IsAuthenticated(), template)
}

// Usage examples
adminPanel := IfPermission(currentUser, "admin", func(b *Builder) Node {
    return b.Div(minty.Class("admin-panel"),
        b.H2("Administration"),
        b.A(minty.Href("/admin/users"), "Manage Users"),
        b.A(minty.Href("/admin/settings"), "Settings"),
    )
})

userGreeting := IfElse(
    currentUser.IsAuthenticated(),
    func(b *Builder) Node {
        return b.P("Welcome back, " + currentUser.Name + "!")
    },
    func(b *Builder) Node {
        return b.P("Welcome! Please sign in to continue.")
    },
)

statusBadge := Switch(order.Status, map[OrderStatus]H{
    OrderPending: func(b *Builder) Node {
        return b.Span(minty.Class("badge badge-warning"), "Pending")
    },
    OrderShipped: func(b *Builder) Node {
        return b.Span(minty.Class("badge badge-info"), "Shipped")
    },
    OrderDelivered: func(b *Builder) Node {
        return b.Span(minty.Class("badge badge-success"), "Delivered")
    },
}, func(b *Builder) Node {
    return b.Span(minty.Class("badge badge-secondary"), "Unknown")
})
```

Conditional rendering patterns provide flexible control flow while maintaining template readability and composability.

### **Iteration and Collection Rendering**

Dynamic applications frequently need to render collections of data with consistent formatting and interactive behavior. Minty's iteration helpers provide type-safe collection rendering with flexible customization options.

```go
// Generic iteration helper with type safety
func Each[T any](items []T, renderer func(T) H) []H {
    var templates []H
    for _, item := range items {
        templates = append(templates, renderer(item))
    }
    return templates
}

// Iteration with index access
func EachWithIndex[T any](items []T, renderer func(int, T) H) []H {
    var templates []H
    for i, item := range items {
        templates = append(templates, renderer(i, item))
    }
    return templates
}

// Filtered iteration
func EachWhere[T any](items []T, predicate func(T) bool, renderer func(T) H) []H {
    var templates []H
    for _, item := range items {
        if predicate(item) {
            templates = append(templates, renderer(item))
        }
    }
    return templates
}

// Grouped iteration with separators
func EachWithSeparator[T any](items []T, renderer func(T) H, separator H) []H {
    if len(items) == 0 {
        return []H{}
    }
    
    var templates []H
    for i, item := range items {
        templates = append(templates, renderer(item))
        
        if i < len(items)-1 && separator != nil {
            templates = append(templates, separator)
        }
    }
    return templates
}

// Usage examples
productList := func(b *Builder) Node {
    productTemplates := Each(products, func(product Product) H {
        return func(b *Builder) Node {
            return b.Div(minty.Class("product-card"),
                b.Img(minty.Src(product.ImageURL), minty.Alt(product.Name)),
                b.H3(product.Name),
                b.P(fmt.Sprintf("$%.2f", product.Price)),
            )
        }
    })
    
    var nodes []Node
    for _, template := range productTemplates {
        nodes = append(nodes, template(b))
    }
    
    return b.Div(minty.Class("product-grid"), nodes...)
}

navigationMenu := func(b *Builder) Node {
    menuItems := EachWithSeparator(
        navigationItems,
        func(item NavigationItem) H {
            return func(b *Builder) Node {
                return b.A(minty.Href(item.URL), item.Title)
            }
        },
        func(b *Builder) Node {
            return b.Span(minty.Class("nav-separator"), " | ")
        },
    )
    
    var nodes []Node
    for _, template := range menuItems {
        nodes = append(nodes, template(b))
    }
    
    return b.Nav(minty.Class("main-navigation"), nodes...)
}

adminUserList := func(b *Builder) Node {
    userRows := EachWithIndex(users, func(index int, user User) H {
        return func(b *Builder) Node {
            rowClass := "user-row"
            if index%2 == 1 {
                rowClass += " alternate"
            }
            
            return b.Tr(minty.Class(rowClass),
                b.Td(user.Name),
                b.Td(user.Email),
                b.Td(string(user.Role)),
                b.Td(
                    IfPermission(currentUser, "user.edit", func(b *Builder) Node {
                        return b.A(
                            minty.Href("/admin/users/"+strconv.Itoa(user.ID)),
                            "Edit",
                        )
                    })(b),
                ),
            )
        }
    })
    
    var tableRows []Node
    for _, template := range userRows {
        tableRows = append(tableRows, template(b))
    }
    
    return b.Table(minty.Class("users-table"),
        b.Thead(
            b.Tr(
                b.Th("Name"),
                b.Th("Email"),
                b.Th("Role"),
                b.Th("Actions"),
            ),
        ),
        b.Tbody(tableRows...),
    )
}
```

Iteration and collection rendering patterns enable efficient processing of dynamic data while maintaining type safety and template composability.

---

## **Template Factories and Component Libraries**

### **Component Library Architecture**

Large applications benefit from organized component libraries that provide consistent, reusable building blocks. Minty's template factory system enables the creation of comprehensive component libraries with clear interfaces and extensibility.

```go
// Component library organization
package components

// Base component interfaces
type Renderable interface {
    Render(*Builder) Node
}

type Configurable[T any] interface {
    WithConfig(T) Renderable
}

// Button component family
type ButtonConfig struct {
    Variant     ButtonVariant
    Size        ButtonSize
    Disabled    bool
    Loading     bool
    Icon        string
    FullWidth   bool
    OnClick     string
}

type Button struct {
    Text   string
    Config ButtonConfig
}

func NewButton(text string) *Button {
    return &Button{
        Text: text,
        Config: ButtonConfig{
            Variant: ButtonPrimary,
            Size:    ButtonMedium,
        },
    }
}

func (b *Button) WithVariant(variant ButtonVariant) *Button {
    b.Config.Variant = variant
    return b
}

func (b *Button) WithSize(size ButtonSize) *Button {
    b.Config.Size = size
    return b
}

func (b *Button) WithIcon(icon string) *Button {
    b.Config.Icon = icon
    return b
}

func (b *Button) Disabled() *Button {
    b.Config.Disabled = true
    return b
}

func (b *Button) Loading() *Button {
    b.Config.Loading = true
    return b
}

func (b *Button) FullWidth() *Button {
    b.Config.FullWidth = true
    return b
}

func (b *Button) OnClick(handler string) *Button {
    b.Config.OnClick = handler
    return b
}

func (btn *Button) Render(b *Builder) Node {
    classes := []string{"btn"}
    
    switch btn.Config.Variant {
    case ButtonPrimary:
        classes = append(classes, "btn-primary")
    case ButtonSecondary:
        classes = append(classes, "btn-secondary")
    case ButtonDanger:
        classes = append(classes, "btn-danger")
    }
    
    switch btn.Config.Size {
    case ButtonSmall:
        classes = append(classes, "btn-sm")
    case ButtonLarge:
        classes = append(classes, "btn-lg")
    }
    
    if btn.Config.Disabled {
        classes = append(classes, "btn-disabled")
    }
    
    if btn.Config.FullWidth {
        classes = append(classes, "btn-full-width")
    }
    
    attributes := []minty.Attr{minty.Class(strings.Join(classes, " "))}
    
    if btn.Config.Disabled {
        attributes = append(attributes, minty.Attr{Name: "disabled", Value: "true"})
    }
    
    if btn.Config.OnClick != "" {
        attributes = append(attributes, minty.Attr{Name: "onclick", Value: btn.Config.OnClick})
    }
    
    children := []Node{}
    
    if btn.Config.Icon != "" {
        children = append(children, 
            b.I(minty.Class("icon icon-"+btn.Config.Icon)))
    }
    
    if btn.Config.Loading {
        children = append(children,
            b.Span(minty.Class("loading-spinner")))
    } else {
        children = append(children, &TextNode{Content: btn.Text})
    }
    
    return b.Button(attributes[0], children...) // Simplified for example
}

// Usage examples
submitButton := NewButton("Submit").
    WithVariant(ButtonPrimary).
    WithIcon("check").
    OnClick("submitForm()")

cancelButton := NewButton("Cancel").
    WithVariant(ButtonSecondary).
    OnClick("cancelAction()")

loadingButton := NewButton("Processing...").
    Loading().
    Disabled()
```

Component library architecture provides a foundation for building consistent, maintainable user interfaces while preserving flexibility for customization and extension.

### **Template Derivation and Inheritance**

Complex component libraries often require template derivation where specialized components build upon base templates while adding specific functionality. Minty's composition system supports elegant template inheritance patterns.

```go
// Base form component
type BaseForm struct {
    Action     string
    Method     string
    EncType    string
    ClassName  string
    Fields     []FormField
    Validation ValidationRules
}

func (bf *BaseForm) Render(b *Builder) Node {
    attributes := []minty.Attr{
        minty.Action(bf.Action),
        minty.Method(bf.Method),
    }
    
    if bf.EncType != "" {
        attributes = append(attributes, minty.Attr{Name: "enctype", Value: bf.EncType})
    }
    
    if bf.ClassName != "" {
        attributes = append(attributes, minty.Class(bf.ClassName))
    }
    
    children := []Node{}
    
    for _, field := range bf.Fields {
        children = append(children, bf.renderField(field, b))
    }
    
    children = append(children, bf.renderSubmitSection(b))
    
    return b.Form(attributes[0], children...) // Simplified
}

func (bf *BaseForm) renderField(field FormField, b *Builder) Node {
    return b.Div(minty.Class("form-field"),
        b.Label(minty.For(field.Name), field.Label),
        field.Render(b),
        bf.renderFieldError(field, b),
    )
}

func (bf *BaseForm) renderSubmitSection(b *Builder) Node {
    return b.Div(minty.Class("form-submit"),
        NewButton("Submit").
            WithVariant(ButtonPrimary).
            Render(b),
    )
}

// Derived login form
type LoginForm struct {
    BaseForm
    RememberMe     bool
    ForgotPassword bool
}

func NewLoginForm(action string) *LoginForm {
    return &LoginForm{
        BaseForm: BaseForm{
            Action:    action,
            Method:    "POST",
            ClassName: "login-form",
            Fields: []FormField{
                TextField{
                    Name:        "email",
                    Label:       "Email Address",
                    Type:        "email",
                    Required:    true,
                    Placeholder: "Enter your email",
                },
                TextField{
                    Name:        "password",
                    Label:       "Password",
                    Type:        "password",
                    Required:    true,
                    Placeholder: "Enter your password",
                },
            },
        },
        RememberMe:     true,
        ForgotPassword: true,
    }
}

func (lf *LoginForm) Render(b *Builder) Node {
    // Start with base form rendering
    baseForm := lf.BaseForm.Render(b)
    
    // Add login-specific elements
    additionalElements := []Node{}
    
    if lf.RememberMe {
        additionalElements = append(additionalElements,
            b.Div(minty.Class("form-field checkbox-field"),
                b.Label(
                    b.Input(
                        minty.Type("checkbox"),
                        minty.Name("remember_me"),
                        minty.Value("1"),
                    ),
                    " Remember me",
                ),
            ),
        )
    }
    
    if lf.ForgotPassword {
        additionalElements = append(additionalElements,
            b.Div(minty.Class("form-links"),
                b.A(minty.Href("/forgot-password"), "Forgot your password?"),
            ),
        )
    }
    
    // Combine base form with additional elements
    return b.Div(minty.Class("login-form-wrapper"),
        baseForm,
        additionalElements...,
    )
}

// Derived contact form with different structure
type ContactForm struct {
    BaseForm
    ShowSubjectField bool
    AllowAttachments bool
    RequiredDisclaimer string
}

func NewContactForm(action string) *ContactForm {
    return &ContactForm{
        BaseForm: BaseForm{
            Action:    action,
            Method:    "POST",
            ClassName: "contact-form",
            Fields: []FormField{
                TextField{
                    Name:     "name",
                    Label:    "Your Name",
                    Required: true,
                },
                TextField{
                    Name:     "email",
                    Label:    "Email Address",
                    Type:     "email",
                    Required: true,
                },
                TextAreaField{
                    Name:     "message",
                    Label:    "Message",
                    Required: true,
                    Rows:     5,
                },
            },
        },
        ShowSubjectField: true,
        AllowAttachments: false,
    }
}

func (cf *ContactForm) Render(b *Builder) Node {
    // Modify fields based on configuration
    if cf.ShowSubjectField {
        subjectField := TextField{
            Name:     "subject",
            Label:    "Subject",
            Required: true,
        }
        
        // Insert subject field after email
        fields := make([]FormField, 0, len(cf.BaseForm.Fields)+1)
        fields = append(fields, cf.BaseForm.Fields[:2]...)
        fields = append(fields, subjectField)
        fields = append(fields, cf.BaseForm.Fields[2:]...)
        cf.BaseForm.Fields = fields
    }
    
    baseForm := cf.BaseForm.Render(b)
    
    additionalElements := []Node{}
    
    if cf.RequiredDisclaimer != "" {
        additionalElements = append(additionalElements,
            b.P(minty.Class("form-disclaimer"), cf.RequiredDisclaimer),
        )
    }
    
    return b.Div(minty.Class("contact-form-wrapper"),
        baseForm,
        additionalElements...,
    )
}
```

Template derivation and inheritance enable the construction of specialized components that build upon common base functionality while adding domain-specific features and customizations.

---

## **Composition Patterns Across the Minty System**

The composition patterns described in this document extend throughout the entire Minty System, enabling sophisticated applications that combine business domains, iterators, and themes:

### **Business Domain Component Composition**

Business domain presentation layers use the same composition patterns:

```go
// Finance domain components compose naturally
func FinancialDashboard(theme Theme, data mifi.DashboardData) mi.H {
    return theme.Dashboard("Financial Dashboard",
        // Sidebar component composition
        FinanceSidebar(theme),
        
        // Main content composition with business components
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Class("dashboard-main"),
                // Composed metric cards
                MetricsSection(theme, data.Metrics),
                
                // Account summary composition
                AccountsSummary(theme, data.Accounts),
                
                // Recent activity composition
                RecentActivity(theme, data.RecentTransactions, data.PendingInvoices),
            )
        },
    )
}

// Individual business components follow same patterns
func AccountsSummary(theme Theme, accounts []mifi.Account) mi.H {
    return theme.Section("Accounts",
        // Iterator-based component composition
        miex.Map(accounts, func(account mifi.Account) mi.H {
            return AccountSummaryCard(theme, account)
        }),
    )
}
```

### **Iterator-Powered Component Composition**

Iterators enable powerful composition patterns for data-driven components:

```go
// Template factories with iterator integration
func TableFactory[T any](headers []string, data []T, 
    rowRenderer func(T) []string) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Table(mi.Class("data-table"),
            b.Thead(
                b.Tr(
                    miex.Map(headers, func(h string) mi.H {
                        return func(b *mi.Builder) mi.Node {
                            return b.Th(h)
                        }
                    })...,
                ),
            ),
            b.Tbody(
                // Iterator composition in table rows
                miex.Map(data, func(item T) mi.H {
                    cells := rowRenderer(item)
                    return func(b *mi.Builder) mi.Node {
                        return b.Tr(
                            miex.Map(cells, func(cell string) mi.H {
                                return func(b *mi.Builder) mi.Node {
                                    return b.Td(cell)
                                }
                            })...,
                        )
                    }
                })...,
            ),
        )
    }
}

// Usage with business domain data
usersTable := TableFactory(
    []string{"Name", "Email", "Status"},
    users,
    func(u User) []string {
        return []string{u.Name, u.Email, u.Status}
    },
)
```

### **Theme-Aware Component Composition**

Theme implementations enable composable styling across the system:

```go
// Cross-domain theme composition
func MultiDomainDashboard(services *ApplicationServices, theme Theme) mi.H {
    return theme.Layout("Business Dashboard",
        // Composed navigation across domains
        theme.Nav([]NavItem{
            {Text: "Finance", URL: "/finance", Icon: "dollar"},
            {Text: "Logistics", URL: "/logistics", Icon: "truck"},
            {Text: "E-commerce", URL: "/ecommerce", Icon: "shopping"},
        }),
        
        // Composed main content with domain-specific components
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Class("dashboard-grid"),
                // Each domain contributes composed components
                mintyfinui.QuickSummary(theme, services.Finance),
                mintymoveui.ShipmentStatus(theme, services.Logistics),
                mintycartui.OrderActivity(theme, services.Ecommerce),
            )
        },
    )
}

// Domain-specific components maintain composition patterns
func QuickSummary(theme Theme, financeService *mifi.FinanceService) mi.H {
    dashboardData := mifi.PrepareDashboardData(financeService)
    
    return theme.Card("Financial Summary",
        func(b *mi.Builder) mi.Node {
            return b.Div(
                // Composed metric components
                MetricRow(theme, "Total Balance", dashboardData.TotalBalance),
                MetricRow(theme, "Pending Invoices", dashboardData.PendingInvoicesCount),
                MetricRow(theme, "Recent Transactions", dashboardData.RecentTransactionsCount),
            )
        },
    )
}
```

### **Complex Application Composition**

Large applications compose all these elements while maintaining clear patterns:

```go
// Complete application composition
func BusinessApplication(services *ApplicationServices, theme Theme, 
    user User, currentRoute string) mi.H {
    
    return theme.Document("Business Management System",
        // Application shell composition
        ApplicationShell(theme, user,
            // Route-based content composition
            RouteContent(services, theme, currentRoute,
                // Domain-specific route handlers
                map[string]mi.H{
                    "/finance":   FinanceModule(services.Finance, theme),
                    "/logistics": LogisticsModule(services.Logistics, theme), 
                    "/ecommerce": EcommerceModule(services.Ecommerce, theme),
                    "/dashboard": MultiDomainDashboard(services, theme),
                },
            ),
        ),
    )
}
```

These composition patterns enable the Minty System to scale from simple components to sophisticated business applications while maintaining consistent, understandable patterns throughout.

---

## **Conclusion**

Minty's template composition system represents a comprehensive approach to building maintainable, scalable web applications through functional programming principles and type-safe component architecture. The patterns explored in this document provide developers with powerful tools for organizing complex user interfaces while maintaining clarity and reusability.

The combination of parameterized templates, layout systems, component composition, fragment rendering, conditional logic, and template factories creates a robust foundation for modern web development. This foundation enables developers to construct sophisticated applications while preserving the simplicity and performance characteristics that make Go an attractive platform for web development.

These composition patterns scale throughout the **entire Minty System**, enabling consistent approaches for business domain presentation, iterator-powered data processing, theme-based styling, and complex multi-domain applications. The same principles that make simple HTML components maintainable also make sophisticated business applications productive and scalable.

As applications grow in complexity, these composition patterns provide the architectural foundation necessary to maintain code quality and developer productivity. The type safety inherent in Minty's design ensures that composition errors are caught at compile time, while the functional approach enables powerful abstractions without sacrificing performance or clarity.

The future of web development lies in combining the simplicity of server-side rendering with the interactivity of modern JavaScript frameworks. Minty's composition system provides the tools necessary to build this future today, enabling developers to create rich, interactive web applications using only Go and HTML.