package inventory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"net/http"
	"strings"
)

type Inventory interface {
	AddHostToGroup(c *gin.Context)
	RenewGroupName(c *gin.Context)
	RemoveHostFromGroup(c *gin.Context)
	RemoveGroupByName(c *gin.Context)
	GetAllHosts(c *gin.Context)
	GetGroups(c *gin.Context)
}

type inventory struct {
	apis.Api
	bo svc.Inventory
}

func NewInventoryApi() (Inventory, error) {
	storage := svc.GetSingletonInventoryFileStorageIns()
	invBO, err := storage.Load()
	if err != nil {
		return nil, err
	}
	invVO := &inventory{
		bo: invBO,
	}
	return Inventory(invVO), nil
}

func (i *inventory) AddHostToGroup(c *gin.Context) {
	// method:POST location: /groups/:group
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}
	//groupName, isExist := c.GetPostForm("group_name")
	//if !isExist {
	//	i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
	//	return
	//}
	hostsStr, isExist := c.GetPostFormArray("hosts")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is null"), "请求的表单不存在hosts")
		return
	}
	if hostsStr == nil || len(hostsStr) == 0 {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is empty"), "hosts列表为空")
		return
	}

	hConverter := newHostVOConverter()
	hostsBO := make([]svc.Host, len(hostsStr))
	for j := 0; j < len(hostsStr); j++ {
		vo := Host(hostsStr[j])
		hostBO, err := hConverter.ConvertToBO(vo)
		if err != nil {
			i.Error(c, http.StatusBadRequest, err, "host解析失败")
			return
		}
		hostsBO[j] = hostBO
	}

	i.bo.AddHostToGroup(groupName, hostsBO...)
	err := i.save()
	if err != nil {
		i.Error(c, http.StatusInternalServerError, err, "保存配置失败")
		return
	}
	i.OK(c, nil, "成功添加主机信息")
}

func (i *inventory) RenewGroupName(c *gin.Context) {
	// method: PATCH location: /groups/:group
	oldGroupName := c.Param("group")
	if strings.TrimSpace(oldGroupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("old group name is null"), "指定的group为空")
		return
	}
	//oldGroupName, isExist := c.GetPostForm("old_group_name")
	//if !isExist {
	//	i.Error(c, http.StatusBadRequest, errors.New("old group name is null"), "请求的表单不存在old_group_name")
	//	return
	//}
	newGroupName, isExist := c.GetQuery("new_group_name")
	if !isExist || strings.TrimSpace(newGroupName) == "" {
		i.Error(c, http.StatusBadRequest, errors.New("new group name is null"), "请求的参数new_group_name为空")
		return
	}
	//newGroupName := c.Param("new_group_name")
	//if strings.TrimSpace(newGroupName) == "" {
	//	i.Error(c, http.StatusBadRequest, errors.New("new group name is null"), "请求的表单不存在new_group_name")
	//	return
	//}
	//newGroupName, isExist := c.GetPostForm("new_group_name")
	//if !isExist {
	//	i.Error(c, http.StatusBadRequest, errors.New("new group name is null"), "请求的表单不存在new_group_name")
	//	return
	//}

	err := i.bo.RenewGroupName(oldGroupName, newGroupName)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "请求的表单参数不规范")
		return
	}
	err = i.save()
	if err != nil {
		i.Error(c, http.StatusInternalServerError, err, "保存配置失败")
		return
	}
	i.OK(c, nil, "成功更改组名")
}

func (i *inventory) RemoveHostFromGroup(c *gin.Context) {
	// method: POST location: /groups/remove/:group/hosts
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}
	//groupName, isExist := c.GetPostForm("group_name")
	//if !isExist {
	//	i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
	//	return
	//}
	hostsStr, isExist := c.GetPostFormArray("hosts")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is null"), "请求的表单不存在hosts")
		return
	}
	if hostsStr == nil || len(hostsStr) == 0 {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is empty"), "hosts列表为空")
		return
	}

	hConverter := newHostVOConverter()
	hostsBO := make([]svc.Host, len(hostsStr))
	for j := 0; j < len(hostsStr); j++ {
		vo := Host(hostsStr[j])
		hostBO, err := hConverter.ConvertToBO(vo)
		if err != nil {
			i.Error(c, http.StatusBadRequest, err, "host解析失败")
			return
		}
		hostsBO[j] = hostBO
	}

	i.bo.RemoveHostFromGroup(groupName, hostsBO...)
	err := i.save()
	if err != nil {
		i.Error(c, http.StatusInternalServerError, err, "保存配置失败")
		return
	}
	i.OK(c, nil, "成功删除主机信息")
}

func (i *inventory) RemoveGroupByName(c *gin.Context) {
	// method: DELETE location: /groups/:group
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}
	//groupName, isExist := c.GetPostForm("group_name")
	//if !isExist {
	//	i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
	//	return
	//}

	i.bo.RemoveGroup(groupName)
	err := i.save()
	if err != nil {
		i.Error(c, http.StatusInternalServerError, err, "保存配置失败")
		return
	}
	i.OK(c, nil, "完成操作")
}

func (i *inventory) GetAllHosts(c *gin.Context) {
	// method: GET location: /hosts
	hConverter := newHostVOConverter()
	hostsVO := make([]Host, 0)
	for _, hostBO := range i.bo.GetAllHosts() {
		hostVO := hConverter.ConvertToVO(hostBO)
		hostsVO = append(hostsVO, hostVO)
	}
	i.OK(c, hostsVO, "成功查询所有主机信息")
}

func (i *inventory) GetGroups(c *gin.Context) {
	// method: GET location: /groups
	gConverter := newGroupVOConverter()
	groupsVO := make(map[string]Group)
	for groupName, groupBO := range i.bo.GetGroups() {
		groupVO := gConverter.ConvertToVO(groupBO)
		groupsVO[groupName] = groupVO
	}
	i.OK(c, groupsVO, "成功查询所有组信息")
}

func (i inventory) save() error {
	storageBO := svc.GetSingletonInventoryFileStorageIns()
	return storageBO.Save(i.bo)
}
