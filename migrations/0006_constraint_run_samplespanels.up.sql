-- check if run_date, run, device, position are all null or all not null
ALTER TABLE samplespanels
ADD CONSTRAINT check_run_date_run_device_position CHECK (
    (run_date IS NULL AND run IS NULL AND device IS NULL AND position IS NULL) OR
    (run_date IS NOT NULL AND run IS NOT NULL AND device IS NOT NULL AND position IS NOT NULL)
);