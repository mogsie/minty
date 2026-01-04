// Package mintytypes provides pure business types, interfaces, and utilities.
// This package has ZERO dependencies on the minty HTML framework and serves
// as the clean foundation for domain packages.
package mintytypes

import (
	"fmt"
	"strings"
	"time"
)

// =====================================================
// MONETARY TYPES
// =====================================================

// Money represents monetary values with precision.
// Uses int64 cents for exact arithmetic, no floating point errors.
type Money struct {
	Amount   int64  `json:"amount"`   // Amount in smallest currency unit (cents for USD)
	Currency string `json:"currency"` // ISO 4217 currency code
}

// MajorUnit returns the major currency unit as float64.
func (m Money) MajorUnit() float64 {
	return float64(m.Amount) / 100.0
}

// Format returns formatted money string based on currency.
func (m Money) Format() string {
	switch strings.ToUpper(m.Currency) {
	case "USD":
		return fmt.Sprintf("$%.2f", m.MajorUnit())
	case "EUR":
		return fmt.Sprintf("€%.2f", m.MajorUnit())
	case "GBP":
		return fmt.Sprintf("£%.2f", m.MajorUnit())
	case "JPY":
		return fmt.Sprintf("¥%.0f", m.MajorUnit()*100) // JPY doesn't use cents
	case "CAD":
		return fmt.Sprintf("CA$%.2f", m.MajorUnit())
	case "AUD":
		return fmt.Sprintf("AU$%.2f", m.MajorUnit())
	default:
		return fmt.Sprintf("%.2f %s", m.MajorUnit(), m.Currency)
	}
}

// Add adds another Money value (must be same currency).
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Subtract subtracts another Money value (must be same currency).
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies: %s and %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}, nil
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive returns true if the amount is positive.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// IsNegative returns true if the amount is negative.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// NewMoney creates a new Money value from a major unit amount.
func NewMoney(majorUnit float64, currency string) Money {
	return Money{
		Amount:   int64(majorUnit * 100), // Convert to cents
		Currency: strings.ToUpper(currency),
	}
}

// =====================================================
// ADDRESS TYPES
// =====================================================

// Address represents a physical or mailing address.
type Address struct {
	Type       string `json:"type"`        // "billing", "shipping", "pickup", "delivery"
	Name       string `json:"name"`        // Recipient name
	Company    string `json:"company"`     // Company name (optional)
	Street1    string `json:"street1"`     // Primary address line
	Street2    string `json:"street2"`     // Secondary address line (optional)
	City       string `json:"city"`        // City name
	State      string `json:"state"`       // State/Province
	PostalCode string `json:"postal_code"` // ZIP/Postal code
	Country    string `json:"country"`     // Country code (ISO 3166-1)
}

// FormatOneLine returns address as single line string.
func (a Address) FormatOneLine() string {
	parts := []string{}
	if a.Street1 != "" {
		parts = append(parts, a.Street1)
	}
	if a.Street2 != "" {
		parts = append(parts, a.Street2)
	}
	if a.City != "" {
		parts = append(parts, a.City)
	}
	if a.State != "" {
		parts = append(parts, a.State)
	}
	if a.PostalCode != "" {
		parts = append(parts, a.PostalCode)
	}
	if a.Country != "" {
		parts = append(parts, a.Country)
	}
	return strings.Join(parts, ", ")
}

// FormatMultiLine returns address as multi-line string.
func (a Address) FormatMultiLine() string {
	lines := []string{}
	if a.Name != "" {
		lines = append(lines, a.Name)
	}
	if a.Company != "" {
		lines = append(lines, a.Company)
	}
	if a.Street1 != "" {
		lines = append(lines, a.Street1)
	}
	if a.Street2 != "" {
		lines = append(lines, a.Street2)
	}
	cityLine := ""
	if a.City != "" {
		cityLine = a.City
	}
	if a.State != "" {
		if cityLine != "" {
			cityLine += ", "
		}
		cityLine += a.State
	}
	if a.PostalCode != "" {
		if cityLine != "" {
			cityLine += " "
		}
		cityLine += a.PostalCode
	}
	if cityLine != "" {
		lines = append(lines, cityLine)
	}
	if a.Country != "" {
		lines = append(lines, a.Country)
	}
	return strings.Join(lines, "\n")
}

