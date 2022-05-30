package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func pictureAsHtmlLink(link string) string {
	return "<a href=\"" + link + "\">&#8205;</a>"
}

func articleFromPost(contest Contest, post Post, id string) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.InlineQueryResultArticle{Type: "article", ID: id}

	article.Title = "Розыгрыш " + contest.ContestName + " от @" + contest.Username

	article.Description = contest.ContestDescription
	article.InputMessageContent = tgbotapi.InputTextMessageContent{
		Text:      pictureAsHtmlLink(post.Image) + post.Message,
		ParseMode: "HTML",
	}
	return article
}

func articleWithMenu(id string, from int64, contest Contest, chatType string) interface{} {
	marshal, err := json.Marshal(InlineData{
		From:    from,
		Contest: contest.ContestID,
	})
	if err != nil {
		return articleFromPost(contest,
			Post{
				Message: "Ошибка",
			}, id)
	}
	var article tgbotapi.InlineQueryResultArticle
	if post, ok := posts[contest.ContestID][chatType]; ok {
		article = articleFromPost(contest, post, id)
		article.ThumbURL = post.Image
	} else if post, ok := posts[contest.ContestID]["post"]; ok {
		article = articleFromPost(contest, post, id)
		article.ThumbURL = post.Image
	} else {
		article = articleFromPost(contest,
			Post{
				Message: "Конкурс от @" + contest.Username,
			}, id)
	}
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить друга", contest.Username),
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "tg://resolve?domain="+contest.Username),
		),
	)
	return article
}

func messageToChannel(contest Contest, chatType string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.Chattable {
	if post, ok := posts[contest.ContestID][chatType]; ok {
		if post.Image != "" {
			cfg := tgbotapi.NewPhotoToChannel("@"+contest.Username, tgbotapi.FileURL(post.Image))
			cfg.Caption = post.Message
			cfg.ReplyMarkup = keyboard
			return cfg
		} else {
			cfg := tgbotapi.NewMessageToChannel("@"+contest.Username, post.Message)
			cfg.ReplyMarkup = keyboard
			return cfg
		}
	}
	cfg := tgbotapi.NewMessageToChannel("@"+contest.Username, contest.ContestName)
	cfg.ReplyMarkup = keyboard
	return cfg
}

func postToChannel(contest Contest) tgbotapi.Chattable {
	marshal, err := json.Marshal(InlineData{
		Contest: contest.ContestID,
	})
	if err != nil {
		log.Println(err)
	}
	return messageToChannel(contest, "post",
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)))))
}

func guideToChannel(contest Contest) tgbotapi.Chattable {
	return messageToChannel(contest, "guide",
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить", ""))))

}
