package ascii

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var (
	ErrBanner = errors.New("invalid banner")
	ErrASCII  = errors.New("invalid input: contains a not ASCII character")
)

type Template struct {
	banner map[rune][8]string
}

func NewTemplate(r io.Reader) (*Template, error) {
	scanner := bufio.NewScanner(r)

	banner := make(map[rune][8]string, 95)

	for r := ' '; r <= '~'; r++ {
		var arr [8]string

		for i := 0; i < len(arr); i++ {
			ok := scanner.Scan()
			if !ok {
				return nil, ErrBanner
			}
			text := scanner.Text()
			if i == 0 && text == "" {
				i--
				continue
			}
			arr[i] = text
		}

		banner[r] = arr
	}

	return &Template{banner}, nil
}

func (t *Template) Execute(in string) (string, error) {
	if in == "" {
		return "", nil
	}
	text := strings.ReplaceAll(in, "\\n", "\n")
	ss := strings.Split(text, "\n")
	var out []string
	for _, s := range ss {
		if s == "" {
			out = append(out, "")
			continue
		}
		var clines [8]strings.Builder
		for _, r := range s {
			char, ok := t.banner[r]
			if !ok {
				return "", ErrASCII
			}
			for i := range char {
				clines[i].WriteString(char[i])
			}
		}
		for i := range clines {
			out = append(out, clines[i].String())
		}
	}
	return strings.Join(out, "\n"), nil
}
