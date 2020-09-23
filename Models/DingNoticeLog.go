package Models

// 钉钉通知记录
type DingNoticeLog struct {
	ID     int    `json:"id"`
	IsPush int    `json:"is_push"`
	From   string `json:"from"`
	Level  string `json:"level"`
	Title  string `json:"title"`
	Msg    string `json:"msg"`
	BaseModel
}

func (l DingNoticeLog) TableName() string {
	return "dingtalk_notice_log"
}

// 写入
func (l *DingNoticeLog) Create() {
	Db.Create(l)
}
