package main

import (
	"fmt"
	"net/http"

	"personal-blog/handlers"
	"personal-blog/middleware"
)

func main() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/article/", handlers.ArticlePage)

	http.HandleFunc("/admin/login", handlers.LoginHandler)
	http.HandleFunc("/admin/logout", handlers.LogoutHandler)

	http.HandleFunc("/admin/dashboard", middleware.RequireAuth(handlers.DashboardHandler))
	http.HandleFunc("/admin/article/add", middleware.RequireAuth(handlers.AddArticleHandler))
	http.HandleFunc("/admin/article/edit/", middleware.RequireAuth(handlers.EditArticleHandler))
	http.HandleFunc("/admin/article/delete/", middleware.RequireAuth(handlers.DeleteArticleHandler))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	fmt.Println("Guest site: http://localhost:8080")
	fmt.Println("Admin site: http://localhost:8080/admin/login")
	fmt.Println("Credentials: admin/admin")

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
