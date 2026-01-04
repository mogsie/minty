// Package main demonstrates basic minty usage
package main

import (
	"fmt"
	"net/http"
	"os"

	mi "github.com/ha1tch/minty"
)

// HomePage creates the home page template
func HomePage(title string, items []string) mi.H {
	return func(b *mi.Builder) mi.Node {
		return mi.Document(title,
			[]mi.Node{
				b.Link(mi.Rel("stylesheet"), mi.Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css")),
			},
			b.Body(mi.Class("container py-4"),
				b.Header(mi.Class("pb-3 mb-4 border-bottom"),
					b.H1(mi.Class("display-4"), title),
				),
				b.Main(
					b.Section(mi.Class("mb-4"),
						b.H2("Welcome to Minty"),
						b.P("A type-safe HTML generation library for Go web applications."),
					),
					b.Section(
						b.H3("Sample Items"),
						b.Ul(mi.Class("list-group"),
							mi.NewFragment(mi.Each(items, func(item string) mi.H {
								return func(b *mi.Builder) mi.Node {
									return b.Li(mi.Class("list-group-item"), item)
								}
							})...),
						),
					),
				),
				b.Footer(mi.Class("pt-3 mt-4 text-muted border-top"),
					b.Small("Built with Minty"),
				),
			),
		)(b)
	}
}

// ContactForm creates a contact form with HTMX
func ContactForm() mi.H {
	return func(b *mi.Builder) mi.Node {
		return b.Form(mi.Class("card p-4"),
			mi.HtmxPost("/contact"),
			mi.HtmxTarget("#result"),
			mi.HtmxSwap("innerHTML"),
			
			b.Div(mi.Class("mb-3"),
				b.Label(mi.Class("form-label"), mi.DataAttr("for", "name"), "Name"),
				b.Input(mi.Type("text"), mi.Name("name"), mi.ID("name"), 
					mi.Class("form-control"), mi.Required()),
			),
			b.Div(mi.Class("mb-3"),
				b.Label(mi.Class("form-label"), mi.DataAttr("for", "email"), "Email"),
				b.Input(mi.Type("email"), mi.Name("email"), mi.ID("email"), 
					mi.Class("form-control"), mi.Required()),
			),
			b.Div(mi.Class("mb-3"),
				b.Label(mi.Class("form-label"), mi.DataAttr("for", "message"), "Message"),
				b.Textarea(mi.Name("message"), mi.ID("message"), 
					mi.Class("form-control"), mi.Rows(4)),
			),
			b.Button(mi.Type("submit"), mi.Class("btn btn-primary"), "Submit"),
			b.Div(mi.ID("result"), mi.Class("mt-3")),
		)
	}
}

func main() {
	// Example 1: Render to string and print
	items := []string{"Type-safe HTML", "HTMX integration", "Theme support", "No templates needed"}
	html := mi.RenderToString(HomePage("Minty Demo", items))
	fmt.Println("Generated HTML length:", len(html), "bytes")
	
	// Example 2: Write to file
	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	
	if err := mi.Render(HomePage("Minty Demo", items), f); err != nil {
		fmt.Println("Error rendering:", err)
		return
	}
	fmt.Println("Wrote output.html")
	
	// Example 3: HTTP handler (commented out - uncomment to run server)
	/*
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		items := []string{"Dynamic content", "Server rendered", "Fast & efficient"}
		mi.Render(HomePage("Minty Server Demo", items), w)
	})
	
	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		mi.Render(ContactForm(), w)
	})
	
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
	*/
	
	_ = http.NewServeMux // Silence import
}
