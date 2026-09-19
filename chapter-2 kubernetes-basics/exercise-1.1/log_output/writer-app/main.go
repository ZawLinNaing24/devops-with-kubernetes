package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/google/uuid"
)

func main() {
	randomString := uuid.New().String()

	// 2. File သိမ်းမည့် နေရာ (Shared Volume နေရာ)
	dir := "/usr/src/app/files"
	os.MkdirAll(dir, os.ModePerm)
	filePath := filepath.Join(dir, "log.txt")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		line := fmt.Sprintf("%s: %s\n", timestamp, randomString)

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(line)
			f.Close()
			fmt.Print("Wrote to file: ", line) 
		} else {
			fmt.Println("Error writing to file:", err)
		}
	}
}