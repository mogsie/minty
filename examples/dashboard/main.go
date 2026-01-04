// Package main demonstrates a self-updating dashboard with HTMX
package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	mi "github.com/ha1tch/minty"
)

// Metrics simulates live system metrics
type Metrics struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	NetworkIn   float64
	NetworkOut  float64
	Requests    int
	Errors      int
	Latency     float64
}

// SalesData for bar chart
type SalesData struct {
	Label string
	Value float64
}

var baseMetrics = Metrics{
	CPUUsage:    45.0,
	MemoryUsage: 62.0,
	DiskUsage:   78.0,
	NetworkIn:   150.0,
	NetworkOut:  89.0,
	Requests:    12450,
	Errors:      23,
	Latency:     45.0,
}

func getMetrics() Metrics {
	// Simulate fluctuating metrics
	return Metrics{
		CPUUsage:    clamp(baseMetrics.CPUUsage+randDelta(15), 0, 100),
		MemoryUsage: clamp(baseMetrics.MemoryUsage+randDelta(8), 0, 100),
		DiskUsage:   clamp(baseMetrics.DiskUsage+randDelta(2), 0, 100),
		NetworkIn:   clamp(baseMetrics.NetworkIn+randDelta(50), 0, 500),
		NetworkOut:  clamp(baseMetrics.NetworkOut+randDelta(30), 0, 300),
		Requests:    baseMetrics.Requests + rand.Intn(100),
		Errors:      baseMetrics.Errors + rand.Intn(5),
		Latency:     clamp(baseMetrics.Latency+randDelta(20), 5, 200),
	}
}

func getSalesData() []SalesData {
	base := []float64{85, 92, 78, 95, 88, 72, 98}
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	data := make([]SalesData, 7)
	for i, day := range days {
		data[i] = SalesData{
			Label: day,
			Value: clamp(base[i]+randDelta(15), 0, 100),
		}
	}
	return data
}

func randDelta(max float64) float64 {
	return (rand.Float64() - 0.5) * 2 * max
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// =====================================================
// COMPONENTS
// =====================================================

// CircularGauge renders an SVG circular gauge
func CircularGauge(id, label string, value float64, maxValue float64, color string) mi.H {
	return func(b *mi.Builder) mi.Node {
		percentage := (value / maxValue) * 100
		circumference := 2 * math.Pi * 40 // radius = 40
		dashOffset := circumference * (1 - percentage/100)

		return b.Div(mi.Class("gauge-container"), mi.ID(id),
			b.Svg(
				mi.Attr("width", "120"),
				mi.Attr("height", "120"),
				mi.Attr("viewBox", "0 0 100 100"),
				// Background circle
				b.Circle(
					mi.Attr("cx", "50"),
					mi.Attr("cy", "50"),
					mi.Attr("r", "40"),
					mi.Attr("fill", "none"),
					mi.Attr("stroke", "#e0e0e0"),
					mi.Attr("stroke-width", "8"),
				),
				// Progress circle
				b.Circle(
					mi.Attr("cx", "50"),
					mi.Attr("cy", "50"),
					mi.Attr("r", "40"),
					mi.Attr("fill", "none"),
					mi.Attr("stroke", color),
					mi.Attr("stroke-width", "8"),
					mi.Attr("stroke-linecap", "round"),
					mi.Attr("stroke-dasharray", fmt.Sprintf("%.1f", circumference)),
					mi.Attr("stroke-dashoffset", fmt.Sprintf("%.1f", dashOffset)),
					mi.Attr("transform", "rotate(-90 50 50)"),
					mi.Style("transition: stroke-dashoffset 0.5s ease"),
				),
				// Value text
				b.SvgText(
					mi.Attr("x", "50"),
					mi.Attr("y", "50"),
					mi.Attr("text-anchor", "middle"),
					mi.Attr("dy", "0.3em"),
					mi.Attr("font-size", "16"),
					mi.Attr("font-weight", "bold"),
					mi.Attr("fill", "#333"),
					fmt.Sprintf("%.0f%%", percentage),
				),
			),
			b.Div(mi.Class("gauge-label"), label),
		)
	}
}

// BarChart renders an SVG bar chart
func BarChart(id string, data []SalesData) mi.H {
	return func(b *mi.Builder) mi.Node {
		barWidth := 30
		gap := 15
		chartHeight := 150
		chartWidth := len(data)*(barWidth+gap) + gap

		bars := make([]mi.Node, 0, len(data)*2)
		for i, d := range data {
			x := gap + i*(barWidth+gap)
			barHeight := int(d.Value / 100 * float64(chartHeight-30))
			y := chartHeight - barHeight - 20

			// Bar
			bars = append(bars, b.Rect(
				mi.Attr("x", fmt.Sprintf("%d", x)),
				mi.Attr("y", fmt.Sprintf("%d", y)),
				mi.Attr("width", fmt.Sprintf("%d", barWidth)),
				mi.Attr("height", fmt.Sprintf("%d", barHeight)),
				mi.Attr("fill", getBarColor(d.Value)),
				mi.Attr("rx", "3"),
				mi.Style("transition: all 0.5s ease"),
			))
			// Label
			bars = append(bars, b.SvgText(
				mi.Attr("x", fmt.Sprintf("%d", x+barWidth/2)),
				mi.Attr("y", fmt.Sprintf("%d", chartHeight-5)),
				mi.Attr("text-anchor", "middle"),
				mi.Attr("font-size", "11"),
				mi.Attr("fill", "#666"),
				d.Label,
			))
		}

		return b.Div(mi.Class("chart-container"), mi.ID(id),
			b.Svg(
				mi.Attr("width", fmt.Sprintf("%d", chartWidth)),
				mi.Attr("height", fmt.Sprintf("%d", chartHeight)),
				mi.Attr("viewBox", fmt.Sprintf("0 0 %d %d", chartWidth, chartHeight)),
				mi.NewFragment(bars...),
			),
		)
	}
}

func getBarColor(value float64) string {
	if value >= 90 {
		return "#22c55e" // green
	} else if value >= 70 {
		return "#3b82f6" // blue
	} else if value >= 50 {
		return "#f59e0b" // amber
	}
	return "#ef4444" // red
}

// StatCard renders a statistic card
func StatCard(id, title, value, subtitle, icon, color string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("stat-card"), mi.ID(id),
			b.Div(mi.Class("stat-icon"), mi.Style("background-color: "+color), icon),
			b.Div(mi.Class("stat-content"),
				b.Div(mi.Class("stat-title"), title),
				b.Div(mi.Class("stat-value"), value),
				b.Div(mi.Class("stat-subtitle"), subtitle),
			),
		)
	}
}

