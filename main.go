package main

import (
	"encoding/json"
	"fmt"
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
	for {
		time.Sleep(time.Duration(time.Millisecond * 1000 * 10))
		orders := readConfirmedOrder(time.Millisecond * 1000 * 10)
		fmt.Println(len(orders))
		for _, order := range orders {
			fmt.Println(order)

			acceptRatingDataJson, _ := json.Marshal(CallbackRatingData{
				Type: AcceptRating,
				Id:   order.Id,
			})

			rejectRatingDataJson, _ := json.Marshal(CallbackRatingData{
				Type: RejectRating,
				Id:   order.Id,
			})

			archiveOrderDataJson, _ := json.Marshal(CallbackRatingData{
				Type: ArchiveOrder,
				Id:   order.Id,
			})

			checkExecutedOrderMessage := tgbotapi.NewMessage(order.CustomerId, "Ваш"+" [заказ](https://t.me/krumos/"+fmt.Sprint(order.MessageId)+") уже выполнен?")

			RatingButtonConfig := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Да", string(acceptRatingDataJson)),
					tgbotapi.NewInlineKeyboardButtonData("Нет", string(rejectRatingDataJson)),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Больше не спрашивать", string(archiveOrderDataJson)),
				))
			checkExecutedOrderMessage.ReplyMarkup = RatingButtonConfig
			checkExecutedOrderMessage.ParseMode = tgbotapi.ModeMarkdownV2
			bot.Send(checkExecutedOrderMessage)
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
