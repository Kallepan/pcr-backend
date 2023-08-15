CREATE OR REPLACE FUNCTION update_position() RETURNS TRIGGER AS $$
-- Update the position of the samplepanel if it is not set and run and device are set by taking the max position of the day and adding 1 with a default of 1
BEGIN
    -- run, device and run_date are provided
    IF NEW.run IS NOT NULL AND NEW.device IS NOT NULL AND NEW.run_date IS NOT NULL AND NEW.position IS NULL THEN
        NEW.position := (
            SELECT COALESCE(MAX(position), 0) + 1
            FROM samplespanels
            WHERE
                DATE(run_date) = NEW.run_date AND
                position IS NOT NULL
        );
    END IF;

    -- run and device are provided but run_date is not
    IF NEW.run IS NOT NULL AND NEW.device IS NOT NULL AND NEW.run_date IS NULL AND NEW.position IS NULL THEN
        NEW.run_date := CURRENT_DATE;
        NEW.position := (
            SELECT COALESCE(MAX(position), 0) + 1
            FROM samplespanels
            WHERE
                DATE(run_date) = CURRENT_DATE AND
                position IS NOT NULL
        );
    END IF;
RETURN NEW;
END;

$$ LANGUAGE plpgsql;
