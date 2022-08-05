package actions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"github.com/quarkcms/quark-go/pkg/ui/component/action"
	"github.com/quarkcms/quark-go/pkg/ui/component/form"
)

type CreateModal struct {
	actions.Modal
}

// 初始化
func (p *CreateModal) Init(name string) *CreateModal {
	// 初始化父结构
	p.ParentInit()

	// 类型
	p.Type = "primary"

	// 图标
	p.Icon = "plus-circle"

	// 文字
	p.Name = "创建" + name

	// 执行成功后刷新的组件
	p.Reload = "table"

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	return p
}

// 内容
func (p *CreateModal) GetBody(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	api := resourceInstance.(interface {
		CreationApi(*fiber.Ctx, interface{}) string
	}).CreationApi(c, resourceInstance)

	fields := resourceInstance.(interface {
		CreationFieldsWithinComponents(*fiber.Ctx, interface{}) interface{}
	}).CreationFieldsWithinComponents(c, resourceInstance)

	// 断言BeforeCreating方法，获取初始数据
	data := resourceInstance.(interface {
		BeforeCreating(*fiber.Ctx) map[string]interface{}
	}).BeforeCreating(c)

	return (&form.Component{}).
		Init().
		SetKey("createModalForm", false).
		SetApi(api).
		SetBody(fields).
		SetInitialValues(data).
		SetLabelCol(map[string]interface{}{
			"span": 6,
		}).
		SetWrapperCol(map[string]interface{}{
			"span": 18,
		})
}

// 弹窗行为
func (p *CreateModal) GetActions(c *fiber.Ctx, resourceInstance interface{}) []interface{} {

	return []interface{}{
		(&action.Component{}).
			Init().
			SetLabel("取消").
			SetActionType("cancel"),

		(&action.Component{}).
			Init().
			SetLabel("提交").
			SetWithLoading(true).
			SetReload("table").
			SetActionType("submit").
			SetType("primary", false).
			SetSubmitForm("createModalForm"),
	}
}
