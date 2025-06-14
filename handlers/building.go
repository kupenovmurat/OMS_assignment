package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BuildingHandler struct {
	db *sql.DB
}

type Building struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
}

type BuildingWithApartments struct {
	Building
	Apartments []Apartment `json:"apartments,omitempty"`
}

func NewBuildingHandler(db *sql.DB) *BuildingHandler {
	return &BuildingHandler{db: db}
}

func (h *BuildingHandler) GetBuildings(c *fiber.Ctx) error {
	includeApartments := c.Query("include_apartments") == "true"

	if includeApartments {
		return h.getBuildingsWithApartments(c)
	}

	rows, err := h.db.Query("SELECT id, name, address FROM building ORDER BY id")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database query failed"})
	}
	defer rows.Close()

	var buildings []Building
	for rows.Next() {
		var building Building
		if err := rows.Scan(&building.ID, &building.Name, &building.Address); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan building"})
		}
		buildings = append(buildings, building)
	}

	return c.JSON(buildings)
}

func (h *BuildingHandler) getBuildingsWithApartments(c *fiber.Ctx) error {
	query := `
		SELECT b.id, b.name, b.address, 
			   a.id, a.number, a.floor, a.sq_meters
		FROM building b
		LEFT JOIN apartment a ON b.id = a.building_id
		ORDER BY b.id, a.id
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch buildings with apartments"})
	}
	defer rows.Close()

	buildingsMap := make(map[int]*BuildingWithApartments)

	for rows.Next() {
		var buildingID int
		var buildingName, buildingAddress string
		var apartmentID sql.NullInt64
		var apartmentNumber sql.NullString
		var apartmentFloor, apartmentSqMeters sql.NullInt64

		err := rows.Scan(&buildingID, &buildingName, &buildingAddress,
			&apartmentID, &apartmentNumber, &apartmentFloor, &apartmentSqMeters)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan building with apartments"})
		}

		if _, exists := buildingsMap[buildingID]; !exists {
			buildingsMap[buildingID] = &BuildingWithApartments{
				Building: Building{
					ID:      buildingID,
					Name:    buildingName,
					Address: buildingAddress,
				},
				Apartments: []Apartment{},
			}
		}

		if apartmentID.Valid {
			apartment := Apartment{
				ID:         int(apartmentID.Int64),
				BuildingID: buildingID,
				Number:     apartmentNumber.String,
				Floor:      int(apartmentFloor.Int64),
				SqMeters:   int(apartmentSqMeters.Int64),
			}
			buildingsMap[buildingID].Apartments = append(buildingsMap[buildingID].Apartments, apartment)
		}
	}

	var buildings []BuildingWithApartments
	for _, building := range buildingsMap {
		buildings = append(buildings, *building)
	}

	return c.JSON(buildings)
}

func (h *BuildingHandler) GetBuilding(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid building ID"})
	}

	var building Building
	err = h.db.QueryRow("SELECT id, name, address FROM building WHERE id = $1", id).
		Scan(&building.ID, &building.Name, &building.Address)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Building not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch building"})
	}

	return c.JSON(building)
}

func (h *BuildingHandler) CreateBuilding(c *fiber.Ctx) error {
	var building Building
	if err := c.BodyParser(&building); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if building.Name == "" || building.Address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and address are required"})
	}

	query := `
		INSERT INTO building (name, address) 
		VALUES ($1, $2) 
		ON CONFLICT (name) 
		DO UPDATE SET address = EXCLUDED.address, updated_at = CURRENT_TIMESTAMP
		RETURNING id, name, address
	`

	err := h.db.QueryRow(query, building.Name, building.Address).
		Scan(&building.ID, &building.Name, &building.Address)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create/update building"})
	}

	return c.Status(201).JSON(building)
}

func (h *BuildingHandler) DeleteBuilding(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid building ID"})
	}

	result, err := h.db.Exec("DELETE FROM building WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete building"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Building not found"})
	}

	return c.JSON(fiber.Map{"message": "Building deleted successfully"})
}
