package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
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
	// Parse the multipart form in the request
	//get the multipart reader for the request.
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy each part to destination.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		// If part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}

		// Copy file into tmp directory
		dst, err := os.Create("./tmp/" + part.FileName())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Display success message.
	fmt.Println("Upload successful")
}

// serveQuery accepts a query string and returns a JSON index of {name:line_number:word_number}
func (h *Handler) serveQuery(w http.ResponseWriter, r *http.Request) {
	// Parse the statement.
	//stmt, err := pieql.NewParser(r.Body).Parse()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	// Get the query param from the request
	q := r.FormValue("query")
	fmt.Println(q)

	// Execute the statement.
	//res, err := h.db.Execute(stmt)
	//if err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}

	// Write the results.
	jw, err := json.Marshal(w)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(jw)
}
