package actions

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"gorm.io/gorm"
)

type ChangeAccount struct {
	actions.Action
}

// 执行行为句柄
func (p *ChangeAccount) Handle(c *fiber.Ctx, model *gorm.DB) error {
	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	data["avatar"], _ = json.Marshal(data["avatar"])

	result := model.Where("id", utils.Admin(c, "id")).Updates(data).Error

	if result != nil {
		return msg.Error(result.Error(), "")
	}

	return msg.Success("操作成功！", "", "")
}
