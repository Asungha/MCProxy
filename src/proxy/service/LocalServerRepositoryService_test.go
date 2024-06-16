package service

import (
	"reflect"
	"testing"
)

func TestLocalServerRepositoryService_Resolve(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		s       *LocalServerRepositoryService
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal 1",
			s: &LocalServerRepositoryService{
				servers: map[string]map[string]string{
					"hostname.a": {
						"target": "target.a",
					},
					"hostname.b": {
						"target": "target.b",
					},
				},
			},
			args: args{
				hostname: "hostname.a",
			},
			want:    "target.a",
			wantErr: false,
		},
		{
			name: "Normal 2",
			s: &LocalServerRepositoryService{
				servers: map[string]map[string]string{
					"hostname.a": {
						"target": "target.a",
					},
					"hostname.b": {
						"target": "target.b",
					},
				},
			},
			args: args{
				hostname: "hostname.b",
			},
			want:    "target.b",
			wantErr: false,
		},
		{
			name: "Hostname not available",
			s: &LocalServerRepositoryService{
				servers: map[string]map[string]string{
					"hostname.a": {
						"target": "target.a",
					},
					"hostname.b": {
						"target": "target.b",
					},
				},
			},
			args: args{
				hostname: "hostname.c",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Hostname not load",
			s: &LocalServerRepositoryService{
				servers: map[string]map[string]string{},
			},
			args: args{
				hostname: "hostname.a",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Resolve(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalServerRepositoryService.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LocalServerRepositoryService.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalServerRepositoryService_List(t *testing.T) {
	tests := []struct {
		name    string
		s       *LocalServerRepositoryService
		want    []ServerList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalServerRepositoryService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalServerRepositoryService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLocalServerRepoService(t *testing.T) {
	tests := []struct {
		name string
		want ServerRepositoryService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalServerRepoService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalServerRepoService() = %v, want %v", got, tt.want)
			}
		})
	}
}
