package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func main() {
	InitConfig()

	router := httprouter.New()
	router.GET("/v1/db", GetDB)
	router.GET("/v1/backup", CreateBackup)
	router.GET("/v1/restore", CreateRestore)
	router.GET("/v1/upload", UploadAws)

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(cfg.Port), router))
}
