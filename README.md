# Building Management System API

A REST API to manage buildings and apartments built with Go, Fiber framework, and PostgreSQL.

## Features

- **Buildings Management**: Create, read, update, and delete buildings
- **Apartments Management**: Manage apartments within buildings
- **Data Integrity**: Foreign key constraints and unique constraints
- **Error Handling**: Comprehensive error responses
- **Auto-migrations**: Database schema setup on startup

## Technology Stack

- **Language**: Go 1.24+
- **Framework**: Fiber v2
- **Database**: PostgreSQL (configurable)
- **ORM**: Raw SQL queries (SQLBoiler ready)

## Project Structure

```
├── config/
│   └── database.go        # Database configuration
├── handlers/
│   ├── building.go        # Building handlers
│   └── apartment.go       # Apartment handlers
├── migrations/
│   └── 001_create_tables.sql  # Database schema
├── main.go                # Application entry point
├── go.mod                 # Go modules
├── config.env             # Environment configuration
└── sqlboiler.toml         # SQLBoiler configuration
```

## Database Schema

### Building Table
- `id`: Primary key, auto-increment
- `name`: Unique building name
- `address`: Building address
- `created_at`, `updated_at`: Timestamps

### Apartment Table
- `id`: Primary key, auto-increment
- `building_id`: Foreign key to building
- `number`: Apartment number (unique per building)
- `floor`: Floor number
- `sq_meters`: Square meters
- `created_at`, `updated_at`: Timestamps

## Setup Instructions

### 1. Prerequisites
- Go 1.24 or higher
- PostgreSQL running locally or remotely

### 2. Clone and Setup
```bash
git clone <repository>
cd building-management-system
go mod tidy
```

### 3. Database Configuration
Update `config.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=building_management
DB_SSLMODE=disable
SERVER_PORT=3000
```

### 4. Create Database
```bash
createdb building_management
```

### 5. Run Application
```bash
go run main.go
```

The API will be available at `http://localhost:3000`

## API Endpoints

### Buildings

#### Get All Buildings
```http
GET /buildings
```
Optional query parameter:
- `include_apartments=true` - Include apartments in response

#### Get Building by ID
```http
GET /buildings/{id}
```

#### Create/Update Building
```http
POST /buildings
Content-Type: application/json

{
  "name": "Building Name",
  "address": "Building Address"
}
```

#### Delete Building
```http
DELETE /buildings/{id}
```

### Apartments

#### Get All Apartments
```http
GET /apartments
```

#### Get Apartment by ID
```http
GET /apartments/{id}
```

#### Get Apartments by Building
```http
GET /apartments/building/{buildingId}
```

#### Create/Update Apartment
```http
POST /apartments
Content-Type: application/json

{
  "building_id": 1,
  "number": "101",
  "floor": 1,
  "sq_meters": 75
}
```

#### Delete Apartment
```http
DELETE /apartments/{id}
```

### Health Check
```http
GET /health
```

## Example Usage

### Create a Building
```bash
curl -X POST http://localhost:3000/buildings \
  -H "Content-Type: application/json" \
  -d '{"name": "Sunrise Tower", "address": "123 Main St"}'
```

### Create an Apartment
```bash
curl -X POST http://localhost:3000/apartments \
  -H "Content-Type: application/json" \
  -d '{"building_id": 1, "number": "101", "floor": 1, "sq_meters": 75}'
```

### Get Buildings with Apartments
```bash
curl "http://localhost:3000/buildings?include_apartments=true"
```

## Error Handling

The API returns consistent error responses:
```json
{
  "error": "Error description"
}
```

HTTP status codes used:
- `200`: Success
- `201`: Created
- `400`: Bad Request
- `404`: Not Found
- `500`: Internal Server Error

## Development

### Running Tests
```bash
go test ./...
```

### Code Generation with SQLBoiler
```bash
sqlboiler psql
```

## Production Deployment

1. Set environment variables instead of using `config.env`
2. Use connection pooling for database
3. Add proper logging and monitoring
4. Set up reverse proxy (nginx)
5. Use HTTPS in production 