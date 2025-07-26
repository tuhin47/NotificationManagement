# LLM API Documentation

This document provides example `curl` requests for all LLM-related endpoints in the NotificationManagement system.

## Base URL
```shell
export base_url=http://localhost:8080
```

---

## 1. Create LLM
**POST** `/api/llm`

```shell
curl -X POST ${base_url}/api/llm \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "model_name": "gpt-4",
    "type": "openai",
    "is_active": true
  }'
```

**Response:**
```json
{
  "id": 1,
  "request_id": 1,
  "model_name": "gpt-4",
  "type": "openai",
  "is_active": true,
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

---

## 2. Get LLM by ID
**GET** `/api/llm/:id`

```shell
curl ${base_url}/api/llm/1
```

**Response:**
```json
{
  "id": 1,
  "request_id": 1,
  "model_name": "gpt-4",
  "type": "openai",
  "is_active": true,
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

---

## 3. Get All LLMs
**GET** `/api/llm?limit=10&offset=0`

```shell
curl "${base_url}/api/llm?limit=10&offset=0"
```

**Response:**
```json
[
  {
    "id": 1,
    "request_id": 1,
    "model_name": "gpt-4",
    "type": "openai",
    "is_active": true,
    "created_at": "2024-06-01T12:00:00Z",
    "updated_at": "2024-06-01T12:00:00Z"
  }
]
```

---

## 4. Update LLM
**PUT** `/api/llm/:id`

```shell
curl -X PUT ${base_url}/api/llm/1 \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "model_name": "gpt-4-turbo",
    "type": "openai",
    "is_active": false
  }'
```

**Response:**
```json
{
  "id": 1,
  "request_id": 1,
  "model_name": "gpt-4-turbo",
  "type": "openai",
  "is_active": false,
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:10:00Z"
}
```

---

## 5. Delete LLM
**DELETE** `/api/llm/:id`

```shell
curl -X DELETE ${base_url}/api/llm/1
```

**Response:**
```json
{
  "message": "LLM deleted successfully"
}
``` 