package models

import "time"

// =============================================================================
// COVERAGE TYPES
// =============================================================================

// CoverageType represents a type of insurance coverage.
type CoverageType string

const (
	CoverageAuto     CoverageType = "auto"
	CoverageHome     CoverageType = "home"
	CoverageLife     CoverageType = "life"
	CoverageBusiness CoverageType = "business"
)

// Coverage represents an insurance coverage option.
type Coverage struct {
	ID          string       `json:"id"`
	Type        CoverageType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	BasePrice   float64      `json:"basePrice"`
	Icon        string       `json:"icon"` // Icon name for UI
}

// =============================================================================
// PLANS
// =============================================================================

// Plan represents an insurance plan with specific coverage levels.
type Plan struct {
	ID           string   `json:"id"`
	CoverageType string   `json:"coverageType"`
	Name         string   `json:"name"`
	Tier         string   `json:"tier"` // basic, standard, premium
	Price        float64  `json:"price"`
	Deductible   float64  `json:"deductible"`
	Coverage     float64  `json:"coverage"` // Coverage amount
	Features     []string `json:"features"`
	Popular      bool     `json:"popular"`
}

// =============================================================================
// QUOTES
// =============================================================================

// QuoteStatus represents the status of a quote.
type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "draft"
	QuoteStatusPending   QuoteStatus = "pending"
	QuoteStatusApproved  QuoteStatus = "approved"
	QuoteStatusDeclined  QuoteStatus = "declined"
	QuoteStatusExpired   QuoteStatus = "expired"
)

// Quote represents an insurance quote request.
type Quote struct {
	ID           string      `json:"id"`
	CustomerName string      `json:"customerName"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	CoverageType string      `json:"coverageType"`
	PlanID       string      `json:"planId"`
	Status       QuoteStatus `json:"status"`
	Premium      float64     `json:"premium"`
	CreatedAt    time.Time   `json:"createdAt"`
	ExpiresAt    time.Time   `json:"expiresAt"`
}

// =============================================================================
// CUSTOMER PROFILE (for form dependencies demo)
// =============================================================================

// CustomerProfile holds customer information with conditional fields.
type CustomerProfile struct {
	// Basic info
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"dateOfBirth"`
	
	// Address
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	
	// Coverage-specific fields (shown/hidden based on coverage type)
	
	// Auto insurance fields
	VehicleMake    string `json:"vehicleMake,omitempty"`
	VehicleModel   string `json:"vehicleModel,omitempty"`
	VehicleYear    int    `json:"vehicleYear,omitempty"`
	VIN            string `json:"vin,omitempty"`
	DrivingYears   int    `json:"drivingYears,omitempty"`
	HasAccidents   bool   `json:"hasAccidents,omitempty"`
	AccidentCount  int    `json:"accidentCount,omitempty"`
	
	// Home insurance fields
	PropertyType   string  `json:"propertyType,omitempty"`   // house, condo, apartment
	YearBuilt      int     `json:"yearBuilt,omitempty"`
	SquareFeet     int     `json:"squareFeet,omitempty"`
	PropertyValue  float64 `json:"propertyValue,omitempty"`
	HasPool        bool    `json:"hasPool,omitempty"`
	HasAlarm       bool    `json:"hasAlarm,omitempty"`
	
	// Life insurance fields
	HealthStatus   string  `json:"healthStatus,omitempty"`   // excellent, good, fair, poor
	Smoker         bool    `json:"smoker,omitempty"`
	CoverageAmount float64 `json:"coverageAmount,omitempty"`
	Beneficiaries  int     `json:"beneficiaries,omitempty"`
	
	// Business insurance fields
	BusinessName   string  `json:"businessName,omitempty"`
	BusinessType   string  `json:"businessType,omitempty"`
	Employees      int     `json:"employees,omitempty"`
	AnnualRevenue  float64 `json:"annualRevenue,omitempty"`
	HasPremises    bool    `json:"hasPremises,omitempty"`
}

// =============================================================================
// CLAIM (for data filtering demo)
// =============================================================================

// ClaimStatus represents the status of a claim.
type ClaimStatus string

const (
	ClaimStatusOpen       ClaimStatus = "open"
	ClaimStatusInProgress ClaimStatus = "in-progress"
	ClaimStatusApproved   ClaimStatus = "approved"
	ClaimStatusDenied     ClaimStatus = "denied"
	ClaimStatusClosed     ClaimStatus = "closed"
)

// Claim represents an insurance claim.
type Claim struct {
	ID           string      `json:"id"`
	PolicyNumber string      `json:"policyNumber"`
	CustomerName string      `json:"customerName"`
	Type         string      `json:"type"` // collision, theft, fire, medical, etc.
	Status       ClaimStatus `json:"status"`
	Amount       float64     `json:"amount"`
	Filed        string      `json:"filed"`
	Description  string      `json:"description"`
}
