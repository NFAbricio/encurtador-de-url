package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"encurtador-go/utils"
)

const requestIDKey = "X-Request-ID"

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

func NewHandler() http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	router := mux.NewRouter()

	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	router.Use(RequestID)
	router.Use(utils.LoggerMiddleware(*logger))

}