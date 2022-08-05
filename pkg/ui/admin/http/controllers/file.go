package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/internal/models"
	"github.com/quarkcms/quark-go/pkg/framework/db"
	"github.com/quarkcms/quark-go/pkg/framework/msg"
	"github.com/quarkcms/quark-go/pkg/framework/rand"
	"github.com/quarkcms/quark-go/pkg/ui/admin/utils"
)

type File struct{}

// 上传文件
func (p *File) Upload(c *fiber.Ctx) error {
	var result error

	if utils.WebConfig("OSS_OPEN") == "1" {
		result = p.OssUpload(c)
	} else {
		result = p.LocalUpload(c)
	}

	return result
}

// 通过base64字符串上传文件
func (p *File) UploadFromBase64(c *fiber.Ctx) error {
	var result error

	if utils.WebConfig("OSS_OPEN") == "1" {
		result = p.OssUploadFromBase64(c)
	} else {
		result = p.LocalUploadFromBase64(c)
	}

	return result
}

// 通过base64字符串上传文件
func (p *File) LocalUploadFromBase64(c *fiber.Ctx) error {
	datasource := c.FormValue("file")

	fileArray := strings.Split(datasource, ",")
	if len(fileArray) != 2 {
		return msg.Error("文件格式错误!", "")
	}

	fileExt := ""
	switch fileArray[0] {
	case "data:image/jpg;base64":
		fileExt = "jpg"
	case "data:image/jpeg;base64":
		fileExt = "jpeg"
	case "data:image/png;base64":
		fileExt = "png"
	case "data:image/gif;base64":
		fileExt = "gif"
	}

	// 限制格式
	if fileExt == "" {
		return msg.Error("只能上传jpg,jpeg,png,gif格式文件!", "")
	}

	base64Buffer, err := base64.StdEncoding.DecodeString(fileArray[1]) //成文件文件并把文件写入到buffer
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	file := bytes.NewBuffer(base64Buffer) // 必须加一个buffer 不然没有read方法就会报错

	filePath := "./storage/app/public/files/" + time.Now().Format("20060102") + "/"
	fileName := rand.MakeAlphanumeric(40) + "." + fileExt
	fileSize := int64(len(datasource))

	// 文件md5值
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(body))

	// 不存在路径，则创建
	if utils.PathExist(filePath) == false {
		err := os.MkdirAll(filePath, 0666)
		if err != nil {
			return msg.Error(err.Error(), "")
		}
	}

	// 保存文件
	err = ioutil.WriteFile(filePath+fileName, base64Buffer, 0666)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	id := (&models.File{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     filePath + fileName,
		"ext":      fileExt,
	})

	if id == 0 {
		return msg.Error("上传失败！", "")
	}

	result := map[string]interface{}{
		"id":   id,
		"name": fileName,
		"url":  strings.Replace(filePath+fileName, "./storage/app/public", "/storage", -1),
		"size": fileSize,
	}

	return msg.Success("上传成功！", "", result)
}

// 通过base64字符串上传文件
func (p *File) OssUploadFromBase64(c *fiber.Ctx) error {
	datasource := c.FormValue("file")

	fileArray := strings.Split(datasource, ",")
	if len(fileArray) != 2 {
		return msg.Error("文件格式错误!", "")
	}

	fileExt := ""
	switch fileArray[0] {
	case "data:image/jpg;base64":
		fileExt = "jpg"
	case "data:image/jpeg;base64":
		fileExt = "jpeg"
	case "data:image/png;base64":
		fileExt = "png"
	case "data:image/gif;base64":
		fileExt = "gif"
	}

	// 限制格式
	if fileExt == "" {
		return msg.Error("只能上传jpg,jpeg,png,gif格式文件!", "")
	}

	base64Buffer, err := base64.StdEncoding.DecodeString(fileArray[1]) //成文件文件并把文件写入到buffer
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	file := bytes.NewBuffer(base64Buffer) // 必须加一个buffer 不然没有read方法就会报错

	filePath := "files/" + time.Now().Format("20060102") + "/"
	fileName := rand.MakeAlphanumeric(40) + "." + fileExt
	fileSize := int64(len(datasource))

	// 文件md5值
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(body))

	accessKeyId := utils.WebConfig("OSS_ACCESS_KEY_ID")
	accessKeySecret := utils.WebConfig("OSS_ACCESS_KEY_SECRET")
	endpoint := utils.WebConfig("OSS_ENDPOINT")
	ossBucket := utils.WebConfig("OSS_BUCKET")
	myDomain := utils.WebConfig("OSS_MYDOMAIN")

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 指定Object访问权限
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	err = bucket.PutObject(filePath+fileName, bytes.NewBuffer(base64Buffer), objectAcl)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	path := ""
	if myDomain != "" {
		path = strings.Replace(filePath+fileName, "files/", "//"+myDomain+"/files/", -1)
	} else {
		path = "//" + ossBucket + "." + endpoint + "/" + filePath + fileName
	}

	id := (&models.File{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     path,
		"ext":      fileExt,
	})

	if id == 0 {
		return msg.Error("上传失败！", "")
	}

	result := map[string]interface{}{
		"id":   id,
		"name": fileName,
		"url":  path,
		"size": fileSize,
	}

	return msg.Success("上传成功！", "", result)
}

