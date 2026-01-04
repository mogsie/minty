// Package mintyfinui provides UI presentation adapters for the mintyfin domain.
// This package converts pure domain data to UI components, handling all theme
// and styling concerns while keeping the domain layer pure.
package mintyfinui

import (
	"fmt"
	"strings"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
	miex "github.com/ha1tch/minty/mintyex"
	mifi "github.com/ha1tch/minty/domains/mintyfin"
)

// Domain identifier for CSS classes and HTML IDs
const Domain = "mifi"

// =====================================================
// ACCOUNT UI COMPONENTS
// =====================================================

// AccountSummaryCard converts domain account to UI card component
func AccountSummaryCard(theme mui.Theme, account mifi.Account) mi.H {
	displayData := mifi.PrepareAccountForDisplay(account)
	
	return mui.DomainCard(theme, Domain, account.Name, func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_account_summary"),
			b.Div(mi.Class("mifi_account_info"),
				b.P(mi.Class("mifi_account_type"), 
					fmt.Sprintf("%s %s", displayData.TypeIcon, displayData.TypeDisplay)),
				b.Div(mi.Class("mifi_account_balance"),
					b.Span("Balance: "),
					b.Strong(displayData.FormattedBalance),
				),
				StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
			),
			
			b.Div(mi.Class("mifi_account_actions"),
				mui.DomainButton(theme, Domain, "View Details", "secondary",
					mi.Href("/accounts/"+account.ID))(b),
				mui.DomainButton(theme, Domain, "Transactions", "view",
					mi.Href("/accounts/"+account.ID+"/transactions"))(b),
			),
		)
	})
}

// AccountBalance displays a simple balance overview
func AccountBalance(theme mui.Theme, account mifi.Account) mi.H {
	displayData := mifi.PrepareAccountForDisplay(account)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_account_balance_item"),
			b.Span(mi.Class("mifi_account_name"), account.Name),
			b.Span(mi.Class("mifi_account_amount"), displayData.FormattedBalance),
			StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
		)
	}
}

// AccountsTable displays accounts in a table format
func AccountsTable(theme mui.Theme, accounts []mifi.Account) mi.H {
	headers := []string{"Account", "Type", "Balance", "Status", "Actions"}
	
	// Use iterator Map function for cleaner code
	rows := miex.Map(accounts, func(account mifi.Account) []string {
		displayData := mifi.PrepareAccountForDisplay(account)
		actions := fmt.Sprintf(`<div class="mifi_account_table_actions">
			<a href="/accounts/%s" class="mifi_view_button">View</a>
			<a href="/accounts/%s/edit" class="mifi_edit_button">Edit</a>
		</div>`, account.ID, account.ID)
		
		return []string{
			account.Name,
			displayData.TypeDisplay,
			displayData.FormattedBalance,
			displayData.StatusDisplay,
			actions,
		}
	})
	
	return theme.Table(headers, rows)
}

// =====================================================
// TRANSACTION UI COMPONENTS  
// =====================================================

// TransactionList displays a list of transactions
func TransactionList(theme mui.Theme, transactions []mifi.Transaction) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_transaction_list"),
			mi.NewFragment(miex.Each(transactions, func(txn mifi.Transaction) mi.H {
				return TransactionItem(theme, txn)
			})...),
		)
	}
}

// TransactionItem displays a single transaction
func TransactionItem(theme mui.Theme, transaction mifi.Transaction) mi.H {
	displayData := mifi.PrepareTransactionForDisplay(transaction)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_transaction_item"),
			b.Div(mi.Class("mifi_transaction_info"),
				b.Div(mi.Class("mifi_transaction_header"),
					b.Span(mi.Class("mifi_transaction_description"), 
						displayData.CategoryIcon + " " + transaction.Description),
					b.Span(mi.Class("mifi_transaction_date"), 
						displayData.FormattedDate),
				),
				miex.If(displayData.DaysAgo > 0,
					func(b *mi.Builder) mi.Node {
						return b.Small(mi.Class("mifi_transaction_age"),
							fmt.Sprintf("%d days ago", displayData.DaysAgo))
					},
				)(b),
			),
			b.Div(mi.Class("mifi_transaction_details"),
				b.Span(mi.Class(displayData.TypeClass), displayData.FormattedAmount),
				StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
			),
		)
	}
}

