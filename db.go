package main

import (
	"fmt"

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
	// defer db.Close()
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

func createUser(user *UserData) {
	fmt.Println(user.Id)
	exists, _ := db.Model(user).Where("id=?", user.Id).Exists()
	if !exists {
		user.CustomerRating = 5
		user.ExecutorRating = 5
		db.Model(user).Insert()
	}
}

func readUser(id int64) UserData {
	user := new(UserData)
	db.Model(user).Where("id=?", id).Select()
	return *user
}

func updateUser(user *UserData) {
	db.Model(user).Where("id=?", user.Id).Update()
}

func deleteUser(user *UserData) {
	db.Model(user).Where("id=?", user.Id).Delete()
}

func createOrder(order *OrderData) {
	exists, _ := db.Model(order).Where("id=?", order.Id).Exists()
	if !exists {
		db.Model(order).Insert()
	}
}

func readOrderById(orderId int64) OrderData {
	order := new(OrderData)
	db.Model(order).Where("id=?", orderId).Select()
	return *order
}

func readOrderByState(customerId int64) (order *OrderData) {
	order = new(OrderData)
	db.Model(order).Where("customer_id=?", customerId).Where("state=?", TitleOrderState).
		WhereOr("state=?", DescriptionOrderState).
		WhereOr("state=?", FilesOrderState).Select()
	fmt.Println(order.CustomerId, order.State)
	return order
}

func updateOrder(order *OrderData) {
	fmt.Println("Бипки")
	db.Model(order).WherePK().Update()
}

func deleteOrder(order *OrderData) {
	db.Model(order).WherePK().Delete()
}
