package inventory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"net/http"
	"strings"
	"sync"
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
	sync.RWMutex
	//bo svc.Inventory
}

func NewInventoryApi() Inventory {
	invVO := &inventory{
		RWMutex: sync.RWMutex{},
	}
	return Inventory(invVO)
}

func (i *inventory) AddHostToGroup(c *gin.Context) {
	// method:POST location: /groups/:group
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}
	hostsStr, isExist := c.GetPostFormArray("hosts")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is null"), "请求的表单不存在hosts")
		return
	}

	hostsC := newHostsVOConverter()
	hostsBO, err := hostsC.ConvertToBOFromString(hostsStr)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "请求的表单格式不正确")
	}

	done := i.boDO(c, true, func(invBO svc.Inventory) (fnDone bool) {
		invBO.AddHostToGroup(groupName, hostsBO...)
		return true
	})

	if done {
		i.OK(c, nil, "成功添加主机信息")
	}
}

func (i *inventory) RenewGroupName(c *gin.Context) {
	// method: PATCH location: /groups/:group
	oldGroupName := c.Param("group")
	if strings.TrimSpace(oldGroupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("old group name is null"), "指定的group为空")
		return
	}

	newGroupName, isExist := c.GetQuery("new_group_name")
	if !isExist || strings.TrimSpace(newGroupName) == "" {
		i.Error(c, http.StatusBadRequest, errors.New("new group name is null"), "请求的参数new_group_name为空")
		return
	}

	done := i.boDO(c, true, func(invBO svc.Inventory) (fnDone bool) {
		err := invBO.RenewGroupName(oldGroupName, newGroupName)
		if err != nil {
			i.Error(c, http.StatusBadRequest, err, "参数不规范")
			return false
		}
		return true
	})

	if done {
		i.OK(c, nil, "成功更改组名")
	}
}

func (i *inventory) RemoveHostFromGroup(c *gin.Context) {
	// method: POST location: /groups/remove/:group/hosts
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}

	hostsStr, isExist := c.GetPostFormArray("hosts")
	if !isExist {
		i.Error(c, http.StatusBadRequest, errors.New("hosts is null"), "请求的表单不存在hosts")
		return
	}

	hostsC := newHostsVOConverter()
	hostsBO, err := hostsC.ConvertToBOFromString(hostsStr)
	if err != nil {
		i.Error(c, http.StatusBadRequest, err, "请求的表单格式不正确")
	}

	done := i.boDO(c, true, func(invBO svc.Inventory) (fnDone bool) {
		invBO.RemoveHostFromGroup(groupName, hostsBO...)
		return true
	})

	if done {
		i.OK(c, nil, "成功删除主机信息")
	}
}

func (i *inventory) RemoveGroupByName(c *gin.Context) {
	// method: DELETE location: /groups/:group
	groupName := c.Param("group")
	if strings.TrimSpace(groupName) == "" {
		i.Error(c, http.StatusNotFound, errors.New("group name is null"), "指定的group为空")
		return
	}

	done := i.boDO(c, true, func(invBO svc.Inventory) (fnDone bool) {
		invBO.RemoveGroup(groupName)
		return true
	})

	if done {
		i.OK(c, nil, "完成操作")
	}
}

func (i *inventory) GetAllHosts(c *gin.Context) {
	// method: GET location: /hosts
	hostsVO := make([]Host, 0)
	done := i.boDO(c, false, func(invBO svc.Inventory) (fnDone bool) {
		hostsC := newHostsVOConverter()
		hostsVO = hostsC.ConvertToVO(invBO.GetAllHosts())
		return true
	})

	if done {
		i.OK(c, hostsVO, "成功查询所有主机信息")
	}
}

func (i *inventory) GetGroups(c *gin.Context) {
	// method: GET location: /groups
	groupsVO := make(map[string]Group)
	done := i.boDO(c, false, func(invBO svc.Inventory) (fnDone bool) {
		gConverter := newGroupVOConverter()
		for groupName, groupBO := range invBO.GetGroups() {
			groupVO := gConverter.ConvertToVO(groupBO)
			groupsVO[groupName] = groupVO
		}
		return true
	})

	if done {
		i.OK(c, groupsVO, "成功查询所有组信息")
	}
}

func (i *inventory) boDO(c *gin.Context, needSave bool, doFn func(invBO svc.Inventory) (fnDone bool)) bool {
	if needSave {
		i.Lock()
		defer i.Unlock()
	} else {
		i.RLock()
		defer i.RUnlock()
	}

	storageBO := svc.GetSingletonInventoryStorageIns()
	invBO, err := storageBO.Load()
	if err != nil {
		i.Error(c, http.StatusInternalServerError, err, "读取配置失败")
		return false
	}

	isDone := doFn(invBO)
	if !isDone {
		return false
	}

	if needSave {
		err = storageBO.Save(invBO)
		if err != nil {
			i.Error(c, http.StatusInternalServerError, err, "保存配置失败")
			return false
		}
		return true
	}

	return true
}
