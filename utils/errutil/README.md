# Error Handling System

This package provides a centralized error handling system for the NotificationManagement application.

## Overview

The error handling system consists of:

1. **ErrorResponse Struct**: Standardized error response structure
2. **ErrorCode Struct**: Predefined error codes with messages and status codes
3. **Helper Functions**: Utility functions for common error handling patterns

## ErrorResponse Struct

```go
type ErrorResponse struct {
    Message     string    `json:"message"`     // Human-readable error message
    Error       string    `json:"error"`       // Actual error message for debugging
    StatusCode  int       `json:"status_code"` // HTTP status code
    Timestamp   time.Time `json:"timestamp"`   // Human-readable timestamp
    ErrorCode   string    `json:"error_code"`  // Internal error code for tracking
}
```

## Usage Examples

### 1. Basic Error Handling

```go
import "your-project/utils/errutil"

// In your controller
func (c *Controller) HandleRequest(w http.ResponseWriter, r *http.Request) {
    if err := someOperation(); err != nil {
        errResp := errutil.CreateErrorResponse(errutil.ErrInvalidInput, err)
        errutil.WriteErrorResponse(w, errResp)
        return
    }
}
```

### 2. Using Helper Functions

```go
// Handle database errors
if err := db.Query(); err != nil {
    errResp := errutil.HandleDatabaseError(err)
    errutil.WriteErrorResponse(w, errResp)
    return
}

// Handle validation errors
if err := validateInput(data); err != nil {
    errResp := errutil.HandleValidationError(err)
    errutil.WriteErrorResponse(w, errResp)
    return
}
```

### 3. Custom Error Messages

```go
errResp := errutil.CreateErrorResponseWithMessage(
    errutil.ErrNotificationFailed,
    err,
    "Failed to send notification to user@example.com",
)
errutil.WriteErrorResponse(w, errResp)
```

### 4. Validation in Services

```go
func (s *Service) ProcessUser(userData map[string]interface{}) error {
    // Validate required fields
    if err := errutil.ValidateRequiredFields(userData); err != nil {
        return err
    }

    // Validate email format
    if email, ok := userData["email"].(string); ok {
        if err := errutil.ValidateEmailFormat(email); err != nil {
            return err
        }
    }

    return nil
}
```

## Predefined Error Codes

### Database Errors
- `ErrDatabaseConnection` (500)
- `ErrRecordNotFound` (404)
- `ErrDatabaseQuery` (500)

### Validation Errors
- `ErrInvalidInput` (400)
- `ErrInvalidEmail` (400)
- `ErrInvalidPassword` (400)

### Authentication Errors
- `ErrInvalidCredentials` (401)
- `ErrUnauthorized` (401)
- `ErrTokenExpired` (401)
- `ErrInvalidToken` (401)

### Business Logic Errors
- `ErrUserAlreadyExists` (409)
- `ErrUserNotFound` (404)
- `ErrEmailAlreadyInUse` (409)

### Server Errors
- `ErrInternalServer` (500)
- `ErrServiceUnavailable` (503)

### Notification Errors
- `ErrNotificationFailed` (500)
- `ErrNotificationNotFound` (404)

### File/Upload Errors
- `ErrFileUploadFailed` (500)
- `ErrInvalidFileType` (400)
- `ErrFileTooLarge` (400)

### Rate Limiting
- `ErrRateLimitExceeded` (429)

### External Service Errors
- `ErrExternalServiceError` (502)

## Adding New Error Codes

To add a new error code, add it to the predefined error codes section in `err.go`:

```go
var (
    // Your new error category
    ErrYourNewError = ErrorCode{
        Code:    "YOUR_ERROR_CODE",
        Message: "Human-readable error message",
        Status:  400, // Appropriate HTTP status code
    }
)
```

## Best Practices

1. **Use Helper Functions**: Use the provided helper functions for common error types
2. **Consistent Error Codes**: Always use predefined error codes for consistency
3. **Meaningful Messages**: Provide clear, human-readable error messages
4. **Proper Status Codes**: Use appropriate HTTP status codes
5. **Log Errors**: Log actual errors for debugging while showing user-friendly messages
6. **Validation**: Use validation functions in services, not controllers

## Error Response Format

The system returns JSON responses in this format:

```json
{
    "message": "Human-readable error message",
    "error": "Actual error message for debugging",
    "status_code": 400,
    "timestamp": "2024-01-15T10:30:45Z",
    "error_code": "INVALID_INPUT"
}
```

## Middleware Integration

You can use the error handling middleware to catch panics:

```go
router.Use(errutil.ErrorHandlingMiddleware)
```

This ensures that any unhandled panics are converted to proper error responses. 