# **Minty Documentation - Part 6**
## **HTMX Integration & Interactivity**
### *"JavaScript-Free Dynamic UIs"*

> **Part of the Minty System**: This document covers HTMX integration patterns that work seamlessly with the entire Minty System. These patterns apply to HTML generation (minty), business domain interactions (mintyfin, mintymove, mintycart), iterator-based data processing (mintyex), theme-based components, and complex multi-domain applications. HTMX enables dynamic behavior while maintaining the server-first philosophy throughout the system.

---

## **Executive Summary**

Modern web applications demand rich interactivity and responsive user experiences, traditionally achieved through complex JavaScript frameworks and extensive client-side state management. Minty challenges this paradigm by demonstrating that sophisticated interactive behaviors can be achieved entirely through server-side rendering combined with declarative HTML attributes provided by HTMX.

This document explores how Minty's HTMX integration transforms the traditional web development model from JavaScript-heavy single-page applications to server-driven hypermedia applications. Rather than treating HTMX as an add-on library, Minty embraces it as a foundational technology, providing first-class integration that makes JavaScript-free development not just possible, but preferable.

The result is a development experience that combines the simplicity of traditional server-side rendering with the interactivity of modern web applications, all while maintaining Go's core values of type safety, performance, and maintainability.

---

## **The JavaScript-Free Philosophy**

### **Redefining Web Application Interactivity**

Traditional web development assumes that interactivity requires JavaScript programming. Minty demonstrates that this assumption is fundamentally incorrect. Through strategic use of HTMX attributes and server-side rendering, applications can achieve sophisticated interactive behaviors while remaining entirely JavaScript-free.

```go
// Traditional JavaScript approach requires:
// - Client-side event handling
// - DOM manipulation logic
// - State synchronization
// - Error handling
// - Loading states
// - API integration

// Minty + HTMX approach achieves the same result:
func liveSearchBox(placeholder string) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("search-container"),
            b.Input(
                minty.Type("text"),
                minty.Name("query"),
                minty.Placeholder(placeholder),
                minty.HtmxGet("/search"),              // Declarative request
                minty.HtmxTarget("#search-results"),   // Declarative target
                minty.HtmxTrigger("keyup changed delay:300ms"), // Declarative timing
                minty.HtmxIndicator("#search-spinner"), // Declarative loading state
            ),
            b.Div(minty.Id("search-spinner"), 
                minty.Class("htmx-indicator"),
                "Searching...",
            ),
            b.Div(minty.Id("search-results"), 
                minty.Class("search-results"),
            ),
        )
    }
}
```

This approach eliminates entire categories of complexity while providing identical user experience to JavaScript-based alternatives.

### **The 98% Rule**

Practical experience demonstrates that approximately 98% of web application functionality can be achieved without writing JavaScript code. The remaining 2% typically involves highly specialized interactions that genuinely require client-side processing.

**JavaScript-Free Capabilities**:
- Form submission and validation with real-time feedback
- Live search with debouncing and result filtering
- Dynamic content updates and partial page replacements
- Modal dialogs and overlay interfaces
- Real-time notifications and status updates
- Interactive data tables with sorting and filtering
- Progressive disclosure and accordion interfaces
- Auto-saving forms and draft persistence
- File upload with progress indication
- Shopping cart and e-commerce interactions

**Edge Cases Requiring Minimal JavaScript**:
- Complex animations requiring frame-by-frame control
- Real-time collaborative editing with conflict resolution
- Advanced drawing or diagramming interfaces
- WebRTC video calling functionality
- Complex mathematical visualizations requiring immediate responsiveness

The key insight is that the vast majority of web applications fall entirely within the JavaScript-free 98%, making Minty's approach broadly applicable.

### **Performance and Simplicity Benefits**

JavaScript-free applications provide significant advantages across multiple dimensions:

**Bundle Size**: No JavaScript bundle means zero download, parse, and execution overhead. Applications become interactive immediately upon HTML load completion.

**Debugging Simplicity**: All application logic executes on the server in a single runtime environment. Debugging involves standard Go debugging tools rather than coordinating between client and server environments.

**Security**: Eliminating client-side code removes entire categories of security vulnerabilities including XSS attacks through third-party dependencies and client-side state manipulation.

**Caching**: Server-rendered HTML fragments cache effectively at multiple levels, providing better performance characteristics than JSON API responses requiring client-side processing.

---

## **Built-in HTMX Attribute Helpers**

### **First-Class HTMX Integration**

Minty treats HTMX as a foundational technology rather than an optional add-on. This philosophical choice manifests in comprehensive built-in helpers that make HTMX attributes as natural to use as standard HTML attributes.

```go
// Manual HTMX attribute approach (verbose and error-prone)
func manualDeleteButton(itemID int) H {
    return func(b *Builder) Node {
        return b.Button(
            minty.Class("btn btn-danger"),
            minty.Attr("hx-delete", fmt.Sprintf("/items/%d", itemID)),
            minty.Attr("hx-target", "#item-list"),
            minty.Attr("hx-swap", "outerHTML"),
            minty.Attr("hx-confirm", "Are you sure you want to delete this item?"),
            minty.Attr("hx-indicator", "#delete-spinner"),
            "Delete",
        )
    }
}

// Minty's built-in HTMX helpers (concise and type-safe)
func deleteButton(itemID int) H {
    return func(b *Builder) Node {
        return b.Button(
            minty.Class("btn btn-danger"),
            minty.HtmxDelete(fmt.Sprintf("/items/%d", itemID)),
            minty.HtmxTarget("#item-list"),
            minty.HtmxSwap("outerHTML"),
            minty.HtmxConfirm("Are you sure you want to delete this item?"),
            minty.HtmxIndicator("#delete-spinner"),
            "Delete",
        )
    }
}
```

The built-in helpers provide type safety, reduce verbosity, and make HTMX attributes discoverable through IDE autocompletion.

### **Comprehensive HTMX Attribute Coverage**

Minty provides helpers for all HTMX attributes, organized by functional category:

#### HTTP Method Helpers

```go
// HTTP method attribute helpers
func (b *Builder) HtmxGet(url string) Attribute {
    return StringAttribute{Name: "hx-get", Value: url}
}

func (b *Builder) HtmxPost(url string) Attribute {
    return StringAttribute{Name: "hx-post", Value: url}
}

func (b *Builder) HtmxPut(url string) Attribute {
    return StringAttribute{Name: "hx-put", Value: url}
}

func (b *Builder) HtmxPatch(url string) Attribute {
    return StringAttribute{Name: "hx-patch", Value: url}
}

func (b *Builder) HtmxDelete(url string) Attribute {
    return StringAttribute{Name: "hx-delete", Value: url}
}

// Usage in interactive components
func crudButtons(item Item) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("action-buttons"),
            b.Button(
                minty.Class("btn btn-primary"),
                minty.HtmxGet(fmt.Sprintf("/items/%d/edit", item.ID)),
                minty.HtmxTarget("#edit-form"),
                "Edit",
            ),
            b.Button(
                minty.Class("btn btn-success"),
                minty.HtmxPut(fmt.Sprintf("/items/%d", item.ID)),
                minty.HtmxTarget("#item-list"),
                "Save",
            ),
            b.Button(
                minty.Class("btn btn-danger"),
                minty.HtmxDelete(fmt.Sprintf("/items/%d", item.ID)),
                minty.HtmxTarget("#item-list"),
                minty.HtmxConfirm("Delete this item?"),
                "Delete",
            ),
        )
    }
}
```

