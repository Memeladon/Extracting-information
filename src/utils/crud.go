package utils

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type ArticleCRUD interface {
	CreateArticle(conn *pgx.Conn, articleDetails Article) error
	ReadAllArticles(conn *pgx.Conn) (*[]Article, error)
	ReadArticleById(conn *pgx.Conn, id int) (*Article, error)
	UpdateArticle(conn *pgx.Conn, id int, articleDetails Article) error
	DeleteArticle(conn *pgx.Conn, id int) error
}

func CreateArticle(conn *pgx.Conn, articleDetails Article) error {
	query := `
        INSERT INTO articles (title, author, abstract, content, source_url) VALUES (@title, @author, @abstract, @content, @source_url)
    `
	args := pgx.NamedArgs{
		"title":      articleDetails.Title,
		"author":     articleDetails.Author,
		"abstract":   articleDetails.Abstract,
		"content":    articleDetails.Content,
		"source_url": articleDetails.URL,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Inserting Article Details")
		return err
	}
	return nil
}

func ReadArticleById(conn *pgx.Conn, ArticleId int) (*Article, error) {
	// Define the SQL query to select article by its ID.
	query := `
        SELECT * FROM articles WHERE id = @ArticleId
    `
	// Define the named argument for the article ID.
	args := pgx.NamedArgs{
		"ArticleId": ArticleId,
	}
	// Execute the query with named arguments to fetch the article details from the database.
	row := conn.QueryRow(context.Background(), query, args)

	var article Article
	// Scan the retrieved row into the article struct.
	err := row.Scan(&article.ID, &article.Title, &article.Author, &article.Abstract, &article.Content, &article.URL)
	if err != nil {
		log.Printf("Error Fetching Article Details")
		return nil, err
	}
	return &article, nil
}

func ReadAllArticles(conn *pgx.Conn) (*[]Article, error) {
	// Define the SQL query to select all articles.
	query := `
        SELECT * FROM articles
    `
	// Execute the query to fetch all article details from the database.
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error Querying the Table")
		return nil, err
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	if err != nil {
		log.Printf("Error Fetching Articles")
		return nil, err
	}
	return &articles, nil
}

func UpdateArticle(conn *pgx.Conn, ArticleId int, articleDetails Article) error {
	// Define the SQL query to update a article's details by its ID.
	query := `
        UPDATE books
        SET title = @title, author = @author, abstract = @abstract, content = @content, source_url = @source_url
        WHERE id = @id
    `
	// Define the named arguments for the query.
	args := pgx.NamedArgs{
		"id":         articleDetails.ID,
		"title":      articleDetails.Title,
		"author":     articleDetails.Author,
		"abstract":   articleDetails.Abstract,
		"content":    articleDetails.Content,
		"source_url": articleDetails.URL,
	}
	// Execute the query with named arguments to update the article details in the database.
	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Updating Article Details")
		return err
	}
	return nil
}

func DeleteArticleById(conn *pgx.Conn, ArticleId int) error {
	// Define the SQL query to delete a article by its ID.
	query := `
        DELETE FROM articles WHERE id = @id
    `
	// Define the named argument for the article ID.
	args := pgx.NamedArgs{
		"id": ArticleId,
	}
	// Execute the query with named arguments to delete the article from the database.
	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Deleting Article")
		return err
	}
	return nil
}
