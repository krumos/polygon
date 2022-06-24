package main

import (
	"encoding/json"
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getConfig() (config *Config, err error) {
	d, e := ioutil.ReadFile("config.json")
	if e != nil {
		config = nil
		return nil, e
	}
	e = json.Unmarshal(d, &config)
	return config, e
}

func connectBot(token string) (bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, err error) {
	bot, e := tgbotapi.NewBotAPI(token)
	if e != nil {
		return nil, nil, e
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates = bot.GetUpdatesChan(updateConfig)
	return bot, updates, e
}

func main() {
	getText("texts.yaml")
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	bot, updates, err := connectBot(config.Token)
	if err != nil {
		panic(err)
	}

	connectDB(":5432", "postgres", "password", "postgres")

	// userBotKeyboard := tgbotapi.NewReplyKeyboard(
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("/start"),
	// 	),
	// 	tgbotapi.NewKeyboardButtonRow(
	// 		tgbotapi.NewKeyboardButton("/new_order"),
	// 	),
	// )

	for update := range updates {
		// Действия с запросом на пост в канал
		if update.CallbackQuery != nil {
			response := CallbackData{}
			json.Unmarshal([]byte(update.CallbackQuery.Data), &response)
			switch response.Type {
			case Approve:
				approveOrderResponse(config, update, bot, &response)
				// TODO: Сделать уведомление юзера об отказе в посте
			case Reject:
				rejectOrderResponse(config, update, bot, &response)
			case Agreement:
				agreementOrderResponse(update, bot, &response)
			}
		}
		// Взаимодействие юзера с постом
		if update.Message != nil {
			switch update.Message.Text {
			case Texts["start_command"]: // TODO: Check if user in DB
				startCommand(update, bot)
			case Texts["new_order_command"]:
				newOrderCommand(update, bot)
			default:
				user := readUser(update.Message.From.ID)

				switch user.State {
				case MakingOrderUserState:
					order := readOrderByState(user.Id)
					switch order.State {
					case TitleOrderState:
						newHeaderOrderCommand(update, bot, &user, order)
					case DescriptionOrderState:
						newDescriptionrOrderCommand(update, bot, &user, order)
					}

				default:

				}

				//	bot.Send(tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text))
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

		}
	}
}
