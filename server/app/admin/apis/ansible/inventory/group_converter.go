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
	hostsBO := make([]svc.Host, 0)
	hConverter := newHostVOConverter()
	for _, hostVO := range vo.Hosts {
		hostBO, err := hConverter.ConvertToBO(hostVO)
		if err != nil {
			return nil, err
		}
		hostsBO = append(hostsBO, hostBO)
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
		hostVO := Host(hostBO.GetIp().IP.String())
		groupVO.Hosts = append(groupVO.Hosts, hostVO)
	}
	return groupVO
}