#### Target and Swapping Helpers

```go
// Targeting and swapping attribute helpers
func (b *Builder) HtmxTarget(selector string) Attribute {
    return StringAttribute{Name: "hx-target", Value: selector}
}

func (b *Builder) HtmxSwap(strategy string) Attribute {
    return StringAttribute{Name: "hx-swap", Value: strategy}
}

func (b *Builder) HtmxSwapOOB(value string) Attribute {
    return StringAttribute{Name: "hx-swap-oob", Value: value}
}

// Predefined swap strategies for common patterns
func HtmxInnerHTML() Attribute {
    return StringAttribute{Name: "hx-swap", Value: "innerHTML"}
}

func HtmxOuterHTML() Attribute {
    return StringAttribute{Name: "hx-swap", Value: "outerHTML"}
}

func HtmxAfterbegin() Attribute {
    return StringAttribute{Name: "hx-swap", Value: "afterbegin"}
}

func HtmxBeforeend() Attribute {
    return StringAttribute{Name: "hx-swap", Value: "beforeend"}
}

// Usage in dynamic content updates
func commentForm(postID int) H {
    return func(b *Builder) Node {
        return b.Form(
            minty.HtmxPost(fmt.Sprintf("/posts/%d/comments", postID)),
            minty.HtmxTarget("#comments-list"),
            HtmxBeforeend(), // Add new comments to end of list
            minty.HtmxSwapOOB("innerHTML:#comment-form"), // Clear form after submit
            
            b.Textarea(
                minty.Name("content"),
                minty.Placeholder("Write your comment..."),
                minty.Required(),
            ),
            b.Button(
                minty.Type("submit"),
                minty.Class("btn btn-primary"),
                "Post Comment",
            ),
        )
    }
}
```

#### Trigger and Timing Helpers

```go
// Event trigger attribute helpers
func (b *Builder) HtmxTrigger(events string) Attribute {
    return StringAttribute{Name: "hx-trigger", Value: events}
}

// Predefined trigger patterns for common interactions
func HtmxTriggerClick() Attribute {
    return StringAttribute{Name: "hx-trigger", Value: "click"}
}

func HtmxTriggerSubmit() Attribute {
    return StringAttribute{Name: "hx-trigger", Value: "submit"}
}

func HtmxTriggerChange() Attribute {
    return StringAttribute{Name: "hx-trigger", Value: "change"}
}

func HtmxTriggerKeyup(delay string) Attribute {
    return StringAttribute{
        Name:  "hx-trigger", 
        Value: fmt.Sprintf("keyup changed delay:%s", delay),
    }
}

func HtmxTriggerEvery(interval string) Attribute {
    return StringAttribute{
        Name:  "hx-trigger", 
        Value: fmt.Sprintf("every %s", interval),
    }
}

// Usage in responsive interfaces
func liveStatusIndicator(userID int) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Class("status-indicator"),
            minty.HtmxGet(fmt.Sprintf("/users/%d/status", userID)),
            minty.HtmxTarget("this"),
            HtmxTriggerEvery("10s"), // Poll every 10 seconds
            minty.HtmxSwap("innerHTML"),
            "Checking status...",
        )
    }
}

func autoSaveForm(formID string) H {
    return func(b *Builder) Node {
        return b.Form(
            minty.Id(formID),
            minty.HtmxPost("/auto-save"),
            minty.HtmxTrigger("input changed delay:2s"), // Auto-save after 2s of inactivity
            minty.HtmxTarget("#save-status"),
            
            b.Textarea(
                minty.Name("content"),
                minty.Placeholder("Start writing..."),
            ),
            b.Div(
                minty.Id("save-status"),
                minty.Class("save-indicator"),
            ),
        )
    }
}
```

### **Fluent Interface Patterns**

Minty's HTMX helpers support fluent interface patterns that enable expressive, readable component definitions:

```go
// Fluent HTMX interface for complex interactions
type HTMXBuilder struct {
    element Node
    attrs   []Attribute
}

func (h *HTMXBuilder) Get(url string) *HTMXBuilder {
    h.attrs = append(h.attrs, minty.HtmxGet(url))
    return h
}

func (h *HTMXBuilder) Target(selector string) *HTMXBuilder {
    h.attrs = append(h.attrs, minty.HtmxTarget(selector))
    return h
}

func (h *HTMXBuilder) Trigger(events string) *HTMXBuilder {
    h.attrs = append(h.attrs, minty.HtmxTrigger(events))
    return h
}

func (h *HTMXBuilder) Swap(strategy string) *HTMXBuilder {
    h.attrs = append(h.attrs, minty.HtmxSwap(strategy))
    return h
}

func (h *HTMXBuilder) Confirm(message string) *HTMXBuilder {
    h.attrs = append(h.attrs, minty.HtmxConfirm(message))
    return h
}

func (h *HTMXBuilder) Build() Node {
    // Apply all accumulated attributes to the element
    return h.element // with attributes applied
}

// Usage with fluent interface
func interactiveButton(text string) H {
    return func(b *Builder) Node {
        return NewHTMXBuilder(b.Button(text)).
            Get("/api/data").
            Target("#content").
            Trigger("click").
            Swap("innerHTML").
            Confirm("Are you sure?").
            Build()
    }
}
```

---

## **Request/Response Patterns**

### **Fragment-Oriented Handler Design**

HTMX applications require handlers that can serve both complete pages (for direct navigation) and HTML fragments (for dynamic updates). Minty provides patterns that make this dual-mode serving elegant and maintainable.

```go
// HTMX request detection helper
func isHTMXRequest(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}

func getHTMXTarget(r *http.Request) string {
    return r.Header.Get("HX-Target")
}

func getHTMXTrigger(r *http.Request) string {
    return r.Header.Get("HX-Trigger")
}

// Fragment-oriented handler pattern
func userListHandler(w http.ResponseWriter, r *http.Request) {
    users, err := getUsersFromDatabase()
    if err != nil {
        http.Error(w, "Failed to load users", http.StatusInternalServerError)
        return
    }
    
    if isHTMXRequest(r) {
        // Return fragment for HTMX updates
        fragment := userListFragment(users)
        renderFragment(w, fragment)
        return
    }
    
    // Return complete page for direct navigation
    page := userListPage(users)
    renderPage(w, page)
}

func userListFragment(users []User) H {
    return func(b *Builder) Node {
        if len(users) == 0 {
            return b.Div(minty.Class("empty-state"),
                b.P("No users found."),
                b.A(minty.Href("/users/new"), "Add the first user"),
            )
        }
        
        var userRows []Node
        for _, user := range users {
            userRows = append(userRows, userCard(user)(b))
        }
        
        return b.Div(
            minty.Id("user-list"),
            minty.Class("user-grid"),
            userRows...,
        )
    }
}

func userListPage(users []User) H {
    return func(b *Builder) Node {
        return pageLayout("Users", func(b *Builder) Node {
            return b.Div(minty.Class("users-page"),
                b.Header(minty.Class("page-header"),
                    b.H1("Users"),
                    b.A(
                        minty.Class("btn btn-primary"),
                        minty.Href("/users/new"),
                        "Add User",
                    ),
                ),
                userListFragment(users)(b),
            )
        })(b)
    }
}
```

