package service

import "testing"

func TestConfigurationService_validate(t *testing.T) {
	tests := []struct {
		name    string
		s       *ConfigurationService
		wantErr bool
	}{
		{
			name: "Minimal config (full ip)",
			s: &ConfigurationService{
				ServerAddress: "0.0.0.0:8088",
			},
			wantErr: false,
		},
		{
			name: "Minimal config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
			},
			wantErr: false,
		},
		{
			name: "With Prometheus config (full ip)",
			s: &ConfigurationService{
				ServerAddress:     ":8088",
				PrometheusAddress: "0.0.0.0:8089",
			},
			wantErr: false,
		},
		{
			name: "With Prometheus config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress:     ":8088",
				PrometheusAddress: ":8089",
			},
			wantErr: false,
		},
		{
			name: "With HTTP API config (full ip)",
			s: &ConfigurationService{
				ServerAddress:  ":8088",
				HTTPApiAddress: "0.0.0.0:8089",
			},
			wantErr: false,
		},
		{
			name: "With HTTP API config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress:  ":8088",
				HTTPApiAddress: ":8089",
			},
			wantErr: false,
		},
		{
			name: "With HTTP Webui config (full ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
				WebuiAddress:  "0.0.0.0:8089",
			},
			wantErr: true,
		},
		{
			name: "With HTTP Webui config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
				WebuiAddress:  ":8089",
			},
			wantErr: true,
		},
		{
			name: "With HTTP Webui config (full ip)",
			s: &ConfigurationService{
				ServerAddress:  ":8088",
				HTTPApiAddress: "0.0.0.0:8089",
				WebuiAddress:   "0.0.0.0:8089",
			},
			wantErr: false,
		},
		{
			name: "With HTTP Webui config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress:  ":8088",
				HTTPApiAddress: ":8089",
				WebuiAddress:   ":8089",
			},
			wantErr: false,
		},
		{
			name: "With GRPC config (full ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
				GRPCAddress:   "0.0.0.0:8090",
			},
			wantErr: false,
		},
		{
			name: "With GRPC config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
				GRPCAddress:   ":8090",
			},
			wantErr: false,
		},
		{
			name:    "Server Address not provide",
			s:       &ConfigurationService{},
			wantErr: true,
		},
		{
			name: "With GRPC config (minimal ip)",
			s: &ConfigurationService{
				ServerAddress: ":8088",
				GRPCAddress:   ":8090",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigurationService.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
