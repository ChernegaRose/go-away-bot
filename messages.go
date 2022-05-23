package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func pictureAsHtmlLink(link string) string {
	return "<a href=\"" + link + "\">&#8205;</a>"
}

func articleFromPost(post Post, id string) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.InlineQueryResultArticle{Type: "article", ID: id}
	article.Title = post.Title
	article.Description = post.Description
	article.InputMessageContent = tgbotapi.InputTextMessageContent{
		Text:      pictureAsHtmlLink(post.Image) + post.Message,
		ParseMode: "HTML",
	}
	return article
}

func articleWithMenu(id string, from int64, contest Contest, chatType string) interface{} {
	marshal, err := json.Marshal(Params{
		From:    from,
		Contest: contest.ContestID,
	})
	if err != nil {
		return articleFromPost(Post{
			Title:       "Ошибка",
			Message:     "Ошибка",
			Description: "Ошибка",
		}, id)
	}
	var article tgbotapi.InlineQueryResultArticle
	if post, ok := posts[contest.ContestID][chatType]; ok {
		article = articleFromPost(post, id)
		article.ThumbURL = post.Image
	} else if post, ok := posts[contest.ContestID]["post"]; ok {
		article = articleFromPost(post, id)
		article.ThumbURL = post.Image
	} else {
		article = articleFromPost(Post{
			Title:       "Пригласить друга",
			Message:     "Конкурс от @" + contest.Username,
			Description: "Приглашение на " + contest.ContestName,
		}, id)
	}
	article.ReplyMarkup = new(tgbotapi.InlineKeyboardMarkup)
	*(article.ReplyMarkup) = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("Пригласить друга", ""),
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "tg://resolve?domain="+contest.Username),
		),
	)
	return article
}
