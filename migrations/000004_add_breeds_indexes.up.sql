CREATE INDEX IF NOT EXISTS breeds_name_idx ON breeds USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS breeds_countries_idx ON breeds USING GIN (countries);