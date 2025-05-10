package main

import (
	"ascii-art-web-stylize/banners"
	ascii "ascii-art-web-stylize/rendering"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type PageData struct {
	Text   string
	Banner string
	Result string
}

func main() {
	templates, err := banners.ParseTemplates()
	if err != nil {
		fmt.Println("error parsing banner templates:", err)
		return
	}
	index, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("error parsing index template: %v\n", err)
		return
	}
	errorTemplate, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		fmt.Printf("error parsing error template: %v\n", err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", Index(index, errorTemplate))
	mux.HandleFunc("POST /ascii-art", ASCIIArt(templates, index, errorTemplate))
	mux.HandleFunc("GET /styles.css", stylesHandler)

	fileServer := http.FileServer(http.Dir("static"))
	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/static/")
		if _, err := os.Stat("static/" + path); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			errorTemplate.Execute(w, "404 Page Not Found")
			return
		}

		http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
	})

	fmt.Println("server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
func stylesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "static/styles.css")
}
func Index(index *template.Template, errorTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			errorTemplate.Execute(w, "404 Page Not Found")
			return
		}
		index.Execute(w, PageData{})
	}
}
func ASCIIArt(templates map[string]*ascii.Template, index *template.Template, errorTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			errorTemplate.Execute(w, "405 Method Not Allowed")
			return
		}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorTemplate.Execute(w, "400 Bad Request - Invalid form data")
			return
		}
		text := r.FormValue("text")
		banner := r.FormValue("banner")
		if containsNonASCII(text) {
			w.WriteHeader(http.StatusBadRequest)
			errorTemplate.Execute(w, "400 Bad Request - Non-ASCII characters detected")
			return
		}
		tmpl, ok := templates[banner]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			errorTemplate.Execute(w, "400 Bad Request - Invalid banner")
			return
		}
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.TrimSpace(text)
		lines := strings.Split(text, "\n")
		var resultLines []string
		for _, line := range lines {
			if line == "" {
				continue
			}
			lineResult, err := tmpl.Execute(line)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorTemplate.Execute(w, "500 Internal Server Error - Art generation failed")
				return
			}
			resultLines = append(resultLines, lineResult)
		}
		data := PageData{
			Text:   text,
			Banner: banner,
			Result: strings.Join(resultLines, "\n"),
		}
		index.Execute(w, data)
	}
}
func containsNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}
