package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	database "github.com/jordannadroj/52_in_52/04_todo_app/db"
	"github.com/jordannadroj/52_in_52/04_todo_app/handler"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	//load .env file
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	db := database.ConnectToDB(&database.Config{})
	defer db.SqlDB.Close()

	httpHandler := handler.NewHttpHandler(db)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Static("/public", "./public")

	app.Get("/", httpHandler.IndexHandler)

	app.Post("/", httpHandler.PostHandler)

	app.Put("/update", httpHandler.PutHandler)

	app.Delete("/delete", httpHandler.DeleteHandler)

	// log.Fatalln will log the output in case of any errors.
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
