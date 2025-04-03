package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/oracle"
	"gorm.io/gorm"
	"github.com/valyala/fasthttp"
	"github.com/fasthttp/router"
)

// Database Model
type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:255" json:"name"`
	Email string `gorm:"unique;size:255" json:"email"`
}

var db *gorm.DB

// Initialize Database
func initDB() {
	var err error
	dsn := "username/password@localhost:1521/XEPDB1" // Change this to your Oracle DB details
	db, err = gorm.Open(oracle.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Oracle database: %v", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
}

// Handlers
func createUser(ctx *fasthttp.RequestCtx) {
	user := User{
		Name:  string(ctx.FormValue("name")),
		Email: string(ctx.FormValue("email")),
	}
	result := db.Create(&user)
	if result.Error != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		fmt.Fprintf(ctx, "Error creating user: %v", result.Error)
		return
	}
	ctx.SetStatusCode(http.StatusCreated)
	fmt.Fprintf(ctx, "User created: %+v", user)
}

func getUsers(ctx *fasthttp.RequestCtx) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		fmt.Fprintf(ctx, "Error fetching users: %v", result.Error)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, "Users: %+v", users)
}

func getUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	var user User
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		fmt.Fprintf(ctx, "User not found")
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, "User: %+v", user)
}

func updateUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		fmt.Fprintf(ctx, "User not found")
		return
	}
	user.Name = string(ctx.FormValue("name"))
	user.Email = string(ctx.FormValue("email"))
	db.Save(&user)
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, "User updated: %+v", user)
}

func deleteUser(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	db.Delete(&User{}, "id = ?", id)
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, "User deleted")
}

// Setup Routes
func setupRouter() *router.Router {
	r := router.New()
	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.GET("/users/{id}", getUser)
	r.PUT("/users/{id}", updateUser)
	r.DELETE("/users/{id}", deleteUser)
	return r
}

func main() {
	initDB()
	r := setupRouter()

	log.Println("Server started on port 8080")
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