This pattern ensures that every endpoint can serve both use cases without duplicating template logic.

### **Context-Aware Response Generation**

Advanced HTMX applications benefit from handlers that can adapt their responses based on the context of the request, including the triggering element and target element.

```go
// HTMX context structure
type HTMXContext struct {
    Request     *http.Request
    IsHTMX      bool
    Target      string
    Trigger     string
    TriggerName string
    Boosted     bool
    HistoryRestore bool
}

func NewHTMXContext(r *http.Request) HTMXContext {
    return HTMXContext{
        Request:        r,
        IsHTMX:         r.Header.Get("HX-Request") == "true",
        Target:         r.Header.Get("HX-Target"),
        Trigger:        r.Header.Get("HX-Trigger"),
        TriggerName:    r.Header.Get("HX-Trigger-Name"),
        Boosted:        r.Header.Get("HX-Boosted") == "true",
        HistoryRestore: r.Header.Get("HX-History-Restore-Request") == "true",
    }
}

// Context-aware handler that adapts response based on HTMX context
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    ctx := NewHTMXContext(r)
    user := getCurrentUser(r)
    
    switch {
    case ctx.Target == "#notifications":
        // Request specifically for notifications section
        notifications := getNotifications(user.ID)
        fragment := notificationsFragment(notifications)
        renderFragment(w, fragment)
        
    case ctx.Target == "#recent-activity":
        // Request specifically for activity feed
        activities := getRecentActivity(user.ID)
        fragment := activityFragment(activities)
        renderFragment(w, fragment)
        
    case ctx.TriggerName == "refresh-stats":
        // Triggered by stats refresh button
        stats := getDashboardStats(user.ID)
        fragment := statsFragment(stats)
        renderFragment(w, fragment)
        
    case ctx.IsHTMX:
        // General HTMX request - return main dashboard content
        fragment := dashboardMainContent(user)
        renderFragment(w, fragment)
        
    default:
        // Direct navigation - return complete page
        page := dashboardPage(user)
        renderPage(w, page)
    }
}

func dashboardMainContent(user User) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("dashboard-main"),
            b.Section(minty.Id("stats"),
                statsFragment(getDashboardStats(user.ID))(b),
                b.Button(
                    minty.Class("btn btn-sm"),
                    minty.HtmxGet("/dashboard"),
                    minty.HtmxTarget("#stats"),
                    minty.HtmxTrigger("click"),
                    minty.Attr("hx-trigger-name", "refresh-stats"),
                    "Refresh",
                ),
            ),
            b.Section(minty.Id("notifications"),
                notificationsFragment(getNotifications(user.ID))(b),
            ),
            b.Section(minty.Id("recent-activity"),
                activityFragment(getRecentActivity(user.ID))(b),
            ),
        )
    }
}
```

Context-aware handlers enable sophisticated interactions while maintaining clean separation of concerns.

### **Response Header Management**

HTMX applications often require specific response headers to control client-side behavior. Minty provides helpers for common HTMX response header patterns.

```go
// HTMX response header helpers
func setHTMXTrigger(w http.ResponseWriter, events string) {
    w.Header().Set("HX-Trigger", events)
}

func setHTMXTriggerAfterSettle(w http.ResponseWriter, events string) {
    w.Header().Set("HX-Trigger-After-Settle", events)
}

func setHTMXTriggerAfterSwap(w http.ResponseWriter, events string) {
    w.Header().Set("HX-Trigger-After-Swap", events)
}

func setHTMXRedirect(w http.ResponseWriter, url string) {
    w.Header().Set("HX-Redirect", url)
}

func setHTMXRefresh(w http.ResponseWriter) {
    w.Header().Set("HX-Refresh", "true")
}

// Usage in handlers for coordinated updates
func addToCartHandler(w http.ResponseWriter, r *http.Request) {
    productID := r.FormValue("product_id")
    quantity, _ := strconv.Atoi(r.FormValue("quantity"))
    
    cart := getCartFromSession(r)
    cart.AddItem(productID, quantity)
    saveCartToSession(r, cart)
    
    // Return updated cart fragment
    cartFragment := cartWidgetFragment(cart)
    renderFragment(w, cartFragment)
    
    // Trigger events for other parts of the page
    setHTMXTrigger(w, "cartUpdated")
    
    // Show success notification
    setHTMXTriggerAfterSettle(w, "showNotification")
}

// Handlers can listen for these events
func cartBadgeHandler(w http.ResponseWriter, r *http.Request) {
    cart := getCartFromSession(r)
    badge := cartBadgeFragment(cart.ItemCount())
    renderFragment(w, badge)
}

// JavaScript (minimal) to handle custom events
/*
document.body.addEventListener('cartUpdated', function(event) {
    // Trigger update of cart badge
    htmx.trigger('#cart-badge', 'htmx:trigger');
});
*/
```

---

## **Fragment Updates and Swapping**

### **Strategic Fragment Design**

Effective HTMX applications require thoughtful fragment design that balances granular updates with maintainable code organization. Minty provides patterns for organizing fragments at different levels of granularity.

```go
// Multi-level fragment organization
type PageFragments struct {
    FullPage    H
    MainContent H
    Sidebar     H
    Components  map[string]H
}

func NewPageFragments(data PageData) *PageFragments {
    return &PageFragments{
        FullPage:    fullPageTemplate(data),
        MainContent: mainContentTemplate(data),
        Sidebar:     sidebarTemplate(data),
        Components: map[string]H{
            "user-profile":   userProfileFragment(data.User),
            "notifications":  notificationsFragment(data.Notifications),
            "recent-posts":   recentPostsFragment(data.RecentPosts),
            "activity-feed":  activityFeedFragment(data.Activities),
        },
    }
}

func (pf *PageFragments) Render(target string) H {
    switch target {
    case "":
        return pf.FullPage
    case "#main-content":
        return pf.MainContent
    case "#sidebar":
        return pf.Sidebar
    default:
        if component, exists := pf.Components[strings.TrimPrefix(target, "#")]; exists {
            return component
        }
        return pf.MainContent
    }
}

// Handler using fragment organization
func profilePageHandler(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromPath(r.URL.Path)
    data := getPageData(userID)
    fragments := NewPageFragments(data)
    
    ctx := NewHTMXContext(r)
    
    template := fragments.Render(ctx.Target)
    if ctx.IsHTMX {
        renderFragment(w, template)
    } else {
        renderPage(w, template)
    }
}
```

