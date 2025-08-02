CREATE TABLE IF NOT EXISTS public.additional_fields
(
    id            bigserial
        PRIMARY KEY,
    property_name varchar(100),
    type          varchar(10),
    description   text,
    request_id    bigint
        CONSTRAINT fk_curl_requests_additional_fields
            REFERENCES public.curl_requests
            ON UPDATE CASCADE ON DELETE CASCADE
);

