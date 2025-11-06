package utils

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	wroteHeader bool
}

// WriteHeader captures the status code and put in struct for logging purposes.
func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int{
	return rw.status
}


// WriteHeader captures the status code before writing header to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	
	// Prevent multiple calls to WriteHeader, its for safety. A header should be written only once.
	//wrote header is set deafult to false before first call to WriteHeader
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

}


// middleware that logs the incoming HTTP requests and their responses.
func LoggerMiddleware(logger slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("panic recovered in logger middleware",
						slog.Any("error", err),
						slog.String("stacktrace", string(debug.Stack())),
				)
				}
			}()

			start := time.Now()
			// Wrap the ResponseWriter to capture the status code
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			// Log the request details after the response is served
			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.status),
				slog.Duration("durantion", time.Since(start)),	
			)
		}
		return http.HandlerFunc(fn)
	}
}

