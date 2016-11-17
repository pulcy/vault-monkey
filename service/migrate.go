package service

import (
	"path"

	logging "github.com/op/go-logging"

	"github.com/pulcy/vault-monkey/service/migration"
)

func Migrate(from, to migration.Backend, log *logging.Logger) error {
	if err := migrate(from, to, "", log); err != nil {
		return maskAny(err)
	}
	return nil
}

func migrate(from, to migration.Backend, baseKey string, log *logging.Logger) error {
	keys, err := from.List(baseKey)
	if err != nil {
		return maskAny(err)
	}
	for _, key := range keys {
		key = path.Join(baseKey, key)
		log.Debugf("Migrating %s", key)
		value, err := from.Get(key)
		if err != nil {
			return maskAny(err)
		}
		if value != nil {
			if err := to.Set(key, value); err != nil {
				return maskAny(err)
			}
		}
		if err := migrate(from, to, key, log); err != nil {
			return maskAny(err)
		}
	}

	return nil
}
