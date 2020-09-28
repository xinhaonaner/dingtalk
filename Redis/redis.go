package Redis

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

var Redis struct {
	Addr     string
	Password string
	Pool     *redigo.Pool
}

func init() {
	NewRedis()
}

// 构造
func NewRedis() {

	cfg, err := goconfig.LoadConfigFile("config/database.ini")
	if err != nil {
		panic(err.Error())
	}

	//redis := new(Redis)
	//Redis := Redis{
	//	Addr:     cfg.MustValue("redis", "host"),
	//	Password: cfg.MustValue("redis", "password"),
	//	Pool:     nil,
	//}
	Redis.Addr = cfg.MustValue("redis", "host") + ":" + cfg.MustValue("redis", "port")
	Redis.Password = cfg.MustValue("redis", "password")
	Redis.Pool = poolInitRedis()

}

// 连接池
func poolInitRedis() *redigo.Pool {
	return &redigo.Pool{
		//空闲数
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		//最大数
		MaxActive: 100,
		Wait: true,
		Dial: func() (redigo.Conn, error) {
			fmt.Println(Redis.Addr)
			c, err := redigo.Dial("tcp", Redis.Addr)
			if err != nil {
				return nil, err
			}
			if Redis.Password != "" {
				if _, err := c.Do("AUTH", Redis.Password); err != nil {
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
