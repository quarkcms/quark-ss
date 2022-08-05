package models

import (
	"strings"
	"time"

	"github.com/go-basic/uuid"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

// 字段
type Admin struct {
	db.Model
	Id            int    `gorm:"autoIncrement"`
	Username      string `gorm:"size:20;index:admins_username_unique,unique;not null"`
	Nickname      string `gorm:"size:200;not null"`
	Sex           int    `gorm:"size:4;not null;default:1"`
	Email         string `gorm:"size:50;index:admins_email_unique,unique;not null"`
	Phone         string `gorm:"size:11;index:admins_phone_unique,unique;not null"`
	Password      string `gorm:"size:255;not null"`
	Avatar        string
	LastLoginIp   string `gorm:"size:255"`
	LastLoginTime time.Time
	Status        int `gorm:"size:1;not null;default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

// 通过用户名获取管理员信息
func (model *Admin) GetAdminViaUsername(username string) *Admin {
	admin := &Admin{}
	model.DB().Where("status = ?", 1).Where("username = ?", username).First(&admin)

	return admin
}

// 通过角色获取管理员权限
func (model *Admin) GetPermissionsViaRoles(id int) []Permission {
	var roleIds []int
	var permissionIds []int

	// 角色id
	(&ModelHasRole{}).DB().Where("model_id", id).Where("model_type", "admin").Pluck("id", &roleIds)

	if roleIds == nil {
		return nil
	}

	// 角色权限id
	(&RoleHasPermission{}).DB().Where("role_id in (?)", roleIds).Pluck("id", &permissionIds)

	if permissionIds == nil {
		return nil
	}

	// 角色权限列表
	var permissions []Permission
	(&Permission{}).DB().Where("id in (?)", &permissionIds).Find(&permissions)

	return permissions
}

// 获取管理员角色
func (model *Admin) GetRoles(id int) *ModelHasRole {
	modelHasRole := &ModelHasRole{}
	modelHasRole.DB().Where("model_id", id).Where("model_type", "admin").First(&modelHasRole)

	return modelHasRole
}

// 获取管理员权限菜单
func (model *Admin) GetMenus(adminId int) interface{} {

	menu := &Menu{}
	var menus []map[string]interface{}
	var menuKey int

	if adminId == 1 {
		menu.DB().Model(&Menu{}).Where("status = ?", 1).Where("guard_name", "admin").Order("sort asc").Find(&menus)
	} else {
		var menuIds []int
		permissions := model.GetPermissionsViaRoles(adminId)

		if permissions != nil {
			for key, v := range permissions {
				menuIds[key] = v.MenuId
			}
		}

		var pids1 []int

		// 三级查询列表
		menu.DB().
			Where("status = ?", 1).
			Where("id in (?)", menuIds).
			Where("pid <> ?", 0).
			Order("sort asc").
			Find(&menus)

		for key, v := range menus {
			if v["pid"] != 0 {
				pids1[key] = v["pid"].(int)
			}
			menuKey = key
		}

		var pids2 []int
		menu2 := []map[string]interface{}{}

		// 二级查询列表
		menu.DB().
			Model(&Menu{}).
			Where("status = ?", 1).
			Where("id in (?)", pids1).
			Where("pid <> ?", 0).
			Order("sort asc").
			Find(&menu2)

		for key, v := range menu2 {
			if v["pid"] != 0 {
				pids2[key] = v["pid"].(int)
			}

			menuKey = menuKey + key
			menus[menuKey] = v
		}

		menu3 := []map[string]interface{}{}

		// 一级查询列表
		menu.DB().
			Model(&Menu{}).
			Where("status = ?", 1).
			Where("id in (?)", pids2).
			Where("pid", 0).
			Order("sort asc").
			Find(&menu3)

		for key, v := range menu3 {
			menuKey = menuKey + key
			menus[menuKey] = v
		}
	}

	result := []map[string]interface{}{}
	getMenu := map[string]interface{}{}

	for _, v := range menus {

		getMenu = v

		getMenu["key"] = uuid.New()
		getMenu["locale"] = "menu" + strings.Replace(v["path"].(string), "/", ".", -1)

		if v["show"] == 1 {
			getMenu["hideInMenu"] = false
		} else {
			getMenu["hideInMenu"] = true
		}

		if v["type"] == "engine" {
			getMenu["path"] = "/index?api=" + v["path"].(string)
		}

		result = append(result, getMenu)
	}

	return utils.ListToTree(result, "id", "pid", "routes", 0)
}
