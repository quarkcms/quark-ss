package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
	"github.com/quarkcms/quark-go/pkg/ui/component/card"
	"github.com/quarkcms/quark-go/pkg/ui/component/descriptions"
	"github.com/quarkcms/quark-go/pkg/ui/component/grid"
	"github.com/quarkcms/quark-go/pkg/ui/component/statistic"
)

// 资源结构体
type Dashboard struct {
	Layout
	Title    string
	SubTitle string
}

// 解析卡片组件
func (p *Dashboard) CardComponentRender(c *fiber.Ctx, dashboard interface{}) interface{} {
	cards := dashboard.(interface{ Cards(*fiber.Ctx) []any }).Cards(c)
	var cols []interface{}
	var rows []interface{}
	var colNum int = 0

	for key, v := range cards {

		// 断言statistic组件类型
		statistic, ok := v.(interface{ Calculate() *statistic.Component })
		item := (&card.Component{}).Init()
		if ok {
			item = item.SetBody(statistic.Calculate())
		} else {

			// 断言descriptions组件类型
			descriptions, ok := v.(interface {
				Calculate() *descriptions.Component
			})
			if ok {
				item = item.SetBody(descriptions.Calculate())
			}
		}

		// struct转换map
		vMap := utils.StructToMap(v).(map[string]interface{})

		// float64转换成int
		col := int(vMap["Col"].(float64))

		colInfo := (&grid.Col{}).Init().SetSpan(col).SetBody(item)
		cols = append(cols, colInfo)

		colNum = colNum + col

		if colNum%24 == 0 {
			row := (&grid.Row{}).Init().SetGutter(8).SetBody(cols)
			if key != 1 {
				row = row.SetStyle(map[string]interface{}{"marginTop": "20px"})
			}
			rows = append(rows, row)
			cols = nil
		}
	}

	if cols != nil {
		row := (&grid.Row{}).Init().SetGutter(8).SetBody(cols)
		if colNum > 24 {
			row = row.SetStyle(map[string]interface{}{"marginTop": "20px"})
		}
		rows = append(rows, row)
	}

	return rows
}

// 仪表盘组件渲染
func (p *Dashboard) DashboardComponentRender(c *fiber.Ctx, resourceInstance interface{}) interface{} {

	return resourceInstance.(interface {
		CardComponentRender(*fiber.Ctx, interface{}) interface{}
	}).CardComponentRender(c, resourceInstance)
}
