package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

func GetEnv(key string) (*string, error) {
	envVar, isExists := os.LookupEnv(key)
	if !isExists {
		return nil, fmt.Errorf("the .env variable, %s, is not present", key)
	}

	return &envVar, nil
}
