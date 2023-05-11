CREATE TABLE IF NOT EXISTS analyses (
    analysis_id SERIAL PRIMARY KEY,

    analyt VARCHAR(5),
    material VARCHAR(50),
    assay VARCHAR(50),
    ready_mix BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT unique_analyt_material_assay UNIQUE (analyt, material, assay)
);

CREATE TABLE IF NOT EXISTS samples (
    sample_id VARCHAR(12) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,

    comment VARCHAR(255) DEFAULT '' NOT NULL,
    sputalysed BOOLEAN NOT NULL DEFAULT FALSE,

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS samplesanalyses (
    sample_id VARCHAR(12) REFERENCES samples(sample_id) ON UPDATE CASCADE ON DELETE CASCADE,
    analysis_id INTEGER REFERENCES analyses(analysis_id) ON UPDATE CASCADE ON DELETE CASCADE,

    run VARCHAR(20) DEFAULT '' NOT NULL,
    device VARCHAR(20) DEFAULT '' NOT NULL,
    position INTEGER DEFAULT NULL,
    
    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT sample_analysis_pk PRIMARY KEY (sample_id, analysis_id) -- composite primary key
);

-- creates an index of certain tables to speed up queries
CREATE INDEX IF NOT EXISTS idx_samplesanalyses_analysis_id ON samplesanalyses (analysis_id);
CREATE INDEX IF NOT EXISTS idx_samplesanalyses_sample_id ON samplesanalyses (sample_id);
CREATE INDEX IF NOT EXISTS idx_samples_sample_id_like ON samples (sample_id varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_analyses_analyt ON analyses (analyt);
CREATE INDEX IF NOT EXISTS idx_analyses_material ON analyses (material);
CREATE INDEX IF NOT EXISTS idx_analyses_assay ON analyses (assay);