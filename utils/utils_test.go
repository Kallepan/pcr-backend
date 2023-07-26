package utils

import (
	"os"
	"testing"
)

var mock_env = map[string]string{
	"POSTGRES_USER":     "test",
	"POSTGRES_PASSWORD": "test",
	"POSTGRES_DB":       "test",
	"POSTGRES_PORT":     "5432",
	"POSTGRES_HOST":     "localhost",
	"expected":          "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
}

var mock_env_empty = map[string]string{
	"expected": "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
}

var testEnvs = []map[string]string{
	mock_env,
	mock_env_empty,
}

func TestGetDBConnectionString(t *testing.T) {
	// Mock environment variables
	for _, env := range testEnvs {
		os.Clearenv()
		for key, value := range env {
			os.Setenv(key, value)
		}

		// Test the function
		actual := GetDBConnectionString()
		expected := env["expected"]
		if actual != expected {
			t.Errorf("GetDBConnectionString() failed, expected %s, got %s", expected, actual)
		}
	}
}

type envTest struct {
	envSrc      string
	defaultVal  string
	expectedVal string
}

var envTests = []envTest{
	{"POSTGRES_USER", "test", "test"},
	{"POSTGRES_PASSWORD", "test", "test"},
	{"POSTGRES_DB", "test", "test"},
	{"test", "test", "test"},
	{"empty", "", ""},
}

func TestGetValueFromEnv(t *testing.T) {
	// Mock environment variables
	for _, env := range envTests {
		os.Clearenv()
		os.Setenv(env.envSrc, env.expectedVal)

		// Test the function
		actual := GetValueFromEnv(env.envSrc, env.defaultVal)
		expected := env.expectedVal
		if actual != expected {
			t.Errorf("GetValueFromEnv() failed, expected %s, got %s", expected, actual)
		}

	}
}
