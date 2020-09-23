package main

import (
	_ "github.com/go-sql-driver/mysql"
	"xinhaonaner-dingtalk/RabbitMQ"
)

func main() {

	//go RabbitMQ.Produce()
	//select {
	//
	//}

	RabbitMQ.Consume("dingtalk", "notice")

}
