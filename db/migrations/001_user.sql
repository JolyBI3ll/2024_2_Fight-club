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

---- create above / drop below ----

DROP TABLE IF EXISTS users CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.