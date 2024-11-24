-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS ad_available_dates (
      id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      ad_id           UUID NOT NULL,
      available_date  DATE,
      FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE IF EXISTS ad_available_dates CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.