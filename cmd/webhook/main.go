package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const theDiscordWebhook = "replaceWithValidDiscordWebhook"

/*
WebhookInfo represents the Discord webhook data entry.
*/
type WebhookInfo struct {
	Content string `json:"content"`
}

func sendDiscordLogEntry(what string) {
	info := WebhookInfo{}
	info.Content = what + "\n"
	raw, _ := json.Marshal(info)
	resp, err := http.Post(theDiscordWebhook, "application/json", bytes.NewBuffer(raw))
	if err != nil {
		fmt.Println(err)
		fmt.Println(ioutil.ReadAll(resp.Body))
	}
}

func main() {
	for {
		text := "Heroku timer test at: " + time.Now().String()
		delay := time.Minute * 15

		sendDiscordLogEntry(text)
		time.Sleep(delay)
	}
}
