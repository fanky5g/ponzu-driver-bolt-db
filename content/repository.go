package content

import (
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu-driver-bolt-db/root"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
)

func New(db *bolt.DB, types map[string]content.Builder) (repositories.ContentRepositoryInterface, error) {
	return root.New(db, types)
}
