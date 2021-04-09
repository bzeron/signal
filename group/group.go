package group

import (
	"context"
	"os"
)

// Group
type Group struct {
	ech    chan error
	sch    chan os.Signal
	ctx    context.Context
	cancel context.CancelFunc
}

// WithContext returns a new Group
func WithContext(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ech:    make(chan error),
		sch:    make(chan os.Signal),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Wait for signal or goroutine error
func (g *Group) Wait() error {
	defer g.cancel()
	defer close(g.ech)
	select {
	case err := <-g.ech:
		return err
	case <-g.ctx.Done():
		return g.ctx.Err()
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
