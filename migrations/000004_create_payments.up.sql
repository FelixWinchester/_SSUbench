CREATE TABLE payments (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id      UUID        NOT NULL REFERENCES tasks(id),
    from_user_id UUID        NOT NULL REFERENCES users(id),
    to_user_id   UUID        NOT NULL REFERENCES users(id),
    amount       BIGINT      NOT NULL CHECK (amount > 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_task_id      ON payments(task_id);
CREATE INDEX idx_payments_from_user_id ON payments(from_user_id);
CREATE INDEX idx_payments_to_user_id   ON payments(to_user_id);