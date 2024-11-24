-- Write your migrate up statements here

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

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.