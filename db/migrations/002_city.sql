-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS cities (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(100),
    description     TEXT CHECK (char_length(description) <= 3000)
);

---- create above / drop below ----

DROP TABLE IF EXISTS cities CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.