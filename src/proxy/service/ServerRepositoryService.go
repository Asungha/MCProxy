package service

type ServerRepositoryService interface {
	Load() error
	Resolve(hostname string) (string, error)
}

type UpdatableRepositoryService interface {
	Upsert(hostname string, address string) error
	Count() (int, error)
	Destroy()
	ServerRepositoryService
}

type ServerList struct {
	Hostname string
	Address  string
}
type ListableRepositoryService interface {
	List() ([]ServerList, error)
	ServerRepositoryService
}
