package pbac

import (
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

type Roles struct {
	Name     string   `yaml:"name"`
	Policies []string `yaml:"policies"`
}

type PbacConfig struct {
	PlatformPolicies     []string `yaml:"platformPolicies"`
	ORGANIZATIONPolicies []string `yaml:"ORGANIZATIONPolicies"`
	Roles                []Roles  `yaml:"roles"`
}

func LoadPbacConfig(file string) (*PbacConfig, error) {
	pConf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(file), file, &PbacConfig{})
	if err != nil {
		return nil, err
	}
	return pConf.(*PbacConfig), nil

}