// TransactionTable displays transactions in table format
func TransactionTable(theme mui.Theme, transactions []mifi.Transaction) mi.H {
	headers := []string{"Date", "Description", "Category", "Amount", "Status"}
	
	// Use iterator Map function for cleaner code
	rows := miex.Map(transactions, func(txn mifi.Transaction) []string {
		displayData := mifi.PrepareTransactionForDisplay(txn)
		return []string{
			displayData.FormattedDate,
			displayData.CategoryIcon + " " + txn.Description,
			strings.Title(txn.Category),
			displayData.FormattedAmount,
			displayData.StatusDisplay,
		}
	})
	
	return theme.Table(headers, rows)
}

// =====================================================
// INVOICE UI COMPONENTS
// =====================================================

// InvoiceCard displays invoice information with payment options
func InvoiceCard(theme mui.Theme, invoice mifi.Invoice) mi.H {
	displayData := mifi.PrepareInvoiceForDisplay(invoice)
	
	return theme.Card(fmt.Sprintf("Invoice #%s", invoice.Number), 
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mifi_invoice_card"),
				b.Div(mi.Class("mifi_invoice_info"),
					b.P("Customer: ", invoice.Customer.Name),
					b.P("Amount: ", b.Strong(displayData.FormattedAmount)),
					b.P("Due Date: ", displayData.FormattedDueDate),
					miex.IfElse(displayData.IsOverdue,
						func(b *mi.Builder) mi.Node {
							return b.P(mi.Class("mifi_invoice_overdue"), "⚠️ OVERDUE")
						},
						func(b *mi.Builder) mi.Node {
							return b.P(mi.Class("mifi_invoice_due_soon"),
								fmt.Sprintf("Due in %d days", displayData.DaysUntilDue))
						},
					)(b),
					StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
				),
				
				miex.If(invoice.Status == miex.StatusPending,
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mifi_invoice_actions"),
							PaymentButton(theme, invoice)(b),
							mui.DomainButton(theme, Domain, "View Details", "secondary",
								mi.Href("/invoices/"+invoice.ID))(b),
						)
					},
				)(b),
			)
		})
}

// PaymentButton creates a payment button for invoices
func PaymentButton(theme mui.Theme, invoice mifi.Invoice) mi.H {
	return mui.DomainButton(theme, Domain, 
		fmt.Sprintf("Pay %s", invoice.Amount.Format()), "payment",
		mi.HxPost("/api/invoices/"+invoice.ID+"/pay"),
		mi.HxTarget("#mifi_payment_result"),
		mi.HxIndicator("#mifi_payment_spinner"),
		mi.HxConfirm(fmt.Sprintf("Pay invoice #%s for %s?", 
			invoice.Number, invoice.Amount.Format())),
	)
}

// InvoicesTable displays invoices in table format
func InvoicesTable(theme mui.Theme, invoices []mifi.Invoice) mi.H {
	headers := []string{"Invoice #", "Customer", "Amount", "Due Date", "Status", "Actions"}
	
	// Use iterator Map function for cleaner code
	rows := miex.Map(invoices, func(invoice mifi.Invoice) []string {
		displayData := mifi.PrepareInvoiceForDisplay(invoice)
		actions := ""
		
		if invoice.Status == miex.StatusPending {
			actions = fmt.Sprintf(`<div class="mifi_invoice_table_actions">
				<button class="mifi_pay_button" data-invoice="%s">Pay</button>
				<a href="/invoices/%s" class="mifi_view_button">View</a>
			</div>`, invoice.ID, invoice.ID)
		} else {
			actions = fmt.Sprintf(`<div class="mifi_invoice_table_actions">
				<a href="/invoices/%s" class="mifi_view_button">View</a>
			</div>`, invoice.ID)
		}
		
		return []string{
			invoice.Number,
			invoice.Customer.Name,
			displayData.FormattedAmount,
			displayData.FormattedDueDate,
			displayData.StatusDisplay,
			actions,
		}
	})
	
	return theme.Table(headers, rows)
}

// =====================================================
// DASHBOARD UI COMPONENTS
// =====================================================

// FinancialDashboard creates a complete financial dashboard
func FinancialDashboard(theme mui.Theme, dashboardData mifi.DashboardData, 
	recentTxns []mifi.Transaction, pendingInvoices []mifi.Invoice) mi.H {
	
	return mui.Dashboard(theme, "Financial Dashboard",
		// Sidebar
		func(b *mi.Builder) mi.Node {
			return theme.Sidebar(func(b *mi.Builder) mi.Node {
				return b.Div(mi.Class("mifi_nav"),
					b.H4("Finance"),
					theme.List([]string{
						"Dashboard", "Accounts", "Transactions", 
						"Invoices", "Reports", "Settings",
					}, false)(b),
				)
			})(b)
		},
		
		// Main content
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mifi_dashboard_main"),
				// Financial metrics
				MetricsSection(theme, dashboardData)(b),
				// Accounts section  
				AccountsSection(theme, dashboardData.TopAccounts)(b),
				// Recent transactions
				RecentTransactionsSection(theme, dashboardData.RecentTransactions)(b),
				// Pending invoices
				miex.If(len(pendingInvoices) > 0, 
					PendingInvoicesSection(theme, pendingInvoices))(b),
			)
		},
	)
}

