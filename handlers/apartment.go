package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ApartmentHandler struct {
	db *sql.DB
}

type Apartment struct {
	ID         int    `json:"id" db:"id"`
	BuildingID int    `json:"building_id" db:"building_id"`
	Number     string `json:"number" db:"number"`
	Floor      int    `json:"floor" db:"floor"`
	SqMeters   int    `json:"sq_meters" db:"sq_meters"`
}

func NewApartmentHandler(db *sql.DB) *ApartmentHandler {
	return &ApartmentHandler{db: db}
}

func (h *ApartmentHandler) GetApartments(c *fiber.Ctx) error {
	rows, err := h.db.Query(`
		SELECT id, building_id, number, floor, sq_meters 
		FROM apartment 
		ORDER BY building_id, floor, number
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch apartments"})
	}
	defer rows.Close()

	var apartments []Apartment
	for rows.Next() {
		var apartment Apartment
		if err := rows.Scan(&apartment.ID, &apartment.BuildingID, &apartment.Number, &apartment.Floor, &apartment.SqMeters); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan apartment"})
		}
		apartments = append(apartments, apartment)
	}

	return c.JSON(apartments)
}

func (h *ApartmentHandler) GetApartment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid apartment ID"})
	}

	var apartment Apartment
	err = h.db.QueryRow(`
		SELECT id, building_id, number, floor, sq_meters 
		FROM apartment 
		WHERE id = $1
	`, id).Scan(&apartment.ID, &apartment.BuildingID, &apartment.Number, &apartment.Floor, &apartment.SqMeters)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Apartment not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch apartment"})
	}

	return c.JSON(apartment)
}

func (h *ApartmentHandler) GetApartmentsByBuilding(c *fiber.Ctx) error {
	buildingID, err := strconv.Atoi(c.Params("buildingId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid building ID"})
	}

	var buildingExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM building WHERE id = $1)", buildingID).Scan(&buildingExists)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check building existence"})
	}
	if !buildingExists {
		return c.Status(404).JSON(fiber.Map{"error": "Building not found"})
	}

	rows, err := h.db.Query(`
		SELECT id, building_id, number, floor, sq_meters 
		FROM apartment 
		WHERE building_id = $1 
		ORDER BY floor, number
	`, buildingID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch apartments"})
	}
	defer rows.Close()

	var apartments []Apartment
	for rows.Next() {
		var apartment Apartment
		if err := rows.Scan(&apartment.ID, &apartment.BuildingID, &apartment.Number, &apartment.Floor, &apartment.SqMeters); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan apartment"})
		}
		apartments = append(apartments, apartment)
	}

	return c.JSON(apartments)
}

func (h *ApartmentHandler) CreateApartment(c *fiber.Ctx) error {
	var apartment Apartment
	if err := c.BodyParser(&apartment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if apartment.BuildingID == 0 || apartment.Number == "" || apartment.Floor == 0 || apartment.SqMeters == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	var buildingExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM building WHERE id = $1)", apartment.BuildingID).Scan(&buildingExists)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check building existence"})
	}
	if !buildingExists {
		return c.Status(400).JSON(fiber.Map{"error": "Building does not exist"})
	}

	query := `
		INSERT INTO apartment (building_id, number, floor, sq_meters) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (building_id, number) 
		DO UPDATE SET floor = EXCLUDED.floor, sq_meters = EXCLUDED.sq_meters, updated_at = CURRENT_TIMESTAMP
		RETURNING id, building_id, number, floor, sq_meters
	`

	err = h.db.QueryRow(query, apartment.BuildingID, apartment.Number, apartment.Floor, apartment.SqMeters).
		Scan(&apartment.ID, &apartment.BuildingID, &apartment.Number, &apartment.Floor, &apartment.SqMeters)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create/update apartment"})
	}

	return c.Status(201).JSON(apartment)
}

func (h *ApartmentHandler) DeleteApartment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid apartment ID"})
	}

	result, err := h.db.Exec("DELETE FROM apartment WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete apartment"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Apartment not found"})
	}

	return c.JSON(fiber.Map{"message": "Apartment deleted successfully"})
}
