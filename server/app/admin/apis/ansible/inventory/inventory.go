package inventory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"net/http"
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

func NewInventoryApi(inventoryDir string) (Inventory, error) {
	storage := svc.NewInventoryFileStorage(inventoryDir)
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
	groupName, isExist := c.GetPostForm("group_name")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
		return
	}
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
	i.OK(c, nil, "成功添加主机信息")
}

func (i *inventory) RenewGroupName(c *gin.Context) {
	oldGroupName, isExist := c.GetPostForm("old_group_name")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("old group name is null"), "请求的表单不存在old_group_name")
		return
	}
	newGroupName, isExist := c.GetPostForm("new_group_name")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("new group name is null"), "请求的表单不存在new_group_name")
		return
	}

	err := i.bo.RenewGroupName(oldGroupName, newGroupName)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "请求的表单参数不规范")
		return
	}
	i.OK(c, nil, "成功更改组名")
}

func (i *inventory) RemoveHostFromGroup(c *gin.Context) {
	groupName, isExist := c.GetPostForm("group_name")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
		return
	}
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
	i.OK(c, nil, "成功删除主机信息")
}

func (i *inventory) RemoveGroupByName(c *gin.Context) {
	groupName, isExist := c.GetPostForm("group_name")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("group name is null"), "请求的表单不存在group_name")
		return
	}

	i.bo.RemoveGroup(groupName)
	i.OK(c, nil, "完成操作")
}

func (i *inventory) GetAllHosts(c *gin.Context) {
	hConverter := newHostVOConverter()
	hostsVO := make([]Host, 0)
	for _, hostBO := range i.bo.GetAllHosts() {
		hostVO := hConverter.ConvertToVO(hostBO)
		hostsVO = append(hostsVO, hostVO)
	}
	i.OK(c, hostsVO, "成功查询所有主机信息")
}

func (i *inventory) GetGroups(c *gin.Context) {
	gConverter := newGroupVOConverter()
	groupsVO := make(map[string]Group)
	for groupName, groupBO := range i.bo.GetGroups() {
		groupVO := gConverter.ConvertToVO(groupBO)
		groupsVO[groupName] = groupVO
	}
	i.OK(c, groupsVO, "成功查询所有组信息")
}
