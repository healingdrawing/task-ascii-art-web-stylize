package main

import (
	ascii_art "ascii-art-web/ascii-art"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*
var fileBox embed.FS

// Data is a struct that will be sent as a respond
type Data struct {
	Output    string
	ErrorNum  int
	ErrorText string
	Warning   string
}

var temp = template.Must(template.ParseFS(fileBox, "templates/*.html"))

// returns a struct filled by error of "erNum" number description
func errorFormatter(erNum int) Data {
	return Data{"", erNum, http.StatusText(erNum), ""}
}

// Handles errors . Writes headers into http.ResponseWriter, according to template structure.
func errorHandler(w http.ResponseWriter, req *http.Request, d Data) {
	w.WriteHeader(d.ErrorNum)

	err := temp.ExecuteTemplate(w, "error.html", d)
	if err != nil {
		log.Printf("ERROR %d. %v\n", http.StatusNotFound, err)
		errorText := fmt.Sprintf("ERROR %d. %s", http.StatusNotFound, http.StatusText(http.StatusNotFound))
		http.Error(w, errorText, http.StatusNotFound)
		http.StatusText(http.StatusNotFound)
		return
	}
}

// Handles requests, and errors 404 405 500 for "/" pattern
func indexHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. Internal Server Error\n", http.StatusInternalServerError)
			errorHandler(w, r, errorFormatter(http.StatusInternalServerError)) // 500 ERROR
		}
	}()

	if r.URL.Path != "/" {
		log.Printf("ERROR %d. r.URL.Path = %s != \"/\"\n", http.StatusNotFound, r.URL.Path)
		errorHandler(w, r, errorFormatter(http.StatusNotFound)) // 404 ERROR
		return
	}
	if r.Method == "GET" {
		err := temp.ExecuteTemplate(w, "index.html", Data{})
		if err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusNotFound, err)
			errorHandler(w, r, errorFormatter(http.StatusNotFound))
			return
		}
	} else { // not a GET method case
		log.Printf("ERROR %d. %s\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method \"%s\" is inappropriate for the URL \"%s\"\n", r.Method, r.URL.Path))
		errorHandler(w, r, errorFormatter(http.StatusMethodNotAllowed)) // 405 error
		return
	}
}

// Handles requests, and errors 400 404 405 500 for "/ascii-art" pattern
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
			errorHandler(w, r, errorFormatter(http.StatusInternalServerError)) // 500 ERROR
		}
	}()

	if r.URL.Path != "/ascii-art" {
		log.Printf("ERROR %d. r.URL.Path = %s != \"/ascii-art\"\n", http.StatusNotFound, r.URL.Path)
		errorHandler(w, r, errorFormatter(http.StatusNotFound)) // 404 ERROR
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("ERROR %d. ParseForm() err: %v\n", http.StatusBadRequest, err)
		errorHandler(w, r, errorFormatter(http.StatusBadRequest)) // 400 ERROR
		return
	}

	input := r.FormValue("input")
	// bannerName := r.FormValue("banner") // will work too
	bannerName := r.Form.Get("banner")
	banner, ok := ascii_art.BannerFiles[bannerName]

	if !ok {
		log.Printf("ERROR %d. Wrong banner name \"%s\"\n", http.StatusBadRequest, bannerName)
		errorHandler(w, r, errorFormatter(http.StatusBadRequest)) // 400 ERROR
		return
	}
	art, warning, err := ascii_art.AsciiToString(input, banner)
	if strings.HasPrefix(art, ascii_art.HeadString) && err == nil {
		log.Println("minor ERROR. Found not supported characters")
	} else if err != nil {
		log.Printf("ERROR %d. %v\n", http.StatusNotFound, err)
		errorHandler(w, r, errorFormatter(http.StatusNotFound)) // 404 ERROR
		return
	}

	// fmt.Print(art) // uncomment to print art to console for every request
	if r.Method == "POST" {
		err := temp.ExecuteTemplate(w, "index.html", Data{art, http.StatusOK, http.StatusText(http.StatusOK), warning})
		if err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusNotFound, err)
			errorHandler(w, r, errorFormatter(http.StatusNotFound))
			return
		}
	} else { // not GET method case
		log.Printf("ERROR %d. %v\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method %s is inappropriate for the URL %s", r.Method, r.URL.Path))
		errorHandler(w, r, errorFormatter(http.StatusMethodNotAllowed)) // 405 error
		return
	}
}

// Starts http server with handlers
func startServer() {

	static, err := fs.Sub(fileBox, "templates/static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(static))

	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	log.Print("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
