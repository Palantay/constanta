package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Transaction struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Email        string    `json:"user_email"`
	Amount       int       `json:"amount"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`
	CancelStatus bool      `json:"cancel_status"`
}

func (t *Transaction) ValidateForPostTransaction() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.UserID, validation.Required),
		validation.Field(&t.Email, validation.Required, is.Email),
		validation.Field(&t.Amount, validation.Required),
		validation.Field(&t.Currency, validation.Required),
	)
}

type TransactionStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

const (
	New       string = "НОВЫЙ"
	Success          = "УСПЕХ"
	Unsuccess        = "НЕУСПЕХ"
	Err              = "ОШИБКА"
)
