CREATE TABLE IF NOT EXISTS breeds (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    description text NOT NULL,
    avg_cost integer NOT NULL,
    countries text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);