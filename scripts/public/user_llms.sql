CREATE TABLE IF NOT EXISTS public.user_llms
(
    id          bigserial
        PRIMARY KEY,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    request_id  bigint
        CONSTRAINT fk_curl_requests_ll_ms
            REFERENCES public.curl_requests,
    is_active   boolean DEFAULT TRUE,
    ai_model_id bigint
        CONSTRAINT fk_user_llms_ai_model
            REFERENCES public.ai_models
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_request_ai_model
    ON public.user_llms (request_id, ai_model_id);

CREATE INDEX IF NOT EXISTS idx_user_llms_deleted_at
    ON public.user_llms (deleted_at);

