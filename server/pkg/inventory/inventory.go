package inventory

import (
	"fmt"
)

type Inventory interface {
	AddHostToGroup(groupName string, hosts ...Host)
	RenewGroupName(oldName, newName string) error
	RemoveHostFromGroup(groupName string, hosts ...Host)
	RemoveGroup(groupName string)
	GetAllHosts() []Host
	GetGroups() map[string]Group
}

type inventory struct {
	groups map[string]Group
}

func newInventory(groups map[string]Group) Inventory {
	inv := &inventory{groups: groups}
	return Inventory(inv)
}

func (i *inventory) AddHostToGroup(groupName string, hosts ...Host) {
	var g Group
	if _, has := i.groups[groupName]; has {
		g = i.groups[groupName]
	} else {
		g = newGroup()
		err := g.setName(groupName)
		if err != nil {
			// todo: log
			return
		}
	}
	err := g.addHost(hosts...)
	if err != nil {
		// todo: log
		return
	}
	i.groups[groupName] = g
}

func (i *inventory) RenewGroupName(oldName, newName string) error {
	if _, has := i.groups[oldName]; !has {
		return fmt.Errorf("nonexistent group by name %s", oldName)
	}
	if _, has := i.groups[newName]; has {
		return fmt.Errorf("duplicate group name %s", newName)
	}
	g := i.groups[oldName]
	//err := g.setName(newName)
	//if err != nil {
	//	return err
	//}
	i.RemoveGroup(oldName)
	i.AddHostToGroup(newName, g.GetHosts()...)
	return nil
}

func (i *inventory) RemoveHostFromGroup(groupName string, hosts ...Host) {
	if _, has := i.groups[groupName]; has {
		i.groups[groupName].removeHost(hosts...)
	}
}

func (i *inventory) RemoveGroup(groupName string) {
	if _, has := i.groups[groupName]; has {
		delete(i.groups, groupName)
	}
}

func (i inventory) GetAllHosts() []Host {
	groupAll := newGroup()
	for _, g := range i.groups {
		_ = groupAll.addHost(g.GetHosts()...)
		// todo: handle error
		//err := groupAll.addHost(g.GetHosts()...)
		//if err != nil {
		//	fmt.Printf("get hosts from group %s failed, cased by: %s\n", g.GetName(), err)
		//}
	}
	return groupAll.GetHosts()
}

func (i inventory) GetGroups() map[string]Group {
	return i.groups
}
