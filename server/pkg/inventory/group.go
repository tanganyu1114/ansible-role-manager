package inventory

import (
	"errors"
	"net"
	"sort"
	"strings"
)

type Group interface {
	addHost(hosts ...Host) error
	removeHost(hosts ...Host)
	setName(newName string) error
	GetName() string
	GetHosts() []Host
	HostsLen() int
}

type group struct {
	groupName string
	hosts     []Host
}

func newGroup() Group {
	g := &group{hosts: make([]Host, 0)}
	return Group(g)
}

func (g *group) addHost(hosts ...Host) error {
	if hosts == nil || len(hosts) == 0 {
		return errors.New("input hosts are null")
	}
	isIn := true
	for _, h := range hosts {
		idx := sort.Search(g.HostsLen(), func(i int) bool {
			return g.hosts[i].GetIp().IP.Equal(h.GetIp().IP)
		})
		if idx < g.HostsLen() && g.hosts[idx].GetIp().IP.Equal(h.GetIp().IP) {
			continue
		}
		isIn = false
		g.hosts = append(g.hosts, h)
		sort.Slice(g.hosts, func(x, y int) bool {
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
	if isIn {
		return errors.New("input hosts are in the group")
	}
	return nil
}

func (g *group) removeHost(hosts ...Host) {
	if hosts == nil || len(hosts) == 0 {
		return
	}
	for _, h := range hosts {
		idx := sort.Search(g.HostsLen(), func(i int) bool {
			return g.hosts[i].GetIp().IP.Equal(h.GetIp().IP)
		})
		if idx < g.HostsLen() && g.hosts[idx].GetIp().IP.Equal(h.GetIp().IP) {
			g.hosts = append(g.hosts[:idx], g.hosts[idx+1:]...)
		}
	}
}

func (g *group) setName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("can not renew group name to null")
	}
	g.groupName = newName
	return nil
}

func (g group) GetName() string {
	return g.groupName
}

func (g group) GetHosts() []Host {
	return g.hosts
}

func (g group) HostsLen() int {
	return len(g.hosts)
}
