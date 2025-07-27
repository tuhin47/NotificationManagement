# AI API Documentation

This document provides example `curl` requests for the AI-related endpoints in the NotificationManagement system.

## Base URL

```shell
export base_url=http://localhost:8080
```

---

## 1. Make AI Request

**POST** `/api/ai/make-request`

This endpoint allows you to make AI requests using Ollama. It takes a CurlRequest ID and a Model ID, executes the curl request to get data, and then uses that data to make an AI request to Ollama.

### Request Flow

1. **Retrieve CurlRequest**: Gets the stored curl request by ID
2. **Execute Curl**: Runs the curl request to get response data
3. **Retrieve Model**: Gets the DeepseekModel configuration by ID
4. **Build AI Request**: Creates an Ollama request with:
   - Assistant message: Content from curl response
   - User message: "Please check the current rate from the json.Is it greater than 125 ? Return Json Response"
5. **Make AI Call**: Sends request to Ollama API
6. **Return Response**: Returns the Ollama response

### Request Body

```json
{
  "curl_request_id": 1,
  "model_id": 1
}
```

### Example Request

```shell
curl -X POST ${base_url}/api/ai/make-request \
  -H "Content-Type: application/json" \
  -d '{
    "curl_request_id": 1,
    "model_id": 1
  }'
```

### Example Response

```json
{
  "model": "deepseek-r1:1.5b",
  "created_at": "2025-07-20T16:35:12Z",
  "message": {
    "role": "assistant",
    "content": "{\"rate_greater_than_125\": false, \"current_rate\": 120.8, \"analysis\": \"The current rate is 120.8, which is less than 125.\"}"
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 1234567,
  "load_duration": 123456,
  "prompt_eval_count": 10,
  "prompt_eval_duration": 123456,
  "eval_count": 50,
  "eval_duration": 987654
}
```

### Prerequisites

Before using this endpoint, ensure you have:

1. **CurlRequest Record**: A stored curl request with the specified `curl_request_id`
2. **DeepseekModel Record**: A stored model configuration with the specified `model_id`
3. **Ollama Service**: Ollama running and accessible at the model's `base_url`

### Example Setup

#### 1. Create a CurlRequest

```shell
curl -X POST ${base_url}/api/curl \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.exchangerate-api.com/v4/latest/USD",
    "method": "GET",
    "headers": {
      "Accept": "application/json"
    }
  }'
```

#### 2. Create a DeepseekModel

```shell
curl -X POST ${base_url}/api/deepseek-model \
  -H "Content-Type: application/json" \
  -d '{
    "name": "deepseek-r1:1.5b",
    "model": "deepseek-r1:1.5b",
    "base_url": "http://localhost:11434",
    "modified_at": "2025-07-20T07:44:29.156565162Z",
    "size": 1117322768
  }'
```

#### 3. Make AI Request

```shell
curl -X POST ${base_url}/api/ai/make-request \
  -H "Content-Type: application/json" \
  -d '{
    "curl_request_id": 1,
    "model_id": 1
  }'
```

### Error Responses

#### 400 Bad Request - Invalid Payload

```json
{
  "error": "INVALID_BODY",
  "message": "Invalid Input",
  "details": "Field validation for 'curl_request_id' failed on the 'required' tag"
}
```

#### 404 Not Found - CurlRequest Not Found

```json
{
  "error": "RECORD_NOT_FOUND",
  "message": "The requested record was not found",
  "details": "curl request with id 999 not found"
}
```

#### 404 Not Found - Model Not Found

```json
{
  "error": "RECORD_NOT_FOUND",
  "message": "The requested record was not found",
  "details": "deepseek model with id 999 not found"
}
```

#### 500 Internal Server Error - Ollama Service Error

```json
{
  "error": "EXTERNAL_SERVICE_ERROR",
  "message": "External service error",
  "details": "failed to make HTTP request to Ollama: connection refused"
}
```

### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `curl_request_id` | uint | Yes | ID of the stored curl request to execute |
| `model_id` | uint | Yes | ID of the DeepseekModel to use for AI processing |

### Response Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `model` | string | Name of the model used for the request |
| `created_at` | string | Timestamp when the response was created |
| `message` | object | The AI-generated response message |
| `message.role` | string | Role of the message (assistant) |
| `message.content` | string | The actual AI response content |
| `done_reason` | string | Reason why the response is complete |
| `done` | boolean | Whether the response is complete |
| `total_duration` | int64 | Total time taken for the request |
| `load_duration` | int64 | Time taken to load the model |
| `prompt_eval_count` | int | Number of prompt evaluations |
| `prompt_eval_duration` | int64 | Time taken for prompt evaluation |
| `eval_count` | int | Number of evaluations |
| `eval_duration` | int64 | Time taken for evaluations |

### Use Cases

1. **Currency Rate Analysis**: Execute a curl request to get currency rates, then use AI to analyze if the rate is above a threshold
2. **Data Processing**: Fetch data via curl and use AI to process or analyze the response
3. **API Integration**: Combine external API calls with AI processing for intelligent data handling

### Notes

- The endpoint automatically handles JSON parsing of curl response headers
- The AI request includes a specific prompt about checking rates against 125
- The endpoint is designed to work with Ollama's chat API format
- All timestamps are in ISO 8601 format
- The service includes comprehensive error handling for various failure scenarios 