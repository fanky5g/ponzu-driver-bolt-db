package root

import (
	"fmt"
	ponzuErrors "github.com/fanky5g/ponzu/errors"
	"github.com/fanky5g/ponzu/util"
	"github.com/gorilla/schema"
)

func (repo *repository) UpdateEntity(entityType, entityId string, update map[string]interface{}) (interface{}, error) {
	target := fmt.Sprintf("%s:%s", entityType, entityId)
	post, err := repo.FindOneByTarget(target)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ponzuErrors.ErrContentNotFound
	}

	u, err := mergeData(post, update)
	if err != nil {
		return nil, err
	}

	if _, err = repo.SetEntity(entityType, u); err != nil {
		return nil, err
	}

	return post, nil
}

func mergeData(post interface{}, update map[string]interface{}) (interface{}, error) {
	v := util.JSONMapToURLValues(update)

	v.Del("id")
	v.Del("uuid")
	v.Del("slug")
	v.Del("timestamp")

	dec := schema.NewDecoder()
	dec.SetAliasTag("json")     // allows simpler struct tagging when creating a entities type
	dec.IgnoreUnknownKeys(true) // will skip over form values submitted, but not in struct
	if err := dec.Decode(post, v); err != nil {
		return nil, err
	}

	return post, nil
}
