-- Write your migrate up statements here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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

---- create above / drop below ----

DROP TABLE IF EXISTS users CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.