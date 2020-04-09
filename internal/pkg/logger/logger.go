package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func OpenLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v", year, int(month), day)
	file, err := os.OpenFile("eventum."+fileName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	return file
}
