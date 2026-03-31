package jwt

import (
	"auth-service/internal/infra/config"
	"libs/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenManager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	log        *logger.Logger
}

var TokenMethod = jwt.SigningMethodHS256

func NewTokenManager(cfg *config.Config, log *logger.Logger) *TokenManager {
	tokenLog := log.Named("TokenManager")

	return &TokenManager{
		secret:     []byte(cfg.JWT.Secret),
		issuer:     cfg.JWT.Issuer,
		accessTTL:  cfg.JWT.AccessTokenTTL,
		refreshTTL: cfg.JWT.RefreshTokenTTL,
		log:        tokenLog,
	}
}

func (tm *TokenManager) GenerateBothTokens(payload TokenPayload) (*TokenPair, error) {
	accessToken, err := tm.GenerateToken(payload, Access)
	if err != nil {
		return nil, err
	}

	refreshToken, err := tm.GenerateToken(payload, Refresh)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tm *TokenManager) GenerateToken(
	payload TokenPayload,
	tokenType TokenType,
) (string, error) {
	now := time.Now()
	var tokenExp time.Time

	switch tokenType {
	case Access:
		tokenExp = now.Add(tm.accessTTL)
	case Refresh:
		tokenExp = now.Add(tm.refreshTTL)
	}

	claims := TokenClaims{
		TokenPayload: payload,
		Type:         tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   payload.UserId.String(),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token, err := jwt.NewWithClaims(TokenMethod, claims).SignedString(tm.secret)
	if err != nil {
		tm.log.Error(FailedToGenerateToken.Error(), zap.Error(err))
		return "", FailedToGenerateToken
	}

	return token, nil
}

func (tm *TokenManager) ParseTokenPayload(
	tokenString string,
	expectedType TokenType,
) (*TokenPayload, *jwt.NumericDate, error) {
	claims := TokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != TokenMethod {
				return nil, InvalidToken
			}

			return tm.secret, nil
		},
		jwt.WithIssuer(tm.issuer),
		jwt.WithValidMethods([]string{TokenMethod.Alg()}),
	)

	if err != nil {
		tm.log.Error(InvalidToken.Error(), zap.Error(err))
		return nil, &jwt.NumericDate{}, InvalidToken
	}

	if !token.Valid {
		return nil, &jwt.NumericDate{}, InvalidToken
	}

	if claims.Type != expectedType {
		return nil, &jwt.NumericDate{}, InvalidTokenType
	}

	return &claims.TokenPayload, claims.ExpiresAt, nil
}
