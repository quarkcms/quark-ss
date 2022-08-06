package actions

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/shadowsocks"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"gorm.io/gorm"
)

type ChangeServerStatus struct {
	actions.Action
}

// 初始化
func (p *ChangeServerStatus) Init() *ChangeServerStatus {
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
func (p *ChangeServerStatus) GetApiParams() []string {
	return []string{
		"id",
		"status",
	}
}

// 执行行为句柄
func (p *ChangeServerStatus) Handle(c *fiber.Ctx, model *gorm.DB) error {
	status := c.Query("status")
	id := c.Query("id")

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

	serverId, err := strconv.Atoi(id)

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	if fieldStatus == 1 {
		shadowsocks.Start(serverId)
	} else {
		shadowsocks.Stop(serverId)
	}

	result = model.Update("status", fieldStatus).Error

	if result == nil {
		return msg.Success("操作成功！", "", "")
	} else {
		return msg.Error("操作失败，请重试！", "")
	}
}
