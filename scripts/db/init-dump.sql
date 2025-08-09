--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Debian 15.13-1.pgdg120+1)
-- Dumped by pg_dump version 15.13 (Debian 15.13-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.telegrams DROP CONSTRAINT IF EXISTS fk_users_telegram;
ALTER TABLE IF EXISTS ONLY public.request_ai_models DROP CONSTRAINT IF EXISTS fk_request_ai_models_ai_model;
ALTER TABLE IF EXISTS ONLY public.curl_requests DROP CONSTRAINT IF EXISTS fk_curl_requests_user;
ALTER TABLE IF EXISTS ONLY public.reminders DROP CONSTRAINT IF EXISTS fk_curl_requests_reminders;
ALTER TABLE IF EXISTS ONLY public.request_ai_models DROP CONSTRAINT IF EXISTS fk_curl_requests_models;
ALTER TABLE IF EXISTS ONLY public.additional_fields DROP CONSTRAINT IF EXISTS fk_curl_requests_additional_fields;
DROP INDEX IF EXISTS public.idx_users_keycloak_id;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_users_deleted_at;
DROP INDEX IF EXISTS public.idx_telegrams_user_id;
DROP INDEX IF EXISTS public.idx_telegrams_deleted_at;
DROP INDEX IF EXISTS public.idx_telegrams_chat_id;
DROP INDEX IF EXISTS public.idx_request_ai_models_deleted_at;
DROP INDEX IF EXISTS public.idx_request_ai_model;
DROP INDEX IF EXISTS public.idx_reminders_upto;
DROP INDEX IF EXISTS public.idx_reminders_triggered_time;
DROP INDEX IF EXISTS public.idx_reminders_next_trigger_time;
DROP INDEX IF EXISTS public.idx_reminders_deleted_at;
DROP INDEX IF EXISTS public.idx_curl_requests_deleted_at;
DROP INDEX IF EXISTS public.idx_ai_models_deleted_at;
DROP INDEX IF EXISTS public.idx_ai_model_model_url;
DROP INDEX IF EXISTS public.idx_ai_model_model_secret;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.telegrams DROP CONSTRAINT IF EXISTS telegrams_pkey;
ALTER TABLE IF EXISTS ONLY public.request_ai_models DROP CONSTRAINT IF EXISTS request_ai_models_pkey;
ALTER TABLE IF EXISTS ONLY public.reminders DROP CONSTRAINT IF EXISTS reminders_pkey;
ALTER TABLE IF EXISTS ONLY public.curl_requests DROP CONSTRAINT IF EXISTS curl_requests_pkey;
ALTER TABLE IF EXISTS ONLY public.ai_models DROP CONSTRAINT IF EXISTS ai_models_pkey;
ALTER TABLE IF EXISTS ONLY public.additional_fields DROP CONSTRAINT IF EXISTS additional_fields_pkey;
ALTER TABLE IF EXISTS public.users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.telegrams ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.request_ai_models ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.reminders ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.curl_requests ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.ai_models ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.additional_fields ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.users;
DROP SEQUENCE IF EXISTS public.telegrams_id_seq;
DROP TABLE IF EXISTS public.telegrams;
DROP SEQUENCE IF EXISTS public.request_ai_models_id_seq;
DROP TABLE IF EXISTS public.request_ai_models;
DROP SEQUENCE IF EXISTS public.reminders_id_seq;
DROP TABLE IF EXISTS public.reminders;
DROP SEQUENCE IF EXISTS public.curl_requests_id_seq;
DROP TABLE IF EXISTS public.curl_requests;
DROP SEQUENCE IF EXISTS public.ai_models_id_seq;
DROP TABLE IF EXISTS public.ai_models;
DROP SEQUENCE IF EXISTS public.additional_fields_id_seq;
DROP TABLE IF EXISTS public.additional_fields;
DROP SCHEMA IF EXISTS public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: additional_fields; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.additional_fields (
    id bigint NOT NULL,
    property_name character varying(100),
    type character varying(10),
    description text,
    request_id bigint
);


ALTER TABLE public.additional_fields OWNER TO "user";

--
-- Name: additional_fields_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.additional_fields_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.additional_fields_id_seq OWNER TO "user";

--
-- Name: additional_fields_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.additional_fields_id_seq OWNED BY public.additional_fields.id;


--
-- Name: ai_models; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.ai_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type character varying(10),
    name character varying(255) NOT NULL,
    model_name character varying(255) NOT NULL,
    base_url character varying(500),
    size bigint,
    api_secret character varying(500),
    CONSTRAINT chk_ai_models_model_name CHECK (((model_name)::text <> ''::text)),
    CONSTRAINT chk_ai_models_type CHECK (((type)::text = ANY (ARRAY[('local'::character varying)::text, ('openai'::character varying)::text, ('gemini'::character varying)::text, ('deepseek'::character varying)::text])))
);


