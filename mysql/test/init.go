package test

import (
	"fmt"
	"github.com/yueja/cf-gorm/mysql/conf"
	"github.com/yueja/cf-gorm/mysql/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func init() {
	config := conf.MysqlInitConfig{
		UserName: `dbuser`,
		Password: "testpassw0rd",
		Host:     "dbi.mshare.cn",
		Port:     4302,
		DbName:   "yja",
		Timeout:  "10s",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		config.UserName, config.Password, config.Host, config.Port, config.DbName, config.Timeout)
	var registryClient = db.RegistryClientConfig{
		MysqlConfig: mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		},
		GormConfig: gorm.Config{
			SkipDefaultTransaction: false, //跳过默认事务
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 复数形式 User的表名应该是users
				TablePrefix:   "t_", // 表名前缀 User的表名应该是t_users
			},
			DisableForeignKeyConstraintWhenMigrating: true, //设置成为逻辑外键(在物理数据库上没有外键，仅体现在代码上)
		},
		MaxIdleConn:     10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
	}
	if err := db.RegistryClient("", registryClient); err != nil {
		panic(err)
	}
}
