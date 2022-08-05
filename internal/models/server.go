package models

import (
	"time"

	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 字段
type Server struct {
	db.Model
	Id          int    `gorm:"autoIncrement"`
	Name        string `gorm:"size:200;not null"`
	EncryptType string `gorm:"size:200;not null"`
	Password    string `gorm:"size:255;not null"`
	Port        string `gorm:"size:200;not null"`
	Plugin      string `gorm:"size:200;not null"`
	PluginOpts  string `gorm:"size:200;not null"`
	Key         string `gorm:"size:200;not null"`
	Status      int    `gorm:"size:1;not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// 返回信息
func (model *Server) Info(id int) Server {
	server := Server{}

	model.DB().Where("id", id).Find(&server)

	return server
}
