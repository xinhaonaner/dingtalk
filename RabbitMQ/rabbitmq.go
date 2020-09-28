package RabbitMQ

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/streadway/amqp"
	"xinhaonaner-dingtalk/Log"
)

// MQURL 格式 amqp://账号：密码@rabbitmq服务器地址：端口号/vhost
//var MQURL = beego.AppConfig.String("rabbitmq_url")
var MQURL string

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key string
	// 连接信息
	Mqurl string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	cfg, err := goconfig.LoadConfigFile("config/rabbitmq.ini")
	if err != nil {
		panic(err.Error())
	}
	MQURL, err = cfg.GetValue("", "mq_url")

	rabbitmq := &RabbitMQ{
		//conn:      nil,
		//channel:   nil,
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}

	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)  //通过amqp.Dial()方法去链接rabbitmq服务端
	rabbitmq.failOnErr(err, "创建连接错误!")    //调用我们自定义的failOnErr()方法去处理异常错误信息
	rabbitmq.channel, err = rabbitmq.conn.Channel() //链接上rabbitmq之后通过rabbitmq.conn.Channel()去设置channel信道
	rabbitmq.failOnErr(err, "获取channel失败!")

	return rabbitmq
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, msg string) {
	if err != nil {
		Log.LogStash.Errorf("%s:%s", msg, err)
		panic(fmt.Errorf("%s:%s", msg, err))
	}
}

// Destory 断开channel和connection
func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

//topic主题模式step1:创建RabbitMQ实例
func NewRabbitMQTopic(exchange string, routingkey string) *RabbitMQ {
	//创建RabbitMQ实例
	return NewRabbitMQ("", exchange, routingkey)
}
