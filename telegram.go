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
