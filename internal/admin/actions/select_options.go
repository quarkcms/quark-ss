package actions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/ui/admin/actions"
	"gorm.io/gorm"
)

type SelectOptions struct {
	actions.Action
}

// 执行行为句柄
func (p *SelectOptions) Handle(c *fiber.Ctx, model *gorm.DB) error {

	resource := c.Params("resource")
	search := c.Query("search")

	lists := []map[string]interface{}{}
	results := []map[string]interface{}{}

	switch resource {
	case "Some Field":
		model.Where("Some Field = ?", search).Find(&lists)
		for _, v := range lists {
			item := map[string]interface{}{
				"label": v["name"],
				"value": v["id"],
			}

			results = append(results, item)
		}
	}

	return msg.Success("操作成功！", "", results)
}
