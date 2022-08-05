package admin

import (
	"reflect"

	"github.com/derekstavis/go-qs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 创建列表查询
func (p *Resource) BuildIndexQuery(c *fiber.Ctx, resourceInstance interface{}, query *gorm.DB, search []interface{}, filters []interface{}, columnFilters map[string]interface{}, orderings map[string]interface{}) *gorm.DB {

	// 初始化查询
	query = p.initializeQuery(c, resourceInstance, query)

	// 执行列表查询，这里使用的是透传的实例
	query = resourceInstance.(interface {
		IndexQuery(*fiber.Ctx, *gorm.DB) *gorm.DB
	}).IndexQuery(c, query)

	// 执行搜索查询
	query = p.applySearch(c, query, search)

	// 执行过滤器查询
	query = p.applyFilters(query, filters)

	// 执行表格列上过滤器查询
	query = p.applyColumnFilters(query, columnFilters)

	// 获取默认排序
	defaultOrder := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("IndexOrder").String()

	if defaultOrder == "" {
		defaultOrder = "id desc"
	}

	// 执行排序查询
	query = p.applyOrderings(query, orderings, defaultOrder)

	return query
}

// 创建详情页查询
func (p *Resource) BuildDetailQuery(c *fiber.Ctx, resourceInstance interface{}, query *gorm.DB) *gorm.DB {
	// 初始化查询
	query = p.initializeQuery(c, resourceInstance, query)

	// 执行列表查询，这里使用的是透传的实例
	query = resourceInstance.(interface {
		DetailQuery(*fiber.Ctx, *gorm.DB) *gorm.DB
	}).DetailQuery(c, query)

	return query
}

// 创建导出查询
func (p *Resource) BuildExportQuery(c *fiber.Ctx, resourceInstance interface{}, query *gorm.DB, search []interface{}, filters []interface{}, columnFilters map[string]interface{}, orderings map[string]interface{}) *gorm.DB {

	// 初始化查询
	query = p.initializeQuery(c, resourceInstance, query)

	// 执行列表查询，这里使用的是透传的实例
	query = resourceInstance.(interface {
		ExportQuery(*fiber.Ctx, *gorm.DB) *gorm.DB
	}).ExportQuery(c, query)

	// 执行搜索查询
	query = p.applySearch(c, query, search)

	// 执行过滤器查询
	query = p.applyFilters(query, filters)

	// 执行表格列上过滤器查询
	query = p.applyColumnFilters(query, columnFilters)

	// 获取默认排序
	defaultOrder := reflect.
		ValueOf(resourceInstance).
		Elem().
		FieldByName("IndexOrder").String()

	if defaultOrder == "" {
		defaultOrder = "id desc"
	}

	// 执行排序查询
	query = p.applyOrderings(query, orderings, defaultOrder)

	return query
}

// 初始化查询
func (p *Resource) initializeQuery(c *fiber.Ctx, resourceInstance interface{}, query *gorm.DB) *gorm.DB {

	return resourceInstance.(interface {
		Query(*fiber.Ctx, *gorm.DB) *gorm.DB
	}).Query(c, query)
}

// 执行搜索表单查询
func (p *Resource) applySearch(c *fiber.Ctx, query *gorm.DB, search []interface{}) *gorm.DB {

	data, error := qs.Unmarshal(c.OriginalURL())

	if error != nil {
		return query
	}

	result, ok := data["search"].(map[string]interface{})

	if ok == false {
		return query
	}

	for _, v := range search {

		// 获取字段
		column := v.(interface {
			GetColumn(search interface{}) string
		}).GetColumn(v) // 字段名，支持数组

		value := result[column]

		if value != nil {
			query = v.(interface {
				Apply(*fiber.Ctx, *gorm.DB, interface{}) *gorm.DB
			}).Apply(c, query, value)
		}
	}

	return query
}

// 执行表格列上过滤器查询
func (p *Resource) applyColumnFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {

	if len(filters) == 0 {
		return query
	}

	for k, v := range filters {
		if v != "" {

			values := []string{}
			for _, subValue := range v.(map[string]interface{}) {
				values = append(values, subValue.(string))
			}

			query = query.Where(k+" IN ?", values)
		}
	}

	return query
}

// 执行过滤器查询
func (p *Resource) applyFilters(query *gorm.DB, filters []interface{}) *gorm.DB {
	// todo
	return query
}

// 执行排序查询
func (p *Resource) applyOrderings(query *gorm.DB, orderings map[string]interface{}, defaultOrder string) *gorm.DB {

	if len(orderings) == 0 {
		return query.Order(defaultOrder)
	}

	var order clause.OrderByColumn

	for key, v := range orderings {

		if v == "descend" {
			order = clause.OrderByColumn{Column: clause.Column{Name: key}, Desc: true}
		} else {
			order = clause.OrderByColumn{Column: clause.Column{Name: key}, Desc: false}
		}

		query = query.Order(order)
	}

	return query
}

// 全局查询
func (p *Resource) Query(c *fiber.Ctx, query *gorm.DB) *gorm.DB {

	return query
}

// 列表查询
func (p *Resource) IndexQuery(c *fiber.Ctx, query *gorm.DB) *gorm.DB {

	return query
}

// 详情查询
func (p *Resource) DetailQuery(c *fiber.Ctx, query *gorm.DB) *gorm.DB {

	return query
}

// 导出查询
func (p *Resource) ExportQuery(c *fiber.Ctx, query *gorm.DB) *gorm.DB {

	return query
}
