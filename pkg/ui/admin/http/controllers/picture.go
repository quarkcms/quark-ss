package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strconv"
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

type Picture struct{}

// 编辑器图片选择
func (p *Picture) GetLists(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	categoryId := c.Query("pictureCategoryId")
	searchName := c.Query("pictureSearchName")
	searchDateStart := c.Query("pictureSearchDate[0]")
	searchDateEnd := c.Query("pictureSearchDate[1]")

	getPage, _ := strconv.Atoi(page)
	model := (&db.Model{}).Model(&models.Picture{})

	if categoryId != "" {
		model.Where("picture_category_id =?", categoryId)
	}

	if searchName != "" {
		model.Where("name LIKE %?%", searchName)
	}

	if searchDateStart != "" && searchDateEnd != "" {
		model.Where("created_at BETWEEN ? AND ?", searchDateStart, searchDateEnd)
	}

	var total int64
	// 获取总数量
	model.Count(&total)

	pictures := []map[string]interface{}{}
	model.Where("status =?", 1).Order("id desc").Limit(8).Offset((getPage - 1) * 8).Find(&pictures)

	for k, v := range pictures {
		if strings.Contains(v["path"].(string), "./") {
			v["path"] = strings.Replace(v["path"].(string), "./storage/app/public", "/storage", -1) + "?timestamp=" + strconv.Itoa(int(time.Now().Unix()))
			pictures[k] = v
		}
	}

	pagination := map[string]interface{}{
		"defaultCurrent": 1,
		"current":        getPage,
		"pageSize":       12,
		"total":          total,
	}

	categorys := []map[string]interface{}{}
	(&db.Model{}).
		Model(&models.PictureCategory{}).
		Where("obj_type = ?", "ADMINID").
		Where("obj_id", utils.Admin(c, "id")).
		Find(&categorys)

	return msg.Success("获取成功", "", map[string]interface{}{
		"pagination": pagination,
		"lists":      pictures,
		"categorys":  categorys,
	})
}

// 上传图片
func (p *Picture) Upload(c *fiber.Ctx) error {
	var result error

	if utils.WebConfig("OSS_OPEN") == "1" {
		result = p.OssUpload(c)
	} else {
		result = p.LocalUpload(c)
	}

	return result
}

// 通过base64字符串上传图片
func (p *Picture) UploadFromBase64(c *fiber.Ctx) error {
	var result error

	if utils.WebConfig("OSS_OPEN") == "1" {
		result = p.OssUploadFromBase64(c)
	} else {
		result = p.LocalUploadFromBase64(c)
	}

	return result
}