ALTER TABLE public.ai_models OWNER TO "user";

--
-- Name: ai_models_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.ai_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ai_models_id_seq OWNER TO "user";

--
-- Name: ai_models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.ai_models_id_seq OWNED BY public.ai_models.id;


--
-- Name: curl_requests; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.curl_requests (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    url text,
    method character varying(10),
    headers text,
    body text,
    raw_curl text,
    response_type character varying(10),
    user_id bigint
);


ALTER TABLE public.curl_requests OWNER TO "user";

--
-- Name: curl_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.curl_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.curl_requests_id_seq OWNER TO "user";

--
-- Name: curl_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.curl_requests_id_seq OWNED BY public.curl_requests.id;


--
-- Name: reminders; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.reminders (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    request_id bigint,
    message text NOT NULL,
    triggered_time timestamp with time zone,
    next_trigger_time timestamp with time zone,
    occurrence bigint DEFAULT 0,
    recurrence character varying(50) NOT NULL,
    after_every bigint NOT NULL,
    task_id text,
    upto timestamp with time zone
);


ALTER TABLE public.reminders OWNER TO "user";

--
-- Name: reminders_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.reminders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reminders_id_seq OWNER TO "user";

--
-- Name: reminders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.reminders_id_seq OWNED BY public.reminders.id;


--
-- Name: request_ai_models; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.request_ai_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    request_id bigint,
    is_active boolean DEFAULT true,
    ai_model_id bigint
);


ALTER TABLE public.request_ai_models OWNER TO "user";

--
-- Name: request_ai_models_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.request_ai_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.request_ai_models_id_seq OWNER TO "user";

--
-- Name: request_ai_models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.request_ai_models_id_seq OWNED BY public.request_ai_models.id;


--
-- Name: telegrams; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.telegrams (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint,
    chat_id bigint NOT NULL,
    otp character varying(255) NOT NULL
);


ALTER TABLE public.telegrams OWNER TO "user";

--
-- Name: telegrams_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.telegrams_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.telegrams_id_seq OWNER TO "user";

--
-- Name: telegrams_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.telegrams_id_seq OWNED BY public.telegrams.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    keycloak_id text NOT NULL,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    roles text
);


ALTER TABLE public.users OWNER TO "user";

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO "user";

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: additional_fields id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.additional_fields ALTER COLUMN id SET DEFAULT nextval('public.additional_fields_id_seq'::regclass);


--
-- Name: ai_models id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.ai_models ALTER COLUMN id SET DEFAULT nextval('public.ai_models_id_seq'::regclass);


--
-- Name: curl_requests id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.curl_requests ALTER COLUMN id SET DEFAULT nextval('public.curl_requests_id_seq'::regclass);


--
-- Name: reminders id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.reminders ALTER COLUMN id SET DEFAULT nextval('public.reminders_id_seq'::regclass);


--
-- Name: request_ai_models id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.request_ai_models ALTER COLUMN id SET DEFAULT nextval('public.request_ai_models_id_seq'::regclass);


