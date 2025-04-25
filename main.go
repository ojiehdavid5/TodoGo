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
	// return c.SendString("Hello")

	var res string
   var todos []string
   rows, err := db.Query("SELECT * FROM todos")
   defer rows.Close()
   if err != nil {
       log.Fatalln(err)
       c.JSON("An error occured")
   }
   for rows.Next() {
       rows.Scan(&res)
       todos = append(todos, res)
   }
   return c.Render("index", fiber.Map{
       "Todos": todos,
   })
}
type todo struct {
	Item string `json:"item"`
}
func postHandler(c *fiber.Ctx,db *sql.DB) error {
		newTodo := todo{}
		if err := c.BodyParser(&newTodo); err != nil {
			log.Printf("An error occured: %v", err)
			return c.SendString(err.Error())
		}
		fmt.Printf("%v", newTodo)
		if newTodo.Item != "" {
			_, err := db.Exec("INSERT into todos VALUES ($1)", newTodo.Item)
			if err != nil {
				log.Fatalf("An error occured while executing query: %v", err)
			}
		}
	 
		return c.Redirect("/")
	 }
	 

func putHandler(c *fiber.Ctx,db *sql.DB) error {olditem := c.Query("olditem")
newitem := c.Query("newitem")
db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
return c.Redirect("/")

}
func deleteHandler(c *fiber.Ctx,db *sql.DB) error {
	item := c.Query("item")
	_, err := db.Exec("DELETE FROM todos WHERE item=$1", item)
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)
	}
	return c.Redirect("/")
}
