package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OrderState int64

const (
	TitleOrderState OrderState = iota + 1
	DescriptionOrderState
	DeadlineOrderState
	PriceOrderState
	FilesOrderState
	ModeratedOrderState
	ConfirmedOrderState
	RejectedOrderState
	ExecutedOrderState
)

type OrderData struct {
	Id           int64               `pg:"id,pk"`
	Title        string              `pg:"title"`
	Description  string              `pg:"description"`
	FilesURL     string              `pg:"files_url"`
	CustomerId   int64               `pg:"customer_id,notnull"`
	ExecutorId   int64               `pg:"executor_id"`
	MessageId    int                 `pg:"message_id"`
	State        OrderState          `pg:"state,notnull"`
	DeadlineDate string              `pg:"deadline_date"`
	Price        string              `pg:"price"`
}

type OrderFile struct {
	OrderId  int64  `pg:"id,nopk"`
	FileId   string `pg:"file_id"`
	FileType string `pg:"file_type"`
}

func createFileOrder(orderFile *OrderFile) {

	exsists, _ := db.Model(orderFile).Where("id=? AND file_id=?", orderFile.OrderId, orderFile.FileId).Exists()
	if !exsists {
		db.Model(orderFile).Insert()
	} 
}

func readFileOrder(orderId int64) []OrderFile {
	var orders []OrderFile
	db.Model((*OrderFile)(nil)).Where("id=?", orderId).Select(&orders)
	return orders
}

func ordersStateMachine(user *UserData, update tgbotapi.Update, config *Config) {
	order := readOrderByState(user.Id)
	switch order.State {
	case TitleOrderState:
		newHeaderOrderCommand(update, user, order)
	case DescriptionOrderState:
		newDescriptionrOrderCommand(update, user, order, config)
	case DeadlineOrderState:
		newDeadlineOrderCommand(update, user, order, config)
	case PriceOrderState:
		newPriceOrderCommand(update, user, order, config)
	case FilesOrderState:
		newFilesOrderCommand(update, user, order, config)
	}
}

func (order *OrderData) toTelegramString() string {
	text := "[ ](" + order.FilesURL + ")\n" + 
		"Дисциплина: *" + order.Title + "* \n\n" +
		"Описание заказа: " + order.Description + "\n\n" +
		"Дедлайн: *" + order.DeadlineDate + "*\n\n" +
		"Цена: *" + order.Price + "*"

	return text
}

func (order *OrderData) toTelegramMediaConfig(chatId int64) tgbotapi.MediaGroupConfig {
	// text := "Дисциплина: *" + order.Title + "* \n\n" +
	// 	"Описание заказа: " + order.Description + "\n\n" +
	// 	"Дедлайн: *" + order.DeadlineDate + "*\n\n" +
	// 	"Цена: *" + order.Price + "*"

	files := readFileOrder(order.Id)

	// f := make([]interface{}, len(files))

	// for i, file := range files {
	// 	fi := tgbotapi.FileID(file.FileId)
	// 	if file.FileType == "photo" {
	// 		photo := tgbotapi.NewInputMediaPhoto(fi)
	// 		if i == 0 {
	// 			photo.Caption = text
	// 		}
	// 		f = append(f, photo)

	// 	}
	// 	if file.FileType == "doc" {
	// 		document := tgbotapi.NewInputMediaDocument(fi)
	// 		if i == 0 {
	// 			document.Caption = text
	// 		}
	// 		f = append(f, document)
	// 	}
	// }
	msg := tgbotapi.NewMediaGroup(chatId, []interface{}{tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(files[0].FileId))})
	return msg
}
