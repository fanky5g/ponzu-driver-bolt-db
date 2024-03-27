package root

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
)

// FindOneBySlug does a lookup in the entities index to find the type and id of
// the requested entities. Subsequently, issues the lookup in the type bucket and
// returns the type and data at that ID or nil if nothing exists.
func (repo *repository) FindOneBySlug(slug string) (string, interface{}, error) {
	val := &bytes.Buffer{}
	var t, id string
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(contentIndexName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		idx := b.Get([]byte(slug))
		if idx != nil && len(idx) > 0 {
			tid := strings.Split(string(idx), ":")

			if len(tid) < 2 {
				return fmt.Errorf("bad data in entities index for slug: %s", slug)
			}

			t, id = tid[0], tid[1]
		}

		if t == "" {
			return nil
		}

		c := tx.Bucket([]byte(t))
		if c == nil {
			return bolt.ErrBucketNotFound
		}
		_, err := val.Write(c.Get([]byte(id)))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return t, nil, err
	}

	if val.Len() == 0 {
		return "", nil, nil
	}

	entity, err := repo.MarshalEntity(t, val)
	if err != nil {
		return "", nil, err
	}

	return t, entity, nil
}
