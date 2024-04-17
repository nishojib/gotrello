CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP(0),
    name TEXT NOT NULL,
    description TEXT,
    sort_order REAL NOT NULL DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 1,
    status_id UUID REFERENCES statuses(id)
);