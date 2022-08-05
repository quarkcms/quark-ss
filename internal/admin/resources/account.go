package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
)

type Account struct {
	admin.Resource
}

// 初始化
func (p *Account) Init() interface{} {

	// 标题
	p.Title = "个人设置"

	// 模型
	p.Model = &models.Admin{}

	return p
}

// 表单接口
func (p *Account) FormApi(c *fiber.Ctx) string {

	return "admin/account/action/change-account"
}

// 字段
func (p *Account) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}

	return []interface{}{

		field.Image("avatar", "头像").OnlyOnForms(),

		field.Text("nickname", "昵称").
			SetEditable(true).
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "昵称必须填写",
				},
			),

		field.Text("email", "邮箱").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "邮箱必须填写",
				},
			).
			SetCreationRules(
				[]string{
					"unique:admins,email",
				},
				map[string]string{
					"unique": "邮箱已存在",
				},
			).
			SetUpdateRules(
				[]string{
					"unique:admins,email,{id}",
				},
				map[string]string{
					"unique": "邮箱已存在",
				},
			),

		field.Text("phone", "手机号").
			SetRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "手机号必须填写",
				},
			).
			SetCreationRules(
				[]string{
					"unique:admins,phone",
				},
				map[string]string{
					"unique": "手机号已存在",
				},
			).
			SetUpdateRules(
				[]string{
					"unique:admins,phone,{id}",
				},
				map[string]string{
					"unique": "手机号已存在",
				},
			),

		field.Radio("sex", "性别").
			SetOptions(map[interface{}]interface{}{
				1: "男",
				2: "女",
			}).SetDefault(1),

		field.Password("password", "密码").
			SetCreationRules(
				[]string{
					"required",
				},
				map[string]string{
					"required": "密码必须填写",
				},
			).OnlyOnForms(),
	}
}

// 行为
func (p *Account) Actions(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&actions.ChangeAccount{}),
		(&actions.FormSubmit{}).Init(),
		(&actions.FormReset{}).Init(),
		(&actions.FormBack{}).Init(),
		(&actions.FormExtraBack{}).Init(),
	}
}

// 创建页面显示前回调
func (p *Account) BeforeCreating(c *fiber.Ctx) map[string]interface{} {
	data := map[string]interface{}{}

	(&db.Model{}).
		Model(p.Model).
		Where("id = ?", utils.Admin(c, "id")).
		First(&data)

	delete(data, "password")

	return data
}
