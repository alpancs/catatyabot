package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
