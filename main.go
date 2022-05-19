package main

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var settings Settings
var members = make(Members)
var inlineQueries = make(InlineQueries)

func main() {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)

	err = settings.load(db)
	if err != nil {
		log.Panic(err)
	}
	err = members.load(db)
	if err != nil {
		log.Panic(err)
	}
	err = inlineQueries.load(db)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(settings.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.InlineQuery != nil {
			err := inlineQueries.insert(db, InlineQuery{
				id:       update.InlineQuery.ID,
				sender:   update.InlineQuery.From.ID,
				name:     update.InlineQuery.From.FirstName + " " + update.InlineQuery.From.LastName,
				chatType: update.InlineQuery.ChatType,
			})
			if err != nil {
				panic(err)
			}
			var results []interface{}

			switch update.InlineQuery.ChatType {
			case "private":
				results = append(results, articlePrivate(update.InlineQuery.ID, update.InlineQuery.From.ID))
			case "channel":
				results = append(results, articleChannel(update.InlineQuery.ID, update.InlineQuery.From.ID))
			default:
				results = append(results, articleChat(update.InlineQuery.ID, update.InlineQuery.From.ID))
			}

			if update.InlineQuery.From.ID == settings.admin {
				results = append(results, adminPost(update.InlineQuery.ID+"a",
					"AgACAgIAAxkBAAPTYoZfgZK9dBg0AkoYgGxrWTwbSvkAApa7MRtnmzBIqVyzvJP8PJ8BAAMCAAN5AAMkBA"))
			}

			inline := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       results,
				CacheTime:     0,
				IsPersonal:    true,
			}
			if _, err = bot.Request(inline); err != nil {
				panic(err)
			}
		} else if (update.Message != nil) && (update.Message.From.ID == settings.admin) {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Message: "+update.Message.Text)

			if len(update.Message.Photo) > 0 {
				msg.Text += "\n" + update.Message.Photo[0].FileID
			}

			msg.ReplyMarkup = nil

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			if checkFollowerByChatName(bot, update.CallbackQuery.From.ID, "gl1ch") {
				if _, ok := members[update.CallbackQuery.From.ID]; ok {
					callback := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Вы уже участвуете")
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}
				} else {
					var msg XMessage
					err := json.Unmarshal([]byte(update.CallbackQuery.Data), &msg)
					if err != nil {
						panic(err)
					}
					m := Member{
						id:   update.CallbackQuery.From.ID,
						from: msg.From,
						date: time.Now().UTC().Format("2006.01.02 15:04:05"),
						post: msg.Post,
					}
					err = members.insert(db, m)
					if err != nil {
						panic(err)
					}
					// Respond to the callback query, telling Telegram to show the user
					// a message with the data received.
					callback := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Вы успешно участвуете")
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}
				}
			} else {
				callback := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Вы не подписаны на канал!")
				if _, err := bot.Request(callback); err != nil {
					panic(err)
				}
			}
		} else if update.MyChatMember != nil {
			//fmt.Println(update.MyChatMember.NewChatMember)
		}
	}
}
