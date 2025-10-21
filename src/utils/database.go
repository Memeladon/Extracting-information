package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func ConnectDB(settings Settings) (*pgx.Conn, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		settings.DBUser, settings.DBPass, settings.DBHost, settings.DBPort, settings.DBName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	return conn, nil
}

func CreateTables(conn *pgx.Conn) error {
	queries := []string{
		`
			CREATE TABLE IF NOT EXISTS articles (
				id SERIAL PRIMARY KEY,
				title VARCHAR(255),
				author VARCHAR(255),
				abstract TEXT,
				content TEXT,
				source_url VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
		`
			CREATE TABLE IF NOT EXISTS article_metrics (
				id SERIAL PRIMARY KEY,
				article_id INTEGER REFERENCES articles(id),
				sentiment_score FLOAT,
				stylistic_features JSONB,
				language_confidence FLOAT,
				detected_language VARCHAR(50),
				analyzed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`,
	}

	// Выполняем все запросы последовательно
	for _, query := range queries {
		// Используем Exec с context.Background()
		_, err := conn.Exec(context.Background(), query)
		if err != nil {
			return fmt.Errorf("error executing query: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Tables created successfully")
	return nil
}
