package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gopkg.in/validator.v2"
)

type ConfigurationService struct {
	ServerAddress        string `json:"listen_address" validate:"regexp=^((\d{1,3}\.){3}\d{1,3}|):\d{2,5}$,required"`
	PrometheusAddress    string `json:"prometheus_address" validate:"regexp=^((\d{1,3}\.){3}\d{1,3}|):\d{2,5}$"`
	WebuiAddress         string `json:"webui_address" validate:"regexp=^((\d{1,3}\.){3}\d{1,3}|):\d{2,5}$"`
	HTTPApiAddress       string `json:"http_api_address" validate:"regexp=^((\d{1,3}\.){3}\d{1,3}|):\d{2,5}$"`
	HTTPHostname         string `json:"http_hostname"`
	GRPCAddress          string `json:"grpc_metric_address" validate:"regexp=^((\d{1,3}\.){3}\d{1,3}|):\d{2,5}$"`
	LoggerMongoDBAddress string `json:"logger_mongodb_address" validate:"regexp=^mongodb:\/\/(?:(?:(?:[\w\d%]+)(?::(?:[\w\d%]+))?@)?(?:[\w\d.-]+)(?::\d+)?(?:,(?:(?:[\w\d.-]+)(?::\d+)?))*)\/?(?:[\w\d-]*)?(?:\?[\w\d&=_-]+)?$"`
	LoggerMongoDBName    string `json:"logger_mongodb_db_name"`
	LoggerMongoColName   string `json:"logger_mongodb_collection_name"`
}

func (s *ConfigurationService) ReadConfig(path string) error {
	config_file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer config_file.Close()

	decoder := json.NewDecoder(config_file)
	err = decoder.Decode(s)
	if err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	return s.validate()
}

func (s *ConfigurationService) validate() error {
	if err := validator.Validate(*s); err != nil {
		return err
	}

	if s.ServerAddress == "" {
		return errors.New("server address no provided")
	}

	if s.WebuiAddress != "" && s.HTTPApiAddress == "" {
		return errors.New("webui frontend required http api to work. Add config 'http_api_address' in the config.json with appropiated address")
	}

	if s.HTTPHostname == "" && s.HTTPApiAddress != "" {
		s.HTTPHostname = s.HTTPApiAddress
	}
	if s.LoggerMongoDBAddress != "" {
		if s.LoggerMongoDBName == "" {
			s.LoggerMongoDBName = "packer_logger"
		}
		if s.LoggerMongoColName == "" {
			s.LoggerMongoColName = "packer_logger"
		}
	}
	return nil
}

func NewConfigurationService(configPath string) (*ConfigurationService, error) {
	c := &ConfigurationService{}
	err := c.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return c, nil
}
