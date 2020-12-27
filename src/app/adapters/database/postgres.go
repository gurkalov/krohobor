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

func (p Postgres) Check() error {
	out, err := p.cmd(p.cfg, "psql", "-l")
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr) + ":" + string(out))
		}
		return err
	}

	return nil
}

func (p Postgres) List() ([]domain.Database, error) {
	var list []domain.Database
	stdout, err := p.cmd(p.cfg, "psql", "-t", "-A", `-F","`,
		"-c", "SELECT datname, pg_database_size(datname) FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1', 'template2');")
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return list, errors.New(string(execErr.Stderr))
		}

		return list, err
	}

	if len(stdout) == 0 {
		return list, nil
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
	_, err := p.cmd(p.cfg, "createdb", dbname)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) Dump(dbname, filename string) error {
	_, err := p.cmd(p.cfg, "pg_dump", dbname, "-f", filename)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) DumpAll(filename string) error {
	_, err := p.cmd(p.cfg, "pg_dumpall","-f", filename)
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

	out, err := p.cmd(p.cfg, "psql", "-f", filename)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr) + ":" + string(out))
		}
		return err
	}

	return nil
}

func (p Postgres) Drop(dbname string) error {
	_, err := p.cmd(p.cfg, "dropdb", dbname)
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

	stdout, err := p.cmd(p.cfg, "psql", dbname, "-t", "-A", `-F","`,
		"-c", "SELECT t.table_name, pg_total_relation_size(t.table_name::text), s.n_live_tup FROM information_schema.tables t JOIN pg_stat_user_tables s ON t.table_name = s.relname WHERE table_schema='public';")
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return list, errors.New(string(execErr.Stderr))
		}

		return list, err
	}

	if len(stdout) == 0 {
		return list, nil
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
	stdout, err := p.cmd(p.cfg, "psql", "-t", "-A",
		"-c", fmt.Sprintf("SELECT count(*) FROM %s;", table))
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p Postgres) cmd(cfg config.PostgresConfig, name string, arg ...string) ([]byte, error) {
	arg = append(arg, "-h", cfg.Host, "-U", cfg.User, "-p", strconv.Itoa(cfg.Port))

	cmd := exec.Command(name, arg...)
	cmd.Env = append(os.Environ(),
		"PGPASSWORD=" + cfg.Password,
	)

	return cmd.Output()
}

func (p Postgres) targetConfig(target string) (config.PostgresConfig, error) {
	cfg := p.cfg
	if target == "" {
		return cfg, nil
	}

	partTarget := strings.Split(target, ":")
	if len(partTarget) > 0 {
		cfg.Host = partTarget[0]
	}
	if len(partTarget) > 1 {
		port, err := strconv.Atoi(partTarget[1])
		if err != nil {
			return cfg, err
		}

		cfg.Port = port
	}

	return cfg, nil
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
