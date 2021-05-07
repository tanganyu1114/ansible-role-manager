package inventory

import (
	"fmt"
	"sort"
	"strings"
)

type Inventory interface {
	getAllGroups() map[string]Group
	getTruncatedGroup() map[string]bool
	generateGroupAll() Group

	AddHostToGroup(groupName string, hosts ...Host) error
	RenewGroupName(oldName, newName string) error
	RemoveHostFromGroup(groupName string, hosts ...Host)
	RemoveGroup(groupName string)
	GetGroups(limit, page int) (totalGroupNum, totalPagesNum int, groups map[string]Group) // DONE: 分页查询机制、反馈主机总数
}

type inventory struct {
	sortedGroupNames []string
	groups           map[string]Group
	isTruncatedGroup map[string]bool
}

func newInventory(groups map[string]Group) Inventory {
	inv := &inventory{
		sortedGroupNames: make([]string, 0),
		groups:           groups,
		isTruncatedGroup: make(map[string]bool),
	}
	for s := range groups {
		inv.sortedGroupNames = append(inv.sortedGroupNames, s)
	}
	inv.sortGroups()
	return Inventory(inv)
}

func (i *inventory) AddHostToGroup(groupName string, hosts ...Host) error {
	// 不允许新增"all"组名的组
	if strings.EqualFold(strings.ToLower(groupName), "all") {
		return fmt.Errorf("can not build a group which named '%s'", groupName)
	}
	var g Group
	if _, has := i.groups[groupName]; has {
		g = i.groups[groupName]
	} else {
		g = newGroup()
		err := g.setName(groupName)
		if err != nil {
			return err
		}
	}
	err := g.addHost(hosts...)
	if err != nil {
		return err
	}
	i.groups[groupName] = g
	i.isTruncatedGroup[groupName] = false
	i.sortedGroupNames = append(i.sortedGroupNames, groupName)
	i.sortGroups()
	return nil
}

func (i *inventory) RenewGroupName(oldName, newName string) error {
	// 不允许新增"all"组名的组
	if strings.EqualFold(strings.ToLower(newName), "all") {
		return fmt.Errorf("rename '%s' is not allowed", newName)
	}
	if _, has := i.groups[oldName]; !has {
		return fmt.Errorf("nonexistent group by name %s", oldName)
	}
	if _, has := i.groups[newName]; has {
		return fmt.Errorf("duplicate group name %s", newName)
	}
	g := i.groups[oldName]
	i.RemoveGroup(oldName)
	return i.AddHostToGroup(newName, g.GetHosts()...)
}

func (i *inventory) RemoveHostFromGroup(groupName string, hosts ...Host) {
	if _, has := i.groups[groupName]; has {
		i.groups[groupName].removeHost(hosts...)
	}
}

func (i *inventory) RemoveGroup(groupName string) {
	if _, has := i.groups[groupName]; has {
		delete(i.groups, groupName)
		idx := i.searchGroup(groupName)
		i.sortedGroupNames = append(i.sortedGroupNames[:idx], i.sortedGroupNames[idx+1:]...)
		i.isTruncatedGroup[groupName] = true
	}
}

func (i inventory) GetGroups(limit, page int) (totalGroupsNum, totalPagesNum int, groups map[string]Group) {
	totalGroupsNum = len(i.sortedGroupNames)
	totalPagesNum = -1
	groups = make(map[string]Group)
	if limit < 0 || page <= 0 {
		return
	}

	// 加入all组进行页数计算
	totalPagesNum = (totalGroupsNum+1)/limit + 1

	// 计算起始索引
	startIdx := (page - 1) * limit
	// 因加入了all组进行计算，起始索引需往前移动一位
	startIdx--
	if startIdx < -1 {
		fmt.Println("invalid groups index", startIdx)
		return
	}

	endIdx := totalGroupsNum - 1
	if totalPagesNum < page {
		return
	} else if totalPagesNum > page {
		// 计算结束索引
		endIdx = page*limit - 1
		// 因加入了all组进行计算，结束索引需往前移动一位
		endIdx--
	}

	for ; startIdx <= endIdx; startIdx++ {
		if startIdx == -1 {
			groups["all"] = i.generateGroupAll()
			continue
		}
		groupName := i.sortedGroupNames[startIdx]
		groups[groupName] = i.groups[groupName]
	}

	return
}

func (i inventory) getAllGroups() map[string]Group {
	return i.groups
}

func (i *inventory) getTruncatedGroup() map[string]bool {
	truncatedGroup := make(map[string]bool)
	for groupName := range i.isTruncatedGroup {
		if i.isTruncatedGroup[groupName] {
			truncatedGroup[groupName] = true
		}
	}
	return truncatedGroup
}

func (i inventory) generateGroupAll() Group {
	groupAll := newGroup()
	_ = groupAll.setName("all")
	count := 0
	for _, g := range i.groups {
		_ = groupAll.addHost(g.GetHosts()...)
		count += g.HostsLen()
		// todo: handle error
		//err := groupAll.addHost(g.GetHosts()...)
		//if err != nil {
		//	fmt.Printf("get hosts from group %s failed, cased by: %s\n", g.GetName(), err)
		//}
	}
	return groupAll
}

func (i *inventory) sortGroups() {
	sort.SliceIsSorted(i.sortedGroupNames, func(x, y int) bool {
		return isLessString(i.sortedGroupNames[x], i.sortedGroupNames[y])
	})
}

func (i *inventory) searchGroup(groupName string) int {
	return sort.Search(len(i.sortedGroupNames), func(idx int) bool {
		// ! idxName < groupName
		return !isLessString(i.sortedGroupNames[idx], groupName)
	})
}

func isLessString(x, y string) bool {
	xLen := len(x)
	yLen := len(y)
	var minLen int
	if xLen < yLen {
		minLen = xLen
	} else {
		minLen = yLen
	}

	for j := 0; j < minLen; j++ {
		if x[j] == y[j] {
			continue
		}
		return x[j] < y[j]
	}
	return xLen == minLen && xLen != yLen
}
