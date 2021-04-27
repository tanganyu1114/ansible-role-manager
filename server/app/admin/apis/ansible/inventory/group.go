package inventory

type Group struct {
	GroupName string	`json:"groupName"`
	Hosts     []Host	`json:"ipAddrs"`
}
