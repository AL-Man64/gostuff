package main

import (
	"database/sql"
	"fmt"
	"log"
    "os"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

type User struct {
    id       uint32
    username string
}

func main() {
    log.Println("starting server")

    connStr := "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    createUserTable(db)

    app := fiber.New()
    app.Use(cors.New(cors.Config{
        AllowOriginsFunc: func(origin string) bool {
            return os.Getenv("ENVIRONMENT") == "development"
        },
    }))

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello World")
    })
    app.Get("/:username", func(c *fiber.Ctx) error {
        query := `SELECT id FROM users WHERE username = $1;`

        var pk int

        err := db.QueryRow(query, c.Params("username")).Scan(&pk)
        if err != nil {
            query := `INSERT INTO users (username) VALUES ($1) RETURNING id;`
            err := db.QueryRow(query, c.Params("username")).Scan(&pk)
            if err != nil {
                log.Fatal(err)
            }
        }

        return c.SendString(fmt.Sprintf("Your user id is %v", pk))
    })

    app.Listen(":3000")
}

func createUserTable(db *sql.DB) {
    query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(20) NOT NULL
    );`

    _, err := db.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}
