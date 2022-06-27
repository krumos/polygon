package main

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandsStateMachine(update tgbotapi.Update, config *Config) {
	switch update.Message.Text {
	case Texts["start_command"]: // TODO: Check if user in DB
		startCommand(update)
	case Texts["new_order_command"]:
		newOrderCommand(update)
	default:
		userStateMachine(update, config)

	}
}

func startCommand(update tgbotapi.Update) {
	user := UserData{Id: update.Message.From.ID}
	createUser(&user)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["start_command_answer"])
	bot.Send(msg)
}

func newOrderCommand(update tgbotapi.Update) {
	user := readUser(update.Message.From.ID)
	user.State = MakingOrderUserState
	updateUser(&user)

	order := OrderData{CustomerId: user.Id, State: TitleOrderState}
	createOrder(&order)

	msg := tgbotapi.NewMessage(user.Id, Texts["new_order_command_answer"])
	bot.Send(msg)
}

func newHeaderOrderCommand(update tgbotapi.Update, user *UserData, order *OrderData) {
	order.Title = update.Message.Text
	order.State = DescriptionOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["new_header_command_an"])
	bot.Send(msg)
}

func newDescriptionrOrderCommand(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.Description = update.Message.Text
	order.State = DeadlineOrderState // создать мметод который будет возвращать следующую фазу заказа
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["new_description_command_an"])
	bot.Send(msg)
}

func newDeadlineOrderCommand(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.DeadlineDate = update.Message.Text
	fmt.Println("============");
	order.State = PriceOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["new_deadline_command_an"])
	bot.Send(msg)
}

func newPriceOrderCommand(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.Price = update.Message.Text
	updateOrder(order)

	complateOrder(update, user, order, config)
}

func complateOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	user.State = DefaultUserState
	updateUser(user)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["order_created_command_an"])
	bot.Send(msg)

	msg = tgbotapi.NewMessage(config.ModeratorChat, order.toTelegramString())
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

