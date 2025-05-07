package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	jwt "github.com/golang-jwt/jwt/v5"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

var DefaultECDSElliptic string = "P-521"

type ECDSAJwt interface {
	GenerateVDPAJWT(claims *VapusDataPlatformAccessClaims) (string, error)
	GenerateVDPARefreshJWT(claims *VapusDataPlatformRefreshTokenClaims) (string, error)
	ParseAndValidateVDPAJWT(tokenString string) (*VapusDataPlatformAccessClaims, error)
	ValidateAccessToken(tokenString string) (map[string]string, error)
}

type ECDSAKeys struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
	EllipticCurve elliptic.Curve
}

type ECDSAManager struct {
	opts        *JWTAuthn
	ParsedPvKey *ecdsa.PrivateKey
	ParsedPbKey *ecdsa.PublicKey
}

var ecdsaSigningAlgoMap = map[string]*jwt.SigningMethodECDSA{
	"P-521": jwt.SigningMethodES512,
	"P-384": jwt.SigningMethodES384,
	"P-256": jwt.SigningMethodES256,
}

var ellipticCurveMap = map[string]elliptic.Curve{
	"P-256": elliptic.P256(),
	"P-384": elliptic.P384(),
	"P-521": elliptic.P521(),
}

func GenerateECDSAKeys(curve string) (*ECDSAKeys, error) {
	eCurve := ellipticCurveMap[curve]
	privKey, err := ecdsa.GenerateKey(eCurve, rand.Reader)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error generating ECDSA private key with elliptic curve %v", curve)
		return nil, err
	}
	return &ECDSAKeys{
		PrivateKey:    privKey,
		PublicKey:     &privKey.PublicKey,
		EllipticCurve: eCurve,
	}, nil
}

// NewECDSAJwtAuthn creates a new ECDSA JWT Authn object with the given options.
// It returns the ECDSAJwt interface. It logs an error if the private key is not parsed.
func NewECDSAJwtAuthn(opts *JWTAuthn) (ECDSAJwt, error) {
	res := &ECDSAManager{
		opts: opts,
	}
	if opts.ForPublicValidation {
		dmlogger.CoreLogger.Info().Msg("Using public key for validation")
		parsedPbKey, err := jwt.ParseECPublicKeyFromPEM([]byte(opts.PublicJWTKey))
		if err != nil || parsedPbKey == nil {
			dmlogger.CoreLogger.Err(err).Msg("Error parsing ECDSA public key")
			return nil, err
		}
		res.ParsedPbKey = parsedPbKey
	} else {
		parsedPvKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(opts.PrivateJWTKey))
		if err != nil || parsedPvKey == nil {
			dmlogger.CoreLogger.Err(err).Msg("Error parsing ECDSA private key")
			return nil, err
		}
		res.ParsedPvKey = parsedPvKey
		res.ParsedPbKey = &parsedPvKey.PublicKey
	}
	return res, nil
}

func (e *ECDSAManager) GenerateVDPAJWT(claims *VapusDataPlatformAccessClaims) (string, error) {
	if e.opts.ForPublicValidation {
		return types.EMPTYSTR, dmerrors.DMError(ErrOnlyPublicJWTKey, nil)
	}
	dmlogger.CoreLogger.Info().Msgf("Generating JWT token for claim %v", claims)
	token := jwt.NewWithClaims(ecdsaSigningAlgoMap[e.ParsedPvKey.Curve.Params().Name], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return types.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *ECDSAManager) GenerateVDPARefreshJWT(claims *VapusDataPlatformRefreshTokenClaims) (string, error) {
	if e.opts.ForPublicValidation {
		return types.EMPTYSTR, dmerrors.DMError(ErrOnlyPublicJWTKey, nil)
	}

	dmlogger.CoreLogger.Info().Msgf("Generating refresh token for claim %v", claims)
	token := jwt.NewWithClaims(ecdsaSigningAlgoMap[e.ParsedPvKey.Curve.Params().Name], claims)

	tokenString, err := token.SignedString(e.ParsedPvKey)
	if err != nil {
		return types.EMPTYSTR, err
	}
	return tokenString, nil
}

func (e *ECDSAManager) ParseAndValidateVDPAJWT(tokenString string) (*VapusDataPlatformAccessClaims, error) {
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

/*
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
	    // token is valid
	} else {

	    // check for clock skew
	    parser := jwt.NewParser(jwt.WithLeeway(5 * time.Second)) // Allow 5 seconds of leeway
	    _, err = parser.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	        return publicKey, nil
	    })

	    if err != nil {
	        log.Printf("Token validation failed even with leeway: %v", err)
	    }
	}
*/

func (e *ECDSAManager) ValidateAccessToken(tokenString string) (map[string]string, error) {
	claim, err := e.ParseAndValidateVDPAJWT(tokenString)
	if err != nil {
		dmlogger.CoreLogger.Err(err).Msgf("error while parsing and validating auth token")
		return nil, err
	}
	dmlogger.CoreLogger.Info().Msgf("parsed ORGANIZATION claims - %v", claim)
	resp := FlatVDPAScopeClaims(claim, "||")
	if resp == nil {
		dmlogger.CoreLogger.Error().Msgf("invalid Claim parsed from the token")
		return nil, dmerrors.DMError(ErrInvalidUserAuthentication, nil)
	}
	return resp, nil
}
