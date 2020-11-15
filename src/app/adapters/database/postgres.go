package database

import (
	"errors"
	"fmt"
	"krohobor/app/adapters/config"
	"krohobor/app/domain"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Postgres struct {
    cfg config.PostgresConfig
}

func NewPostgres(cfg config.PostgresConfig) Postgres {
	return Postgres{cfg}
}

func (p Postgres) List() ([]domain.Database, error) {
	var list []domain.Database

	cmd := exec.Command("psql", "-h", p.cfg.Host, "-U", p.cfg.User,
		"-t", "-A", `-F","`,
		"-c", "SELECT datname, pg_database_size(datname) FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1', 'template2');")
	stdout, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return list, errors.New(string(execErr.Stderr))
		}

		return list, err
	}

	dbs := strings.Split(strings.Trim(string(stdout), "\n"), "\n")
	for _, v := range dbs {
		cols := strings.Split(strings.ReplaceAll(v, `"`, ""), ",")
		// TODO:
		size, err := strconv.Atoi(cols[1])
		if err != nil {
			return list, err
		}

		db := domain.Database{
			Name: cols[0],
			Size: size,
		}
		list = append(list, db)
	}

	return list, nil
}

func (p Postgres) CreateDb(dbname string) error {
	cmd := exec.Command("createdb", dbname)
	_, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) Dump(dbname, filename string) error {
	cmd := exec.Command("pg_dump", "-h", p.cfg.Host, "-U", p.cfg.User,
		dbname, "-f", filename)
	_, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) DumpAll(filename string) error {
	cmd := exec.Command("pg_dumpall", "-h", p.cfg.Host, "-U", p.cfg.User, "-f", filename)
	_, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) Restore(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	cmd := exec.Command("psql",	"-f", filename, p.cfg.DB, "-h", p.cfg.Host, "-p", p.cfg.Port, "-U", p.cfg.User)
	cmd.Env = append(os.Environ(),
		"PGPASSWORD=" + p.cfg.Password,
	)

	out, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr) + ":" + string(out))
		}
		return err
	}

	return nil
}

func (p Postgres) Drop(dbname string) error {
	cmd := exec.Command("dropdb", dbname, "-h", p.cfg.Host, "-U", p.cfg.User)
	_, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) Tables(dbname string) ([]domain.Table, error) {
	var list []domain.Table

	cmd := exec.Command("psql", dbname, "-h", p.cfg.Host, "-U", p.cfg.User,
		"-t", "-A", `-F","`,
		"-c", "SELECT t.table_name, pg_total_relation_size(t.table_name::text), s.n_live_tup FROM information_schema.tables t JOIN pg_stat_user_tables s ON t.table_name = s.relname WHERE table_schema='public';")
	stdout, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return list, errors.New(string(execErr.Stderr))
		}

		return list, err
	}

	dbs := strings.Split(strings.Trim(string(stdout), "\n"), "\n")
	for _, v := range dbs {
		cols := strings.Split(strings.ReplaceAll(v, `"`, ""), ",")
		size, err := strconv.Atoi(cols[1])
		if err != nil {
			return list, err
		}

		count, err := strconv.Atoi(cols[2])
		if err != nil {
			return list, err
		}

		db := domain.Table{
			Name: cols[0],
			Size: size,
			Count: count,
		}
		list = append(list, db)
	}

	return list, nil
}

func (p Postgres) Count(dbname, table string) (int, error) {
	cmd := exec.Command("psql", dbname,
		"-t", "-A",
		"-c", fmt.Sprintf("SELECT count(*) FROM %s;", table))
	stdout, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

//func (p Postgres) Info(dbname string) (map[string]int, error) {
//	result := make(map[string]int, 0)
//	tables := p.Tables(dbname)
//	for _, v := range tables {
//		res, err := p.Count(dbname, v)
//		if err != nil {
//			return result, err
//		}
//		result[v] = res
//	}
//
//	return result, nil
//}
