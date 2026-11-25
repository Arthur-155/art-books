CREATE TABLE login_lock (
    key_kind           TEXT        NOT NULL CHECK (key_kind IN ('email','ip')),
    key_value          TEXT        NOT NULL,
    failure_count      SMALLINT    NOT NULL DEFAULT 0,
    window_started_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_until       TIMESTAMPTZ NULL,
    PRIMARY KEY (key_kind, key_value)
);
CREATE TABLE email_request (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_type  TEXT        NOT NULL CHECK (request_type IN ('confirmation_resend','password_reset','email_change')),
    key_kind      TEXT        NOT NULL CHECK (key_kind IN ('email','ip')),
    key_value     TEXT        NOT NULL,
    requested_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ix_email_request_key ON email_request (request_type, key_kind, key_value, requested_at);
