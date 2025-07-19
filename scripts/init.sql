-- Table for storing CurlRequest
CREATE TABLE IF NOT EXISTS curl_requests (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    method VARCHAR(10) NOT NULL,
    headers TEXT,
    body TEXT,
    raw_curl TEXT,
    created_at BIGINT,
    updated_at BIGINT,
    deleted_at TIMESTAMP NULL
); 