package models

import mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

type TlsConfig struct {
	TlsType        string `json:"tlsType" yaml:"tlsType"`
	CaCertFile     string `json:"caCertFile" yaml:"caCertFile"`
	ServerKeyFile  string `json:"serverKeyFile" yaml:"serverKeyFile"`
	ServerCertFile string `json:"serverCertFile" yaml:"serverCertFile"`
}

func (j *TlsConfig) ConvertToPb() *mpb.TlsConfig {
	if j != nil {
		return &mpb.TlsConfig{
			TlsType:        mpb.TlsType(mpb.TlsType_value[j.TlsType]),
			CaCertFile:     j.CaCertFile,
			ServerKeyFile:  j.ServerKeyFile,
			ServerCertFile: j.ServerCertFile,
		}
	}
	return nil
}

func (j *TlsConfig) ConvertFromPb(pb *mpb.TlsConfig) *TlsConfig {
	if pb == nil {
		return nil
	}
	return &TlsConfig{
		TlsType:        mpb.TlsType_name[int32(pb.GetTlsType())],
		CaCertFile:     pb.GetCaCertFile(),
		ServerKeyFile:  pb.GetServerKeyFile(),
		ServerCertFile: pb.GetServerCertFile(),
	}
}

type PlatformArtifact struct {
	Artifact string `json:"artifact" yaml:"artifact"`
	Tag      string `json:"tag" yaml:"tag"`
	Digest   string `json:"digest" yaml:"digest"`
	IsLatest bool   `json:"isLatest" yaml:"isLatest"`
	AddedOn  int64  `json:"addedOn" yaml:"addedOn"`
}

func (j *PlatformArtifact) ConvertToPb() *mpb.PlatformArtifact {
	if j != nil {
		return &mpb.PlatformArtifact{
			Artifact: j.Artifact,
			Tag:      j.Tag,
			Digest:   j.Digest,
			IsLatest: j.IsLatest,
			AddedOn:  j.AddedOn,
		}
	}
	return nil
}

func (j *PlatformArtifact) ConvertFromPb(pb *mpb.PlatformArtifact) *PlatformArtifact {
	if pb != nil {
		return &PlatformArtifact{
			Artifact: pb.GetArtifact(),
			Tag:      pb.GetTag(),
			Digest:   pb.GetDigest(),
			IsLatest: pb.GetIsLatest(),
			AddedOn:  pb.GetAddedOn(),
		}
	}
	return nil
}

func (pa *PlatformArtifact) GetArtifact() string {
	if pa != nil {
		return pa.Artifact
	}
	return ""
}

func (pa *PlatformArtifact) GetTag() string {
	if pa != nil {
		return pa.Tag
	}
	return ""
}

func (pa *PlatformArtifact) GetDigest() string {
	if pa != nil {
		return pa.Digest
	}
	return ""
}

func (pa *PlatformArtifact) GetIsLatest() bool {
	if pa != nil {
		return pa.IsLatest
	}
	return false
}

func (pa *PlatformArtifact) GetAddedOn() int64 {
	if pa != nil {
		return pa.AddedOn
	}
	return 0
}
