package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alpancs/catatan-belanja-bot/handler"
)

func main() {
	http.HandleFunc("/"+os.Getenv("TELEGRAM_BOT_TOKEN"), handler.Telegram)

	port := getPort()
	fmt.Println("listening port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "9000"
	}
	return port
}
