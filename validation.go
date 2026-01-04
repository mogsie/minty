package minty

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a form validation error.
type ValidationError struct {
	Field   string
	Message string
}

// ValidationResult represents the result of form validation.
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// AddError adds a validation error to the result.
func (vr *ValidationResult) AddError(field, message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// GetError returns the first error for a specific field, or empty string if none.
func (vr *ValidationResult) GetError(field string) string {
	for _, err := range vr.Errors {
		if err.Field == field {
			return err.Message
		}
	}
	return ""
}

// HasError checks if there's an error for a specific field.
func (vr *ValidationResult) HasError(field string) bool {
	return vr.GetError(field) != ""
}

// Basic validation functions

// ValidateRequired checks if a field is not empty.
func ValidateRequired(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	return nil
}

// ValidateEmail checks if a value is a valid email format.
func ValidateEmail(value, fieldName string) *ValidationError {
	if value == "" {
		return nil // Use ValidateRequired separately
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return &ValidationError{
			Field:   fieldName,
			Message: "Please enter a valid email address",
		}
	}
	return nil
}

// ValidateMinLength checks minimum length.
func ValidateMinLength(value, fieldName string, minLen int) *ValidationError {
	if value != "" && len(value) < minLen {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, minLen),
		}
	}
	return nil
}

// ValidateMaxLength checks maximum length.
func ValidateMaxLength(value, fieldName string, maxLen int) *ValidationError {
	if len(value) > maxLen {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be no more than %d characters long", fieldName, maxLen),
		}
	}
	return nil
}

// Form field helpers with validation support

// FormField creates a form field with label, input, and error display.
func FormField(label, name, fieldType, value, errorMsg string, required bool, attributes ...Attribute) H {
	return func(b *Builder) Node {
		attrs := []Attribute{Type(fieldType), Name(name)}
		if value != "" {
			attrs = append(attrs, Value(value))
		}
		if required {
			attrs = append(attrs, Required())
		}
		attrs = append(attrs, attributes...)
		
		var errorElement Node = NewFragment()
		if errorMsg != "" {
			errorElement = ErrorMessage(errorMsg)(b)
		}
		
		return b.Div(Class("form-field"),
			b.Label(For(name), label),
			b.Input(attrs...),
			errorElement,
		)
	}
}

// TextareaField creates a textarea field with label and error display.
func TextareaField(label, name, value, errorMsg string, required bool, rows, cols int) H {
	return func(b *Builder) Node {
		attrs := []interface{}{Name(name)}
		if required {
			attrs = append(attrs, Required())
		}
		if rows > 0 {
			attrs = append(attrs, Rows(rows))
		}
		if cols > 0 {
			attrs = append(attrs, Cols(cols))
		}
		attrs = append(attrs, value) // Text content
		
		var errorElement Node = NewFragment()
		if errorMsg != "" {
			errorElement = ErrorMessage(errorMsg)(b)
		}
		
		return b.Div(Class("form-field"),
			b.Label(For(name), label),
			b.Textarea(attrs...),
			errorElement,
		)
	}
}

// SelectField creates a select field with options.
func SelectField(label, name, value, errorMsg string, required bool, options []SelectOption) H {
	return func(b *Builder) Node {
		selectAttrs := []interface{}{Name(name)}
		if required {
			selectAttrs = append(selectAttrs, Required())
		}
		
		for _, opt := range options {
			optAttrs := []interface{}{Value(opt.Value), opt.Text}
			if opt.Value == value {
				optAttrs = append(optAttrs, Selected())
			}
			selectAttrs = append(selectAttrs, b.Option(optAttrs...))
		}
		
		var errorElement Node = NewFragment()
		if errorMsg != "" {
			errorElement = ErrorMessage(errorMsg)(b)
		}
		
		return b.Div(Class("form-field"),
			b.Label(For(name), label),
			b.Select(selectAttrs...),
			errorElement,
		)
	}
}

// SelectOption represents an option in a select field.
type SelectOption struct {
	Value string
	Text  string
}



// SanitizeInput performs basic input sanitization.
func SanitizeInput(input string) string {
	// Trim whitespace
	sanitized := strings.TrimSpace(input)
	
	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	
	return sanitized
}
