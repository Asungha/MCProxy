package packet

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHandshake_GetHostname(t *testing.T) {
	tests := []struct {
		name string
		h    *Handshake
		want string
	}{
		{
			name: "Vanilla handshake (IP host)",
			h: &Handshake{
				hostname: "192.168.1.123",
			},
			want: "192.168.1.123",
		},
		{
			name: "Vanilla handshake (Domain host)",
			h: &Handshake{
				hostname: "test.com",
			},
			want: "test.com",
		},
		{
			name: "FML handshake (IP host)",
			h: &Handshake{
				hostname: fmt.Sprintf("192.168.1.123%sFML%s", string([]byte{0x00}), string([]byte{0x00})),
			},
			want: "192.168.1.123",
		},
		{
			name: "FML handshake (Domain host)",
			h: &Handshake{
				hostname: fmt.Sprintf("test.com%sFML%s", string([]byte{0x00}), string([]byte{0x00})),
			},
			want: "test.com",
		},
		{
			name: "FML handshake (malformed host)",
			h: &Handshake{
				hostname: fmt.Sprintf("test.com%sFML", string([]byte{0x00})),
			},
			want: "test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.GetHostname(); got != tt.want {
				t.Errorf("Handshake.GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandshake_encode(t *testing.T) {
	tests := []struct {
		name string
		h    Handshake
		want []byte
	}{
		{
			name: "Normal Ping 1",
			h:    SAMPLE_STATUS_PAYLOAD_PACKET_1,
			want: SAMPLE_STATUS_DATA_1[2:],
		},
		{
			name: "Normal Ping 2",
			h:    SAMPLE_STATUS_PAYLOAD_PACKET_2,
			want: SAMPLE_STATUS_DATA_2[2:],
		},
		{
			name: "Length over",
			h:    SAMPLE_STATUS_PAYLOAD_PACKET_LENGTH_OVER,
			want: SAMPLE_STATUS_DATA_LENGTH_OVER[2:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handshake.encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandshake_Encode(t *testing.T) {
	tests := []struct {
		name    string
		h       Handshake
		want    []byte
		wantErr bool
	}{
		{
			name:    "Normal Ping 1",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_1,
			want:    SAMPLE_STATUS_DATA_1,
			wantErr: false,
		},
		{
			name:    "Normal Ping 2",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_2,
			want:    SAMPLE_STATUS_DATA_2,
			wantErr: false,
		},
		{
			name:    "Length over",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_LENGTH_OVER,
			want:    []byte{},
			wantErr: true,
		},
		{
			name:    "tailling data",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET__TAILLING,
			want:    SAMPLE_STATUS_DATA_2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Handshake.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handshake.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandshake_Decode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		h       Handshake
		args    args
		wantErr bool
	}{
		{
			name:    "Normal Ping 1",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_1,
			args:    args{data: SAMPLE_STATUS_DATA_1},
			wantErr: false,
		},
		{
			name:    "Normal Ping 2",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_2,
			args:    args{data: SAMPLE_STATUS_DATA_2},
			wantErr: false,
		},
		{
			name:    "Empty data",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_2,
			args:    args{data: EMPTY_DATA},
			wantErr: true,
		},
		{
			name:    "Non relevant data",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET_2,
			args:    args{data: NON_RELEVANT_DATA},
			wantErr: true,
		},
		{
			name:    "tailling data",
			h:       SAMPLE_STATUS_PAYLOAD_PACKET__TAILLING,
			args:    args{data: SAMPLE_STATUS_DATA_TAILLING},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.Decode(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Handshake.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandshake_String(t *testing.T) {
	tests := []struct {
		name string
		h    Handshake
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.String(); got != tt.want {
				t.Errorf("Handshake.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandshake_Length(t *testing.T) {
	tests := []struct {
		name string
		h    Handshake
		want int
	}{
		{
			name: "Normal Ping 1",
			h:    SAMPLE_STATUS_PAYLOAD_PACKET_1,
			want: int(SAMPLE_STATUS_PACKET_1.Length - 2),
		},
		{
			name: "Normal Ping 2",
			h:    SAMPLE_STATUS_PAYLOAD_PACKET_2,
			want: int(SAMPLE_STATUS_PACKET_2.Length - 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Length(); got != tt.want {
				t.Errorf("Handshake.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandshake_Destroy(t *testing.T) {
	tests := []struct {
		name string
		h    *Handshake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Destroy()
		})
	}
}

func TestNewHandshake(t *testing.T) {
	tests := []struct {
		name string
		want *Handshake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandshake(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandshake() = %v, want %v", got, tt.want)
			}
		})
	}
}
