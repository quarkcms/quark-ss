package fields

import "github.com/quarkcms/quark-go/pkg/ui/component"

type Datetime struct {
	Item
}

// 初始化
func (p *Datetime) Init() *Datetime {
	p.Component = "datetimeField"
	p.InitItem().SetKey(component.DEFAULT_KEY, component.DEFAULT_CRYPT)

	return p
}
