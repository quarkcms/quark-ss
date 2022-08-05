package db

import (
	"github.com/quarkcms/quark-go/pkg/framework/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 结构体
type Model struct{}

var conn *gorm.DB

func init() {
	defaultConn := config.Get("database.default").(string)

	switch defaultConn {
	case "sqlite":

		database := config.Get("database.connections.sqlite.database").(string)
		conn, _ = gorm.Open(sqlite.Open(database), &gorm.Config{})
	case "mysql":

		username := config.Get("database.connections.mysql.username").(string)
		password := config.Get("database.connections.mysql.password").(string)
		host := config.Get("database.connections.mysql.host").(string)
		port := config.Get("database.connections.mysql.port").(string)
		database := config.Get("database.connections.mysql.database").(string)
		charset := config.Get("database.connections.mysql.charset").(string)

		if username != "" && host != "" && port != "" && database != "" && charset != "" {
			dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset + "&parseTime=True&loc=Local"
			conn, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		}
	}
}

// 实例
func (p *Model) DB() *gorm.DB {

	return conn
}

// 模型
func (p *Model) Model(model interface{}) *gorm.DB {

	return conn.Model(model)
}
