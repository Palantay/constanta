package storage

import (
	"log"
	"time"

	"github.com/Palantay/constanta/internal/app/models"
)

type TransactionRepository struct {
	storage *Storage
}

func (tr *TransactionRepository) Create(t *models.Transaction) (*models.Transaction, error) {

	r := Random()

	if r == 3 {
		query := "INSERT INTO user_transaction (user_id, user_email, amount, currency, status) VALUES($1, $2, $3, $4, $5) RETURNING id"
		t.Status = models.Err
		if err := tr.storage.db.QueryRow(query, t.UserID, t.Email, t.Amount, t.Currency, t.Status).Scan(&t.ID); err != nil {
			return nil, err
		}
	} else {
		query := "INSERT INTO user_transaction (user_id, user_email, amount, currency) VALUES($1, $2, $3, $4) RETURNING id"
		if err := tr.storage.db.QueryRow(query, t.UserID, t.Email, t.Amount, t.Currency).Scan(&t.ID); err != nil {
			return nil, err
		}
	}

	query := "SELECT * FROM user_transaction WHERE id=$1"

	if err := tr.storage.db.QueryRow(query, t.ID).Scan(&t.ID, &t.UserID, &t.Email, &t.Amount, &t.Currency, &t.CreatedAt, &t.UpdatedAt, &t.Status, &t.CancelStatus); err != nil {
		return nil, err
	}

	return t, nil
}

func (tr *TransactionRepository) FindUserTransactionByUserId(id int) ([]*models.Transaction, bool, error) {
	founded := true

	query := "SELECT * FROM user_transaction WHERE user_id = $1"
	rows, err := tr.storage.db.Query(query, id)
	if err != nil {
		return nil, founded, err
	}
	defer rows.Close()
	tl := make([]*models.Transaction, 0)

	for rows.Next() {
		t := models.Transaction{}
		err := rows.Scan(&t.ID, &t.UserID, &t.Email, &t.Amount, &t.Currency, &t.CreatedAt, &t.UpdatedAt, &t.Status, &t.CancelStatus)
		if err != nil {
			log.Println(err)
			continue
		}
		tl = append(tl, &t)
	}

	if len(tl) == 0 {
		founded = false
		return nil, founded, nil
	}

	return tl, founded, nil
}

func (tr *TransactionRepository) FindUserTransactionByUserEmail(email string) ([]*models.Transaction, bool, error) {
	founded := true

	query := "SELECT * FROM user_transaction WHERE user_email = $1"
	rows, err := tr.storage.db.Query(query, email)
	if err != nil {
		return nil, founded, err
	}
	defer rows.Close()
	tl := make([]*models.Transaction, 0)

	for rows.Next() {
		t := models.Transaction{}
		err := rows.Scan(&t.ID, &t.UserID, &t.Email, &t.Amount, &t.Currency, &t.CreatedAt, &t.UpdatedAt, &t.Status, &t.CancelStatus)
		if err != nil {
			log.Println(err)
			continue
		}
		tl = append(tl, &t)
	}

	if len(tl) == 0 {
		founded = false
		return nil, founded, nil
	}

	return tl, founded, nil
}

func (tr *TransactionRepository) FindStatusTransactionById(id int) (*models.TransactionStatus, bool, error) {
	founded := true
	var ts models.TransactionStatus

	query := "SELECT id, status FROM user_transaction WHERE id=$1"

	err := tr.storage.db.QueryRow(query, id).Scan(&ts.ID, &ts.Status)
	if err != nil {
		log.Println(err)
	}

	if ts.ID == 0 {
		founded = false
		return nil, founded, nil
	}

	return &ts, founded, nil
}

func (tr *TransactionRepository) FindTransactionById(id int) (*models.Transaction, bool, error) {
	founded := true
	var t models.Transaction

	query := "SELECT * FROM user_transaction WHERE id=$1"

	err := tr.storage.db.QueryRow(query, id).Scan(&t.ID, &t.UserID, &t.Email, &t.Amount, &t.Currency, &t.CreatedAt, &t.UpdatedAt, &t.Status, &t.CancelStatus)
	if err != nil {
		log.Println(err)
	}

	if t.ID == 0 {
		founded = false
		return nil, founded, nil
	}

	return &t, founded, nil
}

func (tr *TransactionRepository) UpdateTransactionStatus(t *models.Transaction) error {

	t.UpdatedAt = time.Now().UTC()

	query := "UPDATE user_transaction SET status = $1, updated_at = $2  WHERE id=$3"
	var err error
	_, err = tr.storage.db.Exec(query, t.Status, t.UpdatedAt, t.ID)
	if err != nil {
		log.Println(err)
	}

	return nil

}

func (tr *TransactionRepository) UpdateCancelStatus(t *models.Transaction) error {
	t.UpdatedAt = time.Now().UTC()

	query := "UPDATE user_transaction SET cancel_status = $1, updated_at = $2 WHERE id = $3"

	var err error
	_, err = tr.storage.db.Exec(query, t.CancelStatus, t.UpdatedAt, t.ID)
	if err != nil {
		log.Println(err)
	}

	return nil
}
