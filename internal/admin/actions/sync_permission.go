package actions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

type SyncPermission struct {
	actions.Action
}

// 初始化
func (p *SyncPermission) Init() *SyncPermission {
	// 初始化父结构
	p.ParentInit()

	// 行为名称
	p.Name = "同步权限"

	// 执行成功后刷新的组件
	p.Reload = "table"

	// 是否具有loading，当action 的作用类型为ajax,submit时有效
	p.WithLoading = true

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	// 行为类型
	p.ActionType = "ajax"

	return p
}

// 执行行为句柄
func (p *SyncPermission) Handle(c *fiber.Ctx, model *gorm.DB) error {
	// 获取当前权限
	permissions := utils.GetPermissions()
	data := []models.Permission{}

	var names []string
	(&db.Model{}).Model(&models.Permission{}).Pluck("name", &names)

	for _, v := range permissions {
		has := false
		for _, nv := range names {
			if nv == v {
				has = true
			}
		}

		if has == false {
			permission := models.Permission{
				MenuId:    0,
				Name:      v,
				GuardName: "admin",
			}
			data = append(data, permission)
		}
	}

	if len(data) == 0 {
		return msg.Success("暂无新增权限！", "", "")
	}

	result := model.Create(data)

	if result.Error != nil {
		return msg.Error("操作失败，请重试！", "")
	}

	result1 := (&db.Model{}).Model(&models.Permission{}).Where("name NOT IN ?", permissions).Delete("")

	if result1.Error != nil {
		return msg.Error("操作失败，请重试！", "")
	}

	return msg.Success("操作成功！", "", "")
}
