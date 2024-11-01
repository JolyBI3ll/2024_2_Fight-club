-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS ads (
    uuid                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id                   INT NOT NULL,
    author_uuid               UUID NOT NULL,
    address                   VARCHAR(255),
    publication_date          DATE,
    distance                  NUMERIC,
    FOREIGN KEY (city_id)     REFERENCES cities(id),
    FOREIGN KEY (author_uuid) REFERENCES users(uuid)
);

---- create above / drop below ----

DROP TABLE IF EXISTS ads CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.