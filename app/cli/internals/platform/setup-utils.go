package plclient

import (
	"log"
	"os"

	guuid "github.com/google/uuid"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
	dmodels "github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/authn"
	encryptions "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

func GetDataSourceParams(filename string) (*dmodels.DataSourceCredsParams, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &dmodels.DataSourceCredsParams{}
	err = filetools.GenericUnMarshaler(bytes, obj, filetools.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func GetAuthnParams(filename string) (*authn.AuthnSecrets, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &authn.AuthnSecrets{}
	err = filetools.GenericUnMarshaler(bytes, obj, filetools.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	log.Println("AuthnSecrets: ", obj)
	return obj, nil
}

func GetJwtAuthnParams(filename string) (*encryptions.JWTAuthn, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &encryptions.JWTAuthn{}
	err = filetools.GenericUnMarshaler(bytes, obj, filetools.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func GetSecretName(secretName string) string {
	if secretName == "" {
		return guuid.New().String()[:5] + "-secret"
	}
	return secretName
}
