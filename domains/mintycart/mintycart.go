// Package mintycart provides pure e-commerce domain logic for the Minty System.
// This package contains NO UI dependencies and focuses solely on business logic.
package mintycart

import (
	"errors"
	"fmt"
	"sort"
	"time"

	mt "github.com/ha1tch/minty/mintytypes"
)

// =====================================================
// PURE BUSINESS TYPES (No UI Dependencies)
// =====================================================

// Product represents a product in the catalog
type Product struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	SKU         string           `json:"sku"`
	Price       mt.Money    `json:"price"`
	Category    string           `json:"category"`
	Brand       string           `json:"brand"`
	Weight      float64          `json:"weight"`
	Dimensions  Dimensions       `json:"dimensions"`
	Inventory   Inventory        `json:"inventory"`
	Images      []ProductImage   `json:"images"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"alt_text"`
	IsPrimary bool  `json:"is_primary"`
	SortOrder int   `json:"sort_order"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // inches, cm
}

// Inventory represents product inventory information
type Inventory struct {
	Quantity      int    `json:"quantity"`
	LowStockLevel int    `json:"low_stock_level"`
	Status        string `json:"status"` // in_stock, low_stock, out_of_stock
	LastUpdated   time.Time `json:"last_updated"`
}

// Cart represents a shopping cart
type Cart struct {
	ID         string           `json:"id"`
	CustomerID string           `json:"customer_id"`
	Items      []CartItem       `json:"items"`
	Subtotal   mt.Money    `json:"subtotal"`
	Tax        mt.Money    `json:"tax"`
	Shipping   mt.Money    `json:"shipping"`
	Total      mt.Money    `json:"total"`
	Status     string           `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	ExpiresAt  time.Time        `json:"expires_at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        string        `json:"id"`
	ProductID string        `json:"product_id"`
	Product   Product       `json:"product"`
	Quantity  int           `json:"quantity"`
	Price     mt.Money `json:"price"`     // Price at time of adding to cart
	Total     mt.Money `json:"total"`     // Price * Quantity
	AddedAt   time.Time     `json:"added_at"`
}

// Order represents a customer order
type Order struct {
	ID              string           `json:"id"`
	Number          string           `json:"number"`
	CustomerID      string           `json:"customer_id"`
	Customer        Customer         `json:"customer"`
	Items           []OrderItem      `json:"items"`
	BillingAddress  mt.Address  `json:"billing_address"`
	ShippingAddress mt.Address  `json:"shipping_address"`
	Payment         Payment          `json:"payment"`
	Subtotal        mt.Money    `json:"subtotal"`
	Tax             mt.Money    `json:"tax"`
	Shipping        mt.Money    `json:"shipping"`
	Discount        mt.Money    `json:"discount"`
	Total           mt.Money    `json:"total"`
	Status          string           `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	ShippedAt       *time.Time       `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time       `json:"delivered_at,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string        `json:"id"`
	ProductID string        `json:"product_id"`
	Product   Product       `json:"product"`
	Quantity  int           `json:"quantity"`
	Price     mt.Money `json:"price"`     // Price at time of order
	Total     mt.Money `json:"total"`     // Price * Quantity
}

// Payment represents payment information
type Payment struct {
	ID            string        `json:"id"`
	Method        string        `json:"method"`        // credit_card, paypal, bank_transfer
	Status        string        `json:"status"`        // pending, completed, failed, refunded
	Amount        mt.Money `json:"amount"`
	TransactionID string        `json:"transaction_id"`
	ProcessedAt   *time.Time    `json:"processed_at,omitempty"`
	CardLast4     string        `json:"card_last_4,omitempty"`
	CardBrand     string        `json:"card_brand,omitempty"`
}

// Customer represents an e-commerce customer
type Customer struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Addresses      []mt.Address  `json:"addresses"`
	Phone          string             `json:"phone"`
	LoyaltyPoints  int                `json:"loyalty_points"`
	TotalSpent     mt.Money      `json:"total_spent"`
	OrderCount     int                `json:"order_count"`
	PreferredPayment string           `json:"preferred_payment"`
	CreatedAt      time.Time          `json:"created_at"`
	LastOrderAt    *time.Time         `json:"last_order_at,omitempty"`
	Status         string             `json:"status"`
	Metadata       map[string]string  `json:"metadata,omitempty"`
}

