package roles

import (
	"errors"
	"fmt"
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
	roleDataFileHeader, err := c.FormFile("file")
	if err != nil {
		r.Error(c, http.StatusBadRequest, err, "错误的role文件")
		return
	}
	roleDataFile, err := roleDataFileHeader.Open()
	if err != nil {
		r.Error(c, http.StatusUnsupportedMediaType, err, "role文件无法打开")
		return
	}
	defer roleDataFile.Close()
	roleData, err := ioutil.ReadAll(roleDataFile)
	if err != nil {
		r.Error(c, http.StatusUnsupportedMediaType, err, "读取role文件失败")
		return
	}
	err = r.vo.ImportRoleData(roleName, roleData)
	if err != nil {
		r.Error(c, http.StatusUnsupportedMediaType, err, "导入role文件失败")
		return
	}
	r.OK(c, nil, "导入role文件成功")
}

func (r *rolesApi) DownloadRoleCompressedData(c *gin.Context) {
	// DONE: 下载role
	// method: GET location: /:role

	roleName := c.Param("role")
	if strings.TrimSpace(roleName) == "" {
		r.Error(c, http.StatusNotFound, errors.New("role name is null"), "指定的role为空")
		return
	}

	roleData, err := r.vo.ExportRoleData(roleName)
	if err != nil {
		r.Error(c, http.StatusForbidden, err, "生成role压缩文件失败")
		return
	}
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", roleName))
	c.Header("Content-Transfer-Encoding", "binary")
	_, err = c.Writer.Write(roleData)
	if err != nil {
		r.Error(c, http.StatusNotImplemented, err, "导出role压缩文件失败")
	}
	//r.OK(c, nil, "导出role文件成功")
}

func (r *rolesApi) RemoveRole(c *gin.Context) {
	// DONE: 删除role
	// method: GET location: /:role

	roleName := c.Param("role")
	if strings.TrimSpace(roleName) == "" {
		r.Error(c, http.StatusNotFound, errors.New("role name is null"), "指定的role为空")
		return
	}

	err := r.vo.RemoveRole(roleName)
	if err != nil {
		r.Error(c, http.StatusNotFound, err, "删除role文件失败")
		return
	}
	r.OK(c, nil, "删除role文件成功")
}
