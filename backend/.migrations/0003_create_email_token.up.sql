CREATE TABLE email_token (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    purpose     TEXT        NOT NULL CHECK (purpose IN ('email_confirmation','password_reset','email_change_approval','email_change_confirmation')),
    token_hash  BYTEA       NOT NULL,
    new_email   TEXT        NULL CHECK (new_email = lower(btrim(new_email))),
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_email_token_user_id_purpose UNIQUE (user_id, purpose),
    CONSTRAINT uq_email_token_token_hash UNIQUE (token_hash),
    CONSTRAINT ck_email_token_new_email CHECK ((purpose LIKE 'email_change_%') = (new_email IS NOT NULL))
);
CREATE INDEX ix_email_token_expires_at ON email_token (expires_at);
