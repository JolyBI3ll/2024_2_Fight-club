-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
    uuid            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(20) UNIQUE NOT NULL,
    password        VARCHAR(20) NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(50) NOT NULL,
    score           NUMERIC,
    avatar          TEXT CHECK (char_length(avatar) <= 1000),
    sex             VARCHAR(1),
    guest_count     INT,
    birthdate       DATE,
    is_host         BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cities (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(100),
    description     TEXT CHECK (char_length(description) <= 3000)
);

CREATE TABLE IF NOT EXISTS ads (
    uuid                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id             INT NOT NULL,
    author_uuid         UUID NOT NULL,
    address             VARCHAR(255),
    publication_date    DATE,
    distance            NUMERIC,
    FOREIGN KEY (city_id) REFERENCES cities(id),
    FOREIGN KEY (author_uuid) REFERENCES users(uuid)
);

CREATE TABLE IF NOT EXISTS ad_available_dates (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    available_date  DATE,
    FOREIGN KEY (ad_id) REFERENCES ads(uuid)
);

CREATE TABLE IF NOT EXISTS ad_positions (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    latitude        NUMERIC,
    longitude       NUMERIC,
    FOREIGN KEY (ad_id) REFERENCES ads(uuid)
);

CREATE TABLE IF NOT EXISTS images (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    image_url       TEXT CHECK (char_length(image_url) <= 1000),
    FOREIGN KEY (ad_id) REFERENCES ads(uuid)
);

CREATE TABLE IF NOT EXISTS requests (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    user_id         UUID NOT NULL,
    status          VARCHAR(255) DEFAULT 'pending' NOT NULL,
    create_date     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_date     TIMESTAMPTZ,
    close_date      TIMESTAMPTZ,
    FOREIGN KEY (ad_id) REFERENCES ads(uuid),
    FOREIGN KEY (user_id) REFERENCES users(uuid)
);

CREATE TABLE IF NOT EXISTS reviews (
    id              SERIAL PRIMARY KEY,
    user_id         UUID NOT NULL,
    host_id         UUID NOT NULL,
    text            TEXT NOT NULL CHECK (char_length(text) <= 1000),
    rating          NUMERIC,
    FOREIGN KEY (user_id) REFERENCES users(uuid),
    FOREIGN KEY (host_id) REFERENCES users(uuid)
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