package inventory

import (
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"reflect"
	"testing"
)

func Test_groupConverter_ConvertToBO(t *testing.T) {
	wantGroupBO, err := svc.NewGroup("test-group", []svc.Host{svc.ParseHost("192.168.0.1"), svc.ParseHost("10.1.[0:128].[129:254]")})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		vo Group
	}
	tests := []struct {
		name    string
		args    args
		want    svc.Group
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{vo: Group{
				GroupName: "test-group",
				Hosts:     []Host{Host("192.168.0.1"), Host("10.1.[0:128].[129:254]")},
			}},
			want: wantGroupBO,
		},
		{
			name: "no group name",
			args: args{vo: Group{
				Hosts: []Host{Host("192.168.0.1")},
			}},
			wantErr: true,
		},
		{
			name: "error host string",
			args: args{vo: Group{
				GroupName: "test-group",
				Hosts:     []Host{Host("192.168.0.1"), Host("10.1.[0..128].{129-254}")},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := groupConverter{}
			got, err := g.ConvertToBO(tt.args.vo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToBO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Logf("ConvertToBO() error: %s", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToBO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupConverter_ConvertToVO(t *testing.T) {
	groupBO, err := svc.NewGroup("test-group", []svc.Host{svc.ParseHost("192.168.0.1"), svc.ParseHost("10.1.[0:100].[1:254]")})
	if err != nil {
		t.Fatal(err)
	}
	wantGroupVO := Group{
		GroupName: "test-group",
		Hosts:     []Host{"10.1.[0:100].[1:254]", "192.168.0.1"},
		HostsLen:  1 + 254*101,
	}
	type args struct {
		bo svc.Group
	}
	tests := []struct {
		name string
		args args
		want Group
	}{
		{
			name: "normal test",
			args: args{bo: groupBO},
			want: wantGroupVO,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := groupConverter{}
			if got := g.ConvertToVO(tt.args.bo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToVO() = %v, want %v", got, tt.want)
			}
		})
	}
}
