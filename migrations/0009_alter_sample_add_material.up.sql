-- Add the material column to the sample material 
ALTER TABLE samples ADD COLUMN material character varying DEFAULT 'NA';