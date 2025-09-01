package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramBot struct {
	token  string
	chatID string
	client *http.Client
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegramBot(token, chatID string) *TelegramBot {
	return &TelegramBot{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

func (t *TelegramBot) SendMessage(text string) error {
	message := TelegramMessage{
		ChatID: t.chatID,
		Text:   text,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API error: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}