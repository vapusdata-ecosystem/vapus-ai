package dmutils

import (
	"log"

	"github.com/blang/semver/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

func GetVersionNumber(current, bumptType string) string {
	if current == "" {
		v, err := semver.Make("1.0.0")
		if err != nil {
			log.Println("error while creating semver", err)
			return ""
		}
		return v.String()

	}
	v, _ := semver.Make(current)

	switch bumptType {
	case mpb.VersionBumpType_MAJOR.String():
		v.Major++
		v.Minor = 0
		v.Patch = 0
		return v.String()
	case mpb.VersionBumpType_MINOR.String():
		v.Minor++
		v.Patch = 0
		return v.String()
	case mpb.VersionBumpType_PATCH.String():
		v.Patch++
		return v.String()
	default:
		return current
	}
}

func SetHelmChartVersion(current, bumptType string) string {
	if current == "" {
		v, err := semver.Make("0.0.1")
		if err != nil {
			log.Println("error while creating semver", err)
			return ""
		}
		return v.String()

	}
	v, _ := semver.Make(current)

	switch bumptType {
	case mpb.VersionBumpType_MAJOR.String():
		v.Major++
		v.Minor = 0
		v.Patch = 0
		return v.String()
	case mpb.VersionBumpType_MINOR.String():
		v.Minor++
		v.Patch = 0
		return v.String()
	case mpb.VersionBumpType_PATCH.String():
		v.Patch++
		return v.String()
	default:
		return current
	}
}