This organization enables efficient partial updates while maintaining clear code structure.

### **Out-of-Band Updates**

HTMX supports out-of-band (OOB) updates that modify multiple page elements in a single response. Minty provides helpers for coordinating complex multi-element updates.

```go
// Out-of-band update helpers
func withOOBUpdate(main H, oobUpdates ...H) H {
    return func(b *Builder) Node {
        nodes := []Node{main(b)}
        
        for _, update := range oobUpdates {
            nodes = append(nodes, update(b))
        }
        
        return minty.Fragment(nodes...)
    }
}

func oobUpdate(id string, content H) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Id(id),
            minty.HtmxSwapOOB("innerHTML"),
            content(b),
        )
    }
}

func oobReplace(id string, content H) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Id(id),
            minty.HtmxSwapOOB("outerHTML"),
            content(b),
        )
    }
}

// Usage in coordinated updates
func likePostHandler(w http.ResponseWriter, r *http.Request) {
    postID := getPostIDFromForm(r)
    userID := getCurrentUserID(r)
    
    post := likePost(postID, userID)
    user := getUserStats(userID)
    
    // Main update: the like button itself
    likeButton := likeButtonFragment(post, userID)
    
    // Out-of-band updates: other affected elements
    postStats := oobUpdate("post-stats-"+postID, postStatsFragment(post))
    userStats := oobUpdate("user-stats", userStatsFragment(user))
    notification := oobUpdate("notification-area", 
        notificationFragment("Post liked!"))
    
    // Combine main update with OOB updates
    response := withOOBUpdate(likeButton, postStats, userStats, notification)
    renderFragment(w, response)
}

func likeButtonFragment(post Post, userID int) H {
    return func(b *Builder) Node {
        isLiked := post.IsLikedBy(userID)
        
        return b.Button(
            minty.Id(fmt.Sprintf("like-button-%d", post.ID)),
            minty.Class(fmt.Sprintf("like-button %s", 
                ternary(isLiked, "liked", "not-liked"))),
            minty.HtmxPost(fmt.Sprintf("/posts/%d/like", post.ID)),
            minty.HtmxTarget("this"),
            minty.HtmxSwap("outerHTML"),
            
            b.I(minty.Class(fmt.Sprintf("icon %s", 
                ternary(isLiked, "icon-heart-filled", "icon-heart")))),
            fmt.Sprintf(" %d", post.LikeCount),
        )
    }
}
```

Out-of-band updates enable sophisticated interactions that update multiple page elements atomically.

### **Optimistic Updates and Error Handling**

Modern web applications benefit from optimistic updates that provide immediate feedback while handling potential server errors gracefully.

```go
// Optimistic update pattern with error handling
func optimisticUpdateButton(action string, optimisticContent H, errorContent H) H {
    return func(b *Builder) Node {
        return b.Button(
            minty.HtmxPost(action),
            minty.HtmxTarget("this"),
            minty.HtmxSwap("outerHTML"),
            
            // Show optimistic state immediately
            minty.Attr("hx-on", fmt.Sprintf(`htmx:beforeRequest: 
                this.innerHTML = '%s'; 
                this.disabled = true;`,
                renderToString(optimisticContent))),
            
            // Handle errors
            minty.Attr("hx-on", fmt.Sprintf(`htmx:responseError: 
                this.innerHTML = '%s'; 
                this.disabled = false;`,
                renderToString(errorContent))),
            
            "Original Content",
        )
    }
}

// Usage for save operations
func saveButton(documentID int) H {
    optimisticState := func(b *Builder) Node {
        return b.Span(
            minty.Class("saving-state"),
            b.I(minty.Class("spinner")),
            " Saving...",
        )
    }
    
    errorState := func(b *Builder) Node {
        return b.Span(
            minty.Class("error-state"),
            b.I(minty.Class("icon-error")),
            " Save Failed - Retry",
        )
    }
    
    return optimisticUpdateButton(
        fmt.Sprintf("/documents/%d/save", documentID),
        optimisticState,
        errorState,
    )
}

// Server-side error handling
func saveDocumentHandler(w http.ResponseWriter, r *http.Request) {
    documentID := getDocumentIDFromPath(r.URL.Path)
    content := r.FormValue("content")
    
    err := saveDocument(documentID, content)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        errorButton := saveButton(documentID) // Return retry button
        renderFragment(w, errorButton)
        return
    }
    
    // Success state
    successButton := func(b *Builder) Node {
        return b.Span(
            minty.Class("success-state"),
            b.I(minty.Class("icon-check")),
            " Saved",
        )
    }
    
    renderFragment(w, successButton)
}
```

---

## **Form Handling and Validation Feedback**

### **Real-Time Form Validation**

Traditional form validation requires page submission to provide feedback. HTMX enables real-time validation that provides immediate feedback as users interact with form fields.

```go
// Real-time validation pattern
func validatedInput(name, label, placeholder string, validator string) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("form-field"),
            b.Label(minty.For(name), label),
            b.Input(
                minty.Type("text"),
                minty.Name(name),
                minty.Id(name),
                minty.Placeholder(placeholder),
                minty.Class("form-input"),
                
                // Validate on blur and change
                minty.HtmxPost("/validate/" + validator),
                minty.HtmxTarget("#" + name + "-error"),
                minty.HtmxTrigger("blur, change"),
                minty.HtmxSwap("innerHTML"),
                
                // Include field name in request
                minty.Attr("hx-include", "[name='"+name+"']"),
            ),
            b.Div(
                minty.Id(name + "-error"),
                minty.Class("error-message"),
            ),
        )
    }
}

// Server-side validation handlers
func validateEmailHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    
    if email == "" {
        w.WriteHeader(http.StatusOK)
        return // No error for empty field during typing
    }
    
    if !isValidEmail(email) {
        errorMsg := func(b *Builder) Node {
            return b.P(
                minty.Class("text-red-500 text-sm"),
                "Please enter a valid email address",
            )
        }
        renderFragment(w, errorMsg)
        return
    }
    
    // Check if email is already taken
    if emailExists(email) {
        errorMsg := func(b *Builder) Node {
            return b.P(
                minty.Class("text-red-500 text-sm"),
                "This email address is already registered",
            )
        }
        renderFragment(w, errorMsg)
        return
    }
    
    // Success state
    successMsg := func(b *Builder) Node {
        return b.P(
            minty.Class("text-green-500 text-sm"),
            b.I(minty.Class("icon-check")),
            " Email is available",
        )
    }
    renderFragment(w, successMsg)
}

// Usage in registration form
func registrationForm() H {
    return func(b *Builder) Node {
        return b.Form(
            minty.Class("registration-form"),
            minty.HtmxPost("/register"),
            minty.HtmxTarget("#registration-result"),
            
            validatedInput("email", "Email Address", 
                "Enter your email", "email")(b),
            validatedInput("username", "Username", 
                "Choose a username", "username")(b),
            
            b.Div(minty.Class("form-field"),
                b.Label(minty.For("password"), "Password"),
                b.Input(
                    minty.Type("password"),
                    minty.Name("password"),
                    minty.Id("password"),
                    minty.Class("form-input"),
                    
                    // Live password strength feedback
                    minty.HtmxPost("/validate/password-strength"),
                    minty.HtmxTarget("#password-strength"),
                    minty.HtmxTrigger("keyup delay:300ms"),
                ),
                b.Div(minty.Id("password-strength")),
            ),
            
            b.Button(
                minty.Type("submit"),
                minty.Class("btn btn-primary"),
                "Create Account",
            ),
            
            b.Div(minty.Id("registration-result")),
        )
    }
}
```

