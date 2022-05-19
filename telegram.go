package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkFollowerByChatID(bot *tgbotapi.BotAPI, userID int64, chatID int64) bool {

	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
		ChatID: chatID,
		UserID: userID,
	}})
	if err != nil {
		return false
	} else {
		return member.Status != "left"
	}
}

func checkFollowerByChatName(bot *tgbotapi.BotAPI, userID int64, chatName string) bool {

	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
		SuperGroupUsername: "@" + chatName,
		UserID:             userID,
	}})
	if err != nil {
		return false
	} else {
		return member.Status != "left"
	}
}

func getLinkToPictureOnTelegramCDN(id string) string {
	return "https://cdn4.telegram-cdn.org/file/F9Yr-bV0Be8ITM_SGsF5dzUZ4nJa3dl-lCtxV2qqwD0kf9FcRWHPUGtP1mKq-YLNSQ-ylagk0jE--Y1UxXRSAsfYqvhVR-rvpkWlC7CmSDvA7hLeeG56B32sDnaXocw-2a_6Up6-kJWJWxT8eC7xFLHKmJp2IxGaZCWqpfDvrOXh5JsB1xZ0F7O7lFw9XOEuryRvNIgWUmaTTA-1ieJvpaPlTKlLkBXZyhqzlhnYm9e2HTnd5dGgWwnySZsGo-c1zfuweuscMINyUkefj5P2sPbcP4srhriOxd9Av5V8EY2zTl2YfS3dGBN_siXllVcibOm68qzB-IPBcqzhRP54SA.jpg"
}
