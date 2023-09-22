package middlewares

import (
	"main/session"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

var hmacSampleSecret = []byte("someSecret") // TODO: put this key in safe place and use proper secret
var Log *slog.Logger

// type for chaining
type Middleware func(http.HandlerFunc) http.HandlerFunc

func init() {
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// basically thisd is middleware chaining
func CompileMiddleware(h http.HandlerFunc, m []Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if r.URL.Path != "/reload" {
			Log.Info("request", "method", r.Method, "url", r.URL.Path, "duration", time.Since(start).Microseconds())
		}
	})
}

func VerifySession(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := session.ReadSigned(r)
		if err != nil {
			Log.Info("no session", "method", r.Method, "url", r.URL.Path, "error", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		h(w, r)
	}
}
