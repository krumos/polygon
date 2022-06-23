package main

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func approveOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = PostedOrderState

	updateOrder(&order)

	msg := tgbotapi.NewMessage(config.ChannelChat, update.CallbackQuery.Message.Text)
	bot.Send(msg)
	
	msg = tgbotapi.NewMessage(order.CustomerId, "Ваш заказ отправлен")
	bot.Send(msg)
	
	msgDel := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msgDel)
}

func rejectOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = PostedOrderState

	updateOrder(&order)

	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}

func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Привет")
	user := UserData{Id: update.Message.From.ID}
	createUser(&user)
	bot.Send(msg)
}

func newOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := readUser(update.Message.From.ID)
	order := OrderData{CustomerId: user.Id, State: TitleOrderState}
	user.State = MakingOrderUserState

	createOrder(&order)
	updateUser(&user)

	msg := tgbotapi.NewMessage(user.Id, "Напиши заголовок")
	bot.Send(msg)
}

func newHeaderOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	order.Title = update.Message.Text
	order.State = DescriptionOrderState

	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Напиши описание заказа")
	bot.Send(msg)
}

func newDescriptionrOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	order.Description = update.Message.Text
	order.State = ModeratedOrderState
	user.State = DefaultUserState

	updateUser(user)
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Заказ на модерации")
	bot.Send(msg)

	config, _ := getConfig()

	msg = tgbotapi.NewMessage(config.ModeratorChat, order.Description)

	approveData := CallbackData{
		Type: Approve,
		Id:   order.Id,
	}
	approveDataJson, _ := json.Marshal(approveData)

	rejectData := CallbackData{
		Type: Reject,
		Id:   order.Id,
	} // , json.Unmarshall()
	rejectDataJson, _ := json.Marshal(rejectData)
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запостить", string(approveDataJson)),
			tgbotapi.NewInlineKeyboardButtonData("Удалить", string(rejectDataJson))))
	msg.ReplyMarkup = btn
	bot.Send(msg)
}
