package controllers

import (
	"context"
    "log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/mahesh060304/go-crud/initializers"
    "github.com/mahesh060304/go-crud/models"
)

func CreateNewUser(c *gin.Context) {

	var user models.User

    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
	
    if err := initializers.CreateUser(&user); err != nil {
        log.Println("Error creating user:", err)
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(201, gin.H{"message": "User created successfully"})
}

func GetAllUsers(c *gin.Context) {
    cursor, err := initializers.UserCollection.Find(context.Background(), bson.M{})
    if err != nil {
        log.Println("Error retrieving users:", err)
        c.JSON(500, gin.H{"error": "Failed to retrieve users"})
        return
    }
    defer cursor.Close(context.Background())

    var users []models.User

    if err := cursor.All(context.Background(), &users); err != nil {
        log.Println("Error decoding users:", err)
        c.JSON(500, gin.H{"error": "Failed to retrieve users"})
        return
    }

    c.JSON(200, users)
}

func UpdateUser(c *gin.Context){
		id :=c.Param("id")

		var body models.User
        c.BindJSON(&body)

		objectID,err:= primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Println("Error decoding users:", err)
			return 
		}


		filter := bson.M{"_id":objectID}

        update := bson.M{
			"$set" : bson.M{	
				"Username": body.Username,
				"Email":body.Email,
				"Password":body.Password,
			},
		}

		result ,err := initializers.UserCollection.UpdateOne(context.Background(), filter,update)
		if err != nil {
			log.Println("Error decoding users:", err)
			return 
		}
	
		c.JSON(200,gin.H{"message":"User Updated successfully!","user":result})

}

func DeleteUser(c *gin.Context){
	    id :=c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)


		result, err := initializers.UserCollection.DeleteMany(context.Background(), bson.M{"_id":objectID})
        log.Println("Delete result:", result)

		if err != nil {
			log.Println("Error deleting users:", err)
			return 
		}

		if result.DeletedCount == 0{
			log.Println("Delete Count:", result.DeletedCount)

			c.JSON(400, gin.H{"error":"User Not Found"})
			return 
		}
		c.JSON(200, gin.H{"message": "User deleted successfully"})

}
	
