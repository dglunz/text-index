package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(0)

	runServer(os.Args[1:])
}

const (
	DefaultBindAddress = ":1134"
	//SplitOn = " ,-,--"
	//DeleteOn = regex./\W|_/
)

func runServer(args []string) {
	fmt.Println(args)
	// Parse command line flags.
	fs := flag.NewFlagSet("dex", flag.ExitOnError)
	p := fs.String("p", DefaultBindAddress, "bind address")
	fs.Parse(args)

	// Initialize handler.
	h := NewHandler()

	// Start HTTP handler.
	log.Printf("Listening on http://localhost%s", *p)
	log.SetFlags(log.LstdFlags)
	http.ListenAndServe(*p, h)
}

// Handler represents the HTTP handler.
type Handler struct {
	ind map[string]string
	mux *mux.Router
}

// NewHandler returns a new instance of Handler
func NewHandler() *Handler {
	// Initialize handler.
	h := &Handler{
		mux: mux.NewRouter(),
	}

	// Setup request multiplexer
	h.mux.HandleFunc("/index", h.serveIndex).Methods("POST")
	h.mux.HandleFunc("/query", h.serveQuery).Methods("POST")

	return h
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// serveIndex accepts a file for upload
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Get the multipart reader for the request.
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mapping := bleve.NewIndexMapping()
	index, _ := bleve.New("classics.bleve")
	index.Index()

	fmt.Println(reader)

	// Display success message.
	fmt.Println("Upload successful")
}

// serveQuery accepts a query string and returns a JSON index of {name:line_number:word_number}
func (h *Handler) serveQuery(w http.ResponseWriter, r *http.Request) {
	// Get the query param from the request
	q := r.FormValue("query")
	fmt.Println(q)

	// Execute the query.
	index, _ := bleve.Open("classics.bleve")
	query := bleve.NewQueryStringQuery(q)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)

	// Write the results.
	jw, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(jw)
}
