package db

import (
	"context"
	"log"
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

func TestPost(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
}
