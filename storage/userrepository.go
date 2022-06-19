package storage

import (
	"log"

	"github.com/Palantay/constanta/internal/app/models"
)

type UserRepository struct {
	storage *Storage
}

func (ur *UserRepository) FindUserByLogin(login string) (*models.User, bool, error) {
	founded := true
	var u models.User

	query := "SELECT * FROM users WHERE login=$1"

	err := ur.storage.db.QueryRow(query, login).Scan(&u.Login, &u.Password)

	if err != nil {
		log.Println(err)
	}

	if u.Login == "" {
		founded = false
		return nil, founded, nil
	}

	return &u, founded, nil

}
