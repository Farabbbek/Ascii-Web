package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"unicode"
)

func containsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false

}

func asciiWebHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	font := r.FormValue("font")

	if containsCyrillic(text) {
		http.Error(w, "Bad request not supported cyrillic", http.StatusBadRequest)
		return

	}

	// Create path to font file in banner directory
	fontPath := filepath.Join("..", "banner", font+".txt")

	ascii := NewASCIIArt()
	if err := ascii.LoadFont(fontPath); err != nil {
		http.Error(w, "Error loading font", http.StatusInternalServerError)
		return
	}

	result := ascii.RenderText(text)
	fmt.Fprint(w, result)
}

func notFoundHandler(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Not Found")
}
func hostCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost:8080" {
			notFoundHandler(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Serve static files from current directory with host check
	http.Handle("/", hostCheckMiddleware(http.FileServer(http.Dir("."))))

	// Update banner path to point to parent directory with host check
	http.Handle("/banner/", hostCheckMiddleware(http.StripPrefix("/banner/", http.FileServer(http.Dir("../banner")))))

	// Wrap asciiWebHandler with host check.
	http.HandleFunc("/ascii-web/", func(w http.ResponseWriter, r *http.Request) {
		// Check host before calling handler.
		if r.Host != "localhost:8080" {
			notFoundHandler(w)
			return
		}
		asciiWebHandler(w, r)
	})

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
