package Models

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
	"time"
	"xinhaonaner-dingtalk/Log"
)

type BaseModel struct {
	gorm.Model
}

var Db *gorm.DB

func init() {
	cfg, err := goconfig.LoadConfigFile("config/database.ini")
	if err != nil {
		Log.LogStash.Errorf("获取数据库配置错误：%s", err)
		return
	}
	//args := cfg.MustValue("mysql", "username") + ":" + cfg.MustValue("mysql", "password") + "@/" + cfg.MustValue("mysql", "dbname") + "?charset=" + cfg.MustValue("mysql", "charset") + "&parseTime=True&loc=" + cfg.MustValue("mysql", "host")
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	args := cfg.MustValue("mysql", "username") + ":" + cfg.MustValue("mysql", "password") + "@tcp(" + cfg.MustValue("mysql", "host") + ":" + cfg.MustValue("mysql", "port") + ")/" + cfg.MustValue("mysql", "dbname") + "?charset=" + cfg.MustValue("mysql", "charset") + "&parseTime=True&loc=Local"
	Log.LogStash.Infof("args：%s", args)
	// user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	//db, err := gorm.Open("mysql", args)
	Db, err = gorm.Open("mysql", args)
	if err != nil {
		Log.LogStash.Errorf("数据库连接失败错误：%s", err)
		defer func() {
			_ = fmt.Errorf("数据库连接,%s", err)
		}()
		return
	}

	//defer func() {
	//	err := Db.Close()
	//	Log.LogStash.Errorf("数据库异常关闭:%s", err)
	//	fmt.Println("数据库异常关闭")
	//}()

	sqlDB := Db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
