CREATE TABLE IF NOT EXISTS public.users
(
    id          bigserial,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    keycloak_id text         NOT NULL,
    username    varchar(255) NOT NULL,
    email       varchar(255) NOT NULL,
    roles       text,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
    ON public.users (email);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_keycloak_id
    ON public.users (keycloak_id);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at
    ON public.users (deleted_at);

