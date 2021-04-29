package inventory

import (
	"fmt"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
)

type HostsVOConverter interface {
	//ConvertToBOFromString(hosts []string) (bo []svc.Host, err error)
	ConvertToBO(vo []Host) (bo []svc.Host, err error)
	ConvertToVO(bo []svc.Host) []Host
}

type hostsConverter struct {
}

func newHostsVOConverter() HostsVOConverter {
	converter := new(hostsConverter)
	return HostsVOConverter(converter)
}

//func (h hostsConverter) ConvertToBOFromString(hosts []string) (bo []svc.Host, err error) {
//	if hosts == nil || len(hosts) == 0 {
//		return nil, errors.New("hosts is empty")
//	}
//
//	hConverter := newHostVOConverter()
//	hostsBO := make([]svc.Host, len(hosts))
//	for j := 0; j < len(hosts); j++ {
//		vo := Host(hosts[j])
//		hostBO, err := hConverter.ConvertToBO(vo)
//		if err != nil {
//			return nil, err
//		}
//		hostsBO[j] = hostBO
//	}
//	return hostsBO, nil
//}

func (h hostsConverter) ConvertToBO(vo []Host) (bo []svc.Host, err error) {
	bo = make([]svc.Host, 0)
	hostC := newHostVOConverter()
	for _, hostVO := range vo {
		hostBO, err := hostC.ConvertToBO(hostVO)
		if err != nil {
			return nil, err
		}
		bo = append(bo, hostBO)
	}
	return bo, nil
}

func (h hostsConverter) ConvertToVO(bo []svc.Host) []Host {
	vo := make([]Host, 0)
	hostC := newHostVOConverter()
	for _, host := range bo {
		vo = append(vo, hostC.ConvertToVO(host))
	}
	return vo
}

type HostVOConverter interface {
	ConvertToBO(vo Host) (bo svc.Host, err error)
	ConvertToVO(bo svc.Host) Host
}

type hostConverter struct {
}

func newHostVOConverter() HostVOConverter {
	converter := new(hostConverter)
	return HostVOConverter(converter)
}

func (h hostConverter) ConvertToBO(vo Host) (svc.Host, error) {
	//ip := net.ParseIP(string(vo))
	//if ip == nil {
	//	return nil, fmt.Errorf("'%s' is not a valid ip address", vo)
	//}
	//return svc.NewHost(ip, "")
	bo := svc.ParseHost(string(vo))
	if bo == nil {
		return nil, fmt.Errorf("'%s' is not a valid ip address", vo)
	}
	return bo, nil
}

func (h hostConverter) ConvertToVO(bo svc.Host) Host {
	vo := Host(bo.GetIPString())
	return vo
}
