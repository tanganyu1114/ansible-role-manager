package inventory

import (
	"fmt"
	"net"
	"sort"
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

// TODO: build Inventory Object

func (i *inventory) AddHostToGroup(groupName string, hosts ...Host) {
	var g Group
	if _, has := i.groups[groupName]; has {
		g = i.groups[groupName]
	} else {
		g = newGroup()
	}
	err := g.addHost(hosts...)
	if err != nil {
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
	return i.groups[oldName].setName(newName)
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
	hosts := make([]Host, 0)
	for _, g := range i.groups {
		for _, h := range g.GetHosts() {
			idx := sort.Search(len(hosts), func(j int) bool {
				return hosts[j].GetIp().IP.Equal(h.GetIp().IP)
			})
			if idx < len(hosts) && hosts[idx].GetIp().IP.Equal(h.GetIp().IP) {
				continue
			}
			hosts = append(hosts, h)
			sort.Slice(hosts, func(x, y int) bool {
				xIPv4 := hosts[x].GetIp().IP.To4()
				yIPv4 := hosts[y].GetIp().IP.To4()
				for j := 0; j < net.IPv4len; j++ {
					if xIPv4[j] < yIPv4[j] {
						return true
					}
				}
				return false
			})
		}
	}
	return hosts
}

func (i inventory) GetGroups() map[string]Group {
	return i.groups
}
