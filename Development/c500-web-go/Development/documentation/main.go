package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// IMPORTANT: Replace 'c500-web-go' with whatever module name you defined in your go.mod file.
	// If you haven't run 'go mod init' yet, do so now: e.g., 'go mod init c500-web-go'
	"c500-web-go/handlers"
)

func main() {
	// 1. Parse Templates
	// We need a robust way to find all .html files in the templates directory,
	// including subdirectories like templates/docs/.
	tmpl, err := parseTemplates("./templates")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	log.Println("Templates parsed successfully.")

	// 2. Initialize Handlers
	// Inject the parsed templates into our documentation handler.
	docsHandler := handlers.NewDocsHandler(tmpl)

	// 3. Create a new router (ServeMux)
	mux := http.NewServeMux()

	// --- Static File Serving ---
	// This tells Go: if a request starts with "/static/", strip that prefix off,
	// and look for the remaining file path inside the local "./static" directory.
	// e.g., Request for "/static/css/style.css" -> serves file "./static/css/style.css"
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// --- Application Routes ---
	// Map URLs to handler functions.
	// Note: In standard Go http, ending a pattern with "/" matches everything under it.
	// So "/docs/" matches "/docs/", "/docs/buyer", etc. We handle specific paths inside the handler if needed, 
    // or define exact matches. For simplicity here, we define exact matches.

	// Main Docs Hub
	mux.HandleFunc("/docs/", docsHandler.Index)

	// Specific Manual Pages
	mux.HandleFunc("/docs/buyer", docsHandler.Buyer)
	mux.HandleFunc("/docs/seller", docsHandler.Seller)
	mux.HandleFunc("/docs/developer", docsHandler.Developer)
	mux.HandleFunc("/docs/admin", docsHandler.Admin)

	// Add a temporary root route just so localhost:8080 doesn't 404 immediately.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/docs/", http.StatusFound)
	})


	// 4. Start the Server
	port := ":8080"
	log.Printf("Starting server on http://localhost%s ...\n", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

// parseTemplates recursively walks the given root directory looking for .html files
// and parses them into a single template set.
func parseTemplates(rootDir string) (*template.Template, error) {
	// Start with a base template.
	rootTmpl := template.New("")
	var files []string

	// Walk through the directory tree.
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If it's not a directory and has an .html extension, add it to our list.
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
			log.Printf("Found template file: %s", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Parse all found files into the root template.
	// This allows templates to reference each other (e.g., {{extends "base.html"}}).
	if len(files) > 0 {
		rootTmpl, err = rootTmpl.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing template files: %w", err)
		}
	}

	return rootTmpl, nil
}
