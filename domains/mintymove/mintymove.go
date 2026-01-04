// Package mintymove provides pure logistics and movement domain logic for the Minty System.
// This package contains NO UI dependencies and focuses solely on business logic.
package mintymove

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

// Shipment represents a shipment with tracking and logistics data
type Shipment struct {
	ID              string           `json:"id"`
	TrackingCode    string           `json:"tracking_code"`
	Origin          mt.Address  `json:"origin"`
	Destination     mt.Address  `json:"destination"`
	Status          string           `json:"status"`
	EstimatedDate   time.Time        `json:"estimated_date"`
	ActualDate      *time.Time       `json:"actual_date,omitempty"`
	Carrier         string           `json:"carrier"`
	Service         string           `json:"service"`
	Weight          float64          `json:"weight"`
	Cost            mt.Money    `json:"cost"`
	Items           []ShipmentItem   `json:"items"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// ShipmentItem represents an item in a shipment
type ShipmentItem struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Quantity    int           `json:"quantity"`
	Weight      float64       `json:"weight"`
	Value       mt.Money `json:"value"`
	SKU         string        `json:"sku"`
	Category    string        `json:"category"`
}

// Route represents a delivery route
type Route struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Origin       mt.Address  `json:"origin"`
	Destination  mt.Address  `json:"destination"`
	Distance     float64          `json:"distance"` // in miles
	Duration     time.Duration    `json:"duration"`
	Cost         mt.Money    `json:"cost"`
	Stops        []RouteStop      `json:"stops"`
	Status       string           `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// RouteStop represents a stop on a route
type RouteStop struct {
	ID            string          `json:"id"`
	Address       mt.Address `json:"address"`
	EstimatedTime time.Time       `json:"estimated_time"`
	ActualTime    *time.Time      `json:"actual_time,omitempty"`
	Type          string          `json:"type"` // pickup, delivery, waypoint
	Status        string          `json:"status"`
	Instructions  string          `json:"instructions"`
}

// Vehicle represents a delivery vehicle
type Vehicle struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Type         string           `json:"type"` // truck, van, car, bike
	LicensePlate string           `json:"license_plate"`
	Capacity     VehicleCapacity  `json:"capacity"`
	Status       string           `json:"status"`
	Location     Location         `json:"location"`
	Driver       Driver           `json:"driver"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// VehicleCapacity represents vehicle capacity constraints
type VehicleCapacity struct {
	Weight     float64 `json:"weight"`      // maximum weight in lbs
	Volume     float64 `json:"volume"`      // maximum volume in cubic feet
	ItemCount  int     `json:"item_count"`  // maximum number of items
}

// Location represents a geographic location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Driver represents a delivery driver
type Driver struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	LicenseNum  string           `json:"license_number"`
	Status      string           `json:"status"`
	Rating      float64          `json:"rating"`
	CreatedAt   time.Time        `json:"created_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Customer represents a logistics customer
type Customer struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Addresses      []mt.Address  `json:"addresses"`
	AccountNumber  string             `json:"account_number"`
	PreferredCarrier string           `json:"preferred_carrier"`
	CreditLimit    mt.Money      `json:"credit_limit"`
	TotalSpent     mt.Money      `json:"total_spent"`
	CreatedAt      time.Time          `json:"created_at"`
	LastActivityAt time.Time          `json:"last_activity_at"`
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

// =====================================================
// STATUS IMPLEMENTATIONS
// =====================================================

// ShipmentStatus implements mt.Status interface
type ShipmentStatus struct {
	status string
}

func NewShipmentStatus(status string) ShipmentStatus {
	return ShipmentStatus{status: status}
}

func (s ShipmentStatus) GetCode() string { return s.status }

func (s ShipmentStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusPending:    return "Pending"
	case "picked_up":    return "Picked Up"
	case "in_transit":   return "In Transit"
	case "out_for_delivery": return "Out for Delivery"
	case "delivered":    return "Delivered"
	case "exception":    return "Exception"
	case mt.StatusCancelled: return "Cancelled"
	default:             return "Unknown"
	}
}

