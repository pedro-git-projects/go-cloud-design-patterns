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

// DebounceLast works like DebounceFirst but
// uses a time.Ticker to determine whether
// enough time has passed since the function
// was last called, calling circuit when it has.
func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	threshold := time.Now()
	ticker := &time.Ticker{}
	result := *new(string)
	err := *new(error)
	once := sync.Once{}
	m := sync.Mutex{}

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()
		threshold = time.Now().Add(d)
		once.Do(func() {
			ticker = time.NewTicker(time.Millisecond * 100)
			go func() {
				defer func() {
					m.Lock()
					ticker.Stop()
					once = sync.Once{}
					m.Unlock()
				}()
				for {
					select {
					case <-ticker.C:
						m.Lock()
						if time.Now().After(threshold) {
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})
		return result, err
	}
}