// 通过base64字符串上传图片
func (p *Picture) LocalUploadFromBase64(c *fiber.Ctx) error {
	datasource := c.FormValue("file")
	limitW := c.Query("limitW")
	limitH := c.Query("limitH")

	fileArray := strings.Split(datasource, ",")
	if len(fileArray) != 2 {
		return msg.Error("图片格式错误!", "")
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
		return msg.Error("只能上传jpg,jpeg,png,gif格式图片!", "")
	}

	base64Buffer, err := base64.StdEncoding.DecodeString(fileArray[1]) //成图片文件并把文件写入到buffer
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	file := bytes.NewBuffer(base64Buffer) // 必须加一个buffer 不然没有read方法就会报错
	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 限制宽高
	if limitW != "" && limitH != "" {
		w, err := strconv.Atoi(limitW)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		h, err := strconv.Atoi(limitH)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		if imageConfig.Width != w || imageConfig.Height != h {
			return msg.Error("请上传 "+limitW+"*"+limitH+" 尺寸的图片", "")
		}
	}

	filePath := "./storage/app/public/pictures/" + time.Now().Format("20060102") + "/"
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

	id := (&models.Picture{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     filePath + fileName,
		"width":    imageConfig.Width,
		"height":   imageConfig.Height,
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

// 通过base64字符串上传图片
func (p *Picture) OssUploadFromBase64(c *fiber.Ctx) error {
	datasource := c.FormValue("file")
	limitW := c.Query("limitW")
	limitH := c.Query("limitH")

	fileArray := strings.Split(datasource, ",")
	if len(fileArray) != 2 {
		return msg.Error("图片格式错误!", "")
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
		return msg.Error("只能上传jpg,jpeg,png,gif格式图片!", "")
	}

	base64Buffer, err := base64.StdEncoding.DecodeString(fileArray[1]) //成图片文件并把文件写入到buffer
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	file := bytes.NewBuffer(base64Buffer) // 必须加一个buffer 不然没有read方法就会报错
	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 限制宽高
	if limitW != "" && limitH != "" {
		w, err := strconv.Atoi(limitW)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		h, err := strconv.Atoi(limitH)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		if imageConfig.Width != w || imageConfig.Height != h {
			return msg.Error("请上传 "+limitW+"*"+limitH+" 尺寸的图片", "")
		}
	}

	filePath := "pictures/" + time.Now().Format("20060102") + "/"
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
		path = strings.Replace(filePath+fileName, "pictures/", "//"+myDomain+"/pictures/", -1)
	} else {
		path = "//" + ossBucket + "." + endpoint + "/" + filePath + fileName
	}

	id := (&models.Picture{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     path,
		"width":    imageConfig.Width,
		"height":   imageConfig.Height,
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

// 图片上传到本地
func (p *Picture) LocalUpload(c *fiber.Ctx) error {
	file, _ := c.FormFile("file")
	limitW := c.Query("limitW")
	limitH := c.Query("limitH")

	limitType := []string{
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/gif",
	}

	typeAllowed := false

	for _, v := range file.Header["Content-Type"] {
		for _, limit := range limitType {
			if v == limit {
				typeAllowed = true
			}
		}
	}

	// 限制格式
	if typeAllowed == false {
		return msg.Error("只能上传jpg,jpeg,png,gif格式图片!", "")
	}

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

	imageConfig, _, err := image.DecodeConfig(f)

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 限制宽高
	if limitW != "" && limitH != "" {
		w, err := strconv.Atoi(limitW)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		h, err := strconv.Atoi(limitH)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		if imageConfig.Width != w || imageConfig.Height != h {
			return msg.Error("请上传 "+limitW+"*"+limitH+" 尺寸的图片", "")
		}
	}

	filePath := "./storage/app/public/pictures/" + time.Now().Format("20060102") + "/"
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

	picture := map[string]interface{}{}
	(&db.Model{}).Model(&models.Picture{}).Where("md5", fileMd5).Where("name", fileName).First(&picture)

	result := map[string]interface{}{}

	if len(picture) > 0 {
		result = map[string]interface{}{
			"id":   picture["id"],
			"name": picture["name"],
			"url":  strings.Replace(picture["path"].(string), "./storage/app/public", "/storage", -1),
			"size": picture["size"],
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

	id := (&models.Picture{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     filePath + fileNewName,
		"width":    imageConfig.Width,
		"height":   imageConfig.Height,
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

// 图片上传到阿里云OSS
func (p *Picture) OssUpload(c *fiber.Ctx) error {

	file, _ := c.FormFile("file")
	limitW := c.Query("limitW")
	limitH := c.Query("limitH")

	limitType := []string{
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/gif",
	}

	typeAllowed := false

	for _, v := range file.Header["Content-Type"] {
		for _, limit := range limitType {
			if v == limit {
				typeAllowed = true
			}
		}
	}

	// 限制格式
	if typeAllowed == false {
		return msg.Error("只能上传jpg,jpeg,png,gif格式图片!", "")
	}

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

	imageConfig, _, err := image.DecodeConfig(f)

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 限制宽高
	if limitW != "" && limitH != "" {
		w, err := strconv.Atoi(limitW)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		h, err := strconv.Atoi(limitH)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		if imageConfig.Width != w || imageConfig.Height != h {
			return msg.Error("请上传 "+limitW+"*"+limitH+" 尺寸的图片", "")
		}
	}

	filePath := "pictures/" + time.Now().Format("20060102") + "/"
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

	picture := map[string]interface{}{}
	(&db.Model{}).Model(&models.Picture{}).Where("md5", fileMd5).Where("name", fileName).First(&picture)

	result := map[string]interface{}{}

	if len(picture) > 0 {
		result = map[string]interface{}{
			"id":   picture["id"],
			"name": picture["name"],
			"url":  picture["path"],
			"size": picture["size"],
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
		path = strings.Replace(filePath+fileNewName, "pictures/", "//"+myDomain+"/pictures/", -1)
	} else {
		path = "//" + ossBucket + "." + endpoint + "/" + filePath + fileNewName
	}

	id := (&models.Picture{}).InsertGetId(map[string]interface{}{
		"obj_type": "ADMINID",
		"obj_id":   utils.Admin(c, "id"),
		"name":     fileName,
		"size":     fileSize,
		"md5":      fileMd5,
		"path":     path,
		"width":    imageConfig.Width,
		"height":   imageConfig.Height,
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

// 图片下载
func (p *Picture) Download(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return msg.Error("参数错误！", "")
	}

	picture := map[string]interface{}{}
	err := (&db.Model{}).Model(&models.Picture{}).Where("id =?", id).First(&picture).Error

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	if len(picture) == 0 {
		return msg.Error("无此数据！", "")
	}

	path, ok := picture["path"].(string)

	if !ok {
		return msg.Error("路径错误！", "")
	}

	if strings.Contains(path, "//") {
		return c.Redirect(path)
	}

	return c.Redirect(strings.Replace(path, "./storage/app/public", "/storage", -1))
}

// 图片删除
func (p *Picture) Delete(c *fiber.Ctx) error {
	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	if data["id"] == "" {
		return msg.Error("参数错误！", "")
	}

	err := (&db.Model{}).Model(&models.Picture{}).Where("id =?", data["id"]).Delete("").Error

	if err != nil {
		return msg.Error(err.Error(), "")
	} else {
		return msg.Success("操作成功！", "", "")
	}
}

// 图片裁剪
func (p *Picture) Crop(c *fiber.Ctx) error {
	var result error

	data := map[string]interface{}{}
	json.Unmarshal(c.Body(), &data)

	if data["id"] == "" {
		return msg.Error("参数错误！", "")
	}

	if data["file"] == "" {
		return msg.Error("参数错误！", "")
	}

	picture := map[string]interface{}{}
	err := (&db.Model{}).Model(&models.Picture{}).Where("id =?", data["id"]).First(&picture).Error

	if err != nil {
		return msg.Error(err.Error(), "")
	}

	datasource := data["file"]
	limitW := c.Query("limitW")
	limitH := c.Query("limitH")

	fileArray := strings.Split(datasource.(string), ",")
	if len(fileArray) != 2 {
		return msg.Error("图片格式错误!", "")
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
		return msg.Error("只能上传jpg,jpeg,png,gif格式图片!", "")
	}

	base64Buffer, err := base64.StdEncoding.DecodeString(fileArray[1]) //成图片文件并把文件写入到buffer
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	file := bytes.NewBuffer(base64Buffer) // 必须加一个buffer 不然没有read方法就会报错
	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}

	// 限制宽高
	if limitW != "" && limitH != "" {
		w, err := strconv.Atoi(limitW)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		h, err := strconv.Atoi(limitH)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		if imageConfig.Width != w || imageConfig.Height != h {
			return msg.Error("请上传 "+limitW+"*"+limitH+" 尺寸的图片", "")
		}
	}

	fileSize := int64(len(datasource.(string)))

	// 文件md5值
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return msg.Error(err.Error(), "")
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(body))

	if utils.WebConfig("OSS_OPEN") == "1" {

		accessKeyId := utils.WebConfig("OSS_ACCESS_KEY_ID")
		accessKeySecret := utils.WebConfig("OSS_ACCESS_KEY_SECRET")
		endpoint := utils.WebConfig("OSS_ENDPOINT")
		ossBucket := utils.WebConfig("OSS_BUCKET")

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

		err = bucket.PutObject(picture["path"].(string), bytes.NewBuffer(base64Buffer), objectAcl)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		result = (&db.Model{}).Model(&models.Picture{}).Where("id", data["id"]).Updates(map[string]interface{}{
			"md5":  fileMd5,
			"size": fileSize,
		}).Error

	} else {

		// 保存文件
		err = ioutil.WriteFile(picture["path"].(string), base64Buffer, 0666)
		if err != nil {
			return msg.Error(err.Error(), "")
		}

		result = (&db.Model{}).Model(&models.Picture{}).Where("id", data["id"]).Updates(map[string]interface{}{
			"md5":  fileMd5,
			"size": fileSize,
		}).Error
	}

	if result != nil {
		return msg.Error(result.Error(), "")
	} else {
		return msg.Success("操作成功", "", "")
	}
}
