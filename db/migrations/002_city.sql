-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS cities (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title           TEXT NOT NULL CONSTRAINT title_length CHECK (char_length(title) <= 100),
    description     TEXT CONSTRAINT description_length CHECK (char_length(description) <= 3000)
    );

---- create above / drop below ----

DROP TABLE IF EXISTS cities CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.