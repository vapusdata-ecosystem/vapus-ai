package encryption

import (
	"crypto/rand"
	"crypto/rsa"

	jwt "github.com/golang-jwt/jwt/v5"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type RSAJwt interface {
	GenerateVDPAJWT(claims *VapusDataPlatformAccessClaims) (string, error)
	GenerateVDPARefreshJWT(claims *VapusDataPlatformRefreshTokenClaims) (string, error)
	ParseAndValidateVDPAJWT(tokenString string) (*VapusDataPlatformAccessClaims, error)
	ValidateAccessToken(tokenString string) (map[string]string, error)
}

type RSAKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Bits       int
}

type RSAManager struct {
	opts        *JWTAuthn
	ParsedPvKey *rsa.PrivateKey
	ParsedPbKey *rsa.PublicKey
}

var rsaSigningAlgoMap = map[string]*jwt.SigningMethodRSA{
	"P-521": jwt.SigningMethodRS512,
	"P-384": jwt.SigningMethodRS384,
	"P-256": jwt.SigningMethodRS256,
}

func GenerateRSAKeys(bits int) (*RSAKeys, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error generating ECDSA private key with bits %v", bits)
		return nil, err
	}
	return &RSAKeys{
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
		Bits:       bits,
	}, nil
}

// NewRSAJwtAuthn creates a new RSA JWT Authn object with the given options.
// It returns the RSAJwt interface. It logs an error if the private key is not parsed.
func NewRSAJwtAuthn(opts *JWTAuthn) (RSAJwt, error) {
	res := &RSAManager{
		opts: opts,
	}
	if opts.ForPublicValidation {
		dmlogger.CoreLogger.Info().Msg("Using public key for validation")
		parsedPbKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(opts.PublicJWTKey))
		if err != nil || parsedPbKey == nil {
			dmlogger.CoreLogger.Err(err).Msg("Error parsing ECDSA public key")
			return nil, err
		}
		res.ParsedPbKey = parsedPbKey
	} else {
		parsedPvKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(opts.PrivateJWTKey))
		if err != nil || parsedPvKey == nil {
			dmlogger.CoreLogger.Err(err).Msg("Error parsing RSA private key")
			return nil, err
		}
		res.ParsedPvKey = parsedPvKey
		res.ParsedPbKey = &parsedPvKey.PublicKey
	}
	return res, nil
}

func (e *RSAManager) GenerateVDPAJWT(claims *VapusDataPlatformAccessClaims) (string, error) {
	if e.opts.ForPublicValidation {
		return types.EMPTYSTR, dmerrors.DMError(ErrOnlyPublicJWTKey, nil)
	}
	dmlogger.CoreLogger.Info().Msgf("Generating JWT token for claim %v", claims)
	token := jwt.NewWithClaims(rsaSigningAlgoMap[e.ParsedPvKey.N.String()], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return types.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *RSAManager) GenerateVDPARefreshJWT(claims *VapusDataPlatformRefreshTokenClaims) (string, error) {
	if e.opts.ForPublicValidation {
		return types.EMPTYSTR, dmerrors.DMError(ErrOnlyPublicJWTKey, nil)
	}

	dmlogger.CoreLogger.Info().Msgf("Generating refresh token for claim %v", claims)
	token := jwt.NewWithClaims(rsaSigningAlgoMap[e.ParsedPvKey.N.String()], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return types.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *RSAManager) ParseAndValidateVDPAJWT(tokenString string) (*VapusDataPlatformAccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &VapusDataPlatformAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return e.ParsedPbKey, nil
	})
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msg(ErrParsingJWT.Error())
		return nil, dmerrors.DMError(ErrParsingJWT, err)
	} else if !token.Valid {
		dmlogger.CoreLogger.Err(err).Msg("Invalid JWT token")
		return nil, dmerrors.DMError(ErrInvalidJWT, nil)
	}

	if claims, ok := token.Claims.(*VapusDataPlatformAccessClaims); !ok {
		dmlogger.CoreLogger.Err(ErrInvalidJWTClaims).Msg("Invalid JWT claims")
		return nil, dmerrors.DMError(ErrInvalidJWTClaims, nil)
	} else {
		return claims, nil
	}
}

func (e *RSAManager) ValidateAccessToken(tokenString string) (map[string]string, error) {
	claim, err := e.ParseAndValidateVDPAJWT(tokenString)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error while parsing and validating auth token")
		return nil, err
	}
	resp := FlatVDPAScopeClaims(claim, "||")
	if resp == nil {
		dmlogger.CoreLogger.Error().Msgf("invalid Claim parsed from the token")
		return nil, dmerrors.DMError(ErrInvalidUserAuthentication, nil)
	}
	dmlogger.CoreLogger.Info().Msgf("flatClaims - %v", resp)
	return resp, nil
}
