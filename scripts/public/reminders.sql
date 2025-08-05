CREATE TABLE IF NOT EXISTS public.reminders
(
    id                bigserial,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    request_id        bigint,
    message           text NOT NULL,
    triggered_time    timestamp with time zone,
    next_trigger_time timestamp with time zone,
    occurrence        bigint DEFAULT 0,
    recurrence        varchar(50),
    PRIMARY KEY (id),
    CONSTRAINT fk_curl_requests_reminders
        FOREIGN KEY (request_id) REFERENCES public.curl_requests,
    CONSTRAINT chk_reminders_recurrence
        CHECK ((recurrence)::text = ANY
               (ARRAY [('once'::character varying)::text, ('minutes'::character varying)::text, ('hour'::character varying)::text, ('daily'::character varying)::text, ('weekly'::character varying)::text]))
);

CREATE INDEX IF NOT EXISTS idx_reminders_deleted_at
    ON public.reminders (deleted_at);

CREATE INDEX IF NOT EXISTS idx_reminders_next_trigger_time
    ON public.reminders (next_trigger_time);

CREATE INDEX IF NOT EXISTS idx_reminders_request_id
    ON public.reminders (request_id);

CREATE INDEX IF NOT EXISTS idx_reminders_triggered_time
    ON public.reminders (triggered_time);

