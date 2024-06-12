package packet

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPacket_Check(t *testing.T) {
	tests := []struct {
		name    string
		rp      Packet
		wantErr bool
	}{
		{
			name:    "Normal Ping 1",
			rp:      SAMPLE_STATUS_PACKET_1,
			wantErr: false,
		},
		{
			name:    "Normal Ping 2",
			rp:      SAMPLE_STATUS_PACKET_2,
			wantErr: false,
		},
		{
			name: "Max allow packet size",
			rp: Packet{
				PacketHeader: PacketHeader{
					Length: 257,
					ID:     0x00,
				},
				Payload: bytes.NewReader(func() []byte {
					mockupArray := []byte{}
					for i := 0; i < 256; i++ {
						mockupArray = append(mockupArray, 0x01)
					}
					return mockupArray
				}()),
			},
			wantErr: false,
		},
		{
			name: "Over sized",
			rp: Packet{
				PacketHeader: PacketHeader{
					Length: 258,
					ID:     0x00,
				},
				Payload: bytes.NewReader(func() []byte {
					mockupArray := []byte{}
					for i := 0; i < 257; i++ {
						mockupArray = append(mockupArray, 0x01)
					}
					return mockupArray
				}()),
			},
			wantErr: true,
		},
		{
			name: "Sized not defined",
			rp: Packet{
				PacketHeader: PacketHeader{
					Length: 0,
					ID:     0x00,
				},
				Payload: bytes.NewReader([]byte{}),
			},
			wantErr: true,
		},
		{
			name:    "Length Over",
			rp:      SAMPLE_STATUS_PACKET_LENGTH_OVER,
			wantErr: true,
		},
		{
			name: "Empty Data",
			rp: Packet{
				PacketHeader: PacketHeader{
					Length: 2,
					ID:     0,
				},
				Payload: bytes.NewReader(nil),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rp.Check(); (err != nil) != tt.wantErr {
				t.Errorf("Packet.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	type args struct {
		packet Packet
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Normal Ping 1",
			args: args{
				packet: SAMPLE_STATUS_PACKET_1,
			},
			want: SAMPLE_STATUS_DATA_1,
		},
		{
			name: "Normal Ping 2",
			args: args{
				packet: SAMPLE_STATUS_PACKET_2,
			},
			want: SAMPLE_STATUS_DATA_2,
		},
		{
			name: "Nil payload",
			args: args{
				packet: Packet{
					Payload: nil,
				},
			},
			want: EMPTY_DATA,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serialize(tt.args.packet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserialize(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Packet
		want1   RemainingData
		wantErr bool
	}{
		{
			name: "Normal ping 1",
			args: args{
				data: SAMPLE_STATUS_DATA_1,
			},
			want:    SAMPLE_STATUS_PACKET_1,
			want1:   []byte{},
			wantErr: false,
		},
		{
			name: "Normal ping 2",
			args: args{
				data: SAMPLE_STATUS_DATA_2,
			},
			want:    SAMPLE_STATUS_PACKET_2,
			want1:   []byte{},
			wantErr: false,
		},
		{
			name: "Over sized",
			args: args{
				data: func() []byte {
					res := []byte{0x81, 0x01}
					for i := 0; i < 258; i++ {
						res = append(res, 0x01)
					}
					return res
				}(),
			},
			want1:   []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Deserialize(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantData := tt.args.data

			// actualData := make([]byte, 32)
			// n, _ := got.Payload.Read(actualData)
			actualData := Serialize(got)
			if !bytes.Equal(wantData, actualData) && (len(wantData) <= 256) {
				t.Errorf("Deserialize() got = %v, want %v", actualData, wantData)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Deserialize() got = %v, want %v", got, tt.want)
			// }
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Deserialize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
