package initializers

import (
    "context"
	"fmt"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/mahesh060304/go-crud/models"

)

var UserCollection *mongo.Collection 	

func ConnectToDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/crud")

    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Error connecting to MongoDB: ", err)
    }
	UserCollection = client.Database("crud").Collection("users")
    log.Println("Connected to MongoDB!")
}


func CreateUser(user *models.User) error {
    // Insert user into the database
    _, err := UserCollection.InsertOne(context.Background(), user)
    if err != nil {
        return err
    }

    fmt.Println("User created successfully!")
    return nil
}


