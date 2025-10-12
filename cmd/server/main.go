package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"delivery_management_system/internal/web"
	"delivery_management_system/internal/middleware"
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/auth"
	"delivery_management_system/internal/handlers"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Prompt the user for the PostgreSQL password
	passwordPrompt := promptui.Prompt{
		Label: "rupupuru@01",
		Mask:  '*', // Mask the input as the user types the password
	}

	password, err := passwordPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	// Set the database connection string with the password
	// Make sure to replace username and dbname with your actual values
	dbConnStr := fmt.Sprintf("user=%s password=%s dbname=delivery_management_system sslmode=disable",
		os.Getenv("DB_USER"), password)

	// Initialize the database with the connection string
	db.Initialize(dbConnStr)

	// Set up authentication
	auth.InitAuth()

	// Set up routes
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/admin", handlers.AdminHandler)
	http.HandleFunc("/order", handlers.OrderHandler)
	http.HandleFunc("/tracking", handlers.TrackingHandler)

	// Apply middleware for authentication
	http.Handle("/admin", middleware.AuthMiddleware(http.HandlerFunc(handlers.AdminHandler)))

	// Start the server
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
