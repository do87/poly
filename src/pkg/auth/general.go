package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/golang-jwt/jwt"
)

// this file handles agent authentication after the Register process -> the "general" authentication
// in this process, the mesh API encodes a JWT token using its own "general" key

// MeshGlobalKey is the env variable name where the global key is stored
const MeshGlobalKey = "MESH_GLOBAL_KEY"

// General handles the mesh API general key
type General struct {
	Key []byte
}

// GenerateKey generates a new key
func (g *General) GenerateKey() []byte {
	g.Key = []byte(uniuri.New())
	return g.Key
}

// KeyExists returns true if the mesh global key is set in the env
func (g *General) KeyExists() bool {
	return os.Getenv(MeshGlobalKey) != ""
}

// SetKey sets the mesh global key in the env
func (g *General) SetKey(key []byte) {
	g.Key = key
	os.Setenv(MeshGlobalKey, string(g.Key))
}

// Token returns general JWT token with hostname claim
func (g *General) Token(hostname string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["hostname"] = hostname
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return t.SignedString(g.Key)
}

// ValidateGeneralTokenHeader fetches the access token from the authorization header
// and validates it
func ValidateGeneralTokenHeader(r *http.Request) (hostname string, err error) {
	t := extractToken(r)
	token, err := validateToken(t)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims")
		return
	}
	h, ok := claims["hostname"]
	if !ok {
		err = errors.New("invalid hostname")
		return
	}
	if hostname, ok = h.(string); !ok {
		err = errors.New("invalid hostname value type")
		return
	}
	return
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(MeshGlobalKey)), nil
	})
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
