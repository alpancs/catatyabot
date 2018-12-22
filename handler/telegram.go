package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alpancs/catatyabot/telegram"
)

var (
	telegramAPI = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("TELEGRAM_BOT_TOKEN"))
)

func Telegram(w http.ResponseWriter, r *http.Request) {
	update, err := parseUpdate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	err = respond(update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
	}
}

func parseUpdate(body io.ReadCloser) (u telegram.Update, err error) {
	defer body.Close()
	err = json.NewDecoder(body).Decode(&u)
	return u, err
}

func respond(update telegram.Update) error {
	response, err := telegram.Respond(update)
	if err != nil {
		return err
	}
	if response == nil {
		return nil
	}

	return sendMessage(update, response)
}

func sendMessage(update telegram.Update, response *telegram.Response) error {
	reqBody, err := json.Marshal(response)
	if err != nil {
		return err
	}
	url := telegramAPI + "sendMessage"
	res, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	resCode := res.StatusCode
	if resCode == 200 {
		return nil
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("URL: %s\nrequest body: %s\nresponse code: %d\nresponse body: %s", url, string(reqBody), resCode, string(resBody)))
}
