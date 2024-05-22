package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TokenPairs struct {
	Token        string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func (auth *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenPairs{}, errors.New("could not assert claims from the token")
	}
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = auth.Audience
	claims["iss"] = auth.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	// Set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(auth.TokenExpiry).Unix()
	// Create a signed Token
	signedAccessToken, err := token.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPairs{}, err
	}
	// Create a refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return TokenPairs{}, errors.New("could not assert claims for refresh token")
	}
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Set expiry for the refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(auth.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create token Pairs and populate with signed tokens.
	pairs := TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return pairs, nil
}

func (auth *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    refreshToken,
		Domain:   auth.CookieDomain,
		Expires:  time.Now().Add(auth.RefreshExpiry),
		MaxAge:   int(auth.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

func (auth *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    "",
		Domain:   auth.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}
