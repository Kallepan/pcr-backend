CREATE TABLE IF NOT EXISTS samples (
    sample_id VARCHAR(12) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,

    comment VARCHAR(255) DEFAULT NULL,
    sputalysed BOOLEAN NOT NULL DEFAULT FALSE,

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS samplespanels (
    sample_id VARCHAR(12) REFERENCES samples(sample_id) ON UPDATE CASCADE ON DELETE CASCADE,
    panel_id VARCHAR(20) REFERENCES panels(panel_id) ON UPDATE CASCADE ON DELETE CASCADE,

    run VARCHAR(20) DEFAULT NULL,
    device VARCHAR(20) DEFAULT NULL,
    position INTEGER DEFAULT NULL,
    run_date DATE DEFAULT NULL,

    deleted BOOLEAN NOT NULL DEFAULT FALSE, -- Keep track of wether the sample-analysis pair was "deleted"

    created_by UUID REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT sample_analysis_pk PRIMARY KEY (sample_id, panel_id), -- composite primary key
    CONSTRAINT unique_position_created_at UNIQUE (position, run_date) -- unique postition in a run
);

-- creates an index of samplespanels
CREATE INDEX IF NOT EXISTS idx_samplespanels_panel_id ON samplespanels (panel_id);
CREATE INDEX IF NOT EXISTS idx_samplespanels_sample_id ON samplespanels (sample_id);
CREATE INDEX IF NOT EXISTS idx_samplespanels_run_date ON samplespanels (run_date);
CREATE INDEX IF NOT EXISTS idx_samplespanels_position ON samplespanels (position);
CREATE INDEX IF NOT EXISTS idx_samplespanels_complete ON samplespanels (run_date, position, device, run);
-- create index of samples 
CREATE INDEX IF NOT EXISTS idx_samples_sample_id_like ON samples (sample_id varchar_pattern_ops);

