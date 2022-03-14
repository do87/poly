package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/api/handlers/mesh/payloads"
)

// KeysRepository is the usecase repo for key keys
type KeysRepository interface {
	GetByName(ctx context.Context, name string) (models.Key, error)
	List(ctx context.Context) ([]models.Key, error)
	Create(ctx context.Context, key models.Key) (models.Key, error)
	Delete(ctx context.Context, id string) error
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
func (u *keysUsecase) Delete(ctx context.Context, r *http.Request, id string) (err error) {
	return u.repo.Delete(ctx, id)
}
