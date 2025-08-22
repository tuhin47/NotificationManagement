CREATE TABLE IF NOT EXISTS public.curl_requests
(
    id            bigserial,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    url           text,
    method        varchar(10),
    headers       text,
    body          text,
    raw_curl      text,
    response_type varchar(10),
    user_id       bigint,
    PRIMARY KEY (id),
    CONSTRAINT fk_curl_requests_user
        FOREIGN KEY (user_id) REFERENCES public.users
);

CREATE INDEX IF NOT EXISTS idx_curl_requests_deleted_at
    ON public.curl_requests (deleted_at);

