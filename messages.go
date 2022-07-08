package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" 
)

func StartMessage(chatId int64) (tgbotapi.MessageConfig) {
	msg := tgbotapi.NewMessage(chatId, Texts["start_command_answer"])
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Texts["new_order_command"])))
	return msg
}

func FilesUploadMessage(chatId int64) (tgbotapi.MessageConfig) {
	msg := tgbotapi.NewMessage(chatId, Texts["new_price_command_an"])
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Texts["stop"])))
	return msg
}

