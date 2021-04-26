package inventory

import (
	svc "github.com/tanganyu1114/ansible-role-manager/pkg/inventory"
	"reflect"
	"testing"
)

func Test_hostsConverter_ConvertToBO(t *testing.T) {
	type args struct {
		vo []Host
	}
	tests := []struct {
		name    string
		args    args
		wantBo  []svc.Host
		wantErr bool
	}{
		{
			name:   "normal test",
			args:   args{vo: []Host{"192.168.0.1", "10.1.0.1"}},
			wantBo: []svc.Host{svc.NewIPv4Host([4]byte{192, 168, 0, 1}), svc.NewIPv4Host([4]byte{10, 1, 0, 1})},
		},
		{
			name:   "nil hosts",
			args:   args{vo: nil},
			wantBo: []svc.Host{},
		},
		{
			name:   "null hosts",
			args:   args{vo: []Host{}},
			wantBo: []svc.Host{},
		},
		{
			name:    "error ipaddr",
			args:    args{vo: []Host{"192.1688.0.1"}},
			wantErr: true,
		},
		{
			name:    "none ipaddr string",
			args:    args{vo: []Host{"none-ipaddr"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hostsConverter{}
			gotBo, err := h.ConvertToBO(tt.args.vo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToBO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBo, tt.wantBo) {
				t.Errorf("ConvertToBO() gotBo = %v, want %v", gotBo, tt.wantBo)
			}
		})
	}
}

func Test_hostsConverter_ConvertToVO(t *testing.T) {
	type args struct {
		bo []svc.Host
	}
	tests := []struct {
		name string
		args args
		want []Host
	}{
		{
			name: "normal test",
			args: args{bo: []svc.Host{svc.NewIPv4Host([4]byte{192, 168, 0, 1}), svc.NewIPv4Host([4]byte{10, 1, 0, 1})}},
			want: []Host{"192.168.0.1", "10.1.0.1"},
		},
		{
			name: "nil hosts",
			args: args{bo: nil},
			want: []Host{},
		},
		{
			name: "null hosts",
			args: args{bo: []svc.Host{}},
			want: []Host{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hostsConverter{}
			if got := h.ConvertToVO(tt.args.bo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToVO() = %v, want %v", got, tt.want)
			}
		})
	}
}
