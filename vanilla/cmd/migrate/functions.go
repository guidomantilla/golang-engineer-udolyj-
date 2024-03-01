package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/environment"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/log"
)

func UpCmdFn(_ *cobra.Command, args []string) {
	var err error
	err = handleMigration(args, func(migration *migrate.Migrate) error {

		if err = migration.Up(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

type MigrationFunction func(migration *migrate.Migrate) error

func handleMigration(args []string, fn MigrationFunction) error {

	var err error
	ctx := context.Background()

	osArgs := os.Environ()
	env := environment.NewDefaultEnvironment(environment.WithArrays(osArgs, args))

	var cfg config.Config
	if err = config.Process(ctx, env, &cfg); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	url := strings.Replace(*cfg.DatasourceUrl, ":username", *cfg.DatasourceUsername, 1)
	url = strings.Replace(url, ":password", *cfg.DatasourcePassword, 1)
	url = strings.Replace(url, ":server", *cfg.DatasourceServer, 1)
	url = strings.Replace(url, ":service", *cfg.DatasourceService, 1)

	var db *sql.DB
	if db, err = sql.Open("mysql", url); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	workingDirectory, _ := os.Getwd()
	log.Info(fmt.Sprintf("working directory: %s", workingDirectory))
	migrationsDirectory := filepath.Join(workingDirectory, "resources/migrations/mysql")

	var driver database.Driver
	if driver, err = mysql.WithInstance(db, &mysql.Config{}); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	var migration *migrate.Migrate
	if migration, err = migrate.NewWithDatabaseInstance("file:///"+migrationsDirectory, *cfg.DatasourceService, driver); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	if err = fn(migration); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	return nil
}
