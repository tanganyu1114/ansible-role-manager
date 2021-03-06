package inventory

import "testing"

func Test_group_addHost(t *testing.T) {
	type fields struct {
		groupName string
		hosts     []Host
	}
	type args struct {
		hosts []Host
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal add host",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{192, 168, 0, 1})}},
		},
		{
			name: "normal add hosts",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{192, 168, 0, 1}), NewIPv4Host([4]byte{10, 1, 0, 1})}},
		},
		{
			name: "add duplicate host",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args:    args{hosts: []Host{NewIPv4Host([4]byte{192, 168, 0, 1})}},
			wantErr: true,
		},
		{
			name: "add duplicate hosts",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args:    args{hosts: []Host{NewIPv4Host([4]byte{192, 168, 0, 1}), NewIPv4Host([4]byte{10, 1, 0, 1})}},
			wantErr: true,
		},
		{
			name: "add null host",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args:    args{hosts: []Host{}},
			wantErr: true,
		},
		{
			name: "add nil host",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args:    args{hosts: []Host{nil}},
			wantErr: true,
		},
		{
			name: "add nil hosts",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args:    args{hosts: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &group{
				groupName: tt.fields.groupName,
				hosts:     tt.fields.hosts,
			}
			if err := g.addHost(tt.args.hosts...); (err != nil) != tt.wantErr {
				t.Errorf("addHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_group_removeHost(t *testing.T) {
	type fields struct {
		groupName string
		hosts     []Host
	}
	type args struct {
		hosts []Host
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "normal remove host",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1})}},
		},
		{
			name: "normal remove hosts",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})}},
		},
		{
			name: "remove not exist host",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{192, 168, 1, 1})}},
		},
		{
			name: "remove null host",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: []Host{}},
		},
		{
			name: "remove nil host",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: []Host{nil}},
		},
		{
			name: "remove nil hosts",
			fields: fields{
				hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1}), NewIPv4Host([4]byte{192, 168, 0, 1})},
			},
			args: args{hosts: nil},
		},
		{
			name: "empty group remove host",
			fields: fields{
				hosts: make([]Host, 0),
			},
			args: args{hosts: []Host{NewIPv4Host([4]byte{10, 1, 0, 1})}},
		},
	}
	defer func() {
		err := recover()
		if err != nil {
			t.Error(err)
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &group{
				groupName: tt.fields.groupName,
				hosts:     tt.fields.hosts,
			}
			g.removeHost(tt.args.hosts...)
		})
	}
}
