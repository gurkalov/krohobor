package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Postgres struct {
	Name string
}

func (p *Postgres) Init(name string) error {
	p.Name = name

	return nil
}

func (p Postgres) Create() error {
	cmd := exec.Command("createdb", p.Name)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) Dump(filename string) error {
	cmd := exec.Command("pg_dump", "-Fc", p.Name)
	res, err := cmd.Output()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, res, 0644); err != nil {
		return err
	}

	return nil
}

func (p Postgres) DumpAll(filename string) error {
	cmd := exec.Command("pg_dumpall")
	res, err := cmd.Output()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, res, 0644); err != nil {
		return err
	}

	return nil
}

func (p Postgres) Restore(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	cmd := exec.Command("pg_restore",	"-d", p.Name, filename)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) RestoreAll(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	cmd := exec.Command("psql",	"-f", filename, p.Name)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) Drop() error {
	cmd := exec.Command("dropdb",	p.Name)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) List() []string {
	cmd := exec.Command("psql",
		"-t", "-A", `-F","`,
		"-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1', 'template2');")
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return []string{}
	}

	return strings.Split(strings.Trim(string(stdout), "\n"), "\n")
}

func (p Postgres) Tables() []string {
	cmd := exec.Command("psql", p.Name,
		"-t", "-A", `-F","`,
		"-c", "SELECT table_name FROM information_schema.tables WHERE table_schema='public';")
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return []string{}
	}

	return strings.Split(strings.Trim(string(stdout), "\n"), "\n")
}

func (p Postgres) Count(table string) int {
	cmd := exec.Command("psql", p.Name,
		"-t", "-A",
		"-c", fmt.Sprintf("SELECT count(*) FROM %s;", table))
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return 0
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		println(err.Error())
		return 0
	}

	return count
}

func (p Postgres) Info() map[string]int {
	result := make(map[string]int, 0)
	tables := p.Tables()
	for _, v := range tables {
		result[v] = p.Count(v)
	}

	return result
}