Real-time validation provides immediate feedback while maintaining server-side validation authority.

### **Progressive Form Enhancement**

Forms benefit from progressive enhancement that adds interactive behaviors while maintaining basic functionality for non-JavaScript environments.

```go
// Progressive form enhancement pattern
func enhancedForm(action string, content H) H {
    return func(b *Builder) Node {
        return b.Form(
            minty.Action(action),
            minty.Method("POST"),
            minty.Class("enhanced-form"),
            
            // HTMX enhancement
            minty.HtmxPost(action),
            minty.HtmxTarget("#form-result"),
            minty.HtmxSwap("innerHTML"),
            
            // Progressive enhancement attributes
            minty.Attr("hx-indicator", "#loading-indicator"),
            minty.Attr("hx-disabled-elt", "find button"),
            
            content(b),
            
            b.Div(
                minty.Id("loading-indicator"),
                minty.Class("htmx-indicator"),
                b.Span(minty.Class("spinner")),
                " Processing...",
            ),
            
            b.Div(minty.Id("form-result")),
        )
    }
}

// Auto-saving draft functionality
func draftEnabledForm(formID string, saveEndpoint string) H {
    return func(b *Builder) Node {
        return b.Form(
            minty.Id(formID),
            minty.Class("draft-enabled-form"),
            
            // Auto-save on input changes
            minty.HtmxPost(saveEndpoint),
            minty.HtmxTarget("#save-status"),
            minty.HtmxTrigger("input changed delay:2s"),
            minty.HtmxSwap("innerHTML"),
            
            b.Div(minty.Class("form-content"),
                b.Textarea(
                    minty.Name("content"),
                    minty.Placeholder("Start writing..."),
                    minty.Class("form-textarea"),
                ),
            ),
            
            b.Div(
                minty.Id("save-status"),
                minty.Class("save-status"),
            ),
            
            b.Div(minty.Class("form-actions"),
                b.Button(
                    minty.Type("submit"),
                    minty.Class("btn btn-primary"),
                    "Publish",
                ),
                b.Button(
                    minty.Type("button"),
                    minty.Class("btn btn-secondary"),
                    minty.HtmxPost(saveEndpoint + "?draft=true"),
                    minty.HtmxTarget("#save-status"),
                    "Save Draft",
                ),
            ),
        )
    }
}

// Server-side draft handling
func saveDraftHandler(w http.ResponseWriter, r *http.Request) {
    content := r.FormValue("content")
    userID := getCurrentUserID(r)
    
    err := saveDraft(userID, content)
    if err != nil {
        errorMsg := func(b *Builder) Node {
            return b.Span(
                minty.Class("save-error"),
                "Failed to save draft",
            )
        }
        renderFragment(w, errorMsg)
        return
    }
    
    timestamp := time.Now().Format("3:04 PM")
    successMsg := func(b *Builder) Node {
        return b.Span(
            minty.Class("save-success"),
            fmt.Sprintf("Draft saved at %s", timestamp),
        )
    }
    renderFragment(w, successMsg)
}
```

Progressive enhancement ensures forms work in all environments while providing enhanced experiences where possible.

---

## **Real-Time Updates with Server-Sent Events**

### **Server-Sent Events Integration**

While HTMX provides excellent support for request-response interactions, real-time updates often require Server-Sent Events (SSE) for efficient server-to-client communication.

```go
// SSE endpoint for real-time updates
func notificationStreamHandler(w http.ResponseWriter, r *http.Request) {
    userID := getCurrentUserID(r)
    
    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }
    
    // Create notification channel for this user
    notifications := make(chan Notification, 10)
    defer close(notifications)
    
    // Register for user notifications
    notificationManager.Subscribe(userID, notifications)
    defer notificationManager.Unsubscribe(userID, notifications)
    
    for {
        select {
        case notification := <-notifications:
            // Render notification as HTML
            html := renderNotificationHTML(notification)
            
            // Send SSE event
            fmt.Fprintf(w, "event: notification\n")
            fmt.Fprintf(w, "data: %s\n\n", html)
            flusher.Flush()
            
        case <-r.Context().Done():
            return
        }
    }
}

func renderNotificationHTML(notification Notification) string {
    template := func(b *Builder) Node {
        return b.Div(
            minty.Class("notification notification-" + notification.Type),
            minty.Attr("hx-swap-oob", "afterbegin:#notifications-list"),
            
            b.Div(minty.Class("notification-content"),
                b.Strong(notification.Title),
                b.P(notification.Message),
                b.Time(
                    minty.Class("notification-time"),
                    notification.CreatedAt.Format("3:04 PM"),
                ),
            ),
            
            b.Button(
                minty.Class("notification-close"),
                minty.HtmxDelete("/notifications/" + notification.ID),
                minty.HtmxTarget("closest .notification"),
                minty.HtmxSwap("outerHTML"),
                "×",
            ),
        )
    }
    
    var buf strings.Builder
    minty.Render(template(minty.NewBuilder()), &buf)
    return buf.String()
}
```

### **Live Data Feeds**

SSE enables efficient live data feeds that update specific page sections without requiring client-side polling.

```go
// Live activity feed component
func liveActivityFeed(userID int) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("activity-feed"),
            b.Header(minty.Class("feed-header"),
                b.H3("Recent Activity"),
                b.Div(
                    minty.Class("connection-status"),
                    minty.Id("connection-status"),
                    "Connected",
                ),
            ),
            
            b.Div(
                minty.Id("activity-list"),
                minty.Class("activity-list"),
                minty.Attr("hx-ext", "sse"),
                minty.Attr("sse-connect", fmt.Sprintf("/activity-stream/%d", userID)),
                minty.Attr("sse-close", "connection-closed"),
                
                // Initial activity items
                initialActivityItems(userID)(b),
            ),
        )
    }
}

// SSE activity stream handler
func activityStreamHandler(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromPath(r.URL.Path)
    
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    flusher, _ := w.(http.Flusher)
    
    activities := make(chan Activity, 10)
    defer close(activities)
    
    activityManager.Subscribe(userID, activities)
    defer activityManager.Unsubscribe(userID, activities)
    
    // Send periodic heartbeats
    heartbeat := time.NewTicker(30 * time.Second)
    defer heartbeat.Stop()
    
    for {
        select {
        case activity := <-activities:
            html := renderActivityHTML(activity)
            fmt.Fprintf(w, "event: activity\n")
            fmt.Fprintf(w, "data: %s\n\n", html)
            flusher.Flush()
            
        case <-heartbeat.C:
            fmt.Fprintf(w, "event: heartbeat\n")
            fmt.Fprintf(w, "data: ping\n\n")
            flusher.Flush()
            
        case <-r.Context().Done():
            return
        }
    }
}

func renderActivityHTML(activity Activity) string {
    template := func(b *Builder) Node {
        return b.Div(
            minty.Class("activity-item new-activity"),
            minty.Attr("hx-swap-oob", "afterbegin:#activity-list"),
            
            b.Div(minty.Class("activity-avatar"),
                b.Img(
                    minty.Src(activity.User.Avatar),
                    minty.Alt(activity.User.Name),
                ),
            ),
            
            b.Div(minty.Class("activity-content"),
                b.P(
                    b.Strong(activity.User.Name),
                    " " + activity.Description,
                ),
                b.Time(
                    minty.Class("activity-time"),
                    activity.CreatedAt.Format("3:04 PM"),
                ),
            ),
        )
    }
    
    var buf strings.Builder
    minty.Render(template(minty.NewBuilder()), &buf)
    return buf.String()
}
```