// 文件上传到本地
func (p *File) LocalUpload(c *fiber.Ctx) error {
	file, _ := c.FormFile("file")

	f, err := file.Open()

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	filePath := "./storage/app/public/files/" + time.Now().Format("20060102") + "/"
	fileName := file.Filename
	fileSize := file.Size

	fileNames := strings.Split(fileName, ".")
	if len(fileNames) <= 1 {
		return msg.Error("无法获取文件扩展名！", "")
	}

	fileExt := fileNames[len(fileNames)-1]
	fileNewName := rand.MakeAlphanumeric(40) + "." + fileExt

	// 文件md5值
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return msg.Error(err.Error(), "")
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(body))

	fileInfo := map[string]interface{}{}
	(&db.Model{}).Model(&models.File{}).Where("md5", fileMd5).Where("name", fileName).First(&fileInfo)

	result := map[string]interface{}{}

	if len(fileInfo) > 0 {
		result = map[string]interface{}{
			"id":   fileInfo["id"],
			"name": fileInfo["name"],
			"url":  strings.Replace(fileInfo["path"].(string), "./storage/app/public", "/storage", -1),
			"size": fileInfo["size"],
		}

		return msg.Success("上传成功！", "", result)
	}

	// 不存在路径，则创建
	if utils.PathExist(filePath) == false {
		err := os.MkdirAll(filePath, 0666)
		if err != nil {
			return msg.Error(err.Error(), "")
		}
	}

	// 保存文件
	err = c.SaveFile(file, filePath+fileNewName)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	id := (&models.File{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     filePath + fileNewName,
		"ext":      fileExt,
	})

	if id == 0 {
		return msg.Error("上传失败！", "")
	}

	result = map[string]interface{}{
		"id":   id,
		"name": fileName,
		"url":  strings.Replace(filePath+fileNewName, "./storage/app/public", "/storage", -1),
		"size": fileSize,
	}

	return msg.Success("上传成功！", "", result)
}

// 文件上传到阿里云OSS
func (p *File) OssUpload(c *fiber.Ctx) error {
	file, _ := c.FormFile("file")

	f, err := file.Open()

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	filePath := "files/" + time.Now().Format("20060102") + "/"
	fileName := file.Filename
	fileSize := file.Size

	fileNames := strings.Split(fileName, ".")
	if len(fileNames) <= 1 {
		return msg.Error("无法获取文件扩展名！", "")
	}

	fileExt := fileNames[len(fileNames)-1]
	fileNewName := rand.MakeAlphanumeric(40) + "." + fileExt

	// 文件md5值
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return msg.Error(err.Error(), "")
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(body))

	fileInfo := map[string]interface{}{}
	(&db.Model{}).Model(&models.File{}).Where("md5", fileMd5).Where("name", fileName).First(&fileInfo)

	result := map[string]interface{}{}

	if len(fileInfo) > 0 {
		result = map[string]interface{}{
			"id":   fileInfo["id"],
			"name": fileInfo["name"],
			"url":  fileInfo["path"],
			"size": fileInfo["size"],
		}

		return msg.Success("上传成功！", "", result)
	}

	accessKeyId := utils.WebConfig("OSS_ACCESS_KEY_ID")
	accessKeySecret := utils.WebConfig("OSS_ACCESS_KEY_SECRET")
	endpoint := utils.WebConfig("OSS_ENDPOINT")
	ossBucket := utils.WebConfig("OSS_BUCKET")
	myDomain := utils.WebConfig("OSS_MYDOMAIN")

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 指定Object访问权限
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	getFile, err := file.Open()

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	defer func() {
		e := getFile.Close()
		if err == nil {
			err = e
		}
	}()

	err = bucket.PutObject(filePath+fileNewName, getFile, objectAcl)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	path := ""
	if myDomain != "" {
		path = strings.Replace(filePath+fileNewName, "files/", "//"+myDomain+"/files/", -1)
	} else {
		path = "//" + ossBucket + "." + endpoint + "/" + filePath + fileNewName
	}

	id := (&models.File{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     path,
		"ext":      fileExt,
	})

	if id == 0 {
		return msg.Error("上传失败！", "")
	}

	result = map[string]interface{}{
		"id":   id,
		"name": fileName,
		"url":  path,
		"size": fileSize,
	}

	return msg.Success("上传成功！", "", result)
}

// 文件下载
func (p *File) Download(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return msg.Error("参数错误！", "")
	}

	fileInfo := map[string]interface{}{}
	err := (&db.Model{}).Model(&models.File{}).Where("id =?", id).First(&fileInfo).Error

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	if len(fileInfo) == 0 {
		return msg.Error("无此数据！", "")
	}

	path, ok := fileInfo["path"].(string)

	if !ok {
		return msg.Error("路径错误！", "")
	}

	if strings.Contains(path, "//") {
		return c.Redirect(path)
	}

	return c.Redirect(strings.Replace(path, "./storage/app/public", "/storage", -1))
}
