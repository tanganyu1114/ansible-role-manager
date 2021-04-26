package inventory

type Group struct {
	GroupName string	`form:"groupName" json:"groupName"`
	Hosts     []Host	`form:"ipAddrs" json:"ipAddrs"`
}
