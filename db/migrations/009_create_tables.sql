-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
    id              UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    username        TEXT UNIQUE NOT NULL CONSTRAINT username_length CHECK (char_length(username) <= 20),
    password        TEXT NOT NULL CONSTRAINT password_length CHECK (char_length(password) <= 255),
    email           TEXT UNIQUE NOT NULL CONSTRAINT email_length CHECK (char_length(email) <= 255),
    name            TEXT NOT NULL CONSTRAINT name_length CHECK (char_length(name) <= 50),
    score           NUMERIC,
    avatar          TEXT DEFAULT '/images/default.png' CONSTRAINT avatar_length CHECK (char_length(avatar) <= 1000),
    sex             TEXT CONSTRAINT sex_length CHECK (char_length(sex) = 1),
    guest_count     INT,
    birthdate       DATE,
    is_host         BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cities (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title           TEXT NOT NULL CONSTRAINT title_length CHECK (char_length(title) <= 100),
    description     TEXT CONSTRAINT description_length CHECK (char_length(description) <= 3000)
);

CREATE TABLE IF NOT EXISTS ads (
    id              UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    city_id         INT NOT NULL,
    author_id       UUID NOT NULL,
    address         TEXT CONSTRAINT address_length CHECK (char_length(address) <= 255),
    publication_date DATE,
    distance        NUMERIC,
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ad_available_dates (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    available_date  DATE,
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ad_positions (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    latitude        NUMERIC,
    longitude       NUMERIC,
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    image_url       TEXT CONSTRAINT image_url_length CHECK (char_length(image_url) <= 1000),
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS requests (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    user_id         UUID NOT NULL,
    status          TEXT DEFAULT 'pending' NOT NULL CONSTRAINT status_length CHECK (char_length(status) <= 255),
    create_date     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_date     TIMESTAMPTZ,
    close_date      TIMESTAMPTZ,
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviews (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         UUID NOT NULL,
    host_id         UUID NOT NULL,
    text            TEXT NOT NULL CONSTRAINT review_text_length CHECK (char_length(text) <= 1000),
    rating          NUMERIC CHECK (rating >= 0 AND rating <= 5),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (host_id) REFERENCES users(id) ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS requests CASCADE;
DROP TABLE IF EXISTS images CASCADE;
DROP TABLE IF EXISTS ad_positions CASCADE;
DROP TABLE IF EXISTS ad_available_dates CASCADE;
DROP TABLE IF EXISTS ads CASCADE;
DROP TABLE IF EXISTS cities CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
