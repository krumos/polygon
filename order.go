package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OrderState int64

const (
	TitleOrderState OrderState = iota + 1
	DescriptionOrderState
	FilesOrderState
	ModeratedOrderState
	ApprovedOrderState
	ExecutedOrderState
)

type OrderData struct {
	Id           int64               `pg:"id,pk"`
	Title        string              `pg:"title"`
	Description  string              `pg:"description"`
	Files        []tgbotapi.Document `pg:"files"` // TODO: How to store files
	CustomerId   int64               `pg:"customer_id,notnull"`
	ExecutorId   int64               `pg:"executor_id"`
	State        OrderState          `pg:"state,notnull"`
	DeadlineDate string              `pg:"deadline_date"`
	Price        int64               `pg:"price"`
}

func ordersStateMachine(user *UserData, update tgbotapi.Update, config *Config) {
	order := readOrderByState(user.Id)
	switch order.State {
	case TitleOrderState:
		newHeaderOrderCommand(update, user, order)
	case DescriptionOrderState:
		newDescriptionrOrderCommand(update, user, order, config)
	}
}

func (order *OrderData) toTelegramString() string {
	return " *" + order.Title + "* " + "\n\n" + order.Description
}
