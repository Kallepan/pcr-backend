CREATE TABLE analyses (
    analysis_id SERIAL PRIMARY KEY,

    analyt VARCHAR(5),
    material VARCHAR(50),
    assay VARCHAR(50),
    ready_mix BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT unique_analyt_material_assay UNIQUE (analyt, material, assay)
);

CREATE TABLE samples (
    sample_id VARCHAR(12) PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sampleanalyses (
    sample_id VARCHAR(12) REFERENCES samples(sample_id) ON UPDATE CASCADE ON DELETE CASCADE,
    analysis_id INTEGER REFERENCES analyses(analysis_id) ON UPDATE CASCADE ON DELETE CASCADE,

    run VARCHAR(20) NOT NULL,
    device VARCHAR(20) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT sample_analysis_pk PRIMARY KEY (sample_id, analysis_id) -- composite primary key
);

-- creates an index of certain tables to speed up queries
CREATE INDEX idx_sampleanalyses_analysis_id ON sampleanalyses (analysis_id);
CREATE INDEX idx_sampleanalyses_sample_id ON sampleanalyses (sample_id);
CREATE INDEX idx_samples_sample_id_like ON samples (sample_id varchar_pattern_ops);

CREATE INDEX idx_analyses_analyt ON analyses (analyt);
CREATE INDEX idx_analyses_material ON analyses (material);
CREATE INDEX idx_analyses_assay ON analyses (assay);