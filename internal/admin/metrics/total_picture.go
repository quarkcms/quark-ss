package metrics

import (
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin/metrics"
	"github.com/quarkcms/quark-go/pkg/ui/component/statistic"
)

type TotalPicture struct {
	metrics.Value
}

// 初始化
func (p *TotalPicture) Init() *TotalPicture {
	p.Title = "图片数量"
	p.Col = 6

	return p
}

// 计算数值
func (p *TotalPicture) Calculate() *statistic.Component {

	return p.
		Init().
		Count((&db.Model{}).Model(&models.Picture{})).
		SetValueStyle(map[string]string{"color": "#cf1322"})
}
