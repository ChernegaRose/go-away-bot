package main

import (
	"bytes"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"mime/multipart"
	"net/http"
)

func RepostToTelegraph(bot *tgbotapi.BotAPI, fileID string) (string, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: fileID,
	})
	if err != nil {
		return "", err
	}
	repost, err := RepostWithForm(file.Link(bot.Token), "https://telegra.ph/upload")
	if err != nil {
		return "", err
	}
	var links []struct {
		Src string `json:"src"`
	}
	if err := json.Unmarshal(repost, &links); err != nil {
		return "", err
	}
	return "https://telegra.ph" + links[0].Src, nil
}

func RepostWithForm(urlFrom string, urlTo string) ([]byte, error) {
	resp, err := http.Get(urlFrom)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(resp.Header.Get("content-type"), "file")

	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, resp.Body); err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", urlTo, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return content, nil
}
