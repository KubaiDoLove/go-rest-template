package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/ra9dev/PROJECTNAME/internal/adapters/database/drivers"
)

type Database struct {
	connURL string
	dbName  string

	client *mongo.Client
	db     *mongo.Database
}

func (d *Database) Name() string { return "Mongo" }

func New(conf drivers.DataStoreConfig) drivers.DataStore {
	return &Database{
		connURL: conf.URL,
		dbName:  conf.DB,
	}
}

func (d *Database) Connect() error {
	if d.connURL == "" {
		return drivers.ErrInvalidConfigStruct
	}

	if d.dbName == "" {
		return drivers.ErrInvalidConfigStruct
	}

	ctx, cancel := context.WithTimeout(context.Background(), drivers.ConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.connURL))
	if err != nil {
		return err
	}
	d.client = client

	if err := d.Ping(); err != nil {
		return err
	}
	d.db = d.client.Database(d.dbName)

	// убеждаемся что созданы все необходимые индексы
	return d.ensureIndexes()
}

func (d *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), drivers.ConnectionTimeout)
	defer cancel()

	return d.client.Ping(ctx, readpref.Primary())
}

func (d *Database) Close(ctx context.Context) error {
	log.Printf("Disconnected from MongoDB: %s", d.dbName)
	return d.client.Disconnect(ctx)
}

// убеждается что все индексы построены
func (d *Database) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), drivers.EnsureIndexesTimeout)
	defer cancel()

	// TODO
	_ = ctx

	return nil
}

// indexExistsByName проверяет существование индекса с именем name.
func (d *Database) indexExistsByName(ctx context.Context, collection *mongo.Collection, name string) (bool, error) {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	for cur.Next(ctx) {
		if name == cur.Current.Lookup("name").StringValue() {
			return true, nil
		}
	}

	return false, nil
}
