package public

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"github.com/google/uuid"

	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	"github.com/tanganyu1114/ansible-role-manager/common/file_store"
)

type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

const path = "static/uploadfile/"

type File struct {
	apis.Api
}

// @Summary 上传图片
// @Description 获取JSON
// @Tags 公共接口
// @Accept multipart/form-data
// @Param type query string true "type" (1：单图，2：多图, 3：base64图片)
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/public/uploadFile [post]
func (e *File) UploadFile(c *gin.Context) {
	tag, _ := c.GetPostForm("type")
	urlPerfix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse
	if tag == "" {
		e.Error(c, 500, nil, "缺少标识")
		//app.Error(c, 200, errors.New(""), "缺少标识")
		return
	} else {
		switch tag {
		case "1": // 单图
			var done bool
			fileResponse, done = e.singleFile(c, fileResponse, urlPerfix)
			if done {
				return
			}
			e.OK(c, fileResponse, "上传成功")
			return
		case "2": // 多图
			multipartFile := e.multipleFile(c, urlPerfix)
			e.OK(c, multipartFile, "上传成功")
			return
		case "3": // base64
			fileResponse = e.baseImg(c, fileResponse, urlPerfix)
			e.OK(c, fileResponse, "上传成功")
		}
	}
}

func (e *File) baseImg(c *gin.Context, fileResponse FileResponse, urlPerfix string) FileResponse {
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
	guid := uuid.New().String()
	fileName := guid + ".jpg"
	err := utils.IsNotExistMkDir(path)
	if err != nil {
		e.Error(c, 500, errors.New(""), "初始化文件路径失败")
	}
	base64File := path + fileName
	_ = ioutil.WriteFile(base64File, ddd, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
	fileResponse = FileResponse{
		Size:     pkg.GetFileSize(base64File),
		Path:     base64File,
		FullPath: urlPerfix + base64File,
		Name:     "",
		Type:     typeStr,
	}
	source, _ := c.GetPostForm("source")
	err = thirdUpload(source, fileName, base64File)
	if err != nil {
		e.Error(c, 200, errors.New(""), "上传第三方失败")
		return fileResponse
	}
	if source != "1" {
		fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
		fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	}
	return fileResponse
}

func (e *File) multipleFile(c *gin.Context, urlPerfix string) []FileResponse {
	files := c.Request.MultipartForm.File["file"]
	source, _ := c.GetPostForm("source")
	var multipartFile []FileResponse
	for _, f := range files {
		guid := uuid.New().String()
		fileName := guid + utils.GetExt(f.Filename)

		err := utils.IsNotExistMkDir(path)
		if err != nil {
			e.Error(c, 500, errors.New(""), "初始化文件路径失败")
		}
		multipartFileName := path + fileName
		err1 := c.SaveUploadedFile(f, multipartFileName)
		fileType, _ := utils.GetType(multipartFileName)
		if err1 == nil {
			err := thirdUpload(source, fileName, multipartFileName)
			if err != nil {
				e.Error(c, 500, errors.New(""), "上传第三方失败")
			} else {
				fileResponse := FileResponse{
					Size:     pkg.GetFileSize(multipartFileName),
					Path:     multipartFileName,
					FullPath: urlPerfix + multipartFileName,
					Name:     f.Filename,
					Type:     fileType,
				}
				if source != "1" {
					fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
					fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
				}
				multipartFile = append(multipartFile, fileResponse)
			}
		}
	}
	return multipartFile
}

func (e *File) singleFile(c *gin.Context, fileResponse FileResponse, urlPerfix string) (FileResponse, bool) {
	files, err := c.FormFile("file")

	if err != nil {
		e.Error(c, 200, errors.New(""), "图片不能为空")
		return FileResponse{}, true
	}
	// 上传文件至指定目录
	guid := uuid.New().String()

	fileName := guid + utils.GetExt(files.Filename)

	err = utils.IsNotExistMkDir(path)
	if err != nil {
		e.Error(c, 500, errors.New(""), "初始化文件路径失败")
	}
	singleFile := path + fileName
	_ = c.SaveUploadedFile(files, singleFile)
	fileType, _ := utils.GetType(singleFile)
	fileResponse = FileResponse{
		Size:     pkg.GetFileSize(singleFile),
		Path:     singleFile,
		FullPath: urlPerfix + singleFile,
		Name:     files.Filename,
		Type:     fileType,
	}
	source, _ := c.GetPostForm("source")
	err = thirdUpload(source, fileName, singleFile)
	if err != nil {
		e.Error(c, 200, errors.New(""), "上传第三方失败")
		return FileResponse{}, true
	}
	fileResponse.Path = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	fileResponse.FullPath = "https://youshikeji.oss-cn-shanghai.aliyuncs.com/img/" + fileName
	return fileResponse, false
}

func thirdUpload(source string, name string, path string) error {
	switch source {
	case "2":
		return ossUpload("img/"+name, path)
	case "3":
		return qiniuUpload("img/"+name, path)
	}
	return nil
}

func ossUpload(name string, path string) error {
	oss := file_store.ALiYunOSS{}
	return oss.UpLoad(name, path)
}

func qiniuUpload(name string, path string) error {
	oss := file_store.ALiYunOSS{}
	return oss.UpLoad(name, path)
}
