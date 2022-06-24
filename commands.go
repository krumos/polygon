package main

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func agreementOrderResponse(update tgbotapi.Update, bot *tgbotapi.BotAPI, response *CallbackData) {
	order := readOrderById(response.Id)

	//Можно откликнуться миллиард раз
	msg := tgbotapi.NewMessage(order.CustomerId, Texts["agreed_order"])
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(Texts["go_to_chat_button"], "https://t.me/" + update.CallbackQuery.From.UserName)))
	msg.ReplyMarkup = btn
	bot.Send(msg)
}

func approveOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ApprovedOrderState

	updateOrder(&order)
	agreementData := CallbackData{
		Type: Agreement,
		Id:   order.Id,
	} // , json.Unmarshall()
	agreementDataJson, _ := json.Marshal(agreementData)

	msg := tgbotapi.NewMessage(config.ChannelChat, order.toString())
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Texts["respond_order"], string(agreementDataJson))))
	msg.ReplyMarkup = btn
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)

	msg = tgbotapi.NewMessage(order.CustomerId, Texts["status_sent"])
	bot.Send(msg)

	msgRej := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msgRej)
}

func rejectOrderResponse(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ApprovedOrderState

	updateOrder(&order)

	msg := tgbotapi.NewMessage(order.CustomerId, Texts["status_rejected"])
	bot.Send(msg)

	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}

func startCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := UserData{Id: update.Message.From.ID}
	createUser(&user)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["hello"])
	bot.Send(msg)
}

func newOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := readUser(update.Message.From.ID)
	user.State = MakingOrderUserState
	updateUser(&user)

	order := OrderData{CustomerId: user.Id, State: TitleOrderState}
	createOrder(&order)

	msg := tgbotapi.NewMessage(user.Id, Texts["header_notifier"])
	bot.Send(msg)
}

func newHeaderOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	order.Title = update.Message.Text
	order.State = DescriptionOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["description_notifier"])
	bot.Send(msg)
}

func newDescriptionrOrderCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, user *UserData, order *OrderData) {
	user.State = DefaultUserState
	updateUser(user)

	order.Description = update.Message.Text
	order.State = ModeratedOrderState
	updateOrder(order)

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["status_moderating"])
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
