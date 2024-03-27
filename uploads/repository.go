package uploads

import (
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu-driver-bolt-db/root"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
)

func New(db *bolt.DB) (repositories.ContentRepositoryInterface, error) {
	return root.New(db, map[string]content.Builder{
		"uploads": func() interface{} {
			return new(entities.FileUpload)
		},
	})
}
