package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/munchiis/portfolio/templates/pages"
)

func getProjects() []pages.Project {
	return []pages.Project{
		{
			Title:        "React Component Library",
			Description:  "A TypeScript-based component library with Storybook documentation, converting Figma mockups into reusable UI components.",
			Technologies: []string{"React", "TypeScript", "Storybook", "Styled Components"},
			GithubURL:    "#",
		},
		{
			Title:        "VSCode Extension -- Saucy",
			Description:  "Extension that integrates gitlab merge request comments into your code editor.",
			Technologies: []string{"TypeScript", "VSCode API", "WebSockets"},
			GithubURL:    "https://github.com/Munchiis/gitlabExtension",
		},
		{
			Title:        "VSCode Extension -- TimeSplitters",
			Description:  "extension that helps developers track and visualize time spent on different Git branches, providing insights into work patterns and improving time management.",
			Technologies: []string{"TypeScript", "VSCode API", "WebSockets"},
			GithubURL:    "https://github.com/Munchiis/vscode-timeSplitters",
		},
		{
			Title:        "Developer Productivity Tools",
			Description:  "CLI tools built with TypeScript & Node.js to automate repetitive tasks and streamline Git workflows.",
			Technologies: []string{"Node.js", "TypeScript", "CLI"},
			GithubURL:    "#",
		},
	}
}

func main() {
	// Command line flags
	port := flag.Int("port", 8080, "Port to serve on")
	dev := flag.Bool("dev", false, "Run in development mode")
	build := flag.Bool("build", false, "Build static site")
	dir := flag.String("dir", "dist", "Directory to serve/build to")
	basePath := flag.String("base-path", "", "Base path for assets (e.g. /portfolio for GitHub Pages)")
	flag.Parse()

	// Build mode: generate static site
	if *build {
		generateStaticSite(*dir, *basePath)
		return
	}

	// Development mode: serve dynamic templ components
	if *dev {
		fmt.Println("Running in development mode with dynamic templates")
		serveDevelopmentMode(*port)
		return
	}

	// Production mode: serve static files
	serveStaticFiles(*port, *dir)
}

// generateStaticSite creates a static website from our templ components
// This function renders all templ components to HTML files
func generateStaticSite(outputDir string, basePath string) {
	fmt.Printf("Building static site to %s with base path %s\n", outputDir, basePath)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create static directory inside output directory
	staticDir := filepath.Join(outputDir, "static")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}

	// Copy static assets
	copyStaticAssets(outputDir)

	// Generate pages from templ components
	generatePage(filepath.Join(outputDir, "index.html"), pages.HomePage(), basePath)

	// Ensure the about directory exists
	aboutDir := filepath.Join(outputDir, "about")
	if err := os.MkdirAll(aboutDir, 0755); err != nil {
		log.Fatalf("Failed to create about directory: %v", err)
	}
	generatePage(filepath.Join(outputDir, "about", "index.html"), pages.AboutPage(), basePath)
	generatePage(filepath.Join(outputDir, "projects", "index.html"), pages.ProjectsPage(getProjects()), basePath)
	generatePage(filepath.Join(outputDir, "contact", "index.html"), pages.ContactPage(), basePath)
	generatePage(filepath.Join(outputDir, "thanks", "index.html"), pages.ThanksPage(), basePath)

	// Create .nojekyll file to prevent GitHub Pages from using Jekyll
	if basePath != "" {
		noJekyllPath := filepath.Join(outputDir, ".nojekyll")
		if err := os.WriteFile(noJekyllPath, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create .nojekyll file: %v", err)
		}
	}

	fmt.Println("Static site built successfully!")
}

// generatePage renders a templ component to an HTML file
// This is used during static site generation
func generatePage(filePath string, component templ.Component, basePath string) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// First render to a buffer so we can adjust paths if needed
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		log.Fatalf("Failed to render component to buffer: %v", err)
	}

	// Adjust paths for GitHub Pages if a base path is provided
	html := buf.String()
	if basePath != "" {
		html = adjustPathsForGitHubPages(html, basePath)
	}

	// Write the adjusted HTML to file
	if err := os.WriteFile(filePath, []byte(html), 0644); err != nil {
		log.Fatalf("Failed to write file %s: %v", filePath, err)
	}

	fmt.Printf("Generated %s\n", filePath)
}

