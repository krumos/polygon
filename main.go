package main

import (
	"encoding/json"
	"io/ioutil"
	//"slices"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db *pg.DB

func connectDB(port string, username string, pword string, dbName string) {
	db = pg.Connect(&pg.Options{
		Addr:     port,
		User:     username,
		Password: pword,
		Database: dbName,
	})

	models := []interface{}{
		(*UserData)(nil),
		(*OrderData)(nil),
	}

	for _, model := range models {
		db.Model(model).CreateTable(&orm.CreateTableOptions{})
	}
	defer db.Close()
}

// func processRequest(db *pg.DB) error {
// 	request := new(transaction)
// 	if err := db.Model(request).Order("created_at ASC").Limit(1).Select(); err == nil {
// 		fromUserData := new(user)
// 		if err := db.Model(fromUserData).Where("nickname = ?", request.From_user).Select(); err == nil {
// 			if fromUserData.Balance >= request.Amount {
// 				toUserData := new(user)
// 				if err := db.Model(toUserData).Where("nickname = ?", request.To_user).Select(); err == nil {

// 					fromUserData.Balance -= request.Amount
// 					toUserData.Balance += request.Amount
// 					db.Model(fromUserData).Set("balance = ?", fromUserData.Balance).Where("nickname = ?", fromUserData.Nickname).Update()
// 					db.Model(toUserData).Set("balance = ?", toUserData.Balance).Where("nickname = ?", toUserData.Nickname).Update()
// 				} else {
// 					return err
// 				}
// 			} else {
// 			}
// 		} else {
// 			return err
// 		}
// 	} else {
// 		return err
// 	}

// 	db.Model(request).Where("id = ?", request.TransactionID).Delete()

// 	return nil
// }

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
	ChannelChat int64
	Token string
}

func main() {
	var config Config
	d, _ := ioutil.ReadFile("config.json")
	err := json.Unmarshal(d, &config)
	connectDB(":5432", "postgres", "password", "postgres")

	var ordersList = []OrderData{}
	var usersList = []UserData{}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.CallbackQuery != nil {
			if strings.Compare(update.CallbackQuery.Data, "post") == 0 {
				msg := tgbotapi.NewMessage(config.ChannelChat, update.CallbackQuery.Message.Text)
				bot.Send(msg)
				msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(msg2)
			}
			// TODO: Сделать уведомление юзера об отказе в посте
			if strings.Compare(update.CallbackQuery.Data, "del") == 0 {
				msg2 := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(msg2)
			}
		}
		if update.Message != nil {
			var msg tgbotapi.MessageConfig

			switch update.Message.Text {
			case "/start": // TODO: Check if user in DB
				msg = tgbotapi.NewMessage(update.Message.From.ID, "Привет")
				user := UserData{id: update.Message.From.ID}
				usersList = append(usersList, user)

			case "/new_order":
				msg = tgbotapi.NewMessage(update.Message.From.ID, "Напиши заголовок")
				user := usersList[0] /* slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })] */
				order := OrderData{customerId: user.id, state: title}
				ordersList = append(ordersList, order)

			default:
				// user := usersList[slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })]
				// switch user.state {

				// case title:
				// 	msg = tgbotapi.NewMessage(update.Message.From.ID, "Напиши описание")
				// 	user := usersList[slices.IndexFunc(usersList, func(u UserData) bool { return u.id == update.Message.From.ID })]
				// 	order := OrderData{customerId: user.id, state: title}
				// 	ordersList = append(ordersList, order)

				// }
				// msg = tgbotapi.NewMessage(update.Message.From.ID, "я тебя не понимать")
			}
			bot.Send(msg)
			// msg := tgbotapi.NewMessage(moderatorChat, update.Message.Text)
			// 	btn := tgbotapi.NewInlineKeyboardMarkup(
			// 		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запостить", "post"),
			// 			tgbotapi.NewInlineKeyboardButtonData("Удалить", "del")))
			// 	msg.ReplyMarkup = btn
			// 	bot.Send(msg)
		}
	}
}
