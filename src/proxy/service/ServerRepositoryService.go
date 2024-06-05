package service

type ServerRepositoryService interface {
	Load() error
	Resolve(hostname string) (string, error)
}
