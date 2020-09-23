package DingTalk

import (
	"fmt"
	"github.com/CatchZeng/dingtalk/client"
	"github.com/CatchZeng/dingtalk/message"
	"github.com/Unknwon/goconfig"
	"github.com/tidwall/gjson"
	"xinhaonaner-dingtalk/Log"
	"xinhaonaner-dingtalk/Models"
	"xinhaonaner-dingtalk/Redis"
)

var (
	access_tokens map[string]string
	dingTalk      client.DingTalk
)

func init() {
	access_tokens = make(map[string]string)

	cfg, err := goconfig.LoadConfigFile("config/dingtalk.ini")
	if err != nil {
		panic(err.Error())
	}
	access_tokens["wms_access_token"], err = cfg.GetValue("", "wms_access_token")
	access_tokens["wms_secret"], err = cfg.GetValue("", "wms_secret")
	access_tokens["admin_access_token"], err = cfg.GetValue("", "admin_access_token")
	access_tokens["admin_secret"], err = cfg.GetValue("", "admin_secret")

}

// MQ 消息内容
type RabbitMQMsg struct {
	From      string
	Msg       string
	Title     string
	Level     string
	Timestamp string
	// 类型
	Type string
	At   []string
}

//var Mobile arr

// 钉钉通知
func Notice(str []byte) error {

	defer protect()

	rabbitMQMsg := RabbitMQMsg{}
	//gjson.Get(string(str),"from")

	rabbitMQMsg.From = gjson.GetBytes(str, "from").String()
	rabbitMQMsg.Msg = gjson.GetBytes(str, "msg").String()
	rabbitMQMsg.Title = gjson.GetBytes(str, "title").String()
	rabbitMQMsg.Timestamp = gjson.GetBytes(str, "timestamp").String()
	rabbitMQMsg.Level = gjson.GetBytes(str, "level").String()
	rabbitMQMsg.Type = gjson.GetBytes(str, "type").String()
	if gjson.GetBytes(str, "at").Exists() {
		at := gjson.Get((string(str)), "at").Array()
		for _, v := range at {
			rabbitMQMsg.At = append(rabbitMQMsg.At, fmt.Sprintf("%s", v))
		}
	}

	isPush := 0
	defer func() {
		dingLog := new(Models.DingNoticeLog)
		dingLog.Msg = rabbitMQMsg.Msg
		dingLog.Title = rabbitMQMsg.Title
		dingLog.From = rabbitMQMsg.From
		dingLog.Level = rabbitMQMsg.Level
		dingLog.IsPush = isPush

		dingLog.Create()
	}()

	result, err := Redis.Notice(rabbitMQMsg.From)
	if result == false || err != nil {
		return err
	}

	access_token := access_tokens[rabbitMQMsg.From+"_access_token"]
	secret := access_tokens[rabbitMQMsg.From+"_secret"]

	//fmt.Println(access_token, secret)
	dingTalk.AccessToken = access_token
	dingTalk.Secret = secret

	at := message.At{}

	if len(rabbitMQMsg.At) > 0 {
		at.AtMobiles = rabbitMQMsg.At
		at.IsAtAll = false
	} else {
		at.IsAtAll = true
	}

	switch rabbitMQMsg.Type {
	case "markdown":
		msg := message.NewMarkdownMessage().SetMarkdown(rabbitMQMsg.Title, rabbitMQMsg.Msg)
		msg.At = at
		_, err = dingTalk.Send(msg)
	default:
		msg := message.NewTextMessage().SetContent(rabbitMQMsg.Msg)
		msg.At = at
		_, err = dingTalk.Send(msg)
	}

	if err != nil {
		Log.LogStash.Errorf("钉钉通知错误 %s ", err.Error())
		fmt.Println(err)
		return err
	}
	isPush = 1

	return nil

}

// 处理错误
func protect() {
	if err := recover(); err != nil {
		Log.LogStash.Errorf("recover-error：%s", err)
	}
}
