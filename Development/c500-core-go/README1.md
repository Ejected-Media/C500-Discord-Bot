That's a great request. In the main.go file we've been building, the most critical helper function—which sits outside the main block but is used directly inside it—is the template parsing function.
Without this function, we would have to manually list every single HTML file in our project, which isn't scalable.
Here is the source code for the parseTemplates function as used in your main.go file.
The Function: parseTemplates
This function resides at the bottom of your main.go file, outside of the main() function itself.
// parseTemplates recursively walks the given root directory looking for .html files
// and parses them into a single template set.
func parseTemplates(rootDir string) (*template.Template, error) {
	// 1. Start with an empty base template set.
	rootTmpl := template.New("")
	var files []string

	// 2. Walk through the directory tree starting from rootDir.
	// filepath.Walk finds files in subdirectories too (like templates/docs/).
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If we can't access a path, return the error to stop walking.
			return err
		}
		// 3. Filter logic:
		// If it's a file (not a directory) AND it ends with ".html", add it to our list.
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
			// Optional log to see what's being loaded during startup
			// log.Printf("Found template file: %s", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 4. Parse all found files into the root template set at once.
	// This step is crucial: parsing them together allows them to "see" each other,
	// enabling features like {{extends "base.html"}}.
	if len(files) > 0 {
		// We use the spread operator (...) to pass the slice of file paths as individual arguments.
		_, err = rootTmpl.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing template files: %w", err)
		}
	}

	// 5. Return the fully assembled template set.
	return rootTmpl, nil
}

How it's used within the main() block:
At the very beginning of your main() function, you call this helper function. It runs once on startup to prepare all your HTML views before the server starts accepting requests.
func main() {
	// ... (other setup like loading .env)

	// --- Template Initialization ---
	// Call the helper function defined below to load all HTML files from the "./templates" folder.
	tmpl, err := parseTemplates("./templates")
	if err != nil {
		log.Fatalf("Critical error: Failed to parse templates: %v", err)
	}
	log.Println("Templates parsed successfully.")

	// ... (The 'tmpl' variable is then passed down to your handlers)
    // storeHandler := handlers.NewStoreHandler(tmpl, productModel)
    // ...
}

