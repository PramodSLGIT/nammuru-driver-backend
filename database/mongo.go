package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDb struct {
	c            *mongo.Client
	rt           time.Duration
	wg           sync.WaitGroup
	shutDownflag int32
}

var (
	_initMongoDBCtx sync.Once
	mdb             *MongoDb
)


func NewDbConn() (*MongoDb, error) {
	_initMongoDBCtx.Do(func() {
		// Set a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Define the URI for local MongoDB connection
		uri := "mongodb://localhost:27019"

		fmt.Println("Connecting to MongoDB")

		// Connect to MongoDB using the URI
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		mdb = &MongoDb{
			c:  client,
			wg: sync.WaitGroup{},
			rt: 5 * time.Second,
		}
	})

	return mdb, nil
}


func (mdb *MongoDb) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), mdb.rt)
	defer cancel()

	return mdb.c.Ping(ctx, readpref.Primary())
}

func (mdb *MongoDb) Use(dbName, collectionName string) *mongo.Collection {
	var ctx = context.TODO()
	err := mdb.c.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return mdb.c.Database(dbName).Collection(collectionName)
}

// Close the db connection
func (mdb *MongoDb) Close() error {
	if mdb.c != nil {
		ctx := context.TODO()

		atomic.StoreInt32(&mdb.shutDownflag, 1)
		err := mdb.c.Disconnect(ctx)

		mdb.wg.Wait()

		return err
	}

	return nil
}

