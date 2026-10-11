package database

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"krohobor/app/adapters/config"
	"krohobor/app/domain"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Postgres struct {
	cfg      config.PostgresConfig
	progress func(done, total int)
}

func NewPostgres(cfg config.PostgresConfig) *Postgres {
	return &Postgres{cfg: cfg}
}

// SetProgress registers a callback invoked while a dump is restored via psql.
func (p *Postgres) SetProgress(fn func(done, total int)) {
	p.progress = fn
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
	selectSql := "SELECT datname, pg_database_size(datname) FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1', 'template2');"
	stdout, err := p.cmd(p.cfg, "psql", "-t", "-A", `-F","`,
		"-c", selectSql, p.cfg.DB)
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
	_, err := p.cmd(p.cfg, "pg_dumpall", "-f", filename)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) Restore(filename, dbname string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	if p.progress != nil {
		return p.restoreWithProgress(filename, dbname)
	}

	out, err := p.cmd(p.cfg, "psql", "-f", filename, dbname)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr) + ":" + string(out))
		}
		return err
	}

	return nil
}

// restoreWithProgress restores a dump by streaming it to psql over stdin
// (byte-identical to "psql -f <file>") and reports the number of bytes sent.
func (p *Postgres) restoreWithProgress(filename, dbname string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	total := int(info.Size())

	c := exec.Command("psql", dbname,
		"-h", p.cfg.Host, "-U", p.cfg.User, "-p", strconv.Itoa(p.cfg.Port))
	c.Env = append(os.Environ(), "PGPASSWORD="+p.cfg.Password)

	var out, errOut bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errOut

	pipe, err := c.StdinPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		pipe.Close()
		c.Wait()
		return err
	}

	buf := make([]byte, 256*1024)
	var done int
	for {
		n, rerr := f.Read(buf)
		if n > 0 {
			if _, err := pipe.Write(buf[:n]); err != nil {
				f.Close()
				pipe.Close()
				c.Wait()
				return err
			}
			done += n
			p.progress(done, total)
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			pipe.Close()
			c.Wait()
			return rerr
		}
	}

	f.Close()
	pipe.Close()

	if err := c.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return errors.New(errOut.String() + ":" + out.String())
		}
		return err
	}

	return nil
}

func (p Postgres) Drop(dbname string, force bool) error {
	var err error
	if force {
		err = p.forceDrop(dbname)
	} else {
		err = p.drop(dbname)
	}

	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(execErr.Stderr))
		}
		return err
	}

	return nil
}

func (p Postgres) drop(dbname string) error {
	_, err := p.cmd(p.cfg, "dropdb", dbname)
	return err
}

func (p Postgres) forceDrop(dbname string) error {
	_, err := p.cmd(p.cfg, "psql", "-t", "-A",
		"-c", fmt.Sprintf("DROP database %s WITH (FORCE);", dbname))
	return err
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
			Name:  cols[0],
			Size:  size,
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
		"PGPASSWORD="+cfg.Password,
	)

	return cmd.Output()
}

//func (p Postgres) targetConfig(target string) (config.PostgresConfig, error) {
//	cfg := p.cfg
//	if target == "" {
//		return cfg, nil
//	}
//
//	partTarget := strings.Split(target, ":")
//	if len(partTarget) > 0 {
//		cfg.Host = partTarget[0]
//	}
//	if len(partTarget) > 1 {
//		port, err := strconv.Atoi(partTarget[1])
//		if err != nil {
//			return cfg, err
//		}
//
//		cfg.Port = port
//	}
//
//	return cfg, nil
//}
//
