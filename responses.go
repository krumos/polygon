package main

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func responseStateMachine(update tgbotapi.Update, config *Config) {
	response := CallbackData{}
	json.Unmarshal([]byte(update.CallbackQuery.Data), &response)

	order := readOrderById(response.Id)

	switch response.Type {
	// -------- результат модерации --------
	case Approve:
		аpproveOrderResponse(config, update, &response, &order)
	case Reject:
		rejectOrderResponse(config, update, &response, &order)
	// -------- отклик на заказ --------
	case Agreement:
		agreementOrderResponse(update, &response, &order)
	// -------- подтверждение исполнителя -------
	case Confirm:
		confirmOrderResponse(update, &response, config, &order)
	}
	updateOrder(&order)
}

//Заказчик выбрал исполнителя
func confirmOrderResponse(update tgbotapi.Update, response *CallbackData, config *Config, order *OrderData) {
	order.ExecutorId = response.ExecutorId
	order.State = ConfirmedOrderState

	//Удаление кнопки "отклика" с поста
	orderWOKeyboardPost := tgbotapi.NewEditMessageText(config.ChannelChat, order.MessageId, order.toTelegramString())
	orderWOKeyboardPost.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Send(orderWOKeyboardPost)

	//Колбэк заказчику
	bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Теперь на Ваш заказ нельзя откликнуться"))
}

//Подает заявку(фрилансер) на выполнение заказа
func agreementOrderResponse(update tgbotapi.Update, response *CallbackData, order *OrderData) {
	orderCallback := OrderCallback{
		Id:          order.Id,
		ResponderId: update.CallbackQuery.From.ID,
	}
	//Проверяем, откликался ли человек ранее
	if isExistsOrderCallback(&orderCallback) {
		bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы уже откликнулись"))
		return
	}

	//InlineButtonData для кнопки подтверждения выбора исполнителя
	confirmDataJson, _ := json.Marshal(CallbackData{
		Type:       Confirm,
		Id:         order.Id,
		ExecutorId: update.CallbackQuery.From.ID,
	})

	//Уведомляем заказчика об отклике
	//Отправляем сообщение создателю заказа с кнопками "перейти в чат" и "выбрать исполнителя"
	orderAgreementMessage := tgbotapi.NewMessage(order.CustomerId, Texts["agreed_order"])
	orderAgreementMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(Texts["go_to_chat_button"], "https://t.me/"+update.CallbackQuery.From.UserName),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Texts["choose_responder_button"], string(confirmDataJson)),
		))
	m, _ := bot.Send(orderAgreementMessage)

	orderCallback.MessageId = m.MessageID
	createOrderCallback(&orderCallback)
}

//Когда пройдена модерация
func аpproveOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData, order *OrderData) {
	order.State = ConfirmedOrderState

	//InlineButtonData для кнопки отклика на заказ
	agreementDataJson, _ := json.Marshal(CallbackData{
		Type: Agreement,
		Id:   order.Id,
	})

	//Постим отмодерированый заказ в канал с клавиатурой отклика
	orderWithKeyboardPost := tgbotapi.NewMessage(config.ChannelChat, order.toTelegramString())
	orderWithKeyboardPost.ParseMode = tgbotapi.ModeMarkdownV2
	orderWithKeyboardPost.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Texts["respond_order"], string(agreementDataJson))))
	m, _ := bot.Send(orderWithKeyboardPost)
	order.MessageId = m.MessageID

	//Сообщаем заказчику о том что заказ прошел модерацию
	orderModeratedMessage := tgbotapi.NewMessage(order.CustomerId, Texts["status_sent"])
	bot.Send(orderModeratedMessage)

	//Удаляем пост из канала модераторов
	bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
}

//Модерация не пройдена
func rejectOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData, order *OrderData) {
	order.State = RejectedOrderState

	//Сообщаем заказчику о том что заказ не прошел модерацию
	bot.Send(tgbotapi.NewMessage(order.CustomerId, Texts["status_rejected"]))

	//Удаляем пост из канала модераторов
	bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
}