// MetricsSection displays financial overview metrics
func MetricsSection(theme mui.Theme, data mifi.DashboardData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mifi_metrics_section"),
			b.H2("Financial Overview"),
			miex.GridLayout(4, "1rem")(
				mui.StatsCard(theme, "Total Balance", 
					data.FormattedTotal, "Across all accounts"),
				mui.StatsCard(theme, "Active Accounts", 
					fmt.Sprintf("%d", data.ActiveAccountCount), "Out of " + fmt.Sprintf("%d", data.AccountCount)),
				mui.StatsCard(theme, "Pending Invoices", 
					fmt.Sprintf("%d", data.PendingInvoices), "Awaiting payment"),
				mui.StatsCard(theme, "Recent Transactions", 
					fmt.Sprintf("%d", len(data.RecentTransactions)), "This week"),
			)(b),
		)
	}
}

// AccountsSection displays account summary cards
func AccountsSection(theme mui.Theme, accounts []mifi.AccountDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mifi_accounts_section"),
			b.H2("Account Summary"),
			miex.GridLayout(3, "1rem")(
				miex.EachH(accounts, func(data mifi.AccountDisplayData) mi.H {
					return AccountSummaryCard(theme, data.Account)
				})...,
			)(b),
			mui.DomainButton(theme, Domain, "View All Accounts", "view",
				mi.Href("/accounts"))(b),
		)
	}
}

// RecentTransactionsSection displays recent transactions
func RecentTransactionsSection(theme mui.Theme, transactions []mifi.TransactionDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mifi_transactions_section"),
			b.H2("Recent Transactions"),
			miex.IfElse(len(transactions) > 0,
				func(b *mi.Builder) mi.Node {
					return b.Div(mi.Class("mifi_transaction_list"),
						mi.NewFragment(miex.Each(transactions, func(data mifi.TransactionDisplayData) mi.H {
							return TransactionItemFromDisplayData(theme, data)
						})...),
					)
				},
				func(b *mi.Builder) mi.Node {
					return b.P(mi.Class("mifi_no_transactions"), "No recent transactions")
				},
			)(b),
			mui.DomainButton(theme, Domain, "View All Transactions", "view",
				mi.Href("/transactions"))(b),
		)
	}
}

// PendingInvoicesSection displays pending invoices
func PendingInvoicesSection(theme mui.Theme, invoices []mifi.Invoice) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mifi_invoices_section"),
			b.H2("Pending Invoices"),
			miex.GridLayout(2, "1rem")(
				miex.EachH(invoices, func(invoice mifi.Invoice) mi.H {
					return InvoiceCard(theme, invoice)
				})...,
			)(b),
		)
	}
}

// =====================================================
// FORM UI COMPONENTS
// =====================================================

// AccountForm creates a form for creating/editing accounts
func AccountForm(theme mui.Theme, account *mifi.Account, isEdit bool) mi.H {
	title := "Create Account"
	action := "/accounts"
	submitText := "Create Account"
	
	if isEdit && account != nil {
		title = "Edit Account"
		action = "/accounts/" + account.ID
		submitText = "Update Account"
	}
	
	accountTypes := []mui.SelectOption{
		{Value: "checking", Text: "Checking Account"},
		{Value: "savings", Text: "Savings Account"},  
		{Value: "investment", Text: "Investment Account"},
		{Value: "credit", Text: "Credit Account"},
	}
	
	return theme.Card(title, func(b *mi.Builder) mi.Node {
		return b.Form(mi.Action(action), mi.Method("POST"),
			mi.Class("mifi_account_form"),
			theme.FormInput("Account Name", "name", "text", 
				mi.Required(), mi.Value(getAccountValue(account, "name")))(b),
			theme.FormSelect("Account Type", "type", accountTypes)(b),
			theme.FormInput("Initial Balance", "balance", "number",
				mi.Step("0.01"), mi.Min("0"),
				mi.Value(getAccountValue(account, "balance")))(b),
			theme.FormTextarea("Description", "description", 
				mi.Value(getAccountValue(account, "description")))(b),
			theme.PrimaryButton(submitText, mi.Type("submit"))(b),
		)
	})
}

