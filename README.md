# syncerr

[![Go Report Card](https://goreportcard.com/badge/github.com/nulab/syncerr)](https://goreportcard.com/report/github.com/nulab/syncerr)
[![Go Reference](https://pkg.go.dev/badge/github.com/nulab/syncerr.svg)](https://pkg.go.dev/github.com/nulab/syncerr)

Alternative to [`errgroup.Group`](https://pkg.go.dev/golang.org/x/sync/errgroup) that is context-independent and returns all errors.

In some cases it may be desirable to run multiple tasks concurrently and inspect or log 
the error returned by each of them. This package provides an alternative to `errgroup.Group` that implements this 
alternate behavior.

## Dependencies

The main module has no dependencies beside the Go standard library.   

## Usage

The main type provided by this package is `Group`. Its API is similar to `errgroup.Group`:

```go
package main

import (
	"errors"
	
	"github.com/nulab/syncerr"
)

func main() {
	// You can initialize a new variable by just using `new`  
	g := new(syncerr.Group)
	
	// run tasks concurrently similarly to what you'd do with errgroup.Group
	g.Go(func() error {
		// do some task
		return errors.New("whoops!")
    })

	g.Go(func() error {
		// do some other task
		return errors.New("ouch!")
	})
	
	// wait for all tasks to complete
	g.Wait()
	
	// this error wraps all errors returned by the child tasks with errors.Join 
	err := g.Error()
	if err != nil {
		// handle error
    }
	// success!
}
```

## Status

This project is actively under development, but it is currently in version 0.
Please be aware that the public API and exported methods may undergo changes.

## Bug reporting

If you encounter a bug, please open a new issue and include the necessary steps to reproduce it. Thank you!

## Authors

* **[Gabriele V.](https://github.com/vibridi/)** - *Main contributor*

## License

This project is licensed under the MIT License. For detailed licensing information, refer to the [LICENSE](LICENSE) file included in the repository.
