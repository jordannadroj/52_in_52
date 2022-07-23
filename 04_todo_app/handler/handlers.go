package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jordannadroj/52_in_52/04_todo_app/db"
	log "github.com/sirupsen/logrus"
)

type todo struct {
	Item string
}

type HttpHandler struct {
	database *db.Database
}

func NewHttpHandler(database *db.Database) *HttpHandler {
	return &HttpHandler{
		database: database,
	}
}

func (h *HttpHandler) IndexHandler(c *fiber.Ctx) error {
	var res string
	var todos []string

	rows, err := h.database.SqlDB.Query("SELECT * FROM todos") // return all rows that are returned by query
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() { // read each row
		rows.Scan(&res)            // Scan() copies the row into a dedicated pointer variable
		todos = append(todos, res) // append each row to the todos array
	}
	return c.Render("index", fiber.Map{
		"Todos": todos, // the key "Todos" is what we will use in the html as our access to the todos array
	})
}

// there is a form in the html file with action POST
func (h *HttpHandler) PostHandler(c *fiber.Ctx) error {
	//	add a new todo to the list in the db
	//	render the todos again with the list
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error occured %v", err)
		return c.SendString(err.Error())
	}

	if newTodo.Item != "" {
		_, err := h.database.SqlDB.Exec("INSERT INTO todos VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Printf("An error occured while executing query: %v", err)
		}
		log.Printf("Row added to database 'todos' with value %q", newTodo.Item)
	}

	return c.Redirect("/")
}

func (h *HttpHandler) PutHandler(c *fiber.Ctx) error {
	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	_, err := h.database.SqlDB.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	if err != nil {
		log.Errorf("An error occured while executing query: %v", err)
		return c.SendString(err.Error())
	}
	log.Infof("Item %q updated to %q", olditem, newitem)
	return nil
}

func (h *HttpHandler) DeleteHandler(c *fiber.Ctx) error {
	deleteItem := c.Query("item")
	_, err := h.database.SqlDB.Exec("DELETE FROM todos WHERE item=$1", deleteItem)
	if err != nil {
		log.Errorf("An error occured while executing query: %v", err)
		return c.SendString(err.Error())
	}
	log.Infof("%v deleted", deleteItem)
	return nil
}