// Implement mt.Customer interface
func (c Customer) GetID() string                { return c.ID }
func (c Customer) GetName() string              { return c.Name }
func (c Customer) GetEmail() string             { return c.Email }
func (c Customer) GetAddresses() []mt.Address { return c.Addresses }

func (c Customer) GetPrimaryAddress() mt.Address {
	for _, addr := range c.Addresses {
		if addr.Type == "primary" {
			return addr
		}
	}
	if len(c.Addresses) > 0 {
		return c.Addresses[0]
	}
	return mt.Address{}
}

func (c Customer) GetBillingAddress() mt.Address {
	for _, addr := range c.Addresses {
		if addr.Type == mt.AddressBilling {
			return addr
		}
	}
	return c.GetPrimaryAddress()
}

func (c Customer) GetShippingAddress() mt.Address {
	for _, addr := range c.Addresses {
		if addr.Type == mt.AddressShipping {
			return addr
		}
	}
	return c.GetPrimaryAddress()
}

// Category represents a product category
type Category struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ParentID    string           `json:"parent_id,omitempty"`
	Children    []Category       `json:"children,omitempty"`
	ProductCount int             `json:"product_count"`
	SortOrder   int              `json:"sort_order"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// =====================================================
// STATUS IMPLEMENTATIONS
// =====================================================

// ProductStatus implements mt.Status interface
type ProductStatus struct {
	status string
}

func NewProductStatus(status string) ProductStatus {
	return ProductStatus{status: status}
}

func (s ProductStatus) GetCode() string { return s.status }

func (s ProductStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusActive:    return "Active"
	case mt.StatusInactive:  return "Inactive"
	case mt.StatusDraft:     return "Draft"
	case "discontinued": return "Discontinued"
	default:             return "Unknown"
	}
}

func (s ProductStatus) IsActive() bool {
	return s.status == mt.StatusActive
}

func (s ProductStatus) GetSeverity() string {
	switch s.status {
	case mt.StatusActive:    return "success"
	case mt.StatusDraft:     return "info"
	case mt.StatusInactive:  return "warning"
	case "discontinued": return "error"
	default:             return "secondary"
	}
}

func (s ProductStatus) GetDescription() string {
	switch s.status {
	case mt.StatusActive:    return "Product is available for purchase"
	case mt.StatusInactive:  return "Product is temporarily unavailable"
	case mt.StatusDraft:     return "Product is in draft mode"
	case "discontinued": return "Product is no longer available"
	default:             return ""
	}
}

// OrderStatus implements mt.Status interface
type OrderStatus struct {
	status string
}

func NewOrderStatus(status string) OrderStatus {
	return OrderStatus{status: status}
}

func (s OrderStatus) GetCode() string { return s.status }

func (s OrderStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusPending:   return "Pending"
	case "processing":  return "Processing"
	case "shipped":     return "Shipped"
	case "delivered":   return "Delivered"
	case mt.StatusCancelled: return "Cancelled"
	case "returned":    return "Returned"
	default:            return "Unknown"
	}
}

func (s OrderStatus) IsActive() bool {
	return s.status == mt.StatusPending || s.status == "processing" || s.status == "shipped"
}

func (s OrderStatus) GetSeverity() string {
	switch s.status {
	case "delivered":   return "success"
	case "shipped":     return "info"
	case mt.StatusPending, "processing": return "warning"
	case mt.StatusCancelled, "returned": return "error"
	default:            return "secondary"
	}
}

func (s OrderStatus) GetDescription() string {
	switch s.status {
	case mt.StatusPending:   return "Order is awaiting processing"
	case "processing":  return "Order is being prepared"
	case "shipped":     return "Order has been shipped"
	case "delivered":   return "Order has been delivered"
	case mt.StatusCancelled: return "Order has been cancelled"
	case "returned":    return "Order has been returned"
	default:            return ""
	}
}

// =====================================================
// PURE BUSINESS LOGIC FUNCTIONS
// =====================================================

// Product Business Logic

// ValidateProduct validates product data
func ValidateProduct(product Product) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("name", product.Name, "Product Name", &errors)
	mt.ValidateRequired("sku", product.SKU, "SKU", &errors)
	mt.ValidateRequired("category", product.Category, "Category", &errors)
	mt.ValidateMoneyAmount("price", product.Price, "Price", &errors)
	
	if product.Weight <= 0 {
		errors.Add("weight", "Weight must be greater than zero")
	}
	
	if product.Inventory.Quantity < 0 {
		errors.Add("inventory.quantity", "Inventory quantity cannot be negative")
	}
	
	return errors
}

// UpdateInventory updates product inventory
func UpdateInventory(product *Product, quantityChange int) error {
	newQuantity := product.Inventory.Quantity + quantityChange
	
	if newQuantity < 0 {
		return errors.New("insufficient inventory")
	}
	
	product.Inventory.Quantity = newQuantity
	product.Inventory.LastUpdated = time.Now()
	product.UpdatedAt = time.Now()
	
	// Update inventory status
	if newQuantity == 0 {
		product.Inventory.Status = "out_of_stock"
	} else if newQuantity <= product.Inventory.LowStockLevel {
		product.Inventory.Status = "low_stock"
	} else {
		product.Inventory.Status = "in_stock"
	}
	
	return nil
}

// CalculateShippingWeight calculates total shipping weight for products
func CalculateShippingWeight(items []CartItem) float64 {
	totalWeight := 0.0
	for _, item := range items {
		totalWeight += item.Product.Weight * float64(item.Quantity)
	}
	return totalWeight
}

// Cart Business Logic

// ValidateCart validates cart data
func ValidateCart(cart Cart) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("customer_id", cart.CustomerID, "Customer ID", &errors)
	
	if len(cart.Items) == 0 {
		errors.Add("items", "Cart must have at least one item")
	}
	
	// Validate cart totals
	calculatedSubtotal := CalculateSubtotal(cart.Items)
	if cart.Subtotal.Amount != calculatedSubtotal.Amount {
		errors.Add("subtotal", "Cart subtotal does not match calculated total")
	}
	
	return errors
}

// AddItemToCart adds an item to the cart or updates quantity if item exists
func AddItemToCart(cart *Cart, product Product, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if product.Inventory.Quantity < quantity {
		return errors.New("insufficient inventory")
	}
	
	// Check if item already exists in cart
	for i, item := range cart.Items {
		if item.ProductID == product.ID {
			cart.Items[i].Quantity += quantity
			cart.Items[i].Total.Amount = cart.Items[i].Price.Amount * int64(cart.Items[i].Quantity)
			cart.UpdatedAt = time.Now()
			RecalculateCartTotals(cart)
			return nil
		}
	}
	
	// Add new item
	cartItem := CartItem{
		ID:        generateID("item"),
		ProductID: product.ID,
		Product:   product,
		Quantity:  quantity,
		Price:     product.Price,
		Total:     mt.Money{Amount: product.Price.Amount * int64(quantity), Currency: product.Price.Currency},
		AddedAt:   time.Now(),
	}
	
	cart.Items = append(cart.Items, cartItem)
	cart.UpdatedAt = time.Now()
	RecalculateCartTotals(cart)
	
	return nil
}

// RemoveItemFromCart removes an item from the cart
func RemoveItemFromCart(cart *Cart, itemID string) error {
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			cart.UpdatedAt = time.Now()
			RecalculateCartTotals(cart)
			return nil
		}
	}
	return errors.New("item not found in cart")
}

// UpdateItemQuantity updates the quantity of an item in the cart
func UpdateItemQuantity(cart *Cart, itemID string, newQuantity int) error {
	if newQuantity <= 0 {
		return RemoveItemFromCart(cart, itemID)
	}
	
	for i, item := range cart.Items {
		if item.ID == itemID {
			if item.Product.Inventory.Quantity < newQuantity {
				return errors.New("insufficient inventory")
			}
			
			cart.Items[i].Quantity = newQuantity
			cart.Items[i].Total.Amount = cart.Items[i].Price.Amount * int64(newQuantity)
			cart.UpdatedAt = time.Now()
			RecalculateCartTotals(cart)
			return nil
		}
	}
	return errors.New("item not found in cart")
}

// RecalculateCartTotals recalculates all cart totals
func RecalculateCartTotals(cart *Cart) {
	cart.Subtotal = CalculateSubtotal(cart.Items)
	cart.Tax = CalculateTax(cart.Subtotal, 0.08) // 8% tax rate
	cart.Shipping = CalculateShipping(cart.Items)
	
	cart.Total.Amount = cart.Subtotal.Amount + cart.Tax.Amount + cart.Shipping.Amount
	cart.Total.Currency = cart.Subtotal.Currency
}

// CalculateSubtotal calculates subtotal from cart items
func CalculateSubtotal(items []CartItem) mt.Money {
	var subtotal mt.Money
	for _, item := range items {
		subtotal.Amount += item.Total.Amount
	}
	return subtotal
}

// CalculateTax calculates tax amount
func CalculateTax(subtotal mt.Money, taxRate float64) mt.Money {
	taxAmount := float64(subtotal.Amount) * taxRate
	return mt.Money{Amount: int64(taxAmount), Currency: subtotal.Currency}
}

// CalculateShipping calculates shipping cost based on weight and value
func CalculateShipping(items []CartItem) mt.Money {
	totalWeight := CalculateShippingWeight(items)
	
	var subtotal mt.Money
	for _, item := range items {
		subtotal.Amount += item.Total.Amount
	}
	
	// Free shipping for orders over $100
	if subtotal.Amount >= 10000 { // $100 in cents
		return mt.Money{Amount: 0, Currency: subtotal.Currency}
	}
	
	// Base shipping cost + weight-based cost
	baseCost := 5.00
	weightCost := totalWeight * 0.50
	totalShipping := baseCost + weightCost
	
	return mt.NewMoney(totalShipping, subtotal.Currency)
}

// Order Business Logic

// ValidateOrder validates order data
func ValidateOrder(order Order) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("number", order.Number, "Order Number", &errors)
	mt.ValidateRequired("customer_id", order.CustomerID, "Customer ID", &errors)
	mt.ValidateRequired("customer.name", order.Customer.Name, "Customer Name", &errors)
	mt.ValidateEmail("customer.email", order.Customer.Email, "Customer Email", &errors)
	
	if len(order.Items) == 0 {
		errors.Add("items", "Order must have at least one item")
	}
	
	// Validate billing address
	if order.BillingAddress.Street1 == "" || order.BillingAddress.City == "" {
		errors.Add("billing_address", "Billing address must have street and city")
	}
	
	// Validate shipping address
	if order.ShippingAddress.Street1 == "" || order.ShippingAddress.City == "" {
		errors.Add("shipping_address", "Shipping address must have street and city")
	}
	
	return errors
}

// CreateOrderFromCart creates an order from a cart
func CreateOrderFromCart(cart Cart, customer Customer, billingAddr, shippingAddr mt.Address) Order {
	orderItems := make([]OrderItem, len(cart.Items))
	for i, cartItem := range cart.Items {
		orderItems[i] = OrderItem{
			ID:        generateID("oi"),
			ProductID: cartItem.ProductID,
			Product:   cartItem.Product,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
			Total:     cartItem.Total,
		}
	}
	
	order := Order{
		ID:              generateID("ord"),
		Number:          generateOrderNumber(),
		CustomerID:      customer.ID,
		Customer:        customer,
		Items:           orderItems,
		BillingAddress:  billingAddr,
		ShippingAddress: shippingAddr,
		Subtotal:        cart.Subtotal,
		Tax:             cart.Tax,
		Shipping:        cart.Shipping,
		Total:           cart.Total,
		Status:          mt.StatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Metadata:        make(map[string]string),
	}
	
	return order
}

// ProcessPayment processes payment for an order
func ProcessPayment(order *Order, paymentMethod string) error {
	// Simulate payment processing
	payment := Payment{
		ID:            generateID("pay"),
		Method:        paymentMethod,
		Status:        "completed", // Simplified - would be "pending" in real scenario
		Amount:        order.Total,
		TransactionID: generateTransactionID(),
	}
	
	if paymentMethod == "credit_card" {
		payment.CardLast4 = "4242"
		payment.CardBrand = "visa"
	}
	
	now := time.Now()
	payment.ProcessedAt = &now
	
	order.Payment = payment
	order.Status = "processing"
	order.UpdatedAt = time.Now()
	
	return nil
}

// ShipOrder marks an order as shipped
func ShipOrder(order *Order, trackingNumber string) error {
	if order.Status != "processing" {
		return errors.New("order must be in processing status to ship")
	}
	
	order.Status = "shipped"
	now := time.Now()
	order.ShippedAt = &now
	order.UpdatedAt = time.Now()
	
	if order.Metadata == nil {
		order.Metadata = make(map[string]string)
	}
	order.Metadata["tracking_number"] = trackingNumber
	
	return nil
}

// =====================================================
// DOMAIN SERVICES
// =====================================================

// EcommerceService provides business operations for the e-commerce domain
type EcommerceService struct {
	products   []Product
	categories []Category
	carts      []Cart
	orders     []Order
	customers  []Customer
}

// NewEcommerceService creates a new e-commerce service
func NewEcommerceService() *EcommerceService {
	return &EcommerceService{
		products:   make([]Product, 0),
		categories: make([]Category, 0),
		carts:      make([]Cart, 0),
		orders:     make([]Order, 0),
		customers:  make([]Customer, 0),
	}
}

// Product Operations

func (es *EcommerceService) CreateProduct(name, description, sku, category string, 
	price mt.Money, weight float64, inventory Inventory) (*Product, error) {
	
	product := Product{
		ID:          generateID("prd"),
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
		Category:    category,
		Weight:      weight,
		Inventory:   inventory,
		Status:      mt.StatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}
	
	if errors := ValidateProduct(product); errors.HasErrors() {
		return nil, errors
	}
	
	es.products = append(es.products, product)
	return &product, nil
}

func (es *EcommerceService) GetProduct(productID string) (*Product, error) {
	for i, product := range es.products {
		if product.ID == productID {
			return &es.products[i], nil
		}
	}
	return nil, errors.New("product not found")
}

func (es *EcommerceService) GetProductsBySKU(sku string) (*Product, error) {
	for i, product := range es.products {
		if product.SKU == sku {
			return &es.products[i], nil
		}
	}
	return nil, errors.New("product not found")
}

func (es *EcommerceService) GetProductsByCategory(category string) []Product {
	var categoryProducts []Product
	for _, product := range es.products {
		if product.Category == category && product.Status == mt.StatusActive {
			categoryProducts = append(categoryProducts, product)
		}
	}
	return categoryProducts
}

func (es *EcommerceService) GetActiveProducts() []Product {
	var activeProducts []Product
	for _, product := range es.products {
		if product.Status == mt.StatusActive {
			activeProducts = append(activeProducts, product)
		}
	}
	return activeProducts
}

func (es *EcommerceService) GetAllProducts() []Product {
	return es.products
}

func (es *EcommerceService) UpdateProductInventory(productID string, quantityChange int) error {
	product, err := es.GetProduct(productID)
	if err != nil {
		return err
	}
	
	return UpdateInventory(product, quantityChange)
}

// Cart Operations

func (es *EcommerceService) CreateCart(customerID string) (*Cart, error) {
	cart := Cart{
		ID:         generateID("cart"),
		CustomerID: customerID,
		Items:      make([]CartItem, 0),
		Status:     mt.StatusActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour), // Cart expires in 24 hours
		Metadata:   make(map[string]string),
	}
	
	es.carts = append(es.carts, cart)
	return &cart, nil
}

func (es *EcommerceService) GetCart(cartID string) (*Cart, error) {
	for i, cart := range es.carts {
		if cart.ID == cartID {
			return &es.carts[i], nil
		}
	}
	return nil, errors.New("cart not found")
}

func (es *EcommerceService) GetCartByCustomer(customerID string) (*Cart, error) {
	for i, cart := range es.carts {
		if cart.CustomerID == customerID && cart.Status == mt.StatusActive {
			return &es.carts[i], nil
		}
	}
	return nil, errors.New("active cart not found for customer")
}

func (es *EcommerceService) AddToCart(cartID, productID string, quantity int) error {
	cart, err := es.GetCart(cartID)
	if err != nil {
		return err
	}
	
	product, err := es.GetProduct(productID)
	if err != nil {
		return err
	}
	
	return AddItemToCart(cart, *product, quantity)
}

func (es *EcommerceService) RemoveFromCart(cartID, itemID string) error {
	cart, err := es.GetCart(cartID)
	if err != nil {
		return err
	}
	
	return RemoveItemFromCart(cart, itemID)
}

func (es *EcommerceService) UpdateCartItemQuantity(cartID, itemID string, quantity int) error {
	cart, err := es.GetCart(cartID)
	if err != nil {
		return err
	}
	
	return UpdateItemQuantity(cart, itemID, quantity)
}

// Order Operations

func (es *EcommerceService) CreateOrder(cartID string, customer Customer, 
	billingAddr, shippingAddr mt.Address, paymentMethod string) (*Order, error) {
	
	cart, err := es.GetCart(cartID)
	if err != nil {
		return nil, err
	}
	
	if len(cart.Items) == 0 {
		return nil, errors.New("cannot create order from empty cart")
	}
	
	order := CreateOrderFromCart(*cart, customer, billingAddr, shippingAddr)
	
	if errors := ValidateOrder(order); errors.HasErrors() {
		return nil, errors
	}
	
	// Process payment
	if err := ProcessPayment(&order, paymentMethod); err != nil {
		return nil, err
	}
	
	// Update inventory
	for _, item := range order.Items {
		if err := es.UpdateProductInventory(item.ProductID, -item.Quantity); err != nil {
			// In a real system, this would require transaction rollback
			return nil, fmt.Errorf("inventory update failed for product %s: %w", item.ProductID, err)
		}
	}
	
	// Mark cart as ordered
	cart.Status = "ordered"
	cart.UpdatedAt = time.Now()
	
	es.orders = append(es.orders, order)
	return &order, nil
}

func (es *EcommerceService) GetOrder(orderID string) (*Order, error) {
	for i, order := range es.orders {
		if order.ID == orderID {
			return &es.orders[i], nil
		}
	}
	return nil, errors.New("order not found")
}

func (es *EcommerceService) GetOrdersByCustomer(customerID string) []Order {
	var customerOrders []Order
	for _, order := range es.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	
	// Sort by created date descending
	sort.Slice(customerOrders, func(i, j int) bool {
		return customerOrders[i].CreatedAt.After(customerOrders[j].CreatedAt)
	})
	
	return customerOrders
}

func (es *EcommerceService) GetRecentOrders(limit int) []Order {
	// Sort all orders by date
	allOrders := make([]Order, len(es.orders))
	copy(allOrders, es.orders)
	
	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].CreatedAt.After(allOrders[j].CreatedAt)
	})
	
	if limit > len(allOrders) {
		limit = len(allOrders)
	}
	
	return allOrders[:limit]
}

func (es *EcommerceService) GetAllOrders() []Order {
	return es.orders
}

func (es *EcommerceService) ShipOrder(orderID, trackingNumber string) error {
	order, err := es.GetOrder(orderID)
	if err != nil {
		return err
	}
	
	return ShipOrder(order, trackingNumber)
}

// Customer Operations

func (es *EcommerceService) CreateCustomer(name, email string) (*Customer, error) {
	customer := Customer{
		ID:        generateID("cust"),
		Name:      name,
		Email:     email,
		Status:    mt.StatusActive,
		CreatedAt: time.Now(),
		Metadata:  make(map[string]string),
	}
	
	var errors mt.ValidationErrors
	mt.ValidateRequired("name", customer.Name, "Customer Name", &errors)
	mt.ValidateEmail("email", customer.Email, "Email", &errors)
	
	if errors.HasErrors() {
		return nil, errors
	}
	
	es.customers = append(es.customers, customer)
	return &customer, nil
}

func (es *EcommerceService) GetCustomer(customerID string) (*Customer, error) {
	for i, customer := range es.customers {
		if customer.ID == customerID {
			return &es.customers[i], nil
		}
	}
	return nil, errors.New("customer not found")
}

func (es *EcommerceService) GetAllCustomers() []Customer {
	return es.customers
}

// =====================================================
// DATA TRANSFER OBJECTS FOR PRESENTATION LAYER
// =====================================================

// ProductDisplayData prepares product data for UI display
type ProductDisplayData struct {
	Product          Product
	FormattedPrice   string
	InventoryStatus  string
	StatusClass      string
	StatusDisplay    string
	PrimaryImageURL  string
	InStock          bool
}

// CartDisplayData prepares cart data for UI display
type CartDisplayData struct {
	Cart              Cart
	ItemCount         int
	FormattedSubtotal string
	FormattedTax      string
	FormattedShipping string
	FormattedTotal    string
	IsEmpty           bool
}

// OrderDisplayData prepares order data for UI display
type OrderDisplayData struct {
	Order            Order
	FormattedTotal   string
	StatusClass      string
	StatusDisplay    string
	DaysAgo          int
	TrackingNumber   string
}

// DashboardData aggregates e-commerce data for dashboard display
type DashboardData struct {
	TotalProducts     int
	ActiveProducts    int
	LowStockProducts  int
	TotalOrders       int
	PendingOrders     int
	RevenueToday      mt.Money
	FormattedRevenue  string
	RecentOrders      []OrderDisplayData
	TopProducts       []ProductDisplayData
}

// =====================================================
// DATA PREPARATION FUNCTIONS (No UI Rendering)
// =====================================================

// PrepareProductForDisplay prepares product data for presentation layer
func PrepareProductForDisplay(product Product) ProductDisplayData {
	status := NewProductStatus(product.Status)
	
	var primaryImageURL string
	for _, img := range product.Images {
		if img.IsPrimary {
			primaryImageURL = img.URL
			break
		}
	}
	if primaryImageURL == "" && len(product.Images) > 0 {
		primaryImageURL = product.Images[0].URL
	}
	
	return ProductDisplayData{
		Product:          product,
		FormattedPrice:   product.Price.Format(),
		InventoryStatus:  product.Inventory.Status,
		StatusClass:      "status-" + status.GetSeverity(),
		StatusDisplay:    status.GetDisplay(),
		PrimaryImageURL:  primaryImageURL,
		InStock:          product.Inventory.Status == "in_stock",
	}
}

// PrepareCartForDisplay prepares cart data for presentation layer
func PrepareCartForDisplay(cart Cart) CartDisplayData {
	itemCount := 0
	for _, item := range cart.Items {
		itemCount += item.Quantity
	}
	
	return CartDisplayData{
		Cart:              cart,
		ItemCount:         itemCount,
		FormattedSubtotal: cart.Subtotal.Format(),
		FormattedTax:      cart.Tax.Format(),
		FormattedShipping: cart.Shipping.Format(),
		FormattedTotal:    cart.Total.Format(),
		IsEmpty:           len(cart.Items) == 0,
	}
}

// PrepareOrderForDisplay prepares order data for presentation layer
func PrepareOrderForDisplay(order Order) OrderDisplayData {
	status := NewOrderStatus(order.Status)
	daysAgo := int(time.Since(order.CreatedAt).Hours() / 24)
	trackingNumber := ""
	if order.Metadata != nil {
		trackingNumber = order.Metadata["tracking_number"]
	}
	
	return OrderDisplayData{
		Order:          order,
		FormattedTotal: order.Total.Format(),
		StatusClass:    "status-" + status.GetSeverity(),
		StatusDisplay:  status.GetDisplay(),
		DaysAgo:        daysAgo,
		TrackingNumber: trackingNumber,
	}
}

// PrepareDashboardData aggregates e-commerce data for dashboard presentation
func PrepareDashboardData(es *EcommerceService) DashboardData {
	allProducts := es.GetAllProducts()
	activeProducts := es.GetActiveProducts()
	allOrders := es.GetAllOrders()
	
	lowStockProducts := 0
	pendingOrders := 0
	
	for _, product := range allProducts {
		if product.Inventory.Status == "low_stock" || product.Inventory.Status == "out_of_stock" {
			lowStockProducts++
		}
	}
	
	var revenueToday mt.Money
	today := time.Now().Truncate(24 * time.Hour)
	
	for _, order := range allOrders {
		if order.Status == mt.StatusPending {
			pendingOrders++
		}
		if order.CreatedAt.After(today) && order.Payment.Status == "completed" {
			revenueToday.Amount += order.Total.Amount
		}
	}
	
	// Prepare recent orders for display
	var recentOrders []OrderDisplayData
	recentOrdersList := es.GetRecentOrders(5)
	for _, order := range recentOrdersList {
		recentOrders = append(recentOrders, PrepareOrderForDisplay(order))
	}
	
	// Get top products (simplified - just first few active products)
	var topProducts []ProductDisplayData
	for i, product := range activeProducts {
		if i < 5 {
			topProducts = append(topProducts, PrepareProductForDisplay(product))
		}
	}
	
	return DashboardData{
		TotalProducts:     len(allProducts),
		ActiveProducts:    len(activeProducts),
		LowStockProducts:  lowStockProducts,
		TotalOrders:       len(allOrders),
		PendingOrders:     pendingOrders,
		RevenueToday:      revenueToday,
		FormattedRevenue:  revenueToday.Format(),
		RecentOrders:      recentOrders,
		TopProducts:       topProducts,
	}
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// generateID generates a unique ID with prefix
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// generateOrderNumber generates a unique order number
func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%d%02d%02d-%d", now.Year(), now.Month(), now.Day(), now.UnixNano()%10000)
}

// generateTransactionID generates a unique transaction ID
func generateTransactionID() string {
	return fmt.Sprintf("TXN_%d", time.Now().UnixNano())
}

// =====================================================
// SAMPLE DATA HELPERS (for demos and testing)
// =====================================================

// SampleCustomer returns sample customer data
func SampleCustomer() Customer {
	return Customer{
		ID:    "cust_001",
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
		Addresses: []mt.Address{
			{
				Type:       mt.AddressBilling,
				Name:       "Jane Smith",
				Street1:    "456 Oak Ave",
				City:       "Springfield",
				State:      "IL",
				PostalCode: "62701",
				Country:    "US",
			},
			{
				Type:       mt.AddressShipping,
				Name:       "Jane Smith",
				Street1:    "456 Oak Ave",
				City:       "Springfield",
				State:      "IL",
				PostalCode: "62701",
				Country:    "US",
			},
		},
		Phone:         "555-0123",
		LoyaltyPoints: 150,
		OrderCount:    5,
		Status:        mt.StatusActive,
		CreatedAt:     time.Now().AddDate(-1, 0, 0),
		Metadata:      make(map[string]string),
	}
}

// SampleProducts returns sample product data for demos
func SampleProducts() []Product {
	return []Product{
		{
			ID:          "prd_001",
			Name:        "Wireless Bluetooth Headphones",
			Description: "High-quality wireless headphones with noise cancellation",
			SKU:         "WBH-001",
			Price:       mt.NewMoney(99.99, mt.CurrencyUSD),
			Category:    "Electronics",
			Brand:       "TechCorp",
			Weight:      1.2,
			Inventory: Inventory{
				Quantity:      50,
				LowStockLevel: 10,
				Status:        "in_stock",
				LastUpdated:   time.Now(),
			},
			Status:    mt.StatusActive,
			CreatedAt: time.Now().AddDate(0, -2, 0),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
		{
			ID:          "prd_002",
			Name:        "Organic Coffee Beans",
			Description: "Premium organic coffee beans from Colombia",
			SKU:         "OCB-001",
			Price:       mt.NewMoney(24.99, mt.CurrencyUSD),
			Category:    "Food & Beverages",
			Brand:       "Mountain Roasters",
			Weight:      2.0,
			Inventory: Inventory{
				Quantity:      25,
				LowStockLevel: 20,
				Status:        "low_stock",
				LastUpdated:   time.Now(),
			},
			Status:    mt.StatusActive,
			CreatedAt: time.Now().AddDate(0, -1, -15),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
		{
			ID:          "prd_003",
			Name:        "Yoga Mat",
			Description: "Non-slip yoga mat for all types of yoga practice",
			SKU:         "YM-001",
			Price:       mt.NewMoney(39.99, mt.CurrencyUSD),
			Category:    "Sports & Fitness",
			Brand:       "ZenFit",
			Weight:      3.5,
			Inventory: Inventory{
				Quantity:      15,
				LowStockLevel: 5,
				Status:        "in_stock",
				LastUpdated:   time.Now(),
			},
			Status:    mt.StatusActive,
			CreatedAt: time.Now().AddDate(0, -3, 0),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
	}
}

// SampleOrders returns sample order data
func SampleOrders() []Order {
	customer := SampleCustomer()
	
	return []Order{
		{
			ID:         "ord_001",
			Number:     generateOrderNumber(),
			CustomerID: customer.ID,
			Customer:   customer,
			Items: []OrderItem{
				{
					ID:        "oi_001",
					ProductID: "prd_001",
					Quantity:  1,
					Price:     mt.NewMoney(99.99, mt.CurrencyUSD),
					Total:     mt.NewMoney(99.99, mt.CurrencyUSD),
				},
			},
			BillingAddress:  customer.GetBillingAddress(),
			ShippingAddress: customer.GetShippingAddress(),
			Payment: Payment{
				ID:     "pay_001",
				Method: "credit_card",
				Status: "completed",
				Amount: mt.NewMoney(107.99, mt.CurrencyUSD),
			},
			Subtotal:  mt.NewMoney(99.99, mt.CurrencyUSD),
			Tax:       mt.NewMoney(8.00, mt.CurrencyUSD),
			Shipping:  mt.NewMoney(0.00, mt.CurrencyUSD),
			Total:     mt.NewMoney(107.99, mt.CurrencyUSD),
			Status:    "delivered",
			CreatedAt: time.Now().AddDate(0, 0, -5),
			UpdatedAt: time.Now().AddDate(0, 0, -2),
			Metadata:  make(map[string]string),
		},
	}
}
