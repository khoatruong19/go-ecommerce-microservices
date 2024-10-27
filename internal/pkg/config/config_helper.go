package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v11"
	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/config/environment"
	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/constants"
	typeMapper "github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/reflection/typemapper"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func BindConfigKey[T any](configKey string, environments ...environment.Environment) (T, error) {
	var configPath string

	currentEnv := environment.Environment("")
	if len(environments) > 0 {
		currentEnv = environments[0]
	} else {
		currentEnv = constants.Dev
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	defaults.SetDefaults(cfg)

	configPathFromEnv := viper.GetString(constants.ConfigPath)
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else {
		appRootPath := viper.GetString(constants.AppRootPath)
		if appRootPath == "" {
			appRootPath = environment.GetProjectRootWorkingDirectory()
		}

		dir, err := searchForConfigFileDir(appRootPath, currentEnv)
		if err != nil {
			return *new(T), err
		}

		configPath = dir
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", currentEnv))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.JSON)

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), errors.WrapIf(err, "viper.ReadInConfig")
	}

	if len(configKey) == 0 {
		// load configs from config file to config object
		if err := viper.Unmarshal(cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	} else {
		if err := viper.UnmarshalKey(configKey, cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	}

	viper.AutomaticEnv()

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg, nil

}

func searchForConfigFileDir(rootDir string, env environment.Environment) (string, error) {
	var dirPath string

	err := filepath.WalkDir(
		rootDir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Check if the file is named "config.%s.json" (replace %s with the env)
			if !info.IsDir() &&
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.json", env),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yaml", env),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yml", env),
				) {
				dir := filepath.Dir(path)
				dirPath = dir

				return filepath.SkipDir // Skip further traversal
			}

			return nil
		},
	)

	if dirPath != "" {
		return dirPath, nil
	}

	return "", errors.WrapIf(err, "No directory with config file found")
}
