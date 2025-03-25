package main

import (
	"Extracting-information/src/classifier"
	"fmt"
)

func main() {
	newsClassifier := classifier.CreateNewsClassifier("news.xlsx")
	if newsClassifier == nil {
		panic("Failed to create news classifier")
	}

	// Пример классификации нового текста
	newText := "В России увеличили возраст призыва"
	words := classifier.PreprocessText(newText) // Предобработка текста
	scores, likely, _ := newsClassifier.LogScores(words)
	fmt.Printf("New Text: %s\n", newText)
	fmt.Printf("Predicted Category: %v\n", likely)
	fmt.Printf("Scores: %v\n", scores)
}
