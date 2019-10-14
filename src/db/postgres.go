package db

import (
	"io/ioutil"
	"os/exec"
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

func (p Postgres) Restore(filename string) error {
	cmd := exec.Command("pg_restore",	"-d", p.Name, filename)
	_, err := cmd.Output()
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

