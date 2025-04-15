-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'booking_type') THEN
        CREATE TYPE booking_type AS ENUM (
            'ROOM',
            'OPEN_SPACE'
        );
    END IF;
END $$;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS entity_floor (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT (now()),
    updated_at TIMESTAMP DEFAULT (now())
);

CREATE TRIGGER update_entity_floor_updated_at
BEFORE UPDATE ON entity_floor
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS booking_entity (
    id UUID PRIMARY KEY,
    type booking_type,
    title VARCHAR(255),
    x INT,
    y INT,
    floor_id UUID,
    width INT,
    height INT,
    capacity INT,
    created_at TIMESTAMP DEFAULT (now()),
    updated_at TIMESTAMP DEFAULT (now()),
    FOREIGN KEY (floor_id) REFERENCES entity_floor (id) ON DELETE CASCADE 
);

CREATE TRIGGER update_booking_entity_updated_at
BEFORE UPDATE ON booking_entity
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS verification_data (
    user_id UUID PRIMARY KEY,
    passport_image VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS booking (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    entity_id UUID NOT NULL,
    user_id UUID NOT NULL,
    time_from TIMESTAMP NOT NULL,
    time_to TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    FOREIGN KEY (entity_id) REFERENCES booking_entity (id) ON DELETE CASCADE
);

CREATE TRIGGER update_booking_updated_at
BEFORE UPDATE ON booking
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    booking_id UUID NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT (false),
    thing VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    FOREIGN KEY (booking_id) REFERENCES booking (id) ON DELETE CASCADE
);

CREATE TRIGGER update_order_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS guest (
    user_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    FOREIGN KEY (booking_id) REFERENCES booking (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS etity_floor;

DROP TABLE IF EXISTS booking_entity;
-- +goose StatementEnd
