package httpserver

import (
	"github.com/julienschmidt/httprouter"
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/config"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/gateways/httpserver/handlers"
	"krohobor/app/usecases"
)

func Router(cfg config.Config) *httprouter.Router {

	dbPostgres := database.NewPostgres(cfg.Postgres)
	zipArchive := archive.Zip{
		Password: cfg.App.Password,
	}
	s3Storage := storage.NewAwsS3(cfg.App.Catalog, zipArchive)

	dbListHandler := handlers.DbList{
		UseCase: usecases.NewDbList(dbPostgres),
	}
	dbReadHandler := handlers.DbRead{
		UseCase: usecases.NewDbRead(dbPostgres),
	}
	dbDumpHandler := handlers.DbDump{
		UseCase: usecases.NewDbDump(dbPostgres, s3Storage),
	}

	dbDumpAllHandler := handlers.DbDumpAll{
		UseCase: usecases.NewDbDumpAll(dbPostgres, s3Storage),
	}

	dbRestoreHandler := handlers.DbRestore{
		UseCase: usecases.NewDbRestore(dbPostgres, s3Storage),
	}

	router := httprouter.New()
	router.GET("/v1/db", dbListHandler.Handle())
	router.GET("/v1/db/:db", dbReadHandler.Handle())
	router.POST("/v1/db/:db/dump", dbDumpHandler.Handle())
	router.POST("/v1/dump", dbDumpAllHandler.Handle())

	router.POST("/v1/db/:db/restore", dbRestoreHandler.Handle())
	router.POST("/v1/db/:db/restore/:name", dbRestoreHandler.Handle())

	router.POST("/v1/restore", dbRestoreHandler.Handle())
	router.POST("/v1/restore/:name", dbRestoreHandler.Handle())

	return router
}
