package actions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"gorm.io/gorm"
)

type ChangeStatus struct {
	actions.Action
}

// 初始化
func (p *ChangeStatus) Init() *ChangeStatus {
	// 初始化父结构
	p.ParentInit()

	// 行为名称，当行为在表格行展示时，支持js表达式
	p.Name = "<%= (status==1 ? '禁用' : '启用') %>"

	// 设置按钮类型,primary | ghost | dashed | link | text | default
	p.Type = "link"

	// 设置按钮大小,large | middle | small | default
	p.Size = "small"

	//  执行成功后刷新的组件
	p.Reload = "table"

	// 设置展示位置
	p.SetOnlyOnIndexTableRow(true)

	// 当行为在表格行展示时，支持js表达式
	p.WithConfirm("确定要<%= (status==1 ? '禁用' : '启用') %>数据吗？", "", "pop")

	return p
}

/**
 * 行为接口接收的参数，当行为在表格行展示的时候，可以配置当前行的任意字段
 *
 * @return array
 */
func (p *ChangeStatus) GetApiParams() []string {
	return []string{
		"id",
		"status",
	}
}

// 执行行为句柄
func (p *ChangeStatus) Handle(c *fiber.Ctx, model *gorm.DB) error {
	status := c.Query("status")

	if status == "" {
		return msg.Error("参数错误！", "")
	}

	var fieldStatus int
	var result interface{}

	if status == "1" {
		fieldStatus = 0
	} else {
		fieldStatus = 1
	}

	result = model.Update("status", fieldStatus).Error

	if result == nil {
		return msg.Success("操作成功！", "", "")
	} else {
		return msg.Error("操作失败，请重试！", "")
	}
}
