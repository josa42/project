package logger

import (
	"log"
	"os"
)

func InitLogger(filePath string) func() {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	return func() {
		f.Close()
	}

}
