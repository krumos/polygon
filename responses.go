package main

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func responseStateMachine(update tgbotapi.Update, config *Config) {
	response := CallbackData{}
	json.Unmarshal([]byte(update.CallbackQuery.Data), &response)
	switch response.Type {
	case Approve:
		аpproveOrderResponse(config, update, &response)
		// TODO: Сделать уведомление юзера об отказе в посте
	case Reject:
		rejectOrderResponse(config, update, &response)
	case Agreement:
		agreementOrderResponse(update, &response)
	}
}

func agreementOrderResponse(update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)

	//Можно откликнуться миллиард раз
	msg := tgbotapi.NewMessage(order.CustomerId, Texts["agreed_order"])
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(Texts["go_to_chat_button"], "https://t.me/"+update.CallbackQuery.From.UserName)))
	msg.ReplyMarkup = btn
	bot.Send(msg)
}

func аpproveOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ApprovedOrderState

	updateOrder(&order)
	agreementData := CallbackData{
		Type: Agreement,
		Id:   order.Id,
	} // , json.Unmarshall()
	agreementDataJson, _ := json.Marshal(agreementData)

	msg := tgbotapi.NewMessage(config.ChannelChat, order.toTelegramString())
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

func rejectOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ApprovedOrderState

	updateOrder(&order)

	msg := tgbotapi.NewMessage(order.CustomerId, Texts["status_rejected"])
	bot.Send(msg)

	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}
