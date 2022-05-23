package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var settings Settings
var queries Queries
var messages Messages
var creators Creators
var members Members
var posts Posts
var contests Contests

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

	if err := DBLoad(db); err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(settings.Token)
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
			go InlineHandler(bot, db, update.InlineQuery)
		} else if update.Message != nil {
			go MessageHandler(bot, db, update.Message)
		} else if update.CallbackQuery != nil {
			go CallbackHandler(bot, db, update.CallbackQuery)
		} else if update.ChosenInlineResult != nil {
			go ChosenInlineResultHandler(bot, db, update.ChosenInlineResult)
		} else if update.MyChatMember != nil {
			go func() {}()
		}
	}
}

func DBLoad(db *sql.DB) error {
	if err := DBInit(db); err != nil {
		return err
	}
	if err := settings.load(db); err != nil {
		return err
	}
	if err := queries.load(db); err != nil {
		return err
	}
	if err := messages.load(db); err != nil {
		return err
	}
	if err := creators.load(db); err != nil {
		return err
	}
	if err := members.load(db); err != nil {
		return err
	}
	if err := posts.load(db); err != nil {
		return err
	}
	if err := contests.load(db); err != nil {
		return err
	}
	return nil
}
