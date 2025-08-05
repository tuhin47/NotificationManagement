CREATE TABLE IF NOT EXISTS public.additional_fields
(
    id            bigserial,
    property_name varchar(100),
    type          varchar(10),
    description   text,
    request_id    bigint,
    PRIMARY KEY (id),
    CONSTRAINT fk_curl_requests_additional_fields
        FOREIGN KEY (request_id) REFERENCES public.curl_requests
            ON UPDATE CASCADE ON DELETE CASCADE
);

