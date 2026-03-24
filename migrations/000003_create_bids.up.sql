CREATE TYPE bid_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE bids (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id     UUID        NOT NULL REFERENCES tasks(id),
    executor_id UUID        NOT NULL REFERENCES users(id),
    comment     TEXT        NOT NULL DEFAULT '',
    status      bid_status  NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (task_id, executor_id)
);

CREATE INDEX idx_bids_task_id     ON bids(task_id);
CREATE INDEX idx_bids_executor_id ON bids(executor_id);