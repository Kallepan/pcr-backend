-- Adds the manual field to sample table
ALTER TABLE samples ADD COLUMN manual boolean NOT NULL DEFAULT false;
