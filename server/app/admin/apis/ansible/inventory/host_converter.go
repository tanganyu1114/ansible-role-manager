package inventory

import (
	"fmt"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"strconv"
	"strings"
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
	if sp := strings.Split(string(vo), "."); len(sp) != 4 {
		return nil, fmt.Errorf("%s is not a ip address", vo)
	} else {
		ip := [4]byte{}
		for k := 0; k < 4; k++ {
			b, err := strconv.Atoi(sp[k])
			if err != nil || b > 255 {
				return nil, fmt.Errorf("%s is not a valid ip address", vo)
			}
			ip[k] = byte(b)
		}
		bo := svc.NewIPv4Host(ip)
		return bo, nil
	}
}

func (h hostConverter) ConvertToVO(bo svc.Host) Host {
	vo := Host(bo.GetIp().IP.String())
	return vo
}
