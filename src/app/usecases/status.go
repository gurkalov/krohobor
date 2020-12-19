package usecases

import (
	"krohobor/app/adapters/config"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
)

type StatusInterface interface {
	Execute(StatusRequest) (StatusResponse, error)
}

type Status struct {
	cfg config.Config
	db database.Interface
	storage storage.Interface
}

type StatusRequest struct {
	Target string
}

type StatusResponse struct {
	Db struct {
		Check bool
		Host  string
	}
	Storage struct {
		Check bool
		Host  string
	}
}

func NewStatus(cfg config.Config, db database.Interface, storage storage.Interface) *Status {
	return &Status{cfg, db, storage}
}

func (s *Status) Execute(request StatusRequest) (StatusResponse, error) {
	response := StatusResponse{}

	response.Db.Check = true
	response.Storage.Check = s.storage.Check() == nil

	response.Db.Host = s.cfg.Postgres.Host + ":" + s.cfg.Postgres.Port
	response.Storage.Host = s.cfg.App.Catalog

	return response, nil
}
