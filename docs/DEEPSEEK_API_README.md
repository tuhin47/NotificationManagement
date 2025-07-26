# DeepseekModel API Documentation

This document provides example `curl` requests for all DeepseekModel-related endpoints in the NotificationManagement system.

## Base URL
```shell
export base_url=http://localhost:8080
```

---

## 1. Create DeepseekModel
**POST** `/api/deepseek-model`

```shell
curl -X POST ${base_url}/api/deepseek-model \
  -H "Content-Type: application/json" \
  -d '{
    "name": "deepseek-r1:1.5b",
    "model": "deepseek-r1:1.5b",
    "modified_at": "2025-07-20T07:44:29.156565162Z",
    "size": 1117322768
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "deepseek-r1:1.5b",
  "model": "deepseek-r1:1.5b",
  "modified_at": "2025-07-20T07:44:29.156565162Z",
  "size": 1117322768,
  "created_at": "2025-07-20T08:00:00Z",
  "updated_at": "2025-07-20T08:00:00Z"
}
```

---

## 2. Get DeepseekModel by ID
**GET** `/api/deepseek-model/:id`

```shell
curl ${base_url}/api/deepseek-model/1
```

**Response:**
```json
{
  "id": 1,
  "name": "deepseek-r1:1.5b",
  "model": "deepseek-r1:1.5b",
  "modified_at": "2025-07-20T07:44:29.156565162Z",
  "size": 1117322768,
  "created_at": "2025-07-20T08:00:00Z",
  "updated_at": "2025-07-20T08:00:00Z"
}
```

---

## 3. Get All DeepseekModels
**GET** `/api/deepseek-model?limit=10&offset=0`

```shell
curl "${base_url}/api/deepseek-model?limit=10&offset=0"
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "deepseek-r1:1.5b",
    "model": "deepseek-r1:1.5b",
    "modified_at": "2025-07-20T07:44:29.156565162Z",
    "size": 1117322768,
    "created_at": "2025-07-20T08:00:00Z",
    "updated_at": "2025-07-20T08:00:00Z"
  }
]
```

---

## 4. Update DeepseekModel
**PUT** `/api/deepseek-model/:id`

```shell
curl -X PUT ${base_url}/api/deepseek-model/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "deepseek-r1:1.5b-v2",
    "model": "deepseek-r1:1.5b-v2",
    "modified_at": "2025-08-01T10:00:00.000000000Z",
    "size": 1200000000
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "deepseek-r1:1.5b-v2",
  "model": "deepseek-r1:1.5b-v2",
  "modified_at": "2025-08-01T10:00:00.000000000Z",
  "size": 1200000000,
  "created_at": "2025-07-20T08:00:00Z",
  "updated_at": "2025-08-01T10:00:00Z"
}
```

---

## 5. Delete DeepseekModel
**DELETE** `/api/deepseek-model/:id`

```shell
curl -X DELETE ${base_url}/api/deepseek-model/1
```

**Response:**
```json
{
  "message": "DeepseekModel deleted successfully"
}
``` 