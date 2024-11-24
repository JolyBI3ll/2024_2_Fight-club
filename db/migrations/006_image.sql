-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS images (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    image_url       TEXT CONSTRAINT image_url_length CHECK (char_length(image_url) <= 1000),
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE IF EXISTS images CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.