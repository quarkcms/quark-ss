package config

import (
	"strings"

	"github.com/quarkcms/quark-go/config"
)

var configs = make(map[string]interface{})

func init() {
	configs["app"] = config.App
	configs["database"] = config.Database
	configs["admin"] = config.Admin
}

// 设置值
func Set(key string, value interface{}) {
	keys := strings.Split(key, ".")

	configs[keys[0]] = map[string]interface{}{
		keys[1]: value,
	}
}

// 获取值
func Get(key string) interface{} {
	keys := strings.Split(key, ".")

	return parseKey(keys, configs[keys[0]])
}

// 解析key
func parseKey(keys []string, value interface{}) interface{} {
	for k, v := range keys {
		if k != 0 {
			value = value.(map[string]interface{})[v]
		}
	}

	return value
}
