package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ascii-art-web/banner" // keep your existing import
)

// data box to be sent to HTML template
type PageData struct {
	Result string
	Error  string
}

func main() {
	// Creates a new ServeMux to register and manage routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/ascii-art", asciiArtHandler)

	// Define the server with timeouts to prevent resource exhaustion
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run server in a goroutine so main thread remains free to listen for shutdown signals
	go func() {
		fmt.Println("Server starting on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("\nServer is shutting down...")

	// Give the server 30 seconds to finish handling active requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server stopped")
}

// serves the homepage with the input form
func homepage(w http.ResponseWriter, r *http.Request) {
	// Reject any other routes that are not exactly "/"
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	// This handler only accepts GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 - Template not found", http.StatusNotFound)
		return
	}

	// Execute template with empty data
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{})
}

// handles ASCII art generation and rendering
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// This handler only accepts POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 - Bad Request: Unable to parse form", http.StatusBadRequest)
		return
	}

	// Read values from form
	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	// Validate input
	if text == "" {
		http.Error(w, "400 - Bad Request: Text input is required", http.StatusBadRequest)
		return
	}
	if bannerName == "" {
		http.Error(w, "400 - Bad Request: Banner selection is required", http.StatusBadRequest)
		return
	}

	// Allow only the three valid banner types
	validBanners := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	if !validBanners[bannerName] {
		http.Error(w, "400 - Bad Request: Invalid banner name. Use 'standard', 'shadow', or 'thinkertoy'", http.StatusBadRequest)
		return
	}

	// Load banner file
	chars, err := banner.Load(bannerName)
	if err != nil {
		http.Error(w, "404 - Banner File Not Found", http.StatusNotFound)
		return
	}

	// Generate ASCII art
	result, err := generateAsciiArt(text, chars)
	if err != nil {
		http.Error(w, "500 - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Load template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 - Template not found", http.StatusNotFound)
		return
	}

	// Send result back to user
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{Result: result})
}

// generates ASCII art from input text and banner data
func generateAsciiArt(input string, chars [][]string) (string, error) {
	var output strings.Builder

	// Replace literal \n with real newlines
	input = strings.ReplaceAll(input, `\n`, "\n")
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		// Handle empty lines (paragraph breaks)
		if line == "" {
			if i < len(lines)-1 {
				output.WriteString("\n")
			}
			continue
		}

		// Validate every character is in the supported ASCII range (32-126)
		for _, ch := range line {
			if ch < ' ' || ch > '~' {
				return "", fmt.Errorf("character '%c' (ASCII %d) is not supported. Only printable ASCII characters (32-126) are allowed", ch, ch)
			}
		}

		// Build the ASCII art row by row (each character is 8 rows tall)
		for row := 0; row < 8; row++ {
			for _, ch := range line {
				index := int(ch) - 32
				// Safety check for index bounds
				if index < 0 || index >= len(chars) {
					return "", fmt.Errorf("character '%c' (position %d) is out of range", ch, index+32)
				}
				if row < len(chars[index]) {
					output.WriteString(chars[index][row])
				}
			}
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}
