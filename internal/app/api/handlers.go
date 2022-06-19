package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Palantay/constanta/internal/app/models"
	"github.com/Palantay/constanta/internal/middleware"
	"github.com/form3tech-oss/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

func (api *API) PostTransaction(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Post Transaction POST /api/transaction")

	var t models.Transaction

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		ResponseJSON(w, 400, true, "Provided json is invalid")
		return
	}

	if err := t.ValidateForPostTransaction(); err != nil {
		ResponseJSON(w, 400, true, fmt.Sprint(err))
		return
	}

	tl, err := api.storage.Transaction().Create(&t)

	if err != nil {
		api.logger.Info("Troubles while creating new transaction:", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(tl)

}

func (api *API) GetTransactionsByUserId(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Get user transactions by user id /api/transaction{id}")

	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		api.logger.Info("Troubles while parsing {id} param: ", err)
		ResponseJSON(w, 400, true, "ID is not format")
		return
	}

	tl, ok, err := api.storage.Transaction().FindUserTransactionByUserId(id)

	if err != nil {
		api.logger.Info("Error while FindUserTransactionByUserId", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	if !ok {
		api.logger.Info("User with this {id} does not exist")
		ResponseJSON(w, 404, true, "User with this {id} does not exist")
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(tl)
}

func (api *API) GetTransactionsByUserEmail(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Get Get user transactions by user email /api/transaction{email}")

	email := r.URL.Query().Get("email")

	if err := validation.Validate(email, is.Email); err != nil {
		ResponseJSON(w, 400, true, fmt.Sprint(err))
		return
	}

	t, ok, err := api.storage.Transaction().FindUserTransactionByUserEmail(email)

	if err != nil {
		api.logger.Info("Error while FindUserTransactionByUserId", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	if !ok {
		api.logger.Info("User with this {email} does not exist")
		ResponseJSON(w, 404, true, "User with this {email} does not exist")
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(t)
}

func (api *API) GetStatusTransaction(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Get status transactions by id /api/transaction/status/{id}")

	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		api.logger.Info("Troubles while parsing {id} param: ", err)
		ResponseJSON(w, 400, true, "ID is not format")
		return
	}

	st, ok, err := api.storage.Transaction().FindStatusTransactionById(id)

	if err != nil {
		api.logger.Info("Error while Get user transactions by user email", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	if !ok {
		api.logger.Info("Transaction with this {id} does not exist")
		ResponseJSON(w, 404, true, "Transaction with this {id} does not exist")
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(st)

}

func (api *API) PostToAuth(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Post to auth POST /api/user/auth")

	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		ResponseJSON(w, 400, true, "Provided json is invalid")
		return
	}

	uInDB, ok, err := api.storage.User().FindUserByLogin(u.Login)

	if err != nil {
		api.logger.Info("Error while FindUserByLogin", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return

	}

	if !ok {
		api.logger.Info("User with this {login} does not exist")
		ResponseJSON(w, 404, true, "User with this {login} does not exist")
		return
	}

	if uInDB.Password != u.Password {
		api.logger.Info("Invalid credentials to auth ")
		ResponseJSON(w, 404, true, "Password is invalid")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	tokenString, err := token.SignedString(middleware.SecretKey)

	if err != nil {
		api.logger.Info("Can not claim jwt-token")
		ResponseJSON(w, 500, true, "We have some troubles. Try again")
		return
	}

	ResponseJSON(w, 201, false, tokenString)
}

func (api *API) SetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Patch status /api/transaction/status")

	var (
		ts    *models.Transaction
		tInDB *models.Transaction
	)

	err := json.NewDecoder(r.Body).Decode(&ts)

	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		ResponseJSON(w, 400, true, "Provided json is invalid")
		return
	}

	if err = validation.Validate(ts.Status, validation.In(models.Err, models.New, models.Success, models.Unsuccess)); err != nil {
		ResponseJSON(w, 400, true, fmt.Sprint(err))
		return
	}

	var ok bool
	tInDB, ok, err = api.storage.Transaction().FindTransactionById(ts.ID)

	if !ok {
		api.logger.Info("Transaction with this {id} does not exist")
		ResponseJSON(w, 404, true, "Transaction with this {id} does not exist")
		return
	}

	if err != nil {
		api.logger.Info("Error while find transaction by id", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	if tInDB.CancelStatus {
		api.logger.Info("Attempt to cancel a transaction that has already been canceled")
		ResponseJSON(w, 400, true, "This transaction has been canceled")
		return
	}

	if tInDB.Status == models.Success || tInDB.Status == models.Unsuccess {
		api.logger.Info("Attempt to change the status with the value 'Success' or 'Unsuccess'")
		ResponseJSON(w, 400, true, "Uncorrected status")
		return
	}

	err = api.storage.Transaction().UpdateTransactionStatus(ts)
	if err != nil {
		api.logger.Info("Error while set transaction status", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	ResponseJSON(w, 204, false, "Status changed")

}

func (api *API) SetCancelStatus(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	api.logger.Info("Patch cancel status /api/transaction/cancel")

	var (
		ts    *models.Transaction
		tInDB *models.Transaction
	)

	err := json.NewDecoder(r.Body).Decode(&ts)

	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		ResponseJSON(w, 400, true, "Provided json is invalid")
		return
	}

	var ok bool
	tInDB, ok, err = api.storage.Transaction().FindTransactionById(ts.ID)
	if !ok {
		api.logger.Info("Transaction with this {id} does not exist")
		ResponseJSON(w, 404, true, "Transaction with this {id} does not exist")
		return
	}

	if err != nil {
		api.logger.Info("Error while find transaction by id", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	if tInDB.CancelStatus {
		api.logger.Info("Attempt to cancel a transaction that has already been canceled")
		ResponseJSON(w, 400, true, "This transaction has already been canceled")
		return
	}

	if tInDB.Status == models.Success || tInDB.Status == models.Unsuccess {
		api.logger.Info("Attempt to change the cancel status with the value 'Success' or 'Unsuccess'")
		ResponseJSON(w, 400, true, "It is not possible to cancel this transaction")
		return
	}

	err = api.storage.Transaction().UpdateCancelStatus(ts)
	if err != nil {
		api.logger.Info("Error while cancel status", err)
		ResponseJSON(w, 501, true, "We have some troubles to accessing data base. Try again")
		return
	}

	ResponseJSON(w, 204, false, "Cancel status changed")

}
