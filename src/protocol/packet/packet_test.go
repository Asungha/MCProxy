package packet

import (
	"bytes"
	"io"
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
		{
			name: "Old protocol",
			args: args{
				packet: Get_LEGACY_STATUS_REQ_PACKET(),
			},
			want: LEGACY_STATUS_REQ,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Serialize(tt.args.packet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserialize(t *testing.T) {
	type args struct {
		data        []byte
		isHandshake bool
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
				data:        SAMPLE_STATUS_DATA_1,
				isHandshake: true,
			},
			want:    SAMPLE_STATUS_PACKET_1,
			want1:   []byte{},
			wantErr: false,
		},
		{
			name: "Normal ping 2",
			args: args{
				data:        SAMPLE_STATUS_DATA_2,
				isHandshake: true,
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
				isHandshake: true,
			},
			want1:   []byte{},
			wantErr: true,
		},
		{
			name: "Legacy status",
			args: args{
				data:        LEGACY_STATUS_REQ,
				isHandshake: true,
			},
			want:    Get_LEGACY_STATUS_REQ_PACKET(),
			want1:   []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, got1, err := Deserialize(tt.args.data, tt.args.isHandshake)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantData := tt.args.data

			// actualData := make([]byte, 32)
			// n, _ := got.Payload.Read(actualData)
			actualData, err := Serialize(got)
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

func TestPacketRound(t *testing.T) {
	tests := []struct {
		name        string
		input       Packet
		isHandshake bool
		payload     []byte
		wantErr     bool
	}{
		{
			name:        "Normal 1",
			input:       Packet{PacketHeader: SAMPLE_STATUS_PACKET_1.PacketHeader, Payload: bytes.NewReader(SAMPLE_STATUS_DATA_1[2:])},
			payload:     SAMPLE_STATUS_DATA_1[2:],
			isHandshake: true,
			wantErr:     false,
		},
		{
			name:        "Normal 2",
			input:       Packet{PacketHeader: SAMPLE_STATUS_PACKET_2.PacketHeader, Payload: bytes.NewReader(SAMPLE_STATUS_DATA_2[2:])},
			payload:     SAMPLE_STATUS_DATA_2[2:],
			isHandshake: true,
			wantErr:     false,
		},
		{
			name:        "Old protocol 1",
			input:       Packet{PacketHeader: PacketHeader{Length: 1, ID: 0x00, IsOldProtocol: true}, Payload: bytes.NewReader(LEGACY_STATUS_REQ)},
			payload:     LEGACY_STATUS_REQ,
			isHandshake: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Serialize(tt.input)
			tt.input.Payload.Reset(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%x", data)
			pac, _, remaining, err := Deserialize(data, tt.isHandshake)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(remaining, []byte{}) {
				t.Errorf("Deserialize() Remaining data not expected: %v", remaining)
				return
			}
			// if !reflect.DeepEqual(pac, tt.input) {
			// 	t.Errorf("Deserialize() got = %v, want %v", pac, tt.input)
			// }
			if pac.ID != tt.input.ID || pac.Length != tt.input.Length || pac.IsOldProtocol != tt.input.IsOldProtocol {
				t.Errorf("Deserialize() got = %v, want %v", pac, tt.input)
			}
			data1 := make([]byte, pac.Payload.Size())
			data2 := make([]byte, tt.input.Payload.Size())

			_, err = pac.Payload.Read(data1)
			if (err != nil && err != io.EOF) != tt.wantErr {
				t.Logf(">> %v", len(data1))
				t.Errorf("Deserialize() error1 = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err = tt.input.Payload.Read(data2)
			if (err != nil && err != io.EOF) != tt.wantErr {
				t.Logf(">> %v", len(data2))
				t.Errorf("Deserialize() error2 = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(data1, data2) {
				t.Errorf("Deserialize() got = %v, want %v", data1, data2)
				return
			}
		})
	}
}
