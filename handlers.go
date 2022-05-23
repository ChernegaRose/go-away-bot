package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InlineHandler(bot *tgbotapi.BotAPI, db *sql.DB, inlineQuery *tgbotapi.InlineQuery) {
	query := Query{
		QueryID:   inlineQuery.ID,
		QueryText: inlineQuery.Query,
		ChatType:  inlineQuery.ChatType,
		UserID:    inlineQuery.From.ID,
		UserName:  inlineQuery.From.FirstName + " " + inlineQuery.From.LastName,
		UserLang:  inlineQuery.From.LanguageCode,
	}
	if inlineQuery.From.UserName != "" {
		query.UserName += " @" + inlineQuery.From.UserName
	}
	if (inlineQuery.ChatType == "sender") || (inlineQuery.ChatType == "") {
		return
	}
	err := queries.insert(db, query)
	if err != nil {
		log.Println(err)
	}

	var results []interface{}
	chatType := inlineQuery.ChatType
	if chatType == "supergroup" {
		chatType = "group"
	}
	results = append(results, articleWithMenu(inlineQuery.ID, inlineQuery.From.ID, contests[settings.Featured], chatType))

	inline := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     0,
		IsPersonal:    true,
	}
	if _, err := bot.Request(inline); err != nil {
		log.Println(err)
	}
}

func MessageHandler(bot *tgbotapi.BotAPI, db *sql.DB, message *tgbotapi.Message) {
	// Construct a new message from the given chat ID and containing
	// the text that we received.
	msg := tgbotapi.NewMessage(message.Chat.ID, "Message: "+message.Text)

	if message.Document != nil {
		url, err := RepostToTelegraph(bot, message.Document.FileID)
		if err != nil {
			log.Println(err)
			return
		}
		msg.Text += "\n" + url
	}

	if len(message.Photo) > 0 {
		url, err := RepostToTelegraph(bot, message.Photo[len(message.Photo)-1].FileID)
		if err != nil {
			log.Println(err)
			return
		}
		msg.Text += "\n" + url
	}

	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}

	// Send the message.
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func CallbackHandler(bot *tgbotapi.BotAPI, db *sql.DB, callback *tgbotapi.CallbackQuery) {
	var params Params
	responseMessage := "Неизвестная ошибка!"
	if err := json.Unmarshal([]byte(callback.Data), &params); err != nil {
		log.Println(err)
	} else {
		if contest, ok := contests[params.Contest]; ok {
			if !isBegin(contest.ContestStart) || isEnd(contest.ContestEnd) || contest.ContestActive == 0 {
				responseMessage = "Конкурс завершен..."
			} else if checkFollowerByChatName(bot, callback.From.ID, contest.Username) {
				if _, ok := members[contest.ContestID][callback.From.ID]; ok {
					responseMessage = "Вы уже участвуете."
				} else {
					//proceed adding
					member := Member{
						UserID:       callback.From.ID,
						UserName:     callback.From.FirstName + " " + callback.From.LastName,
						UserLang:     callback.From.LanguageCode,
						FromID:       params.From,
						MessageID:    callback.InlineMessageID,
						ChatInstance: callback.ChatInstance,
						ContestID:    contest.ContestID,
					}
					if callback.From.UserName != "" {
						member.UserName += " @" + callback.From.UserName
					}
					err = members.insert(db, member)
					if err != nil {
						log.Println(err)
					} else {
						responseMessage = "Вы стали участником!"
					}
				}
			} else {
				responseMessage = "Вы не подписаны на канал!"
			}
		} else {
			responseMessage = "Конкурс не найден..."
		}
	}

	cb := tgbotapi.NewCallbackWithAlert(callback.ID, responseMessage)
	if _, err := bot.Request(cb); err != nil {
		log.Println(err)
	}
}

func ChosenInlineResultHandler(bot *tgbotapi.BotAPI, db *sql.DB, result *tgbotapi.ChosenInlineResult) {
	message := Message{
		MessageID: result.InlineMessageID,
		QueryID:   result.ResultID,
		UserID:    result.From.ID,
		UserName:  result.From.FirstName + " " + result.From.LastName,
	}
	if result.From.UserName != "" {
		message.UserName += " @" + result.From.UserName
	}
	err := messages.insert(db, message)
	if err != nil {
		log.Println(err)
	}
}

func WTFMessageHandler(bot *tgbotapi.BotAPI) { //deprecated
	msg := tgbotapi.NewPhoto(404596828, //-1001749951708,
		tgbotapi.FileID("AgACAgIAAxkBAAPTYoZfgZK9dBg0AkoYgGxrWTwbSvkAApa7MRtnmzBIqVyzvJP8PJ8BAAMCAAN5AAMkBA"))
	//msg := tgbotapi.NewMessage(-1001708457568, "💩")
	msg.Caption = "Приглашаем поучаствовать..."
	marshal, err := json.Marshal(Params{
		From: 0,
	})
	if err != nil {
		log.Println(err)
		return
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Участвовать", string(marshal)),
		),
	)
	if _, err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}
