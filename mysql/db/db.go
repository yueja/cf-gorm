package db

import (
	"errors"
	"github.com/yueja/cf-gorm/mysql/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// dbClientMap 客户端注册表
var dbClientMap = make(map[string]*gorm.DB)

type RegistryClientConfig struct {
	MysqlConfig     mysql.Config
	GormConfig      gorm.Config
	MaxIdleConn     int // 设置空闲连接池中连接的最大数量
	MaxOpenConns    int // 设置打开数据库连接的最大数量
	ConnMaxLifetime int // 设置了连接可复用的最大时间,单位秒
}

// RegistryClient 注册一个新的客户端
// name：客户端名称，config：客户端注册选项
// err：异常
func RegistryClient(name string, config RegistryClientConfig) (err error) {
	if name == "" {
		name = conf.DefaultClientName
	}
	if _, ok := dbClientMap[name]; ok {
		err = errors.New("client already exists")
		return
	}
	db, err := gorm.Open(mysql.New(config.MysqlConfig), &config.GormConfig)
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	dbClientMap[name] = db
	return
}

// GetClient 获取数据库客户端句柄
func GetClient(client []string) *gorm.DB {
	if len(client) != 0 && client[0] != "" {
		return dbClientMap[client[0]]
	}
	return dbClientMap[conf.DefaultClientName]
}
