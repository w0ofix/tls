package main

import (
    "os"
    "fmt"
    "log"

    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "github.com/joho/godotenv"

    "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"

    "github.com/w0ofix/tls/models"
    "github.com/w0ofix/tls/routes"
)

func main() {
    if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Fatal("failed to migrate database: ", err)
    }

    app := fiber.New()
	app.Use(cors.New()) // TODO: Config

    app.Get("/", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "version": "1.0.0",
        })
    })
    routes.RegisterAuthRoutes(app, db)

    app.Get("/*", static.New("./public"))

    log.Fatal(app.Listen(":3000"))
}
