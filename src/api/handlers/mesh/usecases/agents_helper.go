package usecases

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/golang-jwt/jwt"
)

func (u *agentsUsecase) processKey(r *http.Request, key models.Key, encoded, hostname string) (*http.Cookie, error) {
	token, err := u.decodeKey(key, encoded)
	if err != nil {
		return nil, err
	}
	if err := u.validateClaims(token, hostname); err != nil {
		return nil, err
	}
	return u.setKeyAsCookie(r, encoded), nil
}

func (u *agentsUsecase) decodeKey(key models.Key, encoded string) (*jwt.Token, error) {
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(key.PublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(encoded, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return parsedKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *agentsUsecase) validateClaims(token *jwt.Token, hostname string) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("jwt token validatation failed")
	}
	if scope, ok := claims["scope"]; !ok || scope.(string) != "agent" {
		return errors.New("jwt scope payload validatation failed")
	}
	if host, ok := claims["host"]; !ok || host.(string) != hostname {
		return errors.New("jwt host payload validatation failed")
	}
	return nil
}

func (u *agentsUsecase) setKeyAsCookie(r *http.Request, encoded string) *http.Cookie {
	return &http.Cookie{
		Name:     "_token",
		Path:     "/",
		Domain:   r.URL.Host,
		HttpOnly: true,
		Secure:   true,
		Value:    encoded,
	}
}
