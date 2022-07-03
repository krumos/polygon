package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func connectBot(token string) (updates tgbotapi.UpdatesChannel, err error) {
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates = bot.GetUpdatesChan(updateConfig)
	return updates, err
}

func main() {
	s := "...."
	character := s[2]
	fmt.Println(string(character))
	updates, _ := connectBot("5473842943:AAEFw50U83kPXAzxikPu21RaoYrX4diclAY")

	for update := range updates {
		// photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID("AgACAgIAAxkBAAICn2K9hSSzhf_JD_EocdPGTXiTucu-AAISvDEbP9HwSVJcBDeCSlSwAQADAgADeAADKQQ"))
		// photo.Caption = "Привет"

		// photo2 := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID("AgACAgIAAxkBAAICn2K9hSSzhf_JD_EocdPGTXiTucu-AAISvDEbP9HwSVJcBDeCSlSwAQADAgADeAADKQQ"))
		// mediaGroup := tgbotapi.NewMediaGroup(update.Message.From.ID, []interface{}{photo, photo2})

		// bot.Send(mediaGroup)
		text := "[ ](https://t.me/krumos_photo/5)" + update.Message.Text
		msg3 := tgbotapi.NewMessage(-1001468825053, text)
		msg3.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg3)

		// document:= tgbotapi.NewInputMediaDocument(tgbotapi.FileID("BQACAgIAAxkBAAICuGK9lVsj4pSraZqReCMbu73yv6lGAAIiHQACP9HwSV3hKEwmySzPKQQ"))
		// document.Caption = "Привет"
		// tgbotapi.NewMediaGroup(update.Message.From.ID, []interface{}{ document })

		//bot.Send(mediaGroup)
		// media := update.Message.MediaGroupID
		// u := bot.Get

	}
}
