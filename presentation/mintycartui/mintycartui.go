// Package mintycartui provides UI presentation adapters for the mintycart domain.
// This package converts pure e-commerce domain data to UI components.
package mintycartui

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
	miex "github.com/ha1tch/minty/mintyex"
	mica "github.com/ha1tch/minty/domains/mintycart"
)

// Domain identifier for CSS classes and HTML IDs
const Domain = "mica"

// =====================================================
// PRODUCT UI COMPONENTS
// =====================================================

// ProductCard displays product information with add to cart
func ProductCard(theme mui.Theme, product mica.Product) mi.H {
	displayData := mica.PrepareProductForDisplay(product)
	
	return mui.DomainCard(theme, Domain, product.Name,
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mica_product_card"),
				miex.If(displayData.PrimaryImageURL != "",
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mica_product_image"),
							b.Img(mi.Src(displayData.PrimaryImageURL), 
								   mi.Alt(product.Name)),
						)
					},
				)(b),
				
				b.Div(mi.Class("mica_product_info"),
					b.P(mi.Class("mica_product_price"), b.Strong(displayData.FormattedPrice)),
					b.P(mi.Class("mica_product_description"), product.Description),
					StatusBadge(theme, displayData.InventoryStatus, "inventory")(b),
				),
				
				miex.IfElse(displayData.InStock,
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mica_product_actions"),
							AddToCartButton(theme, product)(b),
							mui.DomainButton(theme, Domain, "View Details", "secondary",
								mi.Href("/products/"+product.ID))(b),
						)
					},
					func(b *mi.Builder) mi.Node {
						return b.Div(mi.Class("mica_product_unavailable"),
							b.P("Out of Stock"),
						)
					},
				)(b),
			)
		})
}

// ProductList displays a grid of products
func ProductList(theme mui.Theme, products []mica.Product) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_product_list"),
			miex.GridLayout(3, "1rem")(
				miex.EachH(products, func(product mica.Product) mi.H {
					return ProductCard(theme, product)
				})...,
			)(b),
		)
	}
}

// AddToCartButton creates an add to cart button
func AddToCartButton(theme mui.Theme, product mica.Product) mi.H {
	return mui.DomainButton(theme, Domain, "Add to Cart", "primary",
		mi.HxPost("/api/cart/add"),
		mi.HxVals(fmt.Sprintf(`{"product_id": "%s", "quantity": 1}`, product.ID)),
		mi.HxTarget("#cart-count"),
		mi.HxSwap("innerHTML"),
	)
}

// =====================================================
// CART UI COMPONENTS
// =====================================================

// CartWidget displays cart summary
func CartWidget(theme mui.Theme, cart mica.Cart) mi.H {
	displayData := mica.PrepareCartForDisplay(cart)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_cart_widget"),
			b.A(mi.Href("/cart"), mi.Class("mica_cart_link"),
				b.Span(mi.Class("mica_cart_icon"), "🛒"),
				b.Span(mi.ID("cart-count"), mi.Class("mica_cart_count"), 
					fmt.Sprintf("%d", displayData.ItemCount)),
				b.Span(mi.Class("mica_cart_total"), displayData.FormattedTotal),
			),
		)
	}
}

// CartPage displays full cart with items and checkout
func CartPage(theme mui.Theme, cart mica.Cart) mi.H {
	displayData := mica.PrepareCartForDisplay(cart)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_cart_page"),
			b.Header(mi.Class("mica_page_header"),
				b.H1("Shopping Cart"),
			),
			
			miex.IfElse(displayData.IsEmpty,
				// Empty cart
				func(b *mi.Builder) mi.Node {
					return b.Div(mi.Class("mica_empty_cart"),
						b.P("Your cart is empty"),
						mui.DomainButton(theme, Domain, "Continue Shopping", "primary",
							mi.Href("/products"))(b),
					)
				},
				// Cart with items
				func(b *mi.Builder) mi.Node {
					return b.Div(mi.Class("mica_cart_content"),
						CartItemsList(theme, cart.Items)(b),
						CartSummary(theme, displayData)(b),
						CheckoutButton(theme)(b),
					)
				},
			)(b),
		)
	}
}

// CartItemsList displays cart items
func CartItemsList(theme mui.Theme, items []mica.CartItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_cart_items"),
			mi.NewFragment(miex.Each(items, func(item mica.CartItem) mi.H {
				return CartItemRow(theme, item)
			})...),
		)
	}
}

