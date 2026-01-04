# Minty System Documentation - Part 7
## Business Domain Implementation: Domain-Driven Design with Go

---

### Table of Contents
1. [Domain-Driven Design Principles in the Minty System](#domain-driven-design-principles-in-the-minty-system)
2. [Pure Business Logic Architecture](#pure-business-logic-architecture)
3. [Finance Domain (mintyfin)](#finance-domain-mintyfin)
4. [Logistics Domain (mintymove)](#logistics-domain-mintymove)
5. [E-commerce Domain (mintycart)](#e-commerce-domain-mintycart)
6. [Cross-Domain Integration Patterns](#cross-domain-integration-patterns)
7. [Testing Domain Logic](#testing-domain-logic)
8. [Extending with Custom Domains](#extending-with-custom-domains)

---

## Domain-Driven Design Principles in the Minty System

The Minty System implements true Domain-Driven Design (DDD) principles with a focus on business logic purity and separation of concerns. Each domain package contains complete business functionality with zero dependencies on UI, database, or external frameworks.

### Core Domain Architecture

```go
// Domain packages have ZERO external dependencies
package mintyfin

import (
    "time"
    "errors"
    "github.com/ha1tch/mintyex"  // Only shared utilities
)

// No imports of:
// - UI libraries
// - Database drivers  
// - HTTP frameworks
// - External APIs
```

This strict dependency rule ensures that business logic remains:
- **Testable** without external systems
- **Portable** across different infrastructures  
- **Focused** on business rules and calculations
- **Maintainable** without UI coupling

### Business Entity Design

All domain entities follow consistent patterns with proper encapsulation and validation:

```go
type Account struct {
    ID          string           `json:"id"`
    Name        string           `json:"name"`
    Balance     mintyex.Money    `json:"balance"`
    Type        string           `json:"type"`
    Status      string           `json:"status"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

**Key Design Decisions:**

- **Money Type**: Uses `mintyex.Money` for precise financial calculations with no floating-point errors
- **Status Fields**: Consistent status management across all entities
- **Metadata Maps**: Extensible properties without schema changes
- **Time Tracking**: Created/updated timestamps for audit trails
- **JSON Tags**: Consistent serialization for API compatibility

---

## Pure Business Logic Architecture

### Domain Services Pattern

Each domain implements a service class that encapsulates all business operations:

```go
type FinanceService struct {
    accounts     []Account
    transactions []Transaction
    invoices     []Invoice
    customers    []Customer
}

func NewFinanceService() *FinanceService {
    return &FinanceService{
        accounts:     make([]Account, 0),
        transactions: make([]Transaction, 0),
        invoices:     make([]Invoice, 0),
        customers:    make([]Customer, 0),
    }
}
```

**Service Responsibilities:**
- Entity lifecycle management (create, read, update, delete)
- Business rule validation
- Cross-entity operations and workflows
- Data consistency enforcement
- Business logic encapsulation

### Validation Architecture

Comprehensive validation using shared validation utilities:

```go
func ValidateAccount(account Account) mintyex.ValidationErrors {
    var errors mintyex.ValidationErrors
    
    mintyex.ValidateRequired("name", account.Name, "Account Name", &errors)
    mintyex.ValidateMoneyAmount("balance", account.Balance, "Balance", &errors)
    
    if account.Type != "" {
        validTypes := []string{"checking", "savings", "investment", "credit"}
        if !mintyex.Contains(validTypes, account.Type) {
            errors.Add("type", "Account type must be one of: checking, savings, investment, credit")
        }
    }
    
    return errors
}
```

**Validation Features:**
- Type-safe error collection
- Field-specific error messages
- Business rule validation
- Cross-field validation support
- Internationalization ready

### Money Handling

All monetary calculations use the `mintyex.Money` type for precision:

```go
// Precise monetary calculations
balance := mintyex.NewMoney(1000.00, "USD")
debit := mintyex.NewMoney(150.50, "USD")

newBalance, err := balance.Subtract(debit)
if err != nil {
    return err // Currency mismatch protection
}

// Format for display
formatted := newBalance.Format() // "$849.50"
```

**Money Type Benefits:**
- No floating-point calculation errors
- Currency validation and tracking
- Automatic formatting for different currencies
- Arithmetic operations with error handling
- International currency support

---

## Finance Domain (mintyfin)

The finance domain provides comprehensive financial management capabilities with proper business rule enforcement and calculation accuracy.

### Core Entities

#### Account Management
```go
type Account struct {
    ID          string           `json:"id"`
    Name        string           `json:"name"`
    Balance     mintyex.Money    `json:"balance"`
    Type        string           `json:"type"` // checking, savings, investment, credit
    Status      string           `json:"status"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
    Description string           `json:"description"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

**Account Types:**
- **Checking**: Day-to-day transaction accounts
- **Savings**: Interest-bearing savings accounts  
- **Investment**: Investment and portfolio accounts
- **Credit**: Credit card and loan accounts

#### Transaction Processing
```go
type Transaction struct {
    ID          string        `json:"id"`
    AccountID   string        `json:"account_id"`
    Amount      mintyex.Money `json:"amount"`
    Description string        `json:"description"`
    Date        time.Time     `json:"date"`
    Status      string        `json:"status"`
    Type        string        `json:"type"` // debit, credit
    Category    string        `json:"category"`
    Reference   string        `json:"reference"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

**Transaction Features:**
- Automatic categorization based on description
- Double-entry accounting principles
- Status tracking (pending, completed, failed)
- Reference linking for related transactions
- Audit trail with timestamps

#### Invoice Management
```go
type Invoice struct {
    ID          string           `json:"id"`
    Number      string           `json:"number"`
    Amount      mintyex.Money    `json:"amount"`
    DueDate     time.Time        `json:"due_date"`
    Status      string           `json:"status"`
    Customer    Customer         `json:"customer"`
    Items       []InvoiceItem    `json:"items"`
    CreatedAt   time.Time        `json:"created_at"`
    PaidAt      *time.Time       `json:"paid_at,omitempty"`
    Description string           `json:"description"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

### Business Operations

#### Account Operations
```go
func (fs *FinanceService) CreateAccount(name, accountType string, 
    initialBalance mintyex.Money, customerID string) (*Account, error) {
    
    account := Account{
        ID:        generateID("acc"),
        Name:      name,
        Balance:   initialBalance,
        Status:    mintyex.StatusActive,
        Type:      accountType,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Metadata:  map[string]string{"customer_id": customerID},
    }
    
    // Validate business rules
    if errors := ValidateAccount(account); errors.HasErrors() {
        return nil, errors
    }
    
    fs.accounts = append(fs.accounts, account)
    return &account, nil
}
```

#### Transaction Processing
```go
func (fs *FinanceService) CreateTransaction(accountID string, amount mintyex.Money, 
    description, transactionType string) (*Transaction, error) {
    
    transaction := Transaction{
        ID:          generateID("txn"),
        AccountID:   accountID,
        Amount:      amount,
        Description: description,
        Type:        transactionType,
        Date:        time.Now(),
        Status:      mintyex.StatusPending,
    }
    
    // Auto-categorize transaction
    CategorizeTransaction(&transaction)
    
    // Validate business rules
    if errors := ValidateTransaction(transaction); errors.HasErrors() {
        return nil, errors
    }
    
    // Apply to account balance
    if err := fs.applyTransactionToAccount(&transaction); err != nil {
        return nil, err
    }
    
    transaction.Status = mintyex.StatusCompleted
    fs.transactions = append(fs.transactions, transaction)
    
    return &transaction, nil
}
```

#### Payment Processing
```go
func (fs *FinanceService) PayInvoice(invoiceID string, 
    paymentAmount mintyex.Money) error {
    
    invoice := fs.findInvoice(invoiceID)
    if invoice == nil {
        return errors.New("invoice not found")
    }
    
    if invoice.Status == mintyex.StatusPaid {
        return errors.New("invoice is already paid")
    }
    
    if paymentAmount.Amount != invoice.Amount.Amount {
        return errors.New("payment amount must match invoice amount")
    }
    
    invoice.Status = mintyex.StatusPaid
    now := time.Now()
    invoice.PaidAt = &now
    
    return nil
}
```

### Data Preparation for UI

The finance domain provides functions to prepare business data for display without coupling to UI frameworks:

```go
type AccountDisplayData struct {
    FormattedBalance string
    StatusDisplay    string
    StatusClass      string
    TypeDisplay      string
    TypeIcon         string
}

func PrepareAccountForDisplay(account Account) AccountDisplayData {
    return AccountDisplayData{
        FormattedBalance: account.Balance.Format(),
        StatusDisplay:    formatStatus(account.Status),
        StatusClass:      getStatusClass(account.Status),
        TypeDisplay:      formatAccountType(account.Type),
        TypeIcon:         getAccountTypeIcon(account.Type),
    }
}
```

This pattern maintains clean separation - business logic prepares data, but UI components handle the actual rendering.

---

## Logistics Domain (mintymove)

The logistics domain handles shipment tracking, route planning, and vehicle management with comprehensive business rule enforcement.

### Core Entities

#### Shipment Management
```go
type Shipment struct {
    ID              string           `json:"id"`
    TrackingCode    string           `json:"tracking_code"`
    Origin          mintyex.Address  `json:"origin"`
    Destination     mintyex.Address  `json:"destination"`
    Status          string           `json:"status"`
    Carrier         string           `json:"carrier"`
    ServiceType     string           `json:"service_type"`
    Cost            mintyex.Money    `json:"cost"`
    Weight          Weight           `json:"weight"`
    Dimensions      Dimensions       `json:"dimensions"`
    Items           []ShipmentItem   `json:"items"`
    EstimatedDelivery time.Time      `json:"estimated_delivery"`
    ActualDelivery  *time.Time       `json:"actual_delivery,omitempty"`
    CreatedAt       time.Time        `json:"created_at"`
    UpdatedAt       time.Time        `json:"updated_at"`
    Metadata        map[string]string `json:"metadata,omitempty"`
}
```

#### Vehicle and Driver Management
```go
type Vehicle struct {
    ID              string          `json:"id"`
    Name            string          `json:"name"`
    Type            string          `json:"type"` // truck, van, car
    LicensePlate    string          `json:"license_plate"`
    Capacity        VehicleCapacity `json:"capacity"`
    Status          string          `json:"status"`
    CurrentLocation mintyex.Address `json:"current_location"`
    Driver          Driver          `json:"driver"`
    MaintenanceLog  []MaintenanceRecord `json:"maintenance_log"`
    CreatedAt       time.Time       `json:"created_at"`
    Metadata        map[string]string `json:"metadata,omitempty"`
}

type Driver struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    LicenseNumber   string    `json:"license_number"`
    Phone           string    `json:"phone"`
    Status          string    `json:"status"`
    CurrentLocation mintyex.Address `json:"current_location"`
    AssignedRoutes  []string  `json:"assigned_routes"`
    CreatedAt       time.Time `json:"created_at"`
}
```

### Business Operations

#### Shipment Creation
```go
func (ls *LogisticsService) CreateShipment(trackingCode string, 
    origin, destination mintyex.Address, carrier, serviceType string,
    weight Weight, items []ShipmentItem) (*Shipment, error) {
    
    shipment := Shipment{
        ID:           generateID("shp"),
        TrackingCode: trackingCode,
        Origin:       origin,
        Destination:  destination,
        Carrier:      carrier,
        ServiceType:  serviceType,
        Status:       mintyex.StatusPending,
        Weight:       weight,
        Items:        items,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    // Calculate shipping cost
    cost, err := ls.calculateShippingCost(shipment)
    if err != nil {
        return nil, err
    }
    shipment.Cost = cost
    
    // Estimate delivery time
    estimatedDelivery, err := ls.calculateDeliveryTime(shipment)
    if err != nil {
        return nil, err
    }
    shipment.EstimatedDelivery = estimatedDelivery
    
    // Validate business rules
    if errors := ValidateShipment(shipment); errors.HasErrors() {
        return nil, errors
    }
    
    ls.shipments = append(ls.shipments, shipment)
    return &shipment, nil
}
```

#### Route Optimization
```go
func (ls *LogisticsService) OptimizeRoute(vehicleID string, 
    shipmentIDs []string) (*Route, error) {
    
    vehicle := ls.findVehicle(vehicleID)
    if vehicle == nil {
        return nil, errors.New("vehicle not found")
    }
    
    shipments := ls.getShipmentsByIDs(shipmentIDs)
    
    // Check capacity constraints
    totalWeight := calculateTotalWeight(shipments)
    if totalWeight.Exceeds(vehicle.Capacity.Weight) {
        return nil, errors.New("shipments exceed vehicle weight capacity")
    }
    
    // Optimize delivery sequence
    optimizedStops := ls.optimizeDeliverySequence(shipments)
    
    route := Route{
        ID:        generateID("rte"),
        VehicleID: vehicleID,
        Stops:     optimizedStops,
        Status:    mintyex.StatusPlanned,
        CreatedAt: time.Now(),
    }
    
    return &route, nil
}
```

---

## E-commerce Domain (mintycart)

The e-commerce domain provides complete online commerce functionality with inventory management, order processing, and customer management.

### Core Entities

#### Product Catalog
```go
type Product struct {
    ID          string        `json:"id"`
    Name        string        `json:"name"`
    Description string        `json:"description"`
    SKU         string        `json:"sku"`
    Category    string        `json:"category"`
    Price       mintyex.Money `json:"price"`
    Weight      Weight        `json:"weight"`
    Dimensions  Dimensions    `json:"dimensions"`
    Inventory   Inventory     `json:"inventory"`
    Images      []ProductImage `json:"images"`
    Attributes  map[string]string `json:"attributes"`
    Status      string        `json:"status"`
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type Inventory struct {
    Available   int    `json:"available"`
    Reserved    int    `json:"reserved"`
    Backorder   int    `json:"backorder"`
    Reorder     int    `json:"reorder_level"`
    MaxStock    int    `json:"max_stock"`
    Location    string `json:"location"`
}
```

#### Shopping Cart and Orders
```go
type Cart struct {
    ID         string     `json:"id"`
    CustomerID string     `json:"customer_id"`
    Items      []CartItem `json:"items"`
    Subtotal   mintyex.Money `json:"subtotal"`
    Tax        mintyex.Money `json:"tax"`
    Shipping   mintyex.Money `json:"shipping"`
    Total      mintyex.Money `json:"total"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}

type Order struct {
    ID              string           `json:"id"`
    Number          string           `json:"number"`
    Customer        Customer         `json:"customer"`
    Items           []OrderItem      `json:"items"`
    Subtotal        mintyex.Money    `json:"subtotal"`
    Tax             mintyex.Money    `json:"tax"`
    Shipping        mintyex.Money    `json:"shipping"`
    Total           mintyex.Money    `json:"total"`
    Status          string           `json:"status"`
    PaymentStatus   string           `json:"payment_status"`
    ShippingAddress mintyex.Address  `json:"shipping_address"`
    BillingAddress  mintyex.Address  `json:"billing_address"`
    CreatedAt       time.Time        `json:"created_at"`
    ShippedAt       *time.Time       `json:"shipped_at,omitempty"`
    DeliveredAt     *time.Time       `json:"delivered_at,omitempty"`
    Metadata        map[string]string `json:"metadata,omitempty"`
}
```

### Business Operations

#### Inventory Management
```go
func (es *EcommerceService) UpdateInventory(productID string, 
    quantityChange int, operation InventoryOperation) error {
    
    product := es.findProduct(productID)
    if product == nil {
        return errors.New("product not found")
    }
    
    switch operation {
    case InventoryReceived:
        product.Inventory.Available += quantityChange
    case InventoryReserved:
        if product.Inventory.Available < quantityChange {
            return errors.New("insufficient inventory available")
        }
        product.Inventory.Available -= quantityChange
        product.Inventory.Reserved += quantityChange
    case InventorySold:
        product.Inventory.Reserved -= quantityChange
    }
    
    // Check reorder levels
    if product.Inventory.Available <= product.Inventory.Reorder {
        // Trigger reorder notification
        es.triggerReorderNotification(productID)
    }
    
    product.UpdatedAt = time.Now()
    return nil
}
```

#### Order Processing
```go
func (es *EcommerceService) CreateOrder(cartID string, customer Customer,
    billingAddr, shippingAddr mintyex.Address, paymentMethod string) (*Order, error) {
    
    cart := es.findCart(cartID)
    if cart == nil {
        return nil, errors.New("cart not found")
    }
    
    // Verify inventory availability
    for _, item := range cart.Items {
        if !es.checkInventoryAvailable(item.ProductID, item.Quantity) {
            return nil, fmt.Errorf("insufficient inventory for product %s", item.ProductID)
        }
    }
    
    // Create order
    order := Order{
        ID:              generateID("ord"),
        Number:          es.generateOrderNumber(),
        Customer:        customer,
        Items:           convertCartItemsToOrderItems(cart.Items),
        Subtotal:        cart.Subtotal,
        Tax:            cart.Tax,
        Shipping:       cart.Shipping,
        Total:          cart.Total,
        Status:         mintyex.StatusPending,
        PaymentStatus:  "pending",
        ShippingAddress: shippingAddr,
        BillingAddress: billingAddr,
        CreatedAt:      time.Now(),
    }
    
    // Reserve inventory
    for _, item := range order.Items {
        if err := es.UpdateInventory(item.ProductID, item.Quantity, InventoryReserved); err != nil {
            return nil, err
        }
    }
    
    es.orders = append(es.orders, order)
    return &order, nil
}
```

---

## Cross-Domain Integration Patterns

The Minty System enables sophisticated cross-domain workflows while maintaining clean boundaries between domains.

### Order Fulfillment Workflow

```go
func ProcessCompleteOrder(services *ApplicationServices, orderID string) error {
    // 1. Get order from e-commerce domain
    order := services.Ecommerce.GetOrder(orderID)
    if order == nil {
        return errors.New("order not found")
    }
    
    // 2. Create invoice in finance domain
    invoiceItems := convertOrderItemsToInvoiceItems(order.Items)
    invoice, err := services.Finance.CreateInvoice(
        generateInvoiceNumber(), 
        convertCustomer(order.Customer),
        invoiceItems,
        order.CreatedAt.AddDate(0, 0, 30), // 30 day terms
    )
    if err != nil {
        return fmt.Errorf("failed to create invoice: %w", err)
    }
    
    // 3. Process payment in finance domain  
    err = services.Finance.PayInvoice(invoice.ID, order.Total)
    if err != nil {
        return fmt.Errorf("payment failed: %w", err)
    }
    
    // 4. Create shipment in logistics domain
    shipmentItems := convertOrderItemsToShipmentItems(order.Items)
    shipment, err := services.Logistics.CreateShipment(
        generateTrackingCode(),
        getWarehouseAddress(),
        order.ShippingAddress,
        "FedEx",
        "ground",
        calculateTotalWeight(shipmentItems),
        shipmentItems,
    )
    if err != nil {
        return fmt.Errorf("failed to create shipment: %w", err)
    }
    
    // 5. Update order with shipment info
    err = services.Ecommerce.ShipOrder(order.ID, shipment.TrackingCode)
    if err != nil {
        return fmt.Errorf("failed to update order: %w", err)
    }
    
    return nil
}
```

### Conversion Functions

Cross-domain integration requires conversion functions that translate entities between domains while preserving business integrity:

```go
// Convert e-commerce customer to finance customer
func convertCustomer(ecomCustomer mica.Customer) mifi.Customer {
    return mifi.Customer{
        ID:    ecomCustomer.ID,
        Name:  ecomCustomer.Name,
        Email: ecomCustomer.Email,
        Addresses: ecomCustomer.Addresses,
        PaymentTerms: "NET_30",
        CreditLimit:  mintyex.NewMoney(5000.00, "USD"),
        Status:       mintyex.StatusActive,
        CreatedAt:    ecomCustomer.CreatedAt,
    }
}

// Convert order items to invoice items
func convertOrderItemsToInvoiceItems(orderItems []mica.OrderItem) []mifi.InvoiceItem {
    return mintyex.Map(orderItems, func(item mica.OrderItem) mifi.InvoiceItem {
        return mifi.InvoiceItem{
            ID:          generateID("inv_item"),
            Description: item.ProductName,
            Quantity:    item.Quantity,
            UnitPrice:   item.UnitPrice,
            Total:       item.Total,
            Category:    item.Category,
        }
    })
}
```

---

## Testing Domain Logic

Domain logic testing is simplified by the zero-dependency architecture:

```go
func TestAccountCreation(t *testing.T) {
    service := mifi.NewFinanceService()
    
    account, err := service.CreateAccount(
        "Test Checking", 
        "checking",
        mintyex.NewMoney(1000.00, "USD"), 
        "customer123",
    )
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Checking", account.Name)
    assert.Equal(t, "checking", account.Type)
    assert.Equal(t, int64(100000), account.Balance.Amount) // $1000 in cents
    assert.Equal(t, mintyex.StatusActive, account.Status)
}

func TestInsufficientFundsTransaction(t *testing.T) {
    service := mifi.NewFinanceService()
    
    account, _ := service.CreateAccount("Test", "checking", 
        mintyex.NewMoney(100.00, "USD"), "customer1")
    
    // Try to debit more than available
    _, err := service.CreateTransaction(account.ID, 
        mintyex.NewMoney(200.00, "USD"), "Large withdrawal", "debit")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "insufficient funds")
}

func TestCrossModelValidation(t *testing.T) {
    invoice := mifi.Invoice{
        Number: "INV-001",
        Amount: mintyex.NewMoney(100.00, "USD"),
        DueDate: time.Now().AddDate(0, 0, -1), // Past due date
        Items: []mifi.InvoiceItem{},
    }
    
    errors := mifi.ValidateInvoice(invoice)
    
    assert.True(t, errors.HasErrors())
    assert.Contains(t, errors.Error(), "due date cannot be in the past")
    assert.Contains(t, errors.Error(), "must have at least one item")
}
```

**Testing Benefits:**
- No mocking of external dependencies required
- Fast test execution with no I/O
- Pure function testing enables property-based testing
- Business rule validation is comprehensive and isolated

---

## Extending with Custom Domains

The Minty System architecture makes it straightforward to add new business domains:

### Creating a Custom Domain

```go
// 1. Create domain package with zero dependencies
package inventory

import (
    "time"
    "github.com/ha1tch/mintyex"
)

// 2. Define domain entities
type Warehouse struct {
    ID       string           `json:"id"`
    Name     string           `json:"name"`
    Location mintyex.Address  `json:"location"`
    Capacity int              `json:"capacity"`
    Status   string           `json:"status"`
}

type StockItem struct {
    ID          string    `json:"id"`
    SKU         string    `json:"sku"`
    WarehouseID string    `json:"warehouse_id"`
    Quantity    int       `json:"quantity"`
    Location    string    `json:"location"`
    LastCount   time.Time `json:"last_count"`
}

// 3. Implement domain service
type InventoryService struct {
    warehouses []Warehouse
    stockItems []StockItem
}

func NewInventoryService() *InventoryService {
    return &InventoryService{
        warehouses: make([]Warehouse, 0),
        stockItems: make([]StockItem, 0),
    }
}

// 4. Add business operations
func (is *InventoryService) TransferStock(fromWarehouse, toWarehouse string, 
    sku string, quantity int) error {
    
    // Business logic for stock transfers
    fromStock := is.findStock(fromWarehouse, sku)
    if fromStock == nil {
        return errors.New("source stock not found")
    }
    
    if fromStock.Quantity < quantity {
        return errors.New("insufficient stock for transfer")
    }
    
    // Update source
    fromStock.Quantity -= quantity
    
    // Update or create destination
    toStock := is.findStock(toWarehouse, sku)
    if toStock != nil {
        toStock.Quantity += quantity
    } else {
        newStock := StockItem{
            ID:          generateID("stk"),
            SKU:         sku,
            WarehouseID: toWarehouse,
            Quantity:    quantity,
            LastCount:   time.Now(),
        }
        is.stockItems = append(is.stockItems, newStock)
    }
    
    return nil
}
```

### Creating Presentation Adapter

```go
// Create UI package for the new domain
package inventoryui

import (
    mi "github.com/ha1tch/minty"
    mui "github.com/ha1tch/mintyui"
    "your-app/inventory"
)

func WarehouseCard(theme mui.Theme, warehouse inventory.Warehouse) mi.H {
    return theme.Card(warehouse.Name, func(b *mi.Builder) mi.Node {
        return b.Div(
            b.P("Location: ", warehouse.Location.FormatOneLine()),
            b.P("Capacity: ", fmt.Sprintf("%d items", warehouse.Capacity)),
            StatusBadge(theme, warehouse.Status),
        )
    })
}

func InventoryDashboard(theme mui.Theme, service *inventory.InventoryService) mi.H {
    warehouses := service.GetWarehouses()
    lowStockItems := service.GetLowStockItems()
    
    return mui.Dashboard(theme, "Inventory Management",
        func(b *mi.Builder) mi.Node {
            return b.Div(
                b.Section(
                    b.H2("Warehouses"),
                    b.Div(mi.Class("warehouse-grid"),
                        mintyex.Map(warehouses, func(w inventory.Warehouse) mi.H {
                            return WarehouseCard(theme, w)
                        })...,
                    ),
                ),
                b.Section(
                    b.H2("Low Stock Alerts"),
                    LowStockTable(theme, lowStockItems),
                ),
            )
        },
    )
}
```

This architecture ensures that new domains integrate seamlessly with the existing system while maintaining clean separation of concerns and type safety throughout.

---

## Summary

The Minty System's domain architecture provides:

**Clean Separation**: Business logic with zero UI dependencies enables easy testing and maintenance.

**Rich Domain Models**: Complete business entities with validation, calculations, and proper encapsulation.

**Cross-Domain Integration**: Sophisticated workflows spanning multiple business domains with clear boundaries.

**Type Safety**: Compile-time checking throughout the entire business logic layer.

**Testing Simplicity**: Pure business logic testing without external dependencies or mocking.

**Extensibility**: Clear patterns for adding new business domains without disrupting existing functionality.

The domain layer forms the foundation of sophisticated business applications while maintaining the simplicity and type safety that makes Go development productive and maintainable.
