-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS requests (
    id              SERIAL PRIMARY KEY,
    ad_id           UUID NOT NULL,
    user_id         UUID NOT NULL,
    status          VARCHAR(255) DEFAULT 'pending' NOT NULL,
    create_date     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_date     TIMESTAMPTZ,
    close_date      TIMESTAMPTZ,
    FOREIGN KEY (ad_id) REFERENCES ads(uuid),
    FOREIGN KEY (user_id) REFERENCES users(uuid)
);

---- create above / drop below ----

DROP TABLE IF EXISTS requests CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.