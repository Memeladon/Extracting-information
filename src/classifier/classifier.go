package classifier

import (
	"Extracting-information/src/constants"
	"fmt"
	"github.com/jbrukh/bayesian"
	"github.com/tealeg/xlsx"
	"log"
	"regexp"
	"strings"
)

// PreprocessText Функция для предобработки текста
func PreprocessText(text string) []string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)
	// Удаляем знаки препинания
	reg := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	text = reg.ReplaceAllString(text, "")
	// Разбиваем на слова
	words := strings.Fields(text)
	return words
}

func CreateNewsClassifier(filename string) *bayesian.Classifier {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Создаем классификатор
	classifier := bayesian.NewClassifier(
		constants.Politics,
		constants.Society,
		constants.Economy,
		constants.Sport,
		constants.Business,
		constants.Technology,
	)

	// Проходим по строкам в файле для обучения
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {

			category := row.Cells[4].Value // category-учитель
			//publishDate := row.Cells[6].String()
			//title := row.Cells[5].String()
			//body := row.Cells[7].String()
			text := row.Cells[12].Value // text

			//fmt.Printf("%v : %v\n", category, text)
			fmt.Println(category)
			var class bayesian.Class
			switch category {
			case "Политика":
				class = constants.Politics
			case "Общество":
				class = constants.Society
			case "Экономика":
				class = constants.Economy
			case "Спорт":
				class = constants.Sport
			case "Бизнес":
				class = constants.Business
			case "Технологии и медиа":
				class = constants.Technology
			default:
				continue // Пропускаем неизвестные категории
			}

			// Предобработка текста
			words := PreprocessText(text)

			// Обучаем классификатор
			classifier.Learn(words, class)
		}
	}

	// Сохраняем обученную модель (опционально)
	err = classifier.WriteToFile("classifier.gob")
	if err != nil {
		log.Fatal(err)
	}

	return classifier
}