// CartItemRow displays a single cart item
func CartItemRow(theme mui.Theme, item mica.CartItem) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_cart_item"),
			b.Div(mi.Class("mica_item_info"),
				b.H4(item.Product.Name),
				b.P(item.Product.Description),
				b.P("Price: ", item.Price.Format()),
			),
			
			b.Div(mi.Class("mica_item_quantity"),
				b.Label(mi.DataAttr("for", "qty_"+item.ID), "Quantity: "),
				b.Input(mi.Type("number"), mi.Name("quantity"), 
					mi.Value(fmt.Sprintf("%d", item.Quantity)),
					mi.Min("1"), mi.ID("qty_"+item.ID)),
			),
			
			b.Div(mi.Class("mica_item_total"),
				b.Strong(item.Total.Format()),
			),
			
			b.Div(mi.Class("mica_item_actions"),
				mui.DomainButton(theme, Domain, "Remove", "danger",
					mi.HxDelete("/api/cart/items/"+item.ID),
					mi.HxTarget("closest .mica_cart_item"),
					mi.HxConfirm("Remove this item from cart?"))(b),
			),
		)
	}
}

// CartSummary displays cart totals
func CartSummary(theme mui.Theme, displayData mica.CartDisplayData) mi.H {
	return theme.Card("Order Summary",
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mica_cart_summary"),
				b.Div(mi.Class("mica_summary_line"),
					b.Span("Subtotal:"),
					b.Span(displayData.FormattedSubtotal),
				),
				b.Div(mi.Class("mica_summary_line"),
					b.Span("Tax:"),
					b.Span(displayData.FormattedTax),
				),
				b.Div(mi.Class("mica_summary_line"),
					b.Span("Shipping:"),
					b.Span(displayData.FormattedShipping),
				),
				b.Hr(),
				b.Div(mi.Class("mica_summary_total"),
					b.Strong("Total: ", displayData.FormattedTotal),
				),
			)
		})
}

// CheckoutButton creates a checkout button
func CheckoutButton(theme mui.Theme) mi.H {
	return mui.DomainButton(theme, Domain, "Proceed to Checkout", "primary",
		mi.Href("/checkout"))
}

// =====================================================
// ORDER UI COMPONENTS
// =====================================================

// OrderCard displays order information
func OrderCard(theme mui.Theme, order mica.Order) mi.H {
	displayData := mica.PrepareOrderForDisplay(order)
	
	return mui.DomainCard(theme, Domain, fmt.Sprintf("Order #%s", order.Number),
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mica_order_card"),
				b.Div(mi.Class("mica_order_info"),
					b.P("Total: ", b.Strong(displayData.FormattedTotal)),
					b.P("Placed: ", fmt.Sprintf("%d days ago", displayData.DaysAgo)),
					StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
				),
				
				miex.If(displayData.TrackingNumber != "",
					func(b *mi.Builder) mi.Node {
						return b.P("Tracking: ", displayData.TrackingNumber)
					},
				)(b),
				
				b.Div(mi.Class("mica_order_actions"),
					mui.DomainButton(theme, Domain, "View Order", "secondary",
						mi.Href("/orders/"+order.ID))(b),
				),
			)
		})
}

// OrderList displays a list of orders
func OrderList(theme mui.Theme, orders []mica.Order) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mica_order_list"),
			mi.NewFragment(miex.Each(orders, func(order mica.Order) mi.H {
				return OrderCard(theme, order)
			})...),
		)
	}
}

// =====================================================
// DASHBOARD UI COMPONENTS
// =====================================================

// EcommerceDashboard creates a complete e-commerce dashboard
func EcommerceDashboard(theme mui.Theme, dashboardData mica.DashboardData) mi.H {
	return mui.Dashboard(theme, "E-commerce Dashboard",
		// Sidebar
		func(b *mi.Builder) mi.Node {
			return theme.Sidebar(func(b *mi.Builder) mi.Node {
				return b.Div(mi.Class("mica_nav"),
					b.H4("E-commerce"),
					theme.List([]string{
						"Dashboard", "Products", "Orders", 
						"Customers", "Analytics", "Settings",
					}, false)(b),
				)
			})(b)
		},
		
		// Main content
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mica_dashboard_main"),
				// E-commerce metrics
				MetricsSection(theme, dashboardData)(b),
				// Recent orders
				RecentOrdersSection(theme, dashboardData.RecentOrders)(b),
				// Top products
				TopProductsSection(theme, dashboardData.TopProducts)(b),
			)
		},
	)
}

