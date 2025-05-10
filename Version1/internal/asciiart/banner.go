package asciiart

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func LoadBanner(bannerChoice string) (map[rune][]string, error) {
	fileName := "internal/asciiart/banners/" + bannerChoice + ".txt"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("Banner file not found: %s", fileName)
		return nil, fmt.Errorf("banner file %s not found", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Failed to open banner file %s: %v", fileName, err)
		return nil, fmt.Errorf("failed to open banner file %s: %v", fileName, err)
	}
	defer file.Close()

	banner := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	symbol := []string{}

	for i := ' '; i <= '~'; i++ {
		count := 0
		for scanner.Scan() && count < 8 {
			line := scanner.Text()
			if line == "" {
				if count == 0 {
					continue
				}
				break
			}
			symbol = append(symbol, line)
			count++
		}
		if count == 8 {
			banner[i] = symbol
			symbol = nil
		} else {
			log.Printf("Incomplete character %c: expected 8 lines, got %d", i, count)
			return nil, fmt.Errorf("incomplete character %c in banner %s", i, fileName)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading banner file %s: %v", fileName, err)
		return nil, fmt.Errorf("error reading banner file %s: %v", fileName, err)
	}

	if len(banner) == 0 {
		log.Printf("Banner file %s is empty or malformed", fileName)
		return nil, fmt.Errorf("banner file %s is empty or malformed", fileName)
	}

	return banner, nil
}
