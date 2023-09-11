package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type Mysql struct {
	host                string
	port                int
	user                string
	password            string
	database            string
	mysqldumpExecutable string
	mysqldumpOptions    []string
}

func NewMysql(
	host string,
	port int,
	user string,
	password string,
	database string,
	mysqldumpExecutable string,
	mysqldumpOptions []string,
) Mysql {
	return Mysql{
		host,
		port,
		user,
		password,
		database,
		mysqldumpExecutable,
		mysqldumpOptions,
	}
}

func (m Mysql) Dump(outPutFile *os.File) (err error) {
	command := exec.Command(
		m.mysqldumpExecutable,
		append(
			RemoveSomeValuesFromSliceByRegExp(
				m.mysqldumpOptions,
				regexp.MustCompile(
					`^(-h|--host|-p|--password|-P|--port|-u|--user)((=|\s)(\S+))?$`,
				),
			),
			fmt.Sprintf("--host=%s", m.host),
			fmt.Sprintf("--password=%s", m.password),
			fmt.Sprintf("--port=%d", m.port),
			fmt.Sprintf("--user=%s", m.user),
			m.database,
		)...,
	)

	command.Stdout = outPutFile

	var stderr bytes.Buffer
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		err = errors.New(stderr.String())
	}

	return err
}

type Postgres struct {
	host             string
	port             int
	user             string
	password         string
	database         string
	pgDumpExecutable string
	pdDumpOptions    []string
}

func NewPostgres(
	host string,
	port int,
	user string,
	password string,
	database string,
	pgDumpExecutable string,
	pdDumpOptions []string,
) Postgres {
	return Postgres{
		host,
		port,
		user,
		password,
		database,
		pgDumpExecutable,
		pdDumpOptions,
	}
}

func (pg Postgres) Dump(outPutFile *os.File) (err error) {
	command := exec.Command(
		pg.pgDumpExecutable,
		append(
			RemoveSomeValuesFromSliceByRegExp(
				pg.pdDumpOptions,
				regexp.MustCompile(
					`^(-d|--dbname|-f|--file|-h|--host|-p|--port|-U|--username)((=|\s)(\S+))?$`,
				),
			),
			fmt.Sprintf("--dbname=%s", pg.database),
			fmt.Sprintf("--file=%s", outPutFile.Name()),
			fmt.Sprintf("--host=%s", pg.host),
			fmt.Sprintf("--port=%d", pg.port),
			fmt.Sprintf("--username=%s", pg.user),
		)...,
	)

	command.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", pg.password),
	}

	var stderr bytes.Buffer
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		err = errors.New(stderr.String())
	}

	return err
}
