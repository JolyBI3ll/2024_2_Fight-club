-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS requests (
    id              INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ad_id           UUID NOT NULL,
    user_id         UUID NOT NULL,
    status          TEXT DEFAULT 'pending' NOT NULL CONSTRAINT status_length CHECK (char_length(status) <= 255),
    create_date     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_date     TIMESTAMPTZ,
    close_date      TIMESTAMPTZ,
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

---- create above / drop below ----

DROP TABLE IF EXISTS requests CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.