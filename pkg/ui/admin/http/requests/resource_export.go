package requests

import (
	"reflect"

	"github.com/derekstavis/go-qs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceExport struct {
	Quark
}

// 列表查询
func (p *ResourceExport) IndexQuery(c *fiber.Ctx) interface{} {
	var lists []map[string]interface{}
	resourceInstance := p.Resource(c)
	model := p.NewModel(resourceInstance)

	// 搜索项
	searches := resourceInstance.(interface {
		Searches(c *fiber.Ctx) []interface{}
	}).Searches(c)

	// 过滤项，预留
	filters := resourceInstance.(interface {
		Filters(c *fiber.Ctx) []interface{}
	}).Filters(c)

	query := resourceInstance.(interface {
		BuildExportQuery(c *fiber.Ctx, resourceInstance interface{}, query *gorm.DB, search []interface{}, filters []interface{}, columnFilters map[string]interface{}, orderings map[string]interface{}) *gorm.DB
	}).BuildExportQuery(c, resourceInstance, model, searches, filters, p.columnFilters(c), p.orderings(c))

	query.Find(&lists)

	// 返回解析列表
	return p.performsList(c, resourceInstance, lists)
}

/**
 * Get the column filters for the request.
 *
 * @return array
 */
func (p *ResourceExport) columnFilters(c *fiber.Ctx) map[string]interface{} {
	data, error := qs.Unmarshal(c.OriginalURL())

	if error != nil {
		return map[string]interface{}{}
	}

	result, ok := data["filter"].(map[string]interface{})

	if ok == false {
		return map[string]interface{}{}
	}

	return result
}

/**
 * Get the orderings for the request.
 *
 * @return array
 */
func (p *ResourceExport) orderings(c *fiber.Ctx) map[string]interface{} {
	data, error := qs.Unmarshal(c.OriginalURL())

	if error != nil {
		return map[string]interface{}{}
	}

	result, ok := data["sorter"].(map[string]interface{})

	if ok == false {
		return map[string]interface{}{}
	}

	return result
}

// 处理列表
func (p *ResourceExport) performsList(c *fiber.Ctx, resourceInstance interface{}, lists []map[string]interface{}) []interface{} {
	result := []map[string]interface{}{}

	// 获取列表字段
	exportFields := resourceInstance.(interface {
		ExportFields(c *fiber.Ctx, resourceInstance interface{}) interface{}
	}).ExportFields(c, resourceInstance)

	// 解析字段回调函数
	for _, v := range lists {

		// 给实例的Field属性赋值
		resourceInstance.(interface {
			SetField(fieldData map[string]interface{}) interface{}
		}).SetField(v)

		fields := make(map[string]interface{})
		for _, field := range exportFields.([]interface{}) {

			// 字段名
			name := reflect.
				ValueOf(field).
				Elem().
				FieldByName("Name").String()

			// 获取实例的回调函数
			callback := field.(interface{ GetCallback() interface{} }).GetCallback()

			if callback != nil {
				getCallback := callback.(func() interface{})
				fields[name] = getCallback()
			} else {
				if v[name] != nil {
					fields[name] = v[name]
				}
			}
		}

		result = append(result, fields)
	}

	// 回调处理列表字段值
	return resourceInstance.(interface {
		BeforeExporting(c *fiber.Ctx, result []map[string]interface{}) []interface{}
	}).BeforeExporting(c, result)
}
