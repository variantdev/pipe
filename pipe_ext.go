package pipe

import "context"

func New(ctx context.Context) *State {
	st := NewState(nil, nil)
	st.ctx = ctx

	return st
}

type Result struct {
	State *State
	Error error
}

func Do(parent context.Context, s *State, pipe Pipe) error {
	if err := pipe(s); err != nil {
		return err
	}

	ch := make(chan error)
	go func() {
		ch <- s.Do(parent)
	}()

	return <-ch
}

func (s *State) Do(parent context.Context) error {
	var ctx context.Context
	var cancel context.CancelFunc
	if parent != nil {
		ctx, cancel = context.WithCancel(parent)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	ch := make(chan error)
	go func() {
		ch <- s.RunTasks()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
