package main

import (
	"fmt"
	"os"
	"strings"
	"html/template"
    "net/http"
)

// data box to be sent to HTML template
type PageData struct {
	result      string
	error       string
}

func main() {
	http.HandleFunc("/", hompage)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func hompage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
// to load the template
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	// to execute the template with empty data
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{})
}	
// to load ascii art and render it
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//read value form
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	// validate input
	if text == "" || banner == "" {
		http.Error(w, "Bad Request - missing input", http.StatusBadRequest)
		return
	}
	// load banner file
	bannerPath := "banners/" + banner + ".txt"
	bannerData, err := os.ReadFile(bannerPath)
	if err != nil {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}
	// run ascii art logic
	result, err := generateAsciiArt(text, string(bannerData))
	if err != nil {
		http.Error(w, "Error generating ASCII art", http.StatusInternalServerError)
		return
	}
	// load template
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	// send result back to user
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, PageData{result: result})
}

// to generate ascii art from input text and banner data
func generateAsciiArt(text, bannerData string) (string, error) {
	lines := strings.Split(bannerData, "\n")
	var result strings.Builder

	for _, char := range text {
		index := (int(char) - 32) * 9
		if index < 0 || index >= len(lines) {
			return "", fmt.Errorf("character out of range")
		}
	}

	return result.String(), nil
}