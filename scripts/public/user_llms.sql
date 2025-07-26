CREATE TABLE IF NOT EXISTS public.user_llms
(
    id         bigserial
        PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    request_id bigint
        CONSTRAINT fk_curl_requests_ll_ms
            REFERENCES public.curl_requests,
    is_active  boolean DEFAULT TRUE
);

-- Removed model_name and type columns, constraints, and indexes

CREATE INDEX IF NOT EXISTS idx_user_llms_request_id
    ON public.user_llms (request_id);

CREATE INDEX IF NOT EXISTS idx_user_llms_deleted_at
    ON public.user_llms (deleted_at);

