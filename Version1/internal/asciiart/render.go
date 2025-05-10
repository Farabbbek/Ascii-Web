package asciiart

import (
	"strings"
)

// RenderText обрабатывает и выводит текст в виде ASCII-графики без ограничения ширины
func RenderText(input string, bannerChoice string) (string, error) {
	// Заменяем литеральные строки "\n" на настоящие переносы строк
	// Это нужно, если пользователь вводит "\n" как текст, а не как символ переноса строки
	input = strings.ReplaceAll(input, "\\n", "\n")

	// Загружаем баннер
	banner, err := LoadBanner(bannerChoice)
	if err != nil {
		return "", err
	}

	var output []string

	lines := strings.Split(input, "\n") // Разделяем входной текст по символу новой строки
	for lineIndex, line := range lines {
		if line == "" { // Если строка пустая (например, после \n)
			// Добавляем пустую строку для переноса, но только если это не последняя строка
			if lineIndex < len(lines)-1 {
				output = append(output, "")
			}
			continue
		}

		lineOutput := make([]string, 8) // Массив для строк текущего символа (8 строк на символ)
		for _, char := range line {     // Обрабатываем каждый символ строки
			if lines, found := banner[char]; found { // Если символ найден в баннере
				for j := 0; j < 8; j++ {
					lineOutput[j] += lines[j] // Добавляем строки для текущего символа
				}
			} else { // Если символ не найден в баннере
				for j := 0; j < 8; j++ {
					lineOutput[j] += " " // Добавляем пробелы вместо ASCII-графики
				}
			}
		}

		// Добавляем все строки текущей линии в общий вывод
		for j := 0; j < 8; j++ {
			output = append(output, lineOutput[j])
		}

		// Добавляем пустую строку после каждой обработанной строки ввода, кроме последней
		if lineIndex < len(lines)-1 {
			output = append(output, "")
		}
	}

	return strings.Join(output, "\n"), nil
}