// GaugePanel renders the gauge section (for HTMX updates)
func GaugePanel(m Metrics) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("gauge-panel"), mi.ID("gauge-panel"),
			mi.HtmxGet("/api/gauges"),
			mi.HtmxTrigger("every 2s"),
			mi.HtmxSwap("outerHTML"),
			CircularGauge("cpu-gauge", "CPU", m.CPUUsage, 100, getGaugeColor(m.CPUUsage))(b),
			CircularGauge("mem-gauge", "Memory", m.MemoryUsage, 100, getGaugeColor(m.MemoryUsage))(b),
			CircularGauge("disk-gauge", "Disk", m.DiskUsage, 100, getGaugeColor(m.DiskUsage))(b),
		)
	}
}

func getGaugeColor(value float64) string {
	if value < 50 {
		return "#22c55e" // green
	} else if value < 75 {
		return "#f59e0b" // amber
	}
	return "#ef4444" // red
}

// StatsPanel renders the stats cards (for HTMX updates)
func StatsPanel(m Metrics) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("stats-panel"), mi.ID("stats-panel"),
			mi.HtmxGet("/api/stats"),
			mi.HtmxTrigger("every 3s"),
			mi.HtmxSwap("outerHTML"),
			StatCard("req-stat", "Requests", fmt.Sprintf("%d", m.Requests), "Total today", "↗", "#3b82f6")(b),
			StatCard("err-stat", "Errors", fmt.Sprintf("%d", m.Errors), "Last hour", "⚠", "#ef4444")(b),
			StatCard("lat-stat", "Latency", fmt.Sprintf("%.0fms", m.Latency), "Average", "⏱", "#8b5cf6")(b),
			StatCard("net-stat", "Network", fmt.Sprintf("%.0f MB/s", m.NetworkIn), "Inbound", "↓", "#22c55e")(b),
		)
	}
}

// ChartPanel renders the bar chart (for HTMX updates)
func ChartPanel(data []SalesData) mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Div(mi.Class("chart-panel"), mi.ID("chart-panel"),
			mi.HtmxGet("/api/chart"),
			mi.HtmxTrigger("every 4s"),
			mi.HtmxSwap("outerHTML"),
			b.H3("Weekly Performance"),
			BarChart("sales-chart", data)(b),
		)
	}
}