func (s ShipmentStatus) IsActive() bool {
	return s.status == mt.StatusPending || s.status == "picked_up" || 
		   s.status == "in_transit" || s.status == "out_for_delivery"
}

func (s ShipmentStatus) GetSeverity() string {
	switch s.status {
	case "delivered":    return "success"
	case mt.StatusPending: return "info"
	case "picked_up", "in_transit", "out_for_delivery": return "warning"
	case "exception":    return "error"
	case mt.StatusCancelled: return "secondary"
	default:             return "info"
	}
}

func (s ShipmentStatus) GetDescription() string {
	switch s.status {
	case mt.StatusPending:    return "Shipment is being prepared"
	case "picked_up":    return "Shipment has been picked up"
	case "in_transit":   return "Shipment is in transit"
	case "out_for_delivery": return "Shipment is out for delivery"
	case "delivered":    return "Shipment has been delivered"
	case "exception":    return "There is an issue with the shipment"
	case mt.StatusCancelled: return "Shipment has been cancelled"
	default:             return ""
	}
}

// RouteStatus implements mt.Status interface  
type RouteStatus struct {
	status string
}

func NewRouteStatus(status string) RouteStatus {
	return RouteStatus{status: status}
}

func (s RouteStatus) GetCode() string { return s.status }

func (s RouteStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusPending:   return "Pending"
	case "active":      return "Active"
	case mt.StatusCompleted: return "Completed"
	case "delayed":     return "Delayed"
	case mt.StatusCancelled: return "Cancelled"
	default:            return "Unknown"
	}
}

func (s RouteStatus) IsActive() bool {
	return s.status == mt.StatusPending || s.status == "active"
}

func (s RouteStatus) GetSeverity() string {
	switch s.status {
	case mt.StatusCompleted: return "success"
	case "active":      return "info"
	case mt.StatusPending:   return "warning"
	case "delayed":     return "error"
	case mt.StatusCancelled: return "secondary"
	default:            return "info"
	}
}

func (s RouteStatus) GetDescription() string {
	switch s.status {
	case mt.StatusPending:   return "Route is scheduled but not started"
	case "active":      return "Route is currently in progress"
	case mt.StatusCompleted: return "Route has been completed"
	case "delayed":     return "Route is behind schedule"
	case mt.StatusCancelled: return "Route has been cancelled"
	default:            return ""
	}
}

// =====================================================
// PURE BUSINESS LOGIC FUNCTIONS
// =====================================================

// Shipment Business Logic

// ValidateShipment validates shipment data
func ValidateShipment(shipment Shipment) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("tracking_code", shipment.TrackingCode, "Tracking Code", &errors)
	mt.ValidateRequired("carrier", shipment.Carrier, "Carrier", &errors)
	mt.ValidateRequired("service", shipment.Service, "Service", &errors)
	
	if shipment.Weight <= 0 {
		errors.Add("weight", "Weight must be greater than zero")
	}
	
	if len(shipment.Items) == 0 {
		errors.Add("items", "Shipment must have at least one item")
	}
	
	// Validate origin and destination addresses
	if shipment.Origin.Street1 == "" || shipment.Origin.City == "" {
		errors.Add("origin", "Origin address must have street and city")
	}
	
	if shipment.Destination.Street1 == "" || shipment.Destination.City == "" {
		errors.Add("destination", "Destination address must have street and city")
	}
	
	return errors
}

