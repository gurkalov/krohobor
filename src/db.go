package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func getDbList() []string {
	cmd := exec.Command("psql",
		"-t", "-A", `-F","`,
		"-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('postgres', 'template1', 'template2');")
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return []string{}
	}

	return strings.Split(strings.Trim(string(stdout), "\n"), "\n")
}

func GetDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := getDbList()

	js, err := json.Marshal(list)
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Fprint(w, string(js))
}

func CreateBackup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := getDbList()

	for _, v := range list {
		cmd := exec.Command("pg_dump", "-Fc", v)
		res, err := cmd.Output()
		if err != nil {
			println(err.Error())
		}

		if err := ioutil.WriteFile("/tmp/backup_" + v + ".sql", res, 0644); err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Fprint(w, "OK")
}

func CreateRestore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := getDbList()

	for _, v := range list {
		cmd := exec.Command("dropdb",	v)
		_, err := cmd.Output()
		if err != nil {
			println(err.Error())
		}

		cmd = exec.Command("createdb",	v)
		_, err = cmd.Output()
		if err != nil {
			println(err.Error())
		}

		cmd = exec.Command("pg_restore",	"-d", v, "/tmp/backup_" + v + ".sql")
		_, err = cmd.Output()
		if err != nil {
			println(err.Error())
		}
	}

	fmt.Fprint(w, "OK")
}
