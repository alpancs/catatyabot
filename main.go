package main

import (
	"net/http"
	"os"

	"catatyabot/app"
)

func main() {
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), app.Handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
