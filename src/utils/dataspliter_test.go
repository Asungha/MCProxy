package utils

import (
	// packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
	. "mc_reverse_proxy/src/common"
	"reflect"
	"testing"
)

func TestSplitDataframe(t *testing.T) {
	type args struct {
		buffer []byte
		strict bool
	}
	tests := []struct {
		name    string
		args    args
		want    []PacketFragment
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				buffer: []byte{0x01, 0x00},
				strict: false,
			},
			want:    []PacketFragment{{Data: []byte{0x01, 0x00}, Type: MC_OTHER}},
			wantErr: false,
		},
		{
			name: "Normal case (double)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00},
				strict: false,
			},
			want: []PacketFragment{
				{Data: []byte{0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x02, 0x01, 0x00}, Type: MC_OTHER},
			},
			wantErr: false,
		},
		{
			name: "Normal case (tripple)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00, 0x03, 0x02, 0x01, 0x00},
				strict: false,
			},
			want: []PacketFragment{
				{Data: []byte{0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x02, 0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x03, 0x02, 0x01, 0x00}, Type: MC_OTHER},
			},
			wantErr: false,
		},
		{
			name: "Empty data",
			args: args{
				buffer: []byte{},
				strict: false,
			},
			want:    []PacketFragment{},
			wantErr: true,
		},
		{
			name: "Normal case (strict validate)",
			args: args{
				buffer: []byte{0x01, 0x00},
				strict: true,
			},
			want:    []PacketFragment{{Data: []byte{0x01, 0x00}, Type: MC_OTHER}},
			wantErr: false,
		},
		{
			name: "Normal case (double, strict validate)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00},
				strict: true,
			},
			want: []PacketFragment{
				{Data: []byte{0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x02, 0x01, 0x00}, Type: MC_OTHER},
			},
			wantErr: false,
		},
		{
			name: "Normal case (tripple, strict validate)",
			args: args{
				buffer: []byte{0x01, 0x00, 0x02, 0x01, 0x00, 0x03, 0x02, 0x01, 0x00},
				strict: true,
			},
			want: []PacketFragment{
				{Data: []byte{0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x02, 0x01, 0x00}, Type: MC_OTHER},
				{Data: []byte{0x03, 0x02, 0x01, 0x00}, Type: MC_OTHER},
			},
			wantErr: false,
		},
		{
			name: "Empty data (strict validate)",
			args: args{
				buffer: []byte{},
				strict: true,
			},
			want:    []PacketFragment{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitDataframe(tt.args.buffer, true)
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
