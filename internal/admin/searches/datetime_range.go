package searches

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/searches"
	"gorm.io/gorm"
)

type DateTimeRange struct {
	searches.DatetimeRange
}

// 初始化
func (p *DateTimeRange) Init(column string, name string) *DateTimeRange {
	p.ParentInit()
	p.Column = column
	p.Name = name

	return p
}

// 执行查询
func (p *DateTimeRange) Apply(c *fiber.Ctx, query *gorm.DB, value interface{}) *gorm.DB {
	values, ok := value.(map[string]interface{})

	if ok == false {
		return query
	}

	return query.Where(p.Column+" BETWEEN ? AND ?", values["0"], values["1"])
}
