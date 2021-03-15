package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Completed string    `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DATABASE INSTANCE
var collection *mongo.Collection

func TodoCollection(c *mongo.Database) {
	collection = c.Collection("todos")
}

func GetAllTodos(c *gin.Context) {
	todos := []Todo{}
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}

func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	title := todo.Title
	body := todo.Body
	completed := todo.Completed
	id := guuid.New().String()

	newTodo := Todo{
		ID:        id,
		Title:     title,
		Body:      body,
		Completed: completed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), newTodo)

	if err != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID":      newTodo.ID,
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}

func GetSingleTodo(c *gin.Context) {
	todoId := c.Param("todoId")

	todo := Todo{}
	err := collection.FindOne(context.TODO(), bson.M{"id": todoId}).Decode(&todo)
	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    todo,
	})
	return
}

func EditTodo(c *gin.Context) {
	todoId := c.Param("todoId")
	var todo Todo
	err := collection.FindOne(context.TODO(), bson.M{"id": todoId}).Decode(&todo)
	//c.BindJSON(&todo)
	// newCompleted := todo.Completed
	// newBody := todo.Body
	//updated := time.Now().String()

	_, err = collection.DeleteOne(context.TODO(), bson.M{"id": todoId})

	todo.UpdatedAt = time.Now()

	// newData := bson.M{
	// 	"$set": bson.M{
	// 		"body":       newBody,
	// 		"completed":  newCompleted,
	// 		"updated_at": "$$NOW",
	// 	},
	// }

	_, err = collection.InsertOne(context.TODO(), todo)

	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo Edited Successfully",
	})
	return
}

// getUser function
// func getUser(id string) (Todo, error) {
// 	todo := Todo{}

// 	err := collection.Find(context.TODO(), bson.M{"_id": &id}).One(&todo)
// 	return todo, err
// }

func DeleteTodo(c *gin.Context) {
	todoId := c.Param("todoId")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": todoId})
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}