--
-- Name: telegrams id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.telegrams ALTER COLUMN id SET DEFAULT nextval('public.telegrams_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: additional_fields; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.additional_fields (id, property_name, type, description, request_id) FROM stdin;
1	CurrentRate	number	The current rate from json	1
2	TargetedRate	number	The target rate from statement	1
6	TargetBalace	number	The target balance from statement	3
5	CurrentBalace	number	Current Balance From Json	3
7	ProductName	text	Product Name from html	4
8	Price	number	The price of the product from html	4
10	TargetedRate	number	The target rate from statement	5
9	CurrentRate	number	The current rate from json	5
3	RainProbablity	text	Any Probablity For rain	2
4	MaximumTemperature	text	Maximum Temperature	2
\.


--
-- Data for Name: ai_models; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.ai_models (id, created_at, updated_at, deleted_at, type, name, model_name, base_url, size, api_secret) FROM stdin;
1	2025-08-10 15:44:37.351817+00	2025-08-10 15:44:37.351817+00	\N	gemini	gemini-2.5-flash	gemini-2.5-flash	\N	\N	LS/AYXK4bMT+/L4MCibeCgIOm6AWIpSIZbqgvbKpY8wVtBNkcaZN0Iqa7Q==
\.


--
-- Data for Name: curl_requests; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.curl_requests (id, created_at, updated_at, deleted_at, url, method, headers, body, raw_curl, response_type, user_id) FROM stdin;
1	2025-08-10 15:37:24.435463+00	2025-08-10 15:38:16.088082+00	\N			null	Please check the current rate from the json.Is it greater than 120?	curl "https://api.elevatepay.co/api/v1/fxRate/webSite/transfer?quote=BDT&base=USD&amount=1296.33" -H "accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7" -H "accept-language: en-US,en;q=0.9,bn;q=0.8" -H "cache-control: max-age=0" -H "priority: u=0, i" -H "referer: https://www.google.com/" -H "sec-ch-ua: \\"Not)A;Brand\\";v=\\"8\\", \\"Chromium\\";v=\\"138\\", \\"Microsoft Edge\\";v=\\"138\\"" -H "sec-ch-ua-mobile: ?0" -H "sec-ch-ua-platform: \\"Linux\\"" -H "sec-fetch-dest: document" -H "sec-fetch-mode: navigate" -H "sec-fetch-site: cross-site" -H "sec-fetch-user: ?1" -H "upgrade-insecure-requests: 1" -H "user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0"	json	1
3	2025-08-10 15:37:37.687687+00	2025-08-10 15:40:25.304856+00	\N			null	The balance is below 400?	curl \\'https://prepaid.desco.org.bd/api/unified/customer/getBalance?accountNo=41036730\\'  --insecure	json	1
4	2025-08-10 15:37:40.521678+00	2025-08-10 15:41:43.568686+00	\N			null	Please check the html.Product seems not available	curl -s https://www.startech.com.bd/huawei-mateview-ssn-24bz-se-monitor	html	1
5	2025-08-10 15:42:09.539591+00	2025-08-10 15:43:25.39805+00	\N			null	Within 2 days we got the latest model	curl -s https://ollama.com/library/gemma3n/tags	html	1
2	2025-08-10 15:37:32.454495+00	2025-08-10 15:44:08.102165+00	\N			null	Is there any rain propbaliblity for next 2 hours?	curl \\'https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m\\' 	json	1
\.


--
-- Data for Name: reminders; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.reminders (id, created_at, updated_at, deleted_at, request_id, message, triggered_time, next_trigger_time, occurrence, recurrence, after_every, task_id, upto) FROM stdin;
\.


--
-- Data for Name: request_ai_models; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.request_ai_models (id, created_at, updated_at, deleted_at, request_id, is_active, ai_model_id) FROM stdin;
1	2025-08-10 15:45:37.058681+00	2025-08-10 15:45:37.058681+00	\N	1	t	1
2	2025-08-10 15:45:41.260756+00	2025-08-10 15:45:41.260756+00	\N	2	t	1
3	2025-08-10 15:45:43.990126+00	2025-08-10 15:45:43.990126+00	\N	3	t	1
4	2025-08-10 15:45:47.405311+00	2025-08-10 15:45:47.405311+00	\N	4	t	1
5	2025-08-10 15:45:51.069544+00	2025-08-10 15:45:51.069544+00	\N	5	t	1
\.


--
-- Data for Name: telegrams; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.telegrams (id, created_at, updated_at, deleted_at, user_id, chat_id, otp) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.users (id, created_at, updated_at, deleted_at, keycloak_id, username, email, roles) FROM stdin;
1	2025-08-10 15:37:24.433444+00	2025-08-10 15:46:01.375143+00	\N	07d9c9b3-248b-47c2-b05d-7238ad6938ff	nms	nms@tuhin47.com	curl_create,ai_model_create,ai_make_request,reminder_create,curl_read,llm_delete,curl_delete,ai_model_update,reminder_update,llm_update,ai_model_delete,reminder_delete,curl_update,llm_read,reminder_read,ai_model_read,llm_create
\.