// CalculateShipmentCost calculates shipping cost based on weight, distance, and service
func CalculateShipmentCost(weight float64, distance float64, service string) mt.Money {
	var baseCost float64
	var perMileCost float64
	var weightMultiplier float64
	
	// Service-based pricing
	switch service {
	case "standard":
		baseCost = 5.00
		perMileCost = 0.10
		weightMultiplier = 0.50
	case "express":
		baseCost = 12.00
		perMileCost = 0.15
		weightMultiplier = 0.75
	case "overnight":
		baseCost = 25.00
		perMileCost = 0.25
		weightMultiplier = 1.00
	default:
		baseCost = 5.00
		perMileCost = 0.10
		weightMultiplier = 0.50
	}
	
	totalCost := baseCost + (distance * perMileCost) + (weight * weightMultiplier)
	return mt.NewMoney(totalCost, mt.CurrencyUSD)
}

// EstimateDeliveryTime estimates delivery time based on distance and service
func EstimateDeliveryTime(distance float64, service string) time.Duration {
	var baseHours float64
	var hoursPerMile float64
	
	switch service {
	case "standard":
		baseHours = 48
		hoursPerMile = 0.1
	case "express":
		baseHours = 24
		hoursPerMile = 0.05
	case "overnight":
		baseHours = 12
		hoursPerMile = 0.02
	default:
		baseHours = 48
		hoursPerMile = 0.1
	}
	
	totalHours := baseHours + (distance * hoursPerMile)
	return time.Duration(totalHours) * time.Hour
}

// UpdateShipmentStatus updates shipment status with timestamp
func UpdateShipmentStatus(shipment *Shipment, newStatus string) {
	shipment.Status = newStatus
	shipment.UpdatedAt = time.Now()
	
	// Set actual delivery date if delivered
	if newStatus == "delivered" && shipment.ActualDate == nil {
		now := time.Now()
		shipment.ActualDate = &now
	}
}

// Route Business Logic

// ValidateRoute validates route data
func ValidateRoute(route Route) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("name", route.Name, "Route Name", &errors)
	
	if route.Distance <= 0 {
		errors.Add("distance", "Distance must be greater than zero")
	}
	
	if route.Duration <= 0 {
		errors.Add("duration", "Duration must be greater than zero")
	}
	
	if len(route.Stops) == 0 {
		errors.Add("stops", "Route must have at least one stop")
	}
	
	return errors
}

// OptimizeRoute optimizes route stops for efficiency (simple implementation)
func OptimizeRoute(route *Route) {
	if len(route.Stops) <= 2 {
		return // No optimization needed for routes with 2 or fewer stops
	}
	
	// Simple optimization: sort by estimated time
	sort.Slice(route.Stops, func(i, j int) bool {
		return route.Stops[i].EstimatedTime.Before(route.Stops[j].EstimatedTime)
	})
	
	route.UpdatedAt = time.Now()
}

// CalculateRouteDistance calculates total route distance (simplified)
func CalculateRouteDistance(stops []RouteStop) float64 {
	// Simple implementation - in reality would use mapping service
	totalDistance := 0.0
	
	for i := 0; i < len(stops)-1; i++ {
		// Simplified distance calculation
		totalDistance += 10.0 // Assume 10 miles between stops
	}
	
	return totalDistance
}

// Vehicle Business Logic

// ValidateVehicle validates vehicle data
func ValidateVehicle(vehicle Vehicle) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("name", vehicle.Name, "Vehicle Name", &errors)
	mt.ValidateRequired("type", vehicle.Type, "Vehicle Type", &errors)
	mt.ValidateRequired("license_plate", vehicle.LicensePlate, "License Plate", &errors)
	
	validTypes := []string{"truck", "van", "car", "bike", "motorcycle"}
	isValidType := false
	for _, validType := range validTypes {
		if vehicle.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		errors.Add("type", "Vehicle type must be one of: truck, van, car, bike, motorcycle")
	}
	
	if vehicle.Capacity.Weight <= 0 {
		errors.Add("capacity.weight", "Vehicle weight capacity must be greater than zero")
	}
	
	return errors
}

