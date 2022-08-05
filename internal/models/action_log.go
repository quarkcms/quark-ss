package models

import (
	"time"

	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 字段
type ActionLog struct {
	db.Model
	Id        int `gorm:"autoIncrement"`
	ObjectId  int
	Username  string `gorm:"<-:false"`
	Url       string `gorm:"size:500;not null"`
	Remark    string `gorm:"size:255;not null"`
	Ip        string `gorm:"size:100;not null"`
	Type      string `gorm:"size:100;not null"`
	Status    int    `gorm:"size:1;not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 插入数据
func (model *ActionLog) Insert(data map[string]interface{}) {

	log := ActionLog{
		ObjectId: data["obj_id"].(int),
		Url:      data["url"].(string),
		Ip:       data["ip"].(string),
		Type:     data["type"].(string),
		Status:   1,
	}

	model.DB().Create(&log)
}
