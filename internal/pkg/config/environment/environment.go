package environment

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/spf13/viper"
)

type Environment string

var (
	Development = Environment(constants.Dev)
	Test        = Environment(constants.Test)
	Production  = Environment(constants.Production)
)

func ConfigAppEnv(environments ...Environment) Environment {
	environment := Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = Development
	}

	// setup viper to read from os environment with `viper.Get`
	viper.AutomaticEnv()

	err := loadEnvFilesRecursive()
	if err != nil {
		log.Printf("fail to load .env file, err: %v", err)
	}

	setRootWorkingDirectoryEnvironment()

	manualEnv := os.Getenv(constants.AppRootPath)
	if manualEnv != "" {
		environment = Environment(manualEnv)
	}

	return environment
}

func loadEnvFilesRecursive() error {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envFilePath := filepath.Join(dir, ".env")
		err := godotenv.Load(envFilePath)

		// .env file found and loaded
		if err == nil {
			return nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return errors.New(".env file not found in the project hierarchy")
}

func setRootWorkingDirectoryEnvironment() {
	absoluteRootWorkingDirectory := GetProjectRootWorkingDirectory()

	viper.Set(constants.AppRootPath, absoluteRootWorkingDirectory)
}