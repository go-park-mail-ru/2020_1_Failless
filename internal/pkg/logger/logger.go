package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func OpenLogFile(serviceName string) *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%s.%v-%v-%v.log", serviceName, year, int(month), day)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	return file
}
