package requests

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type ResourceDetail struct {
	Quark
}

// 表单数据
func (p *ResourceDetail) FillData(c *fiber.Ctx) map[string]interface{} {
	result := map[string]interface{}{}
	id := c.Query("id")
	if id == "" {
		return result
	}

	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)
	model.Where("id = ?", id).First(&result)

	// 获取列表字段
	detailFields := resourceInstance.(interface {
		DetailFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).DetailFields(c, resourceInstance)

	// 给实例的Field属性赋值
	resourceInstance.(interface {
		SetField(fieldData map[string]interface{}) interface{}
	}).SetField(result)

	fields := make(map[string]interface{})
	for _, field := range detailFields.([]interface{}) {

		// 字段名
		name := reflect.
			ValueOf(field).
			Elem().
			FieldByName("Name").String()

		// 获取实例的回调函数
		// callback := field.(interface{ GetCallback() interface{} }).GetCallback()

		// if callback != nil {
		// 	getCallback := callback.(func() interface{})
		// 	fields[name] = getCallback()
		// } else {
		// 	if result[name] != nil {
		// 		fields[name] = result[name]
		// 	}
		// }

		if result[name] != nil {
			fields[name] = result[name]
		}
	}

	return fields
}
