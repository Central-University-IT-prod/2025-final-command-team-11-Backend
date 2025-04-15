-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
        CREATE TYPE role AS ENUM (
            'USER', 
            'ADMIN', 
            'SUPER_ADMIN', 
            'SUPPORT'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_field') THEN
        CREATE TYPE account_field AS ENUM (
            'ID',
            'EMAIL',
            'NAME',
            'ROLES',
            'BIRTHDAY',
            'VERIFIED'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth') THEN
        CREATE TYPE oauth AS ENUM (
            'YANDEX',
            'GOOGLE'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id       UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email    VARCHAR(255) UNIQUE NOT NULL,
    name     VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    birthday DATE,
    role    role DEFAULT 'USER',
    verified BOOLEAN DEFAULT false,
    oauth    oauth
);


INSERT INTO users (email, name, password, role, birthday) 
VALUES
('shahov@pobeda.ru', 'Alexandr', '$2a$10$O4PlOiMaBL1LPiODyIQoGuKujyksSCRwZEA3F1F/uaYAa8.3YTtoe', 'ADMIN', '1111-11-11'),
('vasily@mail.ru', 'Opyt', '$2a$10$3O/WbGMHLd5XT7mQj4K7zeUGD4LU.doFuvFseIFMIMghMXaEgjcoi', 'USER', '1111-11-11'),
('loddte@prod.com', 'Lodty', '$2a$10$E71rm3Pnx6Jy3A/8RBEKFuZDQhQyzr3PRL4avn7AXzsNlD6IFhTtK', 'USER', '1111-11-11'),
('shtil@pobeda.ru', 'Shtil', '$2a$10$hjZjKNP0rVtAhI535cngVutlCyzJVFnTzPDVYzAQvV1LOAYxiYADG', 'ADMIN', '1111-11-11'),
('megazord2000@plaki.com', 'Megazord', '$2a$10$Aor9q7fJBH8eXtu.KJf2IuBg0SLUPhNy8QITMXeU4u8iRpefXuDay', 'USER', '1111-11-11'),
('coffee@cookie.ru', 'Coffee', '$2a$10$DsbqT8M19zXYukmc2Oo.UurNVpJPA6ECaw.XtApKSg0G0M8J0DZnCa', 'ADMIN', '1111-11-11'),
('chipi@chapa.com', 'Chipi', '$2a$10$egfCSYbXwRFBxBirYkHRHOECKUHmM82d5/DoL0XoHtuVxZ0HHTa/e', 'USER', '1111-11-11'),
('some@some.com', 'Id', '$2a$10$8HiKlUxwosPz50rv9Kz4beTYCHoXeVFTAQ7/DjnGu9ZAu1Vtvisly', 'USER', '1111-11-11'),
('prodovich@hse.ru', 'Prod', '$2a$10$7mnDHJUiFDRMGNjiGl8/9ursS20lCPA7EvGlsqTWdPyQk9o40xVJq', 'USER', '1111-11-11'),
('admin@admin.ru', 'Admin', '$2a$10$bkhWParYUtO4IT34Gw46tu14/TWX1/XeciEywSpKA6t7WouPefWe.', 'ADMIN', '1111-11-11') ON CONFLICT DO NOTHING
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;

DROP TABLE IF EXISTS yandex;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS role;

DROP TYPE IF EXISTS account_field;

DROP TYPE IF EXISTS oauth;
-- +goose StatementEnd
