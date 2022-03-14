package usecases

import (
	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/api/handlers/mesh/payloads"
	"github.com/do87/poly/src/auth"
)

// MeshGlobalKey is the env variable name where the global key is stored
const MeshGlobalKey = "MESH_GLOBAL_KEY"

func (u *agentsUsecase) processKey(key models.Key, payload payloads.AgentRegister) error {
	meshKey := auth.MeshRegisterKey{
		Name:      key.Name,
		PublicKey: key.PublicKey,
	}
	t, err := meshKey.Decode(payload.EncodedKey.Encoded)
	if err != nil {
		return err
	}
	if err := auth.ValidateRegisterToken(t, payload.Hostname); err != nil {
		return err
	}
	return nil
}
