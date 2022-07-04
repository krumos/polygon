package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OrderState int64

const (
	SubjectInputOrderState OrderState = iota + 1
	DescriptionInputOrderState
	DeadlineInputOrderState
	PriceInputOrderState
	FilesUploadOrderState
	ModeratedOrderState
	ConfirmedOrderState
	RejectedOrderState
	ExecutedOrderState
)

type OrderData struct {
	Id           int64      `pg:"id,pk"`
	Subject      string     `pg:"subject"`
	Description  string     `pg:"description"`
	FilesURL     string     `pg:"files_url"`
	CustomerId   int64      `pg:"customer_id,notnull"`
	ExecutorId   int64      `pg:"executor_id"`
	MessageId    int        `pg:"message_id"`
	State        OrderState `pg:"state,notnull"`
	DeadlineDate string     `pg:"deadline_date"`
	Price        string     `pg:"price"`
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
	case SubjectInputOrderState:
		newHeaderInputOrder(update, user, order)
	case DescriptionInputOrderState:
		newDescriptionInputOrder(update, user, order, config)
	case DeadlineInputOrderState:
		newDeadlineInputOrder(update, user, order, config)
	case PriceInputOrderState:
		newPriceInputOrder(update, user, order, config)
	case FilesUploadOrderState:
		newFilesUploadOrder(update, user, order, config)
	}
}

func toExcapedString(s string) string {
	characters := "_*[]()~`>#+-=|{}.!"
	for _, character := range characters {
		s = strings.ReplaceAll(s, string(character), `\`+string(character))
	}
	return s
}

func (order *OrderData) toTelegramString() string {
	text := "[ ](" + order.FilesURL + ")\n" +
		"Дисциплина: *" + toExcapedString(order.Subject) + "* \n\n" +
		"Описание заказа: " + toExcapedString(order.Description) + "\n\n" +
		"Дедлайн: *" + toExcapedString(order.DeadlineDate) + "*\n\n" +
		"Цена: *" + toExcapedString(order.Price) + "*"

	return text
}
