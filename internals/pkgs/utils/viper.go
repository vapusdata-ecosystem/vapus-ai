package dmutils

import (
	filepath "path/filepath"

	viper "github.com/spf13/viper"
)

// ReadBasicConfig function to read basic config using viper library
func ReadBasicConfig(fileType, path string, obj interface{}) (interface{}, error) {
	// Set the file path and file name when whole file path is provided in the name variable
	path, fileName := filepath.Split(path)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)
	viper.SetConfigType(fileType)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := viper.Unmarshal(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
