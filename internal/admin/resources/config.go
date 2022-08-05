package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/admin/searches"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
)

type Config struct {
	admin.Resource
}

// 初始化
func (p *Config) Init() interface{} {

	// 标题
	p.Title = "配置"

	// 模型
	p.Model = &models.Config{}

	// 分页
	p.PerPage = 10

	return p
}

// 字段
func (p *Config) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}

	return []interface{}{
		field.ID("id", "ID"),

		field.Text("title", "标题").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "标题必须填写",
				},
			),

		field.Text("name", "名称").
			SetEditable(true).
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "名称必须填写",
				},
			).
			SetCreationRules(
				[]string{
					"unique:configs,name",
				},
				map[string]string{
					"unique": "名称已存在",
				},
			).
			SetUpdateRules(
				[]string{
					"unique:configs,name,{id}",
				},
				map[string]string{
					"unique": "名称已存在",
				},
			),

		field.Radio("type", "表单类型").
			SetOptions(map[interface{}]interface{}{
				"text":     "输入框",
				"textarea": "文本域",
				"picture":  "图片",
				"file":     "文件",
				"switch":   "开关",
			}).
			SetDefault("text").
			OnlyOnForms(),

		field.Text("sort", "排序").
			SetEditable(true).
			SetDefault(0).
			SetHelp("值越小越靠前").
			OnlyOnForms(),

		field.Text("group_name", "分组名称").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "分组名称必须填写",
				},
			).OnlyOnForms(),

		field.Text("remark", "备注").
			OnlyOnForms(),

		field.Switch("status", "状态").
			SetTrueValue("正常").
			SetFalseValue("禁用").
			SetEditable(true).
			SetDefault(true),
	}
}

// 搜索
func (p *Config) Searches(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&searches.Input{}).Init("title", "标题"),
		(&searches.Input{}).Init("name", "名称"),
		(&searches.Status{}).Init(),
	}
}

// 行为
func (p *Config) Actions(c *fiber.Ctx) []interface{} {
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