// Dashboard renders the full dashboard
func Dashboard() mi.H {
	return func(b *mi.Builder) mi.Node {
		m := getMetrics()
		sales := getSalesData()

		return mi.NewFragment(
			mi.Raw("<!DOCTYPE html>"),
			b.Html(mi.Lang("en"),
				b.Head(
					b.Title("Minty Dashboard"),
					b.Meta(mi.Charset("UTF-8")),
					b.Meta(mi.Name("viewport"), mi.Content("width=device-width, initial-scale=1")),
					b.Script(mi.Src("https://unpkg.com/htmx.org@1.9.10")),
					b.Style(mi.Raw(dashboardCSS)),
				),
				b.Body(
					b.Div(mi.Class("dashboard"),
						b.Header(mi.Class("dashboard-header"),
							b.H1("System Dashboard"),
							b.Div(mi.Class("header-time"), mi.ID("clock"),
								mi.HtmxGet("/api/time"),
								mi.HtmxTrigger("every 1s"),
								mi.HtmxSwap("innerHTML"),
								time.Now().Format("15:04:05"),
							),
						),
						b.Main(mi.Class("dashboard-content"),
							b.Section(mi.Class("section"),
								b.H2("System Health"),
								GaugePanel(m)(b),
							),
							b.Section(mi.Class("section"),
								b.H2("Key Metrics"),
								StatsPanel(m)(b),
							),
							b.Section(mi.Class("section section-wide"),
								ChartPanel(sales)(b),
							),
						),
						b.Footer(mi.Class("dashboard-footer"),
							b.Small("Built with Minty + HTMX • Auto-refreshing every few seconds"),
						),
					),
				),
			),
		)
	}
}

const dashboardCSS = `
* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
	min-height: 100vh;
	color: #fff;
}

.dashboard {
	max-width: 1200px;
	margin: 0 auto;
	padding: 20px;
}

.dashboard-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20px 0;
	border-bottom: 1px solid rgba(255,255,255,0.1);
	margin-bottom: 30px;
}

.dashboard-header h1 {
	font-size: 1.8rem;
	font-weight: 600;
}

.header-time {
	font-size: 1.5rem;
	font-family: monospace;
	color: #22c55e;
}

.dashboard-content {
	display: grid;
	gap: 30px;
}

.section {
	background: rgba(255,255,255,0.05);
	border-radius: 12px;
	padding: 20px;
}

.section h2 {
	font-size: 1.1rem;
	font-weight: 500;
	margin-bottom: 20px;
	color: rgba(255,255,255,0.7);
}

.section h3 {
	font-size: 1rem;
	font-weight: 500;
	margin-bottom: 15px;
	color: rgba(255,255,255,0.7);
}

.gauge-panel {
	display: flex;
	justify-content: space-around;
	flex-wrap: wrap;
	gap: 20px;
}

.gauge-container {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.gauge-container svg {
	filter: drop-shadow(0 4px 6px rgba(0,0,0,0.3));
}

.gauge-container circle {
	transition: stroke-dashoffset 0.5s ease, stroke 0.5s ease;
}

.gauge-label {
	margin-top: 10px;
	font-size: 0.9rem;
	color: rgba(255,255,255,0.7);
}

.stats-panel {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 15px;
}

.stat-card {
	display: flex;
	align-items: center;
	gap: 15px;
	background: rgba(255,255,255,0.05);
	padding: 15px;
	border-radius: 8px;
	transition: transform 0.2s ease, background 0.2s ease;
}

.stat-card:hover {
	transform: translateY(-2px);
	background: rgba(255,255,255,0.08);
}

.stat-icon {
	width: 45px;
	height: 45px;
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 1.2rem;
}

.stat-content {
	flex: 1;
}

.stat-title {
	font-size: 0.8rem;
	color: rgba(255,255,255,0.5);
}

.stat-value {
	font-size: 1.4rem;
	font-weight: 600;
}

.stat-subtitle {
	font-size: 0.75rem;
	color: rgba(255,255,255,0.4);
}

.chart-panel {
	overflow-x: auto;
}

.chart-container {
	display: flex;
	justify-content: center;
	padding: 10px 0;
}

.chart-container svg {
	filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
}

.section-wide {
	grid-column: 1 / -1;
}

.dashboard-footer {
	text-align: center;
	padding: 30px 0;
	color: rgba(255,255,255,0.4);
}

@media (max-width: 768px) {
	.dashboard-header {
		flex-direction: column;
		gap: 10px;
		text-align: center;
	}
	
	.gauge-panel {
		flex-direction: column;
		align-items: center;
	}
}
`

func main() {
	rand.Seed(time.Now().UnixNano())

	// Main dashboard page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		mi.Render(Dashboard(), w)
	})

	// HTMX endpoints for partial updates
	http.HandleFunc("/api/gauges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		mi.Render(GaugePanel(getMetrics()), w)
	})

	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		mi.Render(StatsPanel(getMetrics()), w)
	})

	http.HandleFunc("/api/chart", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		mi.Render(ChartPanel(getSalesData()), w)
	})

	http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, time.Now().Format("15:04:05"))
	})

	fmt.Println("Dashboard running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
