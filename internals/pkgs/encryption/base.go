package encryption

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	"google.golang.org/grpc/metadata"
)

var encryptLogger zerolog.Logger

// TO:DO Make it more generic to handle different type of claims
type JwtAuthService interface {
	GenerateVDPAJWT(claims *VapusDataPlatformAccessClaims) (string, error)
	GenerateVDPARefreshJWT(claims *VapusDataPlatformRefreshTokenClaims) (string, error)
	ParseAndValidateVDPAJWT(tokenString string) (*VapusDataPlatformAccessClaims, error)
	ValidateAccessToken(tokenString string) (map[string]string, error)
	GenerateKeys(bits int) (string, string, error)
}

type JWTAuthn struct {
	PublicJWTKey        string `validate:"required" yaml:"publicJwtKey" json:"publicJwtKey"`
	PrivateJWTKey       string `yaml:"privateJwtKey" json:"privateJwtKey"`
	SigningAlgorithm    string `validate:"required" yaml:"signingAlgorithm" json:"signingAlgorithm"`
	ForPublicValidation bool   `default:"false" yaml:"forPublicValidation" json:"forPublicValidation"`
	TokenIssuer         string `yaml:"tokenIssuer" json:"tokenIssuer"`
	TokenAudience       string `yaml:"tokenAudience" json:"tokenAudience"`
	Bitsize             int    `default:"2048" yaml:"bitsize" json:"bitsize"`
}

type jwtAuthOpts func(jo *JWTAuthn)

type VapusDataJwtAuthn struct {
	Opts *JWTAuthn
	JwtAuthService
}

var JwtTokenIssuer = "vapusai"
var JwtTokenAudience = "*.vapusdata.com"
var VapusPlatformTokenSubject = "VapusData access token"
var JwtOrganizationScope = "OrganizationScope"
var JwtPlatformScope = "platformScope"
var JwtDataProductScope = "dataProductScope"
var JwtCtxClaimKey = "vapusPlatformJwtClaim"
var JwtDPCtxClaimKey = "vapusPlatformJwtClaim"
var JwtClaimRoleSeparator = "|"
var JWTParser *jwt.Parser

func SetCtxClaim(ctx context.Context, claim map[string]string) context.Context {
	return context.WithValue(ctx, JwtCtxClaimKey, claim)
}

func SetDataProductCtxClaim(ctx context.Context, claim map[string]string) context.Context {
	return context.WithValue(ctx, JwtDPCtxClaimKey, claim)
}

// Checking if the user is authorized or not
func GetCtxClaim(ctx context.Context) (map[string]string, bool) {
	// Here we are trying to find if the JWT is stored in the Context or not
	val, ok := ctx.Value(JwtCtxClaimKey).(map[string]string)
	if !ok {

		// GRPC MetaData is a set of key value pairs that are sent with the request and response.
		// Its like a HTTP header, usefule for sending extra information which are not the part of request/response.
		// There are two types of metadata..
		// Incoming MetaData and Outgoing MetaData
		// Incoming MetaData: MetaData recieved either from the client or server (FromIncomingContext).
		// Outgoing MetaData: MetaData sent either from the server or client

		// md := metadata.New(map[string]string{
		// 		"authorization": "Bearer my-token"
		// })
		// ctx:= metadata.NewOutgoingContext(context.Background(), md)

		// If JWT was not stored in the Context, then we are checking the GRPC MetaData
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ok
		}
		val = make(map[string]string)
		if len(md.Get(JwtDPCtxClaimKey)) != 1 {
			return nil, false
		}
		strval := md.Get(JwtDPCtxClaimKey)[0]

		// Unmarshalling the JSON data and storing it in val(map[string]string)
		err := json.Unmarshal([]byte(strval), &val)
		if err != nil {
			return nil, false
		}
		return val, true

	}
	return val, ok
}

func GetDPtxClaim(ctx context.Context) (map[string]string, bool) {
	val, ok := ctx.Value(JwtDPCtxClaimKey).(map[string]string)
	return val, ok
}

func NewVapusDataJwtAuthnWithConfig(path string) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	jwtAuthnSecrets, err := LoadJwtAuthnSecrets(path)
	if err != nil {
		return nil, err
	}
	return NewVapusDataJwtAuthn(jwtAuthnSecrets)
}

func LoadJwtAuthnSecrets(path string) (*JWTAuthn, error) {
	cf, err := dmutils.ReadBasicConfig(filetools.GetConfFileType(path), path, &JWTAuthn{})
	if err != nil {
		encryptLogger.Info().Msgf("Error loading jwt authn secrets: %v", err)
		return nil, err
	}
	return cf.(*JWTAuthn), err
}

func NewVapusDataJwtAuthn(opts *JWTAuthn) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	obj := &VapusDataJwtAuthn{
		Opts: opts,
	}
	JWTParser = jwt.NewParser(jwt.WithLeeway(2 * time.Second))
	switch opts.SigningAlgorithm {
	case mpb.EncryptionAlgo_ECDSA.String():
		val, err := NewECDSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val.(*ECDSAManager)
		return obj, nil
	case mpb.EncryptionAlgo_RSA.String():
		val, err := NewRSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val.(*RSAManager)
		return obj, nil
	default:
		return nil, ErrInvalidJWT
	}

}
