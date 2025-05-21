package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

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
	GenerateKeys(bits int) (string, string, error)
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

func (e *RSAManager) GenerateKeys(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error generating ECDSA private key with bits %v", bits)
		return "", "", err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPath := "private_key.pem"
	publicKeyPath := "public_key.pem"
	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", fmt.Errorf("failed to open private key file %s for writing: %w", privateKeyPath, err)
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error writing private key to %s", privateKeyPath)
		return "", "", fmt.Errorf("failed to write private key to %s: %w", privateKeyPath, err)
	}

	// --- Derive Public Key and Encode to PEM format (PKIX/SPKI) ---
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error marshalling public key")
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error creating public key file %s", publicKeyPath)
		return "", "", fmt.Errorf("failed to open public key file %s for writing: %w", publicKeyPath, err)
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error writing public key to %s", publicKeyPath)
		return "", "", fmt.Errorf("failed to write public key to %s: %w", publicKeyPath, err)
	}
	pvContent, err := io.ReadAll(privateKeyFile)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error reading private key file %s", privateKeyPath)
		return "", "", fmt.Errorf("failed to read private key file %s: %w", privateKeyPath, err)
	}
	pbContent, err := io.ReadAll(publicKeyFile)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error reading public key file %s", publicKeyPath)
		return "", "", fmt.Errorf("failed to read public key file %s: %w", publicKeyPath, err)
	}
	return string(pvContent), string(pbContent), nil
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
