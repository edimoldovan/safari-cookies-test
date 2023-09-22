package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"main/middlewares"
	"main/session"
	"net/http"
)

var Tmpl *template.Template

func Home(w http.ResponseWriter, r *http.Request) {
	value, err := session.ReadSigned(r)
	log.Println("session value:", value)
	if err != nil {
		middlewares.Log.Error("no session", "method", r.Method, "url", r.URL.Path, "error", err)
	}
	if err := Tmpl.ExecuteTemplate(w, "home", map[string]interface{}{}); err != nil {
		fmt.Printf("ERR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CSSType(w http.ResponseWriter, r *http.Request) {
	value, err := session.ReadSigned(r)
	log.Println("session value:", value)
	if err != nil {
		middlewares.Log.Error("no session", "method", r.Method, "url", r.URL.Path, "error", err)
	}

	w.Header().Add("Content-type", "text/css")
	w.Header().Set("Cache-Control", "public, max-age=0")

	var buf bytes.Buffer
	tmpl := Tmpl.Lookup("css-type")
	_ = tmpl.Execute(&buf, map[string]interface{}{})

	w.Write([]byte(buf.String()))
}

func SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     session.CookieName,
		Value:    "my-email@my-domain.com",
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	session.WriteSigned(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
