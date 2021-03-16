package database

import (
	"github.com/ra9dev/PROJECTNAME/internal/adapters/database/drivers"
	"github.com/ra9dev/PROJECTNAME/internal/adapters/database/drivers/memdb"
	"github.com/ra9dev/PROJECTNAME/internal/adapters/database/drivers/mongo"
)

func New(conf drivers.DataStoreConfig) drivers.DataStore {
	switch conf.Name {
	case "mongo":
		return mongo.New(conf)
	default:
		return memdb.New()
	}
}
