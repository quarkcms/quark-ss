package resources

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/admin/searches"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
	"gorm.io/gorm"
)

type Role struct {
	admin.Resource
}

// 初始化
func (p *Role) Init() interface{} {

	// 标题
	p.Title = "角色"

	// 模型
	p.Model = &models.Role{}

	// 分页
	p.PerPage = 10

	return p
}

// 字段
func (p *Role) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}
	treeData := (&models.Menu{}).Tree()

	return []interface{}{
		field.ID("id", "ID"),

		field.Text("name", "名称").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "名称必须填写",
				},
			),

		field.Text("guard_name", "GuardName").SetDefault("admin"),
		field.Tree("menu_ids", "权限").SetData(treeData).OnlyOnForms(),
		field.Datetime("created_at", "创建时间", func() interface{} {
			if p.Field["created_at"] == nil {
				return p.Field["created_at"]
			}

			return p.Field["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		}).OnlyOnIndex(),
		field.Datetime("updated_at", "更新时间", func() interface{} {
			if p.Field["updated_at"] == nil {
				return p.Field["updated_at"]
			}

			return p.Field["updated_at"].(time.Time).Format("2006-01-02 15:04:05")
		}).OnlyOnIndex(),
	}
}

// 搜索
func (p *Role) Searches(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
	}
}

// 行为
func (p *Role) Actions(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&actions.CreateLink{}).Init(p.Title),
		(&actions.Delete{}).Init("批量删除"),
		(&actions.EditLink{}).Init("编辑"),
		(&actions.Delete{}).Init("删除"),
		(&actions.FormSubmit{}).Init(),
		(&actions.FormReset{}).Init(),
		(&actions.FormBack{}).Init(),
		(&actions.FormExtraBack{}).Init(),
	}
}

// 编辑页面显示前回调
func (p *Role) BeforeEditing(c *fiber.Ctx, data map[string]interface{}) map[string]interface{} {
	id := c.Query("id")
	menus := []map[string]interface{}{}

	(&db.Model{}).Model(&models.Menu{}).Find(&menus)

	checkedMenus := []int{}
	for _, v := range menus {
		var permissionIds []int
		(&db.Model{}).
			Model(&models.Permission{}).
			Where("menu_id", v["id"]).
			Pluck("id", &permissionIds)

		if len(permissionIds) > 0 {
			roleHasPermission := map[string]interface{}{}
			(&db.Model{}).
				Model(&models.RoleHasPermission{}).
				Where("permission_id IN ?", permissionIds).
				Where("role_id", id).
				First(&roleHasPermission)

			if len(roleHasPermission) > 0 {
				checkedMenus = append(checkedMenus, v["id"].(int))
			}
		}
	}

	data["menu_ids"] = checkedMenus

	return data
}

// 保存数据前回调
func (p *Role) BeforeSaving(c *fiber.Ctx, submitData map[string]interface{}) interface{} {

	// 根据菜单id获取所有权限
	var permissionIds []int
	(&db.Model{}).
		Model(&models.Permission{}).
		Where("menu_id IN ?", submitData["menu_ids"]).
		Pluck("id", &permissionIds)

	if len(permissionIds) == 0 {
		return errors.New("获取的权限为空，请在菜单管理中绑定权限")
	}

	delete(submitData, "menu_ids")

	return submitData
}

// 保存后回调
func (p *Role) AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{} {

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	// 根据菜单id获取所有权限
	var permissionIds []int
	(&db.Model{}).
		Model(&models.Permission{}).
		Where("menu_id IN ?", data["menu_ids"]).
		Pluck("id", &permissionIds)

	if len(permissionIds) == 0 {
		return errors.New("获取的权限为空，请先在菜单管理中绑定权限")
	}

	var result interface{}

	if p.IsCreating(c) {
		lastRole := map[string]interface{}{}
		model.Order("id desc").First(&lastRole) // hack

		// 同步权限
		result = p.syncPermissions(lastRole["id"].(int), permissionIds)
	} else {

		// 同步权限
		id := data["id"].(float64)
		result = p.syncPermissions(int(id), permissionIds)
	}

	return result
}

// 保存后回调
func (p *Role) syncPermissions(roleId int, permissionIds []int) *gorm.DB {
	permissionIds = p.arrayFilter(permissionIds)

	// 先清空此角色的权限
	(&db.Model{}).Model(&models.RoleHasPermission{}).Where("role_id", roleId).Delete("")

	data := []map[string]interface{}{}
	for _, v := range permissionIds {
		permission := map[string]interface{}{
			"role_id":       roleId,
			"permission_id": v,
		}
		data = append(data, permission)
	}

	return (&db.Model{}).Model(&models.RoleHasPermission{}).Create(data)
}

// 数组去重
func (p *Role) arrayFilter(list []int) []int {
	var x []int = []int{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}

	return x
}
