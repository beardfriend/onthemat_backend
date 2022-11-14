package utils

import (
	"os"
	"os/exec"

	"onthemat/internal/app/config"
)

func InitMigrationPostgres() *config.Config {
	os.Setenv("GO_ENV", "TEST")
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}

	cmd := exec.Command("make", "docker_postgres_test")
	cmd.Dir = c.Onthemat.PWD
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return c
}

func CloseMigrationPostgres() {
	os.Setenv("GO_ENV", "TEST")
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}

	cmd := exec.Command("docker", "rm", "-f", "psql_repo_test")
	cmd.Dir = c.Onthemat.PWD
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
