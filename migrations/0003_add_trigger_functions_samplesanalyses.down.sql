CREATE OR REPLACE FUNCTION update_position()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.run IS NOT NULL AND NEW.device IS NOT NULL THEN
        IF NEW.position IS NULL THEN
            NEW.position := (
                SELECT COALESCE(MAX(position) + 1, 1)
                FROM samplesanalyses
                WHERE 
                    DATE(created_at) = CURRENT_DATE AND
                    position IS NOT NULL
            );
        END IF;
    END IF;
RETURN NEW;
END;


$$ LANGUAGE plpgsql;

CREATE TRIGGER update_position
BEFORE UPDATE ON samplesanalyses
FOR EACH ROW EXECUTE PROCEDURE update_position();