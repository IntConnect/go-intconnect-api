package exception

// HTTP Error Messages
const (
	ErrInvalidRequestBody  = "Invalid request body"
	ErrUnauthorized        = "Unauthorized access"
	ErrForbidden           = "Forbidden"
	ErrMethodNotAllowed    = "Method not allowed"
	ErrNotAcceptable       = "Not Acceptable"
	ErrProxyAuthRequired   = "Proxy Authentication Required"
	ErrRequestTimeout      = "Request Timeout"
	ErrTooManyRequests     = "Too Many Requests"
	ErrInternalServerError = "Internal Server Error"
	ErrBadRequest          = "Bad Request"
	ErrNotFound            = "Not Found"
)

// Application Error Messages
const (
	// Resource Operations
	ErrSavingResources      = "Error saving resources"
	ErrSomeResourceNotFound = "Some resources not found"
	ErrCreatingResource     = "Error creating resource"
	ErrUpdatingResource     = "Error updating resource"
	ErrDeletingResource     = "Error deleting resource"
	ErrResourceAlreadyExist = "Resource already exists"
	ErrResourceNotFound     = "Resource not found"

	// Database Operations
	ErrDatabaseConnection  = "Database connection error"
	ErrDatabaseQuery       = "Database query error"
	ErrDatabaseTransaction = "Database transaction error"
	ErrDuplicateEntry      = "Duplicate entry"

	// Authentication & Authorization
	ErrInvalidCredentials     = "Invalid credentials"
	ErrTokenExpired           = "Token has expired"
	ErrTokenInvalid           = "Invalid token"
	ErrInsufficientPermission = "Insufficient permission"
	ErrAccountLocked          = "Account is locked"
	ErrAccountNotActivated    = "Account is not activated"

	// Business Logic Errors
	ErrInsufficientBalance = "Insufficient balance"
	ErrInvalidOperation    = "Invalid operation"
	ErrOperationNotAllowed = "Operation not allowed"
	ErrQuotaExceeded       = "Quota exceeded"
	ErrServiceUnavailable  = "Service temporarily unavailable"

	// File Operations
	ErrFileUpload       = "Error uploading file"
	ErrFileNotFound     = "File not found"
	ErrInvalidFileType  = "Invalid file type"
	ErrFileSizeTooLarge = "File size too large"

	// External Service Errors
	ErrExternalServiceError = "External service error"
	ErrThirdPartyAPIError   = "Third party API error"
	ErrPaymentFailed        = "Payment processing failed"

	ErrPayloadInvalid = "Payload invalid or deformed"
)

// Convention Status Code
// Custom API status codes for standardized error handling
const (
	StatusSuccess            = "0"  // Success operation
	StatusValidationError    = "1"  // Validation error (input validation failed)
	StatusAuthError          = "2"  // Authentication/Authorization error
	StatusNotFoundError      = "3"  // Resource not found
	StatusDuplicateError     = "4"  // Duplicate resource/entry
	StatusDatabaseError      = "5"  // Database operation error
	StatusBusinessLogicError = "6"  // Business logic error
	StatusExternalError      = "7"  // External service error
	StatusInternalError      = "8"  // Internal server error
	StatusRateLimitError     = "9"  // Rate limit exceeded
	StatusTimeoutError       = "10" // Request timeout
	StatusMaintenanceMode    = "11" // Service in maintenance mode
)

// Convention Status Code Messages
const (
	MsgSuccess            = "Operation completed successfully"
	MsgValidationError    = "Validation error occurred"
	MsgAuthError          = "Authentication or authorization error"
	MsgNotFoundError      = "Resource not found"
	MsgDuplicateError     = "Duplicate resource detected"
	MsgDatabaseError      = "Database operation failed"
	MsgBusinessLogicError = "Business logic error"
	MsgExternalError      = "External service error"
	MsgInternalError      = "Internal server error"
	MsgRateLimitError     = "Rate limit exceeded"
	MsgTimeoutError       = "Request timeout"
	MsgMaintenanceMode    = "Service under maintenance"
)
