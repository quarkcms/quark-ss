package requests

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type ResourceEdit struct {
	Quark
}

// 表单数据
func (p *ResourceEdit) FillData(c *fiber.Ctx) map[string]interface{} {
	result := map[string]interface{}{}
	id := c.Query("id")
	if id == "" {
		return result
	}

	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)
	model.Where("id = ?", id).First(&result)

	// 获取列表字段
	updateFields := resourceInstance.(interface {
		UpdateFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).UpdateFields(c, resourceInstance)

	// 给实例的Field属性赋值
	resourceInstance.(interface {
		SetField(fieldData map[string]interface{}) interface{}
	}).SetField(result)

	fields := make(map[string]interface{})
	for _, field := range updateFields.([]interface{}) {

		// 字段名
		name := reflect.
			ValueOf(field).
			Elem().
			FieldByName("Name").String()

		if result[name] != nil {
			fields[name] = result[name]
		}
	}

	return fields
}
