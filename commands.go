package main

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandsStateMachine(update tgbotapi.Update, config *Config) {
	switch update.Message.Text {
	case Texts["start_command"]:
		startCommand(update)
	case Texts["new_order_command"]:
		newOrderCommand(update)
	default:
		userStateMachine(update, config)
	}
}

//Команда при запуске бота
func startCommand(update tgbotapi.Update) {
	//Регистрация юзера
	createUser(&UserData{Id: update.Message.From.ID})

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["start_command_answer"]))
}

//Команда создания новго закаща
func newOrderCommand(update tgbotapi.Update) {
	user := readUser(update.Message.From.ID)
	user.State = MakingOrderUserState
	updateUser(&user)

	createOrder(&OrderData{CustomerId: user.Id, State: SubjectInputOrderState})

	bot.Send(tgbotapi.NewMessage(user.Id, Texts["new_order_command_answer"]))
}

//Ввод предметной области заказа
func newHeaderInputOrder(update tgbotapi.Update, user *UserData, order *OrderData) {
	order.Subject = update.Message.Text
	order.State = DescriptionInputOrderState
	updateOrder(order)

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["new_header_command_an"]))
}

//Ввод описания заказа
func newDescriptionInputOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.Description = update.Message.Text
	order.State = DeadlineInputOrderState
	updateOrder(order)

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["new_description_command_an"]))
}

//Ввод даты дедлайна
func newDeadlineInputOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.DeadlineDate = update.Message.Text
	order.State = PriceInputOrderState
	updateOrder(order)

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["new_deadline_command_an"]))
}

//Ввод цены заказа
func newPriceInputOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	order.Price = update.Message.Text
	order.State = FilesUploadOrderState
	updateOrder(order)

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["new_price_command_an"]))
}

//Прикрепление файлов
func newFilesUploadOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	if update.Message.Text == "STOP" { //TODO: забабахать кнопку
		completeOrder(update, user, order, config)
	}

	if update.Message.Photo != nil {
		file := OrderFile{
			FileId:   update.Message.Photo[0].FileID,
			OrderId:  order.Id,
			FileType: "photo",
		}
		createFileOrder(&file)
	}
}

//Получаем медиагруппу из фотографий заказа
func getOrderMediaGroup(order *OrderData, config *Config) tgbotapi.MediaGroupConfig {
	files := readFileOrder(order.Id)
	photos := make([]interface{}, 0)
	for _, file := range files {
		photos = append(photos, tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(file.FileId)))
	}
	return tgbotapi.NewMediaGroup(config.PhotoChat, photos)
}

//Формирование заказа
func completeOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	user.State = DefaultUserState
	updateUser(user)

	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, Texts["order_created_command_an"]))

	OrderMediaPost, _ := bot.SendMediaGroup(getOrderMediaGroup(order, config))

	files_url := "https://t.me/krumos_photo/" + fmt.Sprint(OrderMediaPost[0].MessageID)
	order.FilesURL = files_url
	updateOrder(order)

	approveDataJson, _ := json.Marshal(CallbackData{
		Type: Approve,
		Id:   order.Id,
	})

	rejectDataJson, _ := json.Marshal(CallbackData{
		Type: Reject,
		Id:   order.Id,
	})

	ModeratorButtonConfig := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Texts["approve_button"], string(approveDataJson)),
			tgbotapi.NewInlineKeyboardButtonData(Texts["reject_button"], string(rejectDataJson)),
		))

	OrderToModeratorChatMessage := tgbotapi.NewMessage(config.ModeratorChat, order.toTelegramString())
	OrderToModeratorChatMessage.ParseMode = tgbotapi.ModeMarkdownV2
	OrderToModeratorChatMessage.ReplyMarkup = ModeratorButtonConfig
	bot.Send(OrderToModeratorChatMessage)
}
