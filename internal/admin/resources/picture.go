package resources

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/admin/actions"
	"github.com/quarkcms/quark-go/internal/admin/searches"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/ui/admin"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/table"
)

type Picture struct {
	admin.Resource
}

// 初始化
func (p *Picture) Init() interface{} {

	// 标题
	p.Title = "图片"

	// 模型
	p.Model = &models.Picture{}

	// 分页
	p.PerPage = 10

	return p
}

// 字段
func (p *Picture) Fields(c *fiber.Ctx) []interface{} {
	field := &admin.Field{}

	return []interface{}{
		field.ID("id", "ID"),
		field.Text("path", "显示", func() interface{} {

			return "<img src='" + utils.GetPicture(c, p.Field["id"]) + "' width=50 height=50 />"
		}),
		field.Text("name", "名称").SetColumn(func(column *table.Column) *table.Column {
			return column.SetEllipsis(true)
		}),
		field.Text("size", "大小").
			SetColumn(func(column *table.Column) *table.Column {
				return column.SetSorter(true)
			}),
		field.Text("width", "宽度"),
		field.Text("height", "高度"),
		field.Text("ext", "扩展名"),
		field.Datetime("created_at", "上传时间", func() interface{} {
			if p.Field["created_at"] == nil {
				return p.Field["created_at"]
			}

			return p.Field["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		}),
	}
}

// 搜索
func (p *Picture) Searches(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
		(&searches.DateTimeRange{}).Init("created_at", "上传时间"),
	}
}

// 行为
func (p *Picture) Actions(c *fiber.Ctx) []interface{} {
	return []interface{}{
		(&actions.Delete{}).Init("批量删除"),
		(&actions.Delete{}).Init("删除"),
	}
}
