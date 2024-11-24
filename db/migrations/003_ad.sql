-- Write your migrate up statements here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ads (
    id              UUID DEFAULT uuid_generate_v1mc() PRIMARY KEY,
    city_id         INT NOT NULL,
    author_id       UUID NOT NULL,
    address         TEXT CONSTRAINT address_length CHECK (char_length(address) <= 255),
    publication_date DATE,
    distance        NUMERIC,
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
    );

---- create above / drop below ----

DROP TABLE IF EXISTS ads CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.