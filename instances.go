package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type OrderState int64

const (
	TitleOrderState OrderState = iota + 1
	DescriptionOrderState
	FilesOrderState
	ModeratedOrderState
	ApprovedOrderState
	ExecutedOrderState
)

type UserState int64

const (
	MakingOrderUserState UserState = iota + 1
	DefaultUserState
)

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

type UserData struct {
	Id             int64   `pg:"id,notnull"`
	CustomerRating float32 `pg:"customer_rating,notnull"`
	ExecutorRating float32 `pg:"executor_rating,notnull"`
	State          UserState
}

type OrderData struct {
	Id          int64               `pg:"id,pk"`
	Title       string              `pg:"title"`
	Description string              `pg:"description"`
	Files       []tgbotapi.Document `pg:"files"` // TODO: How to store files
	CustomerId  int64               `pg:"customer_id,notnull"`
	ExecutorId  int64               `pg:"executor_id"`
	State       OrderState          `pg:"state,notnull"`
}

func (order *OrderData) toString() string {
	return " *" + order.Title + "* " + "\n\n" + order.Description
}

type Config struct {
	ModeratorChat int64
	ChannelChat   int64
	Token         string
}
