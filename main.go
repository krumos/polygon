package main

import (
	"time"

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

func spamer() {
	h := 48
	for {
		time.Sleep(time.Duration(time.Hour * time.Duration(h)))
		orders := readConfirmedOrder()
		for _, order := range orders {
			bot.Send(tgbotapi.NewMessage(order.CustomerId, "Ваш заказнейм уже выполнен?")) //TODO клавиатурка
		}
	}
}

func main() {
	getText("texts.yaml")

	var config Config
	err := config.getConfig()
	if err != nil {
		panic(err)
	}

	updates, err := connectBot(config.Token)
	if err != nil {
		panic(err)
	}

	connectDB(":5432", "postgres", "password", "postgres")

	go spamer()

	for update := range updates {
		if update.CallbackQuery != nil {
			// прилетел апдейт с инлайн кнопок
			responseStateMachine(update, &config)
		}
		if update.Message != nil {
			// обычное сообщение
			commandsStateMachine(update, &config)
		}
	}
}
