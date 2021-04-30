package inventory

type Group struct {
	GroupName string `json:"groupName"`
	Hosts     []Host `json:"ipAddrs"`
	HostsLen  int    `json:"hostsLen"`
}

type Host string

type Groups struct {
	GroupsMap map[string]Group `json:"groupsMap"`
	Hosts     []Host           `json:"ipAddrs"`
	HostsLen  int              `json:"hostsLen"`
}
