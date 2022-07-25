package env

import "os"

var env string

func init() {
	env = os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
}

func GetEnv() string {
	return env
}
