package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/zenstats/internal/config"
)

func TestCreatingNewValidation(t *testing.T) {
	envs := []string{
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_PORT", "POSTGRES_DB", "POSTGRES_HOST", "POSTGRES_SSLMODE",
	}

	fullSetup := func() {
		os.Setenv("POSTGRES_USER", "user")
		os.Setenv("POSTGRES_PASSWORD", "password")
		os.Setenv("POSTGRES_PORT", "5432")
		os.Setenv("POSTGRES_DB", "database")
		os.Setenv("POSTGRES_HOST", "database.com")
		os.Setenv("POSTGRES_SSLMODE", "disable")
	}

	clear := func() {
		for _, env := range envs {
			os.Unsetenv(env)
		}
	}

	testCases := []struct {
		Description string
		Setup       func()
		ExpectedCfg *config.Database
		ExpectedErr error
	}{
		{
			Description: "testing a complete env setup",
			Setup: func() {
				fullSetup()
			},
			ExpectedCfg: &config.Database{
				Username: "user",
				Password: "password",
				Port:     5432,
				Host:     "database.com",
				DBName:   "database",
				SSLMode:  "disable",
			},
		},
		{
			Description: "no setup",
			Setup: func() {
				clear()
			},
			ExpectedErr: fmt.Errorf("no POSTGRES_USER env variable set"),
		},
		{
			Description: "no pg user",
			Setup: func() {
				fullSetup()
				os.Unsetenv("POSTGRES_USER")
			},
			ExpectedErr: fmt.Errorf("no POSTGRES_USER env variable set"),
		},
		{
			Description: "no pg password",
			Setup: func() {
				fullSetup()
				os.Unsetenv("POSTGRES_PASSWORD")
			},
			ExpectedErr: fmt.Errorf(
				"loading password: no POSTGRES_PASSWORD or POSTGRES_PASSWORD_FILE env var set",
			),
		},
		{
			Description: "no pg host",
			Setup: func() {
				fullSetup()
				os.Unsetenv("POSTGRES_HOST")
			},
			ExpectedErr: fmt.Errorf("no POSTGRES_HOST env variable set"),
		},
		{
			Description: "no pg port",
			Setup: func() {
				fullSetup()
				os.Unsetenv("POSTGRES_PORT")
			},
			ExpectedErr: fmt.Errorf("no POSTGRES_PORT env variable set"),
		},
		{
			Description: "no pg database",
			Setup: func() {
				fullSetup()
				os.Unsetenv("POSTGRES_DB")
			},
			ExpectedErr: fmt.Errorf("no POSTGRES_DB env variable set"),
		},
		{
			Description: "invalid port setup",
			Setup: func() {
				fullSetup()
				os.Setenv("POSTGRES_PORT", "helloworld")
			},
			ExpectedErr: fmt.Errorf(
				"failed to convert port to int: strconv.Atoi: parsing \"helloworld\": invalid syntax",
			),
		},
		{
			Description: "empty db name",
			Setup: func() {
				fullSetup()
				os.Setenv("POSTGRES_DB", "")
			},
			ExpectedErr: fmt.Errorf(
				"failed to validate config: invalid database name",
			),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.Description, func(t *testing.T) {
			test.Setup()

			cfg, err := config.NewDatabase()

			if test.ExpectedErr != nil {
				assert.EqualError(t, err, test.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.ExpectedCfg, cfg)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	valid := config.Database{
		Username: "pguser",
		Password: "pgpassword",
		Host:     "pghost",
		Port:     5432,
		DBName:   "pgdatabase",
		SSLMode:  "disable",
	}

	t.Run("Test valid", func(t *testing.T) {
		valid := valid
		assert.NoError(t, valid.Validate())
	})

	t.Run("Test invalid username", func(t *testing.T) {
		valid := valid
		valid.Username = ""
		assert.EqualError(t, valid.Validate(), "invalid username")
	})

	t.Run("Test invalid password", func(t *testing.T) {
		valid := valid
		valid.Password = ""
		assert.EqualError(t, valid.Validate(), "invalid password")
	})

	t.Run("Test invalid host", func(t *testing.T) {
		valid := valid
		valid.Host = ""
		assert.EqualError(t, valid.Validate(), "invalid host")
	})

	t.Run("Test invalid port", func(t *testing.T) {
		valid := valid
		valid.Port = 0
		assert.EqualError(t, valid.Validate(), "invalid port")
	})

	t.Run("Test invalid name", func(t *testing.T) {
		valid := valid
		valid.DBName = ""
		assert.EqualError(t, valid.Validate(), "invalid database name")
	})
}
