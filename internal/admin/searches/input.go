package searches

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/searches"
	"gorm.io/gorm"
)

type Input struct {
	searches.Search
}

// 初始化
func (p *Input) Init(column string, name string) *Input {
	p.ParentInit()
	p.Column = column
	p.Name = name

	return p
}

// 执行查询
func (p *Input) Apply(c *fiber.Ctx, query *gorm.DB, value interface{}) *gorm.DB {
	return query.Where(p.Column+" LIKE ?", "%"+value.(string)+"%")
}
