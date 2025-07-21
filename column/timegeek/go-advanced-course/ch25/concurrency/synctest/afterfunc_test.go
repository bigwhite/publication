package concurrency

import (
	"context"
	"testing"
	"testing/synctest"
)

func TestAfterFuncWithSyncTest(t *testing.T) {
	// The testing.T passed to synctest.Test's callback is special.
	// We use the outer 't' to call synctest.Test itself.
	synctest.Test(t, func(st *testing.T) { // st is the synctest-aware *testing.T
		ctx, cancel := context.WithCancel(context.Background())

		called := false
		context.AfterFunc(ctx, func() {
			called = true
		})

		// Assertion 1: AfterFunc should not have been called yet.
		synctest.Wait() // Wait for all goroutines in the bubble to settle.
		// The AfterFunc's goroutine is likely blocked waiting for ctx.Done().
		if called {
			st.Fatal("AfterFunc was called before context was canceled")
		}

		cancel() // Cancel the context, this should trigger AfterFunc.

		// Assertion 2: AfterFunc should now have been called.
		synctest.Wait() // Wait again for AfterFunc's goroutine to run and set 'called'.
		if !called {
			st.Fatal("AfterFunc was not called after context was canceled")
		}
		st.Log("Test with synctest completed successfully.")
	})
}
