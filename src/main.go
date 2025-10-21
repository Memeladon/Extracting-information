package main

import (
	"Extracting-information/src/parser"
	"Extracting-information/src/utils"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// main router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// settings
	err := godotenv.Load("src/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	settings := utils.Settings{
		DBUser: os.Getenv("POSTGRES_USER"),
		DBPass: os.Getenv("POSTGRES_PASSWORD"),
		DBHost: os.Getenv("POSTGRES_HOST"),
		DBPort: os.Getenv("POSTGRES_PORT"),
		DBName: os.Getenv("POSTGRES_NAME"),
		DBUrl:  os.Getenv("POSTGRES_URL"),
	}
	log.Printf("Settings: %v", settings)

	// database
	conn, err := utils.ConnectDB(settings)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	defer conn.Close(context.Background())

	err = utils.CreateTables(conn)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	// scraper & parser
	scraper := parser.NewScraper()
	parser := parser.NewTextParser()

	var templates = template.Must(template.ParseGlob("templates/*.html"))

	// routing
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		articles, err := fetchAndProcessArticles(conn, scraper, parser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "index.html", struct {
			Articles []utils.Article
			Title    string
		}{
			Articles: *articles,
			Title:    "Анализ тональности текста",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

func fetchAndProcessArticles(conn *pgx.Conn, scraper *parser.Scraper, parser *parser.TextParser) (*[]utils.Article, error) {
	// Примеры URL, ты можешь заменить на реальные
	urls := []string{
		"https://www.sciencedirect.com",
		"https://link.springer.com",
		"https://www.mdpi.com",
		"https://arxiv.org",
		"https://dl.acm.org",
	}

	articles, err := utils.ReadAllArticles(conn)
	if err != nil {
		return nil, err
	}

	for _, url := range urls {
		article, err := scraper.ScrapeArticle(url)
		if err != nil {
			continue
		}

		article.Title = parser.CleanText(article.Title)
		article.Abstract = parser.CleanText(article.Abstract)
		article.Content = parser.CleanText(article.Content)

		articleDetails := utils.Article{
			Title:    article.Title,
			Author:   article.Author,
			Abstract: article.Abstract,
			Content:  article.Content,
			URL:      url,
		}

		err = utils.CreateArticle(conn, articleDetails)
		// Сохраняем в базу
		if err != nil {
			fmt.Println("DB error:", err)
			continue
		}

	}

	return articles, nil
}
