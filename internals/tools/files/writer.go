package filetools

import (
	"archive/tar"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	os "os"
	filepath "path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"gopkg.in/yaml.v3"
)

func WriteTomlFile(data interface{}, filename, path string) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return err
	}

	file := filepath.Join(path, filename+types.DOT+types.DEFAULT_CONFIG_TYPE)
	dmlogger.CoreLogger.Info().Msgf("Writing to file: %v", file)
	err = os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(filename, filePath string, data any, base64encoded bool) error {
	if filePath == "" {
		curPath, err := os.Getwd()
		if err != nil {
			curPath = os.TempDir()
		}
		filePath = curPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return err
		}
	}
	fType := GetConfFileType(filename)
	switch strings.ToLower(fType) {
	case "yaml":
		return WriteYAMLFile(data, filepath.Join(filePath, filename), base64encoded)
	case "json":
		return WriteJSONFile(data, filepath.Join(filePath, filename), base64encoded)
	case "toml":
		return WriteTOMLFile(data, filepath.Join(filePath, filename), base64encoded)
	default:
		return dmerrors.ErrInvalidArgs
	}
}

func CreateTarFile(tarFile string, files2Add []string, fileDest string) error {
	tarFileHandle, err := os.Create(tarFile)
	if err != nil {
		return err
	}
	defer tarFileHandle.Close()

	tw := tar.NewWriter(tarFileHandle)
	defer tw.Close()

	for _, fl := range files2Add {
		file, err := os.Open(fl)
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: filepath.Join(fileDest, filepath.Base(file.Name())),
			Mode: 0644,
			Size: fileInfo.Size(),
		}
		log.Println("File to be added to tar - ", filepath.Join(fileDest, filepath.Base(file.Name())))
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
	}

	return err
}

func WriteYAMLFile[T any](data T, fileName string, base64encoded bool) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	log.Println("Writing to file - ", fileName, " data - ", string(bytes))
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteJSONFile[T any](data T, fileName string, base64encoded bool) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteTOMLFile[T any](data T, fileName string, base64encoded bool) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	if base64encoded {
		bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}
