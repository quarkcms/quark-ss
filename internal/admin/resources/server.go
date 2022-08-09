package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/admin/searches"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
)

type Server struct {
	admin.Resource
}

// 初始化
func (p *Server) Init() interface{} {

	// 标题
	p.Title = "服务"

	// 模型
	p.Model = &models.Server{}

	// 分页
	p.PerPage = 10

	return p
}

// 字段
func (p *Server) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}

	return []interface{}{
		field.ID("id", "ID"),

		field.Text("name", "名称").
			SetRules(
				[]string{
					"required",
					"max:20",
				},
				map[string]string{
					"required": "名称必须填写",
					"max":      "用户名不能超过20个字符",
				},
			),

		field.Select("encrypt_type", "加密方式").
			SetOptions([]map[string]interface{}{
				{"label": "不加密", "value": "DUMMY"},
				{"label": "chacha20-ietf-poly1305", "value": "CHACHA20-IETF-POLY1305"},
				{"label": "aes-128-gcm", "value": "AES-128-GCM"},
				{"label": "aes-256-gcm", "value": "AES-256-GCM"},
			}).
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "请选择加密方式",
				},
			).SetHelp("建议选择aes-256-gcm加密方式"),

		field.Text("password", "密码").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "密码必须填写",
				},
			),

		field.Text("port", "端口").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "端口必须填写",
				},
			),

		field.Text("plugin", "插件"),

		field.Text("plugin_opts", "插件参数"),

		field.Switch("status", "状态").
			SetTrueValue("启用").
			SetFalseValue("禁用").
			SetDefault(false).
			OnlyOnIndex(),
	}
}

// 搜索
func (p *Server) Searches(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
		(&searches.Status{}).Init(),
	}
}

// 行为
func (p *Server) Actions(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&actions.Import{}).Init(),
		(&actions.CreateLink{}).Init(p.Title),
		(&actions.Delete{}).Init("批量删除"),
		(&actions.ChangeServerStatus{}).Init(),
		(&actions.EditLink{}).Init("编辑"),
		(&actions.Delete{}).Init("删除"),
		(&actions.FormSubmit{}).Init(),
		(&actions.FormReset{}).Init(),
		(&actions.FormBack{}).Init(),
		(&actions.FormExtraBack{}).Init(),
	}
}
