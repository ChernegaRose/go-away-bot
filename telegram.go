package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func checkChannelPrivileges(bot *tgbotapi.BotAPI, chatName string) (string, bool) {
	admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{SuperGroupUsername: "@" + chatName}})
	if err != nil {
		return "Бот не администратор канала", false
	}
	me, err := bot.GetMe()
	if err != nil {
		return "Внутренняя ошибка", false
	}
	for _, admin := range admins {
		if admin.User.ID == me.ID {
			if admin.CanPostMessages {
				return "Бот настроен правильно!", true
			}
			return "Бот добавлен, но не может создавать посты, добавьте разрешение \"Публикация сообщений\"", false
		} else {
			continue
		}
	}
	return "Бот добавлен, но не может проверить разрешения", false
}

func getUserInfo(bot *tgbotapi.BotAPI, userID int64) (tgbotapi.Chat, error) {
	return bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: userID}})
}
