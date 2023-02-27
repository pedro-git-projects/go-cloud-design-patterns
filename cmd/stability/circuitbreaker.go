// stability stores patterns that address one or more of the assumptions called
// out by the Fallacies of Distributed Computing.
// They’re generally intended to be applied by distributed applications
// to improve their own stability and the stability of
// the larger system they’re a part of.
package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Circuit Breaker automatically degrades service functions
// in response to a likely fault, preventing larger
// or cascading failures by eliminating recurring errors and providing
// reasonable error responses.

// Essentially, the Circuit Breaker is just a specialized Adapter pattern,
// with Breaker wrapping Circuit to add some additional error handling logic.

// The Breaker function accepts any function that conforms to the Circuit type definition,
// and an unsigned integer representing the number of consecutive failures
// allowed before the circuit automatically opens.
// In return it provides another function, which also conforms to the Circuit type definition
// The closure works by counting the number of consecutive errors returned by
// circuit. If that value meets the failure threshold, then it returns the error “service
// unreachable” without actually calling circuit. Any successful calls to circuit cause
// consecutiveFailures to reset to 0, and the cycle begins again.
func Breaker(circuit Circuit, failiureThreshhold uint) Circuit {
	consecutiveFailiures := *new(int)
	lastAttempt := time.Now()
	m := sync.RWMutex{}

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := consecutiveFailiures - int(failiureThreshhold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("unreachable service")
			}
		}

		m.RUnlock()
		response, err := circuit(ctx)
		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailiures++
			return response, err
		}

		consecutiveFailiures = 0
		return response, nil
	}
}
