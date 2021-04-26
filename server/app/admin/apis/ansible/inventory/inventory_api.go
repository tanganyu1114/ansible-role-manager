package inventory

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	"net/http"
	"strings"
)

type InventoryApi interface {
	AddHostToGroup(c *gin.Context)
	RenewGroupName(c *gin.Context)
	RemoveHostFromGroup(c *gin.Context)
	RemoveGroupByName(c *gin.Context)
	GetAllHosts(c *gin.Context)
	GetGroups(c *gin.Context)
}

type inventoryApi struct {
	apis.Api
	vo Inventory
}

func NewInventoryApi() InventoryApi {
	vo := newInventory()
	return newInventoryApi(vo)
}

func newInventoryApi(vo Inventory) InventoryApi {
	api := &inventoryApi{
		vo: vo,
	}
	return InventoryApi(api)
}

func (i *inventoryApi) AddHostToGroup(c *gin.Context) {
	// method:POST location: /groups

	groupVO := Group{}
	err := c.ShouldBindJSON(&groupVO)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的文本格式")
		return
	}

	err = i.vo.AddHostToGroup(groupVO.GroupName, groupVO.Hosts...)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的Group格式")
		return
	}

	i.OK(c, nil, "成功添加主机信息")
}

func (i *inventoryApi) RenewGroupName(c *gin.Context) {
	// method: PUT location: /groups

	params := struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}{}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的文本格式")
		return
	}

	err = i.vo.RenewGroupName(params.OldName, params.NewName)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的参数格式")
		return
	}

	i.OK(c, nil, "成功更改组名")
}

func (i *inventoryApi) RemoveHostFromGroup(c *gin.Context) {
	// method: PATCH location: /groups
	groupVO := Group{}
	err := c.ShouldBindJSON(&groupVO)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的文本格式")
		return
	}

	err = i.vo.RemoveHostFromGroup(groupVO.GroupName, groupVO.Hosts...)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的Group格式")
		return
	}

	i.OK(c, nil, "成功删除主机信息")
}

func (i *inventoryApi) RemoveGroupByName(c *gin.Context) {
	// method: DELETE location: /groups/:group
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}

	err := i.vo.RemoveGroup(groupName)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, fmt.Sprintf("删除[%s]组失败", groupName))
		return
	}

	i.OK(c, nil, "完成操作")
}

func (i *inventoryApi) GetAllHosts(c *gin.Context) {
	// method: GET location: /hosts
	hostsVO := i.vo.GetAllHosts()
	i.OK(c, hostsVO, "成功查询所有主机信息")
}

func (i *inventoryApi) GetGroups(c *gin.Context) {
	// method: GET location: /groups
	groupsVO := i.vo.GetGroups()
	i.OK(c, groupsVO, "成功查询所有组信息")
}
