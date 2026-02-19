package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
	Slug    string `json:"slug"`
}

func HomePage(w http.ResponseWriter, _ *http.Request) {
	files, err := os.ReadDir("./articles")
	if err != nil {
		http.Error(w, "Could not read Articles", http.StatusInternalServerError)
		return
	}
	var articles []Article
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), ".json") {
			data, err := os.ReadFile(filepath.Join("./articles", file.Name()))
			if err != nil {
				continue
			}
			var article Article
			if err := json.Unmarshal(data, &article); err != nil {
				continue
			}
			articles = append(articles, article)
		}
	}
	tmpl, err := template.ParseFiles("./templates/guest/home.html")
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, articles); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func ArticlePage(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/article/")
	data, err := os.ReadFile(filepath.Join("./articles", slug+".json"))
	if err != nil {
		http.Error(w, "Could not read article", http.StatusInternalServerError)
		return
	}
	var article Article
	if err := json.Unmarshal(data, &article); err != nil {
		http.Error(w, "Could not parse article", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./templates/guest/article.html")
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, article); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}