### **Connection State Management**

Real-time applications require handling of connection states and graceful degradation when connections are lost.

```go
// Connection-aware SSE component
func connectionAwareComponent(streamURL string, fallbackURL string) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Class("sse-component"),
            minty.Attr("hx-ext", "sse"),
            minty.Attr("sse-connect", streamURL),
            
            // Handle connection events
            minty.Attr("hx-on", `sse:open: 
                document.getElementById('connection-status').textContent = 'Connected';
                document.getElementById('connection-status').className = 'connected';`),
            
            minty.Attr("hx-on", `sse:error: 
                document.getElementById('connection-status').textContent = 'Disconnected';
                document.getElementById('connection-status').className = 'disconnected';
                // Fallback to polling
                htmx.trigger('#fallback-trigger', 'startPolling');`),
            
            b.Div(
                minty.Id("connection-status"),
                minty.Class("connection-status connected"),
                "Connected",
            ),
            
            b.Div(minty.Id("content-area")),
            
            // Fallback polling trigger
            b.Div(
                minty.Id("fallback-trigger"),
                minty.Style("display: none;"),
                minty.Attr("hx-on", `startPolling: 
                    this.setAttribute('hx-get', '` + fallbackURL + `');
                    this.setAttribute('hx-trigger', 'every 5s');
                    this.setAttribute('hx-target', '#content-area');
                    htmx.process(this);`),
            ),
        )
    }
}
```

---

## **Common Interactive Patterns**

### **Modal Dialog Systems**

Modal dialogs represent one of the most common interactive patterns in web applications. Minty provides elegant patterns for server-rendered modal content.

```go
// Modal dialog helper
func modalDialog(id, title string, content H, actions H) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Id(id),
            minty.Class("modal-overlay"),
            minty.Attr("hx-on", "click: if (event.target === this) this.remove()"),
            
            b.Div(minty.Class("modal-dialog"),
                b.Header(minty.Class("modal-header"),
                    b.H2(title),
                    b.Button(
                        minty.Class("modal-close"),
                        minty.Attr("onclick", "this.closest('.modal-overlay').remove()"),
                        "×",
                    ),
                ),
                
                b.Div(minty.Class("modal-body"),
                    content(b),
                ),
                
                b.Footer(minty.Class("modal-footer"),
                    actions(b),
                ),
            ),
        )
    }
}

// Modal trigger button
func modalTrigger(text, modalURL string) H {
    return func(b *Builder) Node {
        return b.Button(
            minty.Class("btn btn-primary"),
            minty.HtmxGet(modalURL),
            minty.HtmxTarget("body"),
            minty.HtmxSwap("beforeend"),
            text,
        )
    }
}

// Server-side modal handlers
func confirmDeleteModalHandler(w http.ResponseWriter, r *http.Request) {
    itemID := r.URL.Query().Get("id")
    item := getItemByID(itemID)
    
    content := func(b *Builder) Node {
        return b.P(
            fmt.Sprintf("Are you sure you want to delete \"%s\"? This action cannot be undone.", 
                item.Name),
        )
    }
    
    actions := func(b *Builder) Node {
        return b.Div(minty.Class("modal-actions"),
            b.Button(
                minty.Class("btn btn-danger"),
                minty.HtmxDelete("/items/" + itemID),
                minty.HtmxTarget("#item-list"),
                minty.HtmxSwap("outerHTML"),
                minty.Attr("hx-on", "htmx:afterRequest: this.closest('.modal-overlay').remove()"),
                "Delete",
            ),
            b.Button(
                minty.Class("btn btn-secondary"),
                minty.Attr("onclick", "this.closest('.modal-overlay').remove()"),
                "Cancel",
            ),
        )
    }
    
    modal := modalDialog("confirm-delete-modal", "Confirm Delete", content, actions)
    renderFragment(w, modal)
}

// Usage in item list
func itemCard(item Item) H {
    return func(b *Builder) Node {
        return b.Div(
            minty.Class("item-card"),
            minty.Id("item-" + item.ID),
            
            b.H3(item.Name),
            b.P(item.Description),
            
            b.Div(minty.Class("item-actions"),
                b.A(
                    minty.Class("btn btn-primary"),
                    minty.Href("/items/" + item.ID + "/edit"),
                    "Edit",
                ),
                modalTrigger("Delete", "/modals/confirm-delete?id=" + item.ID)(b),
            ),
        )
    }
}
```

### **Live Search Interfaces**

Search interfaces benefit greatly from real-time feedback and result filtering without full page reloads.

