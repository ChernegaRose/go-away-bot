package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func pictureAsHtmlLink(link string) string {
	return "<a href=\"" + link + "\">&#8205;</a>"
}

func articlePrivate(id string, from int64) interface{} {
	marshal, err := json.Marshal(XMessage{
		From:  from,
		Post:  "friend",
		Query: id,
	})
	if err != nil {
		return nil
	}
	article := tgbotapi.NewInlineQueryResultArticleHTML(id, "Пригласить друга",
		pictureAsHtmlLink(getLinkToPictureOnTelegramCDN(""))+
			"Прими участие в розыгрыше %розыгрыш_нейм% и получи повышенный шанс пригласив друзей")
	article.Description = "Отправьте ссылку-приглашение другу и получите повышенный шанс выигрыша"
	article.ThumbURL = getLinkToPictureOnTelegramCDN("")
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить друга", ""),
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "tg://resolve?domain=chernegarose"),
		),
	)
	return article
}

func articleChat(id string, from int64) interface{} {
	marshal, err := json.Marshal(XMessage{
		From:  from,
		Post:  "chat",
		Query: id,
	})
	if err != nil {
		return nil
	}
	article := tgbotapi.NewInlineQueryResultArticleHTML(id, "Пригласить друга",
		pictureAsHtmlLink(getLinkToPictureOnTelegramCDN(""))+
			"Приглашаем пользователей чата принять участие в розыгрыше %розыгрыш_нейм% и получить повышенный шанс пригласив друзей")
	article.Description = "Отправьте ссылку-приглашение другу и получите повышенный шанс выигрыша"
	article.ThumbURL = getLinkToPictureOnTelegramCDN("")
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить друга", ""),
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "tg://resolve?domain=chernegarose"),
		),
	)
	return article
}

func articleChannel(id string, from int64) interface{} {
	marshal, err := json.Marshal(XMessage{
		From:  from,
		Post:  "channel",
		Query: id,
	})
	if err != nil {
		return nil
	}
	article := tgbotapi.NewInlineQueryResultArticleHTML(id, "Пригласить друга",
		pictureAsHtmlLink(getLinkToPictureOnTelegramCDN(""))+
			"Приглашаем подписчиков канала принять участие в розыгрыше %розыгрыш_нейм% и получить повышенный шанс пригласив друзей")
	article.Description = "Отправьте ссылку-приглашение другу и получите повышенный шанс выигрыша"
	article.ThumbURL = getLinkToPictureOnTelegramCDN("")
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить друга", ""),
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "tg://resolve?domain=chernegarose"),
		),
	)
	return article
}

func adminPost(id string, photo string) interface{} {
	marshal, err := json.Marshal(XMessage{
		Post:  "post",
		Query: "chernegarose",
	})
	if err != nil {
		return nil
	}
	article := tgbotapi.NewInlineQueryResultCachedPhoto(id, photo)
	article.Caption = "Отправьте ссылку-приглашение другу и получите повышенный шанс выигрыша"
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
	)
	return article
}
