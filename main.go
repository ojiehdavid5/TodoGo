package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

func main() {
	connStr := "postgresql://user:postgres@localhost/todo?sslmode=disable"

	// Connecting to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close() // Ensure the database connection is closed when done

	app := fiber.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var todos []string
	rows, err := db.Query("SELECT item FROM todos")
	if err != nil {
		log.Println("Error querying todos:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database query error"})
	}
	defer rows.Close() // Close rows after processing

	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			log.Println("Error scanning row:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Error reading data"})
		}
		todos = append(todos, item)
	}

	return c.JSON( fiber.Map{
		"Todos": todos,
	})
}

type todo struct {
	Item string `json:"item"`
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if newTodo.Item != "" {
		_, err := db.Exec("INSERT INTO todos (item) VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatalf("Error executing query: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Database insert error"})
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	oldItem := c.Query("olditem")
	newItem := c.Query("newitem")
	if oldItem == "" || newItem == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Parameters missing"})
	}

	_, err := db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newItem, oldItem)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database update error"})
	}

	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	item := c.Query("item")
	if item == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Item parameter missing"})
	}

	_, err := db.Exec("DELETE FROM todos WHERE item=$1", item)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database delete error"})
	}

	return c.Redirect("/")
}