CREATE TABLE IF NOT EXISTS public.reminders
(
    id                bigserial
        PRIMARY KEY,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    request_id        bigint
        CONSTRAINT fk_curl_requests_reminders
            REFERENCES public.curl_requests,
    message           text NOT NULL,
    triggered_time    timestamp with time zone,
    next_trigger_time timestamp with time zone,
    recurrence        varchar(50)
        CONSTRAINT chk_reminders_recurrence
            CHECK ((recurrence)::text = ANY
                   ((ARRAY ['once'::character varying, 'minutes'::character varying, 'hour'::character varying, 'daily'::character varying, 'weekly'::character varying])::text[]))
);

CREATE INDEX IF NOT EXISTS idx_reminders_next_trigger_time
    ON public.reminders (next_trigger_time);

CREATE INDEX IF NOT EXISTS idx_reminders_triggered_time
    ON public.reminders (triggered_time);

CREATE INDEX IF NOT EXISTS idx_reminders_request_id
    ON public.reminders (request_id);

CREATE INDEX IF NOT EXISTS idx_reminders_deleted_at
    ON public.reminders (deleted_at);