// CheckVehicleCapacity checks if vehicle can handle shipment
func CheckVehicleCapacity(vehicle Vehicle, shipment Shipment) bool {
	totalWeight := shipment.Weight
	totalItems := len(shipment.Items)
	
	return totalWeight <= vehicle.Capacity.Weight && 
		   totalItems <= vehicle.Capacity.ItemCount
}

// AssignDriverToVehicle assigns a driver to a vehicle
func AssignDriverToVehicle(vehicle *Vehicle, driver Driver) error {
	if driver.Status != mt.StatusActive {
		return errors.New("driver is not active")
	}
	
	vehicle.Driver = driver
	vehicle.UpdatedAt = time.Now()
	return nil
}

// =====================================================
// DOMAIN SERVICES
// =====================================================

// LogisticsService provides business operations for the logistics domain
type LogisticsService struct {
	shipments []Shipment
	routes    []Route
	vehicles  []Vehicle
	drivers   []Driver
	customers []Customer
}

// NewLogisticsService creates a new logistics service
func NewLogisticsService() *LogisticsService {
	return &LogisticsService{
		shipments: make([]Shipment, 0),
		routes:    make([]Route, 0),
		vehicles:  make([]Vehicle, 0),
		drivers:   make([]Driver, 0),
		customers: make([]Customer, 0),
	}
}

// Shipment Operations

func (ls *LogisticsService) CreateShipment(trackingCode string, origin, destination mt.Address,
	carrier, service string, weight float64, items []ShipmentItem) (*Shipment, error) {
	
	distance := calculateDistance(origin, destination) // Simplified
	cost := CalculateShipmentCost(weight, distance, service)
	estimatedDelivery := time.Now().Add(EstimateDeliveryTime(distance, service))
	
	shipment := Shipment{
		ID:            generateID("shp"),
		TrackingCode:  trackingCode,
		Origin:        origin,
		Destination:   destination,
		Status:        mt.StatusPending,
		EstimatedDate: estimatedDelivery,
		Carrier:       carrier,
		Service:       service,
		Weight:        weight,
		Cost:          cost,
		Items:         items,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Metadata:      make(map[string]string),
	}
	
	if errors := ValidateShipment(shipment); errors.HasErrors() {
		return nil, errors
	}
	
	ls.shipments = append(ls.shipments, shipment)
	return &shipment, nil
}

func (ls *LogisticsService) GetShipment(shipmentID string) (*Shipment, error) {
	for i, shipment := range ls.shipments {
		if shipment.ID == shipmentID {
			return &ls.shipments[i], nil
		}
	}
	return nil, errors.New("shipment not found")
}

func (ls *LogisticsService) GetShipmentByTracking(trackingCode string) (*Shipment, error) {
	for i, shipment := range ls.shipments {
		if shipment.TrackingCode == trackingCode {
			return &ls.shipments[i], nil
		}
	}
	return nil, errors.New("shipment not found")
}

func (ls *LogisticsService) UpdateShipmentStatus(shipmentID, status string) error {
	shipment, err := ls.GetShipment(shipmentID)
	if err != nil {
		return err
	}
	
	UpdateShipmentStatus(shipment, status)
	return nil
}

func (ls *LogisticsService) GetActiveShipments() []Shipment {
	var activeShipments []Shipment
	for _, shipment := range ls.shipments {
		if NewShipmentStatus(shipment.Status).IsActive() {
			activeShipments = append(activeShipments, shipment)
		}
	}
	return activeShipments
}

func (ls *LogisticsService) GetAllShipments() []Shipment {
	return ls.shipments
}

// Route Operations

