package inventory

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	"net/http"
	"strconv"
	"strings"
)

type InventoryApi interface {
	NewGroup(c *gin.Context)
	ModifyGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
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

func (i *inventoryApi) NewGroup(c *gin.Context) {
	// method:POST location: /groups

	req := NewGroupRequestInfo{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的文本格式")
		return
	}

	err = i.vo.AddHostToGroup(req.GroupName, req.Hosts...)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的Group格式")
		return
	}

	i.OK(c, nil, "成功添加主机信息")
}

func (i *inventoryApi) ModifyGroup(c *gin.Context) {
	// method: PATCH location: /groups
	req := ModifyGroupRequestInfo{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "错误的文本格式")
		return
	}
	var modifyErr error
	defer func() {
		if modifyErr != nil {
			i.Error(c, http.StatusBadRequest, err, "传参异常，修改失败")
		}
	}()
	modifyErr = i.vo.RemoveGroup(req.TargetGroupName)
	if modifyErr != nil {
		return
	}

	modifyErr = i.vo.AddHostToGroup(req.GroupName, req.Hosts...)
	if modifyErr != nil {
		return
	}

	i.OK(c, nil, "成功完成修改")
}

func (i *inventoryApi) DeleteGroup(c *gin.Context) {
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

	i.OK(c, nil, "完成删除操作")
}

func (i *inventoryApi) GetGroups(c *gin.Context) {
	// method: GET location: /groups
	limitStr, ok := c.GetQuery("limit")
	if !ok {
		limitStr = "10"
	}
	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = "1"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		i.Error(c, http.StatusBadRequest, fmt.Errorf("param 'limit' is invalid: %s", err), "传参limit的值不正确")
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		i.Error(c, http.StatusBadRequest, fmt.Errorf("param 'page' is invalid: %s", err), "传参page的值不正确")
	}

	groupsVO := i.vo.GetGroups(limit, page)
	i.OK(c, groupsVO, "成功查询所有组信息")
}
