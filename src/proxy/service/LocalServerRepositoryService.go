package service

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type LocalServerRepositoryService struct {
	ListableRepositoryService

	servers map[string]map[string]string
}

func (s *LocalServerRepositoryService) Load() error {
	// s.servers = map[string]map[string]string{}
	if s.servers == nil {
		host_file, err := os.Open("host.json")
		if err != nil {
			log.Printf("[Repository] Failed to open host config file: %v", err)
			return err
		}
		defer host_file.Close()

		host := make(map[string]map[string]string)
		decoder := json.NewDecoder(host_file)
		err = decoder.Decode(&host)
		if err != nil {
			log.Printf("[Repository] Failed to decode config file: %v", err)
			return err
		}

		backends := map[string]map[string]string{}
		for k, v := range host {
			backends[k] = make(map[string]string)
			backends[k]["target"] = v["ip"] + ":" + v["port"]
			backends[k]["hostname"] = v["hostname"]
		}
		s.servers = backends
	}
	return nil
}

func (s *LocalServerRepositoryService) Resolve(hostname string) (string, error) {
	if data, ok := s.servers[hostname]; ok {
		if target, ok := data["target"]; ok {
			return target, nil
		}
	}
	return "", errors.New("host " + hostname + " not found")
}

func (s *LocalServerRepositoryService) List() ([]ServerList, error) {
	res := []ServerList{}
	for k, v := range s.servers {
		res = append(res, ServerList{Hostname: k, Address: v["target"]})
	}
	return res, nil
}

func NewLocalServerRepoService() ServerRepositoryService {
	return &LocalServerRepositoryService{}
}
