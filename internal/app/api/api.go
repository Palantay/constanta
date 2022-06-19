package api

import (
	"net/http"

	"github.com/Palantay/constanta/internal/middleware"
	"github.com/Palantay/constanta/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type API struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (api *API) Start() error {
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	api.logger.Info("Starting api server at port:", api.config.Port)

	api.configureRouterField()

	if err := api.configureStorageField(); err != nil {
		return err
	}

	return http.ListenAndServe(api.config.Port, api.router)
}

var (
	prefix string = "/api"
)

func (a *API) configureRouterField() {
	a.router.HandleFunc(prefix+"/transaction", a.PostTransaction).Methods("POST")
	a.router.HandleFunc(prefix+"/transaction/{id}", a.GetUserTransactionsByUserId).Methods("GET")
	a.router.HandleFunc(prefix+"/transaction", a.GetUserTransactionsByUserEmail).Methods("GET")
	a.router.HandleFunc(prefix+"/transaction/status/{id}", a.GetStatusTransaction).Methods("GET")
	a.router.HandleFunc(prefix+"/user/auth", a.PostToAuth).Methods("POST")
	a.router.Handle(prefix+"/transaction/status", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(a.SetTransactionStatus),
	)).Methods("PUT")
	a.router.HandleFunc(prefix+"/transaction/cancel", a.SetCancelStatus).Methods("PUT")

}

func (a *API) configureLoggerField() error {
	log_level, err := logrus.ParseLevel(a.config.LoggerLevel)

	if err != nil {
		return err
	}
	a.logger.SetLevel(log_level)
	return nil
}

func (a *API) configureStorageField() error {
	storage := storage.New(a.config.Storage)
	if err := storage.Open(); err != nil {
		return err
	}

	a.storage = storage
	return nil
}
