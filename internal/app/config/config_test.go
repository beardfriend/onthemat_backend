package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("환경설정 없을 때", func(t *testing.T) {
		c := NewConfig()
		if err := c.Load("../../../configs"); err != nil {
			t.Error(err)
		}
		assert.Equal(t, c.MariaDB.Host, "localhost")
	})

	t.Run("DEV 환경일 때", func(t *testing.T) {
		os.Setenv("GO_ENV", "DEV")
		c := NewConfig()
		if err := c.Load("../../../configs"); err != nil {
			t.Error(err)
		}
	})

	t.Run("PROD 환경일 때", func(t *testing.T) {
		os.Setenv("GO_ENV", "PROD")
		c := NewConfig()
		if err := c.Load("../../../configs"); err != nil {
			t.Error(err)
		}
	})

	t.Run("TEST 환경일 때", func(t *testing.T) {
		os.Setenv("GO_ENV", "TEST")
		c := NewConfig()
		if err := c.Load("../../../configs"); err != nil {
			t.Error(err)
		}
	})
}
