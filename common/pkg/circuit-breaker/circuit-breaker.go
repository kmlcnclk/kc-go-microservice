package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func DefaultCircuitBreaker() *gobreaker.CircuitBreaker {
	breakerSettings := gobreaker.Settings{
		Name:        "http-client",
		MaxRequests: 3,                // Number of requests allowed in half-open state
		Interval:    5 * time.Second,  // Time window for counting failures
		Timeout:     10 * time.Second, // Time to wait before switching from open to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			zap.L().Info("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()))
			// You could add logging here
		},
	}

	return gobreaker.NewCircuitBreaker(breakerSettings)
}