// =====================================================
// CUSTOMER INTERFACE
// =====================================================

// Customer interface for entities that can be customers.
type Customer interface {
	GetID() string
	GetName() string
	GetEmail() string
	GetAddresses() []Address
	GetPrimaryAddress() Address
	GetBillingAddress() Address
	GetShippingAddress() Address
}

// =====================================================
// STATUS INTERFACE
// =====================================================

// Status interface for status values across domains.
type Status interface {
	GetCode() string
	GetDisplay() string
	IsActive() bool
	GetSeverity() string // "success", "warning", "danger", "info"
	GetDescription() string
}

// BaseStatus provides a default Status implementation.
type BaseStatus struct {
	Code        string `json:"code"`
	Display     string `json:"display"`
	Active      bool   `json:"active"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

func (s BaseStatus) GetCode() string        { return s.Code }
func (s BaseStatus) GetDisplay() string     { return s.Display }
func (s BaseStatus) IsActive() bool         { return s.Active }
func (s BaseStatus) GetSeverity() string    { return s.Severity }
func (s BaseStatus) GetDescription() string { return s.Description }

// =====================================================
// VALIDATION TYPES
// =====================================================

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Add adds a validation error.
func (v *ValidationErrors) Add(field, message string) {
	*v = append(*v, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are validation errors.
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// GetFieldErrors returns all errors for a specific field.
func (v ValidationErrors) GetFieldErrors(field string) []string {
	var errors []string
	for _, err := range v {
		if err.Field == field {
			errors = append(errors, err.Message)
		}
	}
	return errors
}

// Error implements the error interface.
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var messages []string
	for _, err := range v {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// =====================================================
// VALIDATION FUNCTIONS
// =====================================================

// ValidateRequired validates that a value is not empty.
func ValidateRequired(field, value, fieldName string, errors *ValidationErrors) {
	if strings.TrimSpace(value) == "" {
		errors.Add(field, fmt.Sprintf("%s is required", fieldName))
	}
}

// ValidateEmail validates email format (basic validation).
func ValidateEmail(field, email, fieldName string, errors *ValidationErrors) {
	email = strings.TrimSpace(email)
	if email == "" {
		return // Use ValidateRequired for empty check
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		errors.Add(field, fmt.Sprintf("%s must be a valid email address", fieldName))
	}
}

// ValidateMoneyAmount validates money amount is positive.
func ValidateMoneyAmount(field string, money Money, fieldName string, errors *ValidationErrors) {
	if money.Amount <= 0 {
		errors.Add(field, fmt.Sprintf("%s must be greater than zero", fieldName))
	}
}

// =====================================================
// DATE/TIME UTILITIES
// =====================================================

// FormatDate formats a date string for display.
func FormatDate(date string) string {
	if t, err := time.Parse("2006-01-02", date); err == nil {
		return t.Format("January 2, 2006")
	}
	return date // Return original if parsing fails
}

// DaysAgo calculates how many days ago a date was.
func DaysAgo(date string) int {
	if t, err := time.Parse("2006-01-02", date); err == nil {
		return int(time.Since(t).Hours() / 24)
	}
	return 0
}

// =====================================================
// CONSTANTS
// =====================================================

// Common status codes used across domains.
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
	StatusFailed    = "failed"
	StatusDraft     = "draft"
	StatusPublished = "published"
)

// Common address types.
const (
	AddressBilling  = "billing"
	AddressShipping = "shipping"
	AddressPickup   = "pickup"
	AddressDelivery = "delivery"
	AddressOffice   = "office"
	AddressHome     = "home"
)

// Common currencies.
const (
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	CurrencyGBP = "GBP"
	CurrencyJPY = "JPY"
	CurrencyCAD = "CAD"
	CurrencyAUD = "AUD"
)
