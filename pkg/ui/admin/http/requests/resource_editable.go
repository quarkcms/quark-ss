package requests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type ResourceEditable struct {
	Quark
}

// 执行行为
func (p *ResourceEditable) HandleEditable(c *fiber.Ctx) error {
	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)
	data := map[string]interface{}{}

	c.Context().QueryArgs().VisitAll(func(key, val []byte) {
		var v interface{}
		k := utils.UnsafeString(key)
		v = utils.UnsafeString(val)

		if v == "true" {
			v = 1
		}

		if v == "false" {
			v = 0
		}

		data[k] = v
	})

	return model.Where("id = ?", data["id"]).Updates(data).Error
}
