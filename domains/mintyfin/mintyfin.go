// Package mintyfin provides pure finance domain logic for the Minty System.
// This package contains NO UI dependencies and focuses solely on business logic.
package mintyfin

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	mt "github.com/ha1tch/minty/mintytypes"
)

// =====================================================
// PURE BUSINESS TYPES (No UI Dependencies)
// =====================================================

// Account represents a financial account
type Account struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Balance     mt.Money    `json:"balance"`
	Status      string           `json:"status"`
	Type        string           `json:"type"` // checking, savings, investment, credit
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Description string           `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          string        `json:"id"`
	AccountID   string        `json:"account_id"`
	Amount      mt.Money `json:"amount"`
	Description string        `json:"description"`
	Date        time.Time     `json:"date"`
	Status      string        `json:"status"`
	Type        string        `json:"type"` // debit, credit
	Category    string        `json:"category"`
	Reference   string        `json:"reference"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID          string           `json:"id"`
	Number      string           `json:"number"`
	Amount      mt.Money    `json:"amount"`
	DueDate     time.Time        `json:"due_date"`
	Status      string           `json:"status"`
	Customer    Customer         `json:"customer"`
	Items       []InvoiceItem    `json:"items"`
	CreatedAt   time.Time        `json:"created_at"`
	PaidAt      *time.Time       `json:"paid_at,omitempty"`
	Description string           `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Quantity    int           `json:"quantity"`
	UnitPrice   mt.Money `json:"unit_price"`
	Total       mt.Money `json:"total"`
	Category    string        `json:"category"`
}

// Customer represents a finance customer
type Customer struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Addresses      []mt.Address  `json:"addresses"`
	AccountNumber  string             `json:"account_number"`
	CreditRating   string             `json:"credit_rating"`
	PaymentTerms   string             `json:"payment_terms"`
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

// Portfolio represents an investment portfolio
type Portfolio struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	TotalValue  mt.Money    `json:"total_value"`
	Performance float64          `json:"performance"` // percentage
	Positions   []Position       `json:"positions"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Status      string           `json:"status"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Position represents a position in a portfolio
type Position struct {
	ID        string        `json:"id"`
	Symbol    string        `json:"symbol"`
	Name      string        `json:"name"`
	Quantity  int           `json:"quantity"`
	Price     mt.Money `json:"price"`
	Value     mt.Money `json:"value"`
	Change    float64       `json:"change"` // percentage
	UpdatedAt time.Time     `json:"updated_at"`
}

// =====================================================
// STATUS IMPLEMENTATIONS
// =====================================================

// AccountStatus implements mt.Status interface
type AccountStatus struct {
	status string
}

func NewAccountStatus(status string) AccountStatus {
	return AccountStatus{status: status}
}

func (s AccountStatus) GetCode() string { return s.status }

func (s AccountStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusActive:    return "Active"
	case mt.StatusInactive:  return "Inactive"
	case "suspended": return "Suspended"
	case "closed":    return "Closed"
	default:          return "Unknown"
	}
}

func (s AccountStatus) IsActive() bool {
	return s.status == mt.StatusActive
}

func (s AccountStatus) GetSeverity() string {
	switch s.status {
	case mt.StatusActive:    return "success"
	case mt.StatusInactive:  return "warning"
	case "suspended": return "error"
	case "closed":    return "secondary"
	default:          return "info"
	}
}

func (s AccountStatus) GetDescription() string {
	switch s.status {
	case mt.StatusActive:    return "Account is active and operational"
	case mt.StatusInactive:  return "Account is temporarily inactive"
	case "suspended": return "Account has been suspended due to issues"
	case "closed":    return "Account is permanently closed"
	default:          return ""
	}
}

// TransactionStatus implements mt.Status interface
type TransactionStatus struct {
	status string
}

func NewTransactionStatus(status string) TransactionStatus {
	return TransactionStatus{status: status}
}

func (s TransactionStatus) GetCode() string { return s.status }

func (s TransactionStatus) GetDisplay() string {
	switch s.status {
	case mt.StatusPending:   return "Pending"
	case mt.StatusCompleted: return "Completed"
	case mt.StatusFailed:    return "Failed"
	case mt.StatusCancelled: return "Cancelled"
	default:                      return "Unknown"
	}
}

func (s TransactionStatus) IsActive() bool {
	return s.status == mt.StatusPending || s.status == mt.StatusCompleted
}

func (s TransactionStatus) GetSeverity() string {
	switch s.status {
	case mt.StatusCompleted: return "success"
	case mt.StatusPending:   return "warning"
	case mt.StatusFailed:    return "error"
	case mt.StatusCancelled: return "secondary"
	default:                      return "info"
	}
}

func (s TransactionStatus) GetDescription() string {
	switch s.status {
	case mt.StatusPending:   return "Transaction is being processed"
	case mt.StatusCompleted: return "Transaction completed successfully"
	case mt.StatusFailed:    return "Transaction failed to process"
	case mt.StatusCancelled: return "Transaction was cancelled"
	default:                      return ""
	}
}

// =====================================================
// PURE BUSINESS LOGIC FUNCTIONS
// =====================================================

// Account Business Logic

// CalculateAccountBalance calculates account balance from transactions
func CalculateAccountBalance(transactions []Transaction) mt.Money {
	var balance mt.Money
	for _, txn := range transactions {
		if txn.Type == "credit" {
			balance.Amount += txn.Amount.Amount
		} else if txn.Type == "debit" {
			balance.Amount -= txn.Amount.Amount
		}
	}
	return balance
}

// ValidateAccount validates account data
func ValidateAccount(account Account) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("name", account.Name, "Account Name", &errors)
	mt.ValidateRequired("type", account.Type, "Account Type", &errors)
	
	if account.Type != "" {
		validTypes := []string{"checking", "savings", "investment", "credit"}
		isValid := false
		for _, validType := range validTypes {
			if account.Type == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			errors.Add("type", "Account type must be one of: checking, savings, investment, credit")
		}
	}
	
	if account.Balance.Amount < 0 && account.Type != "credit" {
		errors.Add("balance", "Account balance cannot be negative for this account type")
	}
	
	return errors
}

// ProcessAccountTransaction processes a transaction for an account
func ProcessAccountTransaction(account *Account, transaction Transaction) error {
	if account.ID != transaction.AccountID {
		return errors.New("transaction account ID does not match account")
	}
	
	if transaction.Type == "debit" {
		if account.Balance.Amount < transaction.Amount.Amount && account.Type != "credit" {
			return errors.New("insufficient funds for debit transaction")
		}
		account.Balance.Amount -= transaction.Amount.Amount
	} else if transaction.Type == "credit" {
		account.Balance.Amount += transaction.Amount.Amount
	} else {
		return errors.New("invalid transaction type")
	}
	
	account.UpdatedAt = time.Now()
	return nil
}

// Transaction Business Logic

// ValidateTransaction validates transaction data
func ValidateTransaction(transaction Transaction) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("account_id", transaction.AccountID, "Account ID", &errors)
	mt.ValidateRequired("description", transaction.Description, "Description", &errors)
	mt.ValidateRequired("type", transaction.Type, "Transaction Type", &errors)
	mt.ValidateMoneyAmount("amount", transaction.Amount, "Amount", &errors)
	
	if transaction.Type != "" && transaction.Type != "debit" && transaction.Type != "credit" {
		errors.Add("type", "Transaction type must be either 'debit' or 'credit'")
	}
	
	return errors
}

// CategorizeTransaction automatically categorizes a transaction based on description
func CategorizeTransaction(transaction *Transaction) {
	description := strings.ToLower(transaction.Description)
	
	switch {
	case strings.Contains(description, "grocery") || strings.Contains(description, "food"):
		transaction.Category = "food"
	case strings.Contains(description, "gas") || strings.Contains(description, "fuel"):
		transaction.Category = "transportation"
	case strings.Contains(description, "salary") || strings.Contains(description, "payroll"):
		transaction.Category = "income"
	case strings.Contains(description, "rent") || strings.Contains(description, "mortgage"):
		transaction.Category = "housing"
	case strings.Contains(description, "utility") || strings.Contains(description, "electric") || strings.Contains(description, "water"):
		transaction.Category = "utilities"
	default:
		transaction.Category = "other"
	}
}

// Invoice Business Logic

// ValidateInvoice validates invoice data
func ValidateInvoice(invoice Invoice) mt.ValidationErrors {
	var errors mt.ValidationErrors
	
	mt.ValidateRequired("number", invoice.Number, "Invoice Number", &errors)
	mt.ValidateRequired("customer.name", invoice.Customer.Name, "Customer Name", &errors)
	mt.ValidateMoneyAmount("amount", invoice.Amount, "Amount", &errors)
	
	if invoice.DueDate.IsZero() {
		errors.Add("due_date", "Due date is required")
	} else if invoice.DueDate.Before(time.Now().AddDate(0, 0, -1)) {
		errors.Add("due_date", "Due date cannot be in the past")
	}
	
	if len(invoice.Items) == 0 {
		errors.Add("items", "Invoice must have at least one item")
	}
	
	// Validate that invoice amount matches sum of items
	var itemsTotal mt.Money
	for _, item := range invoice.Items {
		itemsTotal.Amount += item.Total.Amount
	}
	
	if itemsTotal.Amount != invoice.Amount.Amount {
		errors.Add("amount", "Invoice amount must match sum of item totals")
	}
	
	return errors
}

// ProcessPayment processes a payment for an invoice
func ProcessPayment(invoice *Invoice, paymentAmount mt.Money) error {
	if invoice.Status == "paid" {
		return errors.New("invoice is already paid")
	}
	
	if paymentAmount.Currency != invoice.Amount.Currency {
		return fmt.Errorf("payment currency %s does not match invoice currency %s", 
			paymentAmount.Currency, invoice.Amount.Currency)
	}
	
	if paymentAmount.Amount != invoice.Amount.Amount {
		return errors.New("payment amount must match invoice amount")
	}
	
	invoice.Status = "paid"
	now := time.Now()
	invoice.PaidAt = &now
	
	return nil
}

// CalculateInvoiceTotal calculates total from invoice items
func CalculateInvoiceTotal(items []InvoiceItem) mt.Money {
	var total mt.Money
	for _, item := range items {
		total.Amount += item.Total.Amount
	}
	return total
}

// =====================================================
// DOMAIN SERVICES
// =====================================================

// FinanceService provides business operations for the finance domain
type FinanceService struct {
	accounts     []Account
	transactions []Transaction
	invoices     []Invoice
	customers    []Customer
}

// NewFinanceService creates a new finance service
func NewFinanceService() *FinanceService {
	return &FinanceService{
		accounts:     make([]Account, 0),
		transactions: make([]Transaction, 0),
		invoices:     make([]Invoice, 0),
		customers:    make([]Customer, 0),
	}
}

// Account Operations

func (fs *FinanceService) CreateAccount(name, accountType string, initialBalance mt.Money, customerID string) (*Account, error) {
	account := Account{
		ID:        generateID("acc"),
		Name:      name,
		Balance:   initialBalance,
		Status:    mt.StatusActive,
		Type:      accountType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  map[string]string{"customer_id": customerID},
	}
	
	if errors := ValidateAccount(account); errors.HasErrors() {
		return nil, errors
	}
	
	fs.accounts = append(fs.accounts, account)
	return &account, nil
}

func (fs *FinanceService) GetAccount(accountID string) (*Account, error) {
	for i, account := range fs.accounts {
		if account.ID == accountID {
			return &fs.accounts[i], nil
		}
	}
	return nil, errors.New("account not found")
}

func (fs *FinanceService) GetAccountsByCustomer(customerID string) []Account {
	var customerAccounts []Account
	for _, account := range fs.accounts {
		if account.Metadata["customer_id"] == customerID {
			customerAccounts = append(customerAccounts, account)
		}
	}
	return customerAccounts
}

func (fs *FinanceService) UpdateAccountBalance(accountID string, transactions []Transaction) error {
	account, err := fs.GetAccount(accountID)
	if err != nil {
		return err
	}
	
	account.Balance = CalculateAccountBalance(transactions)
	account.UpdatedAt = time.Now()
	return nil
}

func (fs *FinanceService) GetTotalBalance() mt.Money {
	var total mt.Money
	for _, account := range fs.accounts {
		if account.Status == mt.StatusActive {
			total.Amount += account.Balance.Amount
		}
	}
	return total
}

func (fs *FinanceService) GetAllAccounts() []Account {
	return fs.accounts
}

// Transaction Operations

func (fs *FinanceService) CreateTransaction(accountID string, amount mt.Money, 
	description, txnType string) (*Transaction, error) {
	
	transaction := Transaction{
		ID:          generateID("txn"),
		AccountID:   accountID,
		Amount:      amount,
		Description: description,
		Date:        time.Now(),
		Status:      mt.StatusCompleted,
		Type:        txnType,
		Metadata:    make(map[string]string),
	}
	
	// Auto-categorize transaction
	CategorizeTransaction(&transaction)
	
	// Validate transaction
	if errors := ValidateTransaction(transaction); errors.HasErrors() {
		return nil, errors
	}
	
	// Update account balance
	account, err := fs.GetAccount(accountID)
	if err != nil {
		return nil, err
	}
	
	if err := ProcessAccountTransaction(account, transaction); err != nil {
		return nil, err
	}
	
	fs.transactions = append(fs.transactions, transaction)
	return &transaction, nil
}

func (fs *FinanceService) GetTransactionsByAccount(accountID string) []Transaction {
	var accountTransactions []Transaction
	for _, txn := range fs.transactions {
		if txn.AccountID == accountID {
			accountTransactions = append(accountTransactions, txn)
		}
	}
	
	// Sort by date descending
	sort.Slice(accountTransactions, func(i, j int) bool {
		return accountTransactions[i].Date.After(accountTransactions[j].Date)
	})
	
	return accountTransactions
}

func (fs *FinanceService) GetRecentTransactions(limit int) []Transaction {
	// Sort all transactions by date
	allTxns := make([]Transaction, len(fs.transactions))
	copy(allTxns, fs.transactions)
	
	sort.Slice(allTxns, func(i, j int) bool {
		return allTxns[i].Date.After(allTxns[j].Date)
	})
	
	if limit > len(allTxns) {
		limit = len(allTxns)
	}
	
	return allTxns[:limit]
}

func (fs *FinanceService) GetAllTransactions() []Transaction {
	return fs.transactions
}

// Invoice Operations

func (fs *FinanceService) CreateInvoice(number string, customer Customer, 
	items []InvoiceItem, dueDate time.Time) (*Invoice, error) {
	
	invoice := Invoice{
		ID:        generateID("inv"),
		Number:    number,
		Amount:    CalculateInvoiceTotal(items),
		DueDate:   dueDate,
		Status:    mt.StatusPending,
		Customer:  customer,
		Items:     items,
		CreatedAt: time.Now(),
		Metadata:  make(map[string]string),
	}
	
	if errors := ValidateInvoice(invoice); errors.HasErrors() {
		return nil, errors
	}
	
	fs.invoices = append(fs.invoices, invoice)
	return &invoice, nil
}

func (fs *FinanceService) GetPendingInvoices() []Invoice {
	var pendingInvoices []Invoice
	for _, invoice := range fs.invoices {
		if invoice.Status == mt.StatusPending {
			pendingInvoices = append(pendingInvoices, invoice)
		}
	}
	return pendingInvoices
}

func (fs *FinanceService) GetAllInvoices() []Invoice {
	return fs.invoices
}

func (fs *FinanceService) PayInvoice(invoiceID string, paymentAmount mt.Money) error {
	for i, invoice := range fs.invoices {
		if invoice.ID == invoiceID {
			if err := ProcessPayment(&fs.invoices[i], paymentAmount); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("invoice not found")
}

// =====================================================
// DATA TRANSFER OBJECTS FOR PRESENTATION LAYER
// =====================================================

// AccountDisplayData prepares account data for UI display
type AccountDisplayData struct {
	Account          Account
	FormattedBalance string
	StatusClass      string
	StatusDisplay    string
	TypeIcon         string
	TypeDisplay      string
	RecentTxnCount   int
}

// TransactionDisplayData prepares transaction data for UI display
type TransactionDisplayData struct {
	Transaction      Transaction
	FormattedAmount  string
	FormattedDate    string
	StatusClass      string
	StatusDisplay    string
	TypeClass        string
	CategoryIcon     string
	DaysAgo          int
}

// InvoiceDisplayData prepares invoice data for UI display
type InvoiceDisplayData struct {
	Invoice         Invoice
	FormattedAmount string
	FormattedDueDate string
	StatusClass     string
	StatusDisplay   string
	IsOverdue       bool
	DaysUntilDue    int
}

// DashboardData aggregates data for dashboard display
type DashboardData struct {
	TotalBalance       mt.Money
	FormattedTotal     string
	AccountCount       int
	ActiveAccountCount int
	TransactionCount   int
	PendingInvoices    int
	TopAccounts        []AccountDisplayData
	RecentTransactions []TransactionDisplayData
	MonthlySpending    mt.Money
	MonthlyIncome      mt.Money
}

// =====================================================
// DATA PREPARATION FUNCTIONS (No UI Rendering)
// =====================================================

// PrepareAccountForDisplay prepares account data for presentation layer
func PrepareAccountForDisplay(account Account) AccountDisplayData {
	status := NewAccountStatus(account.Status)
	
	return AccountDisplayData{
		Account:          account,
		FormattedBalance: account.Balance.Format(),
		StatusClass:      "status-" + status.GetSeverity(),
		StatusDisplay:    status.GetDisplay(),
		TypeIcon:         getAccountTypeIcon(account.Type),
		TypeDisplay:      getAccountTypeDisplay(account.Type),
		RecentTxnCount:   0, // Would be calculated if transactions were provided
	}
}

// PrepareTransactionForDisplay prepares transaction data for presentation layer
func PrepareTransactionForDisplay(transaction Transaction) TransactionDisplayData {
	status := NewTransactionStatus(transaction.Status)
	
	return TransactionDisplayData{
		Transaction:     transaction,
		FormattedAmount: formatTransactionAmount(transaction),
		FormattedDate:   mt.FormatDate(transaction.Date.Format("2006-01-02")),
		StatusClass:     "status-" + status.GetSeverity(),
		StatusDisplay:   status.GetDisplay(),
		TypeClass:       "transaction-" + transaction.Type,
		CategoryIcon:    getCategoryIcon(transaction.Category),
		DaysAgo:         mt.DaysAgo(transaction.Date.Format("2006-01-02")),
	}
}

// PrepareInvoiceForDisplay prepares invoice data for presentation layer
func PrepareInvoiceForDisplay(invoice Invoice) InvoiceDisplayData {
	isOverdue := time.Now().After(invoice.DueDate) && invoice.Status != "paid"
	daysUntilDue := int(time.Until(invoice.DueDate).Hours() / 24)
	
	return InvoiceDisplayData{
		Invoice:          invoice,
		FormattedAmount:  invoice.Amount.Format(),
		FormattedDueDate: mt.FormatDate(invoice.DueDate.Format("2006-01-02")),
		StatusClass:      getInvoiceStatusClass(invoice.Status, isOverdue),
		StatusDisplay:    getInvoiceStatusDisplay(invoice.Status),
		IsOverdue:        isOverdue,
		DaysUntilDue:     daysUntilDue,
	}
}

// PrepareDashboardData aggregates data for dashboard presentation
func PrepareDashboardData(fs *FinanceService) DashboardData {
	totalBalance := fs.GetTotalBalance()
	activeAccounts := 0
	
	var topAccounts []AccountDisplayData
	for i, account := range fs.accounts {
		if account.Status == mt.StatusActive {
			activeAccounts++
		}
		if i < 3 { // Top 3 accounts
			topAccounts = append(topAccounts, PrepareAccountForDisplay(account))
		}
	}
	
	recentTxns := fs.GetRecentTransactions(5)
	var recentTxnsDisplay []TransactionDisplayData
	for _, txn := range recentTxns {
		recentTxnsDisplay = append(recentTxnsDisplay, PrepareTransactionForDisplay(txn))
	}
	
	return DashboardData{
		TotalBalance:        totalBalance,
		FormattedTotal:      totalBalance.Format(),
		AccountCount:        len(fs.accounts),
		ActiveAccountCount:  activeAccounts,
		TransactionCount:    len(fs.transactions),
		PendingInvoices:     len(fs.GetPendingInvoices()),
		TopAccounts:         topAccounts,
		RecentTransactions:  recentTxnsDisplay,
		MonthlySpending:     calculateMonthlySpending(fs.transactions),
		MonthlyIncome:       calculateMonthlyIncome(fs.transactions),
	}
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// generateID generates a unique ID with prefix
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// getAccountTypeIcon returns icon for account type
func getAccountTypeIcon(accountType string) string {
	switch accountType {
	case "checking": return "💰"
	case "savings":  return "🏦"
	case "investment": return "📈"
	case "credit":   return "💳"
	default:         return "📋"
	}
}

// getAccountTypeDisplay returns display name for account type
func getAccountTypeDisplay(accountType string) string {
	switch accountType {
	case "checking": return "Checking Account"
	case "savings":  return "Savings Account"
	case "investment": return "Investment Account"
	case "credit":   return "Credit Account"
	default:         return "Unknown Account"
	}
}

// formatTransactionAmount formats transaction amount with +/- prefix
func formatTransactionAmount(transaction Transaction) string {
	prefix := ""
	if transaction.Type == "credit" {
		prefix = "+"
	} else if transaction.Type == "debit" {
		prefix = "-"
	}
	return prefix + transaction.Amount.Format()
}

// getCategoryIcon returns icon for transaction category
func getCategoryIcon(category string) string {
	switch category {
	case "food":           return "🍽️"
	case "transportation": return "🚗"
	case "income":         return "💰"
	case "housing":        return "🏠"
	case "utilities":      return "⚡"
	case "healthcare":     return "⚕️"
	case "entertainment":  return "🎬"
	case "shopping":       return "🛒"
	default:               return "📋"
	}
}

// getInvoiceStatusClass returns CSS class for invoice status
func getInvoiceStatusClass(status string, isOverdue bool) string {
	if isOverdue {
		return "status-error"
	}
	switch status {
	case "paid":              return "status-success"
	case mt.StatusPending: return "status-warning"
	case mt.StatusFailed:  return "status-error"
	default:                  return "status-info"
	}
}

// getInvoiceStatusDisplay returns display text for invoice status
func getInvoiceStatusDisplay(status string) string {
	switch status {
	case "paid":              return "Paid"
	case mt.StatusPending: return "Pending"
	case mt.StatusFailed:  return "Failed"
	default:                  return "Unknown"
	}
}

// calculateMonthlySpending calculates total spending for current month
func calculateMonthlySpending(transactions []Transaction) mt.Money {
	var total mt.Money
	now := time.Now()
	for _, txn := range transactions {
		if txn.Date.Year() == now.Year() && txn.Date.Month() == now.Month() && txn.Type == "debit" {
			total.Amount += txn.Amount.Amount
		}
	}
	return total
}

// calculateMonthlyIncome calculates total income for current month
func calculateMonthlyIncome(transactions []Transaction) mt.Money {
	var total mt.Money
	now := time.Now()
	for _, txn := range transactions {
		if txn.Date.Year() == now.Year() && txn.Date.Month() == now.Month() && txn.Type == "credit" {
			total.Amount += txn.Amount.Amount
		}
	}
	return total
}

// =====================================================
// SAMPLE DATA HELPERS (for demos and testing)
// =====================================================

// SampleCustomer returns sample customer data
func SampleCustomer() Customer {
	return Customer{
		ID:    "cust_001",
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Addresses: []mt.Address{
			{
				Type:       mt.AddressBilling,
				Name:       "John Doe",
				Street1:    "123 Main St",
				City:       "Anytown",
				State:      "NY",
				PostalCode: "12345",
				Country:    "US",
			},
		},
		AccountNumber: "ACC001",
		CreditRating:  "excellent",
		PaymentTerms:  "net30",
		CreditLimit:   mt.NewMoney(10000.00, mt.CurrencyUSD),
		Status:        mt.StatusActive,
		CreatedAt:     time.Now().AddDate(-1, 0, 0),
		Metadata:      make(map[string]string),
	}
}

// SampleAccounts returns sample account data for demos
func SampleAccounts() []Account {
	return []Account{
		{
			ID:          "acc_001",
			Name:        "Main Checking",
			Balance:     mt.NewMoney(1250.00, mt.CurrencyUSD),
			Status:      mt.StatusActive,
			Type:        "checking",
			CreatedAt:   time.Now().AddDate(-2, 0, 0),
			UpdatedAt:   time.Now(),
			Description: "Primary checking account",
			Metadata:    map[string]string{"customer_id": "cust_001"},
		},
		{
			ID:          "acc_002", 
			Name:        "Business Savings",
			Balance:     mt.NewMoney(750.00, mt.CurrencyUSD),
			Status:      mt.StatusActive,
			Type:        "savings",
			CreatedAt:   time.Now().AddDate(-1, -6, 0),
			UpdatedAt:   time.Now(),
			Description: "Business savings account",
			Metadata:    map[string]string{"customer_id": "cust_001"},
		},
		{
			ID:          "acc_003",
			Name:        "Investment Portfolio",
			Balance:     mt.NewMoney(5000.00, mt.CurrencyUSD),
			Status:      mt.StatusActive,
			Type:        "investment",
			CreatedAt:   time.Now().AddDate(-3, 0, 0),
			UpdatedAt:   time.Now(),
			Description: "Investment portfolio account",
			Metadata:    map[string]string{"customer_id": "cust_001"},
		},
	}
}

// SampleTransactions returns sample transaction data
func SampleTransactions() []Transaction {
	return []Transaction{
		{
			ID:          "txn_001",
			AccountID:   "acc_001",
			Amount:      mt.NewMoney(25.00, mt.CurrencyUSD),
			Description: "Client Payment - Project Alpha",
			Date:        time.Now().AddDate(0, 0, -1),
			Status:      mt.StatusCompleted,
			Type:        "credit",
			Category:    "income",
			Reference:   "PAY001",
			Metadata:    make(map[string]string),
		},
		{
			ID:          "txn_002",
			AccountID:   "acc_001",
			Amount:      mt.NewMoney(8.50, mt.CurrencyUSD),
			Description: "Office Supplies",
			Date:        time.Now().AddDate(0, 0, -2),
			Status:      mt.StatusCompleted, 
			Type:        "debit",
			Category:    "other",
			Reference:   "PUR001",
			Metadata:    make(map[string]string),
		},
		{
			ID:          "txn_003",
			AccountID:   "acc_001",
			Amount:      mt.NewMoney(12.00, mt.CurrencyUSD),
			Description: "Software Subscription",
			Date:        time.Now().AddDate(0, 0, -3),
			Status:      mt.StatusPending,
			Type:        "debit",
			Category:    "other",
			Reference:   "SUB001",
			Metadata:    make(map[string]string),
		},
	}
}

// SampleInvoices returns sample invoice data
func SampleInvoices() []Invoice {
	customer := SampleCustomer()
	
	return []Invoice{
		{
			ID:      "inv_001",
			Number:  "INV-2025-001",
			Amount:  mt.NewMoney(3500.00, mt.CurrencyUSD),
			DueDate: time.Now().AddDate(0, 0, 7),
			Status:  mt.StatusPending,
			Customer: customer,
			Items: []InvoiceItem{
				{
					ID:          "item_001",
					Description: "Consulting Services",
					Quantity:    20,
					UnitPrice:   mt.NewMoney(175.00, mt.CurrencyUSD),
					Total:       mt.NewMoney(3500.00, mt.CurrencyUSD),
					Category:    "consulting",
				},
			},
			CreatedAt:   time.Now().AddDate(0, 0, -3),
			Description: "Q4 Consulting Services",
			Metadata:    make(map[string]string),
		},
	}
}
