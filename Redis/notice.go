package Redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"xinhaonaner-dingtalk/Log"
)

// 每分钟最大通知次数
var NoticeMaxNum = 5

// 通知计数
func Notice(key string) (bool, error) {
	c := Redis.Pool.Get()
	num, err := redigo.Int(c.Do("SETNX", key, 1))
	if err != nil {
		Log.LogStash.Errorf("err:%s", err)
		return false, err
	}

	defer func() {
		_ = c.Close()
	}()

	if num > 0 {
		// 设置成功,有效期60秒
		if _, err = redigo.Bool(c.Do("EXPIRE", key, 60)); err != nil {
			Log.LogStash.Errorf("redis=%s有效期设置失败：%s", key, err)
			return true, err
		}
	} else {
		if num, err = redigo.Int(c.Do("INCR", key)); err != nil {
			Log.LogStash.Errorf("err:%s", err)
			return false, err
		}
		if num > NoticeMaxNum {
			str := fmt.Sprintf("钉钉通知次数超过%s", NoticeMaxNum)
			Log.LogStash.Error(str)
			return false, nil
		}
	}

	return true, nil
}
