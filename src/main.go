package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func main() {
	InitConfig()
	InitManager()

	router := httprouter.New()
	router.GET("/v1/db", GetDB)
	router.GET("/v1/backup", ListBackup)
	router.POST("/v1/backup", CreateBackup)
	router.POST("/v1/backup/:name", RestoreBackup)
	router.POST("/v1/restore/last", RestoreLastBackup)
	router.DELETE("/v1/backup/:name", DeleteBackup)

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(cfg.App.Port), router))
}
