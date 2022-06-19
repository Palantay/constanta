package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	config                *Config
	db                    *sql.DB
	transactionRepository *TransactionRepository
	userRepository        *UserRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURI)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("Database is connect")
	return nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (s *Storage) Transaction() *TransactionRepository {
	if s.transactionRepository != nil {
		return s.transactionRepository
	}
	s.transactionRepository = &TransactionRepository{
		storage: s,
	}
	return s.transactionRepository
}

func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}

	return s.userRepository
}
