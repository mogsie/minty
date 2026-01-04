# Architecture Guide - Clean Architecture with Minty

This document explains the architectural principles and design decisions behind the Minty System, demonstrating how to build maintainable, testable, and scalable HTML applications in Go.

## 🏗️ Clean Architecture Principles

### The Dependency Rule

**Dependencies point inward.** Source code dependencies must point only inward, toward higher-level policies.

```
┌─────────────────────────────────────────────────┐
│                🌐 Presentation                  │ ──┐
│            UI Components & Themes               │   │
├─────────────────────────────────────────────────┤   │
│               🏢 Application                    │ ──┤ Dependencies
│            Service Orchestration               │   │ point inward
├─────────────────────────────────────────────────┤   │
│                💼 Domain                       │ ──┘
│            Business Logic & Rules              │
├─────────────────────────────────────────────────┤
│              🔧 Infrastructure                 │
│          Framework & External Concerns         │
└─────────────────────────────────────────────────┘
```

### Why This Matters

1. **Business Logic Independence**: Core business rules don't depend on databases, web frameworks, or UI libraries
2. **Testability**: Pure business logic can be tested without external dependencies
3. **Flexibility**: Can switch UI frameworks, databases, or delivery mechanisms without changing business rules
4. **Maintainability**: Clear boundaries make the system easier to understand and modify

## 📦 Layer Breakdown

### 🔧 Infrastructure Layer (Framework)

**Purpose**: Provides the foundation and utilities that all other layers can use.

**Packages**:
- `minty` - Pure HTML generation functions
- `mintyex` - Shared business types and utilities  
- `mintyui` - Theme system and higher-level UI components

**Characteristics**:
- ✅ Framework-agnostic utilities
- ✅ No business logic
- ✅ Can be used by all other layers
- ✅ Highly reusable

**Example**:
```go
// mintyex provides pure business utilities
money1 := mintyex.NewMoney(100.00, "USD")
money2 := mintyex.NewMoney(50.00, "USD")
total, _ := money1.Add(money2)

// minty provides pure HTML generation
html := mi.Render(func(b *mi.Builder) mi.Node {
    return b.Div(mi.Class("container"), "Hello World")
})
```

### 💼 Domain Layer (Business Logic)

**Purpose**: Contains enterprise business rules and domain models.

**Packages**:
- `mintyfin` - Financial domain logic
- `mintymove` - Logistics domain logic
- `mintycart` - E-commerce domain logic

**Characteristics**:
- ✅ **Zero UI dependencies** - No imports of UI packages
- ✅ Pure business logic and rules
- ✅ Framework agnostic
- ✅ Highly testable
- ✅ Self-contained domain services

**Example**:
```go
// Pure business logic - no UI dependencies
func (fs *FinanceService) CreateAccount(name, accountType string, 
    initialBalance mintyex.Money, customerID string) (*Account, error) {
    
    account := Account{
        ID:        generateID("acc"),
        Name:      name,
        Balance:   initialBalance,
        Type:      accountType,
        Status:    mintyex.StatusActive,
        CreatedAt: time.Now(),
    }
    
    // Pure business validation
    if errors := ValidateAccount(account); errors.HasErrors() {
        return nil, errors
    }
    
    fs.accounts = append(fs.accounts, account)
    return &account, nil
}
```

### 🏢 Application Layer (Use Cases)

**Purpose**: Orchestrates domain services and coordinates business workflows.

**Examples**:
- `ApplicationServices` - Coordinates multiple domain services
- HTTP handlers that orchestrate business operations
- Workflow engines that span multiple domains

**Characteristics**:
- ✅ Orchestrates domain services
- ✅ Contains application-specific business rules
- ✅ Can depend on domain layer
- ✅ Independent of presentation layer

**Example**:
```go
type ApplicationServices struct {
    Finance   *mintyfin.FinanceService
    Logistics *mintymove.LogisticsService
    Ecommerce *mintycart.EcommerceService
}

func (app *ApplicationServices) ProcessOrder(orderData OrderRequest) error {
    // 1. Create order in e-commerce domain
    order, err := app.Ecommerce.CreateOrder(...)
    
    // 2. Generate invoice in finance domain
    invoice, err := app.Finance.CreateInvoice(...)
    
    // 3. Create shipment in logistics domain
    shipment, err := app.Logistics.CreateShipment(...)
    
    return nil
}
```

