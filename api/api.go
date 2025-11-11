package api

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"encurtador-go/utils"
)

// requestIDKey is the header key used to store the request ID.
const requestIDKey = "X-Request-ID"

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data any   `json:"data,omitempty"`
}

// RequestID is a middleware that sets a unique request ID for each incoming HTTP request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check if the request already has a request ID header
		RequestID := r.Header.Get(requestIDKey)
		if RequestID == "" {
			RequestID = xid.New().String()
		}

		// Set the request ID in the response header
		ctx = utils.SetRequestID(ctx, RequestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewHandler(db map[string]string) http.Handler {
	
	router := mux.NewRouter()

	// Logging middleware
	logger := handlers.LoggingHandler(os.Stdout, router)

	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	router.Use(RequestID)
	router.Use(utils.JsonMiddleware)

	router.HandleFunc("/api/shorten", handlePost(db)).Methods("POST")
	router.HandleFunc("/{code}", handleGet(db)).Methods("GET")

	return logger

}

// sendJSON sends a JSON response to the client with the given status code.
func sendJSON(w http.ResponseWriter, resp Response, statusCode int){
	// Set the content type to application/json
	data, err := json.Marshal(resp)
	
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(
			w,
			Response{Error: "internal server error"},
			http.StatusInternalServerError, 
		)
		return
	}
	// Write the response
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
	}
}

const caracteres = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortURL() string {
	const length = 8
	bytes := make([]byte, length)

	for i := range length {
		bytes[i] = caracteres[rand.IntN(len(caracteres))]
	}
	return string(bytes)
}

func handlePost(db map[string]string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var body PostBody
		// Decode the request body to URL
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Error: "invalid request body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(w, Response{Error: "invalid url"}, http.StatusBadRequest)
		}

		code := generateShortURL()
		// Store the URL in the database in memory
		db[code] = body.URL

		// Send the response
		sendJSON(w, Response{Data: code}, http.StatusCreated)
}}

func handleGet(db map[string]string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		// Get the code from the URL path
		code := mux.Vars(r)["code"]
		// Look up the original URL in the database
		originalURL, ok := db[code]
		// If the URL is not found, return a 404 error
		if !ok {
			sendJSON(w, Response{Error: "url not found"}, http.StatusNotFound)
		}

		http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)
}}
