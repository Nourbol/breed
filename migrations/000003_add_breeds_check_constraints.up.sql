ALTER TABLE breeds ADD CONSTRAINT breeds_avg_cost_check CHECK (avg_cost >= 0);