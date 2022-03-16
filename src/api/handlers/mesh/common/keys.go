package common

import (
	"errors"
	"time"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/auth"
)

// ProcessRegisterKey processes a registration key using the provided public key
func ProcessRegisterKey(key models.Key, encodedKey, agentHostname string) error {
	if time.Now().After(key.ExpiresAt) {
		return errors.New("given key has expired")
	}
	meshKey := auth.MeshRegisterKey{
		Name:      key.Name,
		PublicKey: key.PublicKey,
	}
	t, err := meshKey.Decode(encodedKey)
	if err != nil {
		return err
	}
	if err := auth.ValidateRegisterToken(t, agentHostname); err != nil {
		return err
	}
	return nil
}
