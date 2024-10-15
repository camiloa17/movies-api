package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func (auth *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	// get auth header
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// get bearer token
	headersParts := strings.Split(authHeader, " ")

	if len(headersParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	//check if Bearer is present.
	if headersParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	maybeToken := headersParts[1]

	// empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(maybeToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.Secret), nil
	})

	// we check for errors of validation
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	if claims.Issuer != auth.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return maybeToken, claims, nil

}
