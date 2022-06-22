package main

import (
	"encoding/json"
	"io/ioutil"

	//"slices"
	//"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getConfig() (config *Config, err error) {
	d, e := ioutil.ReadFile("config.json")
	if e != nil {
		config = nil
		return nil, e
	}
	e = json.Unmarshal(d, &config)
	return config, e
}

func connectBot(token string) (bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, err error) {
	bot, e := tgbotapi.NewBotAPI(token)
	if e != nil {
		return nil, nil, e
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates = bot.GetUpdatesChan(updateConfig)
	return bot, updates, e
}

func main() {

	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	bot, updates, err := connectBot(config.Token)
	if err != nil {
		panic(err)
	}

	connectDB(":5432", "postgres", "password", "postgres")

	for update := range updates {
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "post":
				approveOrderResponse(config, update, bot)
				// TODO: Сделать уведомление юзера об отказе в посте
			case "del":
				rejectOrderResponse(config, update, bot)
			}
		}
		if update.Message != nil {
			switch update.Message.Text {
			case "/start": // TODO: Check if user in DB
				startCommand(update, bot)
			case "/new_order":
				newOrderCommand(update, bot)
			default:
				// user := usersList[slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })]
				// switch user.state {

				// case title:
				// 	msg = tgbotapi.NewMessage(update.Message.From.ID, "Напиши описание")
				// 	user := usersList[slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })]
				// 	order := OrderData{customerId: user.id, state: title}
				// 	ordersList = append(ordersList, order)

				// }
				// msg = tgbotapi.NewMessage(update.Message.From.ID, "я тебя не понимать")
			}
			// msg := tgbotapi.NewMessage(moderatorChat, update.Message.Text)
			// 	btn := tgbotapi.NewInlineKeyboardMarkup(
			// 		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запостить", "post"),
			// 			tgbotapi.NewInlineKeyboardButtonData("Удалить", "del")))
			// 	msg.ReplyMarkup = btn
			// 	bot.Send(msg)
		}
	}
}
