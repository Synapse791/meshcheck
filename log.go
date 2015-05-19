package main

import "log"

func LogFatal(msg string) {
	log.Fatal("[FATAL] " + msg)
}

func LogInfo(msg string) {
	log.Println("[INFO]  " + msg)
}