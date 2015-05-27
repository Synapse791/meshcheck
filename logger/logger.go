package logger

import "log"

func Fatal(msg string) {
	log.Fatal("[FATAL] " + msg)
}

func Warning(msg string) {
	log.Println("[WARN]  " + msg)
}

func Info(msg string) {
	log.Println("[INFO]  " + msg)
}