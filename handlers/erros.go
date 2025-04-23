package handlers

import "log"

func ErrorHandler(err error) {
	log.Printf("[ERROR | go-telegram/bot] %v", err)
}
