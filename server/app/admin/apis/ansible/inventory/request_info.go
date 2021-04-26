package inventory

type NewGroupRequestInfo struct {
	GroupName string `json:"groupName"`
	Hosts     []Host `json:"ipAddrs"`
}

type ModifyGroupRequestInfo struct {
	GroupName       string `json:"groupName"`
	TargetGroupName string `json:"targetGroupName"`
	Hosts           []Host `json:"ipAddrs"`
}
