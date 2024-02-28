package syncerr

import (
	"errors"
	"fmt"
	"sync"
)

// Group is a concurrent error collector. Unlike errgroup.Group
// it is context-independent and returns all errors. It's zero value is valid.
type Group struct {
	sync.Mutex
	wg   sync.WaitGroup
	once sync.Once
	errs chan []error
}

// Go runs fn in a separate goroutine with panic recovery and collection of the return error.
func (e *Group) Go(fn func() error) {
	e.once.Do(func() {
		e.errs = make(chan []error, 1)
		e.errs <- nil // when we extract the errors with Join, nils are ignored
	})
	e.wg.Add(1)
	go func() {
		var err error
		defer e.wg.Done()
		defer e.put(&err)
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Group: panic recovered: %v", r)
			}
		}()
		err = fn()
	}()
}

func (e *Group) put(errp *error) {
	if *errp != nil {
		e.errs <- append(<-e.errs, *errp)
	}
}

// Wait blocks as long as the underlying wait group blocks.
func (e *Group) Wait() {
	e.wg.Wait()
}

// Done wraps the underlying wait group so that it can be used in select statements.
func (e *Group) Done() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(c)
	}()
	return c
}

// Error returns the result of calling errors.Join with all errors returned by the functions supplied to Group.Go calls.
// It should be called after Wait() or Done() to ensure all errors are returned.
// If a program doesn't care about all errors, errgroup.Group should be preferred.
// The returned error implements interface{ Unwrap() []error }.
func (e *Group) Error() error {
	// this case happens when Go was never called
	if len(e.errs) == 0 {
		return nil
	}
	return errors.Join(<-e.errs...)
}
