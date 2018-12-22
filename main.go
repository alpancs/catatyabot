package main

import (
	"net/http"
	"os"

	"github.com/alpancs/catatyabot/handler"
)

func main() {
	http.HandleFunc("/"+os.Getenv("TELEGRAM_BOT_TOKEN"), handler.Telegram)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
