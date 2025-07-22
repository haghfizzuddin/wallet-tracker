package retry

import (
	"context"
	"math"
	"time"
	
	"github.com/aydinnyunus/wallet-tracker/pkg/errors"
)

// Config holds retry configuration
type Config struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	RandomizeFactor float64
}

// DefaultConfig returns default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     3,
		InitialDelay:    1 * time.Second,
		MaxDelay:        30 * time.Second,
		Multiplier:      2.0,
		RandomizeFactor: 0.1,
	}
}

// Func is a function that can be retried
type Func func() error

// Do executes the function with retry logic
func Do(ctx context.Context, fn Func, config Config) error {
	var lastErr error
	
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}
		
		// Check if error is retryable
		if !errors.IsRetryableError(err) {
			return err
		}
		
		lastErr = err
		
		// Check if we've exhausted retries
		if attempt >= config.MaxAttempts-1 {
			break
		}
		
		// Calculate delay with exponential backoff
		delay := calculateDelay(attempt, config)
		
		// Wait or return if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	return lastErr
}

// calculateDelay calculates the delay for the next retry attempt
func calculateDelay(attempt int, config Config) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))
	
	// Cap at max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}
	
	// Add some randomization to prevent thundering herd
	if config.RandomizeFactor > 0 {
		delta := delay * config.RandomizeFactor
		minDelay := delay - delta
		maxDelay := delay + delta
		
		// Simple randomization
		delay = minDelay + (maxDelay-minDelay)*0.5
	}
	
	return time.Duration(delay)
}
