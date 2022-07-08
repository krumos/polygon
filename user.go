package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserState int64

const (
	MakingOrderUserState UserState = iota + 1
	DefaultUserState
)

type UserData struct {
	Id                  int64 `pg:"id,notnull"`
	CustomerRatingSum   int32 `pg:"customer_rating,use_zero"`
	ExecutorRatingSum   int32 `pg:"executor_rating,use_zero"`
	CustomerRatingCount int32 `pg:"customer_rating_count,use_zero"`
	ExecutorRatingCount int32 `pg:"executor_rating_count,use_zero"`
	State               UserState
}

func userStateMachine(update tgbotapi.Update, config *Config) {
	user := readUser(update.Message.From.ID)
	switch user.State {
	case MakingOrderUserState:
		ordersStateMachine(user, update, config)
	case DefaultUserState:
		//TODO:
	}
}
