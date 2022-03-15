package repos

import (
	"context"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"gorm.io/gorm"
)

type keysRepo struct {
	db *gorm.DB
}

// Get returns key by name
func (r *keysRepo) Get(ctx context.Context, name string) (key models.Key, err error) {
	result := r.db.First(&key, "name = ?", name)
	if result.Error != nil {
		return models.Key{}, result.Error
	}
	return
}

// FirstOrCreateGeneralKey returns the global worker key or creates a new one if one doesn't exist
func (r *keysRepo) FirstOrCreateGeneralKey(ctx context.Context, fallbackKey []byte) (key models.Key, err error) {
	result := r.db.First(&key, "name = ?", "general")
	if result.Error == nil {
		return
	}
	key = models.Key{
		Name:      "general",
		PublicKey: fallbackKey,
	}
	if result := r.db.FirstOrCreate(&key); result.Error != nil {
		return key, result.Error
	}
	return
}

// List returns all keys
func (r *keysRepo) List(ctx context.Context) (keys []models.Key, err error) {
	result := r.db.Order("name ASC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// Create a new key
func (r *keysRepo) Create(ctx context.Context, key models.Key) (models.Key, error) {
	if result := r.db.FirstOrCreate(&key, "name = ?", key.Name); result.Error != nil {
		return key, result.Error
	}
	return key, nil
}

// Delete a key by name
func (r *keysRepo) Delete(ctx context.Context, name string) error {
	a, err := r.Get(ctx, name)
	if err != nil {
		return err
	}
	if err := r.db.Delete(&a, "name = ?", name).Error; err != nil {
		return err
	}
	return nil
}
