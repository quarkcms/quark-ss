package database

import (
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
)

type Migrate struct{}

// 执行迁移
func (p *Migrate) Handle() {
	(&db.Model{}).DB().AutoMigrate(
		&models.ActionLog{},
		&models.Admin{},
		&models.Config{},
		&models.Menu{},
		&models.File{},
		&models.FileCategory{},
		&models.Picture{},
		&models.PictureCategory{},
		&models.Permission{},
		&models.Role{},
		&models.ModelHasRole{},
		&models.RoleHasPermission{},
		&models.ModelHasPermission{},
		&models.Server{},
	)
}