// adjustPathsForGitHubPages modifies asset paths to work with GitHub Pages
func adjustPathsForGitHubPages(htmlContent, basePath string) string {
	// Only make these replacements when deploying to GitHub Pages
	if basePath != "" && basePath != "/" {
		// Adjust paths in HTML
		htmlContent = strings.ReplaceAll(htmlContent, `href="/static/`, `href="`+basePath+`/static/`)
		htmlContent = strings.ReplaceAll(htmlContent, `src="/static/`, `src="`+basePath+`/static/`)
		// Adjust navigation links
		htmlContent = strings.ReplaceAll(htmlContent, `href="/"`, `href="`+basePath+`/"`)
		htmlContent = strings.ReplaceAll(htmlContent, `href="/about"`, `href="`+basePath+`/about"`)
		htmlContent = strings.ReplaceAll(htmlContent, `href="/projects"`, `href="`+basePath+`/projects"`)
		htmlContent = strings.ReplaceAll(htmlContent, `href="/contact"`, `href="`+basePath+`/contact"`)
		htmlContent = strings.ReplaceAll(htmlContent, `href="/thanks"`, `href="`+basePath+`/thanks"`)
		htmlContent = strings.ReplaceAll(htmlContent,
			`value="https://munchiis.github.io/portfolio/thanks"`,
			`value="https://munchiis.github.io`+basePath+`/thanks"`)
	}
	return htmlContent
}

// copyStaticAssets copies static files to the output directory
// This preserves the directory structure inside the static directory
func copyStaticAssets(outputDir string) {
	// Walk through the static directory
	err := filepath.Walk("static", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, we'll create them as needed
		if info.IsDir() {
			return nil
		}

		// Read the file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}

		// Create the destination path
		// This keeps the same structure as the source, but inside outputDir
		relPath, err := filepath.Rel("static", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %v", path, err)
		}

		destPath := filepath.Join(outputDir, "static", relPath)

		// Ensure the destination directory exists
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", destDir, err)
		}

		// Write the file
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", destPath, err)
		}

		fmt.Printf("Copied %s to %s\n", path, destPath)
		return nil
	})

	if err != nil {
		log.Fatalf("Error copying static assets: %v", err)
	}
}

// serveDevelopmentMode serves the site with dynamic templ components
// This allows for real-time updates during development
func serveDevelopmentMode(port int) {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes to templ components
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/", "/index.html":
			pages.HomePage().Render(r.Context(), w)
		case "/about", "/about/index.html":
			pages.AboutPage().Render(r.Context(), w)
		case "/projects", "/projects/index.html":
			projects := getProjects()
			pages.ProjectsPage(projects).Render(r.Context(), w)
		case "/contact", "/contact/index.html":
			pages.ContactPage().Render(r.Context(), w)
		case "/thanks", "/thanks/index.html":
			pages.ThanksPage().Render(r.Context(), w)
		default:
			http.NotFound(w, r)
		}
	})

	// Start server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Development server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// serveStaticFiles serves the static site files from the specified directory
// This is used in production mode
func serveStaticFiles(port int, dir string) {
	// Create a file server handler for the static directory
	fs := http.FileServer(http.Dir(dir))

	// Add headers for better performance and security
	handler := addHeaders(fs)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Serving static files from %s at http://localhost%s\n", dir, addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// addHeaders adds performance and security headers to the response
func addHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cache static assets for better performance
		if isStaticAsset(r.URL.Path) {
			// Cache static assets for 1 year (immutable content)
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			// Don't cache HTML by default
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		}

		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// isStaticAsset determines if a path is a static asset that should be cached
func isStaticAsset(path string) bool {
	// List of file extensions that should be cached
	staticExtensions := []string{
		".css", ".js", ".jpg", ".jpeg", ".png",
		".gif", ".svg", ".woff", ".woff2", ".ttf",
		".webp", ".avif",
	}

	ext := filepath.Ext(path)
	for _, staticExt := range staticExtensions {
		if ext == staticExt {
			return true
		}
	}
	return false
}
