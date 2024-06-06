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
	if r.Method != http.MethodPost { // Ensure the request method is POST
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println("Invalid request method")
		return
	}

	text := r.FormValue("text")     // Get the text from the form
	banner := r.FormValue("banner") // Get the banner choice from the form

	if text == "" || banner == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println("Text or banner is empty")
		return
	}

	// Check for invalid characters
	for _, char := range text {
		if int(char) < 32 || int(char) > 126 {
			http.Error(w, "Error 400 Bad Request: Invalid character in input", http.StatusBadRequest)
			fmt.Printf("Invalid character in input: %v\n", char)
			return
		}
	}

	// Clean up the input string
	modifiedText := ModifyString(text)

	asciiArt, err := generateAsciiArt(modifiedText, banner) // Generate ASCII art
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to generate ASCII art: %v\n", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to load template: %v\n", err)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Failed to execute template: %v\n", err)
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
