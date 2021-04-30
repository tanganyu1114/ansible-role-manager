package inventory

import (
	"errors"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"sync"
)

type Inventory interface {
	AddHostToGroup(groupName string, hosts ...Host) error
	RenewGroupName(oldName, newName string) error
	RemoveHostFromGroup(groupName string, hosts ...Host) error
	RemoveGroup(groupName string) error
	GetGroups() Groups
	//GetAllHosts() []Host
	//GetGroups() map[string]Group
}

type inventory struct {
	sync.RWMutex
}

func newInventory() Inventory {
	invVO := &inventory{RWMutex: sync.RWMutex{}}
	return Inventory(invVO)
}

func (i *inventory) AddHostToGroup(groupName string, hosts ...Host) error {
	if hosts == nil || len(hosts) == 0 {
		return errors.New("hosts are null")
	}

	groupVO := Group{
		GroupName: groupName,
		Hosts:     hosts,
	}
	groupC := newGroupVOConverter()
	groupBO, err := groupC.ConvertToBO(groupVO)
	if err != nil {
		return err
	}

	err = i.boDO(true, func(invBO svc.Inventory) error {
		invBO.AddHostToGroup(groupBO.GetName(), groupBO.GetHosts()...)
		return nil
	})

	return err
}

func (i *inventory) RenewGroupName(oldName, newName string) error {
	err := i.boDO(true, func(invBO svc.Inventory) error {
		return invBO.RenewGroupName(oldName, newName)
	})
	return err
}

func (i *inventory) RemoveHostFromGroup(groupName string, hosts ...Host) error {
	if hosts == nil || len(hosts) == 0 {
		return errors.New("hosts are null")
	}

	groupVO := Group{
		GroupName: groupName,
		Hosts:     hosts,
	}
	groupC := newGroupVOConverter()
	groupBO, err := groupC.ConvertToBO(groupVO)
	if err != nil {
		return err
	}

	err = i.boDO(true, func(invBO svc.Inventory) error {
		invBO.RemoveHostFromGroup(groupBO.GetName(), groupBO.GetHosts()...)
		return nil
	})

	return err
}

func (i *inventory) RemoveGroup(groupName string) error {
	err := i.boDO(true, func(invBO svc.Inventory) error {
		invBO.RemoveGroup(groupName)
		return nil
	})
	return err
}

func (i *inventory) GetGroups() Groups {
	groupsVO := &Groups{}
	hostsC := newHostsVOConverter()
	groupC := newGroupVOConverter()
	_ = i.boDO(false, func(invBO svc.Inventory) error {
		hostsBO, l := invBO.GetAllHosts()
		groupsVO.Hosts, groupsVO.HostsLen = hostsC.ConvertToVO(hostsBO), l
		groupsBO := invBO.GetGroups()
		for s, groupBO := range groupsBO {
			groupsVO.GroupsMap[s] = groupC.ConvertToVO(groupBO)
		}
		return nil
	})
	return *groupsVO
}

func (i *inventory) boDO(needSave bool, doFn func(invBO svc.Inventory) error) error {
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
		return err
	}

	err = doFn(invBO)
	if err != nil {
		return err
	}

	if needSave {
		err = storageBO.Save(invBO)
		if err != nil {
			return err
		}
	}

	return nil
}
