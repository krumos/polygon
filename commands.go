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

	bot.Send(StartMessage(update.Message.From.ID))
}

//Команда создания новго закаща
func newOrderCommand(update tgbotapi.Update) {
	user := readUser(update.Message.From.ID)
	user.State = MakingOrderUserState
	updateUser(&user)

	createOrder(&OrderData{CustomerId: user.Id, State: SubjectInputOrderState})

	msg := tgbotapi.NewMessage(user.Id, Texts["new_order_command_answer"])
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	bot.Send(msg)
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

	bot.Send(FilesUploadMessage(update.Message.From.ID))
}

//Прикрепление файлов
func newFilesUploadOrder(update tgbotapi.Update, user *UserData, order *OrderData, config *Config) {
	if update.Message.Text == Texts["stop"] {
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
func getOrderMediaGroup(files []OrderFile, config *Config) tgbotapi.MediaGroupConfig {
	
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

	msg := tgbotapi.NewMessage(update.Message.From.ID, Texts["order_created_command_an"])
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Texts["new_order_command"])))
	bot.Send(msg)
	
	files := readFileOrder(order.Id)
	var files_url string
	if (len(files) != 0) {
		OrderMediaPost, _ := bot.SendMediaGroup(getOrderMediaGroup(files, config))
		files_url = "https://t.me/krumos_photo/" + fmt.Sprint(OrderMediaPost[0].MessageID)
	} 
	order.State = ModeratedOrderState
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
