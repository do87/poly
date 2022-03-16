package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/payloads"
	"github.com/do87/poly/src/mesh/models"
)

// KeysRepository is the usecase repo for agent keys
type KeysRepository interface {
	Get(ctx context.Context, name string) (models.Key, error)
	List(ctx context.Context) ([]models.Key, error)
	Create(ctx context.Context, key models.Key) (models.Key, error)
	FirstOrCreateGeneralKey(ctx context.Context, fallbackKey []byte) (models.Key, error)
	Delete(ctx context.Context, name string) error
}

// NewKeysUsecase creates a new Usecase service
func NewKeysUsecase(keys KeysRepository) *Usecase {
	return &Usecase{
		Keys: &keysUsecase{
			repo: keys,
		},
	}
}

type keysUsecase struct {
	repo KeysRepository
}

// List returns a list of all keys
func (u *keysUsecase) List(ctx context.Context, r *http.Request) ([]models.Key, error) {
	return u.repo.List(ctx)
}

// Create registers an key and returns it
func (u *keysUsecase) Create(ctx context.Context, r *http.Request) (key models.Key, err error) {
	var payload payloads.KeyCreate
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	return u.repo.Create(ctx, payload.ToModel())
}

// Delete unregisters an key
func (u *keysUsecase) Delete(ctx context.Context, r *http.Request, name string) (err error) {
	return u.repo.Delete(ctx, name)
}

// FirstOrCreateGeneralKey returns a general key
// if the key doesn't exist, a new one will be created based on the fallbackKey string
func (u *keysUsecase) FirstOrCreateGeneralKey(ctx context.Context, fallbackKey []byte) (key models.Key, err error) {
	return u.repo.FirstOrCreateGeneralKey(ctx, fallbackKey)
}
