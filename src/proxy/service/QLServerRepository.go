package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cznic/ql"
)

type QLServerRepositoryService struct {
	db *sql.DB

	newDB bool
}

func (s *QLServerRepositoryService) Load() error {
	if s.newDB {
		// localRepo := NewLocalServerRepoService()
		// err := localRepo.Load()
		// if err != nil {
		// 	return err
		// }
		sqlStmt := `
		CREATE TABLE IF NOT EXISTS server (
			server_hostname STRING,
			server_address STRING
		);
		`
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec(sqlStmt)
		if err != nil {
			return err
		}
		// list, err := localRepo.(ListableRepositoryService).List()
		// if err != nil {
		// 	{
		// 		return err
		// 	}
		// }
		// for _, record := range list {
		// 	_, err := tx.Exec("INSERT INTO server (server_hostname, server_address) VALUES ($1, $2)", record.Hostname, record.Address)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		if err = tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (s *QLServerRepositoryService) Resolve(hostname string) (string, error) {
	rows, err := s.db.Query("SELECT server_address FROM server WHERE server_hostname = $1", hostname)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			return "", err
		}
		log.Printf("[Host Resolver] %s %s", hostname, address)
		return address, nil
	}
	return "", errors.New("host " + hostname + " not found")
}

func (s *QLServerRepositoryService) Insert(hostname string, address string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO server (server_hostname, server_address) VALUES ($1, $2)", hostname, address)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *QLServerRepositoryService) Delete(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	row := tx.QueryRow("SELECT count(*) as c FROM server WHERE id() == $1", id)
	var recCount int
	err = row.Scan(&recCount)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if recCount != 0 {
		_, err = tx.Exec("DELETE FROM server WHERE id() == $1", id)
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("not found")
	}
}

func (s *QLServerRepositoryService) Upsert(id int, hostname string, address string) error {
	log.Printf("[QL] Upsert: %s %s", hostname, address)
	if hostname == "" || address == "" {
		return errors.New("value can't be empty")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT count(*) as c FROM server WHERE id() == $1", id)
	var recCount int
	err = row.Scan(&recCount)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if recCount == 0 {
		_, err = tx.Exec("INSERT INTO server (server_hostname, server_address) VALUES ($1, $2)", hostname, address)
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE server SET server_hostname = $1, server_address = $2 WHERE id() == $3", hostname, address, id)
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (s *QLServerRepositoryService) Destroy() {
	log.Printf("Disconnect from database")
	s.db.Close()
}

func (s *QLServerRepositoryService) Count() (int, error) {
	rows, err := s.db.Query("SELECT count(*) as c FROM server")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

func (s *QLServerRepositoryService) List() ([]ServerList, error) {
	rows, err := s.db.Query("SELECT id() as id, server_hostname, server_address FROM server ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ServerList{}
	for rows.Next() {
		var id int
		var hostname string
		var address string
		err = rows.Scan(&id, &hostname, &address)
		if err != nil {
			return nil, err
		}
		res = append(res, ServerList{ID: id, Hostname: hostname, Address: address})
	}
	return res, nil
}

func NewQLServerRepositoryService() ServerRepositoryService {
	ql.RegisterDriver()
	dbFileName := "./server.db"

	// Check if the database file exists
	isNew := false
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		fmt.Printf("Database file %s does not exist, creating...\n", dbFileName)
		isNew = true
	}
	db, err := sql.Open("ql", dbFileName)
	if err != nil {
		panic(err)
	}
	service := &QLServerRepositoryService{db: db, newDB: isNew}
	return service
}
