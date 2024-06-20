package utils

import (
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	type args struct {
		slice [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Normal case 1",
			args: args{
				slice: [][]byte{{0x00, 0x01, 0x02}},
			},
			want: []byte{0x00, 0x01, 0x02},
		},
		{
			name: "Normal case 2",
			args: args{
				slice: [][]byte{{0x00, 0x01, 0x02}, {0x03, 0x04}},
			},
			want: []byte{0x00, 0x01, 0x02, 0x03, 0x04},
		},
		{
			name: "Normal case 3",
			args: args{
				slice: [][]byte{{0x00, 0x01, 0x02}, {0x03, 0x04}, {0x05}},
			},
			want: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		},
		{
			name: "Empty data case 1",
			args: args{
				slice: [][]byte{{}, {0x03, 0x04}, {0x05}},
			},
			want: []byte{0x03, 0x04, 0x05},
		},
		{
			name: "Empty data case 1",
			args: args{
				slice: [][]byte{{0x01, 0x02}, {}, {0x05}},
			},
			want: []byte{0x01, 0x02, 0x05},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Concat(tt.args.slice...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concat() = %v, want %v", got, tt.want)
			}
		})
	}
}
