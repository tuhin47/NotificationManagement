# Reminder API Documentation

This document provides example `curl` requests for all Reminder-related endpoints in the NotificationManagement system.

## Base URL
```
export base_url=http://localhost:8080/api/reminder
```

---

## 1. Create Reminder
**POST** `/api/reminder`

```
curl -X POST ${base_url}/api/reminder \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "message": "Take your medicine",
    "triggered_time": "2024-06-01T09:00:00Z",
    "next_trigger_time": "2024-06-02T09:00:00Z",
    "occurrence": 1,
    "recurrence": "daily"
  }'
```

**Response:**
```
{
  "id": 1,
  "request_id": 1,
  "message": "Take your medicine",
  "triggered_time": "2024-06-01T09:00:00Z",
  "next_trigger_time": "2024-06-02T09:00:00Z",
  "occurrence": 1,
  "recurrence": "daily",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

---

## 2. Get Reminder by ID
**GET** `/api/reminder/:id`

```
curl ${base_url}/api/reminder/1
```

**Response:**
```
{
  "id": 1,
  "request_id": 1,
  "message": "Take your medicine",
  "triggered_time": "2024-06-01T09:00:00Z",
  "next_trigger_time": "2024-06-02T09:00:00Z",
  "occurrence": 1,
  "recurrence": "daily",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

---

## 3. Get All Reminders
**GET** `/api/reminder?limit=10&offset=0`

```
curl "${base_url}/api/reminder?limit=10&offset=0"
```

**Response:**
```
[
  {
    "id": 1,
    "request_id": 1,
    "message": "Take your medicine",
    "triggered_time": "2024-06-01T09:00:00Z",
    "next_trigger_time": "2024-06-02T09:00:00Z",
    "occurrence": 1,
    "recurrence": "daily",
    "created_at": "2024-06-01T12:00:00Z",
    "updated_at": "2024-06-01T12:00:00Z"
  }
]
```

---

## 4. Update Reminder
**PUT** `/api/reminder/:id`

```
curl -X PUT ${base_url}/api/reminder/1 \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "message": "Take your vitamins",
    "triggered_time": "2024-06-01T09:00:00Z",
    "next_trigger_time": "2024-06-03T09:00:00Z",
    "occurrence": 2,
    "recurrence": "daily"
  }'
```

**Response:**
```
{
  "id": 1,
  "request_id": 1,
  "message": "Take your vitamins",
  "triggered_time": "2024-06-01T09:00:00Z",
  "next_trigger_time": "2024-06-03T09:00:00Z",
  "occurrence": 2,
  "recurrence": "daily",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:10:00Z"
}
```

---

## 5. Delete Reminder
**DELETE** `/api/reminder/:id`

```
curl -X DELETE ${base_url}/api/reminder/1
```

**Response:**
```
{
  "message": "Reminder deleted successfully"
}
``` 