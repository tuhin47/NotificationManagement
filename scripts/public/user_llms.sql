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
    model_name varchar(255) NOT NULL,
    type       varchar(50)
        CONSTRAINT chk_user_llms_type
            CHECK ((type)::text = ANY ((ARRAY ['local'::character varying, 'openai'::character varying, 'gemini'::character varying])::text[])),
    is_active  boolean DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_user_llms_model_name
    ON public.user_llms (model_name);

CREATE INDEX IF NOT EXISTS idx_user_llms_request_id
    ON public.user_llms (request_id);

CREATE INDEX IF NOT EXISTS idx_user_llms_deleted_at
    ON public.user_llms (deleted_at);