func (ls *LogisticsService) CreateRoute(name string, origin, destination mt.Address,
	stops []RouteStop) (*Route, error) {
	
	distance := CalculateRouteDistance(stops)
	duration := time.Duration(distance * 6) * time.Minute // 6 minutes per mile
	cost := mt.NewMoney(distance * 0.50, mt.CurrencyUSD) // $0.50 per mile
	
	route := Route{
		ID:          generateID("rte"),
		Name:        name,
		Origin:      origin,
		Destination: destination,
		Distance:    distance,
		Duration:    duration,
		Cost:        cost,
		Stops:       stops,
		Status:      mt.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}
	
	if errors := ValidateRoute(route); errors.HasErrors() {
		return nil, errors
	}
	
	OptimizeRoute(&route)
	
	ls.routes = append(ls.routes, route)
	return &route, nil
}

func (ls *LogisticsService) GetRoute(routeID string) (*Route, error) {
	for i, route := range ls.routes {
		if route.ID == routeID {
			return &ls.routes[i], nil
		}
	}
	return nil, errors.New("route not found")
}

func (ls *LogisticsService) GetActiveRoutes() []Route {
	var activeRoutes []Route
	for _, route := range ls.routes {
		if NewRouteStatus(route.Status).IsActive() {
			activeRoutes = append(activeRoutes, route)
		}
	}
	return activeRoutes
}

func (ls *LogisticsService) GetAllRoutes() []Route {
	return ls.routes
}

// Vehicle Operations

func (ls *LogisticsService) CreateVehicle(name, vehicleType, licensePlate string,
	capacity VehicleCapacity) (*Vehicle, error) {
	
	vehicle := Vehicle{
		ID:           generateID("veh"),
		Name:         name,
		Type:         vehicleType,
		LicensePlate: licensePlate,
		Capacity:     capacity,
		Status:       mt.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     make(map[string]string),
	}
	
	if errors := ValidateVehicle(vehicle); errors.HasErrors() {
		return nil, errors
	}
	
	ls.vehicles = append(ls.vehicles, vehicle)
	return &vehicle, nil
}

func (ls *LogisticsService) GetVehicle(vehicleID string) (*Vehicle, error) {
	for i, vehicle := range ls.vehicles {
		if vehicle.ID == vehicleID {
			return &ls.vehicles[i], nil
		}
	}
	return nil, errors.New("vehicle not found")
}

func (ls *LogisticsService) GetAvailableVehicles() []Vehicle {
	var availableVehicles []Vehicle
	for _, vehicle := range ls.vehicles {
		if vehicle.Status == mt.StatusActive {
			availableVehicles = append(availableVehicles, vehicle)
		}
	}
	return availableVehicles
}

func (ls *LogisticsService) GetAllVehicles() []Vehicle {
	return ls.vehicles
}

// Driver Operations

func (ls *LogisticsService) CreateDriver(name, email, phone, licenseNum string) (*Driver, error) {
	driver := Driver{
		ID:         generateID("drv"),
		Name:       name,
		Email:      email,
		Phone:      phone,
		LicenseNum: licenseNum,
		Status:     mt.StatusActive,
		Rating:     5.0, // Start with perfect rating
		CreatedAt:  time.Now(),
		Metadata:   make(map[string]string),
	}
	
	var errors mt.ValidationErrors
	mt.ValidateRequired("name", driver.Name, "Driver Name", &errors)
	mt.ValidateEmail("email", driver.Email, "Email", &errors)
	mt.ValidateRequired("phone", driver.Phone, "Phone", &errors)
	mt.ValidateRequired("license_num", driver.LicenseNum, "License Number", &errors)
	
	if errors.HasErrors() {
		return nil, errors
	}
	
	ls.drivers = append(ls.drivers, driver)
	return &driver, nil
}

func (ls *LogisticsService) GetDriver(driverID string) (*Driver, error) {
	for i, driver := range ls.drivers {
		if driver.ID == driverID {
			return &ls.drivers[i], nil
		}
	}
	return nil, errors.New("driver not found")
}

func (ls *LogisticsService) GetAvailableDrivers() []Driver {
	var availableDrivers []Driver
	for _, driver := range ls.drivers {
		if driver.Status == mt.StatusActive {
			availableDrivers = append(availableDrivers, driver)
		}
	}
	return availableDrivers
}

