package constants

import "os"

// OSEnv returns environment variables
func OSEnv(s string) string {
	env := os.Getenv(s)
	return env
}
