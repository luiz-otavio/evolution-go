package evolution

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	if flag.Lookup("try-login") == nil {
		flag.Bool("try-login", false, "Try to login with valid credentials")
	}
}

func ShouldTryLogin() bool {
	if !flag.Parsed() {
		flag.Parse()
	}

	return flag.Lookup("try-login").
		Value.
		String() == "true"
}

func findDotenv() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New(".env file not found")
		}
		dir = parent
	}
}

func ReadLocalDotenv() error {
	path, err := findDotenv()
	if err != nil {
		return err
	}

	envFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read .env file: %w", err)
	}

	lines := strings.Split(string(envFile), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s: %w", key, err)
		}
	}

	return nil
}

func RequireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required env variable %s is empty", key)
	}
	return value, nil
}
