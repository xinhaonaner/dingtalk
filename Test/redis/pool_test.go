package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	var addr = "127.0.0.1:6379"
	var password = ""
	pool := PoolInitRedis(addr, password)
	c1 := pool.Get()
	c2 := pool.Get()
	c3 := pool.Get()
	c4 := pool.Get()
	c5 := pool.Get()
	t.Log(c1, c2, c3, c4, c5)
	time.Sleep(time.Second * 5) //redis一共有多少个连接？？
	c1.Close()
	c2.Close()
	c3.Close()
	c4.Close()
	c5.Close()
	time.Sleep(time.Second * 5) //redis一共有多少个连接？？
}

func TestHandle(t *testing.T) {
	var addr = "127.0.0.1:6379"
	var password = ""
	pool := PoolInitRedis(addr, password)
	c1 := pool.Get()
	key := "admin_notice"
	num, err := redigo.Int(c1.Do("SETNX", key, 1))
	if err != nil {
		log.Fatal("err:", err)
		return
	}
	if num > 0 {
		// 设置成功,有效期60秒
		if _, err = redigo.Bool(c1.Do("EXPIRE", key, 60)); err != nil {
			log.Fatal("err:", err)
			return
		}
	} else {
		if num, err = redigo.Int(c1.Do("INCR", key)); err != nil {
			log.Fatal("err:", err)
			return
		}
		if num > 20 {
			log.Fatal("次数超过20次")
			return
		}
	}
	fmt.Println(num)
	//switch reflect.TypeOf(reply) {
	//case "uint8":
	//	log.Println(B2S(reply))
	//}
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}

	return string(ba)
}

// redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   10, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func TestString(t *testing.T) {
	a := "你好"
	t.Log(fmt.Sprintf("测试%s", a))
}