--
-- Name: additional_fields_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.additional_fields_id_seq', 14, true);


--
-- Name: ai_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.ai_models_id_seq', 1, true);


--
-- Name: curl_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.curl_requests_id_seq', 5, true);


--
-- Name: reminders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.reminders_id_seq', 1, false);


--
-- Name: request_ai_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.request_ai_models_id_seq', 6, true);


--
-- Name: telegrams_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.telegrams_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: additional_fields additional_fields_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.additional_fields
    ADD CONSTRAINT additional_fields_pkey PRIMARY KEY (id);


--
-- Name: ai_models ai_models_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.ai_models
    ADD CONSTRAINT ai_models_pkey PRIMARY KEY (id);


--
-- Name: curl_requests curl_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.curl_requests
    ADD CONSTRAINT curl_requests_pkey PRIMARY KEY (id);


--
-- Name: reminders reminders_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.reminders
    ADD CONSTRAINT reminders_pkey PRIMARY KEY (id);


--
-- Name: request_ai_models request_ai_models_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.request_ai_models
    ADD CONSTRAINT request_ai_models_pkey PRIMARY KEY (id);


--
-- Name: telegrams telegrams_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.telegrams
    ADD CONSTRAINT telegrams_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_ai_model_model_secret; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_ai_model_model_secret ON public.ai_models USING btree (model_name, api_secret);


--
-- Name: idx_ai_model_model_url; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_ai_model_model_url ON public.ai_models USING btree (model_name, base_url);


--
-- Name: idx_ai_models_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_ai_models_deleted_at ON public.ai_models USING btree (deleted_at);


--
-- Name: idx_curl_requests_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_curl_requests_deleted_at ON public.curl_requests USING btree (deleted_at);


--
-- Name: idx_reminders_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_reminders_deleted_at ON public.reminders USING btree (deleted_at);


--
-- Name: idx_reminders_next_trigger_time; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_reminders_next_trigger_time ON public.reminders USING btree (next_trigger_time);


--
-- Name: idx_reminders_triggered_time; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_reminders_triggered_time ON public.reminders USING btree (triggered_time);


--
-- Name: idx_reminders_upto; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_reminders_upto ON public.reminders USING btree (upto);


--
-- Name: idx_request_ai_model; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_request_ai_model ON public.request_ai_models USING btree (request_id, ai_model_id);


--
-- Name: idx_request_ai_models_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_request_ai_models_deleted_at ON public.request_ai_models USING btree (deleted_at);


--
-- Name: idx_telegrams_chat_id; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_telegrams_chat_id ON public.telegrams USING btree (chat_id);


--
-- Name: idx_telegrams_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_telegrams_deleted_at ON public.telegrams USING btree (deleted_at);


--
-- Name: idx_telegrams_user_id; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_telegrams_user_id ON public.telegrams USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_keycloak_id; Type: INDEX; Schema: public; Owner: user
--

CREATE UNIQUE INDEX idx_users_keycloak_id ON public.users USING btree (keycloak_id);


--
-- Name: additional_fields fk_curl_requests_additional_fields; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.additional_fields
    ADD CONSTRAINT fk_curl_requests_additional_fields FOREIGN KEY (request_id) REFERENCES public.curl_requests(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: request_ai_models fk_curl_requests_models; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.request_ai_models
    ADD CONSTRAINT fk_curl_requests_models FOREIGN KEY (request_id) REFERENCES public.curl_requests(id);


--
-- Name: reminders fk_curl_requests_reminders; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.reminders
    ADD CONSTRAINT fk_curl_requests_reminders FOREIGN KEY (request_id) REFERENCES public.curl_requests(id);


--
-- Name: curl_requests fk_curl_requests_user; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.curl_requests
    ADD CONSTRAINT fk_curl_requests_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: request_ai_models fk_request_ai_models_ai_model; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.request_ai_models
    ADD CONSTRAINT fk_request_ai_models_ai_model FOREIGN KEY (ai_model_id) REFERENCES public.ai_models(id);


--
-- Name: telegrams fk_users_telegram; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.telegrams
    ADD CONSTRAINT fk_users_telegram FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

