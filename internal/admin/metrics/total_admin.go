package metrics

import (
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin/metrics"
	"github.com/quarkcms/quark-go/pkg/ui/component/statistic"
)

type TotalAdmin struct {
	metrics.Value
}

// 初始化
func (p *TotalAdmin) Init() *TotalAdmin {
	p.Title = "管理员数量"
	p.Col = 6

	return p
}

// 计算数值
func (p *TotalAdmin) Calculate() *statistic.Component {

	return p.
		Init().
		Count((&db.Model{}).Model(&models.Admin{})).
		SetValueStyle(map[string]string{"color": "#3f8600"})
}
