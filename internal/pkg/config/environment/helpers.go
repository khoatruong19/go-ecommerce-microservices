package environment

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"emperror.dev/errors"

	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/spf13/viper"
)

func GetProjectRootWorkingDirectory() string {
	var rootWorkingDirectory string

	pn := viper.GetString(constants.ProjectName)
	if pn != "" {
		rootWorkingDirectory = getProjectRootDirectoryFromProjectName(pn)
	} else {
		wd, _ := os.Getwd()
		dir, err := searchRootDirectory(wd)
		if err != nil {
			log.Fatal(err)
		}
		rootWorkingDirectory = dir
	}

	absoluteRootWorkingDirectory, _ := filepath.Abs(rootWorkingDirectory)

	return absoluteRootWorkingDirectory
}

func getProjectRootDirectoryFromProjectName(pn string) string {
	wd, _ := os.Getwd()

	for !strings.HasSuffix(wd, pn) {
		wd = filepath.Dir(wd)
	}

	return wd
}

func searchRootDirectory(
	dir string,
) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.WrapIf(err, "Error reading directory")
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.EqualFold(fileName, "go.mod") {
				return dir, nil
			}
		}
	}

	// If no config file found in this directory, recursively search its parent
	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		// We've reached the root directory, and no go.mod file was found
		return "", errors.WrapIf(err, "No go.mod file found")
	}

	return searchRootDirectory(parentDir)
}
