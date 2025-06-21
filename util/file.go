package util

import (
	"os"
)

func GetConfigFileContents() ([]byte, error) {
	data, err := os.ReadFile("../config.json")
	return data, err
}
