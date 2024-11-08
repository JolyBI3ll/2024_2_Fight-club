-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS ad_positions (
        id              SERIAL PRIMARY KEY,
        ad_id           UUID NOT NULL,
        latitude        NUMERIC,
        longitude       NUMERIC,
        FOREIGN KEY (ad_id) REFERENCES ads(uuid)
);

---- create above / drop below ----

DROP TABLE IF EXISTS ad_positions CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.