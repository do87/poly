package repos

import (
	"context"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"gorm.io/gorm"
)

type keysRepo struct {
	db *gorm.DB
}

// Get returns key by UUID
func (r *keysRepo) Get(ctx context.Context, id string) (keys models.Key, err error) {
	var key models.Key
	result := r.db.First(&key, "uuid = ?", id)
	if result.Error != nil {
		return models.Key{}, result.Error
	}
	return keys, nil
}

// GetByName returns key by name
func (r *keysRepo) GetByName(ctx context.Context, name string) (keys models.Key, err error) {
	var key models.Key
	result := r.db.First(&key, "name = ?", name)
	if result.Error != nil {
		return models.Key{}, result.Error
	}
	return keys, nil
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
	if result := r.db.FirstOrCreate(&key); result.Error != nil {
		return key, result.Error
	}
	return key, nil
}

// Delete a key by uuid
func (r *keysRepo) Delete(ctx context.Context, id string) error {
	a, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := r.db.Delete(&a, "uuid = ?", id).Error; err != nil {
		return err
	}
	return nil
}
