package stability

import (
	"context"
	"sync"
	"time"
)

// In world of services, we sometimes find ourselves performing
//a cluster of potentially slow or costly operations where
// only one would do. Using the Debounce pattern,
// a series of similar calls that are tightly clustered
// in time are restricted to only one call,
// typically the first or last in a batch.

//	DebounceFirst works by making so that
//
// on each call of the outer function
// —regardless of its outcome— a time interval is set.
// Any subsequent call made before that time interval expires is ignored.
// Any call made afterwards is passed along to the inner function.
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	threshold := time.Time{}
	result := *new(string)
	err := *new(error)
	m := sync.Mutex{}

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)
		return result, err
	}
}
