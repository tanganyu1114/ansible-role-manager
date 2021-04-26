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

func NewGroup(groupName string, hosts []Host) (Group, error) {
	g := newGroup()
	err := g.setName(groupName)
	if err != nil {
		return nil, err
	}
	err = g.addHost(hosts...)
	if err != nil {
		return nil, err
	}
	return g, nil
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
		if h == nil {
			continue
		}
		idx := sort.Search(g.HostsLen(), func(i int) bool {
			if !g.hosts[i].GetIp().IP.Equal(h.GetIp().IP) {
				iIPv4 := g.hosts[i].GetIp().IP.To4()
				hIPv4 := h.GetIp().IP.To4()
				return !isLessIPv4(iIPv4, hIPv4)
			}
			return true
		})
		if idx < g.HostsLen() && g.hosts[idx].GetIp().IP.Equal(h.GetIp().IP) {
			continue
		}
		isIn = false
		g.hosts = append(g.hosts, h)
		sort.Slice(g.hosts, func(x, y int) bool {
			xIPv4 := g.hosts[x].GetIp().IP.To4()
			yIPv4 := g.hosts[y].GetIp().IP.To4()
			return isLessIPv4(xIPv4, yIPv4)
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
		if h == nil {
			continue
		}
		idx := sort.Search(g.HostsLen(), func(i int) bool {
			if !g.hosts[i].GetIp().IP.Equal(h.GetIp().IP) {
				iIPv4 := g.hosts[i].GetIp().IP.To4()
				hIPv4 := h.GetIp().IP.To4()
				return !isLessIPv4(iIPv4, hIPv4)
			}
			return true
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

func isLessIPv4(xIP, yIP net.IP) bool {
	for j := 0; j < net.IPv4len; j++ {
		if xIP[j] == yIP[j] {
			continue
		}
		if xIP[j] < yIP[j] {
			return true
		}
		return false
	}
	return false
}
