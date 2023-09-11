package main

import (
	"errors"
	"fmt"
	"github.com/adebayohountondji/aws-rds-dump/internal"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	databaseDriverMysql        = "mysql"
	databaseDriverPostgres     = "postgres"
	defaultDumpFilenamePattern = "{yyyy}{MM}{dd}-{hh}{mm}{ss}"
)

type rds interface {
	Dump(outPutFile *os.File) (err error)
}

func createRdsFromOsEnv() (rds rds, err error) {
	switch os.Getenv("DATABASE_DRIVER") {
	case databaseDriverMysql:
		port, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
		rds = internal.NewMysql(
			os.Getenv("MYSQL_HOST"),
			port,
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_DATABASE"),
			os.Getenv("MYSQLDUMP_EXECUTABLE"),
			internal.FindAllCommandOptionsFromString(os.Getenv("MYSQLDUMP_OPTIONS")),
		)
	case databaseDriverPostgres:
		port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		rds = internal.NewPostgres(
			os.Getenv("POSTGRES_HOST"),
			port,
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DATABASE"),
			os.Getenv("PG_DUMP_EXECUTABLE"),
			internal.FindAllCommandOptionsFromString(os.Getenv("PG_DUMP_OPTIONS")),
		)
	default:
		err = errors.New(
			fmt.Sprintf(
				"%s is not a valid rds driver. Available values are: %s\n",
				os.Getenv("DATABASE_DRIVER"),
				strings.Join([]string{databaseDriverMysql, databaseDriverPostgres}, ", "),
			),
		)
	}

	return rds, err
}

func createAwsFromOsEnv() (internal.Aws, error) {
	return internal.NewAws(
		internal.AwsConfig{
			AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Region:          os.Getenv("AWS_REGION"),
		},
	)
}

func tryToDump() (err error) {
	rds, err := createRdsFromOsEnv()
	if err != nil {
		return err
	}

	tmpDumpFile, err := os.CreateTemp("", "aws-rds-dump-")
	if err != nil {
		return err
	}
	defer func(filename string) {
		_ = os.Remove(filename)
	}(tmpDumpFile.Name())

	err = rds.Dump(tmpDumpFile)
	if err != nil {
		return err
	}

	tmpGzipDumpFile, err := os.CreateTemp("", "aws-rds-dump-")
	if err != nil {
		return err
	}
	defer func(filename string) {
		_ = os.Remove(filename)
	}(tmpGzipDumpFile.Name())

	err = internal.CreateGzipFile(tmpDumpFile.Name(), tmpGzipDumpFile.Name())
	if err != nil {
		return err
	}

	aws, err := createAwsFromOsEnv()
	if err != nil {
		return err
	}

	dumpFilenamePattern := os.Getenv("DUMP_FILENAME_PATTERN")
	if dumpFilenamePattern == "" {
		dumpFilenamePattern = defaultDumpFilenamePattern
	}

	s3 := aws.NewS3(os.Getenv("AWS_S3_BUCKET_NAME"))
	err = s3.PutObject(
		internal.CreateDumpFilenameFromPatternAndDate(
			fmt.Sprintf("%s.gz", dumpFilenamePattern),
			time.Now(),
		),
		tmpGzipDumpFile,
	)

	return err
}

func main() {
	err := tryToDump()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("aws-rds-dump: main: command executed successfully")
	os.Exit(0)
}
