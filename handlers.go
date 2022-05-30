package main

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"strconv"
	"strings"
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
	results = append(results, articleWithMenu(inlineQuery.ID, inlineQuery.From.ID, contests[0][settings.Featured], chatType))

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
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ParseMode = "MarkdownV2"
	if creator, ok := creators[message.From.ID]; ok {
		creatorOld := creator
		if message.IsCommand() {
			if (strings.HasPrefix(message.Command(), "t") || strings.HasPrefix(message.Command(), "p")) &&
				creator.StateContest != 0 {
				m := map[string]string{"post": "поста", "guide": "инструкции", "private": "ЛС", "group": "чата", "channel": "канала"}
				if v, ok := m[message.Command()[1:]]; ok {
					creator.StateField = message.Command()
					msg.Text = "Пришлите "
					if strings.HasPrefix(message.Command(), "t") {
						msg.Text += "текст для "
					} else {
						msg.Text += "картинку для "
					}
					msg.Text += v
				}
			} else {
				switch message.Command() {
				case "inline":
					msg.Text = "Однажды добавим\\.\\."

				case "s": //select
					if args := strings.Fields(message.CommandArguments()); len(args) > 0 {
						creator.StateContest, _ = strconv.ParseInt(args[0], 10, 0)
					}
					if contest, ok := contests[message.From.ID][creator.StateContest]; ok {
						msg.Text = "Розыгрыш: `" + contest.ContestName + "`"
						creator.StateField = ""
					} else {
						creator.StateContest = 0
						msg.Text = "Розыгрыш не найден"
					}

				case "r": //rename
					if creator.StateContest != 0 {
						creator.StateField = "rename"
						msg.Text = "Введите новое название розыгрыша:"
					} else {
						msg.Text = "Розыгрыш не выбран"
					}

				case "d": //delete
					if creator.StateContest != 0 {
						creator.StateField = "delete"
						msg.Text = "Введите название розыгрыша для подтверждения: `" +
							contests[message.From.ID][creator.StateContest].ContestName + "`"
						if contests[message.From.ID][creator.StateContest].ContestName == "" {
							msg.Text += "\nНазвание розыгрыша пустое, отправьте \\+ для подтверждения"
						}
					} else {
						msg.Text = "Розыгрыш не выбран"
					}

				case "e": //enable
					if creator.StateContest != 0 {
						if args := strings.Fields(message.CommandArguments()); len(args) > 0 {
							if val, err := strconv.ParseInt(args[0], 10, 0); err == nil {
								contest := contests[message.From.ID][creator.StateContest]
								contest.ContestActive = int(val)
								err := contests.update(db, contest)
								if err != nil {
									log.Println(err)
								}
								if contest.ContestActive == 0 {
									msg.Text = "Конкурс остановлен"
								} else {
									msg.Text = "Конкурс запущен"
								}
							}
						} else {
							msg.Text = "Неизвестное состояние"
						}
					} else {
						msg.Text = "Розыгрыш не выбран"
					}

				case "b": //back
					creator.StateContest = 0
					creator.StateField = ""
					msg.Text = "Главное меню:"

				case "l":
					if creator.StateContest != 0 {
						m, ok := checkChannelPrivileges(bot, contests[message.From.ID][creator.StateContest].Username)
						if ok {
							if _, err := bot.Send(postToChannel(contests[message.From.ID][creator.StateContest])); err != nil {
								log.Println(err)
								return
							}
							msg.Text = "Опубликовано"
						} else {
							msg.Text = m
						}
					} else {
						msg.Text = "Розыгрыш не выбран"
					}

				case "n": //new
					contest := Contest{
						CreatorID: message.From.ID,
					}
					contestID, err := contests.insert(db, contest)
					if err != nil {
						log.Println(err)
					}
					creator.StateContest = contestID
					creator.StateField = "name"
					msg.Text = "Конкурс успешно создан"
					msg.Text += "\nВведите название конкурса:"

				default:
					msg.Text = "Команда не распознана"
					creator.StateField = ""
				}
			}
		} else {
			next := ""
			if contest, ok := contests[message.From.ID][creator.StateContest]; ok {
				contestOld := contest

				if (strings.HasPrefix(creator.StateField, "t") || strings.HasPrefix(creator.StateField, "p")) &&
					creator.StateContest != 0 {
					if _, ok := posts[creator.StateContest][creator.StateField[1:]]; !ok {
						post := Post{
							ContestID: creator.StateContest,
							Type:      creator.StateField[1:],
						}
						err := posts.insert(db, post)
						if err != nil {
							log.Println(err)
						}
					}
					post := posts[creator.StateContest][creator.StateField[1:]]
					if strings.HasPrefix(creator.StateField, "t") {
						post.Message = message.Text
						msg.Text = "Значение установлено"
					} else {
						if message.Document != nil {
							url, err := RepostToTelegraph(bot, message.Document.FileID)
							if err != nil {
								log.Println(err)
								return
							}
							post.Image = url
							msg.Text = "Значение установлено"
						} else if message.Photo != nil {
							url, err := RepostToTelegraph(bot, message.Photo[len(message.Photo)-1].FileID)
							if err != nil {
								log.Println(err)
								return
							}
							post.Image = url
							msg.Text = "Значение установлено"
						} else {
							msg.Text = "Повторите отправку"
							next = creator.StateField
						}
					}
					err := posts.update(db, post)
					if err != nil {
						log.Println(err)
					}
				} else {

					switch creator.StateField {
					case "rename":
						contest.ContestName = message.Text
						msg.Text = "Название обновлено: `" + contest.ContestName + "`"

					case "name":
						contest.ContestName = message.Text
						msg.Text = "Название установлено: `" + contest.ContestName + "`"
						msg.Text += "\nВведите имя канала в формате `@channel` или `https://t.me/channel`:"
						next = "username"

					case "username":
						if strings.HasPrefix(message.Text, "@") {
							contest.Username = message.Text[1:]
						}
						if strings.HasPrefix(message.Text, "https://t.me/") {
							contest.Username = message.Text[13:]
						}
						if contest.Username != "" {
							msg.Text = "Имя канала установлено: `" + contest.Username + "`"
						} else {
							msg.Text = "Повторите попытку"
							next = "username"
						}

					case "delete":
						if contest.ContestName == message.Text || (contest.ContestName == "" && message.Text == "+") {
							err := contests.delete(db, contest)
							if err != nil {
								log.Println(err)
							}
							msg.Text = "Розыгрыш удален"
							creator.StateContest = 0
						} else {
							next = "delete"
							msg.Text = "Вы допустили ошибку в названии"
						}

					default:
						msg.Text = "Параметр не выбран"
					}
					if contest != contestOld {
						err := contests.update(db, contest)
						if err != nil {
							log.Println(err)
						}
					}
				}
			} else {
				msg.Text = "Розыгрыш не выбран"
			}
			creator.StateField = next
		}
		if creator != creatorOld {
			err := creators.update(db, creator)
			if err != nil {
				log.Println(err)
			}
		}
		if creator.StateContest != 0 && creator.StateField != "" {
			//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/s " +
						strconv.FormatInt(creator.StateContest, 10) + " Назад"),
				),
			)
			msg.ReplyMarkup = keyboard
		} else {
			if creator.StateContest != 0 {
				row := tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/b Меню"),
					tgbotapi.NewKeyboardButton("/r Переименовать ["+
						contests[message.From.ID][creator.StateContest].ContestName+"]"),
					tgbotapi.NewKeyboardButton("/d Удалить"),
				)
				if contests[message.From.ID][creator.StateContest].ContestActive == 1 {
					row = append(row, tgbotapi.NewKeyboardButton("/e 0 Остановить"))
				} else {
					row = append(row, tgbotapi.NewKeyboardButton("/e 1 Запустить"))
				}
				keyboard := tgbotapi.NewReplyKeyboard(row)
				texts := tgbotapi.NewKeyboardButtonRow()
				pics := tgbotapi.NewKeyboardButtonRow()

				for key, val := range map[string]string{"post": "пост", "guide": "инструкция", "private": "ЛС", "group": "чат", "channel": "канал"} {
					texts = append(texts, tgbotapi.NewKeyboardButton("/t"+key+" Текст: "+val))
					pics = append(pics, tgbotapi.NewKeyboardButton("/p"+key+" Картинка: "+val))
				}
				keyboard.Keyboard = append(keyboard.Keyboard, texts)
				keyboard.Keyboard = append(keyboard.Keyboard, pics)
				keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/l Запостить")))

				msg.ReplyMarkup = keyboard
			} else {
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("/n Создать ➕"),
						//tgbotapi.NewKeyboardButton("/inline mode"),
					),
				)
				if len(contests[message.From.ID]) > 0 {
					var table [][]tgbotapi.KeyboardButton
					counter := 0
					div := int(math.Round(math.Sqrt(float64(len(contests[message.From.ID])))))
					for _, v := range contests[message.From.ID] {
						if len(table) <= counter/div {
							table = append(table, tgbotapi.NewKeyboardButtonRow())
						}
						table[counter/div] = append(table[counter/div], tgbotapi.NewKeyboardButton("/s "+
							strconv.FormatInt(v.ContestID, 10)+" ["+v.ContestName+"]"))
						counter += 1
					}
					keyboard.ResizeKeyboard = true
					keyboard.Keyboard = append(keyboard.Keyboard, table...)

				}
				msg.ReplyMarkup = keyboard
			}
		}
	} else {
		if settings.IsPublic == 2 || message.From.ID == settings.Admin {
			msg.Text = "Вы можете создать реферальный розыгрыш"
			creator := Creator{
				CreatorID: message.From.ID,
				UserName:  message.From.FirstName,
			}
			if message.From.LastName != "" {
				creator.UserName += " " + message.From.LastName
			}
			if message.From.UserName != "" {
				creator.UserName += " @" + message.From.UserName
			}
			err := creators.insert(db, creator)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	if msg.Text != "" {
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
	}
	if false {
		msg := tgbotapi.NewMessage(message.Chat.ID, "")
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		msg.Text += "\nCommand:" + message.Command()
		msg.Text += "\nArguments:" + message.CommandArguments()
		msg.Text += "\nCommandWithAt:" + message.CommandWithAt()
		if message.From.ID == settings.Admin {
			switch message.Text {
			case "/post":
				if _, err := bot.Send(postToChannel(contests[0][settings.Featured])); err != nil {
					log.Println(err)
					return
				}
			case "/check":
				memo, ok := checkChannelPrivileges(bot, "g1pnk")
				msg.Text += "\nResult: " + memo
				if ok {
					msg.Text += "\nPrivileges: True"

				} else {
					msg.Text += "\nPrivileges: False"

				}
			}
		}

		if message.Document != nil {
			url, err := RepostToTelegraph(bot, message.Document.FileID)
			if err != nil {
				log.Println(err)
				return
			}
			msg.Text += "\n" + url
		}

		if message.Photo != nil {
			url, err := RepostToTelegraph(bot, message.Photo[len(message.Photo)-1].FileID)
			if err != nil {
				log.Println(err)
				return
			}
			msg.Text += "\n" + url
		}

		// Send the message.
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func CallbackHandler(bot *tgbotapi.BotAPI, db *sql.DB, callback *tgbotapi.CallbackQuery) {
	var inlineData InlineData
	responseMessage := "Неизвестная ошибка!"
	if err := json.Unmarshal([]byte(callback.Data), &inlineData); err != nil {
		log.Println(err)
	} else {
		if contest, ok := contests[0][inlineData.Contest]; ok {
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
						FromID:       inlineData.From,
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
