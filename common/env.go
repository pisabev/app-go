package common

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort      int
	HttpPortDebug int
}

func DotEnv() {
	if err := godotenv.Load(".env"); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}
}
