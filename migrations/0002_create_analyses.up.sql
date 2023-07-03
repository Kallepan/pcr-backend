CREATE TABLE IF NOT EXISTS panels (
    panel_id VARCHAR(20) PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    ready_mix BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS analyses (
    analysis_id VARCHAR(20) PRIMARY KEY NOT NULL,
    panel_id VARCHAR(20) NOT NULL REFERENCES panels(panel_id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Here we have a foreign key relationship between analyses and panels. Analyses can only be created if the panel exists.

-- creates an index of panels
CREATE INDEX IF NOT EXISTS idx_panels_panel_id ON panels (panel_id);
CREATE INDEX IF NOT EXISTS idx_panels_ready_mix ON panels (ready_mix);

-- creates an index of analyses
CREATE INDEX IF NOT EXISTS idx_analyses_panel_id ON analyses (panel_id, analysis_id);
CREATE INDEX IF NOT EXISTS idx_analyses_analyt_id ON analyses (analysis_id);