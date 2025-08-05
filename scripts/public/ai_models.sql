CREATE TABLE IF NOT EXISTS public.ai_models
(
    id         bigserial,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type       varchar(10),
    name       varchar(255) NOT NULL,
    model_name varchar(255) NOT NULL,
    base_url   varchar(500),
    size       bigint,
    api_secret varchar(500),
    PRIMARY KEY (id),
    CONSTRAINT chk_ai_models_model_name
        CHECK ((model_name)::text <> ''::text),
    CONSTRAINT chk_ai_models_type
        CHECK ((type)::text = ANY (ARRAY [('local'::character varying)::text, ('openai'::character varying)::text, ('gemini'::character varying)::text, ('deepseek'::character varying)::text]))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_model_model_secret
    ON public.ai_models (model_name, api_secret);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_model_model_url
    ON public.ai_models (model_name, base_url);

CREATE INDEX IF NOT EXISTS idx_ai_models_deleted_at
    ON public.ai_models (deleted_at);

