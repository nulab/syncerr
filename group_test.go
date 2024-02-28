package syncerr

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	t.Run("instantiate", func(t *testing.T) {
		s := new(Group)
		assert.Nil(t, s.errs)
		assert.Nil(t, s.Error())
	})

	t.Run("one goroutine with error", func(t *testing.T) {
		s := new(Group)
		s.Go(func() error {
			return errors.New("problem")
		})
		s.Wait()

		require.NotNil(t, s.errs)
		assert.Len(t, s.errs, 1)
		assert.NotNil(t, s.Error())
	})

	t.Run("three goroutines with two errors", func(t *testing.T) {
		s := new(Group)
		s.Go(func() error {
			<-time.NewTimer(100 * time.Millisecond).C
			return errors.New("error 1")
		})
		s.Go(func() error {
			return errors.New("error 2")
		})
		s.Go(func() error {
			<-time.NewTimer(100 * time.Millisecond).C
			return nil
		})

		select {
		case <-s.Done():
			// no-op
		case <-time.NewTimer(500 * time.Millisecond).C:
			assert.FailNow(t, "timeout")
		}
		assert.Len(t, s.errs, 1) // only one item queued

		err := s.Error()
		assert.NotNil(t, err)

		wrapper, ok := err.(interface{ Unwrap() []error })
		require.True(t, ok, "err does not implement `interface{ Unwrap() []error }`")

		errs := wrapper.Unwrap()
		// in this test the order of the errors is predictable
		assert.Equal(t, "error 2", errs[0].Error())
		assert.Equal(t, "error 1", errs[1].Error())
	})

	t.Run("all errors", func(t *testing.T) {
		s := new(Group)
		s.Go(func() error {
			return errors.New("foo")
		})
		s.Go(func() error {
			return errors.New("bar")
		})
		s.Go(func() error {
			return errors.New("baz")
		})
		s.Wait()

		assert.ElementsMatch(t, []string{"foo", "bar", "baz"}, strings.Split(s.Error().Error(), "\n"))
	})
}
