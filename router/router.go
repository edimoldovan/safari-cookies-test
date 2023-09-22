package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

var Routes = []Route{}

type Route struct {
	method  string // use http.MethodGet, http.MethodPost, etc.
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type CtxKey struct{}

func CreateRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(CtxKey{}).([]string)
	return fields[index]
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range Routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), CtxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
	// trying this as a replacement hack shit for a 404
	// http.Redirect(w, r, fmt.Sprintf("/?redirect=%s", r.URL.Path), http.StatusFound)
}
