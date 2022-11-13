package utils

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/pkg/ent"
)

type Tests struct {
	Name   string
	Before func(*testing.T)
	Expect func(*testing.T)
	After  func(*testing.T)
}

func RepoTestInit(t *testing.T) *config.Config {
	os.Setenv("GO_ENV", "TEST")
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(t)
		return nil
	}

	cmd := exec.Command("make", "docker_postgres_test")
	cmd.Dir = c.Onthemat.PWD
	err := cmd.Run()
	if err != nil {
		t.Error(err)
		return nil
	}
	return c
}

func RepoTestClose(t *testing.T) {
	os.Setenv("GO_ENV", "TEST")
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(t)
		return
	}

	cmd := exec.Command("docker", "rm", "-f", "psql_repo_test")
	cmd.Dir = c.Onthemat.PWD
	err := cmd.Run()
	if err != nil {
		t.Error(err)
		return
	}
}

func RepoTestRemoveTable(ctx context.Context, c *ent.Client) error {
	_, err := c.ExecContext(ctx, `
		DO $$ DECLARE
		r RECORD;
		BEGIN
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
		EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
		END LOOP;
		END $$;
	`)
	return err
}
