package inventory

import (
	"testing"
)

func testGroupExample() map[string]Group {
	group1 := newGroup()
	group2 := newGroup()
	_ = group1.setName("test-group")
	_ = group2.setName("test-group2")
	_ = group1.addHost(NewIPv4Host([4]byte{192, 168, 0, 1}), NewIPv4Host([4]byte{10, 1, 0, 1}))
	_ = group2.addHost(NewIPv4Host([4]byte{192, 168, 1, 1}), NewIPv4Host([4]byte{10, 2, 0, 1}))
	groups := make(map[string]Group)
	groups[group1.GetName()] = group1
	groups[group2.GetName()] = group2
	return groups
}

func Test_inventory_AddHostToGroup(t *testing.T) {
	type fields struct {
		groups           map[string]Group
		isTruncatedGroup map[string]bool
	}
	type args struct {
		groupName string
		hosts     []Host
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		addFailed bool
	}{
		{
			name: "add host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{NewIPv4Host([4]byte{192, 168, 1, 2})},
			},
		},
		{
			name: "add hosts to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group1",
				hosts:     []Host{NewIPv4Host([4]byte{192, 168, 2, 2}), NewIPv4Host([4]byte{10, 2, 0, 2})},
			},
		},
		{
			name: "add host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{NewIPv4Host([4]byte{192, 168, 3, 1})},
			},
		},
		{
			name: "add hosts to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group4",
				hosts:     []Host{NewIPv4Host([4]byte{192, 168, 4, 1}), NewIPv4Host([4]byte{10, 4, 0, 1})},
			},
		},
		{
			name: "add nil host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{nil},
			},
			addFailed: true,
		},
		{
			name: "add nil hosts to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     nil,
			},
			addFailed: true,
		},
		{
			name: "add null host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{},
			},
			addFailed: true,
		},
		{
			name: "add nil host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{nil},
			},
			addFailed: true,
		},
		{
			name: "add nil hosts to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     nil,
			},
			addFailed: true,
		},
		{
			name: "add null host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{},
			},
			addFailed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				groups:           tt.fields.groups,
				isTruncatedGroup: tt.fields.isTruncatedGroup,
			}
			var bLen, aLen int
			if _, has := i.GetGroups()[tt.args.groupName]; has {
				bLen = i.GetGroups()[tt.args.groupName].HostsLen()
			}
			i.AddHostToGroup(tt.args.groupName, tt.args.hosts...)
			if _, has := i.GetGroups()[tt.args.groupName]; has {
				aLen = i.GetGroups()[tt.args.groupName].HostsLen()
			}
			if (bLen >= aLen) != tt.addFailed {
				t.Errorf("AddHostToGroup() add failed, before lenghth = %d, after lenghth = %d", bLen, aLen)
			}
		})
	}
}

func Test_inventory_RenewGroupName(t *testing.T) {
	type fields struct {
		groups           map[string]Group
		isTruncatedGroup map[string]bool
	}
	type args struct {
		oldName string
		newName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal rename",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group",
				newName: "test-group1",
			},
		},
		{
			name: "not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group1",
				newName: "test-group3",
			},
			wantErr: true,
		},
		{
			name: "rename to duplicate group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group",
				newName: "test-group2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := make(map[string]Group)
			for s, g2 := range tt.fields.groups {
				g[s] = newGroup()
				err := g[s].setName(g2.GetName())
				if err != nil {
					t.Fatal(err)
				}
				_ = g[s].addHost(g2.GetHosts()...)
			}

			i := &inventory{
				groups:           g,
				isTruncatedGroup: make(map[string]bool),
			}
			if err := i.RenewGroupName(tt.args.oldName, tt.args.newName); (err != nil) != tt.wantErr {
				t.Errorf("RenewGroupName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				t.Logf("old group name = %v, want new group name = %v, got %v", tt.args.oldName, tt.args.newName, i.GetGroups()[tt.args.newName].GetName())
			}
		})
	}
}
