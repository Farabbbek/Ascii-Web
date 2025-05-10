package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Структуры данных:“
type ASCIIChar struct {
	lines [8]string
}

// ASCIIArt - основная структура, которая хранит весь шрифт
type ASCIIArt struct {
	chars map[rune]ASCIIChar
}

func NewASCIIArt() *ASCIIArt {
	return &ASCIIArt{
		chars: make(map[rune]ASCIIChar),
	}
}

// LoadFont читает файл шрифта и загружает все символы в память
func (a *ASCIIArt) LoadFont(filename string) error {
	// Открываем файл шрифта
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()
	supportedChars := []rune{
		' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		':', ';', '<', '=', '>', '?', '@',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'[', '\\', ']', '^', '_', '`',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'{', '|', '}', '~',
	}
	scanner := bufio.NewScanner(file)
	var currentLines [8]string
	lineIndex := 0
	charIndex := 0
	// Читаем файл построчно
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if lineIndex > 0 && charIndex < len(supportedChars) {
				a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
				charIndex++
				lineIndex = 0
				currentLines = [8]string{}
			}
		} else {
			if lineIndex < 8 {
				currentLines[lineIndex] = line
				lineIndex++
			}
		}
	}
	if lineIndex > 0 && charIndex < len(supportedChars) {
		a.chars[supportedChars[charIndex]] = ASCIIChar{lines: currentLines}
	}
	return nil
}

// RenderText преобразует входной текст в ASCII-арт представление
func (a *ASCIIArt) RenderText(input string) string {
	var result strings.Builder
	if input == "" {
		return ""
	}

	if input == "\\n" {
		return "$"
	}

	lines := strings.Split(strings.ReplaceAll(input, "\\n", "\n"), "\n")

	for i, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}

		artLines := [8]string{}
		for _, char := range line {
			if art, exists := a.chars[char]; exists {

				for j := 0; j < 8; j++ {
					artLines[j] += art.lines[j]
				}
			}
		}

		for j := 0; j < 8; j++ {
			result.WriteString(artLines[j])
			result.WriteString("\n")
		}
		if i < len(lines)-1 || (len(input) >= 2 && input[len(input)-2:] == "\\n") {
			result.WriteString("")
		}
	}
	return result.String()
}
