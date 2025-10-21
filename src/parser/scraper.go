package parser

import (
	"Extracting-information/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Scraper struct {
	client *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (s *Scraper) ScrapeArticle(url string) (*utils.Article, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch article: %s", url)
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	article := &utils.Article{}
	err = json.NewDecoder(resp.Body).Decode(article)
	if err != nil {
		return nil, fmt.Errorf("failed to decode article: %w", err)
	}

	return article, nil
}