```go
// Comprehensive live search component
func liveSearchInterface(searchEndpoint string, placeholder string) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("search-interface"),
            b.Div(minty.Class("search-input-container"),
                b.Input(
                    minty.Type("text"),
                    minty.Name("query"),
                    minty.Placeholder(placeholder),
                    minty.Class("search-input"),
                    
                    // Live search with debouncing
                    minty.HtmxGet(searchEndpoint),
                    minty.HtmxTarget("#search-results"),
                    minty.HtmxTrigger("keyup changed delay:300ms"),
                    minty.HtmxSwap("innerHTML"),
                    minty.HtmxIndicator("#search-spinner"),
                    
                    // Include filters in search
                    minty.Attr("hx-include", ".search-filter"),
                ),
                
                b.Div(
                    minty.Id("search-spinner"),
                    minty.Class("htmx-indicator search-spinner"),
                    "Searching...",
                ),
            ),
            
            // Search filters
            b.Div(minty.Class("search-filters"),
                b.Select(
                    minty.Name("category"),
                    minty.Class("search-filter"),
                    minty.HtmxGet(searchEndpoint),
                    minty.HtmxTarget("#search-results"),
                    minty.HtmxTrigger("change"),
                    minty.Attr("hx-include", ".search-input, .search-filter"),
                    
                    b.Option(minty.Value(""), "All Categories"),
                    b.Option(minty.Value("articles"), "Articles"),
                    b.Option(minty.Value("products"), "Products"),
                    b.Option(minty.Value("users"), "Users"),
                ),
                
                b.Label(minty.Class("filter-checkbox"),
                    b.Input(
                        minty.Type("checkbox"),
                        minty.Name("featured"),
                        minty.Value("true"),
                        minty.Class("search-filter"),
                        minty.HtmxGet(searchEndpoint),
                        minty.HtmxTarget("#search-results"),
                        minty.HtmxTrigger("change"),
                        minty.Attr("hx-include", ".search-input, .search-filter"),
                    ),
                    " Featured only",
                ),
            ),
            
            // Search results area
            b.Div(
                minty.Id("search-results"),
                minty.Class("search-results"),
                initialSearchResults()(b),
            ),
        )
    }
}

// Server-side search handler
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("query")
    category := r.FormValue("category")
    featuredOnly := r.FormValue("featured") == "true"
    
    results := performSearch(SearchParams{
        Query:       query,
        Category:    category,
        FeaturedOnly: featuredOnly,
    })
    
    fragment := searchResultsFragment(results, query)
    renderFragment(w, fragment)
}

func searchResultsFragment(results []SearchResult, query string) H {
    return func(b *Builder) Node {
        if len(results) == 0 {
            return b.Div(minty.Class("no-results"),
                b.P(fmt.Sprintf("No results found for \"%s\"", query)),
                b.P("Try adjusting your search terms or filters."),
            )
        }
        
        var resultNodes []Node
        for _, result := range results {
            resultNodes = append(resultNodes, searchResultCard(result)(b))
        }
        
        return b.Div(minty.Class("search-results-list"),
            b.P(
                minty.Class("result-count"),
                fmt.Sprintf("Found %d results for \"%s\"", len(results), query),
            ),
            resultNodes...,
        )
    }
}

func searchResultCard(result SearchResult) H {
    return func(b *Builder) Node {
        return b.Div(minty.Class("search-result-card"),
            b.H3(
                b.A(minty.Href(result.URL), result.Title),
            ),
            b.P(minty.Class("result-description"), result.Description),
            b.Div(minty.Class("result-meta"),
                b.Span(minty.Class("result-category"), result.Category),
                b.Span(minty.Class("result-date"), 
                    result.CreatedAt.Format("Jan 2, 2006")),
            ),
        )
    }
}
```

### **Infinite Scroll and Pagination**

Large datasets benefit from infinite scroll or dynamic pagination that loads content as needed.

```go
// Infinite scroll component
func infiniteScrollList(itemsEndpoint string, initialItems []Item) H {
    return func(b *Builder) Node {
        var itemNodes []Node
        for _, item := range initialItems {
            itemNodes = append(itemNodes, listItemCard(item)(b))
        }
        
        return b.Div(minty.Class("infinite-scroll-container"),
            b.Div(
                minty.Id("items-list"),
                minty.Class("items-list"),
                itemNodes...,
            ),
            
            // Load more trigger
            b.Div(
                minty.Id("load-more-trigger"),
                minty.Class("load-more-trigger"),
                minty.HtmxGet(itemsEndpoint + "?page=2"),
                minty.HtmxTarget("#items-list"),
                minty.HtmxSwap("beforeend"),
                minty.HtmxTrigger("intersect once"),
                minty.HtmxIndicator("#loading-indicator"),
                
                b.Div(
                    minty.Id("loading-indicator"),
                    minty.Class("htmx-indicator"),
                    "Loading more items...",
                ),
            ),
        )
    }
}

// Server-side pagination handler
func itemsListHandler(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    items, hasMore := getItemsPage(page, 20)
    
    ctx := NewHTMXContext(r)
    if ctx.IsHTMX {
        fragment := itemsPageFragment(items, page, hasMore)
        renderFragment(w, fragment)
        return
    }
    
    // Full page for direct access
    allItems, _ := getItemsPage(1, 20)
    page := itemsPage(allItems)
    renderPage(w, page)
}

func itemsPageFragment(items []Item, page int, hasMore bool) H {
    return func(b *Builder) Node {
        var itemNodes []Node
        for _, item := range items {
            itemNodes = append(itemNodes, listItemCard(item)(b))
        }
        
        nodes := itemNodes
        
        if hasMore {
            // Add next load trigger
            loadMoreTrigger := b.Div(
                minty.Id("load-more-trigger"),
                minty.Class("load-more-trigger"),
                minty.HtmxGet(fmt.Sprintf("/items?page=%d", page+1)),
                minty.HtmxTarget("#items-list"),
                minty.HtmxSwap("beforeend"),
                minty.HtmxTrigger("intersect once"),
                minty.Attr("hx-swap-oob", "outerHTML"),
                
                "Loading more items...",
            )
            nodes = append(nodes, loadMoreTrigger)
        } else {
            // Remove load trigger when no more items
            endMarker := b.Div(
                minty.Id("load-more-trigger"),
                minty.Attr("hx-swap-oob", "outerHTML"),
                minty.Class("end-marker"),
                "No more items to load",
            )
            nodes = append(nodes, endMarker)
        }
        
        return minty.Fragment(nodes...)
    }
}
```

---

## **HTMX Integration Across the Minty System**

HTMX integration extends throughout the entire Minty System, enabling dynamic behavior in business domains, iterator-driven components, and theme-based interfaces:

### **Business Domain HTMX Integration**

Business domain operations leverage HTMX for dynamic interactions:

```go
// Finance domain with HTMX integration
func PaymentButton(theme Theme, invoice mifi.Invoice) mi.H {
    return theme.Button(
        fmt.Sprintf("Pay %s", invoice.Amount.Format()), "payment",
        minty.HtmxPost("/api/invoices/"+invoice.ID+"/pay"),
        minty.HtmxTarget("#invoice-"+invoice.ID),
        minty.HtmxIndicator("#payment-spinner"),
        minty.HtmxConfirm(fmt.Sprintf("Pay invoice #%s for %s?", 
            invoice.Number, invoice.Amount.Format())),
    )
}

// Account balance updates with HTMX
func AccountSummaryCard(theme Theme, account mifi.Account) mi.H {
    return theme.Card(account.Name,
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Id("account-"+account.ID),
                b.P("Balance: ", b.Strong(account.Balance.Format())),
                b.Button(
                    minty.HtmxGet("/api/accounts/"+account.ID+"/refresh"),
                    minty.HtmxTarget("#account-"+account.ID),
                    minty.HtmxSwap("outerHTML"),
                    "Refresh Balance",
                ),
            )
        },
    )
}

// Logistics domain with live tracking updates
func ShipmentTracker(theme Theme, shipment mimo.Shipment) mi.H {
    return theme.Card("Shipment "+shipment.TrackingCode,
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Id("shipment-"+shipment.ID),
                b.P("Status: ", b.Span(mi.Id("status"), shipment.Status)),
                b.P("Location: ", b.Span(mi.Id("location"), shipment.CurrentLocation)),
                b.Div(
                    // Auto-refresh tracking status
                    minty.HtmxGet("/api/shipments/"+shipment.ID+"/status"),
                    minty.HtmxTarget("#shipment-"+shipment.ID),
                    minty.HtmxTrigger("every 30s"),
                    minty.HtmxSwap("outerHTML"),
                ),
            )
        },
    )
}
```