// TransactionForm creates a form for creating transactions
func TransactionForm(theme mui.Theme, accountID string) mi.H {
	transactionTypes := []mui.SelectOption{
		{Value: "credit", Text: "Credit (Money In)"},
		{Value: "debit", Text: "Debit (Money Out)"},
	}
	
	return theme.Card("Create Transaction", func(b *mi.Builder) mi.Node {
		return b.Form(mi.Action("/transactions"), mi.Method("POST"),
			mi.Class("mifi_transaction_form"),
			b.Input(mi.Type("hidden"), mi.Name("account_id"), mi.Value(accountID)),
			theme.FormSelect("Transaction Type", "type", transactionTypes)(b),
			theme.FormInput("Amount", "amount", "number", 
				mi.Required(), mi.Step("0.01"), mi.Min("0.01"))(b),
			theme.FormInput("Description", "description", "text", mi.Required())(b),
			theme.FormInput("Date", "date", "date", mi.Required())(b),
			theme.PrimaryButton("Create Transaction", mi.Type("submit"))(b),
		)
	})
}

// InvoiceForm creates a form for creating invoices
func InvoiceForm(theme mui.Theme) mi.H {
	return theme.Card("Create Invoice", func(b *mi.Builder) mi.Node {
		return b.Form(mi.Action("/invoices"), mi.Method("POST"),
			mi.Class("mifi_invoice_form"),
			theme.FormInput("Invoice Number", "number", "text", mi.Required())(b),
			theme.FormInput("Customer Name", "customer_name", "text", mi.Required())(b),
			theme.FormInput("Customer Email", "customer_email", "email", mi.Required())(b),
			theme.FormInput("Amount", "amount", "number", 
				mi.Required(), mi.Step("0.01"), mi.Min("0.01"))(b),
			theme.FormInput("Due Date", "due_date", "date", mi.Required())(b),
			theme.FormTextarea("Description", "description")(b),
			theme.PrimaryButton("Create Invoice", mi.Type("submit"))(b),
		)
	})
}

// =====================================================
// HELPER UI COMPONENTS
// =====================================================

// StatusBadge creates a status badge with appropriate styling
func StatusBadge(theme mui.Theme, statusText, statusClass string) mi.H {
	variant := getStatusVariant(statusClass)
	return theme.Badge(statusText, variant)
}

// TransactionItemFromDisplayData renders transaction from display data
func TransactionItemFromDisplayData(theme mui.Theme, data mifi.TransactionDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_transaction_item"),
			b.Div(mi.Class("mifi_transaction_info"),
				b.Div(mi.Class("mifi_transaction_header"),
					b.Span(mi.Class("mifi_transaction_description"), 
						data.CategoryIcon + " " + data.Transaction.Description),
					b.Span(mi.Class("mifi_transaction_date"), 
						data.FormattedDate),
				),
				miex.If(data.DaysAgo > 0,
					func(b *mi.Builder) mi.Node {
						return b.Small(mi.Class("mifi_transaction_age"),
							fmt.Sprintf("%d days ago", data.DaysAgo))
					},
				)(b),
			),
			b.Div(mi.Class("mifi_transaction_details"),
				b.Span(mi.Class(data.TypeClass), data.FormattedAmount),
				StatusBadge(theme, data.StatusDisplay, data.StatusClass)(b),
			),
		)
	}
}

// MoneyInput creates a money input with currency formatting
func MoneyInput(theme mui.Theme, label, name, currency string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_money_input"),
			theme.FormLabel(label, name)(b),
			b.Div(mi.Class("mifi_input_group"),
				b.Span(mi.Class("mifi_currency_prefix"), getCurrencySymbol(currency)),
				theme.Input(name, "number", mi.Step("0.01"), mi.Min("0"), mi.Required())(b),
			),
		)
	}
}

// AccountSummaryWidget creates a compact account summary
func AccountSummaryWidget(theme mui.Theme, account mifi.Account) mi.H {
	displayData := mifi.PrepareAccountForDisplay(account)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_account_widget"),
			b.Div(mi.Class("mifi_widget_header"),
				b.Span(mi.Class("mifi_account_icon"), displayData.TypeIcon),
				b.Span(mi.Class("mifi_account_name"), account.Name),
			),
			b.Div(mi.Class("mifi_widget_content"),
				b.Div(mi.Class("mifi_balance"), displayData.FormattedBalance),
				StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
			),
		)
	}
}

// =====================================================
// PAGE LAYOUTS
// =====================================================

