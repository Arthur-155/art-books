CREATE TABLE user_session (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    token_hash  BYTEA       NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_user_session_token_hash UNIQUE (token_hash)
);
CREATE INDEX ix_user_session_user_id ON user_session (user_id);
CREATE INDEX ix_user_session_expires_at ON user_session (expires_at);
