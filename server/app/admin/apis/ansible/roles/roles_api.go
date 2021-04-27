package roles

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	"io/ioutil"
	"net/http"
	"strings"
)

type RolesApi interface {
	AddRoleByCompressedData(c *gin.Context)
	DownloadRoleCompressedData(c *gin.Context)
	RemoveRole(c *gin.Context)
}

type rolesApi struct {
	apis.Api
	vo Roles
}

func newRolesApi(vo Roles) RolesApi {
	api := &rolesApi{
		vo: vo,
	}
	return RolesApi(api)
}

func NewRolesApi() RolesApi {
	vo := newRoles()
	return newRolesApi(vo)
}

func (r *rolesApi) AddRoleByCompressedData(c *gin.Context) {
	// method: POST location: /:role
	roleName := c.Param("role")
	if strings.TrimSpace(roleName) == "" {
		r.Error(c, http.StatusNotFound, errors.New("role name is null"), "指定的role为空")
		return
	}
	roleDataFileHeader, err := c.FormFile("compressedData")
	if err != nil {
		r.Error(c, http.StatusBadRequest, err, "错误的role文件")
	}
	roleDataFile, err := roleDataFileHeader.Open()
	if err != nil {
		r.Error(c, http.StatusBadRequest, err, "role文件无法打开")
	}
	defer roleDataFile.Close()
	roleData, err := ioutil.ReadAll(roleDataFile)
	if err != nil {
		r.Error(c, http.StatusBadRequest, err, "读取role文件失败")
	}
	err = r.vo.ImportRoleData(roleName, roleData)
	if err != nil {
		r.Error(c, http.StatusBadRequest, err, "导入role文件失败")
	}
	r.OK(c, nil, "导入role文件成功")
}

func (r *rolesApi) DownloadRoleCompressedData(c *gin.Context) {
	// TODO: 下载role
	panic("implement me")
}

func (r *rolesApi) RemoveRole(c *gin.Context) {
	// TODO: 删除role
	panic("implement me")
}
