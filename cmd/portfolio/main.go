package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/munchiis/go-portfolio/templates/pages"
)

func main() {
	// Determine if we're building a static site
	if len(os.Args) > 1 && os.Args[1] == "build" {
		generateStaticSite("dist")
		return
	}

	// Otherwise serve dynamic site
	serveDevMode(8080)
}

// serveDevMode serves the site with dynamic templ components
func serveDevMode(port int) {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes to templ components
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/", "/index.html":
			pages.HomePage().Render(r.Context(), w)
		case "/about", "/about/index.html":
			pages.AboutPage().Render(r.Context(), w)
		default:
			http.NotFound(w, r)
		}
	})

	// Start server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Development server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// generateStaticSite creates a static website from our templ components
func generateStaticSite(outputDir string) {
	fmt.Printf("Building static site to %s\n", outputDir)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate pages from templ components
	generatePage(filepath.Join(outputDir, "index.html"), pages.HomePage())
	generatePage(filepath.Join(outputDir, "about", "index.html"), pages.AboutPage())

	fmt.Println("Static site built successfully!")
}

// generatePage renders a templ component to an HTML file
func generatePage(filePath string, component templ.Component) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Render component to file
	err = component.Render(context.Background(), file)
	if err != nil {
		log.Fatalf("Failed to render component to %s: %v", filePath, err)
	}

	fmt.Printf("Generated %s\n", filePath)
}
