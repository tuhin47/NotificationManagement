package utils

// Role constants for CURL operations
const (
	RoleCurlCreate = "curl_create"
	RoleCurlRead   = "curl_read"
	RoleCurlUpdate = "curl_update"
	RoleCurlDelete = "curl_delete"
)

// Role constants for LLM operations
const (
	RoleLLMCreate = "llm_create"
	RoleLLMRead   = "llm_read"
	RoleLLMUpdate = "llm_update"
	RoleLLMDelete = "llm_delete"
)

// Role constants for Reminder operations
const (
	RoleReminderCreate = "reminder_create"
	RoleReminderRead   = "reminder_read"
	RoleReminderUpdate = "reminder_update"
	RoleReminderDelete = "reminder_delete"
)

const (
	// Role constants for Deepseek Model operations
	RoleAICreate = "ai_model_create"
	RoleAIRead   = "ai_model_read"
	RoleAIUpdate = "ai_model_update"
	RoleAIDelete = "ai_model_delete"
	// Role constants for AI operations
	RoleMakeRequest = "ai_make_request"
)
