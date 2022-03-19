package common

import (
	"errors"
	"time"

	"github.com/do87/poly/src/mesh/models"
	"github.com/do87/poly/src/pkg/auth"
)

// ProcessRegisterKey processes a registration key using the provided public key
func ProcessRegisterKey(key models.Key, encodedKey, uuid string) error {
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
	if err := auth.ValidateRegisterToken(t, uuid); err != nil {
		return err
	}
	return nil
}
