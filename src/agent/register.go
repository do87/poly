package agent

import "github.com/golang-jwt/jwt"

// EncodeName encodes the key name using the private key
func (a *agent) EncodeName(k Key) (string, error) {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(k.PrivateKey)
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"scope": "agent",
		"host":  a.hostname,
	}).SignedString(parsedKey)
}
