package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"

	"building-management-system/config"
	"building-management-system/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)
	defer db.Close()

	runMigrations(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	setupRoutes(app, db)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}

func setupRoutes(app *fiber.App, db *sql.DB) {
	buildingHandler := handlers.NewBuildingHandler(db)
	apartmentHandler := handlers.NewApartmentHandler(db)

	app.Get("/buildings", buildingHandler.GetBuildings)
	app.Get("/buildings/:id", buildingHandler.GetBuilding)
	app.Post("/buildings", buildingHandler.CreateBuilding)
	app.Delete("/buildings/:id", buildingHandler.DeleteBuilding)

	app.Get("/apartments", apartmentHandler.GetApartments)
	app.Get("/apartments/:id", apartmentHandler.GetApartment)
	app.Get("/apartments/building/:buildingId", apartmentHandler.GetApartmentsByBuilding)
	app.Post("/apartments", apartmentHandler.CreateApartment)
	app.Delete("/apartments/:id", apartmentHandler.DeleteApartment)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Building Management System API is running",
		})
	})
}

func runMigrations(db *sql.DB) {
	migrationFile := "migrations/001_create_tables.sql"

	content, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		log.Printf("Migration file not found: %v", err)
		return
	}

	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("Failed to execute migration: %v", err)
		}
	}

	log.Println("Migrations executed successfully")
}
