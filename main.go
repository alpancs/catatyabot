package main

import (
	"net/http"
	"os"

	"github.com/alpancs/catatyabot/app"
)

func main() {
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), app.Handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
