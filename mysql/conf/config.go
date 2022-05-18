package conf

var DefaultClientName string = "default"

// MysqlInitConfig 配置MySQL连接参数
type MysqlInitConfig struct {
	UserName   string // 账号
	Password   string // 密码
	Host       string // 数据库地址，可以是Ip或者域名
	Port       int    // 数据库端口
	DbName     string // 数据库名
	Timeout    string // 连接超时，10秒
	ClientName string

	MaxIdleConn     int // 设置空闲连接池中连接的最大数量
	MaxOpenConns    int // 设置打开数据库连接的最大数量
	ConnMaxLifetime int // 设置了连接可复用的最大时间,单位秒
}

func (conf MysqlInitConfig) GetClientName() string {
	if conf.ClientName == "" {
		return DefaultClientName
	}
	return conf.ClientName
}
