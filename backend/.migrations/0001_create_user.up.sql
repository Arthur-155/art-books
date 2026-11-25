CREATE TABLE "user" (
    id                  UUID        PRIMARY KEY,
    email               TEXT        NOT NULL CHECK (email = lower(btrim(email)) AND length(email) BETWEEN 3 AND 254),
    username            TEXT        NOT NULL CHECK (username ~ '^[a-z0-9_]{3,30}$'),
    display_name        TEXT        NOT NULL CHECK (length(btrim(display_name)) BETWEEN 1 AND 60),
    password_hash       TEXT        NOT NULL,
    role                TEXT        NOT NULL DEFAULT 'user' CHECK (role IN ('user','admin')),
    email_confirmed_at  TIMESTAMPTZ NULL,
    terms_accepted_at   TIMESTAMPTZ NOT NULL,
    age_declared_at     TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_user_email UNIQUE (email),
    CONSTRAINT uq_user_username UNIQUE (username)
);
