package inventory

import (
	"path/filepath"
	"reflect"
	"testing"
)

//type mockInvFileParser struct {
//}
//
//func (m mockInvFileParser) Parse(data []byte) (Group, error) {
//	panic("implement me")
//}
//
//func (m mockInvFileParser) Dump(g Group) ([]byte, error) {
//	panic("implement me")
//}

func Test_inventoryFileStorage_Load(t *testing.T) {
	dir := "../../test/pkg/inventory/inventory_storage/inventory"
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
		dir    string
		parser InventoryFileParser
	}
	tests := []struct {
		name    string
		fields  fields
		want    Inventory
		wantErr bool
	}{
		{
			name: "normal load",
			fields: fields{
				dir:    dir,
				parser: NewInventoryFileParser(),
			},
			want: newInventory(groups),
		},
		{
			name: "error dir path",
			fields: fields{
				dir:    "testdir",
				parser: NewInventoryFileParser(),
			},
			wantErr: true,
		},
		{
			name: "null inventory file dir",
			fields: fields{
				dir:    filepath.Dir(dir),
				parser: NewInventoryFileParser(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inventoryFileStorage{
				dir:    tt.fields.dir,
				parser: tt.fields.parser,
			}
			got, err := i.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Load() got = %v, want %v", got, tt.want)
			//}
			if got == nil {
				if tt.want == nil {
					return
				}
				t.Errorf("Load() got = %v, want %v", got, tt.want)
				return
			}

			gotGroups := got.GetGroups()
			wantGroups := tt.want.GetGroups()
			if len(gotGroups) != len(wantGroups) {
				t.Errorf("Load() gotGroups = %v, wantGroups %v", gotGroups, wantGroups)
				return
			}

			for gotGroupName, gotGroup := range gotGroups {
				if _, has := wantGroups[gotGroupName]; !has {
					t.Errorf("Load() gotGroups = %v, wantGroups %v", gotGroups, wantGroups)
				}
				if !reflect.DeepEqual(gotGroup.GetHosts(), wantGroups[gotGroupName].GetHosts()) {
					t.Errorf("Load() gotHosts = %v, wantHosts %v", gotGroup.GetHosts(), wantGroups[gotGroupName].GetHosts())
				}
			}

		})
	}
}

func Test_inventoryFileStorage_Save(t *testing.T) {
	dir := "../../test/pkg/inventory/inventory_storage/inventory"
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
		dir    string
		parser InventoryFileParser
	}
	type args struct {
		inv Inventory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal save",
			fields: fields{
				dir:    dir,
				parser: NewInventoryFileParser(),
			},
			args: args{inv: newInventory(groups)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inventoryFileStorage{
				dir:    tt.fields.dir,
				parser: tt.fields.parser,
			}
			if err := i.Save(tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
