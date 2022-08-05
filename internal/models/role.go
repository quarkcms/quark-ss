package models

import (
	"time"

	"github.com/quarkcms/quark-go/pkg/framework/db"
)

// 角色
type Role struct {
	db.Model
	Id        int    `gorm:"autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	GuardName string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 模型角色关联表
type ModelHasRole struct {
	db.Model
	RoleId    int    `gorm:"index:model_has_roles_model_id_model_type_index"`
	ModelType string `gorm:"size:255;not null"`
	ModelId   int    `gorm:"index:model_has_roles_model_id_model_type_index"`
}

// 角色权限关联表
type RoleHasPermission struct {
	db.Model
	PermissionId int
	RoleId       int `gorm:"index:role_has_permissions_role_id_foreign"`
}

// 模型权限关联表
type ModelHasPermission struct {
	db.Model
	PermissionId int    `gorm:"index:model_has_permissions_model_id_model_type_index"`
	ModelType    string `gorm:"size:255;not null"`
	ModelId      int    `gorm:"index:model_has_permissions_model_id_model_type_index"`
}

// 获取角色列表
func (model *Role) List() map[interface{}]interface{} {
	roles := []Role{}
	results := map[interface{}]interface{}{}

	model.DB().Find(&roles)
	for _, v := range roles {
		results[v.Id] = v.Name
	}

	return results
}
