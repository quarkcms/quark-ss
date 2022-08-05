package models

import (
	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 字段
type FileCategory struct {
	db.Model
	Id          int    `gorm:"autoIncrement"`
	ObjType     string `gorm:"size:100"`
	ObjId       int
	Title       string `gorm:"size:255;not null"`
	Sort        int    `gorm:"size:11;default:0"`
	Description string `gorm:"size:255"`
}