func (ls *LogisticsService) GetAllDrivers() []Driver {
	return ls.drivers
}

// =====================================================
// DATA TRANSFER OBJECTS FOR PRESENTATION LAYER
// =====================================================

// ShipmentDisplayData prepares shipment data for UI display
type ShipmentDisplayData struct {
	Shipment        Shipment
	FormattedCost   string
	StatusClass     string
	StatusDisplay   string
	ProgressPercent int
	DaysInTransit   int
}

// RouteDisplayData prepares route data for UI display
type RouteDisplayData struct {
	Route             Route
	FormattedDistance string
	FormattedDuration string
	FormattedCost     string
	StatusClass       string
	StatusDisplay     string
	CompletionPercent int
}

// VehicleDisplayData prepares vehicle data for UI display
type VehicleDisplayData struct {
	Vehicle       Vehicle
	TypeIcon      string
	StatusClass   string
	StatusDisplay string
	CapacityUsed  string
}

// DashboardData aggregates logistics data for dashboard display
type DashboardData struct {
	TotalShipments    int
	ActiveShipments   int
	CompletedToday    int
	ActiveRoutes      int
	AvailableVehicles int
	ActiveDrivers     int
	RecentShipments   []ShipmentDisplayData
	Revenue           mt.Money
	FormattedRevenue  string
}

// =====================================================
// DATA PREPARATION FUNCTIONS (No UI Rendering)
// =====================================================

// PrepareShipmentForDisplay prepares shipment data for presentation layer
func PrepareShipmentForDisplay(shipment Shipment) ShipmentDisplayData {
	status := NewShipmentStatus(shipment.Status)
	daysInTransit := int(time.Since(shipment.CreatedAt).Hours() / 24)
	
	// Calculate progress percentage
	var progressPercent int
	switch shipment.Status {
	case mt.StatusPending:    progressPercent = 10
	case "picked_up":    progressPercent = 25
	case "in_transit":   progressPercent = 50
	case "out_for_delivery": progressPercent = 75
	case "delivered":    progressPercent = 100
	default:             progressPercent = 0
	}
	
	return ShipmentDisplayData{
		Shipment:        shipment,
		FormattedCost:   shipment.Cost.Format(),
		StatusClass:     "status-" + status.GetSeverity(),
		StatusDisplay:   status.GetDisplay(),
		ProgressPercent: progressPercent,
		DaysInTransit:   daysInTransit,
	}
}

// PrepareRouteForDisplay prepares route data for presentation layer
func PrepareRouteForDisplay(route Route) RouteDisplayData {
	status := NewRouteStatus(route.Status)
	
	// Calculate completion percentage
	var completionPercent int
	switch route.Status {
	case mt.StatusPending:   completionPercent = 0
	case "active":      completionPercent = 50
	case mt.StatusCompleted: completionPercent = 100
	default:            completionPercent = 0
	}
	
	return RouteDisplayData{
		Route:             route,
		FormattedDistance: fmt.Sprintf("%.1f miles", route.Distance),
		FormattedDuration: formatDuration(route.Duration),
		FormattedCost:     route.Cost.Format(),
		StatusClass:       "status-" + status.GetSeverity(),
		StatusDisplay:     status.GetDisplay(),
		CompletionPercent: completionPercent,
	}
}

// PrepareVehicleForDisplay prepares vehicle data for presentation layer
func PrepareVehicleForDisplay(vehicle Vehicle) VehicleDisplayData {
	return VehicleDisplayData{
		Vehicle:       vehicle,
		TypeIcon:      getVehicleTypeIcon(vehicle.Type),
		StatusClass:   "status-" + getVehicleStatusSeverity(vehicle.Status),
		StatusDisplay: getVehicleStatusDisplay(vehicle.Status),
		CapacityUsed:  fmt.Sprintf("%.1f / %.1f lbs", 0.0, vehicle.Capacity.Weight), // Would calculate actual usage
	}
}

