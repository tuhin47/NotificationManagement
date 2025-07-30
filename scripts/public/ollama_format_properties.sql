CREATE TABLE IF NOT EXISTS public.ollama_format_properties
(
    id            bigserial
        PRIMARY KEY,
    property_name varchar(100),
    type          varchar(10),
    description   text,
    request_id    bigint
        CONSTRAINT fk_curl_requests_ollama_format_properties
            REFERENCES public.curl_requests
);

