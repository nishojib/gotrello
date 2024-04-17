CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP(0),
    name TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);