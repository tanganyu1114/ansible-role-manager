package roles

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_tgzArchiver_Decompress(t *testing.T) {
	exDir := "../../test/pkg/roles/decompress"
	tgzFileName := "test.tgz"

	tgzData, err := os.ReadFile(filepath.Join(exDir, tgzFileName))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		exDir          string
		compressedData []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal decompress",
			args: args{
				exDir:          exDir,
				compressedData: tgzData,
			},
		},
		{
			name: "unknown compressed data",
			args: args{
				exDir:          exDir,
				compressedData: []byte{1, 1, 1, 1, 1, 0, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tgzArchiver{}
			if err := g.Decompress(tt.args.exDir, tt.args.compressedData); (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
