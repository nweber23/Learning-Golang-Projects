package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AdminData struct {
	Error   string
	Title   string
	Content string
	Date    string
	Slug    string
}

func DashboardHandler(w http.ResponseWriter, _ *http.Request) {
	articles, err := getAllArticles()
	if err != nil {
		http.Error(w, "Couldn't get articles", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./templates/admin/dashboard.html")
	if err != nil {
		http.Error(w, "Couldn't get template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, articles); err != nil {
		http.Error(w, "Couldn't get template", http.StatusInternalServerError)
	}
}

func AddArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		displayAddArticleForm(w)
		return
	}
	if r.Method == http.MethodPost {
		processAddArticle(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func displayAddArticleForm(w http.ResponseWriter, errorMsg ...string) {
	data := AdminData{}
	if len(errorMsg) > 0 {
		data.Error = errorMsg[0]
	}
	tmpl, err := template.ParseFiles("./templates/admin/add_article.html")
	if err != nil {
		http.Error(w, "Couldn't get template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Couldn't get template", http.StatusInternalServerError)
	}
}

func processAddArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if title == "" || content == "" || date == "" {
		displayAddArticleForm(w, "All fields are required")
		return
	}

	slug := generateSlug(title)

	articlePath := filepath.Join("./articles", slug+".json")
	if _, err := os.Stat(articlePath); err == nil {
		displayAddArticleForm(w, "Article with this title already exists")
		return
	}
	article := Article{
		ID:      slug,
		Title:   title,
		Content: content,
		Date:    date,
		Slug:    slug,
	}

	if err := saveArticle(article); err != nil {
		displayAddArticleForm(w, "Couldn't save article")
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/admin/article/edit")
	if r.Method == http.MethodGet {
		displayEditArticleForm(w, slug)
		return
	}
	if r.Method == http.MethodPost {
		processEditArticle(w, r, slug)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func displayEditArticleForm(w http.ResponseWriter, slug string, errorMsg ...string) {
	data, err := os.ReadFile(filepath.Join("./articles", slug+".json"))
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	var article Article
	if err := json.Unmarshal(data, &article); err != nil {
		http.Error(w, "Could not parse article data", http.StatusInternalServerError)
		return
	}

	adminData := AdminData{
		Title:   article.Title,
		Content: article.Content,
		Date:    article.Date,
		Slug:    article.Slug,
	}
	if len(errorMsg) > 0 {
		adminData.Error = errorMsg[0]
	}
	tmpl, err := template.ParseFiles("./templates/admin/edit_article.html")
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, adminData); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func processEditArticle(w http.ResponseWriter, r *http.Request, slug string) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if title == "" || content == "" || date == "" {
		displayEditArticleForm(w, "All fields are required")
		return
	}

	article := Article{
		ID:      slug,
		Title:   title,
		Content: content,
		Date:    date,
		Slug:    slug,
	}
	if err := saveArticle(article); err != nil {
		displayEditArticleForm(w, slug, "Couldn't save article")
		return
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/admin/article/delete")
	articlePath := filepath.Join("./articles", slug+".json")
	if err := os.Remove(articlePath); err != nil {
		http.Error(w, "Could not delete article", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

func saveArticle(article Article) error {
	if err := os.MkdirAll("./articles", os.ModePerm); err != nil {
		return err
	}
	data, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("./articles", article.Slug+".json"), data, 0644)
}

func getAllArticles() ([]Article, error) {
	if err := os.MkdirAll("./articles", os.ModePerm); err != nil {
		return nil, err
	}
	files, err := os.ReadDir("./articles")
	if err != nil {
		return nil, err
	}

	var articles []Article
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
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
	return articles, nil
}
