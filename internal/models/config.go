package models

import (
	"time"

	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 字段
type Config struct {
	db.Model
	Id        int    `gorm:"autoIncrement"`
	Title     string `gorm:"size:255;not null"`
	Type      string `gorm:"size:20;not null"`
	Name      string `gorm:"size:255;not null"`
	Sort      int    `gorm:"size:11;default:0"`
	GroupName string `gorm:"size:255;not null"`
	Value     string
	Remark    string `gorm:"size:100;not null"`
	Status    int    `gorm:"size:1;not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
