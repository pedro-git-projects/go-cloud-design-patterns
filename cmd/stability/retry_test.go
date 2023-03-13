package stability_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pedro-git-projects/go-cloud-design-patterns/cmd/stability"
)

func TestRetry(t *testing.T) {
	var count int
	emulateTransientErr := func(ctx context.Context) (string, error) {
		count++
		if count <= 3 {
			return "Intentional failiure", errors.New("error")
		} else {
			return "success", nil
		}
	}

	r := stability.Retry(emulateTransientErr, 5, 2*time.Second)

	_, err := r(context.Background())
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}

}
