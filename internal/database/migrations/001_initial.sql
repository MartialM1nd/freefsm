-- FreeFSM Initial Schema

-- Users (technicians/admins)
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    name            TEXT NOT NULL,
    phone           TEXT,
    role            TEXT NOT NULL DEFAULT 'technician',
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'technician'))
);

-- Sessions
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           TEXT UNIQUE NOT NULL,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Customers
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    email           TEXT,
    phone           TEXT,
    address         TEXT,
    city            TEXT,
    state           TEXT,
    zip             TEXT,
    notes           TEXT,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Jobs (work orders)
CREATE TABLE jobs (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id             UUID REFERENCES customers(id),
    assigned_to             UUID REFERENCES users(id),
    
    title                   TEXT NOT NULL,
    description             TEXT,
    status                  TEXT NOT NULL DEFAULT 'new',
    priority                TEXT NOT NULL DEFAULT 'medium',
    
    scheduled_date          DATE,
    scheduled_time          TIME,
    estimated_duration      INTEGER, -- minutes
    completed_at            TIMESTAMP,
    
    use_customer_address    BOOLEAN DEFAULT true,
    location_address        TEXT,
    location_city           TEXT,
    location_state          TEXT,
    location_zip            TEXT,
    
    deleted_at              TIMESTAMP,
    created_at              TIMESTAMP DEFAULT NOW(),
    updated_at              TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT jobs_status_check CHECK (status IN (
        'new', 'in_transit', 'in_progress', 'pending', 
        'scheduled_return', 'ready_to_invoice', 'completed', 'cancelled'
    )),
    CONSTRAINT jobs_priority_check CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

CREATE INDEX idx_jobs_customer_id ON jobs(customer_id);
CREATE INDEX idx_jobs_assigned_to ON jobs(assigned_to);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_scheduled_date ON jobs(scheduled_date);

-- Job Notes
CREATE TABLE job_notes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id          UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    content         TEXT NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_job_notes_job_id ON job_notes(job_id);

-- Job History (audit trail)
CREATE TABLE job_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id          UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    changed_by      UUID REFERENCES users(id),
    field           TEXT NOT NULL,
    old_value       TEXT,
    new_value       TEXT,
    changed_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_job_history_job_id ON job_history(job_id);

-- Seed admin user
-- Password: changeme123 (bcrypt hash)
INSERT INTO users (email, password_hash, name, role) VALUES (
    'admin@freefsm.local',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqBuBhRlN5qvpCvzuXlPHBnKGKxHi',
    'Administrator',
    'admin'
);
