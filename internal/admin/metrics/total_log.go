package metrics

import (
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/ui/admin/metrics"
	"github.com/quarkcms/quark-go/pkg/ui/component/statistic"
)

type TotalLog struct {
	metrics.Value
}

// 初始化
func (p *TotalLog) Init() *TotalLog {
	p.Title = "日志数量"
	p.Col = 6

	return p
}

// 计算数值
func (p *TotalLog) Calculate() *statistic.Component {

	return p.
		Init().
		Count((&db.Model{}).Model(&models.ActionLog{})).
		SetValueStyle(map[string]string{"color": "#999999"})
}
