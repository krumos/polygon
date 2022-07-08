package main

import (
	"encoding/json"
	"fmt"
	"time"

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
	case AcceptRating:
		acceptRatingOrderResponse(update, &order)
	case RejectRating:
		rejectRatingOrderResponse(update)
	case RatingExecutor:
		user := readUser(order.ExecutorId)
		ratingOrderResponse(update, &order, user, ratingExecutor)
	case RatingCustomer:
		user := readUser(order.CustomerId)
		ratingOrderResponse(update, &order, user, ratingCustomer)
	}
	updateOrder(&order)
}

func ratingCustomer(user *UserData, mark int32) {
	user.CustomerRatingCount++
	user.CustomerRatingSum += mark
}

func ratingExecutor(user *UserData, mark int32) {
	user.ExecutorRatingCount++
	user.ExecutorRatingSum += mark
}

func ratingOrderResponse(update tgbotapi.Update, order *OrderData, user *UserData, f func(*UserData, int32)) {
	response := CallbackRatingData{}
	json.Unmarshal([]byte(update.CallbackQuery.Data), &response)

	f(user, response.Mark)

	updateUser(user)

	bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID))

	bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Спасибо за оценку"))
}

func getRatingKeyboard(ratingType CallbackDataType, order *OrderData) (keyboard []tgbotapi.InlineKeyboardButton) {
	for i := 0; i < 5; i++ {
		ratingData, _ := json.Marshal(CallbackRatingData{
			Type: ratingType,
			Id:   order.Id,
			Mark: int32(i + 1),
		})
		keyboard[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(i+1), string(ratingData))
	}
	return keyboard
}

func acceptRatingOrderResponse(update tgbotapi.Update, order *OrderData) {
	order.State = ExecutedOrderState
	updateOrder(order)

	RatingButtonConfig := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(getRatingKeyboard(RatingExecutor, order)...))

	ratingOrderMessage := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, RatingButtonConfig)

	bot.Send(ratingOrderMessage)

	RatingButtonConfig = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(getRatingKeyboard(RatingCustomer, order)...))

	ratingOrderTextMessage := tgbotapi.NewMessage(order.ExecutorId, "Оцените работу с заказчиком")
	ratingOrderTextMessage.ReplyMarkup = RatingButtonConfig
	bot.Send(ratingOrderTextMessage)
}

func rejectRatingOrderResponse(update tgbotapi.Update) {
	bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID))

	bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Спасибо, мы спросим Вас позже"))
}

//Заказчик выбрал исполнителя
func confirmOrderResponse(update tgbotapi.Update, response *CallbackData, config *Config, order *OrderData) {
	order.ExecutorId = response.ExecutorId
	order.State = ConfirmedOrderState
	order.ConfirmationTime = time.Now()

	//Удаление кнопки "отклика" с поста
	orderWOKeyboardPost := tgbotapi.NewEditMessageText(config.ChannelChat, order.MessageId, order.toTelegramString())
	orderWOKeyboardPost.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Send(orderWOKeyboardPost)

	//Колбэк заказчику
	bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Теперь на Ваш заказ нельзя откликнуться"))

	callbacks := readCallbacksOrder(order)
	fmt.Println(len(callbacks))
	for _, callback := range callbacks {
		bot.Send(tgbotapi.NewDeleteMessage(order.CustomerId, callback.MessageId))
	}
}

func getStringCustomerRating(user *UserData) string { // TODO расширение над юзером?
	text := "\n\nРейтинг исполнителя: "
	if user.ExecutorRatingCount == 0 {
		return text + "оценок еще нет"
	}
	return text + fmt.Sprintf("%.2f", float64(user.ExecutorRatingSum/user.ExecutorRatingCount)) // вообще не уверен что тут всё ок, но пусть так
}

//Подает заявку(фрилансер) на выполнение заказа
func agreementOrderResponse(update tgbotapi.Update, response *CallbackData, order *OrderData) {
	orderCallback := OrderCallback{
		Id:          order.Id,
		ResponderId: update.CallbackQuery.From.ID,
	}
	if order.CustomerId == update.CallbackQuery.From.ID {
		bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы создатель заказа"))
		return
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
	user := readUser(update.CallbackQuery.From.ID) // TODO: Нужно проверять зареган ли чел в боте. Потому что нам нужно подгружать рейтинг. Возможно его можно сразу регистрировать тут
	if user.Id == 0 {
		user = &UserData{
			Id:    update.CallbackQuery.From.ID,
			State: DefaultUserState,
		}
		createUser(user)
	}
	//Уведомляем заказчика об отклике
	//Отправляем сообщение создателю заказа с кнопками "перейти в чат" и "выбрать исполнителя"
	orderAgreementMessage := tgbotapi.NewMessage(order.CustomerId, "Ваш [заказ]("+"https://t.me/krumos/"+fmt.Sprint(order.MessageId)+") хочет выполнить"+
		" **"+toExcapedString(update.CallbackQuery.From.UserName)+"**"+
		toExcapedString(getStringCustomerRating(user))) //отвратително

	orderAgreementMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(Texts["go_to_chat_button"], "https://t.me/"+update.CallbackQuery.From.UserName),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Texts["choose_responder_button"], string(confirmDataJson)),
		))
	orderAgreementMessage.ParseMode = tgbotapi.ModeMarkdownV2 // нужно экранировать юзернейм
	m, _ := bot.Send(orderAgreementMessage)

	orderCallback.MessageId = m.MessageID
	createOrderCallback(&orderCallback)
}

//Когда пройдена модерация
func аpproveOrderResponse(config *Config, update tgbotapi.Update, response *CallbackData, order *OrderData) {
	if order.State == ConfirmedOrderState {
		return
	}
	order.State = ConfirmedOrderState
	updateOrder(order)

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
