package parser

import (
	"Extracting-information/src/utils"
	"strings"
)

type TextParser struct{}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (p *TextParser) CleanText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	return text
}

func (p *TextParser) FilterEmpty(articles []utils.Article) []utils.Article {
	var filtered []utils.Article
	for _, a := range articles {
		if a.Title != "" && a.Content != "" {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func (p *TextParser) RemoveDuplicates(articles []utils.Article) []utils.Article {
	seen := make(map[string]bool)
	var unique []utils.Article

	for _, article := range articles {
		if !seen[article.URL] {
			unique = append(unique, article)
			seen[article.URL] = true
		}
	}

	return unique
}
