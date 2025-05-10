package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

// computeFileHash вычисляет SHA256-хэш файла.
func ComputeFileHash(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file for hash: %v", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatalf("Error hashing file: %v", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
}
