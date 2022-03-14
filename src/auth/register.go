package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

// this file handles agent authentication during the Register process
// in this process, the agent encodes a JWT token using its RSA private key

// AgentRegisterKey represents key used to register the agent
type AgentRegisterKey struct {
	Name       string // key name set in the mesh API
	PrivateKey []byte // private key
}

// EncodeName encodes the key name using the private key
func (k AgentRegisterKey) EncodeName(hostname string) (string, error) {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(k.PrivateKey)
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"scope":    "agent",
		"hostname": hostname,
	}).SignedString(parsedKey)
}

// MeshRegisterKey represents the mesh register key used to decode the register JWT token
type MeshRegisterKey struct {
	Name      string // key name set in the mesh API
	PublicKey []byte // public key
}

// Decode decodes the register JWT token
func (k MeshRegisterKey) Decode(encodedToken string) (*jwt.Token, error) {
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(k.PublicKey)
	if err != nil {
		return nil, err
	}
	return jwt.Parse(encodedToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return parsedKey, nil
	})
}

// ValidateRegisterToken validates the register token claims
func ValidateRegisterToken(token *jwt.Token, hostname string) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("failed to get jwt claims")
	}
	if scope, ok := claims["scope"]; !ok || scope.(string) != "agent" {
		return errors.New("jwt scope payload validatation failed")
	}
	if host, ok := claims["hostname"]; !ok || host.(string) != hostname {
		return errors.New("jwt host payload validatation failed")
	}
	return nil
}
