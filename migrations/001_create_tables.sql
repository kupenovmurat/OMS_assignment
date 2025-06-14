CREATE TABLE IF NOT EXISTS building (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS apartment (
    id SERIAL PRIMARY KEY,
    building_id INTEGER NOT NULL REFERENCES building(id) ON DELETE CASCADE,
    number VARCHAR(50) NOT NULL,
    floor INTEGER NOT NULL,
    sq_meters INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building_id, number)
);

CREATE INDEX IF NOT EXISTS idx_apartment_building_id ON apartment(building_id);
CREATE INDEX IF NOT EXISTS idx_apartment_floor ON apartment(floor); 