package utils

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/pkg/ent"
)

// ------------------- For Repository -------------------
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
	cmd2 := exec.Command("docker", "volume", "prune", "-f")
	cmd2.Dir = c.Onthemat.PWD
	err = cmd2.Run()

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

func RepoTestTruncateTable(ctx context.Context, c *ent.Client) error {
	_, err := c.ExecContext(ctx, `
	DO $$ DECLARE
	r RECORD;
	BEGIN
	FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
	EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY ' || ' CASCADE;';
	END LOOP;
	END $$;
	`)
	return err
}

// ------------------- For Hanlder -------------------
type TestResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

func MakeRespWithDataForTest[T any](body io.ReadCloser) TestResponse[T] {
	bodyBytes, _ := io.ReadAll(body)
	var result TestResponse[T]
	json.Unmarshal(bodyBytes, &result)
	return result
}

func MakeErrorForTests(body io.ReadCloser) common.HttpError {
	bodyBytes, _ := io.ReadAll(body)
	var result common.HttpError
	json.Unmarshal(bodyBytes, &result)
	return result
}
