CREATE TABLE IF NOT EXISTS render_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_type VARCHAR(50) NOT NULL,
    resolution VARCHAR(20) NOT NULL,
    lighting_profile VARCHAR(100),
    camera_effect VARCHAR(100),
    audio_sensitivity VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- We create an index on the status column because our Go backend workers 
-- will constantly be querying the database specifically for 'queued' jobs.
CREATE INDEX idx_render_jobs_status ON render_jobs(status);