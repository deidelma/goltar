package db

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnection(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)
}

type Post struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func InsertPost(ctx context.Context, client *mongo.Client, title string, body string) (primitive.ObjectID, error) {
	post := Post{title, body}
	collection := client.Database("my_db").Collection("posts")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func GetPost(ctx context.Context, client *mongo.Client, id primitive.ObjectID) (Post, error) {
	var result Post
	filter := bson.D{{"_id", id}}
	collection := client.Database("my_db").Collection("posts")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func DeletePost(
	ctx context.Context,
	client *mongo.Client,
	id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.D{{"_id", id}}
	collection := client.Database("my_db").Collection("posts")
	return collection.DeleteOne(context.TODO(), filter)
}

func UpdatePost(
	ctx context.Context,
	client *mongo.Client,
	id primitive.ObjectID,
	newTitle string) (int64, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"title", newTitle},
		}},
	}
	collection := client.Database("my_db").Collection("posts")
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func TestPost(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)

	id, err := InsertPost(ctx, client, "This is a title", "This is the body")
	if err != nil {
		t.Errorf("Unable to insert post :<%v>", err)
	}
	log.Printf("Successfully posted with id: <%d>", id)

	post, err := GetPost(ctx, client, id)
	if err != nil {
		t.Errorf("Unable to retrieve post: <%v>", err)
	}
	log.Printf("Successfully retrieved post with title: '%s'", post.Title)

	updated, err := UpdatePost(ctx, client, id, "Hello")
	if err != nil {
		t.Errorf("Unable to update post: <%v>", err)
	}
	if updated != 1 {
		t.Errorf("Wrong number of posts updated: %d", updated)
	}
	post, _ = GetPost(ctx, client, id)
	log.Printf("Successfully updated title to: '%s'", post.Title)

	deleted, err := DeletePost(ctx, client, id)
	if err != nil {
		t.Errorf("Unable to delete post: <%v>", err)
	}
	log.Printf("Removed %d records", deleted.DeletedCount)

	post, err = GetPost(ctx, client, id)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "mongo: no documents") {
			return
		}
		t.Errorf("Failed to detect missing element:<%v>", err)
	}
	t.Errorf("Retrieved false record from search")

}

func TestDBCreationAndUse(t *testing.T) {
	db, err := OpenDB("my_db", "posts")
	if err != nil {
		t.Errorf("Unable to open DB: <%v>", err)
	}
	err = Connect(db)
	if err != nil {
		t.Errorf("Unable to connect to database %s.%s:<%v>", db.database, db.coll, err)
	}
	defer Disconnect(db)

	_, err = db.Collection.InsertOne(db.Ctx, Post{"Hello", "Goodbye"})
	if err != nil {
		t.Errorf("Unable to insert into database:<%v>", err)
	}
	filter := bson.D{{"title", "Hello"}}
	var post Post
	err = db.Collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		t.Errorf("Unable to retrieve post from database:<%v>", err)
	}
	if post.Title != "Hello" {
		t.Errorf("Retrieved wrong post!")
	}
	log.Printf("Successfully inserted and retrieved post")
}
