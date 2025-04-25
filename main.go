package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	// "github.com/gravityblast/fresh"
	"database/sql"
	"github.com/lib/pq"
	// "log"
)

func main() {

	connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"

	//connecting to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	app := fiber.New()

	// app.Get("/",indexHandler)
	// app.Post("/",postHandler)
	// app.Put("/update",putHandler)
	// app.Delete("/delete",deleteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

	fmt.Println("Hello, World!")

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)

	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c,db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c,db)
	})
	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c,db)
	})
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello")
}
func postHandler(c *fiber.Ctx,db *sql.DB) error {
	return c.SendString("Hello")
}
func putHandler(c *fiber.Ctx,db *sql.DB) error {
	return c.SendString("Hello")
}
func deleteHandler(c *fiber.Ctx,db *sql.DB) error {
	return c.SendString("Hello")
}
