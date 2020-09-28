package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"xinhaonaner-dingtalk/RabbitMQ"
)

func main() {
	fmt.Printf("钉钉通知-进程：%d \n", os.Getpid())

	//go RabbitMQ.Produce()
	//select {
	//
	//}

	RabbitMQ.Consume("dingtalk", "notice")

}
