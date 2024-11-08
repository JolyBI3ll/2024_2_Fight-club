-- Write your migrate up statements here

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

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.