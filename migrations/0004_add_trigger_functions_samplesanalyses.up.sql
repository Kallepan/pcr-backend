CREATE OR REPLACE FUNCTION update_position()
RETURNS TRIGGER AS $$
-- Update the position of the samplepanel if it is not set and run and device are set by taking the max position of the day and adding 1 with a default of 1
BEGIN
    IF NEW.run IS NOT NULL AND NEW.device IS NOT NULL AND NEW.run_date IS NULL THEN
        IF NEW.position IS NULL THEN
            NEW.run_date := CURRENT_DATE;
            NEW.position := (
                SELECT COALESCE(MAX(position), 0) + 1
                FROM samplespanels
                WHERE
                    DATE(run_date) = CURRENT_DATE AND
                    position IS NOT NULL
            );
        END IF;
    END IF;
RETURN NEW;
END;


$$ LANGUAGE plpgsql;

CREATE TRIGGER update_position
BEFORE UPDATE ON samplespanels
FOR EACH ROW EXECUTE PROCEDURE update_position();