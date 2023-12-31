DROP TABLE IF EXISTS samplespanels;
DROP TABLE IF EXISTS samples;

DROP INDEX IF EXISTS idx_samplespanels_panel_id;
DROP INDEX IF EXISTS idx_samplespanels_sample_id;
DROP INDEX IF EXISTS idx_samplespanels_run_date;
DROP INDEX IF EXISTS idx_samplespanels_position;
DROP INDEX IF EXISTS idx_samplespanels_complete;
DROP INDEX IF EXISTS idx_samples_sample_id_like;
DROP INDEX IF EXISTS idx_samplespanels_run;