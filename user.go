package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserState int64

const (
	MakingOrderUserState UserState = iota + 1
	DefaultUserState
)

type UserData struct {
	Id                  int64   `pg:"id,notnull"`
	CustomerRatingSum   float32 `pg:"customer_rating,notnull"`
	ExecutorRatingSum   float32 `pg:"executor_rating,notnull"`
	CustomerRatingCount int32
	ExecutorRatingCount int32
	State               UserState
}

func userStateMachine(update tgbotapi.Update, config *Config) {
	user := readUser(update.Message.From.ID)
	switch user.State {
	case MakingOrderUserState:
		ordersStateMachine(&user, update, config)
	case DefaultUserState:

	}
}