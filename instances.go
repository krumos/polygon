package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State int64

const (
	title State = iota
	description
	files
	posted
	executed
)

type UserData struct {
	id             int64   `pg:"id"`
	customerRating float32 `pg:"customer_rating"`
	executorRating float32 `pg:"executor_rating"`
}

type OrderData struct {
	id          int64                `pg:"id"`
	title       string               `pg:"title"`
	description string               `pg:"description"`
	files       []tgbotapi.FileBytes `pg:"files"` // TODO: How to store files
	customerId  int64                `pg:"customer_id"`
	executorId  int64                `pg:"executor_id"`
	state       State                `pg:"state"`
}

type Config struct {
	ModeratorChat int64
	ChannelChat   int64
	Token         string
}
