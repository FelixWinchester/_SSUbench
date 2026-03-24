CREATE TYPE task_status AS ENUM (
    'draft',
    'published',
    'in_progress',
    'completed',
    'cancelled'
);

CREATE TABLE tasks (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID        NOT NULL REFERENCES users(id),
    executor_id UUID        REFERENCES users(id),
    title       VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL,
    reward      BIGINT      NOT NULL CHECK (reward > 0),
    status      task_status NOT NULL DEFAULT 'draft',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_customer_id ON tasks(customer_id);
CREATE INDEX idx_tasks_executor_id ON tasks(executor_id);
CREATE INDEX idx_tasks_status      ON tasks(status);