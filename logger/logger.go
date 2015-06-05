package logger

import (
	"fmt"
	"log"
)

func Fatal(msg string) {
	log.Fatal(fmt.Sprintf("[FATAL] %s", msg))
}

func Warning(msg string) {
	log.Printf("[WARN]  %s", msg)
}

func Info(msg string) {
	log.Printf("[INFO]  %s", msg)
}