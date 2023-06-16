DROP TABLE IF EXISTS samplesanalyses;
DROP TABLE IF EXISTS samples;
DROP TABLE IF EXISTS analyses;

DROP INDEX IF EXISTS idx_samplesanalyses_analysis_id;
DROP INDEX IF EXISTS idx_samplesanalyses_sample_id;
DROP INDEX IF EXISTS idx_samples_sample_id_like;

DROP INDEX IF EXISTS idx_analyses_analyt;
DROP INDEX IF EXISTS idx_analyses_material;
DROP INDEX IF EXISTS idx_analyses_assay;
