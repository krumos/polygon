package main

type CallbackDataType int64

const (
	Approve CallbackDataType = iota + 1
	Reject
	Agreement
	Confirm
)

type CallbackData struct {
	Type       CallbackDataType
	Id         int64
	ExecutorId int64
}

type OrderCallback struct {
	Id          int64 `pg:"order_id"`
	ResponderId int64 `pg:"responder_id"`
	MessageId   int   `pg:"message_id"`
}

// userBotKeyboard := tgbotapi.NewReplyKeyboard(
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("/start"),
// 	),
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("/new_order"),
// 	),
// )