### 🌐 Presentation Layer (UI)

**Purpose**: Converts domain data into UI components and handles user interface concerns.

**Packages**:
- `mintyfinui` - Finance UI components
- `mintymoveui` - Logistics UI components
- `mintycartui` - E-commerce UI components
- `bootstrap`, `tailwind` - Theme implementations

**Characteristics**:
- ✅ **Presentation Adapters** - Convert domain data to UI
- ✅ Theme-aware components
- ✅ No business logic
- ✅ Can depend on all other layers

**Example**:
```go
// Presentation adapter - converts domain data to UI
func AccountSummaryCard(theme mui.Theme, account mintyfin.Account) mi.H {
    // Prepare display data (no business logic here)
    displayData := mintyfin.PrepareAccountForDisplay(account)
    
    // Create UI using theme
    return mui.DomainCard(theme, "finance", account.Name, 
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Class("account_summary"),
                b.P("Balance: ", displayData.FormattedBalance),
                b.P("Status: ", displayData.StatusDisplay),
                theme.Button("View Details", "secondary")(b),
            )
        })
}
```

## 🔄 Data Flow Architecture

### Inward Data Flow (User Input)

```
User Input → HTTP Handler → Application Service → Domain Service → Business Logic
```

Example:
```go
// 1. HTTP Handler (Presentation)
func (app *WebApp) createAccountHandler(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    accountType := r.FormValue("type")
    
    // 2. Application Service (Use Case)
    account, err := app.services.Finance.CreateAccount(name, accountType, ...)
    
    // 3. Domain Service does business logic
    // 4. Return result up the chain
}
```

### Outward Data Flow (Display)

```
Domain Data → Data Preparation → Presentation Adapter → Theme → HTML
```

Example:
```go
// 1. Get domain data
account := financeService.GetAccount(accountID)

// 2. Prepare for display (still in domain)
displayData := mintyfin.PrepareAccountForDisplay(account)

// 3. Presentation adapter
card := mintyfinui.AccountSummaryCard(theme, account)

// 4. Render to HTML
html := mi.Render(card)
```

## 🎯 Design Patterns

### Repository Pattern (Implicit)

Each domain service acts as a repository for its entities:

```go
type FinanceService struct {
    accounts     []Account      // In-memory store
    transactions []Transaction  // In real app: database
    invoices     []Invoice
}

func (fs *FinanceService) GetAccount(id string) (*Account, error) {
    // Repository behavior - could be database, API, etc.
    for i, account := range fs.accounts {
        if account.ID == id {
            return &fs.accounts[i], nil
        }
    }
    return nil, errors.New("account not found")
}
```

### Adapter Pattern

Presentation adapters convert domain data to UI components:

```go
// Domain data
type Account struct {
    Balance mintyex.Money
    Status  string
}

// Presentation adapter
func PrepareAccountForDisplay(account Account) AccountDisplayData {
    return AccountDisplayData{
        FormattedBalance: account.Balance.Format(),
        StatusClass:      getStatusClass(account.Status),
        StatusDisplay:    getStatusDisplay(account.Status),
    }
}
```

### Strategy Pattern

Theme system allows different UI strategies:

```go
type Theme interface {
    Button(text, variant string, attrs ...mi.Attribute) mi.H
    Card(title string, content mi.H) mi.H
}

// Different strategies
bootstrapTheme := bootstrap.NewBootstrapTheme()
tailwindTheme := tailwind.NewTailwindTheme()

// Same interface, different implementation
button1 := bootstrapTheme.Button("Click", "primary")
button2 := tailwindTheme.Button("Click", "primary")
```

### Dependency Injection

Services are injected into higher layers:

```go
type WebApplication struct {
    services *ApplicationServices  // Injected dependency
    theme    mui.Theme            // Injected dependency
}

func NewWebApplication(services *ApplicationServices, theme mui.Theme) *WebApplication {
    return &WebApplication{
        services: services,  // DI
        theme:    theme,     // DI
    }
}
```

## ✅ Benefits Achieved

### 1. **Testability**

