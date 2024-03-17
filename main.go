package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Todo struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Done  bool   `gorm:"default:false" json:"done"`
}

func main() {
	// Establish database connection
	initDB()

	router := gin.Default()

	router.POST("/todos", addTodo)
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getOne)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", removeTodo)

	router.Run("0.0.0.0:3333")
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate the schema
	db.AutoMigrate(&Todo{})
}

func addTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add todo item to the database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": todo, "message": "Todo item added successfully"})
}

func getTodos(c *gin.Context) {
	var todos []Todo
	if err := db.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo items from the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todos, "total": len(todos)})
}

func getOne(c *gin.Context) {
	var todo Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func updateTodo(c *gin.Context) {
	var todo Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo item not found"})
		return
	}

	todo.Done = !todo.Done
	if err := db.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo, "message": "Todo item updated successfully"})
}

func removeTodo(c *gin.Context) {
	if err := db.Delete(&Todo{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo item deleted successfully"})
}
