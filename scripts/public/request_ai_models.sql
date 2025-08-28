CREATE TABLE IF NOT EXISTS public.request_ai_models
(
    id          bigserial,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    request_id  bigint,
    is_active   boolean DEFAULT TRUE,
    ai_model_id bigint,
    PRIMARY KEY (id),
    CONSTRAINT fk_request_ai_models_ai_model
        FOREIGN KEY (ai_model_id) REFERENCES public.ai_models,
    CONSTRAINT fk_curl_requests_models
        FOREIGN KEY (request_id) REFERENCES public.curl_requests
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_request_ai_model
    ON public.request_ai_models (request_id, ai_model_id);

CREATE INDEX IF NOT EXISTS idx_request_ai_models_deleted_at
    ON public.request_ai_models (deleted_at);

