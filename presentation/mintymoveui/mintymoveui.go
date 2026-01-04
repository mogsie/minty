// Package mintymoveui provides UI presentation adapters for the mintymove domain.
// This package converts pure logistics domain data to UI components.
package mintymoveui

import (
	"fmt"

	mi "github.com/ha1tch/minty"
	mui "github.com/ha1tch/minty/mintyui"
	miex "github.com/ha1tch/minty/mintyex"
	mimo "github.com/ha1tch/minty/domains/mintymove"
)

// Domain identifier for CSS classes and HTML IDs
const Domain = "mimo"

// =====================================================
// SHIPMENT UI COMPONENTS
// =====================================================

// ShipmentCard displays shipment information with tracking
func ShipmentCard(theme mui.Theme, shipment mimo.Shipment) mi.H {
	displayData := mimo.PrepareShipmentForDisplay(shipment)
	
	return mui.DomainCard(theme, Domain, fmt.Sprintf("Shipment %s", shipment.TrackingCode), 
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mimo_shipment_card"),
				b.Div(mi.Class("mimo_shipment_info"),
					b.P("From: ", shipment.Origin.FormatOneLine()),
					b.P("To: ", shipment.Destination.FormatOneLine()),
					b.P("Carrier: ", shipment.Carrier, " (", shipment.Service, ")"),
					b.P("Cost: ", b.Strong(displayData.FormattedCost)),
					StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
				),
				
				b.Div(mi.Class("mimo_shipment_progress"),
					mui.ProgressBar(displayData.ProgressPercent, 100, "Progress")(b),
				),
				
				b.Div(mi.Class("mimo_shipment_actions"),
					mui.DomainButton(theme, Domain, "Track", "primary",
						mi.Href("/shipments/"+shipment.ID))(b),
				),
			)
		})
}

// ShipmentList displays a list of shipments
func ShipmentList(theme mui.Theme, shipments []mimo.Shipment) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mimo_shipment_list"),
			mi.NewFragment(miex.Each(shipments, func(shipment mimo.Shipment) mi.H {
				return ShipmentCard(theme, shipment)
			})...),
		)
	}
}

// TrackingWidget displays shipment tracking information
func TrackingWidget(theme mui.Theme, shipment mimo.Shipment) mi.H {
	displayData := mimo.PrepareShipmentForDisplay(shipment)
	
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mimo_tracking_widget"),
			b.H3("Tracking: ", shipment.TrackingCode),
			b.Div(mi.Class("mimo_tracking_status"),
				b.Div(mi.Class("mimo_status_icon"), getShipmentStatusIcon(shipment.Status)),
				b.Div(mi.Class("mimo_status_text"),
					b.Strong(displayData.StatusDisplay),
					b.Br(),
					b.Small(fmt.Sprintf("%d days in transit", displayData.DaysInTransit)),
				),
			),
			b.Div(mi.Class("mimo_tracking_progress"),
				mui.ProgressBar(displayData.ProgressPercent, 100, "Delivery Progress")(b),
			),
		)
	}
}

// =====================================================
// VEHICLE UI COMPONENTS
// =====================================================

// VehicleCard displays vehicle information
func VehicleCard(theme mui.Theme, vehicle mimo.Vehicle) mi.H {
	displayData := mimo.PrepareVehicleForDisplay(vehicle)
	
	return mui.DomainCard(theme, Domain, vehicle.Name,
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mimo_vehicle_card"),
				b.Div(mi.Class("mimo_vehicle_info"),
					b.P(mi.Class("mimo_vehicle_type"), displayData.TypeIcon, " ", vehicle.Type),
					b.P("License: ", vehicle.LicensePlate),
					b.P("Capacity: ", displayData.CapacityUsed),
					b.P("Driver: ", vehicle.Driver.Name),
					StatusBadge(theme, displayData.StatusDisplay, displayData.StatusClass)(b),
				),
				
				b.Div(mi.Class("mimo_vehicle_actions"),
					mui.DomainButton(theme, Domain, "View Details", "secondary",
						mi.Href("/vehicles/"+vehicle.ID))(b),
				),
			)
		})
}

// VehicleList displays a list of vehicles
func VehicleList(theme mui.Theme, vehicles []mimo.Vehicle) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("mimo_vehicle_list"),
			miex.GridLayout(3, "1rem")(
				miex.EachH(vehicles, func(vehicle mimo.Vehicle) mi.H {
					return VehicleCard(theme, vehicle)
				})...,
			)(b),
		)
	}
}

