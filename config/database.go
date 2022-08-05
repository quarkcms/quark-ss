package config

import (
	"github.com/quarkcms/quark-go/pkg/framework/env"
)

var Database = map[string]interface{}{

	// Default Database Connection Name
	"default": env.Get("DB_CONNECTION", "mysql"),

	// Database Connections
	"connections": map[string]interface{}{

		// sqlite
		"sqlite": map[string]interface{}{

			// 数据库
			"database": env.Get("DB_DATABASE", "data.db"),
		},

		// mysql
		"mysql": map[string]interface{}{

			// 地址
			"host": env.Get("DB_HOST", "127.0.0.1"),

			// 端口
			"port": env.Get("DB_PORT", "3306"),

			// 数据库
			"database": env.Get("DB_DATABASE", "quarkgo"),

			// 用户名
			"username": env.Get("DB_USERNAME", "root"),

			// 密码
			"password": env.Get("DB_PASSWORD"),

			// 编码
			"charset": "utf8",
		},
	},

	// redis配置
	"redis": map[string]interface{}{

		// 地址
		"host": env.Get("REDIS_HOST", "127.0.0.1"),

		// 密码
		"password": env.Get("REDIS_PASSWORD"),

		// 端口
		"port": env.Get("REDIS_PORT", "6379"),

		// 数据库
		"database": 0,
	},
}
