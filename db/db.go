package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// db exports functionality for connecting with the MongoDB database.
// The goal is to abstract away the specifics of database connections
// and provide atomic operations a la CRUD
// Typical usage:
//		db, err := OpenDB("zoltar","asthma")
// 		if err != nil {
//		log.Fatalf("Unable to open database: <%v>", err)
//		}
//		err = Connect(db)
//		if err != nil {
//			log.Fatalf("Unable to connect to databawse: <%v>", err)
//		}
//		defer Disconnect(db)
//
//		do some work with the database
//		insertResult, _ :=db.Collection.InsertOne(context.TODO(), data)
//		updateResult, _ := db.Collection.Update(context.TODO(), filter, update)

// DB holds the information for the current database, collection and client
type DB struct {
	client     *mongo.Client
	Collection *mongo.Collection
	database   string
	coll       string
	Ctx        context.Context
}

// OpenDB attempts to estblish a client link the provided database and collection.
// It assumes that a localhost MongoDB server.
func OpenDB(database string, collection string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	result := DB{client: client, coll: collection, database: database}
	result.Collection = client.Database(database).Collection(collection)
	return &result, nil
}

// Connect opens a connection to the database
func Connect(db *DB) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := db.client.Connect(ctx)
	if err != nil {
		return err
	}
	db.Ctx = ctx
	return nil
}

// Disconnect closes a connection to the database
// typically used with defer
func Disconnect(db *DB) {
	db.client.Disconnect(db.Ctx)
}
