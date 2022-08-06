package database

import (
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/framework/hash"
)

type Seed struct{}

// 添加初始数据
func (p *Seed) Handle() {

	p.adminTableSeeder()
	p.configTableSeeder()
	p.menuTableSeeder()
}

// 管理员表
func (p *Seed) adminTableSeeder() {
	seeders := []models.Admin{
		{Username: "administrator", Nickname: "超级管理员", Email: "admin@yourweb.com", Phone: "10086", Password: hash.Make("123456"), Sex: 1, Status: 1},
	}

	(&db.Model{}).DB().Create(&seeders)
}

// 配置表
func (p *Seed) configTableSeeder() {
	seeders := []models.Config{
		{Title: "网站名称", Type: "text", Name: "WEB_SITE_NAME", Sort: 0, GroupName: "基本", Value: "QuarkCMS", Remark: "", Status: 1},
		{Title: "关键字", Type: "text", Name: "WEB_SITE_KEYWORDS", Sort: 0, GroupName: "基本", Value: "QuarkCMS", Remark: "", Status: 1},
		{Title: "描述", Type: "textarea", Name: "WEB_SITE_DESCRIPTION", Sort: 0, GroupName: "基本", Value: "QuarkCMS", Remark: "", Status: 1},
		{Title: "Logo", Type: "picture", Name: "WEB_SITE_LOGO", Sort: 0, GroupName: "基本", Value: "", Remark: "", Status: 1},
		{Title: "统计代码", Type: "textarea", Name: "WEB_SITE_SCRIPT", Sort: 0, GroupName: "基本", Value: "", Remark: "", Status: 1},
		{Title: "网站版权", Type: "text", Name: "WEB_SITE_COPYRIGHT", Sort: 0, GroupName: "基本", Value: "© Company 2018", Remark: "", Status: 1},
		{Title: "开启SSL", Type: "switch", Name: "SSL_OPEN", Sort: 0, GroupName: "基本", Value: "0", Remark: "", Status: 1},
		{Title: "开启网站", Type: "switch", Name: "WEB_SITE_OPEN", Sort: 0, GroupName: "基本", Value: "1", Remark: "", Status: 1},
		{Title: "KeyID", Type: "text", Name: "OSS_ACCESS_KEY_ID", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "你的AccessKeyID", Status: 1},
		{Title: "KeySecret", Type: "text", Name: "OSS_ACCESS_KEY_SECRET", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "你的AccessKeySecret", Status: 1},
		{Title: "EndPoint", Type: "text", Name: "OSS_ENDPOINT", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "地域节点", Status: 1},
		{Title: "Bucket域名", Type: "text", Name: "OSS_BUCKET", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "", Status: 1},
		{Title: "自定义域名", Type: "text", Name: "OSS_MYDOMAIN", Sort: 0, GroupName: "阿里云存储", Value: "", Remark: "例如：oss.web.com", Status: 1},
		{Title: "开启云存储", Type: "switch", Name: "OSS_OPEN", Sort: 0, GroupName: "阿里云存储", Value: "0", Remark: "", Status: 1},
	}

	(&db.Model{}).DB().Create(&seeders)
}

// 菜单表
func (p *Seed) menuTableSeeder() {
	seeders := []models.Menu{
		{Id: 1, Name: "控制台", GuardName: "admin", Icon: "icon-home", Type: "default", Pid: 0, Sort: -1, Path: "/dashboard", Show: 1, Status: 1},
		{Id: 2, Name: "主页", GuardName: "admin", Icon: "", Type: "engine", Pid: 1, Sort: 0, Path: "admin/dashboard/index", Show: 1, Status: 1},
		{Id: 3, Name: "管理员", GuardName: "admin", Icon: "icon-admin", Type: "default", Pid: 0, Sort: 0, Path: "/admin", Show: 1, Status: 1},
		{Id: 4, Name: "管理员列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 3, Sort: 0, Path: "admin/admin/index", Show: 1, Status: 1},
		{Id: 5, Name: "权限列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 3, Sort: 0, Path: "admin/permission/index", Show: 1, Status: 1},
		{Id: 6, Name: "角色列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 3, Sort: 0, Path: "admin/role/index", Show: 1, Status: 1},
		{Id: 7, Name: "系统配置", GuardName: "admin", Icon: "icon-setting", Type: "default", Pid: 0, Sort: 0, Path: "/system", Show: 1, Status: 1},
		{Id: 8, Name: "设置管理", GuardName: "admin", Icon: "", Type: "default", Pid: 7, Sort: 0, Path: "/system/config", Show: 1, Status: 1},
		{Id: 9, Name: "网站设置", GuardName: "admin", Icon: "", Type: "engine", Pid: 8, Sort: 0, Path: "admin/webConfig/setting-form", Show: 1, Status: 1},
		{Id: 10, Name: "配置管理", GuardName: "admin", Icon: "", Type: "engine", Pid: 8, Sort: 0, Path: "admin/config/index", Show: 1, Status: 1},
		{Id: 11, Name: "菜单管理", GuardName: "admin", Icon: "", Type: "engine", Pid: 7, Sort: 0, Path: "admin/menu/index", Show: 1, Status: 1},
		{Id: 12, Name: "操作日志", GuardName: "admin", Icon: "", Type: "engine", Pid: 7, Sort: 0, Path: "admin/actionLog/index", Show: 1, Status: 1},
		{Id: 13, Name: "附件空间", GuardName: "admin", Icon: "icon-attachment", Type: "default", Pid: 0, Sort: 0, Path: "/attachment", Show: 1, Status: 1},
		{Id: 14, Name: "文件管理", GuardName: "admin", Icon: "", Type: "engine", Pid: 13, Sort: 0, Path: "admin/file/index", Show: 1, Status: 1},
		{Id: 15, Name: "图片管理", GuardName: "admin", Icon: "", Type: "engine", Pid: 13, Sort: 0, Path: "admin/picture/index", Show: 1, Status: 1},
		{Id: 16, Name: "我的账号", GuardName: "admin", Icon: "icon-user", Type: "default", Pid: 0, Sort: 0, Path: "/account", Show: 1, Status: 1},
		{Id: 17, Name: "个人设置", GuardName: "admin", Icon: "", Type: "engine", Pid: 16, Sort: 0, Path: "admin/account/setting-form", Show: 1, Status: 1},
		{Id: 18, Name: "服务端", GuardName: "admin", Icon: "icon-sever", Type: "default", Pid: 0, Sort: -1, Path: "/server", Show: 1, Status: 1},
		{Id: 19, Name: "服务列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 18, Sort: 0, Path: "admin/server/index", Show: 1, Status: 1},
	}

	(&db.Model{}).DB().Create(&seeders)
}