// MetricsSection displays e-commerce overview metrics
func MetricsSection(theme mui.Theme, data mica.DashboardData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mica_metrics_section"),
			b.H2("Store Overview"),
			miex.GridLayout(4, "1rem")(
				mui.StatsCard(theme, "Total Products", 
					fmt.Sprintf("%d", data.TotalProducts), "In catalog"),
				mui.StatsCard(theme, "Active Products", 
					fmt.Sprintf("%d", data.ActiveProducts), "Available for sale"),
				mui.StatsCard(theme, "Total Orders", 
					fmt.Sprintf("%d", data.TotalOrders), "All time"),
				mui.StatsCard(theme, "Revenue Today", 
					data.FormattedRevenue, "Today's sales"),
			)(b),
		)
	}
}

// RecentOrdersSection displays recent orders
func RecentOrdersSection(theme mui.Theme, orders []mica.OrderDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mica_orders_section"),
			b.H2("Recent Orders"),
			miex.IfElse(len(orders) > 0,
				func(b *mi.Builder) mi.Node {
					return b.Div(mi.Class("mica_order_list"),
						mi.NewFragment(miex.Each(orders, func(data mica.OrderDisplayData) mi.H {
							return OrderCard(theme, data.Order)
						})...),
					)
				},
				func(b *mi.Builder) mi.Node {
					return b.P(mi.Class("mica_no_orders"), "No recent orders")
				},
			)(b),
			mui.DomainButton(theme, Domain, "View All Orders", "view",
				mi.Href("/orders"))(b),
		)
	}
}

// TopProductsSection displays top products
func TopProductsSection(theme mui.Theme, products []mica.ProductDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mica_products_section"),
			b.H2("Top Products"),
			miex.IfElse(len(products) > 0,
				func(b *mi.Builder) mi.Node {
					return miex.GridLayout(3, "1rem")(
						miex.EachH(products, func(data mica.ProductDisplayData) mi.H {
							return ProductCard(theme, data.Product)
						})...,
					)(b)
				},
				func(b *mi.Builder) mi.Node {
					return b.P(mi.Class("mica_no_products"), "No products found")
				},
			)(b),
			mui.DomainButton(theme, Domain, "View All Products", "view",
				mi.Href("/products"))(b),
		)
	}
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// StatusBadge creates a status badge with appropriate styling
func StatusBadge(theme mui.Theme, statusText, statusType string) mi.H {
	variant := getStatusVariant(statusType)
	return theme.Badge(statusText, variant)
}

// getStatusVariant converts status type to theme variant
func getStatusVariant(statusType string) string {
	switch statusType {
	case "inventory":
		return "info"
	case "status-success": 
		return "success"
	case "status-warning": 
		return "warning"  
	case "status-error":   
		return "danger"
	case "status-info":    
		return "info"
	default:               
		return "secondary"
	}
}

// =====================================================
// INTEGRATION HELPERS
// =====================================================

// CreateEcommerceDemoPage creates a complete demo page using sample data
func CreateEcommerceDemoPage(theme mui.Theme) mi.H {
	// Use pure domain functions to create sample data
	service := mica.NewEcommerceService()
	sampleProducts := mica.SampleProducts()
	
	// Add data to service (simplified for demo)
	for _, product := range sampleProducts {
		service.CreateProduct(product.Name, product.Description, product.SKU, 
			product.Category, product.Price, product.Weight, product.Inventory)
	}
	
	// Prepare dashboard data using pure domain functions
	dashboardData := mica.PrepareDashboardData(service)
	
	// Create UI using presentation adapters
	return EcommerceDashboard(theme, dashboardData)
}

// IntegrateWithMainApp shows how to integrate with main application
func IntegrateWithMainApp(theme mui.Theme, ecommerceService *mica.EcommerceService) mi.H {
	// Get business data from pure domain service
	dashboardData := mica.PrepareDashboardData(ecommerceService)
	
	// Use presentation adapters to create UI
	return EcommerceDashboard(theme, dashboardData)
}
