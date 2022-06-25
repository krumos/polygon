package main

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := UserData{Id: update.Message.From.ID}
	createUser(&user)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["start_command_answer"])
	bot.Send(msg)
}

func newOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := readUser(update.Message.From.ID)
	user.State = MakingOrderUserState
	updateUser(&user)

	order := OrderData{CustomerId: user.Id, State: TitleOrderState}
	createOrder(&order)

	msg := tgbotapi.NewMessage(user.Id, Texts["new_order_command_answer"])
	bot.Send(msg)
}

func newHeaderOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	order.Title = update.Message.Text
	order.State = DescriptionOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["new_header_command_an"])
	bot.Send(msg)
}

func newDescriptionrOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	user.State = DefaultUserState
	updateUser(user)

	order.Description = update.Message.Text
	order.State = ModeratedOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["order_created_command_an"])
	bot.Send(msg)

	config, _ := getConfig()

	msg = tgbotapi.NewMessage(config.ModeratorChat, order.toString())
	msg.ParseMode = tgbotapi.ModeMarkdown
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
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Texts["approve_button"], string(approveDataJson)),
			tgbotapi.NewInlineKeyboardButtonData(Texts["reject_button"], string(rejectDataJson))))
	msg.ReplyMarkup = btn
	bot.Send(msg)
}
