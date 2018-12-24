package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/"+os.Getenv("TELEGRAM_BOT_TOKEN"), TelegramHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