// =====================================================
// DASHBOARD UI COMPONENTS
// =====================================================

// LogisticsDashboard creates a complete logistics dashboard
func LogisticsDashboard(theme mui.Theme, dashboardData mimo.DashboardData) mi.H {
	return mui.Dashboard(theme, "Logistics Dashboard",
		// Sidebar
		func(b *mi.Builder) mi.Node {
			return theme.Sidebar(func(b *mi.Builder) mi.Node {
				return b.Div(mi.Class("mimo_nav"),
					b.H4("Logistics"),
					theme.List([]string{
						"Dashboard", "Shipments", "Routes", 
						"Vehicles", "Drivers", "Reports",
					}, false)(b),
				)
			})(b)
		},
		
		// Main content
		func(b *mi.Builder) mi.Node {
			return b.Div(mi.Class("mimo_dashboard_main"),
				// Logistics metrics
				MetricsSection(theme, dashboardData)(b),
				// Recent shipments
				RecentShipmentsSection(theme, dashboardData.RecentShipments)(b),
			)
		},
	)
}

// MetricsSection displays logistics overview metrics
func MetricsSection(theme mui.Theme, data mimo.DashboardData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mimo_metrics_section"),
			b.H2("Logistics Overview"),
			miex.GridLayout(4, "1rem")(
				mui.StatsCard(theme, "Total Shipments", 
					fmt.Sprintf("%d", data.TotalShipments), "All time"),
				mui.StatsCard(theme, "Active Shipments", 
					fmt.Sprintf("%d", data.ActiveShipments), "In transit"),
				mui.StatsCard(theme, "Available Vehicles", 
					fmt.Sprintf("%d", data.AvailableVehicles), "Ready for dispatch"),
				mui.StatsCard(theme, "Revenue", 
					data.FormattedRevenue, "Total earned"),
			)(b),
		)
	}
}

// RecentShipmentsSection displays recent shipments
func RecentShipmentsSection(theme mui.Theme, shipments []mimo.ShipmentDisplayData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Section(mi.Class("mimo_shipments_section"),
			b.H2("Recent Shipments"),
			miex.IfElse(len(shipments) > 0,
				func(b *mi.Builder) mi.Node {
					return b.Div(mi.Class("mimo_shipment_list"),
						mi.NewFragment(miex.Each(shipments, func(data mimo.ShipmentDisplayData) mi.H {
							return ShipmentCard(theme, data.Shipment)
						})...),
					)
				},
				func(b *mi.Builder) mi.Node {
					return b.P(mi.Class("mimo_no_shipments"), "No recent shipments")
				},
			)(b),
			mui.DomainButton(theme, Domain, "View All Shipments", "view",
				mi.Href("/shipments"))(b),
		)
	}
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// StatusBadge creates a status badge with appropriate styling
func StatusBadge(theme mui.Theme, statusText, statusClass string) mi.H {
	variant := getStatusVariant(statusClass)
	return theme.Badge(statusText, variant)
}

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

// getShipmentStatusIcon returns icon for shipment status
func getShipmentStatusIcon(status string) string {
	switch status {
	case miex.StatusPending:    return "📦"
	case "picked_up":    return "🚚"
	case "in_transit":   return "🚛"
	case "out_for_delivery": return "🏃"
	case "delivered":    return "✅"
	case "exception":    return "⚠️"
	default:             return "📋"
	}
}

// =====================================================
// INTEGRATION HELPERS
// =====================================================

// CreateLogisticsDemoPage creates a complete demo page using sample data
func CreateLogisticsDemoPage(theme mui.Theme) mi.H {
	// Use pure domain functions to create sample data
	service := mimo.NewLogisticsService()
	sampleShipments := mimo.SampleShipments()
	
	// Add data to service (simplified for demo)
	for _, shipment := range sampleShipments {
		// In real implementation, would use service methods
		_ = shipment
	}
	
	// Prepare dashboard data using pure domain functions
	dashboardData := mimo.PrepareDashboardData(service)
	
	// Create UI using presentation adapters
	return LogisticsDashboard(theme, dashboardData)
}

// IntegrateWithMainApp shows how to integrate with main application
func IntegrateWithMainApp(theme mui.Theme, logisticsService *mimo.LogisticsService) mi.H {
	// Get business data from pure domain service
	dashboardData := mimo.PrepareDashboardData(logisticsService)
	
	// Use presentation adapters to create UI
	return LogisticsDashboard(theme, dashboardData)
}
