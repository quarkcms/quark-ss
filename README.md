## 介绍
QuarkGO 是一个基于golang管理后台的工具；它提供的丰富组件，能帮助您使用很少的代码就能搭建出功能完善的管理后台。

## 系统特性

**内置功能**
* 管理员管理
* 用户管理
* 权限系统
* 菜单管理
* 系统配置
* 操作日志
* 附件管理

**内置组件**
* Layout组件
* Container组件
* Card组件
* Table组件
* Form组件
* Show组件
* TabForm组件
* ...

## 安装

需要go1.18+，首先确保安装好了环境。

1、重命名.env.example 改为 .env 

2、编辑.env文件，配置数据库信息

3、执行下面的命令完成安装：
``` bash
# 第一步，创建vendor目录
go mod vendor

# 第二步，安装依赖:
go mod tidy

# 第三步，然后运行下面的命令启动服务：
go run main.go
```

后台地址： http://127.0.0.1:3000/admin/

默认用户名：administrator 密码：123456

## 技术支持
为了避免打扰作者日常工作，你可以在Github上提交 [Issues](https://github.com/quarkcms/quark-go/issues)

相关教程，你可以查看 [在线文档](http://www.quarkcms.com/quark-go/)

## License
QuarkAdmin is licensed under The MIT License (MIT).