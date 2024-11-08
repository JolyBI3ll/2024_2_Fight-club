-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS images (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    image_url       TEXT CHECK (char_length(image_url) <= 1000),
    FOREIGN KEY (ad_id) REFERENCES ads(uuid)
);

---- create above / drop below ----

DROP TABLE IF EXISTS images CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.