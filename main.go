package main

import (
	"fmt"
	"os"
	"strings"
	"html/template"
  "net/http"

	"ascii-art-web/banner"
)

// data box to be sent to HTML template
type PageData struct {
	result      string
	error       string
}

// regestrying our to endpoints
func main() {
	http.HandleFunc("/", hompage)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

//  serves the homepage with the input form
func hompage(w http.ResponseWriter, r *http.Request) {
	// to reject any path other than "/"
	if r.URL.Path != "/" {
		http.Error(w, "404 - page not found", http.StatusNotFound)
			return
	}
	// to allow GET method only
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
			return
	}

	// to load the template
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "404 - Template not found", http.StatusNotFound)
		return
	}

	// for pages withn empty data
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{})
	}

// to load ascii art and render it
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// to allow only POST method
	if r.Method != http.MethodPost{
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}
	// read value form
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	// validate input
	if text == "" || banner == "" {
		http.Error(w, "400 - Bad Request - missing input", http.StatusBadRequest)
		return
	}

	// allow only the three valid banner
	validBanners :=  map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	if !validBanners[bannerName] {
		http.Error(w, "400 - Bad Request: Invalid banner name", http.StatusBadRequest)
		return
	}
	// load banner file
	chars, err := banner.Load(bannerName)
	if err != nil {
		http.Error(w, "404 - Banner File Not Found", http.StatusNotFound)
		return
	}
	// run ascii art logic
	result, err := generateAsciiArt(text, chars)
	if err != nil {
		http.Error(w, "500 - Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// load template
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "404 - Template not found", http.StatusNotFound)
		return
	}
	// send result back to user
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{result: result})
}

// to generate ascii art from input text and banner data
func generateAsciiArt(input string, chars [][]string) (string, error) {
	var output strings.Builder

	// replace literal \n with real newlines
	input = strings.ReplaceAll(input, '\n', "\n")
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		if line == "" {
			if i < len(lines)-1 {
				output.WriteString("\n")
			}
			continue
		}

		// validate every character is in the supported ASCII range
			for _, char := range line {
		index := (int(char) - 32)
		if index < 0 || index > 94 {
			return "", fmt.Errorf("character out of range")
	}
}

// build the ascii art row by row
for row := 0; row < 8; row++ {
	for _, ch := range line {
		output.WriteString(chars[int(ch)-32][row])
	}
			output.WriteString("\n")
		}
	}

	return output.string(), nil
}

