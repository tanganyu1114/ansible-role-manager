package inventory

import svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"

type GroupVOConverter interface {
	ConvertToBO(vo Group) (svc.Group, error)
	ConvertToVO(bo svc.Group) Group
}

type groupConverter struct {
}

func newGroupVOConverter() GroupVOConverter {
	converter := new(groupConverter)
	return GroupVOConverter(converter)
}

func (g groupConverter) ConvertToBO(vo Group) (svc.Group, error) {
	groupName := vo.GroupName

	hostsC := newHostsVOConverter()
	hostsBO, err := hostsC.ConvertToBO(vo.Hosts)
	if err != nil {
		return nil, err
	}

	groupBO, err := svc.NewGroup(groupName, hostsBO)
	if err != nil {
		return nil, err
	}

	return groupBO, nil
}

func (g groupConverter) ConvertToVO(bo svc.Group) Group {
	groupName := bo.GetName()
	groupVO := Group{
		GroupName: groupName,
		Hosts:     make([]Host, 0),
	}
	for _, hostBO := range bo.GetHosts() {
		hostVO := Host(hostBO.GetIPString())
		groupVO.Hosts = append(groupVO.Hosts, hostVO)
	}
	return groupVO
}
