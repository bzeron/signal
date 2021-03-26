package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// Error
type Error struct {
	s os.Signal
}

// Error impl error interface
func (e Error) Error() string {
	return fmt.Sprintf("signal: %s", e.s)
}

// Group
type Group struct {
	ech    chan error
	sch    chan os.Signal
	sig    []os.Signal
	cancel context.CancelFunc
}

// WithContext returns a new Group
func WithContext(ctx context.Context, sig ...os.Signal) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ech:    make(chan error),
		sch:    make(chan os.Signal),
		sig:    sig,
		cancel: cancel,
	}, ctx
}

// Wait for signal or goroutine error
func (g *Group) Wait() error {
	signal.Notify(g.sch, g.sig...)
	defer signal.Stop(g.sch)
	defer g.cancel()
	defer close(g.ech)
	select {
	case err := <-g.ech:
		return err
	case sig := <-g.sch:
		return &Error{s: sig}
	}
}

// Go calls the given function in a new goroutine.
func (g *Group) Go(fn func() error) {
	if fn != nil {
		go func() {
			if err := fn(); err != nil {
				g.ech <- err
			}
		}()
	}
}