**Domain Logic** (easy to test):
```go
func TestAccountCreation(t *testing.T) {
    service := mintyfin.NewFinanceService()
    
    account, err := service.CreateAccount("Test", "checking", money, "customer")
    
    assert.NoError(t, err)
    assert.Equal(t, "Test", account.Name)
    // No mocks needed! Pure business logic
}
```

**UI Components** (can be tested separately):
```go
func TestAccountCard(t *testing.T) {
    mockTheme := &MockTheme{}
    account := mintyfin.Account{Name: "Test", Balance: money}
    
    card := mintyfinui.AccountSummaryCard(mockTheme, account)
    html := mi.Render(card)
    
    assert.Contains(t, html, "Test")
}
```

### 2. **Framework Independence**

Same business logic works with any UI framework:

```go
// Business logic (unchanged)
dashboardData := mintyfin.PrepareDashboardData(service)

// Different themes
bootstrapUI := mintyfinui.FinancialDashboard(bootstrapTheme, dashboardData, ...)
tailwindUI := mintyfinui.FinancialDashboard(tailwindTheme, dashboardData, ...)
reactUI := mintyfinui.FinancialDashboard(reactTheme, dashboardData, ...)
```

### 3. **Maintainability**

Clear boundaries make changes safe:

- **Change Business Rules**: Modify domain packages, presentation adapters handle the rest
- **Change UI**: Modify themes/presentation, business logic unaffected
- **Add Features**: Add to domain first, then create presentation adapters

### 4. **Reusability**

Components can be reused across contexts:

```go
// Same business logic
account := service.GetAccount(id)

// Different presentations
summaryCard := mintyfinui.AccountSummaryCard(theme, account)
tableRow := mintyfinui.AccountTableRow(theme, account)  
widget := mintyfinui.AccountWidget(theme, account)
```

## 🚀 Extension Points

### Adding New Domains

1. **Create Domain Package**:
```go
// pkg/inventory/inventory.go
type Product struct {
    ID       string
    Name     string
    Quantity int
}

type InventoryService struct {
    products []Product
}

func (is *InventoryService) UpdateStock(productID string, quantity int) error {
    // Pure business logic
}
```

2. **Create Presentation Adapter**:
```go
// pkg/inventoryui/inventoryui.go  
func ProductCard(theme mui.Theme, product inventory.Product) mi.H {
    return theme.Card(product.Name, func(b *mi.Builder) mi.Node {
        return b.P(fmt.Sprintf("Stock: %d", product.Quantity))
    })
}
```

3. **Integrate in Application**:
```go
type ApplicationServices struct {
    Finance   *mintyfin.FinanceService
    Inventory *inventory.InventoryService  // Add new domain
}
```

### Creating Custom Themes

```go
type MaterialUITheme struct{}

func (t *MaterialUITheme) Button(text, variant string, attrs ...mi.Attribute) mi.H {
    return func(b *mi.Builder) mi.Node {
        class := "MuiButton-root MuiButton-" + variant
        return b.Button(mi.Class(class), text)
    }
}

// Implement all Theme interface methods...
```

### Adding New Presentation Patterns

```go
// Mobile-specific presentation adapter
func MobileAccountCard(theme mui.Theme, account mintyfin.Account) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("mobile-account-card"),
            // Mobile-optimized layout
        )
    }
}
```

## 🎓 Architecture Decision Records

### Why Presentation Adapters?

**Problem**: How to convert domain data to UI without coupling business logic to presentation concerns?

**Decision**: Create presentation adapter packages that depend on domain packages but contain no business logic.

**Benefits**:
- Business logic stays pure
- UI can change without affecting business rules
- Multiple UI representations of same data
- Clear separation of concerns

### Why Theme System?

**Problem**: How to support multiple CSS frameworks without duplicating component logic?

**Decision**: Abstract UI components behind a Theme interface.

**Benefits**:
- Same components work with Bootstrap, Tailwind, etc.
- Easy to switch themes
- Consistent component API
- Custom themes possible

### Why Zero UI Dependencies in Domains?

**Problem**: How to ensure business logic is truly independent and testable?

**Decision**: Domain packages cannot import any UI-related packages.

**Benefits**:
- Pure business logic
- Easy unit testing
- Framework independence
- Clear architectural boundaries

This architecture enables building complex, maintainable HTML applications while keeping business logic clean and testable.
