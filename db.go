package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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

func createUser(user UserData) {

}

func readUser(id int64) {

}

func updateUser(user UserData) {

}

func deleteUser(user UserData) {

}

func createOrder(order OrderData) {

}

func readOrder(id int64) {

}

func updateOrder(order OrderData) {

}

func deleteOrder(order OrderData) {

}
