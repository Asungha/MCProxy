package utils

import (
	"reflect"
	"testing"
)

func TestSplitDataframe(t *testing.T) {
	type args struct {
		buffer []byte
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				buffer: []byte{0x01, 0x00},
			},
			want:    [][]byte{{0x01, 0x00}},
			wantErr: false,
		},
		{
			name: "Normal case (double)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00},
			},
			want:    [][]byte{{0x01, 0x00}, {0x02, 0x01, 0x00}},
			wantErr: false,
		},
		{
			name: "Normal case (tripple)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00, 0x03, 0x02, 0x01, 0x00},
			},
			want:    [][]byte{{0x01, 0x00}, {0x02, 0x01, 0x00}, {0x03, 0x02, 0x01, 0x00}},
			wantErr: false,
		},
		{
			name: "Empty data",
			args: args{
				buffer: []byte{},
			},
			want:    [][]byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitDataframe(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitDataframe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitDataframe() = %v, want %v", got, tt.want)
			}
		})
	}
}