### **Iterator-Powered Dynamic Components**

HTMX works seamlessly with iterator-based data processing:

```go
// Dynamic filtered lists with iterator integration
func DynamicUserList(users []User) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Id("user-list-container"),
            // Filter controls with HTMX
            b.Form(
                minty.HtmxGet("/users/filtered"),
                minty.HtmxTarget("#user-list"),
                minty.HtmxTrigger("change, keyup delay:300ms"),
                
                b.Input(minty.Name("search"), minty.Placeholder("Search users...")),
                b.Select(minty.Name("department"),
                    b.Option(minty.Value(""), "All Departments"),
                    b.Option(minty.Value("engineering"), "Engineering"),
                    b.Option(minty.Value("sales"), "Sales"),
                ),
            ),
            
            // Dynamic list updated by HTMX
            b.Div(mi.Id("user-list"),
                // Iterator-generated user cards
                miex.Map(users, func(u User) mi.H {
                    return UserCard(u)
                })...,
            ),
        )
    }
}

// Server-side filter handler using iterators
func FilterUsersHandler(w http.ResponseWriter, r *http.Request) {
    search := r.URL.Query().Get("search")
    department := r.URL.Query().Get("department")
    
    // Use iterators for filtering
    filteredUsers := miex.ChainSlice(getAllUsers()).
        Filter(func(u User) bool {
            if search != "" && !strings.Contains(strings.ToLower(u.Name), 
                strings.ToLower(search)) {
                return false
            }
            if department != "" && u.Department != department {
                return false
            }
            return true
        }).
        ToSlice()
    
    // Return updated list fragment
    fragment := func(b *mi.Builder) mi.Node {
        return mi.NewFragment(
            miex.Map(filteredUsers, func(u User) mi.H {
                return UserCard(u)
            })...,
        )
    }
    
    html := mi.Render(fragment)
    w.Write([]byte(html))
}
```

### **Theme-Based Dynamic Interface**

Themes support HTMX attributes consistently across different styling frameworks:

```go
// Bootstrap theme with HTMX support
func (t *BootstrapTheme) Modal(title string, content mi.H, 
    trigger mi.H) mi.H {
    modalId := generateId()
    
    return func(b *mi.Builder) mi.Node {
        return b.Div(
            // Trigger element with HTMX
            trigger(b),
            
            // Modal content loaded via HTMX
            b.Div(mi.Id(modalId), mi.Class("modal"),
                minty.HtmxGet("/modals/"+title),
                minty.HtmxTarget("#"+modalId+" .modal-body"),
                minty.HtmxTrigger("click from:#"+modalId+"-trigger"),
                
                b.Div(mi.Class("modal-dialog"),
                    b.Div(mi.Class("modal-content"),
                        b.Div(mi.Class("modal-header"),
                            b.H5(mi.Class("modal-title"), title),
                        ),
                        b.Div(mi.Class("modal-body"),
                            content(b),
                        ),
                    ),
                ),
            ),
        )
    }
}

// Tailwind theme with same HTMX patterns, different styling
func (t *TailwindTheme) Modal(title string, content mi.H, 
    trigger mi.H) mi.H {
    // Same HTMX functionality, different CSS classes
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("relative"),
            trigger(b),
            b.Div(
                mi.Class("fixed inset-0 bg-gray-600 bg-opacity-50"),
                minty.HtmxGet("/modals/"+title),
                minty.HtmxTarget(".modal-content"),
                
                b.Div(mi.Class("flex items-center justify-center min-h-screen"),
                    b.Div(mi.Class("modal-content bg-white rounded-lg shadow-lg"),
                        content(b),
                    ),
                ),
            ),
        )
    }
}
```

### **Cross-Domain Interactive Dashboards**

Complex applications combine all system elements with HTMX:

```go
// Multi-domain dashboard with live updates
func InteractiveDashboard(services *ApplicationServices, 
    theme Theme) mi.H {
    return theme.Layout("Live Business Dashboard",
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Class("dashboard-grid"),
                // Auto-refreshing finance metrics
                b.Div(mi.Id("finance-metrics"),
                    minty.HtmxGet("/api/finance/metrics"),
                    minty.HtmxTarget("#finance-metrics"),
                    minty.HtmxTrigger("every 30s"),
                    mintyfinui.MetricsCard(theme, services.Finance),
                ),
                
                // Real-time shipment updates
                b.Div(mi.Id("logistics-status"),
                    minty.HtmxGet("/api/logistics/status"),
                    minty.HtmxTarget("#logistics-status"),
                    minty.HtmxTrigger("every 15s"),
                    mintymoveui.StatusCard(theme, services.Logistics),
                ),
                
                // Live order activity
                b.Div(mi.Id("order-activity"),
                    minty.HtmxGet("/api/ecommerce/activity"),
                    minty.HtmxTarget("#order-activity"),
                    minty.HtmxTrigger("every 10s"),
                    mintycartui.ActivityCard(theme, services.Ecommerce),
                ),
            )
        },
    )
}
```

These HTMX integration patterns enable the Minty System to deliver sophisticated interactive experiences while maintaining server-side simplicity and avoiding JavaScript complexity.

---

## **Conclusion**

Minty's HTMX integration represents a fundamental shift in web application architecture, moving from JavaScript-heavy client-side applications to server-driven hypermedia applications that achieve equivalent interactivity with dramatically reduced complexity.

The patterns explored in this document demonstrate that modern web application requirements—real-time updates, dynamic forms, interactive interfaces, live search, and responsive feedback—can be achieved entirely through server-side rendering combined with declarative HTML attributes. This approach eliminates the complexity of client-side state management, reduces security vulnerabilities, improves performance, and simplifies debugging.

These HTMX integration patterns scale throughout the **entire Minty System**, enabling dynamic behavior in business domain operations, iterator-driven data processing, theme-based interfaces, and complex multi-domain applications. The same server-first philosophy that makes simple HTML generation clean and maintainable also makes sophisticated interactive applications productive and scalable.

The key to Minty's success lies in making JavaScript-free development not just possible, but preferable. Through comprehensive HTMX integration, thoughtful fragment design, and server-side helper patterns, Minty enables developers to build sophisticated interactive applications while maintaining Go's core values of simplicity, type safety, and performance.

As web development continues to evolve, the combination of server-side rendering with selective interactivity represents a sustainable path forward. Minty's HTMX integration provides the foundation for this future, enabling developers to build modern web applications without sacrificing simplicity or maintainability.

The patterns and techniques presented here form a comprehensive toolkit for JavaScript-free web development. By embracing these approaches, developers can create rich, interactive user experiences while avoiding the complexity and overhead of traditional JavaScript frameworks.