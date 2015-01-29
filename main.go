package main

import (
	"flag"
	"fmt"
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

	return h
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// serveIndex accepts a file for upload
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	// TODO: accept file upload
}
