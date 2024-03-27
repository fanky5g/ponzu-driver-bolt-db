package root

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
)

var contentIndexName = "__contentIndex"

type repository struct {
	db        *bolt.DB
	entityMap map[string]content.Builder
}

// New instantiates common repository functions implemented by all repositories
func New(db *bolt.DB, contentTypes map[string]content.Builder) (repositories.ContentRepositoryInterface, error) {
	repo := &repository{db: db, entityMap: contentTypes}
	for itemName, itemType := range contentTypes {
		if err := repo.CreateEntityStore(itemName, itemType()); err != nil {
			return nil, err
		}
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(contentIndexName))
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to create storage bucket: %v", contentIndexName)
	}

	return repo, nil
}
