CREATE TABLE IF NOT EXISTS public.ai_models
(
    id          bigserial
        PRIMARY KEY,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    type        varchar(50)
        CONSTRAINT chk_ai_models_type
            CHECK ((type)::text = ANY (ARRAY [('local'::character varying)::text, ('openai'::character varying)::text, ('gemini'::character varying)::text])),
    name        varchar(255) NOT NULL,
    model_name  varchar(255) NOT NULL
        CONSTRAINT chk_ai_models_model_name
            CHECK ((model_name)::text <> ''::text),
    modified_at varchar(50),
    size        bigint       NOT NULL,
    base_url    varchar(500)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_models_model_name
    ON public.ai_models (model_name);

CREATE INDEX IF NOT EXISTS idx_ai_models_deleted_at
    ON public.ai_models (deleted_at);

