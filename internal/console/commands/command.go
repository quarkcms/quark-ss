package commands

type Command struct {
	Signature   string
	Description string
}

// 初始化
func (p *Command) Init() *Command {
	return p
}

// 执行命令
func (p *Command) Handle() {

}
