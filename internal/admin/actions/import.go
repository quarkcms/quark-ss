package actions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/action"
	"github.com/quarkcms/quark-go/pkg/ui/component/form"
	"github.com/quarkcms/quark-go/pkg/ui/component/space"
	"github.com/quarkcms/quark-go/pkg/ui/component/tpl"
)

type Import struct {
	actions.Modal
}

// 初始化
func (p *Import) Init() *Import {
	// 初始化父结构
	p.ParentInit()

	// 文字
	p.Name = "导入数据"

	// 设置展示位置
	p.SetOnlyOnIndex(true)

	return p
}

// 内容
func (p *Import) GetBody(c *fiber.Ctx, resourceInstance interface{}) interface{} {
	api := "admin/" + c.Params("resource") + "/import"
	getTpl := (&tpl.Component{}).
		Init().
		SetBody("模板文件: <a href='/api/admin/" + c.Params("resource") + "/import/template?token=" + utils.GetAdminToken(c) + "' target='_blank'>下载模板</a>").
		SetStyle(map[string]interface{}{
			"marginLeft": "50px",
		})

	fields := []interface{}{
		(&space.Component{}).
			Init().
			SetBody(getTpl).
			SetDirection("vertical").
			SetSize("middle").
			SetStyle(map[string]interface{}{
				"marginBottom": "20px",
			}),
		(&admin.Field{}).
			File("fileId", "导入文件").
			SetLimitNum(1).
			SetLimitType([]string{
				"application/vnd.ms-excel",
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			}).
			SetHelp("请上传xls格式的文件"),
	}

	return (&form.Component{}).
		Init().
		SetKey("importModalForm", false).
		SetApi(api).
		SetBody(fields).
		SetLabelCol(map[string]interface{}{
			"span": 6,
		}).
		SetWrapperCol(map[string]interface{}{
			"span": 18,
		})
}

// 弹窗行为
func (p *Import) GetActions(c *fiber.Ctx, resourceInstance interface{}) []interface{} {

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
			SetSubmitForm("importModalForm"),
	}
}
