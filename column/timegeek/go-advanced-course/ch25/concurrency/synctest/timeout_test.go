package concurrency

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestWithTimeoutWithSyncTest(t *testing.T) {
	synctest.Test(t, func(st *testing.T) {
		// Create a context.Context which is canceled after a timeout.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Wait just less than the timeout.
		time.Sleep(timeout - time.Nanosecond)
		synctest.Wait()
		st.Logf("before timeout: ctx.Err() = %v\n", ctx.Err())

		// Wait the rest of the way until the timeout.
		time.Sleep(time.Nanosecond)
		synctest.Wait()
		st.Logf("after timeout:  ctx.Err() = %v\n", ctx.Err())

	})
}
