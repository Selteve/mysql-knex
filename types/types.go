package types

// 定义数据库结构体
type DB struct {
	Host string `json:"host"` // 数据库地址
	Port int `json:"port"` // 数据库端口
	User string `json:"user"` // 数据库用户名
	DBName string `json:"name"` // 数据库名称
	Password string `json:"password"` // 数据库密码
	Database string `json:"database"` // 数据库名称	
}

type DBName struct {
	Table string // 数据库名称
}