// PrepareDashboardData aggregates logistics data for dashboard presentation
func PrepareDashboardData(ls *LogisticsService) DashboardData {
	allShipments := ls.GetAllShipments()
	activeShipments := ls.GetActiveShipments()
	availableVehicles := ls.GetAvailableVehicles()
	activeDrivers := ls.GetAvailableDrivers()
	activeRoutes := ls.GetActiveRoutes()
	
	// Calculate completed shipments today
	completedToday := 0
	today := time.Now().Truncate(24 * time.Hour)
	for _, shipment := range allShipments {
		if shipment.ActualDate != nil && shipment.ActualDate.After(today) {
			completedToday++
		}
	}
	
	// Calculate total revenue
	var revenue mt.Money
	for _, shipment := range allShipments {
		if shipment.Status == "delivered" {
			revenue.Amount += shipment.Cost.Amount
		}
	}
	
	// Prepare recent shipments for display
	var recentShipments []ShipmentDisplayData
	recentShipmentsList := getRecentShipments(allShipments, 5)
	for _, shipment := range recentShipmentsList {
		recentShipments = append(recentShipments, PrepareShipmentForDisplay(shipment))
	}
	
	return DashboardData{
		TotalShipments:    len(allShipments),
		ActiveShipments:   len(activeShipments),
		CompletedToday:    completedToday,
		ActiveRoutes:      len(activeRoutes),
		AvailableVehicles: len(availableVehicles),
		ActiveDrivers:     len(activeDrivers),
		RecentShipments:   recentShipments,
		Revenue:           revenue,
		FormattedRevenue:  revenue.Format(),
	}
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// generateID generates a unique ID with prefix
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// calculateDistance calculates distance between two addresses (simplified)
func calculateDistance(origin, destination mt.Address) float64 {
	// Simple implementation - in reality would use mapping service
	// Return a random distance for demo purposes
	return 50.0 + float64(len(origin.City)+len(destination.City)) // Rough approximation
}

// formatDuration formats duration for display
func formatDuration(duration time.Duration) string {
	hours := duration.Hours()
	if hours < 24 {
		return fmt.Sprintf("%.1f hours", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%.1f days", days)
}

// getVehicleTypeIcon returns icon for vehicle type
func getVehicleTypeIcon(vehicleType string) string {
	switch vehicleType {
	case "truck": return "🚛"
	case "van":   return "🚐"
	case "car":   return "🚗"
	case "bike":  return "🚲"
	case "motorcycle": return "🏍️"
	default:      return "🚐"
	}
}

// getVehicleStatusSeverity returns severity class for vehicle status
func getVehicleStatusSeverity(status string) string {
	switch status {
	case mt.StatusActive:    return "success"
	case mt.StatusInactive:  return "warning"
	case "maintenance": return "error"
	default:            return "info"
	}
}

// getVehicleStatusDisplay returns display text for vehicle status
func getVehicleStatusDisplay(status string) string {
	switch status {
	case mt.StatusActive:    return "Active"
	case mt.StatusInactive:  return "Inactive"
	case "maintenance": return "Maintenance"
	default:            return "Unknown"
	}
}

// getRecentShipments returns the most recent shipments
func getRecentShipments(shipments []Shipment, limit int) []Shipment {
	// Sort by created date descending
	sorted := make([]Shipment, len(shipments))
	copy(sorted, shipments)
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})
	
	if limit > len(sorted) {
		limit = len(sorted)
	}
	
	return sorted[:limit]
}

// =====================================================
// SAMPLE DATA HELPERS (for demos and testing)
// =====================================================

// SampleCustomer returns sample customer data
func SampleCustomer() Customer {
	return Customer{
		ID:    "cust_001",
		Name:  "ABC Logistics Co",
		Email: "shipping@abclogistics.com",
		Addresses: []mt.Address{
			{
				Type:       mt.AddressShipping,
				Name:       "ABC Logistics Co",
				Street1:    "789 Warehouse Blvd",
				City:       "Industrial City",
				State:      "TX",
				PostalCode: "75001",
				Country:    "US",
			},
		},
		AccountNumber:    "LOG001",
		PreferredCarrier: "FedEx",
		CreditLimit:      mt.NewMoney(50000.00, mt.CurrencyUSD),
		Status:           mt.StatusActive,
		CreatedAt:        time.Now().AddDate(-2, 0, 0),
		Metadata:         make(map[string]string),
	}
}

// SampleShipments returns sample shipment data for demos
func SampleShipments() []Shipment {
	origin := mt.Address{
		Type: mt.AddressPickup,
		Name: "Warehouse A", Street1: "123 Origin St",
		City: "Origin City", State: "CA", PostalCode: "90210", Country: "US",
	}
	destination := mt.Address{
		Type: mt.AddressDelivery,
		Name: "Customer B", Street1: "456 Destination Ave",
		City: "Destination City", State: "NY", PostalCode: "10001", Country: "US",
	}
	
	return []Shipment{
		{
			ID:            "shp_001",
			TrackingCode:  "TRK123456789",
			Origin:        origin,
			Destination:   destination,
			Status:        "in_transit",
			EstimatedDate: time.Now().AddDate(0, 0, 2),
			Carrier:       "FedEx",
			Service:       "express",
			Weight:        15.5,
			Cost:          mt.NewMoney(25.50, mt.CurrencyUSD),
			Items: []ShipmentItem{
				{
					ID: "item_001", Description: "Electronic Components",
					Quantity: 10, Weight: 15.5,
					Value: mt.NewMoney(500.00, mt.CurrencyUSD),
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -1),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
		{
			ID:            "shp_002",
			TrackingCode:  "TRK987654321",
			Origin:        origin,
			Destination:   destination,
			Status:        "delivered",
			EstimatedDate: time.Now().AddDate(0, 0, -1),
			Carrier:       "UPS",
			Service:       "standard",
			Weight:        8.2,
			Cost:          mt.NewMoney(15.75, mt.CurrencyUSD),
			Items: []ShipmentItem{
				{
					ID: "item_002", Description: "Office Supplies",
					Quantity: 5, Weight: 8.2,
					Value: mt.NewMoney(150.00, mt.CurrencyUSD),
				},
			},
			CreatedAt: time.Now().AddDate(0, 0, -3),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
	}
}

// SampleVehicles returns sample vehicle data
func SampleVehicles() []Vehicle {
	return []Vehicle{
		{
			ID:           "veh_001",
			Name:         "Delivery Truck 1",
			Type:         "truck",
			LicensePlate: "TRK-001",
			Capacity: VehicleCapacity{
				Weight:    5000.0,
				Volume:    500.0,
				ItemCount: 100,
			},
			Status: mt.StatusActive,
			Driver: Driver{
				ID: "drv_001", Name: "John Driver", Email: "john@example.com",
				Phone: "555-0101", Status: mt.StatusActive, Rating: 4.8,
			},
			CreatedAt: time.Now().AddDate(-1, 0, 0),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
		{
			ID:           "veh_002",
			Name:         "Delivery Van 1",
			Type:         "van",
			LicensePlate: "VAN-001",
			Capacity: VehicleCapacity{
				Weight:    2000.0,
				Volume:    200.0,
				ItemCount: 50,
			},
			Status: mt.StatusActive,
			Driver: Driver{
				ID: "drv_002", Name: "Jane Driver", Email: "jane@example.com",
				Phone: "555-0102", Status: mt.StatusActive, Rating: 4.9,
			},
			CreatedAt: time.Now().AddDate(-6, 0, 0),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]string),
		},
	}
}
