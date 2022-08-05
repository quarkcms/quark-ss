package conf

import (
	"fmt"
	"strings"

	"github.com/Unknwon/goconfig"
)

// 设置值
func Set(key string, value string) bool {
	keys := strings.Split(key, ".")

	cfg, err := goconfig.LoadConfigFile("./" + keys[0] + ".ini")

	if err != nil {
		fmt.Println(err)
	}

	var result bool

	if len(keys) == 2 {
		result = cfg.SetValue(goconfig.DEFAULT_SECTION, keys[1], value)
	} else {
		result = cfg.SetValue(keys[1], keys[2], value)
	}

	return result
}

// 获取值
func Get(key string) string {

	keys := strings.Split(key, ".")

	cfg, err := goconfig.LoadConfigFile("./" + keys[0] + ".ini")

	if err != nil {
		fmt.Println(err)
	}

	var value string
	var valueErr error

	if len(keys) == 2 {
		value, valueErr = cfg.GetValue(goconfig.DEFAULT_SECTION, keys[1])
	} else {
		value, valueErr = cfg.GetValue(keys[1], keys[2])
	}

	if valueErr != nil {
		fmt.Println(valueErr)
	}

	return value
}
