package handlers

import (
	"html/template"
	"net/http"

	"personal-blog/middleware"
)

type LoginData struct {
	Error string
}

const (
	ADMIN_USERNAME = "admin"
	ADMIN_PASSWORD = "admin"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		displayLoginForm(w)
		return
	}
	if r.Method == http.MethodPost {
		processLogin(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func displayLoginForm(w http.ResponseWriter, errorMsg ...string) {
	data := LoginData{}
	if len(errorMsg) > 0 {
		data.Error = errorMsg[0]
	}
	tmpl, err := template.ParseFiles("./templates/admin/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != ADMIN_USERNAME || password != ADMIN_PASSWORD {
		displayLoginForm(w, "Invalid username or password")
		return
	}
	sessionID, err := middleware.CreateSession(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	middleware.SetSessionCookie(w, sessionID)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		middleware.DeleteSession(cookie.Value)
	}
	middleware.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
