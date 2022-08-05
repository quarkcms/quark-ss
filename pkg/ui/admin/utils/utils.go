package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/framework/token"
	"github.com/xuri/excelize/v2"
)

// 获取管理员Token
func GetAdminToken(c *fiber.Ctx) string {
	if c.Query("token") != "" {
		return c.Query("token")
	}

	header := c.GetReqHeaders()
	getToken := strings.Split(header["Authorization"], " ")

	if len(getToken) != 2 {
		return ""
	}

	return getToken[1]
}

// 获取当前登录管理员信息
func Admin(c *fiber.Ctx, field string) interface{} {
	getToken := GetAdminToken(c)
	userInfo, err := token.Parse(getToken)
	if err != nil {
		return nil
	}

	if value, ok := userInfo[field].(float64); ok {
		return int(value)
	}

	return userInfo[field]
}

// 数据集转换成Tree
func ListToTree(list []map[string]interface{}, pk string, pid string, child string, root int) []interface{} {
	var treeList []interface{}
	for _, v := range list {
		if v[pid] == root {
			childNode := ListToTree(list, pk, pid, child, v[pk].(int))
			if childNode != nil {
				v[child] = childNode
			}
			treeList = append(treeList, v)
		}
	}

	return treeList
}

// Tree转换为有序列表
func TreeToOrderedList(tree []interface{}, level int, field string, child string) []interface{} {
	var list []interface{}
	for _, v := range tree {
		if value, ok := v.(map[string]interface{}); ok {
			value[field] = strings.Repeat("—", level) + value[field].(string)
			list = append(list, value)
			if value[child] != nil {
				if childValue, ok := value[child].([]interface{}); ok {
					children := TreeToOrderedList(childValue, level+1, field, child)
					list = append(list, children...)
				}
			}
		}
	}

	return list
}

// struct转map
func StructToMap(v any) any {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	var mapResult any
	json.Unmarshal(jsonBytes, &mapResult)

	return mapResult
}

// 存储权限
var Permissions []string

// 设置权限
func SetPermissions(permissions []string) {
	Permissions = permissions
}

// 获取权限
func GetPermissions() []string {
	return Permissions
}

// 判断路径是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 存储网站配置
var webConfig = make(map[string]string)

// 配置结构体
type Config struct {
	Id        int
	Title     string
	Type      string
	Name      string
	Sort      int
	GroupName string
	Value     string
	Remark    string
	Status    int
}

// 刷新网站配置
func RefreshWebConfig() {
	configs := []map[string]interface{}{}

	(&db.Model{}).Model(&Config{}).Where("status", 1).Find(&configs)

	for _, config := range configs {
		webConfig[config["name"].(string)] = config["value"].(string)
	}
}

// 获取网站配置信息
func WebConfig(key string) string {

	// 刷新网站配置
	if len(webConfig) == 0 {
		RefreshWebConfig()
	}

	return webConfig[key]
}

// 图片结构体
type Picture struct {
	Id                int
	ObjType           string
	ObjId             int
	PictureCategoryId int
	Sort              int
	Name              string
	Size              string
	Width             int
	Height            int
	Ext               string
	Path              string
	Md5               string
	Status            int
}

// 获取图片
func GetPicture(c *fiber.Ctx, id interface{}) string {
	path := ""

	if getId, ok := id.(string); ok {

		// json字符串
		if strings.Contains(getId, "{") {
			var jsonData interface{}
			json.Unmarshal([]byte(getId), &jsonData)

			// 如果为map
			if mapData, ok := jsonData.(map[string]interface{}); ok {
				path = mapData["url"].(string)
			}

			// 如果为数组，返回第一个key的path
			if arrayData, ok := jsonData.([]map[string]interface{}); ok {
				path = arrayData[0]["url"].(string)
			}
		}

		if strings.Contains(path, "//") {
			return path
		}

		if strings.Contains(path, "./") {
			return strings.Replace(path, "./storage/app/public", "/storage", -1)
		}
	}

	picture := map[string]interface{}{}
	(&db.Model{}).Model(&Picture{}).Where("id", id).Where("status", 1).First(&picture)

	if len(picture) > 0 {
		path = picture["path"].(string)

		if strings.Contains(path, "//") {
			return path
		}

		if strings.Contains(path, "./") {
			return strings.Replace(path, "./storage/app/public", "/storage", -1)
		}
	}

	return "/admin/default.png"
}

// 文件结构体
type File struct {
	Id             int
	ObjType        string
	ObjId          int
	FileCategoryId int
	Sort           int
	Name           string
	Size           string
	Ext            string
	Path           string
	Md5            string
	Status         bool
}

// 获取文件
func GetFile(c *fiber.Ctx, id interface{}) string {
	path := ""

	if getId, ok := id.(string); ok {

		// json字符串
		if strings.Contains(getId, "{") {
			var jsonData interface{}
			json.Unmarshal([]byte(getId), &jsonData)

			// 如果为map
			if mapData, ok := jsonData.(map[string]interface{}); ok {
				path = mapData["path"].(string)
			}

			// 如果为数组，返回第一个key的path
			if arrayData, ok := jsonData.([]map[string]interface{}); ok {
				path = arrayData[0]["path"].(string)
			}
		}

		if strings.Contains(path, "//") {
			return path
		}

		if strings.Contains(path, "./") {
			return strings.Replace(path, "./storage/app/public", "/storage", -1)
		}
	}

	file := map[string]interface{}{}
	(&db.Model{}).Model(&File{}).Where("id", id).Where("status", 1).First(&file)

	if len(file) > 0 {
		path = file["path"].(string)

		if strings.Contains(path, "//") {
			return path
		}

		if strings.Contains(path, "./") {
			return strings.Replace(path, "./storage/app/public", "/storage", -1)
		}
	}

	return ""
}

// 获取导入文件数据
func Import(fileId int) ([][]interface{}, error) {
	filePath := ""
	file := map[string]interface{}{}
	importData := [][]interface{}{}

	(&db.Model{}).Model(&File{}).Where("id", fileId).Where("status", 1).First(&file)

	if len(file) > 0 {
		filePath = file["path"].(string)
	}

	if filePath == "" {
		return importData, errors.New("参数错误！")
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return importData, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return importData, err
	}

	for _, row := range rows {
		getRows := []interface{}{}
		for _, colCell := range row {
			getRows = append(getRows, colCell)
		}
		importData = append(importData, getRows)
	}

	return importData, err
}
