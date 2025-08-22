CREATE TABLE IF NOT EXISTS public.telegrams
(
    id         bigserial,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id    bigint,
    chat_id    bigint       NOT NULL,
    otp        varchar(255) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_users_telegram
        FOREIGN KEY (user_id) REFERENCES public.users
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_telegrams_chat_id
    ON public.telegrams (chat_id);

CREATE INDEX IF NOT EXISTS idx_telegrams_user_id
    ON public.telegrams (user_id);

CREATE INDEX IF NOT EXISTS idx_telegrams_deleted_at
    ON public.telegrams (deleted_at);

