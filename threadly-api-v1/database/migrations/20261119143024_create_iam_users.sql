-- +goose Up
CREATE TABLE iam_users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email         TEXT NOT NULL,
    username      TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    is_private    BOOLEAN NOT NULL DEFAULT FALSE,
    status        TEXT NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT iam_users_status_check
        CHECK (status IN ('active', 'deactivated', 'suspended'))
);

CREATE UNIQUE INDEX iam_users_email_lower_uidx
    ON iam_users (LOWER(email))
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX iam_users_username_lower_uidx
    ON iam_users (LOWER(username))
    WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE iam_users;
