package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"asciiart/asciiart" // Ensure this import path matches your project structure
)

func main() {
	http.HandleFunc("/", homeHandler)              // Handle requests to the home page
	http.HandleFunc("/ascii-art", asciiArtHandler) // Handle form submissions

	port := ":8080"
	fmt.Printf("Server is running at http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil) // Start the server on port 8080
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html") // Load the HTML template
	if err != nil {                                          // If there's an error loading the template, handle it
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to load template: %v\n", err)
		return
	}
	err = tmpl.Execute(w, nil) // Render the template and send it to the browser
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to execute template: %v\n", err)
	}
}
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderErrorPage(w, "Bad Request", "Invalid request method", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		renderErrorPage(w, "Bad Request", "Text or banner is empty", http.StatusBadRequest)
		return
	}

	// Check for invalid characters, allow \r (13) and \n (10)
	for _, char := range text {
		if int(char) != 13 && int(char) != 10 && (int(char) < 32 || int(char) > 126) {
			renderErrorPage(w, "Bad Request", fmt.Sprintf("Invalid character in input: %v", char), http.StatusBadRequest)
			return
		}
	}

	// Replace \n with an actual newline character
	text = strings.ReplaceAll(text, "\\n", "\n")

	asciiArt, err := generateAsciiArt(text, banner)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", fmt.Sprintf("Failed to generate ASCII art: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderErrorPage(w, "Internal Server Error", fmt.Sprintf("Failed to load template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		AsciiArt string
		Text     string
		Banner   string
	}{
		AsciiArt: asciiArt, // The generated ASCII art
		Text:     text,
		Banner:   banner,
	}

	err = tmpl.Execute(w, data) // Render the template with the ASCII art
	if err != nil {
		renderErrorPage(w, "Internal Server Error", fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}
}

// renderErrorPage renders the custom error page with a specific message and status code
func renderErrorPage(w http.ResponseWriter, title, message string, statusCode int) {
	w.WriteHeader(statusCode)
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to load error template: %v\n", err)
		return
	}
	data := struct {
		Title   string
		Message string
	}{
		Title:   title,
		Message: message,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to execute error template: %v\n", err)
	}
}

func generateAsciiArt(text, banner string) (string, error) {
	bannerFile := filepath.Join("banners", fmt.Sprintf("%s.txt", banner)) // Determine the file path for the selected banner
	modifiedInput := ModifyString(text)                                   // Clean up the input string
	asciiArt, err := asciiart.AsciiTable(modifiedInput, bannerFile)       // Generate ASCII art
	if err != nil {
		fmt.Printf("Failed to generate ASCII art: %v\n", err)
		return "", err
	}
	return asciiArt, nil
}

func ModifyString(input string) string {
	// Remove carriage returns and replace newlines with \n
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r", ""), "\n", "\\n")
}
