package models

import (
	"strconv"
	"time"

	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 字段
type File struct {
	db.Model
	Id             int    `gorm:"autoIncrement"`
	ObjType        string `gorm:"size:255"`
	ObjId          int
	FileCategoryId int
	Sort           int    `gorm:"size:11;default:0"`
	Name           string `gorm:"size:255;not null"`
	Size           string `gorm:"size:20;default:0"`
	Ext            string `gorm:"size:255"`
	Path           string `gorm:"size:255;not null"`
	Md5            string `gorm:"size:255;not null"`
	Status         int    `gorm:"size:1;not null;default:1"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// 插入数据并返回ID
func (model *File) InsertGetId(data map[string]interface{}) int {
	size := strconv.FormatInt(data["size"].(int64), 10)
	file := File{
		ObjType: data["obj_type"].(string),
		ObjId:   data["obj_id"].(int),
		Name:    data["name"].(string),
		Size:    size,
		Md5:     data["md5"].(string),
		Path:    data["path"].(string),
		Ext:     data["ext"].(string),
		Status:  1,
	}
	model.DB().Create(&file)

	return file.Id
}
