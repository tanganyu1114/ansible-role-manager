package inventory

import "testing"

func Test_inventory_RenewGroupName(t *testing.T) {
	group1 := newGroup()
	group2 := newGroup()
	_ = group1.setName("test-group")
	_ = group2.setName("test-group2")
	_ = group1.addHost(NewIPv4Host([4]byte{192, 168, 0, 1}), NewIPv4Host([4]byte{10, 1, 0, 1}))
	_ = group2.addHost(NewIPv4Host([4]byte{192, 168, 1, 1}), NewIPv4Host([4]byte{10, 2, 0, 1}))
	groups := make(map[string]Group)
	groups[group1.GetName()] = group1
	groups[group2.GetName()] = group2
	type fields struct {
		groups map[string]Group
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
			name:   "normal rename",
			fields: fields{groups: groups},
			args: args{
				oldName: "test-group",
				newName: "test-group1",
			},
		},
		{
			name:   "not exist group",
			fields: fields{groups: groups},
			args: args{
				oldName: "test-group1",
				newName: "test-group3",
			},
			wantErr: true,
		},
		{
			name:   "rename to duplicate group",
			fields: fields{groups: groups},
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
				groups: g,
			}
			if err := i.RenewGroupName(tt.args.oldName, tt.args.newName); (err != nil) != tt.wantErr {
				t.Errorf("RenewGroupName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				t.Logf("new group name = %v", i.GetGroups()[tt.args.newName].GetName())
			}
		})
	}
}
