package usecases

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/golang-jwt/jwt"
)

// MeshGlobalKey is the env variable name where the global key is stored
const MeshGlobalKey = "MESH_GLOBAL_KEY"

func (u *agentsUsecase) processKey(r *http.Request, key models.Key, encoded, hostname string) (*http.Cookie, error) {
	token, err := u.decodeKey(key, encoded)
	if err != nil {
		return nil, err
	}
	if err := u.validateClaims(token, hostname); err != nil {
		return nil, err
	}

	newToken, err := u.newJWTToken(hostname)
	if err != nil {
		return nil, err
	}
	return u.setCookie(r, newToken), nil
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
	if !ok {
		return errors.New("failed to get jwt claims")
	}
	if scope, ok := claims["scope"]; !ok || scope.(string) != "agent" {
		return errors.New("jwt scope payload validatation failed")
	}
	if host, ok := claims["host"]; !ok || host.(string) != hostname {
		return errors.New("jwt host payload validatation failed")
	}
	return nil
}

func (u *agentsUsecase) newJWTToken(hostname string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["hostname"] = hostname
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(MeshGlobalKey)))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *agentsUsecase) setCookie(r *http.Request, token string) *http.Cookie {
	return &http.Cookie{
		Name:     "_token",
		Path:     "/",
		Domain:   r.URL.Host,
		HttpOnly: true,
		Secure:   true,
		Value:    token,
	}
}
