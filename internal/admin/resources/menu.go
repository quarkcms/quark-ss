package resources

import (
	"encoding/json"

	"github.com/derekstavis/go-qs"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/admin/searches"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

type Menu struct {
	admin.Resource
}

// 初始化
func (p *Menu) Init() interface{} {

	// 标题
	p.Title = "菜单"

	// 模型
	p.Model = &models.Menu{}

	// 分页
	p.PerPage = false

	// 默认排序
	p.IndexOrder = "sort asc"

	return p
}

// 字段
func (p *Menu) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}

	// 权限列表
	permissions := (&models.Permission{}).List()

	// 菜单列表
	menus := (&models.Menu{}).OrderedList()

	return []interface{}{
		field.Hidden("id", "ID"), // 列表读取且不展示的字段

		field.Hidden("pid", "PID").OnlyOnIndex(), // 列表读取且不展示的字段

		field.Text("name", "名称").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "名称必须填写",
				},
			),

		field.Text("guard_name", "GuardName").
			SetDefault("admin").
			OnlyOnForms(),

		field.Icon("icon", "图标").OnlyOnForms(),

		field.Radio("type", "渲染组件").
			SetOptions(map[interface{}]interface{}{
				"default": "无组件",
				"engine":  "引擎组件",
			}).SetDefault("engine"),

		field.Text("path", "路由").
			SetEditable(true).
			SetHelp("前端路由或后端api"),

		field.Select("pid", "父节点").
			SetOptions(menus).
			SetDefault(0).
			OnlyOnForms(),

		field.Number("sort", "排序").
			SetEditable(true).
			SetDefault(0),

		field.Select("permission_ids", "绑定权限").
			SetMode("tags").
			SetOptions(permissions).
			OnlyOnForms(),

		field.Switch("status", "状态").
			SetTrueValue("正常").
			SetFalseValue("禁用").
			SetEditable(true).
			SetDefault(true),
	}
}

// 搜索
func (p *Menu) Searches(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
		(&searches.Input{}).Init("path", "路由"),
		(&searches.Status{}).Init(),
	}
}

// 行为
func (p *Menu) Actions(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&actions.CreateDrawer{}).Init(p.Title),
		(&actions.Delete{}).Init("批量删除"),
		(&actions.Disable{}).Init("批量禁用"),
		(&actions.Enable{}).Init("批量启用"),
		(&actions.ChangeStatus{}).Init(),
		(&actions.EditDrawer{}).Init("编辑"),
		(&actions.Delete{}).Init("删除"),
		(&actions.FormSubmit{}).Init(),
		(&actions.FormReset{}).Init(),
		(&actions.FormBack{}).Init(),
		(&actions.FormExtraBack{}).Init(),
	}
}

// 列表页面显示前回调
func (p *Menu) BeforeIndexShowing(c *fiber.Ctx, list []map[string]interface{}) []interface{} {

	data, _ := qs.Unmarshal(c.OriginalURL())

	if search, ok := data["search"].(map[string]interface{}); ok == true && search != nil {
		result := []interface{}{}

		for _, v := range list {
			result = append(result, v)
		}

		return result
	}

	// 转换成树形表格
	return utils.ListToTree(list, "id", "pid", "children", 0)
}

// 编辑页面显示前回调
func (p *Menu) BeforeEditing(c *fiber.Ctx, data map[string]interface{}) map[string]interface{} {
	id := c.Query("id")

	if id != "" {
		menus := []int{}

		(&db.Model{}).
			Model(&models.Permission{}).
			Where("menu_id = ?", id).
			Pluck("id", &menus)

		data["permission_ids"] = menus
	}

	return data
}

// 保存数据前回调
func (p *Menu) BeforeSaving(c *fiber.Ctx, submitData map[string]interface{}) interface{} {

	// 暂时清理permission_ids
	delete(submitData, "permission_ids")

	return submitData
}

// 保存后回调
func (p *Menu) AfterSaved(c *fiber.Ctx, model *gorm.DB) interface{} {

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)
	id := 0

	if p.IsCreating(c) {
		last := map[string]interface{}{}

		// 获取最后一条记录的id
		model.Order("id desc").First(&last)
		id = last["id"].(int)
	} else {
		id = int(data["id"].(float64))

		(&db.Model{}).
			Model(&models.Permission{}).
			Where("menu_id = ?", id).
			Update("menu_id", 0)
	}

	if data["permission_ids"] != nil {
		return (&db.Model{}).
			Model(&models.Permission{}).
			Where("id In ?", data["permission_ids"]).
			Update("menu_id", id)
	} else {
		return model
	}
}
