package main

import (
	"time"

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
		(*OrderCallback)(nil),
		(*OrderFile)(nil),
	}

	for _, model := range models {
		db.Model(model).CreateTable(&orm.CreateTableOptions{})
	}
	// defer db.Close()
}

/*USER_DATA METHODS*/
func createUser(user *UserData) {
	exists, _ := db.Model(user).Where("id=?", user.Id).Exists()
	if !exists {
		user.CustomerRatingSum = 0
		user.ExecutorRatingSum = 0
		db.Model(user).Insert()
	}
}

func readUser(id int64) *UserData {
	user := new(UserData)
	db.Model(user).Where("id=?", id).Select()
	return user
}

func updateUser(user *UserData) {
	db.Model(user).Where("id=?", user.Id).Update()
}

func deleteUser(user *UserData) {
	db.Model(user).Where("id=?", user.Id).Delete()
}

/*ORDER_DATA METHODS*/
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
	db.Model(order).Where("customer_id=? AND (state=? OR state=? OR state=? OR state=? OR state=?)",
		customerId, SubjectInputOrderState, DescriptionInputOrderState, PriceInputOrderState, DeadlineInputOrderState, FilesUploadOrderState).Select()
	return order
}

func readConfirmedOrder(duration time.Duration) (orders []OrderData) {
	currentTime := time.Now()
	currentTime = currentTime.Add(-duration)

	db.Model(&orders).Where("state=? AND confirmation_time<?", ConfirmedOrderState, currentTime).Select()

	//TODO достать из базы все заказы с состоянием ConfirmedOrderState
	return orders
}

func updateOrder(order *OrderData) {
	db.Model(order).WherePK().Update()
}

func deleteOrder(order *OrderData) {
	db.Model(order).WherePK().Delete()
}

/*ORDER_CALLBACK METHODS*/
func isExistsOrderCallback(orderCallback *OrderCallback) (exsists bool) {
	exsists, _ = db.Model(orderCallback).Where("responder_id=? AND order_id=?", orderCallback.ResponderId, orderCallback.Id).Exists()
	return exsists
}

func createOrderCallback(orderCallback *OrderCallback) {
	if !isExistsOrderCallback(orderCallback) {
		db.Model(orderCallback).Insert()

	}
}

func readCallbacksOrder(order *OrderData) (callbacks []OrderCallback) {
	db.Model(&callbacks).Where("order_id=? AND responder_id!=?", order.Id, order.ExecutorId).Select()

	return callbacks
}

func readOrderCallback(responderId, orderCallbackId int64) (orderCallback *OrderCallback) {
	orderCallback = new(OrderCallback)
	db.Model(orderCallback).Where("responder_id=? AND order_id=?", responderId, orderCallbackId).Select()

	return orderCallback
}

func deleteOrdercallback(orderCallback *OrderCallback) {
	db.Model(orderCallback).Where("responder_id=? AND order_id=?", orderCallback.ResponderId, orderCallback.Id).Delete()
}
