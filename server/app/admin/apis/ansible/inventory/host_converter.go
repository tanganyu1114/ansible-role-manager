package inventory

import (
	"fmt"
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"strconv"
	"strings"
)

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
