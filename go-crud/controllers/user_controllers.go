package controllers

import (
	"context"
	"log"
	"time"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mahesh060304/go-crud/initializers"
	"github.com/mahesh060304/go-crud/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	cacheKey := "user:" + user.Username
    userDataJSON, err := json.Marshal(user)
    if err != nil {
        log.Println("Error marshaling user data:", err)
    } else {
        ctx := context.Background()
        if err := initializers.RedisClient.Set(ctx, cacheKey, userDataJSON, 5*time.Minute).Err(); err != nil {
            log.Println("Error caching user data:", err)
        }
    }

    c.JSON(201, gin.H{"message": "User created successfully"})
}

func GetAllUsers(c *gin.Context) {
	cacheKey := "users"
    cachedData, err := initializers.GetFromCache(cacheKey)
    if err == nil && cachedData != "" {
        c.JSON(200, cachedData)
        return
    }

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

	dataJSON, err := bson.MarshalExtJSON(users, false, false)
    if err == nil {
        initializers.SetToCache(cacheKey, string(dataJSON), 10*time.Minute) // Cache for 10 minutes
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

		cacheKey := "user:" + id
		userDataJSON, err := json.Marshal(body)
		if err != nil {
			log.Println("Error marshaling user data:", err)
		} else {
			if err := initializers.RedisClient.Set(context.Background(), cacheKey, userDataJSON, 0).Err(); err != nil {
				log.Println("Error updating cache:", err)
			}
		}
	
		c.JSON(200,gin.H{"message":"User Updated successfully!","user":result})

}

func DeleteUser(c *gin.Context){
	    id :=c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err!=nil{
			log.Println("Error with objectid:", err)
			return 
		}


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

		cacheKey := "user:" + id
		if err := initializers.RedisClient.Del(context.Background(), cacheKey).Err(); err != nil {
			log.Println("Error invalidating cache:", err)
		}
		c.JSON(200, gin.H{"message": "User deleted successfully"})

}
	
