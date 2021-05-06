package inventory

type Group struct {
	GroupName string `json:"groupName"`
	Hosts     []Host `json:"ipAddrs"`
	HostsLen  int    `json:"hostsLen"`
}

type Host string

type Groups struct {
	TotalGroupsNum int              `json:"totalGroupsNum"`
	TotalPagesNum  int              `json:"totalPagesNum"`
	GroupsMap      map[string]Group `json:"groupsMap"`
}
