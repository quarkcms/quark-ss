## 介绍
SS web pannel

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

# 第四步，执行迁移：
install
```

后台地址： http://127.0.0.1:3000/admin/

默认用户名：administrator 密码：123456

## License
QuarkSS is licensed under The MIT License (MIT).