// AccountsPage creates a complete accounts management page
func AccountsPage(theme mui.Theme, accounts []mifi.Account) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_accounts_page"),
			b.Header(mi.Class("mifi_page_header"),
				b.H1("Accounts"),
				mui.DomainButton(theme, Domain, "Create Account", "primary",
					mi.Href("/accounts/new"))(b),
			),
			
			b.Main(mi.Class("mifi_page_content"),
				miex.IfElse(len(accounts) > 0,
					AccountsTable(theme, accounts),
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mifi_empty_state"),
							b.P("No accounts found."),
							mui.DomainButton(theme, Domain, "Create Your First Account", "primary",
								mi.Href("/accounts/new"))(b),
						)
					},
				)(b),
			),
		)
	}
}

// TransactionsPage creates a complete transactions page
func TransactionsPage(theme mui.Theme, transactions []mifi.Transaction, accountID string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_transactions_page"),
			b.Header(mi.Class("mifi_page_header"),
				b.H1("Transactions"),
				mui.DomainButton(theme, Domain, "Add Transaction", "primary",
					mi.Href("/transactions/new?account_id="+accountID))(b),
			),
			
			b.Main(mi.Class("mifi_page_content"),
				miex.IfElse(len(transactions) > 0,
					TransactionTable(theme, transactions),
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mifi_empty_state"),
							b.P("No transactions found."),
						)
					},
				)(b),
			),
		)
	}
}

// InvoicesPage creates a complete invoices management page
func InvoicesPage(theme mui.Theme, invoices []mifi.Invoice) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mifi_invoices_page"),
			b.Header(mi.Class("mifi_page_header"),
				b.H1("Invoices"),
				mui.DomainButton(theme, Domain, "Create Invoice", "primary",
					mi.Href("/invoices/new"))(b),
			),
			
			b.Main(mi.Class("mifi_page_content"),
				miex.IfElse(len(invoices) > 0,
					InvoicesTable(theme, invoices),
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mifi_empty_state"),
							b.P("No invoices found."),
							mui.DomainButton(theme, Domain, "Create Your First Invoice", "primary",
								mi.Href("/invoices/new"))(b),
						)
					},
				)(b),
			),
		)
	}
}

// =====================================================
// UTILITY FUNCTIONS
// =====================================================

// getStatusVariant converts status class to theme variant
func getStatusVariant(statusClass string) string {
	switch statusClass {
	case "status-success": return "success"
	case "status-warning": return "warning"  
	case "status-error":   return "danger"
	case "status-info":    return "info"
	default:               return "secondary"
	}
}

// getCurrencySymbol returns currency symbol for display
func getCurrencySymbol(currency string) string {
	switch strings.ToUpper(currency) {
	case miex.CurrencyUSD: return "$"
	case miex.CurrencyEUR: return "€"
	case miex.CurrencyGBP: return "£"
	case miex.CurrencyJPY: return "¥"
	default:                  return currency + " "
	}
}

// getAccountValue safely extracts account field values
func getAccountValue(account *mifi.Account, field string) string {
	if account == nil {
		return ""
	}
	
	switch field {
	case "name":        return account.Name
	case "description": return account.Description
	case "balance":     return fmt.Sprintf("%.2f", account.Balance.MajorUnit())
	default:            return ""
	}
}

// =====================================================
// INTEGRATION HELPERS
// =====================================================

// CreateFinanceDemoPage creates a complete demo page using sample data
func CreateFinanceDemoPage(theme mui.Theme) mi.H {
	// Use pure domain functions to create sample data
	service := mifi.NewFinanceService()
	sampleAccounts := mifi.SampleAccounts()
	sampleTransactions := mifi.SampleTransactions()
	sampleInvoices := mifi.SampleInvoices()
	
	// Add accounts to service
	for _, account := range sampleAccounts {
		service.CreateAccount(account.Name, account.Type, account.Balance, "demo_customer")
	}
	
	// Prepare dashboard data using pure domain functions
	dashboardData := mifi.PrepareDashboardData(service)
	
	// Create UI using presentation adapters
	return FinancialDashboard(theme, dashboardData, sampleTransactions, sampleInvoices)
}

// IntegrateWithMainApp shows how to integrate with main application
func IntegrateWithMainApp(theme mui.Theme, financeService *mifi.FinanceService) mi.H {
	// Get business data from pure domain service
	dashboardData := mifi.PrepareDashboardData(financeService)
	recentTransactions := financeService.GetRecentTransactions(5)
	pendingInvoices := financeService.GetPendingInvoices()
	
	// Use presentation adapters to create UI
	return FinancialDashboard(theme, dashboardData, recentTransactions, pendingInvoices)
}
