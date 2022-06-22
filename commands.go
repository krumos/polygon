package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func approveOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(config.ChannelChat, update.CallbackQuery.Message.Text)
	bot.Send(msg)
	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}

func rejectOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}

func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Привет")
	user := UserData{id: update.Message.From.ID}
	createUser(user)
	bot.Send(msg)
}

func newOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Напиши заголовок")
	//user := usersList[0] /* slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })] */
	// order := OrderData{customerId: user.id, state: title}
	// ordersList = append(ordersList, order)
	bot.Send(msg)
}
