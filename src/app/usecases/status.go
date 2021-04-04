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
	cfg     config.Config
	db      database.Interface
	storage storage.Interface
}

type StatusRequest struct{}

type StatusResponse struct {
	Db      StatusDb
	Storage StatusStorage
}

type StatusDb struct {
	Check bool
	Error error
}

type StatusStorage struct {
	Check bool
	Error error
}

func NewStatus(cfg config.Config, db database.Interface, storage storage.Interface) *Status {
	return &Status{cfg, db, storage}
}

func (s *Status) Execute(request StatusRequest) (StatusResponse, error) {
	response := StatusResponse{}

	errDb := s.db.Check()
	response.Db.Check = errDb == nil
	response.Db.Error = errDb

	errStorage := s.storage.Check()
	response.Storage.Check = errStorage == nil
	response.Storage.Error = errStorage

	//if request.Target == "" {
	//	response.Db.Host = s.cfg.Postgres.Host + ":" + strconv.Itoa(s.cfg.Postgres.Port)
	//} else {
	//	response.Db.Host = request.Target
	//}
	//response.Storage.Catalog = s.cfg.App.Catalog

	return response, nil
}
