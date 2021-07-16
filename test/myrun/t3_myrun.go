package main

import (
	"TestChat1/db/mysql"
	"fmt"
)

func main() {
	mysql.NewMysqlDB()
	row := mysql.DB.QueryRow("select id, chat_id from message_list where id=10")
	var id int
	var chatId uint64
	err := row.Scan(&id, &chatId)
	fmt.Println(err, id, chatId)
}
