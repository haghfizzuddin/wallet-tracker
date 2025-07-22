package errors

import (
	"errors"
	"fmt"
)

// Custom error types for better error handling
var (
	// API errors
	ErrAPIRateLimit    = errors.New("API rate limit exceeded")
	ErrAPIUnavailable  = errors.New("API service unavailable")
	ErrInvalidResponse = errors.New("invalid API response")
	
	// Wallet errors
	ErrInvalidWallet  = errors.New("invalid wallet address")
	ErrWalletNotFound = errors.New("wallet not found")
	
	// Database errors
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseWrite      = errors.New("database write failed")
	
	// Configuration errors
	ErrInvalidConfig = errors.New("invalid configuration")
	ErrMissingConfig = errors.New("missing required configuration")
)

// WrapError wraps an error with additional context
func WrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	return errors.Is(err, ErrAPIRateLimit) || 
		   errors.Is(err, ErrAPIUnavailable) ||
		   errors.Is(err, ErrDatabaseConnection)
}
