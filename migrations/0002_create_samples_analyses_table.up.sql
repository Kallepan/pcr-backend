CREATE TABLE IF NOT EXISTS analyses (
    analysis_id VARCHAR(20) PRIMARY KEY NOT NULL,

    ready_mix BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS samples (
    sample_id VARCHAR(12) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,

    comment VARCHAR(255) DEFAULT '' NOT NULL,
    sputalysed BOOLEAN NOT NULL DEFAULT FALSE,

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS samplesanalyses (
    sample_id VARCHAR(12) REFERENCES samples(sample_id) ON UPDATE CASCADE ON DELETE CASCADE,
    analysis_id VARCHAR(20) REFERENCES analyses(analysis_id) ON UPDATE CASCADE ON DELETE CASCADE,

    run VARCHAR(20) DEFAULT NULL,
    device VARCHAR(20) DEFAULT NULL,
    position INTEGER DEFAULT NULL,
    
    deleted BOOLEAN DEFAULT FALSE, -- Keep track of wether the sample-analysis pair was "deleted"

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT sample_analysis_pk PRIMARY KEY (sample_id, analysis_id), -- composite primary key
    CONSTRAINT unique_position_created_at UNIQUE (position, created_at) -- unique postition in a run
);

-- creates an index of certain tables to speed up queries
CREATE INDEX IF NOT EXISTS idx_samplesanalyses_analysis_id ON samplesanalyses (analysis_id);
CREATE INDEX IF NOT EXISTS idx_samplesanalyses_sample_id ON samplesanalyses (sample_id);
CREATE INDEX IF NOT EXISTS idx_samples_sample_id_like ON samples (sample_id varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_analyses_analyt_id ON analyses (analysis_id);
CREATE INDEX IF NOT EXISTS idx_analyses_is_active ON analyses (is_active);
CREATE INDEX IF NOT EXISTS idx_analyses_ready_mix ON analyses (ready_mix);
