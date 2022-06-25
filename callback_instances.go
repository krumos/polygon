package main

type CallbackDataType int64

const (
	Approve CallbackDataType = iota + 1
	Reject
	Agreement
)

type CallbackData struct {
	Type CallbackDataType
	Id   int64
}

type OrderCallback struct {
	Id          int64 `pg:"id"`
	ResponderId int64 `pg:"responder_id"`
	MessageId   int64 `pg:"message_id"`
}

// userBotKeyboard := tgbotapi.NewReplyKeyboard(
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("/start"),
// 	),
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("/new_order"),
// 	),
// )
