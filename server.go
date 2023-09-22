package main

import (
	"embed"
	"html/template"
	"log"
	"main/handlers"
	"main/middlewares"
	"main/router"
	"net/http"
)

//go:embed templates
var embededTemplates embed.FS

// public HTML route middleware stack
var public = []middlewares.Middleware{
	middlewares.Logger,
}

// private JSON route middleware stack
var private = []middlewares.Middleware{
	middlewares.VerifySession,
	middlewares.Logger,
}

func init() {
	router.Routes = []router.Route{
		router.CreateRoute(http.MethodGet, "/", middlewares.CompileMiddleware(handlers.Home, public)),
		router.CreateRoute(http.MethodGet, "/set-cookie", middlewares.CompileMiddleware(handlers.SetCookie, public)),
		router.CreateRoute(http.MethodGet, "/private/css/type", middlewares.CompileMiddleware(handlers.CSSType, private)),
	}
}

func main() {
	// pre-parse templates, embedded in server binary
	handlers.Tmpl = template.Must(template.New("").ParseFS(embededTemplates,
		"templates/*.html",
	))

	mux := http.HandlerFunc(router.Serve)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
