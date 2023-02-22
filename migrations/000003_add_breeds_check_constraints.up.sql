ALTER TABLE breeds ADD CONSTRAINT breeds_avg_cost_check CHECK (avg_cost >= 0);
ALTER TABLE breeds ADD CONSTRAINT countries_length_check CHECK (array_length(countries, 1) BETWEEN 1 AND 5);