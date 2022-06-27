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
	case Confirm:
		confirmOrderResponse(update, &response, config)
	}
}

func confirmOrderResponse(update tgbotapi.Update, response *CallbackData, config *Config) {
	order := readOrderById(response.Id)
	order.ExecutorId = response.ExecutorId
	order.State = ConfirmedOrderState

	updateOrder(&order)
	media := tgbotapi.NewEditMessageText(config.ChannelChat, int(order.MessageId), order.toTelegramString()/* tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Уже откликнулись", "23"))) */)
	media.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(media)

	m := tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы молодец")
	bot.Request(m)
}

func agreementOrderResponse(update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)

	orderCallback := OrderCallback{
		Id: order.Id,
		ResponderId: update.CallbackQuery.From.ID,
	}

	if isExistsOrderCallback(&orderCallback) {
		m := tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы уже откликнулись")
		bot.Request(m)
		return 
	}

	confirmData := CallbackData{
		Type:       Confirm,
		Id:         order.Id,
		ExecutorId: update.CallbackQuery.From.ID,
	}

	confirmDataJson, _ := json.Marshal(confirmData)

	msg := tgbotapi.NewMessage(order.CustomerId, Texts["agreed_order"])
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(Texts["go_to_chat_button"], "https://t.me/"+update.CallbackQuery.From.UserName)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Texts["choose_responder_button"], string(confirmDataJson))))
	msg.ReplyMarkup = btn

	m, _ := bot.Send(msg)

	orderCallback.MessageId = int64(m.MessageID)
	createOrderCallback(&orderCallback)
}

func аpproveOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ConfirmedOrderState

	agreementData := CallbackData{
		Type: Agreement,
		Id:   order.Id,
	}
	agreementDataJson, _ := json.Marshal(agreementData)

	msg := tgbotapi.NewMessage(config.ChannelChat, order.toTelegramString())
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Texts["respond_order"], string(agreementDataJson))))
	msg.ReplyMarkup = btn
	msg.ParseMode = tgbotapi.ModeMarkdown
	m, _ := bot.Send(msg)
	order.MessageId = m.MessageID

	updateOrder(&order)

	msg = tgbotapi.NewMessage(order.CustomerId, Texts["status_sent"])
	bot.Send(msg)

	msgRej := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msgRej)
}

func rejectOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData) {
	order := readOrderById(response.Id)
	order.State = ConfirmedOrderState

	updateOrder(&order)

	msg := tgbotapi.NewMessage(order.CustomerId, Texts["status_rejected"])
	bot.Send(msg)

	msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	bot.Send(msg2)
}
